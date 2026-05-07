# モデル設定カード controller 集約

## 状態

- `task_id`: `2026-05-07-model-settings-card-controller`
- `lane`: `implement-lane`
- `target`: モデル設定カードの保存、取得、モデル一覧取得、選択状態管理
- `current_artifact`: `修正後最終検証`
- `source`: 人間指示、直前の調査結果、軽量変更レーンの設計戻し判定

## task 枠

- 人間依頼: `.env` を先に用意し、fake mode で起動できる状態にする。
- 人間依頼: provider、model、モデル一覧、保存、取得、選択状態をモデル設定カード側へ集約する。
- 人間依頼: 集約対象はマスターペルソナと翻訳ジョブ設定の共有カード全体とする。
- 境界: `AIModelSelectionCard.svelte` は表示部品のまま維持し、専用 controller / usecase / store 層へ集約する。
- 禁止範囲: frontend に fake mode 判定や `fake-model` 固有分岐を追加しない。
- 確認したい結果: fake mode で通常 provider ID のままモデル一覧から `fake-model` を選べる。

## 成果物DAG

- `task 枠`: 完了
- `scenario_candidates`: 完了
- `シナリオ設計`: 完了
- `UI設計`: 完了
- `人間設計レビュー`: 完了
- `実装範囲`: 完了
- `実装引き継ぎ入力`: 完了
- `frontend 実装`: 完了
- `frontend 実装後人間レビュー`: 完了
- `backend 実装`: 完了
- `統合境界実装`: 完了
- `シナリオテスト`: 完了
- `単体テスト`: 完了
- `最終検証`: 完了
- `レビュー通過根拠`: 完了
- `レビュー指摘修正`: 完了
- `修正後最終検証`: 完了
- `正本化判断`: 完了
- `詳細仕様正本反映`: 停止
- `作業レポート入力`: 未着手
- `作業計画完了移動`: 未着手

## 進行判断

- モデル設定カード controller 集約は、新しい公開契約判断を含むため新規実装レーンで扱う。
- `Q-MSCC-001` から `Q-MSCC-004` までの人間回答は設計成果物へ反映済みである。
- `scenario-design.md` と `ui-design.md` は承認済みである。
- `frontend-implementation-result.md` は作成済みであり、`frontend 実装` は完了した。
- 人間が「フロントレビュー終わり」と回答したため、frontend 実装後人間レビューは承認済みである。
- `backend-implementation-result.md` は作成済みであり、`backend 実装` は完了した。
- 既存 backend core は承認済み backend 範囲を満たすため、backend プロダクトコード変更は不要と判定した。
- `integration-implementation-result.md` は作成済みであり、`統合境界実装` は完了した。
- 統合境界は既存 Wails / gateway 接続で承認済み範囲を満たすため、プロダクトコード変更は不要と判定した。
- `scenario-test-implementation-result.md` と `unit-test-implementation-result.md` は作成済みであり、`wave-4` の並列テストは完了した。
- `SCN-MSCC-003` は API 受け入れテストで証明済みである。
- provider 切替時の旧 model list / model 混入禁止、遅延応答破棄、保存拒否は frontend usecase 単体テストで証明済みである。
- `final-validation-result.md` は作成済みであり、最終検証は完了した。
- `python3 scripts/harness/run.py --suite all` は通過した。
- 5 観点レビューは完了した。
- `behavior-001`、`contract-001`、`trust-boundary-001` は修正必須である。
- `trust-boundary-001` は hard gate の blocker であるため、`implementation_action` は `fix` とする。
- `review-fix-implementation-result.md` は作成済みであり、レビュー指摘修正は完了した。
- 修正後の最終検証は完了した。
- 修正後の `python3 scripts/harness/run.py --suite coverage` は通過し、Sonar reliability issue は 0 件である。
- 修正後レビュー再実行は完了した。
- 5 観点の `must_fix_open` はすべて false である。
- `implementation_action` は `close` とする。
- 詳細仕様正本反映は `docs/index.md` の規約により、human が直接起動した `updating-docs` へ委ねる。
- 作業レポート入力へ進む。
- backend 実装は `wave-2`、統合境界実装は `wave-3`、テストは `wave-4` の順で扱う。
