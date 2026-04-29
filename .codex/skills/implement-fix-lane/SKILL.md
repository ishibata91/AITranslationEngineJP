---
name: implement-fix-lane
description: Codex implementation レーン 側の fix レーン 恒久修正作業プロトコル。
---
# Implement Fix Lane

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が `承認済み修正範囲` の恒久修正を行う時に、再現条件と矛盾しない変更へ限定する判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-fix-lane` の出力規約で固定する。

## 入力規約

- `task_mode: fix` の 承認済み実装範囲 を実装する時
- 再現根拠 または trace 結果 に基づき修正する時
- 残留リスク と未解消ケースを 終了処理 に残す時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 承認済み修正範囲 を超えない
- 再現条件に関係しない整理を入れない
- trace_or_analysis_result と矛盾しない変更に限る
- 単一引き継ぎ入力 と 承認済み修正範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_tester 出力 も確認する
- 未解消ケースを 終了処理 に残す

- 修正前後で同じ条件の 検証 を比較する
- 残留リスク を明示する
- fix 対象範囲 と touched files を対応づける

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 承認済み修正範囲 と 再現根拠 を確認した。
- trace_or_analysis_result と矛盾しない変更に限定した。
- 残留リスク と未解消ケースを 終了処理 に残した。

## 停止規約

- 新機能や refactor の実装を行う時
- 再現条件が不足している時
- 原因が未確認なのに恒久修正する時
- unrelated cleanup を混ぜない
- 原因断定を 根拠 なしに広げない
- プロダクトテスト、検証データ、スナップショット、test helper を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- unrelated cleanup を混ぜなかった場合は停止する。
- 原因断定を 根拠 なしに広げなかった場合は停止する。
- task_mode が fix であることを確認した。
