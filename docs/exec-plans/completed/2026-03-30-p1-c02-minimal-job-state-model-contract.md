- workflow: impl
- status: completed
- lane_owner: codex
- scope: Define the minimal Phase 1 job state contract and transition policy shared by create and list paths.
- task_id: P1-C02
- task_catalog_ref: docs/tasks/phase-1/tasks/P1-C02.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-C02` to define one stable minimal job state model.
- Fix the boundary used by backend job create, backend job list, and the first import-to-job acceptance anchor.
- Keep the task inside job state domain and job DTO scope without pulling in paused, failure recovery, cancellation, or provider execution controls.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tasks/phase-1/tasks/P1-C02.yaml`
- `docs/tasks/phase-1/tasks/P1-I04.yaml`
- `docs/tasks/phase-1/tasks/P1-I05.yaml`
- `docs/tasks/phase-1/tasks/P1-I06.yaml`
- `docs/tasks/phase-1/tasks/P1-I07.yaml`
- `docs/tasks/phase-1/tasks/P1-V02.yaml`

## Owned Scope

- `src-tauri/src/domain/job_state/`
- `src-tauri/src/application/dto/job/`

## Out Of Scope

- paused, canceled, recoverable failed, and failed state contracts
- provider execution controls
- UI screen wiring
- persistence or query logic outside the shared state contract boundary

## Dependencies / Blockers

- `P1-I04`, `P1-I05`, `P1-I06`, `P1-I07`, and `P1-V02` consume this contract and should not start with a drifting state boundary.
- The contract must stay aligned with the normal-path state progression described in `docs/spec.md` section 7.

## Parallel Safety Notes

- The task is parallel-safe with `P1-C01`, `P1-V01`, and `P1-V02`, but widening scope into future recovery or control states would block the Phase 1 batch.
- DTO naming or serialized shape drift would break the shared contract expected by backend and frontend consumers.

## UI

- No new screen, interaction, or UI wiring is introduced in this task.
- Downstream create and list screens should consume one shared backend state field only and must not require provider control flags, progress detail, or transition metadata in this contract.
- Successful create and list observation should rely on the same serialized state names, while `Draft` remains a backend creation-boundary concern and does not add a Phase 1 UI case in this task.

## Scenario

- Create path: a new job starts as `Draft` inside the shared domain contract, then must transition to `Ready` before the create path returns a successful observable result.
- List path: job list consumers observe the same contract without remapping and only need the minimal normal-path states already fixed by the domain boundary; a leaked `Draft` job is a contract violation, not a new list scenario.
- Execution follow-up tasks can advance the same shared contract from `Ready` to `Running` and from `Running` to `Completed` without introducing new state kinds or alternate Phase 1 branches.
- Pause, cancel, and failure scenarios remain explicitly outside this minimal Phase 1 contract.

## Logic

- `src-tauri/src/domain/job_state/` owns the single Phase 1 source of truth for the minimal job-state model and transition policy.
- The minimal contract contains exactly four states: `Draft`, `Ready`, `Running`, and `Completed`.
- Allowed transitions are exactly `Draft -> Ready`, `Ready -> Running`, and `Running -> Completed`.
- Reverse moves, skip moves, self-transitions, and any transition involving out-of-scope future states are invalid in this contract; `Completed` is terminal.
- `src-tauri/src/application/dto/job/` should expose one shared wire representation that maps 1:1 to the domain state names so create and list paths do not redefine enums or string literals independently.
- The DTO boundary for this task carries only the minimal job state contract needed by create and list. Execution controls, pause or cancel affordances, failure details, and persistence-specific fields stay outside this shared boundary.

## Implementation Plan

- Ordered scope 1: Add `src-tauri/src/domain/job_state/` as the backend-owned minimal contract module, define the single Phase 1 job-state enum there with exactly `Draft`, `Ready`, `Running`, and `Completed`, and keep the transition policy in the same domain module so only `Draft -> Ready`, `Ready -> Running`, and `Running -> Completed` succeed while reverse, skip, self, and any terminal exit from `Completed` fail locally.
- Ordered scope 2: Add `src-tauri/src/application/dto/job/` as the shared wire contract module, mirror the domain state names 1:1 in a job-state DTO, add domain-to-DTO conversion there, and update `src-tauri/src/application/dto/mod.rs` so later create and list paths consume one exported DTO boundary instead of redefining state literals.
- Ordered scope 3: Update `src-tauri/src/domain/mod.rs` and add the minimal coverage that proves the contract from the crate boundary: domain unit tests should own transition-policy checks inside `src-tauri/src/domain/job_state/`, and a new `src-tauri/tests/job_state_contract.rs` integration test should prove the public DTO serialization stays exactly `Draft`, `Ready`, `Running`, and `Completed` for downstream backend consumers.

## Acceptance Checks

- A minimal job state domain contract exists under `src-tauri/src/domain/job_state/`.
- A shared job DTO shape aligned with the minimal state contract exists under `src-tauri/src/application/dto/job/`.
- Tests or acceptance coverage prove the allowed minimal transitions and reject scope creep into future states.
- Backend create and list tasks can depend on the contract without reopening Phase 1 state decisions.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests proving the minimal state shape and transition policy.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src-tauri/src/domain/job_state/` with the minimal Phase 1 contract containing exactly `Draft`, `Ready`, `Running`, and `Completed`.
- Implemented the local transition policy so only `Draft -> Ready`, `Ready -> Running`, and `Running -> Completed` succeed, while reverse, skip, self, and terminal-exit transitions fail with a backend-owned error value.
- Added `src-tauri/src/application/dto/job/` with a 1:1 DTO mirror and domain-to-DTO conversion for the shared wire boundary.
- Updated `src-tauri/src/domain/mod.rs` and `src-tauri/src/application/dto/mod.rs` additively so downstream create and list work can import the shared contract without disturbing parallel `P1-C01` changes.
- Added contract coverage in `src-tauri/tests/job_state_contract.rs` plus local unit tests in `src-tauri/src/domain/job_state/mod.rs`.
- Validation passed: `powershell -File scripts/harness/run.ps1 -Suite structure`, `cargo test --test job_state_contract --all-features`, `cargo test --all-features`, `powershell -File scripts/harness/run.ps1 -Suite execution`, `powershell -File scripts/harness/run.ps1 -Suite all`, and owned-path Sonar open issues = `0`.
- Single-pass review result: `pass`.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not changed because this task reduced local uncertainty with new contract tests but did not materially change the current repository-level quality posture or debt list.
