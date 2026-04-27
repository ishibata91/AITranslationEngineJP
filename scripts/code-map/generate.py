from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from dataclasses import dataclass
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

CODE_EXTENSIONS = {".go", ".svelte", ".ts"}
ROOTS = ["frontend/src", "internal"]
SCHEMA_VERSION = 1

IMPORT_SPEC_PATTERN = re.compile(
    r"""
    (?:
        \bimport\s+(?:type\s+)?(?:[^'"]+?\s+from\s+)?
      | \bexport\s+(?:type\s+)?(?:[^'"]+?\s+from\s+)
    )
    ['"](?P<spec>[^'"]+)['"]
    """,
    re.VERBOSE | re.MULTILINE,
)

FRONTEND_ALIAS_PREFIXES = {
    "@application/": Path("frontend/src/application"),
    "@controller/": Path("frontend/src/controller"),
    "@ui/": Path("frontend/src/ui"),
}


@dataclass(frozen=True)
class LayerDefinition:
    id: str
    name: str
    root: str
    paths: tuple[Path, ...]
    default_next: tuple[str, ...]


LAYER_DEFINITIONS = (
    LayerDefinition(
        id="frontend-bootstrap",
        name="Frontend Bootstrap",
        root="frontend/src",
        paths=(Path("frontend/src/main.ts"),),
        default_next=("frontend-view", "frontend-wails-adapter"),
    ),
    LayerDefinition(
        id="frontend-view",
        name="View",
        root="frontend/src",
        paths=(Path("frontend/src/ui"), Path("frontend/src/test")),
        default_next=("frontend-controller", "frontend-presenter-store"),
    ),
    LayerDefinition(
        id="frontend-controller",
        name="ScreenController",
        root="frontend/src",
        paths=(
            Path("frontend/src/controller/master-dictionary"),
            Path("frontend/src/controller/master-persona"),
            Path("frontend/src/controller/translation-input"),
            Path("frontend/src/controller/translation-job-setup"),
        ),
        default_next=("frontend-usecase", "frontend-presenter-store", "frontend-contract", "frontend-wails-adapter"),
    ),
    LayerDefinition(
        id="frontend-usecase",
        name="Frontend UseCase",
        root="frontend/src",
        paths=(Path("frontend/src/application/usecase"),),
        default_next=("frontend-contract", "frontend-presenter-store", "frontend-wails-adapter"),
    ),
    LayerDefinition(
        id="frontend-presenter-store",
        name="Presenter / Store",
        root="frontend/src",
        paths=(Path("frontend/src/application/presenter"), Path("frontend/src/application/store")),
        default_next=("frontend-contract",),
    ),
    LayerDefinition(
        id="frontend-contract",
        name="Contract",
        root="frontend/src",
        paths=(Path("frontend/src/application/contract"), Path("frontend/src/application/gateway-contract")),
        default_next=("frontend-wails-adapter",),
    ),
    LayerDefinition(
        id="frontend-wails-adapter",
        name="Wails Adapter",
        root="frontend/src",
        paths=(Path("frontend/src/controller/wails"), Path("frontend/src/controller/runtime")),
        default_next=("backend-controller",),
    ),
    LayerDefinition(
        id="backend-bootstrap",
        name="Backend Bootstrap",
        root="internal",
        paths=(Path("internal/bootstrap"),),
        default_next=("backend-controller", "backend-usecase", "backend-service", "backend-repository", "backend-infra-provider"),
    ),
    LayerDefinition(
        id="backend-controller",
        name="Controller",
        root="internal",
        paths=(Path("internal/controller"),),
        default_next=("backend-usecase",),
    ),
    LayerDefinition(
        id="backend-usecase",
        name="UseCase",
        root="internal",
        paths=(Path("internal/usecase"),),
        default_next=("backend-service", "backend-state-jobio"),
    ),
    LayerDefinition(
        id="backend-service",
        name="Service",
        root="internal",
        paths=(Path("internal/service"),),
        default_next=("backend-repository", "backend-state-jobio", "backend-infra-provider"),
    ),
    LayerDefinition(
        id="backend-state-jobio",
        name="State / JobIO",
        root="internal",
        paths=(Path("internal/statemachine"), Path("internal/jobio")),
        default_next=(),
    ),
    LayerDefinition(
        id="backend-repository",
        name="Repository",
        root="internal",
        paths=(Path("internal/repository"),),
        default_next=("backend-infra-provider",),
    ),
    LayerDefinition(
        id="backend-infra-provider",
        name="Infra / Provider",
        root="internal",
        paths=(Path("internal/infra"),),
        default_next=(),
    ),
    # integration test 専用ディレクトリ。
    # 複数コンポーネントをまたぐ SQLite integration test だけを置く場所。
    LayerDefinition(
        id="backend-integration-test",
        name="Integration Test",
        root="internal",
        paths=(Path("internal/integrationtest"),),
        default_next=(),
    ),
)

