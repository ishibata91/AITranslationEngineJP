# Storybook レビューループ画面仕様

## 対象

- 確定した story: Screens/翻訳実行（xAI・状態未確認 / 状態確認中 / 処理中 / 完了段あり / 取り込み済み）、Screens/翻訳実行/進行状況パネル（未確認 / 処理中 / 完了段あり / 全完了）
- 確定した `fixture`: translation-run.fixtures.ts の xAI 進行状況 各状態
- 関連資源: BatchProgressPanel.svelte, TranslationRunScreen.svelte, translation-run-view.ts, translation-run-presentation.ts
- 作業中分類: Screens/翻訳実行（既存 title のまま追加）
- 通常分類: Screens/翻訳実行
- 現在分類: Screens/翻訳実行

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| xAI 操作 | 「反映」を廃し「状態確認」＋「主アクション」の 2 ボタンへ。主アクションは状態でラベル切替（送信して開始 / 取り込んで本文を送信 / 取り込んで完了）、処理待ちありでグレーアウト | TranslationRunScreen.svelte, presentation.ts | 解決（承認済み） |
| 進行状況パネル | 「固有名→本文→完了」の 3 段ステッパー＋現在段の件数（総数・処理待ち・成功・失敗）。未確認は控えめ表示 | BatchProgressPanel.svelte | 解決（承認済み） |

## 現在状態

- 変更ファイル: BatchProgressPanel.svelte(+stories), TranslationRunScreen.svelte, translation-run-view.ts, translation-run-presentation.ts, translation-run.fixtures.ts, TranslationRunScreen.stories.ts
- 通常分類へ戻した story: Screens/翻訳実行（xAI・未入力 / 送信可 / 送信後 / 状態確認中 / 固有名処理中 / 固有名完了 / 本文処理中 / 本文完了 / 取り込み済み）、Screens/翻訳実行/進行状況パネル（未確認 / 固有名処理中 / 固有名完了 / 本文処理中 / 本文完了 / 全完了）。作業中分類は使わず既存 title のまま。
- 承認状態: 承認済み（2026-07-21）。レビューで 2 点（送信/取り込みの主アクション統合、2 段ステッパー化）を反映し再承認。build-storybook 成功、test:frontend 通過。
