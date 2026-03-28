---
name: impl-review
description: 実装差分を単発で照合し、`pass` か `reroute` を返す。
---

# Impl Review

## Review Scope

- 仕様逸脱
- 例外処理
- リソース解放
- テスト不足

## Output

- decision: `pass` or `reroute`
- findings
- recheck
- docs_sync_needed

## Rules

- review は 1 回だけ行う
- score を返さない
- 好みや美しさを主目的にしない
- 重大不足があれば lane へ差し戻す
