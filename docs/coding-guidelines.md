# コーディング規約

関連文書: [`index.md`](./index.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)

本書は、`AITranslationEngineJp` の実装時に守るべきコーディング規約を定義する。
外部の Tauri + Svelte + TypeScript 向けガイドを土台にしつつ、本 repo の `Tauri 2`、`Svelte 5`、`TypeScript`、`Rust`、4 層構成に合わせて正本化する。
外部ガイドの旧版前提や repo 非適合な語彙は持ち込まず、この repo で必要な規約だけを残す。

## 1. 基本原則

- 実装は `UI層`、`アプリケーション層`、`ドメイン層`、`インフラ層` の境界を壊さない
- フロントエンドは `src/ui`、`src/application`、`src/gateway`、`src/shared` の責務分離を維持する
- バックエンドは `src-tauri/src/application`、`src-tauri/src/domain`、`src-tauri/src/infra`、`src-tauri/src/gateway` の責務分離を維持する
- 実装詳細よりも、tests / acceptance checks / validation commands で検証可能な振る舞いを優先する
- コメントは `何をしているか` ではなく、`なぜその判断が必要か` を短く補足する
- 命名は略語より意味を優先し、DTO、Port、Gateway、UseCase の役割が読める名前にする

## 2. Tauri 2 固有ルール

- 権限制御は `src-tauri/capabilities/` と `src-tauri/permissions/` を正本にする
- フロントエンドから Rust を呼ぶ入口は `@tauri-apps/api/core` の `invoke` と event API に限定し、UI から直接呼ばず `src/gateway/tauri/` に閉じ込める
- 新しい native API や plugin を追加する時は、必要最小限の `src-tauri/capabilities/` と `src-tauri/permissions/` を定義し、広い権限を雑に許可しない
- `tauri.conf.json` の `build.beforeDevCommand`、`build.beforeBuildCommand`、`build.devUrl`、`build.frontendDist` を repo root の Vite 設定と整合させる
- Tauri command は transport boundary として扱い、業務ロジックや DB 操作を `#[tauri::command]` 関数へ直接書き込まない
- command の入出力は DTO で固定し、`String` や匿名 object の投げ合いで境界を曖昧にしない
- Tauri plugin を使う時は、採用理由、権限範囲、代替不能性を明示してから追加する

## 3. TypeScript / Svelte 規約

- `Svelte 5` と `TypeScript` を前提にし、`any` と無検証の type assertion を常用しない
- `.svelte` は表示とイベント配線に集中させ、画面フローやデータ取得判断は `Screen UseCase` へ寄せる
- `Presenter / View` は pure に保ち、Tauri API、`invoke`、永続化、外部通信へ直接依存しない
- `Screen Store` は表示状態の保持に専念し、副作用の起点にしない
- 画面初期化や再読込は `Screen UseCase` の public API 経由で起動する
- `src/shared/contracts/` の型を UI と Gateway の共有 DTO とし、Rust 側の内部型を直接前提にしない
- 別 feature の内部 module を直接 import せず、公開 root の `index.ts` か root 直下 file を経由する
- コンポーネント名は `PascalCase`、変数と関数は `camelCase`、定数は用途が明確な時だけ `UPPER_SNAKE_CASE` を使う
- 非同期処理は `async` / `await` を基本とし、UI 側では fire-and-forget を乱用せず、失敗時の表示責務を決める

## 4. Rust / Backend 規約

- `domain` は最も内側の方針として保ち、Tauri、SQLite、ファイル I/O、HTTP 実装へ依存しない
- `application` は UseCase と DTO を持ち、外部依存は port / trait 越しに扱う
- `infra` は DB、ファイル、runtime、外部 API の具体実装だけを持ち、上位層へインフラ都合の型を漏らさない
- `gateway` は Tauri command と application usecase の接続点に限定し、ドメイン判断を持たない
- `Result<T, E>` を使って失敗を明示し、`panic!` や `unwrap()` を通常フローの制御に使わない
- 非同期 command や I/O は `tokio` を前提にし、blocking 処理を UI 応答経路へ混ぜない
- `serde` / `serde_json` で境界 DTO を定義し、暗黙の field 名変換に依存しすぎない
- 永続化は `sqlx` を通し、SQL と repository の責務を混ぜない

## 5. 通信とデータ境界

- この repo の frontend と backend の主通信路は HTTP ではなく Tauri IPC とする
- 外部サービスとの通信が必要な場合でも、UI から直接叩かず Rust 側の application / infra に寄せる
- request / response は DTO と validation で形を固定し、 optional field の意味を曖昧にしない
- versioning が必要な外部連携は Rust 側 adapter に閉じ込め、UI 契約へ直接漏らさない
- エラーは user-facing message と internal diagnostic を分け、表示文言に内部詳細をそのまま流さない

## 6. セキュリティ規約

- native 権限、ファイルアクセス、外部プロセス、shell 実行、HTTP client は最小権限で扱う
- `capabilities` と `permissions` は window / webview ごとに必要最小限へ絞る
- frontend 入力は UI で整形しても、最終 validation は Rust 側で再度行う
- 機密値やローカルパスをログやエラーメッセージへ無加工で出さない
- CSP、HTTP header、plugin permission は `便利だから広く開ける` ではなく、必要な操作だけを明示する

## 7. テストと品質ゲート

- 振る舞いを変える変更では、対応する `Vitest`、`cargo test`、acceptance checks、validation commands を同じ変更で更新する
- UI の画面フローは `Screen UseCase` と `Store` の unit test を優先し、Tauri API は gateway 境界で mock する
- `.svelte` の表示仕様は `@testing-library/svelte` を使い、利用者の見える結果で検証する
- Rust の usecase、repository、importer は `cargo test` で回帰を固定する
- desktop 受け入れ検証は `tauri-driver` と `WebdriverIO` を正本候補とし、重要フローを end-to-end で担保する
- lint、`cargo fmt --check`、`clippy -D warnings`、Sonar gate を通らない変更は完了扱いにしない

## 8. 禁止事項

- UI から `@tauri-apps/api/*` を直接呼び、`src/gateway/tauri/` を迂回する実装
- repo の層境界、採用技術、validation を無視して外部テンプレートをそのまま流用する実装
- `src/ui` から `src-tauri/` の構造や Rust 内部型を直接前提にする実装
- `#[tauri::command]` 関数に DB 接続生成、SQL、ファイル操作、業務判定を直書きする実装
- テストを更新せずに仕様だけをコメントや口頭説明で補う変更

## 9. 参照元

- 参考にした外部ガイド: [`Tauri (Svelte, TypeScript Guide)`](https://cursorrules.org/article/tauri-svelte-typescript-guide-cursorrules-prompt-f)
- Tauri 2 の正本:
  [`Permissions`](https://v2.tauri.app/security/permissions/),
  [`Capabilities`](https://v2.tauri.app/security/capabilities/),
  [`Calling Rust from the Frontend`](https://v2.tauri.app/develop/calling-rust/),
  [`State Management`](https://v2.tauri.app/develop/state-management/),
  [`Vite`](https://v2.tauri.app/start/frontend/vite/)
- repo 固有の正本: [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)
