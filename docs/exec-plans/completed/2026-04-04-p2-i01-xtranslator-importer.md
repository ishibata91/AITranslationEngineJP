- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-I01` from `tasks/phase-2/phase.yaml` by adding the xTranslator-backed master dictionary import path without widening into storage or UI work.
- task_id: P2-I01
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-B2` one task at a time, starting with `P2-I01`.
- Accept xTranslator input for master dictionary ingestion behind the Phase 2 import contract fixed in `P2-C01` and the rebuild validation anchor fixed in `P2-V01`.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-2/phase.yaml`
- `tasks/phase-2/tasks/P2-I01.yaml`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`

## Owned Scope

- `src-tauri/src/application/dictionary_import/`
- `src-tauri/src/infra/xtranslator_importer/`

## Out Of Scope

- dictionary lookup query
- dictionary storage implementation
- persona storage or persona build
- UI observation

## Dependencies / Blockers

- `P2-I01` depends on `P2-C01` and `P2-V01`, both completed in `2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`.
- The importer must preserve the stable dictionary import DTO boundary and remain compatible with the dictionary rebuild validation anchor.

## Parallel Safety Notes

- `P2-I01` is parallel-safe with the other `P2-B2` implementation tasks only while it stays inside import-path ownership.
- Parsing, normalization, or adapter details may be implemented here, but lookup, storage, persona, and UI concerns must remain outside this task.

## UI

- N/A. `P2-I01` は backend import boundary と parser adapter だけを扱う。
- Gateway command、screen state、error presentation はこの task で決めない。

## Scenario

- xTranslator master dictionary import は `DictionaryImportRequestDto` の `sourceType` と `sourceFilePath` を受け、初期実装では `sourceType = "xtranslator-sst"` の 1 系統だけを受理する。lookup、storage、persona、UI は関与しない。
- application use case は request を対応 importer へ委譲し、`DictionaryImportResultDto { dictionaryName, sourceType, entries }` を返す。`sourceType` は request と同じ transport 値を保持し、entry は shared reusable-entry shape をそのまま使う。
- import 成功時の result はそのまま dictionary rebuild validation へ流せることを前提とし、同一 `sourceText` の複数候補と entry 順序を保持する。rebuild 前に dedupe、sort、storage ID 補完は行わない。
- unsupported sourceType、ファイル読取失敗、xTranslator payload の parse failure は import boundary で `Err(String)` として返し、部分成功や UI 向け recovery 情報は持ち込まない。

## Logic

- `src-tauri/src/application/dictionary_import/` は xTranslator dictionary import の use case 境界を持ち、request DTO を受けて sourceType ごとの importer port へ dispatch する。初期 scope では `xtranslator-sst` だけを実装し、他 sourceType は unsupported error に閉じる。
- application 層の責務は boundary orchestration のみに限定する。storage repository、lookup query、gateway command はここで追加せず、戻り値は `DictionaryImportResultDto` に固定する。
- `src-tauri/src/infra/xtranslator_importer/` は parser / adapter 境界を持ち、指定 path から xTranslator source を読み、辞書名と reusable entry 列へ変換する。現時点の実装は deterministic な file stem を `dictionaryName` に使い、metadata 優先規則の固定は後続 task へ残す。
- entry 変換は `source_text` / `dest_text` の語対だけに閉じ、先頭 / 末尾空白を含む文字列を勝手に trim・正規化しない。xTranslator 側の表記差や space-sensitive な語対を壊さないことを優先する。
- importer は reusable entry の重複と並び順を保持する。`DictionaryLookupPort` の shared snapshot 互換性を優先し、application / infra のどちらでも dedupe、sort、merge を行わない。
- 実装時の受け入れは、xTranslator fixture 1 件から `sourceType = "xtranslator-sst"` と shared reusable-entry snapshot を再現できることを基準にする。新しい command surface や storage wiring は importer task の前提にしない。

## Implementation Plan

### 1. Validation anchor (`src-tauri/tests/validation/dictionary-rebuild/`, `src-tauri/tests/dictionary_rebuild_validation.rs`, `src-tauri/tests/xtranslator_importer.rs`)

- Add one xTranslator fixture and the thinnest acceptance-oriented tests that pin the Phase 2 import behavior before broad wiring starts.
- The validation anchor must prove `sourceType = "xtranslator-sst"` can reproduce the shared reusable-entry snapshot without trimming, dedupe, reordering, or storage-side enrichment.
- Keep the test surface backend-only. Do not introduce a gateway command, persistence dependency, or UI recovery contract in this task.

### 2. Infra parser and adapter (`src-tauri/src/infra/xtranslator_importer/`, `src-tauri/src/infra/mod.rs`)

- Add a file-based xTranslator importer module that reads `sourceFilePath`, parses SST payload data, derives a deterministic `dictionaryName`, and returns reusable entries in input order with duplicates preserved.
- Convert file read failure and xTranslator parse failure into `Err(String)` at the importer boundary. Do not add partial-success payloads, lookup semantics, or storage concerns.

### 3. Application orchestration (`src-tauri/src/application/dictionary_import/`, `src-tauri/src/application/mod.rs`)

- Add the dictionary import use case boundary that accepts `DictionaryImportRequestDto`, dispatches only `source_type = "xtranslator-sst"` to the xTranslator importer, and returns `DictionaryImportResultDto` with the transport `source_type` unchanged.
- Reject unsupported `source_type` values at the application boundary, and keep orchestration separate from infra parsing. Do not add repository writes, lookup wiring, persona logic, or new command surface.

### 4. Validation

- `python3 scripts/harness/run.py --suite structure`
- `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test xtranslator_importer -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- dictionary input can be ingested without redesigning the import boundary
- dictionary import path can feed rebuild validation

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Tests or fixtures proving the importer satisfies the existing rebuild validation anchor.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- No diagram update expected unless implementation changes codebase boundaries or execution flow beyond the task catalog scope.

## Outcome

- Added `src-tauri/src/application/dictionary_import/mod.rs` with `DictionaryImporter` and `ImportDictionaryUseCase`, so `DictionaryImportRequestDto` を `source_type = "xtranslator-sst"` の narrow boundary で受けて backend import path へ委譲できるようにした。
- Added `src-tauri/src/infra/xtranslator_importer/mod.rs` and module wiring in `src-tauri/src/application/mod.rs` / `src-tauri/src/infra/mod.rs`, so SST fixture から `dictionaryName`、`sourceType`、ordered `ReusableDictionaryEntryDto` を返す file-based importer が動くようになった。
- Added `src-tauri/tests/xtranslator_importer.rs`, `src-tauri/tests/support/xtranslator_fixture.rs`, and `src-tauri/tests/validation/dictionary-rebuild/fixtures/xtranslator-shared-reusable-entry.sst`, then updated `src-tauri/tests/dictionary_rebuild_validation.rs` and `src-tauri/tests/validation/dictionary-rebuild/mod.rs` so xTranslator import path と shared reusable-entry snapshot の互換が固定された。
- Validation passed with repo-owned commands: `python3 scripts/harness/run.py --suite structure`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test xtranslator_importer -- --nocapture`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `sonar-scanner`, and `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` with `openIssueCount = 0`.
- `python3 scripts/harness/run.py --suite all` remains blocked in this sandbox because `npm run lint:rust:clippy` uses the default Cargo home under `/Users/iorishibata/.cargo`, which is not writable here. The underlying repo-owned Cargo checks passed with `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home`.
- Single-pass implementation review returned `pass`.
