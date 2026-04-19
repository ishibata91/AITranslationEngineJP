# Task Plan: 2026-04-19-legacy-schema-ui-migration-todo

- `workflow`: propose-plans
- `status`: todo-stub
- `lane_owner`: Codex designs later. GitHub Copilot implements only after a future approved implementation-scope.
- `task_id`: `2026-04-19-legacy-schema-ui-migration-todo`
- `task_mode`: backend-frontend-migration
- `request_summary`: 新 canonical schema / repository 作成後に、旧 master 名 schema / repository / service / UI を移行するための TODO plan stub。
- `goal`: 旧 master 名 dictionary / persona の既存資産を canonical schema / repository へ接続し直す後続 task を忘れないようにする。
- `constraints`: この stub では requirements / scenario / implementation-scope を作らない。新 schema / repository 作成 plan の完了後に具体化する。
- `close_conditions`: 後続で正式な requirements-design、ui-design、scenario-design、implementation-scope を作るか、不要判断で close する。

## Artifact Index

- `requirements_design`: `TODO`
- `ui_design`: `TODO`
- `scenario_design`: `TODO`
- `diagramming`: `N/A`
- `implementation_scope`: `TODO`

## TODO Scope Notes

- 旧 `master_dictionary_entries` / `master_persona_entries` 依存を洗い出す。
- 旧 master 名 dictionary / persona repository を削除、置換、adapter 化のどれにするか決める。
- 旧 master 名 dictionary UI を `DICTIONARY_ENTRY` 系 read / command model に接続する。
- 旧 master 名 persona UI を `PERSONA` 系 read / command model に接続する。
- 旧 master 名 persona の AI settings / run status を canonical ER 外の操作状態として残すか、新 schema を足すか判断する。

## Routing Notes

- `required_reading`: `../2026-04-19-sqlite-migration-repositories/`, `internal/controller/wails/master_*`, `internal/service/master_*`, `internal/bootstrap/app_controller.go`, `internal/repository/master_*`, `frontend/src/**/master-*`.
- `source_diagram_targets`: `docs/diagrams/er/combined-data-model-er.d2`.
- `canonicalization_targets`: `N/A` until this TODO becomes an approved task.
- `validation_commands`: `TODO`.

## HITL Status

- `functional_or_design_hitl`: `required-before-design-bundle`
- `approval_record`: `todo-stub-only`

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
- `docs_changes`: N/A

## Closeout Notes

- `canonicalized_artifacts`: N/A
- `follow_up`: 新 schema / repository plan 完了後に具体化する。

## Outcome

- TODO stub only. No implementation handoff exists yet.
