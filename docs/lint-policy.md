# Lint Policy

関連文書: [`index.md`](./index.md), [`tech-selection.md`](./tech-selection.md), [`architecture.md`](./architecture.md), [`executable-specs.md`](./executable-specs.md)

この文書は、lint が何を管理し、何を管理しないかをまとめるための正本とする。
個別ツールの採用自体は `tech-selection.md`、検証入口と失敗条件は `executable-specs.md` を正本とし、本書は lint の責務範囲を一覧化する。

## Lint が管理するもの

- import 境界: 同一層の別 feature / slice / package の internal module への直接 import を禁止する
- import hygiene: 未使用 import、type import の崩れ、静的に検出できる危険な import パターンを検出する
- 未参照コード: 未使用変数、未参照 export、未参照 file、未使用 dependency を検出する
- Rust の静的品質: 未使用要素、dead code、warning 扱いのコード品質問題を失敗扱いにする
- allowlist 管理: tests / spec entrypoint / fixtures / generated code など、明示的に除外を許可する対象を管理する

## Lint が管理しないもの

- runtime の正しさや振る舞い: `cargo test`、`Vitest`、acceptance checks で担保する
- UI の表示仕様や操作結果: component test、screen-level test、end-to-end check で担保する
- 業務フローや受け入れ条件の成立: `docs/executable-specs.md` と対応する tests / acceptance checks で担保する
- formatting の整形結果そのもの: `rustfmt` などの formatter が担当し、lint の主責務には含めない

## Tool Ownership

- `Oxlint`: TypeScript / Svelte の通常 lint。未使用変数、未使用 import、type import hygiene を担当する
- `ESLint Flat Config + repository-local rule`: repo 固有の import 境界と、path-based な architectural lint を担当する
- `Knip`: 未参照 export / file / dependency の検出と cleanup の入口を担当する
- `Semgrep`: `TS` / `Rust` 両方で、責務境界と禁止 API の repository-specific pattern lint を補完する
- `cargo clippy --all-targets --all-features -- -D warnings`: Rust 側の未使用要素と warning を失敗扱いにする

## Initial Config Placement

- frontend lint config は repo root の frontend package に置く
  - `eslint.config.mjs`
  - `knip.json`
  - `.oxlintrc.json`
- Rust lint / format config は `src-tauri/` を基準に置く
  - `Cargo.toml`
  - `clippy.toml` や `rustfmt.toml` は必要になった時だけ `src-tauri/` に追加する
- `Semgrep` config は repo root の `semgrep/` を dedicated path として置く

## Initial Gate Split

- gate に入れるもの:
  - repo root package の `lint`
  - repo root package の `test`
  - repo root package の `build`
  - `src-tauri/` の `cargo fmt --all --check`
  - `src-tauri/` の `cargo clippy --all-targets --all-features -- -D warnings`
  - `src-tauri/` の `cargo test --all-features`
- report-first に留めるもの:
  - `Semgrep` の `semgrep/semgrep.yml`
  - future import graph / cycle 専用解析

初期 gate の責務は `Oxlint` / `ESLint` / `Knip` / `clippy` に固定し、`Semgrep` は補助観測層のままとする。

## Cleanup Policy

- `Knip` の report に出た未参照 export / file / dependency は、同一変更で削除する
- 削除しない場合は、allowlist へ理由付きで追加する
- `Knip --fix` はローカル cleanup 手段として使ってよいが、ファイル削除を伴う結果は review 可能な差分として残す
- allowlist は恒久逃げ道にせず、tests / spec entrypoint / fixtures / generated code のように静的解析で誤検出しやすい対象へ限定する

## Semgrep Role

- `Semgrep` は既存の `Oxlint` / `ESLint` / `Knip` / `clippy` を置き換えず、追加層として使う
- 初期導入の主目的は `責務境界 + 禁止 API` であり、complexity-first lint にはしない
- ルール表現の主軸は命名規約ではなくディレクトリ基準とし、`TS` と `Rust` の両方を対象にする
- import graph や cycle の主担当は `Semgrep` に移さず、専用 lint / 静的解析に残す
- 最初の rule source は `Semgrep Registry` を優先し、不足分だけ repo-local rule を追加する
- 初期運用は `report-first` とし、docs-only フェーズでは gate や CI fail 条件に入れない

## Semgrep First-Wave Targets

- service / usecase / application code での SQL 直書き
- domain code からの DB / HTTP / Tauri / file I/O 実装への直接アクセス
- `Gateway` 以外からの `Tauri invoke` / event API の直接利用
- production code から `test` / `fixture` / `generated` への直接依存
- Registry 既存ルールで検出できる generic code-smell のうち、責務境界と禁止 API に近いもの

## Semgrep Rule Lifecycle

- まず配布ルールを試し、結果を `そのまま使える rule`、`override / tune が必要な rule`、`repo-specific local rule` に分ける
- repo-specific local rule は、配布ルールで埋まらない責務境界違反だけに絞る
- naming-based rule は v1 では採用せず、directory-based rule で不足が明確になった時だけ追加する

## Recommended Next Gates

- 層方向違反: `UI -> infra`、`application -> concrete infra`、`domain -> UI / infra / SDK` の依存を禁止する
- Gateway 境界違反: UI から `Tauri invoke` / event API を直接呼ばず、`Gateway` 経由だけを許可する
- DTO 境界違反: UI が backend / domain の内部型を直接 import せず、DTO / query model だけを見るようにする
- feature 横断 internal import: 同一層の別 feature / slice / package の internal module import を禁止し、公開境界経由だけを許可する
- 循環依存: feature 間、layer 間の cycle を静的解析で失敗扱いにする
- test / fixture / generated の逆流: production code から `*.test.*`、`fixtures/`、`generated/` を直接 import するのを禁止する

## Record Split

- 採用ツールと品質基盤の決定: [`tech-selection.md`](./tech-selection.md)
- import 境界の設計原則: [`architecture.md`](./architecture.md)
- validation command、失敗条件、allowlist 契約: [`executable-specs.md`](./executable-specs.md)
- task-local な一時判断: [`exec-plans/`](./exec-plans/active/README.md)
