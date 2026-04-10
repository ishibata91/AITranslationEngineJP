# 実装計画

- workflow: impl
- status: completed
- lane_owner: codex
- scope: .codex/skills, .codex/workflow.md, docs/exec-plans/templates
- task_id: 2026-04-10-phase-2-skill-split
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- 第2段階を `phase-2-ui`、`phase-2-scenario`、`phase-2-logic` に分割する
- `phase-2-ui` は HTML モック、`phase-2-scenario` はシナリオテスト一覧、`phase-2-logic` は実装計画を担当する
- HTML モックとシナリオテスト一覧は exec-plan 本文とは別 artifact として置く

## 判断根拠

- 現行の `phase-2-design` は `UI`、`Scenario`、`Logic` を 1 skill に抱え、artifact 保存先も exec-plan 内に混在している
- `phase-4-plan` の責務は user 指定の `Logic -> 実装計画作成` と重複する
- workflow、orchestrator、handoff 契約、template を同時にそろえないと live workflow が二重化する

## 対象範囲

- `.codex/workflow.md`
- `.codex/skills/orchestrating-implementation/`
- `.codex/skills/phase-2-design/`
- `.codex/skills/phase-4-plan/`
- `.codex/skills/phase-2.5-design-review/`
- `.codex/skills/phase-5-test-implementation/`
- `.codex/skills/phase-6.5-ui-check/`
- `.codex/skills/phase-7-unit-test/`
- `.codex/skills/diagramming-structure-diff/`
- `docs/exec-plans/templates/impl-plan.md`

## 対象外

- product code
- `docs/` の恒久仕様変更
- fix lane

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 変更対象は `.codex/` と exec-plan template に限定し、product 実装ファイルには触れない

## 実装計画

- `parallel_task_groups`:
  - `group_id`: phase-2-skill-split
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: 新しい phase-2 skill、workflow、contract、template が同じ artifact 前提でそろう
- `tasks`:
  - `task_id`: split-phase-2-skills
  - `owned_scope`: `.codex/skills/phase-2-*`, `.codex/skills/orchestrating-implementation/`
  - `depends_on`: none
  - `parallel_group`: phase-2-skill-split
  - `required_reading`: `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: sync-plan-template-and-downstream-skills
  - `owned_scope`: `docs/exec-plans/templates/impl-plan.md`, `.codex/skills/phase-2.5-design-review/`, `.codex/skills/phase-5-test-implementation/`, `.codex/skills/phase-6.5-ui-check/`, `.codex/skills/phase-7-unit-test/`, `.codex/skills/diagramming-structure-diff/`
  - `depends_on`: split-phase-2-skills
  - `parallel_group`: phase-2-skill-split
  - `required_reading`: `docs/exec-plans/templates/impl-plan.md`, `.codex/skills/phase-5-test-implementation/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `phase-2-ui`、`phase-2-scenario`、`phase-2-logic` が作成され、責務が分離されている
- `phase-4-plan` が live workflow から外れている
- exec-plan template が UI モック path と Scenario path を別 artifact として参照できる
- 後続 skill が別 artifact の UI / Scenario と exec-plan の実装計画を読める前提へ更新されている

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`

## HITL 状態

- user が `skill-modification` を明示起動済み

## 承認記録

- user request at 2026-04-10

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/skills/`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移動する

## 結果

- `.codex/skills/phase-2-ui`、`.codex/skills/phase-2-scenario`、`.codex/skills/phase-2-logic` を追加し、各 skill の `SKILL.md`、`agents/openai.yaml`、`references/permissions.json`、返却契約を作成した
- `.codex/skills/phase-2-design` と `.codex/skills/phase-4-plan` を削除し、orchestrator の handoff 先を新しい 3 skill へ差し替えた
- `.codex/workflow.md` を更新し、第2段階を 3 skill の並列工程として明記し、HTML モックと Scenario テスト一覧を exec-plan 別 artifact とする保存規約を追加した
- `docs/exec-plans/templates/impl-plan.md` を更新し、`UI モック` と `Scenario テスト一覧` を path 参照 section に変更し、`実装計画` を plan 本文の implementation brief として残す形にそろえた
- `.codex/skills/phase-2.5-design-review`、`.codex/skills/phase-5-test-implementation`、`.codex/skills/phase-6.5-ui-check`、`.codex/skills/phase-7-unit-test`、`.codex/skills/diagramming-structure-diff` の参照前提を、新しい artifact 構成へ同期した
- `.codex/agents/task_designer.toml` と `.codex/agents/workplan_builder.toml` の説明を live workflow に合わせて更新した
- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` が通過した
