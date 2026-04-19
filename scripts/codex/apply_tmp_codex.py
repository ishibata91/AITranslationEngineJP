#!/usr/bin/env python3
from __future__ import annotations

import argparse
import difflib
import hashlib
import json
import shutil
import subprocess
import sys
import tomllib
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
STAGING_ROOT = REPO_ROOT / "tmp" / "codex"
FILES_ROOT = STAGING_ROOT / "files"
DELETE_PATHS_FILE = STAGING_ROOT / "delete-paths.txt"
DELETION_RATIONALE_FILE = STAGING_ROOT / "deletion-rationale.md"
ALLOWED_ROOT = ".codex/"


def fail(message: str) -> None:
    print(message, file=sys.stderr)
    raise SystemExit(1)


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def validate_repo_relative_path(relative_path: str) -> None:
    path = Path(relative_path)
    if path.is_absolute():
        fail(f"absolute path is not allowed: {relative_path}")
    if ".." in path.parts:
        fail(f"path traversal is not allowed: {relative_path}")
    if not relative_path.startswith(ALLOWED_ROOT):
        fail(f"only .codex paths are allowed: {relative_path}")


def read_text_lines(path: Path) -> list[str]:
    return path.read_text(encoding="utf-8").splitlines(keepends=True)


def collect_staged_files() -> list[tuple[str, Path, Path]]:
    if not FILES_ROOT.exists():
        return []

    staged: list[tuple[str, Path, Path]] = []
    for source_path in sorted(path for path in FILES_ROOT.rglob("*") if path.is_file()):
        if source_path.is_symlink():
            fail(f"symlink staged source is not allowed: {source_path}")
        relative_path = source_path.relative_to(FILES_ROOT).as_posix()
        validate_repo_relative_path(relative_path)
        staged.append((relative_path, source_path, REPO_ROOT / relative_path))
    return staged


def collect_delete_paths() -> list[tuple[str, Path]]:
    if not DELETE_PATHS_FILE.exists():
        return []

    delete_paths: list[tuple[str, Path]] = []
    for raw_line in DELETE_PATHS_FILE.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        validate_repo_relative_path(line)
        delete_paths.append((line, REPO_ROOT / line))
    return delete_paths


def rationale_text() -> str:
    if not DELETION_RATIONALE_FILE.exists():
        return ""
    return DELETION_RATIONALE_FILE.read_text(encoding="utf-8")


def require_rationale(relative_path: str, reason: str, rationale: str) -> None:
    if relative_path in rationale:
        return
    fail(
        "\n".join(
            [
                f"{reason}: {relative_path}",
                f"write {DELETION_RATIONALE_FILE.relative_to(REPO_ROOT)} with the target path, reason, and replacement reference",
            ]
        )
    )


def removed_content_lines(diff_lines: list[str]) -> list[str]:
    removed: list[str] = []
    for line in diff_lines:
        if not line.startswith("-"):
            continue
        if line.startswith("---"):
            continue
        content = line[1:].strip()
        if content:
            removed.append(content)
    return removed


def print_diff_and_check_deletions(relative_path: str, destination_path: Path, source_path: Path, rationale: str) -> None:
    if not destination_path.exists():
        print(f"new file: {relative_path}")
        return

    diff_lines = list(
        difflib.unified_diff(
            read_text_lines(destination_path),
            read_text_lines(source_path),
            fromfile=str(destination_path.relative_to(REPO_ROOT)),
            tofile=str(source_path.relative_to(REPO_ROOT)),
        )
    )
    if not diff_lines:
        print(f"unchanged: {relative_path}")
        return

    print(f"diff: {relative_path}")
    sys.stdout.writelines(diff_lines)
    removed_lines = removed_content_lines(diff_lines)
    if removed_lines:
        print("removed lines detected:")
        for line in removed_lines:
            print(f"  - {line}")
        require_rationale(relative_path, "content deletion detected without rationale", rationale)


def validate_syntax(relative_path: str, source_path: Path) -> None:
    suffix = source_path.suffix.lower()
    if suffix == ".json":
        with source_path.open("r", encoding="utf-8") as handle:
            json.load(handle)
    elif suffix == ".toml":
        with source_path.open("rb") as handle:
            tomllib.load(handle)
    elif suffix == ".puml":
        plantuml = shutil.which("plantuml")
        if not plantuml:
            print(f"skip PlantUML syntax check, command not found: {relative_path}")
            return
        completed = subprocess.run(
            [plantuml, "--check-syntax", "--no-error-image", str(source_path)],
            cwd=REPO_ROOT,
            check=False,
            text=True,
            capture_output=True,
        )
        if completed.returncode != 0:
            sys.stdout.write(completed.stdout)
            sys.stderr.write(completed.stderr)
            fail(f"PlantUML syntax check failed: {relative_path}")


def copy_staged_files(staged_files: list[tuple[str, Path, Path]], check_only: bool, rationale: str) -> None:
    source_hashes: dict[str, str] = {}
    for relative_path, source_path, destination_path in staged_files:
        source_hashes[relative_path] = sha256_file(source_path)
        validate_syntax(relative_path, source_path)
        print_diff_and_check_deletions(relative_path, destination_path, source_path, rationale)

    if check_only:
        for relative_path, _, _ in staged_files:
            print(f"check-only: {relative_path}")
    else:
        for relative_path, source_path, destination_path in staged_files:
            destination_path.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(source_path, destination_path)
            print(f"updated: {relative_path}")

    for relative_path, source_path, _ in staged_files:
        if source_hashes[relative_path] != sha256_file(source_path):
            fail(f"staged source mutated during apply: {relative_path}")


def delete_requested_paths(delete_paths: list[tuple[str, Path]], check_only: bool, rationale: str) -> None:
    for relative_path, target_path in delete_paths:
        require_rationale(relative_path, "file deletion requested without rationale", rationale)
        if target_path.exists() and target_path.is_dir():
            fail(f"directory deletion is not allowed: {relative_path}")
        if check_only:
            print(f"check-only delete: {relative_path}")
            continue
        if target_path.exists():
            target_path.unlink()
            print(f"deleted: {relative_path}")
        else:
            print(f"already absent: {relative_path}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Apply staged tmp/codex files to .codex.")
    parser.add_argument("--check-only", action="store_true", help="run checks without copying, deleting, or cleaning tmp/codex")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    staged_files = collect_staged_files()
    delete_paths = collect_delete_paths()
    rationale = rationale_text()

    if not staged_files and not delete_paths:
        fail("nothing to apply: put files under tmp/codex/files/.codex or paths in tmp/codex/delete-paths.txt")

    copy_staged_files(staged_files, args.check_only, rationale)
    delete_requested_paths(delete_paths, args.check_only, rationale)

    if args.check_only:
        print("tmp/codex apply check passed")
    else:
        shutil.rmtree(STAGING_ROOT)
        print("tmp/codex applied and removed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
