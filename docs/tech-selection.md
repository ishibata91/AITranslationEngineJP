# 技術選定仕様

関連文書: [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`interface-spec.md`](./external-design/interface-spec.md), [`execution-spec.md`](./external-design/execution-spec.md), [`ui-spec.md`](./external-design/ui-spec.md)

本書は、システム実装のために採用する技術と、その適用対象を定義する。

## 1. アプリケーション基盤

- デスクトップアプリ基盤は `Tauri 2` を採用する
- アプリケーションコア実装言語は `Rust` を採用する
- Tauri 実行基盤として `WebView2 Fixed Version Runtime` を同梱する

## 2. フロントエンド

- UI フレームワークは `Svelte 5` を採用する
- UI 実装言語は `TypeScript` を採用する
- ビルドツールは `Vite` を採用する
- スタイリングは `Tailwind CSS + DaisyUI` を採用する

## 3. バックエンド基盤

- 非同期ランタイムは `tokio` を採用する
- HTTP クライアントは `reqwest` を採用する
- JSON シリアライズは `serde` / `serde_json` を採用する
- XML 出力は `quick-xml` を採用する
- ログ計測は `tracing` を採用する

## 4. 永続化

- ローカルデータベースは `SQLite` を採用する
- DB アクセスは `sqlx` を採用する
- xEdit 抽出 JSON はファイルシステム上の正本として保持する
- `SQLite` は `PLUGIN_EXPORT` 配下の入力データを実行キャッシュとして保持する
- `MASTER_PERSONA` と `MASTER_DICTIONARY` は `SQLite` 上の永続基盤データとして保持する
- DB の内部主キーはシーケンシャル整数を採用し、外部 FormID は別列で保持する

## 5. DI と品質基盤

- DI は `手動 DI` を採用する
- フロントエンドの標準 lint は `Oxlint` を採用する
- UI / `Svelte` 変更時の補助 lint として `ESLint` を採用する
- Rust の品質基盤は `rustfmt` / `clippy` を採用する
