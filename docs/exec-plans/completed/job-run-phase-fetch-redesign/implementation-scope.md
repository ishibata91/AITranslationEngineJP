# Implementation Scope: job-run-phase-fetch-redesign

- `skill`: implementation-scope
- `status`: ready
- `source_plan`: `./plan.md`
- `human_review_status`: 設計3成果物は人間設計レビュー承認済み（2026-05-31）。本 implementation-scope は承認済み設計の分割であり、人間が呼び出し元レーンへ渡す入口にする。
- `approval_record`: 詳細仕様差分 `Q-001`〜`Q-005` 人間回答済み（2026-05-31）、未回答の未決 0 件。画面設計差分・設計差分図 承認済み（2026-05-31）。
- `codex_entry`: `.claude/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`: `./screen-design-diff.job-run.md`
- `component_diagram`: `./design-diff.job-run-phase-fetch-redesign.md`

## Fixed Decisions

- 人間設計レビュー済みの設計判断だけを根拠にする。`detail-spec-diff.md` の未回答の未決は 0 件。
- frontend 引き継ぎは承認済み `screen-design-diff.job-run.md` を根拠にする。
- backend、frontend、統合境界 を別 引き継ぎ に分ける。同一引き継ぎに束ねない。
- `E2E` は UI 人間操作起点だけを指す。`APIテスト` は public seam 起点の system-level test とする。
- secret は本 task で扱わない。全引き継ぎの `secret_boundary` は `not_required`。

### UI 順序ゲート判定（重要）

- 判定: 本 task は UI 変更を含む。初回取得中ローディングレイヤー（selector `<phase-prefix>-processing-target-loading`）の新規追加（`screen-design-diff` 差分 2）と、操作可否・次段階開始可否のフロント導出による表示根拠の変更（差分 6、差分 7）が画面に及ぶ。
- 帰結1: UI がある task のため frontend 引き継ぎを必須にし、backend 引き継ぎより前の wave に置く。
- 帰結2: 可否導出のフロント移設（設計変更 6、7）は、frontend が導出を持った後に統合境界 DTO から可否値を外し、最後に backend service の可否導出を撤去する順序にする。frontend が導出を持つ前に backend / DTO から可否を外すと画面が可否・理由を表示できなくなるため、この順序を不変条件にする。
- 帰結3: 表示フロー作り直し（設計変更 1〜5）は backend / bridge 仕様変更に依存しない（`screen-design-diff` の前提に明記）。UI 入口の frontend 引き継ぎとして最初の wave に置ける。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `FE-fetch-display`, `FE-readiness-derive-term`, `FE-readiness-derive-persona`, `FE-readiness-derive-body` | `なし` | `FE-readiness-derive-term <-> FE-readiness-derive-persona <-> FE-readiness-derive-body` | `FE-fetch-display は JobRunPage / 全段階 store を触るため readiness 系と並列にしない（owned_scope_overlap）` |
| `wave-2` | `INT-readiness-dto-term`, `INT-readiness-dto-persona`, `INT-body-summary-merge` | `wave-1 の対応段階 FE-readiness-derive-*` | `INT-readiness-dto-term <-> INT-readiness-dto-persona <-> INT-body-summary-merge` | `なし` |
| `wave-3` | `BE-fact-only-term`, `BE-fact-only-persona`, `BE-fact-only-body` | `wave-2 の対応段階 INT-*` | `BE-fact-only-term <-> BE-fact-only-persona <-> BE-fact-only-body` | `なし` |
| `wave-4` | `UT-fetch-display`, `UT-readiness-equivalence`, `SCN-target-list-e2e` | `wave-1, wave-2, wave-3 の全実装引き継ぎ` | `UT-fetch-display <-> UT-readiness-equivalence <-> SCN-target-list-e2e` | `なし` |

注: wave-1 の `FE-fetch-display` と `FE-readiness-derive-*` は同 wave だが、`FE-fetch-display` は JobRunPage と全段階 store / viewModel を触るため、段階別 readiness 引き継ぎと owned_scope が重なる。よって `FE-fetch-display` は単独着手とし、`FE-readiness-derive-*` 3 件だけを相互に並列可能とする。

## Handoffs

### `handoff_id`: FE-fetch-display

