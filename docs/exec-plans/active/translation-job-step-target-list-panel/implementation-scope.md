# Implementation Scope: translation-job-step-target-list-panel

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: 2026-05-23 に人間が `approved` として承認
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`:
  - `./screen-design-diff.job-run.md`
  - `./screen-design-diff.term-translation-phase.md`
  - `./screen-design-diff.persona-generation-phase.md`
  - `./screen-design-diff.body-translation-phase.md`
  - `./screen-design-diff.translation-complete.md`
- `component_diagram`: `./design-diff.translation-job-step-target-list-panel.md`

## Fixed Decisions

- `unanswered_questions`: `0`
- UI 変更を含むため、frontend 実装は必須である。
- 実装対象は、`job-run` 共通部品として処理対象一覧表示パネルを追加する範囲に限定する。
- 処理対象一覧表示パネルは、選択ジョブ概要の下、現在段階の画面表示領域の上に配置する。
- 処理対象一覧表示パネルは、現在段階名、処理対象名、処理対象詳細、ページ操作を表示する。
- 処理対象一覧表示パネルは、50 件程度を既定ページサイズとして扱い、現在ページの表示範囲だけを画面要素にする。
- backend、Wails、DTO、gateway、SQLite、docs 正本、`.codex/` は今回の実装範囲に含めない。
- secret は扱わない。

## Target Frontend Boundary

- 対象入口: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- 対象共通部品: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.svelte`
- 対象 props / 表示型: `frontend/src/ui/screens/job-run/job-run-shell-props.ts`
- 対象 fixture: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
- 対象 Storybook: `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- 利用可能な既存部品: `frontend/src/ui/components/PaginationControls.svelte`
- 表示対象段階:
  - 単語翻訳: `共通辞書対象外の用語と固有名詞`
  - NPC ペルソナ生成: `NPC ごとのペルソナ生成入力`
  - 本文翻訳: `辞書置換対象外の翻訳項目`
  - 翻訳結果の確認: `翻訳項目単位の訳文`

## Storybook Confirmation Targets

- レビュー分類: `Review/Changed Screens/Job Run/ProcessingTargetListPanel`
- 通常分類: `Screens/Job Run/ProcessingTargetListPanel`
- 変更部品: `ProcessingTargetListPanel`
- 追加状態:
  - 単語翻訳の処理対象一覧
  - NPC ペルソナ生成の処理対象一覧
  - 本文翻訳の処理対象一覧
  - 翻訳結果の確認の処理対象一覧
  - 1 ページ目
  - 最終ページ
  - 長い処理対象名または長い処理対象詳細
- fixture: `job-run-shell-fixtures.ts` に処理対象一覧の固定データを追加する。

## Forbidden Scope

- `docs/` 正本、`.codex/`、`.codex/skills`、`.codex/agents`、`plan.md` を変更しない。
- backend、Wails bridge、DTO、gateway、repository、SQLite schema を変更しない。
- 翻訳ジョブ実行、開始、一時停止、再開、再試行、取り消し、出力管理への導線の既存動作を変更しない。
- 処理対象一覧のために新しい永続化、API、実データ取得経路を追加しない。
- 未承認のフィルタ、検索、並べ替え、一括操作、空状態、エラー状態を追加しない。
- Storybook review のための fake gateway や長寿命 mock API を追加しない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-processing-target-panel` | `なし` | `なし` | `なし` |
| `wave-2` | `unit-tests-processing-target-panel` | `frontend-processing-target-panel` | `なし` | `depends_on` |
| `wave-3` | `final-frontend-validation-and-review-request` | `frontend-processing-target-panel`, `unit-tests-processing-target-panel` | `なし` | `depends_on` |

## Handoffs

### `frontend-processing-target-panel`

- `implementation_target`: 処理対象一覧表示パネルを `job-run` 共通部品として追加し、現在段階に応じた表示情報を渡す。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.translation-complete.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `ProcessingTargetListPanel.svelte` を追加し、`aria-label="処理対象一覧"` を持つ領域にする。
  - `job-run-shell-props.ts` に表示用処理対象情報とページング props を追加する。
  - `JobRunPage.svelte` で選択ジョブ概要の下、`PhaseHost` の上に共通パネルを配置する。
  - `JobRunPage.svelte` で現在段階に応じた段階名、処理対象名、処理対象詳細を渡す。
  - `ProcessingTargetListPanel` は現在ページの表示範囲だけを描画し、`前へ` と `次へ` でページを切り替える。
  - `job-run-shell-fixtures.ts` に Storybook 用の段階別処理対象一覧 fixture を追加する。
  - `ProcessingTargetListPanel.stories.ts` を `Review/Changed Screens/Job Run/ProcessingTargetListPanel` として追加する。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `size_estimate`:
  - `files`: `5-7`
  - `changed_lines`: `350-650`
  - `classification`: `通常`
