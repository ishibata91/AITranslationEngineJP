# Task Plan: 2026-05-17-all-pages-componentization

- `workflow`: work
- `status`: completed
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
- `ui_design`: `./ui-design.md`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `./ui-design.md#storybook-review` または `./storybook-review.md`
- `ux_review`: `./ux-review.yaml`
- `frontend_human_review`: `approved`
- `approved_frontend_protection`: `approved`
- `scenario_design`: `./scenario-design.md`
- `scenario_candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `scenario_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `design_diff_diagram`: `./design-diff-all-pages-componentization.puml`
- `design_diff_svg`: `./design-diff-all-pages-componentization.svg`, `./design-diff-all-pages-componentization_001.svg`
- `component_split_by_page`: `./component-split-by-page.puml`
- `component_folder_guideline`: `./component-folder-guideline.md`
- `implementation_scope`: `./implementation-scope.md`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

全ページの部品化リファクタと Storybook 化は同じ task で行う。全ページは一括 task として扱う。

### 理由

Storybook 基盤は先行 task で用意済みになる。各ページの部品境界は story fixture と同時に検証した方が、props 境界と見た目レビューがずれにくい。

### 影響

各画面の表示領域棚卸し、共有部品と画面専用部品の分類、props 境界、story fixture、既存表示項目の維持確認が必要になる。ページ component は controller 接続と部品合成へ寄せる。実装後の見た目レビューは Storybook で行う。

### 未決事項

- `StatusPill`、`ConfirmDangerModal`、phase 系共通部品は、props が増えすぎる場合に画面専用へ戻す。

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

- `worktree_checkout`: `not-used-current-repo`
- `branch_ready`: `codex/2026-05-17-all-pages-componentization`
- `commit_hash`: `806dc9cad8ad9ef3421814d70bf6e12b6272c3ac`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `approved`
- `ux_review`: `required-before-frontend-human-review`
- `frontend_human_review`: `approved-for-merge`
- `approval_record`: `2026-05-18 human message: approve`; `2026-05-20 human message: merge-lane として完了させて`

## Codex Implementation Result

- `completed_handoffs`: `APC-FE-01` through `APC-FE-10`, `APC-UT-11`, `APC-ST-12`
- `touched_files`: `pending`
- `implemented_scope`: frontend componentization and Storybook review evidence
- `test_results`: `frontend-local pass`, `build-storybook pass`, `storybook smoke pass`
- `implementation_investigation`: `pending`
- `ui_evidence`: `./storybook-review.md`
- `ux_review_result`: `./ux-review.yaml`
- `approved_frontend_protection`: `approved`
- `codex_review_result`: `not-required-for-this-merge`
- `sonar_gate_result`: `pending`
- `residual_risks`: Storybook は見た目レビュー入口であり、Wails 実画面の exhaustive visual review は別対象。
- `docs_changes`: task-local artifact only

## Merge Readiness

- `merge_ready`: `ready`
- `source_branch`: `codex/2026-05-17-all-pages-componentization`
- `target_branch`: `master`
- `commit_hash`: `806dc9cad8ad9ef3421814d70bf6e12b6272c3ac`
- `validation_evidence`: `python3 scripts/harness/run.py --suite frontend-local` pass; `npm --prefix frontend run build-storybook` pass; `git diff --check` pass; Storybook representative computed style smoke pass.
- `review_evidence`: `./storybook-review.md`; `./ux-review.yaml`; `2026-05-20 human message: merge-lane として完了させて`
- `residual_risks`: Storybook story は backend、Wails runtime、Gateway、DB、secret store を要求しない見た目レビュー入口である。

## Merge Result

- `merge_status`: `merged-to-master`
- `conflict_resolution`: `not-needed`
- `post_merge_validation`: `python3 scripts/harness/run.py --suite frontend-local` pass; `npm --prefix frontend run build-storybook` pass; `git diff --check --cached` pass.
- `completed_move`: `done`
- `merge_commit_hash`: `48829c730b15a3ca30dca33583a3f8c01b059c9f`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `none`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: `Storybook 運用と人間見た目レビュー正本の docs 正本化`

## Outcome

- `branch 準備`、`scenario_candidates`、`シナリオ設計`、`UI設計`、`設計差分図` は完了。
- 人間設計レビューは `2026-05-18 human message: approve` で承認済み。
- `implementation-scope.md` は `handoff-ready`。
- frontend 実装、Storybook 証跡、UX 事前確認、frontend 人間見た目レビューは完了。
- merge-lane は source branch を target branch へ local merge した。
- conflict は発生しなかった。
- merge 後検証は通過した。
- task folder は `docs/exec-plans/completed/2026-05-17-all-pages-componentization/` へ移動した。
