- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I04` batch and single-shot switching through one stable execution-mode path.
- task_id: P4-I04
- task_catalog_ref: tasks/phase-4/tasks/P4-I04.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I04` from `tasks/phase-4/tasks/P4-I04.yaml`.
- Add one stable execution-mode switching path under `src-tauri/src/application/execution_mode_switch/` without changing provider contracts or job state policy.

## Decision Basis

- `tasks/phase-4/tasks/P4-I04.yaml` defines a backend implementation scope limited to `src-tauri/src/application/execution_mode_switch/`.
- `tasks/phase-4/phase.yaml` places `P4-I04` inside Phase 4 implementation batch `P4-B2`, parallel with adapter and UI slices that must keep owned scopes disjoint.
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md` fixed provider-selection and execution-control contracts that this switching path must consume without reshaping.
- Detailed repo facts, design constraints, and required reading are delegated to `distilling-implementation`.

## Owned Scope

- `src-tauri/src/application/execution_mode_switch/`

## Out Of Scope

- provider transport internals
- control UI
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C01`
- `P4-C02`
- existing provider runtime contract and execution-control state contract must stay stable

## Parallel Safety Notes

- Keep changes centered on `src-tauri/src/application/execution_mode_switch/`.
- Do not widen provider adapter scopes or UI scopes owned by sibling Phase 4 tasks.
- If composition touch points outside owned scope are required, justify them narrowly in downstream design and implementation handoff.

## UI

- N/A. この task が固定するのは application 層の execution-mode switching path だけとし、control UI、Gateway command、表示文言、screen state は決めない。

## Scenario

- 呼び出し側は既存の `ProviderSelectionDto` をそのまま `execution_mode_switch` へ渡し、execution mode の分岐判断を application 層の 1 箇所へ集約する。provider 選択、runtime settings、failure vocabulary は既存 contract を再利用し、新しい shared DTO や job state は追加しない。
- `ProviderExecutionModeDto::Batch` は既存 `ProviderRuntimePort::run_provider_step()` へ `ProviderSelectionDto` を変更せず委譲する。batch adapter ごとの `provider_id` 検証、credential/config 検証、transport 開始条件は引き続き各 `src-tauri/src/infra/provider/*` が所有する。
- `ProviderExecutionModeDto::Streaming` はこの task では single-shot 選択経路として扱い、application module に注入された provider-neutral な single-shot delegate へ同じ `ProviderSelectionDto` を渡す。single-shot 側の provider-specific rule や transport detail は switching module に持ち込まない。
- switching path 自体は `ExecutionControlFailureDto` をそのまま透過し、pause / retry / recoverable failure / canceled などの job state policy は既存 execution-control layer に残す。mode を切り替えても state transition vocabulary や retry policy は増やさない。

## Logic

- `src-tauri/src/application/execution_mode_switch/mod.rs` を additive に追加し、1 つの public dispatcher / use case だけを公開する。public surface は `ProviderSelectionDto -> Result<(), ExecutionControlFailureDto>` の薄い形に留め、shared contract の再設計はしない。
- dispatcher は `selection.execution_mode` だけで分岐する。`Batch` branch は `ProviderRuntimePort` を依存先に使い、`Streaming` branch は module-local な single-shot port もしくは同等の injected delegate を依存先に使う。どちらの branch も `ProviderSelectionDto` を clone で受け渡すだけで、mode 以外の field を再解釈しない。
- module 外の touch point は最小に保ち、実装上必要な repo-wide 変更は `src-tauri/src/application/mod.rs` からの export 1 箇所までを正当化上限とする。composition root や integration caller は後続 task で選ぶが、switching rule 自体はこの module 以外へ複製しない。
- validation と test の境界は application module に閉じる。ここで固定するのは「Batch なら batch delegate、Streaming なら single-shot delegate が呼ばれること」と「delegate が返した `ExecutionControlFailureDto` を改変しないこと」までとし、provider adapter ごとの unsupported mode / config / transport 振る舞いは既存 adapter test と acceptance anchor に残す。

