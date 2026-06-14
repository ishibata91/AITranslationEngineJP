# Storybook レビューループ記録（persona-dictionary-pipeline）

確定結果のみを記録する。人間コメントの履歴は残さない。

## 起動

- 固定 URL: `http://localhost:6008/`
- 起動 command: `npm --prefix frontend run storybook`
- build 確認 command: `npm --prefix frontend run build-storybook`

## 対象 story（作業中分類）

| 確認資源 | 作業中分類 story id | 状態名 |
| --- | --- | --- |
| 追加コンポーネント | `review-changed-components-translationprogress--extracting` | 抽出中（不定） |
| 追加コンポーネント | `review-changed-components-translationprogress--translating-start` | 翻訳中（開始直後） |
| 追加コンポーネント | `review-changed-components-translationprogress--translating-mid` | 翻訳中（途中） |
| 追加コンポーネント | `review-changed-components-translationprogress--translating-near-done` | 翻訳中（ほぼ完了） |
| 変更コンポーネント | `review-changed-components-translationresultrow--with-persona-child` | 口調あり（子供） |
| 変更コンポーネント | `review-changed-components-translationresultrow--with-persona-old-woman` | 口調あり（老女） |
| 変更コンポーネント | `review-changed-components-translationresultrow--without-persona` | 口調なし |
| 変更画面 | `review-changed-screens-翻訳実行--running-extract` | 実行中（抽出） |
| 変更画面 | `review-changed-screens-翻訳実行--running-translate` | 実行中（本文翻訳） |
| 変更画面 | `review-changed-screens-翻訳実行--done-persona` | 完了（口調差） |

## 通常分類（承認後の戻し先）

- 追加コンポーネント `TranslationProgress` → `UI Components/TranslationProgress/<状態名>`
- 変更コンポーネント `TranslationResultRow` の口調状態 → `UI Components/TranslationResultRow/<状態名>`（既存ファイルへ統合）
- 変更画面 `翻訳実行` の進捗・口調状態 → `Screens/翻訳実行/<状態名>`（既存ファイルへ統合）

## fixture / 関連資源

- fixture: `translation-run.fixtures.ts`（`RUNNING_EXTRACT_STATE`、`RUNNING_TRANSLATE_STATE`、`DONE_PERSONA_STATE`、`LINE_RESULT_ROWS`）。
- 関連資源: `translation-run-view.ts`（`RunProgress`、`ProgressStage`、`NarrationResultRow.directive`）、`translation-run-presentation.ts`（`PROGRESS_STAGE_LABEL`）。

## frontend 表示実装境界（storybook-module で扱う範囲）

- 扱う: `TranslationProgress.svelte`（進捗バー）、`TranslationResultRow.svelte` の口調指示文併記、`TranslationRunScreen.svelte` の進捗バー配置、story、fixture、表示用 view 型・presentation 値。
- 扱わない（implementation-module へ）: Wails runtime events 購読、gateway の event 経路、container の進捗 state と結果 state、`RunExtractAndTranslate` 呼び出し。

## 画面表示の根拠

- `design-diff.md`（進捗バーと口調併記の追加箇所）、`plan.md` の実装範囲（進捗 event 契約 `{ phase, done, total }`、結果行の口調指示文）。

## 確定結果

- 承認状態: 承認済み。
- 確定した変更画面仕様: 結果行をカード型からコンパクト 1 行（`details/summary`）へ変更し、口調指示は既定で畳んで口調チップ（声質などの短い要約）だけ一覧に出し、行クリックで全文展開。本文翻訳の進捗バーを追加（抽出＝不定／本文翻訳＝done/total）。
- 現在分類: 全 story を通常分類へ復帰済み。`UI Components/TranslationProgress`、`UI Components/TranslationResultRow`、`Screens/翻訳実行`。作業中分類（`Review/...`）と review 一時 story ファイルは削除済み。
- 反映先: `TranslationProgress.svelte`（追加）、`TranslationResultRow.svelte`・`TranslationRunScreen.svelte`・`ResultsPanel.svelte`（変更）、`translation-run-view.ts`・`translation-run-presentation.ts`・`translation-run.fixtures.ts`。
- 検証: `build-storybook` 成功、`frontend-local` suite 通過、`svelte-check` は本 task 変更にエラーなし。
