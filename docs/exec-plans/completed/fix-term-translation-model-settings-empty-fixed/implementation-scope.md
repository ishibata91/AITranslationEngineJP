# Implementation Scope: fix-term-translation-model-settings-empty-fixed

- `skill`: implementation-scope
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `human_review_status`: 詳細仕様差分の Q-001/Q-002/Q-003/Q-004 を 2026-06-01 に人間承認。`detail-spec-diff.md` の未決 0 件。
- `approval_record`: `./detail-spec-diff.md` 回答欄（Q-001/002/003/004、回答日 2026-06-01、回答者: 人間レビュー）
- `module_entry`: `.claude/skills/implementation-module/SKILL.md`
- `handoff_runtime`: `claude-module`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`（採用案 B、`er-REQ-001`、`er-REQ-002`、`term-translation-phase-REQ-002`、`term-translation-phase-REQ-007`、`persona-generation-phase-REQ-002`、`body-translation-phase-REQ-002`）
- `screen_design_diff`: `N/A`（画面構造変更なし。状態別表示は画面設計正本で覆える）
- `design_diff_diagram`: `./design-diff.md`（4 局面 flowchart + 4 場面シーケンス図）
- `fix_decision`: `./fix-decision.md`（確定原因、禁止修正 1〜4 を維持、禁止修正 5 は Q-001 回答で覆る）
- `missing_tests`: `./missing-tests.md`（E2E-UC-FIX-MODEL-001/002/003）
- `data_testid_gaps`: `./data-testid-gaps.md`（SEL-FIX-MODEL-001/002/003）

## Fixed Decisions

- 詳細仕様差分の人間レビュー回答 4 件はすべて承認済み。未決 0 件。
- `JOB_PHASE_AI_SETTINGS` は 3 フェーズ種別共通の汎用テーブル。主キー `phase_type` のみ、`job_id` 列なし、`user_id` 列なし、cascade なし、明示削除 API なし、upsert のみ、3 件のみ存在し得る。保持列は AI サービス、モデル、処理方式、`batchMode` の 4 値であり、`credential_ref` は保持しない。
- `JOB_PHASE_RUN` の AI 設定列（`ai_provider`、`model_name`、`execution_mode`、`credential_ref`）は実行中固定値に責務を限定する。フェーズ開始時に `JOB_PHASE_AI_SETTINGS` から 3 値を、provider-settings 正本から `credential_ref` を都度解決して転写する。
- `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の AI 設定保持責務は廃止する。残責務（観測 snapshot 等）が存在しない場合は table 廃止、残責務がある場合は AI 設定列だけを drop して縮退する。判断は実装モジュールが migration 016 以降の利用箇所走査結果で固定する。
- backend 応答は「設定全体の情報」または「該当状態が存在しない」の二択。派生説明文字列（blocked reason、設定不足の理由）は backend 応答に含めない。
- `SaveAISettings` 入力から `job_id` を抜き、`phase_type` + AI 選択値（provider、model、executionMode、batchMode）のみで構成する。
- `unanswered_questions`: `0`
- 3 フェーズ（term-translation / persona-generation / body-translation）を同時実装範囲とする（Q-001 で禁止修正 5 が覆る）。
- secret 境界: `credential_ref` は参照値（識別子）として `JOB_PHASE_RUN`、UI 表示、DTO に出してよい。secret 本体（API key 等）は provider-settings 正本側で保持され、本 task のいずれの handoff にも持ち込まない。
- 公開接点（Wails DTO 構造）の変更は backend / 統合境界 / frontend の各 handoff の完了条件として共有する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `H-MIG-ER` | なし | なし | なし |
| `wave-2` | `H-BE-TERM`, `H-BE-PERSONA`, `H-BE-BODY` | `H-MIG-ER` | `H-BE-TERM <-> H-BE-PERSONA <-> H-BE-BODY` | なし |
| `wave-3` | `H-INT-PHASE-AI-SETTINGS` | `H-BE-TERM`, `H-BE-PERSONA`, `H-BE-BODY` | なし | DTO contract が 3 サービス共通のため統合境界は単一 |
| `wave-4` | `H-FE-PRESENTER` | `H-INT-PHASE-AI-SETTINGS` | なし | UI ロジックは Wails DTO 形に従う |
| `wave-5` | `H-FE-SELECTOR` | `H-FE-PRESENTER` | なし | E2E selector は UI ロジック完了後 |
| `wave-6` | `H-TEST-UNIT`, `H-TEST-SCENARIO` | `H-FE-SELECTOR` | `H-TEST-UNIT <-> H-TEST-SCENARIO` | なし |

並列実行可能性の根拠:
- `H-BE-TERM`/`H-BE-PERSONA`/`H-BE-BODY` は同じ ER migration 完了を依存とするが、変更対象 service ファイル（`term_translation_phase_service.go` / `persona_generation_phase_service.go` / `body_translation_phase_service.go`）が分離されており、`owned_scope` は重ならない。共有変更となる `internal/service/provider_execution_snapshot.go`（`savePhaseAISettings`）と repository 層 (`JOB_PHASE_AI_SETTINGS` repository) は `H-MIG-ER` で固定する。
- `H-TEST-UNIT` と `H-TEST-SCENARIO` は実装完了後の別成果物として並列実行可能。

