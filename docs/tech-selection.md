# 技術選定仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`core-beliefs.md`](./core-beliefs.md)

本書は、システム実装のために採用する技術と、その適用対象を定義する。
この repo は `Wails + Go + Svelte` を基盤とする。

## 1. アプリケーション基盤

- デスクトップアプリ基盤は `Wails v2` を採用する
- バックエンド実装言語は `Go` を採用する
- frontend asset の配布と desktop packaging は `wails dev` と `wails build` を標準入口とする
- project config は repo root の `wails.json` を正本とする

## 2. フロントエンド

- UI フレームワークは `Svelte 5` を採用する
- UI 実装言語は `TypeScript` を採用する
- frontend build tool は `Vite` を採用する
- frontend package root は `frontend/` を採用する
- 初期スタイリングは repo-local CSS を採用し、CSS framework は初期必須要件に含めない

## 3. バックエンド基盤

- lifecycle と desktop bridge は `Wails Bind` を採用する
- backend から frontend への push 通知は `Wails runtime events` を採用する
- HTTP client は `Go standard library` の `net/http` を標準とする
- JSON シリアライズは `encoding/json` を標準とする
- XML 出力は `encoding/xml` を第一候補とする
- ログ計測は `log/slog` を採用する

## 4. 永続化

- ローカルデータベースは `SQLite` を採用する
- DB access の抽象は `sqlx` を基準にし、query builder language は導入しない
- schema の変更管理は repo-owned SQL migration を採用する
- migration の適用は app 起動時の専用初期化責務へ集約する
- xEdit 抽出 JSON はファイルシステム上の正本として保持する
- `SQLite` は入力データ、基盤マスター、翻訳ジョブの実行キャッシュとして使う
- DB driver の concrete choice は implementation plan で固定するが、上位層へ driver 固有 API を漏らさない

## 5. DI と品質基盤

- DI は `手動 DI` を採用する
- frontend lint は `ESLint` を採用する
- frontend formatterは`prettier`を採用する
- `Svelte` / `TypeScript` の静的検証は `svelte-check` を採用する
- frontend test runner は `Vitest` を採用する
- backend format は `gofmt` を採用する
- backend test は `go test` を採用する
- backend static check の第一波は `go vet` を採用する
- 過度な quality tooling の追加は避け、必要性が明確になってから増やす

## 6. テスト技術

- UI コンポーネントと画面内インタラクションのテストは `Vitest` を採用する
- Svelte UI のユーザー視点テストは `@testing-library/svelte` と `jsdom` を採用する
- backend の unit test / integration test は `go test ./...` を採用する
- system test 用の結合テスト基盤は `Playwright` を採用する
- `Playwright` は frontend の browser automation、主要フローの executable spec、boot smoke の自動化に使う
- system test の実行対象は、初期段階では `wails dev -browser` または frontend dev server で公開した browser surface を正本とする
- Wails の native window 固有挙動は `Playwright` の primary scope に含めず、必要になるまで manual verification または別途専用手段で補う
- `Playwright` の test runner、fixture、web-first assertion を system test の標準入口とする

## 7. 公式参照

- Wails official docs:
  [`Getting Started`](https://wails.io/docs/gettingstarted/firstproject),
  [`Project Config`](https://wails.io/docs/reference/project-config),
  [`Application Development`](https://wails.io/docs/guides/application-development)
- Svelte official docs:
  [`Svelte 5`](https://svelte.dev/docs/svelte/overview),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript),
  [`Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide)
- Vite official docs:
  [`Guide`](https://vite.dev/guide/),
  [`Build`](https://vite.dev/guide/build),
  [`Config`](https://vite.dev/config/)
- Playwright official docs:
  [`Getting Started`](https://playwright.dev/docs/intro),
  [`Test Runner`](https://playwright.dev/docs/test-intro),
  [`Assertions`](https://playwright.dev/docs/test-assertions)
