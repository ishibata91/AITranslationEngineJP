# tools/extractor.Tests（テスト作成時に必ず読む）

このフォルダは C# 抽出器のテストを置く。共有オラクル（`test-oracle/specs.json` の given→spec）を読む結合テストは、下の規約に必ず従う。Go 側（`internal/harness`）も同じ規約に従い、そちらは自身の `CLAUDE.md` を持つ。

## オラクルテストの書き方

目的は「1 テスト = 1 spec の入口→返り」を、誰が見ても同じ形で読めることにする。

- 1 オラクル 1 関数にする。spec 1 件を 1 テスト関数で照合する。複数 spec を switch、表、ループで束ねない。
- オラクル id を引ける形にする。各テスト関数へ対応 spec の id を `[Oracle("<id>")]` 属性で紐づける。網羅番人はこの id 集合と specs.json（自段、非委任）の一致を見る。
- AAA を必須にする。各関数を Arrange（given を用意＝合成入力を読む）、Act（入口を 1 回叩く）、Assert（spec を返りへ照合）の 3 節で書き、節をコメントで示す。
- テストを独立させる。テストは互いに依存しない。共有する可変状態を持たず、各関数が自分で Arrange する。読み取り専用の索きヘルパ（`OracleInput`）は使ってよく、可変 fixture の共有（`IClassFixture` での抽出結果共有など）は禁じる。
- given は入力側、期待値はテスト側に置く。どの record がどの条件を持つかは合成 esm（`test-oracle/fixture/Synthetic.esm`）が持つ。テストは入口の返りへの期待値だけを書く。

## この規約の実装点（C#）

- 属性: `OracleAttribute.cs`（`[Oracle("<id>")]`）。
- Arrange ヘルパ: `OracleInput.cs`（状態なしの索き。`LoadEnv`、`Key`、`Info`、`NewDb` ほか）。
- specs.json loader: `OracleSpecs.cs`。
- 手本: `OracleExtractionTests.cs`（extraction 段。1 関数 1 spec、AAA、網羅番人）。
