---
name: tests-scenario
description: GitHub Copilot 側の scenario artifact を product test に反映する知識 package。
---

# Tests Scenario

## 目的

この skill は知識 package である。
`tester` agent が承認済みシステムテスト設計を product test に落とす時に、主要 outcome を決定的に証明する判断基準を提供する。
この skill の主対象は `UI人間操作E2E` と `APIテスト` である。

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
- `UI人間操作E2E` は、承認済みシナリオの開始操作を模倣する
- UI が入口のシナリオでは、画面操作、ファイル選択、フォーム入力などのユーザー入力を開始点にする
- `APIテスト` は、public seam、request / response contract、外部入力開始、主要観測点を開始点にする
- 裏側の直接呼び出しや fixture 直接投入だけの試験は、明示された補助試験でない限り主 `UI人間操作E2E` にしない

## DO / DON'T

DO:
- Arrange / Act / Assert が body 構造で読めるようにする
- happy path と failure path を別 test case に分ける
- fixture や helper は scenario を支える範囲に限定する
- UI が入口の場合は、ユーザー入力から得られる値を `UI人間操作E2E` の検証対象にする
- `APIテスト` では request / response contract と external input start を検証対象にする

DON'T:
- 新しい要件解釈を足さない
- paid real AI API を呼ばない
- active contract をこの skill に置かない
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しない

## Checklist

- [tests-scenario-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-scenario/references/checklists/tests-scenario-checklist.md) を参照する。
