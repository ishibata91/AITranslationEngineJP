# 実装進行記録

## wave-1

### `OBS-FE-001`

- 状態: 実装完了、影響範囲修正完了、検証通過。
- 担当: `frontend_implementer`
- 変更ファイル:
  - `frontend/src/bootstrap/app-screen-controller-factories.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts`
- 追加ログ:
  - `runtime_event_subscribe`
  - `runtime_event_progress`
  - `runtime_event_completed`
  - `runtime_event_detach`
- 禁止ログ確認: `console.log`、backend 送信、`EventsEmit` は追加していない。
- 初回検証: `python3 scripts/harness/run.py --suite frontend-local` は fail。
- 初回失敗箇所: `frontend/wailsjs/go/wails/AppController.d.ts` が `../models` から `context` を import しているが、`frontend/wailsjs/go/models.ts` に `context` export がない。
- 原因: `AppController` の public lifecycle method が Wails Bind に公開され、`context.Context` が generated frontend API へ漏れた。
- 影響範囲修正:
  - `AppController` から Wails lifecycle hook の public method を外した。
  - Bind されない `AppLifecycle` adapter に startup / shutdown を分離した。
  - generated file は手編集せず、`wails generate module` で再生成した。
  - runtime event adapter の既存 payload 正規化契約を復元した。
- 再検証:
  - `wails generate module` は pass。
  - `python3 scripts/harness/run.py --suite frontend-local` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。

### `OBS-BE-001`

- 状態: 実装完了、検証通過。
- 担当: `observability_implementer`
- 変更ファイル:
  - `internal/service/translation_job_management_service.go`
  - `internal/usecase/translation_job_management_usecase.go`
  - `internal/usecase/term_translation_phase_usecase.go`
  - `internal/usecase/persona_generation_phase_usecase.go`
  - `internal/usecase/body_translation_phase_usecase.go`
- 追加ログ:
  - `translation_job_delete`
  - `translation_job_stop`
  - `translation_job_resume`
  - phase start / pause / resume / retry / cancel / readiness 系 log
- 禁止ログ確認: secret、API key、provider raw payload、prompt 全文、翻訳本文全文、XML 全文、DTO 全体、full path は出していない。
- 検証: `python3 scripts/harness/run.py --suite backend-local` は pass。

## 進行判断

- `wave-2` は `OBS-FE-001` と `OBS-BE-001` の完了後に開始する定義である。
- `OBS-FE-001` と `OBS-BE-001` は検証通過済みである。
- `wave-2` を開始できる。

## 影響範囲修正規約補正

- `implement-frontend` は、`frontend-local` の失敗原因が承認済み実装範囲外でも、frontend プロダクトコードの影響範囲修正なら止まらず直す規約を補強済み。
- `implement-frontend` は、generated file が原因の場合は直接編集せず、生成元または公開境界を直す規約を補強済み。
- `implement-lane` は、`backend-local` / `frontend-local` の範囲外失敗を、対応する実装またはテスト担当 agent へ影響範囲修正として渡す規約を補強済み。

## wave-2

### `OBS-BE-002`

- 状態: 実装完了、検証通過。
- 担当: `observability_implementer`
- 変更ファイル:
  - `internal/service/provider_settings_service.go`
  - `internal/service/translation_job_setup_service.go`
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
- 追加ログ:
  - `provider_settings_validation`
  - `provider_model_list`
  - `provider_execution_settings`
  - `term_translation_provider_*`
  - `persona_generation_provider_*`
  - `body_translation_provider_*`
- 禁止ログ確認: credential 参照実値、secret 本体、API key、endpoint 実値、provider raw payload、prompt 全文、翻訳本文全文は出していない。
- 検証: `go test ./internal/service` は pass。`python3 scripts/harness/run.py --suite backend-local` は pass。

### `OBS-BE-003`

- 状態: 実装完了、影響範囲修正完了、検証通過。
- 担当: `observability_implementer`
- 変更ファイル:
  - `internal/controller/wails/translation_input_controller.go`
  - `internal/controller/wails/translation_output_artifact_controller.go`
  - `internal/service/translation_input_import_service.go`
  - `internal/service/translation_job_setup_service.go`
  - `internal/service/translation_output_artifact_xml_adapter.go`
  - `internal/service/term_translation_phase_service.go`
