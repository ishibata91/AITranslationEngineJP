from __future__ import annotations

import re
import sys
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


def main() -> int:
    parser = build_parser("Run the structure harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    failures = check_index_map(repo_root)
    failures += check_docs_coverage(repo_root)
    return finalize_failures("Structure harness", failures)


if __name__ == "__main__":
    sys.exit(main())
