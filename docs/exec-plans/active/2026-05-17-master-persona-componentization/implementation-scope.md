# Implementation Scope: 2026-05-17-master-persona-componentization

- `skill`: implementation-scope
- `status`: handoff-ready
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: `2026-05-18 human approved scenario-design.md, ui-design.md, and design-diff by instructing: story作成まで進めて`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `design_diff`: `./design-diff-master-persona-componentization.puml`
- `storybook_foundation_scope`: `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`
- `storybook_foundation_review`: `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`

## Fixed Decisions

- `needs_human_decision`: `0`
- UI 変更は承認済み `ui-design.md` を source にする。
- 画面設計差分は `N/A` とする。利用者に見える画面内容の追加、削除、文言変更は扱わない。
- `AIModelSelectionCard` は既存共有 component と既存 story を基準にして再利用する。
- `MasterPersonaPage` は controller 接続、file input bridge、notice / status、表示部品の props 合成だけを持つ。
- screen local component は `viewModel` 全体、`Store`、`Gateway`、generated binding、RuntimeEventAdapter を直接扱わない。
- Storybook story と fixture は fixed props、view model 由来の synthetic fixture、callback stub だけを使う。
- Storybook fixture、story、review 記録は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を含めない。
- backend 実装、統合境界実装、docs 正本化、`.codex` 変更は Codex implementation lane の対象外である。
- `secret_boundary`: `not_required`。ただし story、fixture、review 記録には禁止情報を出さない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `MPC-FE-01-master-persona-panel-stories` | `なし` | `なし` | `なし` |
| `wave-2` | `MPC-UT-02-component-boundary-tests`, `MPC-ST-03-storybook-review-evidence` | `MPC-FE-01-master-persona-panel-stories` | `MPC-UT-02-component-boundary-tests <-> MPC-ST-03-storybook-review-evidence` | `なし` |

## Dependency Table

| handoff_id | depends_on | reason |
| --- | --- | --- |
| `MPC-FE-01-master-persona-panel-stories` | `なし` | panel / modal の props 境界と story が後続検証の前提である。 |
| `MPC-UT-02-component-boundary-tests` | `MPC-FE-01-master-persona-panel-stories` | component props、story、fixture の path と import 境界が確定した後に検査できる。 |
| `MPC-ST-03-storybook-review-evidence` | `MPC-FE-01-master-persona-panel-stories` | Storybook dev URL と story ID は story 作成後に確定する。 |

## Handoffs

### `MPC-FE-01-master-persona-panel-stories`

- `implementation_target`: `MasterPersonaPage` を薄くし、4 つの screen local component を small props と callback だけで story 表示できる境界へ寄せる。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `story_id`、Storybook localhost URL、synthetic fixture の表示値
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: story、fixture、callback stub、review 記録、command 記録に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を出さない。
- `owned_scope`:
  - `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte` は controller mount、dispose、subscribe、file input bridge、notice / status、panel props 合成だけを持つ。
  - `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte` は生成準備、AI 設定、入力 JSON、preview 件数、作成開始 callback を small props で受け取る。
  - `frontend/src/ui/screens/master-persona/RunStatusPanel.svelte` は進行状態、進捗、処理件数、現在対象、一時停止、中止 callback を small props で受け取る。
  - `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte` は検索、plugin 絞り込み、一覧、ページ操作、詳細、編集開始、削除開始 callback を small props で受け取る。
  - `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte` は編集、削除、保存失敗、削除失敗の対象保持を small props と callback で受け取る。
  - `frontend/src/ui/screens/master-persona/__fixtures__/` に synthetic fixture と callback stub を置く。
  - `frontend/src/ui/screens/master-persona/*.stories.ts` に `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` の story を追加する。
  - `AIModelSelectionCard` は `frontend/src/ui/components/AIModelSelectionCard.svelte` を再利用し、手作り代替 card を追加しない。
  - story と fixture は backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を import しない。
