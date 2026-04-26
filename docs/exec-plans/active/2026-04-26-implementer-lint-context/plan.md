# Task Plan: 2026-04-26-implementer-lint-context

- `workflow`: docs-update
- `status`: in-progress
- `lane_owner`: Codex が Copilot 実装 skill の入力規約を更新する
- `task_id`: `2026-04-26-implementer-lint-context`
- `task_mode`: workflow-guidance
- `request_summary`: Copilot implementer が lint gate の規約を知らずに止まりやすいため、事前入力へ lint と architecture boundary を追加する
- `goal`: implementer と implementation-distiller が、handoff ごとに lint-sensitive な規約と境界を明示的に扱う
- `constraints`: product code、product test、docs 正本は変更しない。`.github/skills` と task-local plan だけを更新する
- `close_conditions`: implementer 系 skill が lint policy、architecture boundary、touched-layer validation の事前確認を明記している

## Artifact Index

- `requirements_design`: `N/A`
- `ui_design`: `N/A`
- `scenario_design`: `N/A`
- `diagramming`: `N/A`
- `implementation_scope`: `N/A`

## Workflow State

- `distiller`: `manual-read`; implementation-distill-implement の現行出力粒度を確認済み
- `designer`: `N/A`
- `investigator`: `not-needed`
- `human_review_gate`: `not-required`

## Routing Notes

- `required_reading`: `.github/skills/implementation-distill-implement/SKILL.md`、`.github/skills/implement/SKILL.md`、`.github/skills/implement-frontend/SKILL.md`、`.github/skills/implement-backend/SKILL.md`、`docs/lint-policy.md`、`docs/architecture.md`
- `canonicalization_targets`: `N/A`
- `validation_commands`: `N/A`

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: `approved-by-user-2026-04-26`
