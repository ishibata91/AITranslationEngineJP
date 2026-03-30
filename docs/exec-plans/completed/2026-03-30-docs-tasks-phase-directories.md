- workflow: impl
- status: completed
- lane_owner: codex
- scope: Restructure docs/tasks into per-phase directories with per-task YAML files, then sync roadmap and task references.
- task_id: DOCS-TASKS-DIRECTORIES
- task_catalog_ref: docs/tasks/phase-1/phase.yaml
- parent_phase: cross-phase

## Request Summary

- Split `docs/tasks/` by phase directory instead of one YAML file per phase.
- Split each phase task list into one YAML file per task.
- Keep the task catalog machine-readable and keep roadmap and plan references aligned.

## Decision Basis

- `AGENTS.md`
- `docs/tasks/README.md`
- `4humans/development-roadmap.md`
- `docs/index.md`
- `docs/exec-plans/templates/impl-plan.md`

## Owned Scope

- `docs/tasks/`
- `4humans/development-roadmap.md`
- `docs/index.md`
- `docs/exec-plans/templates/impl-plan.md`

## Out Of Scope

- `.codex/` workflow contracts
- product code under `src/` and `src-tauri/`
- harness scripts
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Dependencies / Blockers

- Existing phase YAML content must be preserved while changing layout.
- References that currently point to `docs/tasks/phase-*.yaml` must move to the new phase directory layout.
- Validation must still prove YAML parseability and Phase 1 parallel-batch `owned_scope` disjointness.

## Parallel Safety Notes

- `phase.yaml` should keep only phase-level metadata and task file references.
- Each task file should contain exactly one task record so ownership and diffs stay isolated.
- The roadmap should continue referencing stable task IDs rather than duplicating task body detail.

## UI

- N/A for end-user screens.

## Scenario

- A human opens `docs/tasks/phase-1/phase.yaml` to inspect phase-level metadata and batch grouping.
- A human or agent opens `docs/tasks/phase-1/tasks/P1-C01.yaml` to inspect a single task in isolation.
- An impl plan can point either to a phase file or a single task file when narrowing scope.

## Logic

- `docs/tasks/<phase>/phase.yaml` stores phase metadata, `parallel_batches`, and task file references.
- `docs/tasks/<phase>/tasks/<task_id>.yaml` stores the full task record for one task.
- Old flat phase files under `docs/tasks/phase-*.yaml` must be removed once references are updated.

## Implementation Plan

- Add per-phase directories and split current phase YAML files into `phase.yaml` plus per-task YAML files.
- Update `docs/tasks/README.md` to describe the new directory layout and reference shape.
- Update `4humans/development-roadmap.md`, `docs/index.md`, and `docs/exec-plans/templates/impl-plan.md` to use the new phase paths.
- Re-run structure/design/full harness and YAML validation against the new layout.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure` passes.
- `powershell -File scripts/harness/run.ps1 -Suite design` passes.
- `powershell -File scripts/harness/run.ps1 -Suite all` passes.
- Every `docs/tasks/phase-*/phase.yaml` and task file parses as YAML.
- Phase 1 parallel batches still have disjoint `owned_scope` values across their referenced task files.

## Required Evidence

- Added `docs/tasks/phase-*/phase.yaml`
- Added `docs/tasks/phase-*/tasks/*.yaml`
- Removed flat `docs/tasks/phase-*.yaml`
- Updated roadmap, index, and impl plan template references
- Validation output for structure, design, full harness, YAML parse, and Phase 1 disjoint-scope check

## 4humans Sync

- `4humans/development-roadmap.md`

## Outcome

- Split `docs/tasks/` into `phase-1/` through `phase-5/`, with `phase.yaml` holding phase metadata and `tasks/<task_id>.yaml` holding one task per file.
- Updated `docs/tasks/README.md`, `4humans/development-roadmap.md`, `docs/exec-plans/templates/impl-plan.md`, and the earlier completed plan so they point at the new phase-directory layout.
- Removed the flat `docs/tasks/phase-*.yaml` files after moving their content into the per-phase directories.
- Validation passed for `powershell -File scripts/harness/run.ps1 -Suite structure`, `powershell -File scripts/harness/run.ps1 -Suite design`, `powershell -File scripts/harness/run.ps1 -Suite all`, Python YAML parsing for `docs/tasks/phase-*/phase.yaml` and `docs/tasks/phase-*/tasks/*.yaml`, and the Phase 1 `owned_scope` disjoint check.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this change reorganized task catalog structure only.