- `estimated_change_size`:
  - `files`: `10-14 files`
  - `changed_lines`: `500-780 lines`
  - `size_class`: `通常`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte` の `Props` を setup 用 small props へ置き換え、`SCN-MPC-002` の props boundary clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - `GenerationSetupPanel` は未選択、JSON 選択済み、preview 成功、AI 設定不足、モデル一覧更新中、長い model 名を story で表示できる。
  - `RunStatusPanel` は生成前、生成中、生成失敗、生成完了を story で表示できる。
  - `PersonaReviewPanel` は空一覧、一覧あり、行選択済み、絞り込み後空、古い選択なしを story で表示できる。
  - `PersonaActionModal` は編集、削除、保存失敗、削除失敗を story で表示できる。
  - `MasterPersonaPage` は screen local component へ `viewModel` 全体を渡さない。
  - screen local component は `Store`、`Gateway`、generated binding、RuntimeEventAdapter を import しない。
  - Storybook build は backend、Wails runtime、AI provider、secret store、DB、実 filesystem なしで成功する。
  - frontend local gate は成功、または task と無関係な環境理由を持つ。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff が最初に実行可能な frontend handoff である。
  - `python3 scripts/harness/run.py --suite frontend-local` は広域 frontend gate だが、この handoff は page と panel の production compile を変えるため実装後 gate として必要である。
  - page 合成 story は任意である。主要表示領域 story の不足を page 合成 story で代替しない。

### `MPC-UT-02-component-boundary-tests`

- `implementation_target`: master-persona の story / fixture / component 境界を lower-level test で保護する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: synthetic fixture の表示値、story file path
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: test fixture、snapshot、test 名、failure message に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を出さない。
- `owned_scope`:
  - 既存 Vitest または boundary test に、master-persona story / fixture の禁止 import 検査を追加する。
  - `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` が `Store`、`Gateway`、generated binding、RuntimeEventAdapter を直接 import しないことを検査する。
  - story / fixture が backend、Wails runtime、AI provider、secret store、DB、実 filesystem flow を import しないことを検査する。
  - component の外から観測できる disabled、callback、表示状態だけを検査対象にする。
- `estimated_change_size`:
  - `files`: `1-3 files`
  - `changed_lines`: `60-180 lines`
  - `size_class`: `通常`
- `depends_on`: `MPC-FE-01-master-persona-panel-stories`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `MPC-ST-03-storybook-review-evidence`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/repository-boundary-plugin.test.mjs` または同等の boundary test に master-persona story / fixture の許可 path と禁止 import case を追加し、`SCN-MPC-006` の forbidden import clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run lint:boundaries`
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - lower-level test は screen local component が禁止境界を import しないことを証明している。
  - lower-level test は story / fixture が外部接続と保存禁止情報を持たない import 境界を証明している。
  - frontend local gate は成功、または task と無関係な環境理由を持つ。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は Storybook dev 表示確認を担当しない。
  - fixture の表示文言全文を snapshot 固定しない。禁止境界と外から見える状態だけを検査する。

### `MPC-ST-03-storybook-review-evidence`

- `implementation_target`: Storybook dev 表示確認と review URL を task-local に残す。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `./ui-design.md`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `./ui-design.md#agent-browser-review`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: Storybook localhost URL、iframe URL、story ID、確認状態
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `storybook-review.md`、screenshot、snapshot、command 記録に secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を出さない。
- `owned_scope`:
  - `docs/exec-plans/active/2026-05-17-master-persona-componentization/storybook-review.md` を作成する。
  - `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` の story ID、review URL、iframe URL、確認状態、未確認理由、再実行 command を記録する。
  - Storybook dev server の確認結果と `build-storybook` 結果を記録する。
  - 必要な screenshot または snapshot は task-local の evidence path に保存する。
  - review URL は Storybook localhost または iframe URL だけを指す。
  - fakeAPI URL、Wails runtime URL、backend API URL を Storybook review URL として扱わない。
- `estimated_change_size`:
  - `files`: `1-6 files`
  - `changed_lines`: `40-140 lines`
  - `size_class`: `通常`
- `depends_on`: `MPC-FE-01-master-persona-panel-stories`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `MPC-UT-02-component-boundary-tests`
- `parallel_blockers`: `なし`
- `first_action`: `docs/exec-plans/active/2026-05-17-master-persona-componentization/storybook-review.md` の review target skeleton を作成し、`SCN-MPC-007` の review record format clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run storybook`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - Storybook dev server が起動し、対象 story の review URL と iframe URL が記録されている。
  - 4 つの必須 story group の story ID と確認状態が `storybook-review.md` にある。
  - 未確認 story がある場合は、未確認理由と再実行 command がある。
  - `npm --prefix frontend run build-storybook` の結果が Storybook 専用 gate として記録されている。
  - review 記録と evidence は禁止情報を含まない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run storybook` は long-running command になる可能性がある。実装者は URL 記録後に process を停止する。
  - `agent-browser` が環境で使えない場合は、未確認理由と代替確認の範囲を `storybook-review.md` に残す。

## Final Validation

全 handoff 完了後に `implement_lane` が最終検証として扱う。
最終検証は handoff ではなく、完了判定である。

- `npm --prefix frontend run lint`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

最終検証の完了条件は次の通りである。

- `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` の story が Storybook registry に存在する。
- `MasterPersonaPage` は production controller 接続と props 合成へ薄くなっている。
- screen local component は small props と callback だけで story 表示できる。
- story / fixture は外部接続、secret、実ユーザーデータ、raw request、raw response、raw prompt を含まない。
- `storybook-review.md` は review URL、iframe URL、story ID、確認状態、未確認理由、再実行 command、Storybook build 結果を持つ。

## Non Targets

- backend 実装は扱わない。
- 統合境界実装は扱わない。
- Wails binding、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB は変更しない。
- docs 正本、`.codex/`、skill、agent、`plan.md`、`scenario-design.md`、`ui-design.md` は変更しない。
- `AIModelSelectionCard` の内部 UI 分割は扱わない。
- page shell、shared list row、shared stat card は新規作成しない。
- Storybook 上で実 AI 生成、実ファイル読み込み、DB 書き込み、provider network を再現しない。
- screen design diff は作成しない。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `./storybook-review.md`
- `final_validation_result`
- `ux_review_result`: `required-before-frontend-human-review`
- `frontend_human_review`: `required-after-frontend-implementation`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`

## Missing Information

- `none`

## Residual Risks

- 既存 screen local component は `viewModel` 全体を受け取っているため、props 型の切り分けで `MasterPersonaPage` の合成処理が一時的に増える可能性がある。
- `GenerationSetupPanel` と `AIModelSelectionCard` の二重枠、長い model 名、長い plugin 名、長い persona 本文は Storybook review で見た目確認が必要である。
- Storybook dev 表示確認は環境に依存する。`agent-browser` が使えない場合は、未確認理由と代替確認の範囲を `storybook-review.md` に残す必要がある。
- page 合成 story は任意である。主要表示領域 story の不足を page 合成 story で代替してはならない。
