# Implementation Scope: 2026-05-17-all-pages-componentization

- `skill`: implementation-scope
- `status`: handoff-ready
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: `2026-05-18 human message: approve`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `human_decision_questionnaire`: `N/A`
- `design_diff`: `./design-diff-all-pages-componentization.puml`
- `component_split_by_page`: `./component-split-by-page.puml`
- `component_folder_guideline`: `./component-folder-guideline.md`
- `storybook_foundation`: `../../completed/2026-05-18-storybook-foundation/implementation-scope.md`
- `storybook_foundation_review`: `../../completed/2026-05-18-storybook-foundation/storybook-review.md`
- `master_persona_poc_ui`: `../../completed/2026-05-17-master-persona-componentization/ui-design.md`

## Fixed Decisions

- `needs_human_decision`: `0`
- この task は frontend refactor と Storybook story 追加を扱う。
- frontend handoff は backend handoff より先に置く。今回は backend handoff は不要である。
- backend 実装は不要である。理由は、プロダクト backend、Wails binding、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を変更しないためである。
- 統合境界実装は不要である。理由は、API、Wails、DTO、gateway、adapter 契約を変更しないためである。
- UI component の public seam は Svelte component props、callback、story args である。
- page component は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
- 共有 component は複数画面で同じ表示規則または操作規則を持つ場合だけ `frontend/src/ui/components/` へ置く。
- 画面固有条件が増える component は `frontend/src/ui/screens/<screen>/` に残す。
- 既存表示項目、状態文、操作、`aria-label` を削る変更は、この implementation-scope の対象外である。
- docs 正本化、`.codex/` 変更、作業流れ変更は implementation handoff に含めない。
- `screen_design_diff`: `N/A`。理由は、画面設計書正本へ反映する恒久的な画面内容差分がないためである。
- Storybook story と fixture は fixed props、view model fixture、callback stub だけを使う。
- Storybook story と fixture は backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。
- Storybook review URL、story ID、確認状態、未確認理由、再実行 command、`build-storybook` 結果は task-local `storybook-review.md` に残す。
- `npm --prefix frontend run build-storybook` は Storybook 専用 gate として扱う。
- `python3 scripts/harness/run.py --suite frontend-local` は全 frontend 変更後の最終 gate として扱う。
- `secret_boundary`: `not_required`。ただし story、fixture、review 記録、検証記録には secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含めない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `APC-FE-01-shared-form-action-primitives` | `なし` | `なし` | `なし` |
| `wave-2` | `APC-FE-02-shared-status-list-primitives` | `APC-FE-01-shared-form-action-primitives` | `なし` | `depends_on` |
| `wave-3` | `APC-FE-03-input-setup-pages`, `APC-FE-04-job-management-page`, `APC-FE-05-job-run-shell`, `APC-FE-06-term-persona-phase-pages`, `APC-FE-07-body-complete-phase-pages`, `APC-FE-08-output-artifact-page`, `APC-FE-09-master-dictionary-page`, `APC-FE-10-provider-dashboard-masterpersona-pages` | `APC-FE-02-shared-status-list-primitives` | `APC-FE-03-input-setup-pages <-> APC-FE-04-job-management-page`, `APC-FE-05-job-run-shell <-> APC-FE-08-output-artifact-page`, `APC-FE-06-term-persona-phase-pages <-> APC-FE-09-master-dictionary-page`, `APC-FE-07-body-complete-phase-pages <-> APC-FE-10-provider-dashboard-masterpersona-pages` | `なし` |
| `wave-4` | `APC-UT-11-component-boundary-regression`, `APC-ST-12-storybook-review-evidence` | `wave-3 の frontend handoff 全件` | `APC-UT-11-component-boundary-regression <-> APC-ST-12-storybook-review-evidence` | `なし` |

## Dependency Table