- `first_action`: `frontend/src/ui/screens/job-run/job-run-shell-props.ts` に `ProcessingTargetListPanel` 用 props 型を追加し、完了条件の「表示用処理対象情報とページング props を追加する」を閉じる。最初に型を固定すると、部品、親画面、fixture、story の入力境界を同じ語彙で実装できる。
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 選択ジョブがある時、選択ジョブ概要、処理対象一覧、現在段階画面、次の作業フッターの順で表示される。
  - 選択ジョブがない時、処理対象一覧は表示されない。
  - 4 段階の現在段階名、処理対象名、処理対象詳細が承認済み画面設計差分と一致する。
  - ページ操作は 1 ページ目で `前へ` を無効にし、最終ページで `次へ` を無効にする。
  - 現在ページの表示範囲だけが DOM に描画される。
  - Storybook review 分類に、段階別、1 ページ目、最終ページ、長い表示文言の確認状態がある。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - backend 契約を追加しないため、統合境界実装 handoff は作らない。
  - 表示用 fixture は Storybook 確認用であり、実データ取得経路の代替ではない。

### `unit-tests-processing-target-panel`

- `implementation_target`: 処理対象一覧表示パネルの表示とページ操作を frontend unit test で固定する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts` を追加する。
  - `aria-label="処理対象一覧"`、段階名、処理対象名、処理対象詳細の表示を検証する。
  - 50 件程度の現在ページ表示範囲、`前へ`、`次へ`、ページ位置表示を検証する。
  - 現在ページ外の処理対象が DOM に存在しないことを検証する。
- `depends_on`: `frontend-processing-target-panel`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `size_estimate`:
  - `files`: `1-2`
  - `changed_lines`: `120-260`
  - `classification`: `通常`
- `first_action`: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts` に「処理対象一覧の landmark と段階別表示文言を描画する」test を追加し、完了条件の「aria-label、段階名、処理対象名、処理対象詳細の表示を検証する」を閉じる。最初の test が表示契約を固定するため、ページ操作 test を同じ入力で追加できる。
- `validation_commands`:
  - `npm --prefix frontend run test -- ProcessingTargetListPanel`
- `completion_signal`:
  - 処理対象一覧の表示名と詳細が承認済み詳細仕様差分を根拠に検証されている。
  - 1 ページ目、次ページ、最終ページのページ操作が検証されている。
  - 現在ページ外の処理対象を DOM に描画しないことが検証されている。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - UI 人間操作 E2E は最終検証と Storybook 人間レビューで扱う。

### `final-frontend-validation-and-review-request`

- `implementation_target`: frontend 実装と単体テスト完了後に、frontend ローカル検証と Storybook 人間レビュー依頼を揃える。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.translation-complete.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - frontend ローカル検証を実行する。
  - Storybook review URL を `http://localhost:6008/` として、人間レビュー対象を返す。
  - 変更部品、追加状態、story、fixture、関連資源を完了報告に含める。
- `depends_on`: `frontend-processing-target-panel`, `unit-tests-processing-target-panel`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `size_estimate`:
  - `files`: `0-1`
  - `changed_lines`: `0-80`
  - `classification`: `通常`
- `first_action`: `python3 scripts/harness/run.py --suite frontend-local` を実行し、完了条件の「frontend ローカル検証を実行する」を閉じる。実装と単体テストが揃った後に、広域 frontend gate の結果を確認できるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - frontend ローカル検証結果が完了報告に記録されている。
  - Storybook review URL が `http://localhost:6008/` として返されている。
  - Review story は `Review/Changed Screens/Job Run/ProcessingTargetListPanel` で確認できる。
  - fixture は `job-run-shell-fixtures.ts` の処理対象一覧データとして確認できる。
  - 人間レビュー後に通常分類へ戻す対象が `Screens/Job Run/ProcessingTargetListPanel` として明示されている。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `notes`:
  - Storybook は `npm --prefix frontend run storybook` で `http://localhost:6008/` に固定して起動する。
  - Storybook 起動後の人間レビュー判断は Codex implementation lane では確定しない。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `storybook_review_request`:
  - `review_url`: `http://localhost:6008/`
  - `changed_components`: `ProcessingTargetListPanel`
  - `added_states`: 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果の確認、1 ページ目、最終ページ、長い表示文言
  - `story`: `Review/Changed Screens/Job Run/ProcessingTargetListPanel`
  - `fixture`: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
  - `normal_category_after_review`: `Screens/Job Run/ProcessingTargetListPanel`
- `residual_risks`
- `completion_evidence`: レーン終了判断で読む実装事実。completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`
