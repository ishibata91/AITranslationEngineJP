# Light Change Planning

- `skill`: `light-change-planning`
- `status`: `ready`
- `decision`: `範囲内修正`
- `return_to`: `light_change_lane`

## 人間要望

- `summary`: 翻訳ジョブ状態関連の追加差分を、stale 廃止を主目的に整理する。
- `expected_result`: 古い設計名、空 package、重複した phase 別 policy wrapper、旧 task-local 参照が減る。
- `forbidden_scope`: 新しい状態遷移、DB 永続値、Wails 公開 DTO、画面仕様、ドメイン上の `stale_*` 理由分類は変えない。

## 根拠参照

- `detail_specs`: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `task_local_artifacts`: `docs/exec-plans/completed/2026-05-10-translation-job-state-machine-redesign/plan.md`, `docs/exec-plans/completed/2026-05-10-translation-job-state-machine-redesign/review-aggregation.md`, `docs/exec-plans/active/observability-log-addition/`
- `docs`: `docs/index.md`, `docs/architecture.md`, `docs/spec.md`, `docs/exec-plans/active/README.md`
- `existing_implementation`: `internal/statemachine/doc.go`, `internal/jobio/doc.go`, `internal/usecase/translationjobpolicy/policy.go`, `internal/usecase/*_phase_usecase.go`, `internal/service/*_phase_service.go`, `.go-arch-lint.yml`
- `validation_logs`: `docs/exec-plans/completed/2026-05-10-translation-job-state-machine-redesign/final-validation.md`

## 突き合わせ結果

- `request_vs_specs`: 仕様は `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の分離を要求する。stale 廃止は状態意味を変えず、実装構造だけを薄くする。
- `request_vs_task_artifacts`: 完了済み state machine 再設計は `TranslationJobPolicy` を採用済みである。active observability 成果物には `StateMachine` と `JobIOService` の古い前提が残る。
- `request_vs_existing_code`: `internal/statemachine` と `internal/jobio` は `doc.go` だけである。UseCase と Service には phase 別の同型分岐が残る。
- `conflicts`: `JobIOService` は architecture 正本に残るため、廃止する場合は docs 正本化判断が必要である。

## 実装入力

- `implementation_skill`: `implement-backend`
- `change_targets`: `internal/statemachine/`, `internal/jobio/`, `.go-arch-lint.yml`, `internal/usecase/phase_policy_helpers.go`, `internal/usecase/*_phase_usecase.go`, `internal/service/*_phase_service.go`, 関連 unit test
- `forbidden_changes`: `docs/exec-plans/completed/**`, domain の `stale_selection` / `validation_stale` / `model_selection_stale`, DB schema, Wails DTO, frontend UI
- `validation_commands`: `gofmt -l internal/usecase internal/service`; `python3 scripts/harness/run.py --suite backend-local`; `python3 scripts/harness/run.py --suite backend-lint`; `python3 scripts/harness/run.py --suite structure`; `python3 scripts/harness/run.py --suite coverage`
- `docs_to_read`: `docs/architecture.md`, `docs/spec.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/coding-guidelines-backend.md`, `docs/coding-guidelines-tests.md`

## 正本化判断材料

- `spec_change`: `no`
- `human_approved_permanent_change`: `unknown`
- `docs_update_target`: `docs/architecture.md`, `docs/diagrams/backend/backend-architecture.puml`, `docs/exec-plans/active/observability-log-addition/*`

## 停止または戻し

- `reason`: `JobIOService` を廃止するか別 task で実装するかの判断が衝突する場合は停止する。
- `missing_information`: `JobIOService` の扱いに関する人間確認。
- `handoff_prompt`: `JobIOService` を architecture 正本から外す場合は docs 正本化へ渡す。実装として作る場合は新規実装レーンへ戻す。

