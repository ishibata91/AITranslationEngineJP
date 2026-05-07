# シナリオテスト結果

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-SCN-01`
- `implementation_skill`: `tests-scenario`
- `result`: `passed`

## 実装結果

- `SCN-PSJD-001`: provider settings API は endpoint と APIキー入力状態を扱い、Job Setup の create / validate / summary DTO は endpoint、credential 参照値、token、raw payload を JSON に出さないことを API テストで確認した。
- `SCN-PSJD-002`: Job Setup の Ready job 作成後要約が provider、model、execution mode、batch mode、APIキー状態分類だけを表示することを frontend scenario-like UI test で確認した。
- `SCN-PSJD-003`: APIキー未設定、model list 取得失敗、model 未選択を画面上で別状態として表示し、APIキー未設定では model list 更新操作が抑止されることを frontend scenario-like UI test で確認した。
- `SCN-PSJD-004`: term、persona、body の phase 開始 API 境界で provider settings を再解決することを API テストで確認した。
- `SCN-PSJD-005`: SQLite runtime snapshot が provider、model、credential 状態分類、execution mode、batch mode だけを永続化し、credential 参照値、endpoint summary、model list token を残さないことを API テストで確認した。
- `SCN-PSJD-006`: retry が既存 phase run API 上で再解決し、attempt history table と raw payload 表示 DTO を増やさないことを API テストで確認した。

## 変更ファイル

- [provider_settings_job_decoupling_scenario_test.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/apitest/provider_settings_job_decoupling_scenario_test.go): `SCN-PSJD-001`、`SCN-PSJD-004`、`SCN-PSJD-005`、`SCN-PSJD-006` の API / integration scenario test を追加した。
- [JobSetupPage.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts): `SCN-PSJD-002`、`SCN-PSJD-003` の frontend scenario-like UI test を追加した。
- [scenario-test-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-test-result.md): シナリオテスト結果と検証証跡を記録した。

## 検証結果

- `go test ./internal/apitest -run 'SCN_PSJD'`: passed。
- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup/JobSetupPage.test.ts`: passed。10 tests passed。
- `python3 scripts/harness/run.py --suite backend-test`: passed。
- `python3 scripts/harness/run.py --suite frontend-test`: passed。57 files、496 tests passed。
- `python3 scripts/harness/run.py --suite system-test`: passed。9 tests passed。

## UI / agent-browser 証跡

- 起動: `npm run dev:wails:agent-browser`
- 成功 URL: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- 成功状態: `セットアップ` tab で 3 phase が `設定済み`、model `fake-model`、作成前確認 `不足はありません。`、`次へ` enabled を表示した。
- 不足 URL: `http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#translation-management`
- 不足状態: `セットアップ` tab で 3 phase が `APIキー未設定`、model list 更新 button disabled、`次へ` disabled を表示した。
- 失敗 URL: `http://localhost:34115/?fakeApi=1&fakeScenario=error#translation-management`
- 失敗状態: `セットアップ` tab で `Job Setup の確認に失敗しました。` と `作成前確認はまだ未完了です。`、`次へ` disabled を表示した。
- `agent-browser errors`: 成功、不足、失敗の各確認で出力なし。

## 残留事項

- `system-test` は環境要因で停止しなかったため、`FAIL_ENVIRONMENT` は発生していない。
- product code 不足は検出していない。
- 有料の実 AI API 呼び出しは実施していない。
