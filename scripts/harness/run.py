from __future__ import annotations

import argparse
import sys
from pathlib import Path

from harness_common import default_repo_root, report_pass, report_section, run_python_script

SUITE_ORDER = {
    "backend-lint": ["check_backend_lint.py"],
    "frontend-lint": ["check_frontend_lint.py"],
    "structure": ["check_structure.py"],
    "design": ["check_design.py"],
    "execution": ["check_execution.py"],
    "all": ["check_structure.py", "check_design.py", "check_execution.py"],
}


def main() -> int:
    parser = argparse.ArgumentParser(description="Run one or more harness suites.")
    parser.add_argument("--suite", choices=sorted(SUITE_ORDER), default="all")
    parser.add_argument("--repo-root", default=str(default_repo_root(__file__)))
    args = parser.parse_args()

    repo_root = Path(args.repo_root).resolve()
    script_root = Path(__file__).resolve().parent

    for script_name in SUITE_ORDER[args.suite]:
        report_section(f"Running {script_name}")
        exit_code = run_python_script(script_root / script_name, repo_root)
        if exit_code != 0:
            return exit_code

    print()
    report_pass("All requested harness suites passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
