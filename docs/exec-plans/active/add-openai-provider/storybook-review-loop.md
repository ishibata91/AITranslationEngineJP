# Storybook レビューループ画面仕様

## 対象

- 確定した story: `OpenAiNoApiKey`、`OpenAiReady`、`OpenAiSubmitted`、`OpenAiBodyProcessing`、`OpenAiBodyReady`、`BodyReadyWithFailures`
- 確定した `fixture`: `OPENAI_NO_API_KEY_STATE`、`OPENAI_READY_STATE`、`OPENAI_SUBMITTED_STATE`、`OPENAI_BODY_PROCESSING_STATE`、`OPENAI_BODY_READY_STATE`
- 関連資源: `translation-run-view.ts`、`translation-run-presentation.ts`
- 作業中分類: `Review/Changed Screens/翻訳実行`、`Review/Changed Components/進行状況パネル`
- 通常分類: `Screens/翻訳実行`、`UI Components/進行状況パネル`
- 現在分類: `Screens/翻訳実行`、`UI Components/進行状況パネル`

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| 翻訳実行画面 | OpenAI（batch）を選択し、公式 endpoint、OpenAI API キー、モデル、二段階の進行状況を表示する。 | `TranslationRunScreen.svelte` | なし。 |
| batch 操作 | OpenAI API キーが空の場合は、送信、状態確認、取り込みを無効にする。 | `TranslationRunScreen.svelte` | なし。 |
| 進行状況 | OpenAI と xAI に共通する固有名、本文、完了の段と件数を表示する。 | `BatchProgressPanel.svelte` | なし。 |

## 現在状態

- 変更ファイル: `TranslationRunScreen.svelte`、`BatchProgressPanel.svelte`、関連する story、`fixture`、表示用の型と値。
- 通常分類へ戻した story: `TranslationRunScreen.stories.ts`、`BatchProgressPanel.stories.ts`
- 承認状態: 2026-08-01 に人間レビュー承認済み。
