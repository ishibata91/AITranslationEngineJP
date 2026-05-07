# Implementation Scope: 2026-05-07-provider-settings-job-decoupling-implement

- `skill`: `implementation-scope`
- `status`: `completed`
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: `./human-design-review.md`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `task_frame`: `./task-frame.md`
- `scenario_design`: `./scenario-design.md`
- `ui_design`: `./ui-design.md`
- `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `diagramming_result`: `./diagramming-result.md`
- `design_diff_components`: `./design-diff-components.puml`
- `design_diff_sequence`: `./design-diff-sequence.puml`
- `human_design_review`: `./human-design-review.md`
- `detail_specs`: `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`
- `canonical_refs`: `docs/er.md`, `docs/architecture.md`, `docs/frontend-fake-api.md`

## Fixed Decisions

- Job 側 DB は endpoint、secret store 参照実値、`credential_ref` 実値を所有しない。
- Job Setup は provider、model、execution mode、batch mode の選択値だけを扱う。
- Ready job は実行開始前に最新 provider settings を再解決する。
- Running phase は開始時の非 secret 要約だけを保存する。
- provider settings revision と更新履歴は Job 側に持たせない。
- `needs_human_decision`: `0`

## Approved Scope

対象:

- Job Setup から credential 参照選択と endpoint 表示を外す。
- Ready job 作成時の永続値を phase ごとの選択値へ寄せる。
- phase 実行開始時に provider settings を再解決する。
- Running phase に保存する値を非 secret 分類要約へ縮める。
- frontend、backend、Wails / DTO / gateway 境界、単体テスト、シナリオテストを別成果物で実装する。

非対象:

- AIサービス設定画面の全面再設計。
- provider settings の更新履歴、revision 履歴、Job 側 revision 所有。
- secret store 実装方式の置き換え。
- 有料の実 AI API を使う検証。
- docs 正本本文の直接更新。

禁止変更:

- docs 正本本文、`.codex/`、`.codex/skills/`、`.codex/agents/` を Codex implementation レーンで変更しない。
- secret 本体、raw request、raw response、raw prompt、復号可能値を UI、DTO、read model、log、error summary、test fixture に出さない。
- `credential_ref`、`secret_ref`、`api_key`、`token` を参照値と secret 本体の両方の意味で使わない。
- backend、frontend、統合境界を 1 つの大きな実装タスクにまとめない。

## Dependency Policy

UI があるため、frontend 実装を先行する。
frontend 実装後に fakeAPI と実画面 URL を使って人間レビューを行う。
人間レビューが承認された後に backend 実装と統合境界実装へ進む。

backend は 2 つに分ける。
1 つ目は Job Setup 永続化と Job 側 snapshot の所有値削減である。
2 つ目は phase 実行開始時の provider settings 再解決と Running 非 secret 要約である。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `PSJD-FE-01` | `なし` | `なし` | `backend_frontend_order` |
| `human-ui-review-gate` | 人間レビュー | `PSJD-FE-01` | `なし` | fakeAPI 証跡未記録 |
| `wave-2` | `PSJD-BE-01` | `PSJD-FE-01`, 人間 UI レビュー承認 | `なし` | `shared_contract_change` |
| `wave-3` | `PSJD-BE-02` | `PSJD-BE-01` | `なし` | `owned_scope_overlap` |
| `wave-4` | `PSJD-INT-01` | `PSJD-BE-02` | `なし` | `shared_contract_change` |
| `wave-5` | `PSJD-UT-01`, `PSJD-SCN-01` | `PSJD-INT-01` | `PSJD-UT-01 <-> PSJD-SCN-01` | `なし` |

## Handoffs

### `PSJD-FE-01`: Job Setup UI と fakeAPI レビュー

- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `ready_wave`: `wave-1`
- `expected_size`: 通常。想定 12 files、650 changed lines。
- `depends_on`: `なし`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`

読むファイル:

