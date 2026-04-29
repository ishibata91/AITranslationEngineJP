---
name: implementation-investigate-observe
description: Codex implementation lane 側の一時観測点作業プロトコル。
---
# Implementation Investigate Observe

## 目的

この skill は作業プロトコルである。
`implementation_investigator` agent が owned_scope 内に一時観測点を add / remove する時の判断基準を提供する。

## 対応ロール

- `implementation_investigator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-investigate-observe` の出力規約で固定する。

## 入力規約

- 一時 log、probe、assertion などを使って観測する時
- temporary_changes と cleanup_status を返す時
- 観測点を返却前に除去する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 一時観測点は owned_scope 内に限る
- 観測目的を明確にする
- 返却前に必ず除去する
- cleanup_status を必ず返す

- temporary_changes に path と目的を残す
- cleanup の validation を行う
- 除去不能なら stop する

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- temporary_changes に path と目的を残した。
- 観測点を返却前に除去した。
- cleanup_status を必ず返した。

## 停止規約

- 恒久修正を行う時
- プロダクトテスト を追加する時
- cleanup 不能な観測変更が必要な時
- 観測点を恒久修正として残さない
- owned_scope 外を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 観測点を恒久修正として残さなかった場合は停止する。
- owned_scope 外を変更しなかった場合は停止する。
- cleanup 不能な時に続行しなかった場合は停止する。
