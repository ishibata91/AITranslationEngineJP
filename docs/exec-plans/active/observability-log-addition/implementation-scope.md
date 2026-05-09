# Implementation Scope: observability-log-addition

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: `./human-design-review.md`
- `approved_at`: 2026-05-09
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `requirement_gate`: `./scenario-design.requirement-gate.md`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `design_diff`: `./design-diff.md`
- `design_diff_component`: `./design-diff.component.puml`
- `design_diff_sequence`: `./design-diff.sequence.puml`
- `observability_spec`: `docs/observability-logging.md`
- `backend_guideline`: `docs/coding-guidelines-backend.md`
- `frontend_guideline`: `docs/coding-guidelines-frontend.md`
- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`

## Fixed Decisions

- 人間レビューは approved である。
- `needs_human_decision`: `0`
- UI 表示、画面文言、layout、style は変更しない。
- backend log は `slog` の JSON log として `stderr` へ出す。
- frontend log は `pino` の browser console 出力に限定する。
- backend log と frontend log を同じ file へ集約しない。
- frontend log を Wails 経由で backend へ送らない。
- trace ID、全 command の start / finish log、DTO 全体 dump は追加しない。
- docs 正本化は Codex implementation レーンへ渡さない。

## Logger Foundation Decision

- backend logger 基盤追加: 不要。
  根拠: `main.go` は `InstallDiagnosticLogger(os.Stderr, "backend", slog.LevelInfo)` を呼び、`internal/infra/runtime/diagnostic_log.go` は `slog` JSON handler を既に提供している。
- frontend logger 基盤追加: 不要。
  根拠: `frontend/package.json` は `pino` を持ち、`frontend/src/application/diagnostic/frontend-diagnostic-logger.ts` は `pino` browser logger を既に提供している。
- frontend runtime event wiring: 必要。
  理由: `frontend/src/bootstrap/app-screen-controller-factories.ts` は `diagnosticLogger` を作るが、`MasterDictionaryRuntimeEventAdapter` へ渡していない。
- 統合境界実装: 不要。
  理由: API、Wails DTO、gateway、adapter 契約を変更しない。frontend log は backend へ送らない。

## Handoff Classification

| 成果物 | 判断 | 理由 |
| --- | --- | --- |
| backend 実装 | 作成しない | backend の `slog` 基盤は既存であり、機能実装ではなく観測ログ追加で閉じる。 |
| frontend 実装 | 作成する | frontend runtime event の `pino` log を使うには既存 logger を runtime adapter へ配線する必要がある。 |
| 統合境界実装 | 作成しない | public API、DTO、Wails gateway を変更しない。 |
| 観測ログ追加 | 作成する | backend 状態遷移、provider / secret、file / DB / output、大量処理の log を slice ごとに追加する。 |
| 単体テスト | 作成する | backend payload と frontend runtime event 分岐を局所テストで証明する。 |
| シナリオテスト | 作成する | backend command 境界から log 分類を APIテストで証明する。 |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `OBS-FE-001`, `OBS-BE-001` | なし | `OBS-FE-001 <-> OBS-BE-001` | なし |
| `wave-2` | `OBS-BE-002`, `OBS-BE-003`, `OBS-UNIT-FE-001` | `OBS-FE-001`, `OBS-BE-001` | `OBS-BE-002 <-> OBS-BE-003`, `OBS-BE-002 <-> OBS-UNIT-FE-001`, `OBS-BE-003 <-> OBS-UNIT-FE-001` | なし |
| `wave-3` | `OBS-BE-004` | `OBS-BE-002`, `OBS-BE-003` | なし | `owned_scope_overlap` |
| `wave-4` | `OBS-UNIT-BE-001`, `OBS-SCN-BE-001` | `OBS-BE-001`, `OBS-BE-002`, `OBS-BE-003`, `OBS-BE-004` | `OBS-UNIT-BE-001 <-> OBS-SCN-BE-001` | なし |

## Handoffs

### `OBS-FE-001`: frontend runtime event の `pino` log 配線

- `implementation_target`: frontend runtime event log
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `agent`: frontend_implementer
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `OBS-BE-001`
- `parallel_blockers`: なし
- `approved_scope`: 既存 `pino` diagnostic logger を `MasterDictionaryRuntimeEventAdapter` へ注入し、subscribe、accepted、dropped、skipped、detached を browser console log へ出す。UI 表示、文言、layout、style は変更しない。
- `expected_files`: 3 files
  - `frontend/src/bootstrap/app-screen-controller-factories.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts`
- `expected_changed_lines`: 110
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_no_change_source`: `./scenario-design.md#UI 設計`, `./human-design-review.md#承認内容`
- `first_action`: `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts` の `MasterDictionaryRuntimeEventAdapter` constructor に任意の diagnostic logger dependency を追加し、no-op 既定値で既存挙動を維持する。対応する完了条件は「runtime adapter が既存 pino logger を受け取れる」である。最初に配線境界を閉じるため、この 1 手目にする。
- `completion_signal`:
  - runtime event subscribe 成功と subscribe 不可を分類して log に出す。
  - progress event と completed event の受信を `accepted` として分類する。
  - payload parse 失敗、stale 相当、detach を `dropped`、`skipped`、`detached` として分類する。
  - log payload は `event`、`where`、`result`、必要な `id`、`reason` だけにする。
  - frontend log を backend へ送らない。
  - UI 表示、画面文言、layout、style を変更しない。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: not_required

