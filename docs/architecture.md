# アーキテクチャ仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`core-beliefs.md`](./core-beliefs.md), [`tech-selection.md`](./tech-selection.md)

本書は、システムの内部構成と責務分割を定義する。
本書では要件ドメインやデータ項目の詳細には入らず、backend と frontend の構造だけを扱う。

関連図:

- [`diagrams/backend/backend-architecture.d2`](./diagrams/backend/backend-architecture.d2)
- [`diagrams/frontend/frontend-architecture.d2`](./diagrams/frontend/frontend-architecture.d2)

## 1. 構造の主語

本 repo の backend は、次の主語で構成する。

- `Controller`: request / response を受け、`UseCase` を起動する入口
- `UseCase`: ジョブ状態の管理と、次に何をやるかの決定を担う
- `Service`: 実処理を担う
- `StateMachine`: 状態遷移規則だけを保持する
- `JobIOService`: job 状態の取得と保存だけを担当する
- `Repository`: 永続化責務を持つ
- `AIProvider`: AI 実行責務を持つ

本書でいう構造図は、この主語同士の依存と呼び出しだけを示す。
DB テーブル、翻訳対象の概念、要件フロー、画面遷移は構造図へ混ぜない。

## 2. 基本構造

backend の基本構造は次の通りとする。

- `Controller -> UseCase`
- `UseCase -> Service`
- `UseCase -> StateMachine`
- `UseCase -> JobIOService`
- `Service -> Repository`
- `Service -> AIProvider`

この構造では、`UseCase` を薄いオーケストレータとして扱う。
`UseCase` は「何をするか」「今それをしてよいか」に集中し、処理の中身を `Service` へ渡す。

## 3. 各責務

### 3.1 Controller

- request / response を受ける
- Wails の request / response DTO を内部入力へ写像する
- `UseCase` を起動する
- 業務判断を持たない

### 3.2 UseCase

- 操作単位を表す
- ジョブ状態を確認する
- 実行可否を判断する
- 次に呼ぶ `Service` を決める
- `JobIOService` を使って job 状態を取得し、保存する
- `StateMachine` を使って状態更新を確定する
- transaction 境界を持つ

`UseCase` は repository 読み取り、AI 呼び出し、コンテキスト構築の詳細を直接持たない。

### 3.3 Service

`Service` は再利用可能な実処理を担う。

- repository から必要データを読む
- 実行用コンテキストを構築する
- AI 実行を依頼する
- 結果を整形して返す

`Service` は広い util 置き場にしない。
責務名で読める単位に分け、`ContextService`、`TranslationService`、`PersonaGenerationService` のように役割を明確にする。

### 3.4 StateMachine

- 状態遷移規則だけを持つ
- 遷移可否を判定する
- 遷移結果を返す

`StateMachine` は I/O を持たない。

### 3.5 Repository

- 永続化を担当する
- `SQLite` 前提の具象依存を許容する
- schema 準備や migration 実行を通常 use case に混ぜない

現時点では repository 全体へ DIP を強制しない。
差し替え需要が明確になるまでは具象のままでよい。

### 3.6 JobIOService

- job 状態の取得と保存だけを担当する
- `StateMachine` の純粋性を担保するためのトレードオフとして導入する
- 状態遷移規則、業務判断、AI 実行は持たない

`JobIOService` は job に関する I/O だけを扱い、汎用 service へ広げない。

### 3.7 AIProvider

- AI 実行の interface を定義する
- `LMStudio`、`Gemini`、`xAI` などの実装差異を吸収する
- 単発実行と batch 実行の差異を吸収する

AI 基盤接続だけは複数実装が前提なので、`AIProvider` 境界で DIP を適用する。

## 4. 依存方針

- `Controller` は `UseCase` に依存する
- `UseCase` は `Service`、`StateMachine`、`JobIOService` に依存する
- `Service` は `Repository` と `AIProvider` に依存する
- `JobIOService` は `Repository` に依存する
- `Repository` は driver や filesystem に依存してよい
- `AIProvider` 実装は HTTP client や SDK に依存してよい

依存の強い制約は次の通りとする。

- `UseCase` から `AIProvider` 実装を直接参照しない
- `UseCase` から `Repository` の細かい問い合わせ手順を直接広げない
- `Controller` に状態遷移や処理順序を持ち込まない
- `Service` を汎用 helper 集合にしない
- `JobIOService` に状態遷移や AI 実行を持ち込まない

## 5. Wails 境界

Wails の `Bind` は frontend から backend への request / response boundary とする。
この境界は出口ではなく入口なので、`Controller` と呼ぶ。

- request / response: `frontend/src/controller/wails/` から generated `wailsjs` を呼ぶ
- backend bind: `internal/controller/wails/` が public method を公開する
- push 通知: backend は `runtime.EventsEmit` で frontend へ進捗や通知を送る

Wails event は progress、notification、background completion のような push 用に限定し、通常の query / command の主経路には使わない。

## 6. 初期レイアウト

`Wails + Go + Svelte` の初期レイアウトは以下を正本とする。

- repo root
  - `main.go`: Wails bootstrap と app 起動
  - `wails.json`: Wails project config
  - `frontend/`: frontend package root
  - `internal/`: Go の backend 実装
- `frontend/`
  - `src/ui/`: App Shell、screen、view、store
  - `src/application/`: frontend 側 use case
  - `src/controller/wails/`: generated binding wrapper と runtime event adapter
  - `src/shared/contracts/`: UI が依存する DTO / query model
  - `wailsjs/`: generated bindings。hand-edit しない
- `internal/`
  - `controller/`: Wails bind と入出力の受け渡し
  - `usecase/`: 操作単位とジョブ状態管理
  - `service/`: 実処理
  - `statemachine/`: 状態遷移規則
  - `jobio/`: job 状態の取得と保存
  - `repository/`: 永続化
  - `aiprovider/`: AI provider interface
  - `infra/ai/`: AI provider 実装
  - `infra/runtime/`: driver、HTTP client、filesystem、SQLite 接続

directory 名は責務を優先して解釈する。

## 7. DTO 境界

frontend と backend のデータ受け渡しは DTO を明示して行う。

- Go 側の public bind method は request / response struct を明示する
- DTO は `json` tag を付け、field 名を暗黙変換に任せない
- frontend は generated `wailsjs` の型を `controller/wails` の中で `src/shared/contracts/` へ写像する
- UI は shared contract だけを前提にし、Go 内部構造を直接前提にしない

## 8. 永続化と AI 接続

- 入力データの raw JSON はファイルシステム上の正本とする
- `SQLite` は入力キャッシュ、基盤マスター、翻訳ジョブの実行状態を保持する
- schema 変更は repo-owned SQL migration で管理し、起動時 bootstrap で一度だけ適用する
- repository は DML と transaction に専念し、DDL 実行や schema 準備を通常 use case へ混ぜない
- AI provider は provider ごとの差異を adapter 側へ閉じ込める
- provider の選択は use case / service 側で決め、接続詳細は provider 実装側へ閉じ込める
