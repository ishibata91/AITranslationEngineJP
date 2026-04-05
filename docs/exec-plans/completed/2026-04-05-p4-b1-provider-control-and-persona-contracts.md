- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Fix Phase 4 batch `P4-B1` provider selection, execution-control, and persona-generation contracts plus failure acceptance anchors.
- task_id: P4-B1
- task_catalog_ref: tasks/phase-4/phase.yaml
- parent_phase: phase-4

## Request Summary

- Implement batch `P4-B1` from `tasks/phase-4/phase.yaml`.
- Land the contract and verification foundation for provider selection, execution-control state expansion, and provider-backed persona-generation runtime before Phase 4 implementation tasks start.

## Decision Basis

- `tasks/phase-4/phase.yaml` defines `P4-B1` as the contract and verification batch that must fix provider, execution-control, and persona-generation boundaries before `P4-B2`.
- `P4-C01`, `P4-C02`, and `P4-C03` own backend DTO and port or state boundaries under `src-tauri/src/application/` and `src-tauri/src/domain/`.
- `P4-V01` and `P4-V02` own acceptance anchors under `src-tauri/tests/acceptance/` and must stay provider-independent at the fixture boundary.
- `tasks/phase-4/phase.yaml` currently points at `tasks/P4-*.yaml`, while the actual task files live under `tasks/phase-4/tasks/`. This mismatch is catalog-only context unless the batch work needs to correct it.

## Owned Scope

- `src-tauri/src/application/dto/provider_selection/`
- `src-tauri/src/application/ports/provider_runtime/`
- `src-tauri/src/domain/execution_control_state/`
- `src-tauri/src/application/dto/execution_control/`
- `src-tauri/src/application/dto/persona_generation_runtime/`
- `src-tauri/src/application/ports/persona_generation_runtime/`
- `src-tauri/tests/acceptance/provider-failure-retry/`
- `src-tauri/tests/acceptance/persona-generation-runtime/`

## Out Of Scope

