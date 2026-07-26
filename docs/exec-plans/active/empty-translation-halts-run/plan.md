# Task Plan: empty-translation-halts-run

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、観測と原因は `investigation.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

`inigo.esp` を実 LLM（LM Studio の `hy-mt2-7b`）で翻訳した実行で人間が見た。

- 画面に次のエラーが出て、翻訳が途中で止まる。`翻訳に失敗: 固有名の翻訳: 翻訳応答解析: 構造化出力の解析失敗: translation が空`
- 固有名 1 件の応答が空だっただけで、その実行の残り全件が訳されないまま終わる。

## 要求

- 要求1: 固有名 1 件の翻訳が skippable な失敗（構造化出力の空・スキーマ違反、応答エンベロープの読み取り失敗、サーバ一時失敗）で終わったとき、その固有名を未訳のまま残して実行を続け、残り全件を訳し切る。本文フェーズ（叙述文・台詞）と同じ扱いに揃える。
- 要求2: 実行完了時に、未訳のまま残した件数を画面へ出す。人間が再実行の必要を判断できるようにする。

## やらないこと

- 空応答に対する再送信（リトライ）の導入。回数と待ち時間の仕様判断が要るため、この task では扱わない。
- 本文フェーズの既存 skip 挙動の変更。
- 未訳固有名の一覧表示。件数の通知までを範囲にする。

## branch 情報

- `execution_branch`: `claude/empty-translation-halts-run`
- `target_branch`: `master`
- `source_commit`: `decdfd9ac9e6f10f665d1a497bf27a806dd80ebd`

