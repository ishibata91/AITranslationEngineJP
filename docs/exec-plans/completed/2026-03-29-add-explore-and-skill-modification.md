# Add Explore And Skill Modification

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/explore/`, `.codex/skills/skill-modification/`, `.codex/README.md`

## Request Summary

- live skill として `explore` と `skill-modification` を追加する。
- `explore` は人間報告用の調査 skill とする。
- `skill-modification` は skill 調整専用 skill とする。

## Decision Basis

- repo 内で常用する調査と skill 調整を live skill として固定すると、routing と権限境界を明示できる。
- `explore` は read-only な事実整理に寄せ、`skill-modification` は `.codex/skills/` 周辺の変更責務を持たせると zero-trust role に合う。

## UI

- N/A

## Scenario

- エージェントは人間向け報告が必要な時に `explore` を使える。
- skill 自体の追加や調整が必要な時に `skill-modification` を使える。

## Logic

- 各 skill は `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` を持つ。
- `explore` は `ctx_loader`、`skill-modification` は `implementer` を agent 契約として使う。

## Implementation Plan

- active plan を追加する。
- `explore` skill 一式を追加する。
- `skill-modification` skill 一式を追加する。
- `.codex/README.md` の helper skill 一覧へ追記する。

## Acceptance Checks

- `explore` と `skill-modification` に `SKILL.md` と `agents/openai.yaml` がある。
- 両 skill に `references/permissions.json` がある。
- `.codex/README.md` から新 skill を辿れる。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/skills/explore/`
- `.codex/skills/skill-modification/`

## Outcome

- `explore` skill を追加し、read-only な人間報告用調査 role を `ctx_loader` で扱う構成にした。
- `skill-modification` skill を追加し、skill 自体の追加や調整を `implementer` で扱う構成にした。
- 両 skill に `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` を追加した。
- `.codex/README.md` の helper skill 一覧へ両 skill を追記した。
- `powershell -File scripts/harness/run.ps1 -Suite structure`、`design`、`all` はすべて pass した。
