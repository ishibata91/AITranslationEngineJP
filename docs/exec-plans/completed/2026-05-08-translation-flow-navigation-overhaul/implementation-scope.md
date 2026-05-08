# Implementation Scope: 2026-05-08-translation-flow-navigation-overhaul

- `skill`: implementation-scope
- `status`: approved-handoff-ready
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: 人間介入状態により、人間設計レビューは approve 済みとして扱う。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `plan`: `./plan.md`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `ui_design`: `./ui-design.md`
- `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `component_diff`: `./navigation-design-diff-components.puml`
- `sequence_diff`: `./navigation-design-diff-sequence.puml`

## Fixed Decisions

- `needs_human_decision`: `0`
- frontend 実装、backend 実装、統合境界実装、シナリオテスト、単体テストを別 handoff にする。
- UI が関係するため、frontend 実装を backend 実装と統合境界実装より前に置く。
- frontend 実装完了後に UX 事前確認と人間レビューを通す。承認前に backend 実装と統合境界実装を開始しない。
- `wave-N` は任意長の連番である。必要な依存が完了した後続 wave を追加できる。
- `E2E` は UI 人間操作起点だけを指す。`APIテスト` は public seam 起点の system-level test とする。
- docs 正本化、`.codex`、`docs/detail-specs/`、プロダクト実装以外の作業流れ変更を handoff に含めない。

## Contract Freeze

- `status`: frozen-after-human-review
- `freeze_source`: `./plan.md`, `./scenario-design.md`, `./ui-design.md`
- `frozen_public_seams`:
  - 翻訳管理初期表示は未完了 job 一覧である。
  - 未完了 job 一覧が新規開始と途中再開の入口である。
  - Job Setup 作成成功後は単語翻訳ページへ進む。
  - フェーズページは job 確定後だけ表示する。
  - sticky footer は移動だけを扱う。
  - 翻訳完了ページは結果確認と出力管理への移動だけを扱う。
  - 出力管理へ移動しても job を自動選択しない。
  - secret、API key 平文、provider raw payload、過剰本文は UI、DTO、summary、log に出さない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-shell-list-navigation` | なし | なし | `backend_frontend_order` |
| `wave-2` | `frontend-phase-pages-footer-complete` | `frontend-shell-list-navigation` | なし | `owned_scope_overlap` |
| `wave-3` | `frontend-ux-precheck-and-human-review` | `frontend-shell-list-navigation`, `frontend-phase-pages-footer-complete` | なし | `backend_frontend_order` |
| `wave-4` | `backend-navigation-readiness-seams` | `frontend-ux-precheck-and-human-review` | なし | `shared_contract_change` |
| `wave-5` | `backend-idempotency-runtime-safety` | `backend-navigation-readiness-seams` | なし | `owned_scope_overlap` |
| `wave-6` | `integration-wails-gateway-navigation` | `backend-navigation-readiness-seams`, `backend-idempotency-runtime-safety`, `frontend-ux-precheck-and-human-review` | なし | `shared_contract_change` |
| `wave-7` | `scenario-tests-navigation-flow`, `unit-tests-navigation-slices` | `integration-wails-gateway-navigation` | `scenario-tests-navigation-flow <-> unit-tests-navigation-slices` | なし |

## Handoffs

### `frontend-shell-list-navigation`

- `implementation_target`: 翻訳管理入口、未完了 job 一覧、新規開始、途中再開、直移動防止の frontend 状態を作る。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `estimated_size`: `8-12 files`, `450-750 changed lines`, 通常
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: credential 状態分類、provider 名、model 名、redacted summary
  - `secret_values_for_provider_external_api_internal_auth`: 扱わない
  - `secret_resolution_owner_layer`: 既存 backend provider 設定層
  - `forbidden_outputs`: API key 平文、secret store key、provider raw payload、過剰本文、URL、DTO、UI、console、error summary
- `owned_scope`:
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
  - `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`
  - `frontend/src/application/contract/translation-job-management/`
  - `frontend/src/application/presenter/translation-job-management/`
  - `frontend/src/application/usecase/translation-job-management/`
  - `frontend/src/controller/translation-job-management/`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `frontend/src/ui/stores/shell-state.ts` の `defaultTranslationManagementViewId` と view 定義を、未完了 job 一覧初期表示を閉じる clause として変更する。初手にする理由は SCN-TFN-001 の入口条件を最小単位で固定できるためである。
