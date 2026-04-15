---
name: design
description: 機能要件、UI モック、Scenario、implementation brief、implementation scope を mode 分岐で固定し、task-local design を work plan へ返す role skill。
---

# Design

## Goal

- 実装前に task-local design を固定する
- `requirements`、`ui-mock`、`scenario`、`implementation-brief`、`implementation-scope` を役割別に作り分ける
- human review 用の判断材料と AI handoff 用の資料を混ぜない

## Modes

- `requirements`: 機能要件と非機能要件を固定する
- `ui-mock`: task-local UI mock working copy と plan 参照を固定する
- `scenario`: Scenario 一覧 working copy と plan 参照を固定する
- `implementation-brief`: human review と実装者 handoff の両方に使う仕様書を固定する
- `implementation-scope`: human review 後の AI handoff 専用資料を固定する

## Common Rules

- 文書は初見の human が判断でき、同時に実装者が補完なしで着手できる粒度で書く
- 各論点は `issue`、`background`、`options`、`recommendation`、`reasoning`、`open_risks` の形にする
- 1 項目 1 論点を守り、複数カテゴリを 1 文に混ぜない
- カテゴリは `Form / UI fields`、`Domain / data model`、`Commands / action semantics`、`State transitions`、`History / operations`、`API / DTO / contracts` を使う
- 各論点には判断基準を明示する
- 固有名詞、既存 field 名、既存 contract 名、mode 名を除き、日本語優先で書く
- `requirements` は active work plan の要件 section を更新する
- `ui-mock` は `docs/exec-plans/active/<task-id>.ui.html` と `docs/mocks/<page-id>/index.html` の最終適用先を plan に更新する
- `scenario` は `docs/exec-plans/active/<task-id>.scenario.md` と `docs/scenario-tests/<topic-id>.md` の最終適用先を plan に更新する
- `implementation-brief` は active work plan の `Request Summary`、`Decision Basis`、`Facts`、`Functional Requirements`、`Work Brief`、`Acceptance Checks`、`Required Evidence`、`HITL Status` と diagram 参照を更新する
- `implementation-scope` は `docs/exec-plans/active/<task-id>.implementation-scope.md` と plan の path 参照だけを更新する
- `implementation-scope` の本文は AI handoff 専用のため英語で圧縮してよい
- 実装コード、product test、docs 正本は変更しない
- human 判断が必要な論点は `open_questions` に分離する
- docs の ER 図や architecture と詳細設計に矛盾がある場合は human に確認する
- 役割を再確定せず、呼び出し元で確定した `design_mode` をそのまま進める

## Output

- `functional_requirements_summary`
- `in_scope`
- `non_functional_requirements`
- `out_of_scope`
- `open_questions`
- `ui_artifact_path`
- `final_mock_path`
- `scenario_artifact_path`
- `final_scenario_path`
- `implementation_brief`
- `implementation_scope_artifact_path`
- `implementation_scope_splits`
- `review_diff_diagrams`
- `source_diagram_targets`
- `canonicalization_targets`

## Detailed Guides

- `references/mode-guides/requirements.md`
- `references/mode-guides/ui-mock.md`
- `references/mode-guides/scenario.md`
- `references/mode-guides/implementation-brief.md`
- `references/mode-guides/implementation-scope.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.design.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.design.<mode>.json` を正本とする
- 返却 contract は `references/contracts/design.to.orchestrate.<mode>.json` を正本とする
