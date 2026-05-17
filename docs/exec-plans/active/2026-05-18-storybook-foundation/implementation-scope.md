# Implementation Scope: 2026-05-18-storybook-foundation

- `skill`: implementation-scope
- `status`: handoff-ready
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: human approved `scenario-design.md` and `design-diff-storybook-foundation.puml` in this turn
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `N/A`
- `screen_design_diff`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`
- `design_diff`: `./design-diff-storybook-foundation.puml`

## Fixed Decisions

- `needs_human_decision`: `0`
- UI 設計は `N/A` とする。画面再設計はこの task の対象外である。
- Storybook 基盤は frontend tooling foundation として扱う。backend 実装は原則不要である。
- 統合境界実装は原則不要である。backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB は変更しない接続先である。
- Storybook build は `npm --prefix frontend run build-storybook` を Storybook 専用 gate として扱う。
- 既存 lint の公開検査入口は `npm --prefix frontend run lint` とする。
- Storybook fixture は component 横の `frontend/src/ui/**/__fixtures__` に置く。
- Storybook review URL と確認状態は task-local の `storybook-review.md` に記録する。
- docs 正本化、`.codex/` 変更、プロダクト backend 変更は implementation handoff に含めない。
- `secret_boundary`: `not_required`。ただし review URL、story ID、検証記録には secret、API key、token、ローカル絶対 path、実ユーザーデータを含めない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `SBF-FE-01-storybook-foundation` | `なし` | `なし` | `なし` |
| `wave-2` | `SBF-UT-02-lint-storybook-boundary`, `SBF-ST-03-storybook-review-evidence` | `SBF-FE-01-storybook-foundation` | `SBF-UT-02-lint-storybook-boundary <-> SBF-ST-03-storybook-review-evidence` | `なし` |

## Dependency Table

| handoff_id | depends_on | reason |
| --- | --- | --- |
| `SBF-FE-01-storybook-foundation` | `なし` | Storybook scripts、config、最小 story、fixture 方針が後続検査の前提である。 |
| `SBF-UT-02-lint-storybook-boundary` | `SBF-FE-01-storybook-foundation` | Storybook 専用 file pattern と fixture path が確定した後に、許可対象と禁止対象を lint に固定する。 |
| `SBF-ST-03-storybook-review-evidence` | `SBF-FE-01-storybook-foundation` | dev URL と story ID の記録は Storybook dev 入口と最小 story が存在した後に行う。 |

## Handoffs

### `SBF-FE-01-storybook-foundation`

- `implementation_target`: Storybook package、scripts、config、最小 story、fixture 方針を追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: review URL、story ID、fixture、command 記録に secret、API key、token、実ユーザーデータを出さない。
- `owned_scope`:
  - `frontend/package.json` に Storybook dev script と build script を追加する。
  - `frontend/package-lock.json` は Storybook devDependencies 追加に伴う lockfile としてだけ更新する。
  - `frontend/.storybook/` に Svelte 5、TypeScript、Vite 前提の最小 config を置く。
  - `frontend/src/ui/**` に最小 sample story を置く。
  - `frontend/src/ui/**/__fixtures__` に fixed props または view model fixture の最小例を置く。
  - story と fixture は generated `wailsjs`、backend DTO、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、filesystem business flow を import しない。
- `estimated_change_size`:
  - `files`: `6-8 files`。lockfile 1 件を含む。
  - `changed_lines`: `120-220 lines`。lockfile 差分は行数見積もりから除外する。
  - `size_class`: `通常`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `frontend/package.json` の `scripts.storybook` と `scripts.build-storybook` を追加して、SCN-SBF-001 と SCN-SBF-002 の script entry clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - Storybook dev script と build script が `frontend/` package root から実行できる。
  - `npm --prefix frontend run build-storybook` が backend、Wails runtime、AI provider、secret store、DB なしで成功する。
  - 最小 story は fixed props または view model fixture だけで表示できる。
  - fixture は `frontend/src/ui/**/__fixtures__` に置かれている。
  - story と fixture は backend mock、backend DTO mock、gateway mock、翻訳実行フロー再現を含まない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run build-storybook` はこの handoff の直接 gate と最終検証の Storybook 専用 gate の両方で実行する。
  - Storybook static output は frontend tooling の成果物であり、既存 app build の `dist` と混同しない。

### `SBF-UT-02-lint-storybook-boundary`

