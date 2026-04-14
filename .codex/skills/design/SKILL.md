---
name: design
description: 機能要件、UI モック、Scenario、implementation brief を mode 分岐で固定し、task-local design を work plan へ返す role skill。
---

# Design

## Goal

- 実装前に task-local design を固定する
- requirements、UI、Scenario、implementation brief を mode ごとに作り分ける
- 実装と docs 正本更新を混ぜず、active plan と task-local artifact に閉じる

## Modes

- `requirements`: 機能要件と非機能要件を固定する
- `ui-mock`: task-local UI モック working copy と plan 参照を固定する
- `scenario`: Scenario テスト一覧 working copy と plan 参照を固定する
- `implementation-brief`: 実装順、ownership、validation、diagram need を brief 化する

## Common Rules

- `requirements` は active work plan の要件 section だけを更新する
- `ui-mock` は `docs/exec-plans/active/<task-id>.ui.html` と plan 参照だけを更新する
- `scenario` は `docs/exec-plans/active/<task-id>.scenario.md` と plan 参照だけを更新する
- `implementation-brief` は active work plan の `Work Brief` と diagram 参照だけを更新する
- 実装コード、product test、docs 正本は変更しない
- human 判断が必要な論点は `open_questions` に分離する
- diagram が必要な時だけ `diagramming` を起動する
- 役割を再確定せず、呼び出し元で確定した design mode をそのまま進める

## Output

- `functional_requirements_summary`
- `in_scope`
- `non_functional_requirements`
- `out_of_scope`
- `open_questions`
- `ui_artifact_path`
- `scenario_artifact_path`
- `implementation_brief`
- `review_diff_diagrams`
- `source_diagram_targets`

## Detailed Guides

- `references/mode-guides/requirements.md`
- `references/mode-guides/ui-mock.md`
- `references/mode-guides/scenario.md`
- `references/mode-guides/implementation-brief.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.design.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.design.<mode>.json` を正本とする
- 返却 contract は `references/contracts/design.to.orchestrate.<mode>.json` を正本とする
