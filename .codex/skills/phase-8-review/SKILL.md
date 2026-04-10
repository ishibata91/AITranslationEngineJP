---
name: phase-8-review
description: 第8段階の実装レビューを担当し、実装差分が詳細設計と整合しているかだけを単発で確認する。
---

# Phase 8 Review

## Review Scope

- 実装が詳細設計と違うことをしていないか
- 設計前提を崩す差分がないか
- 第7段階までの証明で主要不足が残っていないか
- sonar MCP でopen_issueがないか確認する

## Output

- decision: `pass` or `reroute`
- findings
- recheck
- closeout_notes

## Rules

- review は 1 回だけ行う
- 新しい改善提案や新しい要件解釈は追加しない
- 実装差分なら第6段階へ、設計差分なら第2段階または第3段階へ差し戻す

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-8-review.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-8-review.to.orchestrating-implementation.json` を返却契約として使う。
