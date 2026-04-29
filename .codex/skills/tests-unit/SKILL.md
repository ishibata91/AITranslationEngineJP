---
name: tests-unit
description: Codex implementation レーン 側の 単体 test 補強作業プロトコル。
---
# Tests Unit

## 目的

この skill は作業プロトコルである。
`implementation_tester` agent が実装済み責務の分岐と エラー経路 を 単体 test で補う時の判断基準を提供する。

## 対応ロール

- `implementation_tester` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `tests-unit` の出力規約で固定する。

## 入力規約

- public 契約 と主要 分岐 を確認する時
- エラー経路 を 単体 test にする時
- implementation_task_ids 内の責務を証明する時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 各 test method は 1 つの public behavior、分岐、エラー経路 のどれか 1 つを証明する
- setup は決定的にする
- test body に条件分岐を入れない
- implementation_task_ids の外まで広げない

- Arrange / Act / Assert を空行または短いコメントで判別できる状態にする
- 分岐 ごとに test case を分ける
- cロック、random、ID、repository 応答順序を固定する

## 非対象規約

- シナリオ成果物の結果、統合 flow、新しい要件解釈は扱わない。
- test のためだけの広いプロダクトコード変更は扱わない。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 1 test で 1 public behavior / 分岐 / エラー経路 だけを証明した。
- setup の cロック、random、ID、repository 応答順序を固定した。
- implementation_task_ids の外へ広げなかった。

## 停止規約

- シナリオ 成果物 の 結果 を test にする時
- test のためだけに広い プロダクトコード 変更が必要な時
- 統合 flow を証明する時
- 停止時は不足項目、衝突箇所、戻し先を返す。
- test body に条件分岐を入れなかった場合は停止する。
