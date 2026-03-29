# Impl Plan

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `docs/api-refrences/`, `docs/references/`, `docs/index.md`, `4humans/tech-debt-tracker.md`

## Request Summary

- Resolve tech debt item 1 by removing the typoed `docs/api-refrences/` directory and moving its contents under `docs/references/`.

## Decision Basis

- `docs/references/` is already the documented source of truth for external reference material.
- Keeping a second differently named directory continues to create avoidable ambiguity for humans and agents.

## UI

- No UI changes.

## Scenario

- Contributors should find vendor API reference files only via `docs/references/` and its index.

## Logic

- Move vendor API reference assets into a dedicated subdirectory under `docs/references/`.
- Update index and tracker records to point to the canonical location.

## Implementation Plan

- Create `docs/references/vendor-api/` as the canonical location for raw vendor API reference assets.
- Move all files from `docs/api-refrences/` into the new location.
- Update docs and debt tracking to reflect the resolved state.

## Acceptance Checks

- `docs/api-refrences/` no longer exists.
- `docs/references/index.md` links or names the new vendor API location.
- `4humans/tech-debt-tracker.md` marks item 1 as closed.
- Harness checks pass after the documentation move.

## Required Evidence

- Directory listing showing the new location.
- Harness output for structure/design/full verification.

## Docs Sync

- `docs/index.md`
- `docs/references/index.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Moved the vendor API reference assets from the typoed legacy directory into `docs/references/vendor-api/`.
- Added `docs/references/vendor-api/README.md` so the new canonical location is navigable from the references index.
- Updated `docs/index.md` and `4humans/tech-debt-tracker.md`, then passed structure, design, and full harness verification.
