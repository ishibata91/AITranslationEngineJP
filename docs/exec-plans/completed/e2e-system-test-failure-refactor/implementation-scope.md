# Implementation Scope: e2e-system-test-failure-refactor

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: fixed decisions accepted by human input
- `approval_record`: `EF-001` から `EF-005` は全てリファクタ範囲へ含める。仕様乖離整理は不要。実行型は `コード併走型`。
- `codex_entry`: `.codex/skills/refactor-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `plan`: `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
- `structure_quality_investigation`: `docs/exec-plans/active/e2e-system-test-failure-refactor/structure-quality-investigation.md`
- `test_quality_investigation`: `docs/exec-plans/active/e2e-system-test-failure-refactor/test-quality-investigation.md`
- `previous_result`: `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
- `e2e_test_design`: `docs/e2e-test-design/test-design.csv`
- `e2e_test_guidelines`: `docs/e2e-test-guidelines.md`
- `test_coding_guidelines`: `docs/coding-guidelines-tests.md`
- `detail_specs`: `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-output-artifact.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `screen_design`: `docs/screen-design/screens/provider-settings.md`, `docs/screen-design/screens/master-persona.md`, `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/job-run.md`, `docs/screen-design/screens/output-management.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `detail_spec_diff`: `N/A`
- `screen_design_diff`: `N/A`

## Fixed Decisions

- `unanswered_questions`: `0`
- `EF-001` から `EF-005` だけを解決対象にする。
- 新規仕様は作らない。既存詳細仕様、既存画面設計、E2E 観点表、調査結果の範囲で修正する。
- docs 正本文、`.codex/`、プロダクト外の workflow 契約は変更しない。
- 実外部 API、実 secret、実利用者データへ到達する経路は作らない。
- UI 人間操作 E2E は、画面操作を入口にし、DB seed と Wails mock は前提準備または外部境界代替に限定する。
- `contract_freeze.status`: frozen
- `contract_freeze.freeze_source`: fixed decisions in this task input and `./plan.md`
- `contract_freeze.frozen_public_seams`: Wails controller DTO, frontend gateway contract, `data-testid`, system-test seed names, scenario Wails mock names

## EF Scope Matrix

| EF | 証明対象 | 変更候補 | 禁止変更 | 検証 command |
| --- | --- | --- | --- | --- |
| `EF-001` | 不正 endpoint の保存操作で、入力不正が表示され、保存済み状態へ更新されない。 | provider settings の保存前入力検証、gateway wiring、保存通知と接続確認通知の分離、E2E assertion 補強。 | 実 endpoint への接続確認を保存時に強制しない。secret 本体を UI、DTO、log、error summary に出さない。 | `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite system-test` |
| `EF-002` | マスターペルソナで AI サービス、モデル、実行方法を選択でき、`gemini-test` で生成へ進める。 | model select 活性化条件、provider 選択後の model list 反映待機、scenario mock 前提の観測補強。 | JSON fixture を過剰化しない。実 AI provider へ接続しない。 | `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite system-test` |
| `EF-003` | paused job の再開操作で、再開結果 feedback を表示し、必要な場合だけ job run shell へ遷移する。 | resume action read model、frontend 遷移条件、feedback assertion 補強。 | 再開実行 runner を新設しない。状態遷移を根拠なしに永続化しない。 | `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite system-test` |
| `EF-004` | 出力管理で候補行、出力操作、差分 preview が scenario mock または backend gateway 経由で表示される。 | output management gateway wiring、gateway 未接続状態の表示、candidate list 待機、diff row click assertion。 | xTranslator 本体を起動しない。実 filesystem 出力を主証明にしない。AI サービス境界へ接続しない。 | `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite system-test` |
| `EF-005` | 翻訳管理から現在の翻訳段階へ進む操作で、対象 job と action が見つかり、各 phase 画面へ進める。 | system-test seed と scenario mock の job family 整合、open current phase 導線観測、E2E-UC と test 内容の対応整理。 | `EF-005` 以外の翻訳段階仕様を変更しない。completed job を未完了一覧へ恒久復帰させない。実 AI 実行を開始しない。 | `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite system-test` |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `H-BE-001`, `H-BE-002`, `H-FE-001`, `H-FE-002`, `H-FE-003`, `H-FE-004` | なし | `H-BE-001 <-> H-BE-002`; `H-FE-002 <-> H-FE-003`; `H-FE-002 <-> H-FE-004` | `H-FE-001` は `App.svelte` shared composition root を触るため他 frontend handoff と並列不可 |
| `wave-2` | `H-INT-001`, `H-INT-002`, `H-INT-003` | 対応する backend と frontend handoff の完了 | `H-INT-001 <-> H-INT-002` | `H-INT-003` は job family と phase gateway の shared mock を触るため `H-INT-002` と並列不可 |
| `wave-3` | `H-ST-001`, `H-ST-002`, `H-ST-003`, `H-ST-004`, `H-ST-005`, `H-UT-001`, `H-UT-002` | `wave-1` と `wave-2` の完了 | `H-ST-001 <-> H-ST-002`; `H-ST-003 <-> H-ST-004`; `H-UT-001 <-> H-UT-002` | `H-ST-005` は shared scenario mock と phase specs を触るため他 scenario handoff と並列不可 |
| `wave-4` | `H-FINAL-001` | `wave-3` の完了 | なし | 最終検証は広域判定条件共有 |

## Handoffs

### `H-BE-001`: `EF-001` provider settings 保存境界

- `implementation_target`: 不正 endpoint を保存成功として扱わない backend 保存境界を固定する。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `docs/detail-specs/ai-provider-settings-management.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider id, credential state, validation state, saved state, failure kind
  - `secret_values_for_provider_external_api_internal_auth`: API key input, stored credential body
  - `secret_resolution_owner_layer`: backend provider settings service and secret store
  - `forbidden_outputs`: API key body, credential body, external provider raw response, secret reference value that can resolve a credential
