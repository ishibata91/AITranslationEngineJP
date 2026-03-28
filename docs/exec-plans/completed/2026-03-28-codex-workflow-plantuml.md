# Codex Workflow PlantUML Map

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow_spec.puml`, `.codex/workflow_spec.md`, `.codex/README.md`

## Request Summary

- `.codex` の workflow map を Mermaid から PlantUML に切り替え、ファイルとして残す。

## Decision Basis

- Mermaid は読みにくいので、同じ workflow を PlantUML で表現し直す。
- `.codex/workflow_spec.puml` を diagram source にして、`.codex/README.md` から辿れるようにする。
- 既存の `.codex/workflow_spec.md` は、PlantUML への案内板に縮める。

## UI

- N/A

## Scenario

- N/A

## Logic

- N/A

## Implementation Plan

- workflow overview と component map を PlantUML で書き直す。
- skill / agent のクリック可能リンクを図に残す。
- `.codex/README.md` の workflow link を PlantUML source に向ける。
- `.codex/workflow_spec.md` は簡潔な index にする。

## Acceptance Checks

- `.codex/workflow_spec.puml` が存在し、workflow overview と component map を含む。
- skill / agent のリンクが実在ファイルを指す。
- `.codex/README.md` から PlantUML source へ辿れる。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow_spec.md`
- `.codex/workflow_spec.puml`

## Outcome

- `.codex/workflow_spec.puml` を追加し、workflow overview と component map を PlantUML 化した。
- `.codex/workflow_spec.md` を簡潔な index に縮め、README から PlantUML source へ導線を張った。
- workflow overview は legacy activity syntax に切り替え、`start` 系の構文エラーを避ける形へ修正した。
- structure / design harness は再実行して pass した。
