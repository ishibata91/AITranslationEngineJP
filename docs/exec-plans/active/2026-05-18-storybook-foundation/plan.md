# Task Plan: 2026-05-18-storybook-foundation

- `workflow`: work
- `status`: planned
- `lane_owner`: `implement_lane`
- `task_id`: `2026-05-18-storybook-foundation`
- `task_mode`: frontend-tooling-foundation
- `request_summary`: Master Persona 部品化に先行して、Storybook の最小基盤を用意する。
- `goal`: Storybook scripts、config、build 検証、fixture 配置方針、review URL 記録方針を最小単位で固定する。
- `constraints`: Master Persona の部品化や画面再設計をこの task で行わない。gateway mock、backend DTO mock、実行フロー再現を Storybook に入れない。
- `close_conditions`: Storybook dev / build の入口が動き、空または最小サンプル story で基盤確認ができ、後続 Master Persona task が story を追加できる。
- `worktree_path`: `../AITranslationEngineJP-worktrees/2026-05-18-storybook-foundation`
- `source_branch`: `codex/2026-05-18-storybook-foundation`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `N/A`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `pending`
- `ux_review`: `N/A`
- `frontend_human_review`: `not-required`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `pending`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

Master Persona 部品化の前に Storybook 基盤を用意する。

### 理由

部品化と story 化は実装設計中に同時に判断する可能性が高い。先に Storybook の起動、build、fixture 配置、review URL 記録の基盤だけを作ると、Master Persona task は部品境界と story を同じ実装設計で検証できる。

### 影響

`frontend/package.json`、`.storybook/`、Storybook 関連設定、最小 story、build-storybook 検証が対象になる。Master Persona 固有の部品分割、状態 fixture、見た目レビューは後続 task で扱う。

### 未決事項

- Storybook build を `frontend-local` に含めるか、別 gate にするか。
- fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くか。
- Storybook 運用をどの docs 正本へ反映するか。

## Routing Notes

- `required_reading`: `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `canonicalization_targets`: `docs/tech-selection.md`, `docs/lint-policy.md`, Storybook 運用を正本化する対象 docs
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`, `npm --prefix frontend run build-storybook`

## Branch Status

- `worktree_checkout`: `pending`
- `branch_ready`: `pending`
- `commit_hash`: `pending`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `required-after-design-bundle`
- `ux_review`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `pending-after-design-bundle`

## Codex Implementation Result

- `completed_handoffs`: `pending`
- `touched_files`: `pending`
- `implemented_scope`: `pending`
- `test_results`: `pending`
- `implementation_investigation`: `pending`
- `ui_evidence`: `pending`
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: `pending`
- `sonar_gate_result`: `pending`
- `residual_risks`: `pending`
- `docs_changes`: `pending`

## Merge Readiness

- `merge_ready`: `pending`
- `source_branch`: `codex/2026-05-18-storybook-foundation`
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
- `follow_up`: `2026-05-17-master-persona-componentization`

## Outcome

- 未着手。
