# アーキテクチャ仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`core-beliefs.md`](./core-beliefs.md), [`tech-selection.md`](./tech-selection.md)

本書は、システムの内部境界、依存方向、手動 DI の正本を定義する。
本書では backend と frontend のレイヤー関係だけを扱い、要件フローや画面仕様は扱わない。

関連図:

- [`diagrams/backend/backend-architecture.puml`](./diagrams/backend/backend-architecture.puml)
- [`diagrams/frontend/frontend-architecture.puml`](./diagrams/frontend/frontend-architecture.puml)

## 1. 構造の主語

本 repo の構造主語は次の通りとする。

- `Frontend Bootstrap`: `frontend/src/main.ts`。gateway を生成し、root view へ注入する frontend 側の手動 DI 入口
- `View`: Svelte component。表示と DOM event を扱う
- `ScreenController`: 画面操作の入口。screen local な依存を束ね、`UseCase` を起動する
- `Frontend UseCase`: 画面状態の更新手順を決め、`GatewayContract` と `Store` を使う
- `Presenter`: `Store` の状態を view model へ整形する
- `Store`: 画面状態を保持する
- `Gateway`: Wails Bind を呼ぶ frontend adapter
- `RuntimeEventAdapter`: Wails event を購読し、screen local な handler へ流す frontend adapter
- `Backend Bootstrap`: `internal/bootstrap/`。production graph を手動 DI で組み立てる composition root
- `Controller`: backend の入口。Wails Bind の request / response DTO を内部境界へ写像する
- `Backend UseCase`: 操作単位の orchestration を担う
- `Service`: 実処理を担う
- `StateMachine`: 状態遷移規則だけを保持する
- `JobIOService`: job 状態の取得と保存だけを扱う
- `Repository` / `XML adapter` / `Runtime adapter` / `AIProvider`: backend の adapter 群

本書でいう構造図は、この主語同士の依存方向だけを示す。
DB テーブル、DTO 項目、要件フロー、画面遷移は構造図へ混ぜない。

## 2. システム全体の依存方向

全体の依存方向は次の通りとする。

- `frontend/main.ts -> Gateway -> root View`
- `View -> ScreenController`
- `ScreenController -> Frontend UseCase / Presenter / Store / RuntimeEventAdapter`
- `Frontend UseCase -> GatewayContract / Store`
- `Gateway -> generated wailsjs -> backend Controller`
- `Backend Bootstrap -> Controller / UseCase / Service / adapter concrete`
- `Controller -> UseCasePort`
- `Backend UseCase -> ServicePort / StateMachine / JobIOService / RuntimeEventPublisherPort`
- `Service -> RepositoryPort / XMLFilePort / XMLRecordReaderPort / RuntimeContextPort / AIProvider`

`Bootstrap` 以外の層は concrete 実装を new しない。
DI コンテナは使わず、frontend と backend の両方で手動 DI を使う。

## 3. Frontend アーキテクチャ

### 3.1 Frontend Bootstrap

`frontend/src/main.ts` は production 用 gateway を生成し、root view へ注入する。
frontend 全体の composition root はここに置き、DI コンテナは使わない。

### 3.2 View

- 画面を表示する
- DOM event を `ScreenController` へ渡す
- view model だけを前提に描画する

View は backend DTO や generated binding を直接扱わない。

### 3.3 ScreenController

- screen local な composition root として `UseCase`、`Store`、`Presenter`、`RuntimeEventAdapter` を束ねる
- `UseCase` を起動する
- `Store` の状態を `Presenter` で view model へ変換して View へ返す
- gateway 差し替えや mount / dispose を管理する

`ScreenController` は画面境界の制御を持つが、Wails 呼び出しの詳細や DTO 変換は持たない。

### 3.4 Frontend UseCase

- 画面操作ごとの更新手順を決める
- `GatewayContract` を呼ぶ
- `Store` を更新する
- runtime event 完了時の再読込条件を管理する

`Frontend UseCase` は generated `wailsjs` や backend DTO に直接依存しない。

### 3.5 Presenter と Store

- `Store` は screen state の正本を保持する
- `Presenter` は `Store` の state と接続状態から view model を組み立てる
- View は `Store` を直接加工しない

