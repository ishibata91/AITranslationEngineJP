- workflow: impl
- status: completed
- lane_owner: codex
- scope: SQLite schema 管理を ad-hoc DDL から versioned migration へ移行し、backend 起動時の専用初期化責務へ集約する。

## Request Summary

- `sqlx` 標準 migration 機構を導入し、現行の `plugin_exports` と `plugin_export_raw_records` schema を versioned migration へ移す。
- Tauri backend 起動時に DB migration を一度だけ実行する backend-owned 初期化責務を追加する。
- `SqlitePluginExportRepository` から DDL を除去し、repository を DML / transaction 専任に戻す。
- `file_paths: string[]` の command 境界は維持し、migration 前提で persistence test を更新する。

## Decision Basis

- `docs/architecture.md`
- `docs/coding-guidelines.md`
- `docs/tech-selection.md`
- `docs/exec-plans/completed/2026-03-29-per-41-plugin-export-persistence-unit.md`

## Owned Scope

- `src-tauri/Cargo.toml`
- `src-tauri/migrations/`
- `src-tauri/src/lib.rs`
- `src-tauri/src/gateway/commands.rs`
- `src-tauri/src/infra/`
- `src-tauri/tests/`
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Out Of Scope

- DB path の恒久配置変更
- retention policy の定義
- frontend 契約や command DTO の変更
- `bootstrap-status` usecase への DB 初期化責務の混在
- `SqlitePool` 共有導入

## Dependencies / Blockers

- docs 正本の migration 方針は更新済みで、この実装では `docs/` を追加更新しない。
- Sonar gate の確認には `sonar-scanner` と open issue 取得 script の成功が必要。

## Parallel Safety Notes

- `src-tauri/Cargo.toml`、`src-tauri/src/lib.rs`、`src-tauri/src/infra/`、`src-tauri/tests/` は shared scope になりやすいため、worker handoff では allowed scope を限定する。

## UI

- UI 変更なし。frontend から見える command DTO と `file_paths: string[]` 契約は維持する。
- DB schema 準備は backend startup の内部責務へ移し、import / job command に初期化用 parameter や追加 DTO を持ち込まない。
- migration 失敗時は command ごとの schema 自動作成へフォールバックせず、backend 起動失敗として扱う。

## Scenario

- Tauri backend 起動時に `run()` の startup hook が execution cache 用 SQLite DB path を解決し、backend-owned initializer が同一 path に対して versioned migration を一度だけ適用してから command 受付へ進む。
- startup 時に DB open や migration 適用が失敗した場合は backend 起動自体を失敗させ、`bootstrap-status` や import command の通常経路へ schema 準備を逃がさない。
- `import_xedit_export_json` の成功時は既存の `file_paths: string[]` 入力境界を維持したまま、migration 済み DB に対して repository が raw record を transaction 保存する。
- command test や persistence test のうち成功系は、本番と同じ initializer を対象 DB path に先に適用してから repository / command を呼ぶ。未初期化 DB の負系 test は initializer を呼ばず、repository が missing schema をそのまま失敗として返すことを確認する。

## Logic

- execution cache path 解決は `gateway/commands.rs` の private helper から backend 内の共有 helper へ寄せ、startup initializer と per-request repository factory が同じ path 決定規則を使う。
- DB 初期化責務は infra か backend bootstrap 専用モジュールに置き、`src-tauri/migrations/` の versioned `.sql` を `sqlx` 標準 migration 機構で適用する。配布時の相対 path 依存を避けるため、実装は compile 時に `migrations/` を束ねる形を優先する。
- initializer は repository とは別 connection で migration を実行してよく、shared `SqlitePool` や long-lived connection を前提にしない。
- `src-tauri/src/lib.rs` が Tauri startup hook から initializer を呼ぶ唯一の production 入口になり、repository / use case / `bootstrap-status` は migration 実行責務を持たない。
- `SqlitePluginExportRepository` は SQLite open、query、transaction のみを担当し、`CREATE TABLE IF NOT EXISTS` を含む DDL を完全に除去する。migration 未適用 DB では table 不在 error を persistence error として返し、暗黙 schema 作成は行わない。

## Implementation Plan