| handoff_id | depends_on | reason |
| --- | --- | --- |
| `APC-FE-01-shared-form-action-primitives` | `なし` | form、button、inline feedback は後続画面部品と modal の土台である。 |
| `APC-FE-02-shared-status-list-primitives` | `APC-FE-01-shared-form-action-primitives` | search、file 表示、empty、progress、status、pagination、confirm modal は action と form primitive を再利用する可能性がある。 |
| `APC-FE-03-input-setup-pages` | `APC-FE-02-shared-status-list-primitives` | 翻訳入力確認とジョブ作成は shared primitive と既存 `AIModelSelectionCard`、`StickyActionFooter` を前提にする。 |
| `APC-FE-04-job-management-page` | `APC-FE-02-shared-status-list-primitives` | job 管理は status、confirm modal、list、action primitive を前提にする。 |
| `APC-FE-05-job-run-shell` | `APC-FE-02-shared-status-list-primitives` | job run shell は footer、phase navigation、未選択案内の shared 表示規則を前提にする。 |
| `APC-FE-06-term-persona-phase-pages` | `APC-FE-02-shared-status-list-primitives` | 単語翻訳段階と NPC ペルソナ生成段階は phase status、action、progress、failure、metric 候補を前提にする。 |
| `APC-FE-07-body-complete-phase-pages` | `APC-FE-02-shared-status-list-primitives` | 本文翻訳段階と翻訳完了は phase status、progress、pagination、readiness 表示を前提にする。 |
| `APC-FE-08-output-artifact-page` | `APC-FE-02-shared-status-list-primitives` | 出力成果物は status、form、file 表示、action primitive を前提にする。 |
| `APC-FE-09-master-dictionary-page` | `APC-FE-02-shared-status-list-primitives` | マスター辞書は pagination、confirm modal、form、file 表示、search を前提にする。 |
| `APC-FE-10-provider-dashboard-masterpersona-pages` | `APC-FE-02-shared-status-list-primitives` | Provider 設定、Dashboard、Master Persona review 連携は status、form、search、既存 POC story を前提にする。 |
| `APC-UT-11-component-boundary-regression` | `wave-3 の frontend handoff 全件` | component、story、fixture の import 境界は対象 story が揃った後に検査する。 |
| `APC-ST-12-storybook-review-evidence` | `wave-3 の frontend handoff 全件` | review URL と story ID は主要 story が揃った後に確定できる。 |

## Handoffs

### `APC-FE-01-shared-form-action-primitives`

- `implementation_target`: 共有 form、action、feedback primitive を追加し、Storybook story と fixture を用意する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic label、synthetic value、synthetic option、synthetic error message
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: story、fixture、review 記録、test failure に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータを出さない。
- `owned_scope`:
  - `frontend/src/ui/components/ActionButton.svelte`、`ButtonGroup.svelte`、`IconActionButton.svelte` を追加または既存操作表示から抽出する。
  - `frontend/src/ui/components/FormField.svelte`、`TextInputField.svelte`、`TextAreaField.svelte`、`SelectField.svelte`、`CheckboxField.svelte` を追加する。
  - `frontend/src/ui/components/InlineFeedback.svelte` を追加する。
  - 各 shared component の story と fixture を `frontend/src/ui/components/` と `frontend/src/ui/components/__fixtures__/` に置く。
  - component は props、callback、slot 相当の表示入力だけを使う。
  - component は `Store`、Gateway、generated binding、backend DTO、controller factory、RuntimeEventAdapter を import しない。
- `estimated_change_size`:
  - `files`: `18-24 files`
  - `changed_lines`: `850-1400 lines`
  - `size_class`: `注意`
  - `split_reason`: form と action primitive は互いに story 上の組み合わせが必要であり、分けると後続 page handoff が二重に依存するため 1 handoff にする。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/components/ActionButton.svelte` を追加し、primary、secondary、danger、busy、disabled の表示規則を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - shared primitive は props と callback だけで表示できる。
  - shared primitive は production controller、Store、Gateway、generated binding、backend DTO を import しない。
  - action、form、feedback の story は通常、disabled、error、help、required、busy、長文の代表状態を持つ。
  - `npm --prefix frontend run build-storybook` が成功する。
  - 既存 `AIModelSelectionCard` と `StickyActionFooter` を置き換えず、必要な場合だけ story 上で併用できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run build-storybook` は、この handoff が追加する shared story の static build 破綻を即時に検出するため途中 gate として置く。
  - 画面表示文言の恒久変更はこの handoff で扱わない。

