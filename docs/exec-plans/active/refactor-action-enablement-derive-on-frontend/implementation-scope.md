# Implementation Scope: refactor-action-enablement-derive-on-frontend

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: plan.md の `人間設計レビュー確定事項` 節（第 3 版 + H 節 + 責務 4 分割確定）
- `module_entry`: `.claude/skills/implementation-module/SKILL.md`
- `handoff_runtime`: `claude-module`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: 不要（plan.md decision table で仕様変更 N。外部仕様は変えず内部契約リファクタのみ）
- `screen_design_diff`: `N/A`（plan.md decision table で画面変更 N。画面構造・文言・layout を維持）
- `design_diff`: `./design-diff.md`（承認済み。本書の唯一の根拠参照）

## Fixed Decisions

- 設計差分図 `./design-diff.md` を根拠の正本にする。`detail-spec-diff.md` は存在しないため、各引き継ぎの `spec_basis` には `./design-diff.md` を入れ、参照節を明示する。
- 責務 4 分割（ドメイン情報集合 / ドメイン状態射影 / summary / UX 遷移可否）と、backend response 内で projection と summary を別 field group として並べる方針に従う。
- persona → body 移行の readiness 判定は body 側 projection の `personaBodyReadiness` だけに載せる（design-diff G-2-b / H-11）。
- BlockedReason 文字列は frontend 固定文字列とする。backend からは運ばない。
- `unanswered_questions`: `0`
- backend と frontend ロジックは別引き継ぎ。両者を接続する gateway contract / runtime shape validator は統合境界引き継ぎとして独立に切る。
- 画面変更がないため UI 起点の引き継ぎ順序拘束（frontend を backend より先 wave）は適用しない。本 task は契約起点のリファクタであり、契約変更が下流の gateway / presenter に波及する依存方向に従って wave を並べる。本判断は `implementation-scope` の境界規約「UI がある task では frontend を backend より前の wave に置く」の例外として、UI が無い（画面変更 N）ことを根拠に採る。
- テスト設計は本範囲に含めない（plan.md 指示。test-design 側で並列に進む）。
- 各 phase（term / persona / body）の契約は backend ファイルが別、gateway ファイルが別、presenter ファイルが別であり、phase ごとに owned_scope が重ならないため、同一 wave 内 phase 並列を成立条件とする。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `BE-term`, `BE-persona`, `BE-body` | `なし` | `BE-term <-> BE-persona <-> BE-body` | `なし` |
| `wave-2` | `INT-term`, `INT-persona`, `INT-body` | wave-1 の同 phase 引き継ぎが完了済み | `INT-term <-> INT-persona <-> INT-body` | `なし` |
| `wave-3` | `FE-term`, `FE-persona`, `FE-body` | wave-2 の同 phase 引き継ぎが完了済み | `FE-term <-> FE-persona <-> FE-body` | `なし` |

## Handoffs

### `BE-term`:

- `implementation_target`: term phase usecase contract と wails controller DTO から UX 遷移可否 flag と BlockedReason を除去し、ドメイン状態射影 result/DTO を summary と並列に追加する。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./design-diff.md`（C-1、C-5、G-5-a、G-6-a）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `internal/usecase/term_translation_phase_contract.go`: `TermTranslationPhaseProjectionResult`（仮称）を新設。field は `PhaseLifecycle` / `JobLifecycle` / `ErrorKind` / `AISettingsConfigured` / `AITargetCount` / `ConfirmedCount`。summary result の `TermTranslationPhaseActionEnablement` 型と `CanStart` / `CanPause` / `CanResume` / `CanRetry` / 各 `*BlockedReason` を削除。command result の `CanStartNextPhase` / `NextPhaseBlockedReason` を削除。fetch 系 usecase return を `{ summary, projection }` 構造に変える。
  - `internal/controller/wails/term_translation_phase_controller.go`: `TermTranslationPhaseProjectionDTO` を新設し、phase fetch response の field group として summary DTO と並列に並べる。`TermTranslationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field、command response の `canStartNextPhase` / `nextPhaseBlockedReason` を削除。
  - controller の組み立てロジックを usecase の 2 系統 result から JSON へ単純 mapping する形に変える。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `BE-persona`, `BE-body`
