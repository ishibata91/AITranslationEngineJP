---
name: planning-implementation
description: "direction から渡された context summary と active exec-plan をもとに、ordered scope と validation を短い implementation brief に変換する。"
---

# Planning Implementation

## Goal

- 実装順を固定する
- owned scope を固定する
- required reading と validation commands を固定する
- 実装前に確認すべき relevant な repo guardrail を固定する

## Rules

- 必要なら active exec-plan の `Implementation Plan` だけを更新してよい
- `Implementation Plan` はモジュール単位で task section を分ける
- 各 task section は契約に依存し、そのモジュールの責務だけを実装対象にする
- input に `mcp_memory_recall` がある時は MCP memory bucket (`repo_conventions`, `recurring_pitfalls`) の中から今回の task に効く項目だけを work brief に残す
- recall は実装ガードレールとして扱い、`docs/` 正本の代わりにしない
- `tasks.md` を作らない
- frontend / backend の責務境界を曖昧にしない

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.planning-implementation.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/planning-implementation.to.directing-implementation.json` を返却契約として使う。