## Handoffs

### `H-MIG-ER`:

- `implementation_target`: ER migration の新規追加（`JOB_PHASE_AI_SETTINGS` テーブル新設、`TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の AI 設定保持責務廃止）と repository 層の共通実装の固定。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`（`er-REQ-001`、`er-REQ-002`）, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`（migration 範囲は AI 設定参照値 3 値のみ。secret 本体は provider-settings 正本側に残る）
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credential_ref`（参照値）
  - `secret_values_for_provider_external_api_internal_auth`: secret 本体は本 handoff の永続化対象外（provider-settings 正本側で管理）
  - `secret_resolution_owner_layer`: provider-settings 正本 service 層（本 handoff は参照のみ）
  - `forbidden_outputs`: API key 本体、credential 本体値を `JOB_PHASE_AI_SETTINGS` の列、migration の seed、log、error summary に出さない。
- `owned_scope`:
  - `internal/infra/sqlite/dbinit/migrations/` 配下に `017_job_phase_ai_settings.sql`（仮番号、実装時は最新採番）を追加し、`JOB_PHASE_AI_SETTINGS`（PK=`phase_type`、列: `phase_type`, `ai_provider`, `model_name`, `execution_mode`, `batch_mode`, `updated_at`）を作成する。
  - `internal/infra/sqlite/migrations/` 配下にも同名 / 同 DDL の migration を追加する（dbinit と並行構造に従う）。
  - `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の AI 設定列の扱いを固定する。実装初手で `provider_execution_snapshot.go` と repository 層の参照を走査し、(a) 他責務が存在しない → table 廃止 DROP migration を追加、(b) 他責務（観測 snapshot 等）が存在する → AI 設定列のみ DROP する migration を追加、いずれかに分岐する。
  - `internal/repository/` 配下に `JOB_PHASE_AI_SETTINGS` の repository interface と SQLite 実装、InMemory 実装、unit test スケルトン（`SaveAISettings`/`upsert`、`LoadByPhaseType`、`PhaseType` 不在時 `ErrNotFound`）を追加する。
  - `JOB_PHASE_AI_SETTINGS` の DDL は 3 フェーズ共通のため repository 層も汎用化する（`phase_type` を受け取る）。
- `depends_on`: なし
- `execution_group`: wave-1
- `ready_wave`: wave-1
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`:
  - path: `internal/infra/sqlite/dbinit/migrations/017_job_phase_ai_settings.sql`（実装時は最新採番に合わせる）
  - 対象単位: `JOB_PHASE_AI_SETTINGS` テーブルの `CREATE TABLE`
  - 変更種別: 新規追加
  - 対応 `completion_signal` clause: 「`JOB_PHASE_AI_SETTINGS` テーブルが PK=`phase_type` のみで作成されている」
  - 1 手目にする理由: 後続 backend handoff の永続化対象を最初に固定するため、ER 正本の対象テーブルを先に作る。
- `validation_commands`:
  - `go test ./internal/repository/... -run JobPhaseAISettings`
  - `go test ./internal/infra/sqlite/...`
- `completion_signal`:
  - `JOB_PHASE_AI_SETTINGS` が dbinit / migrations の両系列に追加されている。
  - PK が `phase_type` のみであり、`job_id`/`user_id` 列が含まれない。
  - `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の AI 設定列削除または table 廃止が、走査結果に応じて適切に migration として追加されている。
  - repository interface に `JOB_PHASE_AI_SETTINGS` の `Save`/`Load`/`ErrNotFound` が定義され、SQLite 実装と InMemory 実装が揃っている。
  - 単体テスト（upsert と `phase_type` 不在時の `ErrNotFound`）が通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff は 3 フェーズ共通の永続化基盤として一度だけ実行する。3 フェーズの service 側変更は wave-2 で分割する。
  - `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の縮退方針は実装初手で走査して確定する。判断結果は `completion_evidence` に残す。
  - 本番経路: SQLite migration（dbinit 起動時適用） → `JobPhaseAISettingsRepository` → 3 phase service の `savePhaseAISettings` / フェーズ開始時の `JOB_PHASE_RUN` 転写経路（後続 handoff で接続）

### `H-BE-TERM`:

