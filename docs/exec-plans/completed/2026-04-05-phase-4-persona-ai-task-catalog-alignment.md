- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Align Phase 2 to Phase 4 task catalog so provider-backed persona generation work is explicitly represented without collapsing existing phase boundaries.
- task_id: phase-4-persona-ai-alignment
- task_catalog_ref: tasks/phase-4/phase.yaml
- parent_phase: phase-4

## Request Summary

- Add task-catalog coverage for provider-backed persona generation AI flow.
- Keep the existing separation where Phase 3 owns translation-phase orchestration and Phase 4 owns concrete AI runtime and provider execution.

## Decision Basis

- `docs/spec.md` requires AI-backed persona generation for both master persona construction and job-local translation flow.
- `tasks/phase-3/phase.yaml` already owns the translation flow MVP and includes `P3-I03` for NPC persona generation phase orchestration.
- `tasks/phase-4/phase.yaml` owns provider adapters and execution expansion, so concrete provider-backed persona generation belongs there.
- `docs/` permanent records are human-first; this change is limited to task catalog and execution plan alignment.

## Owned Scope

- `tasks/phase-2/tasks/P2-I03.yaml`
- `tasks/phase-3/tasks/P3-I03.yaml`
- `tasks/phase-3/tasks/P3-G01.yaml`
- `tasks/phase-4/phase.yaml`
- `tasks/phase-4/tasks/`

## Out Of Scope

- `docs/spec.md`, `docs/architecture.md`, `docs/tech-selection.md`
- implementation code or tests for persona generation runtime
- `4humans` summary or diagram updates unless the task catalog change forces them

## Dependencies / Blockers

- No upstream blocker for task-catalog edits.
- Existing Phase 2 and Phase 3 task IDs must remain stable.

## Parallel Safety Notes

- Keep new Phase 4 tasks on disjoint `owned_scope` values.
- Do not redefine Phase 2 storage split or Phase 3 handoff contracts beyond wording alignment.

## UI

- No UI change.

## Scenario

- Phase 3 keeps fake-or-abstracted orchestration for NPC persona generation inside the Translation Flow MVP.
- Phase 4 adds concrete provider-backed persona-generation work for both base-game master persona construction and job-local NPC persona generation.

## Logic

- Add a persona-generation runtime contract and acceptance anchor to Phase 4.
- Add one impl task for provider-backed master persona generation and one impl task for provider-backed job-local persona generation runtime.
- Add one integration task that proves provider-backed persona-generation scenarios without merging them into the existing provider/execution-control scenario task.
- Reword Phase 2 and Phase 3 tasks so their responsibilities stay foundation-boundary and orchestration-boundary focused.

## Implementation Plan

- Update `tasks/phase-4/phase.yaml` goal, exit criteria, batch goals, task IDs, and task file list.
- Add `P4-C03`, `P4-V02`, `P4-I07`, `P4-I08`, and `P4-G02` under `tasks/phase-4/tasks/`.
- Reword `P2-I03`, `P3-I03`, and `P3-G01` to remove ambiguity about concrete provider execution ownership.
- Run structure and full harness after the catalog change.

## Acceptance Checks

- Task YAML remains structurally valid and phase/task references stay coherent.
- Phase 4 explicitly covers provider-backed persona-generation work for master and job-local paths.
- Phase 2 and Phase 3 task wording no longer implies that concrete provider execution is already owned there.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- No `4humans` update planned. This change adjusts machine-readable task decomposition only.

## Outcome

- Added Phase 4 task-catalog coverage for provider-backed persona generation by
  introducing `P4-C03`, `P4-V02`, `P4-I07`, `P4-I08`, and `P4-G02`.
- Updated `tasks/phase-4/phase.yaml` so the phase goal, exit criteria, and batch
  goals explicitly include provider-backed master persona generation and job-local
  NPC persona generation.
- Reworded `P2-I03`, `P3-I03`, and `P3-G01` so Phase 2 remains foundation-boundary
  oriented and Phase 3 remains orchestration-boundary oriented.
- Validation passed with `python3 scripts/harness/run.py --suite structure` and
  `python3 scripts/harness/run.py --suite all`.