- `parallel_blockers`: `なし`（term/persona/body の契約ファイルは独立。shared 共通型なし）
- `first_action`:
  - `path`: `internal/usecase/term_translation_phase_contract.go`
  - `symbol`: `TermTranslationPhaseActionEnablement` 型定義
  - `変更種別`: 削除（型本体と参照箇所の連鎖削除）
  - `対応 completion_signal`: 「usecase contract から ActionEnablement 型と CanStartNextPhase が削除されている」
  - `理由`: ActionEnablement 型を起点に依存箇所が controller / fetch result に連鎖し、最初に型を除去することで残る修正範囲が直線的になる。
- `validation_commands`:
  - `go build ./internal/usecase/... ./internal/controller/wails/...`
  - `go vet ./internal/usecase/... ./internal/controller/wails/...`
- `completion_signal`:
  - usecase contract から `TermTranslationPhaseActionEnablement` 型、`Can*` field、`*BlockedReason` field、command result の `CanStartNextPhase` / `NextPhaseBlockedReason` が削除されている。
  - usecase fetch 系 return が summary と projection の 2 系統を並べた構造に変わっており、projection は design-diff C-1 の 6 field を満たす。
  - controller DTO に `TermTranslationPhaseProjectionDTO` が追加され、`TermTranslationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field、command response の `canStartNextPhase` / `nextPhaseBlockedReason` が削除されている。
  - `go build` / `go vet` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: usecase contract 1 + controller 1 + fetch 経路に呼び出す usecase 実装 (内訳 2-3 ファイル) ≈ 4-5 ファイル。
  - 想定変更行数: 削除 200 行前後 + 追加 100 行前後 ≈ 300 changed lines。規模区分は「通常」。
  - 本番経路: `WailsTermTranslationPhaseController.Fetch` → `TermTranslationPhaseUsecase.Fetch` → projection/summary 組み立て → JSON DTO。
  - 単体テスト / シナリオテストは本引き継ぎでは扱わない。test-design 側で並列に進む。

### `BE-persona`:

- `implementation_target`: persona phase usecase contract と wails controller DTO から UX 遷移可否 flag と BlockedReason を除去し、ドメイン状態射影 result/DTO を summary と並列に追加する。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./design-diff.md`（C-2、C-5、G-5-b、G-6-b、G-2-b 注記）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `internal/usecase/persona_generation_phase_contract.go`: `PersonaGenerationPhaseProjectionResult`（仮称）を新設。field は `PhaseLifecycle` / `JobLifecycle` / `ErrorKind` / `AISettingsConfigured` / `TargetCount` / `PreviousPhaseLifecycle`。`PersonaGenerationPhaseActionEnablement` 型と各 `Can*` / `*BlockedReason` を削除。command result の `CanStartBodyPhase` / `BodyReadinessBlockedReason` を削除。`PersonaGenerationBodyReadinessResult` の所属を persona 側から body 側 projection へ移すための削除（design-diff G-2-b の方針）。
  - `internal/controller/wails/persona_generation_phase_controller.go`: `PersonaGenerationPhaseProjectionDTO` を新設し summary DTO と並列に並べる。`PersonaGenerationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field、command response の `canStartBodyPhase` / `bodyReadinessBlockedReason` を削除。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `BE-term`, `BE-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `internal/usecase/persona_generation_phase_contract.go`
  - `symbol`: `PersonaGenerationPhaseActionEnablement` 型定義
  - `変更種別`: 削除
  - `対応 completion_signal`: 「usecase contract から ActionEnablement 型と CanStartBodyPhase が削除されている」
  - `理由`: 削除起点を明確にし、依存箇所の連鎖修正を直線化する。
