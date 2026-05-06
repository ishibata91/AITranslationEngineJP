# Implementation Scope: translation-job-management

- `skill`: implementation-scope
- `status`: ready-for-implement-lane
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: 2026-05-06 に人間が「先進めていい」と回答した。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `spec_reference`: `docs/spec.md`
- `er_reference`: `docs/er.md`
- `code_map`: `tmp/code-map/index.json`

## Approved Implementation Scope

- Completed 以外の job を一覧し、状態、現在フェーズ、進捗、入力出自、操作可否を表示する。
- 選択した未完了 job を Job Run の表示対象にし、表示だけでは job 状態を変更しない。
- 同じ入力データから複数 job を作成できるようにし、同一入力 job 一意制約を削除する。
- Running job の削除を拒否し、停止入口、停止要求中、停止失敗、Paused 収束後の削除可否再判定を表示する。
- 非実行中 job は、DB の job 本体と配下情報を削除し、入力データ、外部入力 JSON、入力ファイルを残す。
- Paused と RecoverableFailed は再開入口を表示し、cache missing、terminal state、状態不整合では再開不可理由を表示する。
- 保存済み AI 設定要約と credential 参照状態だけを表示し、secret 本体、provider 応答原文、過剰な入力本文や翻訳本文を出さない。

## Non Goals

- Completed job の成果物確認、完了履歴画面、再出力導線は扱わない。
- 入力キャッシュ再構築の実行 UI は扱わない。
- 実際の phase 実行、翻訳生成、provider SDK 実装、paid API 呼び出しは扱わない。
- API 設定が現在も使えるかの再開直前確認は、再開実行側の後続 task で扱う。
- 外部通信の停止方式、遅延応答破棄、停止要求後の実行状態不整合防止は、翻訳実行側の後続 task で扱う。
- docs 正本、`tasks/`、`.codex/`、プロダクト外の workflow 契約は変更しない。

## Fixed Decisions

- `needs_human_decision`: `0`
- Q1 は回答済みである。同じファイルまたは同じ入力データから job を何個でも作れる。
- Q2 は回答済みである。非実行中 job 削除では job 本体と配下情報を DB から削除し、入力データと入力ファイルを残す。
- Q3 は回答済みである。Job Management は保存済み AI 設定要約、入力キャッシュ状態、再開不可理由を表示する。
- Q4 は回答済みである。Job Management は停止要求の見え方と削除可否再判定を扱い、停止制御本体は後続境界に送る。
- frontend handoff は backend handoff より先に実行する。
- frontend 実装後人間レビューは、backend handoff と integration handoff の開始条件である。
- backend、frontend、統合境界、シナリオテスト、単体テストは別 handoff にする。

## Handoff Summary

