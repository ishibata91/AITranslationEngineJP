from __future__ import annotations

import re
import sys
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    load_json,
    report_check,
    read_text,
)

CHECKS = [
    { "file": ".codex/README.md", "patterns": ["directing-implementation", "designing-implementation", "directing-fixes", "architecting-tests", "UI", "Scenario", "Logic", "single-pass", "changes/", "tasks.md"] },
    { "file": ".codex/agents/task_designer.toml", "patterns": ["UI", "Scenario", "Logic", "task-local design", "tasks.md"] },
    { "file": ".codex/agents/ctx_loader.toml", "patterns": ["facts", "constraints", "gaps", "exec-plan"] },
    { "file": ".codex/agents/workplan_builder.toml", "patterns": ["UI", "Scenario", "Logic", "Implementation Plan", "tasks.md"] },
    { "file": ".codex/agents/test_architect.toml", "patterns": ["test_architect", "spec-aligned tests", "fixtures", "validation commands"] },
    { "file": ".codex/agents/implementer.toml", "patterns": ["allowed scope", "validation", "touched files"] },
    { "file": ".codex/agents/fault_tracer.toml", "patterns": ["root-cause hypotheses", "observation plan", "temporary logging"] },
    { "file": ".codex/agents/log_instrumenter.toml", "patterns": ["temporary log statements", r"\[tracing-fixes\]", "Remove any temporary instrumentation"] },
    { "file": ".codex/agents/review_cycler.toml", "patterns": ["single-pass", "spec deviation", "exception handling", "resource cleanup", "missing tests", "pass", "reroute"] },
    { "file": ".codex/skills/directing-implementation/SKILL.md", "patterns": ["docs/exec-plans/templates/impl-plan.md", "designing-implementation", "UI", "Scenario", "Logic", "architecting-tests", "reviewing-implementation", "reroute", "changes/", "tasks.md"] },
    { "file": ".codex/skills/designing-implementation/SKILL.md", "patterns": ["UI", "Scenario", "Logic", "task-local design", "changes/", "tasks.md"] },
    { "file": ".codex/skills/directing-fixes/SKILL.md", "patterns": ["docs/exec-plans/templates/fix-plan.md", "distilling-fixes", "tracing-fixes", "logging-fixes", "architecting-tests", "reviewing-fixes", "tasks.md"] },
    { "file": ".codex/skills/architecting-tests/SKILL.md", "patterns": ["Acceptance Checks", "failing tests", "fixtures", "validation commands", "closeout notes", "updating-docs"] },
    { "file": ".codex/skills/reviewing-implementation/SKILL.md", "patterns": ["Review Scope", "pass", "reroute"] },
    { "file": ".codex/skills/reviewing-fixes/SKILL.md", "patterns": ["Review Scope", "pass", "reroute"] },
    { "file": "docs/exec-plans/templates/impl-plan.md", "patterns": ["Decision Basis", "UI", "Scenario", "Logic", "Implementation Plan", "Required Evidence", "4humans Sync", "Outcome"] },
    { "file": "docs/exec-plans/templates/fix-plan.md", "patterns": ["Decision Basis", "Known Facts", "Trace Plan", "Fix Plan", "Required Evidence", "4humans Sync", "Outcome"] },
    { "file": "docs/spec.md", "patterns": ["LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator"] },
    { "file": "docs/architecture.md", "patterns": ["Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust"] },
    { "file": "docs/tech-selection.md", "patterns": ["Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx"] },
    { "file": "docs/core-beliefs.md", "patterns": ["agent-first", "directing-implementation", "directing-fixes", "single-pass", "evidence"] },
    { "file": "docs/index.md", "patterns": ["AGENTS.md", ".codex", "directing-implementation", "directing-fixes", "designing-implementation", "quality-score.md", "overview-manifest.json", "tests / acceptance checks / validation commands"] },
    { "file": "4humans/quality-score.md", "patterns": ["Codex workflow source of truth", "Design harness", "directing-implementation", "directing-fixes"] },
    { "file": "4humans/diagrams/overview-manifest.json", "patterns": ["structures", "processes", "detail_to_overviews", "overviews"] },
]

SEMANTIC_PATH_CHECKS = [
    {"path": "src/application", "should_exist": True, "label": "frontend application layer root"},
    {"path": "src/gateway", "should_exist": True, "label": "frontend gateway root"},
    {"path": "src/shared", "should_exist": True, "label": "frontend shared DTO root"},
    {"path": "src/ui", "should_exist": True, "label": "frontend ui root"},
    {"path": "src/domain", "should_exist": False, "label": "frontend domain forbidden during bootstrap"},
    {"path": "src/infra", "should_exist": False, "label": "frontend infra forbidden during bootstrap"},
    {"path": "src-tauri/src/application", "should_exist": True, "label": "backend application layer root"},
    {"path": "src-tauri/src/domain", "should_exist": True, "label": "backend domain layer root"},
    {"path": "src-tauri/src/infra", "should_exist": True, "label": "backend infra layer root"},
    {"path": "src-tauri/src/gateway", "should_exist": True, "label": "backend gateway root"},
]

SEMANTIC_PHRASE_CHECKS = [
    {"file": "docs/core-beliefs.md", "expected": "tests / acceptance checks / validation commands", "forbidden": ["test / acceptance checks / validation commands"]},
    {"file": "docs/architecture.md", "expected": "tests / acceptance checks / validation commands", "forbidden": []},
    {"file": "docs/index.md", "expected": "tests / acceptance checks / validation commands", "forbidden": []},
]


