# Storybook レビューループ画面仕様

## 対象

- 確定した story: `TranslationRunScreen.stories.ts` の `NotStarted`、`Running`、`Paused`、`Done`、`DoneWithUntranslated`、`Failed`。`BatchProgressPanel.stories.ts` の `NotStarted`、`Starting`、`ProperRunning`、`ProperReady`、`Paused`、`Done`。
- 確定した `fixture`: `translation-run.fixtures.ts` の OpenAI batch 表示用 fixture と `batchRunning`。
- 関連資源: `TranslationRunScreen.svelte`、`BatchProgressPanel.svelte`、`translation-run-presentation.ts`、`frontend/.storybook/main.ts`、`frontend/package.json`、`frontend/package-lock.json`。
- 起動 command: `npm --prefix frontend run storybook`
- 固定 URL: `http://localhost:6008/`
- 作業中分類: `Review/Changed Screens/翻訳実行`、`Review/Changed Components/進行状況パネル`
- 通常分類: `Screens/翻訳実行`、`UI Components/進行状況パネル`
- 現在分類: `Screens/翻訳実行`、`UI Components/進行状況パネル`。人間レビュー承認後に通常分類へ復帰済み。

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| OpenAI・xAI の操作 | 状態確認ボタンと手動取り込みボタンを削除し、batch の主操作を一つだけ表示する | `TranslationRunScreen.svelte` | なし |
| batch の主操作 | 未開始、実行中、再開待ち、未訳再送信、完了をラベルと活性で表示する | `translation-run-presentation.ts`、`TranslationRunScreen.svelte` | なし |
| 進行状況 | 処理中、再開待ち、完了を案内し、人の状態確認・取り込みを促さない | `BatchProgressPanel.svelte`、`translation-run-presentation.ts` | なし |
| story metadata | 各 story の `parameters.docs.description.story` に前提条件と期待値を書き、Docs 画面に表示する | `TranslationRunScreen.stories.ts`、`BatchProgressPanel.stories.ts`、`frontend/.storybook/main.ts` | なし |

## 表示実装とロジック実装の境界

- 本成果物は svelte の表示構造、表示用 props、表示文言、story、fixture を固定する。
- `batchRunning` は、開始処理と自動状態確認中の表示入力である。
- 自動状態確認の予約、gateway 呼び出し、画面終了時の停止、遅延応答の破棄は `implementation-module` が実装する。

## 現在状態

- 変更ファイル: `TranslationRunScreen.svelte`、`BatchProgressPanel.svelte`、`translation-run-presentation.ts`、`TranslationRunScreen.stories.ts`、`BatchProgressPanel.stories.ts`、`translation-run.fixtures.ts`、`frontend/.storybook/main.ts`、`frontend/package.json`、`frontend/package-lock.json`。
- 通常分類へ戻した story: `Screens/翻訳実行`、`UI Components/進行状況パネル`。
- 承認状態: 2026-08-02 に人間レビュー承認済み。文言差し戻し2件を反映し、状態確認ボタンがない画面に合わせて内部処理の説明を外し、完了まで自動で進む結果だけを案内する形で確定した。
