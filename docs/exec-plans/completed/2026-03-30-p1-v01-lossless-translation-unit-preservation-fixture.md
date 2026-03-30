- workflow: impl
- status: completed
- lane_owner: codex
- scope: Create a stable lossless translation-unit preservation fixture and verification checks for Phase 1.
- task_id: P1-V01
- task_catalog_ref: tasks/phase-1/tasks/P1-V01.yaml
- parent_phase: phase-1

## Request Summary

- Implement task `P1-V01` to add one stable verification fixture and checks that prove normalized translation-unit preservation remains lossless for Phase 1 needs.
- Keep the task inside `src-tauri/tests/fixtures/translation-unit-lossless/` and adjacent verification coverage without absorbing UI scenarios or provider execution behavior.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/er.md`
- `tasks/phase-1/tasks/P1-V01.yaml`
- `docs/exec-plans/completed/2026-03-30-p1-c01-translation-unit-canonical-contract.md`
- `src-tauri/tests/translation_unit_contract.rs`
- `src-tauri/tests/xedit_export_importer.rs`

## Owned Scope

- `src-tauri/tests/fixtures/translation-unit-lossless/`
- `src-tauri/tests/translation_unit_contract.rs`
- `src-tauri/tests/xedit_export_importer.rs`

## Out Of Scope

- UI scenarios
- provider execution behavior

## Dependencies / Blockers

- `P1-C01` canonical translation-unit contract must remain the upstream boundary for source-side fields.
- The task catalog acceptance anchor references translated text and output status even though the current canonical translation-unit contract fixed by `P1-C01` only includes pre-job source-side fields plus deterministic ordering metadata.

## Parallel Safety Notes

- `P1-C02` and `P1-V02` are marked parallel-safe, but this task must not absorb job-state modeling or import-to-job scenario wiring while fixing the verification anchor.
- Fixture schema drift would invalidate downstream checks and later integration work.

## UI

- N/A. This task adds backend verification fixtures and checks only, and must not introduce screen, route, or interaction changes.

## Scenario

- Success flow is `fixture source -> normalized translation-unit expectation -> downstream preservation expectation`, where one stable fixture proves the same unit can be reconstructed for Phase 1 output needs without UI or provider execution.
- The fixture should represent at least one normalized translation-unit case with explicit `form_id`, `editor_id`, `record_signature`, `field_name`, `source_text`, and deterministic extraction identity, while also pinning the paired translated-text and output-status expectation needed by later output-facing checks.
- Blank `editor_id` remains a valid preserved case and should be included when it helps prove lossless behavior, but the fixture must stay minimal and avoid turning into a broad importer matrix.
- Verification should consume the fixture from `src-tauri/tests/fixtures/translation-unit-lossless/` rather than inlining ad hoc JSON or expected rows inside each test, so later Phase 1 tasks can reuse one agreed preservation anchor.

## Logic

- `P1-C01` remains the canonical product contract for normalized translation units: source-side identity fields plus deterministic `sort_key`. `P1-V01` must not widen that product contract to include translated text or output status.
- The fixture may carry two aligned layers: a canonical normalized translation-unit snapshot and a preservation expectation snapshot that adds translated text plus output status as downstream verification data. This keeps product ownership unchanged while satisfying the verification anchor.
- Checks should prove lossless preservation by asserting that canonical fields remain stable and that the fixture-owned translated-text and output-status values stay correctly associated with the same normalized unit identity.
- The task should prefer stable serialized fixture data and thin test helpers over new domain types or importer-specific one-off structs, because fixture-shape drift is the main shared risk for downstream work.

## Implementation Plan

- Ordered scope 1: Add one reusable serialized fixture under `src-tauri/tests/fixtures/translation-unit-lossless/` that contains the canonical normalized translation-unit snapshot plus aligned preservation expectation data for translated text and output status without redefining product types.
- Ordered scope 2: Update `src-tauri/tests/translation_unit_contract.rs` to load the fixture and prove the canonical translation-unit identity fields remain lossless while preserving the fixture-owned translated-text and output-status association for the same unit.
- Ordered scope 3: Update `src-tauri/tests/xedit_export_importer.rs` only as needed to reuse the same fixture-backed normalized translation-unit expectation or to prove blank `editor_id` preservation stays compatible with the shared anchor, without turning the file into a broader importer matrix.
- Validation commands:
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_unit_contract -- --nocapture`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test xedit_export_importer -- --nocapture`
  - `powershell -File scripts/harness/run.ps1 -Suite all`

## Acceptance Checks

- Add one stable fixture under `src-tauri/tests/fixtures/translation-unit-lossless/` that downstream checks can consume without importer-specific ad hoc inline JSON.
- Verification coverage proves the normalized translation-unit preservation anchor remains lossless for `FormID`, `EditorID`, record signature, field kind, source text, translated text, and output status at the agreed Phase 1 boundary.
- Validation commands remain narrow enough to exercise fixture-backed preservation checks without expanding into UI or provider behavior.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests or fixtures proving the lossless preservation fixture and checks.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added `src-tauri/tests/fixtures/translation-unit-lossless/lossless-translation-unit-preservation.json` as the shared preservation anchor with one normalized translation-unit snapshot plus aligned translated-text and output-status expectation data.
- Updated `src-tauri/tests/translation_unit_contract.rs` to load the shared fixture and prove canonical identity fields stay lossless while downstream preservation data remains aligned by `extraction_key`.
- Updated `src-tauri/tests/xedit_export_importer.rs` to consume the same fixture-backed normalized expectation for the blank-`editor_id` path and hardened the test temp-path helper with a unique counter so the suite stays stable under parallel `cargo test --all-features`.
- Validation passed: `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_unit_contract -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test xedit_export_importer -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `powershell -File scripts/harness/run.ps1 -Suite all`, and `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths ...` returned `openIssueCount: 0`.
- Single-pass review result: `pass`.
- `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` were not changed because this task added a stable verification anchor and removed a test flake without materially changing repository-level debt posture.
