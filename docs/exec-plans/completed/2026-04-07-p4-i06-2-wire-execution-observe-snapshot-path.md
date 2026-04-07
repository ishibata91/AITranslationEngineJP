- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement Phase 4 integ task `P4-I06-2` by wiring provider-neutral execution observation snapshots from the Tauri command path through the gateway and composition root into the existing execution-observe UI.
- task_id: P4-I06-2
- task_catalog_ref: tasks/phase-4/tasks/P4-I06-2.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I06-2` from `tasks/phase-4/tasks/P4-I06-2.yaml`.
- Replace the placeholder loader in `src/main.ts` with an integrated Tauri-backed observation snapshot path.
- Keep execution-control redesign, control action UI, and writer output observation out of scope.

## Decision Basis

- `tasks/phase-4/tasks/P4-I06-2.yaml` fixes owned scope to `src/gateway/tauri/execution-observe/`, `src-tauri/src/gateway/commands.rs`, `src-tauri/src/lib.rs`, and `src/main.ts`.
- Existing `P4-I06` work already defines the execution-observe frontend slice, so this task should reuse that slice instead of inventing a parallel UI contract.
- Detailed repo facts, shared wiring constraints, test anchors, and 4humans sync candidates are delegated to downstream skills.

## Owned Scope

- `src/gateway/tauri/execution-observe/`
- `src-tauri/src/gateway/commands.rs`
- `src-tauri/src/lib.rs`
- `src/main.ts`

## Out Of Scope

- execution-control contract redesign
- pause, resume, retry, or cancel UI behavior
- writer output observation
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-I06`
- `P4-C02`
- `P4-V01`
- existing execution-observe frontend slice must already expose provider-neutral snapshot consumption points

## Parallel Safety Notes

- Shared frontend files outside owned scope should remain additive wiring only.
- Shared backend registration should expose the existing provider-neutral snapshot contract without introducing adapter-specific payloads.
- If touched paths extend beyond the task catalog scope, justify the minimum additional shared entrypoints before implementation.

## UI

- `src/main.ts` は composition root として、placeholder の `loadSnapshot` を Tauri-backed loader 注入へ置き換えることだけを担う。`ExecutionObserveScreen` / `ExecutionObserveView` の section 構成、`Refresh` 導線、read-only 境界は `P4-I06` のまま維持する。
- execution-observe UI が受け取る入力は既存の `ExecutionObserveSnapshot` 1 形だけに固定し、Tauri command 名、Rust DTO 名、adapter 固有 field は `src/gateway/tauri/execution-observe/` より先へ漏らさない。
- 初回表示、再読込、失敗時の見え方は既存 usecase の `loading` / `error` / confirmed snapshot 振る舞いをそのまま使う。control action、writer output、placeholder 専用 UI は追加しない。

## Scenario

- アプリ起動時に `src/main.ts` が execution-observe 用の Tauri loader を組み立てて usecase へ注入し、screen mount の `initialize()` から frontend gateway -> Tauri command -> backend snapshot DTO -> `ExecutionObserveSnapshot` の順で 1 回読み込む。既存 execution-observe UI はその結果をそのまま描画する。
- `Running` / `Paused` / `Retrying` / `RecoverableFailed` / `Failed` / `Canceled` / `Completed` の全状態は同一の snapshot 読み出し経路で届き、UI は `P4-I06` で定義済みの summary、failure、phase timeline、phase runs、translation progress、selected unit、footer metadata を provider-neutral に表示する。control action や writer 観測には分岐しない。
- 初回読み込み失敗時は既存の observe error 表示だけを出し、確定済み snapshot がある再読込失敗時は最後の snapshot を残したまま error だけ更新する。bootstrap 専用の失敗モードや別 retry path は増やさない。

## Logic

- `src/gateway/tauri/execution-observe/index.ts` をこの task の単一 frontend adapter とし、Tauri command `get_execution_observe_snapshot` を呼んで `Promise<ExecutionObserveSnapshot>` を返す。transport から UI read model への整形はこの gateway 内に閉じ、`src/main.ts` には loader 注入だけを残す。
- `src-tauri/src/gateway/commands.rs` と `src-tauri/src/lib.rs` は同名 command を追加登録し、provider-neutral な observation snapshot DTO を返す。`controlState` と `failure` は既存 `ExecutionControlStateDto` / `ExecutionControlFailureDto` を再利用し、残りの `summary` / `phaseTimeline` / `phaseRuns` / `translationProgress` / `selectedUnit` / `footerMetadata` を 1 payload にまとめて UI の stable read model へ橋渡しする。
- 共有 wiring は additive 最小に留める。`src/main.ts` では placeholder loader の差し替え以外の責務を足さず、backend 側でも observation read command だけを増やす。execution-control の action command、writer output 観測、UI 向けの第 2 snapshot shape は導入しない。

## Implementation Plan

- ordered_scope:
  1. Distill the current execution-observe frontend slice, placeholder loader wiring, and Tauri command registration path.
  2. Fix the task-local design for the integrated snapshot path across Tauri command, gateway adapter, and composition root.
  3. Define failing tests, fixtures, acceptance checks, and validation commands before implementation.
  4. Implement the owned integration scope and only the minimum shared wiring required for the existing UI to consume the snapshot.
- owned_scope:
  - `src/gateway/tauri/execution-observe/`
  - `src-tauri/src/gateway/commands.rs`
  - `src-tauri/src/lib.rs`
  - `src/main.ts`
- required_reading:
  - `tasks/phase-4/tasks/P4-I06-2.yaml`
  - `tasks/phase-4/tasks/P4-I06.yaml`
  - `docs/exec-plans/active/2026-04-07-p4-i06-progress-and-failure-observability-ui.md`
- validation_commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- execution-observe UI loads provider-neutral progress and failure snapshots through the integrated Tauri path
- `src/main.ts` no longer depends on a placeholder execution-observe loader
- control-action and writer-output responsibilities remain outside this task

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- task-local tests and fixtures defined by `architecting-tests`
- owned lint suite defined by downstream planning
- `sonar-scanner`
- Sonar MCP open issue query for project `ishibata91_AITranslationEngineJP` scoped to touched paths
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  Update only if the integrated snapshot path reveals a remaining quality gap or closes an existing one.
- `4humans/tech-debt-tracker.md`
  Update only if placeholder removal leaves a known follow-up or resolves a tracked debt item.
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  Update the relevant structure diagram if the execution-observe gateway or composition root dependency graph changes materially.
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  Update the relevant process diagram if the execution-observe snapshot flow changes materially for human review.
- `4humans/diagrams/overview-manifest.json`
  Update only if a new detail diagram is added.

## Outcome

- Completed on 2026-04-07.
- `src/main.ts` now injects the Tauri-backed execution-observe loader instead of the placeholder path.
- `src/gateway/tauri/execution-observe/` calls `get_execution_observe_snapshot`, and the backend command path returns a provider-neutral observation snapshot built from the provider-failure-retry acceptance fixture while reusing existing execution-control DTOs.
- Execution-observe validation now asserts integrated command behavior, 4humans diagrams were synced for the new stable path, Sonar reported zero open issues on the touched owned scope, and `python3 scripts/harness/run.py --suite all` passed.
