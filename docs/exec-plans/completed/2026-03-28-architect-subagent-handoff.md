# Architect Subagent Handoff

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `.codex/` workflow documents that define architect handoff behavior

## Request Summary

- `architect-direction` が他 skill へ移るときにサブエージェント起動前提であることを明示し、勝手に実装し始めない契約にしたい。

## Decision Basis

- 問題はプロダクト仕様ではなく workflow 契約にある。
- 変更対象は `.codex` の入口 skill と architect role 契約に限定できる。

## Why Light Flow Applies

- 要求は単一責務の workflow 文言修正である。
- blocking unknown はなく、更新箇所と受け入れ条件を短い plan で固定できる。

## Short Plan

- architect 系の正本文書で、他 skill へ handoff する際はサブエージェントを起動する契約を追加する。
- `light-direction` にも同じ handoff 契約を加え、軽量フローでも main agent が直接別 skill を混在実行しないようにする。
- plan 完了後に full harness を実行し、完了済み plan へ移動する。

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `.codex/README.md` に sub-agent handoff の原則があること
- `architect-direction` と `light-direction` に sub-agent 起動ルールがあること
- `architect.toml` に direct implementation 禁止と sub-agent handoff 優先があること

## Reroute Trigger

- 既存 harness が sub-agent wording を前提にしていて追加修正が広がると判明した場合

## Docs Sync

- `.codex/README.md`
- `.codex/skills/architect-direction/SKILL.md`
- `.codex/skills/light-direction/SKILL.md`
- `.codex/agents/architect.toml`

## Record Updates

- 完了後、この plan を `docs/exec-plans/completed/` へ移し結果を追記する

## Outcome

- `architect-direction` に、他 skill へ移る時は必ず専用サブエージェントを起動する契約を追加した。
- `light-direction` にも同じ handoff ルールを追加し、Architect 本体が `light-work` や `gating-workflow` を直接兼務しないようにした。
- `.codex/README.md` と `architect.toml` を同期し、workflow 正本と role 契約の両方で同じ前提を明文化した。

## Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