- `agent_entry_input`:
  - `read_files`: `./plan.md`, `./scenario-design.md`, `./ui-design.md`, `docs/architecture.md`, `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/stores/shell-state.ts`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - `forbidden_scope`: `.codex/`, `docs/detail-specs/`, docs 正本、backend 実装、Wails DTO 接続、product test 以外の広域 gate 修正
  - `expected_outputs`: 翻訳管理初期ページを未完了 job 一覧にする。新規開始は入力データページへ進む。未完了 job 選択は job と current phase を固定する。旧 `Job Run` 直リンクと summary 取得入口を表示しない。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run frontend/src/ui/views/AppShell.test.ts frontend/src/controller/translation-job-management/translation-job-management-screen-controller.test.ts frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.test.ts`
- `completion_signal`:
  - SCN-TFN-001、SCN-TFN-003、SCN-TFN-004 の frontend 入口条件を満たす。
  - `Job Run` 表示名、フェーズ直リンク、セッション取得 UI が翻訳管理入口から消える。
  - job 未確定時は phase summary 取得や runtime event 購読を始めず、未完了 job 一覧へ戻る状態を持つ。
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後
- `notes`:
  - 既存の app shell と translation management 構造を土台にする。
  - 新しい独自 page shell、配色体系、余白体系を作らない。

### `frontend-phase-pages-footer-complete`

- `implementation_target`: 旧 `JobRunPage` をフェーズページ、移動専用 sticky footer、翻訳完了ページへ分解する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `estimated_size`: `12-15 files`, `650-800 changed lines`, 通常
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: phase state、readiness reason、redacted error summary、credential 状態分類
  - `secret_values_for_provider_external_api_internal_auth`: 扱わない
  - `secret_resolution_owner_layer`: 既存 backend provider 設定層
  - `forbidden_outputs`: API key 平文、secret store key、provider raw payload、過剰本文、prompt raw payload、console、error summary
- `owned_scope`:
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
  - `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`
  - `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`
  - `frontend/src/ui/screens/translation-output-artifact/TranslationOutputArtifactPage.svelte`
  - phase screen の controller / usecase / presenter / store
  - 画面専用の新規 phase footer / complete page component
- `depends_on`: `frontend-shell-list-navigation`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `first_action`: `frontend/src/ui/screens/job-run/JobRunPage.svelte` から job id 入力と `summary 取得` UI を削除し、選択済み job だけを phase page entry とする clause を閉じる。初手にする理由は旧セッション取得経路の廃止を最小単位で確認できるためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `./ui-design.md`, `./navigation-design-diff-components.puml`, `./navigation-design-diff-sequence.puml`, `frontend/src/ui/screens/job-run/JobRunPage.svelte`, 各 phase panel、各 phase presenter
  - `forbidden_scope`: `.codex/`, `docs/detail-specs/`, docs 正本、backend 実装、Wails DTO 接続、XML 出力処理の再設計
  - `expected_outputs`: phase 実行操作は本文に残す。sticky footer は `次へ進む`、`一覧へ戻る`、`出力管理へ移動` だけを扱う。翻訳完了ページは原文、訳文、ページング、出力管理移動だけを扱う。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run frontend/src/ui/screens/job-run frontend/src/ui/screens/term-translation-phase frontend/src/ui/screens/persona-generation-phase frontend/src/ui/screens/body-translation-phase`
- `completion_signal`:
  - SCN-TFN-002、SCN-TFN-005、SCN-TFN-006、SCN-TFN-007、SCN-TFN-008、SCN-TFN-009 の UI 条件を満たす。
  - footer 操作で provider request、phase start、retry、cancel、XML 出力を起動しない。
  - `Canceled` と `Failed` は翻訳完了ページへ入らない。
  - 出力管理へ移動後も selected job summary は未選択または一覧選択待ちから始まる。
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後
- `notes`:
  - 原文、訳文、job ID、file path、model 名、error reason は折り返しまたは省略と詳細展開を持つ。
  - mobile で sticky footer が本文を隠さない余白を確保する。

### `frontend-ux-precheck-and-human-review`

- `implementation_target`: frontend 実装後の UX 事前確認と人間レビュー gate を固定する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `estimated_size`: `0 product files`, `0 changed lines`, gate only
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: credential 状態分類
  - `secret_values_for_provider_external_api_internal_auth`: 扱わない
  - `secret_resolution_owner_layer`: 既存 backend provider 設定層
  - `forbidden_outputs`: API key 平文、secret、provider raw payload、console、error summary
- `owned_scope`:
  - `tmp/agent-browser/`
  - `tmp/logs/`
  - `test-results/`
