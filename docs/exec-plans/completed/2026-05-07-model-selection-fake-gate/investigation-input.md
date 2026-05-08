# 修正前調査起動入力

## 起動先

- `agent`: `investigator`
- `skill`: `investigate`
- `investigation_mode`: `修正前調査`

## investigation_goal

fake mode の model list 取得が、frontend の credential preflight または backend の endpoint / credential preflight によって止まる経路を観測事実として整理する。

## known_context

- 人間観測は [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation.md) を正本にする。
- fake mode の正本は backend の test-safe loader / transport である。
- frontend は `fake-model` 文字列を特別扱いしない。

## candidate_paths

- [provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go)
- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go)
- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- [translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)

## 参照成果物

- [fake-fixed-model plan](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md)
- [fake-fixed-model implementation result](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)
- [model refresh fix plan](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/plan.md)

## 禁止事項

- プロダクトコードを変更しない。
- プロダクトテストを変更しない。
- 修正実行入力を作らない。
- frontend fake 判定や `fake-model` 固有分岐を提案しない。

## 期待する成果物

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/investigation.md` に修正前調査を作成する。
- 観測事実、根拠 path、原因候補、影響ファイル候補、残り不足、推奨 next step を分ける。
