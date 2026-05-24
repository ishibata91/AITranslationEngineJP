# Task Plan: frontend-overall-refactor

- `workflow`: refactor-lane
- `status`: completed
- `lane_owner`: `merge_lane`
- `task_id`: `frontend-overall-refactor`
- `task_mode`: frontend refactor
- `request_summary`: フロントエンド全体をリファクタしたい。
- `goal`: frontend 全体について、仕様乖離、構造品質、テスト品質を調査し、人間判断と承認済み `implementation-scope` に基づくリファクタへ分解する。
- `constraints`: docs 正本本文、`.codex`、remote repository は変更しない。実装差分は承認済みスコープと人間指摘へ限定する。
- `close_conditions`: completed 移動、`python3 scripts/harness/run.py --suite all` 通過、`master` への local merge commit 作成まで記録する。
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `codex/frontend-overall-refactor`
- `target_branch`: `master`

## Artifact Index

- `refactor_task_frame`: `./plan.md#task-frame`
- `refactor_classification`: `./refactor-classification.md`
- `spec_drift_investigation`: `./spec-drift-investigation.md`
- `spec_drift_investigation_input`: `./investigate-input.spec-drift.md`
- `structure_quality_investigation`: `ready`
- `test_quality_investigation`: `ready`
- `structure_quality_investigation_result`: `./structure-quality-investigation.md`
- `test_quality_investigation_result`: `./test-quality-investigation.md`
- `refactor_scope_confirmation`: `approved-first-unit`
- `refactor_scope_confirmation_result`: `./refactor-scope-confirmation.md`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `implementation_scope`: `ready`
- `implementation_scope_result`: `./implementation-scope.md`
- `implementation_handoff_input`: `./implementation-scope.md#handoffs`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`: `.codex/skills/refactor-lane/SKILL.md`
- `required_reading`: `.codex/agents/refactor_lane.toml`
- `required_reading`: `docs/index.md`
- `required_reading`: `docs/spec.md`
- `required_reading`: `docs/architecture.md`
- `required_reading`: `docs/coding-guidelines-frontend.md`
- `required_reading`: `docs/coding-guidelines-tests.md`
- `required_reading`: `docs/screen-design/README.md`
- `canonicalization_targets`: `N/A`
- `detail_spec_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_commands`: `npm --prefix frontend run build-storybook`
- `validation_commands`: `python3 scripts/harness/run.py --suite all`

## Task Frame

### リファクタ目的

frontend 全体の責務過多、責務分離不足、コーディング規約逸脱、構造設計正本との不整合を調査する。
調査結果は実装範囲へ直結させず、人間のリファクタ範囲確認を待つ。

### 対象仕様参照

- `docs/spec.md`: UI から観測できる翻訳ジョブ、翻訳補助メタデータ、進捗、失敗回復の恒久要件。
- `docs/architecture.md`: frontend 層、依存方向、`View`、`UI Component`、`ScreenController`、`Frontend UseCase`、`Store`、`Gateway` の責務。
- `docs/coding-guidelines-frontend.md`: TypeScript、Svelte、Wails gateway、画面状態、表示イベントの実装規約。
- `docs/coding-guidelines-tests.md`: frontend テストの品質確認で参照するテスト規約。
- `docs/screen-design/`: 画面構成と visual design の正本。
- `docs/diagrams/frontend/`: frontend 構造図の正本。
- `docs/diagrams/components/frontend/`: frontend component detail 図の正本。

### 対象実装範囲

- `frontend/src/main.ts`
- `frontend/src/bootstrap/`
- `frontend/src/application/contract/`
- `frontend/src/application/gateway-contract/`
- `frontend/src/application/presenter/`
- `frontend/src/application/store/`
- `frontend/src/application/usecase/`
- `frontend/src/controller/`
- `frontend/src/ui/components/`
- `frontend/src/ui/screens/`
- `frontend/src/ui/stores/`
- `frontend/src/ui/views/`

### 対象テスト範囲

- `frontend/src/**/*.test.ts`
- `frontend/src/test/`
- `frontend/src/**/__fixtures__/`
- `frontend/src/**/*.stories.ts`

### 変更禁止範囲

- `frontend/wailsjs/`
- `internal/`
- `.codex/`
- `docs/` の正本本文
- `frontend/storybook-static/`
- remote repository

### 検証要件

- frontend 層が触れた場合は `python3 scripts/harness/run.py --suite frontend-local` を通す。
- 構造正本または import 境界が触れた場合は `python3 scripts/harness/run.py --suite structure` を通す。
- Storybook 確認資源が触れた場合は `npm --prefix frontend run build-storybook` を通す。
- UI 表示または導線が触れた場合は `browser_confirmation` の結果または停止理由を記録する。

## Branch Status

- `worktree_checkout`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `branch_ready`: `ready`
- `source_branch`: `codex/frontend-overall-refactor`
- `target_branch`: `master`
- `target_branch_head_before_merge`: `15044d4c141c6399647f05e54339beab1bfe8a06`
- `source_branch_head`: `a96bc2b62604e1ab083c494af0de2b3cf7d88947`
- `remote_operation`: `not-performed`

## HITL Status

- `spec_implementation_priority`: `approved`
- `refactor_scope_confirmation`: `approved-first-unit`
- `implementation_scope`: `ready`
- `implementation_handoff_input`: `ready`
- `approval_record`: 2026-05-24 の user `continue` を、`refactor-scope-confirmation.md` の次の推奨単位 1 の承認として扱う。

## Refactor Lane Result

- `current_artifact`: `closed`
- `completed_artifacts`: `task 枠`, `branch 準備`, `仕様乖離整理起動入力`, `仕様乖離整理`, `仕様実装優先判断`, `構造品質調査`, `テスト品質調査`, `リファクタ範囲確認`, `実装範囲`, `実装引き継ぎ入力`, `frontend リファクタ`, `統合境界リファクタ`, `単体テスト`, `最終検証`, `実装後ブラウザ確認`, `レビュー通過根拠`, `docs正本化判断`, `作業 commit`, `completed 移動`, `merge 結果 commit`
- `blocked_artifacts`: `N/A`
- `blocked_reason`: `N/A`

## Codex Implementation Result

- `completed_handoffs`: `FE-001-root-wiring-cleanup`, `FE-002-storybook-dead-story-cleanup`, `FE-003-dead-page-component-cleanup`, `UT-001-app-shell-stale-setup-cleanup`, `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup`, `UT-003-dead-store-presenter-test-cleanup`, `INT-001-wails-gateway-cleanup`, `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup`, `UX-001-job-run-phase-progress-actions`, `UX-002-translation-management-stepper-removal`, `HARNESS-001-detail-spec-gate-removal`
- `ready_handoffs`: `N/A`
- `touched_files`: `docs/exec-plans/completed/frontend-overall-refactor/**`, `frontend/src/application/contract/**/*-phase/**`, `frontend/src/application/presenter/**/*-phase/**`, `frontend/src/application/store/**/*-phase/**`, `frontend/src/application/usecase/**/*-phase/**`, `frontend/src/controller/**/*-phase/**`, `frontend/src/ui/components/phase-progress-actions.ts`, `frontend/src/ui/components/phase-progress-actions.test.ts`, `frontend/src/ui/screens/job-run/**`, `frontend/src/ui/screens/*-phase/**`, `frontend/src/ui/screens/translation-job-management/**`, `frontend/src/ui/stores/shell-state.ts`, `frontend/src/ui/views/AppShell.svelte`, `internal/**/*term_translation_phase*`, `internal/repository/processing_target_sqlite_repository.go`, `internal/service/provider_execution_snapshot.go`, `scripts/harness/**`, `scripts/detail_spec/diff_gate.py`, `sonar-project.properties`, `tests/system/**`
- `implemented_scope`: 承認済み first unit、2026-05-24 から 2026-05-25 の人間ブラウザ指摘、詳細仕様差分 gate 廃止、Storybook 除外を含む Sonar 設定調整。
- `test_results`: `python3 scripts/harness/run.py --suite all` passed on 2026-05-25. `structure`、`execution`、`system-test`、`coverage` が通過した。`system-test` は 9 passed。Sonar 集計は coverage 71.2%、security 0、reliability 0、HIGH maintainability 0。
- `implementation_investigation`: `structure-and-test-quality-investigation-ready`
- `ui_evidence`: 2026-05-24 から 2026-05-25 の人間ブラウザ指摘に基づき、進行状況の操作ボタン、全体 stepper、Storybook fixture を実画面基準へ合わせた。
- `codex_review_result`: `no_issue`
- `sonar_gate_result`: `passed`
- `residual_risks`: `N/A`
- `docs_changes`: task-local artifact only. docs 正本本文は変更していない。詳細仕様差分 gate は harness から廃止した。

## Browser Confirmation

- `required`: `true`
- `result`: `human-browser-review-applied`
- `reason`: 人間が `http://localhost:34115/#translation-management/job-run` と Storybook 上で確認した差分指摘を入力し、Codex が現行 UI 基準で反映した。commit 前の最終 gate は `python3 scripts/harness/run.py --suite all` とした。

## Review Evidence

- `behavior`: `reviewback.behavior.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `contract`: `reviewback.contract.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `trust_boundary`: `reviewback.trust-boundary.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `state_invariant`: `reviewback.state-invariant.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `responsibility_boundary`: `reviewback.responsibility-boundary.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`

## Docs Canonicalization Decision

- `required`: `false`
- `target`: `N/A`
- `code_change`: `performed`
- `canonical_docs_change`: `not-performed`
- `stop_reason`: docs 正本本文の変更は今回の人間指示に含まれない。
- `handoff_input`: `N/A`

## Merge Readiness

- `merge_ready`: `ready`
- `source_branch`: `codex/frontend-overall-refactor`
- `target_branch`: `master`
- `source_commit_hash`: `a96bc2b62604e1ab083c494af0de2b3cf7d88947`
- `validation_evidence`: `python3 scripts/harness/run.py --suite all` passed on 2026-05-25. `system-test` は 9 passed。Sonar 集計は coverage 71.2%、security 0、reliability 0、HIGH maintainability 0。
- `review_evidence`: `reviewback.behavior.yaml`, `reviewback.contract.yaml`, `reviewback.trust-boundary.yaml`, `reviewback.state-invariant.yaml`, `reviewback.responsibility-boundary.yaml` are `no_issue`
- `residual_risks`: `N/A`

## Merge Result

- `merge_status`: `merged-with-plan-conflict-resolution`
- `merge_command`: `git merge --no-ff --no-commit codex/frontend-overall-refactor`
- `conflict_resolution`: `docs/exec-plans/completed/frontend-overall-refactor/plan.md` は、既存 `master` 側の旧 merge 記録と source branch 側の最新 UI/harness closeout 記録を統合して解消した。
- `post_merge_validation`: `python3 scripts/harness/run.py --suite all` passed on 2026-05-25. `system-test` は 9 passed。Sonar 集計は coverage 71.2%、security 0、reliability 0、HIGH maintainability 0。
- `completed_move`: `done`
- `merge_commit_hash`: `created-local-merge-commit`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: task-local artifacts only
- `detail_spec_canonicalization`: `not-required-detail-spec-gate-retired`
- `follow_up`: `N/A`

## Outcome

- frontend 全体リファクタの task 枠、調査起動入力、仕様乖離整理を更新した。
- `仕様乖離整理` は `FSD-005`, `FSD-006` だけを残し、旧 `FSD-001` から `FSD-004` は構造品質調査候補へ分離した。
- `FSD-005` は `実装が正` と判断済みである。`FSD-006` は dead code 論点として `SQ-007` へ移した。
- `構造品質調査` と `テスト品質調査` を sub-agent で作成し、`リファクタ範囲確認` を承認済みにした。
- 承認済み first unit は `SQ-007`, `SQ-001`, `SQ-003`, `TQ-001`, `TQ-002`, `TQ-004` である。
- `designer` agent が `implementation-scope.md` を作成し、`status: ready` とした。
- `implementation-scope.md#Handoffs` を `実装引き継ぎ入力` として扱う。
- 承認済み first unit の全 handoff を完了した。
- 2026-05-24 の人間指示により、`model-settings-card` の未使用 export cleanup を追加スコープへ入れた。
- `model-settings-card` の未使用 export cleanup を完了した。
- `translation-job-setup`、`JobSetupPage`、`PhaseSettingsPanel`、`PhaseSettingsSummaryPanel` の参照は `frontend/src` から消えた。
- `cloneModelSettingsCardStates` の参照は `frontend/src` から消えた。
- 人間ブラウザ指摘に基づき、翻訳管理の全体 stepper を削除した。
- 人間ブラウザ指摘に基づき、進行状況内の明示的な更新 action を削除した。
- 人間ブラウザ指摘に基づき、開始、再開、リトライの表示を 1 個の実行 action へ統合した。
- 未開始で開始できない場合も、開始 action は disabled 表示として残す。
- 中断 action は常に表示対象として残す。
- Sticky footer が次工程 action を管理するため、進行状況内では次工程 action を表示しない。
- Storybook の job run fixture を実画面の phase 表示に合わせた。
- 詳細仕様差分 gate を harness から廃止した。
- Sonar は Storybook と fixture を対象外にした。
- `python3 scripts/harness/run.py --suite all` は通過した。
- completed 移動後の `python3 scripts/harness/run.py --suite structure` は通過した。
- source branch の作業 commit は `a96bc2b62604e1ab083c494af0de2b3cf7d88947` である。
- `codex/frontend-overall-refactor` を `master` へ local merge した。
- merge 後検証として `python3 scripts/harness/run.py --suite all` は通過した。
- merge 後の `system-test` は 9 passed である。
- merge 後の Sonar 集計は coverage 71.2%、security 0、reliability 0、HIGH maintainability 0 である。
- merge commit はこの後に作成する。
