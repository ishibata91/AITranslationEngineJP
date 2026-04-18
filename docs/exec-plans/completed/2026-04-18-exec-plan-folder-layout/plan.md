# Task Plan: 2026-04-18-exec-plan-folder-layout

- `workflow`: work
- `status`: completed
- `lane_owner`: skill-modification
- `task_id`: 2026-04-18-exec-plan-folder-layout
- `task_mode`: refactor
- `request_summary`: exec-plan を task ごとの folder にし、skill 対応資料を同じ folder に分ける。
- `goal`: 新規 exec-plan の可読性を上げ、AI が不要資料を読んで context を汚染する状態を避ける。
- `constraints`: 過去の completed / active plan は migration しない。product code と product docs 正本は変更しない。
- `close_conditions`: 新規 workflow と template が folder 形式を正本にし、structure harness が pass する。

## Artifact Index

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `scenario_design`: `N/A`
- `diagramming`: `N/A`
- `implementation_scope`: `N/A`

## Routing Notes

- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `.codex/skills/propose-plans/SKILL.md`, split design skills
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: user requested direct workflow/template change on 2026-04-18

## Copilot Result

- `completed_handoffs`: N/A
- `touched_files`: N/A
- `implemented_scope`: N/A
- `test_results`: N/A
- `ui_evidence`: N/A
- `implementation_review_result`: N/A
- `sonar_gate_result`: N/A
- `residual_risks`: N/A
- `docs_changes`: workflow docs and templates only

## Closeout Notes

- `canonicalized_artifacts`: N/A
- `follow_up`: N/A

## Outcome

- exec-plan の新規正本を `docs/exec-plans/active/<task-id>/` の folder 形式へ変更した。
- `plan.md` は索引と状態管理に限定し、詳細は skill 対応資料へ分離した。
- `templates/task-folder/` に `plan.md`、`requirements-design.md`、`ui-design.html`、`scenario-design.md`、`implementation-scope.md` を追加した。
- `propose-plans` と 4 つの design skill の artifact path を folder 形式へ同期した。
- 過去の flat file plan は legacy として migration しない方針を README に明記した。
- `python3 scripts/harness/run.py --suite structure` は pass した。