### `APC-FE-02-shared-status-list-primitives`

- `implementation_target`: status、list、file、progress、pagination、confirm modal primitive を追加し、Storybook story と fixture を用意する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic status label、synthetic file name、synthetic path label、synthetic hash label、synthetic count
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story args、screenshot 名、review 記録に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータを出さない。
- `owned_scope`:
  - `frontend/src/ui/components/SearchFilterBar.svelte`、`FileSelectionDisplay.svelte`、`EmptyStatePanel.svelte`、`ProgressBar.svelte`、`StatusPill.svelte`、`PaginationControls.svelte`、`ConfirmDangerModal.svelte` を追加する。
  - 各 shared component の story と fixture を `frontend/src/ui/components/` と `frontend/src/ui/components/__fixtures__/` に置く。
  - `ConfirmDangerModal` は対象識別情報、処理中、削除失敗、確定操作を props と callback で表す。
  - `FileSelectionDisplay` は file picker、file read、実 filesystem flow を内部に持たない。
  - `StatusPill` と `ConfirmDangerModal` は props が画面固有条件の塊になる場合、screen local component へ戻す。
- `estimated_change_size`:
  - `files`: `14-20 files`
  - `changed_lines`: `650-1100 lines`
  - `size_class`: `注意`
  - `split_reason`: list、status、modal は複数画面から同時に使うため、画面別 handoff へ散らさず shared primitive handoff に集約する。
- `depends_on`: `APC-FE-01-shared-form-action-primitives`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`: `frontend/src/ui/components/StatusPill.svelte` を追加し、状態 label と tone を少数 props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - list、status、file、progress、pagination、confirm modal の story は通常、空、失敗、長文、処理中、危険操作の代表状態を持つ。
  - component は file picker、file read、secret 保存、runtime event 購読を内部に持たない。
  - component は production controller、Store、Gateway、generated binding、backend DTO を import しない。
  - `npm --prefix frontend run build-storybook` が成功する。
  - props が画面固有条件で肥大化する component は screen local へ戻した理由を実装結果に残す。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run build-storybook` は shared story と shared fixture の孤立描画を検証するため途中 gate として置く。

### `APC-FE-03-input-setup-pages`

- `implementation_target`: 翻訳入力確認とジョブ作成を screen local component へ分け、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic input name、synthetic plugin name、synthetic count、synthetic validation message
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に実 file path、実 input data、raw request、raw response、secret、token を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/translation-input/` の `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail`、入力登録後導線を story 対象にする。
  - `frontend/src/ui/screens/translation-job-setup/` に `JobSetupPurposeHeader`、`InputSourcePanel`、`FoundationDataPanel`、`PhaseSettingsPanel`、`CompatibilityPrecheckPanel`、`CreatedJobSummaryPanel`、`PhaseSettingsSummaryPanel` を追加または抽出する。
  - `InputReviewPage` と `JobSetupPage` は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
  - file input bridge、controller 購読、日時と状態の表示変換、作成可否の集約、phase 別設定の controller 更新、入力削除中 state は page または controller 周辺に残す。
  - screen local story と fixture は `frontend/src/ui/screens/<screen>/stories/` と `frontend/src/ui/screens/<screen>/__fixtures__/` に置く。
- `estimated_change_size`:
  - `files`: `18-24 files`
  - `changed_lines`: `900-1400 lines`
  - `size_class`: `注意`
  - `split_reason`: 翻訳入力確認とジョブ作成は連続する入力準備 use case であり、fixture の synthetic input と作成可否を共有するため 1 handoff にする。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-04-job-management-page`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/translation-job-setup/JobSetupPurposeHeader.svelte` を追加し、ジョブ作成 page の header 表示を props で分離する completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 翻訳入力確認 story は file 未選択、file 選択済み、読み込み中、失敗、登録後導線を確認できる。
  - ジョブ作成 story は入力 source、基盤データ、phase 設定、互換性確認、作成済み summary、固定フッターを確認できる。
  - 既存の `この JSON を登録`、`翻訳設定へ進む`、`単語翻訳へ進む` の操作文言と可否を維持する。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - story と fixture は実 file picker、実 file read、backend DTO、Gateway mock を要求しない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run build-storybook` はこの page group の story 追加直後に不足 import と Svelte compile error を検出するため途中 gate として置く。

