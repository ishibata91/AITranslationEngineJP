# 実装計画

- workflow: impl
- status: completed
- lane_owner: codex
- scope: docs/screen-design, docs/index.md, docs/exec-plans
- task_id: 2026-04-10-ethereal-archive-design-doc-ja
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- 提供された `Design.md` を `docs/` 配下の正本文書として日本語化する
- `screen-design` 配下に配置し、既存 wireframe と矛盾する `App Shell` 表現を最小範囲でそろえる

## 判断根拠

- `docs/index.md` は `screen-design/` を画面設計の入口として扱っている
- 原文には `Top Navigation Only` が含まれ、既存の `docs/screen-design/wireframes/app-shell.md` の左 navigation と衝突する
- source of truth の発見性を保つため、配置先の案内文書も同時に更新する

## 対象範囲

- `docs/screen-design/README.md`
- `docs/screen-design/design-system-ethereal-archive.md`
- `docs/screen-design/wireframes/app-shell.md`
- `docs/index.md`
- `docs/exec-plans/active/2026-04-10-ethereal-archive-design-doc-ja.md`

## 対象外

- product code
- `.codex/`
- `docs/spec.md` の要件変更

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 共有ファイルは `docs/index.md` と `docs/screen-design/wireframes/app-shell.md`
- 本タスクでは `docs/` 以外へ触れない

## UI

- `The Ethereal Archive` の visual design system を日本語で正本化する
- `App Shell` は top navigation 前提へ更新する

## Scenario

- `screen-design/` から visual design と wireframe の双方へ到達できるようにする

## Logic

- なし

## 実装計画

- `parallel_task_groups`:
  - `group_id`: docs-authoring
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: 新規 design doc、案内文、関連 wireframe がそろう
- `tasks`:
  - `task_id`: author-design-doc
  - `owned_scope`: `docs/screen-design/design-system-ethereal-archive.md`
  - `depends_on`: none
  - `parallel_group`: docs-authoring
  - `required_reading`: `docs/index.md`, `docs/spec.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: align-screen-design-entry
  - `owned_scope`: `docs/screen-design/README.md`, `docs/index.md`, `docs/screen-design/wireframes/app-shell.md`
  - `depends_on`: author-design-doc
  - `parallel_group`: docs-authoring
  - `required_reading`: `docs/screen-design/wireframes/README.md`, `docs/screen-design/wireframes/app-shell.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite all`

## 受け入れ確認

- 提供文書の内容が日本語で `docs/` に保存されている
- `screen-design/` から新規文書へ到達できる
- `App Shell` wireframe が top navigation 方針と矛盾しない

## 必要な証跡

- structure harness 実行結果
- full harness 実行結果

## HITL 状態

- human が `docs/` 更新を明示済み

## 承認記録

- user request at 2026-04-10

## review 用差分図

- N/A

## 差分正本適用先

- `docs/screen-design/`

## Closeout Notes

- 追加した plan は完了後に `docs/exec-plans/completed/` へ移動する

## 結果

- `docs/screen-design/design-system-ethereal-archive.md` を追加し、提供された design spec を日本語化した
- `docs/screen-design/README.md` を追加し、visual design と wireframe の入口を整理した
- `docs/screen-design/wireframes/app-shell.md` を top navigation 前提へ更新した
- `docs/index.md` の `screen-design/` 導線を `README.md` 参照へ更新した
- `python3 scripts/harness/run.py --suite structure` と `python3 scripts/harness/run.py --suite all` が通過した
