# Investigation: extractor-test-failures

`investigation.md` は不具合の再現と原因究明だけを持つ。どう直すかの設計は `design.md` が持つ。

## 観測済み問題

`dotnet test tools/extractor.Tests` を実行すると、30 件中 18 件が失敗する。失敗の理由は 2 系統に分かれる。

系統 A（9 件）は plugin ファイルが見つからない失敗である。

- 例外: `System.IO.FileNotFoundException : plugin が data folder に無い: Dawnguard.esm`（7 件）、同 `Update.esm`（2 件）。
- 対象: `ModelInvariantTests` の 9 件すべて。

系統 B（9 件）は一時ファイルを削除できない失敗である。

- 例外: `System.IO.IOException : The process cannot access the file '<temp>\<prefix>-<guid>.sqlite3' because it is being used by another process.`
- 対象: `ExtractedFieldSqliteWriterTests` 3 件、`SpeakerSqliteWriterTests` 2 件、`OracleExtractionTests` 4 件。
- 一時ファイル名の接頭辞は `ef-`・`sp-`・`oracle-` の 3 種で、テストごとに別々の一時 DB を作っている。

系統 B は assert の失敗ではない。テスト本体の検証は通り、後片付けの段で例外が出て失敗と記録されている。

## 画面再現確認

この不具合は画面に現れない。人間観測記録に含まれる操作はテストコマンドの実行だけで、Wails app の画面操作・`data-testid` を伴う導線は関与しない。よって `chrome-devtools` MCP ツールでの画面再現確認は行わず、再現はコマンド実行で行った。

再現手順と結果は次のとおり。

1. repo root で `dotnet test tools/extractor.Tests` を実行する。
2. 結果は `失敗: 18、合格: 12、スキップ: 0、合計: 30`。
3. 失敗理由を集計すると、`FileNotFoundException` が 9 件、`IOException` が 9 件で、上記の系統 A・B に一致する。

再現は安定しており、実行のたびに同じ 18 件が同じ理由で失敗する。

## 原因仮説

系統 A について 2 つの仮説を立てた。検証は A1 を先にする。テストが読む固定パスの実在確認だけで判定でき、A1 が支持されれば A2 の検証範囲が狭まるためである。

- A1: テストが読む実データの場所（`TestPaths.DataFolder` = `<repo>/dictionaries/Data`）がこの機械に存在しない。
- A2: `PluginEnvironment.Load` のパス解決そのものが壊れており、Data フォルダの有無に関係なく plugin を見つけられない。

系統 B について 3 つの仮説を立てた。検証は B3 を先にする。直前の commit（`2eae843a`）で SQLite の依存パッケージを更新しており、更新が原因かどうかを先に切り分ける必要があるためである。

- B1: `Microsoft.Data.Sqlite` の接続プールが `SqliteConnection.Dispose` の後もファイルを開いたまま保持し、`File.Delete` が拒否される。
- B2: writer またはテストのどこかで接続を閉じ忘れている。
- B3: 直前の commit で入れた `Microsoft.Data.Sqlite 10.0.10` と `SQLitePCLRaw.bundle_e_sqlite3 3.0.4` への更新が原因である。

## 観測ログ検証

プロダクトコードへ一時ログは追加していない。仮説はいずれも、既存の実行結果の観測と、scratchpad に置いた最小再現プログラムの実行で判定できたためである。追加した一時ログが無いので、削除も発生しない。

**A1（支持）**: `ls dictionaries/` の結果に `Data` は無い。`dictionaries/` 直下にあるのは xTranslator XML 4 件と抽出 JSON 5 件だけである。一方、実 Data フォルダは同じ機械の `F:\SteamLibrary\steamapps\common\Skyrim Special Edition\Data` に存在し、`Skyrim.esm`・`Update.esm`・`Dawnguard.esm`・`HearthFires.esm`・`Dragonborn.esm` と `Strings` フォルダを含むことを確認した。テストが読む場所と実データの置き場所が食い違っている。

**A2（否定）**: 同じ `PluginEnvironment.Load` を使う `OracleExtractionTests` は、合成 esm（`test-oracle/fixture/Synthetic.esm`）の読み込みに成功している。4 件の失敗理由はいずれも系統 B の `IOException` であり、読み込み段の例外ではない。よってパス解決の実装自体は動いている。

**B3（否定）**: `tools/extractor/extractor.csproj` を更新前（`Microsoft.Data.Sqlite 9.0.0`、`SQLitePCLRaw` の直接参照なし）へ戻して `dotnet test` を実行したところ、同じ 18 件が同じ理由で失敗した。依存更新の前から存在する失敗である。

**B1（支持）**: scratchpad に最小再現プログラムを置き、`TempDb.Dispose` と同じ手順（接続を `using` で閉じた直後に `File.Delete`）を実行した。結果は次のとおり。

| 接続の作り方 | `File.Delete` の結果 |
| --- | --- |
| 既定（プール有効） | 失敗。`IOException`（別プロセスが使用中） |
| 接続文字列に `Pooling=False` | 成功 |
| 閉じたあと `SqliteConnection.ClearPool` を呼ぶ | 成功 |

テストの失敗と同じ例外・同じ文言が、接続プールを有効にした場合だけ再現した。

**B2（否定）**: 一時 DB を触る箇所を走査した。テスト側（`OracleInput.cs:55`、`ExtractedFieldSqliteWriterTests.cs:53`・`104`・`135`、`SpeakerSqliteWriterTests.cs:66`・`87`）と writer 側（`ExtractedFieldSqliteWriter.cs:17`、`InfoConditionSqliteWriter.cs:17`、`InfoEmotionSqliteWriter.cs:15`、`SpeakerSqliteWriter.cs:18`）はすべて `using var conn` で閉じている。閉じ忘れは無く、B1 の観測と整合する。

## 確定原因

**系統 A の確定原因**: テストが読む実データの場所が `TestPaths.DataFolder`（`<repo>/dictionaries/Data`）に固定されており、実データの置き場所を外から指定する手段が無い。この機械の実データは別ドライブの Steam 配下にあるため、`ModelInvariantTests` は実行のたびに必ず失敗する。`dictionaries` は `.gitignore:12` で追跡外の利用者供給データであり、置き場所は機械ごとに異なる。

**系統 B の確定原因**: `Microsoft.Data.Sqlite` の接続プールが、`SqliteConnection.Dispose` の後も一時 DB ファイルを開いたまま保持する。テストの後片付けはプールを解放せずに `File.Delete` を呼ぶため、削除が拒否されて `IOException` になる。該当する後片付けは 6 箇所ある（`OracleInput.cs:68`、`ExtractedFieldSqliteWriterTests.cs:69`・`121`・`143`、`SpeakerSqliteWriterTests.cs:74`・`95`）。