### `APC-FE-04-job-management-page`

- `implementation_target`: 翻訳ジョブ管理を screen local component へ分け、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic job id、synthetic job title、synthetic phase label、synthetic disabled reason
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に実 job data、実 endpoint、secret、token、raw response を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/translation-job-management/` に `JobManagementHeader`、`FeedbackPanel`、`JobListPanel`、`JobCard`、`JobOperationGroup` を追加または抽出する。
  - `TranslationJobManagementDeleteModal` と `TranslationManagementStepper` の story と fixture を追加する。
  - `TranslationJobManagementPage` は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
  - job 選択から job-run へ進む route 同期、停止と再開の実行判断は page または controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `14-20 files`
  - `changed_lines`: `700-1100 lines`
  - `size_class`: `注意`
  - `split_reason`: job 管理は一覧、操作、削除 modal の単一検証意図で閉じられるため 1 handoff にする。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-03-input-setup-pages`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/translation-job-management/JobCard.svelte` を追加し、job 状態、選択状態、操作可否、disabled reason を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - job 管理 story は一覧、選択中、停止、再開、削除、削除失敗、フィードバックを確認できる。
  - 既存の `現在の翻訳段階へ進む`、停止、再開、削除を維持する。
  - job state の表示責務と操作可否を維持する。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - story と fixture は production Gateway、RuntimeEventAdapter、generated binding を要求しない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `TranslationManagementStepper` は翻訳管理画面の状態説明に属するため、この handoff の review 対象に含める。

### `APC-FE-05-job-run-shell`

- `implementation_target`: ジョブ実行 shell を screen local component へ分け、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic job id、synthetic job title、synthetic phase label、synthetic disabled reason
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に実 job data、実 endpoint、secret、token、raw response を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/job-run/` に `JobRunTargetSummary`、`JobUnselectedGuidance`、`PhaseHost`、`PhaseNavigationFooter` の story と fixture を追加する。
  - `JobRunPage` は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
  - job target 同期、phase page 切替、複数 controller の mount は page または controller 周辺に残す。
  - `PhaseNavigationFooter` は共通化できる場合だけ shared component へ移す。phase 固有条件が増える場合は `frontend/src/ui/screens/job-run/` に残す。
- `estimated_change_size`:
  - `files`: `10-16 files`
  - `changed_lines`: `500-900 lines`
  - `size_class`: `通常`
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-08-output-artifact-page`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/job-run/JobRunTargetSummary.svelte` を追加し、選択中 job の target summary を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - job run story は target summary、未選択案内、phase host、phase navigation footer を確認できる。
  - 直リンク相当の job 未選択時は、未完了一覧へ戻る案内を維持する。
  - job state と phase run state の表示責務を混ぜない。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - story と fixture は production Gateway、RuntimeEventAdapter、generated binding を要求しない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - phase の実行詳細は phase page handoff が扱う。この handoff は job run shell の合成境界だけを扱う。

### `APC-FE-06-term-persona-phase-pages`

