- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement the backend job creation usecase from normalized translation units without absorbing job list or UI scope.
- task_id: P1-I04
- task_catalog_ref: tasks/phase-1/tasks/P1-I04.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-I04` to create jobs from normalized translation units through one backend application path.
- Keep the work inside backend job creation usecase and domain creation scope.
- Reuse the completed Phase 1 translation-unit and job-state contracts and stay aligned with the import-to-job acceptance anchor.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/coding-guidelines.md`
- `tasks/phase-1/tasks/P1-I04.yaml`
- `tasks/phase-1/tasks/P1-V02.yaml`
- `docs/exec-plans/completed/2026-03-30-p1-c01-translation-unit-canonical-contract.md`
- `docs/exec-plans/completed/2026-03-30-p1-c02-minimal-job-state-model-contract.md`

## Owned Scope

- `src-tauri/src/application/job/create/`
- `src-tauri/src/domain/job/create/`

## Out Of Scope

- job list query
- UI code
- provider execution
- persistence or query work outside the create path unless minimal ports or fixture adapters are required for this usecase boundary

## Dependencies / Blockers

- `P1-C01` canonical `TRANSLATION_UNIT` contract is required as the stable input boundary.
- `P1-C02` minimal Phase 1 job state contract is required for the create result boundary.
- `P1-V02` acceptance anchor defines the first executable import-to-job path and must remain satisfiable.

## Parallel Safety Notes

- `P1-I05`, `P1-I06`, and `P1-I07` are marked parallel-safe, but widening scope into job list, UI, or execution control would collapse the Phase 1 split.
- Shared DTO names and crate public exports should be added conservatively so parallel tasks can consume the create path without reopening the job-state or translation-unit contracts.

## UI

- No new screen, route, or interaction is introduced in this task.
- Backend-visible output for this task should stay limited to the minimal create result that downstream UI work can observe without inventing a parallel contract: one backend-generated job identifier plus the shared Phase 1 job state, with no provider controls, progress detail, or list-only fields added here.

## Scenario

- Success flow is `validated imported input with canonical translation units -> backend create request grouped by imported source -> domain create in Draft -> transition to Ready -> return one observable created job`.
- One create request may bundle one or more imported source groups into a single job so the Phase 1 path stays compatible with the spec requirement that one translation job can cover multiple input files without losing source provenance.
- Each source group contributes its canonical `TRANSLATION_UNIT` entries exactly once to the created job target set. The create path must preserve source grouping and deterministic unit ordering from the incoming canonical `sort_key`; duplicate `source_text` or repeated record signatures are not merge keys.
- A successful create result must expose the created job only after the internal `Draft -> Ready` transition completes. Returning `Draft` from the backend create boundary is a contract violation for this task.
- Requests with no source groups, or with zero canonical translation units across all source groups, fail locally in the create path instead of creating an empty `Ready` job.
- Provider execution, persona or dictionary binding, job list querying, and UI state handling remain outside this task even when the acceptance anchor later chains import, create, and visibility together.

## Logic

- `src-tauri/src/domain/job/create/` owns the create-specific rule set and minimal create-ready job model. It reuses `crate::domain::translation_unit::TranslationUnit` and `crate::domain::job_state::JobState` directly instead of redefining normalized unit or job-state contracts inside the create slice.
- The create input boundary should preserve imported-source provenance without coupling to importer-private raw-record structures. The create slice should therefore define its own source-group input carrying `source_json_path`, `target_plugin`, and `Vec<TranslationUnit>`, while excluding `raw_records`, raw payload helpers, and persistence-only identifiers.
- The application create usecase under `src-tauri/src/application/job/create/` should accept only already-normalized canonical translation units through that grouped input shape. Revalidation of raw xEdit payloads or reconstruction of translation units is out of scope for this path.
- Domain create logic should require at least one source group and at least one canonical translation unit overall, instantiate the new job in `Draft`, attach all provided source groups and unit targets, then perform the single create-boundary state move `Draft -> Ready`.
- The create result should expose a minimal created-job snapshot built from a backend-generated job identifier and the shared Phase 1 observable state `Ready`. Any default job naming, persistence internals, current-phase values, progress fields, or execution-control metadata stay internal or out of scope for this task.
- Any repository or output port added here must stay local to the create path and save or return the create-ready job plus its source and unit targeting relationship. Query-oriented shaping for job list work and execution-oriented shaping for later tasks must not be absorbed into this slice.

