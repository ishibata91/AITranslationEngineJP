# Task Plan: 2026-05-18-storybook-foundation

- `workflow`: work
- `status`: implementation-reviewed
- `lane_owner`: `implement_lane`
- `task_id`: `2026-05-18-storybook-foundation`
- `task_mode`: frontend-tooling-foundation
- `request_summary`: Master Persona 部品化に先行して、Storybook の最小基盤を用意する。
- `goal`: Storybook scripts、config、build 検証、fixture 配置方針、review URL 記録方針、既存 lint で Storybook への依存混入がないことを確認する検査を最小単位で固定する。
- `constraints`: Master Persona の部品化や画面再設計をこの task で行わない。gateway mock、backend DTO mock、実行フロー再現を Storybook に入れない。
- `close_conditions`: Storybook dev / build の入口が動き、空または最小サンプル story で基盤確認ができ、既存 lint が Storybook への依存混入を検出でき、後続 Master Persona task が story を追加できる。
- `worktree_path`: `/Users/iorishibata/.codex/worktrees/a6a4/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-18-storybook-foundation`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `N/A`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `ux_review`: `N/A`
- `frontend_human_review`: `not-required`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `human-reviewed`
- `implementation_scope`: `handoff-ready`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

Master Persona 部品化の前に Storybook 基盤を用意する。

### 理由

部品化と story 化は実装設計中に同時に判断する可能性が高い。先に Storybook の起動、build、fixture 配置、review URL 記録の基盤だけを作ると、Master Persona task は部品境界と story を同じ実装設計で検証できる。

### 影響

`frontend/package.json`、`.storybook/`、Storybook 関連設定、最小 story、build-storybook 検証、既存 lint の依存境界チェックが対象になる。Master Persona 固有の部品分割、状態 fixture、見た目レビューは後続 task で扱う。

### 回答済み事項

- Storybook build は専用 gate に分ける。
- 既存 lint には、プロダクトコードから Storybook への依存混入がないことを確認する検査を含める。
- fixture は component 横の `frontend/src/ui/**/__fixtures__` に置く。
- Storybook review URL は task-local の `storybook-review.md` に記録する。
- Storybook 運用は、後続 task の plan と POC task が成功した後、skill、agent、docs へ反映する。

## Routing Notes

- `required_reading`: `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `canonicalization_targets`: `docs/tech-selection.md`, `docs/lint-policy.md`, Storybook 運用を正本化する対象 docs
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`, `npm --prefix frontend run lint`, `npm --prefix frontend run build-storybook`

## Branch Status

- `worktree_checkout`: `ready`
- `branch_ready`: `ready`
- `commit_hash`: `see-current-branch-head`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `answered`
- `ux_review`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `scenario-design.md and design-diff-storybook-foundation.puml approved`

## Codex Implementation Result

- `completed_handoffs`: `SBF-FE-01-storybook-foundation`, `SBF-UT-02-lint-storybook-boundary`, `SBF-ST-03-storybook-review-evidence`
- `touched_files`: `frontend/package.json`, `frontend/package-lock.json`, `frontend/.storybook/*`, `frontend/src/ui/components/AIModelSelectionCard.stories.ts`, `frontend/src/ui/components/__fixtures__/ai-model-selection-card-fixture.ts`, `frontend/eslint.config.js`, `scripts/eslint/repository-boundary-plugin.mjs`, `frontend/repository-boundary-plugin.test.mjs`
- `implemented_scope`: Storybook 最小基盤、既存 lint の Storybook 依存混入チェック、Storybook review 記録
- `test_results`: `npm --prefix frontend run lint` pass; `npm --prefix frontend run build-storybook` pass; `python3 scripts/harness/run.py --suite frontend-local` pass; `python3 scripts/harness/run.py --suite coverage` pass
- `implementation_investigation`: `N/A`
- `ui_evidence`: `./storybook-review.md`, `./browser-confirmation.md`, `./browser-confirmation/ai-model-selection-card.png`, `./browser-confirmation/ai-model-selection-card.json`
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: 5 観点 review `no_issue`
- `sonar_gate_result`: repo-local coverage gate pass 71.1%; Sonar server gate not-run
- `residual_risks`: `agent-browser` native snapshot は未取得。headless Playwright 証跡で代替済み。
- `docs_changes`: docs 正本変更なし。task-local 成果物と work_history のみ更新。

## Merge Readiness

- `merge_ready`: `ready-after-commit`
- `source_branch`: `codex/2026-05-18-storybook-foundation`
- `target_branch`: `master`
- `commit_hash`: `see-current-branch-head`
- `validation_evidence`: `npm --prefix frontend run lint`; `npm --prefix frontend run build-storybook`; `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite coverage`; `git diff --check`
- `review_evidence`: `reviewback.behavior.yaml`, `reviewback.contract.yaml`, `reviewback.trust-boundary.yaml`, `reviewback.state-invariant.yaml`, `reviewback.responsibility-boundary.yaml`
- `residual_risks`: `agent-browser` native snapshot 未取得。Storybook global settings warning と Vite chunk-size warning は build 成功に影響なし。

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `pending`
- `post_merge_validation`: `pending`
- `completed_move`: `pending`
- `merge_commit_hash`: `pending`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `none`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: `2026-05-17-master-persona-componentization`

## Outcome

- branch 準備と 6 観点の scenario candidate 作成は完了。
- `scenario-design` は人間回答反映済み。
- `scenario-design` の gate は `finding_count: 0`、`question_count: 0` で通過。
- `implementation-scope.md` は作成済みである。
- implementation handoff 3 件は完了済みである。
- 観測ログ追加は不要と判断済みである。
- 最終検証と実装後ブラウザ確認は完了済みである。
- 5 観点 review はすべて `no_issue` である。
- 次は local commit を作成し、merge lane へ渡す。