- `implementation_target`: 単語翻訳段階と NPC ペルソナ生成段階を component 化し、phase 系 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic phase label、synthetic progress、synthetic failure summary、synthetic metric count、synthetic translation text
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に raw request、raw response、raw prompt、provider 応答原文、実ユーザーデータ、secret、token を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/term-translation-phase/` に `TermExecutionSettingsCard`、`TermResultSummaryCard` と story、fixture を追加または抽出する。
  - `frontend/src/ui/screens/persona-generation-phase/` に `PersonaTargetSummaryCard`、`PersonaExecutionSettingsCard`、`PersonaResultSummaryCard`、`BodyReadinessInputCard` と story、fixture を追加または抽出する。
  - phase 系共通候補 `PhaseStatusPanel`、`PhaseActionPanel`、`PhaseProgressPanel`、`PhaseFailureInfoCard`、`PhaseMetricCounterGrid` は、同じ表示規則を持つ場合だけ shared component へ置く。
  - phase action dispatch、本文翻訳 readiness 判断、phase action の割り当ては page、presenter、controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `16-24 files`
  - `changed_lines`: `800-1300 lines`
  - `size_class`: `注意`
  - `split_reason`: 単語翻訳段階と NPC ペルソナ生成段階は phase 状態、操作、進捗、失敗情報の story 検証意図を共有するため 1 handoff にする。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-09-master-dictionary-page`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/components/PhaseStatusPanel.svelte` または screen local fallback を追加し、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` 相当の状態表示 clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - phase story は `Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` 相当を確認できる。
  - phase story は開始、中断、再開、再試行、キャンセル、次段階確認の可否と disabled reason を確認できる。
  - failure、progress、readiness、result summary、metric counter を fixture で確認できる。
  - job state と phase run state の意味を混ぜない。
  - shared 化で見出し差や readiness 差が失われる場合は screen local component へ戻す。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `PhaseStatusPanel` などの shared 化は、この handoff で最初に判断する。props が肥大化する場合は各 phase screen local に残す。

### `APC-FE-07-body-complete-phase-pages`

- `implementation_target`: 本文翻訳段階と翻訳完了を component 化し、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic phase label、synthetic progress、synthetic translation text、synthetic field key、synthetic readiness reason
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に raw request、raw response、raw prompt、provider 応答原文、実ユーザーデータ、secret、token を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/body-translation-phase/` に `BodyInputSummaryCard`、`BodyExecutionSummaryCard`、`BodyResultSummaryCard`、`FieldResultListPanel`、`OutputReadinessCard` と story、fixture を追加または抽出する。
  - `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte` から `TranslationCompleteSummaryPanel`、`TranslationResultListPanel` を追加または抽出し、story と fixture を追加する。
  - body phase は phase shared component を使える場合だけ再利用する。
  - output readiness 判定、field result key 生成、pageIndex 内部 state は page、presenter、controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `14-22 files`
  - `changed_lines`: `750-1200 lines`
  - `size_class`: `注意`
  - `split_reason`: 本文翻訳段階と翻訳完了は output readiness と結果一覧の検証意図が連続するため 1 handoff にする。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-10-provider-dashboard-masterpersona-pages`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/body-translation-phase/BodyInputSummaryCard.svelte` を追加し、本文翻訳入力 summary と readiness reason を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 本文翻訳 story は入力 summary、実行 summary、結果 summary、field result、出力 readiness、失敗情報を確認できる。
  - 翻訳完了 story は完了 summary、原文訳文一覧、ページング、0 件、長文を確認できる。
  - output readiness 未達では次段階へ進めない理由を表示する。
  - story と fixture は raw request、raw response、raw prompt、provider 応答原文を持たない。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - field result の key 生成は表示部品へ移さない。表示部品は親から渡された key と label だけを表示する。

### `APC-FE-08-output-artifact-page`

- `implementation_target`: 出力成果物を screen local component へ分け、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic output path label、synthetic artifact summary、synthetic completed job summary
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に実 file path、実 output data、raw request、raw response、secret、token を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/translation-output-artifact/` に `OutputSummaryHeader`、`CompletedJobListPanel`、`SelectedJobSummaryCard`、`OutputActionPanel`、`LatestOutputResultCard`、`DiffPreviewPanel` と story、fixture を追加または抽出する。
  - `TranslationOutputArtifactPage` は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
  - path validation 表示と artifact 生成可否の集約は page、presenter、controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `12-18 files`
  - `changed_lines`: `600-1000 lines`
  - `size_class`: `注意`
  - `split_reason`: 出力成果物は出力候補、選択 summary、実行操作、結果、diff preview の単一検証意図で閉じる。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-05-job-run-shell`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/translation-output-artifact/OutputActionPanel.svelte` を追加し、出力操作、path label、実行可否、disabled reason を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 出力成果物 story は出力候補、選択 job summary、出力操作、最新 result、diff preview、0 件、長文 path を確認できる。
  - `XML を出力` と `再出力` を維持する。
  - story と fixture は実 file picker、実 file read、実 output generation を開始しない。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - path は synthetic label として扱う。実 local path は fixture と review 記録へ入れない。

