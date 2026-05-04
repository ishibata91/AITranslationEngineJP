# Task Plan: ai-provider-settings-management

- `workflow`: implement-lane
- `status`: implementation-review-passed
- `lane_owner`: `implement_lane`
- `task_id`: `ai-provider-settings-management`
- `task_mode`: new-feature
- `request_summary`: 各プロバイダ設定画面を作り、app-shell からルーティング可能にする。プロバイダ別にエンドポイントと APIキーを設定できるようにし、翻訳フェーズやマスターペルソナ生成とは別の永続仕様として管理する。Q-AIPSM-003 の人間回答により、AIサービス設定画面は APIキーと endpoint だけを設定対象にし、model、処理方法、利用 provider、Batch API 切り替えは各翻訳フェーズと master-persona 側で扱う。
- `goal`: APIキーとエンドポイントをプロバイダ単位の独立設定として保存し、翻訳ジョブ設定やマスターペルソナ生成が個別に secret や endpoint を持たない構造へ寄せる。
- `constraints`: `implement_lane` はプロダクトコード、プロダクトテスト、docs 正本文を変更しない。UI と DB が関係するため、scenario-design と ui-design の人間レビュー後に implementation-scope を作る。
- `close_conditions`: scenario candidates、scenario-design、ui-design、人間設計レビュー、implementation-scope、implementation handoff、最終検証、5 観点 review、work report、completed 移動が完了している。

## Artifact Index

- `scenario_candidates`: `completed`
- `scenario_design`: `approved`
- `ui_design`: `approved`
- `human_design_review`: `approved`
- `implementation_scope`: `completed`
- `implementation_handoff`: `completed-after-review-fix`
- `detail_spec_target`: `pending-after-review-and-implementation`

## Source Facts

- `docs/spec.md:50`: Gemini と xAI を翻訳 AI として利用できる必要がある。
- `docs/spec.md:51`: Gemini と xAI は BatchAPI を利用できる必要がある。
- `docs/spec.md:56`: 共通ペルソナ構築、共通辞書構築、翻訳フロー、各翻訳フェーズなど、目的に沿った AI を選択可能である必要がある。
- `docs/spec.md:58`: 各フェーズでは、ユーザーがプロバイダとモデルを選択できる必要がある。
- `docs/spec.md:59`: 各フェーズの API 選択と APIKey は再入力不要で保存できる必要がある。
- `docs/spec.md:60`: APIKey は暗号化して保存する必要がある。
- `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/plan.md`: Job Setup は master-persona の provider 設定を使わない設計へ変更済みである。
- `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`: phase 別 provider 設定、model list、secret 非露出、batch mode 限定のシナリオが既に存在する。
- `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`: provider id は `gemini`、`lm_studio`、`xai` を実 provider として扱い、fake は test-only DI に置く判断がある。

## Design Requirements

- プロバイダ設定は app-shell から開ける独立画面にする。
- APIキーとエンドポイントの永続仕様は、翻訳フェーズ、翻訳ジョブ設定、マスターペルソナ生成の設定と別に管理する。
- APIキーは UI、DTO、log、エラー要約へ平文表示しない。
- 画面表示、provider 選択、endpoint 保存では APIキー本文を読み出さない。
- APIキー本文の読み出しは、接続確認、モデル一覧取得、翻訳実行開始、マスターペルソナ生成開始に限定する。
- AIサービス設定画面では、APIキーと endpoint だけを設定できる。
- model、処理方法、利用 provider、Batch API 切り替えは、各翻訳フェーズと master-persona 側で設定する。
- DB 変更が必要かは scenario-design と implementation-scope で repository、migration、secret store の責務を分けて確定する。

## DAG Status

- `task 枠`: completed
- `scenario_candidates`: completed
- `シナリオ設計`: approved
- `UI設計`: approved
- `人間設計レビュー`: approved
- `実装範囲`: completed
- `実装引き継ぎ入力`: completed
- `実装前受け入れテスト`: completed-for-contract-freeze
- `contract_freeze`: completed
- `backend 実装`: completed-after-review-fix
- `統合境界実装`: completed-after-review-fix
- `frontend 実装`: completed-wave-2
- `実装後単体テスト`: no-additional-unit-agent-required
- `最終検証`: completed-after-review-fix
- `レビュー通過根拠`: completed
- `正本化判断`: pending-human-approved-canonicalization
- `詳細仕様正本反映`: blocked-by-human-approved-canonicalization
- `作業レポート入力`: blocked-missing-run-evidence
- `作業計画完了移動`: blocked-by-docs-canonicalization-and-work-report

