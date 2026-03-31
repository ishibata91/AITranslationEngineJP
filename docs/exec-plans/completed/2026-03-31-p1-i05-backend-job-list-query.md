- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement the backend job list query that returns the minimal observable Phase 1 job list without absorbing job creation, execution, or UI scope.
- task_id: P1-I05
- task_catalog_ref: tasks/phase-1/tasks/P1-I05.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-I05` to expose one backend query path for the first minimal job list and status view.
- Keep the work inside backend job list application and domain query scope.
- Reuse the completed minimal job-state contract and stay aligned with the import-to-job acceptance anchor created by Phase 1 backend work.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `tasks/phase-1/tasks/P1-I05.yaml`
- `tasks/phase-1/tasks/P1-C02.yaml`
- `tasks/phase-1/tasks/P1-V02.yaml`
- `docs/exec-plans/completed/2026-03-30-p1-c02-minimal-job-state-model-contract.md`
- `docs/exec-plans/completed/2026-03-31-p1-i04-backend-job-creation-usecase.md`

## Owned Scope

- `src-tauri/src/application/job/list/`
- `src-tauri/src/domain/job/list/`

## Out Of Scope

- job creation mutation path
- execution control or provider progress detail
- UI screen wiring
- persistence work outside the minimal list query boundary unless a local query port or fixture adapter is required

## Dependencies / Blockers

- `P1-C02` minimal Phase 1 job state contract is required for the observable list result boundary.
- `P1-I04` create path establishes the first backend-generated job shape that list work is expected to observe.
- `P1-V02` import-to-job acceptance anchor defines the first backend path that list visibility should stay compatible with.

## Parallel Safety Notes

- `P1-I04`, `P1-I06`, and `P1-I07` are marked parallel-safe, but widening scope into create internals, execution control, or UI state would collapse the Phase 1 split.
- Shared DTO names and crate public exports should be added conservatively so frontend job list work can consume the query path without reopening the state contract.

## UI

- No new screen, route, or interaction is introduced in this task.
- Backend-visible output for this task should stay limited to one list item shape that the first UI list can observe without inventing a parallel contract: backend-generated `job_id` plus the shared Phase 1 job state only.

## Scenario

- Success flow is `saved created jobs -> backend list query -> minimal observable job list result`, where each saved job becomes exactly one list item with `job_id` and the shared observable state.
- The list query must expose the same serialized Phase 1 job-state names already fixed for create and later UI work. `Draft` is not a valid observable list state for this task.
- The first observable list path only needs to prove that a job created through the backend-owned create boundary can be read back through the list boundary without reshaping source provenance, phase progress, provider metadata, or execution-control affordances into the result.
- Empty persistence state should return an empty job list successfully instead of treating the absence of saved jobs as an error path.

## Logic

- `src-tauri/src/domain/job/list/` should own the list-specific minimal observable snapshot and reuse `crate::domain::job_state::JobState` directly instead of redefining a list-only state contract.
- The domain list snapshot should carry only the backend-generated `job_id` and shared observable `JobState`. Stored create-time `source_groups` remain persistence-side evidence for provenance but are not part of the list result contract.
- `src-tauri/src/application/job/list/` should own the query usecase and a list-local read port that returns saved backend jobs or list snapshots from the list boundary only. The list slice must not depend on create-usecase internals beyond stable public contracts needed to observe saved jobs.
- `src-tauri/src/application/dto/job/` should extend the shared job DTO boundary with one minimal list item DTO and one list result DTO that serialize only `job_id` plus the shared `JobStateDto`, so frontend list work can consume the query path without reopening the Phase 1 state decision.

## Implementation Plan

### Domain list slice

- Ordered scope 1: Add `src-tauri/src/domain/job/list/` as the list-owned observation module and define only the minimal observable job snapshot there: backend-generated `job_id` plus shared `JobState`. Reuse `crate::domain::job_state::JobState` directly and keep create-only `source_groups` out of the list result contract.
- Ordered scope 2: Add local domain mapping coverage that proves list observation accepts the minimal created-job state already produced by Phase 1 create work and does not expose `Draft` through a successful observable path.

### Application list slice

- Ordered scope 3: Add `src-tauri/src/application/job/list/` as the query entry for this usecase, with a list-local read port that returns saved backend jobs or equivalent list snapshots needed for observation. Keep the port local to the list slice and avoid create mutation behavior, execution shaping, or provider-facing detail.
- Ordered scope 4: Extend `src-tauri/src/application/dto/job/` with one minimal list item DTO and one list result DTO that serialize only `job_id` plus shared `JobStateDto`, then export the new list slice and DTOs through the existing public roots without widening the crate boundary.

### Tests and validation

- Ordered scope 5: Add backend contract coverage under `src-tauri/tests/` for the list usecase and update `src-tauri/tests/acceptance/import-job/` so one fixture-backed scenario proves `import -> create -> list visibility` through backend-owned paths without UI coupling.
- Validation commands:
  - `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test job_list_contract -- --nocapture`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test acceptance -- --nocapture`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
  - `powershell -File scripts/harness/run.ps1 -Suite all`

## Acceptance Checks

- A backend-owned list query path exists that returns the minimal observable Phase 1 job list.
- The list result reuses the shared job-state contract without redefining state names or exposing create-only `Draft`.
- Backend acceptance or contract coverage proves the list path can observe at least the first created job shape without UI coupling.
- Empty stored state returns an empty list successfully.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests proving the minimal job list query and its observable result boundary.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src-tauri/src/domain/job/list/` with the minimal observable `ListedJob` snapshot carrying only `job_id` and shared `JobState`, plus local validation that rejects blank identifiers and observable `Draft`.
- Added `src-tauri/src/application/job/list/` with `ListJobsUseCase` and a list-local repository port that returns saved observable jobs without absorbing create mutation or execution detail.
- Extended `src-tauri/src/application/dto/job/` with `JobListItemDto` and `ListJobsResultDto`, reusing the shared `JobStateDto` so the list wire contract stays `job_id + state`.
- Added `src-tauri/tests/job_list_contract.rs` and extended `src-tauri/tests/acceptance/import-job/mod.rs` so backend coverage now proves both the minimal list contract and the fixture-backed `import -> create -> list visibility` path.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because the task added local evidence without changing repository-level debt posture.
- Validation passed: `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test job_list_contract -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test acceptance -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `powershell -File scripts/harness/run.ps1 -Suite all`, and the owned-scope Sonar open-issue gate returned `openIssueCount: 0`.
- Single-pass review result: `pass`.