- `depends_on`: `frontend-shell-list-navigation`, `frontend-phase-pages-footer-complete`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `npm run dev:wails:agent-browser` で実画面確認用アプリを起動し、翻訳管理入口を開く clause を閉じる。初手にする理由は backend 開始前の UI gate を実物で確認するためである。
- `agent_entry_input`:
  - `read_files`: `./ui-design.md`, `./scenario-design.md`
  - `forbidden_scope`: product code、product test、docs 正本、`.codex/`, `docs/detail-specs/`
  - `expected_outputs`: desktop と mobile の UI 証跡、console error 確認、secret 平文なし確認、人間レビュー依頼用の確認結果
- `validation_commands`:
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115/#translation-management`
  - `agent-browser screenshot tmp/agent-browser/translation-flow-navigation-frontend-desktop.png`
  - `agent-browser screenshot tmp/agent-browser/translation-flow-navigation-frontend-mobile.png`
  - `agent-browser errors`
- `completion_signal`:
  - 未完了 job 一覧、新規開始、各フェーズページ、翻訳完了ページ、出力管理入口を確認した。
  - desktop と mobile で footer、一覧、phase summary、長文理由が重ならない。
  - 旧 `Job Run`、summary 取得、フェーズ直リンク、前工程戻り導線が残っていない。
  - 人間レビューが approve になった。
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後
- `notes`:
  - この gate が approve になるまで backend 実装と統合境界実装へ進まない。

### `backend-navigation-readiness-seams`

- `implementation_target`: navigation と output readiness の backend public seam を UI 契約に合わせる。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `estimated_size`: `10-15 files`, `650-800 changed lines`, 通常
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: credential 状態分類、provider 名、model 名、redacted error summary、snapshot digest
  - `secret_values_for_provider_external_api_internal_auth`: provider 実行層だけが扱う secret 本体
  - `secret_resolution_owner_layer`: backend provider 設定 service / AI provider adapter
  - `forbidden_outputs`: API key 平文、secret store key、provider raw request / response、過剰本文、URL、DTO、UI、log、error summary、audit、request capture
- `owned_scope`:
  - `internal/usecase/translation_job_management_*`
  - `internal/service/translation_job_management_*`
  - `internal/controller/wails/translation_job_management_*`
  - `internal/usecase/translation_job_setup_*`
  - `internal/usecase/term_translation_phase_*`
  - `internal/usecase/persona_generation_phase_*`
  - `internal/usecase/body_translation_phase_*`
  - `internal/usecase/translation_output_artifact_*`
  - `internal/controller/wails/*phase*`
  - `internal/controller/wails/translation_output_artifact_*`
- `depends_on`: `frontend-ux-precheck-and-human-review`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `internal/usecase/translation_job_management_contract.go` の未完了 job list / detail contract に、current phase と再開不可理由を UI が安全に判定できる clause を閉じる。初手にする理由は SCN-TFN-003 と SCN-TFN-004 の public seam が後続 gateway の前提になるためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `./scenario-design.requirement-coverage.json`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-output-artifact.md`, `docs/detail-specs/body-translation-phase.md`
  - `forbidden_scope`: frontend UI 実装、Wails generated binding 手編集、docs 正本、`.codex/`, `docs/detail-specs/`, migration 追加、XML 出力仕様再設計
  - `expected_outputs`: Completed を未完了 job 一覧から除外する。参照不能 job は一覧理由表示に留める。Job Setup 成功後の単語翻訳 entry に必要な job summary を返す。phase readiness と output readiness を redacted summary で返す。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/service ./internal/controller/wails`
  - `npm run test:backend`
- `completion_signal`:
  - SCN-TFN-002 から SCN-TFN-010 までの backend public seam が成立する。
  - output management の review 取得は selected job なし開始を扱える。
  - secret と provider raw payload が DTO、summary、log に出ない。
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後
- `notes`:
  - Wails は transport boundary であり、domain rule や画面状態の正本にしない。

### `backend-idempotency-runtime-safety`

- `implementation_target`: retry、resume、late response、runtime event の状態不変条件を backend 側で固定する。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `estimated_size`: `8-14 files`, `500-750 changed lines`, 通常
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: phase run id、redacted error kind、ignored event count、late response rejected summary
  - `secret_values_for_provider_external_api_internal_auth`: provider 実行層だけが扱う secret 本体
  - `secret_resolution_owner_layer`: backend provider adapter
  - `forbidden_outputs`: API key 平文、provider raw request / response、過剰本文、log、error summary、runtime event payload
- `owned_scope`:
  - `internal/usecase/term_translation_phase_*`
  - `internal/usecase/persona_generation_phase_*`
  - `internal/usecase/body_translation_phase_*`
  - `internal/service/term_translation_phase_*`
  - `internal/service/persona_generation_phase_*`
  - `internal/service/body_translation_phase_*`
  - `internal/repository/job_*`
  - `internal/jobio/`
  - runtime event publisher / adapter
