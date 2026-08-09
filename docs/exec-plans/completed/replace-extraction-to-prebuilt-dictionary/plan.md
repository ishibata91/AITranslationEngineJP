# Task Plan: replace-extraction-to-prebuilt-dictionary

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

- 現在の `dictionary/dictionary.sqlite3` には、R-5からR-7で作成し、人間レビューを含めて収録を確定した事前作成済みの翻訳辞書が保存されている。
- migrationを適用する中心DBである `db/aitranslation.dev.sqlite3` は削除済みである。中心DBはアプリ起動時に既存migrationから生成する。
- 翻訳実行は中心DBの `master_term` と `proper_noun` から機械置換辞書を毎回組み、`prepareForTranslation` は辞書抽出処理で `DeriveMasterTerms` を実行している。
- 翻訳実行は現在の `dictionary/dictionary.sqlite3` を読んでいない。

## 要求

- **R-1 事前作成済み翻訳辞書へ置き換える**: DBに保存済みの事前作成済み翻訳辞書で翻訳前の機械置換を行い、既存の辞書抽出処理全体を翻訳辞書の供給元として使用しないこと。
- 追加の人間指示: 文字長貪欲マッチ方式は変えず、マッチした単語を参考語として翻訳指示へ載せる。単語に複数の意味がある場合はカテゴリを含む全候補を翻訳指示へ載せる。
- 追加の人間指示: 事前作成済み翻訳辞書DBを `dictionary/` から `db/` へ移動する。
- 追加の人間指示: dictionary viewerを退役し、MCPだけを `db/dictionary.sqlite3` へ移植する。
- 追加の人間指示: `dictionary/` にはMCPの実行に必要なファイルだけを残し、事前作成済み翻訳辞書SQLiteは `db/dictionary.sqlite3` だけを残す。
- 追加の人間指示: modの固有名を先に翻訳し、本文翻訳の参考語へ事前作成済み翻訳辞書DBの候補と翻訳済みの固有名を載せる。
- 追加の人間指示: `dictionary_sense.meaning` は本文翻訳のpromptと翻訳結果UIの参考語から除外する。

追加の人間指示により、R-1の翻訳前機械置換は本文を変更せず参考語を翻訳指示へ渡す動作へ置き換える。

## branch 情報

- `execution_branch`: `codex/replace-extraction-to-prebuilt-dictionary`
- `target_branch`: `master`
- `source_commit`: `da8533a042d9aff0c9a46731698a94f1bc0aca49`

## やらないこと

- `db/dictionary.sqlite3` の収録判断、人間レビュー、辞書項目を変更しない。
- 事前作成済み翻訳辞書をGitへ追加しない。
- 本taskが求める参考語の追加以外のAIの翻訳指示を変更しない。
- 翻訳結果本文の保存形式を変更しない。
- 配布先への事前作成済み翻訳辞書DBの共有を扱わない。実装完了後の開発時はrepository内の `db/dictionary.sqlite3` が存在することを前提にする。
- DBに保存済みの事前作成済み翻訳辞書へない語句を、翻訳実行中に辞書へ追加しない。
