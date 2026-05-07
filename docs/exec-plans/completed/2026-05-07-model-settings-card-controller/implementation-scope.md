# Implementation Scope: 2026-05-07-model-settings-card-controller

- `skill`: implementation-scope
- `status`: ready-for-implement-lane
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: 2026-05-07 に `Q-MSCC-001` から `Q-MSCC-004` まで人間回答済み。回答は `scenario-design.md` と `ui-design.md` へ反映済み。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `ui_agent_browser_review`: `./ui-design.md#Agent-Browser-Review`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `required_reading_basis`:
  - `docs/architecture.md`
  - `docs/detail-specs/translation-job-setup.md`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `frontend/src/ui/components/AIModelSelectionCard.svelte`
  - `frontend/src/application/gateway-contract/provider-settings/provider-settings-screen-model.ts`
  - `internal/usecase/provider_settings_contract.go`
  - `internal/usecase/translation_job_setup_contract.go`

## Fixed Decisions

- `needs_human_decision`: `0`
- 未解決 conflict は 0 件である。
- マスターペルソナと翻訳ジョブ設定は、同じモデル設定カード制御を使う。
- provider と model は参照側ごとに保存し、マスターペルソナと翻訳ジョブ設定は相互に混入しない。
- AIサービス設定は endpoint と credential 参照状態だけを保存し、model と処理方法の保存元にしない。
- 空の model list 成功は取得済み 0 件として表示し、model 未選択として保存と job 作成を拒否する。
- 保存失敗後は未保存変更として残し、保存済み設定として表示または利用しない。
- APIキー未設定時は、共有カードに AIサービス設定を開く導線を出さず、状態表示だけにする。
- fake mode は backend または adapter 境界で扱い、frontend に fake mode 判定や `fake-model` 固有分岐を置かない。
- APIキー本体、secret、raw request、raw response、raw prompt、内部ログ用識別子は UI、DTO、要約、log に出さない。

## Existing Boundary Evidence

