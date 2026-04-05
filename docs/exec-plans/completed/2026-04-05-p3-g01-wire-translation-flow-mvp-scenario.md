- workflow: impl
- status: completed
- lane_owner: codex
- scope: Wire the Phase 3 Translation Flow MVP through integrated Tauri gateway wiring and representative regression.
- task_id: P3-G01
- task_catalog_ref: tasks/phase-3/phase.yaml
- parent_phase: phase-3

## Request Summary

- Implement `tasks/phase-3/tasks/P3-G01.yaml`.
- Prove the first end-to-end Translation Flow MVP through the agreed orchestration boundaries without redesigning upstream contracts.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-3/phase.yaml`
- `tasks/phase-3/tasks/P3-G01.yaml`

## Owned Scope

- `src-tauri/src/lib.rs`
- `src/gateway/tauri/`
- `src-tauri/tests/regression/`

## Out Of Scope

- new contract design
- provider-specific retry policy
- concrete provider-backed persona-generation scenarios

## Dependencies / Blockers

- Depends on `P3-I01`, `P3-I02`, `P3-I03`, `P3-I04`, `P3-I05`, `P3-V01`, and `P3-V02`.
- If integrated wiring requires contract changes across earlier Phase 3 boundaries, stop and reroute instead of widening scope implicitly.

## Parallel Safety Notes

- `P3-G01` is the only task in batch `P3-B3`, but it composes outputs from earlier Phase 3 contract, verification, and implementation tasks.
- Shared boundary risk is concentrated in `src-tauri/src/lib.rs`, `src/gateway/tauri/`, and regression fixtures under `src-tauri/tests/regression/`.

## UI

- N/A. この task では screen、route、frontend state を追加しない。検証対象は backend composition と regression に限定する。
- production で呼ばれる新しい frontend / Tauri endpoint はこの task では増やさない。integration proof は gateway local の orchestration helper と regression で固定する。

## Scenario

- 代表シナリオは既存 fixture の `dialogue_response.text` 1 件を使い、`<Alias=Player>` を含む source text、再利用語、NPC metadata、embedded element descriptor を同じ入力系列でつなぐ。
- MVP の統合入口は gateway module 内の 1 件単位 orchestration helper とし、`word translation -> NPC persona generation -> body translation` を順に合成して、最終的に 1 件の `TranslationPreviewItemDto` 相当の結果を返す。preview 一覧取得や frontend 画面更新、production endpoint 追加まではこの task に含めない。
- persona generation は concrete provider 実行ではなく deterministic な fake で代表経路を証明する。代表 regression では cache miss から persona が補完される経路を優先し、provider retry や複数件 batch は扱わない。
- 既存の `translation-flow-mvp` regression anchor は contract snapshot の役割を維持しつつ、必要なら owned-scope で command 直下または command が使う orchestration helper の回帰を 1 件追加して、translation text・reusable terms・job persona・embedded element policy が同時に揃うことを証明する。
- もし統合のために `TranslationPhaseHandoffDto` や embedded-element preservation contract 自体の shape 変更が必要なら、この task では広げずに reroute する。

## Logic

- `src-tauri/src/gateway/commands.rs` に gateway local の orchestration helper を追加する前提で固める。helper は既存の `TranslationUnitDto`、phase request DTO、`TranslationPhaseHandoffDto` を使って phase を合成し、上流 contract を置き換える新 family は作らない。
- orchestration の責務は gateway 側の helper または private usecase に閉じ込め、phase 順序は `RunWordTranslationPhaseUseCase`、`RunNpcPersonaGenerationPhaseUseCase`、`RunBodyTranslationPhaseUseCase` の固定順とする。instruction 生成は既存どおり body translation phase 内の `build_translation_instruction` に委ねる。
- embedded element policy は代表入力から 1 回だけ組み立て、`TranslationPhaseHandoffDto` から `TranslationPreviewItemDto` まで同じ `unit_key` と descriptor 順序を保って流す。embedded-element 単体の survival は既存 `embedded-elements` regression に任せ、統合 regression では handoff から preview item までの保持を確認する。
- regression が concrete provider や DB 前提にならないよう、command 直下の orchestration は collaborator injection 可能な形に寄せる。回帰では fake dictionary lookup、fake persona generator、fake body translator を差し込み、既存 phase contract を壊さずに統合順序だけを証明する。
- Tauri transport 上で serde 境界が必要な場合は、gateway local wrapper か既存 DTO への transport 向け derive/casing 追加だけで済ませる。field rename や意味変更を伴う contract redesign はしない。
- `src-tauri/src/lib.rs` は新 endpoint 登録を増やさない。MVP proof の完了条件は「代表入力 1 件が gateway orchestration helper から preview item まで到達する」ことであり、preview persistence repository や一覧 query の常設化まではこの task に含めない。

## Implementation Plan

### 1. `src-tauri/src/gateway/commands.rs` と `src-tauri/src/lib.rs`

- translation-flow MVP 用の gateway-local orchestration helper を追加し、既存 phase use case を `word translation -> NPC persona generation -> body translation` の固定順で合成する。
- integration proof 用の fake collaborator は helper の注入点に閉じ込め、production command surface へ seeded endpoint を増やさない。
- `src-tauri/src/lib.rs` は新 invoke handler を追加しない。`src/gateway/tauri/` の TypeScript wrapper もこの task では変更しない。

### 2. `src-tauri/tests/regression/translation-flow-mvp/`

- 既存 contract snapshot は維持しつつ、command 直下または command が使う orchestration helper を deterministic fake collaborator で通す回帰を 1 件追加する。
- 代表入力 1 件について、translation text、reusable terms、job persona、embedded element policy が preview item まで同時に揃うことを証明する。

### 3. `src-tauri/tests/regression/embedded-elements/`

- 既存 embedded-element preservation anchor は既定では変更しない。
- translation-flow regression から参照する保持条件の基準として読み込みに留め、descriptor survival の単体責務は引き続きこの既存回帰に委ねる。

## Acceptance Checks

- Representative translation scenarios run through the agreed integrated path.
- The MVP translation flow is proven through integrated wiring and representative regression using the agreed phase orchestration boundaries.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Integrated regression evidence for the representative translation flow.
- Validation command results, Sonar open-issue status, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- Diagram updates are pending implementation outcome. If execution flow changes materially, update the relevant `4humans/diagrams/processes/*.d2` and matching `.svg`. If structure changes materially, update the relevant `4humans/diagrams/structures/*.d2` and matching `.svg`.

## Outcome

- Added `run_translation_flow_mvp_orchestration` to `src-tauri/src/gateway/commands.rs` as a gateway-local orchestration helper that composes `RunWordTranslationPhaseUseCase`, `RunNpcPersonaGenerationPhaseUseCase`, and `RunBodyTranslationPhaseUseCase` in fixed order without widening the public Tauri command surface.
- Added a representative regression in `src-tauri/tests/regression/translation-flow-mvp/mod.rs` that reuses the existing fixture and injects deterministic fake collaborators to prove translated text, reusable terms, job persona, and embedded-element policy reach one `TranslationPreviewItemDto` together.
- Kept `src-tauri/src/lib.rs` unchanged for command registration after the reroute fix, so seeded proof wiring does not become production transport behavior.
- Updated `4humans/quality-score.md` and `4humans/tech-debt-tracker.md` to reflect the new translation-flow orchestration regression and the remaining transport/provider acceptance-check gaps.
- Added `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.d2` / `.svg` and `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.d2` / `.svg`, then updated `4humans/diagrams/overview-manifest.json`, `4humans/diagrams/processes/processes-overview-robustness.d2` / `.svg`, and `4humans/diagrams/structures/backend-structure-overview.d2` / `.svg`.
- Validation passed: `python3 scripts/harness/run.py --suite structure`, `cargo test --manifest-path src-tauri/Cargo.toml --test translation_flow_mvp_regression`, `cargo test --manifest-path src-tauri/Cargo.toml --test embedded_elements_regression`, `python3 scripts/harness/run.py --suite backend-lint`, `sonar-scanner`, Sonar MCP `OPEN` issue query for project `ishibata91_AITranslationEngineJP` returned `0`, and `python3 scripts/harness/run.py --suite all`.
- Single-pass review returned `reroute`; the seeded public command path was removed and, per lane contract, no second review was run after the reroute fix.
