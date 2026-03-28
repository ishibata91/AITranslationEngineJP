# Research Agent Routing

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `.codex/` documents that define heavy-flow investigation routing

## Request Summary

- Heavy flow の調査委任先を repo 内の `Research` に固定し、組み込み `explorer` を代替に使わないことを明示したい。

## Decision Basis

- 問題は workflow routing の文言不足であり、プロダクト仕様変更ではない。
- 変更対象は architect 系の workflow 正本と role 契約に限定できる。

## Why Light Flow Applies

- 単一責務の文書修正で、受け入れ条件が固定済みである。
- blocking unknown はなく、追加先も明確である。

## Short Plan

- Heavy flow の investigation は `.codex/agents/research.toml` の Research サブエージェントを標準とし、組み込み `explorer` を標準 routing に使わないことを README に追記する。
- `architect-direction` と `architect.toml` に同じ routing 制約を加える。
- design / full harness を実行し、完了後に completed へ移す。

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `.codex/README.md` に Research routing の標準があること
- `architect-direction` に `explorer` 非標準の handoff ルールがあること
- `architect.toml` に Research sub-agent 優先の制約があること

## Reroute Trigger

- 他 skill や harness が `explorer` 前提になっていて、修正範囲が `.codex` 以外へ広がる場合

## Docs Sync

- `.codex/README.md`
- `.codex/skills/architect-direction/SKILL.md`
- `.codex/agents/architect.toml`

## Record Updates

- 完了後、この plan を `docs/exec-plans/completed/` へ移し結果を記録する

## Outcome

- Heavy flow の調査委任先を `.codex/agents/research.toml` の Research に明示的に固定した。
- `architect-direction` に、組み込み `explorer` を標準 investigation lane に使わないことを追加した。
- `architect.toml` と `.codex/README.md` を同期し、Architect の routing 契約と workflow 正本の表現を揃えた。

## Evidence

- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
