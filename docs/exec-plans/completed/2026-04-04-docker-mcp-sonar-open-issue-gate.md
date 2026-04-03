# Docker MCP Sonar Open Issue Gate

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/README.md`, `.codex/workflow.md`, `.codex/workflow_activity_diagram.puml`, `.codex/skills/*implementation*`, `.codex/skills/directing-implementation/scripts/get-open-sonar-issues.py`, `docs/lint-policy.md`, `docs/tech-selection.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- Sonar open-issue gate を `sonar list issues` 依存から Docker MCP 経由へ切り替える。
- helper script、workflow 契約、`docs/` 正本を同じ前提へ同期する。

## Decision Basis

- `workspace-write` サンドボックス下でも `docker mcp tools call search_sonar_issues_in_projects --gateway-arg=--profile --gateway-arg=codexmcps` は実行できる。
- SonarQube MCP のレスポンスには別 project の issue が混ざるため、helper 側で `project == --project` と `status == OPEN` を必須条件にして絞る必要がある。
- `.codex/workflow_activity_diagram.puml` も live workflow の正本なので、文章だけでなく図版も同期する。

## Owned Scope

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/workflow_activity_diagram.puml`
- `.codex/skills/directing-implementation/SKILL.md`
- `.codex/skills/implementing-frontend/SKILL.md`
- `.codex/skills/implementing-backend/SKILL.md`
- `.codex/skills/reviewing-implementation/SKILL.md`
- `.codex/skills/directing-implementation/scripts/get-open-sonar-issues.py`
- `docs/lint-policy.md`
- `docs/tech-selection.md`

## Out Of Scope

- `sonar-scanner` 自体の実行方法変更
- `docs/exec-plans/completed/` の既存履歴の書き換え
- `codexmcps` profile 構成変更

## Dependencies / Blockers

- Docker CLI と `docker mcp` gateway が利用可能であること
- `codexmcps` profile に `mcp/sonarqube` が登録されていること

## UI

- なし

## Scenario

- implementation lane は `sonar-scanner` 後に Docker MCP 経由で Sonar issue を取得する。
- helper script は指定 project の `OPEN` issue だけを返し、必要なら `--owned-paths` で owned scope に絞る。
- open issue が残る間は implementing skill に差し戻し、0 件の時だけ review に進む。

## Logic

- helper script は `docker mcp tools call search_sonar_issues_in_projects` の stdout から JSON payload を抽出して解釈する。
- `component` の `projectKey:path` 形式から repo path を復元して owned scope 判定に使う。
- `.codex` と `docs/` の Sonar gate 記述を Docker MCP 前提へ統一する。

## Implementation Plan

- `get-open-sonar-issues.py` を Docker MCP 呼び出し実装へ差し替える。
- workflow README、overview、diagram、implementation skills の Sonar gate 表現を更新する。
- `docs/lint-policy.md` と `docs/tech-selection.md` を Docker MCP issue gate 前提へ更新する。

## Acceptance Checks

- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/application`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- helper script が Docker MCP 経由で JSON を返すログ
- `workspace-write` サンドボックス下でも helper が動くログ
- `.codex` と `docs/` の Sonar gate 記述が同期した差分

## 4humans Sync

- なし

## Outcome

- `.codex/skills/directing-implementation/scripts/get-open-sonar-issues.py` を `sonar list issues` 依存から `docker mcp tools call search_sonar_issues_in_projects` ベースへ切り替えた。
- implementation lane の `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml`、関連 implementation skills を Docker MCP Sonar open issue gate 前提へ同期した。
- `docs/lint-policy.md` と `docs/tech-selection.md` を `SonarScanner + Docker MCP SonarQube` 前提へ更新した。
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP`、`python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/application`、`python3 scripts/harness/run.py --suite structure`、`python3 scripts/harness/run.py --suite design`、`python3 scripts/harness/run.py --suite all` は pass した。