- `implementation_target`: 翻訳実行画面（`JobRunPage`）の処理対象一覧取得・表示フローを作り直し、初回表示で処理対象一覧が件数分表示される受け入れユースケース「処理対象を確認する」を成立させる。設計変更 1〜5 に対応する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`（`translation-job-management-REQ-006`、`term-translation-phase-REQ-007` の表示・取得規則）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 1〜5）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`（`$effect`/`onMount` の二重起動整理、表示中段階のみ取得起動、段階切り替え時取得、開き直し時の再取得、`currentProcessingTargetPageState` の derived 評価是正）
  - 各段階 frontend usecase の `fetchSummaryAndReadiness` 取得・反映フロー（summary 独立反映、processingTarget の連番ガード `processingTargetListRequestSequence` 適用、旧遅延応答破棄）
  - 各段階 store / viewModel の `initialFetchDone` 保持と購読変換
  - 初回取得中ローディングレイヤー UI（selector `<phase-prefix>-processing-target-loading` 新規追加。term / persona / body 同型）と操作排他（検索・ページ・行展開）
  - 注: 本引き継ぎは可否値（次段階開始可否・操作可否）の取得・保持・導出を変更しない。store から可否値を外す変更は別引き継ぎ（`FE-readiness-derive-*`）が扱う。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `owned_scope_overlap`（JobRunPage と全段階 store / viewModel を触るため `FE-readiness-derive-*` と owned_scope が重なる）
- `first_action`:
  - path: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - symbol: 取得起動を行う `$effect` と `onMount`
  - 変更種別: 二重起動整理（表示中段階のみ `setJobId`/取得起動へ一本化）
  - 対応する `completion_signal` clause: 「初回表示で表示中段階のみ取得し最大 2 本へ削減」
  - 1 手目にする理由: 二重起動の整理が取得本数削減・reopen 再取得・ローディング判定の土台であり、ここを固定しないと後続の反映ガードと derived 評価の前提が定まらない。
- `validation_commands`:
  - `npm --prefix frontend run lint`
  - `npm --prefix frontend run test -- src/application/usecase src/ui/screens/job-run`（本引き継ぎの owned_scope に対応する単体・コンポーネント検証。失敗時は本引き継ぎの実装不足として特定できる）
