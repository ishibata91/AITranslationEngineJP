# Implementation Scope: persona-generation-phase

- `skill`: implementation-scope
- `status`: ready-for-implementation
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human message `approved`
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
- `source_task`: `tasks/usecases/persona-generation-phase.yaml`
- `upstream_task`: `tasks/usecases/term-translation-phase.yaml`
- `downstream_task`: `tasks/usecases/body-translation-phase.yaml`
- `reference_docs`: `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/screen-design/README.md`
- `reference_scope`: `docs/exec-plans/completed/term-translation-phase/implementation-scope.md`
- `code_map`: `tmp/code-map/index.json`

## Fixed Decisions

- `needs_human_decision`: `0`
- `unresolved_conflicts`: `0`
- `requirement_gate`: pass。`finding_count` は `0`、`question_count` は `0`。
- `scenario-gate`: final validation で実行する。
- 単語翻訳フェーズ Completed、非 terminal job、active phase run なしの場合だけ persona phase を開始できる。
- `Completed`、`Failed`、`Canceled` を terminal state とする。`RecoverableFailed` は回復対象として扱う。
- 共通ペルソナ hit 時は新規 `PERSONA` を作らず、job の persona snapshot 参照だけを固定する。
- persona 生成は 1 NPC を 1 request unit とし、NPC 属性と会話文脈を同じ request で扱う。
- Job Setup の persona 専用 provider、model、execution mode を継承する。
- valid provider output は自動採用する。人間確認 UI は含めない。
- 生成対象 0 件は Completed とし、対象 0 件、provider 未実行、snapshot 空を result summary に出す。
- 一部 NPC 失敗時は成功分を維持し、phase は `RecoverableFailed` として未処理 NPC だけ retry する。
- persona phase は pause、resume、retry、cancel を body phase と同じ粒度で許可する。
- UI と DB summary は ID、digest、件数、evidence ref、redacted phase result summary だけを出し、全文と raw prompt は出さない。
- debug log では prompt / request body を確認できる導線を持つ。ただし secret と API key は debug log、structured log、UI、DTO、DB summary に出さない。
- Job Run 再表示用の redacted phase result summary は保持する。直接 DB 保存に限定せず、進行中の job state から復元できる形を許容する。
- backend、frontend、統合境界は別 handoff に分ける。frontend handoff は確定済み `contract_freeze` に依存する。
- `APIテスト` は public seam 起点の system-level test とする。`UI人間操作E2E` は最終検証で証明する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `contract-persona-phase-public-seams` | なし | なし | なし |
| `wave-2` | `backend-persona-phase-state-targets`, `backend-persona-provider-adapter`, `frontend-persona-phase-job-run-ui` | `contract-persona-phase-public-seams` | `backend-persona-phase-state-targets <-> backend-persona-provider-adapter`, `backend-persona-phase-state-targets <-> frontend-persona-phase-job-run-ui`, `backend-persona-provider-adapter <-> frontend-persona-phase-job-run-ui` | なし |
| `wave-3` | `backend-persona-persistence-readiness-retry` | `backend-persona-phase-state-targets`, `backend-persona-provider-adapter` | なし | `depends_on` |
| `wave-4` | `integration-persona-phase-wails-gateway` | `backend-persona-persistence-readiness-retry`, `frontend-persona-phase-job-run-ui` | なし | `depends_on` |
| `wave-5` | `final-validation-and-report` | `integration-persona-phase-wails-gateway` | なし | `broad_gate_shared` |

## Handoffs

### `contract-persona-phase-public-seams`

