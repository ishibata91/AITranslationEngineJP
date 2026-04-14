# Work Plan

- workflow: orchestrate
- status: completed
- lane_owner: orchestrate
- scope: frontend-dip-for-unit-tests
- task_id: frontend-dip-for-unit-tests
- task_catalog_ref: N/A
- parent_phase: implementation-lane

## Request Summary

- `View` より先の frontend layer を差し替え可能な境界へ寄せる。
- frontend の mixed directory をやめ、layer-based directory へ寄せる。
- import 制約と harness を維持したまま coverage target を満たす。

## Decision Basis

- 現在の frontend 実装対象では `master-dictionary` が主要変更面だった。
- `ui -> {ui, application}` と `ui` から `controller` 直接 import 禁止を守るには、Svelte view を薄く保ち、application contract と controller implementation を分ける必要があった。
- user の追加判断により、feature 塊ではなく pure な layer directory を最終形として採用した。

## Task Mode

- `task_mode`: refactor
- `goal`: `View` は Svelte のまま維持し、`controller`、`runtime`、`usecase`、`presenter`、`store`、`contract` を layer ownership へ再配置する。
- `constraints`: docs 正本は human 先行でのみ更新する。`ui` は `controller` を直接 import しない。DI container は導入しない。
- `close_conditions`: review が pass を返すこと。`python3 scripts/harness/run.py --suite all` を通すこと。coverage gate を満たすこと。

## Facts

- 旧 `frontend/src/ui/screens/master-dictionary/` 配下の TypeScript は削除され、`MasterDictionaryPage.svelte` のみが残った。
- application layer は `contract/master-dictionary`、`usecase/master-dictionary`、`presenter/master-dictionary`、`store/master-dictionary` に分離された。
- controller layer は `controller/master-dictionary` と `controller/runtime/master-dictionary` に分離された。
- `main.ts -> App.svelte -> AppShell.svelte -> MasterDictionaryPage.svelte` の factory prop threading へ置き換えた。
- boundary plugin は layer-based public root を追加しつつ、test source 全体の boundary 免責は入れなかった。

## Functional Requirements

- `summary`:
  - `View` は Svelte component に残す。
  - `View` より先は application-owned contract と controller-owned implementation に分ける。
  - `master-dictionary` の TS 資産を layer-based directory に再配置する。
- `in_scope`:
  - `frontend/src/main.ts`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`
  - `frontend/src/test/setup.ts`
  - `frontend/src/application/contract/master-dictionary/`
  - `frontend/src/application/usecase/master-dictionary/`
  - `frontend/src/application/presenter/master-dictionary/`
  - `frontend/src/application/store/master-dictionary/`
  - `frontend/src/controller/master-dictionary/`
  - `frontend/src/controller/runtime/master-dictionary/`
  - `scripts/eslint/repository-boundary-plugin.mjs`
  - `frontend/repository-boundary-plugin.test.mjs`
- `non_functional_requirements`:
  - user-visible behavior、Wails binding 名、runtime event 名、payload 互換を維持する。
  - `ui` は `controller` 実装型を import しない。
  - coverage の主対象は `.ts` seam に置く。
  - boundary plugin の test support 例外は reverse-flow の最小差分に留める。
- `out_of_scope`:
  - `master-dictionary` 以外の screen の一括再編
  - backend 契約変更
  - docs 正本更新
- `open_questions`:
  - なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-14-frontend-dip-for-unit-tests.implementation-scope.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/contract/master-dictionary/index.ts`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-dictionary/index.ts`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-14-frontend-dip-for-unit-tests.implementation-scope.md
- `review_diff_diagrams`: N/A
- `source_diagram_targets`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/frontend/frontend-architecture.d2`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/master-dictionary-management.d2`
- `canonicalization_targets`:
  - `docs/architecture.md`
  - `docs/diagrams/frontend/frontend-architecture.d2`
  - `docs/diagrams/components/frontend/master-dictionary-management.d2`
  - `docs/lint-policy.md`

## Work Brief

- `implementation_target`: frontend
- `accepted_scope`:
  - contract は `frontend/src/application/contract/master-dictionary/` に置く。
  - usecase は `frontend/src/application/usecase/master-dictionary/` に置く。
  - presenter は `frontend/src/application/presenter/master-dictionary/` に置く。
  - store は `frontend/src/application/store/master-dictionary/` に置く。
  - controller は `frontend/src/controller/master-dictionary/` に置く。
  - runtime adapter は `frontend/src/controller/runtime/master-dictionary/` に置く。
  - `MasterDictionaryPage.svelte` は thin view / lifecycle owner に限定する。
- `handoff_targets`: implement -> review(ui-check, implementation-review)
- `validation_commands`:
  - `cd frontend && npm run lint`
  - `cd frontend && npm test`
  - `cd frontend && npm run check`
  - `python3 scripts/harness/run.py --suite all`

## Validation Results

- `frontend`: `npm run lint` pass
- `frontend`: `npm test` pass
- `frontend`: `npm run check` pass
- `repo`: `python3 scripts/harness/run.py --suite all` pass
- `coverage`: frontend statements `73.26%`、lines `73.18%`、functions `77.95%`、branches `60.5%`
- `sonar`: coverage `77.3% >= 70.0%`
- `ui-check`: `http://host.docker.internal:34115/#dashboard` と `#master-dictionary` 表示を確認

## Required Evidence

- layer-based directory への再配置結果
- boundary plugin public root 更新
- `ui` から `controller` 直接 import 不在
- harness / coverage / ui-check の pass

## HITL Status

- `functional_or_design_hitl`: approved
- `design_review_status`: pass
- `implementation_review_status`: pass
- `ui_check_status`: pass
- `approval_record`: 2026-04-14 human request: `View` より先の frontend layer を DIP 化し、feature 塊ではなく layer-based directory へ寄せる。

## Closeout Notes

- `canonicalized_artifacts`: N/A
- docs 正本同期は別 lane のまま残す。
- `favicon.ico` 404 と `vite-plugin-svelte` deprecation warning は本 task の close blocker ではない。

## Outcome

- completed
