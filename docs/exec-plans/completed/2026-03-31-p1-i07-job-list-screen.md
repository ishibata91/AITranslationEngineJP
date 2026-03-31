- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement the first frontend job list screen and observation path without absorbing job creation mutation or provider detail.
- task_id: P1-I07
- task_catalog_ref: tasks/phase-1/tasks/P1-I07.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-I07` to add the first UI path that observes the Phase 1 job list and status view.
- Keep the work inside the frontend job-list screen, view, and screen usecase boundary.
- Reuse the completed backend job list query and minimal job-state contract without pulling in job creation mutation flow or provider progress detail.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `docs/screen-design/wireframes/app-shell.md`
- `tasks/phase-1/tasks/P1-I07.yaml`
- `docs/exec-plans/completed/2026-03-30-p1-c02-minimal-job-state-model-contract.md`
- `docs/exec-plans/completed/2026-03-31-p1-i05-backend-job-list-query.md`
- `docs/exec-plans/active/2026-03-31-p1-i06-job-create-screen.md`

## Owned Scope

- `src/ui/screens/job-list/`
- `src/ui/views/job-list/`
- `src/application/usecases/job-list/`

## Out Of Scope

- job creation mutation handling
- provider progress detail
- backend domain or Tauri command implementation unless an existing frontend-owned boundary needs a local adapter stub

## Dependencies / Blockers

- `P1-C02` defines the observable minimal job-state boundary consumed by the list UI.
- `P1-I05` defines the backend query result that this screen should observe.
- `P1-I06` is parallel-safe but must remain separate so create input state and list observation state do not merge.

## Parallel Safety Notes

- Shared files must stay limited to additive exports and app-shell composition points.
- The list screen must not import job-create internals, mutate job state, or absorb execution-control UI.

## UI

- Add one stable job-list screen under `src/ui/screens/job-list/` and a paired presentational view under `src/ui/views/job-list/`, rendered as one observation panel that fits the existing single-window shell model without adding routing or job-create controls.
- The view should expose a compact header with a refresh action, a primary list region for observable jobs, and a secondary status summary for the currently selected job. Each row and the status summary must show only the minimal observable fields: `jobId` and shared `state`.
- The panel should render four stable presentation states from the same layout: initial loading when no data has been observed yet, successful empty state when the backend returns zero jobs, loaded list state with selectable rows, and retryable failure state. Refresh-time failures may keep the previous successful list visible behind the error block.
- Provider progress detail, execution controls, create inputs, and source provenance detail remain absent from the screen. The first status view is the selected job's minimal state summary only.

## Scenario

- When the screen mounts, it should trigger one initialize path that requests the current backend-owned job list snapshot through the frontend usecase boundary.
- When the first load succeeds with one or more jobs, the screen shows the list and auto-selects the first returned `jobId` so the minimal status summary is populated without a second user action.
- When the user refreshes after a successful load, the current selection should stay active if that `jobId` still exists in the refreshed result. If the selected job disappeared, the screen should fall back to the first returned job or to `null` when the list becomes empty.
- When the backend returns an empty list, the screen should show a successful empty state instead of an error and clear the selected-job summary. When observation fails, the screen should stay on the same panel, show one user-facing error message with retry/refresh affordance, and never mutate jobs or surface provider detail.

## Logic

- `src/application/usecases/job-list/` owns the screen-local observation contract. The query model for this task carries only a list of minimal job items, and each item carries only the frontend-visible job identifier plus the shared minimal job state already fixed by `P1-C02` and returned by `P1-I05`.
- The job-list usecase should reuse the repo-standard query-screen state shape for `data`, `loading`, `error`, and `selection`, and should expose observation-only actions: `initialize`, `refresh`, `retry`, and `select`. Phase 1 list observation does not need local filters or mutation actions.
- The usecase must accept an injected async list executor or gateway boundary and must not call Tauri APIs from Svelte files. Concrete IPC wiring stays outside this task-local design and should plug into the minimal list result from the completed backend query without redesigning the DTO boundary.
- Selection reconciliation is local UI logic: keep the current `jobId` when it is still present in refreshed data, otherwise select the first returned job, otherwise clear selection. Error mapping must stay user-facing and generic and must not leak transport or filesystem details.

## Implementation Plan

- Ordered scope 1: add `src/application/usecases/job-list/` with the minimal job-list query model, frontend-owned screen state built on the repo-standard feature-screen primitives, injected async list executor, selection reconciliation, and generic user-facing error mapping without adding transport wiring in this task.
- Ordered scope 2: add `src/ui/views/job-list/` so the presentational job-list panel renders one stable layout with refresh, loading, empty, loaded, and retryable failure states plus the selected job's minimal status summary.
- Ordered scope 3: add `src/ui/screens/job-list/` and the minimal single-window shell composition needed to mount the new screen through additive prop wiring only, while keeping job-create flow, provider detail, and gateway ownership outside this task.
- Ordered scope 4: add frontend tests that cover initialize, refresh, retry, and select behavior, selection reconciliation across refreshed results, empty-list success handling, generic failure mapping, and server-rendered screen/view output.

## Acceptance Checks

- The app shell can render the first job-list screen.
- The job-list screen can observe the backend-owned Phase 1 job list result through a frontend usecase boundary.
- The screen shows stable loading, empty, loaded, and failure states without absorbing create or provider detail.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated frontend tests that cover list-screen state transitions and rendering.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src/application/usecases/job-list/` with the Phase 1 frontend-owned observation contract, injected async executor, selection reconciliation, and generic user-facing error mapping while keeping transport wiring out of scope for `P1-G01`.
- Added `src/ui/views/job-list/` and `src/ui/screens/job-list/` so the first job-list panel renders stable loading, empty, loaded, and retryable failure states plus a selected-job summary using only `jobId + state`.
- Updated `src/App.svelte`, `src/main.ts`, and `src/ui/app-shell/AppShell.svelte` additively so the single-window shell surfaces the first job-list screen alongside the existing Phase 1 screens with preview observation data.
- Added `src/application/usecases/job-list/index.test.ts` and `src/ui/screens/job-list/index.test.ts`, including regression coverage for refresh-failure retention and a behavioral shell render harness that exercises the real `AppShell.svelte` source with the real `JobListScreen.svelte`.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because the task added local evidence without changing repository-level debt posture.
- Validation passed: `npm test -- src/application/usecases/job-list/index.test.ts src/ui/screens/job-list/index.test.ts`, `npm run build`, `powershell -File scripts/harness/run.ps1 -Suite all`, and owned-path Sonar open issues = `0`.
- Single-pass review result: `pass`.
