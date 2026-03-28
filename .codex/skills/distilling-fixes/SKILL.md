---
name: distilling-fixes
description: bugfix 初期情報を整理し、既知事実、再現条件、関連仕様、関連コードを短くまとめる。
---

# Distilling Fixes

## Output

- known facts
- reproduction status
- related constraints
- related code pointers
- open gaps

## Rules

- 事実と推測を分ける
- packet file を作らない
- `changes/` や `context_board` を前提にしない

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.distilling-fixes.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/distilling-fixes.to.directing-fixes.json` を返却契約として使う。
