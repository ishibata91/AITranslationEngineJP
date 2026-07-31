# Task Plan: add-openai-provider

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

- OpenAI の `gpt-5.6-luna` の API 料金が大幅に安くなった。
- OpenAI の Batch API を翻訳に使いたい。

## 要求

- **R-1 OpenAI Batch API で翻訳できるようにする**: 翻訳実行画面へ OpenAI（batch）を独立した選択肢として追加し、OpenAI の公式 endpoint、OpenAI API キー、モデル一覧を使って `gpt-5.6-luna` で翻訳対象 plugin の固有名と本文を送信、状態確認、取り込みできるようにし、モデル一覧では `gpt-5.6-luna` を初期選択して他のモデルも選べるようにし、成功した訳された本文は結果一覧で確認でき、失敗した文は未訳のまま再送信できるようにし、OpenAI API キーが空の場合は OpenAI batch を操作できないようにする一方、OpenAI 互換 API と xAI batch の選択肢を維持し、選択肢を切り替えた場合は切替前のモデル一覧と進行表示を持ち越さず、進行中の OpenAI batch と xAI batch を取り違えないようにする。

## branch 情報

- `execution_branch`: `codex/add-openai-provider`
- `target_branch`: `master`
- `source_commit`: `8c125e03a3c1fd4cfdce81749d0c5c938bdbf869`

## やらないこと

- OpenAI の同期翻訳方法は変更しない。
- Chat Completions API から Responses API への移行は行わない。
- OpenAI batch の取消操作は追加しない。
- 自動の状態確認は追加しない。
- 変動する API 料金を画面またはソースへ固定値として表示しない。
