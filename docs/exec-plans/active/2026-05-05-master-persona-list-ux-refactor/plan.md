# マスターペルソナ画面 一覧UX調整

## 状態

- `task_id`: `2026-05-05-master-persona-list-ux-refactor`
- `lane`: `ux-refactor-lane`
- `target`: マスターペルソナ画面
- `current_artifact`: `作業レポート入力`
- `source`: 人間依頼

## 依頼要約

- 上側ヒーローへ下側ヒーローのページ説明文だけを移す。
- 下側ヒーローは削除する。
- ペルソナ一覧の行の高さ、余白、名前の文字色をダークモードに合わせる。
- 検索窓とプラグインフィルタの高さと幅の扱いを揃える。

## 成果物DAG

- `task 枠`: 完了
- `frontend 実装`: 停止
- `人間UIレビュー`: 差し戻し
- `レビュー通過根拠`: 未着手
- `作業レポート入力`: 未着手
- `作業計画完了移動`: 未着手

## 停止理由

- `viewModel.runStatus.runState` は軽量な状態表示として復活済みである。
- 詳細カードの不要文言は可視表示と hidden DOM から削除済みである。
- `python3 scripts/harness/run.py --suite frontend-local` は既存 `frontend/src/ui/App.test.ts` の期待値により失敗した。
- UX改善レーンの禁止範囲により、プロダクトテスト変更は実施しない。

## 境界

- 変更許可範囲は `frontend/src/ui/screens/master-persona/` と `frontend/src/ui/stores/shell-state.ts` の表示文言、表示構造、CSSに限定する。
- backend、統合境界、保存先、ログ出力、secret、外部入力は変更しない。
- プロダクトテストは変更しない。
- 表示項目や機能は増やさない。
