# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/orchestrate, .codex/skills/skill-modification, docs/exec-plans/templates
- task_id: 2026-04-15-orchestrate-plan-hitl-gate
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `orchestrate` が active work plan 作成後に human review を必須化する
- human review 前に downstream handoff へ進まないようにする
- HITL の記録先を plan 上で読みやすくする

## Decision Basis

- 現状の `orchestrate` は plan 作成後に human review 未完了でも先へ進みやすい
- `SKILL.md` と `permissions.json` の stop 条件が曖昧で、停止判断が弱い
- active work plan の HITL 欄も pending と approved の境界が明確ではない

## Task Mode

- `task_mode`: refactor
- `goal`: plan 後の human review gate を workflow 契約として固定する
- `constraints`: product code と `docs/` 正本は変更しない。変更は workflow skill と plan template に限定する
- `close_conditions`: `orchestrate` が plan 後に human review 未完了なら停止する。plan template で pending と approval record が明示される

## Facts

- `orchestrate` は `HITL` 管理を役割に含むが、plan 後停止を必須化していない
- `permissions.json` の stop 条件は human decision が必要なのに記録先が不明な場合に限られている
- work plan template には `functional_or_design_hitl` と `approval_record` はあるが、期待値が固定されていない

## Functional Requirements

- `summary`: plan 後の human review gate を明示し、承認前の downstream handoff を禁止する
- `in_scope`: `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/orchestrate/references/permissions.json`, `docs/exec-plans/templates/work-plan.md`
- `non_functional_requirements`: human review 時に判断材料と状態が一読で分かる
- `out_of_scope`: product code, other skill routing redesign, docs 正本変更
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/orchestrate/references/permissions.json`

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
- `tasks`: active plan を追加する; `orchestrate` の plan 後 human review gate を明文化する; permissions と plan template を同期する
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: downstream skill 側 prompt が旧運用を前提にしている場合は別途調整が必要

## Acceptance Checks

- plan 作成直後の `orchestrate` で human review pending が明示される
- human review 未完了では downstream handoff を開始しない
- active work plan の HITL 欄で pending と approval record が読める

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL Status

- `functional_or_design_hitl`: required-after-plan
- `approval_record`: pending

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `orchestrate` の役割と routing に、plan 作成直後の human review gate を追加した
- `orchestrate` の stop 条件と permission に、`required-after-plan` と `pending` のまま停止する条件を追加した
- work plan template の HITL 欄に、期待する状態値を明記した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` が通過した
