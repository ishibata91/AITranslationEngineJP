# Task Plan: extractor-test-failures

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- `dotnet test tools/extractor.Tests` が 30 件中 18 件失敗する状態を解消する。
- 実データ（Skyrim の Data フォルダ）の位置を `.env` で指定できるようにする。
- 実データが無い機械では、実データを要るテストを失敗ではなく skip にする。
- 一時 SQLite ファイルの削除に失敗する後片付けを直す。

## branch 情報

- `execution_branch`: `claude/extractor-test-failures`
- `target_branch`: `master`
- `source_commit`: `2eae843a`

## やらないこと

- Go 側のテスト（`npm run verify:backend`）は扱わない。失敗しているのは C# のテストだけである。
- 抽出ロジック本体（`PluginExtractor`・各 writer の書き込み内容）は変えない。失敗の原因はテストの前提と後片付けにあり、抽出結果の正しさではない。
- 実 Data フォルダの中身は変えない。読み取りだけを行う。
- C# テストを `npm run` から呼ぶ経路の新設は扱わない。現状どおり `dotnet test` を直接叩く。
