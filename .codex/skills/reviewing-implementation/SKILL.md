---
name: reviewing-implementation
description: 実装差分を単発で照合し、`pass` か `reroute` を返す。
---

# Reviewing Implementation

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

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.reviewing-implementation.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/reviewing-implementation.to.directing-implementation.json` を返却契約として使う。
