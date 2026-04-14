# Work Plan

- workflow: orchestrate
- status: planned
- lane_owner: orchestrate
- scope: frontend-dip-for-unit-tests
- task_id: frontend-dip-for-unit-tests
- task_catalog_ref: N/A
- parent_phase: implementation-lane

## Request Summary

- frontend の DIP 化を行い、unit test で concrete 実装や runtime へ直結しない構造へ寄せる。
- 主目的は新機能追加ではなく、既存 frontend の依存方向整理と test seam 導入である。

## Decision Basis

- 2026-04-14 時点の user request は「frontend の DIP 化をしたい。ユニットテストのため」であり、主目的は testability 向上のための構造改善である。
- そのため `task_mode: refactor` を暫定固定し、`distill -> design -> implement/tests -> review(ui-check, implementation-review)` の流れで扱う。
- 具体的な対象画面や module は未確定のため、distill で affected surface と code pointer を最小化してから scope freeze を行う。

## Task Mode

- `task_mode`: refactor
- `goal`: unit test のために frontend 依存を差し替え可能な境界へ整理し、上位層 test が結合テスト化しない構造へ改める。
- `constraints`: docs 正本は human 先行でのみ更新する。orchestrate 自身は詳細調査と実装を行わない。frontend を含むため close 前に `ui-check` と `implementation-review` を必須とする。
- `close_conditions`: review が pass を返すこと。frontend unit test 戦略に必要な seam が固定されること。`python3 scripts/harness/run.py --suite all` を通すこと。

## Facts

- user request は frontend に対する DIP 化であり、unit test 容易性が主要目的である。
- structure harness は着手前に pass 済みである。
- active plan には同名 task がまだ存在しない。

## Functional Requirements

- `summary`: frontend の runtime / store / transport / UI orchestration 依存を見直し、unit test で fake や stub へ差し替えられる境界を定義する。
- `in_scope`: affected frontend surface の特定、dependency seam 候補の整理、implementation scope freeze、必要な unit test 方針の整理。
- `non_functional_requirements`: 既存の user-visible behavior を変えない。不要な docs 正本更新を行わない。最小差分より correctness を優先する。
- `out_of_scope`: backend の恒久仕様変更、docs 正本更新、対象外画面への先回り abstraction 導入。
- `open_questions`: どの page / module を先に DIP 化するか。どの concrete dependency が unit test を阻害しているか。既存 test harness と整合する seam は何か。
- `required_reading`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-14-frontend-dip-for-unit-tests.implementation-scope.md
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Work Brief

- `implementation_target`: frontend DIP for unit tests
- `accepted_scope`: distill 完了前のため暫定。affected surface 特定後に scope freeze する。
- `parallel_task_groups`: distill 後に決定する。
- `tasks`: distill で facts / constraints / gaps / code pointer を回収し、その後 design で implementation scope を固定する。
- `validation_commands`: python3 scripts/harness/run.py --suite structure

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: frontend unit test を阻害する concrete dependency が component か store か bridge 層に残っている可能性が高い。
- `observation_points`: frontend entrypoint、state 管理、Wails bridge、既存 unit test、page-level module。
- `residual_risks`: 対象 surface を広く取りすぎると ownership を安全に切れない。

## Acceptance Checks

- unit test 目的の DIP 対象 surface が plan 上で明確になる。
- scope freeze 後の実装対象が narrow scope として handoff 可能になる。
- frontend close gate に必要な review mode が plan に残る。

## Required Evidence

- distill の facts / constraints / gaps / related_code_pointers
- design の implementation-scope artifact
- frontend 実装後の unit test 証跡
- ui-check と implementation-review の結果

## HITL Status

- `functional_or_design_hitl`: pending
- `approval_record`: 2026-04-14 human request: frontend の DIP 化を行い、unit test 可能な構造へ寄せたい。

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- planned
