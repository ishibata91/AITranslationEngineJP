- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-I02` from `tasks/phase-2/tasks/P2-I02.yaml` by adding master dictionary storage and query behind the agreed lookup port without widening into importer, persona, or UI work.
- task_id: P2-I02
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-I02`.
- Persist and query master dictionary foundation data behind the lookup port fixed in `P2-C02`, while preserving the rebuild validation anchor fixed in `P2-V01`.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-2/phase.yaml`
- `tasks/phase-2/tasks/P2-I02.yaml`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`
- `docs/exec-plans/completed/2026-04-04-p2-i01-xtranslator-importer.md`

## Owned Scope

- `src-tauri/src/application/dictionary_query/`
- `src-tauri/src/infra/dictionary_repository/`

## Out Of Scope

- xTranslator parsing and import adaptation
- persona storage or persona rebuild
- UI observation and rendering
- translation-phase policy beyond the stable lookup contract

## Dependencies / Blockers

- `P2-I02` depends on `P2-C02` and `P2-V01`.
- The implementation must remain compatible with the dictionary import result and reusable-entry contract established by `P2-B1` and consumed by `P2-I01`.

## Parallel Safety Notes

- This task stays parallel-safe with the other `P2-B2` tasks only while it remains inside dictionary storage and query ownership.
- Shared dictionary contracts and rebuild fixtures are upstream inputs and must not be redesigned here.

## UI

- N/A. `P2-I02` は backend storage と query boundary だけを扱う。
- importer command、persona observation、dictionary observation UI は後続 task に残し、この task では transport surface を増やさない。

## Scenario

- `P2-I01` の import 成功結果 `DictionaryImportResultDto { dictionaryName, sourceType, entries }` を、master dictionary 保存の唯一の入力として受ける。1 回の import 結果は 1 件の `MASTER_DICTIONARY` と、その配下の `MASTER_DICTIONARY_ENTRY` 群として永続化し、entry の重複、空文字、先頭末尾空白は importer 出力のまま保持する。
- dictionary rebuild validation は、fixture から import した結果を永続化してから `DictionaryLookupPort` 経由で再読込する流れへ置き換える。validation が固定する request / response shape と shared reusable-entry snapshot は `P2-V01` から変えず、in-memory grouping だけを persistence-backed path に差し替える。
- lookup は `DictionaryLookupRequest { source_texts }` の exact match batch として扱う。request に dictionary selector が存在しないため、query 対象は保存済みの全 `MASTER_DICTIONARY` を横断し、各 requested `source_text` ごとに 1 件の `candidate_group` を返す。未一致語は `candidates = []` で返す。
- `candidate_groups` は request の並び順を保持する。各 group 内の `candidates` は、同一 dictionary 内では import 時の entry 順を保ち、複数 dictionary にまたがる場合は earlier persisted dictionary を先に返すことで deterministic に固定する。

## Logic

- `src-tauri/src/application/dictionary_query/` は application orchestration に限定し、`DictionaryImportResultDto` を master dictionary として保存する use case と、既存 `DictionaryLookupPort` を実装する lookup use case を持つ。ここで新しい transport DTO や UI 向け read model は追加せず、`P2-C02` の request / result shape をそのまま守る。
- application 側には SQLite 非依存の repository trait を置き、責務を `save imported master dictionary batch` と `lookup reusable entries by exact source_texts` に絞る。dictionary selector、translation-phase ranking、provider hint、persona 連携は trait に持ち込まない。
- `src-tauri/src/infra/dictionary_repository/` はその trait の SQLite 実装として、`MASTER_DICTIONARY` と `MASTER_DICTIONARY_ENTRY` への DML と transaction を担当する。保存は `MASTER_DICTIONARY` 親行の insert と配下 entry の insert を 1 transaction にまとめ、途中失敗では rollback する。repository は schema 実行を担当せず、`sqlx` migration 適用は既存 bootstrap 初期化責務に残す。
- lookup SQL は `source_text` の完全一致だけを使い、trim、case-fold、dedupe、fuzzy match を行わない。実装都合で request 内の重複語を問い合わせ前にまとめてもよいが、返却時には request 順へ復元し、同じ `source_text` が request に複数回現れた場合も group 順は入力に従わせる。
- deterministic order は SQL 側で明示し、`MASTER_DICTIONARY.id ASC`、`MASTER_DICTIONARY_ENTRY.id ASC` を基準に candidate 列を返す。これにより importer が保持した entry 順と、later import を後段に積む batch 順を SQLite 上でも再現する。`built_at` は観測用 metadata として保持してよいが、lookup order の基準には使わない。
- `MASTER_DICTIONARY` / `MASTER_DICTIONARY_ENTRY` の naming と ER 上の column 意味は `docs/er.md` に合わせる。query 実装は foundation storage 専用に留め、`JOB_DICTIONARY_ENTRY` や後続 translation-phase の再利用 policy と混線させない。

## Implementation Plan

### ordered_scope

1. Validation anchor (`src-tauri/tests/validation/dictionary-rebuild/`, `src-tauri/tests/dictionary_rebuild_validation.rs`)

