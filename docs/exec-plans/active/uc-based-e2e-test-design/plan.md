# Task Plan: uc-based-e2e-test-design

- `workflow`: test-design
- `status`: planned
- `lane_owner`: `test_designer`
- `task_id`: `uc-based-e2e-test-design`
- `task_mode`: test-design-only
- `request_summary`: UC 正本を根拠に UI 人間操作 E2E のテスト観点表を作る
- `goal`: 画面 UC ごとに、単独実行できる UI 人間操作 E2E の観点を selector レベルで確認できる状態にする
- `constraints`: プロダクトコード、プロダクトテスト、docs 正本本文は変更しない
- `close_conditions`: `test-design.csv` が固定 header を持ち、関連 UC、対象画面、前提条件、手順、期待値、備考を selector レベルで記録している
- `worktree_path`: `/Users/iorishibata/.codex/worktrees/43aa/AITranslationEngineJP`
- `source_branch`: `N/A`
- `target_branch`: `N/A`

## Artifact Index

- `test_design`: `./test-design.csv`
- `data_testid_gaps`: `./data-testid-gaps.md`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`: `N/A`
- `implementation_scope`: `N/A`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`: `docs/usecases/uc-dashboard.md`, `docs/usecases/uc-provider-settings.md`, `docs/usecases/uc-master-dictionary.md`, `docs/usecases/uc-master-persona.md`, `docs/usecases/uc-translation-management.md`, `docs/usecases/uc-output-management.md`, `docs/e2e-test-guidelines.md`, `docs/coding-guidelines-tests.md`
- `canonicalization_targets`: `N/A`
- `detail_spec_id`: `ai-provider-settings-management`, `master-dictionary`, `translation-input-intake`, `translation-job-management`, `term-translation-phase`, `persona-generation-phase`, `body-translation-phase`, `translation-output-artifact`
- `validation_commands`: `python3 csv parse`, `git diff --check --no-index /dev/null <new-file>`

## HITL Status

- `detail_spec_hitl`: `not-required`
- `storybook_review_loop_input`: `not-required`
- `storybook_review_loop_evidence`: `not-required`
- `frontend_human_review`: `not-required`
- `storybook_design_alignment`: `not-required`
- `approval_record`: `2026-05-28 user request: ucベースのe2eテストを設計して active planを作る`

## Codex Implementation Result

- `completed_handoffs`: `N/A`
- `touched_files`: `docs/exec-plans/active/uc-based-e2e-test-design/plan.md`, `docs/exec-plans/active/uc-based-e2e-test-design/detail-spec-diff.md`, `docs/exec-plans/active/uc-based-e2e-test-design/test-design.csv`, `docs/exec-plans/active/uc-based-e2e-test-design/data-testid-gaps.md`
- `implemented_scope`: UC 正本に基づく UI 人間操作 E2E 観点表の作成。正常 27 件、代替 11 件、例外 8 件、境界 2 件を含む。AI リクエスト系正常系は送信中、操作制限、完了件数、結果行、後続導線を期待値へ含める
- `test_results`: `python3 csv parse passed: 48 rows`; all `前提条件` values start with `画面表示:`; repository、DB、fake、seed は `前提条件` に残っていない; classification count is normal=27, alternative=11, exception=8, boundary=2; AI request normal IDs are `E2E-UC-004`, `E2E-UC-013`, `E2E-UC-045`, `E2E-UC-046`, `E2E-UC-047`; 廃止済みジョブ設定画面の固定名、画面名、テスト ID、不足 ID は active plan のテスト設計成果物に残っていない; `git diff --check -- docs/exec-plans/active/uc-based-e2e-test-design` passed
- `implementation_investigation`: `N/A`
- `ui_evidence`: `N/A`
- `codex_review_result`: `N/A`
- `sonar_gate_result`: `N/A`
- `residual_risks`: 子要素単位の `data-testid` が未固定の操作は `data-testid-gaps.md` に資料化した
- `docs_changes`: `docs/usecases/uc-translation-management.md` に `次のフェーズへ進む` を追加し、active plan 内の task artifact を更新した

## Merge Readiness

- `merge_ready`: `ready`
- `source_branch`: `N/A`
- `target_branch`: `master`
- `commit_hash`: `N/A`
- `validation_evidence`: `python3 csv parse passed: 48 rows`; all `前提条件` values start with `画面表示:`; repository、DB、fake、seed は `前提条件` に残っていない; classification count is normal=27, alternative=11, exception=8, boundary=2; AI request normal IDs are `E2E-UC-004`, `E2E-UC-013`, `E2E-UC-045`, `E2E-UC-046`, `E2E-UC-047`; 廃止済みジョブ設定画面の固定名、画面名、テスト ID、不足 ID は active plan のテスト設計成果物に残っていない; `git diff --check -- docs/exec-plans/active/uc-based-e2e-test-design` passed
- `review_evidence`: `N/A`
- `residual_risks`: 子要素 selector 未固定項目は、E2E 実装前に `data-testid-gaps.md` に沿って画面設計または実装側で固定する必要がある

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `N/A`
- `post_merge_validation`: `pending`
- `completed_move`: `pending`
- `merge_commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `N/A`
- `detail_spec_canonicalization`: `not-required`
- `follow_up`: `test-design.csv` を `implement_lane` のシナリオテスト実装判断へ渡す

## Outcome

- UC 正本と画面設計正本を根拠に、UI 人間操作 E2E の観点表を active plan 内へ作成した。
