from __future__ import annotations

import argparse
import json
import re
import shutil
import subprocess
import sys
from pathlib import Path

EXCLUDED_DIR_NAMES = {".git", ".cargo-home", "node_modules", "dist", "build", "coverage", "target"}
ANSI_COLORS = {
    "cyan": "\033[36m",
    "green": "\033[32m",
    "red": "\033[31m",
    "yellow": "\033[33m",
    "reset": "\033[0m",
}


def default_repo_root(script_file: str) -> Path:
    return Path(script_file).resolve().parents[2]


def build_parser(description: str, default_root: Path) -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=description)
    parser.add_argument("--repo-root", default=str(default_root))
    return parser


def colorize(message: str, color: str) -> str:
    if not sys.stdout.isatty():
        return message
    return f"{ANSI_COLORS[color]}{message}{ANSI_COLORS['reset']}"


def print_status(message: str, color: str) -> None:
    print(colorize(message, color), flush=True)


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def load_json(path: Path) -> dict:
    return json.loads(read_text(path))


def run_command_capture(command: str, arguments: list[str], working_directory: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [command, *arguments],
        cwd=working_directory,
        check=False,
        text=True,
        capture_output=True,
    )


def is_repo_owned(path: Path) -> bool:
    return not any(part in EXCLUDED_DIR_NAMES for part in path.parts)


def iter_files(root: Path, pattern: str) -> list[Path]:
    return sorted(path for path in root.rglob(pattern) if path.is_file() and is_repo_owned(path))


def iter_markdown_files(root: Path) -> list[Path]:
    return iter_files(root, "*.md")


def find_command(command: str) -> str | None:
    return shutil.which(command)


def run_command(command: str, arguments: list[str], working_directory: Path) -> int:
    completed = subprocess.run([command, *arguments], cwd=working_directory, check=False)
    return completed.returncode


def run_python_script(script_path: Path, repo_root: Path) -> int:
    return run_command(sys.executable, [str(script_path), "--repo-root", str(repo_root)], script_path.parent)


def report_section(title: str) -> None:
    print()
    print_status(f"== {title} ==", "cyan")


def report_pass(message: str) -> None:
    print_status(message, "green")


def report_fail(message: str) -> None:
    print_status(message, "red")


def report_run(message: str) -> None:
    print_status(message, "cyan")


def report_skip(message: str) -> None:
    print_status(message, "yellow")


def report_check(passed: bool, success_message: str, failure_message: str) -> int:
    if passed:
        report_pass(success_message)
        return 0

    report_fail(failure_message)
    return 1


def finalize_failures(harness_name: str, failures: int) -> int:
    if failures > 0:
        report_fail(f"{harness_name} failed with {failures} issue(s).")
        return 1

    print()
    report_pass(f"{harness_name} passed.")
    return 0


def remove_code_blocks(text: str) -> str:
    return re.sub(r"```.*?```", "", text, flags=re.DOTALL)
