# マスターペルソナ モデル一覧更新時の provider reset 修正

## 状態

- `task_id`: `2026-05-07-master-persona-model-refresh-provider-reset`
- `lane`: `fix-lane`
- `target`: マスターペルソナ画面の AI 設定カード
- `current_artifact`: `作業レポート入力`
- `source_task`: [2026-05-07-fake-fixed-model-closed-path](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md)

## 人間観測

- マスターペルソナ画面でモデル一覧を更新しても、モデルを選択できない。
- モデル一覧更新後に、AI サービスのプルダウン選択状態がリセットされる。
- 失敗は provider 値が意図せず戻る挙動として観測された。

## 成果物DAG

- `人間観測記録`: 完了
- `修正前調査`: 完了
- `修正実行入力`: 完了
- `実装証跡`: 完了
- `回帰テスト証跡`: 完了
- `レビュー通過根拠`: 完了
- `作業レポート入力`: 完了
- `作業計画完了移動`: 未着手

## 参照

- [source plan](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md)
- [light-change-planning.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/light-change-planning.md)
- [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)
- [GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)

## 境界

- `fix_lane` はプロダクトコードとプロダクトテストを直接変更しない。
- 修正対象はマスターペルソナ画面の provider 選択状態保持と model 選択可能化に限定する。
- `fake` provider ID を user-facing provider として追加しない。
- 実装が必要な場合は `implementation_implementer` に `implement-frontend` または `implement-integration` のどちらか 1 つへ固定して渡す。
