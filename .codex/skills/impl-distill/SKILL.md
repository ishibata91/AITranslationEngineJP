---
name: impl-distill
description: 実装 lane 用。active exec-plan、関連 docs、コードから実装判断に必要な facts / constraints / gaps を抽出する。
---

# Impl Distill

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
