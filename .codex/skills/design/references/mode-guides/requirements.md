# Design: requirements

## Goal

- task-local requirement を固定する
- human review で判断しやすく、同時に実装者が補完なしで読める仕様論点をそろえる

## 必須の書き方

- 各論点は `issue`、`background`、`options`、`recommendation`、`reasoning`、`open_risks` の順で書く
- 1 項目 1 論点を守る
- カテゴリは `Form / UI fields`、`Domain / data model`、`Commands / action semantics`、`State transitions`、`History / operations`、`API / DTO / contracts` を使う
- 各論点に判断基準を書く
- 固有名詞、既存 field 名、既存 contract 名、mode 名を除き、日本語優先で書く

## Capture

- `summary`: 何を決める task かを短く書く
- `in_scope`: 初期実装に含める仕様だけを書く
- `non_functional_requirements`: 品質、再実行性、整合性、運用制約を書く
- `out_of_scope`: 将来拡張や今回は扱わない内容を書く
- `open_questions`: human review が必要な判断点だけを書く
- `required_reading`: 根拠になる既存仕様や関連資料を書く

## Work Plan Mapping

- 目的と前提は `Request Summary`、`Decision Basis`、`Facts` に落とす
- 要件の確定内容は `Functional Requirements` に落とす
- 未解消事項は `Functional Requirements.open_questions` と `Investigation.residual_risks` に分ける
- human review が必要な判断点は `HITL Status` に接続できるように書く

## Avoid

- 実装手順の詳細化
- 背景なしの二択提示
- 初期実装範囲と将来構想の混在
- docs 正本の更新
