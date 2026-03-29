# Skill Permissions Zero Trust

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/*/references/permissions.json`, `.codex/README.md`, `.codex/skills/directing-implementation/SKILL.md`, `AGENTS.md`

## Request Summary

- 全 skill に zero-trust 風の role concept を導入する。
- 各 skill 配下の `references/permissions.json` に `allowed_actions`、`forbidden_actions`、`expected_outputs`、`stop_conditions` を追加する。
- `AGENTS.md` に、skill の権限外は実行しないこと、曖昧な場合は停止して handoff することを明記する。

## Decision Basis

- handoff contract と同じ `references/` 配下に権限制約を置くと、入力契約と権限境界を同じ探索導線で確認できる。
- zero-trust 的な運用では、「何をしてよいか」よりも「いつ止まるか」を固定する方が誤動作を減らせる。
- skill ごとに stop condition を明文化しておくと、lane 切り替えや human handoff の判断を再現しやすい。

## UI

- N/A

## Scenario

- 各 skill は着手時に `references/permissions.json` を見れば、許可された操作、禁止された操作、期待される返却、停止条件を判断できる。
- `AGENTS.md` を読んだエージェントは、権限外の作業を始めず、曖昧な依頼では stop and handoff を選ぶ。

## Logic

- `permissions.json` は reference artifact として置き、実行時の role boundary を補助する。
- 既存の handoff JSON は維持し、permissions は skill-local role contract として並置する。
- `directing-implementation` の既存 design harness 不整合も同じ変更で解消する。

## Implementation Plan

- active plan を追加する。
- 各 skill の責務に合わせた `permissions.json` を作成する。
- global guidance として `.codex/README.md` と `AGENTS.md` を更新する。
- `directing-implementation/SKILL.md` に live 正本へ戻さない artifact 制約を補う。

## Acceptance Checks

- 各 skill に `references/permissions.json` がある。
- 各 `permissions.json` に `allowed_actions`、`forbidden_actions`、`expected_outputs`、`stop_conditions` がある。
- `AGENTS.md` に権限外作業禁止と曖昧時の stop / handoff が明記される。
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `AGENTS.md`
- `.codex/README.md`
- `.codex/skills/*/references/permissions.json`

## Outcome

- 16 個の live skill すべてに `references/permissions.json` を追加し、`allowed_actions`、`forbidden_actions`、`expected_outputs`、`stop_conditions` を定義した。
- `diagramming-plantuml` にも `references/` を作成し、他 skill と同じ探索導線にそろえた。
- `AGENTS.md` に skill 権限外の作業禁止と、曖昧時の stop / handoff を追加した。
- `.codex/README.md` に `permissions.json` を role contract として扱うことを追加した。
- `.codex/skills/directing-implementation/SKILL.md` に live 正本へ戻さない artifact 制約と曖昧時 handoff を補い、design harness の既存不整合も解消した。
- `powershell -File scripts/harness/run.ps1 -Suite structure`、`design`、`all` はすべて pass した。
