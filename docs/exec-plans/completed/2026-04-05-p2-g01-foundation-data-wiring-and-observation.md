- workflow: impl
- status: completed
- lane_owner: codex
- scope: Implement task `P2-G01` from `tasks/phase-2/tasks/P2-G01.yaml` by wiring dictionary and persona foundation rebuild plus observation through the agreed Tauri validation paths without redesigning storage boundaries or fixtures.
- task_id: P2-G01
- task_catalog_ref: tasks/phase-2/phase.yaml
- parent_phase: phase-2

## Request Summary

- Implement `P2-G01`.
- Prove dictionary and persona foundation data rebuild and observation through integrated wiring and agreed validation paths.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `docs/screen-design/wireframes/foundation-data.md`
- `tasks/phase-2/phase.yaml`
- `tasks/phase-2/tasks/P2-G01.yaml`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`
- `docs/exec-plans/completed/2026-04-04-p2-i01-xtranslator-importer.md`
- `docs/exec-plans/completed/2026-04-04-p2-i02-master-dictionary-storage-query.md`
- `docs/exec-plans/completed/2026-04-04-p2-i03-master-persona-builder.md`
- `docs/exec-plans/completed/2026-04-04-p2-i04-job-persona-persistence.md`
- `docs/exec-plans/completed/2026-04-05-p2-i05-dictionary-observability-ui.md`
- `docs/exec-plans/completed/2026-04-05-p2-i06-persona-observability-ui.md`

## Owned Scope

- `src-tauri/src/lib.rs`
- `src-tauri/src/infra/master_persona_repository/` as safe adjacency when persona command wiring would otherwise reimplement storage in the transport layer
- `src/gateway/tauri/`
- `src-tauri/tests/validation/`

## Out Of Scope

- new contract design
- new fixture design
- widening frontend observation behavior beyond existing dedicated screens
- storage-boundary redesign for dictionary or persona foundation data

## Dependencies / Blockers

- `P2-G01` depends on `P2-I01`, `P2-I02`, `P2-I03`, `P2-I04`, `P2-I05`, `P2-I06`, `P2-V01`, and `P2-V02`.
- Integrated proof must reuse the already fixed Phase 2 contracts, validation anchors, and observation screens instead of inventing parallel composition paths.

## Parallel Safety Notes

- This task stops being parallel-safe if it reaches into upstream foundation contracts, fixture design, or frontend screen design beyond minimal transport wiring.
- Shared files are expected around `src-tauri/src/lib.rs`, Tauri gateway adapters, and validation tests that compose dictionary and persona paths together.

## UI

- 既存の `dictionary-observe` と `persona-observe` を transport-backed にする。`foundation-data` wireframe の新規 pane、editor、rebuild button、metadata 項目追加は行わず、`P2-I05` / `P2-I06` で固定した read-only observation panel をそのまま使う。
- frontend 側で追加するのは `src/gateway/tauri/` の thin adapter と `src/main.ts` での配線だけに留める。dictionary panel は `DictionaryLookupRequest` / `DictionaryLookupResult`、persona panel は `MasterPersonaReadRequestDto` / `MasterPersonaReadResultDto` をそのまま受ける executor を注入し、client 側で request / response shape を作り替えない。
- rebuild はこの task で UI 操作として露出しない。dictionary と persona の rebuild 実行は Tauri validation path から起動し、観測 UI は rebuild 後の persisted foundation data を読む経路だけを担う。

## Scenario

- dictionary integration proof は `xtranslator-shared-reusable-entry.sst` を使った canonical rebuild fixture を Tauri command で投入し、その直後に既存 dictionary observation と同じ `sourceTexts` request で観測する 2 段構成とする。受け入れ対象は `rebuild -> persist -> observe` の一貫性であり、fixture や lookup semantics 自体は `P2-V01` / `P2-I02` から変えない。
- persona integration proof は `base-game-master-persona-rebuild.fixture.json` を使った canonical rebuild fixture を Tauri command で投入し、その直後に既存 persona observation と同じ `personaName` request で観測する。返却される `personaName`、`sourceType`、entry 順は `P2-I03` / `P2-V02` の rebuild 結果と一致させ、job-local persona は同じ観測 path に混在させない。
- runtime composition では `main.ts` が dictionary と persona の observe screen に Tauri executor を渡し、screen mount 後は各 usecase の既存ルールどおり explicit `observe` / `refresh` / `retry` でのみ backend を呼ぶ。自動 preload、frontend 側 rebuild orchestration、screen redesign はこの task に含めない。

## Logic

- `src-tauri/src/lib.rs` は foundation-data 用の additive Tauri command を登録する。最低限、dictionary rebuild、dictionary observe、persona rebuild、persona observe の 4 command を持ち、各 command は既存 application usecase をそのまま束ねる transport boundary として実装する。`#[tauri::command]` 内に SQL、fixture 解釈、UI 向け整形は書かない。
- dictionary rebuild command は `DictionaryImportRequestDto` を受けて importer と dictionary repository を直列実行し、dictionary observe command は `DictionaryLookupRequest` を受けて既存 lookup port の結果を返す。persona rebuild command は `BaseGameNpcRebuildRequest` を受けて `RebuildMasterPersonaUseCase` を起動し、persona observe command は `MasterPersonaReadRequestDto` を受けて master persona storage の read result を返す。command 名と引数名は camelCase transport を崩さず、frontend adapter から追加変換なしで呼べる形に固定する。
- `src/gateway/tauri/` は既存 `job-create` / `job-list` と同じ thin invoke pattern を踏襲し、dictionary observe 用 executor と persona observe 用 executor を public root として追加する。必要な隣接変更は `src/main.ts` の dependency composition までに留め、UI usecase / screen 側の contract は変更しない。
- `src-tauri/tests/validation/` の proof は既存 `dictionary-rebuild` / `persona-rebuild` anchor を土台に、Tauri-facing command path を通した integration check を追加する。assertion は `rebuild command が canonical fixture を保存できること`、`observe command が既存 observation contract を返すこと`、`dictionary/persona ともに rebuild 後の observation が upstream anchor と同じ deterministic data を返すこと` に絞り、新しい fixture schema や並行する snapshot vocabulary は増やさない。

