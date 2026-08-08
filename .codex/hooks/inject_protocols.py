#!/usr/bin/env python3
"""Inject protocol documents that apply to files in an apply_patch request."""

from __future__ import annotations

import json
import re
import subprocess
import sys
from pathlib import Path


PATCH_FILE_PATTERN = re.compile(
    r"^\*\*\* (?:Add|Delete|Update) File: (.+)$", re.MULTILINE
)
MAX_CONTEXT_CHARACTERS = 24_000


def repository_root(cwd: str) -> Path | None:
    try:
        completed = subprocess.run(
            ["git", "-C", cwd, "rev-parse", "--show-toplevel"],
            check=True,
            capture_output=True,
            text=True,
        )
    except (OSError, subprocess.CalledProcessError):
        return None
    return Path(completed.stdout.strip()).resolve()


def patch_paths(command: str, repo_root: Path) -> list[Path]:
    paths: list[Path] = []
    for raw_path in PATCH_FILE_PATTERN.findall(command):
        candidate = Path(raw_path)
        if candidate.is_absolute():
            candidate = candidate.resolve()
        else:
            candidate = (repo_root / candidate).resolve()
        try:
            candidate.relative_to(repo_root)
        except ValueError:
            continue
        if candidate not in paths:
            paths.append(candidate)
    return paths


def protocol_files(target: Path, repo_root: Path, protocols_root: Path) -> list[Path]:
    try:
        relative_parent = target.parent.relative_to(repo_root)
    except ValueError:
        return []
    if relative_parent.parts and relative_parent.parts[0] == "protocols":
        return []

    directories = [Path()]
    directories.extend(Path(*relative_parent.parts[:index]) for index in range(1, len(relative_parent.parts) + 1))

    files: list[Path] = []
    for directory in directories:
        for protocol in sorted((protocols_root / directory).glob("*.md")):
            if protocol.is_file() and protocol not in files:
                files.append(protocol)
    return files


def main() -> int:
    try:
        payload = json.load(sys.stdin)
    except json.JSONDecodeError:
        return 0

    if payload.get("tool_name") != "apply_patch":
        return 0
    command = payload.get("tool_input", {}).get("command")
    cwd = payload.get("cwd")
    if not isinstance(command, str) or not isinstance(cwd, str):
        return 0

    repo_root = repository_root(cwd)
    if repo_root is None:
        return 0
    protocols_root = repo_root / "protocols"
    if not protocols_root.is_dir():
        return 0

    files: list[Path] = []
    for target in patch_paths(command, repo_root):
        for protocol in protocol_files(target, repo_root, protocols_root):
            if protocol not in files:
                files.append(protocol)
    if not files:
        return 0

    sections: list[str] = ["Apply the following directory protocols before editing:"]
    used = len(sections[0])
    for protocol in files:
        try:
            content = protocol.read_text(encoding="utf-8").strip()
        except OSError:
            continue
        if not content:
            continue
        section = f"\n\n--- {protocol.relative_to(repo_root)} ---\n{content}"
        if used + len(section) > MAX_CONTEXT_CHARACTERS:
            sections.append("\n\n[Protocol context truncated at 24000 characters. Read the remaining protocol files before editing.]")
            break
        sections.append(section)
        used += len(section)

    if len(sections) == 1:
        return 0
    print(json.dumps({
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "additionalContext": "".join(sections),
        }
    }, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