- 追加ログ:
  - `request_invalid`
  - `invalid_json`
  - `source_file_missing`
  - `cache_missing`
  - `db_save_failed`
  - `transaction_failed`
  - `file_write_failed`
  - `response_mapping_failed`
- 影響範囲修正: `internal/service/term_translation_phase_service.go` の `classifyTermTranslationCredentialFailure` 関数分割を復旧した。
- 禁止ログ確認: DTO 全体、full path、XML 全文、翻訳本文全文は出していない。
- 検証: `python3 scripts/harness/run.py --suite backend-local` は pass。

### `OBS-UNIT-FE-001`

- 状態: 実装完了、影響範囲修正完了、検証通過。
- 担当: `implementation_unit_tester`
- 変更ファイル:
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts`
- 追加テスト:
  - runtime 不在時の `skipped` / `runtime_unavailable`
  - progress payload parse 失敗時の `dropped` / `payload_parse_failed`
  - invalid progress 時の `skipped`
  - progress accepted
  - detach 時の `detached`
- 影響範囲修正: progress の非 object payload を `dropped` として log に出し、store を更新しない分岐を追加した。
- 検証: `python3 scripts/harness/run.py --suite frontend-local` は pass。

## wave-2 進行判断

- `OBS-BE-002`、`OBS-BE-003`、`OBS-UNIT-FE-001` は完了した。
- `wave-3` の `OBS-BE-004` を開始できる。

## wave-3

### `OBS-BE-004`

- 状態: 実装完了、検証通過。
- 担当: `observability_implementer`
- 変更ファイル:
  - `internal/service/translation_input_import_service.go`
  - `internal/service/master_dictionary_import_service.go`
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
- 追加ログ:
  - `translation_input_import_bulk_summary`
  - `master_dictionary_import_bulk_summary`
  - `term_translation_provider_bulk_summary`
  - `persona_generation_bulk_summary`
  - `body_translation_provider_bulk_summary`
- 追加しない理由: `internal/service/xtranslator_output_row_builder.go` は 1 row 単位の処理であり、ここへ log を入れると loop 内 1 件ごとの log になるため変更しない。
- 禁止ログ確認: provider raw payload、prompt 全文、翻訳本文全文、XML 全文、secret、API key、DTO 全体、full path は出していない。
- 検証: `go test ./internal/service` は pass。`python3 scripts/harness/run.py --suite backend-local` は pass。

## wave-3 進行判断

- `OBS-BE-004` は完了した。
- `wave-4` の `OBS-UNIT-BE-001` と `OBS-SCN-BE-001` を開始できる。

## wave-4

### `OBS-UNIT-BE-001`

- 状態: 実装完了、影響範囲修正完了、検証通過。
- 担当: `implementation_unit_tester`
- 変更ファイル:
  - `internal/infra/runtime/diagnostic_log_test.go`
  - `internal/service/provider_settings_service_test.go`
  - `internal/service/translation_job_management_service_test.go`
  - `internal/service/translation_input_import_service_test.go`
  - `internal/service/body_translation_phase_service_test.go`
- 追加テスト:
  - diagnostic log payload に禁止語と undefined 相当値が混入しないこと。
  - provider 境界失敗 log が安全な payload だけを持つこと。
  - job 削除拒否 log が安全な payload だけを持つこと。
  - 入力取り込み集約 log が件数中心の安全な payload だけを持つこと。
  - 本文翻訳 provider 集約 log が件数中心の安全な payload だけを持つこと。
- 禁止ログ確認: `api_key`、`endpoint`、`raw_request`、`raw_response`、`prompt`、`dto`、`full_path`、`trace_id` は出していない。
- 影響範囲修正:
  - `internal/infra/runtime/diagnostic_log_test.go` の staticcheck `S1030` を修正した。
  - `string(buffer.Bytes())` を `buffer.String()` へ変更した。
- 検証:
  - `python3 scripts/harness/run.py --suite backend-local` は pass。
  - 初回 `python3 scripts/harness/run.py --suite coverage` は Sonar maintainability HIGH issue 10 件で fail。

### `OBS-SCN-BE-001`

- 状態: 実装完了、検証通過。
- 担当: `implementation_scenario_tester`
- 変更ファイル:
  - `internal/apitest/observability_log_scenario_test.go`
  - `internal/apitest/observability_log_test_helpers_test.go`
- 追加テスト:
  - 削除 command 境界で、状態変更の許可、拒否、状態投影不整合を log で区別する。
  - provider settings command 境界で、credential 不足を provider / secret 失敗として分類する。
  - translation input command 境界で、source file missing を file 境界失敗として分類する。
  - translation input command 境界で、大量処理の集約 count だけを検査する。
- 禁止ログ確認: secret 相当値、API key 風の値、endpoint 実値、provider raw payload、raw request / response、入力 JSON 全文は出していない。
- 検証:
  - `GOCACHE=/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/go-cache go test ./internal/apitest` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。

## coverage 影響範囲修正

### service 層

- 状態: 実装完了、検証通過。
- 担当: `observability_implementer`
- 変更ファイル:
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/provider_settings_service.go`
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/translation_job_setup_service.go`
- 修正内容:
  - service 名の `where` 値を定数化し、Sonar `go:S1192` を解消した。
  - `executeBodyTranslationRun` の pending target 処理を補助関数へ分割し、Sonar `go:S3776` を解消した。
- 保持した契約: 観測ログの `event`、`result`、`where` の実値は変更していない。
- 検証:
  - `go test ./internal/service` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。
  - `python3 scripts/harness/run.py --suite coverage` は pass。

### usecase 層

- 状態: 実装完了、検証通過。
- 担当: `backend_implementer`
- 変更ファイル:
  - `internal/usecase/body_translation_phase_usecase.go`
  - `internal/usecase/persona_generation_phase_usecase.go`
  - `internal/usecase/term_translation_phase_usecase.go`
  - `internal/usecase/translation_job_management_usecase.go`
- 修正内容:
  - usecase の `where` 値を定数化し、Sonar `go:S1192` を解消した。
  - `job:%d` 書式を定数化し、Sonar `go:S1192` を解消した。
- 保持した契約: 観測ログの `event`、`result`、`where`、`id` の実値は変更していない。
- 検証:
  - `go test ./internal/usecase` は pass。
  - service 層修正後の親側 `python3 scripts/harness/run.py --suite backend-local` は pass。
  - service 層修正後の親側 `python3 scripts/harness/run.py --suite coverage` は pass。

## 最終検証

- `git diff --check`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- coverage summary: Sonar coverage `70.5%`、line `71.7%`、branch `61.9%`。
- Sonar issue gate:
  - security issues: `0 <= 0`
  - reliability issues: `0 <= 0`
  - maintainability HIGH issues: `0 <= 0`

## wave-4 進行判断

- `OBS-UNIT-BE-001` と `OBS-SCN-BE-001` は完了した。
- coverage 影響範囲修正は完了した。
- `最終検証` は完了した。
- `実装後ブラウザ確認` は完了した。
- 次成果物は `レビュー通過根拠` である。

## 実装後ブラウザ確認

- 担当: `browser_confirmation`
- 証跡: `browser-confirmation/2026-05-09-browser-confirmation.md`
- 操作確認結果:
  - `agent-browser doctor --offline --quick` は pass。
  - `agent-browser open http://localhost:34115` は pass。
  - 初期画面はダッシュボードを表示した。blank ではない。
  - マスター辞書画面へ遷移できた。
  - マスター辞書画面は辞書一覧と詳細領域を表示した。blank ではない。
