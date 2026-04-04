- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-I04` from `tasks/phase-2/tasks/P2-I04.yaml` by persisting job-local persona data independently from master persona data without widening into master persona build or dictionary storage.
- task_id: P2-I04
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-I04`.
- Add backend persistence for job-local persona state behind the storage split contract defined by `P2-C03`.

## Decision Basis

- `tasks/phase-2/tasks/P2-I04.yaml`
- `tasks/phase-2/tasks/P2-C03.yaml`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`

## Owned Scope

- `src-tauri/src/application/job_persona/`
- `src-tauri/src/infra/job_persona_repository/`

## Out Of Scope

- master persona build
- dictionary storage
- UI observation and rendering

## Dependencies / Blockers

- `P2-I04` depends on `P2-C03`.
- The implementation must preserve the master-persona and job-persona storage split.

## Parallel Safety Notes

- `P2-I04` remains parallel-safe only while edits stay inside the owned backend scope and directly related tests, fixtures, and review diagrams.
- Shared storage contracts must be consumed as-is rather than redefined in this task.

## UI

- N/A. `P2-I04` は backend job persona persistence だけを扱う。
- Tauri command、screen state、UI error presentation、master persona observation はこの task で増やさない。

## Scenario

- job persona 永続化は、呼び出し側が `JobPersonaSaveRequestDto { job_id, source_type, entries }` を application の job persona use case へ渡し、use case が同じ `job_id` で `save_job_persona` と `read_job_persona` を直列実行して `JobPersonaReadResultDto { job_id, entries }` を返す 1 本の job-local path として扱う。read result は既存 contract に合わせて `source_type` を返さず、job-local identity は `job_id` のみに固定する。
- `entries` は対象 `job_id` の job persona 全体 snapshot として扱い、同じ `job_id` への再保存は既存の job persona 行集合を原子的に置き換える。別 job への append、master persona との merge、dictionary 由来データの補完は行わない。
- job persona entry は `npc_form_id`、`race`、`sex`、`voice`、`persona_text` だけを job-local data として保持し、`persona_name`、`npc_name`、master persona build input、UI 向け metadata はこの scenario に入れない。保存順は read 結果でも維持し、validation fixture が示す job persona transport shape を壊さない。
- `job_id` / `source_type` の空文字、空 entry 配列、entry 必須属性の空文字は invalid request として扱い、部分保存を許さない。未保存の `job_id` を読む要求も `Err(String)` で返し、空結果を成功扱いしない。

## Logic

- `src-tauri/src/application/job_persona/` は orchestration 境界として新設し、既存 `JobPersonaStoragePort` を使う job persona use case を持つ。use case は `JobPersonaSaveRequestDto` を入力に受け、`JobPersonaReadRequestDto { job_id }` を内部で組み立てて save/read を束ねるが、master persona module、dictionary module、UI read model には依存しない。
- application 層は transport contract をそのまま使いながら、少なくとも `job_id`、`source_type`、entry 配列、各 entry の required field が空でないことを保存前に検証する。snake_case 列名や SQL 都合は infra 内へ閉じ込め、application からは camelCase DTO / port contract だけを見せる。
- `src-tauri/src/infra/job_persona_repository/` は `JobPersonaStoragePort` 実装に責務を限定し、execution cache への接続、transaction、job-local DML、row-to-DTO mapping だけを持つ。migration・bootstrap・schema 初期化は既存の execution cache 初期化責務へ残し、この repository に混ぜない。
- save は単一 transaction で対象 `job_id` の既存 job persona rows を消してから request 順に insert する replacement semantics を採る。途中で 1 件でも失敗した場合は rollback して `Err(String)` を返し、retry 時に旧データと新データが混ざらないようにする。
- read は 1 つの `job_id` だけを問い合わせ、保存順を再現できる決定的 order で row 群を読み、`JobPersonaReadResultDto { job_id, entries }` に投影する。`source_type` は save 側でのみ受け取り、read result へは再露出しないことで `P2-C03` の split を守る。
- repository query は `JOB_PERSONA_ENTRY` 側の job-local rows だけを対象にし、`MASTER_PERSONA` や dictionary table を参照しない。同じ NPC 属性列を持っていても master/job の保存面を混線させないことを module 構成と query 対象の両方で固定する。

## Implementation Plan

### ordered_scope

1. Validation anchor (`src-tauri/tests/validation/persona-rebuild/`, `src-tauri/tests/persona_rebuild_validation.rs`)

- Extend the existing `P2-V02` anchor from transport non-substitutability into a deterministic job-persona persistence path that proves `save -> read` replacement semantics for a single `job_id`, preserves entry order, and keeps master/job storage separate.
- Reuse the current fixture vocabulary first. Add only the minimum execution-cache-backed fixture or repository-focused harness needed to pin invalid request rejection, missing-job read failure, and no cross-job leakage.

2. Application job persona orchestration (`src-tauri/src/application/job_persona/`, `src-tauri/src/application/mod.rs`)

- Add the application use case that validates `JobPersonaSaveRequestDto`, calls `JobPersonaStoragePort::save_job_persona`, builds `JobPersonaReadRequestDto { job_id }`, and returns the read result for the same `job_id`.
- Keep validation and orchestration in application only: reject empty `job_id`, `source_type`, empty `entries`, and empty required entry fields before storage calls. Do not add SQL, snake_case, master persona logic, or dictionary dependencies here.

3. Infra job persona repository (`src-tauri/src/infra/job_persona_repository/`, `src-tauri/src/infra/mod.rs`)

- Implement the `JobPersonaStoragePort` adapter with execution-cache connection handling, single-transaction delete-and-insert replacement for one `job_id`, deterministic read ordering, and row-to-DTO mapping limited to `JOB_PERSONA_ENTRY`.
- If the current versioned migration set still does not provision `JOB_PERSONA_ENTRY`, add the thinnest migration under `src-tauri/migrations/` in the same task, while keeping migration execution inside the existing bootstrap path rather than the repository module.
- Keep snake_case, SQL, transaction handling, and error formatting inside infra. Do not query `MASTER_PERSONA` or dictionary tables, and do not widen into schema redesign beyond the minimal missing table/index addition needed for `JOB_PERSONA_ENTRY`.

4. Boundary wiring and explicit non-scope decisions

- Wire only the new backend modules through `src-tauri/src/application/mod.rs` and `src-tauri/src/infra/mod.rs`, plus the targeted validation entrypoint.
- Do not change master persona build behavior, UI wiring, or unrelated documentation artifacts in this task.

### owned_scope

- backend

### required_reading

- `docs/exec-plans/active/2026-04-04-p2-i04-job-persona-persistence.md` の `UI` / `Scenario` / `Logic` / `Implementation Plan`
- `tasks/phase-2/tasks/P2-I04.yaml`
- `tasks/phase-2/tasks/P2-C03.yaml`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/er.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `src-tauri/src/application/dto/persona_storage/mod.rs`
- `src-tauri/src/application/ports/persona_storage/mod.rs`
- `src-tauri/src/application/master_persona/mod.rs`
- `src-tauri/src/infra/master_persona_builder/mod.rs`
- `src-tauri/src/application/mod.rs`
- `src-tauri/src/infra/mod.rs`
- `src-tauri/src/infra/execution_cache.rs`
- `src-tauri/tests/support/execution_cache.rs`
- `src-tauri/tests/validation/persona-rebuild/mod.rs`
- `src-tauri/tests/persona_rebuild_validation.rs`
- `src-tauri/migrations/0001_execution_cache_base.sql`
- `src-tauri/migrations/0002_master_dictionary_foundation.sql`
- `src-tauri/migrations/` only if the current versioned migration set does not already create `JOB_PERSONA_ENTRY`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`
- `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
- `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- job-local persona state remains separate from master persona state
- job persona persistence can evolve independently from master persona storage

