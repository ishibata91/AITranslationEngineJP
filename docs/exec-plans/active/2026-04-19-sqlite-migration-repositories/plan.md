# Task Plan: 2026-04-19-sqlite-migration-repositories

- `workflow`: propose-plans
- `status`: design-bundle-draft
- `lane_owner`: Codex が handoff を設計し、人間承認後に GitHub Copilot が実装する。
- `task_id`: `2026-04-19-sqlite-migration-repositories`
- `task_mode`: backend-persistence
- `request_summary`: ER ベースで SQLite の canonical schema migration と repository を作成する計画を作る。
- `goal`: 統合 ER を SQLite migration と repository contract / implementation / test に落とす。既存 service、bootstrap、UI への接続変更は含めない。
- `constraints`: product code はこの plan 作成では変更しない。実装 scope でも旧 master 名 schema / repository / service / UI の修正、削除、切替は行わない。旧資産の移行は別 TODO plan に分ける。
- `close_conditions`: design bundle が active plan に存在する。human review 後に implementation-scope を migration + repository 作成の Copilot handoff として使える。

## Artifact Index

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `diagramming`: `N/A`
- `implementation_scope`: `./implementation-scope.md`

## Routing Notes

- `required_reading`: `docs/er.md`, `docs/diagrams/er/combined-data-model-er.d2`, existing SQLite startup code, existing repository patterns.
- `source_diagram_targets`: `docs/diagrams/er/combined-data-model-er.d2`, `docs/diagrams/er/combined-data-model-er.svg`.
- `canonicalization_targets`: `N/A` for this planning step. Implementation completion may require docs update only after human approval.
- `validation_commands`: required `go test ./internal/infra/sqlite ./internal/repository`; required `python3 scripts/harness/run.py --suite structure`; optional regression `go test ./internal/...`.

## Repository Boundary

- `TranslationSourceRepository`: 翻訳入力元。`X_EDIT_EXTRACTED_DATA`、`TRANSLATION_RECORD`、`NPC_RECORD`、`TRANSLATION_FIELD`、field reference、`NPC_PROFILE` link を扱う。
- `FoundationDataRepository`: 翻訳補助データ。`PERSONA`、`PERSONA_FIELD_EVIDENCE`、`DICTIONARY_ENTRY`、XML import provenance を扱う。
- `JobLifecycleRepository`: job / phase の state machine。`TRANSLATION_JOB`、`JOB_PHASE_RUN`、`PHASE_RUN_*` を扱う。
- `JobOutputRepository`: DB 上の翻訳成果状態。`JOB_TRANSLATION_FIELD` を扱う。
- `TranslationFieldDefinitionRepository`: アプリ固定の lookup / seed data を扱う。
- 実ファイル export と `TranslationArtifactRepository` はこの plan の repository 境界に含めない。

## HITL Status

- `functional_or_design_hitl`: `required-after-design-bundle`
- `approval_record`: `pending-after-design-bundle`

## Related TODO Plan

- `legacy_schema_ui_migration`: `../2026-04-19-legacy-schema-ui-migration-todo/plan.md`
- `purpose`: 新 schema / repository 作成後に、旧 master 名 schema / repository / service / UI を canonical schema へ寄せる作業を別 task として設計する。

## Copilot Result

- `completed_handoffs`: pending
- `touched_files`: pending
- `implemented_scope`: pending
- `test_results`: pending
- `implementation_investigation`: pending
- `ui_evidence`: N/A
- `implementation_review_result`: pending
- `sonar_gate_result`: pending
- `residual_risks`: pending
- `docs_changes`: pending. Copilot must not change docs.

## Closeout Notes

- `canonicalized_artifacts`: pending after implementation review.
- `follow_up`: 旧 schema / repository / service / UI の移行は TODO plan で扱う。

## Outcome

- Active plan was narrowed to canonical SQLite migration and repository creation only.
