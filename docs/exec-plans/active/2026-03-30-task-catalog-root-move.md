- workflow: impl
- status: planned
- lane_owner: codex
- scope: Move the parallel-ready task catalog from docs/tasks to tasks at repository root, then sync all references.
- task_id: TASK-CATALOG-ROOT-MOVE
- task_catalog_ref: tasks/phase-1/phase.yaml
- parent_phase: cross-phase

## Request Summary

- Move the parallel-ready task catalog from `docs/tasks/` to `tasks/` at repository root.
- Keep the per-phase directory layout and per-task YAML files.
- Update roadmap, docs index, impl plan template, and existing completed plans so they reference the new root path.

## Decision Basis

- `AGENTS.md`
- `docs/index.md`
- `4humans/development-roadmap.md`
- current `docs/tasks/` catalog layout
- user request to place the catalog outside `docs/`

## Owned Scope

- `tasks/`
- `docs/index.md`
- `4humans/development-roadmap.md`
- `docs/exec-plans/templates/impl-plan.md`
- affected completed plan references

## Out Of Scope

- `.codex/` workflow contracts
- product code under `src/` and `src-tauri/`
- harness scripts
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Dependencies / Blockers

- The existing task catalog content must remain unchanged apart from the path move.
- Every reference to `docs/tasks/` must move to `tasks/` so no stale links remain.
- Validation must still prove YAML parseability and Phase 1 `owned_scope` disjointness after the move.

## Parallel Safety Notes

- The move must keep `tasks/phase-*/phase.yaml` and `tasks/phase-*/tasks/*.yaml` intact so future task-level diffs stay isolated.
- The roadmap should continue referencing stable task IDs rather than inlining task body detail.

## UI

- N/A for end-user screens.

## Scenario

- A human reads `4humans/development-roadmap.md` and opens `tasks/phase-1/phase.yaml` for the current batch.
- An implementation plan points directly to `tasks/phase-1/tasks/P1-C01.yaml` or another single task file.
- Other skills can update `tasks/` without modifying `docs/`.

## Logic

- `tasks/` becomes the machine-readable task catalog root.
- `docs/index.md` should point to `../tasks/README.md` because the catalog is no longer inside `docs/`.
- Existing completed plans that record `task_catalog_ref` or path evidence should be updated so the links stay valid.

## Implementation Plan

- Move `docs/tasks/` to `tasks/` without changing the internal phase/task layout.
- Update `4humans/development-roadmap.md`, `docs/index.md`, and `docs/exec-plans/templates/impl-plan.md` to use the new root path.
- Update completed plan files that reference the task catalog so the historical records remain navigable.
- Re-run structure, design, full harness, YAML parse checks, and the Phase 1 `owned_scope` disjoint check.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure` passes.
- `powershell -File scripts/harness/run.ps1 -Suite design` passes.
- `powershell -File scripts/harness/run.ps1 -Suite all` passes.
- Every `tasks/phase-*/phase.yaml` and `tasks/phase-*/tasks/*.yaml` parses as YAML.
- Phase 1 parallel batches still have disjoint `owned_scope` values across their referenced task files.

## Required Evidence

- `tasks/` exists at repository root with the phase directory layout preserved
- `docs/tasks/` is removed
- Updated roadmap, index, impl plan template, and completed plan references
- Validation output for structure, design, full harness, YAML parse, and Phase 1 disjoint-scope check

## 4humans Sync

- `4humans/development-roadmap.md`

## Outcome

- Pending
