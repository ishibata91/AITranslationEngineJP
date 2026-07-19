# Task Plan: ai-workflow-fork-orchestrator

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- AI 運用ワークフローをオーケストレーター化する。入口を新規実装フローと修正フローの 2 系統に分け、各入口が branch と plan.md を固定する。
- 実装本体を `fork`（親の文脈とモデルを継承する agent）へ委譲する形に変える。
- ワークフロー説明図を、変更差分の塗り分けでなく AS-IS / TO-BE の 2 図で示す形に統一する。

## branch 情報

- `execution_branch`: `claude/ui-data-testid`
- `target_branch`: `master`
- `source_commit`: `9c993d09`

## やらないこと

- storybook-review-loop.md、reviewback.yaml など Storybook レビュー系テンプレの再設計は扱わない。
- 完了済み exec-plans（`docs/exec-plans/completed/`）の過去記録は書き換えない。
- プロダクトコード、プロダクトテストの変更は扱わない（今回は `.claude` 契約と docs テンプレのみ）。
