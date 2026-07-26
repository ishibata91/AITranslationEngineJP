# Storybook レビューループ画面仕様

## 対象

- 確定した story: `Screens/翻訳実行 → 完了（未訳が残る）`（`DoneUntranslated`）。退行確認に `完了`、`xAI・送信後` を併せて見た。
- 確定した `fixture`: `frontend/src/ui/screens/translation-run/translation-run.fixtures.ts` の `DONE_UNTRANSLATED_STATE`。
- 関連資源: `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`、`translation-run-presentation.ts` の案内文を組む純粋関数。Storybook は `npm --prefix frontend run storybook` で `http://localhost:6008/` に固定して起動した。
- 作業中分類: 使わなかった（既存 story ファイルへ 1 件足す変更のため、通常分類のまま人間レビューを通した）。
- 通常分類: `Screens/翻訳実行`。
- 現在分類: `Screens/翻訳実行`。

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| 案内欄の表示条件 | 案内の文があれば配送方式（同期・xAI）によらず出す。文が空なら出さない | `TranslationRunScreen.svelte` | なし |
| 未訳が残ったときの案内文 | 残った件数と、再実行でその件数だけを翻訳することを 1 文で伝える。件数の内訳（固有名・叙述文・台詞）は出さず、観測ログが持つ | `translation-run-presentation.ts` | なし |
| 未訳が 0 件のときの表示 | 案内欄を出さず、実行を終えた表示だけにする | `TranslationRunScreen.svelte` | なし |

## 現在状態

- 変更ファイル:
    - `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`
    - `frontend/src/ui/screens/translation-run/translation-run-presentation.ts`
    - `frontend/src/ui/screens/translation-run/translation-run.fixtures.ts`
    - `frontend/src/ui/screens/translation-run/TranslationRunScreen.stories.ts`
- 通常分類へ戻した story: なし（作業中分類へ移していないため）。
- 承認状態: 人間承認済み（案内の文言と、案内欄を配送方式で分けない形の 2 点）。
