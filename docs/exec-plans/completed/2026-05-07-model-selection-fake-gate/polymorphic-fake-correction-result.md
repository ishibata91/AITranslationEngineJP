# 責務境界修正実装証跡

## 変更ファイル

- `internal/infra/ai/deterministic_fake_provider.go`
- `internal/infra/ai/provider.go`
- `internal/infra/ai/provider_client.go`
- `internal/infra/ai/provider_models.go`
- `internal/infra/ai/gemini.go`
- `internal/infra/ai/openai_compatible.go`
- `internal/infra/ai/transport.go`
- `internal/bootstrap/app_controller.go`
- `internal/service/provider_settings_service.go`
- `internal/service/master_persona_provider_transport.go`
- `internal/service/master_persona_service.go`
- `internal/service/translation_job_setup_service.go`
- `internal/infra/ai/provider_client_test.go`
- `internal/service/provider_settings_service_test.go`
- `internal/service/master_persona_provider_transport_test.go`
- `internal/service/master_persona_service_test.go`

## 撤去した漏れ

- `ProviderSettingsService` から `ProviderModelListRequestsAreTestSafe` 系 interface、method、分岐を撤去した。
- `ProviderSettingsService` の `test-safe` による credential / endpoint bypass を撤去した。
- `MasterPersonaGenerationService` の `fakeModelDefaults` と `WithMasterPersonaFakeModelDefaults` を撤去した。
- `MasterPersonaGenerationService.LoadSettings()` の `fake-model` 直返しを撤去した。
- `TranslationJobSetupService` に残っていた `ProviderModelListRequestsAreTestSafe` 系分岐を撤去した。
- model list helper の transport marker による `fake-model` 直返しを撤去した。

## 追加した境界

- fake provider は `internal/infra/ai/deterministic_fake_provider.go` に分けた。
- fake provider のファイル名で fake 境界を可視化した。
- fake provider は実 provider と同じ `provider` interface で `Generate` と `ListModels` を実装する。
- bootstrap は fake mode の時に `WithDeterministicProviderRegistry` と `WithProviderModelListDeterministicProviders` だけを差し込む。
- service 層は fake mode、`test-safe`、`fake-model` を判定しない。

## 検証結果

- `rg -n "ProviderModelListRequestsAreTestSafe|providerModelListRequestsAreTestSafe|fakeModelDefaults|WithMasterPersonaFakeModelDefaults|fake-model|testSafe" internal/service internal/infra/ai internal/bootstrap/app_controller.go`
  - 撤去済み: `ProviderSettingsService`、`MasterPersonaGenerationService`、`TranslationJobSetupService`、bootstrap adapter の marker 分岐。
  - 残留: infra の `testSafeHTTPTransport` marker、`ProviderClient` の transport 判定、`FakeModelID`、関連テスト内の固定値。
  - 理由: 残留は infra のテスト用 transport 境界またはテスト fixture であり、service / frontend の fake 判定ではない。
- `go test ./internal/infra/ai ./internal/service ./internal/bootstrap -run 'Provider|MasterPersona|TranslationJobSetup|AI' -count=1`
  - 結果: 通過。
  - 備考: sandbox 内の Go build cache 権限で失敗したため、承認済み通常権限で再実行した。
- `python3 scripts/harness/run.py --suite backend-local`
  - 結果: 通過。

## 残留リスク

- `ProviderClient.ProviderRequestsAreTestSafe()` と `testSafeHTTPTransport` は infra / adapter 系の既存契約として残る。
- `FakeModelID` は fake provider の model list 戻り値として infra に残る。
- user-facing provider list には `fake` provider を追加していない。
