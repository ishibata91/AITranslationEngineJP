# tools/extractor

Skyrim plugin の構造と文字列を読み、翻訳処理へ渡す事実を SQLite に保存する C# コマンド。

- `Program.cs` はコマンドの入力を受け、抽出と保存を実行する。
- `PluginEnvironment.cs` は Data folder、target plugin、master 連鎖、参照解決を構成する。
- `PluginExtractor.cs` は plugin を `ExtractionResult` へ写す。
- `Model.cs` は抽出結果の構造を定義する。
- `TranslationCounts.cs` は抽出結果を翻訳文字列単位へ平坦化する。
- `*SqliteWriter.cs` は抽出した事実を用途別の SQLite table へ保存する。
- `SchemaMigrator.cs` は保存前に SQLite schema を適用する。
- `XTranslatorXml.cs` は xTranslator の XML を読み書きする。