- `completion_signal`:
  - 初回表示で表示中段階のみ取得し、bridge 呼び出しを最大 2 本（summary + processingTarget）へ削減する。
  - summary は独立反映し、processingTarget は連番ガードで反映する。先行取得だけ完了で一覧が空のまま残る非対称を是正する。
  - 開き直し時に旧 sequence を無効化して再取得し、旧遅延応答を破棄する。
  - `initialFetchDone=false` の間、初回取得中ローディングレイヤーを最前面に出し、検索・ページ・行展開を排他する。selector `<phase-prefix>-processing-target-loading` を term / persona / body へ同型に追加する。
  - `currentProcessingTargetPageState` を `initialFetchDone=true` を評価条件に含む素直な derived 形へ是正する。
  - 承認済み画面設計差分の主要区画（処理対象一覧領域）、導線（段階切り替え）、状態表示（取得中・件数あり・空状態）を維持する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`（受け入れ証明は `SCN-target-list-e2e` の UI人間操作E2E が担う）
- `execution_stage`: `実装後`
- `notes`:
  - 本引き継ぎは UI 入口であり、backend / 統合境界に依存しない。最初の wave に置く。
  - 可否値の store 除去は `FE-readiness-derive-*` が担当する。本引き継ぎでは可否値の取得経路を変更しない。

### `handoff_id`: FE-readiness-derive-term

- `implementation_target`: 単語翻訳段階の次段階開始可否と操作可否を、backend 取得値ではなく段階データ事実状態から application 層で導出する受け入れユースケース。設計変更 6（term 分）に対応する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-007`、`term-translation-phase-REQ-008` の導出条件・成立条件・等価性条件）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 6）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/application/presenter/term-translation-phase/` と `frontend/src/application/usecase/term-translation-phase/` の application 層に、段階データ事実状態（フェーズ状態、ジョブ終端状態、対象件数、確認済み件数、エラー種別、実行設定の構成有無）から次段階開始可否（`canStartNextPhase` 相当）と操作可否（`canStart`/`canPause`/`canResume`/`canRetry`/`canCancel` と各 `*BlockedReason`）を導出するロジックを追加する。
  - 等価性条件の担保: 再配置前に backend が返していた可否・理由と、同じ事実入力に対して同じ結果を導く。
  - 画面（store / viewModel / screen-types 経由）が導出値を受け取る配線を整える。
  - 注: 本引き継ぎは backend が事実だけ返す前段階として、frontend に導出を先に持たせる。backend / DTO から可否値を外す変更は `INT-readiness-dto-term`、`BE-fact-only-term` が後続で行う。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `FE-readiness-derive-persona`, `FE-readiness-derive-body`
- `parallel_blockers`: `なし`（段階別ディレクトリで owned_scope が重ならない）
- `first_action`:
  - path: `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
  - symbol: 次段階開始可否導出の追加対象関数
  - 変更種別: 段階データ事実状態から `canStartNextPhase` 相当を導出する純粋関数を追加
  - 対応する `completion_signal` clause: 「次段階開始可否を事実状態から導出し、等価性条件を満たす」
  - 1 手目にする理由: 次段階開始可否は成立条件が `term-translation-phase-REQ-008` に明記され単独で閉じやすく、操作可否導出と等価性検証の基準点になる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/presenter/term-translation-phase src/application/usecase/term-translation-phase`
- `completion_signal`:
  - 次段階開始可否と各操作可否を application 層で事実状態から導出する。
  - 等価性条件: 再配置前後で同じ事実入力に同じ可否・理由が導ける。
  - 画面が application 層の導出値を表示根拠として受け取る（差分 6 規則 14〜19）。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - backend / DTO の可否値除去は後続 wave。本引き継ぎ単独では backend がまだ可否を返していても導出値で上書きする形にし、wave-2/3 完了後に backend 取得値依存を外す。

### `handoff_id`: FE-readiness-derive-persona

- `implementation_target`: NPC ペルソナ生成段階の本文翻訳段階開始可否と操作可否を、段階データ事実状態から application 層で導出する受け入れユースケース。設計変更 6（persona 分）に対応する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`（`persona-generation-phase-REQ-007`、`persona-generation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 6）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/application/presenter/persona-generation-phase/` と `frontend/src/application/usecase/persona-generation-phase/` に、段階データ事実状態（生成対象件数、生成完了件数、ペルソナ参照状態、フェーズ状態、ジョブ終端状態、実行設定の構成有無、エラー種別）から本文翻訳段階開始可否（`bodyReadiness` 相当）と操作可否を導出するロジックを追加する。
  - 等価性条件の担保（再配置前後で同結果）。
  - 画面が導出値を受け取る配線。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `FE-readiness-derive-term`, `FE-readiness-derive-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`
  - symbol: 本文翻訳段階開始可否導出の追加対象関数
  - 変更種別: 事実状態から `bodyReadiness` 相当を導出する純粋関数を追加
  - 対応する `completion_signal` clause: 「本文翻訳段階開始可否を事実状態から導出し、等価性条件を満たす」
  - 1 手目にする理由: term と同型で、段階固有の事実（生成完了件数・ペルソナ参照状態）を最初に固定すると操作可否導出が続けやすい。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/presenter/persona-generation-phase src/application/usecase/persona-generation-phase`
- `completion_signal`:
  - 本文翻訳段階開始可否と各操作可否を application 層で事実状態から導出する。
  - 等価性条件を満たす。
  - 画面が導出値を表示根拠として受け取る。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - backend / DTO の可否値除去は後続 wave。

### `handoff_id`: FE-readiness-derive-body

- `implementation_target`: 本文翻訳段階の成果物出力確認可否と操作可否を、段階要約取得の事実から application 層で導出する受け入れユースケース。設計変更 6（body 分）と設計変更 7 のフロント導出側に対応する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`（`body-translation-phase-REQ-006`、`body-translation-phase-REQ-007`）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 6、差分 7）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/application/presenter/body-translation-phase/`（`getEffectiveReadiness` 等）と `frontend/src/application/usecase/body-translation-phase/` に、段階要約取得の事実（完了項目件数 `completedFieldCount`、状態整合 `statusConsistent`、出力件数 `outputCount`、翻訳項目件数、失敗件数、フェーズ状態、ジョブ終端状態、実行設定の構成有無、エラー種別）から成果物出力確認可否（`ready` 相当・`canCheckOutputReadiness` 相当）と操作可否を導出するロジックを追加する。
  - 専用取得（`GetBodyTranslationOutputReadiness` / `BodyTranslationOutputReadinessResponse`）への frontend 依存を、段階要約取得の事実からの導出へ置き換える。専用取得を呼ぶ frontend gateway / presenter / screen-types 経路を段階要約取得 1 本へ統合する。
  - 等価性条件の担保（専用取得廃止前後で同結果）。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `FE-readiness-derive-term`, `FE-readiness-derive-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`
  - symbol: `getEffectiveReadiness`
  - 変更種別: 専用取得応答ではなく段階要約取得の事実（`completedFieldCount`/`statusConsistent`/`outputCount`）から出力確認可否を導出する形へ変更
  - 対応する `completion_signal` clause: 「成果物出力確認可否を段階要約取得の事実から導出する」
  - 1 手目にする理由: `getEffectiveReadiness` が専用取得依存の中心であり、ここを段階要約事実入力へ切り替えると取得経路統合の起点が定まる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/presenter/body-translation-phase src/application/usecase/body-translation-phase`
