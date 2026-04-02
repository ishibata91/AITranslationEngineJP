# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `AGENTS.md`, `docs/index.md`, `docs/lint-policy.md`, `.codex/workflow.md`, `.codex/skills/*implementation*`, `scripts/harness/`, `.codex/skills/directing-implementation/scripts/`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `pwsh` / `powershell` 前提を廃止し、現行のハーネス入口と Sonar open-issue gate を Python 直呼びへ寄せる。

## Decision Basis

- 現行のハーネス本体はすでに Python 実装であり、`.ps1` は薄い wrapper に留まる。
- 実行環境に `powershell` が存在しないため、現行契約と実行可能入口が乖離している。
- 現行の workflow / lint policy / AGENTS / docs index に PowerShell 前提が残っている。

## Owned Scope

- `AGENTS.md`
- `docs/index.md`
- `docs/lint-policy.md`
- `.codex/workflow.md`
- `.codex/skills/directing-implementation/SKILL.md`
- `.codex/skills/implementing-frontend/SKILL.md`
- `.codex/skills/implementing-backend/SKILL.md`
- `.codex/skills/reviewing-implementation/SKILL.md`
- `scripts/harness/README.md`
- `scripts/harness/check_structure.py`
- `scripts/harness/run.py`
- `scripts/harness/check_structure.py`
- `.codex/skills/directing-implementation/scripts/get-open-sonar-issues.py`

## Out Of Scope

- `docs/exec-plans/completed/` の履歴書き換え
- Sonar CLI 自体の仕様変更

## Dependencies / Blockers

- `sonar` CLI が未導入の環境では open-issue helper の実行確認は限定される。

## Parallel Safety Notes

- 現行契約と実行入口の同期だけを扱う。product code と UI には触れない。

## UI

- なし

## Scenario

- エージェントと human は `powershell -File ...` ではなく Python script を直接呼ぶ。
- implementation lane の Sonar open-issue gate も Python helper を使う。

## Logic

- structure harness の required path を Python source 中心へ更新する。
- PowerShell wrapper を削除し、現行 contract の正本からも参照を外す。

## Implementation Plan

- current source と current contract に残る `powershell` 参照を Python 直呼びへ置換する。
- Sonar open-issue helper を `.ps1` から `.py` へ移植する。
- 不要になった harness / Sonar の `.ps1` wrapper を削除する。

## Acceptance Checks

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- 現行 contract から `powershell -File scripts/harness/run.ps1` が消えていること
- structure / design harness が Python 入口前提で pass すること
- full harness の結果と失敗理由

## 4humans Sync

- なし

## Outcome

- harness の現行入口を `python3 scripts/harness/run.py --suite ...` に統一した。
- Sonar open-issue helper を `.ps1` から `.py` へ移植し、現行 workflow 契約を Python 直呼びへ更新した。
- `scripts/harness/*.ps1` と Sonar helper の `.ps1` を削除した。
- `python3 scripts/harness/run.py --suite structure` と `--suite design` は pass した。
- `python3 scripts/harness/run.py --suite all` は既存どおり `cargo: command not found` で execution suite が fail した。
