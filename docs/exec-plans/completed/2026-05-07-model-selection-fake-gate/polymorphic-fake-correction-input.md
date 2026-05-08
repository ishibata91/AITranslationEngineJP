# 責務境界修正入力

## 起動先

- `agent`: `implementation_implementer`
- `skill`: `implement-integration`
- `artifact`: `責務境界修正実装証跡`

## 必須前提

- fake は実 provider と同じ provider interface を実装する代替 provider である。
- fake mode 判定は DI 構成だけで使う。
- service 層は `fake`、`test-safe`、`fake-model` を知らない。
- credential と endpoint の取得可否は fake 判定と結びつけない。
- fake provider は credential と endpoint を受け取っても provider 実装内で許可または無視する。

## 修正対象

- [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [master_persona_provider_transport.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_provider_transport.go)
- [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
- [provider_client.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider_client.go)
- [provider_models.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider_models.go)
- [persona_generation_provider.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/persona_generation_provider.go)
- [transport.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/transport.go)
- [app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go)

## 撤去対象

- `ProviderSettingsService` の `ProviderModelListRequestsAreTestSafe` 系 interface、method、分岐。
- `ProviderSettingsService` の `test-safe` だけで credential / endpoint gate を bypass する処理。
- `MasterPersonaGenerationService` の `fakeModelDefaults` と `LoadSettings()` の `fake-model` 直返し。
- model list helper 内の transport marker による `fake-model` 直返し。

## 期待する構造

- infra 側に、実 provider と同じ interface を満たす fake provider 実装を置く。
- fake provider の実装ファイルは、実 provider と区別できる名前にする。
- 例: `fake_provider.go`、`deterministic_fake_provider.go`、`fake_model_list_provider.go` のように、fake 境界がファイル名で分かる名前にする。
- `persona_generation_provider.go` のような実 provider と区別しにくい名前へ fake 実装を置き続けない。
- model list も generation も provider 実装の責務に寄せる。
- bootstrap は fake mode の時に provider 実装または registry を差し替えるだけにする。
- service は provider ID、credential reference、endpoint、model list 結果を扱い、fake mode を判定しない。

## 禁止変更

- frontend production code に fake 判定を入れない。
- user-facing provider list に `fake` provider を追加しない。
- `fake-model` 固有分岐を service / frontend に入れない。
- task 外の既存 dirty diff を戻さない。

## 検証観点

- fake mode でも service 層に fake/test-safe 判定が存在しない。
- fake provider 実装のファイル名から、実 provider ではなく fake 境界であることが分かる。
- fake provider は同じ provider interface で model list と generation を返す。
- real provider の credential / endpoint gate は service の通常 preflight と provider 実装で維持される。
- master persona 設定永続化テストで保存済み model を `fake-model` に上書きしない。

## 期待する成果物

- `docs/exec-plans/active/2026-05-07-model-selection-fake-gate/polymorphic-fake-correction-result.md`
- 変更ファイル、撤去した漏れ、追加した fake provider 境界、検証結果、残留リスクを記録する。
