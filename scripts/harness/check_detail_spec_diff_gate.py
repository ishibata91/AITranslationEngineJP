from __future__ import annotations

import subprocess
import sys
from pathlib import Path

from harness_common import build_parser, default_repo_root, finalize_failures, report_fail, report_pass, report_skip


def candidate_detail_spec_diff_files(repo_root: Path) -> list[Path]:
    active_root = repo_root / "docs" / "exec-plans" / "active"
    if not active_root.exists():
        return []
    return sorted(active_root.glob("*/detail-spec-diff.md"))


def is_waiting_for_human_answer(path: Path) -> bool:
    text = path.read_text(encoding="utf-8")
    return "`status`: stopped-for-human-answer" in text or "`status`: `stopped-for-human-answer`" in text


def run_gate(repo_root: Path, detail_spec_diff_path: Path) -> int:
    script = repo_root / "scripts" / "detail_spec" / "diff_gate.py"
    report_path = detail_spec_diff_path.with_suffix(".gate.md")
    completed = subprocess.run(
        [
            sys.executable,
            str(script),
            str(detail_spec_diff_path),
            "--report-out",
            str(report_path),
        ],
        cwd=repo_root,
        check=False,
        text=True,
        capture_output=True,
    )
    if completed.returncode == 0:
        report_pass(f"PASS detail spec diff gate: {detail_spec_diff_path.relative_to(repo_root)}")
    else:
        report_fail(f"FAIL detail spec diff gate: {detail_spec_diff_path.relative_to(repo_root)}")
        if completed.stdout:
            print(completed.stdout)
        if completed.stderr:
            print(completed.stderr, file=sys.stderr)
    return completed.returncode


def main() -> int:
    parser = build_parser("Run the detail spec diff gate.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    detail_spec_diff_files = candidate_detail_spec_diff_files(repo_root)
    if not detail_spec_diff_files:
        report_skip("SKIP no active detail-spec-diff.md files")
        return finalize_failures("Detail spec diff gate", 0)

    failures = 0
    for detail_spec_diff_path in detail_spec_diff_files:
        if is_waiting_for_human_answer(detail_spec_diff_path):
            report_skip(f"SKIP stopped for human answer: {detail_spec_diff_path.relative_to(repo_root)}")
            continue
        if run_gate(repo_root, detail_spec_diff_path) != 0:
            failures += 1

    return finalize_failures("Detail spec diff gate", failures)


if __name__ == "__main__":
    sys.exit(main())
