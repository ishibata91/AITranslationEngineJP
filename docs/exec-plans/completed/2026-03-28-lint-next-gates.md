# Lint Next Gates

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/lint-policy.md`

## Request Summary

- lint 文書に、次に gate 化すべき lint 項目だけを追記する

## Decision Basis

- 既存の `lint-policy.md` は現行責務の整理まではできている
- 次に何を lint gate 化するかを同じ文書内で見えるようにしたい

## Why Light Flow Applies

- 変更対象は `docs/lint-policy.md` のみで、既存方針の補足に限定できる
- blocking unknown がなく、短い plan で判断を固定できる

## Short Plan

- `docs/lint-policy.md` に recommended next lint gates を追記する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- 追記内容が lint の責務範囲に限定されていること
- harness が通ること

## Reroute Trigger

- lint と test / acceptance checks の責務境界まで変更する必要が出る

## Docs Sync

- `docs/lint-policy.md`

## Record Updates

- `docs/lint-policy.md`

## Outcome

- `docs/lint-policy.md` に `Recommended Next Gates` を追加し、次に gate 化すべき lint 項目だけを列挙した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
