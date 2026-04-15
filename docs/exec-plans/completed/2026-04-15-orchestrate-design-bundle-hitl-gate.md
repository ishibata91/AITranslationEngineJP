# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/orchestrate, .codex/skills/design, docs/exec-plans/templates
- task_id: 2026-04-15-orchestrate-design-bundle-hitl-gate
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `orchestrate` の human review gate を plan 作成直後ではなく design bundle 完了後に移す
- design bundle が揃うまでは routing を継続できるようにする
- HITL 状態名を新しい gate 位置に合わせる

## Decision Basis

- 現状の `orchestrate` は plan 作成直後に停止するため、design bundle を揃える前に流れが止まる
- design skill は `requirements`、`ui-mock`、`scenario`、`implementation-brief`、`implementation-scope` に分かれている
- human review は implementation handoff 前の判断として、design bundle 完了後に置く方が意図に合う

## Task Mode

- `task_mode`: refactor
- `goal`: design bundle 完了後に human review gate を固定する
- `constraints`: product code と `docs/` 正本は変更しない。変更は workflow skill と plan template に限定する
- `close_conditions`: design bundle 完了後にだけ HITL pending へ遷移し、approval 前は `implementation-scope` 以降へ進まない

## Facts

- `orchestrate` は現状 `required-after-plan` と `pending` を plan 作成直後に要求する
- design contract 上、`requirements`、`ui-mock`、`scenario`、`implementation-brief` が揃うと review と implementation-scope の判断材料が揃う
- `implementation-scope` は HITL 後の handoff 固定として design 側で定義されている

## Functional Requirements

- `summary`: human review gate を design bundle 完了後へ移し、implementation-scope 前で停止する
- `in_scope`: `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/orchestrate/references/permissions.json`, `docs/exec-plans/templates/work-plan.md`
- `non_functional_requirements`: gate 条件と bundle 構成が一読で分かる
- `out_of_scope`: product code, design skill の mode 再編, docs 正本変更
- `open_questions`: なし
- `required_reading`: `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/orchestrate/references/permissions.json`, `.codex/skills/design/SKILL.md`, `.codex/skills/design/references/contracts/design.to.orchestrate.*.json`

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
- `tasks`: active plan を追加する; orchestrate の human review gate 発火条件を design bundle 完了後へ移す; plan template を同期する
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: 呼び出し側 prompt が旧状態名を固定している場合は別途同期が必要

## Acceptance Checks

- plan 作成直後は human review pending を要求しない
- `requirements` `ui-mock` `scenario` `implementation-brief` 完了後にだけ human review gate が立つ
- human review 未完了では `implementation-scope` と downstream implementation handoff を開始しない

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL Status

- `functional_or_design_hitl`: required-after-design-bundle
- `approval_record`: pending-after-design-bundle

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `orchestrate` の human review gate を plan 作成直後から design bundle 完了後へ移した
- design bundle を `requirements`、`ui-mock`、`scenario`、`implementation-brief` として明文化した
- work plan template の HITL 状態名を `required-after-design-bundle` と `pending-after-design-bundle` に更新した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` が通過した