- 異常記録:
  - `agent-browser errors` では具体的な error text は確認していない。
  - console では `runtime_event_subscribe` を含む frontend log を確認した。
  - console では `wails dev Connected to backend` と `Disconnected from backend` が反復していた。
  - network は初期読込の 200 応答を確認した。
- 未確認理由:
  - master dictionary runtime event の詳細なイベント列は、今回の画面遷移確認の範囲外として深追いしていない。
  - master dictionary の更新や削除は実行していない。

## レビュー指摘修正

### `behavior-001`

- 状態: 修正完了、検証通過。
- 担当: `frontend_implementer`
- 指摘: completed runtime event の malformed payload が `accepted` になり、完了処理へ進む。
- 変更ファイル:
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
- 修正内容:
  - completed event の payload が object でない場合は `dropped` / `payload_parse_failed` として warn log を出す。
  - completed event の payload が object でない場合は `onImportCompleted` を呼ばない。
  - valid object payload は従来どおり accepted として処理する。
- 検証:
  - `npm --prefix frontend run test -- master-dictionary-runtime-event-adapter.test.ts` は pass。1 file、16 tests passed。
  - `python3 scripts/harness/run.py --suite frontend-local` は pass。

### `behavior-002`

- 状態: 修正完了、検証通過。
- 担当: `backend_implementer`
- 指摘: phase state log の `before_state` が実状態ではなく固定前提値になる。
- 変更ファイル:
  - `internal/usecase/translation_job_management_usecase.go`
  - `internal/usecase/term_translation_phase_usecase.go`
  - `internal/usecase/persona_generation_phase_usecase.go`
  - `internal/usecase/body_translation_phase_usecase.go`
  - `internal/usecase/translation_job_management_usecase_test.go`
