- workflow: impl
- status: completed
- lane_owner: codex
- scope: Wire the Phase 1 import-to-job path through Tauri transport and app-shell composition without reopening backend or frontend contracts.
- task_id: P1-G01
- task_catalog_ref: tasks/phase-1/tasks/P1-G01.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-G01` to prove the first integrated import-to-job scenario by wiring the completed backend and frontend Phase 1 slices together.
- Keep the work inside Tauri command registration, frontend Tauri gateway adapters, and the import-to-job acceptance path.
- Reuse the completed create and list contracts without redesigning DTOs, fixtures, or screen-local state models.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-1/tasks/P1-G01.yaml`
- `tasks/phase-1/tasks/P1-V02.yaml`
- `docs/exec-plans/completed/2026-03-31-p1-i04-backend-job-creation-usecase.md`
- `docs/exec-plans/completed/2026-03-31-p1-i05-backend-job-list-query.md`
- `docs/exec-plans/completed/2026-03-31-p1-i06-job-create-screen.md`
- `docs/exec-plans/completed/2026-03-31-p1-i07-job-list-screen.md`

## Owned Scope

- `src-tauri/src/lib.rs`
- `src/gateway/tauri/`
- `src-tauri/tests/acceptance/import-job/`

## Out Of Scope

- new DTO design
- new fixture design
- provider execution
- dictionary or persona reuse
- frontend screen redesign outside the minimum composition needed to replace preview executors with Tauri transport wiring

## Dependencies / Blockers

- `P1-I04` supplies the backend create path and result boundary.
- `P1-I05` supplies the backend list query and observable result boundary.
- `P1-I06` and `P1-I07` supply frontend usecase and screen boundaries that are intended to become transport-backed here.
- `P1-V02` fixes the first agreed acceptance anchor and must stay executable after integration.

## Parallel Safety Notes

- The task must stay inside transport wiring and integrated verification so the earlier slice boundaries remain reusable.
- Shared app-shell composition files may need narrow additive updates, but the owned implementation focus remains `src/gateway/tauri/` and `src-tauri/src/lib.rs`.
- Reopening job DTO names, translation-unit contracts, or fixture shapes in this task would hide unresolved contract drift instead of proving integration.

## UI

- No new screen, route, or shell region is introduced in this task. `AppShell` keeps the existing single-window composition of job list, job create, and bootstrap status, and `src/main.ts` only swaps preview create/list executors for Tauri-backed adapters.
- The visible behavior of the existing screens stays unchanged: job create continues to show pending, success, and failure states; job list continues to show loading, empty, loaded, and retryable failure states. This task does not add optimistic rows, provider detail, or any new cross-screen layout.
- App-shell composition must stay additive and contract-preserving. A successful create does not introduce implicit list mutation or auto-refresh behavior in the shell; list visibility remains owned by the existing list initialize and refresh actions.
- Any transport-specific request wrapping, field-name mapping, or invoke error handling stays inside `src/gateway/tauri/` and must not surface into Svelte files or screen-local state models.

## Scenario

- Desktop startup keeps the current shell flow: bootstrap status loads through its existing Tauri gateway, job list initializes through a real Tauri list adapter, and job create remains idle until the user submits the existing request form.
- When the user submits the existing create form, the same request payload fixed by `P1-I06` is sent through Tauri transport to the backend create boundary. On success, the screen still renders only the returned `jobId` and observable `state`; on failure, the screen keeps the editable request and shows one generic user-facing error.
- Job visibility remains list-owned. After create succeeds, the integrated path is proven by calling the real backend-backed list initialize or refresh path rather than by introducing create-to-list coupling in the app shell.
- The acceptance anchor for this task is the transport-backed chain `fixture import -> job create -> job list visibility` under `src-tauri/tests/acceptance/import-job/`, reusing the agreed fixture path and the existing import, create, and list DTO shapes.

## Logic

