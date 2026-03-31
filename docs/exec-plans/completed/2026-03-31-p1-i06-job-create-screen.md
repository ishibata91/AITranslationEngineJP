- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement the first frontend job create screen and input path without absorbing job list concerns.
- task_id: P1-I06
- task_catalog_ref: tasks/phase-1/tasks/P1-I06.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-I06` to add the first UI path that triggers backend job creation.
- Keep the work inside the frontend job-create screen, view, and screen usecase boundary.
- Reuse the completed backend job creation usecase and minimal job-state contract without pulling in job list rendering or provider controls.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/coding-guidelines.md`
- `tasks/phase-1/tasks/P1-I06.yaml`
- `docs/exec-plans/completed/2026-03-30-p1-c02-minimal-job-state-model-contract.md`
- `docs/exec-plans/completed/2026-03-31-p1-i04-backend-job-creation-usecase.md`

## Owned Scope

- `src/ui/screens/job-create/`
- `src/ui/views/job-create/`
- `src/application/usecases/job-create/`

## Out Of Scope

- job list rendering
- provider controls
- backend domain or Tauri command implementation unless the frontend path proves a missing transport boundary that already belongs to an existing gateway slice

## Dependencies / Blockers

- `P1-C02` defines the observable minimal job-state boundary.
- `P1-I04` defines the backend create input and success result boundary consumed by this screen.
- The frontend path must stay independent from parallel `P1-I05` and `P1-I07` work so UI create and list scopes do not merge.

## Parallel Safety Notes

- Shared files must stay limited to additive exports and app-shell composition points.
- The create flow must not depend on job list state, list query data, or provider execution controls.

## UI

- Add one stable job-create screen under `src/ui/screens/job-create/` and a paired presentational view under `src/ui/views/job-create/`.
- The first screen should stay in the app-shell single-window model and render as one focused create panel instead of adding route infrastructure or multi-screen navigation.
- The view should expose a minimal editable request form that covers one or more source groups with the backend create contract fields `source_json_path`, `target_plugin`, and canonical translation-unit fields already fixed by `P1-C01` and consumed by `P1-I04`.
- The screen should keep presentation narrow: form inputs, one submit action, a pending affordance, a success summary showing `job_id` and observable `state`, and an actionable error block. Job list layout, provider controls, and execution progress remain absent.

## Scenario

- A user opens the app and can reach the job-create screen through the app shell without browser-style routing.
- The user can review or edit the minimal grouped create request payload and submit exactly one create action from the screen.
- While create is in flight, the submit action is disabled and the view shows a stable pending state so duplicate requests are not launched from repeated clicks.
- When create succeeds, the view shows the returned `job_id` and observable `Ready` state from the backend contract without introducing list state or follow-up mutation controls.
- When create fails, the view stays on the same screen, preserves the editable request payload, and shows one user-facing error message instead of raw transport details.

## Logic

- `src/application/usecases/job-create/` owns the screen-local request model, create action, and state transitions for idle, pending, success, and failure.
- The job-create usecase should accept an injected async executor for the create request and must not call Tauri APIs from Svelte files. Concrete `src/gateway/tauri/` wiring stays outside this task and remains for `P1-G01`.
- The UI-facing state should carry only the minimal fields needed for this task: editable create request data, `isSubmitting`, `job_id | null`, observable `state | null`, and one user-facing error message.
- Frontend-side validation should stay light and mechanical: prevent blank required fields from being submitted, but keep canonical contract enforcement in the backend path described by `P1-I04` and `P1-C02`.
- Request and result types should mirror the backend create DTO shape closely enough that `P1-G01` can wire the same usecase to the transport boundary without redesigning the screen contract.

## Implementation Plan

- Ordered scope 1: add `src/application/usecases/job-create/` with the local request and result types, screen state, light client-side validation, and an injected async create executor so the UI slice can exercise the backend-aligned create contract without adding Tauri wiring in this task.
- Ordered scope 2: add `src/ui/views/job-create/` and `src/ui/screens/job-create/` so the first create panel renders the editable request payload, submit action, pending state, success summary, and failure state through one screen/view split.
- Ordered scope 3: make the minimal app-shell composition update required to render the new job-create screen in the single-window shell, while keeping list rendering and provider controls out of scope.
- Ordered scope 4: add frontend tests that prove local validation, success and failure transitions, and screen rendering for the first create path.

## Acceptance Checks

- The app shell can render the first job-create screen.
- The job-create screen can trigger the backend job creation path through a frontend usecase boundary.
- The job-create screen shows stable pending, success, and failure states without absorbing list concerns.
- `npm test -- src/application/usecases/job-create/index.test.ts src/ui/screens/job-create/index.test.ts`
- `npm run build`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated frontend tests that cover create-screen state transitions and screen rendering.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src/application/usecases/job-create/` with the local create request and result contract, a screen-local store, light request validation, and an injected async executor boundary.
- Added `src/ui/screens/job-create/` and `src/ui/views/job-create/` with a single-panel job create form, pending state, success summary, and failure display.
- Updated `src/App.svelte`, `src/main.ts`, and `src/ui/app-shell/AppShell.svelte` to surface the job-create screen in the single-window shell while preserving the existing bootstrap status panel.
- Added frontend coverage in `src/application/usecases/job-create/index.test.ts` and `src/ui/screens/job-create/index.test.ts`, including success, validation, duplicate-submit, and failure-path usecase checks.
- Final validation passed: `npm test -- src/application/usecases/job-create/index.test.ts src/ui/screens/job-create/index.test.ts`, `npm run build`, `powershell -File scripts/harness/run.ps1 -Suite all`, and owned-path Sonar open issues = `0`.
- Single-pass review result: `pass`.
- Remaining integration gap: the current app composition uses a preview executor in `src/main.ts`; concrete frontend-to-backend transport wiring remains for `P1-G01`, which owns `src/gateway/tauri/`.
