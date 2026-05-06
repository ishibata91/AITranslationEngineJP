# translation-job-management plan

## 状態

- `task_id`: `translation-job-management`
- `workflow_state`: `completed`
- `lane_owner`: `implement_lane`
- `source_task`: [`tasks/usecases/translation-job-management.yaml`](../../../../tasks/usecases/translation-job-management.yaml)
- `task_mode`: 新規実装
- `human_review_status`: `approved`

## 目的

Completed 以外の翻訳ジョブを一覧し、選択したジョブの表示、再開入口、停止可否、削除可否、再開不可理由を確認できる状態へ進める。

## 必要判定

- `scenario_candidates`: 必要。新規実装レーンの必須成果物である。
- `designer`: completed。`scenario-design.md` と `ui-design.md` は human review 承認済みである。
- `ui-design`: 必要。`related_screens` に `app-shell.md`、`incomplete-job-list.md`、`job-run.md` がある。
- `implementation-scope`: completed。`implementation-scope.md` は `ready-for-implement-lane` である。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/translation-job-management.yaml`](../../../../tasks/usecases/translation-job-management.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/index.md`](../../../index.md)
- [`docs/screen-design/README.md`](../../../screen-design/README.md)
- [`docs/exec-plans/completed/translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md)
- [`docs/exec-plans/completed/translation-job-setup/ui-design.md`](../../completed/translation-job-setup/ui-design.md)

## 成果物DAG

- `task 枠`: completed。正本はこの `plan.md`。
- `scenario_candidates`: completed。6 観点の成果物は作業計画フォルダに揃っている。
- `シナリオ設計`: completed。`scenario gate` は通過済み。human review 待ちである。
- `UI設計`: completed。human review 待ちである。
- `人間設計レビュー`: completed。2026-05-06 に人間が先へ進めてよいと回答した。
- `実装範囲`: completed。`implementation-scope.md` を作成済みである。
- `実装引き継ぎ入力`: completed。wave-1 の frontend handoff を起動済みである。
- `frontend 実装`: completed。`frontend-job-management-ui` を実装済みである。
- `frontend 実装後人間レビュー`: completed。2026-05-06 に人間レビュー完了と回答済みである。
- `backend 実装`: completed。`backend-job-management-core` は実装済みである。
- `シナリオテスト`: partial_completed。旧仕様の同一入力 duplicate 拒否 integration test は、承認済み仕様へ更新済みである。
- `統合境界実装`: completed。`integration-job-management-wails` は実装済みである。
- `シナリオテスト`: completed。`scenario-test-job-management` は実装済みである。
- `単体テスト`: completed。`unit-test-job-management` は実装済みである。
- `最終検証`: completed。契約レビュー修正後の全体検証は通過済みである。
- `レビュー通過根拠`: completed。5 観点レビューはすべて `no_issue` である。
- `作業レポート入力`: completed。`work_history/runs/2026-05-06-translation-job-management-run/` に作成済みである。

## scenario candidate generator 起動入力

- `context_policy`: `fork_context=false`
- `task`: `translation-job-management` の scenario 候補を 6 観点で作成する。
- `output_directory`: `docs/exec-plans/active/translation-job-management/`
- `required_sources`: `tasks/usecases/translation-job-management.yaml`、`docs/spec.md`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md`
- `must_include`: source requirement、viewpoint、candidate scenario id、actor、trigger、expected outcome、observable point、related detail requirement type、adoption hint
- `forbidden`: final scenario matrix の確定、採否決定、product code、product test、docs 正本、他 agent の起動

## scenario candidate 結果

- [`scenario-candidates.actor-goal.md`](./scenario-candidates.actor-goal.md): 8 件。
- [`scenario-candidates.lifecycle.md`](./scenario-candidates.lifecycle.md): 12 件。
- [`scenario-candidates.state-transition.md`](./scenario-candidates.state-transition.md): 8 件。
- [`scenario-candidates.failure.md`](./scenario-candidates.failure.md): 10 件。
- [`scenario-candidates.external-integration.md`](./scenario-candidates.external-integration.md): 7 件。
- [`scenario-candidates.operation-audit.md`](./scenario-candidates.operation-audit.md): 8 件。

## designer 起動入力

- `context_policy`: `fork_context=false`
- `task`: `translation-job-management` の `scenario-design.md`、coverage JSON、質問票、`ui-design.md` を作成する。
- `output_directory`: `docs/exec-plans/active/translation-job-management/`
- `candidate_sources`: `scenario-candidates.actor-goal.md`、`scenario-candidates.lifecycle.md`、`scenario-candidates.state-transition.md`、`scenario-candidates.failure.md`、`scenario-candidates.external-integration.md`、`scenario-candidates.operation-audit.md`
- `required_sources`: `tasks/usecases/translation-job-management.yaml`、`docs/spec.md`、`docs/index.md`、`docs/screen-design/README.md`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md`、`docs/exec-plans/completed/translation-job-setup/ui-design.md`
- `forbidden`: product code、product test、docs 正本、implementation-scope 作成、候補生成 agent の起動

## 停止条件

- `scenario-candidates.*.md` が 6 件揃わない場合は `designer` へ進めない。
- `scenario-design` の `needs_human_decision` が 1 件以上なら human 質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope.md` を作らない。
- UI がある task のため、frontend 実装後人間レビュー承認前に backend 実装と統合境界実装へ進めない。

## HITL 状態

- `functional_or_design_hitl`: `required-after-design-bundle`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approval_record`: [`scenario-design.questions.md`](./scenario-design.questions.md)

## design bundle 停止結果

- [`scenario-design.md`](./scenario-design.md): human approved。
- [`scenario-design.questions.md`](./scenario-design.questions.md): 未回答質問なし。Q1 から Q4 は回答済み判断として固定済み。Q3 と Q4 は後続タスクへ送る回答済み事項として保持する。
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json): 未解決 conflict 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): `needs_human_decision` 0 件。
- [`ui-design.md`](./ui-design.md): human approved。実画面確認は未実装画面のため未実行。

## 人間設計レビュー結果

- `review_result`: approved。
- `approved_at`: 2026-05-06。
- `approval_note`: 人間が「先進めていい」と回答したため、`scenario-design.md` と `ui-design.md` を implementation-scope 作成の根拠にする。

## implementation-scope 結果

- [`implementation-scope.md`](./implementation-scope.md): `ready-for-implement-lane`。
- `handoff_count`: 5。
- `wave-1`: `frontend-job-management-ui`。
- `wave-2`: `backend-job-management-core`。`frontend 実装後人間レビュー approved` 後に開始する。
- `wave-3`: `integration-job-management-wails`。
- `wave-4`: `scenario-test-job-management` と `unit-test-job-management`。
- `Q4 residual`: `TJM-RES-Q4-translation-execution-stop-control`。翻訳実行側後続 task 化待ちである。

## frontend 実装結果

- `handoff`: `frontend-job-management-ui`。
- `implementation_status`: completed。
- `changed_scope`: `frontend/src/application/gateway-contract/translation-job-management/`、`frontend/src/application/contract/translation-job-management/`、`frontend/src/application/store/translation-job-management/`、`frontend/src/application/presenter/translation-job-management/`、`frontend/src/application/usecase/translation-job-management/`、`frontend/src/controller/translation-job-management/`、`frontend/src/ui/screens/translation-job-management/`、`frontend/src/ui/screens/job-run/JobRunPage.svelte`、`frontend/src/ui/App.svelte`、`frontend/src/ui/views/AppShell.svelte`、`frontend/src/ui/stores/shell-state.ts`。
- `review_url`: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`。
- `ui_evidence`: agent-browser snapshot で `Job Management` tab、6 件の未完了 job 一覧、Completed 非表示、長い path と plugin 名を確認。`agent-browser errors` は空。
- `unverified_ui`: 一覧行 click 後の detail 更新 snapshot と mobile viewport は未確認である。
- `validation`: `npm --prefix frontend run check` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `targeted_test`: `npm --prefix frontend run test -- --run translation-job-management` は対象 test file がなく `No test files found` のため未通過。

## frontend 実装後人間レビュー結果

- `review_result`: approved。
- `approved_at`: 2026-05-06。
- `review_note`: 人間が「人間レビュー終わり」と回答したため、backend 実装へ進める。
- `review_fix`: Ready job は Job Setup、Ready 以外は Job Run へ送る導線分岐を追加済み。
- `review_validation`: `npm --prefix frontend run check` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。

## backend 実装結果

- `handoff`: `backend-job-management-core`。
- `implementation_status`: completed。
- `changed_scope`: `internal/repository/job_lifecycle_*`、`internal/service/translation_job_setup_service.go`、`internal/service/translation_job_management_service.go`、`internal/usecase/translation_job_management_*`、`internal/controller/wails/translation_job_management_controller.go`、`internal/controller/wails/app_controller.go`、`internal/bootstrap/app_controller.go`、`internal/infra/sqlite/**/012_translation_job_multi_input_index.sql`。
- `validation`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'TranslationJobManagement|JobManagement|TranslationJob|JobLifecycle'` pass。
- `backend_local`: `python3 scripts/harness/run.py --suite backend-local` fail。失敗箇所は `internal/integrationtest` の `TestSCN_TJS_CreateTranslationJobRejectsDuplicateInput` である。
- `next_test_action`: 承認済み仕様は「同じ入力から複数 job を作成できる」であるため、旧仕様 test を `scenario-test-job-management` で更新する。

## シナリオテスト部分更新結果

- `handoff`: `scenario-test-job-management` の旧仕様 test 更新。
- `implementation_status`: partial_completed。
- `changed_scope`: `internal/integrationtest/sqlite_integration_test.go`。
- `scenario_proved`: `SCN-TJM-002`。同じ `X_EDIT_EXTRACTED_DATA` へ 2 回 job 作成でき、job 件数が 2 になる。
- `targeted_validation`: `go test ./internal/integrationtest -run TestSCN_TJM_002_CreateTranslationJobAllowsDuplicateInput -count=1` pass。
- `backend_local`: `python3 scripts/harness/run.py --suite backend-local` pass。

## 統合境界停止結果

- `handoff`: `integration-job-management-wails`。
- `implementation_status`: stopped。
- `stop_reason`: production 注入に必要な `frontend/src/bootstrap/app-screen-controller-factories.ts`、`frontend/src/main.ts`、`frontend/src/ui/App.svelte` が統合境界 owned scope から漏れていた。
- `scope_correction`: `implementation-scope.md` の `integration-job-management-wails` に production 注入ファイルを追加済みである。
- `next_action`: 補正済み handoff で `integration-job-management-wails` を再起動する。

## 統合境界実装結果

- `handoff`: `integration-job-management-wails`。
- `implementation_status`: completed。
- `changed_scope`: `frontend/src/controller/wails/translation-job-management.gateway.ts`、`frontend/src/controller/wails/gateway-dto/translation-job-management/*`、`frontend/src/bootstrap/app-screen-controller-factories.ts`、`frontend/src/main.ts`、`frontend/src/ui/App.svelte`。
- `integration_boundary`: Wails binding、frontend gateway DTO、production / fakeApi controller 注入を接続済みである。
- `ui_evidence`: production empty 状態と fake 成功状態の未完了 job 一覧、状態フィルタ、無効理由を agent-browser で確認済みである。
- `validation`: `go test ./internal/bootstrap ./internal/controller/wails -run 'TranslationJobManagement|JobManagement'` pass。`npm --prefix frontend run check` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `targeted_test`: `npm --prefix frontend run test -- --run translation-job-management` は対象 test file がなく `No test files found` のため未通過である。
- `unverified`: production DB に未完了 job がないため、実データでの一覧選択、削除確認、削除結果、再読込後反映は未確認である。

## シナリオテスト結果

- `handoff`: `scenario-test-job-management`。
- `implementation_status`: completed。
- `changed_scope`: `internal/integrationtest/translation_job_management_scenario_test.go`、`tests/system/translation-job-management.spec.ts`。
- `scenario_proved`: `SCN-TJM-001` から `SCN-TJM-009` を API テストと UI E2E で証明済みである。
- `validation`: `go test ./internal/integrationtest -run 'SCN_TJM|TranslationJobManagement'` pass。`npx playwright test --grep 'SCN-TJM|translation-job-management'` pass。`python3 scripts/harness/run.py --suite backend-local` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `residual_risk`: UI E2E の配置は既存 `playwright.config.ts` の `testDir` に合わせて `tests/system` である。

## 単体テスト結果

- `handoff`: `unit-test-job-management`。
- `implementation_status`: completed。
- `changed_scope`: `internal/service/translation_job_management_service_test.go`、`internal/usecase/translation_job_management_usecase_test.go`、`internal/controller/wails/translation_job_management_controller_unit_test.go`、`frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.test.ts`、`frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.test.ts`、`frontend/src/controller/translation-job-management/translation-job-management-screen-controller.test.ts`。
- `unit_proved`: Completed 除外、Running 削除拒否、非実行中削除、reason category、redaction 境界、presenter 表示、controller / usecase DTO を証明済みである。
- `validation`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'TranslationJobManagement|JobManagement|DeleteGuard|Redaction'` pass。`npm --prefix frontend run test -- --run translation-job-management` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `backend_local_note`: 単体テスト agent 実行時点では並列中の `internal/integrationtest/translation_job_management_scenario_test.go` の shadow 警告で fail したが、シナリオテスト agent 修正後に `backend-local` は pass 済みである。
- `residual_risk`: `phase_progress_aggregation_failed`、`terminal_state`、`state_projection_inconsistent` の全分岐 unit 網羅と repository 層の同一入力複数 job 直接 unit 証明は未追加である。

## 最終検証結果

- `final_validation_status`: failed。
- `pass`: `requirement_gate`、`go test ./internal/...`、`npm --prefix frontend run check`、`npm --prefix frontend run test -- --run`、`scenario-gate`、`backend-local`、`frontend-local`、`npx playwright test --grep 'SCN-TJM|translation-job-management'`。
- `all_harness`: `python3 scripts/harness/run.py --suite all` fail。
- `failure`: Sonar coverage summary coverage `69.4%` が閾値 `70.0%` 未満である。line `70.4%`、branch `61.9%`、security / reliability / maintainability high issue は 0 件である。
- `next_action`: `unit-test-job-management` に coverage gap を閉じる追加単体テストを依頼する。

## coverage 修正結果

- `handoff`: `unit-test-job-management` の coverage gap 修正。
- `implementation_status`: completed。
- `changed_scope`: `frontend/src/application/store/translation-job-management/translation-job-management.store.test.ts`、`frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.test.ts`、`frontend/src/controller/wails/translation-job-management.gateway.test.ts`。
- `additional_unit_proved`: store の clone / subscribe、usecase の未接続 / stale / stop_failed / delete guard / resume success / filter、Wails gateway の controller 解決 / binding 未接続 / request mapping。
- `validation`: `npm --prefix frontend run test -- --run translation-job-management` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。`python3 scripts/harness/run.py --suite coverage` pass。
- `coverage_gate`: Sonar coverage summary `70.2%`、line `71.0%`、branch `63.6%`。security / reliability / maintainability high issue は 0 件である。

## 最終検証再実行結果

- `final_validation_status`: passed。
- `all_harness`: `python3 scripts/harness/run.py --suite all` pass。
- `test_counts`: frontend test は 56 files / 476 tests、system test は 9 tests である。
- `coverage_gate`: Sonar coverage summary `70.2%`、line `71.0%`、branch `63.6%`。閾値 `70.0%` を通過済みである。
- `issue_gate`: Sonar security 0、reliability 0、maintainability HIGH 0 である。
- `coverage_manifest`: `test-results/coverage-manifest.json` に出力済みである。

## レビュー集約結果

- `implementation_action`: `fix`。
- `review_trust_boundary`: no_issue。`hard_gate: true`、`must_fix_open: false` である。
- `review_behavior`: issues_open。`behavior-001` は選択詳細と再開不可理由が画面に描画されない major 指摘である。
- `review_contract`: issues_open。`contract-001` は backend 公開 reason category と frontend 型契約の不一致 major 指摘である。
- `review_state_invariant`: issues_open。`state-invariant-001` は phase 実行中の状態不整合 job を削除できる critical 指摘である。
- `review_responsibility_boundary`: issues_open。`responsibility-boundary-001` は frontend 型配置の責務境界 major 指摘である。`responsibility-boundary-002` は承認済み範囲外差分の混入指摘であり、プロダクト修正ではなく再レビュー時の対象差分指定で扱う。

## レビュー修正 wave-1

- `backend_fix`: completed。`state-invariant-001` を修正済みである。phase 実行中または状態投影不整合の削除拒否を action 経路と永続化 guard に揃えた。
- `frontend_fix`: completed。`behavior-001`、`contract-001`、`responsibility-boundary-001` を修正済みである。選択詳細表示、reason category 追加、gateway contract と screen view model の型配置分離を行った。
- `parallel_policy`: backend と frontend の所有ファイルは分かれるため並列実行できる。
- `review_scope_note`: `.codex`、`scripts/scenario`、`tasks` の差分は別レーン由来または人間指示由来として扱い、Job Management product 修正の対象差分から除外して再レビューへ渡す。

## backend レビュー修正結果

- `fix_target`: `state-invariant-001`。
- `changed_scope`: `internal/repository/job_lifecycle_repository.go`、`internal/repository/job_lifecycle_sqlite_repository.go`、`internal/service/translation_job_management_service.go`、`internal/service/translation_job_management_service_test.go`、`internal/integrationtest/translation_job_management_scenario_test.go`。
- `fixed_behavior`: `TRANSLATION_JOB.state` が Running 以外でも、phase run が running / pending または stop_requested 相当なら削除を拒否する。
- `validation`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'TranslationJobManagement|JobManagement|TranslationJob|JobLifecycle'` pass。`go test ./internal/integrationtest -run 'SCN_TJM|TranslationJobManagement' -count=1` pass。`python3 scripts/harness/run.py --suite backend-local` は初回 lint fail 後に修正し pass。
- `residual_risk`: `JOB_PHASE_RUN.latest_error = stop_requested` の専用 service / integration test は未追加である。

## frontend レビュー修正結果

- `fix_target`: `behavior-001`、`contract-001`、`responsibility-boundary-001`。
- `changed_scope`: `frontend/src/application/contract/translation-job-management/**`、`frontend/src/application/gateway-contract/translation-job-management/**`、`frontend/src/application/presenter/translation-job-management/**`、`frontend/src/application/store/translation-job-management/**`、`frontend/src/application/usecase/translation-job-management/**`、`frontend/src/controller/translation-job-management/**`、`frontend/src/ui/screens/translation-job-management/**`、`frontend/src/ui/screens/job-run/JobRunPage.svelte`、`frontend/src/ui/views/AppShell.svelte`。
- `fixed_behavior`: 選択 job 詳細に入力出自、現在フェーズ、進捗、入力キャッシュ状態、AI 設定要約、操作可否、再開不可理由、削除不可理由を表示する。
- `fixed_contract`: frontend reason category に `phase_progress_aggregation_failed` を追加し、表示処理へ通した。
- `fixed_boundary`: gateway request / response は gateway contract に残し、screen state と view model は screen contract 側へ分離した。
- `validation`: `npm --prefix frontend run check` pass。`npm --prefix frontend run test -- --run translation-job-management` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `residual_risk`: なし。

## レビュー修正後 全体検証結果

- `final_validation_status`: failed。
- `all_harness`: `python3 scripts/harness/run.py --suite all` fail。
- `pass`: structure、scenario requirement gate、backend lint、frontend lint、backend test、frontend test、Sonar scan、system test 9 件、frontend coverage、backend coverage、Sonar coverage `70.2%`、line `71.1%`、branch `63.4%`、security 0、reliability 0。
- `failure`: Sonar maintainability HIGH issue が 2 件である。`internal/service/translation_job_management_service.go:199` の cognitive complexity と、同ファイル `:214` の重複文言である。
- `next_action`: backend 修正 agent で `DeleteJob` の分割と重複文言定数化を行い、`backend-local` と `all` を再実行する。

## frontend human feedback

- `feedback_status`: correction_required。
- `feedback`: 人間レビュー済み UI では、詳細パネルや一覧の過剰な表示項目は不要として削除済みである。
- `wrong_fix`: `behavior-001` の修正で、合意済み UI に反して選択 job 詳細パネルと一覧の追加表示項目を復活させた。
- `correction_policy`: stepper と Job Run 連携は人間が指摘して作らせた合意済み UI として維持する。詳細パネルと一覧の過剰表示は削除する。
- `review_policy`: 挙動レビュー指摘は、人間レビュー済み UI 判断で上書きする。再開不可理由と操作可否は、一覧内の操作ボタン、無効理由、Job Run 表示対象で確認できる範囲へ収める。

## human bug report

- `reported_at`: 2026-05-06。
- `report`: データロード画面で同じ JSON を登録すると重複 input と表示され登録できない。
- `expected`: 同じファイルまたは同じ入力データから複数の翻訳 job を作成できる方針に合わせ、入力登録と job 作成の重複拒否を外す。
- `report`: 既存 job があるはずだが Job Management 画面に job が表示されていない。
- `observed_cause_candidate`: `TranslationInputImportService` に `source_content_hash` 重複拒否が残っている。`translation-job-setup.usecase.ts` に既存 job がある入力で作成を拒否する条件が残っている。
- `fix_action`: backend / frontend に分けて、入力重複許可、Job Setup 既存 job block 削除、Job Management 一覧表示の検証を行う。

## data load to job setup policy

- `decision`: Data Load の新規登録は input 作成だけであり、job は自動作成しない。
- `decision`: Job Management は `TRANSLATION_JOB` の Completed 以外だけを表示する。job 未作成 input は混ぜない。
- `decision`: 登録成功後は Data Load に `Job Setup へ進む` ボタンを表示し、利用者が登録結果を確認してから job 作成へ進む。
- `implementation_action`: `InputReviewPage` に Job Setup 遷移 callback を追加し、`AppShell` から `openTranslationJobSetup` を渡す。
- `test_action`: Data Load 登録成功後に Job Setup ボタンが表示され、押すと Job Setup view へ移ることを frontend test で証明する。

## Data Load 導線修正結果

- `implementation_status`: completed。
- `changed_scope`: `frontend/src/ui/screens/translation-input/LoadedInputDetail.svelte`、`frontend/src/ui/screens/translation-input/InputReviewPage.svelte`、`frontend/src/ui/views/AppShell.svelte`、`frontend/src/application/presenter/translation-input/translation-input.presenter.ts`、`frontend/src/application/presenter/translation-input/index.ts`、`frontend/src/ui/screens/translation-input/InputReviewPage.test.ts`、`frontend/src/application/presenter/translation-input/translation-input.presenter.test.ts`。
- `fixed_behavior`: 登録済みまたは警告ありの selected input だけ `Job Setup へ進む` ボタンを表示する。失敗または再構築が必要な selected input では表示しない。
- `validation`: `npm --prefix frontend run test -- src/ui/screens/translation-input/InputReviewPage.test.ts src/application/usecase/translation-input/translation-input.usecase.test.ts src/application/presenter/translation-input/translation-input.presenter.test.ts` pass。`npm --prefix frontend run check` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `residual_risk`: なし。

## Data Load 導線後 全体検証結果

- `final_validation_status`: failed。
- `all_harness`: `python3 scripts/harness/run.py --suite all` fail。
- `pass`: structure、scenario requirement gate、backend lint、frontend lint、backend test、frontend test、Sonar scan。
- `failure`: system test `SCN-TJM-001 translation-job-management lists incomplete jobs and excludes completed` が失敗した。`Job #401` は合意済み UI で見出しではなく一覧カード内表示であるため、テスト期待値が旧 UI 構造を参照している。
- `next_action`: `scenario-test-job-management` にシステムテスト期待値を合意済み UI へ合わせる修正を依頼する。

## system test 期待値修正結果

- `implementation_status`: completed。
- `changed_scope`: `tests/system/translation-job-management.spec.ts`。
- `fixed_behavior`: `SCN-TJM-001` は `Job #401` から `Job #406` を一覧カード内表示として確認する。Completed 非表示、状態バッジ、Job Run 遷移の検証は維持する。
- `validation`: `npm run test:system -- tests/system/translation-job-management.spec.ts` pass。`python3 scripts/harness/run.py --suite system-test` pass。
- `residual_risk`: なし。

## Data Load 導線後 全体検証再実行結果

- `final_validation_status`: passed。
- `all_harness`: `python3 scripts/harness/run.py --suite all` pass。
- `test_counts`: frontend test は 57 files / 483 tests、system test は 9 tests である。
- `coverage_gate`: Sonar coverage summary `70.4%`、line `71.3%`、branch `63.8%`。閾値 `70.0%` を通過済みである。
- `issue_gate`: Sonar security 0、reliability 0、maintainability HIGH 0 である。
- `coverage_manifest`: `test-results/coverage-manifest.json` に出力済みである。

## レビュー wave-2 結果

- `implementation_action`: `fix`。
- `review_behavior`: no_issue。`must_fix_open: false`、`max_level: none` である。
- `review_trust_boundary`: no_issue。`hard_gate: true`、`must_fix_open: false`、`max_level: none` である。
- `review_state_invariant`: no_issue。`must_fix_open: false`、`max_level: none` である。
- `review_responsibility_boundary`: no_issue。`must_fix_open: false`、`max_level: none` である。
- `review_contract`: issues_open。`contract-001` は Job Setup presenter が同じ input の `existingJob` を作成ブロック理由にしている major 指摘である。
- `next_action`: frontend 修正 agent で `buildGlobalBlockedReasons` から `existingJob` 由来の作成ブロックを外す。

## 契約レビュー修正 wave-2

- `fix_target`: `contract-001`。
- `implementation_status`: completed。
- `changed_scope`: `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`、`frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts`。
- `fixed_behavior`: 同じ input の `existingJob` は参考表示だけに残し、`globalBlockedReasons` と create 可否へ使わない。
- `validation`: `npm --prefix frontend run test -- src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts` pass。`python3 scripts/harness/run.py --suite frontend-local` pass。
- `residual_risk`: なし。

## 契約レビュー修正後 全体検証結果

- `final_validation_status`: passed。
- `all_harness`: `python3 scripts/harness/run.py --suite all` pass。
- `test_counts`: frontend test は 57 files / 484 tests、system test は 9 tests である。
- `coverage_gate`: Sonar coverage summary `70.5%`、line `71.3%`、branch `64.0%`。閾値 `70.0%` を通過済みである。
- `issue_gate`: Sonar security 0、reliability 0、maintainability HIGH 0 である。
- `coverage_manifest`: `test-results/coverage-manifest.json` に出力済みである。

## レビュー wave-2 再集約結果

- `implementation_action`: `close`。
- `review_behavior`: no_issue。`must_fix_open: false`、`max_level: none` である。
- `review_contract`: no_issue。`contract-001` は解決済みで、`must_fix_open: false`、`max_level: none` である。
- `review_trust_boundary`: no_issue。`hard_gate: true`、`must_fix_open: false`、`max_level: none` である。
- `review_state_invariant`: no_issue。`must_fix_open: false`、`max_level: none` である。
- `review_responsibility_boundary`: no_issue。`must_fix_open: false`、`max_level: none` である。

## 終了処理

- `作業レポート入力`: `work_history/runs/2026-05-06-translation-job-management-run/` に作成済みである。
- `completion_evidence`: この `plan.md`、5 観点 `reviewback.*.yaml`、`test-results/coverage-manifest.json`、検証コマンド結果である。
- `follow_up`: Q3 API current validity check は resume execution task へ後続送りである。Q4 stop / cancel / late response control は translation execution stop task へ後続送りである。
- `close_action`: 作業計画フォルダを `docs/exec-plans/completed/translation-job-management/` へ移す。
