- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-I03` from `tasks/phase-2/tasks/P2-I03.yaml` by building master persona data from base-game NPC input without widening into job-local persona storage or dictionary logic.
- task_id: P2-I03
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-I03`.
- Build master persona foundation data from base-game NPC input behind the persona storage split fixed in `P2-C03` and the persona rebuild validation anchor fixed in `P2-V02`.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-2/tasks/P2-I03.yaml`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`

## Owned Scope

- `src-tauri/src/application/master_persona/`
- `src-tauri/src/infra/master_persona_builder/`

## Out Of Scope

- job-local persona storage
- dictionary logic
- UI observation and rendering

## Dependencies / Blockers

- `P2-I03` depends on `P2-C03` and `P2-V02`.
- The implementation must preserve the split between master persona data and job-local persona data established by `P2-B1`.

## Parallel Safety Notes

- `P2-I03` stays parallel-safe with the other `P2-B2` tasks only while it remains inside master persona builder ownership.
- Storage-boundary redesign, dictionary behavior, and UI wiring must stay outside this task.

## UI

- N/A. `P2-I03` は backend master persona build path だけを扱う。
- persona observation UI、Gateway command、screen state、error presentation は `P2-I06` と `P2-G01` へ残し、この task で transport surface を増やさない。

## Scenario

- ベースゲーム NPC rebuild は、builder が base-game NPC 入力から `MasterPersonaSaveRequestDto { persona_name, source_type, entries }` を組み立て、`MasterPersonaStoragePort` へ渡す 1 本の foundation path として扱う。初期実装では validation fixture と同じ `source_type = "base-game-rebuild"` を基準にし、entry には `npc_form_id`、`npc_name`、`race`、`sex`、`voice`、`persona_text` だけを含める。
- 同じ base-game NPC 入力を再実行したときは、同じ `persona_name`、`source_type`、entry 順序を再現できることを rebuild の受け入れ基準にする。builder は trim、dedupe、sort、job-local metadata 補完を行わず、入力で与えられた NPC 属性と `persona_text` をそのまま master persona 側 contract へ写す。
- rebuild 後の確認は `MasterPersonaReadRequestDto { persona_name }` で master persona を読み返す経路を前提にし、`P2-V02` の fixture / snapshot が示す master-persona shape を壊さない。job persona 側の `job_id` や `JOB_PERSONA_ENTRY` はこのシナリオへ入れない。
- base-game NPC 入力の decode 失敗、必須属性欠落、unsupported source のような build failure は `Err(String)` で返し、途中まで組み立てた master persona を部分保存しない。

## Logic

- `src-tauri/src/application/master_persona/` は orchestration 境界として、base-game NPC 入力を受ける builder port と既存 `MasterPersonaStoragePort` を束ねる use case を持つ。application 層が固定するのは `base-game NPC input -> MasterPersonaSaveRequestDto -> save/read` の流れだけであり、job persona DTO、dictionary module、UI read model は参照しない。
- builder の入力 shape は、この task の owned scope に閉じた最小構成に留める。少なくとも fixture で既に固定されている `persona_name` / `source_type` と、各 NPC の `npc_form_id`、`npc_name`、`race`、`sex`、`voice`、`persona_text` を表現できればよく、job identifier、foundation 以外の観測 metadata、dictionary 由来情報は追加しない。
- `src-tauri/src/infra/master_persona_builder/` は base-game NPC 入力を master persona entry 群へ正規化する adapter とし、entry 順序を保持したまま `MasterPersonaSaveRequestDto` 互換の出力を返す。base-game NPC 側の入力形式が将来広がっても、この task では fixture が示す 1 行 shape を壊さず、必要最小限の field mapping に閉じる。
- 永続化の責務は既存 `MasterPersonaStoragePort` の向こう側へ残し、この task で `JOB_PERSONA_ENTRY` 側 repository や dictionary persistence を抱き込まない。master persona rebuild の検証に保存先が必要な場合も、境界は master 側 save/read contract だけで表現し、schema / migration 再設計へ広げない。
- `P2-V02` の persona rebuild anchor は downstream validation の基準として維持し、実装は master persona と job persona の非代替性を前提に組む。同じ `npc_form_id` と属性列を共有しても、master 側だけが `persona_name` と `source_type` を持つ、という `P2-C03` の split を module 構成と data flow の両方で守る。

## Implementation Plan

### ordered_scope

1. Validation anchor (`src-tauri/tests/validation/persona-rebuild/`, `src-tauri/tests/persona_rebuild_validation.rs`)