- `completion_signal`:
  - 成果物出力確認可否と各操作可否（`canCheckOutputReadiness` を含む）を application 層で段階要約取得の事実から導出する。
  - 専用取得への frontend 依存を段階要約取得 1 本へ統合する。
  - 等価性条件を満たす。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - frontend 側が段階要約取得 1 本へ統合した後、統合境界（`INT-body-summary-merge`）が DTO 側を整え、backend（`BE-fact-only-body`）が専用取得 endpoint を廃止する。

### `handoff_id`: INT-readiness-dto-term

- `implementation_target`: 単語翻訳段階の次段階開始可否・操作可否の可否値・理由を応答 DTO / gateway-contract から外し、段階データ事実状態だけを伝送する統合境界の接続。設計変更 6（term の DTO 境界分）に対応する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-008` の責務境界）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`（DTO / gateway 境界の伝送形のみ。UI 表示根拠は `FE-readiness-derive-term` が確定済み）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/controller/wails/term_translation_phase_controller.go` の `TermTranslationNextPhaseReadinessResponseDTO`・`ActionEnablement` 系 DTO から可否値（`canStartNextPhase`、`blockedReason`、各 `can*`、各 `*BlockedReason`）を外し、事実状態フィールドへ置き換える。
  - `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts` の `TermTranslationNextPhaseReadinessResponse` 等の型から可否値を外し、事実状態を受け取る型へ変更する。
  - bridge 伝送機構自体は破壊しない（呼び出し本数と応答内容の変更に限る）。
- `depends_on`: `FE-readiness-derive-term`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-readiness-dto-persona`, `INT-body-summary-merge`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`
  - symbol: `TermTranslationNextPhaseReadinessResponse`
  - 変更種別: 可否値フィールドを外し事実状態フィールドへ置換
  - 対応する `completion_signal` clause: 「DTO / gateway-contract が可否値を含まず事実状態を伝送する」
  - 1 手目にする理由: gateway-contract が frontend 導出（`FE-readiness-derive-term`）と backend DTO の接合点であり、ここを先に確定すると backend 側変更の対象が定まる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/gateway-contract/term-translation-phase`
  - `go test ./internal/controller/wails/ -run TermTranslation`
- `completion_signal`:
  - 応答 DTO と gateway-contract が可否値（`canStartNextPhase`、`blockedReason`、各 `can*`、各 `*BlockedReason`）を含まず、段階データ事実状態を伝送する。
  - frontend は伝送された事実状態を `FE-readiness-derive-term` の導出入力として受け取れる。
  - 実画面確認: 統合後に単語翻訳段階を開き、次段階開始可否・操作可否の活性と理由が表示される（実装後ブラウザ確認で確認する経路を `browser_confirmation` へ残す）。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 統合境界 引き継ぎ。backend service の可否導出撤去は `BE-fact-only-term` が行う。本引き継ぎは DTO / gateway 契約の接続に限る。
  - 本番経路: `TermTranslationPhaseController.GetTermTranslationNextPhaseReadiness` → 応答 DTO → bridge → gateway-contract → application 層導出。

### `handoff_id`: INT-readiness-dto-persona

- `implementation_target`: NPC ペルソナ生成段階の本文翻訳段階開始可否・操作可否の可否値・理由を応答 DTO / gateway-contract から外し、段階データ事実状態だけを伝送する統合境界の接続。設計変更 6（persona の DTO 境界分）に対応する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`（`persona-generation-phase-REQ-008` の責務境界）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/controller/wails/persona_generation_phase_controller.go` の `PersonaGenerationBodyReadinessResponseDTO`・`ActionEnablement` 系 DTO から可否値（`bodyReadiness`、`blockedReason`、各 `can*`、各 `*BlockedReason`）を外し、`inputSummary` 相当の事実状態へ置き換える。
  - persona の frontend gateway-contract の対応型から可否値を外し事実状態を受け取る型へ変更する。
  - bridge 伝送機構自体は破壊しない。
