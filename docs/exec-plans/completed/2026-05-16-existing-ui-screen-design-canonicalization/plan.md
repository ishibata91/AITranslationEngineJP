# Task Plan: 2026-05-16-existing-ui-screen-design-canonicalization

- `workflow`: docs-design-canonicalization
- `status`: docs-canonicalized
- `lane_owner`: Codex
- `task_id`: `2026-05-16-existing-ui-screen-design-canonicalization`
- `task_mode`: 既存 UI 画面設計書の作成と docs 正本化
- `request_summary`: 既存の全 UI 画面について、画面単位で screen-design 差分を作り、承認済み docs-only 成果物として画面設計書正本へ反映する。
- `goal`: 実装済み UI 画面の目的、ワイヤーフレーム、表示要素、操作、状態条件を `docs/screen-design/screens/` に画面別正本として揃える。
- `constraints`: プロダクトコード、プロダクトテスト、作業流れ、skill、agent 実行定義は変更しない。
- `constraints`: 画面設計書には、実装指示、テスト手順、agent handoff を書かない。
- `constraints`: 画面設計差分は `docs/screen-design/screens/template.md` の項目に合わせる。
- `close_conditions`: 画面ごとの `screen-design-diff.<screen-id>.md` が active plan 内に揃っている。
- `close_conditions`: `docs/screen-design/screens/<screen-id>.md` が差分成果物と対応している。
- `close_conditions`: docs 構造確認と markdown 形式確認が完了している。
- `worktree_path`: `/Users/iorishibata/.codex/worktrees/207a/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-16-existing-ui-screen-design-canonicalization`
- `target_branch`: `master`

## Artifact Index

- `ui_design`: `N/A`
- `screen_design_diff`: `./screen-design-diff.<screen-id>.md`
- `docs_canonicalization`: `./docs-canonicalization-result.md`
- `screen_design_target`: `docs/screen-design/screens/<screen-id>.md`
- `scenario_design`: `N/A`
- `implementation_scope`: `N/A`
- `detail_spec_target`: `N/A`

## Scope

### 主要ページ

- `dashboard`: ダッシュボード。主要ページへの入口。
- `provider-settings`: AI サービス設定。
- `master-dictionary`: マスター辞書。
- `master-persona`: マスターペルソナ。
- `translation-management`: 翻訳管理シェル。下位画面への入口と現在地表示を扱う。
- `output-management`: 出力管理。

### 翻訳管理の下位画面

- `translation-input-review`: 入力データの確認。
- `translation-job-setup`: 翻訳設定。
- `translation-job-management`: 未完了ジョブ一覧。
- `job-run`: 翻訳実行シェル。段階別画面とフッター遷移を扱う。
- `term-translation-phase`: 単語翻訳。
- `persona-generation-phase`: NPC ペルソナ生成。
- `body-translation-phase`: 本文翻訳。
- `translation-complete`: 翻訳結果の確認。

## Work Plan

1. Codex が既存実装、`docs/spec.md`、`docs/architecture.md`、`docs/UX-standard.md`、`docs/screen-design/` を確認する。
2. Codex が画面一覧、書き込み範囲、差分成果物名を固定する。
3. `designer` agent を画面ごとに起動し、各 agent は担当画面の `screen-design-diff.<screen-id>.md` だけを作る。
4. Codex が designer 成果物を確認し、衝突、重複、template 逸脱を直す。
5. `docs_updater` agent が承認済み docs-only 成果物として、差分を `docs/screen-design/screens/` へ反映する。
6. Codex が構造確認と差分確認を行い、結果を `docs-canonicalization-result.md` に残す。

## Designer Split

- `designer-dashboard`: `screen-design-diff.dashboard.md`
- `designer-provider-settings`: `screen-design-diff.provider-settings.md`
- `designer-master-dictionary`: `screen-design-diff.master-dictionary.md`
- `designer-master-persona`: `screen-design-diff.master-persona.md`
- `designer-translation-management`: `screen-design-diff.translation-management.md`
- `designer-translation-input-review`: `screen-design-diff.translation-input-review.md`
- `designer-translation-job-setup`: `screen-design-diff.translation-job-setup.md`
- `designer-translation-job-management`: `screen-design-diff.translation-job-management.md`
- `designer-job-run`: `screen-design-diff.job-run.md`
- `designer-term-translation-phase`: `screen-design-diff.term-translation-phase.md`
- `designer-persona-generation-phase`: `screen-design-diff.persona-generation-phase.md`
- `designer-body-translation-phase`: `screen-design-diff.body-translation-phase.md`
- `designer-translation-complete`: `screen-design-diff.translation-complete.md`
- `designer-output-management`: `screen-design-diff.output-management.md`

## Approval Record

- `docs-only approval`: この依頼文を、既存 UI 画面設計書の作成と `docs/screen-design/screens/` への反映承認として扱う。
- `implementation approval`: なし。プロダクト実装は行わない。
- `human review`: 完成後に人間が画面設計書を確認する。

## Routing Notes

- `required_reading`: `.codex/skills/design-bundle/SKILL.md`
- `required_reading`: `.codex/skills/ui-design/SKILL.md`
- `required_reading`: `.codex/skills/updating-docs/SKILL.md`
- `required_reading`: `docs/screen-design/README.md`
- `required_reading`: `docs/screen-design/screens/README.md`
- `required_reading`: `docs/screen-design/screens/template.md`
- `canonicalization_targets`: `docs/screen-design/screens/`
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`
- `validation_commands`: `git diff --check -- docs/exec-plans/active/2026-05-16-existing-ui-screen-design-canonicalization docs/screen-design`

## Stop Conditions

- 既存実装だけでは画面目的または主要操作を判断できない画面がある場合は、対象画面を `docs-canonicalization-result.md` の未決事項へ残す。
- 画面が実装上 placeholder だけの場合は、placeholder 画面として設計書を作り、恒久画面仕様を独断で補完しない。
- designer 成果物が docs 正本へ反映できない形式の場合は、Codex が template へ合わせてから正本化する。

## Outcome

- 既存 UI 14 画面の `screen-design-diff.<screen-id>.md` を作成した。
- `docs/screen-design/screens/<screen-id>.md` へ全差分を正本化した。
- `docs/screen-design/screens/README.md` の Records に画面別正本を追加した。
- `python3 scripts/harness/run.py --suite structure` は通過した。
- `git diff --check -- docs/exec-plans/active/2026-05-16-existing-ui-screen-design-canonicalization docs/screen-design` は通過した。
