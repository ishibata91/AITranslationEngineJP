# Sonar Quality Gate Loop

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: `scripts/harness/check-execution.ps1`, `scripts/harness/README.md`, `package.json`, `sonar-project.properties`, `.codex/skills/directing-implementation/SKILL.md`, `.codex/skills/implementing-frontend/SKILL.md`, `.codex/skills/implementing-backend/SKILL.md`, `.codex/skills/reviewing-implementation/SKILL.md`, `docs/lint-policy.md`, `docs/tech-selection.md`

## Request Summary

- `sonar-scanner` をプロジェクトルートで実行し、その後に SonarQube MCP で issue を取得する flow を quality gate にする。
- execution harness では `lint` の次に `sonar-scanner` を実行する。
- 実装 lane では Sonar issue が残る限り implementing skill へ差し戻して修正を継続し、issue がなくなってから review / close に進む。

## Decision Basis

- `sonar-project.properties` に project key `ishibata91_AITranslationEngineJP` が定義済みであり、`sonar-scanner` と SonarQube MCP が使用可能である。
- harness script は repo root 直下の `package.json` / `Cargo.toml` を自動検出して標準 command を実行するため、scanner gate は execution harness へ追加するのが最短である。
- SonarQube MCP は agent workflow からのみ利用できるため、harness 自体には埋め込まず implementation skill 契約で remediation loop を明示する。

## UI

- なし

## Scenario

- 実装後に execution harness が `lint` の次で `sonar-scanner` を実行する。
- `directing-implementation` は scanner 実行後に SonarQube MCP で open issue を取得する。
- open issue がある場合は `implementing-frontend` または `implementing-backend` へ差し戻し、修正後に scanner と issue 取得を再実行する。
- Sonar issue が解消した時だけ `reviewing-implementation` と close に進む。

## Logic

- execution harness に Sonar scanner step と repo-owned path 対象の除外方針を追加する。
- implementation skills の check 契約に Sonar scanner 実行と SonarQube MCP issue 確認を追加する。
- current docs の `SonarQube CLI` / `sonar verify` 前提を `sonar-scanner` + MCP issue gate へ更新する。
- Sonar の役割は repo 固有の責務境界 lint ではなく、code smell / complexity / security hotspot の補完層として表現する。

## Implementation Plan

- execution harness と package script を `sonar-scanner` 前提へ差し替える。
- obsolete な `sonar verify` wrapper を削除し、必要なら scanner wrapper を追加する。
- implementation lane skill docs に Sonar issue remediation loop を追記する。
- lint policy と tech selection を scanner + MCP issue gate 契約へ同期する。

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite execution`

## Required Evidence

- execution harness が `lint` 後に `sonar-scanner` を実行するログ
- SonarQube MCP で project issue を取得できる evidence
- 更新後 docs / skill contract が scanner + MCP loop を示す差分

## Docs Sync

- `docs/lint-policy.md`
- `docs/tech-selection.md`
- `scripts/harness/README.md`
- implementation lane skills

## Outcome

- execution harness は `lint` の直後に `sonar-scanner` を repo root で実行する形へ更新した
- implementation lane は SonarQube MCP の open issue を close 条件として扱う契約へ更新した
- `src/test/setup.ts` の Sonar issue は解消し、SonarQube MCP で open issue 0 件を確認した
- validation 時点で execution harness は `scripts/eslint/repository-boundary-plugin.test.mjs` の既存 failure により完走していない