- `implementation_target`: 単語翻訳フェーズの backend service（`internal/service/term_translation_phase_service.go`）と `provider_execution_snapshot.go` の `savePhaseAISettings`、Wails controller を、新 ER に従い書き換える。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-002`、`term-translation-phase-REQ-007`）, `docs/detail-specs/term-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credential_ref`（参照値識別子）
  - `secret_values_for_provider_external_api_internal_auth`: provider 認証の本体値（API key 等）は provider-settings 正本側で保持
  - `secret_resolution_owner_layer`: provider-settings service（`provider_settings_service.go`）
  - `forbidden_outputs`: API key 本体、認証 token 本体を controller DTO、service 戻り値、log、error summary、URL に出さない。`credential_ref` 自体は参照値として許容する。
- `owned_scope`:
  - `internal/service/term_translation_phase_service.go`:
    - `applyTermTranslationRuntimeSnapshot` の AI 設定空文字上書き経路を廃止する。Ready 期表示は `JobPhaseAISettingsRepository.Load(phase_type="word_translation")` で取得し、record 不在は応答の AI 設定 field 不在で表現する。
    - フェーズ開始経路で `JOB_PHASE_AI_SETTINGS` record の 3 値を取得し、provider-settings から `credential_ref` を都度解決して `JOB_PHASE_RUN` を新規作成する。record 不在時または provider-settings 側解決失敗時は開始拒否を返す（派生理由文字列は返さない）。
    - 開始可否 / 状態判定の派生説明文字列を service 戻り値から削除する。
  - `internal/service/provider_execution_snapshot.go`:
    - `savePhaseAISettings` の書き込み先を `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` から `JobPhaseAISettingsRepository` に切り替える。入力から `job_id` を抜く。upsert 対象は `provider`/`model`/`executionMode`/`batchMode` のみとし、`credential_ref` を含めない。
  - `internal/controller/wails/term_translation_phase_controller.go` 系の DTO 整合（必要な部分）。応答 DTO で AI 設定 field の不在表現（field を含めない、または null）を確定する。
- `depends_on`: `H-MIG-ER`
- `execution_group`: wave-2
- `ready_wave`: wave-2
- `parallelizable_with`: `H-BE-PERSONA`, `H-BE-BODY`
- `parallel_blockers`: なし
- `first_action`:
  - path: `internal/service/term_translation_phase_service.go`
  - 対象単位: `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` ブロック（1468-1473 行付近）
  - 変更種別: 廃止と書き換え（`JobPhaseAISettingsRepository.Load` を呼び、record 不在は応答 AI 設定 field 不在で表現する経路へ置換）
  - 対応 `completion_signal` clause: 「Ready 期表示時、`JOB_PHASE_AI_SETTINGS` record 不在で AI 設定 field を不在として返す」
  - 1 手目にする理由: 確定原因の中核（空文字上書き）を最初に閉じる。
- `validation_commands`:
  - `go test ./internal/service/ -run TermTranslationPhaseService`
  - `go test ./internal/controller/wails/ -run TermTranslationPhaseController`
- `completion_signal`:
  - `applyTermTranslationRuntimeSnapshot` の空文字上書きが廃止されている。
  - Ready 期表示で `JobPhaseAISettingsRepository.Load("word_translation")` を参照し、record 不在は AI 設定 field 不在で表現される。
  - フェーズ開始時に `JOB_PHASE_AI_SETTINGS` record + provider-settings 都度解決から `JOB_PHASE_RUN` へ転写する経路が成立する。
  - `savePhaseAISettings` の書き込み先が `JOB_PHASE_AI_SETTINGS` であり、入力に `job_id` が含まれない。
  - controller DTO で派生説明文字列（blocked reason、設定不足の理由）を返していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `H-BE-PERSONA` / `H-BE-BODY` と同じ ER 基盤を共有するが、変更ファイルは独立しているため並列実行できる。
  - 本番経路: Wails controller → `TermTranslationPhaseService.LoadSummary` → `JobPhaseAISettingsRepository` / `ProviderSettingsService` → DTO

### `H-BE-PERSONA`:

- `implementation_target`: ペルソナ生成フェーズの backend service（`internal/service/persona_generation_phase_service.go`）と controller DTO を、新 ER に従い `H-BE-TERM` と同型に書き換える。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`（`persona-generation-phase-REQ-002`）, `docs/detail-specs/persona-generation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credential_ref`
  - `secret_values_for_provider_external_api_internal_auth`: provider 認証の本体値は provider-settings 正本側
  - `secret_resolution_owner_layer`: provider-settings service
  - `forbidden_outputs`: API key 本体、認証 token 本体を DTO/log/error summary/URL に出さない。
- `owned_scope`:
  - `internal/service/persona_generation_phase_service.go`:
    - `applyPersonaGenerationRuntimeSnapshot` の AI 設定空文字上書き経路を廃止し、`JobPhaseAISettingsRepository.Load(phase_type="npc_persona_generation")` を参照する経路へ置換する。
    - フェーズ開始経路で `JOB_PHASE_AI_SETTINGS` record + provider-settings 解決値から `JOB_PHASE_RUN` を新規作成する。
    - 派生説明文字列を service 戻り値から削除する。
  - `internal/controller/wails/` の persona_generation 関連 DTO の整合。
