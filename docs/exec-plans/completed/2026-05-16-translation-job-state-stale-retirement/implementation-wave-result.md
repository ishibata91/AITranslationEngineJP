# Implementation Wave Result

- `skill`: `implement-lane`
- `status`: `backend_and_tests_completed`
- `source`: `implementation-handoff.md`, `implementation-scope.md`
- `return_to`: `implement_lane`

## Backend 実装

### `BE-TJSR-001`

- `status`: `completed`
- `changed_files`:
  - `.go-arch-lint.yml`
  - `internal/jobio/doc.go`
- `result`: `JobIOService` の product code と architecture lint 定義を削除した。
- `validation`:
  - `python3 scripts/harness/run.py --suite backend-lint`: pass
  - `python3 scripts/harness/run.py --suite structure`: pass
  - `rg -n "internal/jobio|JobIOService|jobio" internal .go-arch-lint.yml --glob '!**/*_test.go'`: exit code `1`, output empty
  - `python3 scripts/harness/run.py --suite backend-local`: pass

### `BE-TJSR-002`

- `status`: `completed`
- `changed_files`:
  - `internal/usecase/persona_generation_phase_contract.go`
- `result`: `PersonaGenerationPhaseContractStub` の cancel fixture response を `canceled` へそろえた。
- `validation`:
  - `go test ./internal/usecase`: pass
  - `rg -n "\"cancelled\"|cancelled" internal/usecase/persona_generation_phase_contract.go`: exit code `1`, output empty
  - `python3 scripts/harness/run.py --suite backend-local`: pass

## 単体テスト

- `status`: `completed`
- `changed_files`:
  - `internal/service/phase_action_enablement_helpers_test.go`
  - `internal/usecase/persona_generation_phase_contract_test.go`
- `result`: 共通操作規則、terminal / 状態不整合時の危険操作無効、`pending` 非正本化、`canceled` spelling をテストで固定した。
- `validation`:
  - `go test ./internal/usecase ./internal/service`: pass
  - `rg -n "\"cancelled\"|cancelled" internal/usecase internal/service --glob '*_test.go'`: output exists only in unrelated master persona `cancelled` variables and text
  - `rg -n "\"cancelled\"" internal/usecase internal/service --glob '*_test.go'`: exit code `1`, output empty

## シナリオテスト

- `status`: `completed_after_fix`
- `changed_files`:
  - `internal/apitest/body_translation_recovery_terminal_readiness_test.go`
  - `internal/integrationtest/translation_job_management_scenario_test.go`
- `result`: Ready job の read-only 確認、start-on-demand、body phase 操作可否、状態不整合時の読み取り非変更と危険削除拒否を追加した。
- `initial_failure`:
  - `go test ./internal/apitest ./internal/integrationtest`: fail
  - `TestSCN_TJSR_001_ReadyJobSummaryDoesNotCreateBodyPhaseRunAndStartCreatesRunningRun` が失敗した。
  - `StartBodyTranslationPhase` が body phase run 未存在時に `find body translation phase run for start: not found` を返した。
- `fix`:
  - `internal/service/body_translation_phase_service.go` で body phase run 未存在時にも `CreateJobPhaseRun` で start-on-demand できるようにした。
  - `internal/service/body_translation_phase_service_test.go` の旧 not found 期待を start-on-demand 期待へ更新した。
  - body phase は start 後に provider 実行まで同期で進むため、API テストは `pending` でない phase run が作成されることを確認する形へ調整した。
- `validation`:
  - `go test ./internal/apitest ./internal/integrationtest`: pass
  - `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`: pass
  - `rg -n "stale_selection|validation_stale|model_selection_stale" internal`: pass, reason category remains
  - `test ! -e docs/exec-plans/active/observability-log-addition`: pass
  - `rg -n "\"cancelled\"|cancelled" internal/usecase/persona_generation_phase_contract.go internal/apitest internal/integrationtest`: exit code `1`, output empty
  - `python3 scripts/harness/run.py --suite backend-local`: pass

## 追加変更ファイル

- `internal/service/body_translation_phase_service.go`: body phase run 未存在時の start-on-demand を実装した。
- `internal/service/body_translation_phase_service_test.go`: body phase run 未存在時の start-on-demand を単体テストへ反映した。
- `internal/service/translation_job_setup_service.go`: `Ready` job 作成時の `pending` phase run 事前作成を削除した。
- `internal/service/translation_job_setup_service_test.go`: setup service 経由の作成で `JOB_PHASE_RUN` を事前作成しない期待へ変更した。
- `internal/service/term_translation_phase_service.go`: 単語翻訳段階の start 時に非 `pending` run を作成するようにした。
- `internal/service/term_translation_phase_service_test.go`: phase run 不在時の start-on-demand と Ready job 読み取りを固定した。
- `internal/service/persona_generation_phase_service.go`: ペルソナ生成段階の start 時に非 `pending` run を作成するようにした。
- `internal/service/persona_generation_phase_service_test.go`: phase run 不在時の start-on-demand を固定した。

## レビュー指摘修正

### `behavior-001`

- `status`: `fixed_by_backend_implementer`
- `source`: `reviewback.behavior.yaml`
- `level`: `major`
- `result`: 実 job setup 経路で `Ready` job 作成時に `JOB_PHASE_RUN` を事前作成しないようにした。
- `result`: 単語翻訳段階、ペルソナ生成段階、本文翻訳段階は start 時に対象 phase run を作成する。
- `validation`:
  - `go test ./internal/service`: pass
  - `go test ./internal/usecase ./internal/service`: pass
  - `go test ./internal/apitest ./internal/integrationtest`: pass
  - `python3 scripts/harness/run.py --suite backend-local`: pass
  - `python3 scripts/harness/run.py --suite backend-lint`: pass
  - `python3 scripts/harness/run.py --suite structure`: pass
  - `python3 scripts/harness/run.py --suite coverage`: pass

## 残留事項

- docs 正本化判断は必要である。
- `docs/architecture.md` と architecture 図の `JobIOService` 参照は、レビュー通過後に `docs_updater` 判断へ渡す。
- completed archive は変更しない。
