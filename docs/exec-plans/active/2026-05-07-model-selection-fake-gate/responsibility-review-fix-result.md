# 責務レビュー差し戻し修正結果

## 判定

結果: `responsibility-boundary-001` と `responsibility-boundary-002` の修正は完了した。

対象: service 層の production code から `fake`、`test-safe`、`fake-model`、`RequestsAreTestSafe` を撤去した。

対象: `deterministicHTTPTransport` から model list 直返しを撤去した。

## 指摘ごとの対応

### responsibility-boundary-001

対応: service 層の provider port から `RequestsAreTestSafe` 系 marker method を削除した。

対応: body translation と persona generation の service 層 supported provider set から `fake` provider ID を削除した。

対応: adapter test は service 側 fake provider ID へ依存しない provider ID に直した。

### responsibility-boundary-002

対応: `deterministicHTTPTransport.Do` の model list request 判定を削除した。

対応: model list 固定値は `deterministic_fake_provider.go` の `deterministicProvider.ListModels` に集約した。

対応: `FakeModelID` 定義を `provider_models.go` から `deterministic_fake_provider.go` へ移した。

## 変更ファイル

- `internal/service/body_translation_provider_adapter.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/persona_generation_provider_adapter_test.go`
- `internal/infra/ai/transport.go`
- `internal/infra/ai/provider_models.go`
- `internal/infra/ai/deterministic_fake_provider.go`

## 検証結果

- `rg -n "fake|test-safe|fake-model|RequestsAreTestSafe" internal/service --glob '!**/*_test.go'`: 該当なし。
- `rg -n "isProviderModelListRequest|deterministicProviderModelListResponse|FakeModelID" internal/infra/ai/transport.go internal/infra/ai/provider_models.go internal/infra/ai/deterministic_fake_provider.go`: `FakeModelID` は `deterministic_fake_provider.go` のみ該当。
- `go test ./internal/infra/ai ./internal/service ./internal/bootstrap -run 'Provider|MasterPersona|TranslationJobSetup|AI' -count=1`: 通過。
- `python3 scripts/harness/run.py --suite backend-local`: 通過。

## 残留リスク

- frontend production code は禁止範囲のため未変更である。
- 既存 worktree には本修正前からの dirty diff がある。
- `ProviderFake` は infra provider registry 側に残る。service 層と user-facing provider list には追加していない。
