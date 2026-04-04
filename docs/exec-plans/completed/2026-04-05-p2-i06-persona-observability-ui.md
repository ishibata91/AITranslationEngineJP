- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-I06` from `tasks/phase-2/tasks/P2-I06.yaml` by adding a dedicated persona observation UI path on top of the existing master persona foundation boundary without widening into dictionary observation or persona generation policy work.
- task_id: P2-I06
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-I06`.
- Expose master persona foundation data through a stable UI observation path that stays split from dictionary observation and generation-policy work.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `docs/screen-design/wireframes/foundation-data.md`
- `tasks/phase-2/phase.yaml`
- `tasks/phase-2/tasks/P2-I06.yaml`
- `docs/exec-plans/completed/2026-04-04-p2-i03-master-persona-builder.md`
- `docs/exec-plans/completed/2026-04-05-p2-i05-dictionary-observability-ui.md`

## Owned Scope

- `src/ui/screens/persona-observe/`
- `src/ui/views/persona-observe/`
- `src/application/usecases/persona-observe/`

## Out Of Scope

- dictionary observation
- persona generation policy
- backend master persona build contract redesign

## Dependencies / Blockers

- `P2-I06` depends on `P2-I03`.
- The UI must consume the stable persona foundation boundary already fixed upstream instead of inventing a frontend-only policy contract.

## Parallel Safety Notes

- This task remains parallel-safe with the rest of `P2-B2` only while edits stay inside persona observation screen, view, and usecase ownership plus minimal shell composition.
- Shared foundation-data navigation, dictionary panes, rebuild actions, and persona generation policy controls must not be absorbed into this slice.

## UI

- `src/ui/screens/persona-observe/` と `src/ui/views/persona-observe/` に、`AppShell` 配下で常に到達できる dedicated persona observation panel を追加する。browser routing は導入せず、shell への追加は `dictionary-observe` と同じ additive composition に限定し、SSR test で安定 path を固定する。
- 画面は `foundation-data` wireframe の persona 側だけを read-only observation 用に切り出し、`personaName` 入力、observe / refresh action、dataset metadata、entry 一覧、選択中 entry detail の 5 要素で構成する。dictionary pane、`New Persona` / `Rebuild` action、entry editor affordance、job-local persona 情報はこの task に含めない。
- dataset metadata は upstream read contract に存在する `personaName` と `sourceType` だけを表示する。wireframe 上の build time は既存 backend contract にないため、この task では追加せず、観測 UI からも表示しない。
- 表示状態は 1 レイアウト内で固定し、初回観測前の案内 empty state、初回実行中 loading、成功後の loaded、retry 可能な failure を扱う。成功結果で `entries` が 0 件でも empty screen に戻さず、metadata を維持したまま一覧と detail で「entry なし」を示す。

## Scenario

- screen mount 時の `initialize` は store 初期化だけを行い、master persona read は自動起動しない。`personaName` が未指定の request は送らず、observe action でだけ query を確定する。
- 利用者は dedicated panel で観測したい `personaName` を入力し、明示的な observe action で `MasterPersonaReadRequestDto { personaName }` 相当の read を起動する。成功時は返却された `personaName`、`sourceType`、`entries` をそのまま表示し、entry が 1 件以上あれば先頭 entry を自動選択する。
- 利用者が別 entry を選ぶと、detail pane はその entry の `npcFormId`、`npcName`、`race`、`sex`、`voice`、`personaText` を read-only で表示する。job-local persona の `jobId` や generation policy 情報は同じ NPC 属性列を共有していても混在させない。
- `refresh` と `retry` は直前に送信した同一 `personaName` を再実行する。再読込成功時は現在選択中の entry index がまだ存在する場合は維持し、存在しない場合は先頭 entry、entry 不在なら `null` へ戻す。再読込失敗時は直前の成功結果を残したまま generic error message と retry affordance だけを重ねて表示する。

## Logic

