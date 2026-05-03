# Implementation Scope: translation-output-artifact

- `skill`: implementation-scope
- `status`: ready-for-implementation
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human instruction `approved`
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
- `source_task`: `tasks/usecases/translation-output-artifact.yaml`
- `upstream_scope`: `docs/exec-plans/completed/body-translation-phase/implementation-scope.md`
- `reference_docs`: `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/screen-design/README.md`
- `code_map`: `tmp/code-map/index.json`

## Fixed Decisions

- 人間レビュー承認を確認した。承認済み design bundle として `scenario-design.md`、`ui-design.md`、coverage / gate 成果物を扱う。
- `needs_human_decision`: `0`。`scenario-design.md`、`scenario-design.candidate-coverage.json`、`scenario-design.requirement-coverage.json`、`scenario-design.questions.md` を確認した。
- `unresolved_conflicts`: `0`。`scenario-design.requirement-gate.md` は `pass`、`finding_count: 0`、`question_count: 0`。
- `coverage_note`: `scenario-design.md` 本文の adopted / merged 集計と `scenario-design.candidate-coverage.json` の集計に 1 件差がある。採用または統合済みの区分差であり、`needs_human_decision: 0` と `unresolved_conflicts: 0` は一致するため、実装範囲は scenario matrix と coverage JSON の未決 0 件を根拠に固定する。
- Output Review は body phase Completed かつ job-level `Completed` の job だけを出力候補にする。
- xTranslator 互換 XML は `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を row 単位で再構成する。
- 内部 `cached` は xTranslator `Status=1` へ写像し、辞書置換の内部観測情報は xTranslator `Status` と分離する。
- 再出力は同じ job に重複 artifact または重複 row を作らない。現行 ER の `TRANSLATION_ARTIFACT.translation_job_id UNIQUE` と `XTRANSLATOR_OUTPUT_ROW.job_translation_field_id UNIQUE` を前提にする。
- 出力処理は provider、network、secret store を必須経路にしない。secret、API key、provider raw payload、過剰な本文全文を UI、DTO、summary、log へ出さない。
- frontend は `contract-output-artifact-public-seams` の `contract_freeze` 完了後にだけ開く。
- `APIテスト` は public seam 起点の system-level test とする。`UI人間操作E2E` は最終検証で証明する。
- docs 正本化、プロダクト仕様変更、`.codex` 変更は Codex implementation レーンへ渡さない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `contract-output-artifact-public-seams` | なし | なし | なし |
| `wave-2` | `backend-output-review-readiness`, `backend-output-row-builder-compatibility`, `frontend-output-review-ui` | `contract-output-artifact-public-seams` | `backend-output-review-readiness <-> frontend-output-review-ui`, `backend-output-row-builder-compatibility <-> frontend-output-review-ui` | backend 同士は `承認済み実装範囲重複` |
| `wave-3` | `backend-output-artifact-write-reoutput` | `backend-output-review-readiness`, `backend-output-row-builder-compatibility` | なし | `依存対象` |
| `wave-4` | `integration-output-artifact-wails-gateway` | `backend-output-review-readiness`, `backend-output-row-builder-compatibility`, `backend-output-artifact-write-reoutput`, `frontend-output-review-ui` | なし | `依存対象` |
| `wave-5` | `final-validation-and-report` | `integration-output-artifact-wails-gateway` | なし | `広域判定条件共有` |

## Handoffs

### `contract-output-artifact-public-seams`

- `承認済み実装範囲`: Output Review の query / command、DTO、gateway contract、controller entry、frontend state contract を固定する。
- `implementation_artifact`: `contract_freeze`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `required`
  - `freeze_source`: `./scenario-design.md` の `SCN-TOA-001` から `SCN-TOA-010`、`./ui-design.md` の `UI Contract`
  - `architecture_layer_basis`: Wails controller / DTO は backend controller 境界、usecase contract は backend usecase 境界、frontend gateway contract は frontend contract 境界に置く。依存方向は frontend gateway -> Wails adapter -> Wails controller -> usecase とする。
  - `frozen_public_seams`:
    - `GetTranslationOutputReview`: completed job list、selected job summary、input provenance、output readiness、拒否理由、result summary、artifact status を返す。
    - `GetTranslationOutputDiffPreview`: selected job と artifact を受け取り、Source、Dest、Status、row reflection summary、stale reason、再出力可否を返す。
    - `GenerateXTranslatorOutputArtifact`: completed job ID、target game、output path を受け取り、artifact status、row count、file path、compatibility summary、redacted error summary を返す。
    - `RegenerateXTranslatorOutputArtifact`: existing artifact を持つ completed job を受け取り、同じ job の artifact / row を置換または更新した summary を返す。
    - error kind は `not_completed`, `canceled`, `status_mismatch`, `missing_required_row_field`, `unknown_output_status`, `xml_serialization_failed`, `file_write_failed`, `artifact_save_failed`, `compatibility_rejected`, `secret_redacted` を区別する。
    - DTO は job ID、artifact ID、field ID、row digest、count、status、file path、target game、error kind、retryable flag を持ち、secret 本体、API key、provider raw payload、過剰な本文全文を持たない。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: job ID, artifact ID, field ID, row digest, artifact status, file path, target game, count, error kind, retryable flag
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, secret store value
  - `secret_resolution_owner_layer`: 出力処理では secret を解決しない。provider 履歴を読む場合も secret store adapter を呼ばない。
  - `forbidden_outputs`: secret 本体、API key、token、authorization header、provider raw request / response、復号可能値、過剰な本文全文を UI、DTO、read model、URL、error summary、structured log、debug log に出さない。
- `依存対象`: なし
- `検証`: `go test ./internal/usecase ./internal/controller/wails -run 'TranslationOutput|OutputArtifact|SCN_TOA_(001|002|003|006|010)'`、`npm --prefix frontend run test -- --run translation-output-contract`
- `初手`: `internal/usecase/translation_output_artifact_contract.go` を追加し、公開 query / command request / response、error kind、redaction field obligation を固定する。対応する完了条件は「downstream が参照する field 名、nullability、error kind が固定される」である。理由は backend、frontend、Wails gateway が同じ public seam に依存するため。
- `実行グループ`: `wave-1`
- `ready_wave`: `wave-1`
- `並列可能対象`: なし
- `並列不可理由`: なし
- `完了条件`:
  - Output Review、diff preview、generate、regenerate の public seam 名と request / response DTO が存在する。
  - field 名、nullability、error kind、retryable flag、redaction obligation が backend controller unit test と frontend contract test で固定される。
  - DTO は secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を表現できない。
  - frontend handoff が参照できる gateway contract と screen state contract が作成される。
- `想定変更ファイル数`: `8-14 files`
- `想定変更行数`: `450-800 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `補足`: contract freeze のみを扱う。永続化実体、XML serializer 実体、Output Review UI 実装は含めない。`本番経路`: Wails controller / DTO -> usecase contract -> frontend gateway contract。

