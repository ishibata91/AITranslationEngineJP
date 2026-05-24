# Implementation Scope: frontend-backend-connection-refactor

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: `./refactor-scope-confirmation.md` の人間入力 `全部承認`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `task_plan`: `./plan.md`
- `spec_implementation_drift`: `./spec-implementation-drift.md`
- `structure_quality_investigation`: `./structure-quality-investigation.md`
- `test_quality_investigation`: `./test-quality-investigation.md`
- `refactor_scope_confirmation`: `./refactor-scope-confirmation.md`
- `detail_spec_diff`: `N/A`
  - 理由: この task は `refactor_lane` 起点であり、要件、詳細仕様、画面仕様に基づく人間判断待ちは空集合として完了済みである。
- `screen_design_diff`: `N/A`
  - 理由: この task は画面表示内容を変更しない。既存画面正本の接続状態表示を維持する。
- `detail_spec_sources`:
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/translation-job-management.md`
- `screen_design_sources`:
  - `docs/screen-design/screens/provider-settings.md`
  - `docs/screen-design/screens/term-translation-phase.md`

## Fixed Decisions

- `unanswered_questions`: `0`
- 人間レビューで `SQ-FBC-001`、`SQ-FBC-002`、`SQ-FBC-003`、`TQI-FBC-001`、`TQI-FBC-002`、`TQI-FBC-003` は承認済みである。
- `frontend/wailsjs/` は generated bindings である。実装系レーンは読む対象にできるが、手編集しない。
- frontend gateway は `GatewayContract` を実装し、generated `wailsjs` と backend DTO を `frontend/src/controller/wails/` に閉じ込める。
- backend bind 公開面は `internal/controller/wails/` とし、frontend の query / command は generated `wailsjs` を経由する。
- docs 正本化、`.codex` 変更、remote repository 変更は実装引き継ぎに含めない。

## Approved Implementation Scope

| ID | 承認済み実装範囲 | 主な成果物 |
| --- | --- | --- |
| `SQ-FBC-001` | frontend gateway の Wails binding 解決経路を `generated wailsjs` と正規 binding 面に寄せる。 | 統合境界実装 |
| `SQ-FBC-002` | Wails bridge 戻り値を gateway 内で runtime shape 検証し、無検証の DTO 型変換を縮小する。 | 統合境界実装 |
| `SQ-FBC-003` | screen controller factory から gateway DTO 型依存を外し、依存方向を application contract 側へ寄せる。 | frontend 実装 |
| `TQI-FBC-001` | frontend gateway test の観測点を `globalThis.go` 探索順から public seam へ寄せる。 | 単体テスト |
| `TQI-FBC-002` | backend controller test を public method 単位へ拡張し、DTO 写像と error wrap の未観測面を埋める。 | 単体テスト |
| `TQI-FBC-003` | fake API ではない接続境界専用の scenario test または integration test を追加する。 | シナリオテスト |

## Split Result

| 分割 | handoff | 扱い |
| --- | --- | --- |
| backend 実装 | `N/A` | 承認済み候補に backend product code の単独変更はない。`internal/controller/` と `internal/bootstrap/` は統合境界と backend 単体テストの対象候補として扱う。 |
| frontend 実装 | `FBC-FE-001` | screen controller factory の DTO 依存を application contract 側へ寄せる。 |
| 統合境界実装 | `FBC-INT-001` | generated `wailsjs` 呼び出し、runtime shape 検証、Wails binding error を gateway 境界で揃える。 |
| 単体テスト | `FBC-UT-FE-001` | frontend gateway test の観測点を public seam へ寄せる。 |
| 単体テスト | `FBC-UT-BE-001` | backend controller public method test を追加する。 |
| シナリオテスト | `FBC-SC-001` | fake API ではない接続境界 test を追加する。 |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `FBC-FE-001` | なし | なし | なし |
| `wave-2` | `FBC-INT-001` | `FBC-FE-001` | なし | `depends_on`, `shared_contract_change` |
| `wave-3` | `FBC-UT-FE-001`, `FBC-UT-BE-001` | `FBC-INT-001` | `FBC-UT-FE-001 <-> FBC-UT-BE-001` | なし |
| `wave-4` | `FBC-SC-001` | `FBC-INT-001`, `FBC-UT-FE-001`, `FBC-UT-BE-001` | なし | `depends_on`, `broad_gate_shared` |

## Handoffs

### `FBC-FE-001`: screen controller factory DTO 依存分離

- `implementation_target`: frontend
- `implementation_artifact`: frontend 実装
- `implementation_skill`: implement-frontend
- `approved_scope`: `SQ-FBC-003`
- `spec_basis`:
  - `./structure-quality-investigation.md`
  - `docs/architecture.md`
  - `docs/coding-guidelines-frontend.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design_sources`: `docs/screen-design/screens/term-translation-phase.md`
- `secret_boundary`:
  - `status`: not_required
- `owned_scope`:
  - 対象ファイル候補:
    - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts`
    - `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`
    - `frontend/src/application/contract/term-translation-phase/`
  - 変更禁止範囲:
    - `frontend/wailsjs/`
    - backend product code
    - docs 正本本文
    - `.codex/`
  - 想定規模: `2-4 files`, `120-220 changed lines`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`: `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts` の `createTermTranslationPhaseScreenControllerFactory` から `@controller/wails/gateway-dto/term-translation-phase` import を外す。変更種別は依存方向修正。対応する完了条件は「screen controller factory が gateway DTO 型へ依存しない」である。理由は `SQ-FBC-003` の責務分離を最小単位で閉じられるためである。
- `completion_signal`:
  - screen controller factory は Wails 呼び出し詳細と gateway DTO 型を参照しない。
  - coverage 用の型が必要な場合は application contract 側、または gateway 境界側の test helper に閉じる。
  - 単語翻訳画面の `Gateway: <接続状態>`、段階状態、操作可否の表示契約を維持する。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `FBC-INT-001`: Wails gateway binding と DTO runtime shape 検証

- `implementation_target`: frontend-backend integration boundary
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: implement-integration
- `approved_scope`: `SQ-FBC-001`, `SQ-FBC-002`
- `spec_basis`:
  - `./structure-quality-investigation.md`
  - `docs/architecture.md`
  - `docs/coding-guidelines-frontend.md`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/term-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design_sources`: `docs/screen-design/screens/provider-settings.md`, `docs/screen-design/screens/term-translation-phase.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `providerId`, `credentialReferenceId`, `credentialState`, `validationState`, `requestToken`, AIサービス名、モデル名、実行方式、一括処理設定
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、`credentialInput` の実値、復号可能な秘密値、認証参照の実値
  - `secret_resolution_owner_layer`: backend usecase / service 側の AIサービス設定解決責務
  - `forbidden_outputs`: API key 平文、復号可能な秘密値、認証参照の実値、外部サービス生データ、翻訳本文全文、実行時に解決した接続先
- `owned_scope`:
  - 対象ファイル候補:
    - `frontend/src/controller/wails/*.gateway.ts`
    - `frontend/src/controller/wails/gateway-dto/**`
    - `frontend/src/application/gateway-contract/**`
    - `frontend/src/bootstrap/app-screen-controller-factories.ts`
    - `frontend/src/main.ts`
    - `frontend/wailsjs/go/wails/AppController.js` は読む対象だけ
    - `internal/controller/wails/app_controller.go` は読む対象を基本とし、bind 面の不整合がある場合だけ変更候補にする
    - `internal/bootstrap/app_controller.go` は読む対象を基本とし、bind wiring の不整合がある場合だけ変更候補にする
  - 変更禁止範囲:
    - `frontend/wailsjs/` の手編集
    - UI 表示追加、画面導線変更、画面文言変更
    - docs 正本本文
    - `.codex/`
  - 想定規模: `6-10 files`, `350-650 changed lines`
- `depends_on`: `FBC-FE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`, `shared_contract_change`
- `first_action`: `frontend/src/controller/wails/provider-settings.gateway.ts` の binding 解決処理を generated `wailsjs/go/wails/AppController.js` の public function import へ寄せる。変更種別は transport adapter 置換。対応する完了条件は「gateway が `globalThis.go.wails.*` の controller 探索順へ依存しない」である。理由は provider settings gateway が接続、secret、runtime shape 検証を同時に代表できるためである。
- `completion_signal`:
  - gateway は generated `wailsjs` の public function を正規 binding 面として使う。
  - `globalThis.go.wails.*` の controller 名探索順は gateway の主経路から外す。
  - Wails bridge 戻り値は `unknown` 相当として受け、gateway 内で DTO shape を絞り込んでから application contract へ返す。
  - runtime shape 検証失敗時は user-facing message と internal diagnostic を分ける。
  - `frontend/wailsjs/` は手編集しない。
  - provider settings と term translation の既存画面は `Gateway` 状態、接続情報、秘匿情報の表示契約を維持する。
  - 接続境界変更を行ったため、実装後ブラウザ確認の対象にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後
- `notes`:
  - 本番経路: `frontend/src/main.ts` -> `frontend/src/bootstrap/app-screen-controller-factories.ts` -> `frontend/src/controller/wails/*.gateway.ts` -> `frontend/wailsjs/go/wails/AppController.js` -> `internal/controller/wails/AppController` -> controller public method。

### `FBC-UT-FE-001`: frontend gateway public seam test

- `implementation_target`: frontend unit test
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `approved_scope`: `TQI-FBC-001`
- `spec_basis`:
  - `./test-quality-investigation.md`
  - `./structure-quality-investigation.md`
  - `docs/coding-guidelines-tests.md`
  - `docs/coding-guidelines-frontend.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `credentialReferenceId`, `credentialState`, `validationState`, `requestToken`
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、`credentialInput` の実値
  - `secret_resolution_owner_layer`: backend usecase / service 側
  - `forbidden_outputs`: API key 平文、復号可能な秘密値、認証参照の実値
- `owned_scope`:
  - 対象ファイル候補:
    - `frontend/src/controller/wails/provider-settings.gateway.test.ts`
    - `frontend/src/controller/wails/term-translation-phase.gateway.test.ts`
    - `frontend/src/controller/wails/translation-job-management.gateway.test.ts`
    - `frontend/src/controller/wails/body-translation-phase.gateway.test.ts`
    - gateway transport adapter test を追加する場合は `frontend/src/controller/wails/`
  - 変更禁止範囲:
    - `frontend/wailsjs/`
    - backend product code
    - docs 正本本文
    - `.codex/`
  - 想定規模: `4-7 files`, `350-650 changed lines`
- `depends_on`: `FBC-INT-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `FBC-UT-BE-001`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/controller/wails/provider-settings.gateway.test.ts` の `ProviderSettingsController が未接続なら AppController の binding を使う` test を public seam 観測へ置き換える。変更種別は test 観測点修正。対応する完了条件は「gateway test が controller 名探索順を主観測点にしない」である。理由は `TQI-FBC-001` の対象がこの test の fallback 観測に現れているためである。
- `completion_signal`:
  - gateway test は request、response、未接続、runtime shape 検証失敗を public seam から観測する。
  - generated binding wrapper または transport adapter を差し替える test seam は gateway 境界内に閉じる。
  - `globalThis.go` の controller 名探索順を gateway suite 全体の前提にしない。
  - secret 値の平文が response、error message、diagnostic の公開値に出ないことを確認する。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `FBC-UT-BE-001`: backend controller public method test

- `implementation_target`: backend unit test
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `approved_scope`: `TQI-FBC-002`
- `spec_basis`:
  - `./test-quality-investigation.md`
  - `docs/coding-guidelines-tests.md`
  - `docs/coding-guidelines-backend.md`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/translation-job-management.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `credentialReferenceId`, `credentialState`, `validationState`, `requestToken`
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、`credentialInput` の実値、認証参照の実値
  - `secret_resolution_owner_layer`: backend usecase / service 側
  - `forbidden_outputs`: API key 平文、復号可能な秘密値、認証参照の実値、外部サービス生データ
- `owned_scope`:
  - 対象ファイル候補:
    - `internal/controller/wails/provider_settings_controller_unit_test.go`
    - `internal/controller/wails/translation_job_management_controller_unit_test.go`
    - `internal/controller/wails/term_translation_phase_controller_unit_test.go`
    - 必要な場合だけ `internal/controller/wails/*_controller.go`
  - 変更禁止範囲:
    - frontend product code
    - `frontend/wailsjs/`
    - docs 正本本文
    - `.codex/`
  - 想定規模: `3-5 files`, `450-750 changed lines`
- `depends_on`: `FBC-INT-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `FBC-UT-FE-001`
- `parallel_blockers`: なし
- `first_action`: `internal/controller/wails/provider_settings_controller_unit_test.go` に `ListProviderSettings` の response DTO 写像を確認する test を追加する。変更種別は controller public method test 追加。対応する完了条件は「ProviderSettingsController の未観測 public method を 1 method ずつ観測する」である。理由は既存 test が `SaveProviderSettings` の trim に偏っているためである。
- `completion_signal`:
  - `ProviderSettingsController` は `ListProviderSettings`、`ResetProviderSettings`、`ValidateProviderSettings` の request / response DTO 写像と error wrap を観測する。
  - `TranslationJobManagementController` は `GetJobDetail`、`RequestStop`、`ResumeJob` の公開応答形と error wrap を観測する。
  - `TermTranslationPhaseController` は pause、resume、retry、save AI settings の DTO 境界を観測する。
  - 成功経路と失敗経路は混ぜず、失敗時に壊れた public method が分かる assertion にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `FBC-SC-001`: fake API ではない接続境界 test

- `implementation_target`: scenario or integration test
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `approved_scope`: `TQI-FBC-003`
- `spec_basis`:
  - `./test-quality-investigation.md`
  - `docs/coding-guidelines-tests.md`
  - `docs/architecture.md`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/term-translation-phase.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
  - `screen_design_sources`: `docs/screen-design/screens/provider-settings.md`, `docs/screen-design/screens/term-translation-phase.md`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: AIサービス名、モデル名、認証状態、実行方式、一括処理設定、`credentialReferenceId`
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、認証参照の実値、復号可能な秘密値
  - `secret_resolution_owner_layer`: backend usecase / service 側
  - `forbidden_outputs`: API key 平文、復号可能な秘密値、認証参照の実値、外部サービス生データ、翻訳本文全文、実行時に解決した接続先
- `owned_scope`:
  - 対象ファイル候補:
    - `tests/system/`
    - `internal/bootstrap/app_controller_test.go`
    - `frontend/src/main.ts`
    - `frontend/src/bootstrap/app-screen-controller-factories.ts`
    - `internal/controller/wails/app_controller.go`
    - `internal/bootstrap/app_controller.go`
  - 変更禁止範囲:
    - `frontend/wailsjs/` の手編集
    - fake API 前提の既存 system test を接続境界の代替にする変更
    - docs 正本本文
    - `.codex/`
  - 想定規模: `1-4 files`, `180-450 changed lines`
- `depends_on`: `FBC-INT-001`, `FBC-UT-FE-001`, `FBC-UT-BE-001`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`, `broad_gate_shared`
- `first_action`: `tests/system/` または `internal/bootstrap/app_controller_test.go` に、fake API を使わず `AppController` bind 面へ到達する最短接続 test の入口を追加する。変更種別は接続境界 test 追加。対応する完了条件は「frontend production factory から backend bind 面までの最短経路を 1 本固定する」である。理由は `TQI-FBC-003` が fake API ではない境界証明の空白を対象にしているためである。
- `completion_signal`:
  - fake API ではない接続境界 test を 1 本追加する。
  - test は UI 表示網羅ではなく、production factory 注入、generated binding 接続、controller response 到達の最短確認に限定する。
  - 既存 `tests/system/translation-job-management.spec.ts` の fake API flow は維持し、接続境界の代替にしない。
  - 接続境界変更のため、実装後ブラウザ確認結果を残す。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite system-test`
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後

## Final Validation

Codex 実装系レーンは全 handoff 完了後に次を確認する。

- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite structure`
- connection boundary を変更したため、実装後ブラウザ確認を行う。

## Completion Packet

Codex 実装系レーンは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `ui_evidence`
- `final_validation_result`
- `harness_gate_result`
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`

## Refactor Lane Handoff Granularity

`refactor_lane` が次に作る実装引き継ぎ入力は、次の ID 単位で切る。

- `FBC-FE-001`
- `FBC-INT-001`
- `FBC-UT-FE-001`
- `FBC-UT-BE-001`
- `FBC-SC-001`

`FBC-INT-001` は shared contract を変更するため、他 handoff と同時開始しない。
`FBC-UT-FE-001` と `FBC-UT-BE-001` は対象ファイルと検証コマンドが分かれるため、`wave-3` で並列可能である。
