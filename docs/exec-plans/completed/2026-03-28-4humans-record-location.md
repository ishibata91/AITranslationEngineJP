# 4humans Record Location

- workflow: heavy
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `AGENTS.md`, `docs/index.md`, `docs/core-beliefs.md`, `docs/references/index.md`, `docs/exec-plans/completed/README.md`, `scripts/harness/check-structure.ps1`, `scripts/harness/check-design.ps1`, `4humans/*.md`

## Request Summary

- `docs/quality-score.md` と `docs/tech-debt-tracker.md` を `4humans/` へ移した前提で、ハーネスと記録契約を直す

## Investigation Summary

- Facts:
- `4humans/quality-score.md` と `4humans/tech-debt-tracker.md` は存在する
- `AGENTS.md`、`docs/index.md`、`docs/core-beliefs.md`、ハーネスはまだ `docs/` 配下を正本として参照している
- `4humans/*.md` の相対リンクは移動前のままで壊れている
- Options:
- `docs/` に stub を戻す方法と、契約ごと `4humans/` に切り替える方法がある
- Risks:
- ハーネスだけ修正すると `AGENTS.md` と `docs/` の契約が食い違う
- Unknowns:
- なし

## Implementation Plan

- `AGENTS.md` と `docs/` の索引文書に `4humans/` の役割を追記する
- 完了計画・参照索引・`4humans/*.md` のリンクを移動後パスへ修正する
- structure/design harness の required paths と pattern checks を `4humans/` 前提へ切り替える

## Delegation Map

- Research: なし
- Coder: codex が実装
- Worker: なし

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Record Updates

- `AGENTS.md`
- `docs/index.md`
- `docs/core-beliefs.md`
- `docs/references/index.md`
- `docs/exec-plans/completed/README.md`
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `scripts/harness/check-structure.ps1`
- `scripts/harness/check-design.ps1`

## Outcome

- `4humans/` を品質状態と負債整理の正本として契約に反映した
- structure/design harness を `4humans/` 前提へ更新した
- `4humans/*.md` と関連索引のリンク切れを解消した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
