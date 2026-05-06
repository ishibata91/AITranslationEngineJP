# 実装証跡

## 判定

- 判定: 一部完了
- 担当: `implementation_implementer`
- skill: `implement-integration`
- 入力: [integration-implementation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/integration-implementation-input.md)

## 変更ファイル

- [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)

## 実装内容

- `ProviderSettingsService` は、test-safe な model list loader の時に secret と endpoint の事前 gate を通さず model list loader を呼ぶ。
- test-safe な model list 結果は `ready` と `not_required` を返し、secret 値を出力しない。
- `TranslationJobSetupUseCase.refreshPhaseModels()` は、credential missing だけで backend 呼び出し前に停止しない。
- model list 応答が 1 件だけなら、model ID の値に依存せず、その 1 件を選択状態へ入れる。

## 検証結果

- `npm --prefix frontend run test -- --run ...`: 成功。2 files、21 tests 通過。
- `go test ./internal/service ./internal/bootstrap ./internal/controller/wails`: 失敗。`internal/bootstrap` の master persona 設定永続化 2 件で停止した。
- `python3 scripts/harness/run.py --suite backend-local`: 失敗。`internal/bootstrap` の master persona 設定永続化 2 件で停止した。
- `python3 scripts/harness/run.py --suite frontend-local`: 失敗。既存の frontend test 型エラーで停止した。

## 残留リスク

- `ResolveProviderExecutionSettings` の test-safe 判定に `Endpoint != nil` 条件が残っている。endpoint が nil の保存状態で実行設定解決が止まる可能性がある。
- 実装担当の返却では task 内成果物ファイルが作成されなかったため、`fix_lane` が返却内容をこのファイルへ転記した。
- 回帰テスト証跡は未完了である。