- `src/application/usecases/persona-observe/` は frontend-owned observation contract を持つが、語彙は upstream master persona read DTO に合わせて `personaName` / `sourceType` / `entries` / `npcFormId` / `npcName` / `race` / `sex` / `voice` / `personaText` をそのまま使う。transport 実装は持ち込まず、`(request) => Promise<result>` 形式の injected executor または gateway port に依存する。
- state shape は `feature-screen` 系の `data` / `loading` / `error` / `selection` / `filters` を踏襲し、`filters` に current `personaName` input と `lastSubmittedRequest` を保持する。observe は keystroke ごとに自動発火せず、明示 action でだけ request を確定する。
- selection は `npcFormId` や `npcName` ではなく entry index で保持する。これにより upstream が返す entry 順序をそのまま使い、同名 NPC や将来の非一意 display label があっても一覧選択と detail 表示を 1 対 1 で安定させる。
- request 生成では `personaName` の alias 解決、case-fold、sort、entry 側の client-side reshape を行わない。空文字 request だけは screen-local validation で防ぎ、error mapping は transport 非依存の user-facing message に固定して、Tauri / filesystem / SQL の詳細を Svelte view へ漏らさない。
- stable dedicated path を成立させるため、実装時は `AppShell` への最小 prop 追加と screen mount だけを隣接変更として許容する。一方で `src/gateway/tauri/`、`src/main.ts`、backend master persona contract の再設計は `P2-G01` 以降または別 task の責務に残し、この task では screen/usecase がそのまま差し替え可能な transport-ready 境界までを固定する。

## Implementation Plan

### ordered_scope

1. Persona observe usecase (`src/application/usecases/persona-observe/`)

- upstream master persona read DTO と同じ `personaName` / `sourceType` / `entries` / `npcFormId` / `npcName` / `race` / `sex` / `voice` / `personaText` をそのまま使う frontend-owned state、store、usecase 契約を追加する。保持対象は current `personaName` input と `lastSubmittedRequest` に限定する。
- `initialize()` は transport-free の初期化だけに留め、read は observe / refresh / retry の明示 action でのみ実行する。selection は entry index で再調停し、再読込後も index が残る限り current selection を維持する。

2. Persona observe view and screen (`src/ui/views/persona-observe/`, `src/ui/screens/persona-observe/`)

- `personaName` 入力、observe / refresh / retry affordance、dataset metadata、entry 一覧、選択 entry detail を持つ read-only dedicated panel を追加する。dictionary pane、`New Persona` / `Rebuild`、editor affordance は含めない。
- screen は initialize、`personaName` 更新、observe、refresh、retry、selection event を usecase へ中継するだけに留める。view は 1 レイアウト内で empty / loading / loaded / failure を描画し、retry 可能な failure でも直前成功結果と metadata を残して表示する。

3. Stable shell composition (`src/ui/app-shell/AppShell.svelte`, `src/App.svelte` only if the current root must thread new props for compile-safe composition)

- `dictionary-observe` と同じ additive composition pattern で persona observation panel を shell 配下へ追加し、stable dedicated path を SSR test で固定する。
- この task では browser routing、`src/main.ts`、`src/gateway/tauri/`、dictionary observation の拡張、backend master persona contract の再設計を行わない。

### owned_scope

- `src/application/usecases/persona-observe/`
- `src/ui/views/persona-observe/`
- `src/ui/screens/persona-observe/`
- `src/ui/app-shell/AppShell.svelte`
- `src/App.svelte` only if additive shell composition needs a thin prop pass-through at the current root
- `src/main.ts`、`src/gateway/tauri/`、dictionary observation、persona generation policy、backend master persona contract redesign は `P2-I06` の scope 外に残す

### required_reading

