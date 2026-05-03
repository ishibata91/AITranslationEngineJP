# Task Plan: translation-job-setup-phase-provider-settings

- `workflow`: implement-lane
- `status`: completed
- `lane_owner`: `implement_lane`
- `task_id`: `translation-job-setup-phase-provider-settings`
- `task_mode`: new-feature
- `request_summary`: Job Setup 画面で、フェーズごとの provider / model / credential / batch mode を設定できるようにする。model 候補は provider の model list API から取得する。LM Studio は API key 入力を出さない。UI は `docs/UX-standard.md` に従い、横比較ではなく上から順に読める画面構造へ整理する。
- `goal`: Job Setup を master-persona の provider 設定から切り離し、単語翻訳、NPC ペルソナ生成、本文翻訳の各フェーズが独立した実行設定を持つ。あわせて、非エンジニアの利用者が入力確認、共通基盤、フェーズ別 AI 設定、作成前不足、作成実行を縦の判断順で追えるようにする。
- `constraints`: product code と product test は `implement_lane` では変更しない。UI 変更があるため、scenario-design と ui-design の human review 後に implementation-scope を作る。
- `close_conditions`: scenario-design、ui-design、人間設計レビュー、implementation-scope、implementation handoff、最終検証、5 観点 review、work report、completed 移動が完了している。

## Artifact Index

- `scenario_candidates`: `completed`
- `scenario_design`: `human-approved`
- `ui_design`: `human-approved`
- `ui_prototype`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/prototype.svelte`
- `ui_prototype_sample_data`: `data-ui-prototype-sample-data-root only; no mock-data directory`
- `ui_agent_browser_review`: `completed-after-UX-standard-overhaul; desktop 1440x1000 and mobile 390x844; official dev:prototype command on 34116 succeeded`
- `ui_prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `ui_prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task translation-job-setup-phase-provider-settings --port 34116`
- `implementation_scope`: `created`
- `detail_spec_target`: `docs/detail-specs/translation-job-setup.md` または `docs/detail-specs/<phase>.md`

## Routing Notes

- `required_reading`:
  - `docs/exec-plans/completed/translation-job-setup/plan.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md`
  - `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `internal/service/translation_job_setup_service.go`
  - `internal/usecase/translation_job_setup_contract.go`
  - `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `canonicalization_targets`:
  - Job Setup の恒久仕様は human 承認後に detail-spec へ反映判断する。
  - phase 実行時の provider 設定参照は各 phase detail-spec と整合させる。
- `detail_spec_upper_scenario_id`: `translation-job-setup`
- `validation_commands`:
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.candidate-coverage.json --json`
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`

## Confirmed Facts

- `internal/service/translation_job_setup_service.go:1197` は Job Setup が `MasterPersonaAISettingsRecord` を読む実装である。
- `internal/service/translation_job_setup_service.go:1210` は runtime option を保存済み master-persona provider / model から 1 件だけ作る。
- `internal/service/translation_job_setup_service.go:1237` と `internal/service/translation_job_setup_service.go:1308` は Job Setup の secret key を `master-persona:<provider>` として解決する。
- `internal/usecase/translation_job_setup_contract.go:123` と `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts:26` は runtime option が単一 provider / model / mode だけを表す。
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:247` は provider / model / execution mode を 1 つの select にまとめている。
- `implementation-scope.md` は承認済み scenario-design と ui-design から作成済みであり、backend、統合境界、frontend、テストの実装引き継ぎ単位を分けている。

## Design Requirements

- Job Setup は master-persona の provider 設定を既定値または保存元として扱わない。
- Job Setup は phase ごとに provider、model、credential 参照、execution mode を持つ。
- model 候補は provider ごとの model list API から取得する。API key が設定済みの場合だけ外部取得を試みる。
- LM Studio は API key を要求しないため、API key 入力、API key 未設定 warning、credential select に出さない。
- batch mode は暗黙推定にしない。対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
- UI は `docs/UX-standard.md` の情報階層、スクロール順、セクション構造に従い、3 フェーズを横並び比較ではなく縦順の設定ブロックとして表示する。

## DAG Status

- `task 枠`: completed
- `scenario_candidates`: completed
- `シナリオ設計`: human-approved
- `UI設計`: human-approved
- `人間設計レビュー`: completed
- `実装範囲`: completed
- `実装引き継ぎ入力`: completed-for-wave-1
- `contract_freeze`: completed
- `backend 実装`: completed
- `frontend 実装`: completed
- `テスト実装`: completed-wave-2-api-acceptance
- `統合境界実装`: completed
- `UI / 回帰テスト`: completed
- `最終検証`: completed-after-review-fix
- `レビュー通過根拠`: completed
- `レビュー後修正`: completed
- `作業レポート入力`: completed
- `作業計画完了移動`: completed

## Implementation Scope

- `artifact`: [`implementation-scope.md`](./implementation-scope.md)
- `status`: `ready-for-implement-lane`
- `human_review_status`: `approved`
- `start_wave`: `wave-1`
- `ready_waves`:
  - `wave-1`: `provider-settings-contract-freeze`
  - `wave-2`: `backend-provider-settings-core`, `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance`
  - `wave-3`: `integration-provider-settings-boundary`
  - `wave-4`: `tests-provider-settings-ui-and-regression`
  - `wave-5`: `final-validation-and-report`
- `parallel_execution`: `wave-2` の backend、frontend、API acceptance test は、contract freeze 完了後に並列実行可能である。
- `dependency_order`: contract freeze、backend / frontend / API test、integration boundary、UI / regression test、final validation の順である。

## Spawn Packet

### scenario candidate generators

- `context_policy`: `fork_context=false`
- `task`: `translation-job-setup-phase-provider-settings` の scenario 候補を 6 観点で作成する。
- `output_files`:
  - `scenario-candidates.actor-goal.md`
  - `scenario-candidates.lifecycle.md`
  - `scenario-candidates.state-transition.md`
  - `scenario-candidates.failure.md`
  - `scenario-candidates.external-integration.md`
  - `scenario-candidates.operation-audit.md`
- `must_include`: phase 別 provider 設定、model list API、secret 非露出、LM Studio の API key 非表示、Gemini / xAI batch mode 明示切替、遅延 model list 結果の混入防止。
- `forbidden`: final scenario matrix の確定、採否決定、product code、product test、docs 正本、他 generator の spawn。

### designer

- `context_policy`: `fork_context=false`
- `task`: 6 件の scenario candidates から scenario-design と ui-design を作成する。
- `required_artifacts`: `scenario-design.md`, `scenario-design.requirement-coverage.json`, `scenario-design.candidate-coverage.json`, `scenario-design.questions.md`
- `conditional_artifacts`: UI 変更があるため `ui-design.md` は必須。
- `must_include`: provider 設定の保存単位、phase 実行時の参照境界、model list 取得失敗時の UI 状態、API key 未設定 provider の候補表示可否、batch mode の対象 provider。
- `forbidden`: product code、product test、docs 正本、implementation-scope。

## HITL Status

- `functional_or_design_hitl`: `completed`
- `approval_record`: `Q-TJSPPS-001 answered: option 1; human UI reviewback on 2026-05-04 reflected; UX-standard layout reviewback expanded plan scope; design review approved by human on 2026-05-04`
- `ui_prototype_server_during_review`: `completed`
- `ui_prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `ui_prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task translation-job-setup-phase-provider-settings --port 34116`
- `open_question`: `none`
- `implementation_scope_review`: `not_required; scenario-design.md と ui-design.md は 2026-05-04 に人間承認済みであり、implementation-scope は承認済み成果物から作成済み`