- `validation_commands`:
  - `go build ./internal/usecase/... ./internal/controller/wails/...`
  - `go vet ./internal/usecase/... ./internal/controller/wails/...`
- `completion_signal`:
  - usecase contract から `PersonaGenerationPhaseActionEnablement` 型、`Can*` field、`*BlockedReason` field、command result の `CanStartBodyPhase` / `BodyReadinessBlockedReason` が削除されている。
  - projection result が design-diff C-2 の 6 field を満たす。
  - controller DTO に `PersonaGenerationPhaseProjectionDTO` が追加され、`PersonaGenerationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field、command response の persona→body 関連 flag が削除されている。
  - `PersonaGenerationBodyReadinessResult` の persona 側出口（persona controller response / persona usecase return）から readiness が落ちている（body 側に集約される前提）。
  - `go build` / `go vet` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: usecase contract 1 + controller 1 + fetch 経路実装 2-3 ファイル ≈ 4-5 ファイル。
  - 想定変更行数: 削除 300 行前後 + 追加 120 行前後 ≈ 400 changed lines。規模区分は「通常」。
  - 本番経路: `WailsPersonaGenerationPhaseController.Fetch` → `PersonaGenerationPhaseUsecase.Fetch` → projection/summary → JSON DTO。

### `BE-body`:

- `implementation_target`: body phase usecase contract と wails controller DTO から UX 遷移可否 flag と BlockedReason を除去し、`personaBodyReadiness` を含むドメイン状態射影 result/DTO を summary と並列に追加する。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./design-diff.md`（C-3、C-5、G-5-c、G-6-c）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `internal/usecase/body_translation_phase_contract.go`: `BodyTranslationPhaseProjectionResult`（仮称）を新設。field は `PhaseLifecycle` / `JobLifecycle` / `ErrorKind` / `AISettingsConfigured` / `TargetCount` / `PreviousPhaseLifecycle` / `PersonaBodyReadiness{ BodyReadiness bool, SnapshotReferenceStatus string }`。`BodyTranslationPhaseActionEnablement` 型と内部 field を削除。
  - `internal/controller/wails/body_translation_phase_controller.go`: `BodyTranslationPhaseProjectionDTO` を新設し summary DTO と並列に並べる（`personaBodyReadiness` を含む）。`BodyTranslationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field を削除。
  - persona snapshot から readiness を解決する箇所を body 側 usecase に移送（persona usecase が出していた readiness を body usecase で解決する）。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `BE-term`, `BE-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `internal/usecase/body_translation_phase_contract.go`
  - `symbol`: `BodyTranslationPhaseActionEnablement` 型定義
  - `変更種別`: 削除
  - `対応 completion_signal`: 「usecase contract から ActionEnablement 型が削除され、projection result に PersonaBodyReadiness が追加されている」
  - `理由`: 削除起点を明確にし、persona 側からの readiness 引き受け実装を続けて行う。
- `validation_commands`:
  - `go build ./internal/usecase/... ./internal/controller/wails/...`
  - `go vet ./internal/usecase/... ./internal/controller/wails/...`
- `completion_signal`:
  - usecase contract から `BodyTranslationPhaseActionEnablement` 型と内部 field が削除されている。
  - projection result が design-diff C-3 の field を満たし、`PersonaBodyReadiness` が `BodyReadiness` と `SnapshotReferenceStatus` を持つ。
  - controller DTO に `BodyTranslationPhaseProjectionDTO` が追加され、`BodyTranslationPhaseActionEnablementDTO` と summary DTO の `actionEnablement` field が削除されている。
  - body usecase が persona snapshot から readiness を解決して projection に載せている。
  - `go build` / `go vet` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: usecase contract 1 + controller 1 + body usecase 実装 1-2 + persona snapshot 参照 helper 1 ≈ 4-5 ファイル。
  - 想定変更行数: 削除 250 行前後 + 追加 150 行前後 ≈ 400 changed lines。規模区分は「通常」。
  - 本番経路: `WailsBodyTranslationPhaseController.Fetch` → `BodyTranslationPhaseUsecase.Fetch` → projection（personaBodyReadiness を含む）/summary → JSON DTO。

