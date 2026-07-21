# Storybook レビューループ画面仕様

## 対象

- 確定した story: Screens/翻訳実行（xAI・未入力 / xAI・送信可 / xAI・送信後 / xAI・反映中 / xAI・反映済み）
- 確定した `fixture`: translation-run.fixtures.ts の xAI 5 状態
- 関連資源: TranslationRunScreen.svelte, translation-run-view.ts, translation-run-presentation.ts
- 作業中分類: Screens/翻訳実行（既存 title のまま追加）
- 通常分類: Screens/翻訳実行
- 現在分類: Screens/翻訳実行

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| provider 選択 | AI サービス節の先頭に「同期 / xAI（batch）」の segmented を置く | TranslationRunScreen.svelte, presentation.ts | 解決（承認済み） |
| xAI 操作 | 実行ボタンを「送信」「反映」の 2 ボタンへ差し替え。送信後 info 案内、反映は全 batch まとめ注記 | TranslationRunScreen.svelte, presentation.ts | 解決（承認済み） |
| xAI 接続欄 | エンドポイント欄の hint/placeholder を xAI 用（既定 https://api.x.ai）へ切替、取得 hint を方式別に | presentation.ts | 解決（承認済み） |

## 現在状態

- 変更ファイル: TranslationRunScreen.svelte, translation-run-view.ts, translation-run-presentation.ts, translation-run.fixtures.ts, TranslationRunScreen.stories.ts
- 通常分類へ戻した story: Screens/翻訳実行（xAI・未入力 / xAI・送信可 / xAI・送信後 / xAI・反映中 / xAI・反映済み）。作業中分類は使わず既存 title のまま。
- 承認状態: 承認済み（2026-07-21）。build-storybook 成功、test:frontend 通過。