- Replace the current in-memory grouping path with a persistence-backed rebuild flow that imports one fixture, saves it as master dictionary foundation data, then re-reads through `DictionaryLookupPort`.
- Keep the existing `P2-V01` request shape, reusable-entry snapshot, duplicate preservation, whitespace preservation, and candidate-group ordering semantics unchanged. Only the backing path changes from in-memory grouping to persisted lookup.

2. Application dictionary query (`src-tauri/src/application/dictionary_query/`, `src-tauri/src/application/mod.rs`)

- Add the application-owned repository trait and use cases needed to save one imported master dictionary batch and to serve the existing lookup port without changing `DictionaryLookupRequest` / `DictionaryLookupResult`.
- Encode the deterministic result rules here: request order is preserved exactly, repeated `source_text` values are restored to input order, and lookup remains an exact-match batch across all persisted master dictionaries with no selector, trim, dedupe, or ranking policy.

3. Infra dictionary repository (`src-tauri/src/infra/dictionary_repository/`, `src-tauri/src/infra/mod.rs`, `src-tauri/migrations/` if schema is still absent)

- Add the SQLite repository implementation that inserts one `MASTER_DICTIONARY` parent row plus ordered `MASTER_DICTIONARY_ENTRY` rows in a single transaction, rolling back on failure.
- Add exact-match lookup SQL that returns candidates ordered by `MASTER_DICTIONARY.id ASC` then `MASTER_DICTIONARY_ENTRY.id ASC` so persisted dictionary order and per-dictionary entry order stay deterministic.
- Keep repository ownership limited to DML and transaction handling. If `MASTER_DICTIONARY` / `MASTER_DICTIONARY_ENTRY` tables are not yet present in the existing versioned migration set, add the thinnest migration under `src-tauri/migrations/`; migration execution itself remains in the existing bootstrap path and must not move into the repository.

### owned_scope

- `src-tauri/src/application/dictionary_query/`
- `src-tauri/src/infra/dictionary_repository/`
- `src-tauri/tests/validation/dictionary-rebuild/`
- `src-tauri/tests/dictionary_rebuild_validation.rs`
- `src-tauri/migrations/` only if the current versioned migration layout does not already create `MASTER_DICTIONARY` and `MASTER_DICTIONARY_ENTRY`

### required_reading

- `docs/exec-plans/active/2026-04-04-p2-i02-master-dictionary-storage-query.md` の `UI` / `Scenario` / `Logic`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`
- `docs/exec-plans/completed/2026-04-04-p2-i01-xtranslator-importer.md`
- `docs/er.md` の `MASTER_DICTIONARY` / `MASTER_DICTIONARY_ENTRY` 定義
- `src-tauri/src/application/ports/dictionary_lookup/mod.rs`
- `src-tauri/tests/validation/dictionary-rebuild/mod.rs`
- `src-tauri/migrations/0001_execution_cache_base.sql`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- master dictionary foundation data is persisted behind the agreed lookup boundary
- stable lookup queries can read persisted reusable-entry candidates through one path
- rebuild validation compatibility remains intact

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Tests or fixtures proving persistence and lookup behavior stay aligned with the Phase 2 dictionary contract.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/backend-dictionary-storage-query-class-diagram.d2`
- `4humans/diagrams/structures/backend-dictionary-storage-query-class-diagram.svg`
- `4humans/diagrams/processes/backend-dictionary-storage-query-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-dictionary-storage-query-sequence-diagram.svg`

## Outcome

- Added `src-tauri/src/application/dictionary_query/mod.rs` and module wiring in `src-tauri/src/application/mod.rs`, so master dictionary import results can be saved and queried behind the existing `DictionaryLookupPort` without widening the request / result DTO shape.
- Added `src-tauri/src/infra/dictionary_repository/mod.rs`, `src-tauri/src/infra/mod.rs`, and `src-tauri/migrations/0002_master_dictionary_foundation.sql`, so `MASTER_DICTIONARY` / `MASTER_DICTIONARY_ENTRY` persist in one SQLite transaction and exact-match lookup returns candidates in deterministic `MASTER_DICTIONARY.id ASC` then `MASTER_DICTIONARY_ENTRY.id ASC` order.
- Updated `src-tauri/tests/dictionary_rebuild_validation.rs` and `src-tauri/tests/validation/dictionary-rebuild/mod.rs`, so the rebuild anchor now verifies import -> persist -> lookup, repeated request order, whitespace-sensitive exact match, empty-request rejection, and rollback on partial entry insert failure.
- Added `4humans` review diagrams for the storage/query slice and its import -> persist -> lookup flow, so the new backend path is reviewable as both structure and sequence.
- Validation passed with `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home python3 scripts/harness/run.py --suite all`, `sonar-scanner`, and `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` with `openIssueCount = 0`.
- Single-pass review returned `reroute`; the reroute fixes above were applied and, per lane contract, no second review was run.
- `4humans` updates were not required because no new long-lived debt item, quality status change, or diagram boundary change was introduced by this task.
