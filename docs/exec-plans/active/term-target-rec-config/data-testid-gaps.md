# data-testid 不足: term-target-rec-config

| 対象画面 | 対象要素 | 必要 selector（候補） | 関連テスト ID | 備考 |
| --- | --- | --- | --- | --- |
| マスター辞書 | XML 取り込み結果のうち 13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）が拒否された旨を確認できる表示 | `master-dictionary-xml-import-rejected-rec-summary`（仮称） | `E2E-TTRC-006` | 現行画面設計（`docs/screen-design/screens/master-dictionary.md`）の `master-dictionary-xml-import-bar` には完了件数表示はあるが、REC 単位の拒否内訳表示 selector は未決。E2E-TTRC-006 は当面 `master-dictionary-dictionary-list` 側で「13 種別由来エントリだけが表示される」観点で証明し、拒否内訳表示は selector 確定後に観点を追記する。|
| 単語翻訳 | 処理対象行の REC を画面 selector から識別する補助属性（`target-id` への REC 埋め込みまたは `data-rec` 等） | `term-translation-phase-processing-target-row.<target-id>` 内に REC を含む `data-rec` 属性（仮称） | `E2E-TTRC-002`, `E2E-TTRC-004`, `E2E-TTRC-005` | 現行画面設計は処理対象行 selector に `<target-id>` プレースホルダだけを定義する。`NPC_:FULL` と `NPC_:SHRT` の区別、13 種別外 REC の不在、共通辞書 hit による AI 対象除外を Playwright から識別するために、REC を判別できる属性が必要。selector 確定までは `target-id` の生成規約（REC を含むか）を実装側で固定し、本資料を更新する。|

## 不足の取り扱い

- selector 不足は本 task 範囲では補完しない。`storybook-module` または画面設計差分の対象となる可能性があるため、人間判断で起動するモジュールへ引き継ぐ。
- 本 task 範囲のシナリオテスト実装時には、画面に表示される情報（一覧、件数、状態）で代替検証する観点を `./test-design.csv` の備考に明記している。
