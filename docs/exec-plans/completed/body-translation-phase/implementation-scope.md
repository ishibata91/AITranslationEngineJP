# Implementation Scope: body-translation-phase

- `skill`: implementation-scope
- `status`: ready-for-implementation
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: user message `approved そのままcloseまで進めて`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `scenario_design`: `./scenario-design.md`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `requirement_gate`: `./scenario-design.requirement-gate.md`
- `source_task`: `tasks/usecases/body-translation-phase.yaml`
- `upstream_tasks`: `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/persona-generation-phase.yaml`
- `downstream_task`: `tasks/usecases/translation-output-artifact.yaml`
- `reference_docs`: `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/screen-design/README.md`
- `reference_scope`: `docs/exec-plans/completed/term-translation-phase/implementation-scope.md`, `docs/exec-plans/completed/persona-generation-phase/implementation-scope.md`
- `code_map`: `tmp/code-map/index.json`

## Fixed Decisions

- `needs_human_decision`: `0`
- `unresolved_conflicts`: `0`
- `requirement_gate`: pass。`finding_count` は `0`、`question_count` は `0`。
- `scenario-gate`: final validation で実行する。
- 本文翻訳フェーズは NPC ペルソナ生成フェーズ Completed、非 terminal job、active phase run なし、前段参照成立の時だけ開始できる。
- provider、model、execution mode は `Job Setup` の本文翻訳用設定を使い、開始時の再選択 UI は作らない。
- 完全一致した辞書 hit は provider request から除外し、部分一致は訳語固定制約として provider request に渡す。
- 本文翻訳対象 0 件は Completed とし、provider 未実行でも output readiness を成立させる。
- 保護要素検証に失敗した訳文は保存前に拒否し、失敗訳文は保持しない。ただし retry 対象にする。
- 一部 field 失敗時は成功済み field result を保持し、phase 全体を `RecoverableFailed` として表示する。
- retry、resume、開始再送は同じ `JOB_PHASE_RUN` を継続し、field result と phase link を重複作成しない。
- cancel は `Paused` からだけ許可し、`Canceled` 後の途中成功結果は output readiness に使わない。
- body phase Completed で job-level `Completed` とし、完了済み job から成果物を出力できる。
- 原文と訳文がローカル UI に見えること自体は許容する。secret、API key 平文、復号可能値は UI、DTO、error summary、structured log、debug log、fake transport log に出さない。
- backend、frontend、統合境界は別 handoff に分ける。frontend handoff は確定済み `contract_freeze` に依存する。
- `APIテスト` は public seam 起点の system-level test とする。`UI人間操作E2E` は final validation で証明する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `contract-body-phase-public-seams` | なし | なし | なし |
| `wave-2` | `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter`, `frontend-body-phase-job-run-ui` | `contract-body-phase-public-seams` | `backend-body-phase-state-input-snapshot <-> backend-body-provider-adapter`, `backend-body-phase-state-input-snapshot <-> frontend-body-phase-job-run-ui`, `backend-body-provider-adapter <-> frontend-body-phase-job-run-ui` | なし |
| `wave-3` | `backend-body-field-result-persistence` | `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter` | なし | `depends_on` |
| `wave-4` | `backend-body-recovery-terminal-readiness` | `backend-body-field-result-persistence` | なし | `depends_on` |
| `wave-5` | `integration-body-phase-wails-gateway` | `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter`, `backend-body-field-result-persistence`, `backend-body-recovery-terminal-readiness`, `frontend-body-phase-job-run-ui` | なし | `depends_on` |
| `wave-6` | `final-validation-and-report` | `integration-body-phase-wails-gateway` | なし | `broad_gate_shared` |

## Handoffs

### `contract-body-phase-public-seams`