- active plan 作成後、distill で関連コードと test 現状を最小限整理する。
- task-local design で migration 実行責務、起動フック、repository 境界、test 前提を固める。
- planning で ordered scope、owned scope、validation commands を短い brief に落とす。
- tests を先に migration 前提へ更新し、未初期化 DB が暗黙 schema 作成しないことも含めて固定する。
- backend 実装で migration file、initializer、起動フック、repository DDL 除去を反映する。
- required validation と Sonar / review / 4humans sync / plan close を完了する。

## Acceptance Checks

- repository 保存時に `CREATE TABLE` や schema 準備を実行しない。
- backend 起動初期化で DB migration が一度だけ走る。
- migration 適用後に既存 import persistence が成功する。
- migration 未適用 DB に対して repository は暗黙に schema を作らない。
- command DTO と frontend 契約は変わらない。

## Required Evidence

- active plan に task-local design と implementation brief が追記されている。
- migration file と backend-owned initializer 実装が追加されている。
- migration 前提の tests / fixtures 更新がある。
- structure / design harness、fmt、test、clippy、Sonar、single-pass review の結果が揃っている。

## 4humans Sync

- `4humans/quality-score.md`
  既存品質評価の範囲で追加更新なし。required validation と startup failure test を追加し、品質状態は改善したが新しい恒久課題は作らなかった。
- `4humans/tech-debt-tracker.md`
  既存の DB path / retention policy debt が継続するだけで、新規 debt 追加は不要。
- `4humans/class-diagrams/*.d2` と対応する `.svg`
  backend startup に initializer を追加したが、公開している境界図の責務分割自体は既存 docs と整合し、図更新は不要と判断した。
- `4humans/sequence-diagrams/*.d2` と対応する `.svg`
  起動時 migration フローは backend 内部初期化で完結し、現行の human-facing sequence 図更新までは不要と判断した。

## Outcome

- `sqlx` の versioned migration を導入し、`src-tauri/migrations/0001_execution_cache_base.sql` へ現行 `plugin_exports` / `plugin_export_raw_records` schema を移した。
- execution cache path 解決と migration 初期化を `src-tauri/src/infra/execution_cache.rs` へ集約し、Tauri backend startup (`src-tauri/src/lib.rs`) が production で migration を実行する唯一の入口になった。
- `SqlitePluginExportRepository` から `CREATE TABLE IF NOT EXISTS` を削除し、repository を DML / transaction 専任に戻した。
- persistence / command / acceptance tests を migration 前提へ更新し、未初期化 DB では repository と command が暗黙 schema 作成せず失敗することを固定した。
- startup failure 契約の自動検証として、execution cache open failure と migration failure の 2 本の startup test を `src-tauri/src/lib.rs` に追加した。
- Validation:
  - PASS `python3 scripts/harness/run.py --suite structure`
  - PASS `python3 scripts/harness/run.py --suite design`
  - PASS `cargo fmt --manifest-path ./src-tauri/Cargo.toml --all --check`
  - PASS `cargo test --manifest-path ./src-tauri/Cargo.toml --test import_xedit_export_persistence -- --nocapture`
  - PASS `cargo test --manifest-path ./src-tauri/Cargo.toml --test xedit_export_importer -- --nocapture`
  - PASS `cargo test --manifest-path ./src-tauri/Cargo.toml --test acceptance -- --nocapture`
  - PASS `cargo clippy --manifest-path ./src-tauri/Cargo.toml --all-targets --all-features -- -D warnings`
  - PASS `cargo test --manifest-path ./src-tauri/Cargo.toml --all-features`
  - PASS `sonar-scanner`
  - PASS `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/Cargo.toml src-tauri/src/lib.rs src-tauri/src/gateway/commands.rs src-tauri/src/infra src-tauri/tests` (`openIssueCount: 0`)
  - FAIL `python3 scripts/harness/run.py --suite all`
    sandbox の cargo cache 書き込み制約による失敗で、required individual checks と code-level validation は通過した。
- Review:
  - single-pass review は startup failure 契約の test 不足で reroute を返した。
  - reroute を反映し、追加 review は実施していない。
- Residual risk:
  - `src-tauri/Cargo.lock` は `sqlx` feature 追加に伴う依存解決差分を含む。
  - startup migration failure test は forced failure 注入であり、実 DB 破損状態の再現までは扱っていない。
