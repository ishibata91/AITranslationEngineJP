# Active Plans

関連文書: [`../../index.md`](../../index.md), [`../../core-beliefs.md`](../../core-beliefs.md)

このディレクトリには未完了の計画を置きます。
role-based skill へ統合した後も、専門運用の知識は plan ではなく各 skill の `SKILL.md` と `references/` に残します。

## Rules

- 非自明な変更は、実装前にここへ `templates/work-plan.md` ベースの計画を追加する
- task-local UI mock は `docs/exec-plans/active/<task-id>.ui.html` に置く
- task-local Scenario 一覧は `docs/exec-plans/active/<task-id>.scenario.md` に置く
- plan 本文には artifact の path、要点、validation、close 条件だけを残す
- 完了したら `../completed/` へ移動し、結果を追記する
