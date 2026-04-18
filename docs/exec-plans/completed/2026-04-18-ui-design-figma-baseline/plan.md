# Task Plan: 2026-04-18-ui-design-figma-baseline

- `workflow`: work
- `status`: completed
- `lane_owner`: skill-modification
- `task_id`: 2026-04-18-ui-design-figma-baseline
- `task_mode`: refactor
- `request_summary`: `ui-design` を HTML mock 前提から Figma artifact 前提へ置き換える。
- `goal`: HTML/CSS の構造負債を設計 artifact へ持ち込まず、Figma の構造を UI review と implementation handoff の主入力にする。
- `constraints`: product code、product test、product docs 正本は変更しない。既存 flat plan は migration しない。
- `close_conditions`: workflow、template、`ui-design`、downstream skill が Figma 前提で一致し、structure harness が pass する。

## Artifact Index

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `scenario_design`: `N/A`
- `diagramming`: `N/A`
- `implementation_scope`: `N/A`

## Routing Notes

- `required_reading`: `.codex/skills/ui-design/SKILL.md`, `docs/exec-plans/templates/task-folder/README.md`, `.codex/workflow.md`
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
- `follow_up`: Figma authoring tool がない session では human-provided Figma file / node が必要

## Outcome

- `result`: `ui-design` を HTML mock 前提から Figma file/node 前提へ変更した。
- `validation`: `python3 scripts/harness/run.py --suite structure` passed.
- `completed_at`: 2026-04-18