def assert_patterns(file_path: Path, patterns: list[str]) -> int:
    failures = 0
    content = read_text(file_path)
    for pattern in patterns:
        failures += report_check(
            re.search(pattern, content) is not None,
            f"PASS pattern '{pattern}' found in {file_path}",
            f"FAIL missing pattern '{pattern}' in {file_path}",
        )
    return failures


def assert_path_state(path: Path, should_exist: bool, label: str) -> int:
    exists = path.exists()
    state = "exists" if should_exist else "is absent"
    expected = "exist" if should_exist else "be absent"
    return report_check(
        exists == should_exist,
        f"PASS semantic check '{label}': '{path}' {state}",
        f"FAIL semantic check '{label}': expected '{path}' to {expected}",
    )


def assert_canonical_phrase(file_path: Path, expected_phrase: str, forbidden_phrases: list[str]) -> int:
    failures = 0
    content = read_text(file_path)
    failures += report_check(
        expected_phrase in content,
        f"PASS semantic phrase '{expected_phrase}' found in {file_path}",
        f"FAIL semantic phrase '{expected_phrase}' missing in {file_path}",
    )

    for forbidden_phrase in forbidden_phrases:
        failures += report_check(
            forbidden_phrase not in content,
            f"PASS semantic phrase '{forbidden_phrase}' absent from {file_path}",
            f"FAIL semantic phrase '{forbidden_phrase}' should not appear in {file_path}",
        )
    return failures


def assert_overview_manifest(repo_root: Path) -> int:
    failures = 0
    manifest_path = repo_root / "4humans/diagrams/overview-manifest.json"
    manifest = load_json(manifest_path)

    for section_name in ("structures", "processes"):
        section = manifest.get(section_name)
        if not isinstance(section, dict):
            failures += report_check(
                False,
                f"PASS overview manifest section '{section_name}' exists",
                f"FAIL overview manifest section '{section_name}' is missing or invalid in {manifest_path}",
            )
            continue

        overviews = section.get("overviews")
        detail_to_overviews = section.get("detail_to_overviews")
        if not isinstance(overviews, list) or not isinstance(detail_to_overviews, dict):
            failures += report_check(
                False,
                f"PASS overview manifest shape for '{section_name}' is valid",
                f"FAIL overview manifest '{section_name}' must contain list 'overviews' and object 'detail_to_overviews' in {manifest_path}",
            )
            continue

        expected_dir = repo_root / "4humans/diagrams" / section_name
        overview_files = sorted(expected_dir.glob("*overview*.d2"))
        declared_overviews = {repo_root / path for path in overviews}

        for overview_file in overview_files:
            failures += report_check(
                overview_file in declared_overviews,
                f"PASS overview '{overview_file}' is registered in {manifest_path}",
                f"FAIL overview '{overview_file}' is not registered in {manifest_path}",
            )

        for declared_overview in declared_overviews:
            failures += report_check(
                declared_overview.exists(),
                f"PASS declared overview '{declared_overview}' exists",
                f"FAIL declared overview '{declared_overview}' does not exist",
            )

        detail_files = [
            path
            for path in sorted(expected_dir.glob("*.d2"))
            if "overview" not in path.name
        ]
        declared_detail_files = {repo_root / path for path in detail_to_overviews}

        for detail_file in detail_files:
            failures += report_check(
                detail_file in declared_detail_files,
                f"PASS detail diagram '{detail_file}' is registered in {manifest_path}",
                f"FAIL detail diagram '{detail_file}' is not registered in {manifest_path}",
            )

        for detail_path_str, overview_refs in sorted(detail_to_overviews.items()):
            detail_path = repo_root / detail_path_str
            failures += report_check(
                detail_path.exists(),
                f"PASS declared detail diagram '{detail_path}' exists",
                f"FAIL declared detail diagram '{detail_path}' does not exist",
            )
            failures += report_check(
                isinstance(overview_refs, list) and len(overview_refs) > 0,
                f"PASS detail diagram '{detail_path}' has linked overview targets",
                f"FAIL detail diagram '{detail_path}' must declare at least one overview target in {manifest_path}",
            )
            if not isinstance(overview_refs, list):
                continue
            for overview_ref in overview_refs:
                overview_path = repo_root / overview_ref
                failures += report_check(
                    overview_path in declared_overviews,
                    f"PASS detail diagram '{detail_path}' points to registered overview '{overview_path}'",
                    f"FAIL detail diagram '{detail_path}' points to undeclared overview '{overview_path}' in {manifest_path}",
                )

    return failures


def main() -> int:
    parser = build_parser("Run the design harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    failures = 0

    for check in CHECKS:
        failures += assert_patterns(repo_root / check["file"], check["patterns"])

    for check in SEMANTIC_PATH_CHECKS:
        failures += assert_path_state(repo_root / check["path"], check["should_exist"], check["label"])

    for check in SEMANTIC_PHRASE_CHECKS:
        failures += assert_canonical_phrase(repo_root / check["file"], check["expected"], check["forbidden"])

    failures += assert_overview_manifest(repo_root)

    return finalize_failures("Design harness", failures)


if __name__ == "__main__":
    sys.exit(main())