- `./scenario-design.md`
- `./ui-design.md`
- `docs/frontend-fake-api.md`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`
- `frontend/src/application/store/translation-job-setup/translation-job-setup.store.ts`

変更許可範囲:

- `frontend/src/ui/screens/translation-job-setup/`
- `frontend/src/application/usecase/translation-job-setup/`
- `frontend/src/application/presenter/translation-job-setup/`
- `frontend/src/application/store/translation-job-setup/`
- `frontend/src/controller/review-fake-api/`
- 上記範囲の frontend test。

禁止範囲:

- `frontend/src/controller/wails/`
- `frontend/wailsjs/`
- `internal/`
- docs 正本本文。

secret 境界:

- UI / DTO / read model に出してよい値: provider、model、execution mode、batch mode、APIキー状態分類。
- secret 本体: APIキー本体、復号可能値。
- secret 解決責務層: provider settings service と secret store adapter。
- 出力禁止値: endpoint 原文、`credential_ref` 実値、secret store key 名、raw request、raw response、raw prompt、`modelListSourceToken`。

初手:

- path: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- 対象: legacy fallback の `credential reference` select
- 変更種別: 表示削除
- 対応する完了条件: Job Setup に endpoint と credential 参照値を表示しない
- 理由: 人間レビューで最初に確認する UI visible 契約を 1 手で閉じられるため。

完了条件:

- Job Setup は 3 つの翻訳段階で AIサービス、モデル、実行方法、一括処理、APIキー状態分類だけを表示する。
- 作成後要約は endpoint、`credential_ref`、secret store 参照実値を表示しない。
- APIキー未設定、モデル未選択、モデル一覧未更新、モデル一覧取得失敗を別状態で表示する。
- fakeAPI で `success`、`config-missing`、`error` を確認できる。
- 人間レビュー入力として review URL、確認状態、未確認状態、未確認理由を task 成果物へ記録する。

検証コマンド:

- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup src/application/usecase/translation-job-setup src/application/presenter/translation-job-setup src/application/store/translation-job-setup src/controller/review-fake-api`
- `python3 scripts/harness/run.py --suite frontend-local`

人間レビュー記録要求:

- 起動: `npm run dev:wails:agent-browser`
- review URL: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- 追加 URL: `http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#translation-management`
- 追加 URL: `http://localhost:34115/?fakeApi=1&fakeScenario=error#translation-management`
- 操作: `翻訳管理` 内の `セットアップ` を選ぶ。
- 記録: 確認した viewport、確認した状態、未確認状態、未確認理由、`agent-browser snapshot`、`agent-browser errors`。

### `PSJD-BE-01`: Job Setup 永続化境界

- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `ready_wave`: `wave-2`
- `expected_size`: 通常。想定 10 files、600 changed lines。
- `depends_on`: `PSJD-FE-01`, 人間 UI レビュー承認
- `parallelizable_with`: `なし`
- `parallel_blockers`: `shared_contract_change`

読むファイル:

- `./scenario-design.md`
- `./diagramming-result.md`
- `internal/infra/sqlite/dbinit/migrations/003_canonical_er_v1_tables.sql`
- `internal/infra/sqlite/dbinit/migrations/009_translation_job_phase_runtime_snapshots.sql`
- `internal/repository/job_lifecycle_repository.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`
- `internal/service/translation_job_setup_service.go`
- `internal/usecase/translation_job_setup_usecase.go`

変更許可範囲:

- `internal/infra/sqlite/dbinit/migrations/`
- `internal/repository/job_lifecycle_repository.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`
- `internal/service/translation_job_setup_service.go`
- `internal/usecase/translation_job_setup_usecase.go`
- 上記範囲の backend test。

禁止範囲:

- Wails controller / DTO / generated binding の公開契約変更。
- frontend 実装。
- docs 正本本文。

secret 境界:

