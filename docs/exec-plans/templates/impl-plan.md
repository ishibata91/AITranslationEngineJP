# 実装計画テンプレート

- workflow: impl
- status: planned
- lane_owner:
- scope:
- task_id:
- task_catalog_ref:
- parent_phase:

## 要求要約

-

## 判断根拠

<!-- Decision Basis -->

-

## 対象範囲

- Prefer repo-root path prefixes or stable scope tokens that match `tasks/phase-*/tasks/*.yaml` when available.

## 対象外

-

## 依存関係・ブロッカー

-

## 並行安全メモ

- Note the shared files, shared fixtures, or upstream `contract` / `verification` tasks that must land first.

## UI モック

- `artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_artifact_path`: `docs/mocks/<page-id>/index.html`
- `summary`:

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`:

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`:
  - `can_run_in_parallel_with`:
  - `blocked_by`:
  - `completion_signal`:
- `tasks`:
  - `task_id`:
  - `owned_scope`:
  - `depends_on`:
  - `parallel_group`:
  - `required_reading`:
  - `validation_commands`:

## 受け入れ確認

-

## 必要な証跡

<!-- Required Evidence -->

-

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第2段階で作った UI モック working copy は、完了前に `docs/mocks/<page-id>/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/<topic-id>.md` へ移す。

## 結果

<!-- Outcome -->

-