- `depends_on`: `backend-navigation-readiness-seams`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `first_action`: `internal/usecase/body_translation_phase_contract.go` または対応 usecase contract に、terminal job と late response を後書きしない clause を閉じる。初手にする理由は SCN-TFN-011 の不変条件を public seam と service の両側で固定するためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `docs/architecture.md`, phase usecase/service/repository tests
  - `forbidden_scope`: frontend UI 実装、Wails generated binding 手編集、docs 正本、`.codex/`, `docs/detail-specs/`, schema migration
  - `expected_outputs`: retry、resume、開始再送で同じ phase run を継続する。成功済み result と artifact row を重複作成しない。別 job と古い phase run の runtime event は画面遷移や provider 再実行を起こさない。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/service ./internal/repository ./internal/jobio`
  - `npm run test:backend`
- `completion_signal`:
  - SCN-TFN-011 と SCN-TFN-012 の backend 側不変条件が成立する。
  - terminal job は状態を変えない。
  - late response は現在 phase run と一致しない場合に後書きされない。
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後
- `notes`:
  - DB schema 変更が必要になった場合は、この handoff では停止して implement_lane へ戻す。

### `integration-wails-gateway-navigation`

- `implementation_target`: frontend gateway、Wails DTO、backend controller を接続し、実画面で navigation flow を確認する。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `estimated_size`: `8-13 files`, `500-750 changed lines`, 通常
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: credential 状態分類、redacted summary、digest、job id
  - `secret_values_for_provider_external_api_internal_auth`: Wails DTO と frontend gateway へ渡さない
  - `secret_resolution_owner_layer`: backend provider 設定 service / AI provider adapter
  - `forbidden_outputs`: API key 平文、secret store key、provider raw payload、URL、DTO、UI、console、log、request capture
- `owned_scope`:
  - `frontend/src/application/gateway-contract/translation-job-management/`
  - `frontend/src/application/gateway-contract/translation-job-setup/`
  - `frontend/src/application/gateway-contract/*phase*/`
  - `frontend/src/application/gateway-contract/translation-output-artifact/`
  - `frontend/src/controller/wails/*gateway*`
  - `frontend/src/controller/wails/gateway-dto/`
  - `internal/controller/wails/`
  - generated `frontend/wailsjs/` は生成結果だけ扱い、手編集しない
- `depends_on`: `backend-navigation-readiness-seams`, `backend-idempotency-runtime-safety`, `frontend-ux-precheck-and-human-review`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `frontend/src/controller/wails/translation-job-management.gateway.ts` の list / detail mapping を backend contract と一致させ、未完了 job 一覧の接続 clause を閉じる。初手にする理由は UI 初期ページの実データ接続が後続接続の入口になるためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `./ui-design.md`, `docs/architecture.md`, frontend gateway contracts, backend Wails controllers
  - `forbidden_scope`: frontend UI 再設計、backend domain rule 再設計、docs 正本、`.codex/`, `docs/detail-specs/`, generated binding の手編集
  - `expected_outputs`: Wails DTO と frontend gateway が redacted public seam を保つ。未完了 job 一覧、Job Setup 作成後、phase readiness、output review が実画面で接続される。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --run frontend/src/controller/wails`
  - `go test ./internal/controller/wails`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115/#translation-management`
  - `agent-browser errors`
- `completion_signal`:
  - SCN-TFN-001 から SCN-TFN-010 までの UI 接続が実画面で確認できる。
  - 統合後も secret 平文と provider raw payload が UI、DTO、console に出ない。
  - 出力管理へ移動しても selected job は自動選択されない。
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E
- `execution_stage`: 実装後
- `notes`:
  - 実画面確認のスクリーンショットは `tmp/agent-browser/` に置く。

### `scenario-tests-navigation-flow`

- `implementation_target`: 承認済みシナリオを APIテストと UI人間操作E2E で証明する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `estimated_size`: `5-10 files`, `600-900 changed lines`, 注意
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: redacted fixture、credential 状態分類
  - `secret_values_for_provider_external_api_internal_auth`: fake だけで扱い、実 secret は使わない
  - `secret_resolution_owner_layer`: fake gateway / test backend fixture
  - `forbidden_outputs`: 実 API key、provider raw payload、過剰本文、test log、screenshot、console
- `owned_scope`:
  - `internal/apitest/`
  - `internal/integrationtest/`
  - `frontend/src/ui/**/*.test.ts`
  - `test-results/`
  - `tmp/agent-browser/`
- `depends_on`: `integration-wails-gateway-navigation`
- `execution_group`: `wave-7`
- `ready_wave`: `wave-7`
- `parallelizable_with`: `unit-tests-navigation-slices`
- `parallel_blockers`: なし
- `first_action`: `internal/apitest/` に SCN-TFN-011 の retry / resume / late response API test を追加し、重複作成防止 clause を閉じる。初手にする理由は UI より先に backend 不変条件を public seam で証明できるためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `./scenario-design.requirement-coverage.json`, `./ui-design.md`, existing `internal/apitest`, existing `internal/integrationtest`
  - `forbidden_scope`: production code、docs 正本、`.codex/`, `docs/detail-specs/`, broad harness 修正、実 AI API 呼び出し
  - `expected_outputs`: SCN-TFN-001 から SCN-TFN-012 までを、APIテスト、UI人間操作E2E、補助 frontend test に対応づける。実 AI API を使わず fake / fixture で通す。
- `validation_commands`:
  - `npm run test:backend`
  - `npm --prefix frontend run test`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115/#translation-management`
  - `agent-browser screenshot tmp/agent-browser/translation-flow-navigation-scenario.png`
  - `agent-browser errors`
