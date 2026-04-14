# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/README.md, .codex/workflow.md, docs/exec-plans
- task_id: 2026-04-14-docs-canonicalization-flow
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- orchestrate の close 手順を、成果物が存在するものだけ `docs/` 正本へ適用する方式に揃える
- 対象は mock、scenario test、architecture の 3 系統とする
- architecture は `docs/architecture.md` と `docs/diagrams/backend|frontend/*.d2` をセットで扱う

## Decision Basis

- 現状の workflow は task-local artifact の作成までは明示しているが、`docs/` 正本への昇格条件が artifact 種別ごとに formalize されていない
- `design bundle` を一括必須にすると、mock / scenario / architecture が独立に存在しうる repo 運用とずれる
- active plan と completed plan の README も、個別成果物の昇格ルールとして読み直せる形に揃える必要がある

## Task Mode

- `task_mode`: implement
- `goal`: workflow contract、guide、template を同期し、artifact 存在時だけ `docs/` 正本へ反映する close フローを formalize する
- `constraints`: product code は変更しない。`docs/` 正本更新の人間承認ルールは破らない。file 操作は MCP 経由のみとする
- `close_conditions`: orchestrate / design / review の contract が個別 artifact 前提で接続する。plan template で最終適用先を記録できる。completed plan に正本化結果を残せる

## Facts

- mock の正本は `docs/mocks/<page-id>/index.html`
- scenario の正本は `docs/scenario-tests/<topic-id>.md`
- architecture の正本は `docs/architecture.md` と `docs/diagrams/backend|frontend/*.d2`
- review 用差分図は正本ではない

## Functional Requirements

- `summary`: artifact 存在時のみ `docs/` 正本へ昇格する workflow を固定する
- `in_scope`: `.codex/skills/orchestrate/`, `.codex/skills/design/`, `.codex/skills/review/`, `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/templates/work-plan.md`, `docs/exec-plans/active/README.md`, `docs/exec-plans/completed/README.md`
- `non_functional_requirements`: quick contract、mode contract、guide、template の語彙を一致させる
- `out_of_scope`: product 実装、既存 docs 正本の内容変更、new skill 追加
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `docs/mocks/README.md`, `docs/scenario-tests/README.md`, `docs/architecture.md`

## Artifacts

- `ui_artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_mock_path`: `docs/mocks/<page-id>/index.html`
- `scenario_artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_scenario_path`: `docs/scenario-tests/<topic-id>.md`
- `review_diff_diagrams`: architecture 変更がある時だけ review 用差分図を持つ
- `source_diagram_targets`: architecture 変更がある時だけ `docs/diagrams/backend|frontend/*.d2`
- `canonicalization_targets`: 存在する artifact だけを列挙する

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: workflow docs、skill docs、handoff contract、plan template
- `parallel_task_groups`: none
- `tasks`: orchestrate close 条件更新; design / review の artifact 契約更新; workflow docs と exec-plan template 更新; plan README の昇格ルール更新
- `validation_commands`: `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null`, `python3 scripts/harness/run.py --suite structure`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: なし

## Acceptance Checks

- `orchestrate -> design -> review -> close` の記述が個別 artifact 前提でつながる
- `docs/exec-plans/templates/work-plan.md` で final path と canonicalization target を記録できる
- active / completed README が artifact 種別ごとの昇格ルールを説明する

## Required Evidence

- JSON parse success
- structure harness result

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user requested implementation on 2026-04-14

## Closeout Notes

- `canonicalized_artifacts`: artifact が存在する時だけ mock / scenario / architecture を `docs/` 正本へ反映する close ルールを formalize した
- review 用差分図は正本ではなく、architecture 正本への反映有無を plan に残す運用へ揃えた

## Outcome

- orchestrate / design / review / implement の語彙を、`design bundle` 必須前提から個別 artifact 前提へ更新した
- `docs/exec-plans/templates/work-plan.md` に `final_mock_path`、`final_scenario_path`、`canonicalization_targets`、`canonicalized_artifacts` を追加した
- active / completed README を、artifact 種別ごとの個別昇格ルールとして読み直せる形に更新した
- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
