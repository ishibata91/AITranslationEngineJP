# data-testid 不足

## 概要

- 対象: fix-phase-ai-model-list-empty で E2E テスト観点 E2E-UC-FAKE-001〜006 の確認に必要な selector
- 判断: モデル select の data-testid が 3 phase の画面設計書（term-translation-phase、persona-generation-phase、body-translation-phase）に未定義である。`credential_missing` 案内メッセージ表示領域（modelStatusText 用 selector）は今回方針では発生しなくなるため追加不要。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| SEL-001 | 単語翻訳 | AI モデル選択カードのモデル select | `term-translation-phase-ai-model-select` | E2E-UC-FAKE-001, E2E-UC-FAKE-004 | モデル選択肢の有無と fake モデル表示を selector で確認するために必要 |
| SEL-002 | NPC ペルソナ生成段階 | AI モデル選択カードのモデル select | `persona-generation-phase-ai-model-select` | E2E-UC-FAKE-002, E2E-UC-FAKE-005 | モデル選択肢の有無と fake モデル表示を selector で確認するために必要 |
| SEL-003 | 本文翻訳段階 | AI モデル選択カードのモデル select | `body-translation-phase-ai-model-select` | E2E-UC-FAKE-003, E2E-UC-FAKE-006 | モデル選択肢の有無と fake モデル表示を selector で確認するために必要 |
