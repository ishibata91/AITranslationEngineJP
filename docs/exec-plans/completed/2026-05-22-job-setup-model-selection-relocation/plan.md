# Task Plan: 2026-05-22-job-setup-model-selection-relocation

- `workflow`: work
- `status`: completed
- `lane_owner`: implement_lane
- `task_id`: `2026-05-22-job-setup-model-selection-relocation`
- `task_mode`: new implementation
- `request_summary`: ジョブセットアップ画面を廃止し、AI モデル選択を翻訳段階ごとの画面へ移動する。
- `goal`: 利用者が入力データ確認後に単語翻訳へ進み、各翻訳段階の開始前にその段階の AI モデルを選べる状態にする。
- `constraints`: 人間設計レビュー前にプロダクトコードを変更しない。AI サービス設定画面は残す。秘密値本体を UI、DTO、read model、log、error summary に出さない。Storybook を frontend 人間レビューの主確認面にする。
- `close_conditions`: 詳細仕様差分、画面設計差分、設計差分図は人間レビューで承認済み。implementation-scope は作成済み。実装後に Storybook 人間レビュー依頼、frontend 人間レビュー、観測ログ判断、最終検証、ブラウザ確認、正本化判断、local commit、マージ準備入力を揃える。
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-22-job-setup-model-selection-relocation`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `./task-frame.md`
- `screen_design_diff`: `./screen-design-diff.translation-management.md`, `./screen-design-diff.translation-input-review.md`, `./screen-design-diff.translation-job-setup.md`, `./screen-design-diff.term-translation-phase.md`, `./screen-design-diff.persona-generation-phase.md`, `./screen-design-diff.body-translation-phase.md`
- `design_diff_diagram`: `./design-diff.job-setup-model-selection-relocation.md`
- `storybook_review_request`: `ready-for-human-review`
- `frontend_human_review`: `approved`
- `approved_frontend_protection`: `ready`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `implementation_scope`: `./implementation-scope.md`
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
- `commit_hash`: `local commit created; see git log -1`
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

- `detail_spec_hitl`: `approved`
- `storybook_review_request`: `completed`
- `frontend_human_review`: `approved`
- `approval_record`: `approved-by-human: 2026-05-22 approve`

## Human Design Review

- `reviewed_artifacts`: `./detail-spec-diff.md`, `./screen-design-diff.translation-management.md`, `./screen-design-diff.translation-input-review.md`, `./screen-design-diff.translation-job-setup.md`, `./screen-design-diff.term-translation-phase.md`, `./screen-design-diff.persona-generation-phase.md`, `./screen-design-diff.body-translation-phase.md`, `./design-diff.job-setup-model-selection-relocation.md`
- `result`: approved
- `reviewer`: human
- `recorded_at`: 2026-05-22
- `comment`: approve
- `next_artifact`: `./implementation-scope.md`

## Codex Implementation Result

- `completed_handoffs`: `implementation-scope`, `frontend-job-setup-relocation-and-storybook`, `backend-job-creation-decoupling`, `integration-phase-ai-settings-save`, `unit-test-input-review-footer`, `scenario-test-job-setup-relocation`, `observability-phase-ai-settings`
- `touched_files`: `frontend/src/ui/stores/shell-state.ts`, `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/App.svelte`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-input/InputReviewPage.test.ts`, `frontend/src/ui/components/AIModelSelectionCard.svelte`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`, `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/FieldResultListPanel.svelte`, `frontend/package.json`, Storybook story / fixture files
- `implemented_scope`: `frontend-job-setup-relocation-and-storybook`, `backend-job-creation-decoupling`, `integration-phase-ai-settings-save`, `unit-and-scenario-tests`, `observability-phase-ai-settings`
- `test_results`: `frontend-local` pass, `backend-local` pass, `structure` pass, `build-storybook` pass, `coverage` failed by existing Sonar issue outside this handoff
- `implementation_investigation`: 未着手
- `ui_evidence`: `/private/tmp/job-setup-relocation-browser-confirmation/snapshot.json`
- `storybook_review_request`: ready
- `storybook_review_resources`: ready
- `approved_frontend_protection`: ready
- `codex_review_result`: 未着手
- `sonar_gate_result`: 未着手
- `residual_risks`: Storybook 起動時に既存の `EMFILE` watcher 警告と settings 書き込み `EPERM` が出る。Wails dev の実データ導線は既存の `ImportTranslationInput` binding 未接続により、入力登録後の実ジョブ作成までは確認していない。
- `docs_changes`: 画面設計正本へ frontend 人間レビュー結果を反映済み。

## Storybook Human Review Request

- `review_status`: approved-and-returned-to-normal-classification
- `storybook_url`: `http://127.0.0.1:6008/`
- `storybook_command`: `npm --prefix frontend run storybook -- --host 127.0.0.1 --port 6008`
- `storybook_start_note`: `6006` と `6007` は使用中だったため `6008` で起動した。起動時に `EMFILE` と `/Users/iorishibata/.storybook/settings.json` の `EPERM` が出たが、Storybook は `http://127.0.0.1:6008/` で ready になった。
- `storybook_build_result`: `npm --prefix frontend run build-storybook` pass
- `frontend_validation_result`: `python3 scripts/harness/run.py --suite frontend-local` pass
- `coverage_note`: `python3 scripts/harness/run.py --suite coverage` は既存 Sonar issue で fail。全体 statements は `70.5%` で条件 `70.0%` 以上を満たした。