## Spawn Packet

### scenario candidate generators

- `context_policy`: `fork_context=false`
- `task`: `ai-provider-settings-management` の scenario 候補を 6 観点で作成する。
- `output_files`:
  - `scenario-candidates.actor-goal.md`
  - `scenario-candidates.lifecycle.md`
  - `scenario-candidates.state-transition.md`
  - `scenario-candidates.failure.md`
  - `scenario-candidates.external-integration.md`
  - `scenario-candidates.operation-audit.md`
- `must_include`: app-shell routing、provider settings page、endpoint persistence、API key secret persistence、translation phase / master-persona との永続仕様分離、参照側の model / 処理方法 / provider 選択、DB migration candidate、secret redaction、fake provider 非表示。
- `forbidden`: final scenario matrix の確定、採否決定、product code、product test、docs 正本、implementation-scope、他 generator の spawn。

### designer

- `context_policy`: `fork_context=false`
- `task`: 6 件の scenario candidates から scenario-design と ui-design を作成する。
- `required_artifacts`: `scenario-design.md`, `scenario-design.requirement-coverage.json`, `scenario-design.candidate-coverage.json`, `scenario-design.questions.md`
- `conditional_artifacts`: UI 変更があるため `ui-design.md` と task-local UIプロトタイプは必須。
- `must_include`: provider settings の保存単位、APIキー secret store と DB 設定値の境界、endpoint 更新、参照側の model / 処理方法 / provider 選択、app-shell route、既存 Job Setup と master-persona からの参照境界。
- `forbidden`: product code、product test、docs 正本、implementation-scope。

## HITL Status

- `functional_or_design_hitl`: `approved`
- `approval_record`: `human approved UI review and requested progression to implementation-scope`
- `open_question`: `none`
- `implementation_scope_review`: `approved-spec-change-2026-05-04`

## Approved Scope Change: APIキー読み出し制御

- `scope_change_status`: `approved-by-human`
- `reason`: 既存レビュー通過扱いを失効させ、秘密情報を通常表示と通常保存で読まない仕様へ差し替える。
- `allowed_secret_read_operations`: 接続確認、モデル一覧取得、翻訳実行開始、マスターペルソナ生成開始。
- `forbidden_secret_read_operations`: 画面表示、provider 選択、endpoint 保存、APIキー入力なし保存、summary 生成。
- `required_rerun`: backend-local、frontend-local、5 観点レビュー。trust-boundary は hard gate として再判定する。

## Current Stop Line

`scenario-design.questions.md` の未回答質問は 0 件である。
人間設計レビューは承認済みである。
APIキー読み出し制御の仕様変更は人間承認済みである。
APIキー読み出し制御の再レビューは 5 観点すべて通過済みである。
次は作業レポート入力と完了移動可否の確認へ進む。

## Review Aggregation: APIキー読み出し制御

- `implementation_action`: `close`
- `behavior`: `no_issue`, `max_level: none`
- `contract`: `no_issue`, `max_level: none`
- `trust_boundary`: `no_issue`, `max_level: none`, `hard_gate: true`
- `state_invariant`: `no_issue`, `max_level: none`
- `responsibility_boundary`: `no_issue`, `max_level: none`

## Validation Notes

