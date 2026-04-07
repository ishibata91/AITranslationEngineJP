# 4humans Overview D2 Verticalization

- workflow: impl
- status: completed
- lane_owner:
- scope: `docs/exec-plans/active/2026-04-07-4humans-overview-d2-verticalization.md`, `4humans/diagrams/processes/*.d2`, `4humans/diagrams/processes/*.svg`, `4humans/diagrams/structures/*.d2`, `4humans/diagrams/structures/*.svg`, `4humans/diagrams/overview-manifest.json`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `4humans` 配下で横に長すぎる overview D2 図を、16:9 以内の横幅で読める構成へ直す。

## Decision Basis

- `processes-overview-robustness.svg` は 20000px 超の横幅で、overview として読めない。
- `frontend-structure-overview.svg` と `backend-structure-overview.svg` も 16:9 を超えている。
- overview は detail への入口を維持しつつ、横方向ではなく縦方向へ情報量を逃がす方が review しやすい。

## Owned Scope

- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/frontend-processes-overview-robustness.d2`
- `4humans/diagrams/processes/backend-processes-overview-robustness.d2`
- `4humans/diagrams/structures/frontend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/overview-manifest.json`
- 対応する `.svg`

## Out Of Scope

- detail 図全体の横長解消
- `docs/` 正本の更新
- 実装コードの変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること

## Parallel Safety Notes

- `.d2` を source of truth として更新する。
- process overview の公開入口ファイルは維持し、分割後 overview へのリンクに使う。

## Logic

- process overview は frontend / backend 主題へ分割し、既存 overview は入口図へ縮退する。
- structure overview は package 間の主導線を優先し、内部詳細は detail 図リンクへ逃がす。
- 各 overview は 16:9 以下の比率を acceptance とする。

## Implementation Plan

- active plan を追加する。
- process overview を 3 枚構成へ再編する。
- frontend / backend structure overview を縦方向ベースに再配置する。
- `overview-manifest.json` を新しい overview 構成へ合わせて更新する。
- `d2 validate` と `d2 -t 201` で対象図を検証し、`.svg` を再生成する。

## Acceptance Checks

- 更新した overview `.svg` がすべて 16:9 以下である。
- overview から detail 図への導線が維持されている。
- `d2 validate` が対象 `.d2` すべてで通る。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `d2 validate <target>.d2`
- `d2 -t 201 <target>.d2 <target>.svg`
- 更新後の `.svg` に対する viewBox 比率確認

## 4humans Sync

- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/frontend-processes-overview-robustness.d2`
- `4humans/diagrams/processes/backend-processes-overview-robustness.d2`
- `4humans/diagrams/structures/frontend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/overview-manifest.json`
- 対応する `.svg`

## Outcome

- `processes-overview-robustness.d2` を入口図へ縮退し、`frontend-processes-overview-robustness.d2` と `backend-processes-overview-robustness.d2` を追加した。
- `frontend-structure-overview.d2` と `backend-structure-overview.d2` は detail 図への縦チェーン入口図へ再編した。
- `overview-manifest.json` は process detail 図を frontend/backend overview へ振り分ける構成へ更新した。
- 対応する `.svg` を `d2 -t 201` で再生成し、今回更新した 5 枚の overview はすべて 16:9 以下になった。
