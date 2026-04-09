---
name: distilling-fixes
description: bugfix 初期情報を整理し、既知事実、再現条件、関連仕様、関連コードを短くまとめる。
---

# Distilling Fixes

## Goal

- active fix plan と入口情報を起点に、次工程に必要な最小限の repo 文脈を探索する
- 既知事実、再現状況、制約、未解明点を分離して返す
- `tracing-fixes` と後続 skill が読むべき path を絞る

## Input

- active fix plan
- bug summary と既知の再現条件
- 入口で明示された seed docs / code / tests
- 必要なら追加で探索する関連仕様、関連コード、関連 tests

## Output

- known facts
- reproduction status
- related constraints
- related code pointers
- open gaps
- required reading

## Rules

- 事実と推測を分ける
- 入口で渡された path をなぞるだけで終わらず、必要な repo 文脈が不足していれば最小限の追加探索を行う
- 追加探索では「なぜ trace や fix 判断に必要か」が説明できる path だけを拾う
- packet file を作らない
- `changes/` や `context_board` を前提にしない

## Reference Use

- 着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.distilling-fixes.json` を参照して入力契約を確認する。
- `orchestrating-fixes` へ返す時は `references/distilling-fixes.to.orchestrating-fixes.json` を返却契約として使う。
