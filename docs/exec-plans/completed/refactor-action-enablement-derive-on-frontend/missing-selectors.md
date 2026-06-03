# 不足セレクタ

## 概要

- 対象: `refactor-action-enablement-derive-on-frontend` task のシナリオテスト観点で必要な selector
- 判断: 各フェーズ画面のボタン操作は `aria-label` ベースの selector が画面設計書に存在する（`開始` / `中断` / `再開` / `リトライ` / `キャンセル`）。一方、無効化時の BlockedReason 文言表示領域と、persona → body 移行の「次段階へ」ボタンの `data-testid` は画面設計書に固定されていない。
- 根拠: docs/screen-design/screens/term-translation-phase.md、persona-generation-phase.md、body-translation-phase.md の操作テーブル（aria-label 列）を確認した。「操作できない理由」表示の selector が定義されていないことを確認した。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| SEL-001 | 単語翻訳 | 開始ボタンの無効理由表示 | `data-testid` または `aria-label` で BlockedReason 文言を識別できる selector | RAEF-E2E-002, RAEF-E2E-003, RAEF-E2E-004 | 画面設計書の「操作できない理由」表示領域に selector が定義されていない。期待値の BlockedReason 文言確認に必要 |
| SEL-002 | NPC ペルソナ生成段階 | 開始ボタンの無効理由表示 | `data-testid` または `aria-label` で BlockedReason 文言を識別できる selector | RAEF-E2E-006 | 同上。persona 段階の「操作できない理由」表示領域の selector が未定義 |
| SEL-003 | 本文翻訳段階 | 開始ボタンの無効理由表示 | `data-testid` または `aria-label` で BlockedReason 文言を識別できる selector | RAEF-E2E-008 | 同上。body 段階の「操作できない理由」表示領域の selector が未定義 |
| SEL-004 | NPC ペルソナ生成段階 | 次段階へボタン（persona → body 移行） | `data-testid` または `aria-label` で「次段階へ」ボタンを識別できる selector | RAEF-E2E-009, RAEF-E2E-010 | 画面設計書に「次段階へ」ボタンの selector が定義されていない。H-11 の有効化条件確認に必要 |
| SEL-005 | 単語翻訳 | 次段階へボタン（term → persona 移行） | `data-testid` または `aria-label` で「次段階へ」ボタンを識別できる selector | RAEF-UNIT-008（シナリオ追加時） | 画面設計書に term 段階の「次段階へ」ボタン selector が定義されていない |
