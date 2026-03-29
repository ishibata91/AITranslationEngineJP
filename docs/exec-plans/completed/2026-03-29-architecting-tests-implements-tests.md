# Architecting Tests Implements Tests

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/architecting-tests/`, `.codex/agents/test_architect.toml`, `.codex/skills/directing-*/`, `.codex/README.md`, `.codex/workflow.md`

## Request Summary

- `architecting-tests` がテスト設計だけで止まらず、必要な test 実装まで担うようにしたい。
- product code の恒久実装責務は既存 implementing skill に残したい。

## Decision Basis

- 現在の `architecting-tests` は「設計」中心の文面で、先行テストを実ファイルへ反映する責務が曖昧だった。
- TDD の入口を lane 契約として明確にするには、test / fixture の最小実装を `architecting-tests` に寄せ、product code の実装責務とは分離したまま handoff を更新するのがよい。

## UI

- N/A

## Scenario

- impl lane では `architecting-tests` が failing tests と fixture を先に実ファイルへ反映してから product 実装へ handoff する。
- fix lane では回帰 test を先に実ファイルへ反映してから恒久修正へ handoff する。

## Logic

- `architecting-tests` の description、overview、workflow、rules、permissions を test 実装込みに更新する。
- `test_architect` agent 契約を test / fixture / helper files の編集責務に合わせる。
- `directing-*` と handoff JSON の completion signal に `touched_test_files` を追加する。
- workflow overview と `.codex/README.md` の lane 説明を同期する。

## Implementation Plan

- `architecting-tests` 本体と permissions を更新する。
- `test_architect` agent と UI metadata を更新する。
- `directing-*` の workflow 記述と architecting-tests handoff JSON を同期する。
- 変更記録を completed plan に残す。

## Acceptance Checks

- `architecting-tests` の責務に test / fixture の最小実装が明記されている。
- `test_architect` agent 契約が test / fixture / helper files の編集責務を持つ。
- `directing-*` と reference JSON が `touched_test_files` を前提に handoff できる。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/skills/architecting-tests/`
- `.codex/agents/test_architect.toml`
- `.codex/skills/directing-implementation/`
- `.codex/skills/directing-fixes/`
- `.codex/README.md`
- `.codex/workflow.md`

## Outcome

- `architecting-tests` を test 設計と test 実装の両方を担う skill として更新した。
- `test_architect` agent 契約を、必要な test / fixture / helper files の編集責務へ広げた。
- impl / fix lane の handoff contract に `touched_test_files` を追加した。

## Validation Results

- `python C:\Users\shiba\.codex\skills\.system\skill-creator\scripts\quick_validate.py .codex\skills\architecting-tests` failed because the validator read `SKILL.md` with `cp932` and raised `UnicodeDecodeError`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all` failed because an existing unrelated test failed in `scripts/eslint/repository-boundary-plugin.test.mjs`
