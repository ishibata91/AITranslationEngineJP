# Task Plan: proper-noun-case-insensitive-replacement

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

- 本文から固有名を見つける処理が大文字小文字を区別するため、機械置換辞書に同じ訳語があっても原文の表記によって置換結果が安定しない。

## 要求

- **R-1 固有名置換の安定化**: 固有名の先頭と末尾に語境界がある本文から固有名を見つける時は大文字小文字を区別せず、大文字小文字が異なる同じ固有名を機械置換辞書の同じ訳語へ置き換える。本文に大文字小文字が異なる同じ固有名が複数ある場合、置換内訳には機械置換辞書の固有名と訳語を1件表示する。機械置換辞書にない固有名の内側で、機械置換辞書の固有名と一致した部分の先頭または末尾に語境界がない場合は置き換えない。

## branch 情報

- `execution_branch`: `codex/proper-noun-case-insensitive-replacement`
- `target_branch`: `master`
- `source_commit`: `7c9937b94e8069cc77a42e625fdfcd62529d52f3`

## やらないこと

- 機械置換辞書の供給源から一般語を除外する条件は変更しない。
- 横断辞書と plugin 内訳語の保存形式は変更しない。
- database に保存済みの言及記録は再計算または削除しない。
- remote repository は変更しない。
