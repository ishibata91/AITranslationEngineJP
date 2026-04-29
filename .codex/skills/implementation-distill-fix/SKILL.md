---
name: implementation-distill-fix
description: Codex implementation lane 側の fix 向け context 圧縮作業プロトコル。
---
# Implementation Distill Fix

## 目的

この skill は作業プロトコルである。
`implementation_distiller` agent が fix handoff を整理する時に、症状、再現済み事実、修正対象、validation entry を分ける判断基準を提供する。

## 対応ロール

- `implementation_distiller` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-distill-fix` の出力規約で固定する。

## 入力規約

- 再現済み症状と trace 結果を整理する時
- `accepted_fix_scope` を実装前 packet に圧縮する時
- 未解消ケースと residual risk を残す時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 症状、再現済み事実、仮説、未確認事項を混ぜない
- 長い log や stack trace は要点と path / command に圧縮する
- 修正対象と validation entry だけを実装者向けに残す
- 失敗を閉じるために必要な fix_ingredients を構造単位で残す
- 再現に似ているだけで修正に不要な context は distracting_context に分ける
- 修正対象は path、symbol、line number、変更種別で返す
- accepted fix scope、決定済み方針、禁止事項は implementation_implementer が再読不要な粒度で要約する
- 再現条件に関係しない整理を入れない

- reproduction evidence を path / command と一緒に残す
- trace_or_analysis_result と accepted_fix_scope を対応づける
- fix_ingredients と distracting_context を分ける
- first_action と change_targets を path、symbol、line number 付きで返す
- requirements_policy_decisions に fix 方針と out of scope を残す
- residual risk と未解消ケースを分ける

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- 症状、再現済み事実、仮説、未確認事項を分けた。
- accepted_fix_scope と修正対象を対応づけた。
- residual risk と未解消ケースを残した。

## 停止規約

- 新規実装の context を整理する時
- refactor の不変条件を整理する時
- 原因を evidence なしで断定する時
- 原因断定を先取りしない
- 実 code を読まず handoff の文章を言い換えない
- 類似 context を required_reading に混ぜない
- fix 方針や決定事項を required_reading に丸投げしない
- プロダクトコード / プロダクトテスト を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 原因断定を先取りしなかった場合は停止する。
- 再現条件に関係しない整理を入れなかった場合は停止する。
- プロダクトコード / プロダクトテスト を変更しなかった場合は停止する。
