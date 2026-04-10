# 実装計画

- workflow: impl
- status: completed
- lane_owner: codex
- scope: .codex/skills/phase-2-design, .codex/skills/orchestrating-implementation/references, docs/exec-plans
- task_id: 2026-04-10-phase-2-design-wireframe-inputs
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- `phase-2-design` が `code.html` と design doc をデザイン指標として参照し、wireframe を作る役割であることを明示する

## 判断根拠

- 現行の `phase-2-design` は `UI` をモック HTML wireframe としか書いておらず、入力ソースと成果物の関係が曖昧である
- orchestrator から渡す `required_reading` に `code.html` と design doc を含める期待を契約に残すと handoff が安定する

## 対象範囲

- `.codex/skills/phase-2-design/SKILL.md`
- `.codex/skills/orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json`
- `docs/exec-plans/active/2026-04-10-phase-2-design-wireframe-inputs.md`

## 対象外

- product code
- `docs/spec.md` の恒久仕様変更
- 他 skill の広範囲な再編

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 既存の未コミット変更は `docs/` 配下にあるため、本タスクは `.codex/skills/phase-2-design` 周辺だけを編集する

## UI

- `phase-2-design` の `UI` section は `docs/screen-design/code.html` と design doc をデザイン指標として読み、そこから wireframe を作る責務だと明記する

## Scenario

- orchestrator が `phase-2-design` へ handoff する時、`required_reading` に `code.html` と design doc を含める期待が読める

## Logic

- skill 本文で role を固定し、handoff 契約で入力期待を補強する

## 実装計画

- `parallel_task_groups`:
  - `group_id`: phase-2-design-skill-docs
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: `phase-2-design` の role と handoff 契約に入力ソースと wireframe 作成責務が明記される
- `tasks`:
  - `task_id`: update-phase-2-design-skill
  - `owned_scope`: `.codex/skills/phase-2-design/SKILL.md`
  - `depends_on`: none
  - `parallel_group`: phase-2-design-skill-docs
  - `required_reading`: `.codex/README.md`, `.codex/skills/phase-2-design/SKILL.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: update-handoff-contract
  - `owned_scope`: `.codex/skills/orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json`
  - `depends_on`: update-phase-2-design-skill
  - `parallel_group`: phase-2-design-skill-docs
  - `required_reading`: `.codex/skills/orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `phase-2-design` に `code.html` と design doc を指標にして wireframe を作る責務が明記されている
- handoff 契約に `required_reading` でその 2 つを含める期待が読める

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure` の通過結果

## HITL 状態

- user が `skill-modification` を明示起動済み

## 承認記録

- user request at 2026-04-10

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/skills/phase-2-design/`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移動する

## 結果

- `.codex/skills/phase-2-design/SKILL.md` に、`UI` は `docs/screen-design/code.html` と `docs/screen-design/design-system-ethereal-archive.md` をデザイン指標として読み、wireframe を作る役割であることを追記した
- `.codex/skills/orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json` に、UI task の `required_reading` は `code.html` と design doc を含める期待を追記した
- `python3 -m json.tool .codex/skills/orchestrating-implementation/references/orchestrating-implementation.to.phase-2-design.json >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