### 変更または追加した部品

- `TranslationManagementStepper`: 翻訳管理段階から `翻訳設定` を除外し、段階番号を詰めた。
- `InputReviewPage`: 固定フッターを `単語翻訳へ進む` へ変更した。
- `AIModelSelectionCard`: 認証不足、モデル一覧状態、実行中編集不可を表示できるようにした。
- `TermTranslationPhasePanel`: 単語翻訳画面へ AI モデル選択領域を追加した。
- `PersonaGenerationPhasePanel`: NPC ペルソナ生成画面へ AI モデル選択領域を追加した。
- `BodyTranslationPhasePanel`: 本文翻訳画面へ AI モデル選択領域を追加した。

### 変更または追加した表示状態

- 翻訳管理段階に `翻訳設定` がない状態。
- 入力データ選択済みから `単語翻訳へ進む` 状態。
- 3 段階の AI 設定未完了状態。
- 認証不足状態。
- モデル一覧 loading / success / failed 状態。
- 実行中に AI モデル選択が編集不可になる状態。

### 確認対象 story

Changed Screens:

- `review-changed-screens-translation-input-inputreviewpage--selected-input-ready`
- `review-changed-screens-term-translation-phase-termtranslationphasepanel--ai-settings-ready`
- `review-changed-screens-term-translation-phase-termtranslationphasepanel--ai-settings-missing`
- `review-changed-screens-term-translation-phase-termtranslationphasepanel--running-locked`
- `review-changed-screens-persona-generation-phase-personagenerationphasepanel--ai-settings-ready`
- `review-changed-screens-persona-generation-phase-personagenerationphasepanel--ai-settings-missing`
- `review-changed-screens-persona-generation-phase-personagenerationphasepanel--running-locked`
- `review-changed-screens-body-translation-phase-bodytranslationphasepanel--ai-settings-ready`
- `review-changed-screens-body-translation-phase-bodytranslationphasepanel--ai-settings-missing`
- `review-changed-screens-body-translation-phase-bodytranslationphasepanel--running-locked`

Changed Components:

- `review-changed-components-translationmanagementstepper--job-management-current`
- `review-changed-components-translationmanagementstepper--term-translation-current`
- `review-changed-components-aimodelselectioncard--fixed-props`
- `review-changed-components-aimodelselectioncard--model-list-loading`
- `review-changed-components-aimodelselectioncard--model-list-failed`
- `review-changed-components-aimodelselectioncard--credential-missing`
- `review-changed-components-aimodelselectioncard--running-locked`
- `review-changed-components-fieldresultlistpanel--default`
- `review-changed-components-fieldresultlistpanel--empty`

### fixture と関連資源

- `frontend/src/ui/screens/translation-job-management/__fixtures__/translation-job-management-fixtures.ts`
- `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts`
- `frontend/src/ui/components/__fixtures__/ai-model-selection-card-fixture.ts`
- `frontend/src/ui/screens/term-translation-phase/__fixtures__/term-phase-card-fixture.ts`
- `frontend/src/ui/screens/persona-generation-phase/__fixtures__/persona-phase-card-fixture.ts`
- `frontend/src/ui/screens/body-translation-phase/__fixtures__/body-phase-card-fixture.ts`
- `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts`
- `frontend/src/ui/screens/term-translation-phase/stories/TermTranslationPhasePanel.stories.ts`
- `frontend/src/ui/screens/persona-generation-phase/stories/PersonaGenerationPhasePanel.stories.ts`
- `frontend/src/ui/screens/body-translation-phase/stories/BodyTranslationPhasePanel.stories.ts`
- `frontend/src/ui/screens/body-translation-phase/stories/FieldResultListPanel.stories.ts`