- `depends_on`: `H-MIG-ER`
- `execution_group`: wave-2
- `ready_wave`: wave-2
- `parallelizable_with`: `H-BE-TERM`, `H-BE-BODY`
- `parallel_blockers`: なし
- `first_action`:
  - path: `internal/service/persona_generation_phase_service.go`
  - 対象単位: `applyPersonaGenerationRuntimeSnapshot`（896-900 行付近）の `ErrNotFound` ブロック
  - 変更種別: 廃止と書き換え
  - 対応 `completion_signal` clause: 「Ready 期表示時、`JOB_PHASE_AI_SETTINGS`（`phase_type="npc_persona_generation"`）record 不在で AI 設定 field を不在として返す」
  - 1 手目にする理由: 同型修正の最初の clause を閉じる。
- `validation_commands`:
  - `go test ./internal/service/ -run PersonaGenerationPhaseService`
  - `go test ./internal/controller/wails/ -run PersonaGeneration`
- `completion_signal`:
  - `applyPersonaGenerationRuntimeSnapshot` の空文字上書きが廃止されている。
  - Ready 期表示で `JobPhaseAISettingsRepository.Load("npc_persona_generation")` を参照する。
  - フェーズ開始時の `JOB_PHASE_RUN` 転写経路が成立する。
  - `savePhaseAISettings` 経路が新 ER 経由になり、入力に `job_id` を含まない。
  - 派生説明文字列が controller DTO から落ちている。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `H-BE-TERM` と同型変更。共有変更（`provider_execution_snapshot.go` の `savePhaseAISettings`）は `H-BE-TERM` の `owned_scope` で更新済みである前提で接続する。`savePhaseAISettings` 自体は `H-MIG-ER` 完了後に汎用化されているため、phase service 側はそれを `phase_type` 引数で呼ぶだけになる。
  - 本番経路: Wails controller → `PersonaGenerationPhaseService.LoadSummary` → `JobPhaseAISettingsRepository` / `ProviderSettingsService` → DTO

### `H-BE-BODY`:

- `implementation_target`: 本文翻訳フェーズの backend service（`internal/service/body_translation_phase_service.go`）と controller DTO を、新 ER に従い `H-BE-TERM` と同型に書き換える。
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`（`body-translation-phase-REQ-002`）, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credential_ref`
  - `secret_values_for_provider_external_api_internal_auth`: provider 認証の本体値は provider-settings 正本側
  - `secret_resolution_owner_layer`: provider-settings service
  - `forbidden_outputs`: API key 本体、認証 token 本体を DTO/log/error summary/URL に出さない。
- `owned_scope`:
  - `internal/service/body_translation_phase_service.go`:
    - 同型の `applyBody...RuntimeSnapshot`（764-768 行付近）の `ErrNotFound` ブロックを廃止し、`JobPhaseAISettingsRepository.Load(phase_type="text_translation")` 経路へ置換する。
    - フェーズ開始経路で `JOB_PHASE_AI_SETTINGS` record + provider-settings 解決値から `JOB_PHASE_RUN` を新規作成する。
    - 派生説明文字列を service 戻り値から削除する。
  - `internal/controller/wails/` の body_translation 関連 DTO の整合。
- `depends_on`: `H-MIG-ER`
- `execution_group`: wave-2
- `ready_wave`: wave-2
- `parallelizable_with`: `H-BE-TERM`, `H-BE-PERSONA`
- `parallel_blockers`: なし
- `first_action`:
  - path: `internal/service/body_translation_phase_service.go`
  - 対象単位: 同型 `applyBody...RuntimeSnapshot`（764-768 行付近）の `ErrNotFound` ブロック
  - 変更種別: 廃止と書き換え
  - 対応 `completion_signal` clause: 「Ready 期表示時、`JOB_PHASE_AI_SETTINGS`（`phase_type="text_translation"`）record 不在で AI 設定 field を不在として返す」
  - 1 手目にする理由: 同型修正の最初の clause を閉じる。
- `validation_commands`:
  - `go test ./internal/service/ -run BodyTranslationPhaseService`
  - `go test ./internal/controller/wails/ -run BodyTranslation`
- `completion_signal`:
  - `applyBody...RuntimeSnapshot` の空文字上書きが廃止されている。
  - Ready 期表示で `JobPhaseAISettingsRepository.Load("text_translation")` を参照する。
  - フェーズ開始時の `JOB_PHASE_RUN` 転写経路が成立する。
  - `savePhaseAISettings` 経路が新 ER 経由で、入力に `job_id` を含まない。
  - 派生説明文字列が controller DTO から落ちている。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本番経路: Wails controller → `BodyTranslationPhaseService.LoadSummary` → `JobPhaseAISettingsRepository` / `ProviderSettingsService` → DTO

### `H-INT-PHASE-AI-SETTINGS`:

