# フロントエンド fakeAPI 仕様

関連文書: [`index.md`](./index.md), [`architecture.md`](./architecture.md), [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md)

本書は、フロントエンド人間レビュー用の fakeAPI 運用仕様である。
fakeAPI は実画面を Wails binding なしで動かし、画面状態を mock data で確認するために使う。
本番 API、永続化、backend 実装の代替ではない。

## 1. 目的

- 人間レビュー前の frontend 実装を、実画面で確認できるようにする。
- Wails binding や backend 起動に依存せず、frontend の controller、usecase、presenter、view を通す。
- 画面ごとの状態を `fakeScenario` で切り替え、空状態、読み込み中、成功、実行中、失敗、設定不足を確認する。
- `implement-lane` と `ux-refactor-lane` の frontend 実装で、レビュー URL と確認状態を task 成果物へ記録できるようにする。

## 2. 起動方法

開発用の Wails / Vite 起動を使う。

```sh
npm run dev:wails:agent-browser
```

fakeAPI を有効にする URL 例は次の通り。

```text
http://localhost:34115/?fakeApi=1&fakeScenario=success#provider-settings
```

`fakeApi` が存在し、かつ `import.meta.env.DEV` が有効な時だけ fakeAPI が有効になる。
本番 build では URL に `fakeApi` を付けても fakeAPI へ切り替わらない。
`fakeScenario` は fakeAPI が有効な時だけ読む。

## 3. 状態パターン

標準の `fakeScenario` は次に固定する。

| 値 | 意味 |
| --- | --- |
| `empty` | 対象データがない状態 |
| `loading` | 初回読み込みが完了しない状態 |
| `success` | 正常データが表示される状態 |
| `running` | 実行中または検証中の状態 |
| `error` | 失敗を表示する状態 |
| `config-missing` | 設定不足または入力不足を表示する状態 |

`fakeApi=1` または `fakeApi=true` の場合、既定状態は `empty` になる。
`fakeApi=success` のように `fakeApi` へ状態値を直接入れた場合、その値を既定状態として扱う。
`fakeScenario` がある場合は、`fakeApi` の既定状態より優先する。

不明な状態値は `config-missing` に寄せる。
理由は、誤ったレビュー URL を開いた時に、設定不足として見える方が原因を追いやすいためである。

## 4. アーキテクチャ

fakeAPI は composition root で gateway を差し替える。
View や usecase は本番と同じ経路を通る。

```text
frontend/src/main.ts
  -> review-fake-api-runtime.ts
  -> app-screen-controller-factories.ts
  -> ReviewFakeApiGatewayRegistry
  -> 画面別 GatewayContract
  -> usecase / presenter / view
```

本番起動では `createProductionAppFactories()` が Wails gateway を作る。
fakeAPI 起動では `createReviewFakeApiAppFactories()` が fake gateway を注入する。
usecase は `*GatewayContract` だけを参照し、Wails gateway か fake gateway かを知らない。

## 5. 追加方法

画面別 fakeAPI を追加する時は、次の順で実装する。

1. 対象画面の `*GatewayContract` を確認する。
2. `frontend/src/controller/review-fake-api/` 配下に、画面別 fake gateway または mock data を置く。
3. fake gateway は対象の `*GatewayContract` を満たす。
4. `ReviewFakeApiGatewayRegistry` に画面別 gateway factory を追加する。
5. `createDefaultReviewFakeApiGatewayRegistry()` から追加した factory を返す。
6. review URL、対象 route、確認した `fakeScenario` を task 成果物へ記録する。

画面固有の状態は、標準の 6 状態へ寄せる。
標準状態で表現できない場合だけ、実装範囲または task 枠で状態追加の理由を固定する。

## 6. 禁止事項

- Wails binding や generated `wailsjs` を直接 mock しない。
- View、ScreenController、Frontend UseCase から fake gateway を生成しない。
- fakeAPI を本番起動、本番初期状態、永続化、backend DTO に混入させない。
- fakeAPI のためだけにレビュー専用 UI や状態選択 UI を追加しない。
- API key、token、ローカル絶対パス、secret を mock data や画面表示へ入れない。
- fakeAPI を backend 実装、統合境界実装、永続化仕様の代替にしない。

## 7. 検証

fakeAPI を追加または変更した時は、次を確認する。

```sh
npm --prefix frontend run test -- src/controller/review-fake-api
npm --prefix frontend run test -- src/ui
python3 scripts/harness/run.py --suite frontend-local
```

実画面確認では `agent-browser` を使い、最低 1 つの成功状態と 1 つの失敗状態を開く。

```sh
agent-browser open "http://localhost:34115/?fakeApi=1&fakeScenario=success#provider-settings"
agent-browser open "http://localhost:34115/?fakeApi=1&fakeScenario=error#provider-settings"
agent-browser snapshot
agent-browser errors
```

task 成果物には、実行 URL、確認した状態、未確認状態、未確認理由を書く。
人間レビュー前の frontend 実装では、fakeAPI 確認結果を frontend 実装結果または人間レビュー入力へ含める。
