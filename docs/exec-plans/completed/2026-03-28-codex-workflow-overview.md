# Codex Workflow Overview

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow.md`, `.codex/README.md`

## Request Summary

- `.codex/workflow_activity_diagram.puml` をベースに、skill の鳥瞰図を説明する `workflow.md` を作る。

## Decision Basis

- 既存の activity diagram は flow の順序を示せるが、各 skill の責務と分岐条件までは読み取りづらい。
- `.codex/README.md` は存在しない `workflow_spec.md` を参照しており、導線が壊れている。
- workflow overview は `.codex/` 内に置き、diagram と live workflow 説明を近接させる。

## UI

- N/A

## Scenario

- User が `.codex/README.md` から workflow の全体像へ辿れる。
- `workflow.md` が impl lane / fix lane / reroute / docs sync の関係を説明する。

## Logic

- diagram の順序を正本として使い、説明文では skill の責務、分岐、close 条件を補足する。

## Implementation Plan

- active plan を追加する。
- `.codex/workflow_activity_diagram.puml` の内容を文章化した `.codex/workflow.md` を作る。
- `.codex/README.md` の参照先を `workflow.md` に更新する。

## Acceptance Checks

- `.codex/workflow.md` が存在し、impl lane / fix lane の両方を説明している。
- `.codex/README.md` から `workflow.md` へ到達できる。
- structure / design / all harness が通る。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow.md`

## Outcome

- `.codex/workflow.md` を新規作成し、diagram ベースの workflow overview を文章化した。
- `.codex/README.md` の壊れた参照先を `workflow.md` に更新した。
- `structure` と `design` harness は pass した。
- `all` harness は execution suite 内の `cargo` コマンド不足で失敗したが、`npm run lint`、`npm run test`、`npm run build` は pass した。
