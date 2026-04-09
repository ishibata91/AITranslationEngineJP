from __future__ import annotations

import sys
from pathlib import Path

from harness_common import build_parser, default_repo_root, finalize_failures, find_command, load_json, report_fail, report_pass, report_run, report_skip, run_command


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
    root_package_json_path = repo_root / "package.json"
    sonar_project_properties = repo_root / "sonar-project.properties"

    if not root_package_json_path.exists():
        report_skip(f"SKIP no package.json found at {root_package_json_path}")
        return finalize_failures("Execution harness", 0)

    root_package = load_json(root_package_json_path)
    package_manager = resolve_package_manager(repo_root)
    ran_anything = False

    if has_script(root_package, "lint:backend"):
        ran_anything = True
        failures += invoke_step(package_manager, ["run", "lint:backend"], repo_root)
    else:
        report_skip(f"SKIP no lint:backend script in {root_package_json_path}")

    if has_script(root_package, "lint:frontend"):
        ran_anything = True
        failures += invoke_step(package_manager, ["run", "lint:frontend"], repo_root)
    else:
        report_skip(f"SKIP no lint:frontend script in {root_package_json_path}")

    if sonar_project_properties.exists():
        ran_anything = True
        if has_script(root_package, "scan:sonar"):
            failures += invoke_step(package_manager, ["run", "scan:sonar"], repo_root)
        else:
            failures += invoke_step("sonar-scanner", [], repo_root)
    else:
        report_skip(f"SKIP no sonar-project.properties found at {sonar_project_properties}")

    if not ran_anything:
        report_skip("SKIP no lint:backend, lint:frontend, or Sonar target found. Execution harness is installed but has nothing to run yet.")

    return finalize_failures("Execution harness", failures)


if __name__ == "__main__":
    sys.exit(main())
