# Task Plan: 2026-05-17-all-pages-componentization

- `workflow`: work
- `status`: planned
- `lane_owner`: `implement_lane`
- `task_id`: `2026-05-17-all-pages-componentization`
- `task_mode`: frontend-refactor
- `request_summary`: Master Persona POC の結果を前提に、全ページをパネル、カード、モーダル単位へ部品化し、Storybook story も同時に作る。
- `goal`: 全ページのページ component を薄くし、人間見た目レビュー用 Storybook story を主要部品へ付ける。
- `constraints`: `2026-05-17-master-persona-componentization` の人間レビュー後に開始する。全ページは一括 task として扱う。部品化候補は事前調査で列挙し、基準は frontend コーディング規約と UI Component 判断表に従う。
- `close_conditions`: 各ページの主要表示領域が props 境界と Storybook story を持ち、実装後の Storybook review と frontend-local が通過または未通過理由を持つ。
- `worktree_path`: `../AITranslationEngineJP-worktrees/2026-05-17-all-pages-componentization`
- `source_branch`: `codex/2026-05-17-all-pages-componentization`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `pending`
- `screen_design_diff`: `pending-per-screen`
- `ui_agent_browser_review`: `./ui-design.md#storybook-review` または `./storybook-review.md`
- `ux_review`: `pending-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approved_frontend_protection`: `pending-after-frontend-human-review`
- `scenario_design`: `pending`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

全ページの部品化リファクタと Storybook 化は同じ task で行う。全ページは一括 task として扱う。

### 理由

Storybook 基盤は先行 task で用意済みになる。各ページの部品境界は story fixture と同時に検証した方が、props 境界と見た目レビューがずれにくい。

### 影響

各画面の表示領域棚卸し、共有部品と画面専用部品の分類、props 境界、story fixture、既存表示項目の維持確認が必要になる。ページ component は controller 接続と部品合成へ寄せる。実装後の見た目レビューは Storybook で行う。

### 未決事項

- 事前調査で列挙する部品化候補の画面順序。
- 共通化候補を `frontend/src/ui/components/` へ上げる具体基準。
- 既存表示項目を削る場合の人間承認粒度。

## Pre-design Investigation

- 全ページの表示領域を事前に調べ、パネル、カード、モーダル単位の部品化候補を `ui-design.md` または `component-candidates.md` に列挙する。
- 部品化候補の判定基準は `docs/coding-guidelines-frontend.md` の component 分割基準と、`docs/architecture.md` の UI Component 判断表に従う。
- 候補ごとに、共有部品へ上げる候補、画面専用部品に残す候補、分けない候補を分ける。
- 分けない候補は、親画面状態を大量に読む、業務フロー全体の進行状態を持つ、props が条件分岐の塊になる、のいずれかの理由を残す。
- story 対象はパネル、カード、モーダルを基本にし、ページ story は密度確認と配置確認に限定する。

## Routing Notes

- `required_reading`: `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, `frontend/src/ui/screens/`, `frontend/src/ui/components/`, `2026-05-18-storybook-foundation` と `2026-05-17-master-persona-componentization` の完了 artifact
- `canonicalization_targets`: `docs/architecture.md` または screen-design 正本に影響が出る場合だけ docs_updater へ渡す
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
- `source_branch`: `codex/2026-05-17-all-pages-componentization`
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
- `follow_up`: `Storybook 運用と人間見た目レビュー正本の docs 正本化`

## Outcome

- 未着手。