LAYER_BY_ID = {layer.id: layer for layer in LAYER_DEFINITIONS}


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Generate an entry-level code map for AI navigation.")
    parser.add_argument("--repo-root", default=".", help="Repository root. Defaults to the current directory.")
    parser.add_argument(
        "--output",
        default="tmp/code-map/index.json",
        help="Output JSON path. Relative paths are resolved from repo root.",
    )
    return parser


def relative_path(path: Path, repo_root: Path) -> str:
    return path.resolve().relative_to(repo_root).as_posix()


def language_for(path: Path) -> str:
    if path.suffix == ".go":
        return "go"
    if path.suffix == ".svelte":
        return "svelte"
    if path.suffix == ".ts":
        return "typescript"
    return "unknown"


def is_test_file(path: Path) -> bool:
    return path.name.endswith("_test.go") or ".test." in path.name


def path_matches_layer(path: Path, layer_path: Path) -> bool:
    if layer_path.suffix:
        return path == layer_path
    return path == layer_path or layer_path in path.parents


def classify_layer(relative: Path) -> str | None:
    for layer in sorted(LAYER_DEFINITIONS, key=lambda definition: max(len(path.parts) for path in definition.paths), reverse=True):
        if any(path_matches_layer(relative, layer_path) for layer_path in layer.paths):
            return layer.id
    return None


def collect_code_files(repo_root: Path) -> list[Path]:
    files: list[Path] = []
    for root_name in ROOTS:
        root = repo_root / root_name
        if not root.exists():
            continue
        files.extend(path for path in root.rglob("*") if path.is_file() and path.suffix in CODE_EXTENSIONS)
    return sorted(files)


def resolve_frontend_import(source_file: Path, spec: str, repo_root: Path) -> str | None:
    candidate: Path | None = None
    if spec.startswith("."):
        candidate = source_file.parent / spec
    else:
        for alias_prefix, alias_target in FRONTEND_ALIAS_PREFIXES.items():
            if spec.startswith(alias_prefix):
                candidate = repo_root / alias_target / spec.removeprefix(alias_prefix)
                break

    if candidate is None:
        return None

    resolved = resolve_frontend_candidate(candidate)
    if resolved is None:
        return None

    try:
        resolved.relative_to(repo_root / "frontend/src")
    except ValueError:
        return None

    return relative_path(resolved, repo_root)


def resolve_frontend_candidate(candidate: Path) -> Path | None:
    if candidate.is_file():
        return candidate
    if candidate.is_dir():
        return candidate

    suffixes = ("", ".ts", ".svelte")
    for suffix in suffixes:
        path = Path(f"{candidate}{suffix}")
        if path.is_file():
            return path

    for index_name in ("index.ts", "index.svelte"):
        path = candidate / index_name
        if path.is_file():
            return candidate

    return None