- `designer` agent が `jq empty` を 2 つの JSON で通過確認済み。
- `designer` agent が `requirement_gate.py` 失敗を確認済み。初回理由は未回答質問 6 件と未解決 conflict 6 件であり、設計停止条件どおりである。
- `designer` agent が `agent-browser` で `http://127.0.0.1:34116/prototype` を確認済み。console error はない。
- UI 確認サーバーは停止済み。人間レビュー時は `npm --prefix frontend run dev:prototype -- --task ai-provider-settings-management --port 34116` で再起動する。
- `implement_lane` が `jq empty` を 2 つの JSON で再確認済み。
- `implement_lane` が `python3 scripts/scenario/requirement_gate.py ... --json` を再実行し、`status: fail`、`finding_count: 38`、`question_count: 16` を確認済み。重複展開前の人間向け質問は `scenario-design.questions.md` の 6 件である。
- `implement_lane` が `git diff --check -- docs/exec-plans/active/ai-provider-settings-management` を確認済み。
- 人間回答 `Q-AIPSM-001=3`、`Q-AIPSM-002=1`、`Q-AIPSM-003=4` を反映した。
- Q-AIPSM-003 の反映により、AIサービス設定画面から model と Batch API を外し、参照側設定へ戻した。
- 未回答質問は `Q-AIPSM-004`、`Q-AIPSM-005`、`Q-AIPSM-006` の 3 件に減った。
- `implement_lane` が反映後に `python3 scripts/scenario/requirement_gate.py ... --json` を再実行し、`status: fail`、`finding_count: 17`、`question_count: 9` を確認済み。重複展開前の人間向け質問は 3 件である。
- 人間回答 `Q-AIPSM-004=2`、`Q-AIPSM-005=1`、`Q-AIPSM-006=4` を反映した。
- Q-AIPSM-004 の反映により、Ready job は最新 provider settings を再解決し、Running phase は開始時 snapshot を使う。
- Q-AIPSM-005 の反映により、未設定へ戻す操作は row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。
- Q-AIPSM-006 の反映により、endpoint は表示でき、secret は伏せ字または存在状態だけ表示し、更新履歴は保存しない。
- `implement_lane` が反映後に `jq empty` を 2 つの JSON で再確認済み。
- `implement_lane` が反映後に `python3 scripts/scenario/requirement_gate.py ... --json` を再実行し、`status: pass`、`finding_count: 0`、`question_count: 0` を確認済み。
- `implement_lane` が反映後に `git diff --check -- docs/exec-plans/active/ai-provider-settings-management` を再確認済み。
- UI 人間レビュー用に `npm --prefix frontend run dev:prototype -- --task ai-provider-settings-management --port 34116` を起動中である。
- `agent-browser open http://127.0.0.1:34116/prototype` は通過した。
- `agent-browser snapshot -i --compact --depth 5` で主要 heading、provider list、接続設定、主要ボタンを確認した。
- `agent-browser errors` は出力なしで、console error は確認されなかった。
- 人間が UI レビュー終了と先行指示を出したため、scenario-design と ui-design を承認済みとして扱う。
- UI 確認サーバーは人間レビュー終了に伴い停止済みである。
- `designer` agent が `implementation-scope.md` を作成し、`ready-for-implement-lane` を確認した。
- `provider-settings-contract-freeze` の実装前 API テスト作成を `implementation_scenario_tester` へ起動済みである。
- `implementation_scenario_tester` が `internal/apitest/provider_settings_contract_freeze_test.go` を追加した。
- `implementation_implementer` が `provider-settings-contract-freeze` を実装した。
- `implementation_scenario_tester` が API テストの repo root path 解決を修正し、`go test ./internal/apitest -run 'SCN_AIPSM|ProviderSettings'` と `python3 scripts/harness/run.py --suite backend-local` を通過確認した。
- `provider-settings-contract-freeze` の局所検証として `go test ./internal/usecase ./internal/controller/wails -run 'ProviderSettings|AIProvider|TranslationJobSetup|MasterPersona|PersonaGeneration'`、`npm --prefix frontend run check`、`python3 scripts/harness/run.py --suite frontend-local` が通過済みである。
- `backend-provider-settings-core` は局所検証 `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'ProviderSettings|Secret|Credential|Endpoint|Redaction'` と `go test ./internal/infra/sqlite/dbinit ./internal/integrationtest -run 'ProviderSettings|Migration|SQLite'` が通過済みである。
- `frontend-provider-settings-route-ui` は `npm --prefix frontend run test -- provider-settings AppShell` と実画面 route snapshot を確認済みである。
- `frontend-reference-model-settings-ui` は `npm --prefix frontend run test -- translation-job-setup master-persona provider-settings` を通過済みである。
- `frontend-provider-settings-route-ui` の secret boundary 補正後、`go test ./internal/apitest -run 'SCN_AIPSM|ProviderSettings'`、`npm --prefix frontend run check`、`python3 scripts/harness/run.py --suite frontend-local`、`python3 scripts/harness/run.py --suite backend-local` は通過済みである。
- `backend-provider-settings-consumer-boundary` の残差修正後、`go test -timeout 60s ./internal/...` と `python3 scripts/harness/run.py --suite backend-local` は通過済みである。
- `internal/bootstrap` の Go テストは provider settings 用 secret store を file backend に差し替え、実 OS keychain を呼ばない境界へ修正済みである。
- `integration-provider-settings-wails-gateway` は `go test ./internal/controller/wails ./internal/bootstrap -run 'ProviderSettings|AppController|TranslationJobSetup|MasterPersona'`、`npm --prefix frontend run check`、`npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona` を通過済みである。
- 最終検証として `go test ./internal/...`、`npm --prefix frontend run check`、`npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona`、`npm --prefix frontend run build`、`python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local` は通過済みである。
- 最終検証中に検出した provider settings 保存 DTO の `APIKeyInput` field は transient `CredentialInput` へ修正し、`SCN-AIPSM-002/005/009` の secret boundary API テストを通過済みである。
- `App.svelte` の default provider settings gateway は test-safe な null gateway に戻し、production entrypoint の `main.ts` から Wails gateway を明示注入する構成を維持した。
- `git diff --check -- docs/exec-plans/active/ai-provider-settings-management internal frontend` は通過済みである。
- 実装後 5 観点レビューで `trust-boundary` が hard gate fail、`state-invariant` が critical を検出したため、最終検証通過状態を取り消して修正へ入る。
- 修正対象は master-persona / Job Setup の APIキー本文 seam 除去、provider settings validate snapshot、Gemini endpoint 反映、default endpoint 表示、OpenAI user-facing provider 除外、DB/secret 保存単位、Running phase snapshot 永続化、provider debug log の raw payload 除去である。
- review fix 後、`go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails -run 'ProviderSettings|Secret|Credential|Endpoint|Redaction|TranslationJobSetup|MasterPersona'`、`go test ./internal/apitest -run 'SCN_AIPSM|ProviderSettings'`、`go test ./internal/infra/ai -run 'ProviderSettings|Model|Gemini|Endpoint|DebugLog|Redaction|Provider'`、`go test ./internal/service ./internal/repository ./internal/integrationtest -run 'Snapshot|ProviderSettings|TermTranslation|BodyTranslation|PersonaGeneration|DebugLog|Redaction'` は通過済みである。
- review fix 後、`go test ./internal/...`、`npm --prefix frontend run check`、`npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona`、`npm --prefix frontend run build`、`python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local` は通過済みである。
- migration 番号の重複を避けるため、phase runtime snapshot endpoint summary migration は `011_translation_job_phase_runtime_snapshot_endpoint_summary.sql` として配置済みである。
- 再レビューで残った provider settings validate token、master-persona / Job Setup 公開 APIキー seam、provider settings frontend contract の参照側 model seam、debug log raw payload、Running phase mutable credential ref を修正済みである。
- 追加修正後、`go test ./internal/...`、`python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local`、`npm --prefix frontend run build`、`git diff --check -- docs/exec-plans/active/ai-provider-settings-management internal frontend` は通過済みである。
- 契約レビューで `SaveProviderSettings` の `credentialInput` が scope 未定義と判定されたため、`implementation-scope.md` へ transient input 例外を追加した。
- `credentialInput` は保存 command input だけで許可し、frontend state、read model、response DTO、DB、log、DOM、保存要約への保持は禁止のまま固定した。
- 契約レビューの追加指摘により、master-persona の `executionMethod` を frontend contract、generated Wails model、Wails DTO、usecase、service、repository、SQLite migration へ反映した。
- `012_master_persona_execution_method.sql` を追加し、既存 DB には `execution_method` default `single_request` を追加する。migration 再適用時の duplicate column は既存 idempotent migration と同じ扱いにした。
- `reviewback.contract.yaml` は `review_status: no_issue`、`must_fix_open: false`、`max_level: none` へ更新済みである。
- `reviewback.trust-boundary.yaml` は hard gate 通過として `review_status: no_issue`、`must_fix_open: false`、`max_level: none` へ更新済みである。
- 最終再検証として `go test ./internal/...`、`npm --prefix frontend run build`、`python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local`、`npm --prefix frontend run check`、`git diff --check -- docs/exec-plans/active/ai-provider-settings-management internal frontend` は通過済みである。
- `work_reporter` は run 全体レポート作成を停止した。理由は `work_history/runs/2026-05-04-ai-provider-settings-management-run/analysis/benchmark-score.json`、`transcript_refs.json`、`workflow-improvement-log.jsonl`、`reviewback.state-invariant.yaml`、`reviewback.responsibility-boundary.yaml` が不足しているためである。