- Extend the existing `P2-V02` anchor from DTO non-substitutability into a deterministic master-persona rebuild path that exercises `base-game NPC input -> build -> save/read` without widening into UI or integration behavior.
- Keep the current fixture vocabulary and snapshot determinism intact. Reuse `base-game-rebuild` and the pinned NPC row shape, then add only the minimum stubbed storage / builder harness needed to prove the builder preserves `persona_name`, `source_type`, entry order, and master-vs-job separation.

2. Application master persona orchestration (`src-tauri/src/application/master_persona/`, `src-tauri/src/application/mod.rs`)

- Add the application-owned input shape, builder port, and use case that translate base-game NPC rebuild input into `MasterPersonaSaveRequestDto`, call `MasterPersonaStoragePort`, and read back through `MasterPersonaReadRequestDto`.
- Keep the application boundary limited to orchestration. Do not add job identifiers, dictionary dependencies, schema concerns, or UI-facing transport models. Failure stays `Err(String)` and partial-save handling remains delegated to the storage port contract.

3. Infra master persona builder (`src-tauri/src/infra/master_persona_builder/`, `src-tauri/src/infra/mod.rs`)

- Add the adapter that maps the minimal base-game NPC input shape into ordered `MasterPersonaEntryDto` values and produces a `MasterPersonaSaveRequestDto` compatible result for the application use case.
- Keep mapping behavior thin and literal: no trim, dedupe, sort, enrichment, dictionary lookup, or job-local metadata synthesis. Unsupported input or missing required fields should fail before save is attempted.

4. Boundary wiring and explicit non-scope decisions

- Wire the new backend modules through `src-tauri/src/application/mod.rs` and `src-tauri/src/infra/mod.rs` only.
- Do not add migration or repository work in this task. `MASTER_PERSONA` schema creation remains outside `P2-I03` ownership unless scope is explicitly expanded later.

### owned_scope

- backend

### required_reading

- `docs/exec-plans/active/2026-04-04-p2-i03-master-persona-builder.md` の `UI` / `Scenario` / `Logic`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`
- `tasks/phase-2/tasks/P2-I03.yaml`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/er.md`
- `src-tauri/src/application/dto/persona_storage/mod.rs`
- `src-tauri/src/application/ports/persona_storage/mod.rs`
- `src-tauri/tests/persona_rebuild_validation.rs`
- `src-tauri/tests/validation/persona-rebuild/mod.rs`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- master persona data can be rebuilt independently
- master persona build does not redefine the storage boundary

## Required Evidence

- Active plan updated with distill, design, implementation brief, and validation commands.
- Tests or fixtures proving master persona rebuild behavior remains independent from job-local persona handling.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.d2`
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.svg`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.svg`
- `4humans/diagrams/overview-manifest.json`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.svg`
- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/processes-overview-robustness.svg`

## Outcome

- Added `src-tauri/src/application/master_persona/mod.rs` and module wiring in `src-tauri/src/application/mod.rs`, so base-game NPC rebuild input can be orchestrated through `MasterPersonaSaveRequestDto -> save/read` without widening into job persona or dictionary modules.
- Added `src-tauri/src/infra/master_persona_builder/mod.rs` and wiring in `src-tauri/src/infra/mod.rs`, so the backend can validate `source_type = "base-game-rebuild"` and map ordered base-game NPC rows into `MasterPersonaEntryDto` values without trim, dedupe, sort, or enrichment.
- Extended `src-tauri/tests/validation/persona-rebuild/mod.rs` with a deterministic rebuild fixture / snapshot path and failure-path checks for unsupported `source_type` and whitespace-only `persona_name`, proving storage calls are skipped on builder rejection while the existing master-vs-job split remains intact.
- Added `src-tauri/tests/validation/persona-rebuild/fixtures/base-game-master-persona-rebuild.fixture.json` and `src-tauri/tests/validation/persona-rebuild/snapshots/base-game-master-persona-rebuild.snapshot.json` to pin the base-game master persona rebuild contract.
- Added `4humans` D2 review diagrams for the master persona builder slice and rebuild flow, then linked both detail diagrams from the backend overview / process overview and registered them in `4humans/diagrams/overview-manifest.json`.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` (`openIssueCount: 0`), and `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home python3 scripts/harness/run.py --suite all`.
- Single-pass review initially returned `reroute` for missing validation coverage on empty/whitespace rejection. The reroute fix added the whitespace-only `persona_name` test above, and per lane contract no second review was run.
