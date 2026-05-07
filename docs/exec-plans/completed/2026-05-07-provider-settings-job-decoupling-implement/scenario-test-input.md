# シナリオテスト入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-SCN-01`
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `source_scope`: `implementation-scope.md`

## 目的

採用シナリオ `SCN-PSJD-001` から `SCN-PSJD-006` を、API テスト、frontend test、fakeAPI 実画面確認で証明する。

## 読むファイル

- `scenario-design.md`
- `ui-design.md`
- `docs/frontend-fake-api.md`
- `backend-implementation-result.md`
- `backend-phase-start-implementation-result.md`
- `integration-implementation-result.md`
- `internal/apitest/`
- `internal/integrationtest/`
- `frontend/src/ui/`
- `frontend/src/controller/review-fake-api/`

## 変更許可範囲

- backend API / integration scenario test
- frontend scenario-like UI test
- `test-results/`
- `tmp/agent-browser/`
- `tmp/logs/`

## 禁止範囲

- product code
- docs 正本本文
- 有料の実 AI API 呼び出し
- secret 本体や raw payload を証跡へ出す変更

## 完了条件

- `SCN-PSJD-001` は provider settings 側だけが endpoint と credential 状態を扱うことを APIテストで確認する。
- `SCN-PSJD-002` は Job Setup から Ready job を作成し、選択値だけが見えることを UI人間操作E2E で確認する。
- `SCN-PSJD-003` は APIキー未設定、model list 取得失敗、model 未選択を UI人間操作E2E で確認する。
- `SCN-PSJD-004`、`SCN-PSJD-005`、`SCN-PSJD-006` は APIテストで確認する。
- system test が OS 権限または Wails 起動で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason と再実行コマンドを残す。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-test`
- `python3 scripts/harness/run.py --suite frontend-test`
- `python3 scripts/harness/run.py --suite system-test`
- `npm run dev:wails:agent-browser`
- `agent-browser open "http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management"`
- `agent-browser snapshot`
- `agent-browser errors`

## 期待出力

- `scenario-test-result.md`

