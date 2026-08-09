# Storybook レビューループ画面仕様

## 対象

- 確定した story: `UI Components/TranslationResultRow/展開（本文・複数参考語）` と `Screens/翻訳実行/完了`。
- 確定した `fixture`: `DIALOGUE_ROW`、`NARRATION_ROW`、`RESULT_ROWS`、`LINE_RESULT_ROWS`。
- 関連資源: `translation-run-view.ts`、`TranslationResultRow.svelte`、`translation-run.fixtures.ts`。
- 作業中分類: `Review/Changed Components/TranslationResultRow`、`Review/Changed Screens/翻訳実行`。
- 通常分類: `UI Components/TranslationResultRow`、`Screens/翻訳実行`。
- 現在分類: 通常分類。

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| 翻訳結果行 | 本文の参考語を原語ごとにまとめ、同じ原語の複数候補へ訳語、品詞、Skyrimカテゴリ、出どころを表示する。`meaning`は表示しない。本文は置換しない。参考語がない結果には節を表示しない。 | `TranslationResultRow.svelte`、`translation-run-view.ts`、story、fixture | 人間の見た目承認待ち。 |
| 翻訳実行画面 | 畳んだ結果行へ参考語の原語数を表示する。 | `TranslationRunScreen.stories.ts`、`translation-run.fixtures.ts` | 人間の見た目承認待ち。 |

## 現在状態

- 変更ファイル: `frontend/src/ui/screens/translation-run/TranslationResultRow.svelte`、`translation-run-view.ts`、`translation-run.fixtures.ts`、`TranslationResultRow.stories.ts`、`TranslationRunScreen.stories.ts`。
- 通常分類へ戻した story: `UI Components/TranslationResultRow`、`Screens/翻訳実行`。
- 承認状態: 承認済み（2026-08-09）。
- 検証: `npm run test:frontend`は80件成功。`npm --prefix frontend run build-storybook`は成功。`npm run lint:frontend`は既存の`react-syntax-highlighter/dist/esm/prism-light`型定義不足で停止した。
