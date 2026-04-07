- workflow: impl
- status: planned
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I05` pause, retry, and cancel UI against the stable execution-control contract.
- task_id: P4-I05
- task_catalog_ref: tasks/phase-4/tasks/P4-I05.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I05` from `tasks/phase-4/tasks/P4-I05.yaml`.
- Expose pause, retry, and cancel through a dedicated execution-control UI path without collapsing UI and provider adapter boundaries.

## Decision Basis

- `tasks/phase-4/tasks/P4-I05.yaml` fixes owned scope to execution-control screen, view, and use case paths.
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md` already fixed the stable execution-control contract that this UI must consume.
- Detailed repo facts, design constraints, and required reading are delegated to `distilling-implementation`.

## Owned Scope

- `src/ui/screens/execution-control/`
- `src/ui/views/execution-control/`
- `src/application/usecases/execution-control/`

## Out Of Scope

- provider adapter transport
- writer output UI
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C02`
- execution-control state contract must remain stable and provider-neutral

## Parallel Safety Notes

- Keep changes centered on the owned execution-control UI path and narrow composition touch points only if required for wiring.
- Do not embed provider-specific behavior, transport rules, or writer output behavior in this task.
- If shared entrypoints outside the owned scope must change, justify the smallest necessary additive touch point in downstream design.

## UI

- `src/ui/screens/execution-control/` に `ExecutionControlScreen` を追加し、既存 screen pattern と同様に injected `store` / `usecase` を受け取り、`onMount` で初期化だけを呼ぶ。screen は view の `pause` / `resume` / `retry` / `cancel` event を usecase へそのまま委譲し、provider 判定や state 判定を持たない。
- `src/ui/views/execution-control/` に `ExecutionControlView` を追加し、表示は provider-neutral な `Control State`、任意の failure panel、`Pause` / `Resume` / `Retry` / `Cancel` の action row で閉じる。wireframe の周辺観測情報のうち、phase timeline、progress、provider run id など execution-control 外の観測パネルはこの task では扱わない。
- action row は 4 ボタンを常時描画し、enabled / disabled と busy label だけを切り替える。`Pause` と `Resume` は相互排他で、state に応じて一方だけを活性にする。
- `RecoverableFailed` では failure panel を主表示し、`ExecutionControlFailureDto.category` と `message` を provider-neutral な失敗表示として出す。その他 state では panel は前回 failure がある時だけ補助表示し、failure が無ければ非表示でよい。

## Scenario

- マウント時に current execution-control snapshot を読み込み、ユーザーは現在 state と action availability を確認できる。
- `Running` では `Pause` と `Cancel` を enable、`Resume` と `Retry` を disable にする。
- `Paused` では `Resume` と `Cancel` を enable、`Pause` と `Retry` を disable にする。
- `RecoverableFailed` では failure panel と `Retry` / `Cancel` を主表示し、`Pause` / `Resume` は disable にする。
- `Retrying` では retry 中表示に切り替え、`Cancel` だけを enable にする。
- `Failed` / `Canceled` / `Completed` では read-only とし、4 操作をすべて disable にする。
- いずれかの action 送信中は全ボタンを一時 disable にし、送信中 action だけ busy label に切り替える。送信失敗時は最後に確定していた snapshot を保持したまま error を出し、action availability は確定 state から再計算する。

## Logic

- usecase は shared contract を再設計せず、最後に確定した `ExecutionControlStateDto` と任意の `ExecutionControlFailureDto` を screen state に保持する。view へは `canPause` / `canResume` / `canRetry` / `canCancel`、`pendingAction`、`error` を含む provider-neutral な view state だけを渡す。
- public surface はこの task で最小限 `initialize()`、`pause()`、`resume()`、`retry()`、`cancel()` を固定対象とし、job / run 識別子の受け取り方は caller 側の injected dependency に閉じる。screen / view から job id や provider id を直接扱わせない。
- action availability は `ExecutionControlStateDto` の stable vocabulary だけから導出する。Svelte view に transition matrix や provider rule を書かず、usecase 側で一元計算する。
- action 実行時の optimistic update は `pendingAction` の付与までに留め、`PausedPending` のような UI 専用 state は増やさない。state 遷移の確定は usecase が受け取った contract response に従う。
- action 実行結果の失敗は 2 系統で扱う。state snapshot 自体が `RecoverableFailed` なら failure panel を更新し、command 実行や再読込の失敗なら generic `error` を更新して最後の確定 snapshot を維持する。これにより recoverable runtime failure と transport / observe failure を混同しない。

## Implementation Plan

