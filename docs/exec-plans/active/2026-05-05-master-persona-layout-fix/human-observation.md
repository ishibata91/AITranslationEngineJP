# 人間観測記録

## 観測

- 生成結果と詳細のレイアウトがプロトタイプと違う。
- 編集モーダルのレイアウトがプロトタイプと違う。
- 編集モーダルのフォントサイズもプロトタイプと違う。

## 期待

- 生成結果と詳細の layout は task-local prototype の `PersonaReviewPanel.svelte` と同じ構造へ寄せる。
- 編集モーダルの layout と font size は task-local prototype の `PersonaActionModal.svelte` と同じ構造へ寄せる。
- 表示項目や機能は増やさない。

## 影響候補

- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
- [PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaActionModal.svelte)
- [App.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts)
