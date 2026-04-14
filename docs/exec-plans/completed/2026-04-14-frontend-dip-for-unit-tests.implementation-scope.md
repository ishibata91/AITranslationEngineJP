# 実装スコープ固定

- `task_id`: `frontend-dip-for-unit-tests`
- `task_mode`: `refactor`
- `design_review_status`: `pass`
- `hitl_status`: `approved`
- `summary`: 現在存在する frontend 実装対象について、`View` を Svelte に残したまま TypeScript 資産を pure な layer directory へ再配置した。

## 共通ルール

- `View` は `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` に残す。
- `ui` は `controller` を直接 import しない。
- `application` は `controller` と `ui` を import しない。
- coverage の主対象は `.ts` seam に置く。
- docs 正本更新はこの artifact に含めない。
- public root は feature 塊ではなく layer 単位で持つ。

## 実装結果

### `frontend-master-dictionary-contract-freeze`

- `implementation_target`: `frontend`
- `owned_scope`: `frontend/src/application/contract/master-dictionary/`
- `completion_signal`: `MasterDictionaryScreenControllerContract` と `CreateMasterDictionaryScreenController` を application-owned contract として固定した。
- `result`: `index.ts`、`master-dictionary-screen-contract.ts`、`master-dictionary-screen-types.ts`、`master-dictionary-screen-constants.ts` を配置した。

### `frontend-master-dictionary-application-relocation`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/application/usecase/master-dictionary/`
  - `frontend/src/application/presenter/master-dictionary/`
  - `frontend/src/application/store/master-dictionary/`
- `completion_signal`: usecase、presenter、store と既存 TS test を layer 別 ownership に移した。
- `result`: 旧 `ui/screens/master-dictionary` 配下の TS は削除し、layer 配下へ移設した。

### `frontend-master-dictionary-controller-composition`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/controller/master-dictionary/`
  - `frontend/src/controller/runtime/master-dictionary/`
- `completion_signal`: controller と runtime adapter を layer 別に分離し、class 本体から concrete new を外した。
- `result`: controller factory は `controller/master-dictionary`、runtime adapter は `controller/runtime/master-dictionary` に置いた。

### `frontend-master-dictionary-ui-threading`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/main.ts`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`
- `completion_signal`: `main.ts -> App.svelte -> AppShell.svelte -> MasterDictionaryPage.svelte` の factory threading のみで screen が起動する。
- `result`: page self-wire と gateway prop 差し替えを廃止した。

### `frontend-master-dictionary-test-boundary`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/test/setup.ts`
  - `scripts/eslint/repository-boundary-plugin.mjs`
  - `frontend/repository-boundary-plugin.test.mjs`
- `completion_signal`: test source 全体の boundary 免責を入れず、reverse-flow の最小例外だけを維持した。
- `result`: `App.test.ts` は `src/test/setup.ts` 経由で factory を受けるように直し、`ui -> controller` 直 import を避けた。

## 検証

- `cd frontend && npm run lint`: pass
- `cd frontend && npm test`: pass
- `cd frontend && npm run check`: pass
- `python3 scripts/harness/run.py --suite all`: pass
- `coverage`: frontend statements `73.26%`、lines `73.18%`

## 補足

- docs 正本影響先は `docs/architecture.md`、`docs/diagrams/frontend/frontend-architecture.d2`、`docs/diagrams/components/frontend/master-dictionary-management.d2`、`docs/lint-policy.md` である。
- docs 正本の同期は human-first の別 lane に残す。
