# Implementation Handoff Input: E2E System Test Failure Refactor

## 起動元

- role: `refactor_lane`
- skill: `refactor-lane`
- task folder: `docs/exec-plans/active/e2e-system-test-failure-refactor/`
- approved scope: `docs/exec-plans/active/e2e-system-test-failure-refactor/implementation-scope.md`
- approval: 人間が `approve` と回答した。

## 共通入力

- `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
- `docs/exec-plans/active/e2e-system-test-failure-refactor/structure-quality-investigation.md`
- `docs/exec-plans/active/e2e-system-test-failure-refactor/test-quality-investigation.md`
- `docs/exec-plans/active/e2e-system-test-failure-refactor/implementation-scope.md`
- `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
- `docs/e2e-test-design/test-design.csv`
- `docs/e2e-test-guidelines.md`
- `docs/coding-guidelines-tests.md`

## 共通禁止事項

- `EF-001` から `EF-005` 以外へ修正範囲を広げない。
- docs 正本文を更新しない。
- `.codex/` を更新しない。
- 実外部 API、実 secret、実利用者データへ到達する経路を作らない。
- secret 本体を UI、DTO、read model、URL、log、error summary、audit、要求捕捉へ出さない。
- 他 agent の変更を revert しない。

## 実装順序

1. `wave-1`: `H-BE-001`, `H-BE-002`, `H-FE-001`
2. `wave-1b`: `H-FE-002`, `H-FE-003`, `H-FE-004`
3. `wave-2`: `H-INT-001`, `H-INT-002`, `H-INT-003`
4. `wave-3`: `H-ST-001`, `H-ST-002`, `H-ST-003`, `H-ST-004`, `H-ST-005`, `H-UT-001`, `H-UT-002`
5. `wave-4`: `H-FINAL-001`

`H-FE-001` は `App.svelte` の共有 composition root を触るため、他 frontend handoff と並列にしない。

## Wave 1 起動入力

### H-BE-001

- agent: `backend_implementer`
- skill: `implement-backend`
- 目的: `EF-001` provider settings 保存境界を修正する。
- owned scope: `internal/service/provider_settings_service.go`, `internal/usecase/provider_settings_usecase.go`, `internal/controller/wails/provider_settings_controller.go`, related backend tests.
- first action: `SaveProviderSettings` の endpoint 形式不正を既存 DTO と test で確認し、不正 endpoint で保存済み状態を更新しない clause を作る。
- validation: `python3 scripts/harness/run.py --suite backend-local`

### H-BE-002

- agent: `backend_implementer`
- skill: `implement-backend`
- 目的: `EF-003` paused job resume action read model を修正する。
- owned scope: `internal/service/translation_job_management_service.go`, `internal/usecase/translation_job_management_usecase.go`, `internal/controller/wails/translation_job_management_controller.go`, related backend tests.
- first action: `ResumeJob` が warning 固定ではなく、再開入口成立を表す result を返す範囲を確認する。
- validation: `python3 scripts/harness/run.py --suite backend-local`

### H-FE-001

- agent: `frontend_implementer`
- skill: `implement-frontend`
- 目的: `EF-001` と `EF-004` root composition gateway wiring を修正する。
- owned scope: `frontend/src/ui/App.svelte`, `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, app shell tests if needed.
- first action: system-test 起動時の root が `main.ts` factory 注入か `App.svelte` fallback かを確認し、provider settings と output management が fallback 経路でも `null gateway` にならないようにする。
- validation: `python3 scripts/harness/run.py --suite frontend-local`

## 完了報告

各 agent は次を返す。

- completed_handoffs
- touched_files
- implemented_scope
- test_results
- blocked_items
- residual_risks
- docs_changes: none
