# Implementation Scope: translation-job-setup

- `skill`: implementation-scope
- `status`: ready-for-copilot
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human message `approvedステータスにする`
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`

## Fixed Decisions

- `needs_human_decision`: `0`
- Draft は UI 未保存状態であり、DB job は validation pass 後の create で初めて作る。
- 同一入力への 2 件目 job 作成は禁止する。過去 job の廃棄手段は別 task で扱う。
- 必須設定不足、参照不能、provider / mode 不整合、credential 参照不能は blocking validation failure にする。
- validation 実行ごとの履歴は business table ではなく structured log に残す。アプリ状態は直近 validation 結果、対象設定断面、失敗カテゴリ、job 作成時の pass 断面だけ保持する。
- credential 解決、provider capability、ネットワーク到達性はすべて blocking にする。fake provider は user-facing provider list に出さず、test は外部 request / SDK transport だけを fake に差し替える。
- cache 欠落時は Job Setup をブロックし、Input Review の再構築導線へ戻す。
- Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集は許可しない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `backend-job-setup-contract-freeze` | `なし` | `なし` | `なし` |
| `wave-2` | `backend-job-setup-core`, `frontend-job-setup-ui` | `backend-job-setup-contract-freeze` | `backend-job-setup-core <-> frontend-job-setup-ui` | `なし` |
| `wave-3` | `final-validation-and-report` | `backend-job-setup-core`, `frontend-job-setup-ui` | `なし` | `broad_gate_shared` |

## Handoffs

### `backend-job-setup-contract-freeze`

- `implementation_target`: Job Setup の Wails public seam、request / response DTO、error kind、validation result shape を固定する。
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-TJS-001` から `SCN-TJS-007`、`./ui-design.md` の `UI Contract`
  - `frozen_public_seams`:
    - `GetTranslationJobSetupOptions`: 入力候補、既存 job 状態、共通辞書、共通ペルソナ、AI runtime 候補、credential 参照状態を返す。
    - `ValidateTranslationJobSetup`: 入力、共通基盤、AI runtime、実行方式を受け取り、pass / fail / warning、blocking failure category、対象設定断面、validation timestamp、create 可否を返す。
    - `CreateTranslationJob`: validation pass 済み断面から `Ready` job を作成し、job ID、状態、入力出自、実行設定要約、validation pass 断面を返す。
    - `GetTranslationJobSetupSummary`: 作成済み job を read-only で再表示するため、job ID、状態、入力出自、設定要約、phase 開始可否を返す。
    - error kind は `required_setting_missing`, `input_not_found`, `cache_missing`, `foundation_ref_missing`, `credential_missing`, `provider_mode_unsupported`, `provider_unreachable`, `duplicate_job_for_input`, `validation_stale`, `partial_create_failed`, `ready_required` を区別する。
