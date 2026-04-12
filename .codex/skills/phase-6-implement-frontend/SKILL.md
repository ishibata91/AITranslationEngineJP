---
name: phase-6-implement-frontend
description: 第6段階の実装と品質通過を担当し、frontend の担当範囲を実装して local validation を返す。
---

# Phase 6 Implement Frontend

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- implementation lane では active exec-plan、承認済み UI モック artifact、承認済み Scenario テスト一覧 artifact、work brief、承認済み required reading を読んでから編集する
- fix lane では active fix plan、accepted fix scope、trace / analysis 根拠、work brief を読んでから編集する
- frontend owned scope だけを変更する
- implementation lane では第2段階で固定した設計をこの段階で作り直さない
- fix lane では accepted fix scope を超えて設計や scope を広げない
- orchestrator から渡された local validation は必要時だけ使い、途中での重い harness 実行を前提にしない
- phase closeout 前に `python3 scripts/harness/run.py --suite all` を実行して問題がないか確認する
- phase closeout 前に Sonar MCP で open issue がないことを確認し、review gate 阻害要因を持ち越さない
- plan の書き換えや lane 切り替えはしない
- validation results には closeout 前の `--suite all` と Sonar MCP による open issue 確認の結果を含める

## Reference Use

- implementation lane では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-6-implement-frontend.json` を参照して入力契約を確認する。
- fix lane では着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.phase-6-implement-frontend.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-6-implement-frontend.to.orchestrating-implementation.json` を返却契約として使う。
- `orchestrating-fixes` へ返す時は `references/phase-6-implement-frontend.to.orchestrating-fixes.json` を返却契約として使う。