def collect_frontend_dependencies(repo_root: Path, files_by_path: dict[str, dict[str, Any]]) -> list[dict[str, Any]]:
    dependencies: dict[tuple[str, str, str], dict[str, Any]] = {}
    for path_string, file_entry in files_by_path.items():
        path = repo_root / path_string
        if not path_string.startswith("frontend/src/") or path.suffix not in {".ts", ".svelte"}:
            continue
        if file_entry["is_test"]:
            continue

        content = path.read_text(encoding="utf-8")
        for match in IMPORT_SPEC_PATTERN.finditer(content):
            spec = match.group("spec")
            target = resolve_frontend_import(path, spec, repo_root)
            if target is None:
                continue

            target_layer = layer_for_dependency_target(Path(target), files_by_path)
            if target_layer is None or target_layer == file_entry["layer"]:
                continue

            key = (path_string, target, spec)
            dependencies[key] = {
                "from": path_string,
                "to": target,
                "kind": "frontend-import",
                "spec": spec,
                "from_layer": file_entry["layer"],
                "to_layer": target_layer,
            }

    return sorted(dependencies.values(), key=lambda entry: (entry["from"], entry["to"], entry["spec"]))


def layer_for_dependency_target(target: Path, files_by_path: dict[str, dict[str, Any]]) -> str | None:
    target_string = target.as_posix()
    if target_string in files_by_path:
        return str(files_by_path[target_string]["layer"])

    layer = classify_layer(target)
    if layer is not None:
        return layer

    for file_path, entry in files_by_path.items():
        if file_path.startswith(f"{target_string}/"):
            return str(entry["layer"])

    return None


def parse_go_list_objects(output: str) -> list[dict[str, Any]]:
    decoder = json.JSONDecoder()
    index = 0
    packages: list[dict[str, Any]] = []
    while index < len(output):
        while index < len(output) and output[index].isspace():
            index += 1
        if index >= len(output):
            break
        package, next_index = decoder.raw_decode(output, index)
        packages.append(package)
        index = next_index
    return packages


def collect_go_dependencies(repo_root: Path, files_by_path: dict[str, dict[str, Any]]) -> list[dict[str, Any]]:
    completed = subprocess.run(
        ["go", "list", "-json", "./internal/..."],
        cwd=repo_root,
        check=False,
        text=True,
        capture_output=True,
    )
    if completed.returncode != 0:
        raise RuntimeError(completed.stderr.strip() or "go list failed")

    packages = parse_go_list_objects(completed.stdout)
    module_path = read_module_path(repo_root)
    dependencies: dict[tuple[str, str], dict[str, Any]] = {}
    for package in packages:
        package_path = str(package["ImportPath"]).removeprefix(f"{module_path}/")
        if not package_path.startswith("internal/"):
            continue

        from_layer = layer_for_dependency_target(Path(package_path), files_by_path)
        for imported in package.get("Imports", []):
            imported_string = str(imported)
            if not imported_string.startswith(f"{module_path}/internal/"):
                continue

            imported_path = imported_string.removeprefix(f"{module_path}/")
            to_layer = layer_for_dependency_target(Path(imported_path), files_by_path)
            if from_layer is None or to_layer is None or from_layer == to_layer:
                continue

            key = (package_path, imported_path)
            dependencies[key] = {
                "from": package_path,
                "to": imported_path,
                "kind": "go-package",
                "from_layer": from_layer,
                "to_layer": to_layer,
            }

    return sorted(dependencies.values(), key=lambda entry: (entry["from"], entry["to"]))


def read_module_path(repo_root: Path) -> str:
    go_mod = repo_root / "go.mod"
    for line in go_mod.read_text(encoding="utf-8").splitlines():
        if line.startswith("module "):
            return line.split(maxsplit=1)[1].strip()
    raise RuntimeError("go.mod module path is missing")


def collect_tests(repo_root: Path, files_by_path: dict[str, dict[str, Any]]) -> list[dict[str, str]]:
    source_paths = {Path(path) for path, entry in files_by_path.items() if not entry["is_test"]}
    tests: list[dict[str, str]] = []
    for path_string, entry in files_by_path.items():
        if not entry["is_test"]:
            continue

        path = Path(path_string)
        target = infer_test_target(path, source_paths)
        tests.append(
            {
                "target": target.as_posix() if target is not None else "",
                "test": path_string,
            }
        )

    return sorted(tests, key=lambda entry: (entry["test"], entry["target"]))


