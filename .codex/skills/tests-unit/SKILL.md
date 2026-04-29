---
name: tests-unit
description: Codex implementation lane 側の unit test 補強作業プロトコル。
---
# Tests Unit

## 目的

この skill は作業プロトコルである。
`implementation_tester` agent が実装済み責務の分岐と error path を unit test で補う時の判断基準を提供する。

## 対応ロール

- `implementation_tester` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `tests-unit` の出力規約で固定する。

## 入力規約

- public contract と主要 branch を確認する時
- error path を unit test にする時
- implementation_task_ids 内の責務を証明する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 各 test method は 1 つの public behavior、branch、error path のどれか 1 つを証明する
- setup は決定的にする
- test body に条件分岐を入れない
- implementation_task_ids の外まで広げない

- Arrange / Act / Assert を空行または短いコメントで判別できる状態にする
- branch ごとに test case を分ける
- clock、random、ID、repository 応答順序を固定する

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- 1 test で 1 public behavior / branch / error path だけを証明した。
- setup の clock、random、ID、repository 応答順序を固定した。
- implementation_task_ids の外へ広げなかった。

## 停止規約

- scenario artifact の outcome を test にする時
- test のためだけに広い プロダクトコード 変更が必要な時
- integration flow を証明する時
- 新しい要件解釈を足さない
- test のためだけの プロダクトコード 変更を広げない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- test body に条件分岐を入れなかった場合は停止する。
- test のためだけの プロダクトコード 変更を広げなかった場合は停止する。
- 新しい要件解釈を足さなかった場合は停止する。
