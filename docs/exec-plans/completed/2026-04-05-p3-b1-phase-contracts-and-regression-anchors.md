- workflow: impl
- status: completed
- lane_owner: codex
- scope: Fix the Phase 3 batch `P3-B1` contracts and regression anchors for translation instruction building, embedded-element preservation, and translation phase handoff.
- task_id: P3-B1
- task_catalog_ref: tasks/phase-3/phase.yaml
- parent_phase: phase-3

## Request Summary

- Implement batch `P3-B1` from `tasks/phase-3/phase.yaml`.
- Fix stable application boundaries for translation instruction building, embedded-element preservation, and word/persona/body translation phase handoff before broader Translation Flow MVP implementation starts.
- Add regression anchors for embedded-element preservation and one representative translation-flow scenario without absorbing provider-specific logic, retry policy, output writing, or preview UI behavior.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-3/phase.yaml`
- `tasks/phase-3/tasks/P3-C01.yaml`
- `tasks/phase-3/tasks/P3-C02.yaml`
- `tasks/phase-3/tasks/P3-C03.yaml`
- `tasks/phase-3/tasks/P3-V01.yaml`
- `tasks/phase-3/tasks/P3-V02.yaml`

## Owned Scope

- `src-tauri/src/application/dto/translation_instruction/`
- `src-tauri/src/application/dto/embedded_element_policy/`
- `src-tauri/src/application/dto/translation_phase_handoff/`
- `src-tauri/tests/regression/embedded-elements/`
- `src-tauri/tests/regression/translation-flow-mvp/`

## Out Of Scope

- provider adapters and provider-specific prompt assembly
- output writer policy and retry policy
- preview UI layout and frontend observation
- later composition and integration wiring reserved for `P3-B2` and `P3-B3`

## Dependencies / Blockers

- `P3-V01` depends on the stable output of `P3-C02`.
- `P3-V02` depends on the stable output of `P3-C01` and `P3-C03`.
- Existing Phase 2 dictionary lookup and persona storage contracts should remain reusable where Phase 3 handoff contracts depend on them.

## Parallel Safety Notes

- `P3-C01`, `P3-C02`, `P3-C03`, `P3-V01`, and `P3-V02` are grouped because they fix backend-owned contracts and regression anchors before provider-facing implementation and end-to-end composition begin.
- The batch must keep contract boundaries free from provider policy, output formatting, UI preview state, or phase-internal helper detail that would reduce independence of later tasks.

## UI

- N/A. `P3-B1` fixes backend application contracts and regression anchors only.

## Scenario

- `embedded-elements` の regression anchor は、`dialogue_response.text` を題材にした 1 件の最小 fixture で固定する。source text は angle-bracket placeholder を 1 つ以上含み、translation flow 内で本文翻訳対象の自然文も同居させるが、provider 応答、retry、output writer の都合は fixture に持ち込まない。
- 代表 Translation Flow MVP scenario は NPC 対話 1 件を基準にし、`dialogue_response.text` の `TranslationUnitDto`、record-type aware instruction、word translation 由来の reusable term 群、job-local persona 1 件、embedded-element preservation expectation を 1 つの deterministic fixture 系列でつなぐ。代表例は `<Alias=Player>` を含む挨拶文とし、body translation が upstream handoff を同時に消費する入口だけを固定する。
- `P3-V01` は preserved element 自体の survival を独立に確認する anchor とし、source text、保護対象 element 群、preserved された最終 text の対応だけを確認する。mask 文字列、provider prompt、内部 parser state は snapshot に含めない。
- `P3-V02` は instruction building と translation phase handoff が body translation 入力へ合流できることを確認する anchor とし、dictionary lookup の candidate 採用方針そのものや persona generation runtime の実行詳細は含めない。後続 task が fake 実装でも real 実装でも同じ fixture を再利用できる粒度に留める。

## Logic

- `src-tauri/src/application/dto/translation_instruction/` は provider-neutral な instruction contract root とし、入力側は既存 `TranslationUnitDto` をそのまま参照して record type 判定に必要な `source_entity_type` / `record_signature` / `field_name` を重複定義しない。出力側は少なくとも `phase_code`、unit を再特定できる stable key、`instruction_text` を持つ 1 unit 単位の instruction payload に閉じる。
- `translation_instruction` contract は instruction 文面の構成責務だけを持ち、provider selection、batch/single 実行方式、retry policy、preview 用装飾文言は持たない。record-type aware 分岐は `dialogue_response.text` を代表ケースとして固定し、他 record type 追加時も同じ payload shape を崩さない前提にする。
- `src-tauri/src/application/dto/embedded_element_policy/` は preserve すべき embedded element の public boundary を持ち、保護対象は source text 内での出現順を保つ ordered descriptor 群として表現する。descriptor には raw element text と unit 内で安定に参照できる識別子だけを持たせ、mask token 生成規則、parser 実装都合の中間値、output format 別 escape 仕様は contract 外に残す。
- `src-tauri/src/application/dto/translation_phase_handoff/` は word translation output と persona generation output が body translation input へ渡る provider-neutral な handoff root とする。dictionary 側は Phase 2 の `ReusableDictionaryEntryDto` を public root 経由で再利用し、persona 側は `JobPersonaEntryDto` を public root 経由で再利用して、同じ field 群を Phase 3 で再定義しない。
- `translation_phase_handoff` は phase ごとの内部 helper を露出せず、body translation が必要とする unit-scoped reusable term 群、optional な job-local persona、preserved embedded element 情報をまとめられる粒度に留める。master persona metadata、dictionary rebuild provenance、provider runtime handle は handoff に入れない。
- `src-tauri/tests/regression/embedded-elements/` と `src-tauri/tests/regression/translation-flow-mvp/` は Phase 2 validation と同じく fixture-backed snapshot 構成を踏襲し、top-level test loader は薄く保つ。`embedded-elements` は preservation contract 単体の snapshot、`translation-flow-mvp` は instruction payload と phase handoff payload から組み立てた representative body-translation input snapshot を固定して、後続 impl task の受け入れ基準にする。

## Implementation Plan

- Ordered scope 1 (`P3-C02`): add `src-tauri/src/application/dto/embedded_element_policy/` and define a provider-neutral ordered descriptor contract for preserved embedded elements only.
- Ordered scope 2 (`P3-C01`): add `src-tauri/src/application/dto/translation_instruction/` and define record-type aware instruction payloads for `dialogue_response.text` by reusing `TranslationUnitDto` and constraining output to `phase_code`, stable unit key, and `instruction_text`.
- Ordered scope 3 (`P3-C03`): add `src-tauri/src/application/dto/translation_phase_handoff/` and update `src-tauri/src/application/dto/mod.rs` so body-translation input can reuse `ReusableDictionaryEntryDto` and `JobPersonaEntryDto` through public roots without redefining equivalent Phase 3 fields.
- Ordered scope 4 (`P3-V01`): add `src-tauri/tests/regression/embedded-elements/` and `src-tauri/tests/embedded_elements_regression.rs`, then fix one fixture-backed snapshot that proves preserved embedded elements survive independently of provider or UI detail.
- Ordered scope 5 (`P3-V02`): add `src-tauri/tests/regression/translation-flow-mvp/` and `src-tauri/tests/translation_flow_mvp_regression.rs`, then fix one representative `<Alias=Player>` scenario snapshot that combines instruction payload and phase handoff payload into a stable body-translation input anchor.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `cargo test --manifest-path src-tauri/Cargo.toml --test embedded_elements_regression`
  - `cargo test --manifest-path src-tauri/Cargo.toml --test translation_flow_mvp_regression`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- `P3-C01`, `P3-C02`, and `P3-C03` each land on one stable contract inside the owned backend DTO scope.
- `P3-V01` proves embedded elements survive the protected translation-flow stages described by the preservation contract.
- `P3-V02` proves one representative Translation Flow MVP scenario can be anchored without depending on preview UI or provider retry policy.
- The batch leaves implementation-only or integration-only behavior out of contract and regression scope.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated regression tests or fixtures for embedded elements and representative translation flow.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- No diagram update expected unless implementation changes codebase boundaries or execution flow beyond the task catalog scope.

## Outcome

- Added `src-tauri/src/application/dto/embedded_element_policy/`, `src-tauri/src/application/dto/translation_instruction/`, and `src-tauri/src/application/dto/translation_phase_handoff/` so Phase 3 has provider-neutral contract roots for preserved embedded elements, record-type-aware instruction payloads, and body-translation handoff payloads.
- Reused `TranslationUnitDto`, `ReusableDictionaryEntryDto`, and `JobPersonaEntryDto` through public DTO roots instead of redefining equivalent Phase 3 shapes, then re-exported the new DTO roots from `src-tauri/src/application/dto/mod.rs`.
- Added fixture-backed regression anchors under `src-tauri/tests/regression/embedded-elements/` and `src-tauri/tests/regression/translation-flow-mvp/` plus top-level entrypoints `src-tauri/tests/embedded_elements_regression.rs` and `src-tauri/tests/translation_flow_mvp_regression.rs`.
- Validation passed: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `cargo test --manifest-path src-tauri/Cargo.toml --test embedded_elements_regression`, `cargo test --manifest-path src-tauri/Cargo.toml --test translation_flow_mvp_regression`, `sonar-scanner`, owned-scope Sonar open issues `0`, single-pass review `pass`, and `python3 scripts/harness/run.py --suite all`.
- `4humans` diagram sync was reviewed and not required because the change fixed DTO contract roots and regression anchors without changing the repo-level structure or documented process flows.