- `owned_scope`: `internal/service/provider_settings_service.go`, `internal/usecase/provider_settings_usecase.go`, `internal/controller/wails/provider_settings_controller.go`, related backend tests.
- `depends_on`: なし。
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `H-BE-002`
- `parallel_blockers`: なし。
- `first_action`: `internal/service/provider_settings_service.go` の save input validation clause を作り、不正 endpoint で保存済み状態を更新しない完了条件を閉じる。
- `implementation_observation`: 実装前に `SaveProviderSettings` が endpoint 形式不正をどの error kind で返すべきかを既存 test と DTO から確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`: 不正 endpoint の保存要求は拒否結果になり、保存済み endpoint と secret 状態を更新しない。error summary は secret を含まない。
- `estimated_size`: 4 files, 180 changed lines
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後

### `H-BE-002`: `EF-003` resume action read model

- `implementation_target`: paused job の再開操作結果を、UI が feedback と遷移可否に使える read model として返す。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `docs/detail-specs/translation-job-management.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: AI service label, model label, credential state label
  - `secret_values_for_provider_external_api_internal_auth`: credential body, credential resolvable reference
  - `secret_resolution_owner_layer`: phase start or retry usecase only
  - `forbidden_outputs`: credential body, credential resolvable reference, external provider raw response
- `owned_scope`: `internal/service/translation_job_management_service.go`, `internal/usecase/translation_job_management_usecase.go`, `internal/controller/wails/translation_job_management_controller.go`, related backend tests.
- `depends_on`: なし。
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `H-BE-001`
- `parallel_blockers`: なし。
- `first_action`: `internal/service/translation_job_management_service.go` の `ResumeJob` success read model clause を作り、paused job の feedback 判定を閉じる。
- `implementation_observation`: 実装前に `Paused` job の detail が `jobRunTarget` 相当の current phase を返しているかを unit test または trace で確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`: paused job の `ResumeJob` は warning 固定ではなく、再開入口成立を表す result を返す。runner 実行や実 AI 実行は行わない。
- `estimated_size`: 5 files, 220 changed lines
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後

### `H-FE-001`: `EF-001` と `EF-004` root composition gateway wiring

- `implementation_target`: root `App.svelte` の既定 factory が production factory と同じ Wails gateway 接続を使う。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `docs/architecture.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/provider-settings.md`, `docs/screen-design/screens/output-management.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider id, credential state, validation state, output artifact state
  - `secret_values_for_provider_external_api_internal_auth`: provider API key body
  - `secret_resolution_owner_layer`: backend provider settings service
  - `forbidden_outputs`: API key body, credential body, external provider raw response
