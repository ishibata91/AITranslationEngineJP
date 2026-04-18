# Task Plan: 2026-04-18-investigate-skill-split

- `workflow`: work
- `status`: completed
- `lane_owner`: skill-modification
- `task_id`: 2026-04-18-investigate-skill-split
- `task_mode`: refactor
- `request_summary`: `investigate` を設計前調査と実装時調査へ分け、実装時調査を Copilot 側 skill に移す。
- `goal`: 設計前再現確認は Codex、実装前再現確認、実装中 trace、修正後再観測、実装 review 補助は Copilot 側で扱う。
- `constraints`: product code、product test、product docs 正本は変更しない。既存 completed plan は migration しない。
- `close_conditions`: Codex / Copilot の investigate 責務が分離し、implementation-orchestrate が実装時調査を route でき、structure harness が pass する。

## Artifact Index

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `scenario_design`: `N/A`
- `diagramming`: `N/A`
- `implementation_scope`: `N/A`

## Routing Notes

- `required_reading`: `.codex/skills/investigate/SKILL.md`, `.github/skills/implementation-orchestrate/SKILL.md`, `.github/skills/review/SKILL.md`
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: user requested direct workflow / skill split on 2026-04-18

## Copilot Result

- `completed_handoffs`: N/A
- `touched_files`: N/A
- `implemented_scope`: N/A
- `test_results`: N/A
- `implementation_investigation`: N/A
- `ui_evidence`: N/A
- `implementation_review_result`: N/A
- `sonar_gate_result`: N/A
- `residual_risks`: N/A
- `docs_changes`: workflow docs and skills only

## Closeout Notes

- `canonicalized_artifacts`: N/A
- `follow_up`: N/A

## Outcome

- `result`: `investigate` を Codex 側の設計前調査に狭め、Copilot 側に `implementation-investigate` を追加した。
- `validation`: `python3 scripts/harness/run.py --suite structure` passed.
- `completed_at`: 2026-04-18