- `frontend-job-management-ui`: 未完了 job 管理 UI、画面状態、presenter、mocked gateway 境界を作る。
- `backend-job-management-core`: job read model、削除 guard、非実行中削除、同一入力 job 複数作成を backend で成立させる。
- `integration-job-management-wails`: backend public seam と frontend gateway を接続し、実画面確認を取る。
- `scenario-test-job-management`: 承認済み scenario を API テストと UI 人間操作 E2E で証明する。
- `unit-test-job-management`: 状態判定、redaction、delete guard、view model の単体テストを補強する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-job-management-ui` | `なし` | `なし` | `なし` |
| `wave-2` | `backend-job-management-core` | `frontend-job-management-ui`, `frontend 実装後人間レビュー approved` | `なし` | `backend_frontend_order` |
| `wave-3` | `integration-job-management-wails` | `frontend-job-management-ui`, `frontend 実装後人間レビュー approved`, `backend-job-management-core` | `なし` | `depends_on` |
| `wave-4` | `scenario-test-job-management`, `unit-test-job-management` | `integration-job-management-wails` | `scenario-test-job-management <-> unit-test-job-management` | `なし` |

## Handoffs

### `frontend-job-management-ui`

- `implementation_target`: 未完了 job 管理 UI を追加し、一覧、選択詳細、操作可否、再開不可理由、secret 非露出表示を mocked gateway で成立させる。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider、model、execution mode、credential 参照状態、credential configured / missing / inaccessible。
  - `secret_values_for_provider_external_api_internal_auth`: API key、token、復号済み credential、provider raw request / response。
  - `secret_resolution_owner_layer`: backend の provider settings / secret store。frontend は解決しない。
  - `forbidden_outputs`: UI、DTO、console、error summary、structured log、request capture、URL、read model への secret 本体出力。
- `owned_scope`:
  - `frontend/src/application/gateway-contract/translation-job-management/*`
  - `frontend/src/application/contract/translation-job-management/*`
  - `frontend/src/application/store/translation-job-management/*`
  - `frontend/src/application/presenter/translation-job-management/*`
  - `frontend/src/application/usecase/translation-job-management/*`
  - `frontend/src/controller/translation-job-management/*`
  - `frontend/src/ui/screens/translation-job-management/*`
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/gateway-contract/translation-job-management/translation-job-management-gateway-contract.ts` を追加し、`completion_signal` の「frontend が未完了 job list、selected job detail、operation availability を型で扱える」を最初に閉じる。理由は store、presenter、UI が同じ画面契約に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run translation-job-management`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - 翻訳管理を開き、Completed 以外の Ready、Running、Paused、RecoverableFailed、Failed、Canceled を一覧できる。
  - 一覧表示、再読込、選択だけでは job 状態が変わらない。
  - 選択詳細で job ID、入力出自、現在フェーズ、進捗、入力キャッシュ状態、保存済み AI 設定要約を表示できる。
  - 停止、再開、削除の有効条件と無効理由を UI Contract に沿って表示できる。
  - Running の削除は無効であり、停止入口と停止後再判定の説明を表示できる。
  - 非実行中削除の確認では、job 以下の DB 情報だけが削除対象で、入力データと抽出 JSON が残ることを表示できる。
  - cache missing、terminal state、state projection inconsistent、stale selection、list load failure を別状態として表示できる。
  - API key 平文、credential 値、provider 応答原文、過剰な入力本文や翻訳本文が UI、console、error summary に出ない。
  - desktop と mobile で長い file path、plugin 名、provider / model 名、reason text が overflow しない実装方針を持つ。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。
  - 1 handoff にする理由は、承認済み UI Contract の一覧、選択詳細、操作 panel が同じ画面状態で閉じるためである。
  - この handoff は mocked gateway までを扱い、generated Wails binding と backend 実装は含めない。
  - `本番経路`: frontend gateway contract -> usecase -> store / presenter -> Translation Job Management screen -> Job Run display target。
  - 完了後、人間が frontend 実装後レビューを行い、backend と integration の開始可否を判断する。

### `backend-job-management-core`

- `implementation_target`: 未完了 job 管理の backend read model、削除 guard、非実行中 job 削除、同一入力 job 複数作成を実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential_ref の存在状態、provider、model、execution mode、endpoint summary。
  - `secret_values_for_provider_external_api_internal_auth`: API key、token、復号済み credential、provider raw request / response。
  - `secret_resolution_owner_layer`: provider settings / secret store。Job Management は secret 本体を解決しない。
  - `forbidden_outputs`: controller DTO、error summary、structured log、audit、request capture、URL、read model への secret 本体出力。
- `owned_scope`:
  - `internal/usecase/translation_job_management_*`
  - `internal/service/translation_job_management_*`
  - `internal/controller/wails/translation_job_management_controller*`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/translation_source_repository.go`
  - `internal/repository/translation_source_sqlite_repository.go`
  - `internal/infra/sqlite/dbinit/migrations/*translation_job*`
  - `internal/infra/sqlite/migrations/*translation_job*`
  - `internal/bootstrap/app_controller.go`
  - `internal/controller/wails/app_controller.go`
- `depends_on`: `frontend-job-management-ui`, `frontend 実装後人間レビュー approved`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `internal/repository/job_lifecycle_repository.go` に未完了 job 一覧 projection と非実行中削除の port を追加し、`completion_signal` の「Completed 以外を read model として取得し、Running 削除を拒否できる」を最初に閉じる。理由は list、detail、delete guard、scenario fixture が同じ repository 境界に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'TranslationJobManagement|JobManagement|TranslationJob|JobLifecycle'`
- `completion_signal`:
  - Completed 以外の job を一覧取得でき、Completed は未完了一覧に出ない。
  - job 状態、現在フェーズ、進捗、入力出自、入力キャッシュ状態、AI 設定要約、credential 参照状態を secret なしで返せる。
  - 同じ `X_EDIT_EXTRACTED_DATA` に複数の `TRANSLATION_JOB` を作成できる。`idx_translation_job_x_edit` 相当の一意制約は削除または非一意 index 化される。
  - Running job の削除要求は拒否され、job、phase、input data は変更されない。
  - 非実行中 job の削除は `TRANSLATION_JOB` と job 配下情報を削除し、`X_EDIT_EXTRACTED_DATA`、抽出 JSON、入力ファイルを削除しない。
  - stale selection、list load failure、phase progress 集約不能、state projection inconsistent を区別した error kind または reason category を返せる。
  - Paused と RecoverableFailed の再開入口用 summary と、cache missing、terminal state、state projection inconsistent の再開不可理由を返せる。
  - 停止要求中、停止失敗、Paused 収束後の削除可否再判定を表示できる状態 projection を返せる。
  - 外部通信 cancel、late response 破棄、停止要求後の実行状態制御は実装しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。
  - 1 handoff にする理由は、一覧 projection、削除 guard、同一入力 job 複数作成が同じ Job Management backend use case と repository 境界に閉じるためである。
  - frontend UI と Wails gateway adapter は含めない。
  - `本番経路`: Wails controller / DTO -> usecase -> service -> repository -> SQLite。

### `integration-job-management-wails`

- `implementation_target`: backend public seam と frontend gateway / Wails adapter を接続し、承認済み UI の実画面確認を完了する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider、model、execution mode、credential 参照状態、reason category。
  - `secret_values_for_provider_external_api_internal_auth`: API key、token、復号済み credential、provider raw request / response。
  - `secret_resolution_owner_layer`: backend provider settings / secret store。
  - `forbidden_outputs`: Wails DTO、frontend gateway DTO、UI、console、error summary、request capture への secret 本体出力。
- `owned_scope`:
  - `frontend/src/controller/wails/translation-job-management.gateway*`
  - `frontend/src/controller/wails/gateway-dto/translation-job-management/*`
  - `frontend/src/bootstrap/app-screen-controller-factories.ts`
  - `frontend/src/main.ts`
  - `frontend/src/ui/App.svelte`
  - `frontend/wailsjs/go/*`
  - `frontend/wailsjs/runtime/*`
  - `internal/controller/wails/translation_job_management_controller*`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `tmp/agent-browser/*`
  - `tmp/logs/*`
- `depends_on`: `frontend-job-management-ui`, `frontend 実装後人間レビュー approved`, `backend-job-management-core`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`: `frontend/src/controller/wails/gateway-dto/translation-job-management/translation-job-management-gateway-dto.ts` を追加し、`completion_signal` の「frontend gateway DTO が backend response と同じ field 名、nullability、reason category を持つ」を最初に閉じる。理由は generated binding と frontend usecase の接続点を先に固定する必要があるため。
- `validation_commands`:
  - `go test ./internal/bootstrap ./internal/controller/wails -run 'TranslationJobManagement|JobManagement'`
  - `npm --prefix frontend run test -- --run translation-job-management`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - frontend gateway が Wails binding を通じて未完了 job list、selected job detail、delete、operation availability を呼べる。
  - generated Wails binding と hand-written gateway DTO の field 名、nullability、reason category が一致する。
  - UI から一覧表示、選択、削除確認、削除結果、再読込、再開不可理由表示を実画面で確認できる。
  - desktop と mobile の agent-browser 確認で、一覧、選択詳細、操作 panel、削除確認、長文 reason が破綻しない。
  - UI、console、error summary に API key 平文、credential 値、provider 応答原文が出ない。
  - paid API は呼ばれず、fixture または local DB だけで確認できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`300-700 changed lines`。
  - backend core と frontend UI の代替実装は含めない。
  - production 注入に必要な `frontend/src/bootstrap/app-screen-controller-factories.ts`、`frontend/src/main.ts`、`frontend/src/ui/App.svelte` は統合境界に含む。
  - 実画面確認は `npm run dev:wails:agent-browser` または repo 定義済み Wails 起動 command と `agent-browser` CLI で行う。
  - `本番経路`: Wails binding -> frontend gateway DTO -> frontend usecase -> UI、Wails controller -> backend usecase。

### `scenario-test-job-management`

- `implementation_target`: `SCN-TJM-001` から `SCN-TJM-009` を API テストと UI 人間操作 E2E で証明する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: fake credential 参照状態、redacted provider / model summary。
  - `secret_values_for_provider_external_api_internal_auth`: fake であっても secret 本体として扱う値、provider raw request / response。
  - `secret_resolution_owner_layer`: test fake secret store。
  - `forbidden_outputs`: screenshot、trace、console、test failure message、fixture dump への secret 本体出力。
- `owned_scope`:
  - `internal/integrationtest/*translation_job_management*`
  - `frontend/e2e/*translation-job-management*`
  - `test-results/*translation-job-management*`
  - `tmp/agent-browser/*translation-job-management*`
  - `tmp/logs/*translation-job-management*`
- `depends_on`: `integration-job-management-wails`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `unit-test-job-management`
- `parallel_blockers`: `なし`
- `first_action`: `internal/integrationtest/translation_job_management_scenario_test.go` に `SCN-TJM-002` の同一入力複数 job 作成 API テストを追加し、`completion_signal` の「同一入力 job 複数作成が API テストで証明される」を最初に閉じる。理由は過去の重複禁止仕様との差分が最大の回帰点であるため。
- `validation_commands`:
  - `go test ./internal/integrationtest -run 'SCN_TJM|TranslationJobManagement'`
  - `npx playwright test --grep 'SCN-TJM|translation-job-management'`
- `completion_signal`:
  - `SCN-TJM-001` は UI 人間操作 E2E で Completed 除外と未完了状態一覧を証明する。
  - `SCN-TJM-002` は API テストで同一入力複数 job 作成を証明する。
  - `SCN-TJM-003` は UI 人間操作 E2E で選択 job が Job Run の表示対象になることを証明する。
  - `SCN-TJM-004` と `SCN-TJM-006` は API テストで Running 削除拒否と非実行中削除を証明する。
  - `SCN-TJM-005` と `SCN-TJM-007` は UI 人間操作 E2E で再開入口と再開不可理由を証明する。
  - `SCN-TJM-008` は API テストで参照不能、読み込み失敗、集約不能の安全側表示を証明する。
  - `SCN-TJM-009` は UI 人間操作 E2E で secret 非露出を証明する。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `6-12 files`、`300-700 changed lines`。
  - APIテストと UI人間操作E2E は承認済み scenario の証明に限定する。
  - system / browser 環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason と再実行コマンドを残す。

### `unit-test-job-management`

- `implementation_target`: 実装済み責務の単体テストを補強し、状態判定、redaction、delete guard、view model を局所的に証明する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: redacted credential status、provider / model summary。
  - `secret_values_for_provider_external_api_internal_auth`: API key、token、復号済み credential。
  - `secret_resolution_owner_layer`: fake secret store または backend secret resolver の test double。
  - `forbidden_outputs`: test assertion message、snapshot、fixture dump、console への secret 本体出力。
- `owned_scope`:
  - `internal/repository/*translation_job_management*_test.go`
  - `internal/service/*translation_job_management*_test.go`
  - `internal/usecase/*translation_job_management*_test.go`
  - `internal/controller/wails/*translation_job_management*_test.go`
  - `frontend/src/application/store/translation-job-management/*.test.ts`
  - `frontend/src/application/presenter/translation-job-management/*.test.ts`
  - `frontend/src/application/usecase/translation-job-management/*.test.ts`
  - `frontend/src/controller/translation-job-management/*.test.ts`
  - `frontend/src/controller/wails/*translation-job-management*.test.ts`
  - `frontend/src/ui/screens/translation-job-management/*.test.ts`
- `depends_on`: `integration-job-management-wails`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `scenario-test-job-management`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.test.ts` に Running、Paused、RecoverableFailed、cache missing の operation availability test を追加し、`completion_signal` の「UI の操作可否と無効理由が状態ごとに安定する」を最初に閉じる。理由は人間レビュー済み UI の主要判断点を局所テストで固定できるため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'TranslationJobManagement|JobManagement|DeleteGuard|Redaction'`
  - `npm --prefix frontend run test -- --run translation-job-management`
- `completion_signal`:
  - Completed 除外、Running 削除拒否、非実行中削除、同一入力複数 job 作成を単体または repository test で証明する。
  - cache missing、terminal state、state projection inconsistent、stale selection、list load failure の reason category を単体 test で証明する。
  - credential 参照状態と secret 非露出の redaction を backend と frontend の境界ごとに証明する。
  - presenter は UI Contract の表示項目、操作可否、disabled reason、長文耐性用の view model を返せる。
  - controller / usecase は DTO、error kind、nullability を実装済み public seam と一致させる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`300-700 changed lines`。
  - scenario test と並列可能である。対象ファイル、検証意図、失敗時の担当が重ならないためである。

## Q4 Residual

- `residual_id`: `TJM-RES-Q4-translation-execution-stop-control`
- `status`: `deferred-task-required`
- `source`: `./scenario-design.questions.md` の `Q-TJM-004`
- `required_follow_up`: `tasks/` に翻訳実行側の後続 task を作る必要がある。
- `not_in_this_scope`:
  - 外部通信を即時 cancel するか、stop request 記録だけにするかの決定。
  - late response を破棄する条件。
  - 停止要求後に `JOB_PHASE_RUN`、翻訳結果、進捗、provider 実行 ID の不整合を防ぐ状態遷移。
  - 実際の provider SDK / transport の停止実装。
- `this_scope_boundary`:
  - Running job の削除拒否を表示する。
  - 停止入口、停止要求中、停止失敗、Paused 収束後の削除可否再判定を表示する。
  - 停止制御本体が未実装の場合は、停止不可理由または後続待ち状態を secret なしで表示する。

## Final Validation Candidates

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-management/scenario-design.md --coverage docs/exec-plans/active/translation-job-management/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-management/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`
- `go test ./internal/...`
- `npm run test:frontend`
- `npm --prefix frontend run check`
- `npx playwright test --grep 'SCN-TJM|translation-job-management'`
- `python3 scripts/harness/run.py --suite all`

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `frontend_human_review_result`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`
- `telemetry_events`
- `docs_changes: none`
