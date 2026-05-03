# Implementation Scope: translation-job-setup-phase-provider-settings

- `skill`: implementation-scope
- `status`: ready-for-implement-lane
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: 2026-05-04 に `scenario-design.md` と `ui-design.md` を人間承認済み。UI 設計レビューは終了済み。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `ui_prototype`: `./prototype.svelte`
- `ui_mock_data`: `N/A`
- `ui_sample_data`: `prototype.svelte` 内の `data-ui-prototype-sample-data-root`
- `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `required_reading_basis`:
  - `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `internal/service/translation_job_setup_service.go`
  - `internal/usecase/translation_job_setup_contract.go`
  - `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`

## Fixed Decisions

- `needs_human_decision`: `0`
- Job Setup は master-persona の provider / model / credential / execution mode を、既定値、保存元、validation 対象、secret 解決元として使わない。
- Job Setup は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 phase ごとに provider、model、credential 参照、execution mode、batch mode、model list 状態を保存、読み出し、表示する。
- phase 実行時は対象 phase の Job Setup 保存設定だけを参照し、他 phase と master-persona provider 設定を参照しない。
- provider 別 model list API は API key 必須 provider の credential 参照が解決できる場合だけ外部取得する。
- LM Studio は credential 不要 provider とし、API key 入力、API key 未設定警告、credential 選択、credential missing failure に出さない。
- batch mode は Gemini と xAI だけで checkbox により明示する。対象外 provider では保存値へ残さない。
- model 候補取得失敗、API key 未設定、credential 参照不能では model 手入力へ進めない。
- 3 phase の不足がなければ、翻訳段階ごとの個別確定ボタンなしで Job Setup 作成を許可する。
- API key 平文、secret 本体、復号可能値、provider raw request / response、raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
- UI 実装は承認済み `ui-design.md` を根拠にし、`prototype.svelte` とサンプル値は product code、fixture、default state、test data へ移植しない。

## Existing Boundary Evidence

