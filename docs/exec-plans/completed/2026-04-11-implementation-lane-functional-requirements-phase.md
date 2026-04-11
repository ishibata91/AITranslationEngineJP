# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/`, `.codex/skills/phase-1.5-functional-requirements/`, `.codex/skills/phase-2-ui/`, `.codex/skills/phase-2.5-design-review/`, `docs/exec-plans/templates/impl-plan.md`
- task_id: 2026-04-11-implementation-lane-functional-requirements-phase
- task_catalog_ref: N/A
- parent_phase: workflow-update

## 要求要約

- 実装レーンに `phase-1.5-functional-requirements` を追加し、`phase-2-ui` を phase-2 の前へ移す。
- `機能要件 + UI モック` を対象にする前段 HITL を追加し、後段の詳細設計 HITL は維持する。

## 判断根拠

<!-- Decision Basis -->

- `phase-1-distill` は facts / constraints / gaps の整理専用であり、機能要件固定の専用 phase がない。
- `phase-2-ui` は現状 detailed design 群の一部として扱われているが、ユーザー要望は UI モックを前段へ引き上げることにある。
- 既存 workflow の人間承認は詳細設計後だけであり、前段の要件合意を明示できない。

## 対象範囲

- `.codex/workflow.md`
- `.codex/skills/orchestrating-implementation/`
- `.codex/skills/phase-1.5-functional-requirements/`
- `.codex/skills/phase-2-ui/`
- `.codex/skills/phase-2.5-design-review/`
- `docs/exec-plans/templates/impl-plan.md`

## 対象外

- product code
- fix lane
- docs/ 配下の恒久仕様文書

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- workflow 記述、skill 契約、plan template の整合だけを扱う。

## 機能要件

- `summary`: 前段で確定した機能境界を active exec-plan に残す。
- `in_scope`: 実装対象に含める項目を列挙する。
- `out_of_scope`: 今回の対象外を列挙する。
- `open_questions`: 前段 HITL が解く論点を残す。
- `required_reading`: UI モックへ渡す最小限の参照を残す。

## UI モック

- `artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_artifact_path`: `docs/mocks/<page-id>/index.html`
- `summary`: `phase-2-ui` は actual skill 名を維持しつつ前段へ移す。

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`: detailed design phase で継続して固定する。

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: `workflow-skill-sync`
  - `can_run_in_parallel_with`: `none`
  - `blocked_by`: `none`
  - `completion_signal`: workflow、orchestrator、skill contracts、plan template が同じ phase 順序を表している
- `tasks`:
  - `task_id`: `add-functional-requirements-phase`
  - `owned_scope`: `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/`, `.codex/skills/phase-1.5-functional-requirements/`, `.codex/skills/phase-2-ui/`, `.codex/skills/phase-2.5-design-review/`, `docs/exec-plans/templates/impl-plan.md`
  - `depends_on`: `none`
  - `parallel_group`: `workflow-skill-sync`
  - `required_reading`: `phase-1-distill`, `phase-2-ui`, `phase-2.5-design-review`, `orchestrating-implementation`, `workflow.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `phase-1.5-functional-requirements` が新設されている。
- `phase-2-ui` が phase-2 前の flow に移されている。
- `機能要件 + UI モック` を対象にする前段 HITL と、詳細設計後の後段 HITL が分離されている。
- plan template が前段承認と後段承認を別々に記録できる。

## 必要な証跡

<!-- Required Evidence -->

- structure harness の通過

## 機能要件 HITL 状態

- N/A

## 機能要件 承認記録

- N/A

## 詳細設計 HITL 状態

- N/A

## 詳細設計 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- workflow と skill 名の対応を崩さない。
- actual skill 名 `phase-2-ui` は維持する。

## 結果

<!-- Outcome -->

- `phase-1.5-functional-requirements` を新設し、SKILL / permissions / agents / handoff contract を追加した。
- `workflow.md` と `orchestrating-implementation` を、前段 `機能要件固定 -> UI モック作成 -> HITL` を含む flow に更新した。
- `phase-2-ui`、`phase-2.5-design-review`、関連 handoff contract、plan template、既存 active plan を新しい承認構造へ同期した。
- `python3 scripts/harness/run.py --suite structure` が通過した。
