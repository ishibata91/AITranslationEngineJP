# Implementation Scope: phase-processing-target-list-refactor

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved by human goal continuation
- `approval_record`: 人間が「この問題のテスト追加，UC差分，修正までを行うこと」と指示したため、implementation-scope を承認済みとして扱う。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`
- `execution_type`: `コード併走型`

## Source Artifacts

- `refactor_plan`: `./plan.md`
- `spec_drift_investigation`: `./spec-drift-investigation.md`
- `structure_quality_investigation`: `./structure-quality-investigation.md`
- `test_quality_investigation`: `./test-quality-investigation.md`
- `detail_spec_sources`: `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `screen_design_sources`: `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `e2e_design_source`: `docs/e2e-test-design/test-design.csv`
- `detail_spec_diff`: N/A。refactor-lane の調査記録と人間固定判断を根拠にする。
- `screen_design_diff`: N/A。docs 正本の画面設計を根拠にし、docs 正本文は更新しない。

## Fixed Decisions

- `unanswered_questions`: `0`
- 3 フェーズで、一覧 total と画面上の処理対象件数は同じ母集団を示す。
- 3 フェーズで、一覧行が表示される。
- 3 フェーズで、検索できる。
- 単語翻訳の処理対象一覧は、共通辞書対象外の用語と固有名詞を母集団にする。
- 単語翻訳の画面上の処理対象件数は、AI 翻訳対象語件数と矛盾しない主語にする。
- NPC ペルソナ生成の処理対象一覧は、現行の対象件数母集団を維持する。
- NPC ペルソナ生成の検索表示は、query 対象である名前、FormID、EditorID、NPC 属性と矛盾しない主語にする。
- 本文翻訳の処理対象一覧は、辞書置換対象外の翻訳項目を母集団にする。
- 本文翻訳の画面上の処理対象件数は、AI 送信対象件数と矛盾しない主語にする。
- `translation_complete`、phase 以外の job lifecycle、実外部 API、実 secret、実利用者データは対象外にする。
- docs 正本文と `.codex/` は変更しない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-count-subject` | なし | なし | なし |
| `wave-2` | `frontend-search-subject` | `frontend-count-subject` | なし | `owned_scope_overlap` |
| `wave-3` | `backend-count-read-model` | `frontend-count-subject`, `frontend-search-subject` | なし | `backend_frontend_order` |
| `wave-4` | `backend-search-read-model` | `backend-count-read-model` | なし | `owned_scope_overlap` |
| `wave-5` | `integration-processing-target-seam` | `frontend-count-subject`, `frontend-search-subject`, `backend-count-read-model`, `backend-search-read-model` | なし | `shared_contract_change` |
| `wave-6` | `unit-count-subject`, `unit-search-subject`, `scenario-page-object`, `scenario-fixture` | `integration-processing-target-seam` | `unit-count-subject <-> unit-search-subject`, `scenario-page-object <-> scenario-fixture` | なし |
| `wave-7` | `scenario-phase-list-search` | `scenario-page-object`, `scenario-fixture` | なし | `depends_on` |

## Handoffs

### `frontend-count-subject`