### `backend-output-review-readiness`

- `承認済み実装範囲`: completed job list、selected job summary、output readiness、拒否理由、result summary、input provenance、artifact status の backend read model を実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-output-artifact-public-seams`
  - `architecture_layer_basis`: repository / SQLite concrete、service / usecase、controller entry までを backend 境界として扱う。frontend UI と Wails gateway 実体は含めない。
  - `frozen_public_seams`: `contract-output-artifact-public-seams` の completion signal を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: job ID, artifact ID, field ID, input provenance digest, output readiness, rejection reason, status distribution, row count, error kind
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。provider 履歴を読む場合も参照値と redacted summary だけを扱う。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を read model、error summary、structured log、debug log に出さない。
- `依存対象`: `contract-output-artifact-public-seams`
- `検証`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/apitest -run 'TranslationOutput|OutputReadiness|SCN_TOA_(001|002|008|010)'`
- `初手`: `internal/service/translation_output_artifact_service.go` を追加し、body phase Completed / job-level `Completed` / status mismatch / `Canceled` の readiness 判定を focused test で固定する。対応する完了条件は「出力候補と拒否理由が同じ job ID に対して返る」である。理由は Output Review 表示、generate command、button enablement が同じ readiness に依存するため。
- `実行グループ`: `wave-2`
- `ready_wave`: `wave-2`
- `並列可能対象`: `frontend-output-review-ui`
- `並列不可理由`: `backend-output-row-builder-compatibility` とは `承認済み実装範囲重複`
- `完了条件`:
  - body phase Completed かつ job-level `Completed` の job だけが出力開始可能として返る。
  - 未完了、失敗中、`Canceled`、field result 不整合、status 不整合では readiness false と拒否理由が返る。
  - target count 0 は output readiness false の理由にならず、row count 0 の result summary を返せる。
  - completed job list、input provenance、訳文件数、output status distribution、artifact status が同じ job ID に対応する。
  - secret、API key、provider raw payload、過剰な本文全文を summary と log に出さない。
