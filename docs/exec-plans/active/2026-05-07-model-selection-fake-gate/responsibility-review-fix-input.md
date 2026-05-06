# 責務レビュー差し戻し修正入力

## 起動先

- `agent`: `implementation_implementer`
- `skill`: `implement-integration`
- `artifact`: `責務レビュー差し戻し実装証跡`

## 入力

- [reviewback.responsibility-boundary.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.responsibility-boundary.yaml)
- [polymorphic-fake-correction-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/polymorphic-fake-correction-input.md)
- [polymorphic-fake-correction-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/polymorphic-fake-correction-result.md)

## 必須前提

- fake は実 provider と同じ provider interface を実装する代替 provider である。
- fake mode 判定は DI 構成だけで使う。
- service 層は `fake`、`test-safe`、`fake-model` を知らない。
- fake provider の実装ファイル名は、実 provider と明確に区別できる。

## 修正対象指摘

- `responsibility-boundary-001`: service 層 provider port に fake provider ID と test-safe marker method が残っている。
- `responsibility-boundary-002`: `deterministicHTTPTransport.Do` に model list の直返しが残っている。

## 修正対象ファイル

- [body_translation_provider_adapter.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/body_translation_provider_adapter.go)
- [persona_generation_provider_adapter.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/persona_generation_provider_adapter.go)
- [term_translation_provider_adapter.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/term_translation_provider_adapter.go)
- [transport.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/transport.go)
- 関連する infra adapter とテスト。

## 期待する修正

- service production code から fake provider ID を撤去する。
- service production code から `RequestsAreTestSafe` 系 marker method を撤去する。
- fake/test-safe の判断を service port へ出さない。
- `deterministicHTTPTransport` は HTTP 応答 seam に閉じ、model list の固定値を返さない。
- model list 固定値は `deterministic_fake_provider.go` の fake provider 実装に集約する。

## 検証

- `rg -n "fake|test-safe|fake-model|RequestsAreTestSafe" internal/service --glob '!**/*_test.go'`
- `rg -n "isProviderModelListRequest|deterministicProviderModelListResponse|FakeModelID" internal/infra/ai/transport.go internal/infra/ai/provider_models.go internal/infra/ai/deterministic_fake_provider.go`
- `go test ./internal/infra/ai ./internal/service ./internal/bootstrap -run 'Provider|MasterPersona|TranslationJobSetup|AI' -count=1`
- `python3 scripts/harness/run.py --suite backend-local`

## 成果物

- `docs/exec-plans/active/2026-05-07-model-selection-fake-gate/responsibility-review-fix-result.md`
- 指摘ごとの対応、変更ファイル、検証結果、残留リスクを書く。
