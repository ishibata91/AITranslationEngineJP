- workflow: impl
- status: completed
- lane_owner: codex
- scope: Validate imported xEdit export input data before importer accepts it and add failure-path test coverage.

## Request Summary

- Implement Linear `PER-40` `[Phase 1][02] Validate imported input data`.
- Validate imported input data structure and required fields at import time so invalid xEdit export input is rejected early.
- Define the importer failure path and error message policy.
- Add corresponding tests.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- Linear `PER-40`
- `docs/exec-plans/completed/2026-03-29-per-39-xedit-export-json-importer.md`

## UI

- UI wiring is out of scope for `PER-40`.
- The existing backend command boundary remains the entrypoint, and validation errors must stay representable through the current `Result<..., String>` shape without adding UI-specific error handling.

## Scenario

- The import scenario remains whole-request import through `file_paths: string[]`.
- Validation runs per input file before an `ImportedPluginExport` is accepted, and the first invalid file aborts the whole request.
- Structural parse failure, missing top-level required data, and missing required record identifiers on a translatable record are all importer failures for this task.
- Blank or absent translatable text continues to mean "no translation unit emitted" rather than a validation error.
- `editor_id` stays lossless and may remain blank when the export provides no value.
- Validation failures must return a deterministic message that names the violated importer precondition closely enough for test assertions and later UI display; partial success stays out of scope.

## Logic

- Validation stays backend-owned along the existing `gateway/commands.rs -> application use case -> infra importer -> domain model` flow.
- `PER-39` already proves JSON parsing and minimum `target_plugin` / `translation_units` preservation. `PER-40` extends that boundary to explicit imported-input validation without changing the command contract.
- Required imported fields for accepted translation data are `target_plugin`, record identity fields needed to build a `TranslationUnit`, and nested indices that participate in extraction keys. Missing or blank required values should fail before the plugin export is returned.
- Optional text-bearing fields remain optional at the raw JSON layer. The importer still skips blank strings instead of manufacturing empty `translation_units`.
- Error policy stays string-based for this task. Messages should remain deterministic and preferably include file path plus a stable field path or precondition label, but should not introduce a new error enum or UI-facing DTO.

## Implementation Plan

- Ordered scope 1: Pin importer-layer validation to malformed or incomplete xEdit export input before `ImportedPluginExport` is returned, while keeping the existing `file_paths: string[]` command contract and whole-request failure semantics.
- Ordered scope 2: Add or tighten Rust tests in importer-owned files for structural parse failure and missing required imported fields, while preserving valid import behavior, blank-text skip behavior, and blank `editor_id` preservation.
- Ordered scope 3: Implement backend validation in `src-tauri/src/infra/xedit_export_importer.rs` and only the smallest supporting domain or application adjustments needed to keep deterministic string errors.
- Ordered scope 4: Run backend validation commands, execution harness, `sonar-scanner`, then hand the resulting diff to single-pass implementation review before closeout.
- Owned scope: `src-tauri/src/infra/xedit_export_importer.rs`, `src-tauri/src/domain/xedit_export.rs`, importer-related tests under `src-tauri/src/infra/` and `src-tauri/src/domain/`, and only importer boundary files in `src-tauri/src/application/` or `src-tauri/src/gateway/` if a supporting change is required.
- Required reading: `docs/exec-plans/active/2026-03-29-per-40-validate-imported-input-data.md`, `docs/exec-plans/completed/2026-03-29-per-39-xedit-export-json-importer.md`, `docs/spec.md`, `docs/architecture.md`, `docs/tech-selection.md`, `src-tauri/src/application/dto/import_xedit_export_dto.rs`, `src-tauri/src/application/importer/import_xedit_export.rs`, `src-tauri/src/infra/xedit_export_importer.rs`, `src-tauri/src/domain/xedit_export.rs`, `src-tauri/src/gateway/commands.rs`.
- Validation commands: `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `powershell -File scripts/harness/run.ps1 -Suite execution`.

## Acceptance Checks

- A Rust test proves structurally invalid xEdit export JSON fails the whole import request with a validation error.
- A Rust test proves an input file that omits a required imported field fails before a plugin export is accepted.
- Existing valid-import coverage still passes without changing the command request contract from `file_paths: string[]`.
- Error messages for validation failures are deterministic enough to assert in tests.

## Required Evidence

- Active plan updated with task-local design, distill facts, and implementation brief.
- Added or updated tests / fixtures for invalid imported input.
- Validation command results, Sonar issue status, and single-pass review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added importer failure-path tests for structurally invalid top-level collection data and blank `quests[].objectives[].objective_index` when objective text is present.
- Updated the xEdit export importer to surface deterministic parse and validation errors with file-path context, and to reject blank quest objective indices before returning an imported plugin export.
- Preserved the existing backend command contract `file_paths: string[]`, whole-request failure semantics, blank-text skip behavior, and blank `editor_id` preservation.
- Validation passed: `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `powershell -File scripts/harness/run.ps1 -Suite execution`, and `powershell -File scripts/harness/run.ps1 -Suite all`.
- `sonar-scanner` uploaded successfully. SonarQube MCP was unavailable due to transport errors, so SonarCloud public API was used to confirm `0` open issues for `src-tauri/src/infra/xedit_export_importer.rs`; project quality gate is currently `NONE`.
- Single-pass implementation review returned `pass`.
