---
name: phase-2-logic
description: 第2段階の Logic 実装計画作成を担当し、承認前の active exec-plan に implementation brief を固定する。
---

# Phase 2 Logic

## Goal

- 実装順と並列実行単位を固定する
- owned scope、task dependency、required reading、validation commands を固定する
- 必要なら review 用差分図と差分正本適用先まで揃える

## Rules

- 更新対象は原則として active exec-plan の `実装計画`、`review 用差分図`、`差分正本適用先` に限定する
- `phase-2-ui` と `phase-2-scenario` の artifact を読み、実装に必要な依存関係と証明対象を崩さない implementation brief を作る
- 各 task は独立したコンテキストで実装できる粒度まで分解する
- 実装コード、HTML モック、Scenario artifact、test file の詳細実装は持たない
- `tasks.md` や別の実装管理 artifact を作らない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2-logic.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2-logic.to.orchestrating-implementation.json` を返却契約として使う。