- `implementation_target`: 既存 lint 入口に Storybook 依存混入チェックを追加する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: lint error、test fixture、test case 名に secret、API key、token、実ユーザーデータを出さない。
- `owned_scope`:
  - `scripts/eslint/repository-boundary-plugin.mjs` または既存 lint 境界 rule に Storybook 禁止 import 判定を追加する。
  - `frontend/repository-boundary-plugin.test.mjs` または同等の lower-level test に、許可対象と禁止対象を追加する。
  - Storybook 専用設定、story、fixture は Storybook package、Storybook runtime、Storybook 専用 module を import できる。
  - プロダクトコードは Storybook package、Storybook runtime、Storybook 専用 module を import できない。
  - `npm --prefix frontend run lint` で依存混入チェックが実行される状態にする。
- `estimated_change_size`:
  - `files`: `2-4 files`
  - `changed_lines`: `80-180 lines`
  - `size_class`: `通常`
- `depends_on`: `SBF-FE-01-storybook-foundation`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `SBF-ST-03-storybook-review-evidence`
- `parallel_blockers`: `なし`
- `first_action`: `scripts/eslint/repository-boundary-plugin.mjs` の import target 判定に Storybook 専用 module の禁止検出を追加して、SCN-SBF-007 の forbidden import clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run lint:boundaries`
  - `npm --prefix frontend run lint`
- `completion_signal`:
  - lower-level test は Storybook 専用 file pattern の許可を証明している。
  - lower-level test はプロダクトコードから Storybook 専用 module への import 禁止を証明している。
  - `npm --prefix frontend run lint` は Storybook 依存混入チェックを含む既存 lint 入口として成功する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は Storybook build gate を担当しない。
  - lint 対象は Storybook 専用設定、story、fixture の許可範囲を誤って止めない。

### `SBF-ST-03-storybook-review-evidence`

- `implementation_target`: Storybook dev 表示確認と review URL 記録を task-local に残す。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `screen_design_diff`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `storybook-review.md` の URL、query、story ID、command 記録に secret、API key、token、ローカル絶対 path、実ユーザーデータを出さない。
- `owned_scope`:
  - `docs/exec-plans/active/2026-05-18-storybook-foundation/storybook-review.md` を作成する。
  - Storybook dev server の URL、story ID、確認状態、未確認理由、起動 command、Storybook build 結果を記録する。
  - review URL は Storybook 表示確認先だけを指す。
  - fakeAPI URL、Wails runtime URL、backend API URL を Storybook review URL として扱わない。
  - command 出力全文、依存 cache の絶対 path、secret、token、長いローカル path は保存しない。
- `estimated_change_size`:
  - `files`: `1 file`
  - `changed_lines`: `30-70 lines`
  - `size_class`: `通常`
- `depends_on`: `SBF-FE-01-storybook-foundation`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `SBF-UT-02-lint-storybook-boundary`
- `parallel_blockers`: `なし`
- `first_action`: `docs/exec-plans/active/2026-05-18-storybook-foundation/storybook-review.md` の記録 skeleton を作成して、SCN-SBF-004 の review record format clause を閉じる。
- `validation_commands`:
  - `npm --prefix frontend run storybook`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - Storybook dev server が起動し、最小 story の表示 URL と story ID が記録されている。
  - Storybook static build の実行結果が Storybook 専用 gate として記録されている。
  - `storybook-review.md` は確認済み状態または未確認理由を持つ。
  - review URL と記録内容は secret、API key、token、ローカル絶対 path、実ユーザーデータを含まない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `notes`:
  - `npm --prefix frontend run storybook` は long-running command になる可能性がある。実装者は起動 URL を記録した後に process を停止する。
  - 実画面確認が環境で止まる場合は、未確認理由と再実行 command を `storybook-review.md` に残す。

## Final Validation

全 handoff 完了後に `implement_lane` が最終検証として扱う。

- `npm --prefix frontend run lint`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

最終検証の完了条件は次の通りである。

- `npm --prefix frontend run lint` は Storybook 依存混入チェックを含む既存 lint 入口として成功する。
- `npm --prefix frontend run build-storybook` は Storybook 専用 gate として成功する。
- `python3 scripts/harness/run.py --suite frontend-local` は既存 frontend local gate として成功する。
- `storybook-review.md` は URL、story ID、確認状態、未確認理由、再実行 command、残留リスクを持つ。
- 最終検証記録には secret、API key、token、ローカル絶対 path、実ユーザーデータを含めない。

## Non Targets

- Master Persona の部品化は扱わない。
- 画面再設計は扱わない。
- backend 実装は扱わない。
- Wails、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB の接続実装は扱わない。
- docs 正本化、`.codex/` 変更、skill 変更、agent 変更は扱わない。
- gateway mock、backend DTO mock、翻訳実行フロー再現は Storybook に入れない。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `storybook-review.md` または `N/A`
- `final_validation_result`
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
