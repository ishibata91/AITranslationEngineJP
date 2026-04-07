---
name: designing-implementation
description: active exec-plan の `UI` / `Scenario` / `Logic` を task-local design として固める。
---

# Designing Implementation

## Goal

- task-local design が必要な task だけ、direction から渡された context summary を前提に active exec-plan の `UI` / `Scenario` / `Logic` を埋める
- downstream skill が読める粒度まで設計判断を短く固定する
- task-local design を active exec-plan の外へ逃がさない

## Rules

- 更新対象は原則として active exec-plan の `UI` / `Scenario` / `Logic` だけに限定する
- direction から渡された facts / constraints / gaps / required reading を起点に task-local design を固める
- section が不要な task では `N/A` を維持し、不要な設計を増やさない
- 実装コード、`Implementation Plan` の詳細、tests の詳細設計は持たない
- `changes/`、`context_board`、`tasks.md`、別の design artifact を作らない
- repo の恒久仕様や境界が不足していて task-local design を安全に決められない時は停止して direction へ返す

## Reference Use

- 着手前に `../proposing-implementation/references/proposing-implementation.to.designing-implementation.json` を参照して入力契約を確認する。
- `proposing-implementation` へ返す時は `references/designing-implementation.to.proposing-implementation.json` を返却契約として使う。
