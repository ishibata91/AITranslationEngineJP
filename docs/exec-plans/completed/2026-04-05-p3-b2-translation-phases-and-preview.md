- workflow: impl
- status: completed
- lane_owner: codex
- scope: Complete Phase 3 batch `P3-B2` by landing the remaining translation phases and preview path on top of the already completed instruction builder.
- task_id: P3-B2
- task_catalog_ref: tasks/phase-3/phase.yaml
- parent_phase: phase-3

## Request Summary

- Implement batch `P3-B2` from `tasks/phase-3/phase.yaml`.
- Treat `P3-I01` as already completed and use the existing instruction-builder boundary as an upstream dependency rather than reopening it.
- Land `P3-I02`, `P3-I03`, `P3-I04`, and `P3-I05` so Translation Flow MVP can translate reusable terms, orchestrate job-local persona generation, run body translation with embedded-element preservation, and expose the first preview path.

## Decision Basis

- `.codex/README.md`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- `docs/coding-guidelines.md`
- `tasks/phase-3/phase.yaml`
- `tasks/phase-3/tasks/P3-I01.yaml`
- `tasks/phase-3/tasks/P3-I02.yaml`
- `tasks/phase-3/tasks/P3-I03.yaml`
- `tasks/phase-3/tasks/P3-I04.yaml`
- `tasks/phase-3/tasks/P3-I05.yaml`
- `docs/exec-plans/completed/2026-04-05-p3-b1-phase-contracts-and-regression-anchors.md`
- `docs/exec-plans/completed/2026-04-05-p3-i01-record-type-instruction-builder.md`

## Owned Scope

- `src-tauri/src/application/word_translation_phase/`
- `src-tauri/src/application/npc_persona_generation_phase/`
- `src-tauri/src/application/body_translation_phase/`
- `src-tauri/src/application/translation_preview/`
- `src/ui/screens/translation-preview/`
- `src/ui/views/translation-preview/`
- minimal shared wiring, tests, fixtures, and exports required to compose the batch without reopening upstream Phase 2 or `P3-B1` contracts

## Out Of Scope

- provider adapter implementation details
- output writer logic and file export behavior
- provider selection UI
- master persona rebuild behavior
- Phase 3 integrated proof work reserved for `P3-B3`

## Dependencies / Blockers

- `P3-I01` is already completed and should be consumed as-is unless a verified contract gap forces a separate follow-up.
- `P3-I02` depends on `P3-C03` and `P2-C02`.
- `P3-I03` depends on `P3-C03` and `P2-C03`.
- `P3-I04` depends on `P3-C01`, `P3-C02`, `P3-C03`, `P3-V01`, and `P3-V02`.
- `P3-I05` depends on `P3-I04`.

## Parallel Safety Notes

- The batch remains parallel-safe only while each task stays inside its declared owned scope and shared wiring stays minimal.
- `P3-I05` cannot finalize before `P3-I04` exposes a stable preview-facing output shape.
- Shared fixtures or exports introduced for body translation and preview must not absorb provider runtime or writer policy that belongs to later phases.

## UI

- `P3-I05` は `src/ui/screens/translation-preview/` と `src/ui/views/translation-preview/` に 1 画面を追加し、既存 `Screen -> View` と `FeatureScreenState` 系の observe 画面パターンを踏襲する。preview 用 route は増やさず、`AppShell.svelte` に optional な store / usecase props を足して single-window composition のまま差し込む。
- preview 画面は read-only observation に限定し、provider 選択、再翻訳、writer/export 操作は持たない。最小 UI は `jobId` を指定して preview を取得する入力、preview item 一覧、選択 item の詳細で構成する。
- 一覧は `unitKey` を基準に body translation 済み item を並べ、詳細は `translationUnit.sourceText`、translated text、reusable terms、job-local persona、preserved embedded elements を表示する。preview は Phase 3 の観測入口であり、writer 向け format や provider diagnostics は露出しない。
- `AppShell` 側の共有状態は増やしすぎず、初期実装では preview 画面を self-contained に保つ。job list 選択との連動は後続 task に残し、この batch では preview 画面自身の filter / selection state だけで成立させる。

## Scenario

