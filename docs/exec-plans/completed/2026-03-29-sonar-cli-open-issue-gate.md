# Sonar CLI Open Issue Gate

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/README.md`, `.codex/workflow.md`, `.codex/workflow_activity_diagram.puml`, `.codex/skills/directing-implementation/`, `.codex/skills/implementing-frontend/`, `.codex/skills/implementing-backend/`, `.codex/skills/reviewing-implementation/`

## Request Summary

- SonarQube MCP issue read が不安定なので、implementation lane の issue gate を Sonar CLI へ移行する。
- `sonar list issues --project ishibata91_AITranslationEngineJP --format json` の結果には `status: CLOSED` が含まれるため、gate 対象は `OPEN` のみへ明示的に絞る。

## Decision Basis

- live workflow の正本は `.codex/` であり、今回の変更対象は workflow 契約と skill 契約である。
- `docs/` の恒久仕様更新は human-first なので、今回の lane 変更では触らない。
- Sonar CLI の生 JSON を skill ごとに都度解釈させるより、workflow 専用 helper script へ集約した方が再利用しやすい。

## UI

- Use only when the task changes screen structure, presentation, or interaction flow.

## Scenario

- implementation lane は `sonar-scanner` 実行後に helper script で open issue 一覧を取得する。
- helper script は Sonar CLI JSON から `status == OPEN` だけを返す。
- review は open issue が解消した後だけ開始する。

## Logic

- helper script は `component` または `components[].path` から repo path を補完し、owned scope の判定に使える形で返す。
- Sonar CLI の historical issue を gate 対象へ混ぜないよう、`status` を唯一の gate 条件として扱う。

## Implementation Plan

- active plan を追加する。
- Sonar CLI の open issue 抽出 helper script を `.codex/skills/directing-implementation/scripts/` に追加する。
- impl lane の `.codex` 文書を `SonarQube MCP` から `Sonar CLI open issue gate` へ更新する。

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- helper script が `status: CLOSED` issue を除外する実行結果
- `.codex` workflow 文書の Sonar gate 記述差分
- harness 実行結果

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- `.codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1` を追加し、`sonar list issues --project ishibata91_AITranslationEngineJP --format json` から `status == OPEN` だけを返す形へ統一した。
- `.codex` の impl lane workflow を `SonarQube MCP` から `Sonar CLI open issue gate` へ更新した。
- validation passed: `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP`, `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths src/application`, `powershell -File scripts/harness/run.ps1 -Suite all`.
- current Sonar evidence: source JSON の `totalIssues` は 2 件だったが、どちらも `status: CLOSED` のため helper script は `openIssueCount: 0` を返した。
- `4humans/quality-score.md` と `4humans/tech-debt-tracker.md` は既存差分があるため今回は更新していない。