- ordered_scope:
  1. `src/application/usecases/execution-control/` に provider-neutral な screen model と `pause` / `resume` / `retry` / `cancel` action handler の薄い frontend interface を追加し、既存 DTO / state vocabulary だけで表示可否と action availability を決める。
  2. `src/ui/screens/execution-control/` に execution-control 専用 screen を追加し、usecase 呼び出し、初期表示用 state 受け取り、単一 action 実行中の busy 制御、view event から usecase への委譲を実装する。
  3. `src/ui/views/execution-control/` に presentational view を追加し、現在 state、provider-neutral failure summary、state-aware な `Pause` / `Resume` / `Retry` / `Cancel` controls を表示して event だけを screen へ返す。
  4. mount/export のために owned scope 外の接続が不可避な場合だけ、最小の additive touch point を 1 箇所に限定して wiring する。
- owned_scope:
  - `src/application/usecases/execution-control/`
  - `src/ui/screens/execution-control/`
  - `src/ui/views/execution-control/`
- required_reading:
  - `tasks/phase-4/tasks/P4-I05.yaml`
  - `docs/exec-plans/active/2026-04-06-p4-i05-pause-retry-cancel-ui.md`
  - `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md`
  - `docs/screen-design/wireframes/job-run.md`
  - `src-tauri/src/application/dto/execution_control/mod.rs`
  - `src-tauri/src/domain/execution_control_state/mod.rs`
  - `src/ui/screens/job-list/JobListScreen.svelte`
  - `src/ui/views/job-list/JobListView.svelte`
  - `src/ui/screens/translation-preview/TranslationPreviewScreen.svelte`
  - `src/ui/views/translation-preview/TranslationPreviewView.svelte`
- validation_commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- The UI can trigger pause, retry, and cancel against the stable control contract.
- The execution-control UI path stays provider-neutral.
- mounted execution-control path derives action availability only from stable `ExecutionControlStateDto` vocabulary for `Running`, `Paused`, `RecoverableFailed`, `Retrying`, `Failed`, `Canceled`, and `Completed`
- during a submitted control action, all four controls are temporarily disabled, the pending action is surfaced without inventing a UI-only pseudo-state, and a command failure restores the last confirmed snapshot with a generic error
- failure panel renders provider-neutral `ExecutionControlFailureDto.category` and `message`, with recoverable failure making `Retry` and `Cancel` the only enabled recovery actions
- The mounted execution-control path derives action availability only from the stable `ExecutionControlStateDto` vocabulary for `Running`, `Paused`, `RecoverableFailed`, `Retrying`, `Failed`, `Canceled`, and `Completed`.
- During a submitted control action, all four controls are temporarily disabled, the pending action is surfaced without inventing a UI-only pseudo-state, and a command failure restores the last confirmed snapshot with a generic error.
- The failure panel renders provider-neutral `ExecutionControlFailureDto.category` and `message`, with recoverable failure making `Retry` and `Cancel` the only enabled recovery actions.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `npm test -- src/application/usecases/execution-control/index.test.ts src/ui/screens/execution-control/index.test.ts`
- `python3 scripts/harness/run.py --suite frontend-lint`
- `sonar-scanner`
- Sonar MCP `search_sonar_issues_in_projects` for project `ishibata91_AITranslationEngineJP` scoped to `src/application/usecases/execution-control`, `src/ui/screens/execution-control`, and `src/ui/views/execution-control`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。今回の変更は frontend の action path 追加であり、既存の品質評価軸や未充足項目の整理対象は増やしていない。
- `4humans/tech-debt-tracker.md`
  更新なし。新規の恒久負債項目は追加しなかった。
- `4humans/diagrams/structures/frontend-execution-control-slice-class-diagram.d2`
- `4humans/diagrams/structures/frontend-execution-control-slice-class-diagram.svg`
- `4humans/diagrams/structures/frontend-structure-overview.d2`
- `4humans/diagrams/structures/frontend-structure-overview.svg`
- `4humans/diagrams/processes/frontend-execution-control-sequence-diagram.d2`
- `4humans/diagrams/processes/frontend-execution-control-sequence-diagram.svg`
- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/processes-overview-robustness.svg`
- `4humans/diagrams/overview-manifest.json`

## Outcome

- Added `src/application/usecases/execution-control/` with a provider-neutral screen store and usecase that derive action availability from stable execution-control state vocabulary, preserve the last confirmed snapshot during command failure, and surface generic initialization errors.
- Added `src/ui/screens/execution-control/` and `src/ui/views/execution-control/` with one dedicated execution-control UI path for pause, resume, retry, and cancel plus provider-neutral failure rendering.
- Added targeted tests for the new usecase and screen/view roots, and updated `4humans` structure/process diagrams plus overview manifest for the new execution-control slice.
- Validation passed for `python3 scripts/harness/run.py --suite structure`, `npm test -- src/application/usecases/execution-control/index.test.ts src/ui/screens/execution-control/index.test.ts`, `python3 scripts/harness/run.py --suite frontend-lint`, `sonar-scanner`, `d2 validate`, and `d2 -t 201`.
- `python3 scripts/harness/run.py --suite all` still fails in repo-wide Rust doctests because `rustdoc` cannot resolve crates such as `sqlx`, `serde`, `tauri`, `reqwest`, and `tokio` under `src-tauri`. This blocker is outside the owned frontend scope.
