# Implementation Scope: 2026-05-16-translation-job-state-stale-retirement

- `skill`: `implementation-scope`
- `status`: `ready_for_implement_lane`
- `source_plan`: `./implement-lane-task-frame.md`
- `human_review_status`: `approved`
- `approval_record`: `./human-design-review-request.md`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `design_diff`: `./design-diff.md`
- `design_diff_component`: `./design-diff.component.puml`
- `design_diff_sequence`: `./design-diff.sequence.puml`
- `existing_backend_result`: `./backend-implementation-result.md`
- `human_review`: `./human-design-review-request.md`
- `canonical_requirements`: `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`

## Fixed Decisions

- `needs_human_decision`: `0`
- `pending` は canonical state へ昇格しない。
- `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
- read model の `CanPause`、`CanResume`、`CanRetry`、`CanCancel` は、`TranslationJobPolicy` と同じ state 事実から導く。
- `PolicyResult`、rule 名、policy 判定履歴は DB、DTO、repository 永続契約、read model の永続値へ出さない。
- `JobIOService` は stale として扱い、architecture 正本から外す。
- `observability-log-addition` は completed archive であるため、今回の active task-local 更新対象にしない。
- `docs/exec-plans/completed/**` は変更しない。
- `cancelled` fixture spelling は今回範囲に含め、`canceled` へそろえる。
- `stale_selection`、`validation_stale`、`model_selection_stale` は削除対象にしない。
- UI、DB schema、Wails DTO は変更しない。
- provider raw payload、prompt 全文、翻訳本文全文、credential 実値、API key、endpoint 実値は新しく保存、表示、ログ出力しない。

## Existing Backend Scope To Preserve

既存 backend 実装差分は保持してよい。
この差分は `backend-implementation-result.md` の検証済み証跡として扱い、追加 handoff で巻き戻さない。

- `.go-arch-lint.yml`: `statemachine` component と許可依存の削除。
- `internal/statemachine/doc.go`: `doc.go` だけの旧設計 package の削除。
- `internal/usecase/phase_policy_helpers.go`: phase 別 policy input 生成 helper の統合先。
- `internal/usecase/term_translation_phase_usecase.go`: 単語翻訳段階の policy input 生成を共通 helper へ寄せた差分。
- `internal/usecase/persona_generation_phase_usecase.go`: ペルソナ生成段階の policy input 生成を共通 helper へ寄せた差分。
- `internal/usecase/body_translation_phase_usecase.go`: 本文翻訳段階の policy input 生成を共通 helper へ寄せた差分。
- `internal/service/phase_action_enablement_helpers.go`: phase 共通の操作可否 helper。
- `internal/service/term_translation_phase_service.go`: 単語翻訳段階の操作可否判定を共通 helper へ寄せた差分。
- `internal/service/persona_generation_phase_service.go`: ペルソナ生成段階の操作可否判定を共通 helper へ寄せた差分。
- `internal/service/body_translation_phase_service.go`: 本文翻訳段階の操作可否判定を共通 helper へ寄せた差分。

保持条件:

- `TranslationJobPolicy` の状態意味を変更しない。
- `commonPhaseActionAvailability` は `TranslationJobPolicy` package を service 層へ直接 import しない。
- `PolicyResult`、rule 名、policy 判定履歴を DTO、DB、repository 永続契約へ出さない。
- `stale_selection`、`validation_stale`、`model_selection_stale` を削除しない。

## Required Artifacts

| 成果物 | 必要有無 | 担当 agent | 理由 |
| --- | --- | --- | --- |
| `backend 実装` | 必要 | `backend_implementer` | 承認済み差分に対し、`JobIOService` stale 廃止と `cancelled` spelling 統一が追加で残っている。 |
| `frontend 実装` | 不要 | なし | UI、画面文言、layout、style を変更しない。 |
| `統合境界実装` | 不要 | なし | DB schema、Wails DTO、gateway、adapter 契約を変更しない。 |
| `単体テスト` | 必要 | `implementation_unit_tester` | 共通操作規則、`pending` 非正本化、spelling 統一を lower-level で固定する。 |
| `シナリオテスト` | 必要 | `implementation_scenario_tester` | API 境界で state、操作可否、stale reason、terminal guard を確認する。 |
| `docs 正本化判断` | 必要 | `implement_lane` | `JobIOService` を architecture 正本から外す承認済み仕様変更がある。docs 正本本文の更新は `docs_updater` 判断へ渡す。 |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `BE-TJSR-001`, `BE-TJSR-002` | なし | `BE-TJSR-001 <-> BE-TJSR-002` | なし |
| `wave-2` | `UNIT-TJSR-001`, `SCN-TJSR-001` | `BE-TJSR-001`, `BE-TJSR-002` | `UNIT-TJSR-001 <-> SCN-TJSR-001` | なし |
| `wave-3` | `DOCS-TJSR-DECISION` | `BE-TJSR-001`, `BE-TJSR-002`, `UNIT-TJSR-001`, `SCN-TJSR-001` | なし | `depends_on` |

## Handoffs

### `BE-TJSR-001`: `JobIOService` stale 実体と lint 定義を削除する

- `implementation_target`: `JobIOService` を状態保存境界として扱わないようにし、`internal/jobio/` と architecture lint の `jobio` component を削除する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `provider raw payload`, `prompt 全文`, `翻訳本文全文`, `credential 実値`, `API key`, `endpoint 実値`
- `owned_scope`:
  - `.go-arch-lint.yml`
  - `internal/jobio/doc.go`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `BE-TJSR-002`
- `parallel_blockers`: なし
- `estimated_size`: `2 files`, `40 changed lines 以下`, `通常`
- `first_action`: `.go-arch-lint.yml` の `jobio` component 定義を削除する。対応する完了条件は `jobio component と許可依存が lint 設定から消えること`。最初に行う理由は、削除対象の公開依存境界を先に閉じるためである。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `python3 scripts/harness/run.py --suite structure`
  - `rg -n "internal/jobio|JobIOService|jobio" internal .go-arch-lint.yml --glob '!**/*_test.go'`。期待結果は exit code `1` と出力なしである。
