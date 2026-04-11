from __future__ import annotations

import sys
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    find_command,
    load_json,
    report_fail,
    report_run,
    report_skip,
    run_command,
)


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

    return 0


def main() -> int:
    parser = build_parser("Run the backend test harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    package_json_path = repo_root / "package.json"
    if not package_json_path.exists():
        report_skip(f"SKIP no package.json found at {package_json_path}")
        return finalize_failures("Backend test harness", 0)

    package = load_json(package_json_path)
    if not has_script(package, "test:backend"):
        report_skip(f"SKIP no test:backend script in {package_json_path}")
        return finalize_failures("Backend test harness", 0)

    package_manager = resolve_package_manager(repo_root)
    failures = invoke_step(package_manager, ["run", "test:backend"], repo_root)
    return finalize_failures("Backend test harness", failures)


if __name__ == "__main__":
    sys.exit(main())
