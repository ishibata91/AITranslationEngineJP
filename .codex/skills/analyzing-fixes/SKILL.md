---
name: analyzing-fixes
description: 観測結果を圧縮し、修正対象か docs sync 対象かを整理する。
---

# Analyzing Fixes

## Output

- observed facts
- disproved hypotheses
- remaining gaps
- recommended next step

## Rules

- ログや観測結果を事実へ圧縮する
- packet file を作らない
- 推測を事実として扱わない

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.analyzing-fixes.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/analyzing-fixes.to.directing-fixes.json` を返却契約として使う。
