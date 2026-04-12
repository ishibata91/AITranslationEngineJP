# 実装計画

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/workflow.md
- task_id: 2026-04-12-fix-lane-implementation-reuse
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- `orchestrating-fixes` が後段で implementation lane の skill を基本再利用できるようにする
- 特に実装、UI確認、レビューを `phase-6-*`、`phase-6.5-ui-check`、`phase-8-review` へ寄せる
- `reporting-risks` は fix lane の live workflow から外す

## 判断根拠

- 現行の fix lane は `implementing-fixes`、`reviewing-fixes`、`reporting-risks` に分かれ、implementation lane と closeout の観点が二重化している
- `phase-6.5-ui-check` と `phase-8-review` は implementation lane 専用 contract のため、fix lane から再利用するには skill と contract の両方を同期する必要がある
- live workflow は `.codex/skills/` と `.codex/workflow.md` を同時にそろえないと崩れる

## 対象範囲

- `.codex/skills/orchestrating-fixes/`
- `.codex/skills/phase-6-implement-backend/`
- `.codex/skills/phase-6-implement-frontend/`
- `.codex/skills/phase-6.5-ui-check/`
- `.codex/skills/phase-8-review/`
- `.codex/workflow.md`

## 対象外

- product code
- `docs/` の恒久仕様変更
- implementation lane の前段詳細設計フロー

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 変更対象は `.codex/` と exec-plan に限定する
- `implementing-fixes` と `reviewing-fixes` は削除せず、まずは fix lane の live 参照だけを切り替える

## 実装計画

- `parallel_task_groups`:
  - `group_id`: fix-lane-implementation-reuse
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: fix lane の skill、contract、workflow が同じ後段 skill 前提でそろう
- `tasks`:
  - `task_id`: repoint-fix-orchestrator
  - `owned_scope`: `.codex/skills/orchestrating-fixes/`
  - `depends_on`: none
  - `parallel_group`: fix-lane-implementation-reuse
  - `required_reading`: `.codex/workflow.md`, `.codex/skills/orchestrating-fixes/SKILL.md`, `.codex/skills/orchestrating-implementation/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: widen-shared-phase-contracts
  - `owned_scope`: `.codex/skills/phase-6-implement-backend/`, `.codex/skills/phase-6-implement-frontend/`, `.codex/skills/phase-6.5-ui-check/`, `.codex/skills/phase-8-review/`
  - `depends_on`: repoint-fix-orchestrator
  - `parallel_group`: fix-lane-implementation-reuse
  - `required_reading`: `.codex/skills/phase-6-implement-backend/SKILL.md`, `.codex/skills/phase-6.5-ui-check/SKILL.md`, `.codex/skills/phase-8-review/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: sync-fix-workflow-doc
  - `owned_scope`: `.codex/workflow.md`
  - `depends_on`: widen-shared-phase-contracts
  - `parallel_group`: fix-lane-implementation-reuse
  - `required_reading`: `.codex/workflow.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `orchestrating-fixes` が `phase-6-implement-*`、`phase-6.5-ui-check`、`phase-8-review` を live handoff 先として明記している
- 共用 skill 側が fix lane 用 contract を参照できる
- `reporting-risks` が fix lane の live workflow から外れている
- `.codex/workflow.md` が新しい fix lane 後段を説明している

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`

## HITL 状態

- user が `skill-modification` を明示起動済み

## 承認記録

- user request at 2026-04-12

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/skills/`
- `.codex/workflow.md`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移動する

## 結果

- `.codex/skills/orchestrating-fixes/` を更新し、fix lane の live handoff 先を `phase-6-implement-backend`、`phase-6-implement-frontend`、`phase-6.5-ui-check`、`phase-8-review` へ切り替えた
- `.codex/skills/phase-6-implement-backend/`、`.codex/skills/phase-6-implement-frontend/`、`.codex/skills/phase-6.5-ui-check/`、`.codex/skills/phase-8-review/` を更新し、implementation lane と fix lane の両方の contract を受けられるようにした
- fix lane 用の双方向 contract JSON を追加し、shared phase skill との handoff を live workflow に合わせた
- `.codex/workflow.md` を更新し、修正レーンの後段を implementation phase / UI check phase / implementation review phase 再利用として明記した
- `python3 scripts/harness/run.py --suite structure` が通過した
