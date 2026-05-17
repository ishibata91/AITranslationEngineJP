# Task Plan: 2026-05-17-fake-api-retirement

- `workflow`: work
- `status`: completed
- `lane_owner`: `light_change_lane`
- `task_id`: `2026-05-17-fake-api-retirement`
- `task_mode`: cleanup-and-retirement
- `request_summary`: 人間見た目レビュー用 fakeAPI を即削除し、Storybook へ役割を移す前提を固定する。
- `goal`: fakeAPI の runtime、gateway mock、review URL 運用、関連 test、関連 docs 参照を廃止する。
- `constraints`: 並列 worktree 分岐 task として進める。Storybook 基盤、Master Persona 部品化、Storybook POC の実装判断を先取りしない。
- `close_conditions`: fakeAPI 参照が runtime、test、docs、workflow artifact から除去または廃止記録化され、frontend-local が通過または未通過理由を持つ。
- `worktree_path`: `../AITranslationEngineJP-worktrees/2026-05-17-fake-api-retirement`
- `source_branch`: `codex/2026-05-17-fake-api-retirement`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./plan.md`
- `ui_design`: `N/A`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `ux_review`: `N/A`
- `frontend_human_review`: `N/A`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `N/A`
- `implementation_scope`: `N/A`
- `detail_spec_target`: `N/A`

## Decision Record

### 決定

fakeAPI は即削除する。

### 理由

fakeAPI は見た目レビューのために gateway mock を使った仕組みである。Storybook は同じ目的を、UI 部品へ固定 props を渡す形で満たすため、fakeAPI を残すとレビュー正本が二重化する。

### 影響

`frontend/src/controller/review-fake-api/`、fakeAPI 起動 URL、fakeScenario、fakeAPI 関連 test、fakeAPI docs 参照が削除または廃止記録の対象になる。人間見た目レビューの次期正本は Storybook 基盤と Master Persona POC で固定する。

### 未決事項

- `docs/frontend-fake-api.md` は削除する。理由は、廃止済みの起動手順を正本として残すと参照先が二重化するため。
- completed の過去 task artifact に残る fakeAPI 証跡は履歴として残す。理由は、当時の検証証跡とレビュー入力を改変しないため。
- 現行 workflow 文書の fakeAPI 参照は、実画面確認の状態確認へ置き換える。Storybook の入口正本化は後続 task で扱う。

## Routing Notes

- `required_reading`: `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`
- `canonicalization_targets`: `docs/index.md`, fakeAPI を参照する現行 workflow artifact
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `python3 scripts/harness/run.py --suite frontend-local`, `python3 scripts/harness/run.py --suite structure`

## Branch Status

- `worktree_checkout`: `/Users/iorishibata/.codex/worktrees/9da0/AITranslationEngineJP`
- `branch_ready`: `codex/2026-05-17-fake-api-retirement`
- `commit_hash`: `342ad69`; source branch head `e7924693a1d12ad6a6f123f6da473e9645113adc`
- `remote_operation`: `not-performed`

## HITL Status

- `functional_or_design_hitl`: `approved`
- `ux_review`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `2026-05-17 chat: fakeAPI は即削除し、並列 worktree 分岐 task で扱う`

## Codex Implementation Result

- `completed_handoffs`: `light-change direct implementation`
- `touched_files`: `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, `frontend/src/controller/review-fake-api/`, `frontend/src/controller/translation-job-management/translation-job-management-review-gateway.ts`, `frontend/src/ui/review-fake-api-scenario.test.ts`, `docs/frontend-fake-api.md`, `docs/index.md`, `.codex/skills/implement-lane/SKILL.md`, `.codex/skills/implement-frontend/SKILL.md`, `.codex/skills/ux-review/SKILL.md`, `.codex/agents/ux_review.toml`
- `implemented_scope`: fakeAPI runtime、review gateway、fakeAPI 専用 test、現行 docs / workflow 参照を削除または実画面状態確認へ置換した。
- `test_results`: `python3 scripts/harness/run.py --suite frontend-local` pass, `python3 scripts/harness/run.py --suite structure` pass
- `implementation_investigation`: `rg` で `frontend/src`、`.codex/skills`、`.codex/agents`、`docs/index.md` に fakeAPI 現行参照が残らないことを確認した。completed の過去 task artifact は履歴として残した。
- `ui_evidence`: `N/A`
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: `human-reviewed`
- `sonar_gate_result`: `pending`
- `residual_risks`: Storybook 入口正本化は後続 task で扱う。
- `docs_changes`: `docs/frontend-fake-api.md` を削除し、`docs/index.md` から参照を外した。

## Merge Readiness

- `merge_ready`: `ready`
- `source_branch`: `codex/2026-05-17-fake-api-retirement`
- `target_branch`: `master`
- `commit_hash`: `342ad69`; source branch head `e7924693a1d12ad6a6f123f6da473e9645113adc`
- `validation_evidence`: `frontend-local pass`, `structure pass`
- `review_evidence`: `2026-05-18 chat: レビューはした`
- `residual_risks`: completed artifact の fakeAPI 証跡は履歴として残る。

## Merge Result

- `merge_status`: `merged-to-master-local-commit-created`
- `conflict_resolution`: `none`
- `post_merge_validation`: `git diff --check --cached` pass; `python3 scripts/harness/run.py --suite frontend-local` pass; `python3 scripts/harness/run.py --suite structure` pass
- `completed_move`: `docs/exec-plans/active/2026-05-17-fake-api-retirement/` から `docs/exec-plans/completed/2026-05-17-fake-api-retirement/` へ移動済み
- `merge_commit_hash`: `673639d`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `docs/index.md`, `.codex/skills/implement-lane/SKILL.md`, `.codex/skills/implement-frontend/SKILL.md`, `.codex/skills/ux-review/SKILL.md`, `.codex/agents/ux_review.toml`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: `2026-05-18-storybook-foundation` と `2026-05-17-master-persona-componentization` が人間見た目レビュー正本を固定する`

## Outcome

- fakeAPI runtime、review gateway、fakeAPI 専用 test、現行 docs / workflow 参照を削除または廃止後の実画面状態確認へ置換した。
- `frontend-local` と `structure` は通過した。
- 作業 commit は `342ad69` で作成済み。
- マージ準備入力は `Merge Readiness` に記録済み。
- 人間レビューは `2026-05-18 chat: レビューはした` で完了済みである。
- `codex/2026-05-17-fake-api-retirement` は `master` へ conflict なしで local merge 済みである。
- merge 後検証は通過済みである。
- task folder は completed archive へ移動済みである。
- remote repository は変更していない。