- `owned_scope`: `frontend/src/ui/App.svelte`, `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, app shell tests if needed.
- `depends_on`: なし。
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし。
- `parallel_blockers`: shared_contract_change
- `first_action`: `frontend/src/ui/App.svelte` の provider settings fallback factory を Wails gateway 接続へ変え、`EF-001` の gateway 未接続 clause を閉じる。
- `implementation_observation`: 実装前に system-test 起動時の root が `main.ts` から注入された factory を使うか、`App.svelte` fallback を使うかを trace で確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: provider settings と output management は fallback 経路でも `null gateway` にならない。gateway 未接続表示は test fixture や明示的 null 注入時だけ残る。
- `estimated_size`: 3 files, 120 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後
- `notes`: 2 つの EF を含む理由は、同じ composition root の同じ欠陥を分割すると `App.svelte` の shared seam を同時変更するためである。

### `H-FE-002`: `EF-002` master persona model select readiness

- `implementation_target`: provider 選択後に model list が反映され、`gemini-test` を選択できる状態へ進む。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-013`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/master-persona.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider option, model option, credential status
  - `secret_values_for_provider_external_api_internal_auth`: provider credential body
  - `secret_resolution_owner_layer`: backend provider settings service
  - `forbidden_outputs`: credential body, external provider raw response
- `owned_scope`: `frontend/src/application/usecase/master-persona/`, `frontend/src/application/presenter/master-persona/`, `frontend/src/ui/screens/master-persona/`, `tests/system/master-persona.spec.ts`, `tests/system/support/master-persona-page.ts`.
- `depends_on`: なし。
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `H-FE-003`, `H-FE-004`
- `parallel_blockers`: なし。
- `first_action`: `frontend/src/application/presenter/master-persona/master-persona.presenter.ts` の model select disabled 判定 clause を、model list status と option 有無を分けて閉じる。
- `implementation_observation`: 実装前に provider 選択後の `MasterPersonaListProviderModels` 応答が store に反映される時系列を trace する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: `gemini` 選択後、model list success と `gemini-test` option が view model へ反映される。生成ボタンの有効化条件は provider と model の選択完了に閉じる。
- `estimated_size`: 5 files, 220 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-FE-003`: `EF-003` resume feedback and job run transition

- `implementation_target`: resume action result を見て feedback を表示し、成立時だけ job run shell へ遷移する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `docs/detail-specs/translation-job-management.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/job-run.md`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`: `frontend/src/ui/screens/translation-job-management/`, `frontend/src/application/usecase/translation-job-management/`, `frontend/src/application/presenter/translation-job-management/`, `tests/system/translation-job-management.spec.ts`, `tests/system/support/translation-job-management-page.ts`.
- `depends_on`: `H-BE-002`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `H-FE-002`, `H-FE-004`
- `parallel_blockers`: なし。
- `first_action`: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte` の resume branch を action result 成否 clause に合わせる。
- `implementation_trace`: 実装中に `requestResume()` 後の store feedback、selected job detail、job run target を trace し、shell 非表示の原因を分離する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: resume 結果 feedback が表示される。warning または失敗 result では shell へ進まない。成立 result では job run target が渡される。
- `estimated_size`: 4 files, 180 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-FE-004`: `EF-005` current phase opening state

- `implementation_target`: translation management から現在段階を開く操作で、対象 job と action 可用性を画面状態から判別できる。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `docs/detail-specs/translation-job-management.md`, phase detail specs
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/job-run.md`, phase screen design files
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider label, model label, credential state label
  - `secret_values_for_provider_external_api_internal_auth`: credential body
  - `secret_resolution_owner_layer`: phase usecase when start or retry is invoked
  - `forbidden_outputs`: credential body, provider raw response, prompt body
- `owned_scope`: `frontend/src/ui/screens/translation-job-management/`, `frontend/src/ui/screens/job-run/`, `frontend/src/controller/body-translation-phase/`, `frontend/src/application/usecase/body-translation-phase/`, phase screen presenters if needed.
- `depends_on`: なし。
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `H-FE-002`, `H-FE-003`
- `parallel_blockers`: なし。
- `first_action`: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte` の `handleOpenJob` clause を、card 表示不備と action 不可を分離できる状態へ閉じる。
- `implementation_observation`: 実装前に `system-test-term`, `system-test-persona`, `system-test-body-*` の card が表示されるか、button 名が `現在の翻訳段階へ進む` かを trace する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: job card 未表示、action 不可、phase shell 未表示が別々に診断できる。正常な ready/running/failed job では current phase を開ける。
- `estimated_size`: 6 files, 260 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-INT-001`: `EF-001` provider settings public seam

- `implementation_target`: provider settings の Wails DTO、frontend gateway、backend controller の error / success contract を揃える。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `docs/detail-specs/ai-provider-settings-management.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/provider-settings.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider id, saved state, credential state, validation state, error kind
  - `secret_values_for_provider_external_api_internal_auth`: API key body
  - `secret_resolution_owner_layer`: backend provider settings service
  - `forbidden_outputs`: API key body, credential body, external provider raw response