- `想定変更ファイル数`: `10-15 files`
- `想定変更行数`: `500-850 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `補足`: 想定規模は normal から caution 境界である。1 受け入れユースケースは「出力候補の確認と開始可否判定」で閉じる。row build、XML write、frontend UI は別 handoff に分けるため、分割必須にはしない。`本番経路`: usecase -> service -> repository / transactor -> SQLite。

### `backend-output-row-builder-compatibility`

- `承認済み実装範囲`: output-ready field から xTranslator row 必須列、status mapping、row validation、compatibility summary、diff preview 用 row digest を生成する backend rule を実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-output-artifact-public-seams`
  - `architecture_layer_basis`: row builder、compatibility validator、diff preview read rule は backend service 境界に置く。filesystem concrete と frontend UI は含めない。
  - `frozen_public_seams`: row required columns、one field one row、`cached` -> `Status=1`、unknown status rejection、compatibility summary。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: field ID, row digest, EDID, REC, FIELD, FORMID, Source excerpt, Dest excerpt, Status, compatibility kind, severity, affected row count
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: secret 解決は行わない。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な Source / Dest 全文を validator summary、diff summary、structured log、debug log に出さない。
- `依存対象`: `contract-output-artifact-public-seams`
- `検証`: `go test ./internal/service ./internal/usecase -run 'XTranslator|OutputRow|Compatibility|DiffPreview|SCN_TOA_(003|005|009|010)'`
- `初手`: `internal/service/xtranslator_output_row_builder.go` を追加し、`cached` field を xTranslator row `Status=1` に写像する focused test を固定する。対応する完了条件は「row 必須列と status mapping が成功 row に入る」である。理由は XML serializer、diff preview、artifact persistence が同じ row set に依存するため。
- `実行グループ`: `wave-2`
- `ready_wave`: `wave-2`
- `並列可能対象`: `frontend-output-review-ui`
- `並列不可理由`: `backend-output-review-readiness` とは `承認済み実装範囲重複`
- `完了条件`:
  - 各 `XTRANSLATOR_OUTPUT_ROW` 候補は 1 つの `JOB_TRANSLATION_FIELD` に対応し、`EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を持つ。
  - `cached` は row `Status=1` になり、辞書置換である事実は内部 summary に分離される。
  - 未定義 status、欠損識別子、row 重複候補は成功 row に混入しない。
  - diff preview は Source、Dest、Status、row reflection summary、stale reason、再出力可否を返せる。
  - compatibility validator は大きすぎる文字列、翻訳非推奨 field、先頭または末尾空白などを検出し、summary に出す。
  - 致命的な構造違反は success artifact にならない入力として返る。
- `想定変更ファイル数`: `8-14 files`
- `想定変更行数`: `450-800 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `補足`: warning / reject の細分類は人間レビュー済み範囲では「検出と表示」までを固定し、詳細閾値の拡張は別 task に残す。XML file write と artifact 保存は次 handoff に分ける。

