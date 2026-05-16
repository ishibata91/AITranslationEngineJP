# State Knowledge Investigation Input

- `caller`: `light_change_lane`
- `target_agent`: `investigator`
- `skill`: `investigate`
- `status`: `ready`

## 調査目的

翻訳ジョブ状態関連の stale 廃止について、現在の削減範囲だけで足りるかを調べる。
重複した状態知識、旧設計名、廃止候補の状態値、未使用 package、古い task-local 参照を観測事実として列挙する。

## 既知文脈

- `light-change-planning.md` は `範囲内修正` と判定している。
- `design-diff.md` は `internal/statemachine/`、phase 別 wrapper、phase service action enablement 分岐を主な差分としている。
- `backend-implementation-result.md` は backend 実装完了を記録している。
- 人間確認へ進める前に、調査範囲が狭すぎないかを確認する必要がある。

## 調査対象

- `internal/`
- `.go-arch-lint.yml`
- `docs/architecture.md`
- `docs/spec.md`
- `docs/detail-specs/`
- `docs/diagrams/backend/`
- `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/`
- `docs/exec-plans/active/observability-log-addition/`

## 観測してほしい観点

- 同じ job state / phase state の知識が複数箇所に重複していないか。
- 廃止済みまたは旧設計名の `StateMachine`、`JobIOService`、`statemachine`、`jobio` が残っていないか。
- `pending`、`ready`、`running`、`paused`、`recoverable_failed`、`completed`、`failed`、`canceled` の扱いが、仕様、policy、service、read model、docs で食い違っていないか。
- domain の `stale_selection`、`validation_stale`、`model_selection_stale` と、今回廃止したい stale の意味が混ざっていないか。
- 今回の実装差分によって新しく重複または責務違反が生まれていないか。

## 禁止事項

- プロダクトコードを変更しない。
- プロダクトテストを変更しない。
- docs 正本本文を変更しない。
- 実装方針と変更ファイルを確定しない。
- 恒久修正を行わない。
- 下位 agent を起動しない。

## 期待する成果物

- `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/state-knowledge-investigation.md`

## 返却してほしい内容

- 判断結果
- 調査 mode
- 観測点
- 観測事実
- 仮説
- 影響ファイル候補
- 残り不足
- 残留リスク
- 推奨 next step