### `OBS-BE-001`: backend 状態遷移 log

- `implementation_target`: backend state transition log
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `agent`: observability_implementer
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `OBS-FE-001`
- `parallel_blockers`: なし
- `approved_scope`: job state と phase run state の開始、停止、再開、retry、削除、後続 phase readiness の境界に、変更前、変更後、拒否理由を残す backend JSON log を追加する。
- `expected_files`: 5 files
  - `internal/service/translation_job_management_service.go`
  - `internal/usecase/translation_job_management_usecase.go`
  - `internal/usecase/term_translation_phase_usecase.go`
  - `internal/usecase/persona_generation_phase_usecase.go`
  - `internal/usecase/body_translation_phase_usecase.go`
- `expected_changed_lines`: 230
- `first_action`: `internal/service/translation_job_management_service.go` の `DeleteJob` 経路に、削除許可または削除拒否を分類する `slog.InfoContext` / `slog.WarnContext` を 1 箇所追加する。対応する完了条件は「削除の許可と拒否が before / after / reason で確認できる」である。状態遷移 log の代表 pattern を先に固定するため、この 1 手目にする。
- `completion_signal`:
  - 許可された状態変更は変更前状態、変更後状態、代表 `id` を含む。
  - 拒否された状態変更は状態を変えず、`reason` を含む。
  - `event`、`where`、`result` は全 log に入る。
  - terminal state、active phase 既存、前段 phase 未完了、再開不可、削除拒否を区別できる。
  - trace ID と全 command start / finish log を追加しない。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: not_required

### `OBS-BE-002`: provider、secret、credential 再解決 log

