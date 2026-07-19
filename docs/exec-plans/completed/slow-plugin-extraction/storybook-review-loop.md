# Storybook レビューループ: 翻訳前区間の進捗サブ段表示

## 確定した story

`frontend/src/ui/screens/translation-run/TranslationProgress.stories.ts`。抽出段の4サブ段（台詞抽出・辞書準備・既存訳取込・取込段）、サブ段未指定フォールバック、翻訳中2件。

## 変更された画面仕様

実行中の「抽出」段（不定バー）の見出しを、内部の4サブ段に合わせて出し分ける。見出しだけがサブ段ごとに切り替わり、不定バーは据え置き。翻訳段の done/total 確定バーは不変。

サブ段と見出し:

| サブ段（backend の処理） | 見出しラベル |
| --- | --- |
| 台詞抽出（C# 抽出子） | 台詞を抽出しています |
| 辞書準備（DeriveMasterTerms） | 辞書を準備しています |
| 既存訳取込（LoadReferenceTranslations） | 既存訳を取り込んでいます |
| 取込段（Ingest） | 翻訳対象を仕分けています |

人間レビューでの修正: 当初「固有名を派生しています」は内部用語で意味が通らないため「辞書を準備しています」へ変更（承認済み）。

## 反映先

- `frontend/src/ui/screens/translation-run/translation-run-view.ts`: `ExtractStep` union と `RunProgress.step?`。
- `frontend/src/ui/screens/translation-run/translation-run-presentation.ts`: `EXTRACT_STEP_LABEL`。
- `frontend/src/ui/screens/translation-run/TranslationProgress.svelte`: 見出し導出（extract 段でサブ段があればサブ段ラベル）。
- `frontend/src/ui/screens/translation-run/TranslationProgress.stories.ts`: story。

## 現在分類

通常分類 `UI Components/TranslationProgress` へ復帰済み（レビュー中は `Review/Changed Components/TranslationProgress`）。

## 承認状態

承認済み。検証: `npm run test:frontend`（2 passed）、`npm --prefix frontend run build-storybook`（成功）。

## 表示範囲外の残課題（implementation-module へ渡す）

- backend Go `ProgressEvent` に段内ラベル（`step`）を足す。
- `RunExtractAndTranslate` の4境界（C#抽出/辞書準備/既存訳取込/取込段の完了）で `step` 付き進捗を emit。
- `TranslationRunContainer.svelte` の進捗購読で `step` を `RunProgress` へ流し込む。