### `backend-output-artifact-write-reoutput`

- `承認済み実装範囲`: xTranslator XML serializer、filesystem adapter、artifact / row 永続化、再出力、失敗回復、provider 非接続 assertion を backend で実装する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-output-artifact-public-seams`, `backend-output-review-readiness`, `backend-output-row-builder-compatibility`
  - `architecture_layer_basis`: service / adapter port / repository / SQLite concrete を backend 境界として扱う。filesystem concrete は adapter に閉じ、frontend UI は含めない。
  - `frozen_public_seams`: XML UTF-8、Skyrim root mapping、file write result、unique artifact per job、unique row per field、failure kind、retryable flag。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: job ID, artifact ID, row count, file path, generated_at, status, failed stage, retryable flag, digest, error kind
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header, secret store value
  - `secret_resolution_owner_layer`: secret 解決は行わない。fake provider と fake secret store は fail-on-call にできることを検証する。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を XML comment、artifact summary、operation summary、structured log、debug log、error summary に出さない。
- `依存対象`: `backend-output-review-readiness`, `backend-output-row-builder-compatibility`
- `検証`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'TranslationOutput|XTranslator|Artifact|Reoutput|SCN_TOA_(004|006|007|008|010)'`
- `初手`: `internal/repository/job_output_repository.go` に `TranslationArtifact` / `XTranslatorOutputRow` の upsert または replace 境界を追加し、同じ job と field で重複 row を作らない focused test を固定する。対応する完了条件は「同じ job の再出力で artifact と row が重複しない」である。理由は XML 出力、再出力、失敗回復が同じ永続化境界に依存するため。
- `実行グループ`: `wave-3`
- `ready_wave`: `wave-3`
- `並列可能対象`: なし
- `並列不可理由`: `依存対象`
- `完了条件`:
  - XML は UTF-8 として local parser で parse でき、Skyrim SE / LE の root element と `<String>` 子要素を持つ。
  - XML 特殊文字と日本語 text は壊れない。real xTranslator 起動は必須にしない。
  - 同じ job の再出力は現行 artifact を更新または置換し、同一 field の row を重複作成しない。
  - row validation、XML serialization、file write、artifact 保存失敗は success artifact として公開されない。
  - retryable flag、failed stage、redacted error summary、operation summary が返る。
  - output artifact の success、failure、re-output で provider call count と secret store call count は 0 になる。
- `想定変更ファイル数`: `14-22 files`
- `想定変更行数`: `800-1300 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装前`
- `補足`: 想定規模は caution。1 受け入れユースケースは「xTranslator 成果物の生成と再出力」で閉じる。row builder と readiness を依存完了に分けているため、分割必須にはしない。schema 追加が必要な場合は artifact summary と failure summary に限定し、docs 正本化は含めない。`本番経路`: usecase -> service -> XML adapter / repository / transactor -> SQLite / filesystem。

### `frontend-output-review-ui`