- `completion_signal`:
  - SCN-TFN-001 から SCN-TFN-012 までの証明先が test result と UI 証跡で追える。
  - UI人間操作E2E は裏側 API 直接投入だけで完了扱いにしない。
  - system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` と再実行条件を残す。
- `acceptance_test`: required
- `execution_test_classification`: UI人間操作E2E / APIテスト
- `execution_stage`: 実装後
- `notes`:
  - 注意規模だが、全シナリオを 1 つの受け入れテスト証明成果物として集約するため 1 handoff にする。
  - production code 修正が必要になった場合は停止し、fix lane または実装 handoff へ戻す。

### `unit-tests-navigation-slices`

- `implementation_target`: frontend と backend の unit-level 振る舞いを、実装成果物ごとに補強する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `estimated_size`: `8-14 files`, `650-900 changed lines`, 注意
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: redacted fixture、credential 状態分類
  - `secret_values_for_provider_external_api_internal_auth`: test fake にも実 secret を入れない
  - `secret_resolution_owner_layer`: fake gateway / backend test fixture
  - `forbidden_outputs`: 実 API key、provider raw payload、過剰本文、test log、snapshot
- `owned_scope`:
  - `frontend/src/application/**/*.test.ts`
  - `frontend/src/controller/**/*.test.ts`
  - `frontend/src/ui/**/*.test.ts`
  - `internal/usecase/**/*_test.go`
  - `internal/service/**/*_test.go`
  - `internal/controller/wails/**/*_test.go`
- `depends_on`: `integration-wails-gateway-navigation`
- `execution_group`: `wave-7`
- `ready_wave`: `wave-7`
- `parallelizable_with`: `scenario-tests-navigation-flow`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.test.ts` に、Completed 除外と phase page target 作成の view model clause を追加する。初手にする理由は SCN-TFN-001 と SCN-TFN-003 の UI 状態を production UI なしで検証できるためである。
- `agent_entry_input`:
  - `read_files`: `./scenario-design.md`, `./ui-design.md`, touched frontend and backend tests
  - `forbidden_scope`: production code、docs 正本、`.codex/`, `docs/detail-specs/`, broad harness 修正、実 AI API 呼び出し
  - `expected_outputs`: presenter、usecase、controller、gateway mapping、phase readiness、output readiness、runtime event filtering の単体テストを追加する。
- `validation_commands`:
  - `npm --prefix frontend run test`
  - `go test ./internal/usecase ./internal/service ./internal/controller/wails`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - frontend の footer enablement、禁止導線、長文理由、secret 非表示が unit-level で検証される。
  - backend の readiness、redaction、terminal state、late response reject が unit-level で検証される。
  - scenario test と同じ shared fixture を同時編集していない。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 注意規模だが、production code を変更しない test-only handoff であり、scenario test と並列可能である。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `frontend_ux_precheck_result`
- `frontend_human_review_result`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate を指す。Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `docs_changes: none`

## Stop Conditions

- frontend UX 事前確認または人間レビューが approve にならない場合は、backend 実装と統合境界実装へ進まない。
- DB schema migration、docs 正本化、`.codex` 変更、`docs/detail-specs/` 変更が必要になった場合は、対象 handoff を停止して implement_lane へ戻す。
- 実 AI API、実 secret、provider raw payload がテストに必要になった場合は停止する。
- backend と frontend を同一 handoff に束ねる必要が出た場合は停止する。
- generated `frontend/wailsjs/` を手編集する必要が出た場合は停止する。
