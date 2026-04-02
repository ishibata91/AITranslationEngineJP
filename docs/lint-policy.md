# Lint Policy

関連文書: [`index.md`](./index.md), [`tech-selection.md`](./tech-selection.md), [`architecture.md`](./architecture.md)

この文書は、lint が何を管理し、何を管理しないかをまとめるための正本とする。
個別ツールの採用自体は `tech-selection.md`、検証入口と失敗条件は対応する tests / acceptance checks / validation commands を正本とし、本書は lint の責務範囲を一覧化する。

## Lint が管理するもの

- import 境界: layer boundary と production code から `test` / `fixtures` / `generated` への逆流を管理し、同一層の別 public root への import は `target root` の `index` file または `target root` 直下 file だけを許可する
- import hygiene: 未使用 import、type import の崩れ、静的に検出できる危険な import パターンを検出する
- 未参照コード: 未使用変数、未参照 export、未参照 file、未使用 dependency を検出する
- Rust の静的品質: 未使用要素、dead code、warning 扱いのコード品質問題を失敗扱いにする
- allowlist 管理: tests / spec entrypoint / fixtures / generated code など、明示的に除外を許可する対象を管理する

## Lint が管理しないもの

- runtime の正しさや振る舞い: `cargo test`、`Vitest`、acceptance checks で担保する
- UI の表示仕様や操作結果: component test、screen-level test、end-to-end check で担保する
- 業務フローや受け入れ条件の成立: 対応する tests / acceptance checks / validation commands で担保する
- formatting の整形結果そのもの: `rustfmt` などの formatter が担当し、lint の主責務には含めない

## Tool Ownership

- `Oxlint`: TypeScript / Svelte の通常 lint。未使用変数、未使用 import、type import hygiene を担当する
- `ESLint Flat Config + repository-local rule`: repo 固有の import 境界と、path-based な architectural lint を担当する
- `Knip`: 未参照 export / file / dependency の検出と cleanup の入口を担当する
- `SonarScanner + Sonar CLI`: `TS` / `Svelte` / `Rust` を対象に、code smell、complexity、security hotspot などの server-side issue gate を担当する
- `cargo fmt --all --check`: `src-tauri/` の formatter 出力が正本どおりかを gate で確認する
- `cargo clippy --all-targets --all-features -- -D warnings`: Rust 側の未使用要素と warning を失敗扱いにする

## Initial Config Placement

- frontend lint config は repo root の frontend package に置く
  - `eslint.config.mjs`
  - `knip.config.mjs`
  - `.oxlintrc.json`
  - path allowlist と tool-specific ignore は `config/lint/allowlists.json` を正本にする
- Rust lint / format config は `src-tauri/` を基準に置く
  - 現在の repo-local 正本は `src-tauri/Cargo.toml`
  - `clippy.toml` や `rustfmt.toml` は repo-local override が必要になった時だけ `src-tauri/` に追加する
- Sonar scanner config は repo root の `sonar-project.properties` を正本とする

## Initial Gate Split

- gate に入れるもの:
  - repo root package の `gate:execution`
  - `gate:execution` の内訳として `lint`、`src-tauri/` の `cargo fmt --all --check`、`src-tauri/` の `cargo clippy --all-targets --all-features -- -D warnings`、repo root の `scan:sonar`、`test`、`src-tauri/` の `cargo test --all-features`、`build`
- report-first に留めるもの:
  - future import graph / cycle 専用解析

初期 gate の責務は `Oxlint` / `ESLint` / `Knip` / `rustfmt check` / `SonarScanner + Sonar CLI` / `clippy` に固定する。

## Cleanup Policy

- `Knip` の report に出た未参照 export / file / dependency は、同一変更で削除する
- 削除しない場合は、allowlist へ理由付きで追加する
- `Knip --fix` はローカル cleanup 手段として使ってよいが、ファイル削除を伴う結果は review 可能な差分として残す
- allowlist は恒久逃げ道にせず、tests / spec entrypoint / fixtures / generated code のように静的解析で誤検出しやすい対象へ限定する
- repo root frontend lint の path ignore と `Knip` binary ignore は `config/lint/allowlists.json` を正本とし、`ESLint` / `Knip` / `Oxlint runner` はその値を参照する

## Sonar Gate Role

- `SonarScanner + Sonar CLI` は既存の `Oxlint` / `ESLint` / `Knip` / `clippy` を置き換えず、構造解析系 lint で取り切れない issue を補完する追加層として gate に入れる
- 初期導入の主目的は code smell、complexity、security hotspot などの server-side issue を継続的に潰すことであり、repo 固有の責務境界や禁止 API の主担当にはしない
- repo root の `scan:sonar` は `sonar-project.properties` を正本にして `sonar-scanner` を実行し、execution harness はその script を lint 後段で呼ぶ
- implementation lane は `sonar list issues --project ishibata91_AITranslationEngineJP --format json` を helper script 経由で読み、`status == OPEN` の issue が残る限り implementing skill に差し戻す
- import graph や cycle の主担当は Sonar に移さず、専用 lint / 静的解析に残す
- project context は `sonar-project.properties` の project key を正本とし、Sonar CLI query でも同じ key を使う

## Sonar First-Wave Targets

- high / medium severity の code smell
- cognitive complexity や nested control flow の肥大化
- duplicated logic や保守性低下を示す generic issue
- security hotspot と unsafe pattern
- 構造解析系 lint では取り切れない generic quality issue

## Sonar Gate Lifecycle

- execution harness は repo root の `gate:execution` を優先実行し、fallback の時だけ `lint` / `test` / `build` と `Cargo.toml` command を個別実行する
- `scan:sonar` は `lint` と Rust lint の後段で実行し、`sonar-project.properties` に定義した project context で server-side analysis を更新する
- scanner 実行後は implementation lane が `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP` を使って `status == OPEN` issue を取得する
- open issue が残る間は owned scope に応じて `implementing-frontend` または `implementing-backend` に戻して修正する
- false positive や severity tune は repo-local ルール増設ではなく Sonar project 側の設定を優先する
- gate failure は同一変更で修正するか、Sonar project 側で理由付き suppression を行う

## Recommended Next Gates

- 層方向違反: `UI -> infra`、`application -> concrete infra`、`domain -> UI / infra / SDK` の依存を禁止する
- Gateway 境界違反: UI から `Tauri invoke` / event API を直接呼ばず、`Gateway` 経由だけを許可する
- DTO 境界違反: UI が backend / domain の内部型を直接 import せず、DTO / query model だけを見るようにする
- feature 横断 internal import: 同一層の別 feature / slice / package の internal module import を禁止し、公開境界経由だけを許可する
- 循環依存: feature 間、layer 間の cycle を静的解析で失敗扱いにする
- test / fixture / generated の逆流: production code から `*.test.*`、`fixtures/`、`generated/` を直接 import するのを禁止する

ESLint gate は `src/ui` / `src/application` / `src/gateway` / `src/shared` の layer boundary・reverse-flow・same-layer public/internal import boundary を対象にする。

## Record Split

- 採用ツールと品質基盤の決定: [`tech-selection.md`](./tech-selection.md)
- import 境界の設計原則: [`architecture.md`](./architecture.md)
- validation command、失敗条件、allowlist 契約: 対応する tests / acceptance checks / validation commands
- task-local な一時判断: [`exec-plans/`](./exec-plans/active/README.md)
