# Flow Light, Gate Heavy, Plan Stabilization Loop

- workflow: heavy
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: workflow contracts, plan templates, harness, and quality records

## Request Summary

- Introduce `flow light, gate heavy`, add `workflow-gate`, and standardize a heavy-only plan stabilization loop.

## Decision Basis

- This is heavy because it changes repository workflow contracts, plan templates, and harness behavior across `.codex/`, `docs/`, and `scripts/`.

## Investigation Summary

- Facts:
  - Current workflow uses `Architect -> Research -> Plan -> Coder -> Architect review` for heavy and `Architect -> Short plan -> Coder -> Architect review` for light.
  - `light-review` is currently the only explicit review checklist skill.
  - Structure and design harnesses do not validate gate-specific workflow contracts.
- Options:
  - Add more review stages.
  - Keep flow light and move enforcement into plan stabilization, evidence, and gate checks.
- Risks:
  - Workflow docs can drift if `.codex`, `AGENTS.md`, templates, and harness checks are not updated together.
  - Adding a new gate skill without updating harness leaves the contract unenforced.
- Unknowns:
  - None blocking after plan approval.

## Unknown Classification

- Blocking:
  - None.
- Non-blocking:
  - Final wording can be tightened during implementation as long as the approved behavior stays intact.

## Assumptions / Defaults

- `workflow-gate` is a read-only Architect skill, not a new role.
- `light-review` remains available as a supplemental checklist, not the default gate.

## Plan Ready Criteria

- Workflow source-of-truth docs describe stabilization loop, gate, and reroute behavior consistently.
- Plan templates capture evidence, reroute triggers, and docs sync requirements.
- Harness checks fail if the new contracts are missing.

## Implementation Plan

- Add an active plan, then update `.codex` workflow docs and role contracts to define `flow light, gate heavy`.
- Add `workflow-gate`, reframe `light-review` as supplemental, and update plan templates around evidence and unknown classification.
- Extend structure and design harnesses to enforce the new contracts.
- Update human-facing quality posture and executable-spec rules to reflect gate-driven quality.

## Delegation Map

- Research: none
- Coder: apply the approved workflow/document/harness updates and run validation
- Worker: none

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- Structure harness output showing `workflow-gate` is required.
- Design harness output showing stabilization loop, gate, and template fields are enforced.
- Full harness output showing no regressions after docs and script changes.

## Reroute Trigger

- If implementation uncovers an unresolved blocking conflict between `.codex` workflow contracts and existing repository policy, pause and revise the heavy plan before proceeding.

## Docs Sync

- `.codex/README.md`
- `.codex/agents/`
- `.codex/skills/`
- `AGENTS.md`
- `docs/core-beliefs.md`
- `docs/executable-specs.md`
- `docs/index.md`
- `4humans/quality-score.md`

## Record Updates

- Update workflow source-of-truth docs, plan templates, and harness scripts.
- Move this plan to `completed/` with outcome notes after validation passes.

## Outcome

- Added `workflow-gate` as the standard read-only gate for heavy and light flows.
- Updated `.codex`, `AGENTS.md`, plan templates, and human-facing records to use `flow light, gate heavy`.
- Standardized heavy-only `Plan Stabilization Loop` with blocking vs non-blocking unknown handling.

## Validation Results

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
- Execution harness still reports `SKIP no Cargo.toml or package.json targets found`, which matches the current repository state.
