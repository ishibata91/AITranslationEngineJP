- workflow: impl
- status: planned
- lane_owner: codex
- scope: Fix the Phase 3 batch `P3-B1` contracts and regression anchors for translation instruction building, embedded-element preservation, and translation phase handoff.
- task_id: P3-B1
- task_catalog_ref: tasks/phase-3/phase.yaml
- parent_phase: phase-3

## Request Summary

- Implement batch `P3-B1` from `tasks/phase-3/phase.yaml`.
- Fix stable application boundaries for translation instruction building, embedded-element preservation, and word/persona/body translation phase handoff before broader Translation Flow MVP implementation starts.
- Add regression anchors for embedded-element preservation and one representative translation-flow scenario without absorbing provider-specific logic, retry policy, output writing, or preview UI behavior.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-3/phase.yaml`
- `tasks/phase-3/tasks/P3-C01.yaml`
- `tasks/phase-3/tasks/P3-C02.yaml`
- `tasks/phase-3/tasks/P3-C03.yaml`
- `tasks/phase-3/tasks/P3-V01.yaml`
- `tasks/phase-3/tasks/P3-V02.yaml`

## Owned Scope

- `src-tauri/src/application/dto/translation_instruction/`
- `src-tauri/src/application/dto/embedded_element_policy/`
- `src-tauri/src/application/dto/translation_phase_handoff/`
- `src-tauri/tests/regression/embedded-elements/`
- `src-tauri/tests/regression/translation-flow-mvp/`

## Out Of Scope

- provider adapters and provider-specific prompt assembly
- output writer policy and retry policy
- preview UI layout and frontend observation
- later composition and integration wiring reserved for `P3-B2` and `P3-B3`

## Dependencies / Blockers

- `P3-V01` depends on the stable output of `P3-C02`.
- `P3-V02` depends on the stable output of `P3-C01` and `P3-C03`.
- Existing Phase 2 dictionary lookup and persona storage contracts should remain reusable where Phase 3 handoff contracts depend on them.

## Parallel Safety Notes

- `P3-C01`, `P3-C02`, `P3-C03`, `P3-V01`, and `P3-V02` are grouped because they fix backend-owned contracts and regression anchors before provider-facing implementation and end-to-end composition begin.
- The batch must keep contract boundaries free from provider policy, output formatting, UI preview state, or phase-internal helper detail that would reduce independence of later tasks.

## UI

- N/A. `P3-B1` fixes backend application contracts and regression anchors only.

## Scenario

- To be refined after `designing-implementation`.

## Logic

- To be refined after `designing-implementation`.

## Implementation Plan

- To be refined after `planning-implementation`.

## Acceptance Checks

- `P3-C01`, `P3-C02`, and `P3-C03` each land on one stable contract inside the owned backend DTO scope.
- `P3-V01` proves embedded elements survive the protected translation-flow stages described by the preservation contract.
- `P3-V02` proves one representative Translation Flow MVP scenario can be anchored without depending on preview UI or provider retry policy.
- The batch leaves implementation-only or integration-only behavior out of contract and regression scope.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated regression tests or fixtures for embedded elements and representative translation flow.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- No diagram update expected unless implementation changes codebase boundaries or execution flow beyond the task catalog scope.

## Outcome

- Pending.
