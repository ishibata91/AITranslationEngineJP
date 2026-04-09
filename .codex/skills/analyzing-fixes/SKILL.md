---
name: analyzing-fixes
description: 観測結果を圧縮し、結果を前SKILLにhandoffするSKILL。
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

- 着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.analyzing-fixes.json` を参照して入力契約を確認する。
- `orchestrating-fixes` へ返す時は `references/analyzing-fixes.to.orchestrating-fixes.json` を返却契約として使う。