### `APC-FE-09-master-dictionary-page`

- `implementation_target`: マスター辞書を screen local component へ分け、主要 story と fixture を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic dictionary entry、synthetic search text、synthetic import summary
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に実 file path、実 dictionary data、raw request、raw response、secret、token を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/master-dictionary/` に `DictionaryHeader`、`DictionaryImportPanel`、`DictionaryListPanel`、`DictionaryDetailPanel`、`DictionaryEditModal`、`DictionaryDeleteModal` と story、fixture を追加または抽出する。
  - `MasterDictionaryPage` は controller 接続、購読、dispose、通知、表示部品合成へ寄せる。
  - file input bridge、一覧検索と詳細選択の一体更新は page、presenter、controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `14-22 files`
  - `changed_lines`: `700-1200 lines`
  - `size_class`: `注意`
  - `split_reason`: 辞書一覧、詳細、作成編集、削除 modal は辞書管理の単一検証意図で閉じる。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-06-term-persona-phase-pages`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/master-dictionary/DictionaryListPanel.svelte` を追加し、一覧、検索、選択、空状態を props で表す completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - マスター辞書 story は XML 取り込み、辞書一覧、詳細、作成編集モーダル、削除モーダル、保存失敗、削除失敗を確認できる。
  - 辞書の作成、保存、削除、確認操作を維持する。
  - modal 失敗時は dialog を閉じず、対象識別情報と入力値を保持する。
  - story と fixture は実 file picker、実 file read、実 dictionary import を開始しない。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 一覧検索と詳細選択の一体更新は表示部品へ移さない。表示部品は選択済み id と表示 props を受け取る。

### `APC-FE-10-provider-dashboard-masterpersona-pages`

- `implementation_target`: Provider 設定、Dashboard、Master Persona review 連携を screen local component と story 対象に整理する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: credential reference label、masked credential state、synthetic provider name、synthetic route label、synthetic persona label
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: fixture、story、review 記録に API key、token、secret 本体、実 endpoint、raw request、raw response、raw prompt、provider 応答原文を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/provider-settings/` に `ProviderSettingsSummaryPanel`、`ProviderListPanel`、`ProviderDetailPanel`、`ApiKeyPanel`、`ConnectionCheckPanel`、`SettingsActionPanel` と story、fixture を追加または抽出する。
  - Dashboard 表示がある実装 path に `AppHeader`、`GlobalNavigation`、`CurrentPageHero`、`DashboardEntryGrid`、`DashboardEntryCard` を screen local component として追加または抽出し、story と fixture を追加する。
  - Master Persona は POC 済み `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` の story を review 対象一覧へ含める。
  - Master Persona の部品境界を再設計しない。必要な変更は shared primitive 互換または review 対象一覧への接続に限定する。
  - secret 入力 draft の破棄、credential input の controller 同期、AppShell hash 同期、mobile nav 開閉、主要 route orchestration は page、shell、controller 周辺に残す。
- `estimated_change_size`:
  - `files`: `18-25 files`
  - `changed_lines`: `850-1450 lines`
  - `size_class`: `注意`
  - `split_reason`: Provider 設定と Dashboard は app shell と設定導線の review 対象であり、Master Persona は POC story を参照するだけのため 1 handoff にする。26 files 以上の見込みが出た場合は Provider 設定と Dashboard で再分割する。
