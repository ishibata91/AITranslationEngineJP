---
name: implementation-distill-refactor
description: Codex implementation lane 側の refactor 向け context 圧縮作業プロトコル。
---
# Implementation Distill Refactor

## 目的

この skill は作業プロトコルである。
`implementation_distiller` agent が refactor handoff を整理する時に、不変条件、依存境界、preserved behavior を抽出する判断基準を提供する。

## 対応ロール

- `implementation_distiller` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-distill-refactor` の出力規約で固定する。

## 入力規約

- refactor handoff を実装前 packet に圧縮する時
- 変更してはいけない振る舞いを整理する時
- affected package / component / tests をまとめる時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- agent runtime と tool policy は [implementation_distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 不変条件、依存境界、変更候補を別々に圧縮する
- preserved behavior を守るための fix_ingredients を構造単位で残す
- refactor に似ているだけの周辺 context は distracting_context に分ける
- 似た責務の file は cluster としてまとめる
- 実装手順ではなく、守る制約を残す
- refactor 開始点は path、symbol、line number、変更種別で返す
- preserved behavior、決定済み方針、禁止事項は implementation_implementer が再読不要な粒度で要約する
- validation command と completion signal を明示する

- preserved behavior を先に固定する
- 代表 path と差分だけを残す
- dependency direction を明示する
- fix_ingredients と distracting_context を分ける
- first_action と change_targets を path、symbol、line number 付きで返す
- requirements_policy_decisions に preserved behavior と out of scope を残す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- 不変条件、依存境界、変更候補を分けた。
- preserved behavior を明示した。
- affected package / component / tests を整理した。

## 停止規約

- 新機能実装の facts を整理する時
- fix 再現条件を整理する時
- broad refactor を新たに提案する時
- 追加の設計判断をしない
- 実 code を読まず handoff の文章を言い換えない
- 類似 context を required_reading に混ぜない
- refactor 方針や決定事項を required_reading に丸投げしない
- owned_scope 外の broad refactor を広げない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- 追加の設計判断をしなかった場合は停止する。
- owned_scope 外の broad refactor を広げなかった場合は停止する。
- product code / product test を変更しなかった場合は停止する。
