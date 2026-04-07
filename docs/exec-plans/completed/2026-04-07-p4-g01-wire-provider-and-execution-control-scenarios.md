- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Phase 4 integ task `P4-G01` として provider selection、execution mode choice、execution control を既存 Phase 4 実装へ通した統合シナリオ証跡を固定する。
- task_id: P4-G01
- task_catalog_ref: tasks/phase-4/tasks/P4-G01.yaml
- parent_phase: phase-4

## 要求要約

- `tasks/phase-4/tasks/P4-G01.yaml` を implementation proposal として具体化する。
- provider failure、retry、pause、recovery を含む execution-control シナリオを、合意済みの integrated path 上で証明する方針を固める。
- provider contract や state policy を再設計せず、既存 Phase 4 実装を束ねる task-local design を human review 可能な形にする。

## 判断根拠

- `tasks/phase-4/tasks/P4-G01.yaml` は owned scope を `src-tauri/src/lib.rs`、`src/gateway/tauri/`、`src-tauri/tests/acceptance/provider-failure-retry/` に固定している。
- 依存 task として `P4-I01` から `P4-I06`、`P4-V01` が列挙されており、個別実装済みの provider adapter、execution mode、execution control UI/observe を integrated path へ接続する task と読める。
- acceptance anchor は provider failure、retry、pause、recovery を agreed integrated path で通すことを要求しており、task-local design では UI / gateway / backend command registration / acceptance fixture の責務分離を維持する必要がある。
- `P4-I06-2` により execution-observe の read path は `src/gateway/tauri/execution-observe/`、`src-tauri/src/gateway/commands.rs`、`src-tauri/src/lib.rs` を通る integrated route として既に成立している。一方で provider selection と execution-control command の integrated proof はまだ `P4-G01` で固定する必要がある。
- `P4-V01` の acceptance fixture は `Running -> RecoverableFailed -> Retrying -> Running -> Completed` と `Running -> Paused -> Running -> Canceled` を既に外部固定しており、`P4-G01` ではこの fixture を provider-neutral な command / observe path の証跡へ昇格させるのが最短である。

## 対象範囲

- `src-tauri/src/lib.rs`
- `src/gateway/tauri/`
- `src-tauri/tests/acceptance/provider-failure-retry/`

## 対象外

- new provider contract design
- new state policy
- 永続仕様や `docs/` 正本の更新
- Phase 4 既存 task の contract 再オープン

## 依存関係・ブロッカー

- `P4-I01`
- `P4-I02`
- `P4-I03`
- `P4-I04`
- `P4-I05`
- `P4-I06`
- `P4-V01`
- 既存 Phase 4 の provider selection、execution-control、execution-observe 契約が stable であること

## 並行安全メモ

- `src/gateway/tauri/` は共有 entrypoint のため、既存 gateway slice を再利用し、adapter 固有 detail を UI へ漏らさない。
- `src-tauri/src/lib.rs` は command registration の共有 composition root なので、追加 wiring は additive 最小に留める。
- acceptance fixture 追加や更新は `provider-failure-retry` anchor に閉じ、別 acceptance suite の責務を取り込まない。

## UI

- provider selection、execution mode choice、execution-control の表示構成は既存 screen / usecase を前提とし、`P4-G01` では新しい screen や provider-specific panel を増やさない。必要なのは既存 UI が integrated Tauri path を使って同一 vocabulary を読めることの証明だけとする。
- `src/gateway/tauri/` では execution-observe と同じ責務分離で provider selection と execution-control command 用の thin gateway slice を追加または補完し、`src/main.ts` や screen 層には injected function だけを渡す。Tauri command 名、Rust DTO 名、adapter 固有 detail は `src/gateway/tauri/` を越えて露出しない。
- provider name、execution mode、control state、failure / retry 情報は既存の provider-neutral / UI-stable vocabulary を継続利用し、pause、retry、recovery 後の見え方も既存 execution-control / execution-observe UI の state transition に従う。`P4-G01` 専用の UI 分岐や補助 vocabulary は導入しない。

## Scenario