- `implementation_target`: 3 フェーズ共通の Wails DTO 形と gateway 接続（応答 DTO の「AI 設定 field 不在」「`execution` field 不在」の二択化、`SaveAISettings` 入力からの `job_id` 削除）の境界実装。frontend gateway 型の更新を含む。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-002`、`persona-generation-phase-REQ-002`、`body-translation-phase-REQ-002` の「backend 応答の構造に関する仕様」「`SaveAISettings` 入力構造に関する仕様」）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credential_ref`
  - `secret_values_for_provider_external_api_internal_auth`: なし（境界 DTO は参照値のみ）
  - `secret_resolution_owner_layer`: provider-settings service（境界では解決しない）
  - `forbidden_outputs`: API key 本体、認証 token 本体を DTO/log/URL に出さない。境界 DTO に派生説明文字列（blocked reason、設定不足理由）を含めない。
- `owned_scope`:
  - Wails DTO 定義（`internal/controller/wails/*phase_controller.go`、`frontend/wailsjs/go/*`、`frontend/src/infrastructure/gateway/` 配下の term/persona/body 用 gateway 型）の境界整合。
  - 応答 DTO で `execution` field と AI 設定 field を optional（不在で「未設定」を表現）にする。
  - `SaveAISettings` 入力 DTO 形を 3 フェーズ共通で `phase_type` + AI 選択値（provider/model/executionMode/batchMode）に揃える（`job_id` を削除）。
  - 実画面確認: `npm run dev:wails:run` で起動した実 app で、未保存ジョブを開いた時に AI 設定 field 不在の DTO が frontend に届くこと、保存後に field が含まれることを `chrome-devtools` MCP で観測する。
- `depends_on`: `H-BE-TERM`, `H-BE-PERSONA`, `H-BE-BODY`
- `execution_group`: wave-3
- `ready_wave`: wave-3
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`（3 フェーズの DTO 形を共通で揃えるため単一 handoff にまとめる）
- `first_action`:
  - path: `frontend/src/infrastructure/gateway/`（term-translation 用 gateway 型の起点ファイル）
  - 対象単位: 応答 DTO 型の `execution` field を optional にする型修正
  - 変更種別: 型変更（optional 化）
  - 対応 `completion_signal` clause: 「frontend gateway 型で `execution` field と AI 設定 field が optional であり、不在で『未設定』を表現する」
  - 1 手目にする理由: DTO 境界の最も狭い接点を先に閉じ、後続の presenter 修正の依存を確定する。
- `validation_commands`:
  - `go test ./internal/controller/wails/...`
  - `npm --prefix frontend run test:unit -- --run gateway`
- `completion_signal`:
  - 3 フェーズの応答 DTO で AI 設定 field と `execution` field が optional であり、不在で「未設定」を表現できる。
  - `SaveAISettings` 入力 DTO に `job_id` が含まれていない。
  - frontend gateway の型定義と Wails 生成型が一致する。
  - 実画面確認: 未保存ジョブの単語翻訳画面で、応答 DTO の `execution` field が不在で返ることを `chrome-devtools` で確認できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`（境界 DTO 観測）
- `execution_stage`: `実装後`
- `notes`:
  - 派生説明文字列（blocked reason、設定不足理由）を境界 DTO に含めない。
  - 本番経路: Wails controller DTO → `frontend/wailsjs/go/*` 生成型 → `frontend/src/infrastructure/gateway/` 型 → presenter 入力

### `H-FE-PRESENTER`:

- `implementation_target`: 単語翻訳 / ペルソナ生成 / 本文翻訳の 3 フェーズ presenter で、判定根拠を `execution` field の有無 + AI 設定 field の有無へ切り替える。派生説明（「設定未完了」「AI 設定不足」「モデルを選択してください」）の組み立てを frontend 側に集約する。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-002`/`-007`、`persona-generation-phase-REQ-002`、`body-translation-phase-REQ-002`）, `docs/screen-design/screens/term-translation-phase.md:91-129`（既存正本）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`（画面構造変更なし。既存状態別表示で覆える）
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: `provider`, `model`, `executionMode`, `batchMode`, `credentialRefLabel`（参照値表示）
  - `secret_values_for_provider_external_api_internal_auth`: なし（frontend は本体値を扱わない）
  - `secret_resolution_owner_layer`: provider-settings service（frontend は呼び出すだけ）
  - `forbidden_outputs`: API key 本体、認証 token 本体を viewModel/log/error UI に出さない。
- `owned_scope`:
  - `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`:
    - `isExecutionConfigured` の空文字 `trim()` 判定を廃止し、`execution` field と AI 設定 field の有無で判定する形に置き換える。
    - `providerLabel` / `modelLabel` / `executionModeLabel` の `?? "-"` フォールバックは、field 不在時に「設定未完了」相当の派生表示語へ移す。
    - `aiSettingsBlockedReason` 相当の派生語の組み立てを frontend 側で行う。
  - `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`: 同型修正。
  - `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`: 同型修正。
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` の `viewModel.modelLabel === "-"` 系判定（61-73 行付近）を、presenter の新しい `isExecutionConfigured` / 派生 boolean に切り替える形で受け取る。同様に persona / body の Panel も追従させる（範囲が同型で軽微であれば本 handoff に含める）。
  - 画面構造、layout、文言の新規追加はしない。状態 pill の「固定済み」「設定未完了」の 2 値を維持する。
- `depends_on`: `H-INT-PHASE-AI-SETTINGS`
- `execution_group`: wave-4
- `ready_wave`: wave-4
- `parallelizable_with`: なし
- `parallel_blockers`: `backend_frontend_order`（DTO 形確定後に着手）
- `first_action`:
  - path: `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
  - 対象単位: `isExecutionConfigured` 関数（346 行付近）
  - 変更種別: 判定ロジック書き換え（空文字 `trim()` を `execution` field 有無に置換）
  - 対応 `completion_signal` clause: 「presenter は `execution` field の有無で `isExecutionConfigured` を判定し、空文字 trim 判定を廃止する」
  - 1 手目にする理由: 確定原因の症状連鎖の起点を最初に閉じる。
