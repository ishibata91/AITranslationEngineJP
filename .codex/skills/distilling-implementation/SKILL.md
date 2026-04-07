---
name: distilling-implementation
description: 実装 lane 用。active exec-plan、関連 docs、コードから実装判断に必要な facts / constraints / gaps を抽出する。
---

# Distilling Implementation

## Goal

- active exec-plan と入口情報を起点に、次工程に必要な最小限の repo 文脈を探索する
- 実装判断に必要な facts / constraints / gaps を分離して返す
- `designing-implementation` と `planning-implementation` が読むべき path を絞る
- 4humans/diagrams/以下のd2で、実装工程で変更が必要な図を特定する

## Input

- active exec-plan
- request summary と入口で明示された seed docs / code / tests
- 必要なら追加で探索する relevant docs under `docs/`
- 必要なら追加で探索する related code and tests

## Output

- facts
- constraints
- gaps
- closeout notes
- required reading

## Rules

- `UI` / `Scenario` / `Logic` は active plan を正本として読む
- 入口で渡された path をなぞるだけで終わらず、必要な repo 文脈が不足していれば最小限の追加探索を行う
- 追加探索では「なぜ次工程に必要か」が説明できる path だけを拾う
- packet file を作らない
- `changes/` や `context_board` を前提にしない
- 実装判断が未確定なら gap として返す

## Reference Use

- 着手前に `../proposing-implementation/references/proposing-implementation.to.distilling-implementation.json` を参照して入力契約を確認する。
- `proposing-implementation` へ返す時は `references/distilling-implementation.to.proposing-implementation.json` を返却契約として使う。