### `INT-term`:

- `implementation_target`: term phase の frontend gateway contract（型と shape）と runtime shape validator を、backend の projection / summary 2 系統に合わせて更新する。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `./design-diff.md`（C-1、C-5、G-4-a）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`: `TermTranslationPhaseProjection` 型（design-diff C-1 の 6 field）を追加。fetch 系 response 型に projection field group を追加。`TermTranslationPhaseActionEnablement` interface、`actionEnablement` field、command response の `canStartNextPhase` / `nextPhaseBlockedReason` を削除。
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts`: runtime shape validator から `actionEnablement` 必須検証と `canStart` / `canPause` / `canResume` / `canRetry` / `canStartNextPhase` 必須検証を削除。projection 必須検証（`phaseLifecycle: string` / `jobLifecycle: string` / `errorKind: string` / `aiSettingsConfigured: boolean` / `aiTargetCount: number` / `confirmedCount: number`）を追加。summary 必須検証は既存維持。
- `depends_on`: `BE-term`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-persona`, `INT-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`
  - `symbol`: `TermTranslationPhaseActionEnablement` interface
  - `変更種別`: 削除（および projection 型追加）
  - `対応 completion_signal`: 「gateway contract から ActionEnablement 型と canStartNextPhase が削除され、TermTranslationPhaseProjection が追加されている」
  - `理由`: gateway contract 型を起点に validator / presenter 入力経路の波及を直線化する。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - gateway contract から `TermTranslationPhaseActionEnablement` 型、fetch response の `actionEnablement` field、command response の `canStartNextPhase` / `nextPhaseBlockedReason` が削除されている。
  - `TermTranslationPhaseProjection` 型が design-diff C-1 の field を満たす。
  - runtime shape validator が projection 必須 field を検証し、`actionEnablement` 必須検証が消えている。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: gateway contract 1 + gateway 1 + 既存 index re-export 1 ≈ 3 ファイル。
  - 想定変更行数: 削除 100 行前後 + 追加 80 行前後 ≈ 180 changed lines。規模区分は「通常」。
  - 本番経路: Wails JSON → `term-translation-phase.gateway.ts` validator → presenter（次 wave）。
  - motivating bug（`canStartNextPhase` 必須検証）が本引き継ぎの validator 改修で解消することを引き継ぎ内で確認する。

### `INT-persona`:

- `implementation_target`: persona phase の frontend gateway contract と runtime shape validator を 2 系統に合わせて更新する。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `./design-diff.md`（C-2、C-5、G-4-b）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-phase-gateway-contract.ts`: `PersonaGenerationPhaseProjection` 型（design-diff C-2 の field）を追加。fetch 系 response 型に projection field group を追加。`PersonaGenerationPhaseActionEnablement` interface、`actionEnablement` field、command response の `canStartBodyPhase` / `bodyReadinessBlockedReason` を削除。persona 側 summary から readiness を運ぶ field を削除（body 側に集約）。
  - `frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-contract.test.ts`: 型 contract に対する static type 確認部分を新 contract に追従。テスト振る舞いの新規設計はしない（追従のみ）。
  - `frontend/src/controller/wails/persona-generation-phase.gateway.ts`: runtime shape validator から `actionEnablement` 必須検証と `canStartBodyPhase` 必須検証を削除。projection 必須検証（`phaseLifecycle` / `jobLifecycle` / `errorKind` / `aiSettingsConfigured` / `targetCount` / `previousPhaseLifecycle`）を追加。
