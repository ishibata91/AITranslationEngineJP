# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/agents, .codex/skills/*/agents, .codex/skills/orchestrate, docs/exec-plans/active, docs/exec-plans/completed
- task_id: 2026-04-14-agent-noun-alignment
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- primary agent 名を名詞形にする
- `design -> designer` と同じ規則で他の primary agent もそろえる
- 直前の agent naming 変更結果を名詞形へ再同期する

## Decision Basis

- skill 名そのままの agent 名は routing では分かるが、agent role としては名詞形の方が読みやすい
- `design -> designer` が基準なら、他の primary agent も同じ命名規則へそろえるのが自然である
- 参照点は `.codex/agents`、skill 側 `agents/openai.yaml`、`orchestrate`、completed plan に限定される

## Task Mode

- `task_mode`: implement
- `goal`: primary agent naming を noun-based rule に統一する
- `constraints`: live skill 数は変えない。secondary agent は今回の対象外とする
- `close_conditions`: primary agent 8 件の file name、`name =`、skill 参照、orchestrate mapping が名詞形でそろう

## Facts

- 変更前の agent 名は `distill`、`design`、`investigate`、`implement`、`tests`、`review`、`diagramming`、`updating-docs` であった
- `skill-modification` は direct agent として `implementer` を参照し続ける
- completed plan 2 件が current naming を記録していた

## Functional Requirements

- `summary`: primary agent 名を noun-based naming へ再統一する
- `in_scope`: `.codex/agents/*.toml`, `.codex/skills/*/agents/openai.yaml`, `.codex/skills/orchestrate/SKILL.md`, `docs/exec-plans/completed/2026-04-14-orchestrate-agent-skill-uniqueness.md`, `docs/exec-plans/completed/2026-04-14-agent-name-alignment.md`
- `non_functional_requirements`: 命名規則が一貫し、参照切れがない
- `out_of_scope`: product code, docs 正本, secondary agent rename
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/completed/2026-04-14-agent-name-alignment.md`

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
- `tasks`: active plan を追加する; primary agent file を noun-based name へ rename する; ref 名を名詞形へ更新する; completed plan を追従更新する
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: `docs_updater` は skill 名の語幹と完全一致しないが、noun rule を優先する

## Acceptance Checks

- primary agent 8 件の file name が名詞形に一致する
- 対象 `.toml` の `name =` が名詞形に一致する
- 対象 skill の `agents/openai.yaml` が名詞形 agent 名を指す
- `orchestrate` と completed plan が名詞形 agent 名を指す

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user requested noun-based agent naming on 2026-04-14

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- primary agent file を `distiller.toml`、`designer.toml`、`investigator.toml`、`implementer.toml`、`tester.toml`、`reviewer.toml`、`diagrammer.toml`、`docs_updater.toml` にそろえた
- 各 primary agent の `name =` を `distiller`、`designer`、`investigator`、`implementer`、`tester`、`reviewer`、`diagrammer`、`docs_updater` に更新した
- `distill`、`design`、`investigate`、`implement`、`tests`、`review`、`diagramming`、`updating-docs` の `agents/openai.yaml` を名詞形 agent 名へ更新した
- `orchestrate` の `Primary Skill-Agent Mapping` を `distill -> distiller`、`design -> designer`、`investigate -> investigator`、`implement -> implementer`、`tests -> tester`、`review -> reviewer`、`diagramming -> diagrammer`、`updating-docs -> docs_updater` に更新した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` が通過した
