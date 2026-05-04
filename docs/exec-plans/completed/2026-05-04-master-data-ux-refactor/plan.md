# マスターデータ画面 UX 改善

## 状態

- `task_id`: `2026-05-04-master-data-ux-refactor`
- `lane`: `ux-refactor-lane`
- `target_screen`: マスターペルソナ生成画面
- `current_artifact`: `完了`
- `human_review`: 承認済み

## 依頼要約

マスターペルソナ生成画面を、既存 UI に強く縛られずに作り直す。
既存 UI の表示項目は、利用者視点で不要なら削ってよい。
新規表示項目の追加は、人間確認を必要とする。

モデル選択カードは例外として既存の `AIModelSelectionCard.svelte` を維持する。
task-local UIプロトタイプでも同じ部品を import して使う。

## 成果物DAG

- `task 枠`: 完了
- `UI改善契約`: 完了
- `人間UIレビュー`: 承認済み
- `UX実装修正入力`: 完了
- `frontend 実装`: 完了
- `実装後単体テスト`: 完了
- `実装後確認`: 完了
- `レビュー通過根拠`: 完了
- `作業レポート入力`: 完了
- `作業計画完了移動`: 完了

## 作業レポート

- `run_folder`: [2026-05-05-2026-05-04-master-data-ux-refactor-run](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-2026-05-04-master-data-ux-refactor-run)
- `run_summary`: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-2026-05-04-master-data-ux-refactor-run/README.md)
- `codex_report`: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-2026-05-04-master-data-ux-refactor-run/codex.md)

## 実行境界

- `ux_refactor_lane` はプロダクトコードとプロダクトテストを直接変更しない。
- UI改善契約は `ui-design` に従い、既存画面根拠と既存仕様を根拠にする。
- `implementation_implementer` へ進む条件は、人間UIレビュー承認済みの `UI改善契約` が存在すること。
- docs 正本本文は変更しない。

## 既存画面根拠

- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)
- [master-persona-gateway-contract.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/gateway-contract/master-persona/master-persona-gateway-contract.ts)
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md)
- [UX-standard.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/UX-standard.md)

## 検証予定

- UIプロトタイプ確認サーバー: `npm --prefix frontend run dev:prototype -- --task 2026-05-04-master-data-ux-refactor --port 34118`
- 確認 URL: `http://127.0.0.1:34118/prototype`
- 確認観点: 目的、主要操作、状態表示、表示項目の削減、文言、レスポンシブ、モデル選択カードの維持

## 人間UIレビュー

- `review_status`: 承認済み
- `designer_agent_id`: `019df3c4-6e0c-74b3-9304-def64a25eda1`
- `feedback_route`: `designer` へ戻す
- `review_result`: 承認
- `review_note`: 追加表示項目なしで frontend 実装へ進める
- `server_status`: レビュー終了
- `server_note`: 人間UIレビュー終了後に `designer` agent を閉じる
