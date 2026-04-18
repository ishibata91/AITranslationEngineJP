---
name: implementation-investigate-reobserve
description: GitHub Copilot 側の修正後再観測知識 package。
---

# Implementation Investigate Reobserve

## 目的

この skill は知識 package である。
`investigator` agent が実装後または test 後に同じ条件で再観測する時の判断基準を提供する。

## いつ参照するか

- 修正後に同じ条件で再観測する時
- reproduction_status と validation_results を比較する時
- remaining_gaps と residual_risks を返す時

## 参照しない場合

- 初回再現を行う時
- 実装中 trace が必要な時
- review の判定だけを行う時

## 原則

- 事前の reproduction condition と同じ条件で観測する
- 変更前後の差を observed_facts として返す
- 未解消ケースを remaining_gaps に残す
- 実装修正を同時に行わない

## DO / DON'T

DO:
- command、入力、期待、実際を比較する
- residual risk を根拠付きで残す
- recommended_next_step を返す

DON'T:
- 条件を変えて pass にしない
- product test を追加しない
- active contract をこの skill に置かない

## Checklist

- [implementation-investigate-reobserve-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-reobserve/references/checklists/implementation-investigate-reobserve-checklist.md) を参照する。
