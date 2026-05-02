# Implementation Scope: term-translation-phase

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
- `source_task`: `tasks/usecases/term-translation-phase.yaml`
- `reference_docs`: `docs/spec.md`, `docs/er.md`, `docs/architecture.md`
- `vendor_api_reference`: `docs/references/vendor-api/xai_openapi_full.json`, `docs/references/vendor-api/gemini batch ref.md`

## Fixed Decisions

- `needs_human_decision`: `0`
- `scenario-gate`: pass
- Ready job だけが単語翻訳フェーズを開始できる。Ready 以外、terminal job、既存 active phase run は開始を拒否する。
- 共通辞書完全一致語は provider request へ送らず、phase 開始時 snapshot で除外判定を固定する。
- AI provider の用語翻訳結果は自動で確定訳語にする。人間確認 UI は含めない。
- 初期実装は 1 対象語 1 request unit とする。Batch API を使う場合も batch item は 1 対象語単位にする。
- 再開、リトライ、再実行では既存 entry を維持し、未処理 term だけ進める。
- secret、API key 平文、provider raw request / response、翻訳フィールド本文の全文は UI、error summary、structured log に出さない。
- backend、frontend、統合境界は別 handoff に分ける。frontend handoff は確定済み `contract_freeze` に依存する。
- `APIテスト` は public seam 起点の system-level test とする。`UI人間操作E2E` は最終検証で証明する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `contract-term-phase-public-seams` | なし | なし | なし |
| `wave-2` | `backend-term-phase-state-dictionary`, `backend-term-provider-adapter`, `frontend-term-phase-job-run-ui` | `contract-term-phase-public-seams` | `backend-term-phase-state-dictionary <-> backend-term-provider-adapter`, `backend-term-provider-adapter <-> frontend-term-phase-job-run-ui`, `backend-term-phase-state-dictionary <-> frontend-term-phase-job-run-ui` | なし |
| `wave-3` | `integration-term-phase-wails-gateway` | `backend-term-phase-state-dictionary`, `backend-term-provider-adapter`, `frontend-term-phase-job-run-ui` | なし | `depends_on` |
| `wave-4` | `final-validation-and-report` | `integration-term-phase-wails-gateway` | なし | `broad_gate_shared` |

## Handoffs

### `contract-term-phase-public-seams`

