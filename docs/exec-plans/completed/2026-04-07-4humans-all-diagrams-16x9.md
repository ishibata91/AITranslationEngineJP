# 4humans Structure Diagrams 16:9 Rewrite

- workflow: impl
- status: planned
- lane_owner:
- scope: `docs/exec-plans/active/2026-04-07-4humans-all-diagrams-16x9.md`, `4humans/diagrams/structures/*.d2`, `4humans/diagrams/structures/*.svg`, `4humans/diagrams/overview-manifest.json`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `4humans/diagrams/structures/` 配下の全図に 16:9 制約を適用し、情報量を減らさずに読み直せる D2 / SVG へ再設計する。

## Decision Basis

- structure overview と detail class 図の多くが 16:9 を超えており、review 時に横スクロールが必須になっている。
- frontend slice class 図と backend class 図では必要な再配置密度が異なる。
- information loss を避けるため、単純なラベル削減ではなく、縦積み・section 分割・layer 分割で対応する。

## Owned Scope

- `4humans/diagrams/structures/*.d2` と対応 `.svg`
- `4humans/diagrams/overview-manifest.json`

## Out Of Scope

- `docs/` 正本の更新
- 実装コードの変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること
- backend 大型 class 図で section 分割とリンク整合を保つこと

## Parallel Safety Notes

- `.d2` を source of truth として更新する。
- 既存パスと manifest 整合を維持する。

## Logic

- overview は縦 stack を維持しつつ、layer/package 情報を持たせる。
- frontend class 図は `dagre` と縦 layer 配置で 16:9 に寄せる。
- backend class 図は package の責務単位を縦 section 化し、重い図は source 内で section 分割する。

## Implementation Plan

- 横長図の分類結果を固定する。
- overview 図を情報量維持のまま 16:9 に収める。
- frontend class 図を縦 layer 配置へ置換する。
- backend class 図を section 化して再配置する。
- 全 `.d2` を validate し、`.svg` を再生成する。

## Acceptance Checks

- `4humans/diagrams/structures/*.svg` の全件で `width / height <= 16/9` を満たす。
- 各図で既存の責務、参加者、主要 edge / message が欠落していない。
- `d2 validate` が全対象で通る。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`
- 比率確認スクリプトの結果
- `d2 validate <target>.d2`
- `d2 -t 201 <target>.d2 <target>.svg`

## 4humans Sync

- `4humans/diagrams/structures/*.d2`
- `4humans/diagrams/structures/*.svg`
- `4humans/diagrams/overview-manifest.json`

## Outcome

- pending