- `承認済み実装範囲`: Output Review の job list、summary、diff preview、action rail、state variants、disabled reason、redacted summary を frontend 状態と UI に実装する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-output-artifact-public-seams`
  - `architecture_layer_basis`: frontend contract / gateway contract / store / usecase / presenter / controller / View を frontend 境界として扱う。generated `wailsjs` と backend DTO は View が直接扱わない。
  - `frozen_public_seams`: `contract-output-artifact-public-seams` の gateway contract と screen state contract を参照する。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: job ID, artifact ID, field ID, file path, row count, status distribution, redacted error kind, compatibility kind, retryable flag
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: frontend は secret を解決しない。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を UI、store、debug output、test fixture snapshot に出さない。
- `依存対象`: `contract-output-artifact-public-seams`
- `検証`: `npm --prefix frontend run test -- --run translation-output`、`npm --prefix frontend run lint`
- `初手`: `frontend/src/application/contract/translation-output-artifact/translation-output-artifact-screen-types.ts` を追加し、loading / empty / not-ready / ready / generating / success / failed / stale state と disabled reason の型を固定する。対応する完了条件は「UI state variants と button enablement が contract として表現される」である。理由は presenter、store、View が同じ状態語彙に依存するため。
- `実行グループ`: `wave-2`
- `ready_wave`: `wave-2`
- `並列可能対象`: `backend-output-review-readiness`, `backend-output-row-builder-compatibility`
- `並列不可理由`: なし
- `完了条件`:
  - completed job list、selected job summary、readiness、拒否理由、result summary、artifact status が表示される。
  - diff preview は Source、Dest、Status、row reflection summary、stale reason、再出力可否を translation unit 単位で表示する。
  - 出力 action は readiness true、row validation pass、出力先 path valid の時だけ有効になる。
  - invalid job、`Canceled` job、field result 不整合、status 不整合では出力 action が disabled になり、理由が隣接表示または tooltip で確認できる。
  - target count 0、row count 0、readonly path、XML parse failure、compatibility warning が区別される。
  - secret、API key、provider raw payload、過剰な本文全文を UI と test fixture に出さない。
- `想定変更ファイル数`: `10-15 files`
- `想定変更行数`: `600-900 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `補足`: `UI人間操作E2E` の証明は最終検証で行う。この handoff は fake gateway と frontend local test で状態と表示 contract を閉じる。`本番経路`: screen controller -> frontend usecase -> gateway contract -> store / presenter -> Output Review UI。

### `integration-output-artifact-wails-gateway`

- `承認済み実装範囲`: backend Wails controller、generated binding wrapper、frontend Wails gateway、DTO mapping、bootstrap wiring を接続する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `contract_freeze`:
  - `status`: `done`
  - `freeze_source`: `contract-output-artifact-public-seams`
  - `architecture_layer_basis`: Wails は transport boundary として扱う。frontend gateway -> generated binding wrapper -> backend Wails controller -> usecase の接続だけを扱い、backend rule と frontend UI の代替にしない。
  - `frozen_public_seams`: `GetTranslationOutputReview`, `GetTranslationOutputDiffPreview`, `GenerateXTranslatorOutputArtifact`, `RegenerateXTranslatorOutputArtifact`。
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: contract handoff の許可参照値だけを通す。
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: Wails 境界では secret を解決しない。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を DTO、gateway log、runtime event に出さない。
- `依存対象`: `backend-output-review-readiness`, `backend-output-row-builder-compatibility`, `backend-output-artifact-write-reoutput`, `frontend-output-review-ui`
- `検証`: `go test ./internal/controller/wails ./internal/bootstrap -run 'TranslationOutput|OutputArtifact'`、`npm --prefix frontend run test -- --run translation-output-gateway`
- `初手`: `internal/controller/wails/translation_output_artifact_controller.go` を追加し、`GetTranslationOutputReview` の DTO mapping test を固定する。対応する完了条件は「contract freeze 済み response が Wails DTO と frontend gateway contract に欠落なく写る」である。理由は query / command の接続不備を最小 seam で先に検出できるため。
- `実行グループ`: `wave-4`
- `ready_wave`: `wave-4`
- `並列可能対象`: なし
- `並列不可理由`: `依存対象`
- `完了条件`:
  - backend controller は usecase port だけに依存し、service concrete を直接 new しない。
  - frontend gateway は generated `wailsjs` を adapter 内でだけ扱い、View へ backend DTO を漏らさない。
  - request / response の field 名、nullability、error kind、redaction obligation が contract と一致する。
  - bootstrap wiring は translation output artifact usecase、repository、adapter、controller を接続する。
  - Wails DTO、frontend gateway DTO、backend usecase response の変換で secret 禁止値を表現しない。
