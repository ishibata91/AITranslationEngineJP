from __future__ import annotations

import re
import sys
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    iter_markdown_files,
    report_fail,
    report_pass,
    report_section,
    read_text,
    remove_code_blocks,
)

LINK_PATTERN = re.compile(r"!?\[[^\]]*\]\((?P<target>[^)]+)\)")
URI_SCHEME_PATTERN = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:")
WINDOWS_DRIVE_PATTERN = re.compile(r"^[a-zA-Z]:[\\\\/]")

REQUIRED_PATHS = [
    "AGENTS.md",
    ".codex/README.md",
    ".codex/agents/task_designer.toml",
    ".codex/agents/ctx_loader.toml",
    ".codex/agents/workplan_builder.toml",
    ".codex/agents/test_architect.toml",
    ".codex/agents/implementer.toml",
    ".codex/agents/fault_tracer.toml",
    ".codex/agents/log_instrumenter.toml",
    ".codex/agents/review_cycler.toml",
    ".codex/skills/directing-implementation/SKILL.md",
    ".codex/skills/directing-implementation/agents/openai.yaml",
    ".codex/skills/designing-implementation/SKILL.md",
    ".codex/skills/designing-implementation/agents/openai.yaml",
    ".codex/skills/distilling-implementation/SKILL.md",
    ".codex/skills/distilling-implementation/agents/openai.yaml",
    ".codex/skills/planning-implementation/SKILL.md",
    ".codex/skills/planning-implementation/agents/openai.yaml",
    ".codex/skills/architecting-tests/SKILL.md",
    ".codex/skills/architecting-tests/agents/openai.yaml",
    ".codex/skills/implementing-frontend/SKILL.md",
    ".codex/skills/implementing-frontend/agents/openai.yaml",
    ".codex/skills/implementing-backend/SKILL.md",
    ".codex/skills/implementing-backend/agents/openai.yaml",
    ".codex/skills/reviewing-implementation/SKILL.md",
    ".codex/skills/reviewing-implementation/agents/openai.yaml",
    ".codex/skills/directing-fixes/SKILL.md",
    ".codex/skills/directing-fixes/agents/openai.yaml",
    ".codex/skills/distilling-fixes/SKILL.md",
    ".codex/skills/distilling-fixes/agents/openai.yaml",
    ".codex/skills/tracing-fixes/SKILL.md",
    ".codex/skills/tracing-fixes/agents/openai.yaml",
    ".codex/skills/analyzing-fixes/SKILL.md",
    ".codex/skills/analyzing-fixes/agents/openai.yaml",
    ".codex/skills/logging-fixes/SKILL.md",
    ".codex/skills/logging-fixes/agents/openai.yaml",
    ".codex/skills/implementing-fixes/SKILL.md",
    ".codex/skills/implementing-fixes/agents/openai.yaml",
    ".codex/skills/reviewing-fixes/SKILL.md",
    ".codex/skills/reviewing-fixes/agents/openai.yaml",
    ".codex/skills/reporting-risks/SKILL.md",
    ".codex/skills/reporting-risks/agents/openai.yaml",
    ".codex/skills/updating-docs/SKILL.md",
    ".codex/skills/updating-docs/agents/openai.yaml",
    "docs/index.md",
    "docs/core-beliefs.md",
    "docs/spec.md",
    "docs/architecture.md",
    "docs/tech-selection.md",
    "docs/er.md",
    "4humans/tech-debt-tracker.md",
    "4humans/quality-score.md",
    "docs/references/index.md",
    "docs/exec-plans/active/README.md",
    "docs/exec-plans/completed/README.md",
    "docs/exec-plans/templates/impl-plan.md",
    "docs/exec-plans/templates/fix-plan.md",
    "scripts/harness/run.py",
    "scripts/harness/check_structure.py",
    "scripts/harness/check_design.py",
    "scripts/harness/check_execution.py",
    ".codex/skills/directing-implementation/scripts/get-open-sonar-issues.py",
]


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


def collect_markdown_files(repo_root: Path) -> list[Path]:
    markdown_files: list[Path] = []
    for path in iter_markdown_files(repo_root):
        path_text = path.as_posix()
        if "/docs/exec-plans/completed/" in path_text:
            continue
        if "/.codex/.codex/" in path_text:
            continue
        markdown_files.append(path)
    return markdown_files


def check_required_paths(repo_root: Path) -> int:
    failures = 0
    report_section("Required paths")
    for relative_path in REQUIRED_PATHS:
        path = repo_root / relative_path
        if path.exists():
            report_pass(f"PASS path exists: {path}")
        else:
            report_fail(f"FAIL missing path: {path}")
            failures += 1
    return failures


def check_markdown_links(repo_root: Path) -> int:
    failures = 0
    report_section("Markdown links")
    for file_path in collect_markdown_files(repo_root):
        content = remove_code_blocks(read_text(file_path))
        for match in LINK_PATTERN.finditer(content):
            target = match.group("target").strip()
            if target.startswith("#"):
                continue

            resolved = resolve_link_target(file_path, target)
            if resolved is None:
                continue

            if not resolved.exists():
                report_fail(f"FAIL broken link: {file_path} -> {target}")
                failures += 1
    return failures


def main() -> int:
    parser = build_parser("Run the structure harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    failures = check_required_paths(repo_root)
    failures += check_markdown_links(repo_root)
    return finalize_failures("Structure harness", failures)


if __name__ == "__main__":
    sys.exit(main())