- `validation_commands`:
  - `npm --prefix frontend run test:unit -- --run term-translation-phase.presenter`
  - `npm --prefix frontend run test:unit -- --run persona-generation-phase.presenter`
  - `npm --prefix frontend run test:unit -- --run body-translation-phase.presenter`
- `completion_signal`:
  - 3 presenter で `isExecutionConfigured` の判定が `execution` field 有無に基づく。
  - field 不在時に「設定未完了」相当の派生表示語が presenter で組み立てられる（backend から派生語を受け取らない）。
  - 単語翻訳 / ペルソナ生成 / 本文翻訳の Panel 側で `viewModel.modelLabel === "-"` 依存が廃止され、presenter の派生 boolean を参照する形になる。
  - 既存の画面構造、状態 pill 2 値、layout、文言、style に新規追加が無い。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 表示変更（layout、style、文言、状態値追加）は実装外。画面設計差分なしで覆える。
  - 本番経路: gateway DTO → presenter → viewModel → Svelte Panel

### `H-FE-SELECTOR`:

- `implementation_target`: 不足セレクタ SEL-FIX-MODEL-001/002/003 の `data-testid` を 3 フェーズの Panel に追加する。E2E 観点の安定特定基盤を整える。
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `spec_basis`: `./data-testid-gaps.md`、`./missing-tests.md`、`docs/screen-design/screens/term-translation-phase.md:91-129`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: N/A
  - `secret_values_for_provider_external_api_internal_auth`: N/A
  - `secret_resolution_owner_layer`: N/A
  - `forbidden_outputs`: N/A