- `implementation_target`: backend provider and secret boundary log
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `agent`: observability_implementer
- `depends_on`: `OBS-BE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `OBS-BE-003`, `OBS-UNIT-FE-001`
- `parallel_blockers`: なし
- `approved_scope`: provider 設定、credential 再解決、secret store、model list、phase 実行時 provider 失敗を、credential 未設定、secret store 失敗、provider timeout、不正応答、correlation error、provider skipped に分類する。
- `expected_files`: 5 files
  - `internal/service/provider_settings_service.go`
  - `internal/service/translation_job_setup_service.go`
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
- `expected_changed_lines`: 300
- `first_action`: `internal/service/provider_settings_service.go` の credential 解決または validation 失敗経路に、secret 本体を含めない失敗分類 log を 1 箇所追加する。対応する完了条件は「credential 未設定と secret store 失敗を別 reason で確認できる」である。secret 境界の安全な payload を先に固定するため、この 1 手目にする。
- `completion_signal`:
  - credential 未設定、secret store 失敗、provider timeout、不正応答、correlation error、provider skipped を区別する。
  - provider 種別は原因分離に必要な分類値だけにする。
  - credential 参照実値、secret 本体、API key、endpoint 実値、provider raw payload を出さない。
  - fake provider と fake secret store で再現できる。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider ID、credential 状態分類、credential 参照の存在有無、model list 状態分類。
  - `secret_values_for_provider_external_api_internal_auth`: API key、secret store から読み出した credential 本体、provider endpoint 実値。
  - `secret_resolution_owner_layer`: repository secret store と provider settings service の解決境界。
  - `forbidden_outputs`: log、error summary、audit、request capture、URL、DTO、UI、read model へ secret 本体、API key、endpoint 実値、credential 参照実値を出さない。

### `OBS-BE-003`: file、DB、Wails Bind、成果物出力 log

- `implementation_target`: backend file DB Wails output boundary log
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `agent`: observability_implementer
- `depends_on`: `OBS-BE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `OBS-BE-002`, `OBS-UNIT-FE-001`
- `parallel_blockers`: なし
- `approved_scope`: 入力登録、cache rebuild、job 作成、DB 保存、Wails request / response 変換、xTranslator 互換 XML 出力の失敗 stage を分類する backend JSON log を追加する。
- `expected_files`: 7 files
  - `internal/controller/wails/translation_input_controller.go`
  - `internal/controller/wails/translation_output_artifact_controller.go`
  - `internal/service/translation_input_import_service.go`
  - `internal/service/translation_job_setup_service.go`
  - `internal/service/translation_output_artifact_service.go`
  - `internal/service/translation_output_artifact_xml_adapter.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
- `expected_changed_lines`: 380
- `first_action`: `internal/controller/wails/translation_input_controller.go` の request invalid 経路に、DTO 全体と full path を含めない `request_invalid` 分類 log を 1 箇所追加する。対応する完了条件は「Wails request 変換失敗が provider 失敗や UI 表示失敗へ混ざらない」である。公開接点境界の安全 payload を先に固定するため、この 1 手目にする。
- `completion_signal`:
  - invalid JSON、source file missing、cache missing、DB save failed、transaction failed、request invalid、response mapping failed、file write failed を区別する。
  - file path は full path ではなく stage と分類に留める。
  - DTO 全体、XML 全文、翻訳本文全文を出さない。
  - retry 可否分類が分かる場合だけ `reason` に含める。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: not_required

### `OBS-BE-004`: 大量処理の集約 log

- `implementation_target`: backend bulk aggregate log
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `agent`: observability_implementer
- `depends_on`: `OBS-BE-002`, `OBS-BE-003`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `approved_scope`: 入力取り込み、辞書反映、provider 実行、保護要素検証、成果物 row validation の大量処理へ、件数、失敗分類、最初の失敗分類、最後の失敗分類だけを集約 log として追加する。
- `expected_files`: 6 files
  - `internal/service/translation_input_import_service.go`
  - `internal/service/master_dictionary_import_service.go`
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/xtranslator_output_row_builder.go`
- `expected_changed_lines`: 320
- `first_action`: `internal/service/translation_input_import_service.go` の import summary 作成後に、input count、skipped count、failed count を 1 回だけ出す集約 log を追加する。対応する完了条件は「loop 内 1 件ごとの log を出さず、件数と失敗分類を確認できる」である。大量処理の出力頻度を先に固定するため、この 1 手目にする。
- `completion_signal`:
  - loop 内で 1 件ごとの log を出さない。
  - input count、output count、skipped count、failed count を集約する。
  - 最初の失敗分類と最後の失敗分類だけを必要時に残す。
  - provider raw payload、prompt 全文、翻訳本文全文、XML 全文を出さない。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: not_required

### `OBS-UNIT-FE-001`: frontend runtime event log の単体テスト

- `implementation_target`: frontend runtime event log tests
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `agent`: implementation_unit_tester
- `depends_on`: `OBS-FE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `OBS-BE-002`, `OBS-BE-003`
- `parallel_blockers`: なし
- `approved_scope`: runtime event adapter の subscribed、accepted、dropped、skipped、detached と、backend へ送らない境界を単体テストで証明する。
- `expected_files`: 2 files
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/application/diagnostic/frontend-diagnostic-logger.test.ts`
- `expected_changed_lines`: 180
- `first_action`: `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts` に、payload parse 失敗が `dropped` log になり store を更新しない test を 1 件追加する。対応する完了条件は「破棄理由を分類し、画面状態を壊さない」である。SCN-OBSLOG-004 の失敗分岐を最初に閉じるため、この 1 手目にする。
- `completion_signal`:
  - subscribe 成功、subscribe 不可、payload parse 失敗、detach の各分岐を 1 テスト 1 分岐で証明する。
  - `console` 出力 payload に `event`、`where`、`result`、必要な `reason` がある。
  - frontend log を backend へ送る呼び出しを追加していない。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: not_required

### `OBS-UNIT-BE-001`: backend log payload と禁止 payload の単体テスト