- 代表シナリオは `P3-V02` の `<Alias=Player>` を含む `dialogue_response.text` をそのまま使う。word translation phase が reusable term 群を確定し、NPC persona generation phase が同一 NPC の job-local persona を返し、body translation phase が `P3-I01` の instruction builder 出力と upstream handoff をまとめて消費して translated text を生成する。
- body translation の結果は preview query がそのまま読める stable shape で保持し、preview UI は同じ representative item を一覧 1 件 + 詳細 1 件として観測できるようにする。preview で確認するのは translated text と preservation / reuse / persona の反映結果であり、output file 書き出し結果ではない。
- 初期 acceptance は representative 1 件で十分だが、phase orchestration は `job_persona` が absent、`reusable_terms` が empty のケースでも shape を変えない。UI は optional section を空状態で描画し、body translation や preview contract を分岐で増やさない。
- refresh / retry は preview query の再読込だけを行い、word translation・persona generation・body translation の再実行起点にはしない。Phase 実行と preview 観測の責務を分離したまま `P3-B3` の統合 wiring に引き渡す。

## Logic

- `src-tauri/src/application/word_translation_phase/` は application 層の phase entrypoint として、body translation 対象 unit に対して再利用する `ReusableDictionaryEntryDto` 群を確定する責務だけを持つ。dictionary candidate の取得は既存 `DictionaryLookupPort` 越しに行い、選定結果は `TranslationPhaseHandoffDto.reusable_terms` にそのまま載る shape に閉じる。body translation instruction、provider adapter detail、preview shaping は持ち込まない。
- `src-tauri/src/application/npc_persona_generation_phase/` は job-local persona の read / generate / save orchestration を持ち、public handoff は `Option<JobPersonaEntryDto>` に固定する。保存は既存 `JobPersonaStoragePort` を使い、master persona rebuild や master persona persistence には触れない。concrete provider runtime は phase module 内の抽象境界越しに呼び、保存先分離だけをこの task で守る。
- `src-tauri/src/application/body_translation_phase/` は `P3-I01` の `build_translation_instruction` を upstream dependency として再利用し、`TranslationUnitDto` と phase handoff から本文翻訳を実行する。module 境界では `instruction.unit_key == handoff.translation_unit.extraction_key` の整合、embedded element preservation contract の適用、reusable terms / optional persona の注入を扱い、provider retry policy や writer format 生成は扱わない。
- `P3-I04` で固定する preview-facing output shape は、少なくとも `job_id`、`unit_key`、`translation_unit`、`translated_text`、`reusable_terms`、`job_persona`、`embedded_element_policy` を持つ body-translation result / preview item 単位の application DTO とする。`P3-I05` はこの shape を read するだけにし、preview query 側で instruction rebuild や provider 再実行をしない。
- `src-tauri/src/application/translation_preview/` は preview item 群を `job_id` 単位で返す read-only query root とし、初期実装では body translation result source を抽象ポートまたは最小 shared wiring で受ける。preview module 自体は writer/export boundary を知らず、返却順は body translation result の `sort_key` / `unit_key` に従って安定化させる。

## Implementation Plan

