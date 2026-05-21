# Task Plan: 2026-05-22-job-setup-model-selection-relocation

- `workflow`: work
- `status`: waiting-human-design-review
- `lane_owner`: implement_lane
- `task_id`: `2026-05-22-job-setup-model-selection-relocation`
- `task_mode`: new implementation
- `request_summary`: ジョブセットアップ画面を廃止し、AI モデル選択を翻訳段階ごとの画面へ移動する。
- `goal`: 利用者が入力データ確認後に単語翻訳へ進み、各翻訳段階の開始前にその段階の AI モデルを選べる状態にする。
- `constraints`: 人間設計レビュー前にプロダクトコードを変更しない。AI サービス設定画面は残す。秘密値本体を UI、DTO、read model、log、error summary に出さない。Storybook を frontend 人間レビューの主確認面にする。
- `close_conditions`: 詳細仕様差分、画面設計差分、設計差分図が人間レビューで承認される。承認後に implementation-scope を作る。実装後に Storybook 人間レビュー依頼、frontend 人間レビュー、観測ログ判断、最終検証、ブラウザ確認、正本化判断、local commit、マージ準備入力を揃える。
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-22-job-setup-model-selection-relocation`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./task-frame.md`
- `screen_design_diff`: `./screen-design-diff.translation-management.md`, `./screen-design-diff.translation-input-review.md`, `./screen-design-diff.translation-job-setup.md`, `./screen-design-diff.term-translation-phase.md`, `./screen-design-diff.persona-generation-phase.md`, `./screen-design-diff.body-translation-phase.md`
- `design_diff_diagram`: `./design-diff.job-setup-model-selection-relocation.md`
- `storybook_review_request`: `required-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approved_frontend_protection`: `pending-after-frontend-human-review`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`

## Routing Notes

- `required_reading`:
  - `docs/index.md`
  - `docs/detail-specs/translation-job-setup.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `docs/screen-design/screens/translation-management.md`
  - `docs/screen-design/screens/translation-input-review.md`
  - `docs/screen-design/screens/translation-job-setup.md`
  - `docs/screen-design/screens/term-translation-phase.md`
  - `docs/screen-design/screens/persona-generation-phase.md`
  - `docs/screen-design/screens/body-translation-phase.md`
  - `frontend/src/ui/stores/shell-state.ts`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/components/AIModelSelectionCard.svelte`
- `canonicalization_targets`:
  - human 承認後に `docs/detail-specs/` と `docs/screen-design/screens/` の正本反映を判断する。
- `detail_spec_id`: `translation-job-setup`, `term-translation-phase`, `persona-generation-phase`, `body-translation-phase`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite structure`
  - `npm --prefix frontend run build-storybook`
  - `npm run dev:wails:agent-browser`

## Branch Status

- `worktree_checkout`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `branch_ready`: `created: codex/2026-05-22-job-setup-model-selection-relocation`
- `commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Design Decisions

### D-01 ジョブセットアップ画面を削除する

決定: 翻訳管理の下位画面からジョブセットアップ画面を削除する。

理由: 画面名はセットアップだが、現在の主操作は翻訳段階ごとの AI モデル選択に寄っている。

影響: 入力データ確認の次の導線は、翻訳設定画面ではなく単語翻訳画面へ進む。

### D-02 AI モデル選択を段階画面へ移す

決定: 単語翻訳、NPC ペルソナ生成、本文翻訳の各画面に AI モデル選択を置く。

理由: 各段階で使う AI モデルは、段階開始前に判断する情報である。

影響: 3 段階すべての AI 設定をジョブ作成前にまとめて固定しない。

### D-03 ジョブ作成条件と段階開始条件を分ける

決定: 翻訳ジョブ作成は入力データ確認に閉じ、AI モデル選択は各段階の開始条件にする。

理由: ジョブ作成と AI 実行条件を同じ画面で扱うと、開始前の段階まで先回りして選ばせる。

影響: backend では、ジョブ作成、段階 AI 設定保存、段階開始を別の責務として扱う必要がある。

### D-04 翻訳段階画面の常時表示要素を減らす

決定: 単語翻訳、NPC ペルソナ生成、本文翻訳の画面は、開始判断、実行状態、結果判断に必要な情報へ整理する。

理由: 現在の段階画面は、診断情報、入力 snapshot、後続 readiness、失敗情報を常時並べており、利用者が次の操作を判断しにくい。

影響: 詳細な診断情報は、常時表示ではなく、失敗時、完了時、または詳細確認が必要な状態の表示へ寄せる。

## HITL Status

- `detail_spec_hitl`: `required-after-design-bundle`
- `storybook_review_request`: `required-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approval_record`: `pending-human-design-review`

## Codex Implementation Result

- `completed_handoffs`: 未着手
- `touched_files`: 未着手
- `implemented_scope`: 未着手
- `test_results`: 未着手
- `implementation_investigation`: 未着手
- `ui_evidence`: 未着手
- `storybook_review_request`: 未着手
- `storybook_review_resources`: 未着手
- `approved_frontend_protection`: 未着手
- `codex_review_result`: 未着手
- `sonar_gate_result`: 未着手
- `residual_risks`: 実装前である。既存 backend はジョブ作成時に phase runtime snapshot を保存するため、backend と統合境界の変更範囲が大きくなる可能性がある。
- `docs_changes`: task 内設計成果物のみ

## Merge Readiness

- `merge_ready`: `blocked`
- `source_branch`: `codex/2026-05-22-job-setup-model-selection-relocation`
- `target_branch`: `master`
- `commit_hash`: `N/A`
- `validation_evidence`: 設計成果物のみ。プロダクト検証は未実行。
- `review_evidence`: `pending-human-design-review`
- `residual_risks`: 人間設計レビュー未承認のため実装へ進めない。

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `N/A`
- `post_merge_validation`: `N/A`
- `completed_move`: `N/A`
- `merge_commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: 未着手
- `detail_spec_canonicalization`: human 承認後に判断する。
- `follow_up`: human 承認後、implementation-scope を frontend、backend、統合境界、シナリオテスト、単体テストに分ける。

## Outcome

- 設計成果物を作成した。
- 人間設計レビュー待ちで停止する。
