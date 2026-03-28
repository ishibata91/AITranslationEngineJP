---
name: tracing-fixes
description: 原因仮説を順位付けし、最小の trace 方針を返す。
---

# Tracing Fixes

## Output

- hypotheses
- observation points
- whether temporary logging is needed
- next investigation step

## Rules

- 観測計画は狭く保つ
- legacy logger helper を前提にしない
- 恒久修正を混ぜない

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.tracing-fixes.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/tracing-fixes.to.directing-fixes.json` を返却契約として使う。
