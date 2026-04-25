from __future__ import annotations

import subprocess
import sys
from pathlib import Path

from harness_common import build_parser, default_repo_root, finalize_failures, report_fail, report_pass, report_skip


def candidate_scenario_files(repo_root: Path) -> list[Path]:
    active_root = repo_root / "docs" / "exec-plans" / "active"
    if not active_root.exists():
        return []
    return sorted(active_root.glob("*/scenario-design.md"))


def has_requirement_coverage(path: Path) -> bool:
    coverage_path = path.with_suffix(".requirement-coverage.json")
    if coverage_path.exists():
        return True
    return "```json requirement-coverage" in path.read_text(encoding="utf-8")


def run_gate(repo_root: Path, scenario_path: Path) -> int:
    script = repo_root / "scripts" / "scenario" / "requirement_gate.py"
    report_path = scenario_path.with_suffix(".requirement-gate.md")
    questionnaire_path = scenario_path.with_suffix(".questions.md")
    completed = subprocess.run(
        [
            sys.executable,
            str(script),
            str(scenario_path),
            "--report-out",
            str(report_path),
            "--questionnaire-out",
            str(questionnaire_path),
        ],
        cwd=repo_root,
        check=False,
        text=True,
        capture_output=True,
    )
    if completed.returncode == 0:
        report_pass(f"PASS scenario requirement gate: {scenario_path.relative_to(repo_root)}")
    else:
        report_fail(f"FAIL scenario requirement gate: {scenario_path.relative_to(repo_root)}")
        if completed.stdout:
            print(completed.stdout)
        if completed.stderr:
            print(completed.stderr, file=sys.stderr)
    return completed.returncode


def main() -> int:
    parser = build_parser("Run the scenario requirement coverage gate.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    scenario_files = candidate_scenario_files(repo_root)
    if not scenario_files:
        report_skip("SKIP no active scenario-design.md files")
        return finalize_failures("Scenario requirement gate", 0)

    failures = 0
    for scenario_path in scenario_files:
        if not has_requirement_coverage(scenario_path):
            report_fail(f"FAIL missing requirement coverage JSON: {scenario_path.with_suffix('.requirement-coverage.json').relative_to(repo_root)}")
            failures += 1
            continue
        if run_gate(repo_root, scenario_path) != 0:
            failures += 1

    return finalize_failures("Scenario requirement gate", failures)


if __name__ == "__main__":
    sys.exit(main())
