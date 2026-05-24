# Task Plan: frontend-overall-refactor

- `workflow`: refactor-lane
- `status`: ready-for-merge-lane
- `lane_owner`: `refactor_lane`
- `task_id`: `frontend-overall-refactor`
- `task_mode`: frontend refactor
- `request_summary`: フロントエンド全体をリファクタしたい。
- `goal`: frontend 全体について、仕様乖離、構造品質、テスト品質を調査し、人間判断と承認済み `implementation-scope` に基づくリファクタへ分解する。
- `constraints`: `refactor-lane` はプロダクトコード、プロダクトテスト、docs 正本本文を直接変更しない。
- `close_conditions`: `refactor-lane` の完了規約に従い、local commit と `マージ準備入力` まで記録する。
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
- `docs/exec-plans/completed/`
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
- `target_branch_head`: `e301d3ef0ccc2a9d0ebf95841c3da771c442eca4`
- `commit_hash`: `pending`
- `remote_operation`: `not-performed`

## HITL Status

- `spec_implementation_priority`: `approved`
- `refactor_scope_confirmation`: `approved-first-unit`
- `implementation_scope`: `ready`
- `implementation_handoff_input`: `ready`
- `approval_record`: 2026-05-24 の user `continue` を、`refactor-scope-confirmation.md` の次の推奨単位 1 の承認として扱う。

## Refactor Lane Result

- `current_artifact`: `マージ準備入力`
- `completed_artifacts`: `task 枠`, `branch 準備`, `仕様乖離整理起動入力`, `仕様乖離整理`, `仕様実装優先判断`, `構造品質調査`, `テスト品質調査`, `リファクタ範囲確認`, `実装範囲`, `実装引き継ぎ入力`, `frontend リファクタ`, `統合境界リファクタ`, `単体テスト`, `最終検証`, `実装後ブラウザ確認`, `レビュー通過根拠`, `docs正本化判断`, `作業 commit`
- `blocked_artifacts`: `N/A`
- `blocked_reason`: `N/A`

## Codex Implementation Result

- `completed_handoffs`: `FE-001-root-wiring-cleanup`, `FE-002-storybook-dead-story-cleanup`, `FE-003-dead-page-component-cleanup`, `UT-001-app-shell-stale-setup-cleanup`, `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup`, `UT-003-dead-store-presenter-test-cleanup`, `INT-001-wails-gateway-cleanup`, `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup`
- `ready_handoffs`: `N/A`
- `touched_files`: `docs/exec-plans/active/frontend-overall-refactor/plan.md`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md`, `docs/exec-plans/active/frontend-overall-refactor/investigate-input.spec-drift.md`, `docs/exec-plans/active/frontend-overall-refactor/spec-drift-investigation.md`, `docs/exec-plans/active/frontend-overall-refactor/structure-quality-investigation.md`, `docs/exec-plans/active/frontend-overall-refactor/test-quality-investigation.md`, `docs/exec-plans/active/frontend-overall-refactor/refactor-scope-confirmation.md`, `docs/exec-plans/active/frontend-overall-refactor/detail-spec-diff.md`, `docs/exec-plans/active/frontend-overall-refactor/implementation-scope.md`, `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, `frontend/src/application/contract/translation-job-setup/**`, `frontend/src/application/gateway-contract/model-settings-card/index.ts`, `frontend/src/application/gateway-contract/model-settings-card/model-settings-card-policy.ts`, `frontend/src/application/gateway-contract/translation-job-setup/**`, `frontend/src/application/presenter/translation-job-setup/**`, `frontend/src/application/store/translation-job-setup/**`, `frontend/src/application/usecase/translation-job-setup/**`, `frontend/src/controller/translation-job-setup/**`, `frontend/src/controller/wails/translation-job-setup.gateway.ts`, `frontend/src/controller/wails/translation-job-setup.gateway.test.ts`, `frontend/src/controller/wails/gateway-dto/translation-job-setup/**`, `frontend/src/ui/screens/__fixtures__/screen-page-controller-fixtures.ts`, `frontend/src/ui/screens/translation-job-setup/**`, `frontend/src/ui/views/AppShell.test.ts`
- `implemented_scope`: `FE-001-root-wiring-cleanup`, `FE-002-storybook-dead-story-cleanup`, `FE-003-dead-page-component-cleanup`, `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup`, `INT-001-wails-gateway-cleanup`, `UT-001-app-shell-stale-setup-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup`, `UT-003-dead-store-presenter-test-cleanup`, `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup`, `INT-002-model-settings-card-unused-export-cleanup`
- `test_results`: `rg -n "translation-job-setup|TranslationJobSetup|createTranslationJobSetup|JobSetupPage|PhaseSettingsPanel|PhaseSettingsSummaryPanel|cloneModelSettingsCardStates" frontend/src` は一致なし; `python3 scripts/harness/run.py --suite structure` passed on 2026-05-24; `npm --prefix frontend run build-storybook` passed on 2026-05-24; `python3 scripts/harness/run.py --suite frontend-local` passed on 2026-05-24; `npm --prefix frontend run test -- src/ui/views/AppShell.test.ts` passed in `UT-001`; `npm --prefix frontend run test -- src/application/store src/application/presenter` passed in `UT-003`; `npm --prefix frontend run test -- src/application/usecase` passed in `UT-004`; `npm --prefix frontend run test -- src/controller/wails` passed in `UT-005`.
- `implementation_investigation`: `structure-and-test-quality-investigation-ready`
- `ui_evidence`: `not-run-dead-code-cleanup-only`
- `codex_review_result`: `no_issue`
- `sonar_gate_result`: `not-run`
- `residual_risks`: `FSD-005` は `実装が正` の docs 正本化候補として残る。code 修正対象ではない。
- `docs_changes`: task-local artifact only. `FSD-005` の docs 正本化は必要だが、docs 正本文言の人間承認が未取得のため正本本文は未変更。