- `depends_on`: `APC-FE-02-shared-status-list-primitives`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `APC-FE-07-body-complete-phase-pages`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/provider-settings/ApiKeyPanel.svelte` を追加し、credential draft、masked state、保存禁止値を表示へ出さない completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - Provider 設定 story は設定 summary、AI サービス一覧、設定詳細、API キー panel、接続確認、保存操作を確認できる。
  - Provider 設定 story は secret 本体を表示、fixture、review 記録へ出さない。
  - Dashboard story は app header、global navigation、current page hero、entry grid、entry card、狭い幅 navigation を確認できる。
  - Master Persona POC story は全ページ review 対象一覧に含まれ、再設計しない。
  - 既存表示項目、状態文、操作、`aria-label` を削らない。
  - `npm --prefix frontend run build-storybook` が成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `credential reference label` は UI に出してよい参照値である。API key、token、secret 本体は UI、DTO、fixture、review 記録に出さない。

### `APC-UT-11-component-boundary-regression`

- `implementation_target`: component、story、fixture の import 境界と型境界を lower-level test で確認する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic fixture value
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: test fixture、test name、test failure に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータを出さない。
- `owned_scope`:
  - 既存 boundary lint または unit test に、Storybook story と fixture が許可対象であることを確認する test を追加または更新する。
  - product component から Storybook package、Storybook runtime、story 専用 module への import 禁止を確認する。
  - shared component と screen local component が `Store`、Gateway、generated binding、backend DTO、RuntimeEventAdapter を直接 import しないことを確認する。
  - props 型、callback 型、story args 型が `any` や過剰な type assertion に依存しないことを既存 lint と type check で確認する。
- `estimated_change_size`:
  - `files`: `2-5 files`
  - `changed_lines`: `80-220 lines`
  - `size_class`: `通常`
- `depends_on`: `APC-FE-03-input-setup-pages`, `APC-FE-04-job-management-page`, `APC-FE-05-job-run-shell`, `APC-FE-06-term-persona-phase-pages`, `APC-FE-07-body-complete-phase-pages`, `APC-FE-08-output-artifact-page`, `APC-FE-09-master-dictionary-page`, `APC-FE-10-provider-dashboard-masterpersona-pages`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `APC-ST-12-storybook-review-evidence`
- `parallel_blockers`: `なし`
- `first_action`: 既存 boundary lint test に Storybook story と fixture の許可 pattern を追加し、story だけが Storybook runtime を import できる completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run lint:boundaries`
  - `npm --prefix frontend run lint:types`
- `completion_signal`:
  - Storybook story と fixture は Storybook 専用 import の許可対象として扱われる。
  - product component は Storybook package、Storybook runtime、story 専用 module を import できない。
  - shared component と screen local component は `Store`、Gateway、generated binding、backend DTO、RuntimeEventAdapter を直接 import しない。
  - TypeScript type check が通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は `npm --prefix frontend run build-storybook` を担当しない。Storybook render と static build は frontend handoff と final validation で扱う。

### `APC-ST-12-storybook-review-evidence`

- `implementation_target`: 全ページの Storybook review 証跡を `storybook-review.md` に残し、UI 人間操作 E2E の入口を固定する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: Storybook localhost URL、story ID、confirmation status
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `storybook-review.md` の URL、story ID、screenshot 名、command 記録に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を出さない。
- `owned_scope`:
  - `docs/exec-plans/active/2026-05-17-all-pages-componentization/storybook-review.md` を作成する。
  - review 対象の story ID、review URL、iframe URL、確認状態、未確認理由、再実行 command、`build-storybook` 結果を記録する。
  - Page Review Order に沿って、翻訳入力確認、ジョブ作成、翻訳ジョブ管理、ジョブ実行、単語翻訳段階、NPC ペルソナ生成段階、本文翻訳段階、翻訳完了、出力成果物、マスター辞書、Provider 設定、Dashboard、Master Persona を確認対象にする。
  - page 合成 story は密度確認と配置確認に限定し、主要部品 story の不足を代替しないことを記録する。
  - UX review と frontend human review は、この handoff の実装後に人間レビューへ進める前提として未承認または確認待ち状態を記録する。