- provider adapter internals
- provider-specific transport detail
- output writer behavior
- dictionary or persona foundation rebuild validation
- translation body output verification
- broad UI layout changes outside acceptance-anchor needs
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C03` depends on `P4-C01`, `P2-C03`, and `P3-C03`.
- `P4-V01` depends on `P4-C01` and `P4-C02`.
- `P4-V02` depends on `P4-C01`, `P4-C02`, and `P4-C03`.
- Relevant Phase 2 and Phase 3 contract artifacts must already exist and remain stable.

## Parallel Safety Notes

- Keep provider selection, execution-control state, persona-generation runtime, and acceptance fixtures separated by the catalog-owned scopes above.
- Do not let provider-specific prompt, transport, or snapshot detail leak into contracts or acceptance fixtures.
- Coordinate backend contract files and acceptance fixtures, but keep them isolated from provider adapter implementation scopes in `P4-B2`.

## UI

- N/A. `P4-B1` は backend contract と acceptance anchor の固定だけを扱う。後続 UI は provider selection / execution-control の語彙を再利用してよいが、この plan では画面構成、操作フロー、表示文言を決めない。

## Scenario

- 実行系は `provider_selection` 契約で provider 識別子、execution mode、provider 非依存の runtime 設定だけを受け渡し、adapter 固有の transport / credential / prompt / snapshot detail は `provider_runtime` 実装側へ閉じ込める。
- Job の create/list が依存する Phase 1 の最小 `JobState` は `Draft` / `Ready` / `Running` / `Completed` のまま維持し、pause / retry / recoverable failure / failed / canceled は Phase 4 所有の `execution_control_state` 契約で追加表現する。
- master persona rebuild と job-local NPC persona generation は同じ `persona_generation_runtime` 契約を使うが、前者の保存先は Phase 2 の `persona_storage`、後者の handoff は Phase 3 の `translation_phase_handoff` に残す。acceptance fixture は success / failure / retry を共通 runtime 境界で固定し、provider 固有 detail を埋め込まない。

## Logic

- `src-tauri/src/application/dto/` と `src-tauri/src/application/ports/` には `provider_selection`、`execution_control`、`persona_generation_runtime` を additive に追加し、既存の `persona_storage`、`translation_phase_handoff`、`job` DTO shape は再設計しない。downstream 実装が必要とする共通 surface は root `mod.rs` から再 export できる前提で揃える。
- provider selection contract は execution-control と persona-generation runtime の両方が共有できる最小 DTO に留める。固定するのは provider 選択、実行方式、retry/pause 可否判定に必要な provider 非依存設定までとし、adapter 固有 error payload や API request body は shared contract に入れない。
- execution-control contract は `domain/job_state` の widening ではなく、Phase 4 所有の `domain/execution_control_state/` と `application/dto/execution_control/` で表現する。state / transition vocabulary は pause、resume、retry、recoverable failure、failed、canceled、completed の外部観測を満たしつつ、Phase 1 DTO へ必要以上に逆流させない。
- `P4-V01` は provider failure / retry / recovery を generic fixture で固定し、assertion は state transition と user-visible failure category に限る。`P4-V02` は master / job-local の 2 経路を同じ persona-generation runtime fixture family で固定し、差分は upstream source envelope と downstream sink (`persona_storage` または `translation_phase_handoff`) だけに留める。

## Implementation Plan

- Ordered scope 1 (`P4-C01`): scaffold `src-tauri/src/application/dto/provider_selection/` and `src-tauri/src/application/ports/provider_runtime/`, then extend `src-tauri/src/application/dto/mod.rs` and `src-tauri/src/application/ports/mod.rs` with a provider-neutral selection and runtime-configuration contract. Keep only provider identifiers, execution mode, and shared runtime knobs; exclude transport, credential, prompt, snapshot, and adapter error payload detail.
- Ordered scope 2 (`P4-C02`): scaffold `src-tauri/src/domain/execution_control_state/` and `src-tauri/src/application/dto/execution_control/`, then extend `src-tauri/src/domain/mod.rs` and DTO exports with Phase 4-only control states and transition vocabulary for pause, resume, retry, recoverable failure, failed, canceled, and completed without widening `src-tauri/src/domain/job_state/`.
- Ordered scope 3 (`P4-C03`): scaffold `src-tauri/src/application/dto/persona_generation_runtime/` and `src-tauri/src/application/ports/persona_generation_runtime/`, then define one provider-independent runtime contract that reuses the new provider-selection boundary plus existing Phase 2 `persona_storage` and Phase 3 `translation_phase_handoff` roots instead of redefining either DTO family.
- Ordered scope 4 (`P4-V01`): scaffold `src-tauri/tests/acceptance/provider-failure-retry/` and update `src-tauri/tests/acceptance.rs` with the thinnest loader needed for Cargo discovery, then fix one generic fixture family that anchors provider failure, retry, pause, and recovery through execution-control transitions and provider-independent failure categories only.
- Ordered scope 5 (`P4-V02`): scaffold `src-tauri/tests/acceptance/persona-generation-runtime/` under the same acceptance loader, then fix one shared fixture family that covers master-persona rebuild and job-local persona generation with success, failure, and retry paths while varying only the upstream source envelope and downstream sink (`persona_storage` or `translation_phase_handoff`).
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `cargo test --manifest-path src-tauri/Cargo.toml --test acceptance -- --nocapture`
  - `cargo test --manifest-path src-tauri/Cargo.toml --all-features`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/dto/provider_selection src-tauri/src/application/ports/provider_runtime src-tauri/src/domain/execution_control_state src-tauri/src/application/dto/execution_control src-tauri/src/application/dto/persona_generation_runtime src-tauri/src/application/ports/persona_generation_runtime src-tauri/tests/acceptance/provider-failure-retry src-tauri/tests/acceptance/persona-generation-runtime`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- Provider selection and persona-generation runtime contracts remain provider-independent.
- Execution-control states cover pause, retry, recoverable failure, failure, and cancel without embedding provider transport details.
- Acceptance fixtures prove provider failure, retry, pause, recovery, and provider-backed persona-generation behavior through stable external anchors.
- Phase 4 implementation tasks can consume the new contracts and acceptance anchors without redesigning Phase 2 or Phase 3 boundaries.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。既存品質記録の評価軸を変える変更ではなかった。
- `4humans/tech-debt-tracker.md`
  更新なし。新規の恒久負債項目は残さなかった。
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.d2`
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.svg`
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.svg`

## Outcome

- Added Phase 4 contract modules for `provider_selection`, `execution_control`, and `persona_generation_runtime`, including root exports under `src-tauri/src/application/dto/`, `src-tauri/src/application/ports/`, and `src-tauri/src/domain/`.
- Fixed provider-independent acceptance anchors under `src-tauri/tests/acceptance/provider-failure-retry/` and `src-tauri/tests/acceptance/persona-generation-runtime/`, plus `json_contract_guard` support for forbidden provider-specific fixture keys.
- Kept `src-tauri/src/domain/job_state/mod.rs` unchanged while adding Phase 4-only execution-control transitions and unit tests in `src-tauri/src/domain/execution_control_state/mod.rs`.
- Unified provider and persona runtime failure boundaries on `ExecutionControlFailureDto` instead of opaque strings.
- Updated the backend translation-flow review diagrams to show `ProviderRuntimePort`, `PersonaGenerationRuntimePort`, and `ExecutionControlState`.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` reporting `openIssueCount: 0`, `d2 validate`, `d2 -t 201`, and `python3 scripts/harness/run.py --suite all`.