def infer_test_target(test_path: Path, source_paths: set[Path]) -> Path | None:
    if test_path.name.endswith("_test.go"):
        stem = test_path.name.removesuffix("_test.go")
        candidates = [test_path.with_name(f"{stem}.go")]
        parts = stem.split("_")
        while len(parts) > 1:
            parts.pop()
            candidates.append(test_path.with_name(f"{'_'.join(parts)}.go"))
        return first_existing_candidate(candidates, source_paths)

    if ".test." in test_path.name:
        target_name = test_path.name.replace(".test.", ".", 1)
        candidates = [test_path.with_name(target_name)]
        if target_name.endswith(".ts"):
            candidates.append(test_path.with_name(target_name.removesuffix(".ts") + ".svelte"))
        return first_existing_candidate(candidates, source_paths)

    return None


def first_existing_candidate(candidates: list[Path], source_paths: set[Path]) -> Path | None:
    for candidate in candidates:
        if candidate in source_paths:
            return candidate
    return None


def build_next_entries(dependencies: list[dict[str, Any]]) -> tuple[list[dict[str, str]], dict[str, set[str]]]:
    next_entries: dict[tuple[str, str, str], dict[str, str]] = {}
    next_by_layer: dict[str, set[str]] = {layer.id: set(layer.default_next) for layer in LAYER_DEFINITIONS}

    for layer in LAYER_DEFINITIONS:
        for target in layer.default_next:
            next_entries[(layer.id, target, "default")] = {
                "from_layer": layer.id,
                "to_layer": target,
                "kind": "default",
            }

    for dependency in dependencies:
        from_layer = str(dependency["from_layer"])
        to_layer = str(dependency["to_layer"])
        next_by_layer.setdefault(from_layer, set()).add(to_layer)
        next_entries[(from_layer, to_layer, "observed")] = {
            "from_layer": from_layer,
            "to_layer": to_layer,
            "kind": "observed",
        }

    return sorted(next_entries.values(), key=lambda entry: (entry["from_layer"], entry["to_layer"], entry["kind"])), next_by_layer


def build_code_map(repo_root: Path) -> dict[str, Any]:
    files: list[dict[str, Any]] = []
    unmapped: list[str] = []
    for path in collect_code_files(repo_root):
        relative = Path(relative_path(path, repo_root))
        layer = classify_layer(relative)
        if layer is None:
            unmapped.append(relative.as_posix())
            layer = "unmapped"
        files.append(
            {
                "path": relative.as_posix(),
                "language": language_for(path),
                "layer": layer,
                "is_test": is_test_file(path),
            }
        )

    if unmapped:
        raise RuntimeError(f"unmapped code files: {', '.join(unmapped)}")

    files_by_path = {entry["path"]: entry for entry in files}
    dependencies = collect_frontend_dependencies(repo_root, files_by_path)
    dependencies.extend(collect_go_dependencies(repo_root, files_by_path))
    dependencies = sorted(dependencies, key=lambda entry: (entry["kind"], entry["from"], entry["to"]))
    tests = collect_tests(repo_root, files_by_path)
    next_entries, next_by_layer = build_next_entries(dependencies)

    layers = [
        {
            "id": layer.id,
            "name": layer.name,
            "root": layer.root,
            "paths": [path.as_posix() for path in layer.paths],
            "next": sorted(next_by_layer.get(layer.id, set())),
        }
        for layer in LAYER_DEFINITIONS
    ]

    return {
        "version": SCHEMA_VERSION,
        "generated_at": datetime.now(UTC).replace(microsecond=0).isoformat().replace("+00:00", "Z"),
        "roots": ROOTS,
        "layers": layers,
        "files": files,
        "dependencies": dependencies,
        "tests": tests,
        "next": next_entries,
    }


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()

    repo_root = Path(args.repo_root).resolve()
    output_path = Path(args.output)
    if not output_path.is_absolute():
        output_path = repo_root / output_path

    try:
        code_map = build_code_map(repo_root)
    except RuntimeError as error:
        print(f"code-map: {error}", file=sys.stderr)
        return 1

    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(
        json.dumps(code_map, ensure_ascii=False, indent=2, sort_keys=True) + "\n",
        encoding="utf-8",
    )
    print(f"Wrote {output_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
