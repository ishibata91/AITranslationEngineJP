# 実装スコープ固定

- `task_id`: `backend-frontend-unit-test-abstraction-validation`
- `task_mode`: `refactor`
- `design_review_status`: `not_run`
- `hitl_status`: `approved`
- `summary`: frontend unit test のうち `App.test.ts` を abstraction-first へ戻し、concrete controller factory 依存と Wails runtime event 詳細依存を App-level test から外す。

## 共通ルール

- `App.test.ts` は `CreateMasterDictionaryScreenController` と `MasterDictionaryScreenControllerContract` だけを境界として扱う。
- `@controller/master-dictionary` と `window.runtime` の詳細は App-level test へ持ち込まない。
- runtime event 名、`EventsOnMultiple` 呼び出し形、payload parse は controller/runtime focused test へ寄せる。
- coverage gate を維持するため、必要なら controller / factory の product seam に対する lower-level unit test を追加してよい。
- product code は変更しない。
- 既存の UI シナリオ観点は維持し、差分は test support と test file に閉じる。

## 実装 handoff 一覧

### `frontend-app-contract-only-tests`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/test/setup.ts`
- `depends_on`: `none`
- `validation_commands`:
  - `cd frontend && npm test -- App.test.ts`
  - `cd frontend && npm test`
- `completion_signal`: `App.test.ts` が concrete factory helper を使わず、contract fake から画面状態と action 呼び出しを制御できる。
- `notes`:
  - fake は `MasterDictionaryScreenControllerContract` を実装し、`getViewModel`、`subscribe`、主要 action の呼び出し記録、view model 更新 push を持つ。
  - `setup.ts` は `jest-dom` 初期化だけへ縮退させるか、contract fake の export だけに限定する。
  - XML import 完了待ちや progress 表示の App-level 検証は fake の `startImport` と view model push で表現し、runtime bridge mock を使わない。

### `frontend-runtime-and-controller-focused-tests`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts`
- `depends_on`: `frontend-app-contract-only-tests`
- `validation_commands`:
  - `cd frontend && npm test -- master-dictionary-runtime-event-adapter`
  - `cd frontend && npm test -- master-dictionary-screen-controller`
  - `cd frontend && npm test`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: progress/completed event 名、`EventsOnMultiple` 登録、detach、payload forwarding の検証が App render を介さず完結し、controller / factory seam の主要分岐を lower-level unit test で補完して coverage gate を回復する。
- `notes`:
  - test helper 自体のテストは増やさない。
  - `maxCallbacks = -1` や invalid payload の扱いは runtime adapter test でのみ確認する。
  - controller / factory test は fake gateway / runtime / listener を使い、product seam の責務だけを検証する。
