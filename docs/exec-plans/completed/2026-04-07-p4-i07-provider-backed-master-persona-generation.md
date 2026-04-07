- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement `P4-I07` provider-backed master persona generation without rebuilding Phase 2 persistence semantics or job-local persona orchestration.
- task_id: P4-I07
- task_catalog_ref: tasks/phase-4/tasks/P4-I07.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I07`.
- Generate and persist master persona data through the provider runtime contract for the base-game NPC path.

## Decision Basis

- `tasks/phase-4/tasks/P4-I07.yaml`
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md`
- `docs/exec-plans/completed/2026-04-04-p2-i03-master-persona-builder.md`

## Owned Scope

- `src-tauri/src/application/master_persona_generation_runtime/`

## Out Of Scope

- job-local persona generation phase orchestration
- provider selection UI
- redesign of Phase 2 storage semantics
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-I07` depends on `P4-C01`, `P4-C03`, `P4-V02`, and `P2-I03`.
- The provider-backed runtime path must reuse the stable provider selection and persona-generation runtime contracts from `P4-B1`.
- The persistence path must reuse the existing master persona foundation boundary from `P2-I03`.

## Parallel Safety Notes

- Keep implementation isolated to `src-tauri/src/application/master_persona_generation_runtime/`.
- Reuse existing provider runtime and master persona boundaries instead of reopening adapter internals or storage schema work.
- Do not pull in job-local persona generation or provider selection UI ownership from sibling Phase 4 tasks.

## UI

- N/A. `P4-I07` は backend-only とし、provider selection UI、master persona command surface、execution observation transport は既存 owner に残す。
- backend 観測は既存の `persona-rebuild` validation と `persona-generation-runtime` acceptance anchor を前提にし、この task で新しい UI / transport anchor は増やさない。

## Scenario

- `P4-I07` は `BaseGameNpcRebuildRequest` を base-game NPC 入力の正本として維持したまま、同一実行で `PersonaGenerationRuntimeRequestDto` を併走させる薄い application path とする。runtime request は `source.kind = MasterPersonaSeed` と `sink = PersonaStorage` に固定し、`source.source_key` は provider 実行の trace key としてのみ扱う。
- 実行順は `provider runtime -> existing master persona rebuild/save/read` の 2 段に限定する。provider 成功後だけ既存 `RebuildMasterPersonaUseCase` を呼び、永続化と read-back は Phase 2 の `MasterPersonaStoragePort` へそのまま委譲する。
- `translation_unit` source や `translation_phase_handoff` sink はこの module で拒否し、job-local persona orchestration へ分岐しない。provider failure 時は rebuild/save/read を呼ばず、master persona 側 failure は既存 save/read semantics を変えずに扱う。

## Logic

- `src-tauri/src/application/master_persona_generation_runtime/` には、既存 `RebuildMasterPersonaUseCase` を置き換えない薄い orchestrator を置く。依存は `ProviderRuntimePort` と master persona rebuild boundary だけに限定し、dictionary、provider adapter internals、job-local persona phase は参照しない。
- task-local request は少なくとも `PersonaGenerationRuntimeRequestDto` と `BaseGameNpcRebuildRequest` を束ねる。module 内で `source/sink` guard を行い、`ProviderRuntimePort` へ渡す surface は `ProviderSelectionDto` だけに留める。provider 固有の prompt / credential / transport detail は application 層へ入れない。
- 戻り値は persisted master persona read result を維持し、失敗語彙は Phase 4 の `ExecutionControlFailureDto` へ寄せる。provider failure は reshape せず返し、rebuild / save / read の `String` failure だけ `ValidationFailure` 系として包んで observation vocabulary を揃える。
- 実装アンカーは 3 つに固定する。`src-tauri/tests/validation/persona-rebuild/mod.rs` の既存 fixture / snapshot で persisted master persona shape と entry 順序を守り、`src-tauri/tests/acceptance/persona-generation-runtime/mod.rs` の shared runtime fixture で master-persona source/sink の vocabulary を守り、新規 module-local test では provider success 後だけ rebuild が走る順序、master-persona 以外の route rejection、provider failure の passthrough を確認する。

## Implementation Plan

### ordered_scope

1. Application orchestrator (`src-tauri/src/application/master_persona_generation_runtime/`, `src-tauri/src/application/mod.rs`)

- Add the new backend-only module that bundles `BaseGameNpcRebuildRequest` with `PersonaGenerationRuntimeRequestDto` and orchestrates `provider runtime -> existing master persona rebuild/save/read`.
- Keep the route guard local to this module. Accept only `MasterPersonaSeed -> PersonaStorage`, pass only `ProviderSelectionDto` into `ProviderRuntimePort`, and do not branch into job-local persona orchestration or provider adapter detail.

2. Failure normalization and persistence-boundary preservation (`src-tauri/src/application/master_persona_generation_runtime/`)

- Return the persisted master persona read result without redefining Phase 2 storage behavior. Reuse `RebuildMasterPersonaUseCase`, `MasterPersonaBuilderPort`, and `MasterPersonaStoragePort` as-is instead of reopening persistence semantics.
- Preserve provider failures as `ExecutionControlFailureDto` passthrough, and wrap rebuild / save / read `String` failures into `ValidationFailure` only inside this module so observation vocabulary stays aligned with Phase 4 contracts.

3. Verification anchors (`src-tauri/src/application/master_persona_generation_runtime/`, `src-tauri/tests/validation/persona-rebuild/`, `src-tauri/tests/acceptance/persona-generation-runtime/`)

- Add the minimum module-local tests needed to prove execution order and route rejection: provider success triggers rebuild once, provider failure skips rebuild, and non-master-persona source/sink combinations fail before persistence.
- Let `architecting-tests` choose the exact minimal file touch list inside the existing validation / acceptance anchors, while keeping those anchors responsible for persisted master persona shape and shared persona-generation runtime vocabulary.

### owned_scope

- backend

### required_reading

- `docs/exec-plans/active/2026-04-07-p4-i07-provider-backed-master-persona-generation.md` の `UI` / `Scenario` / `Logic`
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md`
- `docs/exec-plans/completed/2026-04-04-p2-i03-master-persona-builder.md`
- `src-tauri/src/application/master_persona/mod.rs`
- `src-tauri/src/application/dto/persona_generation_runtime/mod.rs`
- `src-tauri/src/application/dto/execution_control/mod.rs`
- `src-tauri/src/application/ports/provider_runtime/mod.rs`
- `src-tauri/src/application/ports/persona_generation_runtime/mod.rs`
- `src-tauri/src/application/ports/persona_storage/mod.rs`
- `src-tauri/src/application/mod.rs`
- `src-tauri/tests/validation/persona-rebuild/mod.rs`
- `src-tauri/tests/acceptance/persona-generation-runtime/mod.rs`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test acceptance -- --nocapture`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/master_persona_generation_runtime src-tauri/tests/validation/persona-rebuild src-tauri/tests/acceptance/persona-generation-runtime`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- base-game NPC input can produce persisted master persona data through the stable provider runtime path
- provider-backed master persona generation runs without redesigning the Phase 2 persistence and observation path

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/master_persona_generation_runtime`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。品質評価軸や残留リスクの正本を追加で変える変更ではなかった。
- `4humans/tech-debt-tracker.md`
  更新なし。新規の恒久負債項目は残さなかった。
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.d2`
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.svg`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.svg`

