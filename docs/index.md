# Repository Index

`docs/` はプロダクト仕様と設計の記録系であり、作業方法の正本は `.codex/` にある。
新規参加者とエージェントは `AGENTS.md` の後に `.codex/README.md` を読み、その後にこのページを使う。

## Read Order

1. [`../.codex/README.md`](../.codex/README.md)
2. [`../.codex/skills/architect-direction/SKILL.md`](../.codex/skills/architect-direction/SKILL.md)
3. [`core-beliefs.md`](./core-beliefs.md)
4. [`spec.md`](./spec.md)
5. [`architecture.md`](./architecture.md)
6. [`tech-selection.md`](./tech-selection.md)
7. [`er-draft.md`](./er-draft.md)
8. [`executable-specs.md`](./executable-specs.md)
11. Relevant file under [`exec-plans/`](./exec-plans/)

## Directory Contract

- [`core-beliefs.md`](./core-beliefs.md): agent-first principles, hard rules, and repository habits
- [`../.codex/`](../.codex/README.md): multi-agent workflow, role contracts, and workflow skills
- [`../.codex/skills/architect-direction/`](../.codex/skills/architect-direction/SKILL.md): standard architect entrypoint
- [`../.codex/skills/light-direction/`](../.codex/skills/light-direction/SKILL.md): lightweight fix / tweak entrypoint
- [`spec.md`](./spec.md): permanent requirements and glossary
- [`architecture.md`](./architecture.md): layers, ports, dependency direction, and boundaries
- [`tech-selection.md`](./tech-selection.md): chosen technologies and quality tooling
- [`er-draft.md`](./er-draft.md): concept-level data model draft
- [`executable-specs.md`](./executable-specs.md): testable constraints, acceptance checks, and executable-spec policy
- [`exec-plans/active/`](./exec-plans/active/README.md): plans that are not yet complete
- [`exec-plans/completed/`](./exec-plans/completed/README.md): finished plans and outcomes
- [`quality-score.md`](./quality-score.md): current quality posture and missing coverage
- [`tech-debt-tracker.md`](./tech-debt-tracker.md): unresolved debt and cleanup backlog
- [`references/`](./references/index.md): curated reference index and external material policy

## Choose The Right Record

- Requirement or product boundary changed: update [`spec.md`](./spec.md)
- Dependency rule or layering changed: update [`architecture.md`](./architecture.md)
- Technology decision changed: update [`tech-selection.md`](./tech-selection.md)
- Data structure or entity relationship changed: update [`er-draft.md`](./er-draft.md)
- Detailed behavior or constraint changed: update [`executable-specs.md`](./executable-specs.md) and the corresponding tests or acceptance checks
- Work is non-trivial and not yet finished: create a plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work is finished: move the plan into [`exec-plans/completed/`](./exec-plans/completed/README.md)
- Workflow or role confusion keeps recurring: update [`../.codex/`](../.codex/README.md) or the relevant file in `../.codex/agents/` or `../.codex/skills/`
- Product-level confusion keeps recurring: add a rule to [`core-beliefs.md`](./core-beliefs.md) or [`AGENTS.md`](../AGENTS.md)
- The repository is missing coverage or confidence: update [`quality-score.md`](./quality-score.md)
- The problem is known but not resolved yet: update [`tech-debt-tracker.md`](./tech-debt-tracker.md)

## Repository Checks

- Structure harness: `powershell -File scripts/harness/run.ps1 -Suite structure`
- Design harness: `powershell -File scripts/harness/run.ps1 -Suite design`
- Execution harness: `powershell -File scripts/harness/run.ps1 -Suite execution`
- Full pass: `powershell -File scripts/harness/run.ps1 -Suite all`

## Notes

- New external references should be added under [`references/`](./references/index.md).
- The legacy raw API dumps still live under `docs/api-refrences/`. Keep them until they are migrated,
  but do not add new files there.
- Workflow skills live under [`../.codex/skills/`](../.codex/skills/).
