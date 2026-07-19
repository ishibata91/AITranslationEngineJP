# internal/harness（テスト作成時に必ず読む）

このフォルダは Go 翻訳パイプラインの結合テストを置く。共有オラクル（`test-oracle/specs.json` の given→spec）を読む結合テストは、下の規約に必ず従う。C# 側（`tools/extractor.Tests`）も同じ規約に従い、そちらは自身の `CLAUDE.md` を持つ。

## これは結合テストであって e2e ではない

入口から出口までを通し、出口の値を spec ごとに照合する。実 LLM も UI も通らないので e2e とは呼ばない。

- 入口: `SyntheticRun`（合成 seed → 取込 → 翻訳）。ツール単位の入口→出口にする（C# 抽出とはプロセスを分ける）。
- 出口: `Capture`（送信プロンプト列 ＋ 最終 DB の訳・状態・配置）。
- 途中の package 継ぎ目（取込→注入→後処理）を外から単体のように確かめない。入口→出口を通せば、途中の結合点は自ずと担保される。単段で純粋に閉じるルールは core package の単体テストが担保する。

## オラクルテストの書き方

- 1 オラクル 1 関数にする。spec 1 件を 1 テスト関数で照合する。複数 spec を switch、ループで束ねない。
- オラクル id を引ける形にする。Go には属性が無いため、id→テスト関数の登録（`goOracles` map）で id を持つ。map の各値が 1 spec の 1 関数。網羅番人はこの map のキー集合と specs.json（go 段、非委任）の一致を見る。
- AAA を守る。入口→出口（`SyntheticRun`）を 1 回通し、その read-only の出口を spec ごとに照合する（パラメタライズド）。Arrange＝合成入力（fixture）、Act＝入口を 1 回通す、Assert＝各 spec 関数が出口の値を照合、の 3 節を保つ。
- テストを独立させる。spec 関数は互いに依存しない。出口 `Capture` は read-only なので共有してよい。可変状態は共有しない。
- given は入力側、期待値はテスト側に置く。どの record がどの条件を持つかは合成 fixture（`SyntheticFixture`）が持つ。テストは出口への期待値だけを書く。

## この規約の実装点（Go）

- 登録＋網羅番人＋入口ヘルパ: `oracle_test.go`。
- 合成入力: `fixture.go`（`SyntheticFixture`）。given が足りない spec はここへ最小追加する。
- 出口の観測: `Capture`（`golden.go`）＝送信プロンプトと最終 DB 状態。
