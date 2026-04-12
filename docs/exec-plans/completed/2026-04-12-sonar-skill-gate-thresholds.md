# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: `.codex/skills/` 内の Sonar 通過条件文言の更新
- task_id: sonar-skill-gate-thresholds
- task_catalog_ref: N/A
- parent_phase: skill-modification

## 要求要約

- Sonar の通過条件を `open issue がない` から、明示閾値へ置き換える。
- 対象は `phase-6-implement-frontend`、`phase-6-implement-backend`、`phase-8-review` と関連 permissions に限定する。
- 正本条件は `Security: 0`、`Reliability: 0`、`Maintainability の HIGH/BLOCKER: 0` とする。

## 判断根拠

<!-- Decision Basis -->

- 現在の `open issue がない` は厳しすぎ、Maintainability の低優先度 issue まで closeout / review gate を止める。
- ユーザー意図は Sonar サーバ側 gate 変更ではなく、repo 内 skill 契約の通過条件変更である。
- `.codex/` の role 契約変更なので `skill-modification` の責務内で完結する。

## 対象範囲

- `.codex/skills/phase-6-implement-frontend/SKILL.md`
- `.codex/skills/phase-6-implement-backend/SKILL.md`
- `.codex/skills/phase-8-review/SKILL.md`
- `.codex/skills/phase-6-implement-frontend/references/permissions.json`
- `.codex/skills/phase-6-implement-backend/references/permissions.json`

## 対象外

- Sonar サーバ側 Quality Gate の変更
- `sonar-project.properties` の変更
- product code や tests の変更

## 依存関係・ブロッカー

- `skill-modification` の権限境界
- structure harness の通過

## 並行安全メモ

- 3 skill の Sonar 文言は同一表現へ揃える。
- permissions JSON は frontend/backend で同一ルールへ揃える。

## 機能要件

- `summary`: Sonar の review / closeout 条件を severity と software quality ベースの明示条件へ更新する。
- `in_scope`: skill 文言、permissions 文言、active/completed exec-plan の記録。
- `non_functional_requirements`: 3 skill 間で表現を揃え、曖昧な `quality gate 阻害要因` 文言を残さない。
- `out_of_scope`: Sonar MCP の呼び出し方法変更、workflow 再編、docs 正本更新。
- `open_questions`: N/A
- `required_reading`: `docs/coding-guidelines.md`, `.codex/skills/skill-modification/SKILL.md`

## UI モック

- N/A

## Scenario テスト一覧

- N/A

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: sonar-skill-wording
  - `can_run_in_parallel_with`: none
  - `blocked_by`: structure harness precheck
  - `completion_signal`: 3 skill と 2 permissions の Sonar 条件が一致する
- `tasks`:
  - `task_id`: update-phase6-skills
  - `owned_scope`: phase-6 frontend/backend の Sonar closeout 文言と validation results 文言
  - `depends_on`: structure harness precheck
  - `parallel_group`: sonar-skill-wording
  - `required_reading`: `docs/coding-guidelines.md`
  - `validation_commands`: `rg -n "Security 0|Reliability 0|Maintainability issue" .codex/skills`
  - `task_id`: update-phase8-review
  - `owned_scope`: review checklist の Sonar pass 条件
  - `depends_on`: structure harness precheck
  - `parallel_group`: sonar-skill-wording
  - `required_reading`: `docs/coding-guidelines.md`
  - `validation_commands`: `rg -n "Security 0|Reliability 0|Maintainability issue" .codex/skills/phase-8-review`
  - `task_id`: update-phase6-permissions
  - `owned_scope`: frontend/backend permissions.json の allowed_actions と contract_version
  - `depends_on`: structure harness precheck
  - `parallel_group`: sonar-skill-wording
  - `required_reading`: `.codex/skills/skill-modification/references/permissions.json`
  - `validation_commands`: `rg -n "Security 0|Reliability 0|Maintainability issue" .codex/skills/phase-6-implement-*/references/permissions.json`

## 受け入れ確認

- `open issue がない` が対象 3 skill から消えている。
- Sonar 条件が `Security 0 / Reliability 0 / Maintainability issue(HIGH,BLOCKER) 0` に揃っている。
- frontend/backend permissions の allowed_actions が同じ条件を指している。

## 必要な証跡

<!-- Required Evidence -->

- `python3 scripts/harness/run.py --suite structure`
- `rg` による残存文言確認
- 完了後の `python3 scripts/harness/run.py --suite all`

## 機能要件 HITL 状態

- N/A

## 機能要件 承認記録

- N/A

## 詳細設計 HITL 状態

- N/A

## 詳細設計 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- 完了時は plan を `docs/exec-plans/completed/` へ移し、結果を追記する。

## 結果

<!-- Outcome -->

- `phase-6-implement-frontend` と `phase-6-implement-backend` の Sonar closeout 条件を `Security 0 / Reliability 0 / Maintainability(HIGH,BLOCKER) 0` へ更新した。
- `phase-8-review` の Sonar review 条件を同じ閾値へ揃えた。
- frontend/backend の `references/permissions.json` も同じ条件へ同期した。
- `rg` により、対象 skill から旧 `open issue がない` 文言が除去され、新条件が残っていることを確認した。
- `python3 scripts/harness/run.py --suite structure` は通過した。
- `python3 scripts/harness/run.py --suite all` は SonarCloud 側の `Another SonarQube analysis is already in progress for this project` により失敗した。今回差分による lint / test failure は確認されていない。