- `implementation_target`: backend log payload tests
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `agent`: implementation_unit_tester
- `depends_on`: `OBS-BE-001`, `OBS-BE-002`, `OBS-BE-003`, `OBS-BE-004`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `OBS-SCN-BE-001`
- `parallel_blockers`: なし
- `approved_scope`: backend log が `event`、`where`、`result`、必要な `id`、`count`、`reason` を持ち、secret、API key、endpoint 実値、raw payload、全文本文、full path、trace ID を含まないことを局所テストで証明する。
- `expected_files`: 6 files
  - `internal/infra/runtime/diagnostic_log_test.go`
  - `internal/service/provider_settings_service_test.go`
  - `internal/service/translation_job_management_service_test.go`
  - `internal/service/translation_input_import_service_test.go`
  - `internal/service/body_translation_phase_service_test.go`
  - `internal/service/translation_output_artifact_service_test.go`
- `expected_changed_lines`: 420
- `first_action`: `internal/infra/runtime/diagnostic_log_test.go` に、代表 log payload から undefined 相当の値や禁止 payload が出ないことを確認する test を 1 件追加する。対応する完了条件は「禁止 log が出ないことを局所的に検査できる」である。共通安全条件を先に閉じるため、この 1 手目にする。
- `completion_signal`:
  - 状態遷移、provider / secret、file / DB / output、大量処理の代表 payload をテストする。
  - secret、API key、endpoint 実値、raw request、raw response、prompt 全文、翻訳本文全文、XML 全文、DTO 全体、full path、trace ID を含まないことを検査する。
  - 網羅率検証の結果または未実行理由を返す。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite coverage`
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider ID、credential 状態分類、代表 ID、集約 count。
  - `secret_values_for_provider_external_api_internal_auth`: API key、secret store credential、endpoint 実値、provider raw payload。
  - `secret_resolution_owner_layer`: repository secret store、provider settings service、provider adapter。
  - `forbidden_outputs`: log、error summary、audit、request capture、URL、DTO、UI、read model へ secret 本体と raw payload を出さない。

### `OBS-SCN-BE-001`: backend command 境界の観測 log APIテスト

- `implementation_target`: backend observability API scenario tests
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `agent`: implementation_scenario_tester
- `depends_on`: `OBS-BE-001`, `OBS-BE-002`, `OBS-BE-003`, `OBS-BE-004`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `OBS-UNIT-BE-001`
- `parallel_blockers`: なし
- `approved_scope`: Wails Bind 相当の backend command 境界から、SCN-OBSLOG-001、SCN-OBSLOG-002、SCN-OBSLOG-003、SCN-OBSLOG-005 の代表経路を APIテストで証明する。
- `expected_files`: 2 files
  - `internal/apitest/observability_log_scenario_test.go`
  - `internal/apitest/observability_log_test_helpers_test.go`
- `expected_changed_lines`: 300
- `first_action`: `internal/apitest/observability_log_scenario_test.go` を追加し、削除拒否の command 境界から `state transition` log の `result` と `reason` を検査する test を 1 件作る。対応する完了条件は「公開接点起点で状態遷移の拒否理由を確認できる」である。APIテストの capture pattern を先に固定するため、この 1 手目にする。
- `completion_signal`:
  - 状態変更の許可、拒否、失敗を APIテストで区別する。
  - provider / secret 失敗と file / DB / output 失敗を APIテストで区別する。
  - 大量処理は件数と失敗分類の集約 log だけを検査する。
  - 実 AI API を使わず fake provider、fake secret store、temp file、SQLite test DB で検証する。
  - UI人間操作E2E は本 task では作らない。理由は UI 表示変更がなく、SCN-OBSLOG-004 は lower-level only であるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: fake provider ID、credential 状態分類、代表 ID、集約 count。
  - `secret_values_for_provider_external_api_internal_auth`: fake secret、実 secret、API key、endpoint 実値。
  - `secret_resolution_owner_layer`: fake secret store、provider settings service。
  - `forbidden_outputs`: log capture、error summary、request capture へ secret 本体、endpoint 実値、provider raw payload を出さない。

## Final Validation And Review

全 handoff 完了後に `implement_lane` が最終検証とレビューを判断する。

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite coverage`
- repo-local Sonar issue gate
- Codex review

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `N/A` または UI 表示変更なしの根拠
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes`: none

## Residual Risks

- backend log 追加箇所は既存の責務境界に従う必要がある。logger のために constructor 引数を広げる必要が出た場合は停止する。
- frontend runtime event log は UI 表示変更なしで閉じる。画面文言や表示状態の変更が必要になった場合は停止する。
- provider / secret 境界では credential 参照実値も log に出さない。分類値だけで原因分離できない場合は停止する。
- `observability-logger-lightweight` は既存基盤として扱い、docs 正本化や完了扱いの変更は本 implementation-scope へ含めない。
