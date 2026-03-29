# Replace Docs Sync With 4humans Sync

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/`, `AGENTS.md`, `docs/index.md`, `docs/exec-plans/templates/`, `scripts/harness/`

## Request Summary

- 通常の impl / fix lane から `docs/` を close 条件として更新しない運用へ切り替える。
- close 条件の既定を `4humans sync` に変更する。
- `docs/` 更新は human が先に行い、専用 skill を human が直接起動した時だけ許可する。

## Decision Basis

- 現在の live workflow では `docs sync` が close 条件として広く書かれていた。
- 既存契約では `docs/` は恒久仕様と設計、`.codex/` は workflow、詳細挙動は tests / acceptance checks / validation commands が持つ。
- 通常 lane の close 条件として `docs/` 更新を既定にすると、`docs/` の責務より広く触る圧力が残る。

## UI

- N/A

## Scenario

- 通常の impl / fix close は `4humans sync` のみを既定とする。
- workflow 契約変更は `skill-modification` が扱う。
- `docs/` 正本更新は `updating-docs` を human が直接起動した時だけ許可する。

## Logic

- live workflow の `docs sync` 文言を `4humans sync` または `closeout notes` に整理した。
- plan template の `Docs Sync` を `4humans Sync` に変更した。
- docs 更新専用 skill `updating-docs` を追加し、権限を `docs/` の正本更新だけに制限した。
- design / structure harness も新契約を検証するよう更新した。

## Implementation Plan

- live な workflow / skill / reference JSON の `docs sync` 文言を置換した。
- `AGENTS.md` と `docs/index.md` に human-first の docs 更新運用を明記した。
- `.codex/skills/updating-docs/` を追加した。
- validation を実行し、plan を completed へ移した。

## Acceptance Checks

- live 契約の通常 close 条件が `4humans sync` で表現される。
- `Docs Sync` section が `4humans Sync` に置き換わる。
- docs 更新専用 skill が追加され、`docs/` 更新が human 直接起動時だけ許可される。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
- `rg -n "docs sync|Docs Sync|docs_sync" AGENTS.md .codex docs 4humans`
- `rg -n "4humans sync|closeout notes|updating-docs|closeout_" AGENTS.md .codex docs 4humans`

## 4humans Sync

- なし。今回の変更は workflow 契約と harness の更新に限定した。

## Outcome

- live な impl / fix lane の close 条件を `docs sync` から `4humans sync` に切り替えた。
- handoff / permissions / agent contract の `docs_sync_*` を `closeout_*` 系へ寄せた。
- `docs/` 更新専用 skill `updating-docs` を追加し、human-first 運用を `AGENTS.md` と `docs/index.md` に明記した。
- structure / design harness を新契約へ同期した。

## Validation Results

- `python C:\Users\shiba\.codex\skills\.system\skill-creator\scripts\quick_validate.py .codex\skills\updating-docs` failed because the validator read `SKILL.md` with `cp932` and raised `UnicodeDecodeError`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all` failed because an existing unrelated test failed in `scripts/eslint/repository-boundary-plugin.test.mjs`
- `rg -n "docs sync|Docs Sync|docs_sync" AGENTS.md .codex docs 4humans` still matches historical records under `docs/exec-plans/completed/` and this task's plan history, but no live workflow contract under `AGENTS.md` or `.codex/` remains on the old term