- `implementation_target`: 3 フェーズ画面の件数主語を処理対象一覧 total と矛盾しない表示へそろえる。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./plan.md`, `./structure-quality-investigation.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
  - `screen_design_sources`: `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `frontend/src/ui/screens/*PhasePanel.svelte`、phase presenter、phase screen type、phase view-model test。単語翻訳は AI 翻訳対象語件数、NPC ペルソナ生成は現行対象件数、本文翻訳は AI 送信対象件数を処理対象件数の主語にする。
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `estimated_size`: `8-12 files`, `250-500 changed lines`, 通常
- `first_action`: `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` の `phaseMetrics` と `progressDetails` を、単語翻訳の処理対象件数が AI 翻訳対象語件数を指す 1 clause へ変更する。完了条件 `単語翻訳の件数主語が一覧 total と矛盾しない` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - 単語翻訳の処理対象件数は、AI 翻訳対象語件数と同じ主語で表示される。
  - NPC ペルソナ生成の処理対象件数は、現行対象件数の主語を維持する。
  - 本文翻訳の処理対象件数は、AI 送信対象件数と同じ主語で表示される。
  - `ProcessingTargetListPanel.svelte` のページング計算本体は変更されていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `forbidden_changes`: `translation_complete` 画面、phase 以外の job lifecycle、docs 正本、`.codex/`

### `frontend-search-subject`

- `implementation_target`: 3 フェーズ画面の検索表示を query 対象と矛盾しない主語へそろえる。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./plan.md`, `./structure-quality-investigation.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
  - `screen_design_sources`: `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: 3 フェーズ画面の `ProcessingTargetListWrapper` 呼び出し、検索 placeholder、`searchTestId`、検索入力結線。NPC ペルソナ生成は名前だけに限定しない。
- `depends_on`: `frontend-count-subject`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `estimated_size`: `3-6 files`, `80-180 changed lines`, 通常
- `first_action`: `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte` の `searchPlaceholder` を、名前、FormID、EditorID、NPC 属性を検索対象として読める文言へ変更する。完了条件 `NPC ペルソナ生成の検索表示が query 対象と矛盾しない` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - 3 フェーズの検索入力に stable な `data-testid` がある。
  - 単語翻訳と本文翻訳の検索表示は、名前、原文、訳語と矛盾しない。
  - NPC ペルソナ生成の検索表示は、名前、FormID、EditorID、NPC 属性と矛盾しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `forbidden_changes`: 検索対象と無関係な画面文言、docs 正本、`.codex/`

### `backend-count-read-model`

- `implementation_target`: processing target read model の一覧 total を、3 フェーズの処理対象件数と同じ母集団へそろえる。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./plan.md`, `./spec-drift-investigation.md`, `./structure-quality-investigation.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `internal/service/processing_target_read_model_service.go`、`internal/repository/processing_target_sqlite_repository.go`、必要な service / repository test。単語翻訳は共通辞書対象外の用語と固有名詞、NPC ペルソナ生成は現行対象件数、本文翻訳は辞書置換対象外の翻訳項目を total の母集団にする。
- `depends_on`: `frontend-count-subject`, `frontend-search-subject`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `backend_frontend_order`
- `estimated_size`: `4-8 files`, `250-650 changed lines`, 通常
- `first_action`: `internal/repository/processing_target_sqlite_repository.go` の `processingTargetTermCountSQL` と `processingTargetTermListSQL` を、翻訳ジョブ内辞書ではなく単語翻訳の AI 翻訳対象語母集団を返す clause へ変更する。完了条件 `単語翻訳の一覧 total と処理対象件数が同じ母集団を示す` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - 単語翻訳の一覧 total は AI 翻訳対象語件数と矛盾しない。
  - NPC ペルソナ生成の一覧 total は現行対象件数母集団を維持する。
  - 本文翻訳の一覧 total は辞書置換対象外の翻訳項目件数と矛盾しない。
  - `translation_complete` の query と page state は変更されていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `forbidden_changes`: phase 以外の job lifecycle、`translation_complete`、実外部 API、実 secret、実利用者データ

### `backend-search-read-model`

- `implementation_target`: processing target read model の検索対象を、3 フェーズの検索表示と矛盾しない query へそろえる。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./plan.md`, `./spec-drift-investigation.md`, `./structure-quality-investigation.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `processingTargetSearchPattern`、phase 別 count/list SQL の検索条件、必要な repository test。Page Object と scenario test は含めない。
- `depends_on`: `backend-count-read-model`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `estimated_size`: `2-5 files`, `120-320 changed lines`, 通常
- `first_action`: `internal/repository/processing_target_sqlite_repository.go` の `processingTargetPersonaCountSQL` と `processingTargetPersonaListSQL` の検索条件を、NPC ペルソナ生成の表示文言と一致する検索主語へ整理する。完了条件 `NPC ペルソナ生成の検索表示と query 対象が矛盾しない` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - 単語翻訳の検索は、処理対象名、原文相当、訳語候補相当を対象にできる。
  - NPC ペルソナ生成の検索は、名前、FormID、EditorID、NPC 属性を対象にできる。
  - 本文翻訳の検索は、名前、原文、訳語を対象にできる。
  - 検索語なし、検索一致 1 件、検索結果 0 件を repository test で区別できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `forbidden_changes`: phase 以外の検索、`translation_complete`、docs 正本、`.codex/`

### `integration-processing-target-seam`

- `implementation_target`: Wails DTO、frontend gateway、phase usecase の接続境界で、一覧 total、items、searchQuery が 3 フェーズへ正しく届く状態にする。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./plan.md`, `./structure-quality-investigation.md`, `frontend/src/application/gateway-contract/processing-target/processing-target-gateway-contract.ts`, `internal/controller/wails/processing_target_controller.go`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `internal/controller/wails/processing_target_controller.go`、`frontend/src/application/gateway-contract/processing-target/processing-target-gateway-contract.ts`、`frontend/src/controller/wails/*phase.gateway.ts`、phase usecase の `getProcessingTargetList` 呼び出し境界、必要な gateway test。
- `depends_on`: `frontend-count-subject`, `frontend-search-subject`, `backend-count-read-model`, `backend-search-read-model`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `estimated_size`: `6-12 files`, `220-600 changed lines`, 通常
- `first_action`: `frontend/src/application/gateway-contract/processing-target/processing-target-gateway-contract.ts` の `ProcessingTargetListResponse` を確認し、既存 field のまま `totalCount` が phase 固有の処理対象母集団件数を示す clause を gateway test で固定する。完了条件 `public seam の totalCount 主語が downstream へ再解釈なしで渡る` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - `GetProcessingTargetList` の `phase`、`page`、`pageSize`、`searchQuery` が backend から frontend page state まで保持される。
  - `totalCount` は phase 固有の処理対象母集団件数として、frontend page state へ届く。
  - 新規 DTO field を追加した場合は、Go DTO、TypeScript contract、gateway runtime validation、mock が同じ shape にそろう。
  - 実画面または system-test 証跡で、3 フェーズの一覧行、件数表示、検索入力が同じ接続境界を通ることを示す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 既存 `metadata` field の単独整理は今回の主目的にしない。
  - DTO shape を増やす場合は、この handoff だけで public seam を閉じる。
- `forbidden_changes`: `translation_complete`、phase 以外の Wails controller、実 secret、実利用者データ、docs 正本

### `unit-count-subject`

- `implementation_target`: 件数主語の局所分岐と変換を単体テストで保護する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./plan.md`, `./structure-quality-investigation.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: phase presenter / usecase test、backend service / repository test、必要な gateway contract test。scenario test は含めない。
- `depends_on`: `integration-processing-target-seam`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: `unit-search-subject`, `scenario-page-object`, `scenario-fixture`
- `parallel_blockers`: なし
- `estimated_size`: `5-10 files`, `250-650 changed lines`, 通常
- `first_action`: `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.test.ts` に、単語翻訳の処理対象件数表示が AI 翻訳対象語件数を使う test を追加する。完了条件 `単語翻訳の件数主語が単体テストで固定される` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - 単語翻訳、NPC ペルソナ生成、本文翻訳の件数主語を test 名で説明できる。
  - repository または service test が、一覧 total の母集団差を検出できる。
  - frontend presenter / usecase test が、画面上の処理対象件数と page state total の主語差を検出できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `forbidden_changes`: product code、scenario test、docs 正本、`.codex/`

### `unit-search-subject`

- `implementation_target`: 検索主語の局所分岐と query 変換を単体テストで保護する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./plan.md`, `./structure-quality-investigation.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: phase usecase search event test、repository search test、gateway request forwarding test。Page Object と scenario fixture は含めない。
- `depends_on`: `integration-processing-target-seam`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: `unit-count-subject`, `scenario-page-object`, `scenario-fixture`
- `parallel_blockers`: なし
- `estimated_size`: `4-8 files`, `180-500 changed lines`, 通常
- `first_action`: `frontend/src/application/usecase/persona-generation-phase/persona-generation-phase.usecase.test.ts` に、検索入力が page 1 の `GetProcessingTargetList` request として流れる test を追加する。完了条件 `NPC ペルソナ生成の検索 request が単体テストで固定される` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - 検索語なし、検索一致 1 件、検索結果 0 件の単体テストがある。
  - 3 フェーズの検索 query が page 1 へ戻ることを test で確認できる。
  - NPC ペルソナ生成の検索対象が UI 表示と query で矛盾しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `forbidden_changes`: product code、scenario test、docs 正本、`.codex/`

### `scenario-page-object`

- `implementation_target`: phase system-test Page Object に、処理対象一覧の領域、件数、検索、空状態を観測する入口を追加する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./test-quality-investigation.md`, `docs/e2e-test-design/test-design.csv`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `tests/system/support/translation-phase-pages.ts`。fixture と scenario spec は含めない。
- `depends_on`: `integration-processing-target-seam`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: `unit-count-subject`, `unit-search-subject`, `scenario-fixture`
- `parallel_blockers`: なし
- `estimated_size`: `1-2 files`, `80-180 changed lines`, 通常
- `first_action`: `tests/system/support/translation-phase-pages.ts` に `processingTargetSearchInput` locator を追加する。完了条件 `Page Object が検索操作を持つ` を最初に閉じるため。
- `validation_commands`:
  - `npx playwright test --config ./playwright.config.ts --list`
- `completion_signal`:
  - Page Object から検索 input、一覧領域、一覧件数表示、空状態、行 locator を参照できる。
  - 追加 locator は 3 フェーズ共通 prefix で組み立てられる。
  - Page Object は product 実装詳細へ過剰に依存しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `forbidden_changes`: product code、fixture、scenario spec、docs 正本

### `scenario-fixture`

- `implementation_target`: system-test fixture で、phase 差、件数差、検索一致、検索結果 0 件を表現できるようにする。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./test-quality-investigation.md`, `docs/e2e-test-design/test-design.csv`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `tests/system/support/scenario-wails-mocks.ts`。Page Object と scenario spec は含めない。
- `depends_on`: `integration-processing-target-seam`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: `unit-count-subject`, `unit-search-subject`, `scenario-page-object`
- `parallel_blockers`: なし
- `estimated_size`: `1-2 files`, `180-420 changed lines`, 通常
- `first_action`: `tests/system/support/scenario-wails-mocks.ts` の `GetProcessingTargetList` を phase 別 fixture へ分岐する。完了条件 `fixture が 3 フェーズ別の処理対象一覧を返す` を最初に閉じるため。
- `validation_commands`:
  - `npx playwright test --config ./playwright.config.ts --list`
- `completion_signal`:
  - 単語翻訳、NPC ペルソナ生成、本文翻訳で別の処理対象行を返せる。
  - `request.searchQuery` により検索一致 1 件と検索結果 0 件を返せる。
  - `request.phase`、`request.page`、`request.pageSize`、`request.searchQuery` が fixture 応答へ反映される。
  - 実 secret、実外部 API、実利用者データを fixture に含めない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `forbidden_changes`: product code、Page Object、scenario spec、docs 正本

### `scenario-phase-list-search`

- `implementation_target`: `E2E-UC-045/046/047` または追加 scenario test で、3 フェーズの一覧表示、件数一致、検索を証明する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./test-quality-investigation.md`, `docs/e2e-test-design/test-design.csv`
- `frontend_required_sources`:
  - `screen_design_diff`: N/A
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: 実 secret、実外部 API 応答、実利用者データ
- `owned_scope`: `tests/system/job-run-shell.spec.ts`、必要なら `tests/system/translation-phases.spec.ts`。Page Object と fixture の新規 helper を使い、helper 本体は変更しない。
- `depends_on`: `scenario-page-object`, `scenario-fixture`
- `execution_group`: `wave-7`
- `ready_wave`: `wave-7`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `estimated_size`: `1-3 files`, `180-450 changed lines`, 通常
- `first_action`: `tests/system/job-run-shell.spec.ts` の `E2E-UC-045` を、単語翻訳の処理対象行、一覧件数、検索一致を確認する assertion へ置き換える。完了条件 `単語翻訳の一覧表示、件数一致、検索が scenario test で証明される` を最初に閉じるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`:
  - `E2E-UC-045` は単語翻訳の一覧行、処理対象件数、検索一致、検索結果 0 件を検証する。
  - `E2E-UC-046` は NPC ペルソナ生成の一覧行、処理対象件数、検索一致、検索結果 0 件を検証する。
  - `E2E-UC-047` は本文翻訳の一覧行、処理対象件数、検索一致、検索結果 0 件を検証する。
  - `E2E-UC-048/049/050` と `E2E-UC-051/052/053` の意図は維持される。
  - system-test が Wails readiness または Chromium 権限で止まる場合は、環境起因と product failure を分けて記録する。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `forbidden_changes`: product code、Page Object helper 本体、fixture helper 本体、docs 正本、`.codex/`

## Final Validation

- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite system-test`

`system-test` が Wails readiness または Chromium 権限で止まる場合は、`FAIL_ENVIRONMENT` として blocked reason、再実行環境、再実行コマンドを残す。

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
- `sonar_gate_result`: repo-local Sonar issue gate を指す。Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とする。
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`

## Human Review Items

- この implementation-scope を承認するか。
- DTO shape を既存 field 維持で進めるか。新規 field が必要になった場合は `integration-processing-target-seam` で閉じる。
- scenario test で `E2E-UC-045/046/047` を上書きするか、追加 test 名で証明するか。