## Implementation Plan

### ordered_scope

1. Tauri command wiring (`src-tauri/src/lib.rs` と既存 transport 置き場の `src-tauri/src/gateway/commands.rs`)

- `generate_handler!` に foundation-data 用の 4 command を追加し、既存 command と同じ single `request` payload pattern に固定する。
- command 名は既存 snake_case 慣例と upstream usecase 名に合わせて `rebuild_dictionary`、`lookup_dictionary`、`rebuild_master_persona`、`read_master_persona` とする。
- `rebuild_dictionary` は `DictionaryImportRequestDto` を受けて import と persist を直列実行し、成功時は既存 import contract を返す。`lookup_dictionary` は `DictionaryLookupRequest` を受けて `DictionaryLookupResult` を返す。`rebuild_master_persona` は `BaseGameNpcRebuildRequest` を transport からそのまま受け、`MasterPersonaReadResultDto` を返す。`read_master_persona` は `MasterPersonaReadRequestDto` を受けて同じ read contract を返す。
- SQL、fixture 解釈、snapshot 整形は command 内へ持ち込まず、既存 application / infra module を束ねる transport boundary に留める。persona 側 storage 実装が未整備なら `src-tauri/src/infra/master_persona_repository/` を safe adjacency として追加し、transport module へ repository 実装を置かない。

2. Tauri frontend adapters (`src/gateway/tauri/`、必要なら `src/main.ts` のみ)

- `job-create` と同じ thin invoke pattern で `dictionary-observe/` と `persona-observe/` の executor root を追加し、それぞれ `lookup_dictionary` と `read_master_persona` を `{ request }` payload で呼ぶ。
- adapter test は既存 gateway test と同じ粒度に揃え、command 名と named payload が崩れていないことだけを固定する。
- runtime composition は既存 observe screen の executor 注入に閉じる。screen/usecase contract、empty or loading or retry の UI state、rebuild UI action は変更しない。

3. Validation proof (`src-tauri/tests/validation/`)

- 既存 `dictionary-rebuild` / `persona-rebuild` anchor を土台に、Tauri command 関数を直接通す integration check を 1 系統ずつ追加する。
- dictionary 側は canonical SST fixture で `rebuild_dictionary` を呼んだ後、同じ request shape で `lookup_dictionary` を呼び、既存 snapshot の `dictionaryImportResult` と `dictionaryLookupResult` に一致することを確認する。
- persona 側は canonical rebuild fixture で `rebuild_master_persona` を呼んだ後、同じ `personaName` で `read_master_persona` を呼び、返却される `personaName`、`sourceType`、entry 順が既存 snapshot と一致することを確認する。
- assertion は command path の wiring に限定し、fixture schema、snapshot vocabulary、job-local persona path は増やさない。

### owned_scope

- `src-tauri/src/lib.rs`
- `src-tauri/src/gateway/commands.rs`
- `src-tauri/src/infra/master_persona_repository/` as safe adjacency if persona command wiring needs a dedicated infra repository instead of transport-local SQL
- `src/gateway/tauri/`
- `src/main.ts` only if executor injection is still missing after adapter追加
- `src-tauri/tests/validation/`

### required_reading