- `estimated_change_size`:
  - `files`: `1-3 files`
  - `changed_lines`: `80-180 lines`
  - `size_class`: `通常`
- `depends_on`: `APC-FE-03-input-setup-pages`, `APC-FE-04-job-management-page`, `APC-FE-05-job-run-shell`, `APC-FE-06-term-persona-phase-pages`, `APC-FE-07-body-complete-phase-pages`, `APC-FE-08-output-artifact-page`, `APC-FE-09-master-dictionary-page`, `APC-FE-10-provider-dashboard-masterpersona-pages`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `APC-UT-11-component-boundary-regression`
- `parallel_blockers`: `なし`
- `first_action`: `docs/exec-plans/active/2026-05-17-all-pages-componentization/storybook-review.md` の skeleton を作成し、review URL、story ID、確認状態、未確認理由を記録する completion clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run storybook`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - Storybook dev server または static preview で主要 story を開ける。
  - `storybook-review.md` は story ID、review URL、確認状態、未確認理由、再実行 command、`build-storybook` 結果を持つ。
  - 各画面の主要 panel、card、modal が review 対象として列挙される。
  - story 欠落は未確認理由として残り、確認済み扱いにしない。
  - review 証跡は fakeAPI URL、Wails runtime URL、backend API URL を Storybook review URL として扱わない。
  - review 証跡は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含まない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run storybook` は long-running command になる可能性がある。実装者は起動 URL を記録した後に process を停止する。
  - `agent-browser` が使える場合は Storybook review URL を開いて screenshot、snapshot、errors を取得する。環境で失敗する場合は未確認理由と代替確認を `storybook-review.md` に残す。
  - この handoff は UI 人間操作 E2E と Storybook review 証跡の入口である。frontend human review の承認自体は人間が行う。

## Final Validation

全 handoff 完了後に `implement_lane` が最終検証として扱う。

- `npm --prefix frontend run lint`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

最終検証の完了条件は次の通りである。

- `npm --prefix frontend run lint` は Storybook import 境界、型、eslint、export 検査を含む既存 frontend lint 入口として成功する。
- `npm --prefix frontend run build-storybook` は全 story の static build gate として成功する。
- `python3 scripts/harness/run.py --suite frontend-local` は既存 frontend local gate として成功する。
- `storybook-review.md` は review URL、story ID、確認状態、未確認理由、再実行 command、`build-storybook` 結果、残留リスクを持つ。
- UX review と frontend human review は、実装後に別途人間確認へ進める状態になっている。
- 最終検証記録には secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含めない。

広域検証を最終検証へ置く理由は次の通りである。

- `python3 scripts/harness/run.py --suite frontend-local` は frontend 全体の lint と test を束ねるため、途中 handoff の単独責任にできない。
- `npm --prefix frontend run lint` は shared component、screen local component、story、fixture の全体 import 境界を確認するため、frontend 実装全件後に実行する。
- `npm --prefix frontend run build-storybook` は途中 handoff でも局所 gate として使うが、最終検証では全 story の統合 gate として再実行する。

## Non Targets

- backend 実装は扱わない。
- 統合境界実装は扱わない。
- API、Wails binding、DTO、gateway、adapter 契約の変更は扱わない。
- DB、secret store、AI provider、filesystem flow の変更は扱わない。
- docs 正本化、`.codex/` 変更、skill 変更、agent 変更は扱わない。
- screen design 正本の変更は扱わない。
- 未承認の表示文言変更、導線変更、画面削除は扱わない。
- Storybook 上で実 AI 生成、実ファイル読み込み、DB 書き込み、provider network、xTranslator 出力生成を再現しない。
- Master Persona POC の部品境界を再設計しない。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `storybook-review.md`
- `ux_review_result`: 実装後 UX review の状態または未実施理由
- `frontend_human_review_result`: 人間見た目レビューの状態または未実施理由
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、ui_evidence、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`

## Missing Information

- `none`
