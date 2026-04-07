# 実装計画テンプレート

- workflow: impl
- status: planned
- lane_owner:
- scope:
- task_id:
- task_catalog_ref:
- parent_phase:

## 要求要約

-

## 判断根拠

<!-- Decision Basis -->

-

## 対象範囲

- Prefer repo-root path prefixes or stable scope tokens that match `tasks/phase-*/tasks/*.yaml` when available.

## 対象外

-

## 依存関係・ブロッカー

-

## 並行安全メモ

- Note the shared files, shared fixtures, or upstream `contract` / `verification` tasks that must land first.

## UI

- Use only when the task changes screen structure, presentation, or interaction flow.

## Scenario

- Use only when the task changes user-visible behavior, state transitions, or execution flow.

## Logic

- Use only when the task changes domain logic, contracts, validation, or dependency boundaries.

## 実装計画

<!-- Implementation Plan -->

-

## 受け入れ確認

-

## 必要な証跡

<!-- Required Evidence -->

-

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  クラス追加、依存追加、責務分割変更などで構造が変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  処理追加、相互作用順序変更、主要シナリオ変更などで実行フローが変わる時は、対象 diagram の修正または追加を同じ変更に含め、更新対象ファイルを明記する。
- `4humans/diagrams/overview-manifest.json`
  new detail `.d2` を `4humans/diagrams/structures/` または `4humans/diagrams/processes/` に追加する時は、manifest を同じ変更で更新し、manifest で紐づく overview `.d2` / `.svg` も更新対象へ含める。
- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`4humans` 正本同期後に削除し、completed plan へ持ち越さない。

## 結果

<!-- Outcome -->

-
