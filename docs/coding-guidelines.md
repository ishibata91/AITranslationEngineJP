# コーディング規約

関連文書: [`index.md`](./index.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)

本書は、`AITranslationEngineJp` の実装時に守るべきコーディング規約を定義する。
外部の一般論をそのまま持ち込まず、本 repo の `Wails`、`Go`、`Svelte 5`、`TypeScript`、4 層構成に合わせて必要な規約だけを残す。

## 1. 基本原則

- 実装は `UI層`、`アプリケーション層`、`ドメイン層`、`インフラ層` の境界を壊さない
- フロントエンドは `frontend/src/ui`、`frontend/src/application`、`frontend/src/gateway`、`frontend/src/shared` の責務分離を維持する
- バックエンドは `internal/application`、`internal/domain`、`internal/infra`、`internal/gateway/wails` の責務分離を維持する
- コメントは `何をしているか` ではなく、`なぜその判断が必要か` を短く補足する
- 命名は略語より意味を優先し、UseCase、Gateway、Port、Repository の役割が読める名前にする
- generated file は hand-edit せず、source of truth から再生成する

## 2. Wails 固有ルール

- `main.go` は Wails bootstrap と wiring に集中させる
- `Bind` する public method は transport boundary として扱い、業務ロジックや DB 操作を直書きしない
- bound struct の責務は request validation、DTO mapping、use case 呼び出しまでに留める
- frontend から backend を呼ぶ入口は generated `wailsjs` に限定し、UI から直接呼ばず `frontend/src/gateway/wails/` に閉じ込める
- `wailsjs` は generated output として扱い、`frontend/src/gateway/wails/` の外から直接 import しない
- backend から frontend への通知は `runtime.EventsEmit` を使ってよいが、通常の query / command を event へ逃がさない
- lifecycle hook (`OnStartup`, `OnShutdown`) は薄く保ち、重い初期化や domain 判断を抱え込まない

## 3. TypeScript / Svelte 規約

- `Svelte 5` と `TypeScript` を前提にし、`any` と無検証の type assertion を常用しない
- component state は `\$state`、派生値は `\$derived`、副作用は `\$effect` を基本にする
- component event は callback prop を優先し、`createEventDispatcher` の新規採用は避ける
- event handler は `onclick` などの標準 event 属性を優先する
- `.svelte` は表示とイベント配線に集中させ、画面フローや取得判断は use case へ寄せる
- `Presenter / View` は pure に保ち、Wails runtime、永続化、外部通信へ直接依存しない
- `Screen Store` は表示状態の保持に専念し、副作用の起点にしない
- `frontend/src/shared/contracts/` の型を UI と Gateway の共有 DTO とし、Go 側の内部型や generated type を直接前提にしない
- 別 feature の内部 module を直接 import せず、公開 root を経由する

## 4. Go / Backend 規約

- `domain` は最も内側の方針として保ち、Wails、SQLite driver、ファイル I/O、HTTP 実装へ依存しない
- `application` は UseCase と DTO を持ち、外部依存は port 越しに扱う
- `infra` は DB、ファイル、runtime、外部 API の具体実装だけを持ち、上位層へインフラ都合の型を漏らさない
- `internal/gateway/wails` は Wails binding と application usecase の接続点に限定し、ドメイン判断を持たない
- `error` を返して失敗を明示し、`panic` を通常フローの制御に使わない
- `context.Context` が必要な runtime 操作は gateway と infra で明示的に扱う
- 永続化は repository に閉じ込め、use case や bound method に SQL を直書きしない
- migration 実行は backend 起動初期化で一度だけ行い、通常 request 経路へ混ぜない

## 5. 通信とデータ境界

- frontend と backend の主通信路は `Wails bindings` とする
- backend からの push 通知が必要な場合だけ `Wails runtime events` を使う
- request / response は DTO と validation で形を固定し、optional field の意味を曖昧にしない
- versioning が必要な外部連携は Go 側 adapter に閉じ込め、UI 契約へ直接漏らさない
- user-facing message と internal diagnostic は分け、内部詳細をそのまま表示しない

## 6. セキュリティ規約

- file path、外部 URL、provider 設定値は backend 側で再度検証する
- frontend 入力は UI で整形しても、最終 validation は Go 側で再実行する
- 機密値、API key、ローカル絶対パスをログや UI へ無加工で出さない
- filesystem、shell、外部プロセス、HTTP 通信は必要最小限の責務に閉じ込める
- Wails の public bind surface は広げすぎず、必要な method だけを公開する

## 7. テストと品質ゲート

- 振る舞いを変える変更では、対応する `Vitest`、`go test`、acceptance checks、validation commands を同じ変更で更新する
- UI の画面フローは use case と store の unit test を優先し、gateway は mock する
- `.svelte` の表示仕様は `@testing-library/svelte` を使い、利用者の見える結果で検証する
- backend の use case、repository、provider adapter は `go test ./...` で回帰を固定する
- `eslint`、`svelte-check`、`gofmt`、`go vet`、`go test` を通らない変更は完了扱いにしない

## 8. 禁止事項

- UI から generated `wailsjs` や runtime API を直接呼び、`frontend/src/gateway/wails/` を迂回する実装
- `main.go` や bound method に DB 接続生成、SQL、ファイル操作、業務判定を直書きする実装
- `frontend/src/ui` から Go 内部型や backend directory 構造を直接前提にする実装
- テストを更新せずに仕様だけをコメントや口頭説明で補う変更
- repo の層境界、採用技術、validation を無視して外部テンプレートをそのまま流用する実装

## 9. 参照元

- Wails official docs:
  [`Application Development`](https://wails.io/docs/guides/application-development),
  [`How Does It Work`](https://wails.io/docs/howdoesitwork),
  [`Project Config`](https://wails.io/docs/reference/project-config)
- Svelte official docs:
  [`Svelte 5 Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript)
- repo 固有の正本: [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)
