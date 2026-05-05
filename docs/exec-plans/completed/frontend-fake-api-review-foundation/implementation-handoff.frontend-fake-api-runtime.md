# 実装引き継ぎ入力: frontend-fake-api-runtime

## 状態

- `handoff_id`: `frontend-fake-api-runtime`
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `ready_wave`: `wave-1`
- `depends_on`: なし
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `human_review_status`: 承認済み

## 目的

frontend の composition root に、レビュー起動用 fakeAPI の DI 差し替えを追加する。
実画面は Wails バインディングと backend なしで開けるようにする。

## 承認済み実装範囲

- `frontend/src/main.ts` で本番 gateway と fakeAPI gateway の選択を閉じる。
- fakeAPI 起動中だけ `fakeScenario` を URL パラメータから読む。
- 本番起動相当では `fakeScenario` を無視する。
- `frontend/src/controller/review-fake-api/` など controller 配下に fakeAPI 境界を置く。
- 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態の状態パターン ID を登録できる形にする。
- 後続 task が画面固有モックデータを追加できる登録境界を用意する。

## 範囲外

- backend、本番 API、永続化、本番初期状態を変更しない。
- 生成済み `frontend/wailsjs/` を変更しない。
- レビュー専用 UI、状態パターン選択 UI、表示文言設計を作らない。
- docs 正本、`.codex`、プロダクトテストを変更しない。

## 初手

`frontend/src/main.ts` の gateway 生成を production factory 作成関数へ切り出す。
対応する完了条件は、composition root で本番 gateway と fakeAPI gateway の選択を閉じることである。

理由: fakeAPI DI の入口は frontend bootstrap である。
URL パラメータ解決と fake factory は、同じ入口に依存する。

## 完了条件

- レビュー起動条件が有効な時だけ fake ゲートウェイが選ばれる。
- 本番起動相当では Wails gateway だけが選ばれ、`fakeScenario` は無視される。
- View、ScreenController、Frontend UseCase は生成済み `wailsjs` を直接参照しない。
- fakeAPI とモックデータは本番 API、永続化、本番初期状態へ接続されない。
- 未登録状態パターンと欠落モックデータを成功状態に見せない。
- 画面固有モックデータを後続 task 側で追加できる登録境界がある。

## 検証コマンド

- `npm --prefix frontend run lint:types`
- `npm --prefix frontend run lint:boundaries`

## 禁止事項

- 依存していない別 wave のテスト実装を先取りしない。
- backend と統合境界の実装を混ぜない。
- 状態パターン ID 以外のモックデータ本文や内部診断を URL、log、UI 文言へ出さない。
