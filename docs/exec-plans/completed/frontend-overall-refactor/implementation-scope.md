# Implementation Scope: frontend-overall-refactor

- `skill`: implementation-scope
- `status`: ready
- `source_plan`: `./plan.md`
- `human_review_status`: `approved-first-unit`
- `approval_record`: `./refactor-scope-confirmation.md` の `status: approved-first-unit`
- `codex_entry`: `.codex/skills/refactor-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`
- `frontend_guideline_reference`: `docs/coding-guidelines-frontend.md`
- `test_guideline_reference`: `docs/coding-guidelines-tests.md`

## Source Artifacts

- `refactor_scope_confirmation`: `./refactor-scope-confirmation.md`
- `structure_quality_investigation`: `./structure-quality-investigation.md`
- `test_quality_investigation`: `./test-quality-investigation.md`
- `spec_drift_investigation`: `./spec-drift-investigation.md`
- `refactor_classification`: `./refactor-classification.md`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`: `N/A`
- `screen_design_reference`: `docs/screen-design/README.md`

## Fixed Decisions

- 人間承認済み first unit は `SQ-007`, `SQ-001`, `SQ-003`, `TQ-001`, `TQ-002`, `TQ-004` である。
- 2026-05-24 の人間指示により、`model-settings-card` の未使用 export cleanup を first unit の追加スコープに含める。
- `JobSetupPage.svelte` と関連 wiring / story / test は dead code cleanup として扱う。
- `InputReviewPage.svelte` から job 作成後に `job-run` へ進む導線は変更しない。
- `AppShell.svelte` と `shell-state.ts` の現行 step 並びは変更しない。
- `TranslationJobManagementPage.svelte`, `JobRunPage.svelte`, `TranslationOutputArtifactPage.svelte` のユーザー向け振る舞いは変更しない。
- `FSD-005` は `実装が正` であるため、code 修正対象に入れない。
- `FSD-005` は docs 正本化候補として後続の `docs正本化判断` へ分離する。
- `frontend/wailsjs/`、`internal/`、docs 正本本文は変更禁止とする。
- backend と frontend は同一引き継ぎにしない。
- docs 正本化は実装引き継ぎへ混ぜない。

## 承認済み実装範囲

| ID | 扱い | 実装範囲 | 根拠 |
| --- | --- | --- | --- |
| `SQ-007` | 実装対象 | `translation-job-setup` の dead code cleanup。対象は page、screen local component、screen controller、usecase、presenter、store、gateway、DTO、story、test、wiring。 | `structure-quality-investigation.md` の `SQ-007` |
| `SQ-001` | 実装対象 | dead code 化した `translation-job-setup` usecase / presenter を削除し、残す live 導線を変えない。 | `structure-quality-investigation.md` の `SQ-001` |
| `SQ-003` | 実装対象 | root View 側に残る `translation-job-setup` fallback wiring と production wiring 参照を削除する。 | `structure-quality-investigation.md` の `SQ-003` |
| `TQ-001` | 単体テスト対象 | `AppShell.test.ts` の未使用 `createTranslationJobSetupScreenController` setup を削除する。 | `test-quality-investigation.md` の `TQ-001` |
| `TQ-002` | 単体テスト対象 | `JobSetupPage.test.ts` を dead code page test として削除する。 | `test-quality-investigation.md` の `TQ-002` |
| `TQ-004` | 単体テスト対象 | live shell 導線 test と dead code page test の証明対象を分け、live shell 側の期待値を維持する。 | `test-quality-investigation.md` の `TQ-004` |
| `FV-001` | 実装対象 | `frontend-local` の `lint:exports` で検出された `model-settings-card` の未使用 export を削除する。 | 2026-05-24 の最終検証失敗と人間承認 |

## 除外範囲と理由

| ID / 範囲 | 除外理由 | 後続扱い |
| --- | --- | --- |
| `SQ-002` | 人間判断が未承認である。shared component 分割は今回の first unit ではない。 | 別 slice 候補 |
| `TQ-003` | 人間判断が未承認である。provider settings page test の mock 境界整理は今回の first unit ではない。 | 別 slice 候補 |
| `SQ-004` | directory 正本と実装配置の優先判断が未確定である。 | docs 正本優先判断後に再計画 |
| `SQ-005` | runtime adapter 配置と transport 種別の判断が未確定である。 | 別判断後に再計画 |
| `SQ-006` | `SQ-007` の dead code cleanup 後に対象が変わる可能性がある。 | cleanup 後に再評価 |
| `FSD-005` code 修正 | 人間判断で `実装が正` である。 | docs 正本化候補のみ |
| `InputReviewPage.svelte` の job 作成導線 | 現行実装が正であり、変更不要範囲である。 | 変更禁止 |
| `AppShell.svelte` と `shell-state.ts` の step 並び | 現行の翻訳管理 step 並びを守る必要がある。 | 変更禁止 |
| `TranslationJobManagementPage.svelte`, `JobRunPage.svelte`, `TranslationOutputArtifactPage.svelte` | ユーザー向け振る舞いを変えない条件がある。 | 変更禁止 |
| backend と `internal/` | first unit は frontend refactor であり、backend 変更は承認されていない。 | 変更禁止 |
| `frontend/wailsjs/` | Wails generated output であり hand-edit しない。 | 変更禁止 |
| docs 正本本文 | docs 正本化は `docs正本化判断` と `updating-docs` の担当である。 | 実装引き継ぎに混ぜない |

## 要否判定

| 成果物 | 要否 | 理由 |
| --- | --- | --- |
| frontend 実装 | 必要 | dead code page、root wiring、screen local object、Storybook story を整理するため。 |
| 統合境界実装 | 必要 | `frontend/src/controller/wails/translation-job-setup.gateway.ts`、gateway DTO、gateway contract を削除対象に含むため。 |
| 単体テスト | 必要 | stale setup と dead code page test を整理し、live shell test の期待値を維持するため。 |
| backend 実装 | 不要 | `internal/` と backend binding は今回の承認済み first unit に含まれないため。 |
| シナリオテスト | 不要 | ユーザー向け導線を変更しない refactor であり、UI 人間操作 E2E の新規証明対象がないため。 |
| docs 正本化 | 実装引き継ぎには不要 | `FSD-005` は docs 正本化候補だが、実装系レーンへ混ぜないため。 |

## 依存表

| handoff | depends_on | 依存理由 |
| --- | --- | --- |
| `FE-001-root-wiring-cleanup` | なし | root wiring の dead reference を先に外す。 |
| `FE-002-storybook-dead-story-cleanup` | なし | story 資源だけを削除し、production code へ依存しない。 |
| `FE-003-dead-page-component-cleanup` | `FE-002-storybook-dead-story-cleanup` | story が dead component を import しない状態を先に作る。 |
| `FE-004-dead-controller-store-contract-cleanup` | `FE-001-root-wiring-cleanup`, `FE-003-dead-page-component-cleanup` | root wiring と page component から contract / controller 参照を外した後に削除する。 |
| `FE-005-dead-usecase-cleanup` | `FE-004-dead-controller-store-contract-cleanup` | controller factory から usecase 参照を外した後に削除する。 |
| `FE-006-dead-presenter-cleanup` | `FE-003-dead-page-component-cleanup` | page component から presenter type 参照を外した後に削除する。 |
| `INT-001-wails-gateway-cleanup` | `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup` | frontend 側の gateway contract 参照を外した後に Wails adapter を削除する。 |
| `UT-001-app-shell-stale-setup-cleanup` | `FE-001-root-wiring-cleanup` | production prop forwarding を削除した後に stale setup を削る。 |
| `UT-002-dead-page-test-and-fixture-cleanup` | `FE-003-dead-page-component-cleanup`, `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup` | dead page と fixture の参照元が削除済みである必要がある。 |
| `UT-003-dead-store-presenter-test-cleanup` | `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup` | deleted module の単体 test を削除する。 |
| `UT-004-dead-usecase-test-cleanup` | `FE-005-dead-usecase-cleanup`, `INT-001-wails-gateway-cleanup` | usecase と gateway contract の削除後に単体 test を削除する。 |
| `UT-005-dead-gateway-test-cleanup` | `INT-001-wails-gateway-cleanup` | Wails gateway 削除後に gateway test を削除する。 |
| `INT-002-model-settings-card-unused-export-cleanup` | `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup` | 最終検証で `translation-job-setup` cleanup 後に残った export gate の失敗を解消する。 |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `FE-001-root-wiring-cleanup`, `FE-002-storybook-dead-story-cleanup` | なし | `FE-001-root-wiring-cleanup` <-> `FE-002-storybook-dead-story-cleanup` | なし |
| `wave-2` | `FE-003-dead-page-component-cleanup`, `UT-001-app-shell-stale-setup-cleanup` | `FE-001-root-wiring-cleanup`, `FE-002-storybook-dead-story-cleanup` | `FE-003-dead-page-component-cleanup` <-> `UT-001-app-shell-stale-setup-cleanup` | なし |
| `wave-3` | `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup` | `FE-001-root-wiring-cleanup`, `FE-003-dead-page-component-cleanup` | `FE-004-dead-controller-store-contract-cleanup` <-> `FE-006-dead-presenter-cleanup` | なし |
| `wave-4` | `FE-005-dead-usecase-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup`, `UT-003-dead-store-presenter-test-cleanup` | `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup` | `FE-005-dead-usecase-cleanup` <-> `UT-002-dead-page-test-and-fixture-cleanup`, `FE-005-dead-usecase-cleanup` <-> `UT-003-dead-store-presenter-test-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup` <-> `UT-003-dead-store-presenter-test-cleanup` | なし |
| `wave-5` | `INT-001-wails-gateway-cleanup` | `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup` | なし | `shared_contract_change` |
| `wave-6` | `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup` | `FE-005-dead-usecase-cleanup`, `INT-001-wails-gateway-cleanup` | `UT-004-dead-usecase-test-cleanup` <-> `UT-005-dead-gateway-test-cleanup` | なし |
| `wave-7` | `INT-002-model-settings-card-unused-export-cleanup` | `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup` | なし | なし |

## Handoffs

### `FE-001-root-wiring-cleanup`

- `implementation_target`: root wiring から `translation-job-setup` の dead controller factory 参照を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-1`
- `depends_on`: なし
- `parallelizable_with`: `FE-002-storybook-dead-story-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/main.ts` の `createTranslationJobSetupScreenController` prop forwarding を削除し、completion の root wiring clause を閉じる。最初にする理由は production entry の unused dependency を最上流で切れるためである。
- `owned_scope`: `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, `frontend/src/ui/App.svelte`
- `change_forbidden_scope`: `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/stores/shell-state.ts`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/wailsjs/`, `internal/`, docs 正本本文
- `size_estimate`: 3 files、180 changed lines 以下
- `completion_signal`: `main.ts` と `bootstrap` から `createTranslationJobSetupScreenController` と `createTranslationJobSetupGateway` の production wiring が消えている。`InputReviewPage.svelte` から `job-run` へ進む導線は変わっていない。`AppShell.svelte` と `shell-state.ts` の step 並びは変わっていない。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が import 境界の破綻を直せる。
- `execution_test_classification`: lower-level only
- `notes`: `frontend-local` は unit test cleanup 完了後の全体検証で実行する。

