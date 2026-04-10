# 実装計画

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/`, `docs/index.md`, `docs/exec-plans/templates/impl-plan.md`, `docs/mocks/`, `docs/scenario-tests/`
- task_id: workflow-artifact-destination-sync
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- phase2 artifact の最終正本置き場を `docs/mocks/` と `docs/scenario-tests/` に切り替える
- orchestrator が完了時に最終置き場へ移動する流れを workflow に追加する
- `phase-2.5-design-review` の対象に、固定 HTML だけでなく主要なページの動きの再現を含める

## 判断根拠

- task-local working copy と最終正本置き場を分けた方が、phase2 の設計中 artifact と完了後の docs 正本が混ざらない
- `phase-2.5-design-review` が静的 HTML だけを見ると、導線や状態変化の設計漏れを拾いにくい

## 対象範囲

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/skills/orchestrating-implementation/`
- `.codex/skills/phase-2-ui/`
- `.codex/skills/phase-2-scenario/`
- `.codex/skills/phase-2.5-design-review/`
- `docs/index.md`
- `docs/exec-plans/templates/impl-plan.md`
- `docs/mocks/README.md`
- `docs/scenario-tests/README.md`

## 対象外

- product code
- 実際の page mock 実装
- 実際の scenario test 文書の新規追加

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- `docs/index.md` は docs 導線に影響する
- `.codex/workflow.md` と orchestrator skill の表現を揃える必要がある

## UI モック

- `artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_artifact_path`: `docs/mocks/<page-id>/index.html`
- `summary`:

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `summary`:

## 実装計画

- `parallel_task_groups`:
  - `group_id`: workflow-sync
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: workflow、skill、docs 導線、template が新しい正本置き場へ同期される
- `tasks`:
  - `task_id`: sync-workflow
  - `owned_scope`: `.codex/README.md`, `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/`
  - `depends_on`: none
  - `parallel_group`: workflow-sync
  - `required_reading`: `.codex/README.md`, `.codex/workflow.md`
  - `validation_commands`: `git diff --check`
  - `task_id`: sync-phase-skills
  - `owned_scope`: `.codex/skills/phase-2-ui/`, `.codex/skills/phase-2-scenario/`, `.codex/skills/phase-2.5-design-review/`
  - `depends_on`: sync-workflow
  - `parallel_group`: workflow-sync
  - `required_reading`: `.codex/skills/orchestrating-implementation/SKILL.md`
  - `validation_commands`: `git diff --check`
  - `task_id`: sync-doc-entry
  - `owned_scope`: `docs/index.md`, `docs/exec-plans/templates/impl-plan.md`, `docs/mocks/README.md`, `docs/scenario-tests/README.md`
  - `depends_on`: sync-workflow
  - `parallel_group`: workflow-sync
  - `required_reading`: `docs/index.md`, `docs/exec-plans/templates/impl-plan.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- workflow で active working copy と最終正本置き場が分離されている
- orchestrator の close 手順に `docs/mocks/` と `docs/scenario-tests/` への移動が入っている
- `phase-2.5-design-review` がページの動き再現も review 対象として扱う

## 必要な証跡

- `git diff --check`
- `python3 scripts/harness/run.py --suite structure`

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/`
- `docs/mocks/`
- `docs/scenario-tests/`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移す

## 結果

- `.codex/workflow.md` と `.codex/README.md` に、phase2 artifact の working copy と最終正本置き場の分離を反映した
- `orchestrating-implementation` に、完了前に `docs/mocks/<page-id>/index.html` と `docs/scenario-tests/<topic-id>.md` へ移す close 手順を追加した
- `phase-2-ui` に、固定 HTML だけでなく主要導線と状態変化をある程度再現する page mock working copy を作る責務を追加した
- `phase-2.5-design-review` に、主要導線と状態変化の再現不足を review 対象として追加した
- `docs/mocks/README.md` と `docs/scenario-tests/README.md` を追加し、`docs/index.md` の導線を更新した
- `git diff --check` と `python3 scripts/harness/run.py --suite structure` が通過した