- `completion_signal`:
  - `internal/jobio/` が product code から削除されている。
  - `.go-arch-lint.yml` に `jobio` component と `jobio` 許可依存が残っていない。
  - `JobIOService` を実体化する新 package、service、repository、DTO が追加されていない。
  - 状態事実の取得と保存は既存 usecase、service、repository 境界で追える。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - docs 正本本文の `JobIOService` 参照は、この handoff では変更しない。
  - docs 正本化は `DOCS-TJSR-DECISION` へ分離する。

### `BE-TJSR-002`: `cancelled` fixture spelling を `canceled` へそろえる

- `implementation_target`: `PersonaGenerationPhaseContractStub` の cancel fixture spelling を正本 spelling の `canceled` へそろえる。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `provider raw payload`, `prompt 全文`, `翻訳本文全文`, `credential 実値`, `API key`, `endpoint 実値`
- `owned_scope`:
  - `internal/usecase/persona_generation_phase_contract.go`
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `BE-TJSR-001`
- `parallel_blockers`: なし
- `estimated_size`: `1 file`, `10 changed lines 以下`, `通常`
- `first_action`: `internal/usecase/persona_generation_phase_contract.go` の `PersonaGenerationPhaseContractStub` cancel fixture response を `cancelled` から `canceled` へ変更する。対応する完了条件は `fixture response に cancelled spelling が残らないこと`。最初に行う理由は、承認済み spelling 差分を 1 clause で閉じられるためである。
- `validation_commands`:
  - `go test ./internal/usecase`
  - `rg -n "\"cancelled\"|cancelled" internal/usecase/persona_generation_phase_contract.go`。期待結果は exit code `1` と出力なしである。
