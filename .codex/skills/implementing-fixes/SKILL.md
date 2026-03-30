---
name: implementing-fixes
description: bugfix scope を直接修正し、必要な checks と残留リスクを返す。
---

# Implementing Fixes

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- active fix plan と trace 結果を読んでから編集する
- narrow fix を優先する
- temporary logging cleanup が必要なら明示する
- 指定 checks を実行する
- broad refactor を混ぜない

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.implementing-fixes.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/implementing-fixes.to.directing-fixes.json` を返却契約として使う。
