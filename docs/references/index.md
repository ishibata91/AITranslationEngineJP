# References Index

関連文書: [`../index.md`](../index.md), [`../../4humans/tech-debt-tracker.md`](../../4humans/tech-debt-tracker.md)

このディレクトリは、外部仕様やベンダー資料の参照方針をまとめる。

## Rules

- 新しい参照資料は、原則として `docs/references/` 配下に追加する
- 生の仕様ダンプを置く場合でも、用途と対象を短く説明する
- 実装や仕様判断に使う場合は、どの文書から参照するかを明示する

## Current Sources

- [`./xtranslator_ref.md`](./xtranslator_ref.md): xTranslator の入出力形式整理
- [`./vendor-api/README.md`](./vendor-api/README.md): ベンダー API の生参照とダンプ置き場

## Migration Note

`docs/references/` を外部参照資料の正本とし、ベンダー API 参照も `vendor-api/` 配下へ統一する。
