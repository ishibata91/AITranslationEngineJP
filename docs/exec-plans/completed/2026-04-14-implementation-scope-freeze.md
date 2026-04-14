# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/workflow.md, docs/exec-plans/templates, docs/exec-plans/active, docs/exec-plans/completed
- task_id: 2026-04-14-implementation-scope-freeze
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `mixed` の特別ルールを消す
- 実装前に実装スコープを確定する独立資料を active plan 配下に持つ
- design 後、design-review と HITL 後に orchestrate が狭い `owned_scope` で implement へ handoff できるようにする

## Decision Basis

- `mixed` は本質ではなく、実装 handoff の 1 形態にすぎない
- 実装スコープの確定は design bundle と HITL 後でないと安全に決められない
- work plan の `Work Brief` だけでは、handoff 単位の固定先として弱い

## Task Mode

- `task_mode`: implement
- `goal`: implementation-scope artifact と design mode を追加し、実装手前の scope freeze を formalize する
- `constraints`: product code は変更しない。file 操作は MCP 経由のみとする。既存 task_mode は増やさない
- `close_conditions`: orchestrate / design / implement / review / workflow docs が implementation-scope artifact 前提で接続する。template と README が独立資料の置き場を説明する

## Facts

- 現状の `implementation-brief` は `Work Brief` 更新が中心で、HITL 後の handoff 固定責務を持たない
- `implementation_target: mixed` は orchestrate と implement の双方で narrow task として特別扱いされている
- active plan 配下には scenario template はあるが、implementation scope 用の独立 template はない

## Functional Requirements

- `summary`: implementation-scope を design の独立 mode と独立 artifact へ切り出す
- `in_scope`: `.codex/skills/orchestrate/`, `.codex/skills/design/`, `.codex/skills/implement/`, `.codex/skills/review/`, `.codex/workflow.md`, `docs/exec-plans/templates/`, `docs/exec-plans/active/README.md`, `docs/exec-plans/completed/README.md`
- `non_functional_requirements`: quick contract、mode contract、guide、template の語彙を一致させる
- `out_of_scope`: product 実装、docs 正本内容変更、task_mode 追加
- `open_questions`: なし
- `required_reading`: `.codex/workflow.md`, `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/design/SKILL.md`, `.codex/skills/implement/SKILL.md`

## Artifacts

- `ui_artifact_path`: existing
- `scenario_artifact_path`: existing
- `implementation_scope_artifact_path`: `docs/exec-plans/active/<task-id>.implementation-scope.md`
- `source_diagram_targets`: architecture 変更がある時だけ使う
- `canonicalization_targets`: existing

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: workflow docs、skill docs、handoff contract、exec-plan template
- `parallel_task_groups`: none
- `tasks`: add design mode and contracts; add implementation-scope template; route orchestrate through scope freeze after HITL; remove mixed-special wording; sync implement/review docs
- `validation_commands`: `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null`, `python3 scripts/harness/run.py --suite structure`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: なし

## Acceptance Checks

- `design` に implementation-scope mode が追加される
- orchestrate が `implementation-scope` から `owned_scope` と `implementation_target` を受けて implement へ渡す前提になる
- active plan 配下に implementation-scope 独立資料の template と運用説明が追加される
- `mixed` の特別扱い文言が orchestrate / implement から消える

## Required Evidence

- JSON parse success
- structure harness result

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user requested implementation on 2026-04-14

## Closeout Notes

- `implementation-scope` を design の独立 mode とし、HITL 後の handoff 固定責務を `implementation-brief` から分離した
- `mixed` の特別扱いは削除し、必要なら scope freeze 結果として通常の handoff として扱う形に揃えた

## Outcome

- `docs/exec-plans/templates/implementation-scope.md` を追加し、active plan 配下の独立資料 template を用意した
- orchestrate / design / implement / review / workflow docs を `implementation-scope` artifact 前提へ更新した
- `docs/exec-plans/templates/work-plan.md` に `implementation_scope_artifact_path` を追加した
- active / completed README に implementation-scope artifact の置き場と履歴ルールを追加した
- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
