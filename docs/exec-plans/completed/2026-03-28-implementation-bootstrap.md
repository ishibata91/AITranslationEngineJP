# Implementation Bootstrap

- workflow: heavy
- status: completed
- architect: codex
- research: Volta
- coder: codex
- reviewer: architect
- scope: minimal Rust + Tauri + Svelte initialization, folder skeleton, lint config bootstrap, execution harness integration

## Request Summary

- language / project の最小初期化を行う
- 責務境界を表現できる最小フォルダ構成を作る
- 既存 lint 契約に沿う初期設定を追加する
- execution harness から lint / test / build を呼べる状態にする

## Decision Basis

- `docs/tech-selection.md` で `Tauri 2`、`Rust`、`Svelte 5`、`TypeScript`、`Oxlint`、`ESLint`、`Knip`、`Semgrep`、`clippy` が採用済み
- `docs/architecture.md` で backend の 4 層と UI 内の `Gateway` 境界が定義済み
- `docs/lint-policy.md` と `docs/executable-specs.md` で `Semgrep` は report-first、`Oxlint` / `ESLint` / `Knip` / `clippy` は gate 候補として定義済み
- 現在の repo には `Cargo.toml` と `package.json` が存在せず、execution harness は対象なしで skip する状態

## Investigation Summary

- Facts:
  - 現在の repo に実装コード、frontend package、Rust crate はまだ存在しない
  - execution harness は `Cargo.toml` と `package.json` を自動検出して `cargo fmt` / `cargo clippy` / `cargo test` と `lint` / `test` / `build` script を実行する
  - Tauri 2 の標準構成は repo root 側に frontend、`src-tauri/` 側に Rust crate を置く
- Options:
  - frontend を repo root package にする
  - frontend を `frontend/` 等の subdirectory package にする
  - Rust を単一 crate で始める
  - Rust workspace を初期から導入する
- Risks:
  - frontend root を過剰分割すると lint / build / Tauri path が早期に複雑化する
  - `Semgrep` を gate 化すると既存 docs 契約と矛盾する
  - Rust workspace を早期導入すると初期変更量が増える
- Unknowns:
  - frontend root 名
  - Rust 単一 crate / workspace の選択
  - `Semgrep` config directory の具体名

## Unknown Classification

- Blocking:
  - frontend root を repo root に置くか subdirectory package に置くか
  - Rust を単一 crate で始めるか workspace にするか
- Non-blocking:
  - `Semgrep` config file 名
  - 初期サンプル UI / command の命名

## Assumptions / Defaults

- frontend root は repo root の単一 package とする
- Rust backend は `src-tauri/` 配下の単一 crate とする
- `Semgrep` config は repo root の `semgrep/` ディレクトリに置く
- TS 側は `ui` と `gateway` を主軸にし、backend core の `domain` / `infra` は Rust 側で表現する

## Plan Ready Criteria

- 初期レイアウト、config 配置、gate/report-first 境界を固定できている
- docs sync 先が明示されている
- coder が追加実装時に仕様判断を増やさず進められる

## Implementation Plan

- repo root に frontend package を追加し、`src/` と `src-tauri/` を持つ Tauri 2 + Vite + Svelte 5 の最小構成を作る
- frontend 側は `src/ui/`、`src/application/`、`src/gateway/`、`src/shared/` の最小構成を置き、`Gateway` 経由で Tauri command を呼ぶ
- Rust 側は `src-tauri/src/application/`、`src-tauri/src/domain/`、`src-tauri/src/infra/`、`src-tauri/src/gateway/` の最小構成を置き、Tauri command から application usecase を呼ぶ
- `Oxlint` / `ESLint` / `Knip` の config を frontend root に置き、path-based boundary lint を初期化する
- `Semgrep` report-first config を repo root `semgrep/` に置き、package script から別 command として実行可能にする
- architecture / lint / executable specs / quality record に concrete bootstrap layout と gate policy を同期する

## Delegation Map

- Research:
  - repo 現況、既存 docs 契約、blocking unknown の調査
- Coder:
  - scaffold 作成、config 実装、docs sync、validation 実行
- Worker:
  - なし

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
- 実行可能なら `powershell -File scripts/harness/run.ps1 -Suite execution`
- repo root package の `lint` / `test` / `build` が通る
- `src-tauri` の `cargo fmt --all --check`、`cargo clippy --all-targets --all-features -- -D warnings`、`cargo test --all-features` が通る

## Required Evidence

- 実際に `Cargo.toml` と `package.json` が追加され execution harness から検出されること
- TS / Rust の最小 directory layout が docs 契約と一致していること
- `Semgrep` が report-first のままで `lint` gate に入っていないこと
- docs sync と quality record が追従していること

## Reroute Trigger

- frontend を multi-package workspace にしないと boundary lint を表現できないと判明した場合
- Tauri 2 + Svelte 5 の最小初期化が現行 toolchain で成立せず、技術選定の更新が必要になった場合
- layer / gateway 切り方が既存 architecture と矛盾すると判明した場合

## Docs Sync

- `docs/architecture.md`
- `docs/executable-specs.md`
- `docs/lint-policy.md`
- `4humans/quality-score.md`
- 必要なら `4humans/tech-debt-tracker.md`

## Record Updates

- `docs/exec-plans/active/2026-03-28-implementation-bootstrap.md`

## Outcome

- repo root に frontend package を追加し、`src/` と `src-tauri/` を持つ最小 Tauri 2 + Svelte 5 + TypeScript + Rust bootstrap を作成した
- TS 側は `src/ui/`、`src/application/`、`src/gateway/`、`src/shared/` に分割し、`Gateway` 経由の bootstrap status 取得で責務境界を可視化した
- Rust 側は `src-tauri/src/application/`、`domain/`、`infra/`、`gateway/` を作成し、Tauri command から application usecase を呼ぶ最小配線を追加した
- `Oxlint` / `ESLint` / `Knip` の gate 用 config と `Semgrep` report-first config を追加した
- execution harness を toolchain 未導入でも failure evidence を返す形に改善した

## Verification

- `npm run lint`
- `npm run test`
- `npm run build`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
- `powershell -File scripts/harness/run.ps1 -Suite execution`

## Remaining Gaps

- `cargo` / `rustc` が実行環境に未導入のため、`cargo fmt` / `cargo clippy` / `cargo test` は未実行
- `Semgrep` は契約通り report-first であり、execution harness の gate にはまだ入れていない
