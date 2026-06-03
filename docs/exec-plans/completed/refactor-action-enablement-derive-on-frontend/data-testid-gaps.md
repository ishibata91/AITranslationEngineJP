# data-testid 不足資料

## 対象 task

`refactor-action-enablement-derive-on-frontend`

## 判断

各フェーズ画面のアクションボタン（`開始` / `中断` / `再開` / `リトライ` / `キャンセル`）は画面設計書の `aria-label` 列に selector が存在する。
無効化時の BlockedReason 文言表示領域と「次段階へ」ボタンは `data-testid` が画面設計書に存在しない。
本表に未決として記録し、実装者が画面設計書へ追記するまでシナリオテスト観点の期待値確認を保留する。

## 不足 selector 一覧

| 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 備考 |
| --- | --- | --- | --- | --- |
| 単語翻訳 | 開始ボタンの無効理由（BlockedReason）表示領域 | `data-testid=term-translation-phase-action-blocked-reason` または同等の selector | RAEF-E2E-002, RAEF-E2E-003, RAEF-E2E-004 | 画面設計書 E-02「操作できない理由」表示に selector 定義なし |
| NPC ペルソナ生成段階 | 開始ボタンの無効理由（BlockedReason）表示領域 | `data-testid=persona-generation-phase-action-blocked-reason` または同等の selector | RAEF-E2E-006 | 同上（persona 段階） |
| 本文翻訳段階 | 開始ボタンの無効理由（BlockedReason）表示領域 | `data-testid=body-translation-phase-action-blocked-reason` または同等の selector | RAEF-E2E-008 | 同上（body 段階） |
| NPC ペルソナ生成段階 | 次段階へボタン（persona → body 移行） | `data-testid=persona-generation-phase-start-next-phase-button` または同等の selector | RAEF-E2E-009, RAEF-E2E-010 | H-11 有効化条件（¬terminal ∧ COMPLETED_PHASE）の確認に必要。persona 側に personaBodyReadiness ガードがないことも合わせて確認する |
| 単語翻訳 | 次段階へボタン（term → persona 移行） | `data-testid=term-translation-phase-start-next-phase-button` または同等の selector | シナリオテスト追加時 | H-5 有効化条件（¬terminal ∧ COMPLETED_PHASE ∧ confirmedCount≥aiTargetCount）の確認に必要 |
