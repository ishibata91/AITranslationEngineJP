# UX改善レーン plan

- `task_id`: `ux-dashboard-refactor-20260504`
- `workflow`: `ux-refactor-lane`
- `status`: `completed`
- `lane_owner`: `ux_refactor_lane`
- `target_screen`: `ダッシュボード`
- `request`: ダッシュボード画面を対象にリファクタする。既存画面を確認し、全項目が維持されていることを保証する。

## 状態

- `task 枠`: 完了
- `UI改善契約`: 完了
- `人間UIレビュー`: 完了
- `UX実装修正入力`: 完了
- `frontend 実装`: 完了
- `実装後単体テスト`: 完了
- `実装後確認`: 完了
- `レビュー通過根拠`: 完了
- `作業レポート入力`: 完了
- `作業計画完了移動`: 完了

## 影響範囲

- 対象画面は `AppShell` の `#dashboard` 表示に限定する。
- UI改善契約は既存仕様、既存画面、実物確認結果を根拠にする。
- 既存表示項目、導線、状態表示、禁止表示を維持対象として固定する。
- プロダクトコードとプロダクトテストは `ux_refactor_lane` では変更しない。

## 実行計画

- 既存画面根拠を `existing-screen-evidence.md` に固定する。
- `designer` に文脈継承なしで `UI改善契約` を作らせる。
- 人間UIレビュー後に、承認済み `UI改善契約` だけを根拠に実装修正入力を作る。
- 実装後は実画面または検証コマンドで、維持対象の欠落がないことを確認する。

## 検証方法

- `agent-browser open http://127.0.0.1:34115/#dashboard`
- `agent-browser snapshot`
- `agent-browser screenshot tmp/agent-browser/ux-dashboard-refactor-before.png`
- 実装後に `npm --prefix frontend run test -- AppShell` を候補にする。

## HITL

- `人間UIレビュー`: 承認
- `人間UIレビュー記録`: `準備中` の文言変更を追加すれば承認。
- `review_target`: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/ux-dashboard-refactor-20260504/ui-design.md)

## 完了移動

- `completed_path`: [docs/exec-plans/completed/ux-dashboard-refactor-20260504](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/ux-dashboard-refactor-20260504)
- `work_report`: [work_history/runs/2026-05-04-ux-dashboard-refactor-20260504-run/README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-04-ux-dashboard-refactor-20260504-run/README.md)