- UI / DTO / read model に出してよい値: provider、model、execution mode、batch mode、credential 状態分類、接続確認状態、再解決結果分類、再解決時刻。
- secret 本体: APIキー本体、外部 API へ渡す credential。
- secret 解決責務層: provider settings service と secret store adapter。
- 出力禁止値: `credential_ref` 実値、endpoint 原文、endpoint summary、secret store 参照実値、raw payload。

初手:

- path: `internal/repository/job_lifecycle_repository.go`
- 対象: `TranslationJobPhaseRuntimeSnapshot`
- 変更種別: Job 側 snapshot の所有 field 削減
- 対応する完了条件: phase runtime snapshot が非 secret 要約だけを表す
- 理由: repository model を先に縮めると migration、SQL、service の変更漏れを検出しやすい。

完了条件:

- `JOB_PHASE_RUN` と `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は Job Setup 作成時に `credential_ref` 実値と endpoint 系値を所有しない。
- Ready job 作成時に保存する phase 値は provider、model、execution mode、batch mode、credential 状態分類だけである。
- `modelListSourceToken` は Job 側 DB、Job Setup summary、利用者向け表示に残さない。
- provider settings revision と更新履歴を Job 側へ保存しない。
- 既存の Job Setup 作成経路は provider settings を fallback として Job にコピーしない。

検証コマンド:

- `go test ./internal/repository ./internal/infra/sqlite/dbinit ./internal/service ./internal/usecase -run 'JobLifecycle|TranslationJobSetup|ProviderSettings|Migration'`
- `python3 scripts/harness/run.py --suite backend-local`

### `PSJD-BE-02`: phase 実行開始時の再解決

- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `ready_wave`: `wave-3`
- `expected_size`: 通常。想定 12 files、750 changed lines。
- `depends_on`: `PSJD-BE-01`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `owned_scope_overlap`

読むファイル:

- `./scenario-design.md`
- `./design-diff-sequence.puml`
- `internal/service/provider_settings_service.go`
- `internal/service/provider_settings_consumer.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/provider_execution_snapshot.go`

変更許可範囲:

- `internal/service/provider_settings_service.go`
- `internal/service/provider_settings_consumer.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/provider_execution_snapshot.go`
- 上記範囲の backend test。

禁止範囲:

- Job 側 DB へ endpoint 原文、endpoint summary、`credential_ref` 実値を戻す変更。
- provider settings revision を Job 側へ保存する変更。
- raw request / raw response を保存または出力する変更。
- frontend 実装。

secret 境界:

- UI / DTO / read model に出してよい値: credential 状態分類、接続確認状態、再解決結果分類、短い失敗理由。
- secret 本体: provider adapter へ渡す APIキー本体。
- secret 解決責務層: phase 開始時の provider settings 再解決と secret store adapter。
- 出力禁止値: APIキー本体、復号可能値、secret snapshot ref、endpoint 原文、raw payload。

初手:

- path: `internal/service/term_translation_phase_service.go`
- 対象: `resolveExecutionSnapshotForStart`
- 変更種別: 再解決結果の永続保存値を非 secret 要約へ変更
- 対応する完了条件: Ready job 実行開始前に最新 provider settings を再解決する
- 理由: 最初の phase 開始経路で再解決契約を固定し、残り 2 phase を同じ境界へ揃えられる。

完了条件:

- Ready job 実行開始前に最新 provider settings を再解決する。
- provider settings が未設定または参照不能なら Running phase を開始しない。
- Running phase は開始時の非 secret 要約だけを保存する。
- 実行中に provider settings が更新されても、Running phase は途中で設定由来を混在させない。
- Completed phase は provider settings 更新だけで再評価されない。
- Failed phase の再実行は開始時に最新 provider settings を再解決する。

検証コマンド:

- `go test ./internal/service ./internal/usecase -run 'TermTranslation|PersonaGeneration|BodyTranslation|ProviderSettings'`
- `python3 scripts/harness/run.py --suite backend-local`

### `PSJD-INT-01`: Wails / DTO / gateway 接続

- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `ready_wave`: `wave-4`
- `expected_size`: 通常。想定 12 files、700 changed lines。
- `depends_on`: `PSJD-BE-02`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `shared_contract_change`

読むファイル:

- `./scenario-design.md`
- `./ui-design.md`
- `docs/architecture.md`
- `internal/usecase/translation_job_setup_contract.go`
- `internal/controller/wails/translation_job_setup_controller.go`
- `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
- `frontend/src/controller/wails/translation-job-setup.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/translation-job-setup-gateway-dto.ts`

