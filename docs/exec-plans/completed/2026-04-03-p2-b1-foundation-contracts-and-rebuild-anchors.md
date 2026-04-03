- workflow: impl
- status: completed
- lane_owner: codex
- scope: Fix the Phase 2 batch `P2-B1` contracts and rebuild validation anchors for master dictionary and persona foundation data.
- task_id: P2-B1
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement batch `P2-B1` from `tasks/phase-2/phase.yaml`.
- Fix the contract boundary for dictionary import, dictionary lookup and reuse, and master-persona versus job-persona storage split before broad foundation-data implementation starts.
- Add the first validation anchors for dictionary rebuild and persona rebuild without absorbing UI observation implementation or later integration wiring.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/er.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-2/phase.yaml`
- `tasks/phase-2/tasks/P2-C01.yaml`
- `tasks/phase-2/tasks/P2-C02.yaml`
- `tasks/phase-2/tasks/P2-C03.yaml`
- `tasks/phase-2/tasks/P2-V01.yaml`
- `tasks/phase-2/tasks/P2-V02.yaml`

## Owned Scope

- `src-tauri/src/application/dto/dictionary_import/`
- `src-tauri/src/application/ports/dictionary_lookup/`
- `src-tauri/src/application/dto/persona_storage/`
- `src-tauri/src/application/ports/persona_storage/`
- `src-tauri/tests/validation/dictionary-rebuild/`
- `src-tauri/tests/validation/persona-rebuild/`

## Out Of Scope

- dictionary storage implementation detail beyond the stable import and lookup boundaries
- persona generation prompt or translation-phase policy
- UI observation implementation and rendering
- cross-layer integration wiring reserved for later Phase 2 implementation and integration tasks

## Dependencies / Blockers

- `P2-V01` depends on the stable output of `P2-C01` and `P2-C02`.
- `P2-V02` depends on the stable output of `P2-C03`.
- Existing Phase 1 importer, DTO, and verification conventions should remain compatible where foundation-data contracts reuse the same backend layering patterns.

## Parallel Safety Notes

- `P2-C01`, `P2-C02`, `P2-C03`, `P2-V01`, and `P2-V02` are grouped into the same batch because they can be fixed before later backend and frontend slices start, but storage details, UI concerns, or translation-phase policy must not leak into this batch.
- Verification anchors must stay small and reusable so later implementation and integration tasks can consume them without fixture-schema redesign.

## UI

- N/A. `P2-B1` は backend application contract と rebuild validation anchor の固定だけを扱う。
- 後続 UI 観測 task が参照できるのは backend-owned DTO / port の項目名と意味だけであり、screen state、Gateway command / event、polling、rendering、error presentation はこの batch で決めない。

## Scenario

- `xTranslator` 形式入力からの辞書 rebuild は、`dictionary_import` DTO で import source を受け取り、`MASTER_DICTIONARY` 再構築に必要な基盤メタデータと正規化済み `source_text` / `dest_text` entry 群だけを返す。保存先 ID、UI 表示補助、翻訳フェーズ都合の情報は import contract に入れない。
- 後続の単語翻訳フェーズや本文翻訳フェーズから見る辞書再利用は、`dictionary_lookup` port を通じて 1 つの stable reusable-entry shape を読む。lookup は辞書 rebuild 結果と同じ語対境界を共有するが、どの候補を採用するかという translation-phase policy は持ち込まない。
- ベースゲーム NPC 由来の rebuild は `master persona` 側 contract に保存され、translation job 中に生成した NPC persona は `job persona` 側 contract に保存される。両者は同じ NPC 属性列を参照しても保持スコープと参照起点が異なるため、同一 save/read DTO に畳み込まない。
- `dictionary-rebuild` と `persona-rebuild` の validation anchor は、それぞれ 1 つの最小 fixture から deterministic な rebuild snapshot を確認する。fixture が依存してよいのは contract 境界だけであり、UI 観測結果や後続 integration wiring は前提に含めない。

## Logic

- `src-tauri/src/application/dto/dictionary_import/` は xTranslator import request / result DTO の所有境界とし、request は import source 識別に閉じ、result は foundation rebuild に必要な dictionary metadata と reusable entry 群へ閉じる。runtime 生成 ID、永続化内部キー、validation fixture を不安定にする時刻依存値は stable contract の必須項目に含めない。
- `src-tauri/src/application/ports/dictionary_lookup/` は後続 translation workflow が参照する application port とし、入力は 1 件以上の `source_text` 群で表現できる形に固定する。出力は exact `source_text` に紐づく reusable entry 候補列を deterministic に返し、優先順位決定、provider hint、prompt 文脈、UI 用ラベルは port 外へ残す。
- `dictionary_lookup` port が返す reusable entry shape は import result の entry shape と整合させ、少なくとも `source_text` と `dest_text` の語対境界を共有する。これにより master dictionary rebuild と後続再利用が別 module でも同じ contract を前提にできるようにする。
- `src-tauri/src/application/dto/persona_storage/` と `src-tauri/src/application/ports/persona_storage/` は `MASTER_PERSONA` / `MASTER_PERSONA_ENTRY` 用の基盤 contract と、`JOB_PERSONA_ENTRY` 用の job-local contract を分離して持つ。共有してよいのは `npc_form_id`、NPC 属性、`persona_text` などの persona 実体だけであり、job identifier や foundation metadata は必要な側にだけ置く。
- `src-tauri/tests/validation/dictionary-rebuild/` の anchor は import result と lookup port が同じ reusable entry contract を共有していることを確認し、rerun しても同じ snapshot を組み立てられることを固定する。`src-tauri/tests/validation/persona-rebuild/` の anchor は master rebuild snapshot と job-local persona snapshot が相互代用できないことを固定し、同一 NPC 属性でも保持先 contract が混線しないことを示す。
- この batch が固定するのは boundary semantics と rebuild validation anchor だけである。master dictionary と job-generated dictionary reuse の採用順、persona generation prompt、UI observation wiring は後続 task へ明示的に残す。

## Implementation Plan

- Ordered scope 1 (`P2-C01`): add `src-tauri/src/application/dto/dictionary_import/` and extend `src-tauri/src/application/dto/mod.rs` with foundation-oriented dictionary import request / result DTOs plus one reusable-entry DTO. `dictionary_name` and `source_type` stay opaque `String` fields in this batch so the contract does not bake in derivation policy.
- Ordered scope 2 (`P2-C02`): add `src-tauri/src/application/ports/dictionary_lookup/` and minimal module wiring so later translation workflow can query exact `source_text` groups and receive deterministic reusable-entry candidates that share the same `source_text` / `dest_text` boundary as the dictionary import result.
- Ordered scope 3 (`P2-C03`): add `src-tauri/src/application/dto/persona_storage/` and `src-tauri/src/application/ports/persona_storage/` with split master-persona and job-persona save/read contracts. `persona_name` also stays an opaque `String`, while job identifier fields remain only on the job-local side.
- Ordered scope 4 (`P2-V01`): add `src-tauri/tests/validation/dictionary-rebuild/` plus the thinnest top-level test loader needed for Cargo discovery, and fix one deterministic fixture / snapshot pair that proves dictionary import and dictionary lookup share the same reusable-entry contract.
- Ordered scope 5 (`P2-V02`): add `src-tauri/tests/validation/persona-rebuild/` plus the thinnest top-level test loader needed for Cargo discovery, and fix one deterministic fixture / snapshot pair that proves master-persona rebuild data and job-persona data remain non-substitutable even with matching NPC attributes.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
  - `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`
  - `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- `P2-C01`, `P2-C02`, and `P2-C03` each land on one stable boundary inside their owned scope.
