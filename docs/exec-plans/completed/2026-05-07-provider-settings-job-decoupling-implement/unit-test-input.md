# 単体テスト入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-UT-01`
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `source_scope`: `implementation-scope.md`

## 目的

実装済み責務を単体テストで証明する。
Job Setup の選択値永続化、phase 開始時の provider settings 再解決、Running 非 secret 要約、frontend 表示禁止値を確認する。

## 読むファイル

- `scenario-design.md`
- `ui-design.md`
- `backend-implementation-result.md`
- `backend-phase-start-implementation-result.md`
- `integration-implementation-result.md`
- `internal/service/*provider_settings*`
- `internal/service/*phase_service*`
- `internal/repository/job_lifecycle*`
- `frontend/src/application/*/translation-job-setup/`
- `frontend/src/ui/screens/translation-job-setup/`

## 変更許可範囲

- backend unit test
- frontend unit test
- test fixture と fake object

## 禁止範囲

- product code
- docs 正本本文
- secret 本体、raw request、raw response、raw prompt を fixture に入れる変更

## 完了条件

- Job Setup の選択値永続化を単体テストで確認する。
- phase 開始時の provider settings 再解決を 3 phase の単体テストで確認する。
- Running 非 secret 要約の保存禁止値を単体テストで確認する。
- frontend presenter / usecase / view は credential 参照値を表示しないことを確認する。
- product code 不足が見つかった場合は修正せず、`implement_lane` へ戻す。

## 検証コマンド

- `go test ./internal/...`
- `npm --prefix frontend run test -- src/application src/ui/screens/translation-job-setup src/controller`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

## 期待出力

- `unit-test-result.md`