変更許可範囲:

- `internal/usecase/translation_job_setup_contract.go`
- `internal/usecase/translation_job_setup_usecase.go`
- `internal/controller/wails/translation_job_setup_controller.go`
- `frontend/src/application/gateway-contract/translation-job-setup/`
- `frontend/src/controller/wails/translation-job-setup.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/`
- 必要な generated binding 更新。
- 上記範囲の integration / gateway test。

禁止範囲:

- frontend UI layout の再設計。
- backend service の新規仕様追加。
- docs 正本本文。
- secret 本体を DTO に出す変更。

secret 境界:

- UI / DTO / read model に出してよい値: provider、model、execution mode、batch mode、credential 状態分類、接続確認状態、再解決分類。
- secret 本体: APIキー本体、外部 API に送る credential。
- secret 解決責務層: backend provider settings service。
- 出力禁止値: `credential_ref` 実値、secret store key 名、endpoint 原文、raw request、raw response、raw prompt。

初手:

- path: `internal/usecase/translation_job_setup_contract.go`
- 対象: `TranslationJobSetupPhaseRuntimeSelection`
- 変更種別: public request field から Job 所有の credential 参照を外す
- 対応する完了条件: frontend と backend の公開接点が選択値だけを受け渡す
- 理由: Go 側公開契約を先に固定すると Wails DTO と TypeScript gateway の差分を一方向に揃えられる。

完了条件:

- create / validate / summary の公開接点は provider、model、execution mode、batch mode、credential 状態分類だけを扱う。
- `credentialRef`、endpoint、provider settings revision は Job Setup DTO に出ない。
- frontend gateway と Wails controller は同じ field 境界へ揃う。
- fakeAPI と Wails gateway の表示結果で secret と raw payload が出ない。
- 実画面確認で Job Setup の成功状態と不足状態を確認する。

検証コマンド:

