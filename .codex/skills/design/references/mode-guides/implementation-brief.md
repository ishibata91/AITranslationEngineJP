# Design: implementation-brief

## Goal

- human review と実装者 handoff の両方に使う仕様書を固定する
- 実装前提、制約、依存、分割方針、判断理由を明示する

## 必須の書き方

- 各論点は `issue`、`background`、`options`、`recommendation`、`reasoning`、`open_risks` の順で書く
- 1 項目 1 論点を守る
- カテゴリは `Form / UI fields`、`Domain / data model`、`Commands / action semantics`、`State transitions`、`History / operations`、`API / DTO / contracts` を使う
- 各論点に判断基準を書く
- 固有名詞、既存 field 名、既存 contract 名、mode 名を除き、日本語優先で書く

## Capture

- 実装前提
- 制約
- 依存
- 分割方針
- design-review 前に確認したい論点
- human review で確認してほしい判断点
- diagram need
- architecture 変更時だけ `docs/architecture.md` と対象 D2 を `source_diagram_targets` に載せる
- `implementation-scope` へ渡す前提

## Work Plan Mapping

- 目的と前提は `Request Summary`、`Decision Basis`、`Facts` に落とす
- 要件の確定内容は `Functional Requirements` に落とす
- 初期実装に含める範囲は `Work Brief.accepted_scope` に落とす
- 実装対象と分割単位は `Work Brief.implementation_target`、`parallel_task_groups`、`tasks` に落とす
- 確認手段は `Acceptance Checks`、`Required Evidence`、`validation_commands` に落とす
- 未解消事項は `Functional Requirements.open_questions` と `Investigation.residual_risks` に分ける
- human review が必要な判断点は `HITL Status` に接続できる形で書く

## Avoid

- `implementation-scope` 相当の narrow scope 確定
- human 未承認の仕様追加
- architecture 成果物がないのに `source_diagram_targets` を埋めること
