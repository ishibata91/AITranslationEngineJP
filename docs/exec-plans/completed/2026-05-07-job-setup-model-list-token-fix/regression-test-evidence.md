# 回帰テスト証跡

## 結果

- 結果: 完了
- テスト担当: `implementation_unit_tester`

## 変更ファイル

- `internal/service/translation_job_setup_service_test.go`

## 証明内容

- provider 設定 consumer に渡す `RequestToken` が provider 設定 summary の token であることを証明した。
- Job Setup response の `RequestToken` が画面操作用 token のまま維持されることを証明した。

## 検証

- `go test ./internal/service -run 'TranslationJobSetupServiceProviderSettingsTestSafeModelListAllowsMissingCredential'`: 通過
- `go test ./internal/service`: 通過

## 根拠参照

- `internal/service/translation_job_setup_service_test.go:602`
- `internal/service/translation_job_setup_service_test.go:610`