- `internal/service/translation_job_setup_service.go` は `MasterPersonaAISettingsRecord` から単一 runtime option を作り、`master-persona:<provider>` の secret key を解決している。
- `internal/usecase/translation_job_setup_contract.go` と `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts` は runtime selection を単一 provider / model / execution mode として扱っている。
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` は provider / model / execution mode と credential reference を単一 select として表示している。
- phase 詳細仕様は各 phase が Job Setup の provider、model、execution mode 要約を使う前提である。今回の実装は、その参照元を phase 別 Job Setup 保存設定へ固定する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `provider-settings-contract-freeze` | `なし` | `なし` | `なし` |
| `wave-2` | `backend-provider-settings-core`, `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance` | `provider-settings-contract-freeze` | `backend-provider-settings-core <-> frontend-provider-settings-ui`, `backend-provider-settings-core <-> tests-provider-settings-api-acceptance`, `frontend-provider-settings-ui <-> tests-provider-settings-api-acceptance` | `なし` |
| `wave-3` | `integration-provider-settings-boundary` | `backend-provider-settings-core`, `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance` | `なし` | `backend_frontend_order` |
| `wave-4` | `tests-provider-settings-ui-and-regression` | `integration-provider-settings-boundary` | `なし` | `backend_frontend_order` |
| `wave-5` | `final-validation-and-report` | `backend-provider-settings-core`, `frontend-provider-settings-ui`, `integration-provider-settings-boundary`, `tests-provider-settings-api-acceptance`, `tests-provider-settings-ui-and-regression` | `なし` | `broad_gate_shared` |

## Handoffs

### `provider-settings-contract-freeze`

- `implementation_target`: phase 別 provider 設定、provider model list、credential 状態、batch mode、create / summary の公開接点を固定する。
- `implementation_artifact`: `contract_freeze`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_prototype`: `./prototype.svelte`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-TJSPPS-001` から `SCN-TJSPPS-008`、`./ui-design.md` の `UI Contract`
  - `architecture_layer_basis`: Wails public API、backend usecase contract、frontend gateway contract、frontend state contract の transport boundary。
  - `frozen_public_seams`:
    - `GetTranslationJobSetupOptions`: 入力、共通辞書、共通ペルソナ、provider capability、phase 別 draft 初期値、credential 参照状態を返す。master-persona provider 設定は返さない。
    - `ListTranslationJobSetupProviderModels`: phase、provider、credential 参照状態、request token を受け取り、model list の `not_updated | loading | success | failed | credential_missing | credential_not_required` 相当の状態、候補、redacted failure を返す。
    - `ValidateTranslationJobSetup`: 3 phase の provider、model、credential ref、execution mode、batch mode、model list source token を受け取り、phase 別 blocking failure、create 可否、stale model list 判定を返す。
    - `CreateTranslationJob`: validation pass 済みの 3 phase 設定断面から `Ready` job と phase runtime snapshot を保存する。
    - `GetTranslationJobSetupSummary`: 作成済み job の phase 別 provider、model、credential 参照状態、execution mode、batch mode、model list source summary を read-only で返す。
    - phase 開始境界は対象 phase の runtime snapshot だけを読み、master-persona と他 phase の provider 設定を参照しない。
    - error kind は `phase_runtime_missing`, `model_list_credential_missing`, `model_list_failed`, `model_selection_stale`, `credential_missing`, `provider_mode_unsupported`, `provider_unreachable`, `validation_stale`, `ready_required` を区別できる。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、phase id、credential ref、credential configured / missing / not_required、execution mode、batch mode、model list request token、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、復号可能値、provider SDK token、provider request authorization。
  - `secret_resolution_owner_layer`: backend service / infra provider boundary。frontend、DTO、read model は secret 本体を受け取らない。
  - `forbidden_outputs`: log、error summary、audit、request capture、URL、DTO、UI、read model、fake transport log。
- `owned_scope`:
  - `internal/usecase/translation_job_setup_contract.go`
  - `internal/controller/wails/translation_job_setup_controller.go`
  - `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
  - `frontend/src/application/contract/translation-job-setup/*`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `internal/usecase/translation_job_setup_contract.go` に phase 別 runtime selection と provider model list response 型を追加し、`completion_signal` の「3 phase runtime DTO と model list DTO が固定される」を最初に閉じる。理由は backend、frontend、test が同じ公開接点に依存するため。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/controller/wails -run 'TranslationJobSetup|JobSetup|ProviderSettings|ModelList'`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - 3 phase の provider、model、credential ref、execution mode、batch mode を表す DTO shape が backend と frontend gateway contract で一致している。
  - provider model list API の request / response、失敗状態、stale response 判定用 token が固定されている。
  - LM Studio は credential ref を required field にせず、credential missing error kind の対象外である。
  - Gemini と xAI 以外の batch mode は contract 上 unsupported として扱える。
  - secret 本体を出さない field obligation、nullability、error kind が controller unit test または type check で固定されている。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `6-10 files`、`300-650 changed lines`。
  - この handoff は contract freeze だけを扱い、永続化、UI 実装、外部 API 実体接続、test 実装を含めない。
  - `本番経路`: Wails controller / DTO -> usecase contract -> frontend gateway contract。

### `backend-provider-settings-core`

- `implementation_target`: phase 別 provider 設定の永続化、読み出し、validation、provider model list 取得、phase 実行時参照を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_prototype`: `N/A`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `provider-settings-contract-freeze`
  - `architecture_layer_basis`: backend service / repository / infra provider boundary。
  - `frozen_public_seams`: `provider-settings-contract-freeze` の completion signal を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential ref、credential configured / missing / not_required、provider、model、execution mode、batch mode。
  - `secret_values_for_provider_external_api_internal_auth`: secret store から解決した API key、provider authorization header。
  - `secret_resolution_owner_layer`: `internal/service` から `internal/infra/ai` へ渡す境界。
  - `forbidden_outputs`: DTO、UI、structured log、fake transport log、error summary、URL、保存要約。