- `depends_on`: `BE-persona`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-term`, `INT-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-phase-gateway-contract.ts`
  - `symbol`: `PersonaGenerationPhaseActionEnablement` interface
  - `変更種別`: 削除（および projection 型追加）
  - `対応 completion_signal`: 「gateway contract から ActionEnablement 型と canStartBodyPhase が削除され、PersonaGenerationPhaseProjection が追加されている」
  - `理由`: 削除起点を明確にし、validator / presenter 入力経路の波及を直線化する。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - gateway contract から `PersonaGenerationPhaseActionEnablement` 型、`actionEnablement` field、command response の `canStartBodyPhase` / `bodyReadinessBlockedReason` が削除されている。
  - `PersonaGenerationPhaseProjection` 型が design-diff C-2 の field を満たす。
  - runtime shape validator が projection 必須 field を検証し、`actionEnablement` / `canStartBodyPhase` 必須検証が消えている。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: gateway contract 1 + 既存 contract test 1（追従）+ gateway 1 + index 1 ≈ 4 ファイル。
  - 想定変更行数: 削除 130 行前後 + 追加 90 行前後 ≈ 220 changed lines。規模区分は「通常」。
  - 本番経路: Wails JSON → `persona-generation-phase.gateway.ts` validator → presenter（次 wave）。

### `INT-body`:

- `implementation_target`: body phase の frontend gateway contract と runtime shape validator を 2 系統に合わせて更新し、`personaBodyReadiness` を projection に取り込む。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `./design-diff.md`（C-3、C-5、G-4-c）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/gateway-contract/body-translation-phase/body-translation-phase-gateway-contract.ts`: `BodyTranslationPhaseProjection` 型（design-diff C-3 の field、`personaBodyReadiness` 含む）を追加。fetch 系 response 型に projection field group を追加。`BodyTranslationPhaseActionEnablement` interface と `actionEnablement` field を削除。
  - `frontend/src/application/gateway-contract/body-translation-phase/body-translation-contract.test.ts`: 型 contract に対する static type 確認部分を追従。
  - `frontend/src/controller/wails/body-translation-phase.gateway.ts`: runtime shape validator から `actionEnablement` 必須検証を削除。projection 必須検証（`phaseLifecycle` / `jobLifecycle` / `errorKind` / `aiSettingsConfigured` / `targetCount` / `previousPhaseLifecycle` / `personaBodyReadiness: { bodyReadiness: boolean, snapshotReferenceStatus: string }`）を追加。
