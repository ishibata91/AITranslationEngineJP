# 別プロセス

別プロセスの実装は `tools/` 配下に置く。

- `extractor/` は C#/.NET と Mutagen による抽出を担当する。
- `extractor/` は SQLite writer を持つ。
- `extractor.Tests/` は抽出を検証する。
- `extractor.Tests/` の検証群には `ModelInvariantTests`、`ExtractedFieldSqliteWriterTests`、`OracleExtractionTests` などを含める。