- `owned_scope`: `frontend/src/controller/wails/provider-settings.gateway.ts`, `frontend/src/controller/wails/gateway-dto/provider-settings/`, `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts`, `internal/controller/wails/provider_settings_controller.go`, generated seam if required by existing workflow.
- `depends_on`: `H-BE-001`, `H-FE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `H-INT-002`
- `parallel_blockers`: なし。
- `first_action`: frontend provider settings gateway の save error mapping clause を、backend error kind と UI error state に接続する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`: 不正 endpoint の save failure は UI の入力不正表示へ写像され、保存成功 notice へ写像されない。
- `estimated_size`: 5 files, 220 changed lines
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後

### `H-INT-002`: `EF-004` output management gateway seam

- `implementation_target`: output management の frontend gateway、Wails mock、screen controller factory の接続を揃える。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `docs/detail-specs/translation-output-artifact.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/output-management.md`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`: `frontend/src/controller/wails/translation-output-artifact.gateway.ts`, `frontend/src/controller/wails/gateway-dto/translation-output-artifact/`, `tests/system/support/scenario-wails-mocks.ts`, output controller tests if needed.
- `depends_on`: `H-FE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `H-INT-001`
- `parallel_blockers`: なし。
- `first_action`: `tests/system/support/scenario-wails-mocks.ts` の `TranslationOutputArtifactController` seam を、画面が呼ぶ controller 名と method 名に一致させる。
- `implementation_observation`: 実装前に output management screen が `GetTranslationOutputReview` を呼んでいる controller 名を Playwright console trace または frontend diagnostic trace で確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: output management 画面から scenario mock の completed jobs と diff rows が取得される。
- `estimated_size`: 4 files, 180 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-INT-003`: `EF-005` job family and phase gateway seam

- `implementation_target`: system-test seed、scenario Wails mock、job run / phase gateway の job family を揃える。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `docs/detail-specs/translation-job-management.md`, phase detail specs
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design`: `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/job-run.md`, phase screen design files
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider label, model label, credential state label
  - `secret_values_for_provider_external_api_internal_auth`: credential body
  - `secret_resolution_owner_layer`: phase usecase on start or retry
  - `forbidden_outputs`: credential body, provider raw response, prompt body
- `owned_scope`: `scripts/test/seed-system-test-db/main.go`, `tests/system/support/scenario-wails-mocks.ts`, frontend phase gateway DTO if needed.
- `depends_on`: `H-FE-004`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: なし。
- `parallel_blockers`: shared_contract_change
- `first_action`: `tests/system/support/scenario-wails-mocks.ts` の `ListIncompleteJobs` clause を、`job-run-shell.spec.ts` と `translation-phases.spec.ts` が探す job family に一致させる。
- `implementation_observation`: 実装前に real backend seed の job family と scenario mock の job family を一覧化し、completed job は output management 側の候補として扱うかを確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: `system-test-term`, `system-test-persona`, `system-test-body-failed`, `system-test-body-running` は current phase action を持つ。completed body job は未完了一覧へ恒久復帰させない。
- `estimated_size`: 4 files, 220 changed lines
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後

### `H-ST-001`: `EF-001` provider settings scenario test

