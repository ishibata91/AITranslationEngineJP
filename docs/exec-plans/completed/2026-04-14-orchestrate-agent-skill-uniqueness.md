# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/orchestrate, .codex/skills/*/agents, .codex/agents, docs/exec-plans/active
- task_id: 2026-04-14-orchestrate-agent-skill-uniqueness
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `orchestrate` の skill と agent の組み合わせを一意にする
- 重複している skill 側 agent 指定を整理する
- agent の description と developer instructions を単一責務に寄せる

## Decision Basis

- 現状の `orchestrate` は 1 agent が複数 skill を兼務する書き方になっている
- `agents/openai.yaml` 側でも `default` と複数 agent コメントが混在し、正本が読みにくい
- skill と agent の主担当を 1 対 1 に固定した方が routing と保守の基準が明確になる

## Task Mode

- `task_mode`: implement
- `goal`: workflow 契約上の primary skill-agent mapping を一意化する
- `constraints`: product code と `docs/` 正本は変更しない。live workflow の skill 数は増やさない
- `close_conditions`: `orchestrate`、各 skill の `agents/openai.yaml`、対応 agent 定義の責務が同じ組み合わせを指す

## Facts

- live skill は `orchestrate`、`distill`、`investigate`、`design`、`implement`、`tests`、`review`、`diagramming`、`skill-modification`、`updating-docs` の 10 個である
- downstream skill のうち `design`、`investigate`、`review`、`diagramming`、`updating-docs` は `default` や複数 agent コメントを含んでいた
- `orchestrate` の旧 `Handoff Agents` section は同一 agent を複数 skill に割り当てていた

## Functional Requirements

- `summary`: primary handoff を skill ごとに 1 agent へ固定する
- `in_scope`: `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/*/agents/openai.yaml`, `.codex/agents/*.toml`
- `non_functional_requirements`: naming と責務の説明を検索しやすく保つ
- `out_of_scope`: product code, tests, `docs/` 正本, legacy skill の復活
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/completed/2026-04-14-skill-role-compression.md`, `docs/exec-plans/completed/2026-04-14-skill-role-rehydration.md`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: workflow contract files only
- `parallel_task_groups`: none
- `tasks`: active plan を追加する; `orchestrate` の primary mapping を更新する; 対応する skill agent 定義と agent metadata を同期する
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: `skill-modification` と `orchestrate` 自体の direct agent は別途再整理の余地がある

## Acceptance Checks

- `orchestrate` に primary skill-agent mapping が 1 対 1 で明記される
- 対象 skill の `agents/openai.yaml` が primary agent を明示する
- 対応 agent の `description` と `developer_instructions` が single-role で読める

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user requested unique skill-agent mapping on 2026-04-14

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `orchestrate` の `Handoff Agents` section を廃止し、`Primary Skill-Agent Mapping` として `distill -> distiller`、`design -> designer`、`investigate -> investigator`、`implement -> implementer`、`tests -> tester`、`review -> reviewer`、`diagramming -> diagrammer`、`updating-docs -> docs_updater` を明示した
- `design`、`investigate`、`tests`、`review`、`diagramming`、`updating-docs` の `agents/openai.yaml` を primary agent 指定へ統一した
- primary agent の `description` と `developer_instructions` を single-role に整理した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` が通過した
