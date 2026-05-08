# マスターペルソナ画面 layout fix

## 状態

- `task_id`: `2026-05-05-master-persona-layout-fix`
- `lane`: `fix-lane`
- `target`: マスターペルソナ生成画面
- `current_artifact`: `修正前調査`
- `source_task`: [2026-05-04-master-data-ux-refactor](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/plan.md)

## 人間観測

- 生成結果と詳細のレイアウトがプロトタイプと違う。
- 編集モーダルのレイアウトがプロトタイプと違う。
- 編集モーダルのフォントサイズもプロトタイプと違う。

## 成果物DAG

- `人間観測記録`: 完了
- `修正前調査`: 着手中
- `修正実行入力`: 未着手
- `実装証跡`: 未着手
- `回帰テスト証跡`: 未着手
- `レビュー通過根拠`: 未着手
- `作業レポート入力`: 未着手

## 参照

- [prototype/PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte)
- [prototype/PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte)
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
- [PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaActionModal.svelte)

## 境界

- `fix_lane` はプロダクトコードとプロダクトテストを直接変更しない。
- 実装が必要な場合は `implementation_implementer` に `implement-frontend` 固定で渡す。
- 修正対象は layout / typography の不一致に限定する。