- `depends_on`: `FE-readiness-derive-persona`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-readiness-dto-term`, `INT-body-summary-merge`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-phase-gateway-contract.ts`
  - symbol: 本文翻訳段階開始可否応答型
  - 変更種別: 可否値フィールドを外し事実状態フィールドへ置換
  - 対応する `completion_signal` clause: 「DTO / gateway-contract が可否値を含まず事実状態を伝送する」
  - 1 手目にする理由: term と同型で gateway-contract が接合点になる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/gateway-contract/persona-generation-phase`
  - `go test ./internal/controller/wails/ -run PersonaGeneration`
- `completion_signal`:
  - 応答 DTO と gateway-contract が可否値（`bodyReadiness`、`blockedReason`、各 `can*`、各 `*BlockedReason`）を含まず事実状態を伝送する。
  - frontend が事実状態を `FE-readiness-derive-persona` の導出入力として受け取れる。
  - 実画面確認: 統合後に NPC ペルソナ生成段階を開き本文翻訳段階開始可否・操作可否が表示される経路を `browser_confirmation` へ残す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - backend service の可否導出撤去は `BE-fact-only-persona` が行う。
  - 本番経路: `PersonaGenerationPhaseController.GetPersonaGenerationBodyReadiness` → 応答 DTO → bridge → gateway-contract → application 層導出。

### `handoff_id`: INT-body-summary-merge

- `implementation_target`: 本文翻訳段階の成果物出力確認専用取得を廃止し、必要な事実（`completedFieldCount`、`statusConsistent`、`outputCount`）を段階要約取得応答へ集約する統合境界の接続。可否値を DTO から外す。設計変更 7 の DTO 境界分に対応する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`（`body-translation-phase-REQ-007` の取得経路統合・責務境界）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/controller/wails/body_translation_phase_controller.go` の `BodyTranslationPhaseSummaryResponseDTO` に `completedFieldCount`、`statusConsistent`、`outputCount` を事実状態として集約する（既存 `resultSummary` 内項目の扱いを段階要約の事実として frontend が参照できる形に確定する）。
  - `BodyTranslationOutputReadinessResponseDTO`（可否値 `canCheckOutputReadiness`、`outputReadiness.ready`、`blockedReason` を含む）と専用取得 endpoint `GetBodyTranslationOutputReadiness` の DTO / gateway 経路を廃止する。
  - body の frontend gateway-contract（`body-translation-phase-gateway-contract.ts` の `BodyTranslationOutputReadinessResponse`）を廃止し、段階要約取得型へ統合する。
  - bridge 伝送機構自体は破壊しない。
- `depends_on`: `FE-readiness-derive-body`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `INT-readiness-dto-term`, `INT-readiness-dto-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/gateway-contract/body-translation-phase/body-translation-phase-gateway-contract.ts`
  - symbol: `BodyTranslationOutputReadinessResponse`
  - 変更種別: 専用取得応答型の廃止と段階要約取得型への事実集約
  - 対応する `completion_signal` clause: 「専用取得を廃止し段階要約取得へ事実を集約する」
  - 1 手目にする理由: `FE-readiness-derive-body` が段階要約事実への依存へ切り替え済みのため、gateway-contract の専用取得型廃止が次に閉じる境界になる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/gateway-contract/body-translation-phase`
  - `go test ./internal/controller/wails/ -run BodyTranslation`
- `completion_signal`:
  - 段階要約取得応答 DTO に `completedFieldCount`、`statusConsistent`、`outputCount` が事実として含まれ、frontend が `FE-readiness-derive-body` の導出入力として参照できる。
  - 専用取得 endpoint と `BodyTranslationOutputReadinessResponse(DTO)` が廃止され、可否値（`canCheckOutputReadiness`、`outputReadiness.ready`、`blockedReason`）を DTO が含まない。
  - 実画面確認: 統合後に本文翻訳段階を開き成果物出力確認可否・理由が表示される経路を `browser_confirmation` へ残す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - backend service / usecase 側の専用取得実装（`ReadOutputReadiness`、`buildOutputReadiness`、`GetBodyTranslationOutputReadiness` usecase）の撤去と段階要約への事実集約は `BE-fact-only-body` が行う。本引き継ぎは DTO / gateway 契約の統合に限る。
  - 注意: 出力成果物側の最終出力可否（`translation-output-artifact` の `outputReadiness`）は対象外（`body-translation-phase-REQ-007` の注意）。本引き継ぎは本文翻訳段階画面の成果物出力確認可否に限る。
  - 本番経路: `BodyTranslationPhaseController.GetBodyTranslationPhaseSummary` → 段階要約応答 DTO → bridge → gateway-contract → application 層導出。

### `handoff_id`: BE-fact-only-term

