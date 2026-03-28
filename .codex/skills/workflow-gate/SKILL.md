---
name: workflow-gate
description: AITranslationEngineJp 専用。Architect が approved plan、実装差分、checks、docs sync を照合し、accept 前の gate 判定を返したいときに使う。
---

# Workflow Gate

この skill は Architect が使う read-only gate です。
approved plan、実装差分、checks 結果、docs 更新有無、残る non-blocking unknown を照合し、accept 前に必要な判定だけを返します。

## 使う場面

- heavy / light のどちらでも Coder 実装後に使う
- plan 適合性、証跡不足、docs 同期漏れ、reroute 要否を機械的に確認したい
- 好みや美しさではなく、契約準拠だけを見たい

## 入力契約

- approved plan
- 実装差分の要約
- 実行した checks と結果
- docs 更新有無
- 残っている non-blocking unknown
- 未解消点

## 出力

- decision
- missing_evidence
- contract_breaks
- docs_sync
- recheck

## Gate Rules

- `decision` は `pass` / `block` / `reroute` のいずれかにする
- plan に書かれた checks や required evidence が不足していれば `missing_evidence` に出す
- plan 外の仕様判断、docs sync 漏れ、reroute trigger 該当を `contract_breaks` に出す
- blocking unknown が新たに見つかったら `reroute` を返す
- 美しさや好みは判定しない

## 禁止

- 実装しない
- 好みの改善提案を主目的にしない
- `light-review` と同じ観点を標準 gate として重複運用しない
