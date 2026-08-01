# Storybook レビューループ記録: batch-retry-untranslated-records

確定結果だけを記録する。人間コメントの履歴は残さない。

## 起動

- 固定 URL: `http://localhost:6008/`
- 起動 command: `npm --prefix frontend run storybook`
- build 確認 command: `npm --prefix frontend run build-storybook`

## 対象 story

| 確認資源 | story | 状態名 |
| --- | --- | --- |
| 変更画面 | `Screens/翻訳実行/OpenAI・未訳が残る` | batch 完了後に未訳件数と再送信表示がある |
| 変更画面 | `Screens/翻訳実行/OpenAI・未訳なし` | batch 完了後に未訳表示がない |
| 変更画面 | `Screens/翻訳実行/xAI・未訳が残る` | batch 完了後に未訳件数と再送信表示がある |
| 変更画面 | `Screens/翻訳実行/xAI・未訳のみ` | 未訳だけに絞った結果一覧 |
| 変更コンポーネント | `UI Components/ResultsPanel/未訳のみ` | チェックボックス選択中 |
| 変更コンポーネント | `UI Components/ResultsPanel/未訳なし（書き出し可能）` | 未訳 0 件でも書き出し操作がある |

## 分類

- 作業中分類: 既存 story ファイルへ状態を追加するため、通常分類のまま確認する。
- 通常分類: `Screens/翻訳実行`、`UI Components/ResultsPanel`。
- 現在分類: 人間承認済み。通常分類。

## fixture / 関連資源

- fixture: `translation-run.fixtures.ts` の `OPENAI_BATCH_UNTRANSLATED_STATE`、`OPENAI_NO_UNTRANSLATED_STATE`、`XAI_BATCH_UNTRANSLATED_STATE`、`XAI_UNTRANSLATED_ONLY_STATE`。
- 表示用型: `translation-run-view.ts` の `BatchProgressView.untranslatedCount`。
- 表示規則: `translation-run-presentation.ts` の batch 未訳案内と主操作表示。

## frontend 表示実装境界

- 扱った範囲: `ResultsPanel.svelte` と `TranslationRunScreen.svelte` の props、template、表示文言、story、fixture、表示用型、表示用純関数。
- `implementation-module` へ渡す範囲: backend の未訳件数、batch 再送信、Wails binding、gateway、Container のチェック状態とページ取得、取得失敗時の状態維持。

## 画面表示の根拠

- `design.md` の R-1 と R-3。
- `spec.md` の R-1-1〜R-1-4、R-3-1〜R-3-5。

## 確定結果

- 承認状態: 人間承認済み。
- 確定した変更画面仕様: OpenAI と xAI の batch 完了後に未訳件数と「未訳だけを再送信」を表示する。結果一覧に「未訳のみ」のチェックボックスを置く。選択時は未訳だけを表示する。未訳 0 件では「未訳はありません」と表示し、xTranslator への書き出し操作は残す。
- 現在分類: 通常分類。
- 反映先: `ResultsPanel.svelte`、`TranslationRunScreen.svelte`、`translation-run-view.ts`、`translation-run-presentation.ts`、`translation-run.fixtures.ts`、`ResultsPanel.stories.ts`、`TranslationRunScreen.stories.ts`、`BatchProgressPanel.stories.ts`。