## Implementation Plan

### Domain create slice

- Ordered scope 1: Add `src-tauri/src/domain/job/create/` as the create-owned rule module and define only the local create contract there: a source-group input carrying `source_json_path`, `target_plugin`, and `Vec<TranslationUnit>`; a created-job snapshot carrying the backend-generated job identifier and shared `JobState`; and the local validation that rejects empty source groups and zero translation units overall without reopening canonical `TranslationUnit` or shared `JobState`.
- Ordered scope 2: Keep create-time behavior inside the same domain slice by instantiating the new job in `Draft`, attaching every source group and its canonical units in incoming grouped order, and completing the single allowed create-boundary transition to `Ready` before any success result can exist. Domain coverage in this slice should prove grouped provenance is preserved, empty inputs fail locally, and `Draft` never survives the successful create boundary.

### Application create slice

- Ordered scope 3: Add `src-tauri/src/application/job/create/` as the application entry for this usecase, with a create-local repository or output port that saves or returns the create-ready job snapshot plus its source and unit targeting relationship. Keep that port local to the create slice and avoid query-facing, execution-facing, or importer-private shaping.
- Ordered scope 4: Extend `src-tauri/src/application/dto/job/` with the create request and result DTOs needed by this slice, reusing canonical translation-unit DTO mapping and shared `JobStateDto` instead of introducing parallel contracts. The successful result should expose only the created job identifier and observable `Ready` state, while the grouped request shape preserves `source_json_path`, `target_plugin`, and canonical translation units.

### Public exports and acceptance anchor

- Ordered scope 5: Update `src-tauri/src/application/mod.rs`, `src-tauri/src/domain/mod.rs`, and the new slice `mod.rs` files so downstream tests can import the create public roots without reaching into internal modules. Do not widen this task into job list work, UI work, provider execution, or Tauri command wiring unless implementation proves the existing acceptance anchor cannot compile without that boundary.
- Ordered scope 6: Add the first backend-owned import-to-job coverage under `src-tauri/tests/acceptance/import-job/` and any narrowly scoped create-slice tests needed to prove the Phase 1 path: canonical grouped input creates one observable job only after `Draft -> Ready`, empty grouped input fails locally, and source provenance survives the create boundary.

## Acceptance Checks

- A backend-owned create path exists that accepts normalized translation units and returns one stable Phase 1 job result.
- Job creation aligns with the Phase 1 minimal job-state contract and does not leak `Draft` beyond the backend create boundary unless the contract explicitly requires it.
- The import-to-job acceptance anchor remains executable against the create path without UI coupling.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests and fixtures proving job creation from normalized translation units.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src-tauri/src/domain/job/create/` with the create-local source-group input, created-job snapshot, empty-input validation, and internal `Draft -> Ready` transition so successful create never leaks `Draft`.
- Added `src-tauri/src/application/job/create/` with `CreateJobUseCase`, a create-local repository port, and backend-generated `job-{n}` identifiers for the minimal Phase 1 create path.
- Extended `src-tauri/src/application/dto/job/` with create request and result DTOs that preserve grouped canonical translation-unit input and expose only `job_id` plus observable `Ready` on success.
- Added backend acceptance and contract coverage in `src-tauri/tests/job_create_contract.rs` and `src-tauri/tests/acceptance/import-job/mod.rs`, including failure-path evidence for repository save errors and malformed canonical DTO input.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this task added local evidence without changing the repository-level quality posture or debt inventory.
- Final validation passed: `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `powershell -File scripts/harness/run.ps1 -Suite all`, and the owned-scope Sonar open-issue gate returned `openIssueCount: 0`.
- Single-pass review initially returned `reroute` for missing create-usecase failure-path evidence; the required tests were added in the same lane, and per workflow no second review pass was run after the reroute fix landed.
