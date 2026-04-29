---
name: implementation-investigate-trace
description: Codex implementation lane 側の実装中 trace 作業プロトコル。
---
# Implementation Investigate Trace

## 目的

この skill は作業プロトコルである。
`implementation_investigator` agent が実装中の原因候補、観測点、不足情報を整理する時の判断基準を提供する。

## 対応ロール

- `implementation_investigator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-investigate-trace` の出力規約で固定する。

## 入力規約

- observed facts と hypotheses を分ける時
- 次の observation point を整理する時
- implement、tests、review、reroute の次 action を判断材料として返す時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 観測済み事実と仮説を混ぜない
- trace は owned_scope 内に限定する
- 不足情報を remaining_gaps に残す
- evidence のない結論を固定しない

- hypotheses に根拠と未確認点を付ける
- observation_points を最小にする
- recommended_next_step を根拠付きで返す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- observed facts と hypotheses を分けた。
- observation_points を最小にした。
- recommended_next_step を根拠付きで返した。

## 停止規約

- 実装前再現だけを行う時
- 一時観測点の add / remove が主目的の時
- 恒久修正が主目的の時
- 恒久修正をしない
- プロダクトテスト を追加しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- evidence のない結論を固定しなかった場合は停止する。
- 恒久修正をしなかった場合は停止する。
- プロダクトテスト を追加しなかった場合は停止する。
