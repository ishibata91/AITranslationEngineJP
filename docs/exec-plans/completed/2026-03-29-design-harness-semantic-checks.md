# Impl Plan

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `scripts/harness/check-design.ps1`, `4humans/tech-debt-tracker.md`, `4humans/quality-score.md`

## Request Summary

- Resolve tech debt item 3 by extending the design harness beyond keyword presence checks.

## Decision Basis

- The current design harness only proves that contract words exist in key documents.
- `docs/spec.md` and `docs/architecture.md` already define canonical terminology and boundary rules that can be checked more semantically.

## UI

- No UI changes.

## Scenario

- Contributors should get a failing design harness when the repository layout drifts from the architecture contract or when canonical design records lose required semantic guarantees.

## Logic

- Add semantic checks that compare architecture rules to the current repository layout.
- Add terminology-oriented checks that validate the canonical executable-spec wording across source-of-truth docs.

## Implementation Plan

- Extend `check-design.ps1` with semantic assertions in addition to pattern presence checks.
- Validate the initial architecture layout for frontend and backend directories and guard forbidden frontend directories.
- Validate canonical wording for executable-spec ownership in core docs, then update human-facing records to reflect the stronger harness.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite design` passes with the new semantic checks enabled.
- `powershell -File scripts/harness/run.ps1 -Suite all` passes after the harness update.
- `4humans/tech-debt-tracker.md` marks item 3 as closed.

## Required Evidence

- Design harness output showing the semantic checks pass.
- Full harness output after the change.

## Docs Sync

- `4humans/tech-debt-tracker.md`
- `4humans/quality-score.md`

## Outcome

- Extended `check-design.ps1` with semantic checks for architecture bootstrap boundaries and canonical executable-spec phrasing.
- Synced `docs/core-beliefs.md`, `4humans/quality-score.md`, and `4humans/tech-debt-tracker.md` to reflect the stronger design harness.
- Re-ran design and full harness validation successfully.
