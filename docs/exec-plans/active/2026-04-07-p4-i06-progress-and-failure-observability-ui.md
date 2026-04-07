- workflow: impl
- status: planned
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I06` progress and failure observability UI against the stable execution observation contract.
- task_id: P4-I06
- task_catalog_ref: tasks/phase-4/tasks/P4-I06.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I06` from `tasks/phase-4/tasks/P4-I06.yaml`.
- Expose execution progress and failure reasons through a dedicated observation UI path without relying on adapter internals.

## Decision Basis

- `tasks/phase-4/tasks/P4-I06.yaml` fixes owned scope to execution-observe screen, view, and use case paths.
- `P4-I06` depends on `P4-C02` and `P4-V01`; detailed repo facts, constraints, and required reading are delegated to `distilling-implementation`.
- Existing phase-4 work already separated execution-control from observation, so this task must keep that UI boundary intact.

## Owned Scope

- `src/ui/screens/execution-observe/`
- `src/ui/views/execution-observe/`
- `src/application/usecases/execution-observe/`

## Out Of Scope

- pause or retry action handling
- writer output observation
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C02`
- `P4-V01`
- stable execution-control state contract and provider failure acceptance anchor must remain provider-neutral

## Parallel Safety Notes

- Keep changes centered on the owned execution-observe UI path and only use additive shared wiring if mounting requires it.
- Do not merge control actions into the observation path.
- If shared entrypoints outside the owned scope must change, justify the smallest necessary touch point in downstream design.

## UI

- `src/ui/screens/execution-observe/` に `ExecutionObserveScreen` を追加し、既存 screen pattern と同様に injected `store` / `usecase` を受け取り、`onMount` で `initialize()` だけを呼ぶ。screen は view からの `refresh` event だけを usecase へ委譲し、control action、provider 判定、phase 集約ロジックを持たない。
- `src/ui/views/execution-observe/` に `ExecutionObserveView` を追加し、`docs/screen-design/wireframes/job-run.md` の観測面を read-only dashboard として再構成する。表示ブロックは `Job Summary`、`Failure Summary`、`Phase Timeline`、`Phase Runs`、`Translation Progress`、`Selected Unit Detail`、`Footer Metadata` に限定し、`execution-control` view の action row は入れない。
- `Failure Summary` は `ExecutionControlFailureDto.category` と `message` を provider-neutral に出す。`RecoverableFailed` / `Failed` / `Canceled` では主表示、その他 state では failure がある時だけ補助表示に留める。失敗時も操作系 UI は置かず、必要なら `Refresh` だけを許可する。
- 観測 path の mount は owned scope 外で `App.svelte` / `AppShell.svelte` / `main.ts` への additive な props 追加だけに留める。`execution-control` は別 screen として並存させ、navigation や shell 構造の再設計はしない。

## Scenario

- マウント時に observation snapshot を 1 回読み込み、ユーザーは current control state、現在 phase、phase runs、translation progress、最新 failure reason を 1 画面で確認できる。
- `Running` と `Retrying` では phase timeline の current marker、phase runs の最新 status、translation progress counters、selected unit detail を最新 snapshot で表示する。live streaming や polling はこの task では固定せず、再読込は明示的な `Refresh` で行う。
- `Paused` では最後に確定した progress を保持したまま pause 中であることを示し、footer metadata で provider run id や last event timestamp を確認できる。
- `RecoverableFailed` と `Failed` では failure summary を主表示にしつつ、phase runs と translation progress は最後に確定した snapshot を残す。これにより失敗理由と停止位置を同時に観測できるが、retry / cancel などの回復操作は `execution-control` path に委ねる。
- `Canceled` と `Completed` では read-only の最終状態として観測し、progress と footer metadata を保持する。`ExecutionControlTransitionDto` は action path 用 vocabulary として扱い、この画面では transition button や transition matrix を描画しない。
- 読み込み失敗時は最後に確定していた observation snapshot があれば表示を維持し、その上に generic observe error を出す。初回読み込み前に失敗した場合だけ empty state 相当の error panel を出し、`Refresh` で再試行できる。

## Logic

- usecase は stable control vocabulary を核にした provider-neutral な `ExecutionObserveSnapshot` を frontend 契約として固定する。最低限の必須項目は `controlState` と `failure` とし、画面表示用に `summary`、`phaseTimeline`、`phaseRuns`、`translationProgress`、`selectedUnit`、`footerMetadata` を 1 つの read model に集約して view へ渡す。
- progress 表示に使う値は adapter 内部 payload を直接見せず、executor が返す provider-neutral な count / status / timestamp / display text に閉じる。provider request body、transport status、raw error payload、writer output detail は screen state に入れない。
- public surface は `initialize()` と `refresh()` を最小固定対象とする。`initialize()` は初回 snapshot 読み込みだけを担当し、`retry()` / `pause()` / `resume()` / `cancel()` のような action API は持たせない。refresh 中は前回の確定 snapshot を保持し、error だけを上書きして view のちらつきを避ける。
- phase timeline と phase runs の並びは usecase が observation snapshot から整形し、view は `isCurrent`、`statusLabel`、開始 / 終了時刻などの表示済みフィールドだけを受け取る。selected unit detail も screen-local selection は持たず、snapshot が示す current または latest unit をそのまま表示する。
- 最小 shared wiring は `executionObserveStore` / `executionObserveUsecase` を `main.ts` から `App.svelte` と `AppShell.svelte` へ optional prop として通す点で固定する。実データ取得手段は injected executor に閉じ、`src/gateway/tauri/` や adapter 実装の shape をこの task の設計に持ち込まない。

## Implementation Plan

- ordered_scope:
  1. Distill task-local facts, dependencies, and relevant existing UI patterns for execution observation.
  2. Fix task-local UI, scenario, and logic design for progress and failure observation.
  3. Define minimum tests, fixtures, and validation commands before implementation.
  4. Implement the owned frontend scope plus minimal additive wiring only if required.
- owned_scope:
  - `src/application/usecases/execution-observe/`
  - `src/ui/screens/execution-observe/`
  - `src/ui/views/execution-observe/`
- required_reading:
  - `tasks/phase-4/tasks/P4-I06.yaml`
  - `docs/exec-plans/active/2026-04-07-p4-i06-progress-and-failure-observability-ui.md`
- validation_commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- users can observe execution progress and failure state through a stable UI path
- execution progress and failure state are exposed without relying on adapter internals
- control actions remain outside the execution-observe UI path

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- task-local test command to be fixed by `architecting-tests`
- owned lint suite to be fixed by downstream planning
- `sonar-scanner`
- Sonar MCP open issue query for project `ishibata91_AITranslationEngineJP` scoped to touched paths
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新要否を実装結果で判断する。現時点では未確定。
- `4humans/tech-debt-tracker.md`
  更新要否を実装結果で判断する。現時点では未確定。
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  execution-observe の構造追加や依存追加があれば対象を plan に具体化する。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  execution-observe の観測フロー追加や変更があれば対象を plan に具体化する。
- `4humans/diagrams/overview-manifest.json`
  new detail `.d2` を追加する場合だけ更新対象へ含める。

## Outcome

- In progress.
