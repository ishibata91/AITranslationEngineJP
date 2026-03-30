- workflow: impl
- status: completed
- lane_owner: codex
- scope: Define the canonical TRANSLATION_UNIT contract and DTO/result shape before job creation and UI observation start.
- task_id: P1-C01
- task_catalog_ref: docs/tasks/phase-1/tasks/P1-C01.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-C01` to define one stable canonical `TRANSLATION_UNIT` contract.
- Fix the boundary used after raw plugin export import and before later job creation and UI observation tasks.
- Keep the task inside application DTO and domain translation-unit scope without absorbing job persistence or job UI wiring.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/er.md`
- `docs/tech-selection.md`
- `docs/tasks/phase-1/tasks/P1-C01.yaml`
- `docs/exec-plans/completed/2026-03-29-per-39-xedit-export-json-importer.md`
- `docs/exec-plans/completed/2026-03-29-per-41-plugin-export-persistence-unit.md`

## Owned Scope

- `src-tauri/src/application/dto/translation_unit/`
- `src-tauri/src/domain/translation_unit/`

## Out Of Scope

- job persistence
- job UI wiring

## Dependencies / Blockers

- `P1-I03` translation-unit persistence groundwork must already be available as an upstream dependency.
- Existing importer and persistence behavior from completed Phase 1 tasks should remain compatible while this task fixes the canonical boundary.

## Parallel Safety Notes

- `P1-C02`, `P1-V01`, and `P1-V02` are marked parallel-safe in the task catalog, but helper-only fields or fixture schema drift could collapse that safety if this contract grows beyond owned scope.
- The task must not redefine raw `PLUGIN_EXPORT` storage truth or later job result shape beyond the canonical translation-unit boundary.

## UI

- No new screen, route, or interaction is introduced in this task.
- Later UI observation should see `TRANSLATION_UNIT` through one backend-owned DTO module under `src-tauri/src/application/dto/translation_unit/`, not through importer-local DTO definitions.
- The UI-observable unit shape for this task is limited to stable translation-unit fields: `source_entity_type`, `form_id`, `editor_id`, `record_signature`, `field_name`, `extraction_key`, `source_text`, and `sort_key`.
- Raw-record payload, raw-record helper metadata, persistence identifiers such as `plugin_export_id` or `source_entity_id`, and any job/progress fields remain outside this DTO boundary.

## Scenario

- Success flow is `validated raw plugin export import -> derive canonical domain translation units once -> serialize the same units through the canonical translation-unit DTO/result shape before any job creation starts`.
- Each non-blank translatable source field yields one canonical `TRANSLATION_UNIT`; blank or absent source text still yields no unit, while blank `editor_id` remains preserved losslessly.
- Later job-create, job-list, and verification tasks should be able to consume the same canonical unit shape without depending on importer-local structs, raw-record payloads, or storage-only identifiers.
- The canonical shape must remain lossless for the `docs/er.md` translation-unit fields needed for xTranslator reconstruction and deterministic downstream ordering, without leaking importer-private helper fields.

## Logic

- `TRANSLATION_UNIT` remains the canonical normalized unit created from translatable raw fields at import time and used as the stable boundary before later workflow stages.
- The backend-owned domain contract should move under `src-tauri/src/domain/translation_unit/` and define the canonical field set as `source_entity_type`, `form_id`, `editor_id`, `record_signature`, `field_name`, `extraction_key`, `source_text`, and `sort_key`.
- `extraction_key` is part of the canonical contract because it identifies the exact translatable field within the source-side record graph; it must not be treated as an importer-local helper.
- `sort_key` is also canonical and must be deterministic from source-side ordering inputs such as nested indices and field identity so later UI observation, fixtures, and job targeting can rely on stable ordering without redesign.
- Storage-owned identifiers such as `TRANSLATION_UNIT.id`, `plugin_export_id`, and `source_entity_id` remain outside the pre-persistence canonical DTO/domain boundary for this task. Later persistence work should attach those identifiers by mapping from raw records plus the canonical unit contract rather than widening the canonical shape now.
- The application DTO under `src-tauri/src/application/dto/translation_unit/` should mirror the canonical unit fields directly, and importer result mapping should depend on the translation-unit module instead of `domain::xedit_export::TranslationUnit`.

## Implementation Plan

- Ordered scope 1: Add `src-tauri/src/domain/translation_unit/` as the backend-owned canonical contract module, move `TranslationUnit` there with the full field set including `sort_key`, and keep `ImportedPluginExport` / `ImportedRawRecord` in `src-tauri/src/domain/xedit_export.rs` as consumers of that contract rather than duplicate owners.
- Ordered scope 2: Add `src-tauri/src/application/dto/translation_unit/` as the canonical wire module, mirror the domain field set directly there, and update `src-tauri/src/application/dto/import_xedit_export_dto.rs` so importer result mapping depends on the translation-unit DTO module instead of importer-local struct definitions.
- Ordered scope 3: Update `src-tauri/src/infra/xedit_export_importer.rs` to construct canonical translation units from the new domain module, preserve blank `editor_id`, continue dropping blank or absent `source_text`, and set the first backend-owned `sort_key` from the same deterministic extraction identity currently used for `extraction_key` unless an equally local and safer ordering rule is found from existing importer inputs during implementation.
- Ordered scope 4: Update `src-tauri/tests/xedit_export_importer.rs` plus affected unit coverage so the canonical boundary proves `sort_key` is present and deterministic, the canonical field set stays lossless, and compatibility re-exports through `xedit_export` are added only if compilation shows downstream coupling outside the new module boundary.

## Acceptance Checks

- A canonical translation-unit domain contract exists under `src-tauri/src/domain/translation_unit/`.
- A canonical translation-unit DTO and result mapping exists under `src-tauri/src/application/dto/translation_unit/`.
- Import-related result coverage proves the stable shape preserves `source_entity_type`, `form_id`, `editor_id`, `record_signature`, `field_name`, `extraction_key`, `source_text`, and `sort_key` losslessly.
- Import-related result coverage proves blank `editor_id` remains preserved and deterministic `sort_key` is emitted from source-side extraction identity inputs.
- Later Phase 1 job and verification tasks can depend on the contract without redesign or importer-local type coupling, and importer result mapping depends on the canonical translation-unit modules rather than importer-local structs.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests or fixtures proving the canonical translation-unit shape is stable and lossless.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src-tauri/src/domain/translation_unit/` as the canonical backend `TRANSLATION_UNIT` contract and included canonical `sort_key` alongside the existing lossless source fields.
- Added `src-tauri/src/application/dto/translation_unit/` as the canonical DTO surface and moved importer result mapping to that module instead of importer-local struct ownership.
- Updated `src-tauri/src/domain/xedit_export.rs` and `src-tauri/src/infra/xedit_export_importer.rs` so imported plugin exports consume the canonical translation-unit contract, preserve blank `editor_id`, skip blank source text, and set `sort_key` from the same deterministic extraction identity as `extraction_key`.
- Added contract coverage in `src-tauri/tests/translation_unit_contract.rs`, strengthened `src-tauri/tests/xedit_export_importer.rs`, and added a reroute follow-up unit test that pins blank `sort_key` rejection in `src-tauri/src/domain/translation_unit/mod.rs`.
- Validation passed on the current working tree: `powershell -File scripts/harness/run.ps1 -Suite structure`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo test --manifest-path ./src-tauri/Cargo.toml --tests`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `powershell -File scripts/harness/run.ps1 -Suite all`, and `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths ...` returned `openIssueCount: 0`.
