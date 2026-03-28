# ER Decisions

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/er-draft.md`

## Request Summary

- `docs/er-draft.md` の未解消論点 4 件を確定事項へ更新する

## Why Light Flow Applies

- 更新対象は `docs/er-draft.md` のみで、既存の正本文書内の判断を固定する作業に限定できる
- 低リスクで、短い plan だけで編集方針を固定できる

## Short Plan

- `quest_id` と `previous_id` を正規化後の FK 表現へ変更する
- `DIALOGUE_RESPONSE.voicetype` を廃止し、`NPC.voice` に統一する
- `JOB_RECORD` のポリモーフィック参照を採用方針として明記する
- `cells` は独立エンティティにせず `LOCATION` に集約する方針へ更新する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Record Updates

- `docs/er-draft.md`

## Outcome

- `quest_id` と `previous_id` を抽出時正規化の FK として確定した
- `JOB_RECORD` のポリモーフィック参照を採用方針として確定した
- ボイス種別を `NPC.voice` に統一し、`cells` を `LOCATION` 集約として確定した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