## Implementation Plan

### ordered_scope

1. `src-tauri/src/application/execution_mode_switch/`

- `ProviderSelectionDto` を受けて `selection.execution_mode` だけで分岐する 1 つの public dispatcher / use case を追加する。
- `Batch` branch は既存 `ProviderRuntimePort::run_provider_step()` へ `ProviderSelectionDto` をそのまま委譲する。
- `Streaming` branch は module-local な single-shot port または同等の injected delegate へ同じ `ProviderSelectionDto` をそのまま委譲する。
- switching module は provider 固有の validation、runtime settings の再解釈、job state policy の分岐を持たず、delegate が返した `ExecutionControlFailureDto` を改変せず透過する。

2. `src-tauri/src/application/execution_mode_switch/` の module-local tests

- batch delegate が選ばれること、single-shot delegate が選ばれること、非選択 branch が呼ばれないことを spy で固定する。
- batch / single-shot の各 branch で返された `ExecutionControlFailureDto` が変更されず返ることを固定する。
- test 境界は application switching logic に閉じ、provider transport や adapter 固有 validation は再検証しない。

3. `src-tauri/src/application/mod.rs`

- downstream caller が additive に参照できる最小 export として `execution_mode_switch` module を公開する。
- composition root、Gateway command、provider adapter module、shared DTO / runtime contract へは広げない。

### owned_scope

- `src-tauri/src/application/execution_mode_switch/`
- `src-tauri/src/application/mod.rs`

### required_reading

- `docs/exec-plans/completed/2026-04-06-p4-i04-batch-and-single-shot-switching.md` の `UI` / `Scenario` / `Logic`
- `tasks/phase-4/tasks/P4-I04.yaml`
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md`
- `src-tauri/src/application/dto/provider_selection/mod.rs`
- `src-tauri/src/application/ports/provider_runtime/mod.rs`
- `src-tauri/src/application/dto/execution_control/mod.rs`
- `src-tauri/src/application/mod.rs`
- `src-tauri/src/infra/provider/lmstudio/mod.rs`
- `src-tauri/src/infra/provider/gemini/mod.rs`
- `src-tauri/src/infra/provider/xai/mod.rs`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml execution_mode_switch -- --nocapture`
- `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/execution_mode_switch src-tauri/src/application/mod.rs`
- `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`

### implementation_updates

- `Implementation Plan` を `ordered_scope` / `owned_scope` / `required_reading` / `validation_commands` の brief 形式へ更新し、touch point を `src-tauri/src/application/mod.rs` までに固定した。

## Acceptance Checks

- Batch and single-shot execution can be selected through one stable path.
- Provider runtime and execution-control boundaries remain unchanged.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml execution_mode_switch -- --nocapture`
- `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/execution_mode_switch src-tauri/src/application/mod.rs`
- `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。既存の品質評価軸や未充足項目は今回の additive な switching module 追加では変わらなかった。
- `4humans/tech-debt-tracker.md`
  更新なし。恒久負債として残すべき追加項目は発生しなかった。
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  更新なし。`execution_mode_switch` module の追加と `application/mod.rs` の export 追加に留まり、既存の構造説明を差し替える必要はないと single-pass review で確認した。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  更新なし。caller や process narration の既存フローは変えていないため、diagram sync は不要と判断した。

## Outcome

- Added `src-tauri/src/application/execution_mode_switch/mod.rs` with `SwitchExecutionModeUseCase`, a provider-neutral `SingleShotExecutionDelegate`, and module-local tests that prove batch versus single-shot delegate selection plus `ExecutionControlFailureDto` passthrough only.
- Updated `src-tauri/src/application/mod.rs` to export `execution_mode_switch` as the single justified touch point outside the owned module.
- Single-pass review found no code reroute items; the only reroute was missing close evidence, which was resolved by running the required `structure` and full harness commands.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml execution_mode_switch -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/execution_mode_switch src-tauri/src/application/mod.rs` with `openIssueCount: 0`, and `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`.