### `FE-002-storybook-dead-story-cleanup`

- `implementation_target`: `translation-job-setup` の dead story を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-1`
- `depends_on`: なし
- `parallelizable_with`: `FE-001-root-wiring-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/ui/screens/translation-job-setup/stories/JobSetupPurposeHeader.stories.ts` を削除し、completion の dead story clause を閉じる。最初にする理由は story import を production component 削除前に減らせるためである。
- `owned_scope`: `frontend/src/ui/screens/translation-job-setup/stories/*.stories.ts`
- `change_forbidden_scope`: Storybook 設定、通常カテゴリの別 story、`frontend/storybook-static/`
- `size_estimate`: 7 files、180 changed lines 以下
- `completion_signal`: `translation-job-setup` directory 配下に Storybook story が残っていない。Storybook の別 screen story は変わっていない。
- `validation_commands`: `npm --prefix frontend run build-storybook`
- `validation_owner`: この handoff が Storybook import 破綻を直せる。
- `execution_test_classification`: lower-level only
- `notes`: Storybook story を削除する scope を含むため `build-storybook` を必須にする。

### `FE-003-dead-page-component-cleanup`

- `implementation_target`: dead page と screen local component を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-2`
- `depends_on`: `FE-002-storybook-dead-story-cleanup`
- `parallelizable_with`: `UT-001-app-shell-stale-setup-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` を削除し、completion の dead page component clause を閉じる。最初にする理由は page 本体の削除が関連 component 削除の起点になるためである。
- `owned_scope`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`, `CompatibilityPrecheckPanel.svelte`, `CreatedJobSummaryPanel.svelte`, `FoundationDataPanel.svelte`, `InputSourcePanel.svelte`, `JobSetupPurposeHeader.svelte`, `PhaseSettingsPanel.svelte`, `PhaseSettingsSummaryPanel.svelte`, `job-setup-panel-props.ts`
- `change_forbidden_scope`: `InputReviewPage.svelte`, `TranslationJobManagementPage.svelte`, `JobRunPage.svelte`, `TranslationOutputArtifactPage.svelte`, `AppShell.svelte`, `shell-state.ts`
- `size_estimate`: 9 files、約 1408 changed lines。注意規模だが 1 画面の dead component 削除に閉じるため 1 handoff とする。
- `completion_signal`: `translation-job-setup` page component と専用 component が消えている。live 画面の導線、文言、状態表示は変わっていない。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が deleted component の import 破綻を直せる。
- `execution_test_classification`: lower-level only
- `notes`: `frontend-local` は `UT-002` 以降で deleted component test を整理してから実行する。

### `FE-004-dead-controller-store-contract-cleanup`

- `implementation_target`: dead page 用の controller、store、screen contract を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-3`
- `depends_on`: `FE-001-root-wiring-cleanup`, `FE-003-dead-page-component-cleanup`
- `parallelizable_with`: `FE-006-dead-presenter-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/controller/translation-job-setup/translation-job-setup-screen-controller-factory.ts` を削除し、completion の controller factory clause を閉じる。最初にする理由は usecase / store / presenter 参照の束ね口を先に消せるためである。
- `owned_scope`: `frontend/src/controller/translation-job-setup/`, `frontend/src/application/store/translation-job-setup/`, `frontend/src/application/contract/translation-job-setup/`
- `change_forbidden_scope`: `frontend/src/controller/wails/`, `frontend/src/application/gateway-contract/translation-job-setup/`, `frontend/wailsjs/`
- `size_estimate`: 7 files、450 changed lines 以下
- `completion_signal`: dead screen controller factory、controller、store、screen contract が production source から消えている。root wiring は `FE-001` の completion を維持している。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が controller / contract import 破綻を直せる。
- `execution_test_classification`: lower-level only

### `FE-005-dead-usecase-cleanup`

- `implementation_target`: dead page 用 usecase を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-4`
- `depends_on`: `FE-004-dead-controller-store-contract-cleanup`
- `parallelizable_with`: `UT-002-dead-page-test-and-fixture-cleanup`, `UT-003-dead-store-presenter-test-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts` を削除し、completion の usecase clause を閉じる。最初にする理由は controller factory の参照削除後に単独で削除できるためである。
- `owned_scope`: `frontend/src/application/usecase/translation-job-setup/`
- `change_forbidden_scope`: live input review usecase、job run usecase、backend usecase
- `size_estimate`: 2 files、約 1311 changed lines。注意規模だが 1 usecase directory の dead code 削除に閉じるため 1 handoff とする。
- `completion_signal`: `translation-job-setup` usecase directory が消えている。live job 作成導線の usecase は変わっていない。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が usecase import 破綻を直せる。
- `execution_test_classification`: lower-level only

### `FE-006-dead-presenter-cleanup`

- `implementation_target`: dead page 用 presenter を削除する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `ready_wave`: `wave-3`
- `depends_on`: `FE-003-dead-page-component-cleanup`
- `parallelizable_with`: `FE-004-dead-controller-store-contract-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts` を削除し、completion の presenter clause を閉じる。最初にする理由は page component の type import を消した後に単独で削除できるためである。
- `owned_scope`: `frontend/src/application/presenter/translation-job-setup/`
- `change_forbidden_scope`: live input review presenter、job run presenter、shared presenter 以外の application directory
- `size_estimate`: 2 files、850 changed lines 以下
- `completion_signal`: `translation-job-setup` presenter directory が消えている。live screen の view model 整形は変わっていない。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が presenter import 破綻を直せる。
- `execution_test_classification`: lower-level only

### `INT-001-wails-gateway-cleanup`

- `implementation_target`: dead page 用 Wails gateway、gateway DTO、gateway contract を削除する。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `ready_wave`: `wave-5`
- `depends_on`: `FE-004-dead-controller-store-contract-cleanup`, `FE-005-dead-usecase-cleanup`, `FE-006-dead-presenter-cleanup`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `frontend/src/controller/wails/translation-job-setup.gateway.ts` を削除し、completion の Wails gateway clause を閉じる。最初にする理由は generated `wailsjs` を触らず frontend adapter だけを閉じられるためである。
- `owned_scope`: `frontend/src/controller/wails/translation-job-setup.gateway.ts`, `frontend/src/controller/wails/gateway-dto/translation-job-setup/`, `frontend/src/application/gateway-contract/translation-job-setup/`
- `change_forbidden_scope`: `frontend/wailsjs/`, `internal/`, backend Wails bind、他 gateway
- `size_estimate`: 5 files、700 changed lines 以下
- `completion_signal`: frontend source から `TranslationJobSetupGatewayContract` と `createTranslationJobSetupGateway` の参照が消えている。generated Wails binding は変更していない。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が frontend Wails adapter import 境界の破綻を直せる。
- `execution_test_classification`: lower-level only
- `notes`: backend binding の削除は承認済み first unit に含めない。

### `INT-002-model-settings-card-unused-export-cleanup`

- `implementation_target`: `model-settings-card` の未使用 export を削除し、`frontend-local` の `lint:exports` を通す。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `ready_wave`: `wave-7`
- `depends_on`: `UT-004-dead-usecase-test-cleanup`, `UT-005-dead-gateway-test-cleanup`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/gateway-contract/model-settings-card/index.ts` から未使用の `cloneModelSettingsCardStates` export を削除する。最初にする理由は `knip --production --include-entry-exports` の失敗点が barrel export だからである。
- `owned_scope`: `frontend/src/application/gateway-contract/model-settings-card/index.ts`, `frontend/src/application/gateway-contract/model-settings-card/model-settings-card-policy.ts`
- `change_forbidden_scope`: `frontend/src/application/store/master-persona/`, `frontend/src/application/usecase/master-persona/`, `frontend/src/controller/wails/`, `frontend/wailsjs/`, `internal/`, docs 正本本文
- `size_estimate`: 2 files、40 changed lines 以下
- `completion_signal`: `cloneModelSettingsCardState` は live store 用に残り、未使用の `cloneModelSettingsCardStates` だけが消えている。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `validation_owner`: この handoff が `model-settings-card` の unused export gate 失敗を直せる。
- `execution_test_classification`: lower-level only

### `UT-001-app-shell-stale-setup-cleanup`

- `implementation_target`: `AppShell.test.ts` の stale setup を削除する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-2`
- `depends_on`: `FE-001-root-wiring-cleanup`
- `parallelizable_with`: `FE-003-dead-page-component-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/ui/views/AppShell.test.ts` の `createTranslationJobSetupScreenController: null` を削除し、completion の stale setup clause を閉じる。最初にする理由は assertion を変えずに不要前提だけを消せるためである。
- `owned_scope`: `frontend/src/ui/views/AppShell.test.ts`
- `change_forbidden_scope`: `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/stores/shell-state.ts`
- `size_estimate`: 1 file、80 changed lines 以下
- `completion_signal`: shell test の期待値は `未完了ジョブ一覧`、`job-run` hash 正規化、`現在の翻訳段階へ進む` を維持している。未使用 prop setup は残っていない。
- `validation_commands`: `npm --prefix frontend run test -- src/ui/views/AppShell.test.ts`
- `validation_owner`: この handoff が AppShell test の失敗を直せる。
- `execution_test_classification`: lower-level only

### `UT-002-dead-page-test-and-fixture-cleanup`

- `implementation_target`: dead page test と shared test fixture の job setup 部分を削除する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-4`
- `depends_on`: `FE-003-dead-page-component-cleanup`, `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup`
- `parallelizable_with`: `FE-005-dead-usecase-cleanup`, `UT-003-dead-store-presenter-test-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts` を削除し、completion の dead page test clause を閉じる。最初にする理由は live 導線に存在しない page の保護を先に外せるためである。
- `owned_scope`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`, `frontend/src/ui/screens/translation-job-setup/__fixtures__/translation-job-setup-panel-fixtures.ts`, `frontend/src/ui/screens/__fixtures__/screen-page-controller-fixtures.ts` の job setup fixture 部分
- `change_forbidden_scope`: live screen fixture、provider settings fixture、AppShell assertion
- `size_estimate`: 3 files、約 1200 changed lines
- `completion_signal`: `JobSetupPage` の単体 test と job setup 専用 fixture が消えている。live shell 導線 test は `UT-001` の期待値を維持している。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_owner`: この handoff が deleted page test と fixture import の構造破綻を直せる。
- `execution_test_classification`: lower-level only
- `notes`: deleted file 自体は直接 test target にできないため、局所 Vitest は置かない。最終検証の `frontend-local` で live shell test と残存 unit test をまとめて確認する。

### `UT-003-dead-store-presenter-test-cleanup`

- `implementation_target`: dead store / presenter の単体 test を削除する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-4`
- `depends_on`: `FE-004-dead-controller-store-contract-cleanup`, `FE-006-dead-presenter-cleanup`
- `parallelizable_with`: `FE-005-dead-usecase-cleanup`, `UT-002-dead-page-test-and-fixture-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/store/translation-job-setup/translation-job-setup.store.test.ts` を削除し、completion の dead store test clause を閉じる。最初にする理由は deleted store の保護を局所的に消せるためである。
- `owned_scope`: `frontend/src/application/store/translation-job-setup/translation-job-setup.store.test.ts`, `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts`
- `change_forbidden_scope`: live application tests
- `size_estimate`: 2 files、500 changed lines 以下
- `completion_signal`: dead store / presenter の単体 test が消えている。live application test は変わっていない。
- `validation_commands`: `npm --prefix frontend run test -- src/application/store src/application/presenter`
- `validation_owner`: この handoff が deleted store / presenter test の破綻を直せる。
- `execution_test_classification`: lower-level only

### `UT-004-dead-usecase-test-cleanup`

- `implementation_target`: dead usecase の単体 test を削除する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-6`
- `depends_on`: `FE-005-dead-usecase-cleanup`, `INT-001-wails-gateway-cleanup`
- `parallelizable_with`: `UT-005-dead-gateway-test-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts` を削除し、completion の dead usecase test clause を閉じる。最初にする理由は 1 ファイルで usecase test の削除を完了できるためである。
- `owned_scope`: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- `change_forbidden_scope`: live usecase tests、backend tests
- `size_estimate`: 1 file、約 1075 changed lines。注意規模だが 1 deleted usecase test に閉じるため 1 handoff とする。
- `completion_signal`: dead usecase の単体 test が消えている。live job creation usecase test は変わっていない。
- `validation_commands`: `npm --prefix frontend run test -- src/application/usecase`
- `validation_owner`: この handoff が deleted usecase test の破綻を直せる。
- `execution_test_classification`: lower-level only

### `UT-005-dead-gateway-test-cleanup`

- `implementation_target`: dead Wails gateway の単体 test を削除する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-6`
- `depends_on`: `INT-001-wails-gateway-cleanup`
- `parallelizable_with`: `UT-004-dead-usecase-test-cleanup`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/controller/wails/translation-job-setup.gateway.test.ts` を削除し、completion の dead gateway test clause を閉じる。最初にする理由は Wails adapter 削除後の test 破綻を局所的に閉じられるためである。
- `owned_scope`: `frontend/src/controller/wails/translation-job-setup.gateway.test.ts`
- `change_forbidden_scope`: other gateway tests、`frontend/wailsjs/`
- `size_estimate`: 1 file、350 changed lines 以下
- `completion_signal`: dead Wails gateway test が消えている。他 gateway tests は変わっていない。
- `validation_commands`: `npm --prefix frontend run test -- src/controller/wails`
- `validation_owner`: この handoff が deleted gateway test の破綻を直せる。
- `execution_test_classification`: lower-level only

## Final Validation For Refactor Lane

全 handoff 完了後に、呼び出し元 `refactor_lane` は次を最終検証として扱う。

- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite structure`
- `npm --prefix frontend run build-storybook`

## Completion Packet

Codex 実装系レーンは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `not_required` または停止理由。ユーザー向け導線を変えないため新規 UI 人間操作 E2E は不要。
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `harness_gate_result`
- `residual_risks`
- `completion_evidence`
- `docs_changes`: `none`

## 未決事項

- 実装範囲に関する未決事項はなし。
- `detail-spec-diff.md` の `Q-001` は docs 正本化候補の扱いであり、今回の実装引き継ぎには混ぜない。`FSD-005` は code 修正対象外として固定済みである。
