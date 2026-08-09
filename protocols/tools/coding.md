# C# tools の境界

- C# の抽出結果と Go backend は SQLite schema を契約にする。
- schema は `db/migrations/` を正本とし、C# 固有の中間 schema を作らない。
- C# は plugin から観測できる事実を保存し、翻訳方針と分類の最終判断を持たない。
- CLI の通常結果は stdout、入力不備と失敗は stderr へ出し、失敗時は非 zero exit code を返す。
- tools の変更後は `dotnet test tools/extractor.Tests` を実行する。
