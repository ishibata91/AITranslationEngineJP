---
name: plan-compaction
description: plan_compactor agent が exec-plan の規約違反と重複を、設計判断をせずに整理する時に使う。
---

# Plan compaction

## 入力

- 指定された task folder の `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md`。
- `protocols/docs/exec-plans/coding.md`。

`log.jsonl` を読まない。

## 整理する

次の変更だけを直接行う。

- 同じ意味を重複して書いた文を、一方へ集約する。
- source の path、symbol、外部資料の所在を `references.md` へ移す。
- 明示的に未決定またはブロッカーと書かれた内容を `pending.md` へ移す。
- 文書責務に反する内容を、意味を変えずに正しい文書へ移す。

移動先が一意に決まらない場合は書き換えない。要求に要るかを設計判断しない。要約によって対象、条件、例外、観測可能な振る舞いを失わせない。

## 書かないこと

- 要求、as-is、to-be、仕様の意味を変えない。
- 新しい要求、設計、仕様、pending、参照を作らない。
- 人間の設計HITL後に task folder を書き換えない。
- `log.jsonl` を読む、書く、または入力へ渡す。
- product code、test、task folder 外の file を変更しない。

## 返す成果物

- 直接行った整理の差分。
- 整理しなかった箇所と、整理できない理由。

## 完了条件

- 直接変更が文書責務の移動または重複の集約だけである。
- task folder 外の変更がない。
- `log.jsonl` を読まない。
