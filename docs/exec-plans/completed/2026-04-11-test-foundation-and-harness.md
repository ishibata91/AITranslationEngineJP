# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: test-foundation-and-harness
- task_id: 2026-04-11-test-foundation-and-harness
- task_catalog_ref:
- parent_phase: quality-foundation

## 要求要約

- `docs/tech-selection.md` に合わせて frontend、backend、system test の基盤を追加し、harness から実行できる状態にする。

## 判断根拠

<!-- Decision Basis -->

- 現状の repo は lint と structure check まではあるが、`Vitest`、`go test`、`Playwright` の共通入口が未整備である。
- `execution` suite は root `package.json` を入口にしているため、test 系も root script に揃える必要がある。
- system test は user 指定により `wails dev -browser` を正本の起動対象とする。

## 対象範囲

- `/package.json`
- `/frontend/package.json`
- `/frontend/vite.config.ts`
- `/frontend/tsconfig.json`
- `/frontend/src/**/*.test.ts`
- `/internal/**/*_test.go`
- `/tests/system/**/*.spec.ts`
- `/playwright.config.ts`
- `/scripts/harness/`
- `/scripts/test/`

## 対象外

- `docs/tech-selection.md` と `docs/lint-policy.md` の正本更新
- CI workflow の追加
- Wails native window 固有挙動の自動化

## 依存関係・ブロッカー

- `Playwright` は browser binary install が必要である。
- `wails dev -browser` を使うため、local に `wails` CLI が必要である。

## 並行安全メモ

- shared file は root `package.json`、`frontend/package.json`、`scripts/harness/*` に集中する。

## 機能要件

- `summary`: test を追加できる最低限の foundation を repo に常設する。
- `in_scope`: frontend unit test、backend unit test、system test、harness 統合、sample test。
- `out_of_scope`: coverage upload、CI、複数 browser matrix。
- `open_questions`: なし。
- `required_reading`: `docs/tech-selection.md`、`docs/lint-policy.md`、`wails.json`、`scripts/harness/README.md`。

## UI モック

- `artifact_path`: `docs/exec-plans/active/2026-04-11-test-foundation-and-harness.ui.html`
- `final_artifact_path`: `docs/mocks/test-foundation/index.html`
- `summary`: N/A

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/2026-04-11-test-foundation-and-harness.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/test-foundation-and-harness.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`: app shell の boot smoke と最小 render を executable spec として固定する。

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: `frontend-test-foundation`
  - `can_run_in_parallel_with`: `backend-test-foundation`
  - `blocked_by`: `active-plan-created`
  - `completion_signal`: `npm --prefix frontend run test` が成功する。
  - `group_id`: `backend-test-foundation`
  - `can_run_in_parallel_with`: `frontend-test-foundation`
  - `blocked_by`: `active-plan-created`
  - `completion_signal`: backend test shell が成功する。
  - `group_id`: `system-and-harness`
  - `can_run_in_parallel_with`:
  - `blocked_by`: `frontend-test-foundation`, `backend-test-foundation`
  - `completion_signal`: `python3 scripts/harness/run.py --suite all` が成功する。
- `tasks`:
  - `task_id`: `frontend-vitest-foundation`
  - `owned_scope`: `frontend/`
  - `depends_on`: `active-plan-created`
  - `parallel_group`: `frontend-test-foundation`
  - `required_reading`: `docs/tech-selection.md`, `frontend/vite.config.ts`
  - `validation_commands`: `npm --prefix frontend run test`
  - `task_id`: `backend-go-test-foundation`
  - `owned_scope`: `internal/`, `scripts/test/`
  - `depends_on`: `active-plan-created`
  - `parallel_group`: `backend-test-foundation`
  - `required_reading`: `docs/tech-selection.md`, `scripts/lint/run-go-backend-lint.sh`
  - `validation_commands`: `sh ./scripts/test/run-go-backend-test.sh`
  - `task_id`: `playwright-and-harness-integration`
  - `owned_scope`: `package.json`, `playwright.config.ts`, `scripts/harness/`, `tests/system/`
  - `depends_on`: `frontend-vitest-foundation`, `backend-go-test-foundation`
  - `parallel_group`: `system-and-harness`
  - `required_reading`: `wails.json`, `scripts/harness/README.md`
  - `validation_commands`: `npm run test:system`, `python3 scripts/harness/run.py --suite all`

## 受け入れ確認

- frontend sample test が成功する。
- backend sample test が成功する。
- `wails dev -browser` を起動する system test が成功する。
- harness から frontend/backend/system test を実行できる。

## 必要な証跡

<!-- Required Evidence -->

- `npm --prefix frontend run test`
- `sh ./scripts/test/run-go-backend-test.sh`
- `npm run test:system`
- `python3 scripts/harness/run.py --suite execution`
- `python3 scripts/harness/run.py --suite all`

## 機能要件 HITL 状態

- N/A

## 機能要件 承認記録

- N/A

## 詳細設計 HITL 状態

- N/A

## 詳細設計 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第1.6段階で作った UI モック working copy は、完了前に `docs/mocks/<page-id>/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/<topic-id>.md` へ移す。

## 結果

<!-- Outcome -->

- 導入済み:
  - root `package.json` に `test:frontend`、`test:backend`、`test:system`、`test:system:install` を追加した。
  - `frontend/` に `Vitest + @testing-library/svelte + jsdom` の実行設定と sample UI test を追加した。
  - backend test 用に `scripts/test/run-go-backend-test.sh` と `internal/controller/wails/app_controller_test.go` を追加した。
  - system test 用に `playwright.config.ts`、`tests/system/app-shell.spec.ts`、`scripts/test/run-system-test.sh` を追加した。
  - harness に `frontend-test`、`backend-test`、`system-test` を追加し、`execution` と `all` へ test suite を統合した。
- 証跡:
  - `npm --prefix frontend run test`
  - `sh ./scripts/test/run-go-backend-test.sh`
  - `npm run test:system`
  - `python3 scripts/harness/run.py --suite execution`
  - `python3 scripts/harness/run.py --suite all`