- `owned_scope`:
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`:
    - `data-testid="term-translation-phase-ai-settings-status-pill"`（SEL-FIX-MODEL-001）
    - `data-testid="term-translation-phase-start-blocked-reason"`（SEL-FIX-MODEL-002）
    - `data-testid="term-translation-phase-ai-provider-select"` / `term-translation-phase-ai-model-select` / `term-translation-phase-ai-execution-mode-select`（SEL-FIX-MODEL-003）
  - persona-generation / body-translation の Panel にも同型の `data-testid` を追加する（名前 prefix は `persona-generation-phase-` / `body-translation-phase-`）。
  - 状態 pill / 開始禁止理由領域 / 3 つの select 要素以外には `data-testid` を追加しない。
- `depends_on`: `H-FE-PRESENTER`
- `execution_group`: wave-5
- `ready_wave`: wave-5
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`:
  - path: `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
  - 対象単位: 状態 pill 要素
  - 変更種別: 属性追加（`data-testid="term-translation-phase-ai-settings-status-pill"`）
  - 対応 `completion_signal` clause: 「単語翻訳 Panel の状態 pill に `data-testid=term-translation-phase-ai-settings-status-pill` が付与されている」
  - 1 手目にする理由: 後続シナリオテストの最も多用される selector を先に固定する。
- `validation_commands`:
  - `npm --prefix frontend run test:unit -- --run TermTranslationPhasePanel`
  - `npm --prefix frontend run check`
- `completion_signal`:
  - 3 Panel に SEL-FIX-MODEL-001/002/003 相当の `data-testid` が付与されている。
  - layout / 文言 / style / 状態値追加を伴っていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff は `H-TEST-SCENARIO` の前提となるため、selector 命名は data-testid-gaps.md の例示に従う。

### `H-TEST-UNIT`:

- `implementation_target`: 新規 repository、3 phase service、3 presenter の単体テストを追加し、空文字代理表現の廃止と field 不在表現を期待値として固定する。
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `spec_basis`: `./detail-spec-diff.md`（採用案 B 全 REQ）、`docs/detail-specs/term-translation-phase.md`、`docs/detail-specs/persona-generation-phase.md`、`docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`（テスト fixture では参照値のみを使い、本体値は使わない）
  - `reference_values_allowed_in_ui_dto_read_model`: 参照値 fixture のみ
  - `secret_values_for_provider_external_api_internal_auth`: なし
  - `secret_resolution_owner_layer`: provider-settings fake
  - `forbidden_outputs`: API key 本体や認証 token 本体をテスト fixture に含めない。
- `owned_scope`:
  - 新規 `JobPhaseAISettingsRepository` 単体テスト（upsert、`phase_type` 不在時 `ErrNotFound`、cascade なし、3 phase_type 並存）。
  - 3 phase service 単体テスト（Ready 期 record 不在で AI 設定 field 不在応答、record 存在で値応答、フェーズ開始時の `JOB_PHASE_RUN` 転写、開始拒否の二択（record 不在 / provider-settings 解決不可）、派生説明文字列を返さない）。
  - 3 presenter 単体テスト（`execution` field 不在で `isExecutionConfigured=false`、`isExecutionConfigured=false` 時の派生表示語が presenter で組み立てられる）。
- `depends_on`: `H-FE-SELECTOR`
- `execution_group`: wave-6
- `ready_wave`: wave-6
- `parallelizable_with`: `H-TEST-SCENARIO`
- `parallel_blockers`: なし
- `first_action`:
  - path: `internal/repository/job_phase_ai_settings_sqlite_repository_test.go`（仮）
  - 対象単位: 「`phase_type` 不在時に `ErrNotFound`」テスト
  - 変更種別: 新規追加
  - 対応 `completion_signal` clause: 「`JobPhaseAISettingsRepository` の `phase_type` 不在時に `ErrNotFound` が返るテストが通過する」
  - 1 手目にする理由: 永続化基盤の不在表現を最初に証明する。
- `validation_commands`:
  - `go test ./internal/repository/... ./internal/service/...`
  - `npm --prefix frontend run test:unit`
- `completion_signal`:
  - repository の upsert と `ErrNotFound` テストが通過する。
  - 3 phase service の field 不在応答と転写経路テストが通過する。
  - 3 presenter の判定切り替えテストが通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 期待結果の元ネタは `spec_basis` の本差分 REQ と既存 detail-specs。

### `H-TEST-SCENARIO`:

- `implementation_target`: missing-tests.md の E2E-UC-FIX-MODEL-001/002/003 をシナリオテスト（UI 人間操作 E2E）として追加し、`H-FE-SELECTOR` の `data-testid` を使って状態 pill と select の挙動を証明する。
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `./missing-tests.md`、`./detail-spec-diff.md`（`term-translation-phase-REQ-002`/`-007`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: 参照値 fixture のみ
  - `secret_values_for_provider_external_api_internal_auth`: なし
  - `secret_resolution_owner_layer`: provider-settings fake / mock
  - `forbidden_outputs`: API key 本体、認証 token 本体を fixture / 観測ログに残さない。
- `owned_scope`:
  - 単語翻訳画面の E2E-UC-FIX-MODEL-001/002/003 を Playwright 等の既存 E2E 基盤に追加する。
  - 状態 pill の「設定未完了」表示、開始ボタンの `disabled` 維持と禁止理由テキスト、AI サービス / モデル / 処理方式の選択と「固定済み」遷移を証明する。
- `depends_on`: `H-FE-SELECTOR`
- `execution_group`: wave-6
- `ready_wave`: wave-6
- `parallelizable_with`: `H-TEST-UNIT`
- `parallel_blockers`: なし
- `first_action`:
  - path: 既存 E2E スイート配下（実装時に確認）
  - 対象単位: E2E-UC-FIX-MODEL-001（未実行ジョブで状態 pill が「設定未完了」を表示する）
  - 変更種別: 新規シナリオ追加
  - 対応 `completion_signal` clause: 「E2E-UC-FIX-MODEL-001 が通過する」
  - 1 手目にする理由: 確定原因の最も外側からの証明を最初に置く。
- `validation_commands`:
  - 既存 E2E 実行コマンド（実装時に確認）
- `completion_signal`:
  - E2E-UC-FIX-MODEL-001/002/003 がすべて通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - 観点はすべて単語翻訳画面に閉じる。persona / body の同型 E2E 追加は本 task の範囲外として `missing-tests.md` の補足どおりに扱う。

## 実装外（明示除外）

- 画面構造変更（layout、領域追加、新規カード、文言の新規追加）は実装外。
- 状態 pill の新規状態値（例: `"未設定"`）追加は実装外（fix-decision 禁止修正 3 を維持）。
- 空文字代理表現（空文字フィールドで「未設定」を表す）は実装外（fix-decision 禁止修正全件と新規方針が一致）。
- presenter の `?? "-"` を `?.trim() || "-"` に置き換える対症修正は実装外（fix-decision 禁止修正 1 を維持）。
- Svelte 判定を `=== ""` に変更する対症修正は実装外（fix-decision 禁止修正 2 を維持）。
- backend 戻り値型を `*string` 等に変える型変更は実装外（fix-decision 禁止修正 4 を維持）。
- persona-generation / body-translation 画面側の E2E 観点追加は本 task の範囲外（missing-tests.md の補足どおり）。

## 禁止修正の再掲

fix-decision.md の禁止修正のうち、1〜4 を本 task でも維持する。5 は Q-001 回答で覆る。新規禁止を追加する。

- 禁止 1（維持）: presenter の `?? "-"` を `?.trim() || "-"` 等に置き換える対症修正。
- 禁止 2（維持）: Svelte の判定を `viewModel.modelLabel === ""` に変更する対症修正。
- 禁止 3（維持）: `aiSettingsStatusLabel` に新しい状態値を追加する修正。
- 禁止 4（維持）: backend が `null` を返すように型を変更する修正（field 不在による表現に統一する）。
- 禁止 5（覆る）: persona-generation / body-translation の同一パターン同時修正の禁止は本 task で覆る。3 フェーズ同時実装が承認済み実装範囲である（Q-001）。
- 新規禁止 A: `credential_ref` を `JOB_PHASE_AI_SETTINGS` の保持列に含める修正（Q-003 で provider-settings 都度解決に分離）。
- 新規禁止 B: backend が派生説明文字列（blocked reason、設定不足の理由）を組み立てて DTO に含める修正（responses は「全体」または「不在」の二択）。
- 新規禁止 C: `SaveAISettings` 入力に `job_id` を含める修正、または `JOB_PHASE_AI_SETTINGS` の主キーに `job_id`/`user_id` を含める修正（Q-004 で `phase_type` 単独主キーに確定）。

## 検証単位と E2E 観点との対応

| 検証単位 | 対応 handoff | 対応 E2E 観点 |
| --- | --- | --- |
| `JOB_PHASE_AI_SETTINGS` の永続化（upsert / `ErrNotFound` / 3 phase_type 並存） | `H-MIG-ER`, `H-TEST-UNIT` | なし（lower-level） |
| 3 phase service の Ready 期 record 不在 → AI 設定 field 不在応答 | `H-BE-TERM`/`H-BE-PERSONA`/`H-BE-BODY`, `H-TEST-UNIT` | E2E-UC-FIX-MODEL-001 |
| 3 phase service のフェーズ開始時 `JOB_PHASE_RUN` 転写と開始拒否二択 | `H-BE-TERM`/`H-BE-PERSONA`/`H-BE-BODY`, `H-TEST-UNIT` | E2E-UC-FIX-MODEL-002 |
| 境界 DTO の field 不在表現と `SaveAISettings` 入力からの `job_id` 削除 | `H-INT-PHASE-AI-SETTINGS` | E2E-UC-FIX-MODEL-001/002/003（境界 DTO 越しに観測） |
| presenter の `isExecutionConfigured` 切り替えと派生語 frontend 集約 | `H-FE-PRESENTER`, `H-TEST-UNIT` | E2E-UC-FIX-MODEL-001/002 |
| `data-testid` 安定特定基盤 | `H-FE-SELECTOR` | E2E-UC-FIX-MODEL-001/002/003（selector 経由で観測） |
| 設定固定経路（保存後に「固定済み」へ遷移） | `H-FE-PRESENTER`, `H-INT-PHASE-AI-SETTINGS` | E2E-UC-FIX-MODEL-003 |

## 規模見積もり

| handoff | 想定変更ファイル数 | 想定変更行数 | 判定 |
| --- | --- | --- | --- |
| H-MIG-ER | 6-8（migration 2 系列各 2 ファイル + repository + interface + test） | 400-600 | 通常 |
| H-BE-TERM | 3-4 | 200-300 | 通常 |
| H-BE-PERSONA | 3-4 | 200-300 | 通常 |
| H-BE-BODY | 3-4 | 200-300 | 通常 |
| H-INT-PHASE-AI-SETTINGS | 5-8（DTO / gateway 型 / wails 生成系の追従） | 300-500 | 通常 |
| H-FE-PRESENTER | 6（3 presenter + 3 Panel の最小追従） | 300-500 | 通常 |
| H-FE-SELECTOR | 3（3 Panel） | 30-60 | 通常 |
| H-TEST-UNIT | 5-8 | 400-600 | 通常 |
| H-TEST-SCENARIO | 1-2 | 150-300 | 通常 |

各 handoff は通常範囲（15 files / 800 changed lines 以下）に収まる見込み。

## 残留リスク

- `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の縮退方針（廃止 / AI 設定列のみ DROP）は `H-MIG-ER` 初手の走査で確定するため、走査結果次第で migration 内容と repository コードの差分が増減する。`completion_evidence` に走査結果を残す。
- E2E 実行環境は Wails / 既存 E2E 基盤の整備状況で `FAIL_ENVIRONMENT` になる可能性があり、`H-TEST-SCENARIO` 完了判定は環境状況を `harness_gate_result` に明示する。
- 既存 detail-specs（`docs/detail-specs/term-translation-phase.md` 等）への docs 反映は本 implementation-scope の対象外。finalization-module の `updating-docs` で扱う。

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
- `sonar_gate_result`
- `harness_gate_result`
- `residual_risks`
- `completion_evidence`: 特に `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の縮退方針判断根拠を含める。
- `telemetry_events`
- `docs_changes: none`
