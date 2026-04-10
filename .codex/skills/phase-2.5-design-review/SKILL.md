---
name: phase-2.5-design-review
description: task-local design と review 用差分図を単発で照合し、詳細設計 AI review として `pass` か `reroute` を返す。
---

# Phase 2.5 Design Review

## Review Scope

- 要件取りこぼし
- 責務腐敗
- 検証不足
- 構造差分の不整合

## Output

- decision: `pass` or `reroute`
- findings
- human_open_questions
- closeout_notes

## Rules

- review は 1 回だけ行う
- 対象は HTML モック artifact、Scenario テスト一覧 artifact、active exec-plan の `実装計画`、必要時の review 用差分図に限定する
- 要件取りこぼし、責務腐敗、検証不足、構造差分の不整合だけを見る
- 実装改善やコード改善は提案しない
- human 判断が必要な論点は `human_open_questions` として切り出す

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2.5-design-review.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2.5-design-review.to.orchestrating-implementation.json` を返却契約として使う。
