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
- `UI Component`: `View` から呼ばれる Svelte component。表示、入力部品、表示規則を画面状態から分離する
- `ScreenController`: 画面操作の入口。screen local な依存を束ね、`UseCase` を起動する
- `Frontend UseCase`: 画面状態の更新手順を決め、`GatewayContract` と `Store` を使う
- `Presenter`: `Store` の状態を view model へ整形する
- `Store`: 画面状態を保持する
- `Gateway`: Wails Bind を呼ぶ frontend adapter
- `RuntimeEventAdapter`: Wails event を購読し、screen local な handler へ流す frontend adapter
- `Backend Bootstrap`: `internal/bootstrap/`。production graph を手動 DI で組み立てる composition root
- `Controller`: backend の入口。Wails Bind の request / response DTO を内部境界へ写像する
- `UseCasePort`: `Controller` から `Backend UseCase` へ向かう controller 依存境界
- `Backend UseCase`: query / command / import を orchestrate し、job 状態と application result を扱う
- `ServicePort`: `Backend UseCase` から `Service` へ向かう usecase 依存境界
- `TranslationJobPolicy`: 翻訳ジョブ操作の共通操作規則と phase 開始前提を評価する UseCase 専用の純粋な規則オブジェクト
- `NotificationSinkPort`: 実行側の複数主体が進捗事実、完了事実、破棄事実を横から渡す通知入口
- `NotificationDispatcher`: `NotificationSinkPort` の実装として、通知事実を Wails 非依存の通知へ整形する
- `NotificationPort`: `Notification` から transport adapter へ渡す送信境界
- `Service`: CRUD / import の実処理を担い、concrete API を直接持たない
- `RepositoryPort`: `Service` から永続化 adapter へ向かう永続化境界
- `XMLFilePort`: `Service` から XML file adapter へ向かう path 解決と file open の境界
- `XMLRecordReaderPort`: `Service` から XML reader adapter へ向かう XML record 読み出し境界
- `Repository`: SQLite などの concrete 永続化を持つ adapter
- `XML adapter: file`: file path と file open の concrete 実装
- `XML adapter: reader`: XML decoder の concrete 実装
- `Runtime adapter`: Wails runtime event の送信だけを扱う transport adapter
- `AIProvider`: provider ごとの差異を吸収する adapter

本書でいう構造図は、この主語同士の依存方向だけを示す。
DB テーブル、DTO 項目、要件フロー、画面遷移は構造図へ混ぜない。

## 2. システム全体の依存方向

全体の依存方向は次の通りとする。

- `frontend/main.ts -> Gateway -> root View`
- `View -> UI Component`
- `View -> ScreenController`
- `UI Component -> View`
- `ScreenController -> Frontend UseCase / Presenter / Store / RuntimeEventAdapter`
- `Frontend UseCase -> GatewayContract / Store`
- `Gateway -> generated wailsjs -> backend Controller`
- `Backend Bootstrap -> Controller`
- `Controller -> UseCasePort`
- `UseCasePort -> Backend UseCase`
- `Backend UseCase -> ServicePort / TranslationJobPolicy / NotificationSinkPort`
- `ServicePort -> Service`
- `Service -> RepositoryPort / XMLFilePort / XMLRecordReaderPort / NotificationSinkPort / AIProvider`
- `NotificationSinkPort -> NotificationDispatcher`
- `NotificationDispatcher -> NotificationPort`
- `RepositoryPort -> Repository`
- `XMLFilePort -> XML adapter: file`
- `XMLRecordReaderPort -> XML adapter: reader`
- `NotificationPort -> Runtime adapter`

`Bootstrap` 以外の層は concrete 実装を new しない。
DI コンテナは使わず、frontend と backend の両方で手動 DI を使う。

## 3. Frontend アーキテクチャ

### 3.1 Frontend Bootstrap

`frontend/src/main.ts` は production 用 gateway を生成し、root view へ注入する。
frontend 全体の composition root はここに置き、DI コンテナは使わない。

### 3.2 View

- 画面を表示する
- `UI Component` を組み合わせる
- DOM event を `ScreenController` へ渡す
- view model だけを前提に描画する

View は backend DTO や generated binding を直接扱わない。

### 3.3 UI Component

UI Component は `View` の下位に置く表示部品である。
UI Component は画面単位の部品と共有部品の二層で扱う。

- 画面専用部品は `frontend/src/ui/screens/<screen>/` に置く
- 複数画面で使う部品だけ `frontend/src/ui/components/` に置く
- 業務フロー全体の進行状態を持たない
- backend DTO、generated binding、`Store`、`Gateway` を直接扱わない
- 状態変更は event として `View` へ返し、`View` から `ScreenController` へ渡す

UI Component は部品化できるものを部品化する。
ただし、画面専用の大きなレイアウトや親画面の状態を大量に読む部品は分けない。

UI Component の部品化判断は次の表に従う。

