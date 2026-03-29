- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement a PLUGIN_EXPORT-equivalent persistence unit that stores imported xEdit export input without losing provenance.

## Request Summary

- Implement Linear `PER-41` `[Phase 1][03] Implement PLUGIN_EXPORT-equivalent persistence unit`.
- Add a persistence unit that can retain imported xEdit export data in execution cache form.
- Preserve input provenance and source data shape closely enough for later translation workflow stages.
- Add the minimum persistence-focused test coverage.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/er-draft.md`
- Linear `PER-41`
- `docs/exec-plans/completed/2026-03-29-per-39-xedit-export-json-importer.md`
- `docs/exec-plans/completed/2026-03-29-per-40-validate-imported-input-data.md`

## UI

- N/A for new screens or interaction changes.
- The existing backend command boundary remains the only entrypoint in this task, and persistence is added behind that boundary without introducing UI-specific DTO or state changes.
- If persistence needs an internal identifier for later phases, keep that identifier backend-owned for now rather than widening the frontend contract in `PER-41`.

## Scenario

- The success path for `PER-41` is `read validated xEdit export JSON -> map imported raw structures -> persist one execution-cache-owned PLUGIN_EXPORT unit plus its raw child records -> return success through the existing backend entrypoint`.
- The persisted `PLUGIN_EXPORT` unit must retain imported provenance at minimum through `target_plugin`, `source_json_path`, and an import timestamp, and the persisted child records must keep the source-side record identity needed to trace later normalization work back to the raw input.
- `PER-41` should persist raw imported input categories under `PLUGIN_EXPORT`; `PER-42` remains responsible for adding `TRANSLATION_UNIT`-equivalent canonical records, so this task must not treat the current in-memory translation-unit list as the persisted source of truth.
- Persistence failures should abort the whole request with deterministic backend-owned errors; partial persistence is out of scope.

## Logic

- Persistence stays backend-owned along `gateway/commands.rs -> application use case -> domain persistence model / port -> infra SQLite adapter`, keeping UI and permanent docs unchanged in this task.
- `PER-41` should introduce a `PLUGIN_EXPORT`-equivalent persistence model that can store the raw imported execution-cache graph described in `docs/er-draft.md`, beginning with the parent `PLUGIN_EXPORT` row and the minimum child tables needed to preserve imported source data categories without loss of provenance.
- `PER-41` must keep raw imported structures and canonical `TRANSLATION_UNIT` records conceptually separate. If current importer code already derives translation units in memory, that derivation may remain for compatibility, but persisted execution-cache truth for this task should be the raw imported `PLUGIN_EXPORT` graph rather than prematurely persisted `TRANSLATION_UNIT` rows.
- The implementation should use the repository-standard SQLite direction from `docs/tech-selection.md` and keep storage concerns behind an infra adapter so later phases can add normalization and job linkage without replacing the import boundary again.

## Implementation Plan

- Ordered scope 1: Expand the backend xEdit export domain model from translation-unit-only import results to a raw `PLUGIN_EXPORT` graph that preserves source-side record identity and can still derive the current translation-unit DTO output without making `TRANSLATION_UNIT` persistence the source of truth.
- Ordered scope 2: Add a backend-owned persistence port and SQLite adapter that create the minimum execution-cache schema for `PLUGIN_EXPORT` plus raw child record tables, persist imported graphs transactionally, and keep request failure whole-import on persistence errors.
- Ordered scope 3: Wire the existing import use case and Tauri command so `file_paths: string[]` still enters through the current backend boundary while successful imports are also written into the execution cache.
- Ordered scope 4: Add or update backend tests that prove one valid import persists a `PLUGIN_EXPORT` row plus representative raw child data and that provenance fields remain queryable after persistence.
- Owned scope: `src-tauri/Cargo.toml`, `src-tauri/src/application/dto/import_xedit_export_dto.rs`, `src-tauri/src/application/importer/`, `src-tauri/src/domain/xedit_export.rs`, `src-tauri/src/gateway/commands.rs`, `src-tauri/src/lib.rs` if async command registration changes are required, new or updated persistence modules under `src-tauri/src/infra/`, and importer / persistence tests under `src-tauri/tests/` or `src-tauri/src/infra/`.
- Required reading: `docs/exec-plans/active/2026-03-29-per-41-plugin-export-persistence-unit.md`, `docs/spec.md`, `docs/architecture.md`, `docs/tech-selection.md`, `docs/er-draft.md`, `docs/exec-plans/completed/2026-03-29-per-39-xedit-export-json-importer.md`, `docs/exec-plans/completed/2026-03-29-per-40-validate-imported-input-data.md`, `src-tauri/src/application/importer/import_xedit_export.rs`, `src-tauri/src/domain/xedit_export.rs`, `src-tauri/src/infra/xedit_export_importer.rs`, `src-tauri/src/gateway/commands.rs`, `src-tauri/Cargo.toml`.
- Validation commands: `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `powershell -File scripts/harness/run.ps1 -Suite execution`.

## Acceptance Checks

- A backend integration test proves a valid xEdit export import persists one `PLUGIN_EXPORT` execution-cache record with `target_plugin`, `source_json_path`, and non-empty `imported_at`.
- A backend integration test proves representative raw imported child data is persisted under the stored `PLUGIN_EXPORT` as raw-record storage, without making `TRANSLATION_UNIT` persistence the source of truth.
- A backend use-case test proves repository failure aborts the whole import request with a deterministic error instead of returning partial success.
- The import boundary remains `file_paths: string[]`, and valid import behavior observed in existing importer coverage stays intact while persistence is added behind the same backend entrypoint.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated persistence-focused tests / fixtures.
- Validation command results, Sonar issue status, and single-pass review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added backend-owned persistence for imported xEdit export input by extending `ImportedPluginExport` with provenance-preserving `raw_records` while keeping the frontend DTO output shape unchanged.
- Added a SQLite execution-cache adapter and application persistence port that store `plugin_exports` plus generic raw child records transactionally, preserving `target_plugin`, `source_json_path`, import time, and source-side record identity.
- Updated the Tauri import command to use the async command path and added command-boundary persistence failure coverage without reintroducing nested runtime bootstrapping.
- Added persistence-focused Rust tests for use-case success, use-case failure, command-boundary success, and command-boundary persistence failure.
- Validation passed: `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `sonar-scanner`, `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths ...`, and `powershell -File scripts/harness/run.ps1 -Suite all`.
- Updated `4humans/tech-debt-tracker.md` with the remaining debt that execution-cache SQLite location and retention policy are still temp-based and need formalization in a follow-up.