- `想定変更ファイル数`: `8-14 files`
- `想定変更行数`: `450-800 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `補足`: Wails generated files は必要な生成手順で更新するが、hand edit しない。`本番経路`: frontend Wails gateway -> generated binding wrapper -> backend Wails controller -> usecase。

### `final-validation-and-report`

- `承認済み実装範囲`: 全 handoff 完了後の最終検証、UI 人間操作確認、レビュー入力、作業レポート入力を集約する。
- `implementation_artifact`: `最終検証`
- `implementation_skill`: `implement-lane`
- `contract_freeze`:
  - `status`: `not_required`
  - `freeze_source`: `contract-output-artifact-public-seams`
  - `architecture_layer_basis`: 全 wave 完了後の検証境界であり、実装変更は扱わない。
  - `frozen_public_seams`: なし
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 最終検証で表示される redacted summary と operation summary だけを確認する。
  - `secret_values_for_provider_external_api_internal_auth`: decrypted API key, provider token, authorization header
  - `secret_resolution_owner_layer`: 最終検証は secret を解決しない。
  - `forbidden_outputs`: secret 本体、API key、provider raw payload、復号可能値、過剰な本文全文を UI 証跡、test log、work report 入力に出さない。
- `依存対象`: `integration-output-artifact-wails-gateway`
- `検証`: `python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local`、`python3 scripts/harness/run.py --suite scenario-gate`、`python3 scripts/harness/run.py --suite system-test`
- `初手`: `python3 scripts/harness/run.py --suite backend-local` を実行し、backend handoff の局所検証が全体 harness でも通ることを先に確認する。対応する完了条件は「backend local validation が全 handoff の backend 変更を通過させる」である。理由は integration と UI 証跡の前に永続化、usecase、controller の破損を切り分けるため。
- `実行グループ`: `wave-5`
- `ready_wave`: `wave-5`
- `並列可能対象`: なし
- `並列不可理由`: `広域判定条件共有`
- `完了条件`:
  - backend-local、frontend-local、scenario-gate の結果を記録する。
  - system-test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
  - Output Review を開き、completed job、not-ready job、diff preview、XML 出力、再出力、失敗 summary、redaction を UI 人間操作起点で確認する。
  - final validation 後に Codex review 入力として changed files、test results、UI evidence、residual risks を渡せる。
  - docs 正本化は実装後レビュー完了後の Codex docs lane 判断へ残し、この handoff では変更しない。
- `想定変更ファイル数`: `0 files`
- `想定変更行数`: `0 changed lines`
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `最終検証`
- `補足`: real xTranslator import smoke は必須にしない。必要なら任意手動確認として `manual_optional` に分け、close 必須条件にしない。

## Implement Lane Entry

人間が `implement_lane` へ渡す入口は次の形にする。

- `task_id`: `translation-output-artifact`
- `source_scope`: `docs/exec-plans/active/translation-output-artifact/implementation-scope.md`
- `start_ready_wave`: `wave-1`
- `handoff_order`: `contract-output-artifact-public-seams` -> `wave-2` handoffs -> `backend-output-artifact-write-reoutput` -> `integration-output-artifact-wails-gateway` -> `final-validation-and-report`
- `approved_design_bundle`: `scenario-design.md`, `ui-design.md`, coverage / gate 成果物
- `human_review`: `approved`

## Implementation Prohibitions

- プロダクトコード実装を `designer` から直接依頼しない。
- backend、frontend、統合境界を同一引き継ぎに束ねない。
- `contract-output-artifact-public-seams` 完了前に frontend handoff を開始しない。
- docs 正本、`.codex`、agent 実行定義、tool 権限を Codex implementation レーンへ渡さない。
- real paid AI API、real xTranslator 自動操作、secret store 解決を必須検証にしない。

## Expected Completion Report

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate の結果。Sonar server の Quality Gate ではない。
- `harness_gate_result`: system test が環境で止まる場合は `FAIL_ENVIRONMENT`、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: work reporter が読む実装事実。report 文面ではなく、completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `docs_changes`: `none`
