# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/diagramming-structure-diff/`
- task_id: 2026-04-11-diagramming-structure-diff-component-scope
- task_catalog_ref: N/A
- parent_phase: phase-2

## 要求要約

- `diagramming-structure-diff` に、`docs/diagrams/components/backend` または `docs/diagrams/components/frontend` 以外を触らない制約を明記する。

## 判断根拠

<!-- Decision Basis -->

- 現在の skill 記述は `diagrams/backend/` を対象にしており、architecture 図まで触る余地が残っている。
- ユーザー要望は対象ディレクトリを component 図配下へ限定することにある。

## 対象範囲

- `.codex/skills/diagramming-structure-diff/references/permissions.json`
- `.codex/skills/diagramming-structure-diff/SKILL.md`

## 対象外

- `docs/` 正本ディレクトリ構成そのものの変更
- orchestrator 側の handoff 文面変更

## 依存関係・ブロッカー

- `docs/diagrams/components/` が未作成なら、この skill は停止条件として扱う必要がある。

## 並行安全メモ

- skill 契約の明文化だけを扱う。

## UI モック

- `artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_artifact_path`: `docs/mocks/<page-id>/index.html`
- `summary`: N/A

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`: N/A

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: `skill-scope-clarification`
  - `can_run_in_parallel_with`: `none`
  - `blocked_by`: `none`
  - `completion_signal`: `permissions.json` と `SKILL.md` に component 図配下限定が明記される
- `tasks`:
  - `task_id`: `restrict-component-scope`
  - `owned_scope`: `.codex/skills/diagramming-structure-diff/`
  - `depends_on`: `none`
  - `parallel_group`: `skill-scope-clarification`
  - `required_reading`: `permissions.json`, `SKILL.md`, `docs/index.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- 許可された source path が `docs/diagrams/components/backend` または `docs/diagrams/components/frontend` に限定されている。
- それ以外の図ディレクトリへ触らないことが forbidden または rule として明示されている。

## 必要な証跡

<!-- Required Evidence -->

- structure harness の通過

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- repo 側の正本ディレクトリが未整備な場合、skill は path 不整合を停止条件として返す。

## 結果

<!-- Outcome -->

- `diagramming-structure-diff` の allowed / forbidden / stop conditions を更新し、`docs/diagrams/components/backend/` と `docs/diagrams/components/frontend/` 以外の図ディレクトリを対象外にした。
- `SKILL.md` の goal / workflow / rules を更新し、component 図配下以外を読まない、書かない、更新対象に含めないことを明記した。
- `python3 scripts/harness/run.py --suite structure` が通過した。
