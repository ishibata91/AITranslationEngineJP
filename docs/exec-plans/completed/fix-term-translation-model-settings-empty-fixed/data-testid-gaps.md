# 不足セレクタ

## 概要

- 対象: 単語翻訳画面（`term-translation-phase-screen`）の AI モデル選択カード（E-03）配下の状態 pill、select、開始ボタン禁止理由領域。
- 判断: 不足あり。
- 根拠: 画面設計（`docs/screen-design/screens/term-translation-phase.md` 91-129 行）は AI モデル選択カードを `aria-label` ベースで規定しており、`data-testid` が固定されていない。`missing-tests.md` の E2E-UC-FIX-MODEL-001/002/003 は状態 pill のテキスト、開始ボタン禁止理由テキスト、select の値変化を確認するため、対象要素を `data-testid` で安定特定したい。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| SEL-FIX-MODEL-001 | 単語翻訳 | AI モデル設定パネルの状態 pill（「固定済み」または「設定未完了」を表示する要素） | `data-testid=term-translation-phase-ai-settings-status-pill` 相当 | E2E-UC-FIX-MODEL-001 | 状態 pill のテキスト確認に `aria-label` がないため、表示テキストへの依存を減らし状態遷移を安定確認したい。 |
| SEL-FIX-MODEL-002 | 単語翻訳 | 「開始」ボタンの禁止理由を表示する領域（モデル未選択の警告を含む） | `data-testid=term-translation-phase-start-blocked-reason` 相当 | E2E-UC-FIX-MODEL-002 | 禁止理由テキストの差分を確認するが、現状はクラス階層に依存する。固定 selector が必要。 |
| SEL-FIX-MODEL-003 | 単語翻訳 | AI モデル選択カードの `AIサービス` / `モデル` / `処理方式` の各 select | `data-testid=term-translation-phase-ai-provider-select` / `term-translation-phase-ai-model-select` / `term-translation-phase-ai-execution-mode-select` 相当 | E2E-UC-FIX-MODEL-001, E2E-UC-FIX-MODEL-003 | `aria-label=AIサービス` は画面設計に存在するが、E2E ではカード横断で同名 label が衝突する可能性があるため、画面固有 `data-testid` で識別したい。 |

## 補足

- 本 task は test_designer による判断固定のみで、フロントエンドへ `data-testid` を追加する実装は本 task の範囲外。実装に進む場合は別途 implementation-scope を作る必要がある。
- 修正方針が presenter 正規化のみで完結する場合、上記 selector 追加は E2E 安定化のための後続 task として扱える。最小限の修正範囲を守るため、本 task では「不足記録」までに留める。
