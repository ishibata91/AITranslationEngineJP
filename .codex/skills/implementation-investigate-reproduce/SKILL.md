---
name: implementation-investigate-reproduce
description: Codex implementation lane 側の実装前再現作業プロトコル。
---
# Implementation Investigate Reproduce

## 目的

この skill は作業プロトコルである。
`implementation_investigator` agent が実装前に再現可否と観測事実を確認する時の判断基準を提供する。

## 対応ロール

- `implementation_investigator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-investigate-reproduce` の出力規約で固定する。

## 入力規約

- 実装前に症状や対象挙動を再現する時
- reproduction_status と observed_facts を返す時
- validation command の現状を確認する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- agent runtime と tool policy は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 再現条件と観測結果を分ける
- evidence のない原因断定をしない
- 再現できない場合も条件と不足情報を返す
- 実装や test 追加を混ぜない

- command、入力、期待、実際を残す
- reproduction_status を明確にする
- remaining_gaps を次 action へつなげる

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- command、入力、期待、実際を残した。
- reproduction_status を明確にした。
- observed facts と hypotheses を分けた。

## 停止規約

- 実装中の原因 trace を行う時
- 一時観測点を入れる時
- 修正後の再観測を行う時
- 恒久修正をしない
- design 不足を実装側で補わない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- 原因断定を evidence なしにしなかった場合は停止する。
- 実装や test 追加を混ぜなかった場合は停止する。
- design 不足を実装側で補わなかった場合は停止する。
