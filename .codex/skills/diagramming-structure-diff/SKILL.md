---
name: diagramming-structure-diff
description: active exec-plan の `実装計画` と関連 artifact、既存 backend 図から更新対象を特定し、必要なら new component detail 図を判断して、review 用構造差分 D2 / SVG を active exec-plan 配下へ作る。承認後は同じ差分を `diagrams/backend/` 正本へ適用する。
---

# Diagramming Structure Diff

## Goal

- active exec-plan の `実装計画`、HTML モック artifact、Scenario テスト一覧 artifact と既存 `diagrams/backend/` を読み、どの source 図を更新するかを特定する
- 既存 detail 図で足りない時は、new component detail 図を作るべきかを判断し、source path を決める
- `proposal_diff` では active exec-plan 配下に review 用構造差分 `.d2` / `.svg` を作る
- `apply_to_source` では承認済み差分を `diagrams/backend/components.d2` と `diagrams/backend/<component>/<component>.d2` へ適用する

## Workflow

1. 入力契約を確認し、active exec-plan の `要求要約`、`UI モック`、`Scenario テスト一覧`、`実装計画`、`review 用差分図`、`差分正本適用先` を読む。
2. `diagram_mode` が `proposal_diff` か `apply_to_source` かを確認する。
3. task-local design を backend component 単位へ写像し、まず `diagrams/backend/components.d2` の更新有無を判定する。
4. 各 component について、既存 detail 図を更新するか、`diagrams/backend/<component>/<component>.d2` を新規作成するかを決める。
5. `proposal_diff` では active exec-plan 配下へ review 用差分 `.d2` / `.svg` を出力し、追加を緑、削除を赤で読める状態にする。
6. `apply_to_source` では承認済み差分を source `.d2` へ反映し、対応する `.svg` を更新する。
7. すべての出力で `d2 validate`、`d2 -t 201`、必要時 class 図の縦横比確認まで終える。

## Rules

- `proposal_diff` では `diagrams/backend/` 正本を変更しない
- 更新対象の特定は、active exec-plan の `実装計画` と関連 artifact、既存 source 図の対応だけで説明できる状態にする
- new component detail 図は、既存 detail 図へ追記すると主題が混ざる時だけ作る
- component map は cross-component の依存と責務境界を主題にし、component detail 図は 1 component を主題に保つ
- review 用差分図は active exec-plan 配下の一時成果物であり、source of truth にしない
- 承認されていない境界変更や component 分割を `apply_to_source` で追加しない
- validate や render が失敗したまま完了扱いにしない
- `d2` の新しい layout / routing / style 構文を使う時は、最小例で検証してから本図へ入れる

## Reference Use

- proposal phase では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.diagramming-structure-diff.proposal.json` を参照し、返却時は `references/diagramming-structure-diff.to.orchestrating-implementation.proposal.json` を使う。
- close phase では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.diagramming-structure-diff.json` を参照し、返却時は `references/diagramming-structure-diff.to.orchestrating-implementation.close.json` を使う。
