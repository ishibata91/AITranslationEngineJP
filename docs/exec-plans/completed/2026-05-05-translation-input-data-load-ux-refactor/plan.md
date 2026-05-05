# 翻訳入力データロード画面 UX改善

## 状態

- `task_id`: `2026-05-05-translation-input-data-load-ux-refactor`
- `lane`: `ux-refactor-lane`
- `target`: 翻訳管理の入力データ画面
- `current_artifact`: `completed`
- `source`: 人間依頼

## 依頼要約

- 翻訳管理の `Input Review` 画面を作り直す。
- 画面概念を `input review` ではなくデータロード画面へ変える。
- UX観点を適用し、作業順、状態、結果、再構築判断を見やすくする。
- コンポーネント化を考慮し、画面専用部品へ分ける。

## 成果物DAG

- `task 枠`: 完了
- `frontend 実装`: 完了
- `人間UIレビュー`: 承認
- `テスト修正証跡`: 完了
- `レビュー通過根拠`: 完了
- `作業レポート入力`: 完了
- `作業計画完了移動`: 完了

## 境界

- 変更許可範囲は `frontend/src/ui/screens/translation-input/`、`frontend/src/ui/stores/shell-state.ts`、`frontend/src/ui/views/AppShell.svelte` の翻訳管理セクション表示文言、表示構造、CSS、画面専用部品に限定する。
- backend、Wails gateway、controller、usecase、store、gateway contract は変更しない。
- 保存先、ログ出力、secret、外部入力の扱いは変更しない。
- プロダクトテスト、検証データ、snapshot、test helper は原則変更しない。
- 人間UIレビュー後の `テスト修正証跡` に限り、旧 `Input Review` 表示名を期待する frontend 単体テストだけを変更してよい。

## 検証

- frontend agent は `python3 scripts/harness/run.py --suite frontend-local` を実行する。
- 失敗原因がプロダクトテスト更新を必要とする場合は、UX改善レーンの禁止範囲として停止理由に残す。
- 実画面確認は `agent-browser open http://localhost:34115` を基本入口にする。

## 実装結果

- 画面 source と実 DOM から `Input Review` 見出しを削除し、`翻訳入力データロード` へ置き換えた。
- 翻訳管理セクション見出しとタブ表示を `データロード` へ置き換えた。
- 画面専用 component を `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail` に分割した。
- `agent-browser snapshot` で `翻訳入力データロード`、`ロード準備`、`読み込み済みデータ`、`選択データの内容` を確認した。

## テスト修正対象

- `python3 scripts/harness/run.py --suite frontend-local` は旧 `Input Review` 期待で失敗した。
- 対象は [InputReviewPage.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.test.ts) と [translation-input-app.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/translation-input-app.test.ts) に限定する。
- テスト修正は UI 表示文言と構造変更への追従だけを扱う。

## レビュー通過根拠

- [reviewback.responsibility-boundary.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-translation-input-data-load-ux-refactor/reviewback.responsibility-boundary.yaml)
- `review_status`: `no_issue`
- `must_fix_open`: `false`
- `max_level`: `none`

## 作業レポート入力

- [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-translation-input-data-load-ux-refactor-run/README.md)
- [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-translation-input-data-load-ux-refactor-run/codex.md)
- [transcript_refs.json](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/runs/2026-05-05-translation-input-data-load-ux-refactor-run/transcript_refs.json)

## 作業計画完了移動

- `docs/exec-plans/completed/2026-05-05-translation-input-data-load-ux-refactor/` へ移動済み。