- `depends_on`: `BE-body`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-term`, `INT-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/gateway-contract/body-translation-phase/body-translation-phase-gateway-contract.ts`
  - `symbol`: `BodyTranslationPhaseActionEnablement` interface
  - `変更種別`: 削除（および projection 型追加、personaBodyReadiness 含む）
  - `対応 completion_signal`: 「gateway contract から ActionEnablement 型が削除され、BodyTranslationPhaseProjection（personaBodyReadiness 含む）が追加されている」
  - `理由`: 削除起点を明確にし、validator / presenter 入力経路の波及を直線化する。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - gateway contract から `BodyTranslationPhaseActionEnablement` 型と `actionEnablement` field が削除されている。
  - `BodyTranslationPhaseProjection` 型が design-diff C-3 の field を満たし、`personaBodyReadiness` が `bodyReadiness: boolean` と `snapshotReferenceStatus: string` を持つ。
  - runtime shape validator が projection 必須 field（`personaBodyReadiness` を含む）を検証し、`actionEnablement` 必須検証が消えている。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: gateway contract 1 + 既存 contract test 1（追従）+ gateway 1 + index 1 ≈ 4 ファイル。
  - 想定変更行数: 削除 130 行前後 + 追加 120 行前後 ≈ 250 changed lines。規模区分は「通常」。
  - 本番経路: Wails JSON → `body-translation-phase.gateway.ts` validator → presenter（次 wave）。

### `FE-term`:

- `implementation_target`: term phase presenter の `derive*` 系を projection 入力に切り替え、BlockedReason 固定文字列を frontend 内に置く。表示選択子は summary 入力のまま維持し、ActionCard の出力形と文言を維持する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `./design-diff.md`（B-2、G-1、H-1〜H-5）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`（画面変更 N。design-diff 全体を frontend 引き継ぎ根拠に充てる）
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`: `deriveTermActionEnablement` と `deriveCanStartNextPhase` / `deriveNextPhaseBlockedReason` の入力をドメイン状態射影 (`projection`) に切り替える。判定論理は H-1〜H-5 の論理式に揃え、enum 集合 `TERMINAL_JOB` / `RUNNING_PHASE` / `IDLE_READY_PHASE` / `PAUSED_PHASE` / `RECOVERABLE_FAILED_PHASE` / `COMPLETED_PHASE` を frontend 内定数として持つ。BlockedReason は H 節の固定文字列をそのまま使う。出力 `{ canStart, canPause, canResume, canRetry, *BlockedReason }` の形は維持。
  - 表示選択子（progress / counts / errorMessage / execution）は summary を入力に維持し、ActionCard 生成（`buildActionCards`）の形と文言は変えない。
- `depends_on`: `INT-term`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `FE-persona`, `FE-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
  - `symbol`: `deriveTermActionEnablement`
  - `変更種別`: 入力差し替え（summary → projection）と判定論理整列
  - `対応 completion_signal`: 「`deriveTermActionEnablement` が projection だけを入力に H-1〜H-4 の論理式で導出している」
  - `理由`: 起点関数の入力経路を切り替えることで `deriveCanStartNextPhase` の波及修正へ進める。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - `deriveTermActionEnablement` の入力が projection（`jobLifecycle` / `phaseLifecycle` / `errorKind` / `aiSettingsConfigured` / `aiTargetCount` / `confirmedCount`）だけになっている。summary は読まない。
  - `deriveCanStartNextPhase` / `deriveNextPhaseBlockedReason` の入力が projection だけになっている。
  - BlockedReason 文字列が H-1〜H-5 の固定文字列と一致する。
  - 出力 `{ canStart, canPause, canResume, canRetry }` と次段階遷移 flag の形を維持する。ActionCard 生成の出力形と表示文言は変えない。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: presenter 1 + 共通 enum 定数 helper 1（新規または既存ファイル拡張）≈ 1-2 ファイル。
  - 想定変更行数: 修正 200 行前後（presenter 内の derive 関数まわり）。規模区分は「通常」。
  - frontend 引き継ぎ根拠: 画面変更が無いため `screen_design_diff` は `N/A`。design-diff の H 節を frontend 引き継ぎ根拠に充てる旨を notes に明示する例外運用とする。
  - 単体テスト引き継ぎ（test-design 側）が並列で進む。本引き継ぎでは presenter 既存テストの追従修正が発生する場合があるが、新規テスト設計は test-design 側に渡す。

### `FE-persona`:

