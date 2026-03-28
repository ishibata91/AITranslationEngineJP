# Executable Specs

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md)

この文書は、細かな仕様や制約を「後で実行して確かめられる形」に寄せるための入口とする。
詳細仕様は長い説明文で増やすのではなく、テスト、acceptance checks、fixture、検証コマンドへ落とす。

## Principles

- 細かな振る舞いは、可能な限りテストで分かる形にする
- 文書はテストや acceptance checks を作るための最小限の契約だけを書く
- 仕様変更では、必要なら対応する test case や acceptance checks も同時に更新する
- 実装がまだない領域では、先に期待結果、失敗条件、観測点を記録する
- review で繰り返し見つかる指摘は、可能な限り harness や acceptance checks に昇格する

## Record Here

- どの種類のテストで何を担保するか
- acceptance checks に必ず入れるべき観点
- fixture や sample input / output の扱い
- 実行可能仕様に昇格させるべき制約

## Current Policy

- task-local な詳細設計や一時的な実装判断は `docs/exec-plans/` に置いてよいが、完了後も保持すべき詳細仕様の正本にはしない
- 実装 task で必要になる `UI` / `Scenario` / `Logic` は active plan の section に置き、別の `changes/` や packet artifact に分けない
- 完了後も残すべき詳細な振る舞い、制約、受け入れ条件は、ここで入口を定義し、対応する tests / acceptance checks / validation commands を正本として管理する
- UI、実行、DTO、状態遷移の細かな制約は、将来的に対応する test と validation command で表現する
- plan には `Acceptance Checks` を必須で持たせ、`UI` / `Scenario` / `Logic` を含む詳細仕様の一時的な置き場にする
- plan には `Required Evidence` を持たせ、lane close が pass/reroute を判定できるようにする
- 永続ルールだけを `spec.md` と `architecture.md` に残し、細かな分岐条件はここからテストへ寄せる
- TypeScript / Svelte の通常 lint は `Oxlint` を基本入口とし、未使用変数、未使用 import、type import hygiene を検出対象にする
- repo 固有の import 境界は `ESLint` の `no-restricted-imports` または repository-local rule で検証し、同一層の別 feature / slice / package の internal module import を失敗扱いにする
- 未参照 export / file / dependency は `Knip` を基本入口として検出し、tests / spec entrypoint / fixtures / generated code の明示 allowlist を除いて削除対象として扱う
- `Semgrep` は将来の責務境界 / 禁止 API 向け validation command として採用し、初期は `TS` と `Rust` の両方を report-first で観測する
- `Semgrep` の将来の config entrypoint は repo 管理下の専用パスに置き、`Semgrep Registry` の配布ルールと repo-local rule を併用できる構成にする
- `Semgrep` の将来の実行では `tests` / `fixtures` / `generated` などの path exclude を持てるようにし、誤検知を分類した後にだけ gate へ昇格させる
- `Semgrep` の rule 結果は、`そのまま採用`、`override / tune`、`repo-specific local rule` に分類して整理する
- Rust 側の未参照コードと未使用要素は `cargo clippy --all-targets --all-features -- -D warnings` を基本入口として失敗扱いにする
- `Knip` の report に出た未参照 export / file / dependency は同一変更で削除または allowlist 理由を追加し、放置を許可しない
- `Knip --fix` はローカル cleanup 手段として許可するが、ファイル削除を伴う結果は review 可能な差分として残す
- ドメイン / アプリケーション層のルールは Rust の `cargo test` を基本入口として表現する
- UI コンポーネントと画面内の振る舞いは `Vitest` と `@testing-library/svelte` で表現する
- デスクトップ統合の acceptance checks は `tauri-driver` と `WebdriverIO` を基本入口とする
- review は runtime 品質の総合点ではなく、`仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 観点だけを判定する
- xEdit importer の acceptance checks では、`extractData.pas` の raw JSON に含まれる `cells` 空配列と `voicetype` 互換項目を読み込んでも canonical DB モデルが崩れないことを確認する
- 複数入力ファイルの acceptance checks では、1 つの `TRANSLATION_JOB` が複数 `PLUGIN_EXPORT` を参照しつつ、出自情報を保持したまま `TRANSLATION_UNIT` を生成できることを確認する
- 出力 writer の acceptance checks では、複数フィールドを持つ 1 レコードから `TRANSLATION_UNIT` を field 単位で生成し、xTranslator XML の `EDID` / `REC` / `FIELD` / `FORMID` / `Source` / `Dest` / `Status` を lossless に再構成できることを確認する
- ペルソナ生成の acceptance checks では、`MASTER_PERSONA` と `JOB_PERSONA_ENTRY` が混線せず、mod 追加 NPC のジョブ内ペルソナを UI 観測用 DTO へ分けて出せることを確認する
- 実装初期化フェーズでは、repo root package の `lint` / `test` / `build` と `src-tauri/` の `cargo fmt` / `cargo clippy` / `cargo test` が execution harness から呼ばれる状態を最小成立条件とする
- 実装初期化フェーズの boundary lint は、frontend の `eslint.config.mjs` で `Gateway` 以外からの Tauri API 利用を禁止し、UI view から gateway 直接依存を禁止する
- 実装初期化フェーズの `Semgrep` は repo root `semgrep/semgrep.yml` を entrypoint とし、`report-first` の観測コマンドに留める