- `src/gateway/tauri/` should add concrete job-create and job-list adapters that translate between the existing frontend usecase request/result shapes and Tauri `invoke` calls only. The adapters must not introduce new UI contracts, screen state, or domain rules.
- The create adapter must send the existing create request through the backend command's named request boundary and map the returned `jobId + state` result back into the unchanged frontend usecase contract. Transport-only serde derives or casing annotations on existing Rust DTOs are acceptable when required to make the unchanged DTO shape invokable.
- `src-tauri/src/gateway/commands.rs` should expose thin `create_job` and `list_jobs` commands that compose the completed backend usecases and return the existing DTOs. `src-tauri/src/lib.rs` should register those commands alongside the existing bootstrap and import commands so desktop composition and acceptance coverage use the same transport entrypoints.
- The create and list commands must share one concrete repository boundary so a job created through the transport path is observable from the subsequent list command in the same integrated scenario. The storage choice stays internal as long as the existing usecase traits, DTO names, and fixture contracts remain unchanged.
- `src-tauri/tests/acceptance/import-job/` should move from usecase-direct composition to transport-owned composition: import through the existing agreed import boundary, build the existing create request from imported plugin exports, call the create command, then call the list command and assert stable `jobId` plus `Ready` visibility without redesigning fixtures or DTOs.

## Implementation Plan

### 1. Acceptance anchor

- Update `src-tauri/tests/acceptance/import-job/` so the import-to-job scenario calls transport-owned command functions instead of composing create/list usecases directly.
- Keep the fixture path, import request, create request construction, and list assertions unchanged in shape; prove only that created jobs become visible through the same transport-backed repository.

### 2. Backend transport and shared store

- Add thin `create_job` and `list_jobs` commands in `src-tauri/src/gateway/commands.rs`, keeping usecase composition and invoke-facing error mapping inside the command boundary.
- Register both commands in `src-tauri/src/lib.rs`.
- Add the minimal concrete job repository/store under `src-tauri/src/infra/` so create and list share one backend-owned in-memory state within the integrated scenario.
- Add only transport-required serde derives or casing annotations to existing job DTOs if invoke cannot deserialize the unchanged frontend request shape.

### 3. Frontend Tauri adapters and shell composition

- Add concrete job-create and job-list adapters under `src/gateway/tauri/`, reusing the existing feature-screen helper and keeping invoke request wrapping inside the gateway layer.
- Replace preview executors in `src/main.ts` with the new Tauri adapters without adding create-to-list coupling or changing screen-local state contracts.

### 4. Validation

- Run the owned acceptance and frontend gateway checks first, then the task-level harness and review commands required by this plan.

## Acceptance Checks

- One end-to-end import-to-job scenario runs through the agreed fixture path.
- Frontend create and list usecases call Tauri gateway adapters instead of preview executors.
- Tauri app registration exposes the backend commands required by the integrated Phase 1 path.
- The existing acceptance anchor remains executable without contract redesign.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests proving the integrated import-to-job path and transport wiring.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added thin backend transport commands `create_job` and `list_jobs` in `src-tauri/src/gateway/commands.rs`, registered them in `src-tauri/src/lib.rs`, and backed both with a shared in-memory repository under `src-tauri/src/infra/` so one created job becomes observable through the next list call in the same desktop process.
- Added transport-only serde support for the existing backend create request DTO so the frontend's camelCase request shape remains invokable without reopening the Phase 1 DTO contract.
- Added concrete Tauri adapters under `src/gateway/tauri/job-create/` and `src/gateway/tauri/job-list/`, then replaced the preview create/list executors in `src/main.ts` without adding create-to-list coupling or changing screen-local contracts.
- Updated `src-tauri/tests/acceptance/import-job/` to prove `fixture import -> create command -> list command` through the agreed fixture path, and added `src-tauri/tests/import_job_transport_contract.rs` to lock the request transport shape.
- Single-pass review returned `reroute` for missing cleanup of acceptance cache artifacts and in-memory job-store entries; the reroute fix was applied in `src-tauri/tests/acceptance/import-job/mod.rs` and `src-tauri/src/infra/job_repository.rs`, and per workflow no second review pass was run afterward.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this task closed a planned Phase 1 integration gap without changing repository-level debt posture.
- Final validation passed: `cargo test --manifest-path ./src-tauri/Cargo.toml import_job`, `npm run test -- src/gateway/tauri`, `sonar-scanner`, `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths ...` with `openIssueCount: 0`, and `powershell -File scripts/harness/run.ps1 -Suite all`.