### 3.6 Gateway と RuntimeEventAdapter

- `Gateway` は `GatewayContract` を実装する
- `Gateway` は `GatewayDTO` と generated `wailsjs` を `frontend/src/controller/wails/` に閉じ込める
- `RuntimeEventAdapter` は Wails runtime event を購読し、screen local handler へ写像する
- query / command の主経路は Bind call とし、event は push 通知専用に使う

## 4. Backend アーキテクチャ

### 4.1 Backend Bootstrap

`internal/bootstrap/` は backend の唯一の composition root とする。
`internal/bootstrap/` だけが concrete 実装を生成し、手動 DI で依存グラフを接続する。

### 4.2 Controller

- Wails Bind の入口になる
- request / response DTO を usecase 境界へ写像する
- caller-owned の `UseCasePort` を起動する
- runtime context を受け取り、必要な emitter state へ橋渡しする

`Controller` は service concrete や repository concrete を直接 new しない。

### 4.3 Backend UseCase

- 操作単位の orchestration を担う
- `ServicePort` を使って query / command / import を起動する
- `StateMachine` と `JobIOService` を使って job 状態を扱う
- runtime event 完了 payload を組み立てる

`Backend UseCase` は adapter concrete を直接参照しない。

### 4.4 Service

- 永続化 port を通して master data を読む、書く
- XML file / reader port を通して import を実行する
- runtime port を通して進捗を通知する
- AI 実行が必要な機能では `AIProvider` を使う

`Service` core は filesystem、Wails runtime、XML decoder、driver 固有 API を直接参照しない。

### 4.5 Adapter 群

- `Repository` は SQLite などの永続化実装を持つ
- `XML adapter` は path 解決、file open、record 読み出しを持つ
- `Runtime adapter` は Wails runtime event の具体送信を持つ
- `AIProvider` は provider ごとの差異を吸収する

adapter concrete は `internal/repository/`、`internal/service/`、`internal/infra/` に閉じ込める。

## 5. 強い制約

- frontend / backend ともに DI コンテナを使わない
- `Bootstrap` 以外の層で concrete 実装を new しない
- `View` は generated `wailsjs` と backend DTO を直接扱わない
- `Frontend UseCase` は `GatewayContract` と `Store` だけに依存する
- `Backend Controller` は caller-owned `UseCasePort` だけに依存する
- `Backend UseCase` は caller-owned `ServicePort` と純粋な rule object に依存する
- `Service` core は concrete driver や runtime API を直接参照しない
- Wails event は push 通知専用に限定し、通常の query / command を置き換えない

## 6. 現在のディレクトリ正本

- `frontend/src/main.ts`: frontend bootstrap
- `frontend/src/ui/`: View と screen local な controller / usecase / presenter / store
- `frontend/src/application/`: shared な gateway contract などの frontend 境界定義
- `frontend/src/controller/wails/`: gateway、DTO、generated binding wrapper、runtime 連携 adapter
- `frontend/wailsjs/`: generated bindings。hand-edit しない
- `internal/bootstrap/`: backend bootstrap と default wiring
- `internal/controller/`: backend bind と入出力の受け渡し
- `internal/usecase/`: 操作単位の orchestration
- `internal/service/`: 実処理と adapter port
- `internal/statemachine/`: 状態遷移規則
- `internal/jobio/`: job 状態の取得と保存
- `internal/repository/`: 永続化 adapter
- `internal/aiprovider/`: AI provider 境界
- `internal/infra/`: runtime、HTTP client、filesystem、database driver などの concrete 実装

現在の frontend は screen local な application object を `frontend/src/ui/screens/` に置いている。
shared contract と Wails adapter だけを別 directory に分ける構成を正本とする。

## 7. Wails 境界

- frontend の query / command は `frontend/src/controller/wails/` から generated `wailsjs` を呼ぶ
- backend の bind 公開面は `internal/controller/wails/` とする
- backend から frontend への push は runtime event adapter 経由で送る
- runtime の concrete handle は bootstrap と adapter に閉じ込める

Wails は transport boundary であり、domain rule や画面状態の正本ではない。
