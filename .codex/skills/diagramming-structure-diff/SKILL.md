---
name: diagramming-structure-diff
description: active exec-plan の `実装計画` と関連 artifact、既存 component 図から更新対象を特定し、必要なら new component detail 図を判断して、review 用構造差分 D2 / SVG を active exec-plan 配下へ作る。承認後は同じ差分を `docs/diagrams/components/backend/` または `docs/diagrams/components/frontend/` 正本へ適用する。
---

# Diagramming Structure Diff

## Goal

- active exec-plan の `実装計画`、HTML モック artifact、Scenario テスト一覧 artifact と既存 `docs/diagrams/components/backend/` または `docs/diagrams/components/frontend/` を読み、どの source 図を更新するかを特定する
- 既存 detail 図で足りない時は、new component detail 図を作るべきかを判断し、source path を決める
- `proposal_diff` では active exec-plan 配下に review 用構造差分 `.d2` / `.svg` を作る
- `apply_to_source` では承認済み差分を `docs/diagrams/components/backend/` または `docs/diagrams/components/frontend/` 配下の component 図へ適用する

## Workflow

1. 入力契約を確認し、active exec-plan の `要求要約`、`UI モック`、`Scenario テスト一覧`、`実装計画`、`review 用差分図`、`差分正本適用先` を読む。
2. `diagram_mode` が `proposal_diff` か `apply_to_source` かを確認する。
3. task-local design を backend または frontend component 単位へ写像し、まず `docs/diagrams/components/backend/` または `docs/diagrams/components/frontend/` の既存 component map 更新有無を判定する。
4. `proposal_diff` では、対象が backend か frontend かを判定し、対応する architecture 図があるかを確認して layer taxonomy をその正本図から固定する。現在の正本は backend なら `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/backend/backend-architecture.d2`、frontend なら `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/frontend/frontend-architecture.d2` とする。
5. 各 component について、既存 detail 図を更新するか、許可された component 図ディレクトリ配下に new detail 図を新規作成するかを決める。
6. `proposal_diff` では active exec-plan 配下へ review 用差分 `.d2` / `.svg` を出力し、追加を緑、削除を赤で読める状態にする。
7. `proposal_diff` では、package で layer 所属を示し、主要 component は class/member-shaped node で表現する。
8. `apply_to_source` では承認済み差分を source `.d2` へ反映し、対応する `.svg` を更新する。
9. すべての出力で `d2 validate`、`d2 -t 201`、必要時 class 図の縦横比確認まで終える。

## Rules

- `proposal_diff` では component 図正本を変更しない
- 更新対象の特定は、active exec-plan の `実装計画` と関連 artifact、既存 source 図の対応だけで説明できる状態にする
- new component detail 図は、既存 detail 図へ追記すると主題が混ざる時だけ作る
- component map は cross-component の依存と責務境界を主題にし、component detail 図は 1 component を主題に保つ
- `docs/diagrams/components/backend/` と `docs/diagrams/components/frontend/` 以外の図ディレクトリは読まない、書かない、更新対象に含めない
- backend / frontend の `proposal_diff` では、対応する architecture 図の layer 名を source of truth とし、task-local 名称をそのまま layer 名へ昇格させない
- backend / frontend の `proposal_diff` では、layer は package / container で表し、component は class/member-shaped node で表す
- backend / frontend の `proposal_diff` では、package は layer 所属を示す役割、class は component 責務と state surface を示す役割として分ける
- backend / frontend の `proposal_diff` では、task scope に存在しない layer を無理に追加しない。対象 layer だけを architecture 図へ写像して使う
- backend / frontend の `proposal_diff` では、主要 component node に公開 state、責務、操作面を member として載せる。低レベル実装 detail までは落とし込まない
- backend / frontend の `proposal_diff` では、背景色と文字色を明示し、diff 色と可読性が衝突しない状態を render 前提で保つ
- class/member-shaped 表現が component 粒度を壊す時だけ簡略表現へ落としてよいが、その時も package による layer 所属は維持する
- review 用差分図は active exec-plan 配下の一時成果物であり、source of truth にしない
- 承認されていない境界変更や component 分割を `apply_to_source` で追加しない
- validate や render が失敗したまま完了扱いにしない
- `d2` の新しい layout / routing / style 構文を使う時は、最小例で検証してから本図へ入れる
- 許可された component 図ディレクトリが存在しない、または対象 path を安全に対応付けできない時は停止して orchestrator へ返す
- backend / frontend の `proposal_diff` で architecture 図の layer taxonomy に安全に写像できない時は停止して orchestrator へ返す

## Reference Use

- proposal phase では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.diagramming-structure-diff.proposal.json` を参照し、返却時は `references/diagramming-structure-diff.to.orchestrating-implementation.proposal.json` を使う。
- close phase では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.diagramming-structure-diff.json` を参照し、返却時は `references/diagramming-structure-diff.to.orchestrating-implementation.close.json` を使う。
