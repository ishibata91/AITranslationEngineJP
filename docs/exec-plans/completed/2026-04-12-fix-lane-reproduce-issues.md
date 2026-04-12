# 実装計画

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/workflow.md
- task_id: 2026-04-12-fix-lane-reproduce-issues
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- `analyzing-fixes` を live workflow から外す
- `logging-fixes` の直後に `reproduce-issues` を追加する
- `reproduce-issues` は Playwright MCP の console 確認と Wails ログの直接確認を担当する
- `reproduce-issues` の agent は `ui_checker` を使う
- 回帰防止の test 実装は修正後、`phase-8-review` の前に置く

## 判断根拠

- 現行の `analyzing-fixes` は観測結果圧縮に寄っているが、今回欲しいのは console と Wails log を実際に確認する再現エージェントである
- `logging-fixes` の直後に再現確認を置くと、temporary logging で埋めた観測点をそのまま確認できる
- fix lane では修正後に UI で治ったことを確認してから回帰防止テストを足す方が、review 前の証明として自然である
- live workflow は skill、contract、workflow doc を同時に更新しないと崩れる

## 対象範囲

- `.codex/skills/orchestrating-fixes/`
- `.codex/skills/logging-fixes/`
- `.codex/skills/reproduce-issues/`
- `.codex/skills/analyzing-fixes/`
- `.codex/skills/phase-5-test-implementation/`
- `.codex/workflow.md`

## 対象外

- product code
- implementation lane
- `docs/` の恒久仕様変更

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 変更対象は `.codex/` と exec-plan に限定する
- `analyzing-fixes` は削除ではなく retired 表記に留める

## 実装計画

- `parallel_task_groups`:
  - `group_id`: fix-lane-reproduce-issues
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: fix lane の観測後段が `reproduce-issues` と修正後 test 実装前提でそろう
- `tasks`:
  - `task_id`: add-reproduce-issues-skill
  - `owned_scope`: `.codex/skills/reproduce-issues/`
  - `depends_on`: none
  - `parallel_group`: fix-lane-reproduce-issues
  - `required_reading`: `.codex/skills/logging-fixes/SKILL.md`, `.codex/skills/phase-6.5-ui-check/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: repoint-fix-orchestrator
  - `owned_scope`: `.codex/skills/orchestrating-fixes/`, `.codex/skills/logging-fixes/`, `.codex/skills/analyzing-fixes/`, `.codex/skills/phase-5-test-implementation/`
  - `depends_on`: add-reproduce-issues-skill
  - `parallel_group`: fix-lane-reproduce-issues
  - `required_reading`: `.codex/skills/orchestrating-fixes/SKILL.md`, `.codex/workflow.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: sync-fix-workflow-doc
  - `owned_scope`: `.codex/workflow.md`
  - `depends_on`: repoint-fix-orchestrator
  - `parallel_group`: fix-lane-reproduce-issues
  - `required_reading`: `.codex/workflow.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `orchestrating-fixes` が `analyzing-fixes` ではなく `reproduce-issues` を参照する
- `reproduce-issues` に `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` がある
- `reproduce-issues` の agent が `ui_checker` になっている
- fix lane が `phase-6 -> phase-6.5 -> phase-5 -> phase-8` の順で記述されている
- `.codex/workflow.md` が `logging-fixes` 後の再現確認と、修正後の回帰防止 test 実装を説明している

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`

## HITL 状態

- user が `skill-modification` 相当の変更を依頼済み

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

- `.codex/skills/reproduce-issues/` を追加し、`ui_checker` を agent に割り当てた
- `.codex/skills/orchestrating-fixes/` を更新し、`logging-fixes` の後に `reproduce-issues` を置き、fix lane の後段順序を `phase-6 -> phase-6.5 -> phase-5 -> phase-8` へ変更した
- `.codex/skills/phase-5-test-implementation/` の fix lane 説明と contract を更新し、修正後・review 前の回帰防止 test 実装として扱うようにした
- `.codex/skills/logging-fixes/` の返却 contract を更新し、`reproduce-issues` へ渡せる観測情報にそろえた
- `.codex/skills/analyzing-fixes/` を retired 表記へ更新し、live workflow から外した
- `.codex/workflow.md` を更新し、修正レーンの順序と差し戻し条件を新しい flow に合わせた
- `python3 scripts/harness/run.py --suite structure` が通過した