- `P2-V01` proves dictionary rebuild behavior can be validated independently.
- `P2-V02` proves persona rebuild behavior can be validated independently.
- The batch leaves no UI or integration-only behavior embedded in contract or verification scope.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests or fixtures proving dictionary and persona rebuild anchors.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- No diagram update expected unless implementation changes codebase boundaries or execution flow beyond the task catalog scope.

## Outcome

- Added `src-tauri/src/application/dto/dictionary_import/` with `DictionaryImportRequestDto`, `DictionaryImportResultDto`, and `ReusableDictionaryEntryDto` so Phase 2 dictionary foundation work has one backend-owned import boundary. The request now preserves both `source_type` and a stable `source_file_path` input handle for later xTranslator ingestion without widening the contract.
- Added `src-tauri/src/application/dto/persona_storage/`, `src-tauri/src/application/ports/dictionary_lookup/`, and `src-tauri/src/application/ports/persona_storage/`, then wired them through `src-tauri/src/application/dto/mod.rs`, `src-tauri/src/application/ports/mod.rs`, and `src-tauri/src/application/mod.rs`. Master persona and job persona remain split, and job-local save requests now carry `source_type` for later persistence alignment.
- Added rebuild validation entrypoints and fixtures under `src-tauri/tests/dictionary_rebuild_validation.rs`, `src-tauri/tests/persona_rebuild_validation.rs`, and `src-tauri/tests/validation/` so dictionary rebuild and persona rebuild each have deterministic snapshot anchors that later Phase 2 tasks can reuse unchanged.
- Added camelCase transport-boundary checks in `src-tauri/tests/import_job_transport_contract.rs` and the rebuild validation modules so `sourceType`, `sourceFilePath`, and `jobId` remain pinned at the DTO boundary.
- Validation passed after the reroute fixes: `python3 scripts/harness/run.py --suite structure`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test import_job_transport_contract -- --nocapture`, `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `sonar-scanner`, and `python3 scripts/harness/run.py --suite all`.
- Single-pass review initially returned `reroute`; the reroute fixes above were applied and per lane contract no second review was run.
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` could not be completed in this sandbox because `sonar auth login` is not configured. `sonar-scanner` itself completed successfully, so the remaining gap is Sonar CLI authentication rather than scanner execution.
