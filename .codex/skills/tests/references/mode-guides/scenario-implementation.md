# Tests: scenario-implementation

## Goal

- Scenario artifact または fix 再現条件をそのまま証明する

## Rules

- happy path と主要 failure path を含める
- 各 test method は Arrange / Act / Assert を body 構造で判別できる状態にする
- 各 test method は 1 つの scenario outcome だけを証明し、検証対象は 1 つに絞る
- assertion を複数置いてよいのは、1 つの completion event payload、1 つの response body、1 つの rendered section の中身を確認する場合だけに限る
- setup は決定的にする。fixture の入力値、clock、runtime 応答、seed は固定する
- test body に条件分岐を入れない。happy path と failure path は別 test case に分ける
- fixture や helper は scenario を支える範囲に限定する
- runtime event 完了は completion event を観測点にする