- `implementation_target`: 不正 endpoint 保存時の error 表示と保存済み detail 非更新を E2E で証明する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-028`
- `frontend_required_sources`: `docs/screen-design/screens/provider-settings.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider id, credential state, validation state
  - `secret_values_for_provider_external_api_internal_auth`: API key body
  - `secret_resolution_owner_layer`: backend provider settings service
  - `forbidden_outputs`: API key body, credential body, external provider raw response
- `owned_scope`: `tests/system/frontend-backend-connection.spec.ts`, `tests/system/support/provider-settings-page.ts`.
- `depends_on`: `H-INT-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-ST-002`
- `parallel_blockers`: なし。
- `first_action`: `tests/system/frontend-backend-connection.spec.ts` の `E2E-UC-028` assertion clause に detail region 非更新確認を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`: `E2E-UC-028` は summary error と detail non-update の両方を観測する。
- `estimated_size`: 2 files, 80 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-ST-002`: `EF-002` master persona scenario test

- `implementation_target`: model select が選択可能になる前提を明示し、`gemini-test` で生成結果を確認する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-013`
- `frontend_required_sources`: `docs/screen-design/screens/master-persona.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider option, model option, credential status
  - `secret_values_for_provider_external_api_internal_auth`: provider credential body
  - `secret_resolution_owner_layer`: scenario Wails mock for fake boundary
  - `forbidden_outputs`: credential body, external provider raw response
- `owned_scope`: `tests/system/master-persona.spec.ts`, `tests/system/support/master-persona-page.ts`, `tests/system/support/scenario-wails-mocks.ts` if needed.
- `depends_on`: `H-FE-002`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-ST-001`
- `parallel_blockers`: なし。
- `first_action`: `tests/system/support/master-persona-page.ts` の `selectAISettings` clause に model select enabled / option visible wait を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`: `E2E-UC-013` は model select readiness と generation result の両方を観測する。
- `estimated_size`: 3 files, 120 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-ST-003`: `EF-003` resume scenario test

- `implementation_target`: paused job の再開結果 feedback と job run shell 表示条件を分けて証明する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-019`
- `frontend_required_sources`: `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/job-run.md`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`: `tests/system/translation-job-management.spec.ts`, `tests/system/support/translation-job-management-page.ts`, `tests/system/support/job-run-shell-page.ts`.
- `depends_on`: `H-FE-003`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-ST-004`
- `parallel_blockers`: なし。
- `first_action`: `tests/system/translation-job-management.spec.ts` の `E2E-UC-019` assertion clause に feedback notification の確認を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`: resume action failure、feedback absence、shell transition failure を別 assertion で診断できる。
- `estimated_size`: 3 files, 100 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-ST-004`: `EF-004` output management scenario test

