# Direction Skills Commit Close

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/directing-implementation/`, `.codex/skills/directing-fixes/`, `.codex/README.md`, `.codex/workflow.md`, `.codex/workflow_activity_diagram.puml`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `directing-implementation` と `directing-fixes` の close 条件に commit を含める。
- lane の正本と workflow 鳥瞰図で、review pass 後は `4humans sync` だけでなく commit まで同じ lane の責務として扱う。

## Decision Basis

- 現在の direction 系 skill は `close` までしか明示しておらず、commit を lane 責務として読めない。
- close 条件だけを変えると `.codex/README.md` と `.codex/workflow.md` と図版がずれるため、同じ変更で同期が必要になる。
- 変更は workflow 契約の更新であり product 実装ではないため、`.codex/` 配下と exec-plan の範囲に閉じる。

## Owned Scope

- `directing-implementation` と `directing-fixes` の Required Workflow と description
- direction skill の `references/permissions.json`
- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/workflow_activity_diagram.puml`

## Out Of Scope

- product code、tests、docs 正本の更新
- direction skill 以外の skill に commit 実行責務を広げること

## Dependencies / Blockers

- なし

## Parallel Safety Notes

- `.codex/README.md` と workflow docs は shared contract なので、direction skill の文言と同時に揃える。

## UI

- N/A

## Scenario

- review が `pass` の lane は `4humans sync` を済ませた後、その変更を commit してから close する。

## Logic

- commit は direction 系 skill の close 条件に含める。
- commit 失敗や commit 不可の状態では close 完了とみなさない。

## Implementation Plan

- active plan を追加する。
- `directing-implementation` / `directing-fixes` の Required Workflow、description、permissions を commit 前提へ更新する。
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml` を同じ close 契約へ同期する。
- structure / design / all harness を実行する。
- plan を completed へ移し、結果を記録する。

## Acceptance Checks

- direction skill を読むと review pass 後に commit まで行う責務が分かる。
- permissions 契約でも commit が許可と close 条件に含まれる。
- `.codex/README.md`、`.codex/workflow.md`、図版が同じ close 契約を示す。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## 4humans Sync

- なし。workflow 契約変更のみで、品質評価や未解決負債の更新は不要なら触らない。

## Outcome

- `directing-implementation` と `directing-fixes` の close 条件を commit まで含む契約へ更新した。
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml` を同じ close 契約へ同期した。
- 4humans 系ファイルは更新不要だったため変更していない。

## Validation Results

- `powershell -File scripts/harness/run.ps1 -Suite structure`: pass
- `powershell -File scripts/harness/run.ps1 -Suite design`: pass
- `powershell -File scripts/harness/run.ps1 -Suite all`: pass
