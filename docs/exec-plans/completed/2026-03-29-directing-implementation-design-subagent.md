# Directing Implementation Design Subagent

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/agents/task_designer.toml`, `.codex/skills/designing-implementation/`, `.codex/skills/directing-implementation/`, `.codex/README.md`, `.codex/workflow.md`, `.codex/workflow_activity_diagram.puml`, `docs/index.md`, `docs/core-beliefs.md`, `scripts/harness/check-structure.ps1`, `scripts/harness/check-design.ps1`

## Request Summary

- `directing-implementation` が active plan の `UI` / `Scenario` / `Logic` を自前で埋める責務を別 skill へ切り出す。
- 切り出し先は sub-agent として呼び出せるようにする。
- task-local design を active plan の中に残したまま、impl lane の live workflow を更新する。

## Decision Basis

- `UI` / `Scenario` / `Logic` の補完は task-local design 固定であり、実装順整理を担う `planning-implementation` とは責務が異なる。
- `directing-implementation` が直接設計を書くより、専用 skill と agent role に切り出した方が role boundary と stop condition を明示しやすい。
- task-local design の正本は引き続き active exec-plan に置き、別 artifact へ戻さない。

## UI

- N/A

## Scenario

- `directing-implementation` は active plan 作成直後に task-local design が必要か判断し、必要時だけ `designing-implementation` を sub-agent で呼び出す。
- `designing-implementation` は active plan の `UI` / `Scenario` / `Logic` を補完して direction へ返す。

## Logic

- 新しい agent role `task_designer` を追加する。
- 新しい skill `designing-implementation` は active plan の `UI` / `Scenario` / `Logic` だけを更新対象にする。
- `directing-implementation` から `designing-implementation` への handoff 契約と返却契約を追加する。

## Implementation Plan

- active plan を追加する。
- `task_designer` agent 契約と `designing-implementation` skill 一式を追加する。
- `directing-implementation` と workflow docs を新しい handoff 順序へ更新する。
- structure harness と design harness が新 skill / agent の存在を確認できるようにする。

## Acceptance Checks

- `.codex/agents/task_designer.toml` がある。
- `.codex/skills/designing-implementation/` に `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` がある。
- `directing-implementation` が `designing-implementation` を sub-agent handoff 先として参照している。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/workflow_activity_diagram.puml`
- `.codex/skills/directing-implementation/`
- `.codex/skills/designing-implementation/`
- `docs/index.md`
- `docs/core-beliefs.md`

## Outcome

- `task_designer` agent を追加し、task-local design 専用の sub-agent role を定義した。
- `designing-implementation` skill を追加し、active exec-plan の `UI` / `Scenario` / `Logic` を補完する責務を切り出した。
- `directing-implementation` は `designing-implementation` を先に spawn してから distill / plan / test / implement に進む flow へ更新した。
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml`、`docs/index.md`、`docs/core-beliefs.md` を新 workflow に同期した。
- `scripts/harness/check-structure.ps1` と `scripts/harness/check-design.ps1` を更新し、新 skill / agent の存在と参照を検証対象に加えた。
- `powershell -File scripts/harness/run.ps1 -Suite structure`、`design`、`all` はすべて pass した。
