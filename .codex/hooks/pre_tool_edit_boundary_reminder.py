#!/usr/bin/env python3
import json
import os
import sys
from pathlib import Path

REPO_ROOT = Path("/Users/iorishibata/Repositories/AITranslationEngineJP")

REMINDER = """ファイル編集前に、責務境界と変更可否を確認する。必要な場合だけ skill / agent を読む。"""


def main() -> int:
    payload = read_payload()
    cwd = Path(payload.get("cwd") or os.getcwd()).resolve()

    if not is_under_repo(cwd):
        print(json.dumps({"continue": True}, ensure_ascii=False))
        return 0

    output = {
        "continue": True,
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "additionalContext": REMINDER,
        },
    }
    print(json.dumps(output, ensure_ascii=False))
    return 0


def read_payload() -> dict:
    try:
        raw = sys.stdin.read()
        if not raw.strip():
            return {}
        value = json.loads(raw)
        return value if isinstance(value, dict) else {}
    except json.JSONDecodeError:
        return {}


def is_under_repo(path: Path) -> bool:
    try:
        path.relative_to(REPO_ROOT)
        return True
    except ValueError:
        return path == REPO_ROOT


if __name__ == "__main__":
    raise SystemExit(main())