## Required Evidence

- Active plan updated with distill facts, design decisions, ordered implementation scope, and validation commands.
- Tests or fixtures that prove job persona persistence remains separated from master persona persistence.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/backend-job-persona-persistence-class-diagram.d2`
- `4humans/diagrams/structures/backend-job-persona-persistence-class-diagram.svg`
- `4humans/diagrams/processes/backend-job-persona-persistence-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-job-persona-persistence-sequence-diagram.svg`
- `4humans/diagrams/overview-manifest.json`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.svg`
- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/processes-overview-robustness.svg`

## Outcome

- Added `src-tauri/src/application/job_persona/mod.rs` and module wiring in `src-tauri/src/application/mod.rs`, so job persona save requests are validated at the application boundary and run through a single `save -> read` orchestration path without widening into master persona or dictionary modules.
- Added `src-tauri/src/infra/job_persona_repository/mod.rs` and module wiring in `src-tauri/src/infra/mod.rs`, so job persona persistence now uses repository-owned transaction handling, replacement semantics per `job_id`, deterministic read ordering, and a temporary bridge from transport identity to ER-shaped storage (`translation_job.job_name -> id`, `npc.form_id -> id`).
- Added `src-tauri/migrations/0003_job_persona_entry_foundation.sql`, so fresh execution-cache bootstrap now provisions the minimal bridge tables and `job_persona_entry` storage needed for the current repository path. The migration and repository include explicit TODO comments marking this as a temporary bridge rather than the final canonical `TRANSLATION_JOB` / `NPC` schema.
- Extended `src-tauri/tests/validation/persona-rebuild/mod.rs` and `src-tauri/tests/persona_rebuild_validation.rs`, so validation now proves application request rejection, repository replacement semantics, missing-job read failure, rollback on insert failure, master/job separation, and bootstrap-path viability through `TempExecutionCache::initialize_base_schema()`.
- Added `4humans` detail diagrams for the job persona persistence slice and linked them from `4humans/diagrams/overview-manifest.json`, `4humans/diagrams/structures/backend-structure-overview.d2`, and `4humans/diagrams/processes/processes-overview-robustness.d2`.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`, `CARGO_HOME=/tmp/aitranslationenginejp-cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `sonar-scanner`, and `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` with `openIssueCount: 0`.
- `python3 scripts/harness/run.py --suite all` remains blocked by an unrelated dirty change in `.codex/skills/directing-fixes/SKILL.md` that fails the design harness on the missing `tasks.md` pattern; this is outside the `P2-I04` owned scope.
- Single-pass review continued to prefer widening the task into canonical `TRANSLATION_JOB` / `NPC` schema work. Per user direction, the implementation is closed here with explicit TODO markers documenting that the bridge schema and identity lookups are temporary and intentionally narrower than the final ER target.
