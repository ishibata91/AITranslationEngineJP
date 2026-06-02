# 不足セレクタ

## 概要

- 対象: E2E-UC-048-B1（翻訳実行シェル、境界観点追加候補）
- 判断: 既存 selector で対応可能
- 根拠: E2E-UC-048-B1 が使う selector は `[data-testid=job-run-phase-screen-region]` と `[data-testid=persona-generation-phase-screen]` であり、いずれも既存 E2E-UC-048 と E2E-UC-046 が利用している確定済み selector。新規 selector の追加は不要。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| - | 翻訳実行シェル | - | なし | E2E-UC-048-B1 | E2E-UC-048-B1 が必要とする selector はすべて既存 CSV に記載済み |
