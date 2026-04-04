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
- `4humans` D2 sync 要否と実施有無

## Output

- decision: `pass` or `reroute`
- findings
- recheck
- closeout_notes

## Rules

- review は 1 回だけ行う
- 好みや美しさを主目的にしない
- `4humans/diagrams/processes/` と `4humans/diagrams/structures/` の D2 sync は close 後の任意作業ではなく、差分に応じて review 中に要否と実施有無を判定する
- 重大不足があれば lane へ差し戻す

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.reviewing-implementation.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/reviewing-implementation.to.directing-implementation.json` を返却契約として使う。
