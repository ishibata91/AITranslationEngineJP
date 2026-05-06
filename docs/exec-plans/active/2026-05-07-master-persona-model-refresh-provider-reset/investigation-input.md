# 修正前調査起動入力

## 起動先

- `agent`: `investigator`
- `skill`: `investigate`
- `investigation_mode`: `修正前調査`

## investigation_goal

マスターペルソナ画面で、モデル一覧更新後に AI サービス選択状態がリセットされ、モデルを選択できない原因候補を観測事実として整理する。

## known_context

- 人間観測は [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/human-observation.md) を正本にする。
- `2026-05-07-fake-fixed-model-closed-path` は、fake mode で通常 provider ID のまま `fake-model` を扱う変更である。
- 同 task の残留リスクは、マスターペルソナのモデルカード実表示が未確認である点だった。

## candidate_paths

- [GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)
- [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
- [master-persona.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts)
- [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)

## 参照成果物

- [source plan](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md)
- [light-change-planning.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/light-change-planning.md)
- [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)
- [review-summary.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/review-summary.md)

## 禁止事項

- プロダクトコードを変更しない。
- プロダクトテストを変更しない。
- 修正実行入力を作らない。
- 実装 skill を確定しない。
- `fake` provider ID を user-facing provider に追加する方針を置かない。

## 期待する成果物

- `修正前調査`: 観測事実、根拠 path、UI 証跡または未確認理由、原因候補、影響ファイル候補、残り不足を分けて返す。
- `推奨 next step`: 修正実行入力へ進めるか、追加調査で止めるかを返す。