- `docs/exec-plans/active/2026-04-05-p2-i06-persona-observability-ui.md` の `UI` / `Scenario` / `Logic`
- `docs/exec-plans/completed/2026-04-05-p2-i05-dictionary-observability-ui.md`
- `docs/screen-design/wireframes/foundation-data.md`
- `src/application/ports/input/feature-screen/index.ts`
- `src/application/usecases/feature-screen/index.ts`
- `src/application/usecases/dictionary-observe/index.ts`
- `src/ui/screens/dictionary-observe/DictionaryObserveScreen.svelte`
- `src/ui/views/dictionary-observe/DictionaryObserveView.svelte`
- `src/ui/app-shell/AppShell.svelte`
- `src/App.svelte`
- `src/application/usecases/job-list/index.test.ts`
- `src/ui/screens/job-list/index.test.ts`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite frontend-lint`
- `npm run test -- src/application/usecases/persona-observe src/ui/screens/persona-observe`
- `npm run build`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- the UI can observe persona foundation state only after an explicit `observe` action and keeps `lastSubmittedRequest` plus entry-index selection stable across `refresh` or `retry`
- the UI can observe persona foundation data through a stable shell path, and the shell or App root still render safely before persona transport wiring is provided

## Required Evidence

- Active plan updated with task-local UI, scenario, and logic decisions.
- Tests covering the persona observation usecase request lifecycle, screen rendering for empty or loading or loaded or retryable-failure states, and stable shell or App-root composition states.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md` update not required unless repository-level quality posture changes
- `4humans/tech-debt-tracker.md` update not required unless new unresolved debt is introduced
- `4humans/diagrams/structures/frontend-persona-observe-slice-class-diagram.d2`
- `4humans/diagrams/structures/frontend-persona-observe-slice-class-diagram.svg`
- `4humans/diagrams/structures/frontend-structure-overview.d2`
- `4humans/diagrams/structures/frontend-structure-overview.svg`
- `4humans/diagrams/processes/frontend-persona-observe-sequence-diagram.d2`
- `4humans/diagrams/processes/frontend-persona-observe-sequence-diagram.svg`
- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/processes-overview-robustness.svg`
- `4humans/diagrams/overview-manifest.json`

## Outcome

- Added `src/application/usecases/persona-observe/index.ts` with a transport-free persona observation store and usecase that keeps initialization side-effect free, runs explicit observe / refresh / retry requests only, stores `personaName` plus `lastSubmittedRequest`, and preserves selection by entry index.
- Added `src/ui/views/persona-observe/PersonaObserveView.svelte`, `src/ui/views/persona-observe/index.ts`, `src/ui/screens/persona-observe/PersonaObserveScreen.svelte`, and `src/ui/screens/persona-observe/index.ts`, so the dedicated persona observation panel now renders input, observe / refresh / retry actions, dataset metadata, entry list, selected entry detail, empty guidance, loading, and retryable failure states.
- Updated `src/ui/app-shell/AppShell.svelte` and `src/App.svelte` additively so persona observation composes through the shell path only when dependencies are provided, without adding browser routing, gateway wiring, or local noop transport.
- Added `4humans` D2 sync with `frontend-persona-observe-slice-class-diagram.d2/.svg`, `frontend-persona-observe-sequence-diagram.d2/.svg`, and linked updates to `frontend-structure-overview.d2/.svg`, `processes-overview-robustness.d2/.svg`, and `overview-manifest.json`.
- Reroute after single-pass review was resolved by adding the required `4humans` sync, and final harness reroutes were resolved by extending SSR AppShell test stubs in `src/ui/screens/job-list/index.test.ts` and `src/ui/screens/dictionary-observe/index.test.ts` to replace the new persona-observe screen import during server-render compilation.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite frontend-lint`, `npm run test -- src/application/usecases/persona-observe src/ui/screens/persona-observe`, `npm run build`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/application/usecases/persona-observe src/ui/screens/persona-observe src/ui/views/persona-observe src/ui/app-shell/AppShell.svelte src/App.svelte` reporting `openIssueCount = 0`, and `python3 scripts/harness/run.py --suite all`.