| 条件 | 見るもの | 部品化しやすい例 | 分けないほうがよい例 |
| --- | --- | --- | --- |
| 意味が独立している | その部品を一言で説明できるか | `UserStatusBadge`、`SearchForm`、`Pagination` | 右上にある灰色の箱 |
| 入力が明確 | props や引数にできるか | `status`、`label`、`onClick` | 親画面の状態を大量に直接読む |
| 出力が明確 | event や表示結果が限定されるか | `onSubmit(query)`、`onSelect(id)` | 内部で複数の画面状態を勝手に更新する |
| 状態を閉じ込められる | 内部状態と外部状態を分けられるか | 開閉状態、入力中テキスト | 業務フロー全体の進行状態 |
| 変更理由がまとまる | 仕様変更時に同じ理由で変わるか | 日付表示規則、状態表示 | A画面では契約都合、B画面では権限都合で変わる |
| 使用箇所が複数ある | 再利用されるか | ボタン、カード、一覧行 | 1画面専用の大きなレイアウト |
| バリエーションが制御可能 | variant で表現できるか | `primary`、`secondary`、`danger` | props が増えすぎて条件分岐の塊になる |
| テスト単位になる | 単体で期待値を書けるか | `status=pending` なら未確認表示 | 画面全体を起動しないと意味がない |
| デザイン規則を担う | 余白、色、文言規則を統一できるか | `FormField`、`ErrorMessage` | 個別画面の例外スタイル |
| ドメイン概念に対応する | 業務上の概念名を持てるか | `LicenseLimitSummary`、`TenantRoleTable` | 単なる `BoxWithIconAndText` |

### 3.4 ScreenController

- screen local な composition root として `UseCase`、`Store`、`Presenter`、`RuntimeEventAdapter` を束ねる
- `UseCase` を起動する
- `Store` の状態を `Presenter` で view model へ変換して View へ返す
- gateway 差し替えや mount / dispose を管理する

`ScreenController` は画面境界の制御を持つが、Wails 呼び出しの詳細や DTO 変換は持たない。

### 3.5 Frontend UseCase

- 画面操作ごとの更新手順を決める
- `GatewayContract` を呼ぶ
- `Store` を更新する
- runtime event 完了時の再読込条件を管理する

`Frontend UseCase` は generated `wailsjs` や backend DTO に直接依存しない。

### 3.6 Presenter と Store

- `Store` は screen state の正本を保持する
- `Presenter` は `Store` の state と接続状態から view model を組み立てる
- View は `Store` を直接加工しない

### 3.7 Gateway と RuntimeEventAdapter

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
- synchronous response を返す

`Controller` は service concrete、repository concrete、`NotificationDispatcher` を直接 new しない。
`Controller` は実行中の途中経過通知の戻り先にならない。

### 4.3 Backend UseCase

- query / command / import を orchestrate する
- job 状態と application result を扱う
- `ServicePort` を使って query / command / import を起動する
- `TranslationJobPolicy` を使って job / phase run 状態の操作可否を判断する
- job / phase run 状態の取得と保存は既存の `ServicePort`、`Service`、`RepositoryPort`、`Repository` 境界を通す
- `TranslationJobPolicy` から操作可否、拒否理由、状態作用、呼び出す service method の種類を得る
- UseCase が確定した状態事実だけを既存の service / repository 経由で保存する
- 必要な通知事実を `NotificationSinkPort` へ渡す
- synchronous response に必要な操作結果を返す
- Wails runtime event payload を組み立てない

`Backend UseCase` は adapter concrete、`NotificationDispatcher`、`NotificationPort` を直接参照しない。

### 4.4 TranslationJobPolicy

`TranslationJobPolicy` は `internal/usecase/translationjobpolicy/` に置く。
`TranslationJobPolicy` は DB を読まず、DB へ保存せず、Service を呼ばない。

