# Task Plan: ai-workflow-fork-orchestrator

`plan.md` は branch 情報とこの task でやらないことの要点だけを持つ。
設計、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## branch 情報

- `execution_branch`: `claude/ui-data-testid`
- `target_branch`: `master`
- `source_commit`: `9c993d09`

## やらないこと

- storybook-review-loop.md、reviewback.yaml など Storybook レビュー系テンプレの再設計は扱わない。
- 完了済み exec-plans（`docs/exec-plans/completed/`）の過去記録は書き換えない。
- プロダクトコード、プロダクトテストの変更は扱わない（今回は `.claude` 契約と docs テンプレのみ）。
