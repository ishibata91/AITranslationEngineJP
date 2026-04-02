from __future__ import annotations

import sys
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    find_command,
    iter_files,
    load_json,
    report_fail,
    report_pass,
    report_run,
    report_skip,
    run_command,
)


def invoke_step(command: str, arguments: list[str], working_directory: Path) -> int:
    if find_command(command) is None:
        report_fail(f"FAIL missing command: {command}")
        return 1

    rendered_command = f"{command} {' '.join(arguments)}"
    report_run(f"RUN {rendered_command}")
    exit_code = run_command(command, arguments, working_directory)
    if exit_code != 0:
        report_fail(f"FAIL {rendered_command}")
        return 1

    report_pass(f"PASS {rendered_command}")
    return 0


def resolve_package_manager(directory: Path) -> str:
    if (directory / "pnpm-lock.yaml").exists():
        return "pnpm"
    if (directory / "package-lock.json").exists():
        return "npm"
    if (directory / "yarn.lock").exists():
        return "yarn"
    return "npm"


def has_script(package: dict, script_name: str) -> bool:
    scripts = package.get("scripts")
    return isinstance(scripts, dict) and script_name in scripts


def main() -> int:
    parser = build_parser("Run the execution harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    failures = 0
    ran_anything = False
    sonar_scanned = False
    root_execution_gate_ran = False
    root_package: dict | None = None

    sonar_project_properties = repo_root / "sonar-project.properties"
    root_package_json_path = repo_root / "package.json"
    root_cargo_toml_path = repo_root / "src-tauri/Cargo.toml"

    package_jsons = iter_files(repo_root, "package.json")
    cargo_tomls = iter_files(repo_root, "Cargo.toml")

    root_package_json = next((path for path in package_jsons if path == root_package_json_path), None)
    if root_package_json is not None:
        root_package = load_json(root_package_json)
        if has_script(root_package, "gate:execution"):
            ran_anything = True
            root_execution_gate_ran = True
            package_manager = resolve_package_manager(repo_root)
            failures += invoke_step(package_manager, ["run", "gate:execution"], repo_root)
            if has_script(root_package, "scan:sonar"):
                sonar_scanned = True

    for cargo_toml in cargo_tomls:
        if root_execution_gate_ran and cargo_toml == root_cargo_toml_path:
            continue
        ran_anything = True
        cargo_dir = cargo_toml.parent
        failures += invoke_step("cargo", ["fmt", "--all", "--check"], cargo_dir)
        failures += invoke_step("cargo", ["clippy", "--all-targets", "--all-features", "--", "-D", "warnings"], cargo_dir)
        failures += invoke_step("cargo", ["test", "--all-features"], cargo_dir)

    for package_json in package_jsons:
        if root_execution_gate_ran and package_json == root_package_json_path:
            continue

        ran_anything = True
        package_dir = package_json.parent
        package_manager = resolve_package_manager(package_dir)
        package = load_json(package_json)
        scripts_to_run = [script for script in ("lint", "test", "build") if has_script(package, script)]

        if not scripts_to_run:
            report_skip(f"SKIP no lint/test/build scripts in {package_json}")
            continue

        for script_name in scripts_to_run:
            failures += invoke_step(package_manager, ["run", script_name], package_dir)
            if (
                not sonar_scanned
                and script_name == "lint"
                and sonar_project_properties.exists()
                and package_dir.resolve() == repo_root
            ):
                if has_script(package, "scan:sonar"):
                    failures += invoke_step(package_manager, ["run", "scan:sonar"], repo_root)
                else:
                    failures += invoke_step("sonar-scanner", [], repo_root)
                sonar_scanned = True

    has_package_lint = any(has_script(load_json(package_json), "lint") for package_json in package_jsons)
    if sonar_project_properties.exists() and not sonar_scanned and not has_package_lint:
        ran_anything = True
        if root_package_json is not None:
            root_package = root_package or load_json(root_package_json)
            if has_script(root_package, "scan:sonar"):
                failures += invoke_step(resolve_package_manager(repo_root), ["run", "scan:sonar"], repo_root)
                sonar_scanned = True
        if not sonar_scanned:
            failures += invoke_step("sonar-scanner", [], repo_root)
            sonar_scanned = True

    if not ran_anything:
        report_skip(
            "SKIP no Cargo.toml, package.json, or sonar-project.properties targets found. Execution harness is installed but has nothing to run yet.",
        )

    return finalize_failures("Execution harness", failures)


if __name__ == "__main__":
    sys.exit(main())
