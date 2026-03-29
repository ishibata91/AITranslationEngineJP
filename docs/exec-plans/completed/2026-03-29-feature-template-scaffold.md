# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: frontend feature template scaffold for screen / store / usecase / gateway

## Request Summary

- Add reusable frontend feature templates so future business features can be added without redesigning screen / store / usecase structure from zero.
- Keep the existing architecture contract between UI layer, application layer, and gateway layer intact.
- Show how loading, error, data, and selection state should be modeled in reusable code.

## Decision Basis

- User requested reusable templates under the existing frontend and backend skeleton.
- `docs/architecture.md` defines `Screen UseCase`, `Screen Store`, `Presenter / View`, and `Gateway` as distinct UI-layer responsibilities.
- `docs/` source-of-truth must not be edited in this task.

## UI

- `src/ui/screens/<feature>/` に screen entry を置き、`<Feature>Screen.svelte` は view model 取得と UI event の binding だけを担う。
- `src/ui/views/<feature>/` に表示専用の `<Feature>View.svelte` を置き、props と event dispatch だけを扱う。
- `src/ui/stores/<feature>/` に feature store root を置き、`loading`、`error`、`data`、`selection`、必要なら `filters` を 1 つの screen state として保持する。
- App Shell から screen を差し替え可能な構成を保ち、bootstrap 固有の表示文言や DTO 名を template root に持ち込まない。

## Scenario

- UI event は screen から usecase へ渡し、usecase が store 更新順序を決める。
- 初期表示では `initialize` 系 entry で `loading` を開始し、gateway 取得成功時は `data` を反映して `error` を空にし、失敗時は既存 data / selection を壊さず `error` を更新する。
- ユーザー選択は `select` 系 entry で store に閉じ込め、再取得や refresh 後も selection の再適用可否を usecase で判断できる形にする。
- feature 追加時は `screen -> usecase -> gateway/store` の呼び順を踏襲し、view から gateway や store へ直接触れない。

## Logic

- `src/application/usecases/<feature>/` に feature usecase root を置き、UI が使う input contract、store contract、gateway port contract を明示する。
- usecase は orchestration に専念し、request DTO の組み立て、`loading / error / data / selection` の state transition、refresh / retry / select の公開 entry を持つ。
- `src/application/ports/input/` には feature 公開面の再 export root を置き、screen は usecase 実装詳細ではなく input port から参照する。
- `src/application/ports/gateway/` には feature ごとの gateway interface を置き、`src/gateway/tauri/` 側 adapter がその interface を実装する。
- feature store は UI state の setter / getter に留め、gateway call や usecase call を持たない。gateway adapter は Tauri invoke 専用に留め、画面状態や選択再適用の判断を持たない。

## Implementation Plan

- Keep existing `feature-screen` primitives as the reusable base for store / usecase / gateway orchestration.
- Add a canonical `feature-template` copy source across `shared/contracts`, `application/ports`, `application/usecases`, `ui/stores`, `ui/screens`, `ui/views`, and `gateway/tauri`.
- Preserve `bootstrap-status` as the mounted reference feature without widening runtime scope.
- Add minimal tests for generic gateway transport and template usecase orchestration.
- Run structure harness, targeted Vitest scope, ESLint, `sonar-scanner`, and single-pass implementation review.

## Acceptance Checks

- Structure harness passes.
- Added template files show standard locations and responsibility split for screen / store / usecase / gateway.
- Added or updated tests cover the reusable state and orchestration entry points in the changed scope.
- Single-pass implementation review returns `pass`.

## Required Evidence

- Validation output for `powershell -File scripts/harness/run.ps1 -Suite structure`
- Validation output for the minimal test / lint scope touched by the change

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `feature-template` as a compile-backed copy source for future business features.
- Kept the responsibility boundary explicit as `screen -> usecase -> gateway/store`, with view remaining display-only.
- Added gateway transport coverage and template usecase coverage.
- Validation passed for structure harness, targeted Vitest, ESLint, and `sonar-scanner`.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not updated because this change introduced no new unresolved debt or quality posture change that needed human tracking.
