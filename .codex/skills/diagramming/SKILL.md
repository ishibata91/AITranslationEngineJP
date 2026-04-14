---
name: diagramming
description: D2、PlantUML、review 用構造差分図を mode 分岐で作成・更新・検証する role skill。
---

# Diagramming

## Goal

- 図の source of truth を保ったまま review 用差分図と正本更新を扱う
- D2、PlantUML、structure diff を 1 skill で切り替える
- format-specific な知識を失わず、live 名だけを統一する

## Modes

- `structure-diff`: active work plan の design と既存 component 図から review 用差分図または正本反映を行う
- `d2`: D2 source を作成・更新し、validate と SVG render を行う
- `plantuml`: PlantUML source を作成・更新し、構文確認を行う

## Operation Modes

- `proposal_diff`: review 用の一時差分図を作る
- `apply_to_source`: 承認済み差分を正本へ反映する
- `standalone`: 単独の図を作成または更新する

## Common Rules

- `.d2` と `.puml` を source of truth にする
- review 用差分図は source of truth にしない
- `structure-diff` は active work plan と component 図の対応だけで更新対象を説明できる状態にする
- D2 では `shape: class` と `shape: sql_table` を優先し、package や class style は最新 docs に沿って使う
- D2 の SVG は `d2 validate` と `d2 -t 201` を通してから返す
- PlantUML は構文確認を通してから返す
- 未検証の layout、routing、style 構文は小さい例で確認してから使う
- 役割を再確定せず、呼び出し元で確定した `diagram_mode` と `operation_mode` を前提に進める

## D2 Notes

- Context7 で `shape: class`、`shape: sql_table`、class style、label 記法を確認済み
- class 図では public / private などの記号を label に保持してよい
- ER では table field と constraint を同じ source に残す

## Detailed Guides

- `references/mode-guides/structure-diff.md`
- `references/mode-guides/d2.md`
- `references/mode-guides/plantuml.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.diagramming.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.diagramming.<mode>.json` を正本とする
- 返却 contract は `references/contracts/diagramming.to.orchestrate.<mode>.json` を正本とする
