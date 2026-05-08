# 修正実行入力

## 起動先

- `agent`: `implementation_implementer`
- `skill`: `implement-integration`
- `artifact`: `実装証跡`

## 人間観測

- fake mode なら、credential と endpoint の取得可否や内容に関係なく AI model を使えるべきである。
- fake 判定は frontend へ波及させない。
- frontend は `fake-model` 文字列を特別扱いしない。
- 人間観測の正本は [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation.md) とする。

## 修正前調査

- 調査結果は [investigation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/investigation.md) とする。
- 観測: `ProviderModelListRequestsAreTestSafe()` は backend service へ伝搬している。
- 観測: `TranslationJobSetupUseCase.refreshPhaseModels()` は credential missing の時点で backend 呼び出し前に止まる。
- 観測: `ProviderSettingsService.listProviderModelsWithTestSafeLoader()` は test-safe path でも endpoint 空なら `endpoint_missing` で止まる。
- 観測: `translation_job_setup_service.go` の direct path は test-safe の時点で credential gate より前に model loader を呼べる。

## 実装対象

- [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)

## 対象変更範囲

- backend の test-safe model list path では endpoint が空でも model list loader を呼ぶ。
- backend の test-safe model list path では credential / secret の不足や内容不正で停止しない。
- frontend の `refreshPhaseModels()` は credential missing だけを理由に backend 呼び出し前で停止しない。
- frontend は backend 応答の `credential_missing`、`endpoint_missing`、`success` など通常 status だけを表示判断に使う。
- model list が 1 件だけで現 selection の model が候補外なら、その 1 件を汎用的に選択してよい。

## 禁止変更範囲

- frontend へ fake mode 判定を追加しない。
- frontend へ `fake-model` 固有分岐を追加しない。
- `fake` provider ID を user-facing provider として追加しない。
- provider catalog、Job Setup provider list、Wails 公開 method、新規 DTO、DB schema は変更しない。
- プロダクトテストは変更しない。
- docs 正本本文は変更しない。

## 回帰確認観点

- fake mode の backend model list は endpoint 空、secret 不在、secret 不正でも成功して model を返す。
- real mode の backend model list は endpoint 空または credential 不足で従来通り停止する。
- frontend は credential missing 状態でも model list request を gateway へ渡す。
- frontend は backend が `credential_missing` を返した場合だけ未完了表示へ進む。
- frontend production code に `fake-model` 固有分岐が残らない。

## 検証コマンド

- `go test ./internal/service ./internal/bootstrap ./internal/controller/wails`
- `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/integration-implementation-result.md` に作成する。
- 変更ファイル、実装内容、検証結果、残留リスクを分けて記録する。