- `completion_signal`:
  - `PersonaGenerationPhaseContractStub` の cancel fixture response は `canceled` を返す。
  - 正本 spelling の `Canceled` / `canceled` と、残留 spelling の `cancelled` が fixture 内で混在しない。
  - master persona など別文脈の変数名や文言は、この handoff で無関係に変更しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `cancelled` という英語一般語が別ドメインの変数名として存在する場合は、正本 state spelling ではないため対象外にする。

### `UNIT-TJSR-001`: 状態事実と spelling の単体テストを固定する

- `implementation_target`: 共通操作規則、`pending` 非正本化、`canceled` spelling 統一を単体テストで固定する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `provider raw payload`, `prompt 全文`, `翻訳本文全文`, `credential 実値`, `API key`, `endpoint 実値`
- `owned_scope`:
  - `internal/usecase/*_test.go`
  - `internal/service/*_test.go`
  - 必要な既存 test helper
- `depends_on`: `BE-TJSR-001`, `BE-TJSR-002`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `SCN-TJSR-001`
- `parallel_blockers`: なし
- `estimated_size`: `6 files 以下`, `250 changed lines 以下`, `通常`
- `first_action`: `internal/service/phase_action_enablement_helpers.go` の helper を直接または既存 service test 経由で検証する table-driven test を追加する。対応する完了条件は `Running`, `Paused`, `RecoverableFailed`, terminal state の操作可否が phase type 横断で一致すること。最初に行う理由は、read model 操作可否の共通規則が今回の中心仕様であるためである。
- `validation_commands`:
  - `go test ./internal/usecase ./internal/service`
  - `rg -n "\"cancelled\"|cancelled" internal/usecase internal/service --glob '*_test.go'`。期待結果は cancel state fixture に関する出力なしである。
- `completion_signal`:
  - `Running` は pause だけを許可する。
  - `Paused` は resume と cancel を許可する。
  - `RecoverableFailed` は retry を許可する。
  - terminal job または状態不整合では危険操作を無効にする。
  - `pending` は成功 state、操作可能 state、正本 DTO state として扱われない。
  - `PersonaGenerationPhaseContractStub` の cancel fixture test は `canceled` を期待する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - product code を直接壊した test helper だけを更新する。
  - DB migration、Wails DTO、frontend fixture は対象にしない。

### `SCN-TJSR-001`: API 境界で state stale 廃止を確認する

- `implementation_target`: 承認済みシナリオ `SCN-TJSR-001` から `SCN-TJSR-008` までを、API 境界、lower-level 境界、補助検索に分けて確認する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `credential 状態分類だけ`
  - `secret_values_for_provider_external_api_internal_auth`: `実 API key は使わない`
  - `secret_resolution_owner_layer`: `既存 provider fake と credential stub`
  - `forbidden_outputs`: `provider raw payload`, `prompt 全文`, `翻訳本文全文`, `credential 実値`, `API key`, `endpoint 実値`
- `owned_scope`:
  - `internal/apitest/*_test.go`
  - `internal/integrationtest/*_test.go`
  - 必要な scenario fixture