- user が provider と execution mode を選択して execution を開始した後、integrated path は frontend gateway、Tauri command registration、backend execution runtime、`provider-failure-retry` acceptance fixture を通じて 1 本で進行する。proof の主対象は「selection / mode / control が同じ route へ乗ること」であり、contract 自体の再定義ではない。
- provider failure が発生した時、既存 execution-observe path は fixture-backed snapshot を provider-neutral に返し、既存 execution-control path は pause / resume / retry / cancel の command を同じ Tauri wiring 経路で backend へ渡す。failure state、retry affordance、recovery 後の resumed path は同じ integrated route 上で再読込される。
- acceptance anchor は `Running -> RecoverableFailed -> Retrying -> Running -> Completed` と `Running -> Paused -> Running -> Canceled` を外部固定済みの正本として再利用し、`P4-G01` では command path の追加後にこの anchor を integrated scenario proof として成立させる。新しい state machine、adapter-specific bypass、fixture 専用 vocabulary は導入しない。

## Logic

- `src/gateway/tauri/` は provider selection、execution-control、execution-observe の integrated adapter surface を provider-neutral に揃える最小共有点だけを持つ。`P4-G01` で追加するのは thin Tauri gateway と mapping に限り、selection logic、state policy、provider validation は既存 usecase / backend contract の owner に残す。
- `src-tauri/src/lib.rs` は `P4-I06-2` で成立した observe registration を壊さず、selection / execution-control command の Tauri registration を additive 最小で統合する。もし shared touch point が `src-tauri/src/gateway/commands.rs` まで必要なら、その 1 箇所だけを正当化上限とする。
- `src-tauri/tests/acceptance/provider-failure-retry/` は integrated proof の正本 anchor として、provider failure、retry、pause、recovery、完了までの主要経路を既存 contract 上で検証する。fixture は provider-neutral な state / failure vocabulary を維持し、adapter internals や新しい state machine を持ち込まない。
- `P4-G01` の proof は observe-only の既存 route と command path を結び、UI action から backend state transition、再観測までが 1 つの acceptance narrative に入ることを示す。contract gap が見つかった場合は redesign へ進まず gap として記録する。

## 実装計画

- ordered_scope:
  1. `P4-I04`、`P4-I05`、`P4-I06-2`、`P4-V01` の既存成果を前提に、provider selection / execution mode / execution-control command / execution-observe が 1 本の integrated route に入る最小 touch point を固定する。
  2. `src/gateway/tauri/` に provider-neutral な thin gateway slice を追加し、execution-control command と provider-selection proof を既存 observe と同じ Tauri integration style に揃える。frontend 側は injected function の差し替えだけに留め、selection logic や provider-specific rule は持ち込まない。
  3. `src/main.ts`、必要なら `src/App.svelte` と `src/ui/app-shell/AppShell.svelte` の shared composition touch point を additive 最小で更新し、既存 execution-control screen / usecase と execution-observe screen / usecase が同じ integrated gateway 群を受け取れるようにする。
  4. `src-tauri/src/lib.rs` と必要最小限の backend gateway command entrypoint に command registration を追加し、pause / resume / retry / cancel と observe 再読込が同一 route で backend fixture state を読み書きできるようにする。
  5. `src-tauri/tests/acceptance/provider-failure-retry/` を integrated proof の正本として更新し、provider failure、retry、pause、recovery、完了までの主要経路と provider selection / execution mode 証跡を UI-stable vocabulary のまま検証する。
- owned_scope:
  - `src-tauri/src/lib.rs`
  - `src/gateway/tauri/`
  - `src-tauri/tests/acceptance/provider-failure-retry/`
  - shared touch point が strictly required な場合のみ `src/main.ts`、`src/App.svelte`、`src/ui/app-shell/AppShell.svelte`
