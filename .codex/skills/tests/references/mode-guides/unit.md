# Tests: unit

## Goal

- 実装済み責務の分岐と error path を補う

## Rules

- public contract と主要 branch を優先する
- 各 test method は Arrange / Act / Assert を空行または短いコメントで判別できる状態にする
- 各 test method は 1 つの public behavior、branch、error path のどれか 1 つだけを証明し、検証対象は 1 つに絞る
- assertion を複数置いてよいのは、1 つの戻り値 object や 1 つの state object の中身を確認する場合だけに限る
- setup は決定的にする。clock、random、ID、repository 応答順序は stub や fixture で固定する
- test body に条件分岐を入れない。同じ責務でも branch ごとに test case を分ける
- implementation_task_ids の外まで広げない
- test のためだけの product code 変更は最小にする
