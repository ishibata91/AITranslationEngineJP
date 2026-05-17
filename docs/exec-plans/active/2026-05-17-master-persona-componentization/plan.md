# Task Plan: 2026-05-17-master-persona-componentization

- `workflow`: work
- `status`: planned
- `lane_owner`: `implement_lane`
- `task_id`: `2026-05-17-master-persona-componentization`
- `task_mode`: frontend-refactor
- `request_summary`: Storybook 基盤を前提に、マスターペルソナ生成画面の部品化と Storybook POC を同じ実装設計で行う。
- `goal`: `MasterPersonaPage` を薄くし、画面 view model 由来の小さい props を持つ表示部品と story を作る。
- `constraints`: `2026-05-18-storybook-foundation` の完了後に開始する。fakeAPI 削除 task と並列に進めてよい。見た目レビュー対象はパネル、カード、モーダルを基本にする。
- `close_conditions`: Master Persona の主要表示領域が story 化可能な props 境界を持ち、Storybook 上で主要状態を確認でき、frontend-local と build-storybook が通過または未通過理由を持つ。
- `worktree_path`: `../AITranslationEngineJP-worktrees/2026-05-17-master-persona-componentization`
- `source_branch`: `codex/2026-05-17-master-persona-componentization`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `pending`
- `screen_design_diff`: `./screen-design-diff.master-persona.md`
- `ui_agent_browser_review`: `./ui-design.md#storybook-review` または `./storybook-review.md`
- `ux_review`: `pending-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approved_frontend_protection`: `pending-after-frontend-human-review`
- `scenario_design`: `pending`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

最初の Storybook 前提リファクタ対象はマスターペルソナ生成画面にする。部品化と Storybook POC は同じ task で扱う。

### 理由

`AIModelSelectionCard` と既存の画面専用部品があるため、パネル、カード、モーダル単位の境界を検証しやすい。Storybook 基盤を先に用意すると、実装設計時に props 境界と story fixture を同時に検証できる。

### 影響

`MasterPersonaPage` は controller 接続と部品合成へ寄せる。生成設定、実行状態、ペルソナレビュー、操作モーダル、モデル選択カードは story を持つ。実装後の人間見た目レビューは Storybook の review URL で行う。

### 未決事項

- 現行 view model だけで、未設定、生成中、生成成功、生成失敗、編集中をすべて表現できるか。
- 既存部品を残す範囲と、新規部品として切り出す範囲。
- ページ story 用の最小合成単位を page component にするか、review-only wrapper にするか。
- Storybook review URL、確認状態、未確認状態を `ui-design.md` と別 artifact のどちらに残すか。

## Routing Notes

- `required_reading`: `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `frontend/src/ui/screens/master-persona/`, `2026-05-18-storybook-foundation` の完了 artifact
- `canonicalization_targets`: `N/A`
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`, `npm --prefix frontend run build-storybook`

## Branch Status

- `worktree_checkout`: `pending`
- `branch_ready`: `pending`
- `commit_hash`: `pending`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `required-after-design-bundle`
- `ux_review`: `required-before-frontend-human-review`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approval_record`: `pending-after-design-bundle`

## Codex Implementation Result

- `completed_handoffs`: `pending`
- `touched_files`: `pending`
- `implemented_scope`: `pending`
- `test_results`: `pending`
- `implementation_investigation`: `pending`
- `ui_evidence`: `pending`
- `ux_review_result`: `pending`
- `approved_frontend_protection`: `pending`
- `codex_review_result`: `pending`
- `sonar_gate_result`: `pending`
- `residual_risks`: `pending`
- `docs_changes`: `pending`

## Merge Readiness

- `merge_ready`: `pending`
- `source_branch`: `codex/2026-05-17-master-persona-componentization`
- `target_branch`: `master`
- `commit_hash`: `pending`
- `validation_evidence`: `pending`
- `review_evidence`: `pending`
- `residual_risks`: `pending`

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `pending`
- `post_merge_validation`: `pending`
- `completed_move`: `pending`
- `merge_commit_hash`: `pending`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `pending`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: `3-1 全ページ部品化と Storybook 化`

## Outcome

- 未着手。