- `go test ./internal/controller/wails ./internal/usecase -run 'TranslationJobSetup|ProviderSettings'`
- `npm --prefix frontend run test -- src/controller/wails/translation-job-setup.gateway.test.ts src/controller/wails/gateway-dto/translation-job-setup src/application/gateway-contract/translation-job-setup`
- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`

### `PSJD-UT-01`: 単体テスト

- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `ready_wave`: `wave-5`
- `expected_size`: 通常。想定 15 files、800 changed lines。
- `depends_on`: `PSJD-INT-01`
- `parallelizable_with`: `PSJD-SCN-01`
- `parallel_blockers`: `なし`

読むファイル:

- `./scenario-design.md`
- `./ui-design.md`
- `internal/service/*provider_settings*`
- `internal/service/*phase_service*`
- `internal/repository/job_lifecycle*`
- `frontend/src/application/*/translation-job-setup/`
- `frontend/src/ui/screens/translation-job-setup/`

変更許可範囲:

- backend unit test。
- frontend unit test。
- test fixture と fake object。

禁止範囲:

- product code。
- docs 正本本文。
- secret 本体、raw request、raw response、raw prompt を fixture に入れる変更。

secret 境界:

- test fixture に出してよい値: provider、model、状態分類、短い失敗分類。
- secret 本体: 使用禁止。
- secret 解決責務層: fake secret store。
- 出力禁止値: 実 API key、token、endpoint 原文、復号可能値、raw payload。

初手:

- path: `internal/service/translation_job_setup_service_test.go`
- 対象: Ready job 作成時の phase runtime 保存検証
- 変更種別: 期待値更新
- 対応する完了条件: Job Setup は選択値だけを永続化する
- 理由: backend の永続境界が崩れた場合に最初に赤くなる単体テストである。

完了条件:

- Job Setup の選択値永続化を単体テストで確認する。
- phase 開始時の provider settings 再解決を 3 phase の単体テストで確認する。
- Running 非 secret 要約の保存禁止値を単体テストで確認する。
- frontend presenter / usecase / view は credential 参照値を表示しないことを確認する。
- product code 不足が見つかった場合は修正せず、`implement_lane` へ戻す。

検証コマンド:

- `go test ./internal/...`
- `npm --prefix frontend run test -- src/application src/ui/screens/translation-job-setup src/controller`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

### `PSJD-SCN-01`: シナリオテスト

- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `ready_wave`: `wave-5`
- `expected_size`: 通常。想定 8 files、500 changed lines。
- `depends_on`: `PSJD-INT-01`
- `parallelizable_with`: `PSJD-UT-01`
- `parallel_blockers`: `なし`

読むファイル:

- `./scenario-design.md`
- `./ui-design.md`
- `docs/frontend-fake-api.md`
- `internal/apitest/`
- `internal/integrationtest/`
- `frontend/src/ui/`
- `frontend/src/controller/review-fake-api/`

変更許可範囲:

- backend API / integration scenario test。
- frontend scenario-like UI test。
- `test-results/`
- `tmp/agent-browser/`
- `tmp/logs/`

禁止範囲:

- product code。
- docs 正本本文。
- 有料の実 AI API 呼び出し。
- secret 本体や raw payload を証跡へ出す変更。

secret 境界:

- test 証跡に出してよい値: provider、model、状態分類、再解決分類、短い失敗分類。
- secret 本体: 使用禁止。
- secret 解決責務層: fake secret store と fake transport。
- 出力禁止値: APIキー本体、secret store key 名、endpoint 原文、raw request、raw response、raw prompt。

初手:

- path: `internal/apitest/`
- 対象: `SCN-PSJD-004`
- 変更種別: APIテスト追加
- 対応する完了条件: Ready job 実行開始前に最新 provider settings を再解決する
- 理由: backend 実行境界の中核であり、UI 証跡と独立して判定できる。

完了条件:

- `SCN-PSJD-001` は provider settings 側だけが endpoint と credential 状態を扱うことを APIテストで確認する。
- `SCN-PSJD-002` は Job Setup から Ready job を作成し、選択値だけが見えることを UI人間操作E2E で確認する。
- `SCN-PSJD-003` は APIキー未設定、model list 取得失敗、model 未選択を UI人間操作E2E で確認する。
- `SCN-PSJD-004`、`SCN-PSJD-005`、`SCN-PSJD-006` は APIテストで確認する。
- system test が OS 権限または Wails 起動で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason と再実行コマンドを残す。

検証コマンド:

- `python3 scripts/harness/run.py --suite backend-test`
- `python3 scripts/harness/run.py --suite frontend-test`
- `python3 scripts/harness/run.py --suite system-test`
- `npm run dev:wails:agent-browser`
- `agent-browser open "http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management"`
- `agent-browser snapshot`
- `agent-browser errors`

## Final Validation

全 handoff 完了後に `implement_lane` が最終検証を判断する。
途中 handoff の担当外で失敗した広域検証は、原因を該当 handoff に戻す。

最終候補:

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite all`
- repo-local Sonar issue gate。Sonar server の Quality Gate とは扱わない。

## Canonical Docs Decision

仕様変更を含むため、実装完了後に docs 正本化判断が必要である。
Codex implementation レーンには docs 正本本文を変更させない。
人間承認後、`updating-docs` で次を確認する。

- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-setup.md`
- `docs/er.md`
- `docs/diagrams/er/combined-data-model-er.puml`
- `docs/frontend-fake-api.md`

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `ui_evidence`
- `human_ui_review_result`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate の結果。Sonar server の Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とする。
- `residual_risks`
- `docs_changes: none`
