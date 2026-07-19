# test-oracle

処理段の境界で、抽出した属性が正しい値のまま次段へ合流するかの仕様適合を、C# 抽出と Go 翻訳が同じ JSON から照合するためのオラクル。単体テストは段の内側の純粋ルールを見る。ここは段の境界の合流を見る。既存 golden はこのオラクルで置き換える。

## 粒度

- 中心は処理段（stage）。UC は現状 1 本なので、割る軸にしない。
- 1 エントリ = 処理段 × 属性。段ごと・属性ごとに失敗様式が違うので、交点で割る。

## フィールド

| key | 役割 | 値の例 |
|---|---|---|
| `id` | 対応テスト存在確認の join key。C# と Go のテストが同じ id を参照する | `persona-extracted-from-flags` |
| `stage` | どちらのツールの継ぎ目か。担当と観測点が一意に決まる | `extraction`（C# 抽出ツール・esm→抽出結果／staging） / `integration`（Go 翻訳ツール・seed→翻訳出力の継ぎ目） |
| `attribute` | その段で運ばれる属性 | 話者 / ペルソナ / 基底口調 / 台詞感情 / 役割語 / 固有名 / 言及 / 箱振り分け / 出力 |
| `category` | 分類。`given` のリッチさから従属的に決まる | `正常` / `異常` |
| `given` | 入力 esp のリッチさの一断面（下記） | 台詞レコードが TRDT 感情型を持つ |
| `spec` | `given` のもとで観測できる、実現する仕様 | 各台詞の感情を台詞単位でプロンプトへ加算する |

## given は esp のリッチさ

- `given` は抽象的な前提でなく、入力 esp（synthetic Mod）のリッチさの一断面を名指す。どの属性が揃うか、どの属性を欠くか。
- `category`（正常/異常）は `given` から従属する。属性が揃えば正常、欠ければ fallback・保持（異常）。
- fixture は 1 本の synthetic esp に、リッチさの違うレコードを混載する。各エントリの `given` が、その中のどのレコードを見るかを名指す。C# 抽出も Go 構築も同じ 1 fixture を読む。

## coverage フィールド（例外時のみ）

- 既定（フィールド無し）: その spec は合成オラクルの入口から assert する。
- `coverage: "existing-unit"`: master や localized 形式が要り合成 esm では再現できない spec を、既存の実 esm 単体テストへ委ねる印。網羅番人はこの印がある spec を「担保済み」と扱い、取りこぼしと区別する。現状は `override-identical-not-counted` / `container-owns-child-counted` / `localized-header-excluded` の 3 件。

## 書き方の規約

- ドメイン語彙で書く。table・列・status・訳値・prompt 文字列・file・method を `spec` と `given` に入れない。
- 固定名（Female flag、ElderRace、TRDT など）は残してよい。残すときは日本語で意味を補う。

## 例

`specs.json` を参照する。オラクルを読むテストの書き方は、各テストフォルダの `CLAUDE.md`（`tools/extractor.Tests/CLAUDE.md` ほか）に置く。