- `depends_on`: `BE-TJSR-001`, `BE-TJSR-002`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `UNIT-TJSR-001`
- `parallel_blockers`: なし
- `estimated_size`: `8 files 以下`, `450 changed lines 以下`, `通常`
- `first_action`: `internal/apitest` に Ready job の読み取りと phase start を確認する API test を追加する。対応する完了条件は `Ready job の読み取りでは phase run を作成せず、start 許可時だけ Running phase run を作成すること`。最初に行う理由は、`pending` 非正本化と start-on-demand の受け入れ条件を最初に閉じるためである。
- `validation_commands`:
  - `go test ./internal/apitest ./internal/integrationtest`
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`
  - `rg -n "stale_selection|validation_stale|model_selection_stale" internal`
  - `test ! -e docs/exec-plans/active/observability-log-addition`
- `completion_signal`:
  - Ready job の query だけでは `JOB_PHASE_RUN` が作成されない。
  - phase start 許可時だけ `Running` の `JOB_PHASE_RUN` が作成される。
  - read model の操作可否は phase type 横断で共通操作規則と矛盾しない。
  - job state と phase state が食い違う fixture で、読み取りだけでは永続 state を変更せず危険操作を無効にする。
  - terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。
  - `stale_selection`、`validation_stale`、`model_selection_stale` は reason category として残る。
  - `observability-log-addition` は active task-local として存在しない。
  - `cancelled` fixture spelling は lower-level 確認で `canceled` へそろっている。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - UI 人間操作 E2E は不要である。
  - 実 API、実 credential、実 provider endpoint は使わない。
  - `SCN-TJSR-008` は API 境界の新規 DTO 変更ではなく、fixture spelling の lower-level 確認として扱う。

### `DOCS-TJSR-DECISION`: docs 正本化判断を分離する

- `implementation_target`: `JobIOService` を architecture 正本から外す承認済み仕様変更を、docs 正本化判断へ渡す。
- `implementation_artifact`: `docs 正本化判断`
- `implementation_skill`: `N/A`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `provider raw payload`, `prompt 全文`, `翻訳本文全文`, `credential 実値`, `API key`, `endpoint 実値`
- `owned_scope`:
  - `implementation-scope.md` では docs 正本本文を変更しない。
  - `implement_lane` がレビュー通過後に `正本化判断` を作る。
  - 必要と判断された場合だけ `docs_updater` が `docs/architecture.md` と関連 architecture 図を扱う。
- `depends_on`: `BE-TJSR-001`, `BE-TJSR-002`, `UNIT-TJSR-001`, `SCN-TJSR-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `estimated_size`: `docs updater 判断対象`, `implementation handoff ではない`
- `first_action`: `implement_lane` が最終レビュー後に `正本化判断` を作り、`JobIOService` を architecture 正本から外す必要有無を明示する。対応する完了条件は `docs 正本化が必要か不要かをレビュー証跡に基づいて固定すること`。最初に行う理由は、docs 正本本文の変更権限が Codex implementation lane にはないためである。
- `validation_commands`:
  - `rg -n "JobIOService|internal/jobio|jobio" docs/architecture.md docs/diagrams .go-arch-lint.yml internal`
- `completion_signal`:
  - docs 正本化判断は必要である。
  - Codex implementation lane には docs 正本本文の編集を渡さない。
  - completed archive は変更対象にしない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `notes`:
  - `docs/spec.md` は `pending` 非正本化と `Canceled` spelling を既に持つため、変更要否は docs 正本化判断で再確認する。
  - `docs/detail-specs/*` は今回の product 差分で仕様を変える必要がある場合だけ `docs_updater` の判断対象にする。

## Final Validation After Handoffs

`implement_lane` は全 handoff 完了後に次を確認する。

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite coverage`
- `git diff --check`
- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`
- `rg -n "internal/jobio|JobIOService|jobio" internal .go-arch-lint.yml --glob '!**/*_test.go'`
- `rg -n "\"cancelled\"|cancelled" internal/usecase/persona_generation_phase_contract.go`

上記 2 件の `rg` は、期待結果を exit code `1` と出力なしとして扱う。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `N/A`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`
- `docs_changes`: `none`

## Return To Implement Lane

`implementation-scope.md` は承認済み人間設計レビュー後の実装範囲として固定済みである。
`implement_lane` は、この成果物から実装引き継ぎ入力を作る。

未決事項:

- `observability-log-addition` completed archive の旧名参照は今回変更しない。
- docs 正本化判断は必要である。実行時期と `docs_updater` 起動判断は `implement_lane` が行う。
