# Add Working Light Skill

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/working-light/`, `.codex/README.md`

## Request Summary

- 軽いと分かっているタスク向けの `working-light` skill を追加する。
- 既存 live workflow を壊さず、bounded な軽作業を扱う補助 skill として定義する。

## Decision Basis

- 現在の live workflow には軽作業専用の補助 skill がなく、軽い依頼でも毎回 role を読み替える必要がある。
- 旧 `light-*` artifact を復活させるのではなく、`skill-modification` や `explore` と同じ helper skill として限定的な責務を与えると zero-trust 境界を保てる。

## UI

- N/A

## Scenario

- 軽いと確定している docs / workflow / skill / 小規模コード変更で、scope と validation が明確な依頼に `working-light` を使える。
- 変更が重くなった時は lane owner へ戻して `directing-implementation` か `directing-fixes` へ reroute する。

## Logic

- `working-light` は bounded scope の直接変更と最小 validation 実行だけを扱う。
- 新 skill には `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` をそろえる。
- `.codex/README.md` の helper skill 一覧へ `working-light` を追加する。

## Implementation Plan

- active plan を追加する。
- `working-light` skill の雛形を生成する。
- `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` を repo 方針に合わせて埋める。
- `.codex/README.md` の helper skill 一覧へ追記する。

## Acceptance Checks

- `.codex/skills/working-light/` に `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` がある。
- `working-light` の責務が bounded な軽作業に限定され、重い task の reroute 条件が明記されている。
- `.codex/README.md` から `working-light` を辿れる。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/skills/working-light/`

## Outcome

- `working-light` skill を helper skill として追加した。
- `SKILL.md` で light task の利用条件、出力、reroute 条件を定義した。
- `agents/openai.yaml` と `references/permissions.json` を追加し、UI metadata と role contract をそろえた。
- `.codex/README.md` の helper skill 一覧へ `working-light` を追加した。

## Validation Results

- `python C:\Users\shiba\.codex\skills\.system\skill-creator\scripts\quick_validate.py .codex\skills\working-light`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all` failed because an existing unrelated test failed in `scripts/eslint/repository-boundary-plugin.test.mjs`