### Codex 内蔵ブラウザ確認

- `url`: `http://127.0.0.1:6008/iframe.html`
- `checked_stories`: `Review/Changed Screens`, `Review/Changed Components`, `TranslationManagementStepper`, `InputReviewPage`, `AIModelSelectionCard`, `TermTranslationPhasePanel`, `PersonaGenerationPhasePanel`, `BodyTranslationPhasePanel`
- `result`: pass
- `console_error`: 0
- `checked_conditions`: `翻訳設定` 非表示、`単語翻訳へ進む` 表示、AI モデル選択表示、実行中編集不可表示、削除対象の常時診断表示なし

### 人間レビュー受付条件

- Codex 内蔵ブラウザでコメントする場合は、コメント本文、対象 story、対象 selector、frame URL、marker screenshot を 1 件ずつ記録する。
- 承認する場合は `approve` と返す。
- 差し戻す場合は、対象 story と修正対象の部品または表示状態を明示する。

## Frontend Human Review Result

- `review_status`: approved
- `reviewer`: human
- `recorded_at`: 2026-05-23
- `approval_comment`: フロントレビュー終わり
- `storybook_url`: `http://localhost:6008/`
- `review_result`: Storybook 上のフロント見た目レビューは完了した。
- `reviewback_result`: 単語翻訳、本文翻訳、NPC ペルソナ生成のページ構成を `状態概要`、`進行状況兼操作`、`AI モデル選択` に統一した。上部状態概要は `PhaseStatusPanel`、進行状況兼操作は `PhaseProgressPanel`、AI 設定は `AIModelSelectionCard` を共通利用する。
- `docs_result`: frontend 人間レビュー結果を `docs/screen-design/screens/term-translation-phase.md`、`docs/screen-design/screens/persona-generation-phase.md`、`docs/screen-design/screens/body-translation-phase.md` へ反映済み。
- `validation_result`: `python3 scripts/harness/run.py --suite frontend-local` pass。`npm --prefix frontend run build-storybook` pass。
- `storybook_normalization`: approved 後に確認対象 story を通常分類へ戻した。

### 通常分類へ戻した Storybook 資源

Screens:

- `screens-translation-input-inputreviewpage--selected-input-ready`
- `screens-term-translation-phase-termtranslationphasepanel--ai-settings-ready`
- `screens-term-translation-phase-termtranslationphasepanel--ai-settings-missing`
- `screens-term-translation-phase-termtranslationphasepanel--running-locked`
- `screens-persona-generation-phase-personagenerationphasepanel--ai-settings-ready`
- `screens-persona-generation-phase-personagenerationphasepanel--ai-settings-missing`
- `screens-persona-generation-phase-personagenerationphasepanel--running-locked`
- `screens-body-translation-phase-bodytranslationphasepanel--ai-settings-ready`
- `screens-body-translation-phase-bodytranslationphasepanel--ai-settings-missing`
- `screens-body-translation-phase-bodytranslationphasepanel--running-locked`

UI Components:

- `ui-components-translationmanagementstepper--job-management-current`
- `ui-components-translationmanagementstepper--term-translation-current`
- `ui-components-aimodelselectioncard--fixed-props`
- `ui-components-aimodelselectioncard--model-list-loading`
- `ui-components-aimodelselectioncard--model-list-failed`
- `ui-components-aimodelselectioncard--credential-missing`
- `ui-components-aimodelselectioncard--running-locked`
- `ui-components-phase-status-panel-phasestatuspanel--persona-ai-settings-ready`
- `ui-components-phase-status-panel-phasestatuspanel--term-ai-settings-ready`
- `ui-components-phase-status-panel-phasestatuspanel--body-ai-settings-ready`
- `ui-components-phase-progress-panel-phaseprogresspanel--default`

### 人間コメント証跡

- `開始` ボタンと操作パネルが大きすぎるため、操作を進行状況カード内へ統合した。
- 進行状況内の重複項目である `処理済み` と `現在の処理段階` を整理した。
- AI モデル選択と重複する `実行設定` パネルを削除した。
- `結果 summary` と独立した `失敗情報` パネルを削除した。結果は後続の表化対象とする。
- 各翻訳段階を同じページ構成へ揃えた。
- 上部状態概要を単体 component story で確認できるようにした。
- 上部状態概要の metrics を `対象`、`処理済み`、`成功`、`失敗`、`スキップ` に統一し、通常 review 幅では 1 行表示にした。

## Approved Frontend Protection

