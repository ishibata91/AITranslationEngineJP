# Task Plan: 2026-05-17-master-persona-componentization

- `workflow`: completed
- `status`: completed
- `lane_owner`: `implement_lane`
- `task_id`: `2026-05-17-master-persona-componentization`
- `task_mode`: frontend-refactor
- `request_summary`: Storybook 基盤を前提に、マスターペルソナ生成画面の部品化と Storybook POC を同じ実装設計で行う。
- `goal`: `MasterPersonaPage` を薄くし、画面 view model 由来の小さい props を持つ表示部品と story を作る。
- `constraints`: `2026-05-18-storybook-foundation` の完了後に開始する。fakeAPI 削除 task と並列に進めてよい。見た目レビュー対象はパネル、カード、モーダルを基本にする。
- `close_conditions`: Master Persona の主要表示領域が story 化可能な props 境界を持ち、Storybook 上で主要状態を確認でき、frontend-local と build-storybook が通過または未通過理由を持つ。
- `worktree_path`: `current-worktree`
- `source_branch`: `codex/2026-05-17-master-persona-componentization`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `./ui-design.md`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `representative-story-playwright-confirmed`
- `ux_review`: `not-run`
- `frontend_human_review`: `human-requested-close-after-story-folder-adjustment`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `design_diff`: `./design-diff-master-persona-componentization.puml`
- `implementation_scope`: `./implementation-scope.md`
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

- `worktree_checkout`: `current-worktree`
- `branch_ready`: `confirmed-local-branch`
- `commit_hash`: `cb5270f`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `approved`
- `ux_review`: `required-before-frontend-human-review`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approval_record`: `2026-05-18 human approved scenario-design.md, ui-design.md, and design-diff by instructing: story作成まで進めて`

## Codex Implementation Result

- `completed_handoffs`: `MPC-FE-01-master-persona-panel-stories`
- `touched_files`: `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte`, `frontend/src/ui/screens/master-persona/RunStatusPanel.svelte`, `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte`, `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte`, `frontend/src/ui/screens/master-persona/master-persona-panel-props.ts`, `frontend/src/ui/screens/master-persona/__fixtures__/master-persona-panel-fixtures.ts`, `frontend/src/ui/screens/master-persona/*.stories.ts`
- `implemented_scope`: `MasterPersonaPage` を controller 接続と panel props 合成へ寄せ、4 つの screen local component の story と synthetic fixture を追加した。
- `test_results`: `npm --prefix frontend run build-storybook` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。`git diff --check` pass。`format:check` は既存広域未整形 63 件で失敗したため、今回変更ファイルだけ `npx prettier --check` pass。
- `implementation_investigation`: Storybook dev は `6006` 使用中と環境 warning のため `6007` で代表 4 story を Playwright 確認した。
- `ui_evidence`: representative Storybook story confirmation via Playwright on port `6007`; no task-local `storybook-review.md` was created before close.
- `ux_review_result`: not run. Human requested commit, merge, completed move, and close after story folder adjustment.
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: not run for 5-viewpoint review. Human requested merge-lane close.
- `sonar_gate_result`: not run.
- `residual_risks`: Full Storybook visual review across all variants is not recorded. UX review and 5-viewpoint Codex review are not recorded.
- `docs_changes`: task-local design artifacts only; docs canonicalization `N/A`

## Merge Readiness

- `merge_ready`: `yes`
- `source_branch`: `codex/2026-05-17-master-persona-componentization`
- `target_branch`: `master`
- `commit_hash`: `f1fdee3`
- `validation_evidence`: `npm --prefix frontend run build-storybook` pass; `python3 scripts/harness/run.py --suite frontend-local` pass; `git diff --check` pass
- `review_evidence`: human requested close after design approval and story folder adjustment
- `residual_risks`: UX review, full visual review across all variants, and 5-viewpoint Codex review are not recorded

## Merge Result

- `merge_status`: `local-merged`
- `conflict_resolution`: `none`
- `post_merge_validation`: `npm --prefix frontend run build-storybook` pass; `python3 scripts/harness/run.py --suite frontend-local` pass; `git diff --check HEAD~1..HEAD` pass
- `completed_move`: `docs/exec-plans/active/2026-05-17-master-persona-componentization` moved to `docs/exec-plans/completed/2026-05-17-master-persona-componentization`
- `merge_commit_hash`: `8f9b06b`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `none`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: `3-1 全ページ部品化と Storybook 化`

## Outcome

- story 作成、story folder 配置、local merge、merge 後検証、completed 移動まで完了。
