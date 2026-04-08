---
name: working-light
description: Handle clearly bounded small tasks. Use when scope, affected files, and validation are obvious upfront and the change can be completed directly without extra planning. Do not use when the task grows, needs spec interpretation, or requires lane-level coordination.
---

# Working Light

## Output

- touched_files
- completed_scope
- validation_results
- reroute_needed

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- Handle only narrow tasks that are already known to be light.
- Confirm scope, affected files, and required validation before editing.
- Keep `4humans sync` in the same change when the narrow task affects quality records or review diagrams.
- When the narrow task changes review diagrams, explicitly update or confirm `4humans/class-diagrams/` and `4humans/sequence-diagrams/`.
- If the narrow task changes codebase boundaries or execution flow, use `diagramming-d2` and update the required `4humans/class-diagrams/` or `4humans/sequence-diagrams/` `.d2` / `.svg` in the same change.
- When a task note or plan exists, list the required `4humans/...diagrams` updates in its `4humans Sync` section.
- Use `skill-modification` for `.codex/` updates and `updating-docs` for `docs/` source-of-truth updates.
- Stop and reroute when extra investigation or design starts to grow.
- Do not take on broad refactors, architecture changes, or new lane decisions.

## Reroute

- Specification interpretation is unclear.
- Touched files or impact area become wider than expected.
- Multiple skills or lane owners need to coordinate.
- Minimal validation is no longer enough and test design or a new plan is required.
