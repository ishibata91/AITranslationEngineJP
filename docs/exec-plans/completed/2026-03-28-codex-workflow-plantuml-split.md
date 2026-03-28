# Codex Workflow PlantUML Split

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow_spec.md`, `.codex/workflow_activity.puml`, `.codex/workflow_components.puml`, `.codex/README.md`

## Request Summary

- workflow 図を activity diagram と component diagram に分離し、縦方向で読みやすくする。

## Decision Basis

- 1 ファイル 2 図だと読みづらく、PlantUML の図種も混ざる。
- activity diagram と component diagram を別ファイルに分けると、レンダリングと参照が明確になる。
- 既存の workflow index は、2 つの図への導線として残す。

## UI

- N/A

## Scenario

- N/A

## Logic

- N/A

## Implementation Plan

- activity diagram を縦向きで新規ファイル化する。
- component diagram を縦向きで新規ファイル化する。
- `.codex/workflow_spec.md` は 2 図の索引にする。
- `.codex/README.md` から workflow index へ辿れるように保つ。

## Acceptance Checks

- activity diagram と component diagram が別ファイルで存在する。
- どちらも縦方向で描画される前提になっている。
- `.codex/workflow_spec.md` が 2 図へのリンクを持つ。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow_spec.md`
- `.codex/workflow_activity.puml`
- `.codex/workflow_components.puml`

## Outcome

- `.codex/workflow_activity.puml` と `.codex/workflow_components.puml` に分離した。
- どちらも縦方向になるようにし、`README` は `workflow_spec.md` の索引に戻した。
- `workflow_spec.md` は activity / component の 2 図への入口に整理した。
- structure / design harness は再実行して pass した。
