- workflow: impl
- status: completed
- lane_owner: codex
- scope: Introduce a YAML task catalog under docs/tasks, compress the human roadmap summary, and sync task references in docs/index.md and the impl plan template.
- task_id: DOCS-TASKS-CATALOG
- task_catalog_ref: docs/tasks/phase-1/phase.yaml
- parent_phase: cross-phase

## Request Summary

- Add a machine-readable task catalog under `docs/tasks/` so roadmap details can be split into `contract -> verification -> impl -> integ` tasks.
- Keep `4humans/development-roadmap.md` as the human-facing summary and next-batch guide.
- Add enough documentation so task authors can understand `owned_scope`, dependency edges, and how to keep `integ` small.
- Sync the repository index and impl plan template with the new task catalog.

## Decision Basis

- `AGENTS.md`
- `docs/index.md`
- `4humans/development-roadmap.md`
- `docs/exec-plans/templates/impl-plan.md`
- user-approved plan for `docs/tasks/` YAML task catalog

## Owned Scope

- `docs/tasks/`
- `4humans/development-roadmap.md`
- `docs/index.md`
- `docs/exec-plans/templates/impl-plan.md`

## Out Of Scope

- `.codex/` workflow contracts
- product code under `src/` and `src-tauri/`
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- harness script changes or YAML lint automation

## Dependencies / Blockers

- `4humans/development-roadmap.md` must remain a summary and not duplicate the task catalog.
- `docs/tasks/phase-1/phase.yaml` needs stable task IDs so the roadmap can reference the immediate next batch.
- `docs/index.md` must point readers to the new task catalog without changing permanent product requirements.

## Parallel Safety Notes

- `contract` tasks must stay at layer boundaries and must not absorb implementation detail.
- `verification` tasks must fix fixture shape and public behavior before multi-scope implementation begins.
- `impl` tasks must keep `owned_scope` disjoint enough that the same parallel batch can land independently.
- `integ` tasks must only do composition, wiring, and scenario proof; if a task needs new API or fixture design it must move back to `contract` or `verification`.

## UI

- N/A for end-user screens.
- Documentation should still make the `src/` and `src-tauri/` ownership split easy to read for future task authors.

## Scenario

- A human reads `4humans/development-roadmap.md` to understand the current phase, risks, and immediate batches.
- A human or agent opens `docs/tasks/phase-*/phase.yaml` and `docs/tasks/phase-*/tasks/*.yaml` to see the detailed task graph, `owned_scope`, and dependency edges.
- An implementation plan can reference a stable `task_id` and reuse `owned_scope` from the YAML task catalog instead of redefining parallel-safety rules from zero.

## Logic

- `docs/tasks/` is a machine-readable task catalog, not a replacement for `spec.md` or other permanent product docs.
- Each phase YAML must expose `phase_id`, `phase_name`, `goal`, `exit_criteria`, `task_types`, `parallel_batches`, and `tasks`.
- Each task must expose `id`, `title`, `type`, `status`, `goal`, `owned_scope`, `out_of_scope`, `depends_on`, `produces`, `consumes`, `acceptance_anchor`, `parallel_safe_with`, `shared_risks`, and `done_when`.
- `status` values stay ASCII and machine-friendly so humans can skim and scripts can adopt the files later without redesign.

## Implementation Plan

- Add `docs/tasks/README.md` to define the YAML model, task type semantics, `owned_scope` rules, and `integ` limits.
- Add `docs/tasks/phase-1/phase.yaml` through `docs/tasks/phase-5/phase.yaml` with phase goals, exit criteria, detailed task-file references, and immediate parallel batches.
- Compress `4humans/development-roadmap.md` into a summary that links to `docs/tasks/` and names the next batches by task ID.
- Update `docs/index.md` to mention the task catalog as the place for parallel-ready task details.
- Update `docs/exec-plans/templates/impl-plan.md` so active plans can reference `task_id`, `task_catalog_ref`, and `owned_scope`.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure` passes after the new docs and links are added.
- `powershell -File scripts/harness/run.ps1 -Suite design` passes after the task catalog references are added.
- `docs/tasks/phase-1/phase.yaml` plus its task files show no overlapping `owned_scope` values inside the same `parallel_batches` entry.
- `4humans/development-roadmap.md` references Phase 1 task IDs that exist in `docs/tasks/phase-1/tasks/*.yaml`.

## Required Evidence

- Added `docs/tasks/README.md`
- Added `docs/tasks/phase-1/phase.yaml` through `docs/tasks/phase-5/phase.yaml`
- Updated `4humans/development-roadmap.md`, `docs/index.md`, and `docs/exec-plans/templates/impl-plan.md`
- Validation output for structure and design harness
- Validation output for `powershell -File scripts/harness/run.ps1 -Suite all`
- YAML parse output for `docs/tasks/phase-*.yaml`
- `owned_scope` disjoint check output for `docs/tasks/phase-1/phase.yaml` and its task files

## 4humans Sync

- `4humans/development-roadmap.md`

## Outcome

- Added a machine-readable task catalog under `docs/tasks/` with phase-level YAML files and `contract -> verification -> impl -> integ` task typing.
- Reduced `4humans/development-roadmap.md` to a human-facing summary and immediate-batch guide.
- Synced the repository index and impl plan template so future tasks can reference `task_id`, `task_catalog_ref`, and `owned_scope`.
- Validation passed for `powershell -File scripts/harness/run.ps1 -Suite structure`, `powershell -File scripts/harness/run.ps1 -Suite design`, `powershell -File scripts/harness/run.ps1 -Suite all`, Python YAML parsing for `docs/tasks/phase-*/phase.yaml` and `docs/tasks/phase-*/tasks/*.yaml`, and the Phase 1 `owned_scope` disjoint check.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this change introduced a roadmap and task-catalog structure but no new unresolved product debt or quality posture change.
