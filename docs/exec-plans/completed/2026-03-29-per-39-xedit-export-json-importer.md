- workflow: impl
- status: completed
- lane_owner: codex
- scope: Add an xEdit export JSON importer entrypoint with room for multiple input files and minimal validation.

## Request Summary

- Implement Linear `PER-39` `[Phase 1][01] Add xEdit export JSON importer`.
- Add an importer entrypoint so xEdit exported structured JSON can be accepted as translation workflow input data.
- Keep the structure extensible for multiple input files.
- Add the minimum validation required to prove the importer contract.

## Decision Basis

- `docs/spec.md`
- `docs/architecture.md`
- `docs/tech-selection.md`
- Linear `PER-39`

## UI

- `AppShell` と既存画面構成は今回の scope に含めない。
- importer の入口は backend command と application use case を先に追加し、UI からは後続 task で接続できる input port / gateway 形に寄せる。
- 多ファイル対応は UI 部品を先に作らず、command / DTO が単一 file import から複数 file import へ拡張できる contract shape を優先する。

## Scenario

- ユーザー起点の最初の導線は「xEdit export JSON の path 群を importer 入口へ渡す」で統一する。
- Phase 1 では最小成立を優先し、最初の accepted path は 1 file でも、入口 contract は file collection を受けられる形にする。
- import 成功時は、translation workflow が次段で利用できる入力データ表現を返す。invalid JSON または importer 前提を満たさない入力は error として返し、部分成功は扱わない。

## Logic

- 実装責務は backend 中心とし、`gateway/commands.rs -> application use case -> domain/infrastructure boundary` の依存方向を守る。
- xEdit JSON は load 時に型検証し、translation workflow の canonical な入力単位へ正規化する前段として importer request / result DTO を定義する。
- 最小 validation は「xEdit export JSON として読めること」「translation workflow へ渡す必須識別情報と翻訳対象データを欠落なく保持できること」を確認対象にする。
- 原本 JSON を正本として保持する前提を崩さず、今回の task では永続化全体を完結させるより importer 入口と validation 境界を固定する。

## Implementation Plan

- Ordered scope 1: Define backend-owned importer request / result contract around a file-path collection, and pin the minimum xEdit schema subset needed to preserve `target_plugin`, record identifiers, and translatable source fields.
- Ordered scope 2: Add backend modules for the importer along `gateway/commands.rs -> application -> domain / infra`, register the new Tauri command, and keep failure behavior whole-request on invalid input.
- Ordered scope 3: Add the smallest Rust test coverage that proves one valid xEdit export JSON import succeeds and one invalid input fails, using inline sample JSON or a minimal fixture if reuse is cleaner.
- Ordered scope 4: Add a TypeScript shared contract only if the Tauri boundary cannot stay backend-local without it; do not add UI wiring, screen state, or app-shell changes in this ticket.
- Owned scope: `src-tauri/src/gateway/commands.rs`, `src-tauri/src/gateway/mod.rs`, `src-tauri/src/lib.rs`, new importer modules under `src-tauri/src/application/`, `src-tauri/src/domain/`, `src-tauri/src/infra/`, and only boundary contract files in `src/shared/` or `src/gateway/tauri/` if the command signature requires them.
- Required reading: `docs/exec-plans/active/2026-03-29-per-39-xedit-export-json-importer.md`, `docs/spec.md`, `docs/architecture.md`, `docs/tech-selection.md`, `docs/er-draft.md`, `src-tauri/src/gateway/commands.rs`, `src-tauri/src/lib.rs`, `src-tauri/src/application/bootstrap/get_bootstrap_status.rs`, `src-tauri/Cargo.toml`.
- Validation commands: `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `powershell -File scripts/harness/run.ps1 -Suite execution`.

## Acceptance Checks

- Structure harness passes before implementation starts.
- A backend integration test proves a valid xEdit export JSON file path can be imported through the new entrypoint and returns one `plugin_export` plus `translation_units` that preserve `target_plugin`, `form_id`, `editor_id`, `record_signature`, `field_name`, and `source_text`.
- A backend integration test proves an input file that violates the importer preconditions fails the whole request and surfaces a validation error instead of partial success.
- The importer request contract remains `file_paths: string[]` so later multi-file support can extend the same entrypoint shape without redesign.

## Required Evidence

- Active plan updated with task-local design and implementation brief.
- Added or updated tests / validation commands.
- Harness and implementation validation results.
- Sonar remediation status and single-pass implementation review result.

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- Added a backend-owned xEdit export JSON importer entrypoint at the Tauri command boundary with request shape `file_paths: string[]`.
- Added importer application/domain/infra modules that load xEdit JSON, preserve `target_plugin`, and return `translation_units` for translatable fields without widening persistence scope.
- Added importer fixtures and Rust tests for valid import, invalid missing `target_plugin`, and valid import with blank `editor_id` to preserve lossless input behavior.
- Validation passed: `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`, `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`, `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`, `sonar-scanner`, and SonarQube MCP returned `0` open issues for the touched backend files.
- Remaining repo blocker outside owned scope: `scripts/eslint/repository-boundary-plugin.test.mjs` still fails in the existing execution harness, so `powershell -File scripts/harness/run.ps1 -Suite execution` is not green yet for unrelated reasons.
