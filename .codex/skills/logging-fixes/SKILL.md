---
name: logging-fixes
description: trace に必要な一時観測ログの add/remove だけを行う。
---

# Logging Fixes

## Rules

- 一時観測だけを目的にする
- 既存モジュールの logging style を再利用する
- メッセージには `[tracing-fixes]` を含める
- 恒久修正や logger 設計変更を混ぜない
- 調査後に除去する

## Reference Use

- 着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.logging-fixes.json` を参照して入力契約を確認する。
- `orchestrating-fixes` へ返す時は `references/logging-fixes.to.orchestrating-fixes.json` を返却契約として使う。