- `status`: ready
- `approved_screens`: `InputReviewPage`, `TermTranslationPhasePanel`, `PersonaGenerationPhasePanel`, `BodyTranslationPhasePanel`
- `approved_components`: `TranslationManagementStepper`, `AIModelSelectionCard`, `PhaseStatusPanel`, `PhaseProgressPanel`
- `approved_display_rules`: 翻訳管理段階に `翻訳設定` を出さない。入力確認後の主操作は `単語翻訳へ進む` とする。3 phase 画面は `状態概要`、`進行状況兼操作`、`AI モデル選択` の順に表示する。状態概要 metrics は `対象`、`処理済み`、`成功`、`失敗`、`スキップ` の 5 項目を日本語で同順に表示する。操作見出しは視覚表示しない。実行設定、結果 summary、独立失敗情報は常時表示しない。
- `approved_storybook_states`: 通常分類へ戻した Storybook 資源に記録した story。
- `protected_files`: `frontend/src/ui/components/AIModelSelectionCard.svelte`, `frontend/src/ui/components/PhaseStatusPanel.svelte`, `frontend/src/ui/components/PhaseProgressPanel.svelte`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`, `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`, 該当 story / fixture。
- `change_prohibition_for_next_agents`: 後続の backend、integration、test agent は、承認済み画面、承認済み表示規則、通常分類へ戻した Storybook 資源を変更しない。画面、文言、余白、密度、共通 component の変更が必要な場合は frontend 実装へ戻す。

## Merge Readiness

- `merge_ready`: `ready-after-local-commit`
- `source_branch`: `codex/2026-05-22-job-setup-model-selection-relocation`
- `target_branch`: `master`
- `commit_hash`: `local commit created; see git log -1`
- `validation_evidence`: `python3 scripts/harness/run.py --suite backend-local` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。`python3 scripts/harness/run.py --suite structure` pass。`npm --prefix frontend run build-storybook` pass。
- `review_evidence`: `approved-by-human: 2026-05-22 approve`; `frontend-review-approved-by-human: 2026-05-23 フロントレビュー終わり`
- `browser_confirmation`: `/private/tmp/job-setup-relocation-browser-confirmation/snapshot.json` pass。確認対象は `InputReviewPage`、`TranslationManagementStepper`、`TermTranslationPhasePanel`、`PersonaGenerationPhasePanel`、`BodyTranslationPhasePanel`。console error、page error、network failure は 0 件。
- `coverage_result`: `python3 scripts/harness/run.py --suite coverage` は既存 Sonar issue で fail。unit-test handoff 時点の値は backend statements `69.7%`、Sonar line `70.7%`、branch `64.8%`、Sonar Reliability 8、Maintainability HIGH 13。
- `residual_risks`: Wails dev の実データ導線は既存の `ImportTranslationInput` binding 未接続により、入力登録後の実ジョブ作成までは未確認である。Storybook 起動時に既存の `EMFILE` watcher 警告と settings 書き込み `EPERM` が出る。

## Merge Result

- `merge_status`: `merged-to-master`
- `merge_command`: `git merge --no-ff --no-commit codex/2026-05-22-job-setup-model-selection-relocation`
- `conflict_resolution`: `none`
- `post_merge_validation`: `python3 scripts/harness/run.py --suite frontend-local` pass。`python3 scripts/harness/run.py --suite backend-local` pass。`python3 scripts/harness/run.py --suite structure` pass。`npm --prefix frontend run build-storybook` pass。
- `completed_move`: `docs/exec-plans/active/2026-05-22-job-setup-model-selection-relocation/` から `docs/exec-plans/completed/2026-05-22-job-setup-model-selection-relocation/` へ移動済み。
- `source_commit_hash`: `f4ab8385679e528f53850f8943edc92cb95ad750`
- `merge_commit_hash`: `pending-closeout-commit`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `detail_spec_canonicalization`: 実装内容は active plan の `detail-spec-diff.md` と `implementation-scope.md` に保持した。恒久 detail spec への反映は docs 正本化 lane で扱う。
- `follow_up`: なし。

## Outcome

- 設計成果物は人間レビューで承認済みである。
- implementation-scope を作成した。
- frontend 実装と Storybook 人間レビュー依頼を固定した。
- frontend 人間レビューは承認済みである。
- 確認対象 Storybook story は通常分類へ戻した。
- 合意済み frontend 保護を固定した。
- backend 実装、統合境界実装、観測ログ追加、単体テスト、シナリオテストを完了した。
- 最終検証と実装後ブラウザ確認は通過した。
- source branch は master へ local merge 済みである。
- conflict は発生していない。
- merge 後検証は通過した。
- active plan folder は completed archive へ移動済みである。
- remote repository は変更していない。
