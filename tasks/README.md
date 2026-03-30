# Task Catalog

関連文書: [`../index.md`](../index.md), [`../../4humans/development-roadmap.md`](../../4humans/development-roadmap.md), [`../exec-plans/templates/impl-plan.md`](../exec-plans/templates/impl-plan.md)

このディレクトリは、並列実行しやすい task 分解を YAML で保持する machine-readable task catalog である。
人間向けの要約と現在地は [`../../4humans/development-roadmap.md`](../../4humans/development-roadmap.md) を正本とし、phase metadata は `phase-*/phase.yaml`、task の詳細は `phase-*/tasks/*.yaml` を正本とする。

## File Layout

- `phase-1/` から `phase-5/`: phase ごとの task catalog directory
- `phase-*/phase.yaml`: phase metadata、`parallel_batches`、task file references
- `phase-*/tasks/<task_id>.yaml`: 1 task 1 file の task record

## Phase File Fields

- `phase_id`: stable phase identifier
- `phase_name`: phase label
- `goal`: phase goal
- `exit_criteria`: phase completion conditions
- `task_types`: definitions for `contract`, `verification`, `impl`, `integ`
- `parallel_batches`: recommended parallel launch groups with disjoint `owned_scope`
- `task_files`: relative paths to the phase's task files

## Task Fields

- `id`: stable task identifier referenced by roadmap and impl plans
- `title`: short task name
- `type`: one of `contract`, `verification`, `impl`, `integ`
- `status`: one of `done`, `ready`, `planned`, `blocked`
- `goal`: what the task must establish
- `owned_scope`: repo-root path prefixes or stable subsystem tokens owned by the task
- `out_of_scope`: explicit non-goals that must not be absorbed into the task
- `depends_on`: task IDs that must land first
- `produces`: artifacts or decisions produced by the task
- `consumes`: upstream artifacts or decisions the task assumes
- `acceptance_anchor`: checks, fixtures, or behaviors that prove the task is complete
- `parallel_safe_with`: task IDs that can run in parallel without scope collisions
- `shared_risks`: risks that can collapse parallel safety if the task grows
- `done_when`: concrete completion criteria

## Task Type Rules

- `contract`
  - Holds only layer-boundary contracts such as DTO shape, port interfaces, state enums, result shapes, fixture schema, and external validation policy.
  - Must not absorb helper-only structure, SQL detail, logging detail, retry detail, or private module decomposition.
- `verification`
  - Fixes public behavior and fixture or acceptance anchor before broad implementation starts.
  - Should prove what must remain true, not how an implementation is internally structured.
- `impl`
  - Stays inside a closed `owned_scope`.
  - Should be the default place for repository code, adapter code, UI code, and use case code.
- `integ`
  - Exists for composition root work, adapter wiring, cross-layer scenario proof, and end-to-end confirmation.
  - Must not become the place where new API shape, fixture shape, or state policy is invented.

## Owned Scope Rules

- Prefer repo-root path prefixes such as `src-tauri/src/application/job/create/` or `src/ui/screens/job-list/`.
- If a path does not exist yet, use the narrowest stable scope token that can later map to one implementation area.
- Keep `owned_scope` disjoint within the same `parallel_batches` entry.
- If two tasks need the same path prefix, split the shared decision into a smaller `contract` or `verification` task first.

## `depends_on` And `parallel_safe_with`

- `depends_on` means a task cannot safely start until the upstream task has fixed a required artifact or decision.
- `parallel_safe_with` means the tasks can progress at the same time because the required decisions are already fixed and the `owned_scope` values do not collide.
- A task can reference the same ID in both lists only when a future iteration is expected; do not do that in the initial catalog.

## Keeping `integ` Small

- `integ` must stay limited to DI registration, boundary wiring, scenario proof, and cross-layer acceptance.
- If an `integ` task needs to redefine DTO shape, fixture shape, error policy, or state transition, move that work back into `contract` or `verification`.
- Backend work still needs `integ` even with DIP, because composition, persistence wiring, runtime configuration, and end-to-end failure propagation remain cross-scope concerns.

## Example

```yaml
phase_id: phase-1
phase_name: "Input Cache And Job Skeleton"
goal: "Move from raw imported input toward a job-oriented workflow."
parallel_batches:
  - id: P1-B1
    goal: "Fix contracts and verification anchors before parallel implementation."
    task_ids:
      - "P1-C01"
      - "P1-C02"
      - "P1-V01"
      - "P1-V02"
task_files:
  - "tasks/P1-C01.yaml"
  - "tasks/P1-C02.yaml"
```

```yaml
id: P1-C01
title: "Define TRANSLATION_UNIT canonical contract"
type: contract
status: ready
owned_scope:
  - "src-tauri/src/application/dto/translation_unit/"
  - "src-tauri/src/domain/translation_unit/"
```