- `implementation_target`: persona phase presenter の `derive*` 系を projection 入力に切り替え、BlockedReason 固定文字列を frontend 内に置く。表示選択子は summary 入力のまま維持。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `./design-diff.md`（B-2、B-3、G-2、H-6〜H-11）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`: `derivePersonaActionEnablement` と `derivePersonaCanStartBodyPhase` / `derivePersonaBodyReadinessBlockedReason` の入力を persona projection に切り替える。persona → body 移行の有効化条件は H-11 の論理式（terminal でない ∧ phaseLifecycle ∈ COMPLETED_PHASE）に限定し、persona snapshot 参照可否（personaBodyReady）は本 presenter で判定しない（body 側 presenter に集約）。BlockedReason は H-6〜H-11 の固定文字列を使う。出力 `{ canStart, canPause, canResume, canRetry, canCancel, *BlockedReason }` の形は維持。
- `depends_on`: `INT-persona`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `FE-term`, `FE-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`
  - `symbol`: `derivePersonaActionEnablement`
  - `変更種別`: 入力差し替えと判定論理整列
  - `対応 completion_signal`: 「`derivePersonaActionEnablement` が persona projection だけを入力に H-6〜H-10 の論理式で導出している」
  - `理由`: 起点関数の入力経路を切り替えることで `derivePersonaCanStartBodyPhase` 等の波及修正に進める。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - `derivePersonaActionEnablement` の入力が persona projection（`jobLifecycle` / `phaseLifecycle` / `errorKind` / `aiSettingsConfigured` / `targetCount` / `previousPhaseLifecycle`）だけになっている。summary は読まない。
  - `derivePersonaCanStartBodyPhase` の入力が persona projection だけになり、有効化条件が H-11 の論理式と一致する（personaBodyReady を含まない）。
  - BlockedReason 文字列が H-6〜H-11 の固定文字列と一致する。
  - 出力 `{ canStart, canPause, canResume, canRetry, canCancel }` と次段階遷移 flag の形を維持する。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: presenter 1 + 共通 enum 定数 helper 利用 1 ≈ 1-2 ファイル。
  - 想定変更行数: 修正 250 行前後。規模区分は「通常」。
  - 画面変更なしのため `screen_design_diff` は `N/A`。design-diff の H 節を frontend 引き継ぎ根拠に充てる。

### `FE-body`:

- `implementation_target`: body phase presenter の `derive*` 系を projection 入力に切り替え、`personaBodyReadiness` を本 presenter で判定する。BlockedReason 固定文字列を frontend 内に置く。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `./design-diff.md`（B-2、B-3、G-3、H-12〜H-16）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: not_required
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`: `deriveBodyActionEnablement` の入力を body projection に切り替える。`personaBodyReady` 判定を本 presenter で実施（H-12 の論理式: `personaBodyReadiness.bodyReadiness = true ∨ personaBodyReadiness.snapshotReferenceStatus = "available"`）。BlockedReason は H-12〜H-16 の固定文字列を使う。出力 `{ canStart, canPause, canResume, canRetry, canCancel }` の形は維持。
- `depends_on`: `INT-body`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `FE-term`, `FE-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - `path`: `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`
  - `symbol`: `deriveBodyActionEnablement`
  - `変更種別`: 入力差し替えと判定論理整列（personaBodyReady 判定の body 側集約）
  - `対応 completion_signal`: 「`deriveBodyActionEnablement` が body projection だけを入力に H-12〜H-16 の論理式で導出し、personaBodyReady を本関数内で判定している」
  - `理由`: 起点関数の入力経路と readiness 判定の所属切り替えを最初に閉じる。
- `validation_commands`:
  - `npm --prefix frontend run check`
- `completion_signal`:
  - `deriveBodyActionEnablement` の入力が body projection（`jobLifecycle` / `phaseLifecycle` / `errorKind` / `aiSettingsConfigured` / `targetCount` / `previousPhaseLifecycle` / `personaBodyReadiness`）だけになっている。summary は読まない。
  - `personaBodyReady` 判定が本 presenter 内で実施され、persona 側 presenter からは消えている（FE-persona と整合）。
  - BlockedReason 文字列が H-12〜H-16 の固定文字列と一致する。
  - 出力 `{ canStart, canPause, canResume, canRetry, canCancel }` の形を維持する。startNextPhase 系は持たない。
  - `npm --prefix frontend run check` が通過する。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`:
  - 想定変更ファイル数: presenter 1 + 共通 enum 定数 helper 利用 1 ≈ 1-2 ファイル。
  - 想定変更行数: 修正 280 行前後。規模区分は「通常」。
  - 画面変更なしのため `screen_design_diff` は `N/A`。design-diff の H 節を frontend 引き継ぎ根拠に充てる。

## 範囲外（本 implementation-scope に含めない）

- テスト設計と単体テスト / シナリオテスト引き継ぎ: plan.md 指示で test-design 側に分離。
- docs 正本化（`docs/coding-guidelines-backend.md` 第 7 節以外の追加正本化）: finalization-module で扱う。
- Storybook 表示変更: 画面変更が無いため不要。

## Completion Packet

実装モジュールは完了時に次を返す。

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
- `completion_evidence`
- `telemetry_events`
- `docs_changes: none`
