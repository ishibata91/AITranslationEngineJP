# Repository Index

`docs/` はプロダクト仕様と設計の記録系であり、作業方法の正本は `.codex/` にある。
新規参加者とエージェントは `AGENTS.md` の後に `.codex/README.md` を読み、その後にこのページを使う。
workflow の入口は `directing-implementation` と `directing-fixes` で、必要な task-local design は `designing-implementation` を通す。
詳細な振る舞いと制約は tests / acceptance checks / validation commands を正本として扱う。
`docs/` 正本は human が先に更新し、agent は human が直接起動した `../.codex/skills/updating-docs/SKILL.md` でだけ同期する。

## Read Order

1. [`../.codex/README.md`](../.codex/README.md)
2. [`core-beliefs.md`](./core-beliefs.md)
3. [`spec.md`](./spec.md)
4. [`architecture.md`](./architecture.md)
5. [`tech-selection.md`](./tech-selection.md)
6. [`er.md`](./er.md)
7. [`coding-guidelines.md`](./coding-guidelines.md)
8. Relevant file under [`screen-design/`](./screen-design/)
9. [`lint-policy.md`](./lint-policy.md)
10. Relevant file under [`exec-plans/`](./exec-plans/)
11. Relevant file under [`../tasks/`](../tasks/README.md) when parallel-ready task decomposition matters

## Directory Contract

- [`core-beliefs.md`](./core-beliefs.md): agent-first principles, hard rules, and repository habits
- [`../.codex/`](../.codex/README.md): multi-agent workflow, role contracts, and workflow skills
- [`spec.md`](./spec.md): permanent requirements and glossary
- [`architecture.md`](./architecture.md): layers, ports, dependency direction, and boundaries
- [`tech-selection.md`](./tech-selection.md): chosen technologies and quality tooling
- [`coding-guidelines.md`](./coding-guidelines.md): repository coding conventions for Tauri 2, Svelte 5, TypeScript, and Rust
- [`screen-design/`](./screen-design/): screen map, wireframes, and UI layout references
- [`er.md`](./er.md): canonical data model and ER specification
- [`lint-policy.md`](./lint-policy.md): what lint manages, what it does not manage, and tool ownership
- [`exec-plans/active/`](./exec-plans/active/README.md): plans that are not yet complete
- [`exec-plans/completed/`](./exec-plans/completed/README.md): finished plans and outcomes
- [`../tasks/`](../tasks/README.md): machine-readable task catalog for parallel-safe decomposition and batch planning
- [`../4humans/quality-score.md`](../4humans/quality-score.md): human-facing quality posture and missing coverage
- [`../4humans/tech-debt-tracker.md`](../4humans/tech-debt-tracker.md): human-facing unresolved debt and cleanup backlog
- [`references/`](./references/index.md): curated reference index and external material policy

## Choose The Right Record

- Requirement or product boundary changed: human-first update [`spec.md`](./spec.md) via [`../.codex/skills/updating-docs/SKILL.md`](../.codex/skills/updating-docs/SKILL.md)
- Dependency rule or layering changed: human-first update [`architecture.md`](./architecture.md) via [`../.codex/skills/updating-docs/SKILL.md`](../.codex/skills/updating-docs/SKILL.md)
- Technology decision changed: human-first update [`tech-selection.md`](./tech-selection.md) via [`../.codex/skills/updating-docs/SKILL.md`](../.codex/skills/updating-docs/SKILL.md)
- Coding conventions or Tauri 2 implementation rules changed: update [`coding-guidelines.md`](./coding-guidelines.md)
- Screen map, wireframe, or UI layout reference changed: update the relevant file under [`screen-design/`](./screen-design/)
- Data structure or entity relationship changed: human-first update [`er.md`](./er.md) via [`../.codex/skills/updating-docs/SKILL.md`](../.codex/skills/updating-docs/SKILL.md)
- Detailed behavior or constraint changed: update the corresponding tests or acceptance checks and validation commands
- Lint の責務範囲、allowlist 方針、tool ownership を変えた: update [`lint-policy.md`](./lint-policy.md)
- Work is non-trivial and not yet finished: create a plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work needs parallel-safe task decomposition or batch planning: update the relevant file in [`../tasks/`](../tasks/README.md)
- Work is finished: move the plan into [`exec-plans/completed/`](./exec-plans/completed/README.md)
- Workflow or role confusion keeps recurring: update [`../.codex/`](../.codex/README.md) or the relevant file in `../.codex/agents/` or `../.codex/skills/` via `skill-modification`
- Product-level confusion keeps recurring: add a rule to [`core-beliefs.md`](./core-beliefs.md) or [`AGENTS.md`](../AGENTS.md)
- The repository is missing coverage or confidence: update [`../4humans/quality-score.md`](../4humans/quality-score.md)
- The problem is known but not resolved yet: update [`../4humans/tech-debt-tracker.md`](../4humans/tech-debt-tracker.md)

## Repository Checks

- Structure harness: `powershell -File scripts/harness/run.ps1 -Suite structure`
- Design harness: `powershell -File scripts/harness/run.ps1 -Suite design`
- Execution harness: `powershell -File scripts/harness/run.ps1 -Suite execution`
- Full pass: `powershell -File scripts/harness/run.ps1 -Suite all`

## Notes

- New external references should be added under [`references/`](./references/index.md).
- Raw vendor API dumps live under [`references/vendor-api/`](./references/vendor-api/).
- Parallel-ready task catalogs live under [`../tasks/`](../tasks/README.md).
- Workflow skills live under [`../.codex/skills/`](../.codex/skills/).