- 修正内容:
  - phase command log の入力から `idle_ready`、`running`、`paused_or_recoverable_failed`、`recoverable_failed` の固定前提値を除去した。
  - accepted 経路では、実 read model の `PhaseState` を `after_state` に出し、操作前実状態が取れない `before_state` は `unknown` にする。
  - rejected 経路では、実 read model の `PhaseState` が取れている場合に `before_state` と `after_state` を同じ実状態にする。
  - 実状態が空の場合だけ、両方を `unknown` にする。
- 検証:
  - `go test ./internal/usecase` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。

### `contract-001`

- 状態: 修正完了、検証通過。
- 担当: workflow 契約修正 worker
- 指摘: 影響範囲修正の許可範囲が承認済み UI 契約を越えている。
- 変更ファイル:
  - `.codex/skills/implement-frontend/SKILL.md`
  - `.codex/skills/implement-lane/SKILL.md`
- 修正内容:
  - `implement-frontend` の影響範囲修正を、generated file、生成元、公開境界、検証を壊した frontend プロダクトコードに限定した。
  - 実画面と UI 根拠の差分が承認済み実装範囲外の場合は、停止して人間承認へ戻す規約にした。
  - UI 表示、画面文言、`layout`、`style`、承認済み UI 根拠を越える変更は、影響範囲修正に使わないと明記した。
- 検証:
  - `git diff --check -- .codex/skills/implement-frontend/SKILL.md .codex/skills/implement-lane/SKILL.md` は pass。

## 修正後検証

- `git diff --check`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- coverage summary: Sonar coverage `70.5%`、line `71.7%`、branch `61.9%`。
- Sonar issue gate:
  - security issues: `0 <= 0`
  - reliability issues: `0 <= 0`
  - maintainability HIGH issues: `0 <= 0`

## 再レビュー指摘修正

### `state-invariant-001`

- 状態: 修正完了、検証通過。
- 担当: `frontend_implementer`
- 指摘: object 型の malformed completed payload が store 完了へ進む。
- 変更ファイル:
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
- 修正内容:
  - completed payload は object であるだけでは accepted にしない。
  - `summary` または `page` が有効な shape の場合だけ accepted にする。
  - 空 object、未知 key だけの object、不正な `page` / `summary` は `invalid_payload` として dropped にする。
  - invalid object payload では `onImportCompleted` を呼ばない。
- 検証:
  - `npm --prefix frontend run test -- master-dictionary-runtime-event-adapter.test.ts` は pass。21 tests passed。
  - `python3 scripts/harness/run.py --suite frontend-local` は pass。

### `behavior-002` 再修正

- 状態: 修正完了、検証通過。
- 担当: `backend_implementer`
- 指摘: 許可経路の `before_state` が `unknown` になり、変更前状態を残す要件を満たしていない。
- 変更ファイル:
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
  - `internal/usecase/translation_job_management_usecase.go`
  - `internal/usecase/translation_job_management_usecase_test.go`
  - `internal/usecase/term_translation_phase_usecase.go`
  - `internal/usecase/persona_generation_phase_usecase.go`
  - `internal/usecase/body_translation_phase_usecase.go`
- 修正内容:
  - service の command read model から `BeforePhaseState` / `AfterPhaseState` を usecase へ渡す。
  - start は service が判断した開始前 state を `before_state` として出す。
  - pause / resume / retry / cancel は更新前 `run.State` を `before_state` として出す。
  - rejected は実状態が取れる場合に `before_state` と `after_state` を同じ実状態にする。
  - 実状態が取れない service error だけ `unknown` に戻す。
