---
name: implementation-investigate-reproduce
description: GitHub Copilot 側の実装前再現知識 package。
---

# Implementation Investigate Reproduce

## 目的

この skill は知識 package である。
`investigator` agent が実装前に再現可否と観測事実を確認する時の判断基準を提供する。

## いつ参照するか

- 実装前に症状や対象挙動を再現する時
- reproduction_status と observed_facts を返す時
- validation command の現状を確認する時

## 参照しない場合

- 実装中の原因 trace を行う時
- 一時観測点を入れる時
- 修正後の再観測を行う時

## 原則

- 再現条件と観測結果を分ける
- evidence のない原因断定をしない
- 再現できない場合も条件と不足情報を返す
- 実装や test 追加を混ぜない

## DO / DON'T

DO:
- command、入力、期待、実際を残す
- reproduction_status を明確にする
- remaining_gaps を次 action へつなげる

DON'T:
- 恒久修正をしない
- design 不足を実装側で補わない
- active contract をこの skill に置かない

## Checklist

- [implementation-investigate-reproduce-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-reproduce/references/checklists/implementation-investigate-reproduce-checklist.md) を参照する。