`TranslationJobPolicy` は UseCase だけが呼び出す。
`TranslationJobPolicy` は共通操作規則を先に評価し、`start` の時だけ phase 別開始前提を評価する。
`retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。

`PolicyResult` は UseCase 内の一時値である。
`PolicyResult`、適用 rule 名、policy 判定履歴は DB、DTO、repository 永続契約、read model の永続値へ出さない。

### 4.5 Job / phase run 状態事実

job / phase run 状態の取得と保存は、既存の UseCase、Service、Repository 境界で扱う。
専用の状態 IO service を構造主語として置かない。
状態遷移可否、terminal guard、provider response validation、UI 表示文言は永続化 adapter で判断しない。

保存する対象は、UseCase が確定した `TRANSLATION_JOB.state`、`JOB_PHASE_RUN.state`、継続または作成された `JOB_PHASE_RUN` id、進捗、開始時刻、終了時刻、失敗 reason category などの状態事実だけである。
operation summary、provider raw payload、secret、API key、credential 参照実値は保存しない。

### 4.6 Notification

`NotificationSinkPort` は実行側から `Notification` へ入る横接続の入口である。
UseCase、Service、将来の Runner / Worker は、進捗事実、完了事実、破棄事実を `NotificationSinkPort` へ渡せる。
途中経過通知は `Controller` へ戻さない。

`NotificationDispatcher` は `internal/notification/` に置く。
`NotificationDispatcher` は `NotificationSinkPort` を実装する。
`NotificationDispatcher` は通知種別、redaction、送信可否、送信失敗の扱いを決める。

`NotificationDispatcher` は状態遷移可否、terminal guard、provider response validation を判断しない。
`NotificationDispatcher` は operation summary、Wails event payload、通知結果を DB に永続化しない。

`NotificationPort` は `Notification` から transport adapter への境界である。
`Runtime adapter` は `NotificationPort` を実装し、Wails runtime event の実送信だけを扱う。

Service は `NotificationSinkPort` へ進捗事実と完了事実を渡す。
Service は Wails runtime event payload を組み立てず、runtime handle も扱わない。

### 4.7 Service

- CRUD / import の実処理を担う
- concrete API を直接持たない
- 永続化 port を通してデータを読む、書く
- XML file / reader port を通して import を実行する
- `NotificationSinkPort` を通して進捗事実、完了事実、破棄事実を渡す
- AI 実行が必要な機能では `AIProvider` を使う

`Service` core は filesystem、Wails runtime、XML decoder、driver 固有 API を直接参照しない。

### 4.8 Adapter 群

- `Repository` は SQLite などの concrete 永続化を持つ
- `XML adapter: file` は file path と file open の concrete 実装を持つ
- `XML adapter: reader` は XML decoder の concrete 実装を持つ
- `Runtime adapter` は `NotificationPort` を実装し、Wails runtime event の送信だけを持つ
- `AIProvider` は provider ごとの差異を吸収する

adapter concrete は `internal/repository/`、`internal/service/`、`internal/infra/` に閉じ込める。

## 5. 強い制約

- frontend / backend ともに DI コンテナを使わない
- `Bootstrap` 以外の層で concrete 実装を new しない
- `View` は generated `wailsjs` と backend DTO を直接扱わない
- `UI Component` は backend DTO、generated binding、`Store`、`Gateway` を直接扱わない
- `Frontend UseCase` は `GatewayContract` と `Store` だけに依存する
- `Backend Controller` は caller-owned `UseCasePort` だけに依存する
- `Backend Controller` は途中経過通知の経路にならない
- `Backend UseCase` は caller-owned `ServicePort`、`TranslationJobPolicy`、`NotificationSinkPort` に依存する
- `Backend UseCase` は `NotificationDispatcher`、`NotificationPort`、`Runtime adapter` に依存しない
- `TranslationJobPolicy` を呼べる層は `Backend UseCase` だけにする
- policy 判断結果、rule 名、判定履歴は永続化 adapter へ渡さない
- `NotificationDispatcher` は状態遷移可否を判断しない
- `Service` core は concrete driver、runtime API、`NotificationDispatcher` を直接参照しない
- Wails event は push 通知専用に限定し、通常の query / command を置き換えない

## 6. 現在のディレクトリ正本

- `frontend/src/main.ts`: frontend bootstrap
- `frontend/src/ui/`: View、UI Component、screen local な controller / usecase / presenter / store
- `frontend/src/application/`: shared な gateway contract などの frontend 境界定義
- `frontend/src/controller/wails/`: gateway、DTO、generated binding wrapper、runtime 連携 adapter
- `frontend/wailsjs/`: generated bindings。hand-edit しない
- `internal/bootstrap/`: backend bootstrap と default wiring
- `internal/controller/`: backend bind と入出力の受け渡し
- `internal/usecase/`: 操作単位の orchestration
- `internal/usecase/translationjobpolicy/`: 翻訳ジョブ操作の共通操作規則と phase 開始前提
- `internal/notification/`: `NotificationSinkPort`、通知種別、redaction、送信可否、送信失敗の扱い
- `internal/service/`: 実処理と adapter port
- `internal/repository/`: 永続化 adapter
- `internal/aiprovider/`: AI provider 境界
- `internal/infra/`: runtime、HTTP client、filesystem、database driver などの concrete 実装

現在の frontend は screen local な application object を `frontend/src/ui/screens/` に置いている。
shared contract と Wails adapter だけを別 directory に分ける構成を正本とする。

## 7. Wails 境界

- frontend の query / command は `frontend/src/controller/wails/` から generated `wailsjs` を呼ぶ
- backend の bind 公開面は `internal/controller/wails/` とする
- backend から frontend への push は実行側から `NotificationSinkPort` へ入り、`NotificationDispatcher` から `Runtime adapter` 経由で送る
- runtime の concrete handle は bootstrap と adapter に閉じ込める

Wails は transport boundary であり、domain rule や画面状態の正本ではない。
