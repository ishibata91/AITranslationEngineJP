# 不足セレクタ

## 概要

- 対象: 翻訳実行画面の単語翻訳・NPC ペルソナ生成・本文翻訳の各段階が表示する処理対象一覧領域の固定 selector。
- 判断: 件数表示・空状態・検索入力欄の 3 selector は現行 frontend 実装に存在し、screen-design-diff 差分5 で画面設計正本へ登録済みのため不足ではない。初回取得中ローディングレイヤーの selector は現行実装に対応要素がなく実装時に追加する新規 selector であり、関連テストが依存するため記録する。
- 根拠: `./screen-design-diff.job-run.md` 差分5（注意2: ローディングレイヤーの data-testid は新規確定で実装時に要素を追加）。`./test-design.csv`（E2E-LTLE-004 ほか）。`./test-design-unit.md`（UT-LOAD-001）。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| GAP-LOADING-TERM | 単語翻訳 | 初回取得中ローディングレイヤー（上部の進行状況兼操作区画と処理対象一覧領域を含むフェーズ画面全体を覆うオーバーレイ） | `term-translation-phase-processing-target-loading` | E2E-LTLE-004, UT-LOAD-001 | 現行実装に対応要素がない。初回取得完了までフェーズ画面全体を覆うオーバーレイとして表示し全体操作を排他する新規エレメントを実装時に追加する。selector 値は処理対象一覧由来の名称のまま据え置く（screen-design-diff 注意2）。 |
| GAP-LOADING-PERSONA | NPC ペルソナ生成 | 初回取得中ローディングレイヤー（同上） | `persona-generation-phase-processing-target-loading` | E2E-PGTL-001 | term と同型の新規エレメントを実装時に追加する。 |
| GAP-LOADING-BODY | 本文翻訳 | 初回取得中ローディングレイヤー（同上） | `body-translation-phase-processing-target-loading` | E2E-BTTL-001 | term と同型の新規エレメントを実装時に追加する。 |

## 補足（不足ではない既存 selector）

- `term-translation-phase-processing-target-total` / `-empty` / `-search-input` / `-row` は現行 frontend 実装に存在し、残置 E2E と page object（`tests/system/support/translation-phase-pages.ts`）が参照する。persona / body も prefix 切り替えで同型に存在する。screen-design-diff 差分5 で画面設計正本へ登録済みのため新規追加は不要。
- 空状態の文言は現行実装の `処理対象がありません` を正とする。前タスクの `data-testid-gaps.md` にあった `処理対象が見つかりません。` は誤記として扱う（screen-design-diff 注意3）。
