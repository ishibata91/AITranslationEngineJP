# Skill Gerund Renaming

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/`, `.codex/README.md`, `.codex/workflow.md`, `.codex/workflow_activity_diagram.puml`, `AGENTS.md`, `docs/`, `4humans/`, `scripts/harness/`

## Request Summary

- skill 名をすべて動名詞ベースへ統一する。

## Decision Basis

- 現行 skill 名は `direction`、`review`、`report` など名詞終わりが混在しており、命名規則が揃っていない。
- 参照は `.codex/README.md`、`AGENTS.md`、`docs/`、`4humans/`、harness に跨るため、directory 名と文書参照を同時に更新する必要がある。
- 役割の意味は維持しつつ、`<gerund>-<domain>` の形に揃える。

## UI

- N/A

## Scenario

- リポジトリ内の skill 参照が動名詞ベースの新名称へ揃う。
- workflow と harness が rename 後も同じ skill を指せる。

## Logic

- 実 skill が存在するディレクトリだけを rename 対象にする。
- historical plan は記録として残しつつ、リンク切れや current source-of-truth の混乱を避けるため本文中の skill 参照は更新する。

## Implementation Plan

- rename 対応表を定義する。
- `.codex/skills/*` の実 skill ディレクトリ名と `SKILL.md` frontmatter を更新する。
- README / workflow / AGENTS / docs / 4humans / harness の参照を一括更新する。
- structure / design / all harness を再実行して整合を確認する。

## Acceptance Checks

- すべての実 skill directory が動名詞ベース名になっている。
- `SKILL.md` の `name:` が directory 名と一致している。
- current source-of-truth と harness が新名称を参照している。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/workflow_activity_diagram.puml`
- `AGENTS.md`
- `docs/index.md`
- `docs/core-beliefs.md`
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `scripts/harness/check-structure.ps1`
- `scripts/harness/check-design.ps1`

## Outcome

- 実 skill directory を動名詞ベース名へ rename し、`SKILL.md` の `name:` と見出しを揃えた。
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml`、`AGENTS.md`、`docs/`、`4humans/`、harness の参照を新名称へ更新した。
- fixture 内の旧 skill 名も current naming に合わせて更新した。
- `structure` と `design` harness は pass した。
- `all` harness は execution suite 内の `cargo` コマンド不足で失敗したが、`npm run lint`、`npm run test`、`npm run build` は pass した。