- `implementation_target`: 単語翻訳段階の backend service / usecase / controller から次段階開始可否・操作可否の導出を撤去し、段階データ事実状態だけ返す責務へ確定する。設計変更 6（term の backend 分）に対応する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-008` の責務境界・等価性条件）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/service/term_translation_phase_service.go` の `readinessFromState`、`termTranslationStartBlockedReason`、`termTranslationPauseBlockedReason`、`termTranslationResumeBlockedReason`、`termTranslationRetryBlockedReason` などの可否導出を撤去し、段階データ事実状態（フェーズ状態、ジョブ終端、対象件数、確認済み件数、エラー種別、実行設定の構成有無）を read model として返す形へ確定する。
  - `internal/usecase/term_translation_phase_usecase.go` の `GetTermTranslationNextPhaseReadiness` を、可否ではなく事実状態を返す責務へ整える。
  - `internal/controller/wails/term_translation_phase_controller.go` を `INT-readiness-dto-term` で確定した事実状態 DTO に合わせる。
  - backend テスト（`term_translation_phase_service_test.go` の readiness 系）を事実状態返却の検証へ更新する。
- `depends_on`: `INT-readiness-dto-term`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `BE-fact-only-persona`, `BE-fact-only-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/term_translation_phase_service.go`
  - symbol: `readinessFromState`
  - 変更種別: 可否導出（`CanStartNextPhase` 算出）を撤去し事実状態返却へ変更
  - 対応する `completion_signal` clause: 「backend は可否を導出せず事実状態だけ返す」
  - 1 手目にする理由: `readinessFromState` が次段階開始可否導出の中心であり、ここを事実返却へ変えると操作可否導出関数群の撤去対象が定まる。
- `validation_commands`:
  - `go test ./internal/service/ -run TermTranslation`
  - `go test ./internal/usecase/ -run TermTranslation`
  - `go test ./internal/controller/wails/ -run TermTranslation`
- `completion_signal`:
  - backend service / usecase / controller が次段階開始可否・操作可否を導出・返却せず、段階データ事実状態だけ返す。
  - `INT-readiness-dto-term` が確定した事実状態 DTO 形と一致する。
  - 等価性条件: 撤去前に backend が返していた可否・理由を、frontend が同じ事実入力から再現できる（`FE-readiness-derive-term` と `UT-readiness-equivalence` で証明する）。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 本番経路: repository → service（事実状態） → usecase → controller DTO（事実状態） → bridge。

### `handoff_id`: BE-fact-only-persona

- `implementation_target`: NPC ペルソナ生成段階の backend から本文翻訳段階開始可否・操作可否の導出を撤去し、事実状態だけ返す責務へ確定する。設計変更 6（persona の backend 分）に対応する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`persona-generation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/service/persona_generation_phase_service.go` の `ReadBodyReadiness`、`personaBlockedReason` 等の可否導出を撤去し、生成対象件数・生成完了件数・ペルソナ参照状態・フェーズ状態・ジョブ終端・実行設定の構成有無・エラー種別の事実状態を read model として返す。
  - `internal/usecase/persona_generation_phase_usecase.go` の `GetPersonaGenerationBodyReadiness` を事実状態返却へ整える。
  - `internal/controller/wails/persona_generation_phase_controller.go` を `INT-readiness-dto-persona` の事実状態 DTO に合わせる。
  - backend テスト（`persona_generation_phase_service_test.go` の readiness 系）を事実状態返却の検証へ更新する。
- `depends_on`: `INT-readiness-dto-persona`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `BE-fact-only-term`, `BE-fact-only-body`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/persona_generation_phase_service.go`
  - symbol: `ReadBodyReadiness`
  - 変更種別: 可否導出を撤去し事実状態返却へ変更
  - 対応する `completion_signal` clause: 「backend は可否を導出せず事実状態だけ返す」
  - 1 手目にする理由: `ReadBodyReadiness` が本文翻訳段階開始可否導出の中心。
- `validation_commands`:
  - `go test ./internal/service/ -run PersonaGeneration`
  - `go test ./internal/usecase/ -run PersonaGeneration`
  - `go test ./internal/controller/wails/ -run PersonaGeneration`
- `completion_signal`:
  - backend が本文翻訳段階開始可否・操作可否を導出・返却せず事実状態だけ返す。
  - `INT-readiness-dto-persona` の事実状態 DTO 形と一致する。
  - 等価性条件を frontend が再現できる（`FE-readiness-derive-persona` と `UT-readiness-equivalence` で証明する）。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 本番経路: repository → service（事実状態） → usecase → controller DTO（事実状態） → bridge。

### `handoff_id`: BE-fact-only-body

- `implementation_target`: 本文翻訳段階の backend から成果物出力確認可否・操作可否の導出を撤去し、専用取得 endpoint を廃止して段階要約取得へ事実を集約する責務へ確定する。設計変更 7 の backend 分に対応する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`body-translation-phase-REQ-007` の責務境界・取得経路統合・等価性条件）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/service/body_translation_phase_service.go` の `ReadOutputReadiness`、`buildOutputReadiness` を撤去し、段階要約 read model に `completedFieldCount`、`statusConsistent`、`outputCount` を含む事実状態を集約する。
  - `internal/usecase/body_translation_phase_usecase.go` の `GetBodyTranslationOutputReadiness`、`toBodyTranslationOutputReadinessSummary` を撤去し、段階要約取得 usecase へ事実を集約する。
  - `internal/controller/wails/body_translation_phase_controller.go` を `INT-body-summary-merge` の段階要約事実集約 DTO に合わせ、専用取得 endpoint を撤去する。
  - backend テスト（`body_translation_phase_service_test.go`、`body_translation_phase_usecase_test.go`、`body_translation_phase_scenario_test.go` の output readiness 系、`body_translation_phase_contract.go` のスタブ）を段階要約集約・専用取得廃止に合わせて更新する。