## Browser Confirmation

- `required`: `false`
- `result`: `not-run`
- `reason`: 実装差分は `translation-job-setup` の dead code cleanup と `model-settings-card` の unused export cleanup である。live UI 導線、表示、generated Wails binding、backend 接続は変更していない。確認対象 URL と操作経路を定義できる live 画面差分がないため、実装後ブラウザ確認は起動しない。

## Review Evidence

- `behavior`: `reviewback.behavior.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `contract`: `reviewback.contract.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `trust_boundary`: `reviewback.trust-boundary.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `state_invariant`: `reviewback.state-invariant.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`
- `responsibility_boundary`: `reviewback.responsibility-boundary.yaml` は `review_status: no_issue`, `must_fix_open: false`, `max_level: none`

## Docs Canonicalization Decision

- `required`: `true`
- `target`: `FSD-005`
- `code_change`: `not-required`
- `canonical_docs_change`: `not-performed`
- `stop_reason`: docs 正本文言の人間承認が未取得である。`FSD-005` は後続 `updating-docs` の docs-only 候補として扱う。
- `handoff_input`: `detail-spec-diff.md#docs-正本化判断`

## Merge Readiness

- `merge_ready`: `ready`
- `source_branch`: `codex/frontend-overall-refactor`
- `target_branch`: `master`
- `commit_hash`: `created-local-commit`
- `validation_evidence`: `python3 scripts/harness/run.py --suite structure` passed on 2026-05-24; `npm --prefix frontend run build-storybook` passed on 2026-05-24; `python3 scripts/harness/run.py --suite frontend-local` passed on 2026-05-24
- `review_evidence`: `reviewback.behavior.yaml`, `reviewback.contract.yaml`, `reviewback.trust-boundary.yaml`, `reviewback.state-invariant.yaml`, `reviewback.responsibility-boundary.yaml` are `no_issue`
- `residual_risks`: `FSD-005` の docs 正本化は未実施である。code 修正対象ではない。

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `N/A`
- `post_merge_validation`: `N/A`
- `completed_move`: `N/A`
- `merge_commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: task-local artifacts only
- `detail_spec_canonicalization`: `stopped-docs-wording-approval-needed`
- `follow_up`: `merge-lane` へ渡す。`FSD-005` の docs 正本化は後続 `updating-docs` の docs-only 候補として扱う。

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
- `python3 scripts/harness/run.py --suite structure` は通過した。
- `npm --prefix frontend run build-storybook` は通過した。
- `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 5 観点レビューは全て `no_issue` である。
- `FSD-005` の docs 正本化判断を task-local に記録した。docs 正本本文は変更していない。
- local commit を作成した。commit hash は最終応答で報告する。
