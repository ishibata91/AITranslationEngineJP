# `.codex` Lane Reset

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/`, workflow docs, plan templates, harness, quality records

## Request Summary

- Replace the live `.codex` workflow with `impl-direction` and `fix-direction`, keep `UI` / `Scenario` / `Logic` inside exec-plans, and remove current-repo-incompatible legacy assumptions.

## Decision Basis

- This is non-trivial because it changes the live workflow contract across `.codex/`, `AGENTS.md`, `docs/`, and harness scripts.
- Current structure and design harness had already drifted from the intended live workflow.

## UI

- Not required for the workflow rewrite itself.

## Scenario

- `impl-direction` owns task-local design sections, implementation handoff, single-pass review, docs sync, and closeout.
- `fix-direction` owns bugfix framing, optional tracing/logging, single-pass review, docs sync, and closeout.

## Logic

- `UI` / `Scenario` / `Logic` now live inside exec-plans, not in `changes/` artifacts.
- Review is single-pass and only checks spec deviation, exception handling, resource cleanup, and missing tests.
- Legacy packet, `context_board`, `tasks.md`, and score-loop contracts are removed from the live workflow.

## Implementation Plan

- Rewrite workflow source-of-truth docs around `impl-direction` and `fix-direction`.
- Replace live `.codex/agents` with specialized contracts for distill, work planning, implementation, tracing, logging, and review.
- Replace live `.codex/skills` with repo-fitted `impl-*`, `fix-*`, and `risk-report` skills.
- Replace heavy/light templates with `impl-plan` and `fix-plan`, then update structure and design harness expectations.

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- Structure harness confirms only the new live `.codex` files and template names.
- Design harness confirms the new lane workflow, embedded `UI` / `Scenario` / `Logic`, and single-pass review contract.
- Full harness was attempted after the rewrite.

## Docs Sync

- `.codex/README.md`
- `AGENTS.md`
- `docs/index.md`
- `docs/core-beliefs.md`
- `docs/executable-specs.md`
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Replaced the live workflow contract with `impl-direction` and `fix-direction`.
- Kept task-local `UI` / `Scenario` / `Logic` inside exec-plans instead of reviving `changes/` artifacts.
- Replaced Architect/Research/Coder contracts with specialized helper agents.
- Rewrote structure and design harnesses to validate the new lane-based workflow.

## Validation Results

- `powershell -File scripts/harness/run.ps1 -Suite structure`: pass
- `powershell -File scripts/harness/run.ps1 -Suite design`: pass
- `powershell -File scripts/harness/run.ps1 -Suite all`: fail because `cargo` is not installed in the current environment, so the execution harness cannot run the Rust-side commands