- `depends_on`: `INT-body-summary-merge`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `BE-fact-only-term`, `BE-fact-only-persona`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/body_translation_phase_service.go`
  - symbol: `buildOutputReadiness`
  - 変更種別: 出力確認可否導出を撤去し段階要約 read model へ事実集約
  - 対応する `completion_signal` clause: 「専用取得を廃止し段階要約へ事実を集約する」
  - 1 手目にする理由: `buildOutputReadiness` が出力確認可否導出の中心であり、ここを事実集約へ変えると専用取得 usecase / endpoint の撤去対象が定まる。
- `validation_commands`:
  - `go test ./internal/service/ -run BodyTranslation`
  - `go test ./internal/usecase/ -run BodyTranslation`
  - `go test ./internal/controller/wails/ -run BodyTranslation`
- `completion_signal`:
  - backend が成果物出力確認可否・操作可否を導出・返却せず、段階要約 read model に `completedFieldCount`、`statusConsistent`、`outputCount` を含む事実だけ返す。
  - 専用取得 endpoint `GetBodyTranslationOutputReadiness` と関連 read model / DTO が廃止される。
  - `INT-body-summary-merge` の段階要約事実集約 DTO 形と一致する。
  - 等価性条件を frontend が再現できる（`FE-readiness-derive-body` と `UT-readiness-equivalence` で証明する）。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 本番経路: repository → service（段階要約事実） → usecase → controller 段階要約 DTO → bridge。

### `handoff_id`: UT-fetch-display

- `implementation_target`: 取得・表示フロー作り直し（`FE-fetch-display`）の責務を単体テストで証明する。取得本数削減、summary 独立反映・processingTarget 連番ガードの対称化、開き直し時の旧応答破棄、`initialFetchDone` による derived 評価と操作排他。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`translation-job-management-REQ-006`、`term-translation-phase-REQ-007`）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 1〜5）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 各段階 frontend usecase の取得・反映フロー単体テスト（summary 独立反映、processingTarget 連番ガード、旧 sequence 破棄、`initialFetchDone` 更新）
  - `JobRunPage` の `currentProcessingTargetPageState` derived 評価と初回取得中ローディングレイヤーの操作排他のコンポーネント／単体テスト
- `depends_on`: `FE-fetch-display`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `UT-readiness-equivalence`, `SCN-target-list-e2e`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.test.ts`
  - symbol: 連番ガードと summary 独立反映の検証ケース
  - 変更種別: 先行取得だけ完了で一覧が空に残らないことを証明するテスト追加
  - 対応する `completion_signal` clause: 「summary 独立反映と processingTarget 連番ガードの対称化を証明する」
  - 1 手目にする理由: 本不具合の主軸（反映取りこぼし防止）であり、最初に固定すべき検証意図。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/usecase src/ui/screens/job-run`
- `completion_signal`:
  - 取得本数削減、summary 独立反映・processingTarget 連番ガードの対称化、旧応答破棄、`initialFetchDone` による derived 評価と操作排他を単体・コンポーネントテストで証明する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `notes`:
  - 実装成果物（`FE-fetch-display`）完了後に着手する。

### `handoff_id`: UT-readiness-equivalence

- `implementation_target`: 可否導出のフロント移設（`FE-readiness-derive-*`）と backend 事実化（`BE-fact-only-*`）について、等価性条件（再配置前後で同じ事実入力に同じ可否・理由）を単体テストで証明する。term / persona / body 同型。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-008`、`persona-generation-phase-REQ-008`、`body-translation-phase-REQ-007` の等価性条件）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 各段階 presenter / usecase の可否導出単体テスト（次段階開始可否・本文翻訳段階開始可否・成果物出力確認可否と各操作可否、各 `*BlockedReason`）。
  - 等価性検証: 移設前に backend が返していた可否・理由を期待値とし、同じ事実入力でフロント導出が同結果になることを証明する。期待値の元ネタは `spec_basis` の成立条件・区別理由を使う。