- `implementation_target`: Job Run から persona phase の summary、開始、pause、resume、retry、cancel、body readiness を扱う public seam と DTO を固定する。
- `implementation_artifact`: `contract_freeze`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-PGP-001` から `SCN-PGP-010`、`./ui-design.md` の `UI Contract`
  - `architecture_layer_basis`: Wails controller / DTO は backend controller 境界、usecase contract は backend usecase 境界、frontend gateway contract は frontend contract 境界に置く。依存方向は frontend gateway -> Wails adapter -> Wails controller -> usecase とする。
  - `frozen_public_seams`:
    - `GetPersonaGenerationPhaseSummary`: job ID を受け取り、current phase、phase state、progress、target count、generated count、failed count、skipped count、snapshot summary、AI execution summary、redacted error summary、button enablement を返す。
    - `StartPersonaGenerationPhase`: term phase Completed、非 terminal job、active phase run なしの job を受け取り、phase run ID、state、progress、target snapshot digest、開始不可理由を返す。
    - `PausePersonaGenerationPhase`: running persona phase run を pause し、phase state と resume / cancel 可否を返す。
    - `ResumePersonaGenerationPhase`: paused または recoverable failed の同じ phase run を再開し、target snapshot digest と progress を維持して返す。
    - `RetryPersonaGenerationPhase`: retryable failure の同じ phase run で未処理 NPC だけを再実行し、latest error と progress を更新する。
    - `CancelPersonaGenerationPhase`: cancel 可能な non-terminal job だけを中止し、後続 phase readiness を false にする。
    - `GetPersonaGenerationBodyReadiness`: persona phase Completed かつ snapshot 参照成立時だけ body phase readiness を true にする。
    - error kind は `term_phase_incomplete`, `terminal_job`, `active_phase_exists`, `target_snapshot_failed`, `provider_failure`, `invalid_provider_response`, `save_failed`, `snapshot_missing`, `body_readiness_blocked`, `secret_redacted` を区別する。
    - DTO は `credential_ref`、provider、model、execution mode、target snapshot digest、prompt digest、input count、output count、evidence ref を持ち、secret 本体、API key、provider raw request / response、原文発話全文、会話文脈全文を持たない。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, phase run ID, snapshot ID, snapshot digest, prompt digest, evidence ref, input count, output count, error kind
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, secret store value, raw credential material
  - `secret_resolution_owner_layer`: backend secret store adapter と AI provider transport adapter。usecase、controller、frontend gateway は `credential_ref` だけを扱う。
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、provider raw request / response、raw prompt、原文発話全文、会話文脈全文を UI、DTO、DB summary、error summary、structured log、URL、read model に出さない。debug log では prompt / request body を許可するが、secret と API key は redacted 済みにする。
- `owned_scope`:
  - `internal/usecase/persona_generation_phase_contract.go`
  - `internal/controller/wails/persona_generation_phase_controller.go`
  - `internal/controller/wails/persona_generation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/application/gateway-contract/persona-generation-phase/*`
  - `frontend/src/controller/wails/gateway-dto/persona-generation-phase/*`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`: `internal/usecase/persona_generation_phase_contract.go` に summary / command DTO、public error kind、redaction field obligation を追加し、`completion_signal` の「downstream が参照する field 名、nullability、error kind が固定される」を最初に閉じる。理由は backend、frontend、Wails gateway が同じ seam に依存するため。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/controller/wails -run 'PersonaGeneration|JobRun|SCN_PGP_(001|005|006|009)'`
  - `npm --prefix frontend run test -- --run persona-generation-contract`
- `completion_signal`:
  - Job Run が参照する persona phase public seam 名、request / response DTO、error kind が存在する。
  - field 名、nullability、phase state、retryable flag、body readiness、progress summary、redaction obligation が controller unit test と frontend contract test で固定される。
  - DTO は secret 本体、API key、provider raw request / response、raw prompt、原文発話全文、会話文脈全文を表現できない。
  - frontend handoff が参照できる gateway contract が作成される。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`450-800 changed lines`。
  - contract freeze のみを扱い、永続化実体、provider adapter 実体、Job Run UI 実装は含めない。
  - `本番経路`: Wails controller / DTO -> usecase contract -> frontend gateway contract。

### `backend-persona-phase-state-targets`

- `implementation_target`: persona phase の開始条件、terminal guard、target snapshot、common persona hit / miss、生成対象 0 件、phase progress の backend 状態を実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-persona-phase-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、controller entry までを backend 境界として扱う。frontend UI と Wails gateway 実体は含めない。
  - `frozen_public_seams`: `contract-persona-phase-public-seams` の completion signal を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: phase run ID, target snapshot digest, NPC count, common persona hit / miss count, skipped count, evidence ref, `credential_ref`
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。`credential_ref` は downstream provider adapter へ渡す参照値として扱う。
  - `forbidden_outputs`: secret 本体、API key、raw prompt、原文発話全文、会話文脈全文を phase result、DB summary、structured log、error summary に出さない。
- `owned_scope`:
  - `internal/usecase/persona_generation_phase_*`
  - `internal/service/persona_generation_phase_*`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/foundation_data_repository.go`
  - `internal/repository/foundation_data_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/infra/sqlite/dbinit/migrations/*persona_generation*`
  - `internal/integrationtest/*persona_generation*`
- `depends_on`: `contract-persona-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-persona-provider-adapter`, `frontend-persona-phase-job-run-ui`
- `parallel_blockers`: なし
- `first_action`: `internal/service/persona_generation_phase_service.go` に term phase Completed / terminal job / active phase run の開始判定と focused test を追加し、`completion_signal` の「開始可能条件と開始拒否理由が固定される」を最初に閉じる。理由は target snapshot、progress、後続 readiness が同じ phase state に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest -run 'PersonaGeneration|PhaseRun|TargetSnapshot|SCN_PGP_(001|002|010)'`
- `completion_signal`:
  - term phase Completed、非 terminal job、active phase run なしの場合だけ persona phase run が作成される。
  - term 未完了、active phase あり、terminal job では phase run が作成されず、拒否理由が返る。
  - target snapshot は NPC record、translation field reference、common persona hit / miss、対象外理由、snapshot digest を持つ。
  - common persona hit 時は新規 `PERSONA` を作らず、job の persona snapshot 参照だけを固定する。
  - 生成対象 0 件は Completed とし、provider 未実行、snapshot 空、対象 0 件が redacted summary に残る。
  - target count と progress total が一致する。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `14-22 files`、`700-1300 changed lines`。
  - 1 受け入れユースケースは「persona phase 開始と target snapshot 固定」で閉じる。provider 実行、persona 保存、frontend UI は別 handoff に分けるため、分割必須にはしない。
  - schema 追加が必要な場合は `JOB_PHASE_RUN` の persona phase summary と target snapshot digest に限定する。docs 正本化は含めない。
  - `本番経路`: usecase -> service -> repository / transactor -> SQLite。

### `backend-persona-provider-adapter`

- `implementation_target`: 1 NPC 1 request unit の prompt builder、fake provider、response validation、provider failure handling、debug log redaction を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-persona-phase-public-seams`
  - `architecture_layer_basis`: AIProvider / prompt builder / response adapter boundary は backend service と infra / provider 境界に置き、service は provider-agnostic contract だけを見る。
  - `frozen_public_seams`: provider adapter output は NPC correlation、persona result、retryable failure、redacted error summary、debug log redaction status を返す。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, request unit ID, target snapshot digest, prompt digest, input count, output count, error kind
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, provider secret material
  - `secret_resolution_owner_layer`: AI provider transport adapter が provider 呼び出し直前に secret 本体を解決する。prompt builder と response adapter は secret 本体を受け取らない。
  - `forbidden_outputs`: secret 本体、API key、authorization header、provider raw response、raw provider payload を UI、DTO、DB summary、structured log、error summary、URL、read model に出さない。debug log の prompt / request body は secret と API key を redacted した内容だけ許可する。
- `owned_scope`:
  - `internal/infra/ai/provider.go`
  - `internal/infra/ai/provider_client.go`
  - `internal/infra/ai/gemini.go`
  - `internal/infra/ai/openai_compatible.go`
  - `internal/infra/ai/transport.go`
  - `internal/infra/ai/*persona_generation*`
  - `internal/service/persona_generation_provider_*`
  - `internal/service/*persona_generation*_test.go`
  - `internal/infra/ai/*persona_generation*_test.go`
- `depends_on`: `contract-persona-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-persona-phase-state-targets`, `frontend-persona-phase-job-run-ui`
- `parallel_blockers`: なし
- `first_action`: `internal/service/persona_generation_provider_adapter.go` に 1 NPC request unit の provider-agnostic input / output contract と invalid response test を追加し、`completion_signal` の「valid response が NPC correlation を保持して persona result へ写像される」を最初に閉じる。理由は fake transport、Gemini / xAI adapter、persistence handoff が同じ correlation に依存するため。
- `validation_commands`:
  - `go test ./internal/infra/ai ./internal/service -run 'PersonaGenerationProvider|PersonaProviderAdapter|Redaction|SCN_PGP_(003|009)'`
- `completion_signal`:
  - provider、model、execution mode は Job Setup の persona 専用設定から継承される。
  - 1 NPC が 1 request unit になり、NPC 属性と会話文脈を同じ request で扱う。
  - valid response は NPC correlation を保持し、persona result へ写像される。
  - invalid response、対象 NPC 欠落 response、余分な response、timeout は persona として保存されない adapter output になる。
  - provider failure で別 provider へ暗黙 fallback しない。
  - fake provider と fixed response で paid real API を呼ばずに検証できる。
  - debug log では prompt / request body を確認できるが、secret と API key は redacted される。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`450-800 changed lines`。
  - provider adapter と prompt builder だけを扱い、phase state、persona 永続化、Job Run UI は含めない。
  - `本番経路`: service provider adapter -> internal/infra/ai provider client -> HTTP transport。

### `frontend-persona-phase-job-run-ui`

- `implementation_target`: Job Run に persona phase の current phase、progress、target summary、phase result、AI execution summary、error summary、body readiness、主要 action を表示する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-persona-phase-public-seams`
  - `architecture_layer_basis`: frontend contract、store / presenter / usecase / controller、Job Run screen を frontend 境界として扱う。Wails 接続実体は統合境界 handoff に分ける。
  - `frozen_public_seams`: `contract-persona-phase-public-seams` の frontend gateway contract。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, target count, generated count, failed count, skipped count, snapshot ID, snapshot digest, prompt digest, error kind, retryable flag
  - `secret_values_for_provider_external_api_internal_auth`: なし。frontend は secret 本体を受け取らない。
  - `secret_resolution_owner_layer`: frontend では解決しない。backend secret store adapter と AI provider transport adapter に限定する。
  - `forbidden_outputs`: secret 本体、API key、token、raw prompt、raw response、原文発話全文、会話文脈全文を画面、console、error summary、frontend store、frontend DTO に出さない。
- `owned_scope`:
  - `frontend/src/application/contract/persona-generation-phase/*`
  - `frontend/src/application/store/persona-generation-phase/*`
  - `frontend/src/application/presenter/persona-generation-phase/*`
  - `frontend/src/application/usecase/persona-generation-phase/*`
  - `frontend/src/controller/persona-generation-phase/*`
  - `frontend/src/ui/screens/persona-generation-phase/*`
  - `frontend/src/ui/screens/job-run/*`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
- `depends_on`: `contract-persona-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-persona-phase-state-targets`, `backend-persona-provider-adapter`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/contract/persona-generation-phase/persona-generation-phase-screen-contract.ts` に Job Run 表示状態と action enablement contract を追加し、`completion_signal` の「phase state、progress、result summary、button enablement を型で表せる」を最初に閉じる。理由は store、presenter、UI が同じ表示状態に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run persona-generation-phase`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - Job Run で `NPC ペルソナ生成`、phase state、progress、target count、generated count、failed count、skipped count を表示できる。
  - target summary に NPC count、common persona hit / miss、対象外理由を表示できる。
  - phase result に persona snapshot ID または snapshot digest、snapshot 参照状態、missing count、body readiness を表示できる。
  - AI execution summary に provider、model、execution mode、credential ref、input count、output count、短い error kind を表示できる。
  - 開始、pause、resume、retry、cancel、body readiness 確認、body phase 開始の enablement が UI contract と一致する。
  - loading、not started、running、paused、recoverable failed、failed、completed、empty completed、blocked、snapshot missing の状態差分を表示できる。
  - secret、API key 平文、raw prompt、raw response、原文発話全文、会話文脈全文を画面、error summary、console に出さない。
  - 長い NPC 名、provider 名、model 名、error reason、snapshot digest、credential ref は desktop / mobile で overflow しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`850-1400 changed lines`。
  - UI人間操作E2E の最終証明は `final-validation-and-report` に寄せる。この handoff の local validation は mocked gateway の frontend tests と frontend check に限定する。
  - Wails gateway 実体は `integration-persona-phase-wails-gateway` に分ける。
  - `本番経路`: frontend gateway contract -> usecase -> store / presenter -> Job Run screen。

### `backend-persona-persistence-readiness-retry`

- `implementation_target`: valid persona result を job-scoped persona、evidence、phase link、snapshot summary へ保存し、failure、resume、retry、body readiness を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-persona-phase-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、state machine、JobIOService を backend 境界として扱う。provider adapter output と target snapshot を入力にし、frontend UI は含めない。
  - `frozen_public_seams`: `contract-persona-phase-public-seams` の completion signal と `backend-persona-provider-adapter` の provider adapter output を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: phase run ID, `PERSONA` ID, `PHASE_RUN_PERSONA` ID, evidence ref, snapshot digest, persona count, missing count, input count, output count, error kind, retryable flag
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。保存対象は provider output の redacted result と summary だけに限定する。
  - `forbidden_outputs`: secret 本体、API key、raw prompt、provider raw request / response、原文発話全文、会話文脈全文を `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA`、phase result summary、structured log、error summary に出さない。
- `owned_scope`:
  - `internal/usecase/persona_generation_phase_*`
  - `internal/service/persona_generation_phase_*`
  - `internal/repository/persona_repository.go`
  - `internal/repository/persona_sqlite_repository.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/statemachine/*persona*`
  - `internal/jobio/*persona*`
  - `internal/infra/sqlite/dbinit/migrations/*persona_generation*`
  - `internal/integrationtest/*persona_generation*`
- `depends_on`: `backend-persona-phase-state-targets`, `backend-persona-provider-adapter`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `first_action`: `internal/service/persona_generation_phase_service.go` に valid persona result の atomic save と save failure test を追加し、`completion_signal` の「partial save failure は Completed にならない」を最初に閉じる。理由は body readiness、retry、phase result が保存の atomicity に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest -run 'PersonaGeneration|PersonaPersistence|BodyReadiness|Retry|SCN_PGP_(004|006|007|008|010)'`
- `completion_signal`:
  - valid provider output は人間確認なしで job-scoped persona または snapshot ref へ反映される。
  - `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA` が phase run、target snapshot、evidence ref と整合する。
  - common persona hit 時は新規 `PERSONA` を作らず、job の snapshot ref だけを保持する。
  - save failure、invalid response、provider failure、input missing では partial state を Completed にしない。
  - 一部 NPC 失敗時は成功分を維持し、phase は `RecoverableFailed` として未処理 NPC だけ retry 対象にする。
  - resume / retry / 開始再送では同じ `JOB_PHASE_RUN` を使い、成功済み persona と `PHASE_RUN_PERSONA` を重複作成しない。
  - persona phase Completed かつ snapshot 参照成立時だけ body readiness が true になる。
  - terminal job では persona phase start、persona save、body readiness update が拒否され、state snapshot は変更されない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`900-1500 changed lines`。
  - 1 受け入れユースケースは「persona result の永続化と body readiness」で閉じる。provider transport と Job Run UI は別 handoff に分けるため、分割必須にはしない。
  - `PERSONA` 周辺 schema 追加が必要な場合は job-scoped persona、evidence、phase link、snapshot summary に限定する。docs 正本化は含めない。
  - `本番経路`: usecase -> service -> repository / transactor -> SQLite -> body readiness query。

### `integration-persona-phase-wails-gateway`

- `implementation_target`: backend usecase、Wails controller、frontend Wails gateway / DTO を接続し、Job Run UI が実 backend の persona phase public seam を呼べる状態にする。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-persona-phase-public-seams`
  - `architecture_layer_basis`: Wails controller、bootstrap binding、frontend Wails adapter、gateway DTO の接続だけを統合境界として扱う。
  - `frozen_public_seams`: `GetPersonaGenerationPhaseSummary`, `StartPersonaGenerationPhase`, `PausePersonaGenerationPhase`, `ResumePersonaGenerationPhase`, `RetryPersonaGenerationPhase`, `CancelPersonaGenerationPhase`, `GetPersonaGenerationBodyReadiness`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, provider, model, execution mode, phase run ID, snapshot digest, prompt digest, count summary, error kind, retryable flag, body readiness
  - `secret_values_for_provider_external_api_internal_auth`: なし。Wails DTO と frontend gateway は secret 本体を受け取らない。
  - `secret_resolution_owner_layer`: backend secret store adapter と AI provider transport adapter。Wails controller は redacted DTO だけを返す。
  - `forbidden_outputs`: secret 本体、API key、token、raw prompt、raw response、原文発話全文、会話文脈全文を Wails DTO、frontend gateway DTO、console、UI、read model に出さない。
- `owned_scope`:
  - `internal/controller/wails/persona_generation_phase_controller.go`
  - `internal/controller/wails/persona_generation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/controller/wails/persona-generation-phase.gateway.ts`
  - `frontend/src/controller/wails/persona-generation-phase.gateway.test.ts`
  - `frontend/src/controller/wails/gateway-dto/persona-generation-phase/*`
  - `frontend/src/application/gateway-contract/persona-generation-phase/*`
- `depends_on`: `backend-persona-persistence-readiness-retry`, `frontend-persona-phase-job-run-ui`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `first_action`: `internal/controller/wails/persona_generation_phase_controller.go` に `StartPersonaGenerationPhase` の controller method と DTO mapping test を追加し、`completion_signal` の「Wails controller が persona phase usecase を呼び、redacted response を返す」を最初に閉じる。理由は bootstrap と frontend gateway mapping が controller method 名に依存するため。
- `validation_commands`:
  - `go test ./internal/controller/wails ./internal/bootstrap -run 'PersonaGeneration|JobRun'`
  - `npm --prefix frontend run test -- --run persona-generation-phase.gateway`
- `completion_signal`:
  - Wails controller が backend usecase の summary / command DTO を frontend gateway DTO へ lossless に写像する。
  - bootstrap と app controller に persona phase controller が接続される。
  - frontend Wails gateway が frozen gateway contract を満たす。
  - error kind、retryable flag、phase state、progress、target summary、AI execution summary、snapshot summary、body readiness が backend / frontend DTO で一致する。
  - secret、API key、raw payload、raw prompt、原文発話全文、会話文脈全文は DTO mapping で落とされる。
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
  - `reference_values_allowed_in_ui_dto_read_model`: final evidence に載せる redacted summary、test result、UI screenshot note、debug log redaction assertion result
  - `secret_values_for_provider_external_api_internal_auth`: なし。final validation は secret 本体を収集しない。
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: secret 本体、API key、token、raw provider payload、原文発話全文、会話文脈全文を work report、UI evidence、final validation summary に出さない。
- `owned_scope`:
  - `work_history/runs/YYYY-MM-DD-persona-generation-phase-run/codex.md`
  - 実装で必要になった task-local residual note
- `depends_on`: `integration-persona-phase-wails-gateway`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: なし
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `work_history/runs/YYYY-MM-DD-persona-generation-phase-run/codex.md` を作成し、final validation 欄と completion evidence 欄を先に固定する。理由は gate 結果、UI 証跡、環境 blocker を completion packet の根拠にするため。
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/persona-generation-phase/scenario-design.md --coverage docs/exec-plans/active/persona-generation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/persona-generation-phase/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `go test ./internal/...`
  - `npm run test:frontend`
  - `npm --prefix frontend run check`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - requirement gate と scenario gate が pass する。
  - relevant backend、integration、frontend tests が pass する。
  - Job Run の desktop / mobile UI 証跡で current phase、progress、phase result、target summary、AI execution summary、error summary、button state が重ならない。
  - UI、console、error summary、structured log、debug log、fake transport log に secret と API key がない。
  - UI と DB summary に provider raw request / response、full prompt、原文発話全文、会話文脈全文がない。
  - debug log では prompt / request body を確認でき、secret と API key が redacted されている。
  - paid API が呼ばれていない証跡を fake transport log または test result で確認できる。
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
- `task_id`: `persona-generation-phase`
- `source_scope`: `docs/exec-plans/active/persona-generation-phase/implementation-scope.md`
- `human_review_status`: approved
- `approval_record`: human message `approved`
- `ready_for_implementation`: true
- `forbidden_changes`: product docs 正本、`.codex/`、`.codex/skills`、`.codex/agents`
- `docs_changes`: none
- `implementation_start_decision`: implement_lane が判断する
- `agent_launch_decision`: implement_lane が判断する

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
- `docs_changes`: none

## Open Items

- なし。