- `implementation_target`: Job Run から本文翻訳フェーズの summary、開始、pause、resume、retry、cancel、output readiness を扱う public seam と DTO を固定する。
- `implementation_artifact`: `contract_freeze`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-BTP-001` から `SCN-BTP-011`、`./ui-design.md` の `UI Contract`
  - `architecture_layer_basis`: Wails controller / DTO は backend controller 境界、usecase contract は backend usecase 境界、frontend gateway contract は frontend contract 境界に置く。依存方向は frontend gateway -> Wails adapter -> Wails controller -> usecase とする。
  - `frozen_public_seams`:
    - `GetBodyTranslationPhaseSummary`: job ID を受け取り、current phase、phase state、progress、target count、translated count、skipped count、provider summary、input summary、field result summary、error summary、button enablement、output readiness を返す。
    - `StartBodyTranslationPhase`: persona phase Completed、非 terminal job、active phase run なし、前段参照成立の job を受け取り、phase run ID、state、progress、input snapshot digest、開始不可理由を返す。
    - `PauseBodyTranslationPhase`: running body phase run を pause し、phase state と resume / cancel 可否を返す。
    - `ResumeBodyTranslationPhase`: paused または recoverable failed の同じ phase run を再開し、input snapshot digest と progress を維持して返す。
    - `RetryBodyTranslationPhase`: retryable failure の同じ phase run で未処理または失敗 field だけを再実行し、latest error と progress を更新する。
    - `CancelBodyTranslationPhase`: Paused の body phase run だけを cancel し、output readiness を false にする。
    - `GetBodyTranslationOutputReadiness`: body phase Completed、field result 整合、output status 整合の時だけ output readiness を true にする。
    - error kind は `persona_phase_incomplete`, `terminal_job`, `active_phase_exists`, `input_snapshot_failed`, `provider_failure`, `invalid_provider_response`, `protection_validation_failed`, `save_failed`, `output_readiness_blocked`, `late_response_rejected`, `secret_redacted` を区別する。
    - DTO は `credential_ref`、provider、model、execution mode、target count、dictionary digest、persona digest、metadata digest、prompt digest、request unit count、output count、error kind を持ち、secret 本体、API key、provider raw request / response、raw prompt を持たない。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, phase run ID, input snapshot digest, dictionary digest, persona digest, metadata digest, prompt digest, request unit count, output count, error kind, retryable flag
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, secret store value, raw credential material
  - `secret_resolution_owner_layer`: backend secret store adapter と AI provider transport adapter。usecase、controller、frontend gateway は `credential_ref` だけを扱う。
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、provider raw request / response、raw prompt、復号可能値を UI、DTO、read model、URL、error summary、structured log、debug log、fake transport log に出さない。
- `owned_scope`:
  - `internal/usecase/body_translation_phase_contract.go`
  - `internal/controller/wails/body_translation_phase_controller.go`
  - `internal/controller/wails/body_translation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/application/gateway-contract/body-translation-phase/*`
  - `frontend/src/controller/wails/gateway-dto/body-translation-phase/*`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `初手`: `internal/usecase/body_translation_phase_contract.go` に summary / command DTO、public error kind、redaction field obligation を追加し、`completion_signal` の「downstream が参照する field 名、nullability、error kind が固定される」を最初に閉じる。理由は backend、frontend、Wails gateway が同じ public seam に依存するため。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/controller/wails -run 'BodyTranslation|JobRun|SCN_BTP_(001|006|010|011)'`
  - `npm --prefix frontend run test -- --run body-translation-contract`
- `completion_signal`:
  - Job Run が参照する body phase public seam 名、request / response DTO、error kind が存在する。
  - field 名、nullability、phase state、retryable flag、output readiness、progress summary、redaction obligation が controller unit test と frontend contract test で固定される。
  - DTO は secret 本体、API key、provider raw request / response、raw prompt、復号可能値を表現できない。
  - frontend handoff が参照できる gateway contract が作成される。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`450-800 changed lines`。
  - contract freeze のみを扱い、永続化実体、provider adapter 実体、Job Run UI 実装は含めない。
  - `本番経路`: Wails controller / DTO -> usecase contract -> frontend gateway contract。

### `backend-body-phase-state-input-snapshot`

- `implementation_target`: body phase の開始条件、terminal guard、input snapshot、辞書 / persona / metadata 参照、完全一致除外、部分一致制約、対象 0 件 Completed を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、controller entry までを backend 境界として扱う。frontend UI と Wails gateway 実体は含めない。
  - `frozen_public_seams`: `contract-body-phase-public-seams` の completion signal を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: phase run ID, input snapshot digest, target count, skipped count, dictionary digest, persona digest, metadata digest, prompt digest, `credential_ref`
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。`credential_ref` は downstream provider adapter へ渡す参照値として扱う。
  - `forbidden_outputs`: secret 本体、API key、raw prompt、provider raw request / response、復号可能値を phase result、DB summary、structured log、error summary に出さない。
- `owned_scope`:
  - `internal/usecase/body_translation_phase_*`
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/body_translation_input_snapshot.go`
  - `internal/service/body_translation_phase_service_test.go`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/foundation_data_repository.go`
  - `internal/repository/foundation_data_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/infra/sqlite/dbinit/migrations/*body_translation*`
  - `internal/apitest/*body_translation*`
- `depends_on`: `contract-body-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-body-provider-adapter`, `frontend-body-phase-job-run-ui`
- `parallel_blockers`: なし
- `初手`: `internal/service/body_translation_phase_service.go` に persona phase Completed / terminal job / active phase run の開始判定と focused test を追加し、`completion_signal` の「開始可能条件と開始拒否理由が固定される」を最初に閉じる。理由は input snapshot、progress、output readiness が同じ phase state に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|BodyInputSnapshot|SCN_BTP_(001|002|008)'`
- `completion_signal`:
  - persona phase Completed、非 terminal job、active phase run なし、前段参照成立の場合だけ body phase run が作成される。
  - persona phase 未完了、active phase あり、terminal job、辞書 / persona 参照不能では phase run が作成されず、拒否理由が返る。
  - input snapshot は target field count、対象外理由、dictionary digest、persona digest、metadata digest、prompt digest を持つ。
  - 完全一致した辞書 hit は provider request 対象から除外され、部分一致は訳語固定制約として request summary に残る。
  - target count 0 は provider 未実行の Completed phase result になり、job-level Completed と output readiness へ進める summary を返す。
  - retry / resume / 開始再送で input snapshot digest が差し替わらず、同じ `JOB_PHASE_RUN` を使う。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `14-22 files`、`700-1300 changed lines`。
  - 1 受け入れユースケースは「body phase 開始と input snapshot 固定」で閉じる。provider 実行、field result 保存、Job Run UI は別 handoff に分けるため、分割必須にはしない。
  - schema 追加が必要な場合は body phase summary、input snapshot digest、phase link に限定する。docs 正本化は含めない。
  - `本番経路`: usecase -> service -> repository / transactor -> SQLite。

### `backend-body-provider-adapter`

- `implementation_target`: 翻訳レコード種別と field type に応じた prompt builder、fake provider、response validation、field correlation、provider failure handling、redaction / audit summary を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: AIProvider / prompt builder / response adapter boundary は backend service と infra / provider 境界に置き、service は provider-agnostic contract だけを見る。
  - `frozen_public_seams`: provider adapter output は field correlation、translated candidate、protection source digest、retryable failure、redacted error summary、audit summary を返す。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, request unit ID, input snapshot digest, prompt digest, request unit count, output count, error kind
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, provider secret material
  - `secret_resolution_owner_layer`: AI provider transport adapter が provider 呼び出し直前に secret 本体を解決する。prompt builder と response adapter は secret 本体を受け取らない。
  - `forbidden_outputs`: secret 本体、API key、authorization header、provider raw request / response、raw provider payload、raw prompt、復号可能値を UI、DTO、DB summary、structured log、debug log、fake transport log、error summary、URL、read model に出さない。
- `owned_scope`:
  - `internal/infra/ai/provider.go`
  - `internal/infra/ai/provider_client.go`
  - `internal/infra/ai/gemini.go`
  - `internal/infra/ai/openai_compatible.go`
  - `internal/infra/ai/transport.go`
  - `internal/infra/ai/*body_translation*`
  - `internal/service/body_translation_provider_adapter.go`
  - `internal/service/body_translation_prompt_builder.go`
  - `internal/service/body_translation_response_adapter.go`
  - `internal/service/body_translation_provider_adapter_test.go`
  - `internal/infra/ai/*body_translation*_test.go`
- `depends_on`: `contract-body-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-body-phase-state-input-snapshot`, `frontend-body-phase-job-run-ui`
- `parallel_blockers`: なし
- `初手`: `internal/service/body_translation_provider_adapter.go` に field correlation を保持する provider-agnostic input / output contract と invalid response test を追加し、`completion_signal` の「valid response が field correlation を保持して translated candidate へ写像される」を最初に閉じる。理由は fake transport、Gemini / xAI adapter、field result persistence が同じ response correlation に依存するため。
- `validation_commands`:
  - `go test ./internal/infra/ai ./internal/service -run 'BodyTranslationProvider|BodyProviderAdapter|Redaction|SCN_BTP_(003|011)'`
- `completion_signal`:
  - provider、model、execution mode は `Job Setup` の本文翻訳用設定から継承される。
  - record type、field type、field correlation key、保護要素 digest を失わず provider 境界を通過できる。
  - 完全一致辞書 hit は provider request に含まれず、部分一致は訳語固定制約として request summary に残る。
  - valid response は field correlation を保持し、translated candidate と保護要素検証対象へ写像される。
  - invalid response、field 欠落 response、余分な response、correlation error、timeout は translated field として保存されない adapter output になる。
  - provider failure で別 provider へ暗黙 fallback しない。
  - fake provider と fixed response で paid real API を呼ばずに検証できる。
  - API key、secret 本体、復号可能値は fake transport log、structured log、debug log、error summary に出ない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `10-15 files`、`500-800 changed lines`。
  - provider adapter と prompt builder だけを扱い、phase state、field result 永続化、Job Run UI は含めない。
  - `本番経路`: service provider adapter -> internal/infra/ai provider client -> HTTP transport。

### `backend-body-field-result-persistence`

- `implementation_target`: provider 応答の訳文を保護要素検証し、訳文、出力ステータス、検証結果、phase link を field 単位で保存する backend 永続化を実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、state machine、JobIOService を backend 境界として扱う。provider adapter output と input snapshot を入力にし、frontend UI は含めない。
  - `frozen_public_seams`: `contract-body-phase-public-seams` の completion signal と `backend-body-provider-adapter` の provider adapter output を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: phase run ID, `JOB_TRANSLATION_FIELD` ID, `PHASE_RUN_TRANSLATION_FIELD` ID, field correlation key, output status, validation status, validation summary, translated text, retry count
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。保存対象は provider output の訳文と redacted summary だけに限定する。
  - `forbidden_outputs`: secret 本体、API key、authorization header、provider raw request / response、raw prompt、復号可能値を `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase result summary、structured log、debug log、error summary に出さない。
- `owned_scope`:
  - `internal/usecase/body_translation_phase_*`
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/body_translation_field_result.go`
  - `internal/service/body_translation_protection_validator.go`
  - `internal/service/body_translation_field_result_test.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/statemachine/*body*`
  - `internal/jobio/*body*`
  - `internal/infra/sqlite/dbinit/migrations/*body_translation*`
  - `internal/apitest/*body_translation*`
- `depends_on`: `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `初手`: `internal/service/body_translation_protection_validator.go` に保護要素欠落を validation failed にする focused test を追加し、`completion_signal` の「保護要素検証失敗は保存前に拒否される」を最初に閉じる。理由は訳文保存、出力ステータス、retry target が検証結果に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|Protection|FieldResult|SCN_BTP_(004|005)'`
- `completion_signal`:
  - 保護要素検証を通過した field だけが `JOB_TRANSLATION_FIELD` に訳文、output status、retry count を保存できる。
  - `PHASE_RUN_TRANSLATION_FIELD` が phase run と field result を重複なく紐づける。
  - 保護要素欠落、改変、重複、順序違い、余分追加は validation failed として観測できる。
  - 保護要素検証失敗の訳文は保存前に拒否され、successful translated field として保存されない。
  - 保存失敗または検証失敗では phase state が Completed にならない。
  - 後続 output artifact に必要な output status 語彙だけを保持する。
  - field correlation key と保存先 field が一致する。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-24 files`、`900-1400 changed lines`。
  - 1 受け入れユースケースは「保護要素検証後の field result 保存」で閉じる。recovery、terminal guard、output readiness、frontend UI は別 handoff に分けるため、分割必須にはしない。
  - schema 追加が必要な場合は field result と phase link に限定する。docs 正本化は含めない。
  - `本番経路`: usecase -> service -> protection validator -> repository / transactor -> SQLite。

### `backend-body-recovery-terminal-readiness`

- `implementation_target`: 部分失敗、retry、resume、開始再送、Paused からの cancel、terminal guard、late response rejection、output readiness を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: service / usecase、state machine、JobIOService、repository を backend 境界として扱う。frontend UI と Wails gateway 実体は含めない。
  - `frozen_public_seams`: `contract-body-phase-public-seams` の completion signal と `backend-body-field-result-persistence` の field result persistence を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: phase run ID, phase state, retryable flag, affected field count, output readiness, blocked reason, error kind, progress, field status summary
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。retry では `credential_ref` を provider adapter へ渡す参照値として扱う。
  - `forbidden_outputs`: secret 本体、API key、provider raw request / response、raw prompt、復号可能値を phase result、readiness summary、structured log、debug log、error summary に出さない。
- `owned_scope`:
  - `internal/usecase/body_translation_phase_*`
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/body_translation_recovery.go`
  - `internal/service/body_translation_readiness.go`
  - `internal/service/body_translation_recovery_test.go`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/statemachine/*body*`
  - `internal/jobio/*body*`
  - `internal/apitest/*body_translation*`
- `depends_on`: `backend-body-field-result-persistence`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `初手`: `internal/service/body_translation_recovery.go` に同じ `JOB_PHASE_RUN` で retry target を算出する focused test を追加し、`completion_signal` の「retry / resume / 開始再送で field result と phase link を重複作成しない」を最初に閉じる。理由は partial failure、terminal guard、readiness が同じ phase run invariant に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|Recovery|Readiness|Terminal|SCN_BTP_(007|008|009|010)'`
- `completion_signal`:
  - provider failure、invalid response、correlation error、save failure では successful translated field として保存されない。
  - 部分失敗では成功済み field result を保持し、phase 全体は `RecoverableFailed` になり、失敗 field だけ retry target になる。
  - retry、resume、開始再送では同じ `JOB_PHASE_RUN` を使い、`JOB_TRANSLATION_FIELD` と `PHASE_RUN_TRANSLATION_FIELD` を重複作成しない。
  - Running から直接 cancel できず、Paused からだけ `Canceled` へ遷移する。
  - terminal job では body phase run 作成、field save、readiness update、late response 後書きを拒否する。
  - Canceled 後の途中成功結果は output readiness に使われない。
  - body phase Completed かつ field result と output status が整合する時だけ job-level `Completed` と output readiness が成立する。
  - 未完了、失敗中、status 不整合では output artifact readiness が false になり、拒否理由が返る。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-24 files`、`900-1500 changed lines`。
  - 1 受け入れユースケースは「body phase recovery と output readiness」で閉じる。provider transport、field result 保存、Job Run UI は別 handoff に分けるため、分割必須にはしない。
  - `本番経路`: usecase -> service -> state machine / JobIOService -> repository / transactor -> SQLite -> output readiness query。

### `frontend-body-phase-job-run-ui`

- `implementation_target`: Job Run に本文翻訳フェーズの current phase、progress、input summary、execution summary、field result、recovery panel、output readiness、主要 action を表示する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: frontend contract、store / presenter / usecase / controller、Job Run screen を frontend 境界として扱う。Wails 接続実体は統合境界 handoff に分ける。
  - `frozen_public_seams`: `contract-body-phase-public-seams` の frontend gateway contract。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, target count, processed count, translated count, failed count, skipped count, dictionary digest, persona digest, metadata digest, prompt digest, output readiness, error kind, retryable flag, translated text
  - `secret_values_for_provider_external_api_internal_auth`: なし。frontend は secret 本体を受け取らない。
  - `secret_resolution_owner_layer`: frontend では解決しない。backend secret store adapter と AI provider transport adapter に限定する。
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、raw prompt、raw response、復号可能値を画面、console、error summary、frontend store、frontend DTO に出さない。
- `owned_scope`:
  - `frontend/src/application/contract/body-translation-phase/*`
  - `frontend/src/application/store/body-translation-phase/*`
  - `frontend/src/application/presenter/body-translation-phase/*`
  - `frontend/src/application/usecase/body-translation-phase/*`
  - `frontend/src/controller/body-translation-phase/*`
  - `frontend/src/ui/screens/body-translation-phase/*`
  - `frontend/src/ui/screens/job-run/*`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
- `depends_on`: `contract-body-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter`
- `parallel_blockers`: なし
- `初手`: `frontend/src/application/contract/body-translation-phase/body-translation-phase-screen-contract.ts` に Job Run 表示状態と action enablement contract を追加し、`completion_signal` の「phase state、progress、field result summary、button enablement を型で表せる」を最初に閉じる。理由は store、presenter、UI が同じ表示状態に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run body-translation-phase`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - Job Run で current phase、phase state、progress、target count、processed count、success count、failure count、skipped count を表示できる。
  - input summary に dictionary digest、persona digest、metadata digest、provider setting、credential ref 状態、prompt digest を表示できる。
  - execution summary に provider、model、execution mode、request unit count、output count、provider skipped、late response rejected を表示できる。
  - field result list に field identity、source excerpt boundary、translated text、output status、protection validation result を表示できる。
  - error summary は error kind、retryable flag、affected field count、redacted summary を表示し、secret と API key 平文を出さない。
  - start、pause、resume、retry、cancel、output readiness の enablement が UI contract と一致する。
  - loading、not-ready、ready、running、paused、recoverable failed、validation failed、empty completed、completed、canceled、failed の状態差分を表示できる。
  - 長い source / translated text、FormID、EditorID、provider model 名、error kind が desktop / mobile で overflow しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`900-1500 changed lines`。
  - UI人間操作E2E の最終証明は `final-validation-and-report` に寄せる。この handoff の local validation は mocked gateway の frontend tests と frontend check に限定する。
  - Wails gateway 実体は `integration-body-phase-wails-gateway` に分ける。
  - `本番経路`: frontend gateway contract -> usecase -> store / presenter -> Job Run screen。

### `integration-body-phase-wails-gateway`

- `implementation_target`: backend usecase、Wails controller、frontend Wails gateway / DTO を接続し、Job Run UI が実 backend の body phase public seam を呼べる状態にする。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-body-phase-public-seams`
  - `architecture_layer_basis`: Wails controller、bootstrap binding、frontend Wails adapter、gateway DTO の接続だけを統合境界として扱う。
  - `frozen_public_seams`: `GetBodyTranslationPhaseSummary`, `StartBodyTranslationPhase`, `PauseBodyTranslationPhase`, `ResumeBodyTranslationPhase`, `RetryBodyTranslationPhase`, `CancelBodyTranslationPhase`, `GetBodyTranslationOutputReadiness`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, phase run ID, input snapshot digest, count summary, field summary, error kind, retryable flag, output readiness
  - `secret_values_for_provider_external_api_internal_auth`: なし。Wails DTO と frontend gateway は secret 本体を受け取らない。
  - `secret_resolution_owner_layer`: backend secret store adapter と AI provider transport adapter。Wails controller は redacted DTO だけを返す。
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、raw prompt、raw response、復号可能値を Wails DTO、frontend gateway DTO、console、UI、read model に出さない。
- `owned_scope`:
  - `internal/controller/wails/body_translation_phase_controller.go`
  - `internal/controller/wails/body_translation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/controller/wails/body-translation-phase.gateway.ts`
  - `frontend/src/controller/wails/body-translation-phase.gateway.test.ts`
  - `frontend/src/controller/wails/gateway-dto/body-translation-phase/*`
  - `frontend/src/application/gateway-contract/body-translation-phase/*`
- `depends_on`: `backend-body-phase-state-input-snapshot`, `backend-body-provider-adapter`, `backend-body-field-result-persistence`, `backend-body-recovery-terminal-readiness`, `frontend-body-phase-job-run-ui`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `初手`: `internal/controller/wails/body_translation_phase_controller.go` に `StartBodyTranslationPhase` の controller method と DTO mapping test を追加し、`completion_signal` の「Wails controller が body phase usecase を呼び、redacted response を返す」を最初に閉じる。理由は bootstrap と frontend gateway mapping が controller method 名に依存するため。
- `validation_commands`:
  - `go test ./internal/controller/wails ./internal/bootstrap -run 'BodyTranslation|JobRun'`
  - `npm --prefix frontend run test -- --run body-translation-phase.gateway`
- `completion_signal`:
  - Wails controller が backend usecase の summary / command DTO を frontend gateway DTO へ lossless に写像する。
  - bootstrap と app controller に body phase controller が接続される。
  - frontend Wails gateway が frozen gateway contract を満たす。
  - error kind、retryable flag、phase state、progress、input summary、execution summary、field result summary、output readiness が backend / frontend DTO で一致する。
  - secret、API key、raw prompt、raw response、復号可能値は DTO mapping で落とされる。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`400-800 changed lines`。
  - backend core と frontend UI の代替実装は含めない。Wails / DTO / gateway 接続だけを扱う。
  - `本番経路`: Wails controller -> usecase -> frontend Wails gateway -> frontend usecase。

### `final-validation-and-report`

- `implementation_target`: 全 handoff 完了後に scenario、broad gate、UI 証跡、Codex work report 入力を確認する。
- `implementation_artifact`: `final validation / report input`
- `implementation_skill`: `implement-lane`
- `contract_freeze`:
  - `status`: `not_required`
  - `freeze_source`: `N/A`
  - `architecture_layer_basis`: `N/A`
  - `frozen_public_seams`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: final evidence に載せる redacted summary、test result、UI screenshot note、redaction assertion result
  - `secret_values_for_provider_external_api_internal_auth`: なし。final validation は secret 本体を収集しない。
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、raw prompt、raw response、復号可能値を work report、UI evidence、final validation summary に出さない。
- `owned_scope`:
  - `work_history/runs/body-translation-phase/codex.md`
  - 実装で必要になった task-local residual note
- `depends_on`: `integration-body-phase-wails-gateway`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: なし
- `parallel_blockers`: `broad_gate_shared`
- `初手`: `work_history/runs/body-translation-phase/codex.md` を作成し、final validation 欄と completion evidence 欄を先に固定する。理由は gate 結果、UI 証跡、環境 blocker を completion packet の根拠にするため。
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/body-translation-phase/scenario-design.md --coverage docs/exec-plans/completed/body-translation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/body-translation-phase/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `go test ./internal/...`
  - `npm run test:frontend`
  - `npm --prefix frontend run check`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - requirement gate と scenario gate が pass する。
  - relevant backend、integration、frontend tests が pass する。
  - Job Run の desktop / mobile UI 証跡で current phase、progress、input summary、execution summary、field result、error summary、button state、output readiness が重ならない。
  - UI、console、error summary、structured log、debug log、fake transport log に secret、API key、復号可能値がない。
  - paid API が呼ばれていない証跡を fake transport log または test result で確認できる。
  - body phase Completed で job-level Completed と output readiness が成立し、translation-output-artifact の開始条件が確認できる。
  - system / harness が環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason、再実行環境、再実行コマンドを report に残す。
  - Codex work report が completion packet の schema を満たす。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `最終検証`
- `notes`:
  - 想定規模は normal。想定 `1-3 files`、`100-300 changed lines`。
  - product 実装は含めない。broad validation owner は全 handoff 完了後だけに置く。
  - Sonar を使う場合は `/tmp` cache を使い、Sonar server の Quality Gate と repo-local issue gate を混同しない。

## Codex Implementation Handoff Packet

- `entry`: `.codex/skills/implement-lane/SKILL.md`
- `task_id`: `body-translation-phase`
- `source_scope`: `docs/exec-plans/completed/body-translation-phase/implementation-scope.md`
- `human_review_status`: approved
- `approval_record`: user message `approved そのままcloseまで進めて`
- `ready_for_implementation`: true
- `start_wave`: `wave-1`
- `ready_now`: `contract-body-phase-public-seams`
- `forbidden_changes`: product docs 正本、`.codex/`、`.codex/skills`、`.codex/agents`
- `docs_changes`: none

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
- `completion_evidence`: Codex 側 `work_reporter` が読む実装事実。completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`

## Open Items

- なし。
