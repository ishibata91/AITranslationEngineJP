# Work Plan

- workflow: orchestrate
- status: completed
- lane_owner: orchestrate
- scope: backend-frontend-unit-test-abstraction-validation
- task_id: backend-frontend-unit-test-abstraction-validation
- task_catalog_ref: N/A
- parent_phase: refactor-lane

## Request Summary

- backend と frontend の現存ユニットテストが、DIP 後の抽象境界を前提に成立しているかを検証する。
- 具体実装への過剰依存、層逆流、DI seam 不足が残っていないかを確認する。

## Decision Basis

- user は architecture 全般を DIP 化した後の妥当性確認を求めている。
- 直近 completed plan に frontend DIP for unit tests があり、今回は frontend 限定ではなく backend も含めた横断確認が必要である。
- 現時点では product 挙動変更より evidence 収集と評価が主目的である。

## Task Mode

- `task_mode`: refactor
- `goal`: backend は abstraction-first 判定を維持しつつ、frontend の `App.test.ts` 系を contract-only な unit test へ戻す。
- `constraints`: docs 正本は更新しない。product code は変更しない。`App.test.ts` は `CreateMasterDictionaryScreenController` と `MasterDictionaryScreenControllerContract` だけを境界として扱う。review_mode: ui-check と implementation-review を close 前に必須とする。
- `close_conditions`: frontend test support と focused runtime test への分離が完了し、review が pass を返し、`cd frontend && npm test` と `python3 scripts/harness/run.py --suite all` が通ること。

## Facts

- `docs/exec-plans/completed/2026-04-14-frontend-dip-for-unit-tests.md` が存在する。
- structure harness は 2026-04-15 に pass 済みである。
- active plan には別件として `2026-04-12-master-dictionary-category-and-count-bug.md` が存在する。

## Functional Requirements

- `summary`:
  - backend の abstraction-first 判定結果を保持しつつ、frontend の unit test を concrete factory / runtime detail 依存から切り離す。
  - `App.test.ts` は contract fake 駆動に戻し、runtime event 詳細は focused test へ寄せる。
- `in_scope`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/test/setup.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts`
  - 直近の DIP 関連 completed plan
- `non_functional_requirements`:
  - App-level test から `@controller/master-dictionary` と `window.runtime` 詳細を消す。
  - 既存の UI シナリオ観点は維持する。
  - 差分は test support と test file に閉じる。
- `out_of_scope`:
  - product code の変更
  - backend 側 unit test の修正
  - docs 正本更新
- `open_questions`:
  - なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-14-frontend-dip-for-unit-tests.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-15-backend-frontend-unit-test-abstraction-validation.implementation-scope.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/test/setup.ts`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-15-backend-frontend-unit-test-abstraction-validation.implementation-scope.md
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/test/setup.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`

## Work Brief

- `implementation_target`: frontend-tests
- `accepted_scope`:
  - `frontend-app-contract-only-tests`
  - `frontend-runtime-and-controller-focused-tests`
- `parallel_task_groups`:
  - frontend app contract-only tests
  - frontend runtime and controller focused tests
- `tasks`:
  - design
  - review
  - implement
  - tests
  - review
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `cd frontend && npm test`
  - `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: review により frontend `App.test.ts` が concrete factory と runtime detail に依存していると判定済み。
- `trace_hypotheses`:
  - contract fake 化で App-level unit test は abstraction-first に戻せる。
  - runtime event 詳細は adapter focused test へ分離できる。
- `observation_points`:
  - `App.test.ts` の import 先
  - fake controller contract の表現力
  - runtime event focused test の独立性
- `residual_risks`:
  - App-level test の contract-only 化に伴い、coverage gate を lower-level unit test で補完する必要がある。

## Acceptance Checks

- frontend App-level test が contract-only になっていること。
- runtime event 詳細が focused test へ分離されていること。
- lower-level controller / factory unit test で coverage gate を維持していること。
- backend は no-change で abstraction-first 判定を維持すること。

## Required Evidence

- backend no-change 判定の evidence
- frontend App-level test の contract-only 化 evidence
- runtime focused test の追加 evidence
- controller / factory focused test の追加 evidence

## Validation Results

- `frontend`: `npm run lint` pass
- `frontend`: `npm test` pass (`6 files / 24 tests`)
- `repo`: `python3 scripts/harness/run.py --suite all` pass
- `sonar`: coverage `72.1% >= 70.0%`、line `75.3%`、branch `49.0%`
- `backend`: code change なし、既存 abstraction-first 判定を維持

## HITL Status

- `functional_or_design_hitl`: not_required
- `design_review_status`: pass
- `implementation_review_status`: pass
- `ui_check_status`: pass
- `approval_record`: 2026-04-15 user request: backend と frontend の両方で、現存 unit test が抽象前提になっているか確認したい。

## Closeout Notes

- `canonicalized_artifacts`:
  - `docs/exec-plans/completed/2026-04-15-backend-frontend-unit-test-abstraction-validation.md`
  - `docs/exec-plans/completed/2026-04-15-backend-frontend-unit-test-abstraction-validation.implementation-scope.md`
- backend の reviewed unit test は no-change で abstraction-first 判定を維持した。
- frontend の `App.test.ts` は contract-only へ移し、runtime / controller / factory の lower-level test で coverage gate を回復した。
- UI live check では `#master-dictionary` と `#translation-management` の表示維持を確認し、console error は `favicon.ico` 404 のみだった。

## Outcome

- completed