- `owned_scope`:
  - `internal/usecase/translation_job_setup_contract.go`
  - `internal/controller/wails/translation_job_setup_controller.go`
  - `internal/controller/wails/translation_job_setup_controller_unit_test.go`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `internal/usecase/translation_job_setup_contract.go` に request / response / error kind を追加し、`completion_signal` の「frontend が依存する public seam 名と DTO shape が固定される」を最初に閉じる。理由は backend core と frontend UI の両方が同じ seam に依存するため。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/controller/wails -run 'TranslationJobSetup|JobSetup'`
- `completion_signal`:
  - Wails controller に Job Setup 用 method 名と DTO shape が存在する。
  - DTO は secret 平文を含まず、credential は参照状態だけを表す。
  - validation response は pass / fail / warning、blocking failure category、対象設定断面、create 可否を表せる。
  - create response は `Ready` job、入力出自、実行設定要約、validation pass 断面を表せる。
  - frontend handoff が参照できる field 名、nullability、error kind が controller unit test で固定される。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は normal。想定 `6-10 files`、`300-600 changed lines`。
  - この handoff は contract freeze のみを扱い、永続化、provider validation 実体、frontend UI は含めない。
  - `本番経路`: Wails controller / DTO -> usecase contract。

### `backend-job-setup-core`

- `implementation_target`: Job Setup validation と `Ready` job 作成を backend で実装し、atomic create、重複禁止、cache 欠落、provider validation fake transport を検証する。
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `backend-job-setup-contract-freeze`
  - `frozen_public_seams`: `backend-job-setup-contract-freeze` の completion signal を参照する。
- `owned_scope`:
  - `internal/usecase/translation_job_setup_*`
  - `internal/service/translation_job_setup_*`
  - `internal/controller/wails/translation_job_setup_controller*`
  - `internal/repository/job_lifecycle_repository.go`
  - `internal/repository/job_lifecycle_sqlite_repository.go`
  - `internal/repository/translation_source_repository.go`
  - `internal/repository/translation_source_sqlite_repository.go`
  - `internal/repository/foundation_data_repository.go`
  - `internal/repository/foundation_data_sqlite_repository.go`
  - `internal/repository/master_persona_repository.go`
  - `internal/infra/ai/provider*`
  - `internal/infra/ai/transport*`
  - `internal/bootstrap/app_controller.go`
  - `internal/integrationtest/*translation_job_setup*`
- `depends_on`: `backend-job-setup-contract-freeze`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `frontend-job-setup-ui`
- `parallel_blockers`: `なし`
- `first_action`: `internal/service/translation_job_setup_service.go` に validation service skeleton と focused test を追加し、`completion_signal` の「blocking validation failure が create を禁止する」を最初に閉じる。理由は create、UI state、provider fake transport がすべて validation result に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/integrationtest -run 'TranslationJobSetup|JobSetup|SCN_TJS|JobLifecycle'`
- `completion_signal`:
  - 取り込み済み input と valid foundation / AI runtime から validation pass を作れる。
  - validation pass かつ未失効の断面だけ `Ready` job 作成を許可する。
  - `TRANSLATION_JOB` は 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照し、同一入力への 2 件目 job 作成を拒否する。
  - 必須設定不足、参照不能、provider / mode 不整合、credential 参照不能、cache 欠落は blocking failure として返る。
  - cache 欠落時は Job Setup 内で再構築せず、Input Review の再構築導線へ戻せる error kind を返す。
  - provider list は real provider を返し、test は external request / SDK transport だけを fake に差し替えて paid API を呼ばない。
  - create 途中の保存失敗は transaction 全体を rollback し、partial `TRANSLATION_JOB` や欠けた `JOB_PHASE_RUN` を残さない。
  - validation 実行ごとの履歴は business table に保存せず、structured log と直近アプリ状態だけで観測する。
  - Ready 未成立の job または setup から Running / phase run を開始できない error surface を持つ。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。
  - validation、create、provider fake transport は 1 つの Job Setup backend use case に閉じるため 1 handoff にする。
  - frontend UI は含めない。
  - `本番経路`: Wails controller / DTO -> usecase -> service -> repository / AI transport -> SQLite。
  - 過去 job の廃棄手段、phase 実行中の common foundation lock、翻訳 phase 実行は含めない。

### `frontend-job-setup-ui`

- `implementation_target`: Translation Management 内に Job Setup UI を追加し、入力、共通基盤、AI runtime、validation、create result を 1 画面で操作できるようにする。
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `backend-job-setup-contract-freeze`
  - `frozen_public_seams`: Job Setup Wails DTO と error kind に合わせた frontend gateway contract。
- `owned_scope`:
  - `frontend/src/application/gateway-contract/translation-job-setup/*`
  - `frontend/src/application/contract/translation-job-setup/*`
  - `frontend/src/application/store/translation-job-setup/*`
  - `frontend/src/application/presenter/translation-job-setup/*`
  - `frontend/src/application/usecase/translation-job-setup/*`
  - `frontend/src/controller/translation-job-setup/*`
  - `frontend/src/controller/wails/gateway-dto/translation-job-setup/*`
  - `frontend/src/controller/wails/translation-job-setup.gateway*`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
- `depends_on`: `backend-job-setup-contract-freeze`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `backend-job-setup-core`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts` に frozen DTO と同じ gateway contract を追加し、`completion_signal` の「frontend が validation / create / summary response を型で受け取れる」を最初に閉じる。理由は store、presenter、UI が同じ contract に依存するため。
- `validation_commands`:
  - `npm --prefix frontend run test -- --run translation-job-setup`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - Translation Management から Job Setup を開ける。
  - 入力データ名、出自、登録日時、翻訳レコード件数、既存 job 状態を表示できる。
  - 共通辞書、共通ペルソナ、provider、model、credential 参照状態、実行方式を選択または確認できる。
  - validation pass / fail / warning、dirty state、失敗理由、対象設定断面を表示できる。
  - blocking failure 中、validation 未実行、validation 失効、同一入力の既存 job ありでは create job を無効にする。
  - cache missing は Job Setup 内で再構築せず、Input Review の再構築導線を表示する。
  - API key 平文、secret 本体、復号可能な値を画面、error summary、console へ表示しない。
  - create 成功後は `Ready` job、入力出自、設定要約、validation pass 断面を read-only で表示し、再編集 action を出さない。
  - 長い file path、plugin 名、provider / model 名、foundation 名、failure reason は desktop / mobile で overflow しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。
  - UI人間操作E2E の最終証明は `final-validation-and-report` に寄せる。この handoff の local validation は mocked gateway の frontend tests と type / check に限定する。
  - backend core と同時に実行できるが、frozen contract から外れる field 変更が必要になった場合は replan ではなく `backend-job-setup-contract-freeze` の修正完了を待つ。
  - `本番経路`: Wails gateway -> frontend usecase -> store / presenter -> Job Setup screen。

### `final-validation-and-report`

- `implementation_target`: 全 handoff 完了後に scenario、broad gate、UI 証跡、Copilot report をまとめて確認する。
- `contract_freeze`:
  - `status`: `not_required`
  - `freeze_source`: `N/A`
  - `frozen_public_seams`: `N/A`
- `owned_scope`:
  - `work_history/runs/YYYY-MM-DD-translation-job-setup-run/copilot.md`
  - 実装で必要になった task-local residual note
- `depends_on`: `backend-job-setup-core`, `frontend-job-setup-ui`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `work_history/runs/YYYY-MM-DD-translation-job-setup-run/copilot.md` を作成し、実行する final validation 欄を先に固定する。理由は report path と gate 結果を completion packet の正本にするため。
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `go test ./internal/...`
  - `npm run test:frontend`
  - `npm --prefix frontend run check`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - scenario gate が pass する。
  - backend と frontend の relevant test が pass する。
  - Job Setup の desktop / mobile UI 証跡で selector、validation summary、create result、error reason が重ならない。
  - paid API が呼ばれていない証跡を test result または fake transport log で確認できる。
  - system / harness が環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason、再実行環境、再実行コマンドを report に残す。
  - Copilot work report が completion packet の schema を満たす。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `notes`:
  - product 実装はここで追加しない。
  - broad validation owner は全 handoff 完了後だけに置く。
  - Sonar を使う場合は `/tmp` cache を使い、Sonar server の Quality Gate と repo-local issue gate を混同しない。

## Human Copilot Handoff Packet

- `entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `task_id`: `translation-job-setup`
- `scope_source`: `docs/exec-plans/active/translation-job-setup/implementation-scope.md`
- `start_wave`: `wave-1`
- `do_not_change`:
  - `docs/`
  - `.codex/`
  - `.github/skills`
  - `.github/agents`
- `do_not_implement`:
  - 過去 job の廃棄手段
  - 共通辞書と共通ペルソナの管理 UI
  - 翻訳 phase 実行
  - 訳文生成
  - 成果物出力
  - phase 実行中の common foundation lock
  - fake provider の user-facing provider list 追加
- `must_return`:
  - `completed_handoffs`
  - `touched_files`
  - `implemented_scope`
  - `test_results`
  - `ui_evidence`
  - `pre_review_gate_result`
  - `implementation_review_result`
  - `harness_gate_result`
  - `completion_evidence`
  - `telemetry_events`

## Completion Packet

Copilot は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `pre_review_gate_result`
- `implementation_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: Codex 側 `work_reporter` が読む実装事実。report 文面ではなく、completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: copilot` の `assistant_response` event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`