- Ordered scope 1 (`src-tauri/src/application/word_translation_phase/`): implement the word-translation phase entrypoint against `DictionaryLookupPort` so one `TranslationUnitDto` can resolve reusable terms and hand back a `TranslationPhaseHandoffDto.reusable_terms`-compatible result without absorbing instruction building, body translation, provider detail, or preview shaping. Cover both selected-term and empty-term paths so downstream shape stays stable.
- Ordered scope 2 (`src-tauri/src/application/npc_persona_generation_phase/`): implement job-local persona read / generate / save orchestration behind the existing storage and provider-neutral boundaries, and fix the public output to `Option<JobPersonaEntryDto>` only. Cover cache-hit, generation-save, and absent-persona paths without reopening master persona rebuild or writer behavior.
- Ordered scope 3 (`src-tauri/src/application/body_translation_phase/`): implement body-translation orchestration on top of `build_translation_instruction(...)` and `TranslationPhaseHandoffDto`, and fix one stable preview-facing backend result shape with `job_id`, `unit_key`, `translation_unit`, `translated_text`, `reusable_terms`, `job_persona`, and `embedded_element_policy`. Keep embedded-element preservation application, unit-key alignment, and representative regression fixture reuse inside this module; exclude retry policy, provider adapter detail, and output formatting.
- Ordered scope 4 (`src-tauri/src/application/translation_preview/`): implement the read-only preview query root that accepts `job_id`, reads the Ordered scope 3 result shape without rebuilding instructions or rerunning translation, and returns preview items in stable `translation_unit.sort_key` then `unit_key` order. Treat this module as the only preview observation boundary that `P3-I05` reads.
- Ordered scope 5 (`src-tauri/src/application/mod.rs`, `src-tauri/src/application/dto/mod.rs`, `src-tauri/tests/`): add only the exports, minimal shared wiring, and regression / contract coverage needed to compose Ordered scope 1 through Ordered scope 4 on top of the existing `P3-I01` boundary and `P3-B1` anchors. Do not widen DTO roots or reopen provider adapters, output writers, or master persona ownership.
- Ordered scope 6 (`src/ui/screens/translation-preview/`, `src/ui/views/translation-preview/`, minimal preview wiring): after Ordered scope 4 freezes the preview-facing output shape, add the self-contained translation-preview usecase / screen / view path and optional `AppShell.svelte` composition hook. UI owns `jobId` input, preview list selection, selected-item detail, refresh / retry, and empty optional sections for `reusable_terms` / `job_persona`; UI must not own phase execution, provider choice, or writer/export actions.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_flow_mvp_regression -- --nocapture`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
  - `python3 scripts/harness/run.py --suite frontend-lint`
  - `npm run test -- src/application/usecases/translation-preview src/ui/screens/translation-preview`
  - `npm run build`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- `translation_preview` backend query は `job_id` 指定で preview item を返し、item shape は `job_id`、`unit_key`、`translation_unit`、`translated_text`、`reusable_terms`、`job_persona`、`embedded_element_policy` を保ったまま `translation_unit.sort_key`、次に `unit_key` で安定順序化される。
- preview observation 用 frontend usecase は `jobId` 入力を request に反映し、初回 observe で先頭 item を選択し、refresh/retry では同じ `unitKey` が残る限り selection を維持する。
- preview view は representative item について `translationUnit.sourceText`、`translatedText`、`reusableTerms`、`jobPersona`、preserved embedded elements を read-only で表示し、`reusableTerms=[]` と `jobPersona=null` でも空状態描画に落ちて shape を分岐させない。
- `AppShell.svelte` は preview store / usecase が渡された時だけ preview screen を加算的に compose し、既存 screen 順序と single-window 構成を崩さない。

## Required Evidence

- Active plan updated with distill facts, task-local design, ordered work brief, and acceptance checks.
- Added or updated implementation code, tests, and fixtures that cover the Phase 3 translation phases and preview path.
- Validation command results, Sonar open-issue status for touched paths, and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- Review whether Phase 3 execution-flow changes require updates under `4humans/diagrams/processes/`.
- Review whether new composition or boundary changes require updates under `4humans/diagrams/structures/`.
- No new detail diagram is assumed upfront; if one is added, update `4humans/diagrams/overview-manifest.json` and the linked overview `.d2` / `.svg` in the same change.

## Outcome

- Implemented the remaining backend Phase 3 application modules for word-translation handoff, job-local NPC persona generation, body translation orchestration, and read-only translation preview query composition on top of the existing `P3-I01` instruction-builder boundary.
- Added the frontend translation-preview usecase, screen, and view wiring, plus additive `AppShell.svelte` composition so preview can be observed without introducing provider selection or export actions.
- Added direct backend coverage for the preview contract and for `P3-I03` / `P3-I04` phase behavior, and updated existing AppShell SSR tests to stub the new preview screen import path.
- Synced the affected `4humans` process / structure overview diagrams and regenerated their SVG artifacts after the new preview and phase boundaries landed.
- Validation passed: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --test translation_flow_mvp_regression -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `python3 scripts/harness/run.py --suite frontend-lint`, `npm run test -- src/application/usecases/translation-preview src/ui/screens/translation-preview`, `npm run build`, and `python3 scripts/harness/run.py --suite all`.
- Sonar verification for the owned scope returned `openIssueCount: 0`.
