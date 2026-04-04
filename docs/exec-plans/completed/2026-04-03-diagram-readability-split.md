# Diagram Readability Split

- workflow: impl
- status: completed
- lane_owner:
- scope: `4humans/class-diagrams/`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- review 用 class diagram の線重なりを減らし、読みやすい粒度へ分割する。

## Decision Basis

- overview に詳細ノードと詳細 edge を同居させると交差が増え、review 速度が落ちる。
- `diagramming-d2` は主題が混ざるときの分割を許可している。

## Owned Scope

- `docs/exec-plans/active/2026-04-03-diagram-readability-split.md`
- `4humans/class-diagrams/`

## Out Of Scope

- product code の変更
- sequence diagram の主題変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること

## Logic

- 既存の frontend/backend class diagram は overview に落とす。
- 詳細は slice / flow ごとの追加図へ分割する。

## Implementation Plan

- frontend overview と backend overview を簡素化する。
- frontend detail と backend detail の追加図を作る。
- 全 `.d2` を validate し、`.svg` を再生成する。

## Acceptance Checks

- overview 図の edge 数が減り、主題が一つに絞られている。
- 追加した detail 図が slice / flow ごとの説明に使える。
- 全 `.d2` が `d2 validate` を通る。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `d2 validate 4humans/class-diagrams/*.d2`
- `d2 4humans/class-diagrams/*.d2 *.svg`

## 4humans Sync

- `4humans/class-diagrams/`

## Outcome

- `frontend-structure-overview.d2` と `backend-structure-overview.d2` を overview 用へ簡素化した。
- 線の交差を減らすため、job-create、job-list、backend create-job、backend import-xedit の detail class diagram を追加した。
- `4humans/class-diagrams/*.d2` はすべて `d2 validate` を通し、対応する `.svg` を再生成した。
- `python3 scripts/harness/run.py --suite structure` は通過した。
