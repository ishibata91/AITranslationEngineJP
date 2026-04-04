# Fix Plan Template

- workflow: fix
- status: planned
- lane_owner:
- scope:

## Request Summary

-

## Decision Basis

-

## Known Facts

-

## Trace Plan

-

## Fix Plan

-

## Acceptance Checks

-

## Required Evidence

-

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  クラス追加、依存追加、責務分割変更などで構造が変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  処理追加、相互作用順序変更、主要シナリオ変更などで実行フローが変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。
- `4humans/diagrams/overview-manifest.json`
  new detail `.d2` を `4humans/diagrams/structures/` または `4humans/diagrams/processes/` に追加する時は、manifest を同じ変更で更新し、manifest で紐づく overview `.d2` / `.svg` も更新対象へ含める。

## Outcome

-
