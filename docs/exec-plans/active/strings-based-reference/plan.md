# Task Plan: strings-based-reference

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- 翻訳前処理の参照訳（既存訳の再利用）と固有名辞書を、Data フォルダの Strings（english / japanese）由来で解決できるようにする。
- xTranslator 英日 XML（`dictionaries/xTranslatorXMLs`）を必須依存から外す。XML が無くても翻訳前処理が破綻せず進むようにする。
- 翻訳対象を Data フォルダで指定すると、その english / japanese Strings から参照訳を解決する経路を、実際に参照訳が保存・利用されるところまで通す。

## branch 情報

- `execution_branch`: `claude/strings-based-reference`
- `target_branch`: `master`
- `source_commit`: `2aedca6a`

## やらないこと

- 翻訳本体（LLM 呼び出し・口調生成・batch 送信）の変更。
- 抽出対象 plugin の選択方式（ファイルピッカー）の変更。
- Windows ビルド手順（別 task `build-windows` の成果物）の再設計。相対パス解決の見直しは別途扱う。
