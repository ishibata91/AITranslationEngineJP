# data-testid 不足: fix-lucien-target-list-empty

## 目的

追加候補テスト観点(`e2e-diff.md` の E2E-LTLE-001 から E2E-LTLE-003)が参照する画面要素のうち、画面設計正本「E2E 固定 selector」に未登録の固定 selector を記録する。

## 参照した正本

- 画面設計正本: `docs/screen-design/screens/term-translation-phase.md` の「E-04 処理対象一覧」「E2E 固定 selector」

## 既に固定済みの selector

- 画面全体領域: `term-translation-phase-screen`
- 処理対象行: `term-translation-phase-processing-target-row.<target-id>`
- 処理対象一覧領域: `aria-label=処理対象一覧`

## 不足する selector

| 対象画面 | 対象要素 | 必要 selector(候補) | 関連テスト ID |
| --- | --- | --- | --- |
| 単語翻訳 | 処理対象件数の表示 | 画面設計で未確定。値の確定は画面設計正本の更新を要する | E2E-LTLE-001, E2E-LTLE-002 |
| 単語翻訳 | 処理対象一覧の空状態(処理対象が見つかりません。) | 画面設計で未確定。値の確定は画面設計正本の更新を要する | E2E-LTLE-001, E2E-LTLE-002, E2E-LTLE-003 |
| 単語翻訳 | 処理対象一覧内の検索入力欄 | 画面設計で未確定。値の確定は画面設計正本の更新を要する | E2E-LTLE-002 |

## 注意

- 本作業では固定 selector の値を独断で確定しない。
- 値の確定と画面設計正本「E2E 固定 selector」への登録は、画面設計を扱う担当の判断とする。
- プロダクトコードと docs 正本本文は変更していない。
