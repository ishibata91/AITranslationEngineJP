---
name: tests-unit
description: GitHub Copilot 側の unit test 補強知識 package。
---

# Tests Unit

## 目的

この skill は知識 package である。
`tester` agent が実装済み責務の分岐と error path を unit test で補う時の判断基準を提供する。

## いつ参照するか

- public contract と主要 branch を確認する時
- error path を unit test にする時
- implementation_task_ids 内の責務を証明する時

## 参照しない場合

- scenario artifact の outcome を test にする時
- test のためだけに広い product code 変更が必要な時
- integration flow を証明する時

## 原則

- 各 test method は 1 つの public behavior、branch、error path のどれか 1 つを証明する
- setup は決定的にする
- test body に条件分岐を入れない
- implementation_task_ids の外まで広げない

## DO / DON'T

DO:
- Arrange / Act / Assert を空行または短いコメントで判別できる状態にする
- branch ごとに test case を分ける
- clock、random、ID、repository 応答順序を固定する

DON'T:
- 新しい要件解釈を足さない
- test のためだけの product code 変更を広げない
- active contract をこの skill に置かない

## Checklist

- [tests-unit-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-unit/references/checklists/tests-unit-checklist.md) を参照する。
