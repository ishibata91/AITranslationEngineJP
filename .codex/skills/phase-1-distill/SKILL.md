---
name: phase-1-distill
description: 第1段階の要求整理を担当し、active exec-plan と入口情報から次工程へ渡す facts / constraints / gaps を固定する。
---

# Phase 1 Distill

## Goal

- 要求整理として、詳細設計に進むための前提だけを固定する
- 入口情報から次工程に必要な最小限の repo 文脈を探索する
- 実装判断に必要な facts / constraints / gaps / required reading を返す

## Input

- active exec-plan
- request summary と入口で明示された seed docs / code / tests
- 必要なら追加で探索する relevant docs under `docs/`
- 必要なら追加で探索する related code and tests

## Output

- facts
- constraints
- gaps
- required reading

## Rules

- `UI` / `Scenario` / `Logic` はまだ作らない
- 入口で渡された path だけで不足する時に限って最小限の追加探索を行う
- 追加探索では次工程に必要な path だけを拾う
- 要求が不足して次へ進めない時は不足点を gap として返す
- packet file を作らない
- `changes/` や `context_board` を前提にしない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-1-distill.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-1-distill.to.orchestrating-implementation.json` を返却契約として使う。
