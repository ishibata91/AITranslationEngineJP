---
name: reviewing-fixes
description: bugfix 差分を単発で照合し、`pass` か `reroute` を返す。
---

# Reviewing Fixes

## Review Scope

- 仕様逸脱
- 例外処理
- リソース解放
- テスト不足

## Output

- decision: `pass` or `reroute`
- findings
- recheck
- residual_risk
- closeout_notes

## Rules

- review は 1 回だけ行う
- 好みの改善提案を主目的にしない
- 重大不足があれば lane へ差し戻す

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.reviewing-fixes.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/reviewing-fixes.to.directing-fixes.json` を返却契約として使う。
