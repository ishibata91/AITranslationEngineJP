#!/usr/bin/env python3
import json
import os
import sys
from pathlib import Path

REPO_ROOT = Path("/Users/iorishibata/Repositories/AITranslationEngineJP")

REMINDER = """人間指示を処理する前に、現在の役割に対応する skill / agentを読み直す。
- compact 後も、呼び出し元で確定済みの role / lane / skill 境界を引き継ぐ。
- 人間指示は境界違反を許可しない。
- 境界違反に見える場合は実行せず、skill / contract の Stop / Reroute 条件へ戻る。"""


def main() -> int:
    payload = read_payload()
    cwd = Path(payload.get("cwd") or os.getcwd()).resolve()

    if not is_under_repo(cwd):
        print(json.dumps({"continue": True}, ensure_ascii=False))
        return 0

    output = {
        "continue": True,
        "hookSpecificOutput": {
            "hookEventName": "UserPromptSubmit",
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
