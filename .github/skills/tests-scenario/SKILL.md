---
name: tests-scenario
description: GitHub Copilot 側の scenario artifact を product test に反映する知識 package。
---

# Tests Scenario

## 目的

この skill は知識 package である。
`tester` agent が承認済み scenario artifact を product test に落とす時に、主要 outcome を決定的に証明する判断基準を提供する。

## いつ参照するか

- scenario artifact の観点を product test にする時
- happy path と主要 failure path を整理する時

## 参照しない場合

- unit branch だけを補う時
- scenario artifact が未承認の時
- 原因未確定の regression test を書く時
- product code の修正が主目的の時

## 原則

- 各 test method は 1 つの scenario outcome だけを証明する
- setup は決定的にする
- test body に条件分岐を入れない
- runtime event 完了は completion event を観測点にする

## DO / DON'T

DO:
- Arrange / Act / Assert が body 構造で読めるようにする
- happy path と failure path を別 test case に分ける
- fixture や helper は scenario を支える範囲に限定する

DON'T:
- 新しい要件解釈を足さない
- paid real AI API を呼ばない
- active contract をこの skill に置かない

## Checklist

- [tests-scenario-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-scenario/references/checklists/tests-scenario-checklist.md) を参照する。