- `implementation_target`: candidate row 表示、XML 出力 / 再出力、diff row click、disabled 条件を E2E で証明する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-023` から `E2E-UC-025`, `E2E-UC-042` から `E2E-UC-044`
- `frontend_required_sources`: `docs/screen-design/screens/output-management.md`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`: `tests/system/output-management.spec.ts`, `tests/system/support/output-management-page.ts`.
- `depends_on`: `H-INT-002`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-ST-003`
- `parallel_blockers`: なし。
- `first_action`: `tests/system/output-management.spec.ts` の `E2E-UC-025` clause に diff row click と selected job 維持確認を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`: output candidate list 接続成立を先に観測し、各 test が候補行存在を前提に失敗原因を切り分けられる。
- `estimated_size`: 2 files, 120 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-ST-005`: `EF-005` job run and phase scenario test

- `implementation_target`: translation management から各 phase を開く導線と、`E2E-UC-045` から `E2E-UC-053` の証明対象を揃える。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `docs/e2e-test-design/test-design.csv` の `E2E-UC-045` から `E2E-UC-053`
- `frontend_required_sources`: `docs/screen-design/screens/job-run.md`, phase screen design files
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider label, model label, credential state label
  - `secret_values_for_provider_external_api_internal_auth`: credential body
  - `secret_resolution_owner_layer`: phase usecase on start or retry
  - `forbidden_outputs`: credential body, provider raw response, prompt body
- `owned_scope`: `tests/system/job-run-shell.spec.ts`, `tests/system/translation-phases.spec.ts`, `tests/system/support/translation-job-management-page.ts`, `tests/system/support/job-run-shell-page.ts`, `tests/system/support/translation-phase-pages.ts`, `tests/system/support/scenario-wails-mocks.ts`.
- `depends_on`: `H-INT-003`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし。
- `parallel_blockers`: shared_contract_change
- `first_action`: `tests/system/job-run-shell.spec.ts` の `openJobRun` clause に card visibility と open action availability の観測を追加する。
- `implementation_observation`: 実装前に `E2E-UC-045` から `E2E-UC-047` が現行 test 内容と不一致である箇所を、修正方針として test 名、コメント、assertion の対応で閉じる。
- `validation_commands`: `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`: 対象 job 不在、action 不在、phase screen 不在を別 failure として診断できる。`E2E-UC-045` から `E2E-UC-053` は CSV の期待値または承認済み現行読み替えに対応する。
- `estimated_size`: 6 files, 320 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後

### `H-UT-001`: frontend unit protection

- `implementation_target`: frontend の分岐変更を unit test で保護する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `spec_basis`: source artifacts listed above
- `frontend_required_sources`: relevant screen design files
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider state, credential state label, model option
  - `secret_values_for_provider_external_api_internal_auth`: credential body
  - `secret_resolution_owner_layer`: fake gateway only in frontend unit tests
  - `forbidden_outputs`: credential body, provider raw response
- `owned_scope`: frontend tests adjacent to changed frontend usecase, presenter, controller, or page files.
- `depends_on`: `H-FE-001`, `H-FE-002`, `H-FE-003`, `H-FE-004`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-UT-002`
- `parallel_blockers`: なし。
- `first_action`: 最初に変更された frontend usecase または presenter の新規分岐に対し、1 つの failing unit test を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`: provider settings rejection、master persona readiness、resume transition、current phase open state のうち、変更した frontend 分岐が unit test で守られる。
- `estimated_size`: 6 files, 300 changed lines
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `H-UT-002`: backend unit protection

- `implementation_target`: backend の provider settings と resume read model 分岐を unit test で保護する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `spec_basis`: `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-management.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider id, credential state, validation state, action result reason
  - `secret_values_for_provider_external_api_internal_auth`: credential body
  - `secret_resolution_owner_layer`: backend fake secret store in tests
  - `forbidden_outputs`: credential body, external provider raw response
- `owned_scope`: backend tests adjacent to changed service, usecase, controller files.
- `depends_on`: `H-BE-001`, `H-BE-002`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-UT-001`
- `parallel_blockers`: なし。
- `first_action`: `internal/service/provider_settings_service_test.go` または related test に、不正 endpoint 保存拒否の test clause を追加する。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`: backend 変更分岐は unit test で保護される。secret 本体は response と log 相当値に現れない。
- `estimated_size`: 4 files, 220 changed lines
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `H-FINAL-001`: final validation and review input

- `implementation_target`: `EF-001` から `EF-005` の修正結果を system-test と local harness で確認する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: this implementation-scope
- `frontend_required_sources`: all related screen design files
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: all non-secret read model values
  - `secret_values_for_provider_external_api_internal_auth`: any credential body used by fake boundary only
  - `secret_resolution_owner_layer`: fake boundary or backend secret store
  - `forbidden_outputs`: credential body, external provider raw response, prompt body
- `owned_scope`: validation only. No product code changes.
- `depends_on`: all previous handoffs.
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし。
- `parallel_blockers`: broad_gate_shared
- `first_action`: `python3 scripts/harness/run.py --suite system-test` を実行し、`EF-001` から `EF-005` の結果を completion packet に分けて記録する。
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`; `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite system-test`; `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: system-test が通過する。sandbox や Wails readiness で止まる場合は `FAIL_ENVIRONMENT` と product failure を分けて記録する。
- `estimated_size`: 0 files, 0 changed lines
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: final validation

## Completion Packet

Codex 実装系レーンは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate として扱う
- `harness_gate_result`: system-test が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` とする
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`

## Review Notes

- 人間レビューは `approve` により完了した。
- `EF-003` の resume は runner 新設ではなく、既存 read model と UI 遷移の範囲に閉じる。
- `EF-005` の completed job は、未完了一覧の恒久対象にしない。completed job の確認は output management または translation complete の導線で扱う。
- `EF-004` の XML 出力証明は fake 境界へ閉じる。実 filesystem 出力を主証明へ昇格しない。
