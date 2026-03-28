# Repository Index

`docs/` はプロダクト仕様と設計の記録系であり、作業方法の正本は `.codex/` にある。
新規参加者とエージェントは `AGENTS.md` の後に `.codex/README.md` を読み、その後にこのページを使う。

## Read Order

1. [`../.codex/README.md`](../.codex/README.md)
2. 実装なら [`../.codex/skills/impl-direction/SKILL.md`](../.codex/skills/impl-direction/SKILL.md)
3. バグ修正なら [`../.codex/skills/fix-direction/SKILL.md`](../.codex/skills/fix-direction/SKILL.md)
4. [`core-beliefs.md`](./core-beliefs.md)
5. [`spec.md`](./spec.md)
6. [`architecture.md`](./architecture.md)
7. [`tech-selection.md`](./tech-selection.md)
8. [`er-draft.md`](./er-draft.md)
9. [`executable-specs.md`](./executable-specs.md)
10. [`lint-policy.md`](./lint-policy.md)
11. Relevant file under [`exec-plans/`](./exec-plans/)

## Directory Contract

- [`core-beliefs.md`](./core-beliefs.md): agent-first principles, hard rules, and repository habits
- [`../.codex/`](../.codex/README.md): multi-agent workflow, role contracts, and workflow skills
- [`../.codex/skills/impl-direction/`](../.codex/skills/impl-direction/SKILL.md): implementation lane entrypoint that can also settle task-local design
- [`../.codex/skills/fix-direction/`](../.codex/skills/fix-direction/SKILL.md): bugfix lane entrypoint with optional tracing and logging
- [`spec.md`](./spec.md): permanent requirements and glossary
- [`architecture.md`](./architecture.md): layers, ports, dependency direction, and boundaries
- [`tech-selection.md`](./tech-selection.md): chosen technologies and quality tooling
- [`er-draft.md`](./er-draft.md): canonical data model and ER specification
- [`executable-specs.md`](./executable-specs.md): testable constraints, acceptance checks, and executable-spec policy
- [`lint-policy.md`](./lint-policy.md): what lint manages, what it does not manage, and tool ownership
- [`exec-plans/active/`](./exec-plans/active/README.md): plans that are not yet complete
- [`exec-plans/completed/`](./exec-plans/completed/README.md): finished plans and outcomes
- [`../4humans/quality-score.md`](../4humans/quality-score.md): human-facing quality posture and missing coverage
- [`../4humans/tech-debt-tracker.md`](../4humans/tech-debt-tracker.md): human-facing unresolved debt and cleanup backlog
- [`references/`](./references/index.md): curated reference index and external material policy

## Choose The Right Record

- task-local な詳細設計や一時的な実装判断は plan に置く。実装 task で必要な `UI` / `Scenario` / `Logic` も active plan の section に含める
- 完了後も保持すべき詳細な振る舞い、制約、受け入れ条件は [`executable-specs.md`](./executable-specs.md) と対応する tests / acceptance checks / validation commands に昇格する
- Requirement or product boundary changed: update [`spec.md`](./spec.md)
- Dependency rule or layering changed: update [`architecture.md`](./architecture.md)
- Technology decision changed: update [`tech-selection.md`](./tech-selection.md)
- Data structure or entity relationship changed: update [`er-draft.md`](./er-draft.md)
- Detailed behavior or constraint changed: update [`executable-specs.md`](./executable-specs.md) and the corresponding tests or acceptance checks
- Lint の責務範囲、allowlist 方針、tool ownership を変えた: update [`lint-policy.md`](./lint-policy.md)
- Work is non-trivial and not yet finished: create a plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work is finished: move the plan into [`exec-plans/completed/`](./exec-plans/completed/README.md)
- Workflow or role confusion keeps recurring: update [`../.codex/`](../.codex/README.md) or the relevant file in `../.codex/agents/` or `../.codex/skills/`
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
- The legacy raw API dumps still live under `docs/api-refrences/`. Keep them until they are migrated,
  but do not add new files there.
- Workflow skills live under [`../.codex/skills/`](../.codex/skills/).
