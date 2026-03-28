# Readiness And Repo-Owned Harness

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `.codex/README.md`, `.codex/agents/architect.toml`, `.codex/skills/architect-direction/SKILL.md`, `.codex/skills/light-direction/SKILL.md`, `.codex/skills/workflow-gate/SKILL.md`, `scripts/harness/check-design.ps1`, `scripts/harness/check-structure.ps1`, `scripts/harness/check-execution.ps1`

## Request Summary

- スキルと harness 契約に `validation readiness check` と `repo-owned files only` を追加する

## Decision Basis

- 修正対象は workflow 契約と harness 契約に限定される
- readiness は gate 後追いではなく Architect 入口で固定したい
- repo-owned files only は gitignore ではなく harness 自身の対象選定契約として持たせる

## Why Light Flow Applies

- 変更対象は `.codex` と harness に限定され、product 仕様や実装設計の再判断を含まない
- 契約の方向性は固定済みで、short plan で実装判断を固定できる

## Short Plan

- `.codex/README.md` と Architect / light / gate skill に readiness と repo-owned harness ルールを追加する
- `architect.toml` に implementation handoff 前の readiness 責務を追加する
- design harness に新契約の存在確認を追加する
- structure / execution harness に external / generated path 除外を明示する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- Architect 系契約で readiness の初回判定責務が gate から分離されていること
- gate 契約で prerequisite gap と通常の check failure が区別されていること
- structure / execution harness が repo-owned files only の方針を反映していること

## Reroute Trigger

- product docs や plan template まで同時に再設計する必要が出る
- readiness を Architect ではなく別 role に持たせる再設計が必要になる

## Docs Sync

- `.codex/README.md`
- `.codex/agents/architect.toml`
- `.codex/skills/architect-direction/SKILL.md`
- `.codex/skills/light-direction/SKILL.md`
- `.codex/skills/workflow-gate/SKILL.md`

## Record Updates

- `docs/exec-plans/active/2026-03-28-readiness-and-repo-owned-harness.md`

## Outcome

- `.codex/README.md` に `validation readiness check` と `repo-owned files only` の共通ルールを追加した
- `architect-direction`、`light-direction`、`workflow-gate` に readiness の責務分離と prerequisite gap の扱いを追加した
- `architect.toml` に implementation handoff 前の readiness 責務を追加した
- design harness に新契約の存在確認を追加した
- structure / execution harness に external / generated path の除外方針を追加した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Remaining Gaps

- `powershell -File scripts/harness/run.ps1 -Suite all` は引き続き `cargo` 未導入の環境差分で失敗する
- repo-owned files only の方針は structure / execution harness には反映したが、今後 harness を増やす時も同じ除外方針を継続する必要がある
