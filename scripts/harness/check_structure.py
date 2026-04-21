from __future__ import annotations

import json
import re
import subprocess
import sys
import time
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    report_fail,
    report_pass,
    report_section,
    read_text,
    remove_code_blocks,
)

LINK_PATTERN = re.compile(r"!?\[[^\]]*\]\((?P<target>[^)]+)\)")
URI_SCHEME_PATTERN = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:")
WINDOWS_DRIVE_PATTERN = re.compile(r"^[a-zA-Z]:[\\\\/]")
STRUCTURE_INDEX_PATH = Path("docs/index.md")
DOCS_ROOT = Path("docs")
CODE_MAP_SCRIPT_PATH = Path("scripts/code-map/generate.py")
CODE_MAP_OUTPUT_PATH = Path("tmp/code-map/index.json")
CODE_MAP_TIME_BUDGET_SECONDS = 5.0
CODE_MAP_ROOTS = [Path("frontend/src"), Path("internal")]
CODE_MAP_EXTENSIONS = {".go", ".svelte", ".ts"}
CODE_MAP_REQUIRED_KEYS = {"version", "generated_at", "roots", "layers", "files", "dependencies", "tests"}
CODE_MAP_EXPECTED_LAYER_IDS = {
    "frontend-bootstrap",
    "frontend-view",
    "frontend-controller",
    "frontend-usecase",
    "frontend-presenter-store",
    "frontend-contract",
    "frontend-wails-adapter",
    "backend-bootstrap",
    "backend-controller",
    "backend-usecase",
    "backend-service",
    "backend-state-jobio",
    "backend-repository",
    "backend-infra-provider",
    # integration test 専用ディレクトリ。production code を置かない。
    "backend-integration-test",
}
ALWAYS_ALLOWED_DOC_PATHS = {
    Path("docs/index.md"),
}


def resolve_link_target(source_file: Path, target: str) -> Path | None:
    clean_target = target.split("#", 1)[0].strip()
    if not clean_target:
        return None
    if URI_SCHEME_PATTERN.match(clean_target) and not WINDOWS_DRIVE_PATTERN.match(clean_target):
        return None
    if clean_target.startswith("/"):
        return None
    if WINDOWS_DRIVE_PATTERN.match(clean_target):
        return Path(clean_target)
    return (source_file.parent / clean_target).resolve()


def check_index_map(repo_root: Path) -> int:
    failures = 0
    index_path = repo_root / STRUCTURE_INDEX_PATH
    report_section("Index map")

    if index_path.exists():
        report_pass(f"PASS index map exists: {index_path}")
    else:
        report_fail(f"FAIL missing index map: {index_path}")
        return 1

    content = remove_code_blocks(read_text(index_path))
    for match in LINK_PATTERN.finditer(content):
        target = match.group("target").strip()
        if target.startswith("#"):
            continue

        resolved = resolve_link_target(index_path, target)
        if resolved is None:
            continue

        if resolved.exists():
            report_pass(f"PASS index link: {index_path} -> {target}")
        else:
            report_fail(f"FAIL broken index link: {index_path} -> {target}")
            failures += 1

    return failures


def collect_index_targets(repo_root: Path) -> tuple[set[Path], set[Path]]:
    index_path = repo_root / STRUCTURE_INDEX_PATH
    content = remove_code_blocks(read_text(index_path))
    allowed_files: set[Path] = set()
    allowed_directories: set[Path] = set()

    for match in LINK_PATTERN.finditer(content):
        target = match.group("target").strip()
        if target.startswith("#"):
            continue

        resolved = resolve_link_target(index_path, target)
        if resolved is None or not resolved.exists():
            continue
        try:
            resolved.relative_to(repo_root / DOCS_ROOT)
        except ValueError:
            continue

        if resolved.is_dir():
            allowed_directories.add(resolved)
        else:
            allowed_files.add(resolved)
            if resolved.name.lower() == "readme.md":
                allowed_directories.add(resolved.parent)

    for relative_path in ALWAYS_ALLOWED_DOC_PATHS:
        allowed_files.add(repo_root / relative_path)

    return allowed_files, allowed_directories


def is_allowed_doc_file(file_path: Path, allowed_files: set[Path], allowed_directories: set[Path]) -> bool:
    if file_path in allowed_files:
        return True
    return any(directory in file_path.parents for directory in allowed_directories)


def check_docs_coverage(repo_root: Path) -> int:
    failures = 0
    docs_root = repo_root / DOCS_ROOT
    allowed_files, allowed_directories = collect_index_targets(repo_root)
    report_section("Docs coverage")

    for file_path in sorted(path for path in docs_root.rglob("*") if path.is_file()):
        if is_allowed_doc_file(file_path, allowed_files, allowed_directories):
            report_pass(f"PASS docs file is mapped by index: {file_path}")
        else:
            report_fail(f"FAIL docs file is not mapped by index: {file_path}")
            failures += 1

    return failures


