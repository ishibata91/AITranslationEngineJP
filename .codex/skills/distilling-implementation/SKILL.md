---
name: distilling-implementation
description: 実装 lane 用。active exec-plan、関連 docs、コードから実装判断に必要な facts / constraints / gaps を抽出する。
---

# Distilling Implementation

## Input

- active exec-plan
- relevant docs under `docs/`
- related code and tests

## Output

- facts
- constraints
- gaps
- docs sync candidates
- required reading

## Rules

- `UI` / `Scenario` / `Logic` は active plan を正本として読む
- packet file を作らない
- `changes/` や `context_board` を前提にしない
- 実装判断が未確定なら gap として返す

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.distilling-implementation.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/distilling-implementation.to.directing-implementation.json` を返却契約として使う。
