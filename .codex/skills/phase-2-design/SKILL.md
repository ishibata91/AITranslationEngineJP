---
name: phase-2-design
description: 第2段階の詳細設計を担当し、active exec-plan の `UI` / `Scenario` / `Logic` を task-local design として固定する。
---

# Phase 2 Design

## Goal

- `UI`、`Scenario`、`Logic` を詳細設計として固める
- downstream skill が読める粒度まで設計判断を短く固定する
- task-local design を active exec-plan の外へ逃がさない
- `Logic` の review に図が必要な時は、review 用差分図までこの段階で揃える

## Rules

- 更新対象は原則として active exec-plan の `UI` / `Scenario` / `Logic` だけに限定する
- `UI` はモック HTML wireframe として扱い、画面構造と操作配置を固定する
- `Scenario` はホワイトボックステスト一覧として扱い、後続工程の証明対象を固定する
- `Logic` は component responsibility map として扱い、責務、主要な振る舞い、member 対応を固定する
- `Logic` を図で review したい時は、この段階の責務として `diagramming-structure-diff` を使い、review 用差分図と差分正本適用先を揃える
- 設計判断が揺れている間は次工程へ渡さない
- 実装コード、`Implementation Plan`、test file の詳細実装は持たない
- `changes/`、`context_board`、`tasks.md`、別の design artifact を作らない
- repo の恒久仕様や境界が不足していて task-local design を安全に決められない時は停止して orchestrator へ返す

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2-design.to.orchestrating-implementation.json` を返却契約として使う。
