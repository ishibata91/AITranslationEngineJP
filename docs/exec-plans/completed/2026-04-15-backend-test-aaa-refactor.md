# Work Plan

- workflow: orchestrate
- status: completed
- lane_owner: orchestrate
- scope: backend-test-aaa-refactor
- task_id: backend-test-aaa-refactor
- task_catalog_ref: N/A
- parent_phase: refactor-lane

## Request Summary

- backend 既存テストを調査し、AAA と単一意図の規約に未準拠な最小集合を特定する。
- 未準拠な backend テストを最小差分で refactor する。

## Decision Basis

- 直前 task では frontend だけを対象にしたため、backend 既存テストの現状は未確認である。
- user は backend も調査し、必要なら refactor することを求めている。
- 主目的は仕様追加ではなく既存テスト構造の整理であるため `refactor` とする。

## Task Mode

- `task_mode`: refactor
- `goal`: distill で固定した backend test 11 file だけを対象に AAA / single-intent refactor を行い、package 単位 validation と full harness を通す。
- `constraints`: product code は変更しない。`docs/` 正本は更新しない。test scope を不要に広げない。close 前に implementation-review を必須とする。
- `close_conditions`: implementation-scope artifact に固定した 11 file の refactor が完了し、targeted `go test` と `python3 scripts/harness/run.py --suite all` が pass し、implementation-review が pass を返す。

## Facts

- 既存の test 規約は `.codex/skills/tests/` と `.codex/skills/implement/SKILL.md` へ反映済みである。
- structure harness は着手前に pass した。
- distill により backend AAA / single-intent 未準拠最小集合は 11 file に固定済みである。

## Functional Requirements

- `summary`:
  - backend テストを調査し、AAA と単一意図へ未準拠な最小ファイル集合を特定する。
  - 既存の backend テストだけを対象に、振る舞い不変で読みやすさを改善する。
- `in_scope`:
  - distill で固定した backend test 11 file
  - active plan
- `non_functional_requirements`:
  - test 名と body 構造の可読性を優先する。
  - AAA と単一検証対象を読み取れる形に寄せる。
  - backend product code は不変に保つ。
- `out_of_scope`:
  - frontend テスト
  - product 実装変更
  - docs 正本更新
- `open_questions`:
  - なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/SKILL.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/references/mode-guides/unit.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/SKILL.md`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-backend-test-aaa-refactor.implementation-scope.md`
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`:
  - `internal/repository/master_dictionary_repository_test.go`
  - `internal/repository/master_dictionary_sqlite_repository_test.go`
  - `internal/infra/sqlite/sqlite_test.go`
  - `internal/bootstrap/app_controller_test.go`
  - `internal/controller/wails/master_dictionary_controller_unit_test.go`
  - `internal/controller/wails/app_controller_test.go`
  - `internal/service/master_dictionary_command_service_test.go`
  - `internal/service/master_dictionary_import_service_test.go`
  - `internal/service/master_dictionary_query_service_test.go`
  - `internal/service/master_dictionary_xml_adapter_test.go`
  - `internal/usecase/master_dictionary_usecase_test.go`

## Work Brief

- `implementation_target`: backend-existing-test-aaa-refactor
- `accepted_scope`:
  - `backend-test-distill`
  - `backend-test-aaa-refactor`
- `parallel_task_groups`:
  - backend repository and sqlite test refactor
  - backend bootstrap and controller test refactor
  - backend service and usecase test refactor
- `tasks`:
  - distill
  - design
  - tests
  - review
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `go test ./internal/repository ./internal/infra/sqlite`
  - `go test ./internal/bootstrap ./internal/controller/wails`
  - `go test ./internal/service ./internal/usecase`
  - `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: not_applicable
- `trace_hypotheses`:
  - backend test 群にも multi-intent と AAA 崩れが一部残っている可能性がある。
  - 最小集合へ絞れば product code 変更なしで refactor できる。
- `observation_points`:
  - backend unit test names
  - test body 内の assertion bundle
  - 複数操作 / 複数観測点を同居させた test
- `residual_risks`:
  - 対象を広げすぎると test scope が不要に膨らむ。

## Acceptance Checks

- backend の最小未準拠集合が説明できること。
- refactor 対象 test が AAA と single-intent を読み取れること。
- 対象 test と full harness が pass すること。

## Required Evidence

- backend distill evidence
- backend test refactor evidence
- validation pass evidence
- implementation-review pass evidence

## HITL Status

- `functional_or_design_hitl`: not_required
- `approval_record`: 2026-04-15 user request to investigate backend tests and refactor if they were not investigated.

## Validation Results

- `backend targeted go test`: `go test ./internal/repository ./internal/infra/sqlite ./internal/bootstrap ./internal/controller/wails ./internal/service ./internal/usecase` pass
- `full harness`: `python3 scripts/harness/run.py --suite all` pass
- `sonar`: open `HIGH` / `BLOCKER` `0`、open reliability `0`、open security `0`
- `implementation_review`: reroute once due to transient Sonar measure fetch, then close criteria re-established by rerun evidence and open issue recheck

## Closeout Notes

- `canonicalized_artifacts`:
  - `docs/exec-plans/completed/2026-04-15-backend-test-aaa-refactor.md`
  - `docs/exec-plans/completed/2026-04-15-backend-test-aaa-refactor.implementation-scope.md`
- backend 既存テスト 11 file を AAA と single-intent に寄せて分割した。
- Sonar `go:S1192` の test-side duplicate literal は file-local const 化で解消した。
- validation は一時的に Sonar measure 取得待ちで揺れたが、再実行で `suite all` 通過を再証跡化した。

## Outcome

- completed