- 検証:
  - `go test ./internal/usecase` は pass。
  - `go test ./internal/service` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。

### coverage 影響範囲修正

- 状態: 修正完了、検証通過。
- 担当: `backend_implementer`
- 指摘: `internal/service/body_translation_phase_service.go:1340` の Sonar `go:S3776`。
- 変更ファイル:
  - `internal/service/body_translation_phase_service.go`
- 修正内容:
  - `persistBodyTranslationRunStateTransition` から、状態遷移の事前検証を `validateBodyTranslationRunStateTransition` へ切り出した。
  - 状態遷移結果、error 文言、transaction 内の更新順序は変更していない。
  - `beforePhaseState` の返却と error 時の戻り値は維持した。
- 検証:
  - `go test ./internal/service` は pass。
  - `python3 scripts/harness/run.py --suite backend-local` は pass。
  - `python3 scripts/harness/run.py --suite coverage` は pass。

## 再修正後検証

- `git diff --check`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- coverage summary: Sonar coverage `70.6%`、line `71.8%`、branch `62.3%`。
- Sonar issue gate:
  - security issues: `0 <= 0`
  - reliability issues: `0 <= 0`
  - maintainability HIGH issues: `0 <= 0`

## workflow 契約修正

### `trust-boundary-001`

- 状態: 修正完了、レビュー通過。
- 担当: workflow 契約修正 worker。
- 指摘: 承認済み範囲外の影響範囲修正一般許可が複数 skill と agent TOML に残っていた。
- 変更ファイル:
  - `.codex/README.md`
  - `.codex/agents/backend_implementer.toml`
  - `.codex/agents/frontend_implementer.toml`
  - `.codex/agents/implementation_scenario_tester.toml`
  - `.codex/agents/implementation_unit_tester.toml`
  - `.codex/agents/integration_implementer.toml`
  - `.codex/agents/observability_implementer.toml`
  - `.codex/skills/implement-backend/SKILL.md`
  - `.codex/skills/implement-frontend/SKILL.md`
  - `.codex/skills/implement-integration/SKILL.md`
  - `.codex/skills/implement-lane/SKILL.md`
  - `.codex/skills/implementation-investigate/SKILL.md`
  - `.codex/skills/observability-implementer/SKILL.md`
  - `.codex/skills/tests-scenario/SKILL.md`
  - `.codex/skills/tests-unit/SKILL.md`
- 修正内容:
  - 影響範囲修正を、今回変更が直接壊した生成物、生成元、公開境界、検証経路、担当責務内成果物に限定した。
  - UI 表示、画面文言、layout、style、承認済み UI 根拠を越える変更は停止または人間返却にした。
  - secret、trust boundary、API / DTO / DB / schema の意味拡張、docs 正本化、`.codex` 作業流れ変更は停止または禁止にした。
- 検証:
  - `git diff --check`: pass。
  - 編集 TOML の `tomllib` 構文確認: pass。
  - `reviewback.trust-boundary.yaml`: `no_issue`。

### `contract-002`

- 状態: 修正完了、レビュー通過。
- 担当: workflow 契約修正 worker。
- 指摘: `implement-lane` の非対象規約が、review agent へ検証証跡を渡す契約と衝突していた。
- 変更ファイル:
  - `.codex/skills/implement-lane/SKILL.md`
- 修正内容:
  - review agent に渡さない資料の例から `ハーネス結果など` を削除した。
  - review agent 起動入力には、レビュー対象差分、実装目的、承認済み実装範囲、実装結果、検証証跡、変更ファイル、review YAML path を含める契約として整理した。
- 検証:
  - `git diff --check`: pass。
  - `reviewback.contract.yaml`: `no_issue`。

## レビュー通過根拠

- `reviewback.behavior.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.contract.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.trust-boundary.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.state-invariant.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.responsibility-boundary.yaml`: `issues_open`、`must_fix_open: false`、`max_level: minor`。
- `implementation_action`: `close`。
- 残留 minor: provider settings service に横断 provider log helper と persona generation 固有変換が同居している。修正必須ではない。