- `docs/exec-plans/active/2026-04-05-p2-g01-foundation-data-wiring-and-observation.md` の `UI` / `Scenario` / `Logic`
- `docs/exec-plans/completed/2026-04-03-p2-b1-foundation-contracts-and-rebuild-anchors.md`
- `docs/exec-plans/completed/2026-04-04-p2-i02-master-dictionary-storage-query.md`
- `docs/exec-plans/completed/2026-04-04-p2-i03-master-persona-builder.md`
- `docs/exec-plans/completed/2026-04-05-p2-i05-dictionary-observability-ui.md`
- `docs/exec-plans/completed/2026-04-05-p2-i06-persona-observability-ui.md`
- `src-tauri/src/lib.rs`
- `src-tauri/src/gateway/commands.rs`
- `src-tauri/src/application/dictionary_query/mod.rs`
- `src-tauri/src/application/master_persona/mod.rs`
- `src/gateway/tauri/job-create/index.ts`
- `src/gateway/tauri/job-create/index.test.ts`
- `src/application/usecases/dictionary-observe/index.ts`
- `src/application/usecases/persona-observe/index.ts`
- `src/main.ts`
- `src-tauri/tests/validation/dictionary-rebuild/mod.rs`
- `src-tauri/tests/validation/persona-rebuild/mod.rs`

### validation_commands

- `python3 scripts/harness/run.py --suite structure`
- `npm run test -- src/gateway/tauri`
- `npm run build`
- `CARGO_HOME=.cargo-home cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
- `CARGO_HOME=.cargo-home cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
- `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`
- `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`
- `python3 scripts/harness/run.py --suite all`

### implementation_updates

- backend-first split でも command 本体は `lib.rs` 単独では置けないため、既存 transport module の `src-tauri/src/gateway/commands.rs` を safe owned adjacency として含める。
- `src/main.ts` は strict には task catalog scope 外だが、transport-backed observe screen の executor 注入が未配線なら 1 file adjacency としてのみ開く。未配線であっても screen/usecase 側 contract 変更へは広げない。

## Acceptance Checks

- dictionary and persona rebuild plus observation run through agreed validation paths

## Required Evidence

- Active plan updated with task-local UI, scenario, and logic decisions for integrated Tauri wiring.
- Tests or validation coverage proving rebuild and observation behavior through owned gateway and validation paths.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md` update not required unless repository-level quality posture changes
- `4humans/tech-debt-tracker.md` update not required unless new unresolved debt is introduced
- `4humans/diagrams/structures/backend-dictionary-storage-query-class-diagram.d2`
- `4humans/diagrams/structures/backend-dictionary-storage-query-class-diagram.svg`
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.d2`
- `4humans/diagrams/structures/backend-master-persona-builder-class-diagram.svg`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.svg`
- `4humans/diagrams/processes/backend-dictionary-storage-query-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-dictionary-storage-query-sequence-diagram.svg`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-master-persona-builder-sequence-diagram.svg`
- `4humans/diagrams/processes/processes-overview-robustness.d2`
- `4humans/diagrams/processes/processes-overview-robustness.svg`

## Outcome

- Added foundation-data Tauri command wiring in `src-tauri/src/lib.rs` and `src-tauri/src/gateway/commands.rs` so dictionary rebuild or observe and persona rebuild or observe now run through additive command boundaries named `rebuild_dictionary`, `lookup_dictionary`, `rebuild_master_persona`, and `read_master_persona`.
- Added thin frontend Tauri adapters in `src/gateway/tauri/dictionary-observe/` and `src/gateway/tauri/persona-observe/`, then wired them through `src/main.ts` so the existing observation screens consume transport-backed executors without changing screen-local contracts.
- Added `src-tauri/src/infra/master_persona_repository/mod.rs` and updated `src-tauri/src/infra/mod.rs` so master persona persistence stays in infra rather than being reimplemented inside the Tauri transport layer.
- Extended `src-tauri/tests/validation/dictionary-rebuild/mod.rs` and `src-tauri/tests/validation/persona-rebuild/mod.rs` to prove `rebuild -> persist -> observe` through Tauri command paths while reusing the existing canonical fixtures and snapshots.
- Synced `4humans` backend diagrams in the listed structure and process files, validated the updated `.d2` sources, and regenerated the paired `.svg` outputs with `d2 -t 201`.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `npm run test -- src/gateway/tauri`, `npm run build`, `CARGO_HOME=.cargo-home cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `CARGO_HOME=.cargo-home cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test dictionary_rebuild_validation -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test persona_rebuild_validation -- --nocapture`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` with `openIssueCount = 0`, and `python3 scripts/harness/run.py --suite all`.
- Single-pass review returned `reroute`; the reroute was resolved by removing transport-local persona SQL from `src-tauri/src/gateway/commands.rs`, restoring `BaseGameNpcRebuildRequest` as the rebuild command contract, moving persistence into infra, and adding the required `4humans` diagram sync. Per lane contract, no second review was run.