- required_reading:
  - `tasks/phase-4/tasks/P4-G01.yaml`
  - `tasks/phase-4/tasks/P4-V01.yaml`
  - `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md`
  - `docs/exec-plans/completed/2026-04-06-p4-i04-batch-and-single-shot-switching.md`
  - `docs/exec-plans/completed/2026-04-06-p4-i05-pause-retry-cancel-ui.md`
  - `docs/exec-plans/completed/2026-04-07-p4-i06-2-wire-execution-observe-snapshot-path.md`
  - `docs/exec-plans/completed/2026-04-07-p4-i07-provider-backed-master-persona-generation.md`
  - `src-tauri/src/gateway/commands.rs`
  - `src-tauri/src/lib.rs`
  - `src/main.ts`
  - `src/App.svelte`
  - `src/ui/app-shell/AppShell.svelte`
  - `src-tauri/tests/acceptance/provider-failure-retry/mod.rs`
  - `src-tauri/tests/acceptance/provider-failure-retry/fixtures/provider-failure-retry.fixture.json`
- validation_commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite frontend-lint`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `npm test -- src/gateway/tauri/execution-observe/index.test.ts src/gateway/tauri/execution-control/index.test.ts src/ui/screens/execution-observe/index.test.ts src/ui/screens/execution-control/index.test.ts`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml provider_failure_retry -- --nocapture`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/gateway/tauri src-tauri/src/lib.rs src-tauri/src/gateway/commands.rs src-tauri/tests/acceptance/provider-failure-retry src/main.ts src/App.svelte src/ui/app-shell/AppShell.svelte`
  - `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`

## 受け入れ確認

- provider selection、execution mode choice、execution control が agreed integrated path 上で接続される proposal になっている
- provider failure、retry、pause、recovery を acceptance anchor で証明する方針が task-local design に落ちている
- provider contract と state policy を reopen しない境界が明記されている

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`
- distilling-implementation からの facts / constraints / gaps / required_reading
- designing-implementation による `UI` / `Scenario` / `Logic` 固定結果
- 必要時 review 用差分図

## HITL 状態

- human review 完了
- LGTM 済み

## 承認記録

- 2026-04-07: human LGTM 受領

## review 用差分図

- N/A
- 現時点では integrated path の shared touch point が `src/gateway/tauri/`、`src-tauri/src/lib.rs`、`src-tauri/tests/acceptance/provider-failure-retry/` に閉じており、active plan 本文だけで review 可能と判断する。implementation 具体化で共有 wiring が増える場合だけ差分図を追加する。

## 4humans Sync

- `4humans/quality-score.md`
  今回は更新しない。品質スコアの評価軸自体は変更していない。
- `4humans/tech-debt-tracker.md`
  今回は更新しない。follow-up を新規 debt として固定する事項は残らなかった。
- `4humans/diagrams/structures/*.d2` と対応する `.svg`
  `4humans/diagrams/structures/frontend-execution-control-slice-class-diagram.d2` と対応する `.svg` を更新した。
- `4humans/diagrams/processes/*.d2` と対応する `.svg`
  `4humans/diagrams/processes/frontend-execution-control-sequence-diagram.d2` と対応する `.svg` を更新した。
- `4humans/diagrams/overview-manifest.json`
  new detail diagram は追加していないため未更新。

## 結果

- `src-tauri/src/gateway/commands.rs` で fixture-backed execution-control / observe runtime を統合し、pause / resume / retry / cancel と provider-neutral observe 再読込を 1 本の route に揃えた。
- `src-tauri/tests/acceptance/provider-failure-retry/mod.rs` で `Running -> RecoverableFailed -> Retrying -> Running -> Completed` と `Running -> Paused -> Running -> Canceled` を integrated scenario proof として固定した。
- `src/ui/screens/*/index.test.ts` の shell-path tests と `src-tauri/tests/validation/execution-observe-snapshot/mod.rs` を最小更新し、shared composition / observe validation を現行 wiring に合わせた。
- `4humans/diagrams/structures/frontend-execution-control-slice-class-diagram.d2` と `4humans/diagrams/processes/frontend-execution-control-sequence-diagram.d2` を更新し、対応 `.svg` を再生成した。
- `SONAR_USER_HOME=.sonar-local sonar-scanner` は成功し、owned scope の open issue は 0 件だった。
- `SONAR_USER_HOME=.sonar-local CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all` は成功した。
