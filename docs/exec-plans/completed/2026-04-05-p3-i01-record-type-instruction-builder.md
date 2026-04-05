- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P3-I01` by adding a backend record-type instruction builder through the agreed Phase 3 contract.
- task_id: P3-I01
- task_catalog_ref: tasks/phase-3/tasks/P3-I01.yaml
- parent_phase: phase-3

## Request Summary

- Implement `tasks/phase-3/tasks/P3-I01.yaml`.
- Build translation instructions per record type through the agreed contract without absorbing provider transport or output writing.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-3/tasks/P3-I01.yaml`
- `tasks/phase-3/tasks/P3-C01.yaml`
- `tasks/phase-3/tasks/P3-V02.yaml`
- `docs/exec-plans/completed/2026-04-05-p3-b1-phase-contracts-and-regression-anchors.md`

## Owned Scope

- `src-tauri/src/application/translation_instruction_builder/`

## Out Of Scope

- provider transport
- output writing
- unrelated translation phase composition outside the builder owned scope

## Dependencies / Blockers

- `P3-C01` must already define the translation instruction contract used by this builder.
- `P3-V02` fixture and regression anchor must stay reusable as the representative scenario acceptance reference.

## Parallel Safety Notes

- Task catalog marks `P3-I01` parallel-safe with `P3-I02`, `P3-I03`, `P3-I04`, and `P3-I05`, but owned scope must stay inside `src-tauri/src/application/translation_instruction_builder/`.
- Shared risks explicitly exclude embedding provider selection into this task because that would collapse the Phase 3 and Phase 4 boundary.

## UI

- N/A. Backend-only implementation inside `src-tauri/src/application/translation_instruction_builder/`; preview UI, Tauri command surface, and provider-facing presentation stay unchanged.

## Scenario

- Initial scenario is the `P3-V02` representative anchor only: one `TranslationUnitDto` for `dialogue_response` / `INFO` / `text` must pass through one builder entrypoint and yield one `TranslationInstructionDto`.
- The builder output for the anchored case must keep `phase_code=body_translation`, `unit_key=translation_unit.extraction_key`, the original `translation_unit`, and the instruction text aligned with the existing regression snapshot for `dialogue_response.text`.
- The anchored fixture with `<Alias=Player>` remains the acceptance reference, but the builder responsibility stops at instruction construction; embedded-element preservation, reusable terms, persona handoff, provider transport, retry, and output writing stay outside this task.
- Broader record-type expansion is deferred. Until another task fixes additional wording rules, unsupported record-type combinations should stay behind the same entrypoint and fail explicitly instead of widening the contract ad hoc.

## Logic

- Add one application-rooted builder module under `src-tauri/src/application/translation_instruction_builder/` with a single stable entrypoint that consumes `TranslationUnitDto` and returns `Result<TranslationInstructionDto, ...>`.
- The entrypoint should derive record-type identity only from `translation_unit.source_entity_type`, `translation_unit.record_signature`, and `translation_unit.field_name`; do not duplicate those fields into a separate request DTO inside this task.
- Fix `phase_code` to the Phase 3 body-translation contract value and derive `unit_key` from `translation_unit.extraction_key` so the builder output stays aligned with `TranslationInstructionDto` and the regression snapshot shape.
- Keep instruction wording provider-neutral by resolving a small internal rule set keyed by record type. For this task, only the `dialogue_response` / `INFO` / `text` rule is required, and it must produce the anchored Skyrim NPC dialogue instruction that preserves embedded elements such as `<Alias=Player>` exactly.
- Keep the rule selection and instruction text assembly inside the application layer module. Do not add provider selection, transport concerns, output formatting, phase handoff composition, or persistence branching here.
- If the entrypoint receives an unsupported record-type combination, return an explicit builder error from the module boundary so later tasks can add rules without changing the DTO contract or silently emitting generic instructions.

## Implementation Plan

- Ordered scope 1 (`src-tauri/tests/`): add or update backend validation so one stable builder entrypoint is exercised through the anchored `dialogue_response` / `INFO` / `text` success case and one explicit unsupported-record-type failure case without absorbing provider transport or phase handoff composition.
- Ordered scope 2 (`src-tauri/src/application/translation_instruction_builder/`): add the builder root, internal record-type rule matching on `source_entity_type` / `record_signature` / `field_name`, fixed `body_translation` phase code, `unit_key = translation_unit.extraction_key`, anchored dialogue instruction text, and an explicit module-boundary error for unsupported combinations.
- Ordered scope 3 (`src-tauri/src/application/mod.rs`): expose the builder through the application root only if the backend test or downstream application path needs a stable public import; keep `src-tauri/src/lib.rs`, gateway transport, provider selection, retry policy, output writing, and later phase composition out of this task.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_flow_mvp_regression -- --nocapture`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/application/translation_instruction_builder src-tauri/src/application/mod.rs src-tauri/tests`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- Representative record types can produce translation instructions through one stable backend path.
- Builder behavior stays inside the application layer owned scope and does not absorb provider transport or output writing.
- Regression fixtures continue to anchor representative record-type instruction building.

## Required Evidence

- Active plan updated with distill, design, implementation brief, and test plan outcomes.
- Backend tests or fixtures proving representative record-type instruction building.
- Validation results for backend lint, targeted tests, Sonar owned-scope issues, single-pass review, and full harness.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- Diagram updates are not expected unless implementation adds or changes repo-level structure or documented execution flow beyond the current Phase 3 task boundary.

## Outcome

- Added `src-tauri/src/application/translation_instruction_builder/mod.rs` with a stable application entrypoint `build_translation_instruction(TranslationUnitDto) -> Result<TranslationInstructionDto, String>`.
- Supported the anchored `dialogue_response` / `INFO` / `text` record type only, returning `phase_code=body_translation`, `unit_key=translation_unit.extraction_key`, the input `translation_unit`, and the agreed instruction text for the representative Skyrim dialogue case.
- Added explicit unsupported-record-type failure at the builder boundary and covered both success and failure paths in `src-tauri/tests/translation_instruction_builder_contract.rs`.
- Exposed the builder from `src-tauri/src/application/mod.rs` without widening into provider transport, retry, output writing, or phase handoff composition.
- Validation passed: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_instruction_builder_contract -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_flow_mvp_regression -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `sonar-scanner`, owned-scope Sonar open issues `0`, single-pass review `pass`, and `python3 scripts/harness/run.py --suite all`.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were reviewed and did not require updates for this change. `4humans` diagram sync was not required because repo-level structure and documented process flow stayed unchanged.