## Codex Implementation Result

- `completed_handoffs`: `provider-settings-contract-freeze`, `backend-provider-settings-core`, `frontend-provider-settings-ui`, `tests-provider-settings-api-acceptance`, `integration-provider-settings-boundary`, `tests-provider-settings-ui-and-regression`
- `touched_files`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/plan.md`, `scenario-candidates.*.md`, `scenario-design.md`, `scenario-design.requirement-coverage.json`, `scenario-design.candidate-coverage.json`, `scenario-design.questions.md`, `ui-design.md`, `prototype.svelte`, `implementation-scope.md`
- `implemented_scope`: `provider-settings-contract-freeze`, `backend-provider-settings-core`, `frontend-provider-settings-ui`
- `test_results`: `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.candidate-coverage.json --json` は pass。`finding_count: 0`、`question_count: 0`。`go test ./internal/usecase ./internal/controller/wails -run 'TranslationJobSetup|JobSetup|ProviderSettings|ModelList'` は pass。`npm --prefix frontend run check` は pass。`go test ./internal/service ./internal/usecase ./internal/repository ./internal/infra/ai -run 'TranslationJobSetup|ProviderSettings|ModelList|PhaseRuntime'` は pass。`npm --prefix frontend run test -- --run translation-job-setup` は pass。`go test ./internal/integrationtest ./internal/service ./internal/infra/ai -run 'TJSPPS|TranslationJobSetup|ProviderSettings|ModelList|Redaction'` は pass。`go test ./internal/controller/wails -run 'TranslationJobSetup|ProviderSettings|ModelList' -count=1 -v` は pass。`python3 scripts/harness/run.py --suite scenario-gate` は pass。`python3 scripts/harness/run.py --suite backend-local` は pass。`python3 scripts/harness/run.py --suite frontend-local` は pass。`python3 scripts/harness/run.py --suite all` は pass。最終再実行では structure、scenario gate、execution、system test 5 passed、frontend coverage statements 65.5% / lines 65.5%、backend coverage statements 69.9% / lines 69.7%、Sonar coverage 71.1%、Sonar security 0、Sonar reliability 0、Sonar maintainability HIGH 0。
- `implementation_investigation`: `initial fact check only`
- `ui_evidence`: `completed; agent-browser snapshot/errors/viewport review and CDP interaction review recorded in ui-design.md`
- `codex_review_result`: `pass; behavior no_issue, contract no_issue, trust-boundary no_issue, state-invariant no_issue, responsibility-boundary no_issue`
- `implementation_action`: `close`
- `sonar_gate_result`: `pass; coverage 71.1%; security 0; reliability 0; maintainability_high 0`
- `residual_risks`: UI は人間指摘により後続デザインオーバーホール予定。今回の実装範囲とレビューゲート上の未解決修正必須問題はない。
- `docs_changes`: `none`

## Closeout Notes

- `canonicalized_artifacts`: `none`
- `detail_spec_canonicalization`: `pending-after-review-and-implementation`
- `follow_up`: UI デザイン差分は後続タスクのデザインオーバーホールで扱う。詳細仕様正本反映は human 承認済み docs-only 作業として別途判断する。

## Outcome

- scenario candidates、scenario-design、ui-design を作成し、`Q-TJSPPS-001` と 2026-05-04 の人間 UI 差し戻しを反映した。さらに UX-standard に基づく画面全体の縦順整理を task 目的へ追加した。人間設計レビューは 2026-05-04 に承認済みであり、implementation-scope も作成済みである。implementation handoff はすべて完了済みである。最終検証は `python3 scripts/harness/run.py --suite all` で通過した。
