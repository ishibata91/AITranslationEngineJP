# 修正前調査

## 判定

- 判定: 完了
- 調査 mode: `修正前調査`
- 対象 task: `2026-05-07-model-selection-fake-gate`
- 推奨引き継ぎ先: `fix_lane`

## 人間観測

- 人間観測では、AIサービス設定が設定済み表示でも、「AIサービス設定が未完了です。」へ入る経路があると記録されている。根拠: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation.md)
- 人間観測では、fake mode の model list 取得は backend の test-safe loader / transport で完結し、frontend は fake mode や `fake-model` を特別扱いしない前提である。根拠: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation.md), [investigation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/investigation-input.md)

## 観測事実

- 観測: fake mode の model list loader は `ProviderModelListRequestsAreTestSafe()` を返す transport で wiring される。translation job setup と provider settings の両 adapter は、その test-safe 判定を service へ伝搬する。根拠: [app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go), [provider_models.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider_models.go)
internal/bootstrap/app_controller.go:333
internal/bootstrap/app_controller.go:359
internal/bootstrap/app_controller.go:452
internal/bootstrap/app_controller.go:455
internal/infra/ai/provider_models.go:119

- 観測: translation job setup frontend の `refreshPhaseModels()` は、provider が credential 必須で `selection.credentialStatus !== "configured"` の時点で backend 呼び出しを行わず、画面 state を `status: "credential_missing"` に更新して終了する。根拠: [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:829
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:841
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:846

- 観測: frontend の `selection.credentialStatus` は、provider capability と credential reference の有無から計算される。test-safe 判定は考慮されない。credential 必須 provider で usable な credential reference を解決できない時は `missing` になる。根拠: [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:248
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:256
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:281
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:304

- 観測: translation job setup presenter は、`selection.credentialStatus === "missing"` の時点で「AIサービス設定が未完了です。設定が必要です。」を返す。更新ボタンも `selection?.credentialStatus !== "missing"` の時だけ有効化する。根拠: [translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts), [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:198
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:231
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:318
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:325
frontend/src/ui/components/AIModelSelectionCard.svelte:150
frontend/src/ui/components/AIModelSelectionCard.svelte:178

- 観測: translation job setup backend は、`providerModelListRequestsAreTestSafe()` が true の時、direct path では credential gate より先に model loader を呼び、`CredentialStatus = "not_required"` で成功を返す。根拠: [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
internal/service/translation_job_setup_service.go:610
internal/service/translation_job_setup_service.go:612
internal/service/translation_job_setup_service.go:620
internal/service/translation_job_setup_service.go:623

- 観測: translation job setup backend が provider settings 経由の path を使う時は、frontend から受けた `CredentialStatus` ではなく、provider settings summary の `CredentialState` を result へ上書きして返す。根拠: [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
internal/service/translation_job_setup_service.go:557
internal/service/translation_job_setup_service.go:568
internal/service/translation_job_setup_service.go:583
internal/service/translation_job_setup_service.go:595

- 観測: provider settings backend は、test-safe path でも `summary.Endpoint` が空なら `endpoint_missing` を返して止まる。非 test-safe path では snapshot 不一致と endpoint 不足を preflight で先に判定する。根拠: [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
internal/service/provider_settings_service.go:669
internal/service/provider_settings_service.go:672
internal/service/provider_settings_service.go:699
internal/service/provider_settings_service.go:702
internal/service/provider_settings_service.go:733
internal/service/provider_settings_service.go:744
internal/service/provider_settings_service.go:747

- 観測: provider settings summary は user-facing provider catalog の default endpoint を優先して持つ。現在の catalog では `gemini`、`lm_studio`、`xai` に default endpoint がある。根拠: [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
internal/service/provider_settings_service.go:67
internal/service/provider_settings_service.go:68
internal/service/provider_settings_service.go:69
internal/service/provider_settings_service.go:1110
internal/service/provider_settings_service.go:1115

- 観測: provider settings frontend の接続確認は、credential 必須 provider で `credentialState !== "configured"` の時点、または endpoint 未入力の時点で backend 呼び出し前に止まる。根拠: [provider-settings.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts)
frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:331
frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:338
frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:350
frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:357
frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:366

- 観測: fake-fixed-model 実装結果では、「test-safe loader の場合は外部 secret を要求しないようにした」と記録されている。一方で現行 frontend には test-safe 判定を使った分岐はなく、credential 状態だけで停止する。根拠: [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md), [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md:11
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:829

## 根拠 path

- [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation.md)
- [investigation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/investigation-input.md)
- [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)
- [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
- [app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go)
- [provider_models.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider_models.go)
- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- [translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts)
- [provider-settings.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)

## 原因候補

- 候補1: translation job setup の停止点は frontend credential preflight である可能性が高い。backend には test-safe path で credential gate を外す分岐があるが、frontend が backend 呼び出し前に `credential_missing` を確定している。根拠: [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts), [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:829
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:841
internal/service/translation_job_setup_service.go:612
internal/service/translation_job_setup_service.go:620

- 候補2: frontend の credential 状態計算が、provider capability と credential reference だけに依存し、fake mode の test-safe 実行可否を知らないため、fake mode でも required provider を `missing` と表示し続ける可能性が高い。根拠: [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts), [translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts)
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:248
frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:256
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:198
frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:231

- 候補3: provider settings 側には、test-safe path でも endpoint 判定が残っている。現在 catalog 既定値と summary 生成では通常は埋まるが、実行時 summary が空 endpoint になる条件がある場合は backend 側で停止する可能性がある。根拠: [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
internal/service/provider_settings_service.go:699
internal/service/provider_settings_service.go:702
internal/service/provider_settings_service.go:1110
internal/service/provider_settings_service.go:1115

## 影響ファイル候補

- [frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- [frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts)
- [frontend/src/ui/components/AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)
- [internal/service/translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
- [internal/service/provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts)

## 残り不足

- 不足: 人間観測にある「AIサービス設定が設定済み表示でも未完了へ入る」事象を、現行ローカル UI で再現した証跡はまだない。
- 不足: `providerSettings.ListProviderModels()` が fake mode 実行時に live でどの state を返しているかの実レスポンス記録がない。
- 不足: provider settings summary の endpoint が live 実行時に default endpoint 由来か、保存済み値由来かの記録がない。
- 不足: translation job setup と AIサービス設定のどちらが、人間観測の主経路だったかを画面操作単位では固定できていない。
- 不足: console、Wails binding、backend log の追加証跡は未取得である。

## 推奨 next step

- 推奨: `fix_lane` は、修正実行入力の前に translation job setup 画面で model refresh 操作を 1 回だけ live 観測し、backend 呼び出し前停止かどうかを UI 証跡で固定する。
- 推奨: 同じ turn で `ListTranslationJobSetupProviderModels` または provider settings model list の実レスポンスを 1 回だけ採取し、`credential_missing` と `endpoint_missing` のどちらが live path に出ているかを分けて記録する。
- 推奨: 修正実行入力では、frontend fake 判定や `fake-model` 固有分岐を前提にせず、観測済み停止点だけを narrow して渡す。