## Outcome

- Added `src-tauri/src/application/master_persona_generation_runtime/mod.rs` and `src-tauri/src/application/mod.rs` wiring so a backend-only orchestrator now validates the route, runs `ProviderRuntimePort`, and then delegates to existing `RebuildMasterPersonaUseCase` for persisted master persona read-back.
- Kept the owned route limited to `MasterPersonaSeed -> PersonaStorage`, returned provider failures as `ExecutionControlFailureDto` passthrough, and normalized rebuild/save/read `String` failures into `ValidationFailure` without reopening Phase 2 persistence semantics.
- Added `src-tauri/tests/validation/master-persona-generation-runtime/mod.rs`, wired it through `src-tauri/tests/persona_rebuild_validation.rs`, and tightened `src-tauri/tests/acceptance/persona-generation-runtime/mod.rs` so both invalid pairings and the canonical job-local route stay outside this module.
- Updated `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.d2` / `.svg` and `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.d2` / `.svg` to show the provider-runtime gate, route rejection, and provider-success-only rebuild flow.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `sonar-scanner`, Sonar owned-scope `OPEN` issues = `0`, `d2 validate`, `d2 -t 201`, and `python3 scripts/harness/run.py --suite all`.
- Single-pass review returned `reroute` for missing job-local rejection coverage and missing `4humans` sync. Both were fixed, and per lane contract no second review was run.
