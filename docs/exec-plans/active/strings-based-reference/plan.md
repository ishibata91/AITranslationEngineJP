# Task Plan: strings-based-reference

plan は branch 情報と、やること・やらないことの要点を持つ。設計は `design.md`、恒久的に残す判断は `docs/changelog.md`。

## やること

- 参照訳（`reference_translation`）と固有名の確定訳語（`master_term`）の日本語 dest 供給を、xTranslator 英日 XML から Data フォルダの Strings（english / japanese）へ移し、XML 依存を廃止する。
- 入力読み込みを C# 抽出器へ集約する。抽出器が japanese も解決し `extracted_field` に dest 列を書く。engine は DB（`extracted_field`・`record_type_master`）だけから `reference_translation`・`master_term` を組む。
- english / japanese の片方しか無い時に画面警告を出す（既存訳を再利用できず全文 AI 翻訳になる旨）。

## やらないこと

- テーブルの作り直し。追加は `extracted_field` への dest 列のみ（ALTER）。
- 翻訳段・termderive の派生ロジックの変更。入力元だけ替える。
- 抽出対象 plugin の選択方式、翻訳本体（LLM 呼び出し・口調生成・batch 送信）の変更。

## branch 情報

- `execution_branch`: `claude/strings-based-reference`
- `target_branch`: `master`
- `source_commit`: `2aedca6a`