- `owned_scope`:
  - `internal/service/translation_job_setup_service.go`
  - `internal/service/translation_job_setup_*`
  - `internal/usecase/translation_job_setup_*`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/infra/ai/*`
  - `internal/bootstrap/app_controller.go`
  - 必要な SQLite migration と schema-first backend test
- `depends_on`: `provider-settings-contract-freeze`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance`
- `parallel_blockers`: `なし`
- `first_action`: `internal/service/translation_job_setup_service.go` の master-persona AI settings 読み込み経路を phase 別 Job Setup runtime source へ置き換える focused test を追加し、`completion_signal` の「master-persona provider 設定を参照しない」を最初に閉じる。理由は今回の責務境界違反を先に消す必要があるため。
- `validation_commands`:
  - `go test ./internal/service ./internal/usecase ./internal/repository ./internal/infra/ai -run 'TranslationJobSetup|ProviderSettings|ModelList|PhaseRuntime'`
- `completion_signal`:
  - `MasterPersonaAISettingsRecord` と `master-persona:<provider>` は Job Setup の provider / model 既定値、保存元、secret 解決元にならない。
  - create 時に 3 phase の provider、model、credential ref、execution mode、batch mode、model list source summary が永続化される。
  - summary と phase 開始境界は対象 phase の保存済み runtime snapshot だけを読む。
  - API key 必須 provider は credential 参照が解決できる場合だけ external model list request を送る。
  - credential missing、credential unavailable、credential expired、model list failed では external request を送らない、または secret 非露出の failure として返す。
  - LM Studio は secret store 参照なしで credential not_required とし、endpoint success / failure を API key missing と別カテゴリにする。
  - Gemini と xAI だけ batch mode を pass にし、対象外 provider の stale batch 値は保存しない。
  - paid real API は local tests で呼ばず、provider list は real provider のまま、transport / SDK seam だけ fake に差し替える。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-24 files`、`850-1400 changed lines`。
  - backend 永続化、secret gate、provider model list、phase runtime snapshot は同じ backend use case の整合条件であり、1 handoff にまとめる。
  - frontend UI、Wails 接続、frontend test は含めない。
  - `本番経路`: usecase -> service -> repository / secret store / infra provider -> SQLite。

### `frontend-provider-settings-ui`

- `implementation_target`: 承認済み UI 契約に従い、Job Setup 画面を phase 別 provider 設定 UI へ差し替える。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_prototype`: `./prototype.svelte`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `provider-settings-contract-freeze`
  - `architecture_layer_basis`: frontend gateway contract、controller、usecase、presenter、store、Svelte screen。
  - `frozen_public_seams`: phase 別 DTO、model list DTO、validation response、create response、summary response。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: AIサービス名、model id、APIキー状態、credential ref、not_required、batch mode、failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: backend service。frontend は secret 本体を扱わない。
  - `forbidden_outputs`: API key 平文、復号可能値、provider raw request / response、raw prompt、internal log id。
- `owned_scope`:
  - `frontend/src/application/store/translation-job-setup/*` の product code。`*test.ts` は除く。
  - `frontend/src/application/presenter/translation-job-setup/*` の product code。`*test.ts` は除く。
  - `frontend/src/application/usecase/translation-job-setup/*` の product code。`*test.ts` は除く。
  - `frontend/src/controller/translation-job-setup/*` の product code。`*test.ts` は除く。
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
  - `frontend/src/ui/screens/translation-job-setup/` の product CSS / helper。`*test.ts` は除く。
- `depends_on`: `provider-settings-contract-freeze`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-provider-settings-core`, `tests-provider-settings-api-acceptance`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts` を参照する screen state を phase 別 draft structure へ変え、`completion_signal` の「3 phase を別々の draft state として保持する」を最初に閉じる。理由は model list、validation、create button state が phase 別 draft に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run translation-job-setup`
- `completion_signal`:
  - 入力確認、共通基盤、翻訳段階別 AI 設定、作成前確認、作成実行の順で読める。
  - 単語翻訳、NPC ペルソナ生成、本文翻訳ごとに AIサービス、モデル、APIキー状態、一括処理、設定状態を表示する。
  - API key 必須 provider で API key 未設定の場合だけ APIキー登録モーダルとモデル一覧更新不可状態を表示する。
  - LM Studio では API key 入力、API key 未設定警告、credential 選択を表示しない。
  - Gemini と xAI だけ batch checkbox を表示または有効化し、他 provider では選択できない。
  - model list 失敗、API key 未設定、credential 参照不能では model 手入力欄を出さず、create job を有効にしない。
  - provider または APIキー状態変更後は model 選択を未選択へ戻し、遅延 model list response を現在 phase へ混入させない。
  - 3 phase に不足がない時だけ、個別確定ボタンなしで翻訳ジョブ作成を実行できる。
  - `prototype.svelte` とサンプル値を product code、fixture、default state、test data へ移植していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `12-20 files`、`700-1300 changed lines`。
  - UI人間操作E2E の証明は `tests-provider-settings-ui-and-regression` と `final-validation-and-report` に寄せる。
  - backend core と同時に実行できるが、frozen contract 変更が必要な場合は `provider-settings-contract-freeze` の更新完了を待つ。
  - `本番経路`: frontend gateway contract -> frontend usecase / store / presenter -> Job Setup screen。

### `tests-provider-settings-api-acceptance`

- `implementation_target`: API 起点の受け入れテストを先行して固定し、master-persona 非依存、model list secret gate、phase 実行時参照、secret 非露出を検証する。
- `implementation_artifact`: `シナリオテスト実装`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_prototype`: `N/A`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `provider-settings-contract-freeze`
  - `architecture_layer_basis`: public API seam、backend fake transport、secret store seam。
  - `frozen_public_seams`: `SCN-TJSPPS-002`, `SCN-TJSPPS-003`, `SCN-TJSPPS-007`, `SCN-TJSPPS-008` に必要な public API。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential ref と credential state。
  - `secret_values_for_provider_external_api_internal_auth`: fake secret store 内の API key value。
  - `secret_resolution_owner_layer`: backend service / fake transport seam。
  - `forbidden_outputs`: test log、fake transport log、DTO、error summary、structured log。
- `owned_scope`:
  - `internal/integrationtest/*translation_job_setup*`
  - `internal/service/*translation_job_setup*_test.go`
  - `internal/infra/ai/*_test.go`
  - 必要最小限の test helper / fixture
- `depends_on`: `provider-settings-contract-freeze`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-provider-settings-core`, `frontend-provider-settings-ui`
- `parallel_blockers`: `なし`
- `first_action`: `internal/integrationtest` に master-persona provider 設定が存在しても phase runtime missing で create できない failing acceptance test を追加し、`completion_signal` の「master-persona provider 設定を参照しない」を最初に閉じる。理由は境界違反を API 起点で固定するため。
- `validation_commands`:
  - `go test ./internal/integrationtest ./internal/service ./internal/infra/ai -run 'TJSPPS|TranslationJobSetup|ProviderSettings|ModelList|Redaction'`
- `completion_signal`:
  - `SCN-TJSPPS-002` の master-persona 非依存を API 起点で検証する。
  - `SCN-TJSPPS-003` の provider model list secret gate を fake secret store と fake transport で検証する。
  - `SCN-TJSPPS-007` の phase 実行時に対象 phase 専用設定だけを読む条件を検証する。
  - `SCN-TJSPPS-008` の secret 非露出と paid API 非依存を検証する。
  - fake provider を user-facing provider list へ追加せず、external request / SDK transport だけを fake に差し替える。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `5-10 files`、`250-650 changed lines`。
  - product code は変更しない。必要な helper は test helper に限定する。
  - backend core と並列可能だが、shared helper の重複が出た場合は `owned_scope_overlap` として implement_lane が順序調整する。

### `integration-provider-settings-boundary`

- `implementation_target`: backend contract、Wails controller、frontend gateway / adapter を接続し、phase 別 provider 設定を実画面から backend へ通す。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_prototype`: `./prototype.svelte`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `provider-settings-contract-freeze`
  - `architecture_layer_basis`: Wails controller、generated DTO mapping、frontend gateway adapter。
  - `frozen_public_seams`: phase 別 options、model list、validation、create、summary の Wails binding。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential ref、credential state、redacted error kind。
  - `secret_values_for_provider_external_api_internal_auth`: backend secret store が解決する API key。
  - `secret_resolution_owner_layer`: backend service。Wails DTO は secret 本体を運ばない。
  - `forbidden_outputs`: Wails DTO、frontend gateway log、browser console、request capture。
- `owned_scope`:
  - `internal/controller/wails/translation_job_setup_controller.go`
  - `internal/controller/wails/translation_job_setup_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/controller/wails/gateway-dto/translation-job-setup/*`
  - `frontend/src/controller/wails/translation-job-setup.gateway*`
  - 必要な generated binding 更新
- `depends_on`: `backend-provider-settings-core`, `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `internal/controller/wails/translation_job_setup_controller.go` に `ListTranslationJobSetupProviderModels` の controller method と redacted DTO mapping を接続し、`completion_signal` の「model list API が Wails public seam から呼べる」を最初に閉じる。理由は frontend のモデル一覧更新がこの seam に依存するため。
- `validation_commands`:
  - `go test ./internal/controller/wails -run 'TranslationJobSetup|ProviderSettings|ModelList'`
  - `npm --prefix frontend run test -- --run translation-job-setup`
- `completion_signal`:
  - frontend gateway から options、provider model list、validation、create、summary を同じ DTO shape で呼べる。
  - provider model list の request token と phase id が backend response と frontend state で対応する。
  - credential missing または LM Studio not_required が Wails DTO で secret 非露出のまま表現できる。
  - create response と summary response が phase 別 provider、model、credential state、execution mode、batch mode を返す。
  - gateway / controller の error mapping は user-facing failure kind と debug 用 redacted summary だけを返す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`350-800 changed lines`。
  - backend 永続化と frontend UI の代替実装は含めない。
  - `本番経路`: Wails controller / DTO -> frontend gateway adapter。

### `tests-provider-settings-ui-and-regression`

- `implementation_target`: UI 状態、frontend presenter、Wails gateway 接続後の回帰テストを実装する。
- `implementation_artifact`: `シナリオテスト実装`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_prototype`: `./prototype.svelte`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `provider-settings-contract-freeze`
  - `architecture_layer_basis`: frontend presenter / store / UI screen と Wails gateway fake。
  - `frozen_public_seams`: frontend gateway contract、UI-visible state、model list response。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: APIキー状態、credential ref、not_required、redacted failure。
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: backend service。frontend tests は secret 本体を fixture に入れない。
  - `forbidden_outputs`: frontend fixture、snapshot、console、error summary、DOM text。
- `owned_scope`:
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`
  - `frontend/src/application/presenter/translation-job-setup/*test.ts`
  - `frontend/src/application/usecase/translation-job-setup/*test.ts`
  - `frontend/src/controller/wails/translation-job-setup.gateway*test.ts`
  - 必要最小限の frontend test helper
- `depends_on`: `integration-provider-settings-boundary`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts` に LM Studio 選択時の API key UI 非表示テストを追加し、`completion_signal` の「LM Studio に API key 入力、未設定警告、credential 選択が出ない」を最初に閉じる。理由は UI 仕様の provider 別差分が回帰しやすいため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run translation-job-setup`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - `SCN-TJSPPS-001` の不足なし作成導線を frontend test または mocked gateway test で検証する。
  - `SCN-TJSPPS-004` の LM Studio credential 不要 UI を検証する。
  - `SCN-TJSPPS-005` の Gemini / xAI だけ batch checkbox を検証する。
  - `SCN-TJSPPS-006` の provider / APIキー状態変更後の model 未選択化と遅延 response 混入防止を検証する。
  - UI test fixture、DOM、console に API key 平文や secret 本体が出ない。
  - `prototype.svelte` のサンプル値を product test data へ移植していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `4-8 files`、`250-650 changed lines`。
  - UI人間操作E2E の実ブラウザ証跡は `final-validation-and-report` で扱う。
  - product code は変更しない。test helper 以外の fixtures へ prototype sample を移植しない。

### `final-validation-and-report`

- `implementation_target`: 全 handoff 完了後に scenario gate、backend / frontend local harness、UI 証跡、5 観点 review 入力、completion packet を確認する。
- `implementation_artifact`: `最終検証`
- `implementation_skill`: `implement-lane`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_prototype`: `./prototype.svelte`
  - `ui_mock_data`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#Agent Browser Review`
- `contract_freeze`:
  - `status`: `not_required`
  - `freeze_source`: `N/A`
  - `architecture_layer_basis`: `N/A`
  - `frozen_public_seams`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: redacted completion evidence。
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: completed implementation evidence。
  - `forbidden_outputs`: work report、review input、UI evidence、test log summary。
- `owned_scope`:
  - `work_history/runs/YYYY-MM-DD-translation-job-setup-phase-provider-settings-run/codex.md`
  - task-local residual note が必要な場合だけ追加
- `depends_on`: `backend-provider-settings-core`, `frontend-provider-settings-ui`, `integration-provider-settings-boundary`, `tests-provider-settings-api-acceptance`, `tests-provider-settings-ui-and-regression`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `work_history/runs/YYYY-MM-DD-translation-job-setup-phase-provider-settings-run/codex.md` の final validation 欄を作成し、`completion_signal` の「検証結果と blocked reason の記録場所が固定される」を最初に閉じる。理由は closeout と work_reporter が source_ref 付き証跡を読むため。
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - scenario gate が pass する。
  - backend local と frontend local が pass する。
  - UI 証跡で desktop と mobile の入力確認、共通基盤、3 phase 設定、下部固定 CTA、作成後 summary が重ならない。
  - LM Studio に API key 入力、API key 未設定警告、credential 選択が出ない。
  - Gemini と xAI だけ batch checkbox を使える。
  - model list 失敗または API key 未設定時に model 手入力欄が存在しない。
  - paid real API が呼ばれていない証跡を fake transport log または test result で確認できる。
  - system / harness が環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason、再実行環境、再実行コマンドを report に残す。
  - completion packet が `completed_handoffs`、`touched_files`、`implemented_scope`、`test_results`、`ui_evidence`、`final_validation_result`、`codex_review_result`、`completion_evidence` を含む。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `notes`:
  - product 実装はここで追加しない。
  - broad validation owner は全 handoff 完了後だけに置く。
  - Sonar を使う場合は repo-local issue gate とし、Sonar サーバ側 Quality Gate と混同しない。

## Human Codex Handoff Packet

- `entry`: `.codex/skills/implement-lane/SKILL.md`
- `task_id`: `translation-job-setup-phase-provider-settings`
- `scope_source`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/implementation-scope.md`
- `start_wave`: `wave-1`
- `do_not_change`:
  - `docs/`
  - `.codex/`
  - `.codex/skills`
  - `.codex/agents`
- `do_not_implement`:
  - master-persona provider 設定への依存復活
  - master-persona secret namespace の Job Setup secret 解決元化
  - 手動 model 入力欄
  - LM Studio の API key 入力、API key 未設定 warning、credential select
  - Gemini と xAI 以外の batch mode
  - provider list への fake provider 追加
  - phase 実行 UI の再選択導線
  - docs 正本化
  - task-local UI prototype や sample data の product code / fixture 移植
- `must_return`:
  - `completed_handoffs`
  - `touched_files`
  - `implemented_scope`
  - `test_results`
  - `implementation_investigation`
  - `ui_evidence`
  - `final_validation_result`
  - `codex_review_result`
  - `coverage_gate_result`
  - `sonar_gate_result`
  - `harness_gate_result`
  - `completion_evidence`
  - `telemetry_events`

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: Codex 側 `work_reporter` が読む実装事実。report 文面ではなく、completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`

## Open Items

- none
