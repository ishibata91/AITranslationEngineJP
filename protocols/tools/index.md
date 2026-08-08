# tools

本体から別プロセスで動く C#/.NET の補助ツールを置く。

- `extractor/` は Skyrim plugin から翻訳に必要な事実を抽出して SQLite へ保存する。
- `extractor.Tests/` は抽出結果と SQLite 出力を検証する。
- `synthetic-fixture/` は抽出機の検証用 plugin を生成する。