def collect_code_map_target_files(repo_root: Path) -> list[Path]:
    files: list[Path] = []
    for root in CODE_MAP_ROOTS:
        absolute_root = repo_root / root
        if not absolute_root.exists():
            continue
        files.extend(path for path in absolute_root.rglob("*") if path.is_file() and path.suffix in CODE_MAP_EXTENSIONS)
    return sorted(files)


def is_test_code_file(path: Path) -> bool:
    return path.name.endswith("_test.go") or ".test." in path.name


def check_code_map(repo_root: Path) -> int:
    failures = 0
    script_path = repo_root / CODE_MAP_SCRIPT_PATH
    output_path = repo_root / CODE_MAP_OUTPUT_PATH
    report_section("Code map")

    if not script_path.exists():
        report_fail(f"FAIL missing code map generator: {script_path}")
        return 1

    started_at = time.perf_counter()
    completed = subprocess.run(
        [sys.executable, str(script_path), "--repo-root", str(repo_root), "--output", str(output_path)],
        cwd=repo_root,
        check=False,
        text=True,
        capture_output=True,
    )
    elapsed_seconds = time.perf_counter() - started_at

    if completed.returncode != 0:
        report_fail(f"FAIL code map generator exited with {completed.returncode}: {completed.stderr.strip()}")
        return 1

    if elapsed_seconds > CODE_MAP_TIME_BUDGET_SECONDS:
        report_fail(f"FAIL code map generator exceeded budget: {elapsed_seconds:.3f}s")
        failures += 1
    else:
        report_pass(f"PASS code map generator completed in {elapsed_seconds:.3f}s")

    if not output_path.exists():
        report_fail(f"FAIL missing code map output: {output_path}")
        return failures + 1

    try:
        code_map = json.loads(read_text(output_path))
    except json.JSONDecodeError as error:
        report_fail(f"FAIL code map output is invalid JSON: {error}")
        return failures + 1

    missing_keys = sorted(CODE_MAP_REQUIRED_KEYS - set(code_map))
    if missing_keys:
        report_fail(f"FAIL code map missing required keys: {', '.join(missing_keys)}")
        failures += 1
    else:
        report_pass("PASS code map has required schema keys")

    layer_ids = {str(layer.get("id")) for layer in code_map.get("layers", []) if isinstance(layer, dict)}
    if layer_ids == CODE_MAP_EXPECTED_LAYER_IDS:
        report_pass(f"PASS code map has expected layers: {len(layer_ids)}")
    else:
        missing_layers = sorted(CODE_MAP_EXPECTED_LAYER_IDS - layer_ids)
        extra_layers = sorted(layer_ids - CODE_MAP_EXPECTED_LAYER_IDS)
        report_fail(f"FAIL code map layer mismatch: missing={missing_layers}, extra={extra_layers}")
        failures += 1

    target_files = collect_code_map_target_files(repo_root)
    expected_paths = {path.relative_to(repo_root).as_posix() for path in target_files}
    indexed_paths = {str(entry.get("path")) for entry in code_map.get("files", []) if isinstance(entry, dict)}
    if indexed_paths == expected_paths:
        report_pass(f"PASS code map indexes target files: {len(indexed_paths)}")
    else:
        missing_files = sorted(expected_paths - indexed_paths)
        extra_files = sorted(indexed_paths - expected_paths)
        report_fail(f"FAIL code map file mismatch: missing={missing_files[:10]}, extra={extra_files[:10]}")
        failures += 1

    unmapped_files = sorted(
        str(entry.get("path"))
        for entry in code_map.get("files", [])
        if isinstance(entry, dict) and entry.get("layer") == "unmapped"
    )
    if unmapped_files:
        report_fail(f"FAIL code map has unmapped files: {unmapped_files[:10]}")
        failures += 1
    else:
        report_pass("PASS code map has no unmapped files")

    dependency_kinds = {str(entry.get("kind")) for entry in code_map.get("dependencies", []) if isinstance(entry, dict)}
    if {"frontend-import", "go-package"}.issubset(dependency_kinds):
        report_pass("PASS code map has frontend and backend dependencies")
    else:
        report_fail(f"FAIL code map missing dependency kinds: {sorted(dependency_kinds)}")
        failures += 1

    expected_test_paths = {path.relative_to(repo_root).as_posix() for path in target_files if is_test_code_file(path)}
    indexed_test_paths = {str(entry.get("test")) for entry in code_map.get("tests", []) if isinstance(entry, dict)}
    if indexed_test_paths == expected_test_paths:
        report_pass(f"PASS code map indexes tests: {len(indexed_test_paths)}")
    else:
        missing_tests = sorted(expected_test_paths - indexed_test_paths)
        extra_tests = sorted(indexed_test_paths - expected_test_paths)
        report_fail(f"FAIL code map test mismatch: missing={missing_tests[:10]}, extra={extra_tests[:10]}")
        failures += 1

    return failures


def main() -> int:
    parser = build_parser("Run the structure harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    failures = check_index_map(repo_root)
    failures += check_docs_coverage(repo_root)
    failures += check_code_map(repo_root)
    return finalize_failures("Structure harness", failures)


if __name__ == "__main__":
    sys.exit(main())
