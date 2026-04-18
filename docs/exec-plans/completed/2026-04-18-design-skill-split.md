# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/README.md, .codex/workflow.md, .codex/agents/designer.toml, docs/exec-plans/templates
- task_id: 2026-04-18-design-skill-split
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `implementation-brief` を廃止し、旧 `design` の mode を 4 つの独立 skill に分ける。
- `../everything-claude-code` の設計系 skill から、互換性のある考え方だけを吸収する。

## Decision Basis

- `skill-modification` は skill の追加、更新、整理を許可している。
- 新 skill には `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` が必要である。
- 既存 `design` は `requirements`、`ui-mock`、`scenario`、`implementation-brief`、`implementation-scope` の mode 切り替えを持つ。
- `product-capability`、`frontend-design`、`design-system`、`tdd-workflow`、`verification-loop`、`agentic-engineering`、`blueprint`、`architecture-decision-records` の一部が互換性を持つ。

## Task Mode

- `task_mode`: refactor
- `goal`: mode 切り替え型の `design` を廃止し、要件、UI、シナリオ、実装スコープを独立 skill にする。
- `constraints`: product code、product test、docs product 正本は変更しない。AI design review は導入しない。
- `close_conditions`: live workflow が 4 skill を参照し、旧 `design` と `implementation-brief` が live から外れている。

## Facts

- `.codex/README.md` と `.codex/workflow.md` は live skill と設計 flow の正本である。
- `propose-plans` は design bundle 完了後に human review で停止する。
- `implementation-scope` は human review 後にだけ作る。

## Functional Requirements

- `summary`: `implementation-brief` を廃止し、4 つの設計 skill へ責務を分割する。
- `in_scope`: `requirements-design`, `ui-design`, `scenario-design`, `implementation-scope` の追加と workflow 同期。
- `non_functional_requirements`: 日本語優先、責務境界の明確化、Codex/Copilot 境界維持、AI 駆動前提の設計品質向上。
- `out_of_scope`: product code、product test、docs product 正本、AI design review、Copilot 実装 skill の変更。
- `open_questions`: なし。
- `required_reading`: `.codex/skills/skill-modification/SKILL.md`, `.codex/skills/design/SKILL.md`, `.codex/README.md`, `.codex/workflow.md`.

## Design Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Copilot Handoff

- `implementation_scope_artifact_path`: N/A
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`
- `parallel_task_groups`: none
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: session の available skill list は再読込まで旧一覧を表示する可能性がある。

## Acceptance Checks

- `.codex/README.md` の live skill が 4 つの新設計 skill を示す。
- `.codex/workflow.md` が `implementation-brief` なしの flow になっている。
- `propose-plans` が 4 skill を呼び分ける。
- 新 skill に `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` がある。
- 旧 `design` は live workflow から外れている。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: user requested direct skill split on 2026-04-18

## Copilot Result

- `completed_handoffs`: N/A
- `touched_files`: N/A
- `implemented_scope`: N/A
- `test_results`: N/A
- `ui_evidence`: N/A
- `implementation_review_result`: N/A
- `sonar_gate_result`: N/A
- `residual_risks`: N/A
- `docs_changes`: workflow docs only

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `implementation-brief` を live workflow から廃止した。
- 旧 `design` を `requirements-design`、`ui-design`、`scenario-design`、`implementation-scope` の 4 skill に分割した。
- `product-capability`、`frontend-design`、`design-system`、`tdd-workflow`、`verification-loop`、`agentic-engineering`、`blueprint`、`architecture-decision-records` の互換部分だけを新 skill に吸収した。
- 旧 `design` directory を `.codex/.trash/2026-04-18-design-retired-after-split` へ退避した。
- `python3 scripts/harness/run.py --suite structure` は pass した。
