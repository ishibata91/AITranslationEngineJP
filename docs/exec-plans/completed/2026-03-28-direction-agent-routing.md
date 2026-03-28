# Direction Agent Routing

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/directing-implementation/SKILL.md`, `.codex/skills/directing-fixes/SKILL.md`

## Request Summary

- `directing-implementation` の few-shot に合わせて、direction 系 skill 呼び出しへ agent 指定を明示する。

## Decision Basis

- 現状は `directing-implementation` の一部 step だけ agent 指定が入り、lane 全体の handoff 契約としては不揃い。
- `.codex/README.md` の agent 契約に対応する形で、direction skill 側の handoff 先も明示したほうが再現しやすい。

## UI

- N/A

## Scenario

- `directing-implementation` と `directing-fixes` を読めば、各 helper skill をどの agent role で呼ぶかが分かる。

## Logic

- 実装 lane は `ctx_loader` / `workplan_builder` / `test_architect` / `implementer` / `review_cycler` を使う。
- fix lane は `ctx_loader` / `fault_tracer` / `log_instrumenter` / `test_architect` / `implementer` / `review_cycler` を使う。

## Implementation Plan

- active plan を追加する。
- `directing-implementation` の Required Workflow を agent role 付きで統一する。
- `directing-fixes` の Required Workflow も同じ書式で agent role を追加する。

## Acceptance Checks

- 両 direction skill の helper skill 呼び出しに agent role が書かれている。
- role 名が `.codex/README.md` の agent 契約と整合している。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Docs Sync

- `.codex/skills/directing-implementation/SKILL.md`
- `.codex/skills/directing-fixes/SKILL.md`

## Outcome

- `directing-implementation` の helper skill 呼び出しに `ctx_loader`、`workplan_builder`、`test_architect`、`implementer`、`review_cycler` を明記した。
- `directing-fixes` の helper skill 呼び出しに `ctx_loader`、`fault_tracer`、`log_instrumenter`、`test_architect`、`implementer`、`review_cycler` を明記した。
- `design` harness は pass した。