- `depends_on`: `FE-readiness-derive-term`, `FE-readiness-derive-persona`, `FE-readiness-derive-body`, `BE-fact-only-term`, `BE-fact-only-persona`, `BE-fact-only-body`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `UT-fetch-display`, `SCN-target-list-e2e`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.test.ts`
  - symbol: 次段階開始可否導出の等価性ケース
  - 変更種別: 成立条件（終端でない・完了・確認済み件数≥対象件数）と区別理由の等価性検証を追加
  - 対応する `completion_signal` clause: 「term の次段階開始可否導出が等価性条件を満たすことを証明する」
  - 1 手目にする理由: term の成立条件が `term-translation-phase-REQ-008` に最も具体的に書かれ、等価性検証の基準になる。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/application/presenter src/application/usecase`
- `completion_signal`:
  - term / persona / body の次段階開始可否相当と各操作可否を、`spec_basis` の成立条件・区別理由を期待値として導出が満たすことを証明する。
  - 等価性条件（再配置前後で同じ事実入力に同じ結果）を 3 段階で証明する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `notes`:
  - 単体テスト引き継ぎ。期待結果の元ネタは `spec_basis`（承認済み詳細仕様差分）。
  - frontend 導出と backend 事実化の両方が揃った後に着手する。

### `handoff_id`: SCN-target-list-e2e

- `implementation_target`: 残置 E2E（`tests/system/fix-lucien-target-list-empty.spec.ts` の E2E-LTLE-001/002/003）を green にし、受け入れユースケース「処理対象を確認する」の初回表示挙動を UI 人間操作起点で証明する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`（`translation-job-management-REQ-006`、`term-translation-phase-REQ-007`）
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`（差分 5 の固定 selector）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `tests/system/fix-lucien-target-list-empty.spec.ts`（E2E-LTLE-001/002/003）を作り直し後の挙動に合わせて green 化する。
  - `tests/system/support/scenario-wails-mocks.ts` の `termZeroAITargetJobId` オプションと page object（`tests/system/support/translation-phase-pages.ts` / `system-test-pages.ts`）が差分 5 の固定 selector（`-total`、`-empty`、`-search-input`、`-row`、新規 `-loading`）と一致することを確認する。
- `depends_on`: `FE-fetch-display`, `FE-readiness-derive-term`, `FE-readiness-derive-persona`, `FE-readiness-derive-body`, `INT-readiness-dto-term`, `INT-readiness-dto-persona`, `INT-body-summary-merge`, `BE-fact-only-term`, `BE-fact-only-persona`, `BE-fact-only-body`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `UT-fetch-display`, `UT-readiness-equivalence`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `tests/system/fix-lucien-target-list-empty.spec.ts`
  - symbol: `E2E-LTLE-001`
  - 変更種別: 初回表示で処理対象行が件数分表示される正常系の green 化確認
  - 対応する `completion_signal` clause: 「E2E-LTLE-001/002/003 が green になる」
  - 1 手目にする理由: E2E-LTLE-001 が本不具合の正常系（初回0件にならない）の中心で、ここを通すと境界系の前提が定まる。
- `validation_commands`:
  - `npm run test:system -- fix-lucien-target-list-empty.spec.ts`（system test 実行は elevated 権限が必要。harness で止まる場合は `harness_gate_result` を `FAIL_ENVIRONMENT` とし再実行環境・コマンドを残す）
- `completion_signal`:
  - E2E-LTLE-001（初回表示で処理対象行が件数分・空状態にならない）、E2E-LTLE-002（検索後リロードで初回0件へ戻らない）、E2E-LTLE-003（母数0で空状態保持）が green になる。
  - 差分 5 の固定 selector と page object 参照が一致する。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `notes`:
  - 限界: モック E2E は同期解決のため実機 IPC 飽和（約15秒遅延）を再現しない。実機固有症状の確認は `browser_confirmation` の実装後ブラウザ確認で扱う。
  - 全実装引き継ぎ（wave-1〜3）完了後に着手する。

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
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`: 実機 15 秒遅延（IPC 飽和の疑い、真因未確定）の分離は観測ログ追加の要否を呼び出し元レーンが再判定する（`detail-spec-diff` 注記1）。
- `completion_evidence`
- `telemetry_events`: `runtime: codex`
- `docs_changes: none`
