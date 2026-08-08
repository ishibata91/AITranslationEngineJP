# tools/extractor.Tests

抽出機のモデル、SQLite 出力、実データ照合を検証する C# テスト。

- `*SqliteWriterTests.cs` は SQLite table に保存される行を検証する。
- `ModelInvariantTests.cs` は抽出結果の構造上の不変条件を検証する。
- `OracleExtractionTests.cs` は既知の plugin データとの抽出結果を照合する。
- `TempSqliteDb.cs` と `ExtractionCache.cs` はテスト用の一時資源を提供する。
- `RealDataFactAttribute.cs` は実データが必要な検証を明示する。
