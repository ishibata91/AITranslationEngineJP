# Impl Plan Template

- workflow: impl
- status: planned
- lane_owner:
- scope:
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

-

## Decision Basis

-

## Owned Scope

- Prefer repo-root path prefixes or stable scope tokens that match `tasks/phase-*/tasks/*.yaml` when available.

## Out Of Scope

-

## Dependencies / Blockers

-

## Parallel Safety Notes

- Note the shared files, shared fixtures, or upstream `contract` / `verification` tasks that must land first.

## UI

- Use only when the task changes screen structure, presentation, or interaction flow.

## Scenario

- Use only when the task changes user-visible behavior, state transitions, or execution flow.

## Logic

- Use only when the task changes domain logic, contracts, validation, or dependency boundaries.

## Implementation Plan

-

## Acceptance Checks

-

## Required Evidence

-

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  クラス追加、依存追加、責務分割変更などで構造が変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  処理追加、相互作用順序変更、主要シナリオ変更などで実行フローが変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。

## Outcome

-
