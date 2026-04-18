---
name: implementation-investigate-review-support
description: GitHub Copilot 側の review 補助証跡知識 package。
---

# Implementation Investigate Review Support

## 目的

この skill は知識 package である。
`investigator` agent が implementation review や UI check に必要な証跡を補う時の判断基準を提供する。

## いつ参照するか

- reviewer が追加証跡不足で reroute した時
- review_evidence を補う時
- validation_results と observed_facts を review に渡す時

## 参照しない場合

- design review を行う時
- 実装修正を行う時
- review 判定そのものを行う時

## 原則

- reviewer が判定できる最小証跡だけを集める
- observed_facts と review_evidence を分ける
- remaining_gaps を明示する
- 修正や test 追加を混ぜない

## DO / DON'T

DO:
- review target と evidence の対応を残す
- console、screenshot、command result の出所を残す
- next step を reviewer または orchestrate へ戻す

DON'T:
- pass / fail を investigator が確定しない
- design review をしない
- active contract をこの skill に置かない

## Checklist

- [implementation-investigate-review-support-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-review-support/references/checklists/implementation-investigate-review-support-checklist.md) を参照する。
