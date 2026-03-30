- workflow: impl
- status: completed
- lane_owner: codex
- scope: Restructure the root task catalog into per-phase directories with per-task YAML files, then sync roadmap and task references.
- task_id: DOCS-TASKS-DIRECTORIES
- task_catalog_ref: tasks/phase-1/phase.yaml
- parent_phase: cross-phase

## Request Summary

- Split the task catalog by phase directory instead of one YAML file per phase.
- Split each phase task list into one YAML file per task.
- Keep the task catalog machine-readable and keep roadmap and plan references aligned.

## Decision Basis

- `AGENTS.md`
- `tasks/README.md`
- `4humans/development-roadmap.md`
- `docs/index.md`
- `docs/exec-plans/templates/impl-plan.md`

## Owned Scope

- `tasks/`
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

- A human opens `tasks/phase-1/phase.yaml` to inspect phase-level metadata and batch grouping.
- A human or agent opens `tasks/phase-1/tasks/P1-C01.yaml` to inspect a single task in isolation.
- An impl plan can point either to a phase file or a single task file when narrowing scope.

## Logic

- `tasks/<phase>/phase.yaml` stores phase metadata, `parallel_batches`, and task file references.
- `tasks/<phase>/tasks/<task_id>.yaml` stores the full task record for one task.
- Old flat phase files must be removed once references are updated.

## Implementation Plan

- Add per-phase directories and split current phase YAML files into `phase.yaml` plus per-task YAML files.
- Update `tasks/README.md` to describe the new directory layout and reference shape.
- Update `4humans/development-roadmap.md`, `docs/index.md`, and `docs/exec-plans/templates/impl-plan.md` to use the new phase paths.
- Re-run structure/design/full harness and YAML validation against the new layout.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure` passes.
- `powershell -File scripts/harness/run.ps1 -Suite design` passes.
- `powershell -File scripts/harness/run.ps1 -Suite all` passes.
- Every `tasks/phase-*/phase.yaml` and task file parses as YAML.
- Phase 1 parallel batches still have disjoint `owned_scope` values across their referenced task files.

## Required Evidence

- Added `tasks/phase-*/phase.yaml`
- Added `tasks/phase-*/tasks/*.yaml`
- Removed flat phase YAML files
- Updated roadmap, index, and impl plan template references
- Validation output for structure, design, full harness, YAML parse, and Phase 1 disjoint-scope check

## 4humans Sync

- `4humans/development-roadmap.md`

## Outcome

- Split the task catalog into `phase-1/` through `phase-5/`, with `phase.yaml` holding phase metadata and `tasks/<task_id>.yaml` holding one task per file.
- Updated `tasks/README.md`, `4humans/development-roadmap.md`, `docs/exec-plans/templates/impl-plan.md`, and the earlier completed plan so they point at the phase-directory layout.
- Removed the flat phase YAML files after moving their content into the per-phase directories.
- Validation passed for `powershell -File scripts/harness/run.ps1 -Suite structure`, `powershell -File scripts/harness/run.ps1 -Suite design`, `powershell -File scripts/harness/run.ps1 -Suite all`, Python YAML parsing for `tasks/phase-*/phase.yaml` and `tasks/phase-*/tasks/*.yaml`, and the Phase 1 `owned_scope` disjoint check.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this change reorganized task catalog structure only.
