# Plan: dictionary-reference-glossary-evaluation

## 事象

- DBの参考用語を提示すると、指定訳との一致は47語中12語から44語へ増えた。
- DBの参考用語を提示すると、意味を維持した文は50件中48件から37件へ減った。
- 一般語への誤適用、長い項目の部分一致、原文情報の欠落が発生した。

## 対象

- 変える対象: 原文に対してDBから参考用語を選び、翻訳AIへ提示するまでの候補選別と提示方法。
- 変えない対象: `dictionary.sqlite3`の収録内容。参考用語の使い方を先に評価するため、DBの採否は変えない。
- 変えない対象: 翻訳モデル`gpt-5.6-luna`、`reasoning_effort=medium`、Batch API。候補選別以外による差を混ぜないため固定する。

## 砂場

- 道具と標本: `tmp/dictionary-reference-glossary-experiment/`
- 根拠: `.gitignore`が`tmp/`を除外していることを`git check-ignore`で確認する。

## 接続先

- OpenAI Batch APIの`/v1/chat/completions`を使う。
- API keyはmacOS Keychainのservice名`OPENAI_API_KEY`から取得する。
- 辞書は`dictionary/dictionary.sqlite3`を読み取り専用で使う。
- 開発用は既存50件、評価用は別に作る未見100件を使う。

## branch情報

- 作業branch: `codex/dictionary-reference-glossary`
- 統合先branch: `master`
- 分岐元commit: `5e971fcc9ffda72a69b4e52737af94fc6d2063f1`

## やらないこと

- 辞書DBの採否をやり直さない。
- プロダクトの翻訳処理へ組み込まない。採用後の実装taskへ回す。
- 開発用50件の結果だけで達成を宣言しない。
- 評価用標本を改善ループ中に閲覧しない。