- `AIModelSelectionCard.svelte` は provider、model、状態、操作を props と event で扱う表示部品である。
- `provider_settings_contract.go` は provider 3 種、credential 状態、model list 状態、redacted failure を持つ。
- `translation_job_setup_contract.go` は 3 翻訳段階の phase runtime draft と model list status を持つ。
- `architecture.md` は `UI Component` が Store、Gateway、backend DTO、generated binding を直接扱わないと定める。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-shared-model-card-controller` | `なし` | `なし` | `なし` |
| `wave-2` | `backend-reference-model-settings-core` | `frontend-shared-model-card-controller` | `なし` | `backend_frontend_order` |
| `wave-3` | `integration-model-settings-wails-gateway` | `frontend-shared-model-card-controller`, `backend-reference-model-settings-core` | `なし` | `shared_contract_change` |
| `wave-4` | `tests-model-settings-scenario`, `tests-model-settings-unit` | `integration-model-settings-wails-gateway` | `tests-model-settings-scenario <-> tests-model-settings-unit` | `なし` |
| `wave-5` | `final-validation-and-review-input` | `tests-model-settings-scenario`, `tests-model-settings-unit` | `なし` | `broad_gate_shared` |

## Handoffs

### `frontend-shared-model-card-controller`

- `implementation_target`: 承認済み UI 契約に従い、マスターペルソナと Job Setup が使う共有モデル設定カード制御を frontend 側で作る。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#Agent-Browser-Review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、参照側 ID、credential 状態、credential 参照 ID、model list 状態、model 件数、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: backend service / infra AI provider。frontend は secret 本体を扱わない。
  - `forbidden_outputs`: APIキー本文、復号可能値、provider raw request、provider raw response、raw prompt、内部 request 識別子。
- `owned_scope`:
  - `frontend/src/application/gateway-contract/*` のモデル設定カード用 contract または既存 provider settings / Job Setup contract の参照側 model list 拡張。
  - `frontend/src/application/store/*` の共有モデル設定状態。
  - `frontend/src/application/usecase/*` の provider 変更、model list 更新、model 選択、保存状態制御。
  - `frontend/src/application/presenter/*` の共有カード view model。
  - `frontend/src/controller/master-persona/*` と `frontend/src/controller/translation-job-setup/*` の screen controller 接続。
  - `frontend/src/ui/components/AIModelSelectionCard.svelte` と両画面の利用箇所。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/store/` 配下に共有モデル設定状態の最小 state 型を追加し、`completion_signal` の「参照側ごとの provider / model / model list / 保存状態を 1 つの状態規則で保持する」を最初に閉じる。理由は両画面の controller と presenter が同じ状態規則に依存するためである。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run 'master-persona|translation-job-setup|AIModelSelectionCard|model'`
- `completion_signal`:
  - マスターペルソナと Job Setup は同じ frontend 状態規則で provider、model、model list、保存状態を扱う。
  - provider 変更後は旧 provider の model list と model を現在 provider の保存済み状態へ混入しない。
  - 空の model list 成功は取得済み 0 件として表示される。
  - 保存失敗後は未保存変更として残り、再試行できる。
  - APIキー未設定時は共有カード内に AIサービス設定導線を出さず、更新不可状態だけを表示する。
  - fake mode 判定、`fake` provider ID、`fake-model` 固有分岐を frontend に追加していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `10-15 files`、`500-800 changed lines`。
  - backend 実装と同時に開始しない。UI がある task のため frontend を先行 wave に置く。
  - backend public seam の実接続は `integration-model-settings-wails-gateway` に残す。

### `backend-reference-model-settings-core`

- `implementation_target`: 参照側ごとの provider / model 保存、再取得、model list 取得、secret 非露出、fake transport 境界を backend core で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、参照側 ID、credential 参照 ID、credential configured / missing / not_required、endpoint summary、model list status、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: APIキー平文、復号可能値、provider SDK token、provider request authorization。
  - `secret_resolution_owner_layer`: `internal/service` から `internal/infra/ai` と secret store へ渡す境界。
  - `forbidden_outputs`: DB row、DTO、UI、structured log、fake transport log、error summary、URL、保存要約、request capture。
- `owned_scope`:
  - `internal/usecase/provider_settings_contract.go` または参照側 model 設定用 contract。
  - `internal/service/provider_settings_service.go` と provider settings consumer 境界。
  - `internal/usecase/master_persona_*` と `internal/service/master_persona_*` の参照側保存取得。
  - `internal/usecase/translation_job_setup_*` と `internal/service/translation_job_setup_*` の phase 別保存取得。
  - `internal/repository/*` の参照側 provider / model 永続化が必要な範囲。
  - `internal/infra/ai/*` の model list fake transport / provider adapter 境界。
- `depends_on`: `frontend-shared-model-card-controller`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `internal/usecase/provider_settings_contract.go` または参照側 model 設定 contract に、参照側 ID、provider、model、model list status、redacted failure を持つ read / save / list models の最小 DTO を追加し、`completion_signal` の「参照側ごとの保存取得 contract が secret 本体を含まない」を最初に閉じる。理由は repository、service、controller が同じ field obligation に依存するためである。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/service ./internal/repository ./internal/infra/ai -run 'ProviderSettings|Model|MasterPersona|TranslationJobSetup|Fake'`
- `completion_signal`:
  - provider と model は参照側ごとに保存され、マスターペルソナと Job Setup へ相互混入しない。
  - AIサービス設定は model 保存元にならず、endpoint と credential 参照状態だけを提供する。
  - APIキー必須 provider は credential 参照が解決できる場合だけ model list 取得へ進む。
  - LM Studio は APIキー不要 provider として扱い、APIキー不足に分類しない。
  - 空の model list 成功は成功応答かつ 0 件として扱われる。
  - 保存失敗は保存済み状態を更新せず、redacted failure として返る。
  - paid real API は local tests で呼ばず、provider list は real provider のまま transport / SDK seam だけ fake に差し替える。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `14-22 files`、`750-1300 changed lines`。
  - backend 永続化、provider settings consumer、model list fake transport は同じ secret 不変条件を共有するため 1 handoff にまとめる。
  - frontend UI と Wails gateway 接続は含めない。

### `integration-model-settings-wails-gateway`

- `implementation_target`: backend core と frontend shared controller を Wails / DTO / gateway / bootstrap 境界で接続する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#Agent-Browser-Review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、参照側 ID、credential 状態、credential 参照 ID、model list 状態、model 件数、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: APIキー平文、provider authorization。
  - `secret_resolution_owner_layer`: backend service / infra AI provider。Wails response と frontend gateway response は secret 本体を受け取らない。
  - `forbidden_outputs`: Wails DTO、frontend gateway DTO、UI、console、structured log、error summary、fake transport log。
- `owned_scope`:
  - `internal/controller/wails/*` のモデル設定カード用 bind または既存 controller の接続。
  - `internal/bootstrap/*` の手動 DI 接続。
  - `frontend/src/controller/wails/*` と `frontend/src/controller/wails/gateway-dto/*` の gateway 接続。
  - `frontend/wailsjs/` の生成物が必要な場合は生成結果だけを扱う。
  - frontend screen controller factory の production gateway wiring。
- `depends_on`: `frontend-shared-model-card-controller`, `backend-reference-model-settings-core`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `internal/controller/wails/` の公開 bind と `frontend/src/controller/wails/` の gateway DTO の field 対応を 1 つ固定し、`completion_signal` の「Wails DTO と frontend gateway DTO が同じ redaction 境界を持つ」を最初に閉じる。理由は backend と frontend の接続不一致を早く検出するためである。
- `validation_commands`:
  - `go test ./internal/controller/wails ./internal/bootstrap -run 'ProviderSettings|Model|MasterPersona|TranslationJobSetup'`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run 'gateway|master-persona|translation-job-setup|model'`
- `completion_signal`:
  - Wails bind、backend controller DTO、frontend gateway DTO、frontend gateway contract が provider / model / credential state / model list status を同じ意味で扱う。
  - generated binding を手編集していない。
  - APIキー本体と raw payload を Wails DTO、frontend gateway DTO、UI、console へ出していない。
  - fake mode で通常 provider ID のまま `fake-model` が取得結果として伝わる。
  - 遅延した model list 応答は現在 provider と現在要求へ反映されない。
  - 実装後に `agent-browser` で Job Setup とマスターペルソナの共有カード表示を確認する材料が揃っている。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-12 files`、`350-700 changed lines`。
  - backend 実装と frontend UI 実装の代替にしない。接続と redaction 境界だけを扱う。

### `tests-model-settings-scenario`

- `implementation_target`: 承認済みシナリオを APIテストと UI人間操作E2E の証明へ落とす。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#Agent-Browser-Review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential 状態、credential 参照 ID、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: fake secret store 内の APIキー値。
  - `secret_resolution_owner_layer`: backend fake secret store / fake provider transport。
  - `forbidden_outputs`: test log、fake transport log、DTO、UI、error summary。
- `owned_scope`:
  - `internal/apitest/*` と `internal/integrationtest/*` のモデル設定カード受け入れテスト。
  - `frontend/src/ui/screens/master-persona/*test.ts` または関連 UI scenario test。
  - `frontend/src/ui/screens/translation-job-setup/*test.ts` の共有カード状態回帰。
  - 必要最小限の fake gateway / fixture。
- `depends_on`: `integration-model-settings-wails-gateway`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `tests-model-settings-unit`
- `parallel_blockers`: `なし`
- `first_action`: `internal/apitest/` または `internal/integrationtest/` に SCN-MSCC-003 の fake mode model list 受け入れテストを追加し、`completion_signal` の「fake provider ID を表示または保存せず `fake-model` を取得結果として扱う」を最初に閉じる。理由は人間期待結果の中心であり、frontend fake 分岐禁止も同時に確認できるためである。
- `validation_commands`:
  - `go test ./internal/apitest ./internal/integrationtest -run 'ModelSettings|ProviderSettings|TranslationJobSetup|MasterPersona'`
  - `npm --prefix frontend run test -- --run 'JobSetupPage|MasterPersona|AIModelSelectionCard|model'`
- `completion_signal`:
  - `SCN-MSCC-001` から `SCN-MSCC-010` までの受け入れ条件が APIテストまたは UI人間操作E2E のどちらかで証明される。
  - 有料の実 AI API を呼ばない。
  - 空の model list、保存失敗、遅延応答、APIキー未設定、fake mode を fixture で再現できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`450-800 changed lines`。
  - system test 全体や Sonar は最終検証へ寄せる。

### `tests-model-settings-unit`

- `implementation_target`: frontend 状態規則、backend redaction、保存 namespace、provider 種別分岐を単体テストで固定する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、credential 状態、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: fake secret value。
  - `secret_resolution_owner_layer`: backend unit test fake。
  - `forbidden_outputs`: unit test assertion output、DTO、UI、error summary、fake transport log。
- `owned_scope`:
  - `frontend/src/application/store/*test.ts`
  - `frontend/src/application/usecase/*test.ts`
  - `frontend/src/application/presenter/*test.ts`
  - `internal/usecase/*_test.go`
  - `internal/service/*_test.go`
  - `internal/repository/*_test.go`
  - `internal/infra/ai/*_test.go`
- `depends_on`: `integration-model-settings-wails-gateway`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `tests-model-settings-scenario`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/usecase/` の共有モデル設定 usecase test で、provider 変更時に旧 model list と model を現在 provider へ混入しない clause を閉じる。理由は遅延応答破棄と保存拒否の基本不変条件になるためである。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run 'model|provider|master-persona|translation-job-setup'`
  - `go test ./internal/usecase ./internal/service ./internal/repository ./internal/infra/ai -run 'Model|ProviderSettings|MasterPersona|TranslationJobSetup|Redaction'`
- `completion_signal`:
  - provider 変更、model list 更新、model 選択、保存失敗、空一覧、遅延応答破棄が単体で検証される。
  - backend redaction と secret 非露出が単体で検証される。
  - 参照側ごとの保存 namespace が単体で検証される。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`400-750 changed lines`。
  - シナリオテストの代替ではなく、失敗箇所を狭く特定するための補助検証である。

### `final-validation-and-review-input`

- `implementation_target`: 全 handoff 完了後の最終検証、UI 証跡、レビュー入力を集める。
- `implementation_artifact`: `最終検証`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#Agent-Browser-Review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider id、model id、credential 状態、redacted failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: 既存実装済み backend 境界。
  - `forbidden_outputs`: final log、UI 証跡、console、report input。
- `owned_scope`: `tmp/agent-browser/`、`tmp/logs/`、`test-results/` の証跡出力だけ。
- `depends_on`: `tests-model-settings-scenario`, `tests-model-settings-unit`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `npm --prefix frontend run check` の最終結果を取得し、`completion_signal` の「frontend type / Svelte check が通る」を最初に閉じる。理由は UI 実装後の広域破綻を最初に確認するためである。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test`
  - `go test ./internal/...`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - frontend、backend、integration、scenario test、unit test の完了根拠が揃っている。
  - `agent-browser` でマスターペルソナと Job Setup の主要状態を確認した証跡がある。
  - system test が環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason と再実行コマンドを残す。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `notes`:
  - 広域検証だけを扱い、プロダクトコードやプロダクトテストの修正は行わない。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`