- `implementation_target`: Job Run から単語翻訳フェーズを開始、再開、リトライ、結果取得、後続 phase 開始可否確認を行う public seam と DTO を固定する。
- `implementation_artifact`: `contract_freeze`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-TTP-001` から `SCN-TTP-009`、`./ui-design.md` の `UI Contract`
  - `architecture_layer_basis`: Wails controller / DTO は backend controller 境界、usecase contract は backend usecase 境界、frontend gateway contract は frontend contract 境界に置く。依存方向は frontend gateway -> Wails adapter -> Wails controller -> usecase とする。
  - `frozen_public_seams`:
    - `GetTermTranslationPhaseSummary`: job ID を受け取り、current phase、phase state、progress、対象語件数、共通辞書 hit 件数、AI 実行対象語件数、result summary、error summary、button enablement を返す。
    - `StartTermTranslationPhase`: Ready job を受け取り、term phase run を作成または既存 active run を拒否し、phase run ID、state、progress、開始不可理由を返す。
    - `PauseTermTranslationPhase`: active term phase run を中断し、phase state と retry / resume 可否を返す。
    - `ResumeTermTranslationPhase`: paused または recoverable failed の同じ phase run を再開し、既存 entry を維持した progress を返す。
    - `RetryTermTranslationPhase`: retryable failure の同じ phase run で未処理 term だけを再実行し、latest error と progress を更新する。
    - `GetTermTranslationNextPhaseReadiness`: term phase 完了と辞書参照成立後だけ後続 phase 可を返す。
    - error kind は `ready_required`, `terminal_job`, `active_phase_exists`, `dictionary_snapshot_failed`, `provider_failure`, `invalid_provider_response`, `save_failed`, `term_phase_incomplete`, `secret_redacted` を区別する。
    - DTO は `credential_ref`、provider、model、execution mode、snapshot digest / version を要約として持ち、secret と raw payload を持たない。
- `owned_scope`:
  - `internal/usecase/term_translation_phase_contract.go`
  - `internal/controller/wails/term_translation_phase_controller.go`
  - `internal/controller/wails/term_translation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/application/gateway-contract/term-translation-phase/*`
  - `frontend/src/controller/wails/gateway-dto/term-translation-phase/*`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`: `internal/usecase/term_translation_phase_contract.go` に summary / command DTO と public error kind を追加し、`completion_signal` の「downstream が参照する field 名、nullability、error kind が固定される」を最初に閉じる。理由は backend core、frontend UI、Wails gateway が同じ seam に依存するため。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/controller/wails -run 'TermTranslation|JobRun|SCN_TTP_001|SCN_TTP_005'`
  - `npm --prefix frontend run test -- --run term-translation-contract`
- `completion_signal`:
  - Job Run が参照する term phase public seam 名、request / response DTO、error kind が存在する。
  - DTO は secret 平文、API key、provider raw request / response、翻訳フィールド本文全文を表現できない。
  - field 名、nullability、phase state、retryable flag、後続 phase 可否、progress summary が controller unit test と frontend contract test で固定される。
  - frontend handoff が参照できる gateway contract が作成される。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `8-12 files`、`400-700 changed lines`。
  - contract freeze のみを扱い、永続化実体、provider adapter 実体、Job Run UI 実装は含めない。
  - `本番経路`: Wails controller / DTO -> usecase contract -> frontend gateway contract。

### `backend-term-phase-state-dictionary`

- `implementation_target`: Ready job から term phase run を開始し、共通辞書 snapshot、完全一致除外、ジョブ内辞書保存、phase state、後続 phase guard、再開 / リトライの冪等性を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-term-phase-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、controller entry までを backend 境界として扱う。frontend UI は含めない。
  - `frozen_public_seams`: `contract-term-phase-public-seams` の completion signal を参照する。
- `owned_scope`:
  - `internal/usecase/term_translation_phase_*`
  - `internal/service/term_translation_phase_*`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/job_output_repository.go`
  - `internal/repository/job_output_sqlite_repository.go`
  - `internal/repository/foundation_data_repository.go`
  - `internal/repository/foundation_data_sqlite_repository.go`
  - `internal/repository/transactor*`
  - `internal/infra/sqlite/dbinit/migrations/*term_translation*`
  - `internal/integrationtest/*term_translation*`
- `depends_on`: `contract-term-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-term-provider-adapter`, `frontend-term-phase-job-run-ui`
- `parallel_blockers`: なし
- `first_action`: `internal/service/term_translation_phase_service.go` に Ready job / active phase run の開始判定と focused test を追加し、`completion_signal` の「Ready job だけが term phase run を開始できる」を最初に閉じる。理由は辞書反映、再開、後続 phase guard が同じ phase state に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest -run 'TermTranslation|PhaseRun|Dictionary|SCN_TTP_(001|002|004|007|008)'`
- `completion_signal`:
  - Ready job だけが単語翻訳フェーズを開始でき、Ready 以外、terminal job、既存 active phase run は拒否される。
  - 共通辞書 snapshot が phase 開始時に固定され、完全一致語だけが provider request 対象から除外される。
  - 共通辞書除外後に対象語が 0 件でも Completed phase result になり、provider 未実行が result summary に残る。
  - 確定訳語は対象 job の `DICTIONARY_ENTRY` と `PHASE_RUN_DICTIONARY_ENTRY` に atomic に反映される。
  - 同一 job、同一 record type、同一 source term は重複 entry を作らず、別 record type は別 entry として扱える。
  - 保存失敗では partial dictionary state を成功扱いにしない。
  - 再開、リトライ、開始再送では同じ `JOB_PHASE_RUN` を再利用し、既存 entry を維持して未処理 term だけ進める。
  - term phase 未完了、失敗、辞書参照不能では後続 phase run を作成しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-24 files`、`900-1400 changed lines`。
  - 1 受け入れユースケースは「term phase の状態遷移と job dictionary 反映」で閉じる。provider HTTP 実装と frontend UI は別 handoff に分けるため、分割必須にはしない。
  - schema 追加が必要な場合は job-local dictionary の dedup index と phase result に限定する。docs 正本化は含めない。
  - `本番経路`: usecase -> service -> repository / transactor -> SQLite。

### `backend-term-provider-adapter`

- `implementation_target`: 1 対象語 1 request unit の provider adapter、fake transport、response validation、provider failure handling、redaction / audit summary を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-term-phase-public-seams`
  - `architecture_layer_basis`: AIProvider / response adapter boundary は backend infra / provider 境界に置き、service は provider-agnostic contract だけを見る。
  - `frozen_public_seams`: provider adapter output は source term correlation、translated term、retryable failure、redacted error summary を返す。
- `owned_scope`:
  - `internal/infra/ai/provider.go`
  - `internal/infra/ai/provider_client.go`
  - `internal/infra/ai/gemini.go`
  - `internal/infra/ai/openai_compatible.go`
  - `internal/infra/ai/transport.go`
  - `internal/infra/ai/*term_translation*`
  - `internal/service/term_translation_provider_*`
  - `internal/service/*term_translation*_test.go`
  - `internal/infra/ai/*term_translation*_test.go`
- `depends_on`: `contract-term-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-term-phase-state-dictionary`, `frontend-term-phase-job-run-ui`
- `parallel_blockers`: なし
- `first_action`: `internal/service/term_translation_provider_adapter.go` に 1 対象語 request unit の provider-agnostic input / output contract と invalid response test を追加し、`completion_signal` の「provider 応答が source term と translated term の対応を保持する」を最初に閉じる。理由は Gemini / xAI の実 adapter と fake transport が同じ response correlation に依存するため。
- `validation_commands`:
  - `go test ./internal/infra/ai ./internal/service -run 'TermTranslationProvider|ProviderAdapter|SCN_TTP_(003|006|009)'`
- `completion_signal`:
  - valid response は source term と translated term の対応を保持し、自動で確定訳語候補になる。
  - Batch API を使う場合も batch item は 1 対象語単位として扱われる。
  - response 欠落、余分な応答、空訳語、invalid shape は対象語単位の failed / retryable として扱われる。
  - provider failure で別 provider へ暗黙 fallback しない。
  - fake transport で paid API を呼ばずに provider path を検証できる。
  - API key、secret、raw request / response、本文全文は request log、error summary、structured log に出ない。
  - provider、model、execution mode、input count、output count、prompt version または digest を audit summary で確認できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は normal。想定 `8-14 files`、`450-800 changed lines`。
  - vendor API 資料は request / response shape と batch result mapping の参照だけに使う。paid API 実行は検証条件にしない。
  - `本番経路`: service provider adapter -> internal/infra/ai provider client -> HTTP transport。

### `frontend-term-phase-job-run-ui`

- `implementation_target`: Job Run に単語翻訳フェーズの current phase、progress、phase result、error summary、開始 / 中断 / 再開 / リトライ / 後続 phase 操作を表示する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-term-phase-public-seams`
  - `architecture_layer_basis`: frontend contract、store / presenter / usecase / controller、Job Run screen を frontend 境界として扱う。Wails 接続実体は統合境界 handoff に分ける。
  - `frozen_public_seams`: `contract-term-phase-public-seams` の frontend gateway contract。
- `owned_scope`:
  - `frontend/src/application/contract/term-translation-phase/*`
  - `frontend/src/application/store/term-translation-phase/*`
  - `frontend/src/application/presenter/term-translation-phase/*`
  - `frontend/src/application/usecase/term-translation-phase/*`
  - `frontend/src/controller/term-translation-phase/*`
  - `frontend/src/ui/screens/term-translation-phase/*`
  - `frontend/src/ui/screens/job-run/*`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
- `depends_on`: `contract-term-phase-public-seams`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-term-phase-state-dictionary`, `backend-term-provider-adapter`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/contract/term-translation-phase/term-translation-phase-screen-contract.ts` に Job Run 表示状態と action enablement contract を追加し、`completion_signal` の「phase state、progress、result summary、button enablement を型で表せる」を最初に閉じる。理由は store、presenter、UI が同じ表示状態に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run term-translation-phase`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - Job Run で current phase、phase state、progress、開始時刻、完了時刻、対象語件数、共通辞書 hit 件数、AI 実行対象語件数を表示できる。
  - phase result に確定訳語件数、ジョブ内辞書反映件数、置換対象件数、未一致件数、provider / model / execution mode の要約を表示できる。
  - 開始は Ready job かつ active term phase run がない時だけ有効になる。
  - 後続 phase へ進む操作は term phase 完了と辞書参照成立後だけ有効になる。
  - retry は retryable failure の時だけ有効になる。
  - loading、empty completed、completed、paused、recoverable failed、blocked の状態差分を表示できる。
  - secret、API key 平文、provider raw request / response、翻訳フィールド本文全文を画面、error summary、console に出さない。
  - 長い source term、provider 名、model 名、error reason は desktop / mobile で overflow しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`850-1400 changed lines`。
  - UI人間操作E2E の最終証明は `final-validation-and-report` に寄せる。この handoff の local validation は mocked gateway の frontend tests と type / check に限定する。
  - Wails gateway 実体は `integration-term-phase-wails-gateway` に分ける。
  - `本番経路`: frontend gateway contract -> usecase -> store / presenter -> Job Run screen。

### `integration-term-phase-wails-gateway`

- `implementation_target`: backend usecase、Wails controller、frontend Wails gateway / DTO を接続し、Job Run UI が実 backend の term phase public seam を呼べる状態にする。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-term-phase-public-seams`
  - `architecture_layer_basis`: Wails controller、bootstrap binding、frontend Wails adapter、gateway DTO の接続だけを統合境界として扱う。
  - `frozen_public_seams`: `GetTermTranslationPhaseSummary`, `StartTermTranslationPhase`, `PauseTermTranslationPhase`, `ResumeTermTranslationPhase`, `RetryTermTranslationPhase`, `GetTermTranslationNextPhaseReadiness`
- `owned_scope`:
  - `internal/controller/wails/term_translation_phase_controller.go`
  - `internal/controller/wails/term_translation_phase_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts`
  - `frontend/src/controller/wails/term-translation-phase.gateway.test.ts`
  - `frontend/src/controller/wails/gateway-dto/term-translation-phase/*`
  - `frontend/src/application/gateway-contract/term-translation-phase/*`
- `depends_on`: `backend-term-phase-state-dictionary`, `backend-term-provider-adapter`, `frontend-term-phase-job-run-ui`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `first_action`: `internal/controller/wails/term_translation_phase_controller.go` に `StartTermTranslationPhase` の controller method と DTO mapping test を追加し、`completion_signal` の「Wails controller が term phase usecase を呼び、redacted response を返す」を最初に閉じる。理由は bootstrap と frontend gateway mapping が controller method 名に依存するため。
- `validation_commands`:
  - `go test ./internal/controller/wails ./internal/bootstrap -run 'TermTranslation|JobRun'`
  - `npm --prefix frontend run test -- --run term-translation-phase.gateway`
- `completion_signal`:
  - Wails controller が backend usecase の summary / command DTO を frontend gateway DTO へ lossless に写像する。
  - bootstrap と app controller に term phase controller が接続される。
  - frontend Wails gateway が frozen gateway contract を満たす。
  - error kind、retryable flag、phase state、progress、dictionary summary、next phase readiness が backend / frontend DTO で一致する。
  - secret、raw payload、翻訳フィールド本文全文は DTO mapping で落とされる。
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
- `owned_scope`:
  - `work_history/runs/YYYY-MM-DD-term-translation-phase-run/codex.md`
  - 実装で必要になった task-local residual note
- `depends_on`: `integration-term-phase-wails-gateway`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `work_history/runs/YYYY-MM-DD-term-translation-phase-run/codex.md` を作成し、final validation 欄と completion evidence 欄を先に固定する。理由は gate 結果、UI 証跡、環境 blocker を completion packet の根拠にするため。
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/term-translation-phase/scenario-design.md --coverage docs/exec-plans/active/term-translation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/term-translation-phase/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `go test ./internal/...`
  - `npm run test:frontend`
  - `npm --prefix frontend run check`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - scenario gate が pass する。
  - relevant backend、integration、frontend tests が pass する。
  - Job Run の desktop / mobile UI 証跡で current phase、progress、phase result、error summary、button state が重ならない。
  - UI、console、error summary、structured log、fake transport log に secret、raw response、本文全文がない。
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
- `task_id`: `term-translation-phase`
- `source_scope`: `docs/exec-plans/active/term-translation-phase/implementation-scope.md`
- `human_review_status`: approved
- `approval_record`: human message `approved`
- `ready_for_implementation`: true
- `start_wave`: `wave-1`
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
- `docs_changes`: none

## Open Items

- なし。
