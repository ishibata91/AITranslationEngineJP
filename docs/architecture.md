# アーキテクチャ仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`core-beliefs.md`](./core-beliefs.md), [`tech-selection.md`](./tech-selection.md)

本書は、システムの内部境界、依存方向、手動 DI の一般原則を定義する。
本書では `Wails + Go + Svelte` 前提の一般メタ原則だけを扱い、具体的な feature 構造、画面遷移、DTO 項目、DB schema は扱わない。

具体構造は backend greenfield reset 後に追加する。本書は再構築の骨格として残す。

## 1. 全体方針

- frontend と backend の両方で手動 DI を使う。
- DI コンテナを使わない。
- `Bootstrap` 以外の層で concrete 実装を new しない。
- Wails は transport boundary であり、domain rule や画面状態の正本ではない。

## 2. 構造の主語（骨格）

層名は再構築時に確定する。次は典型的な分類例として残す。

### Frontend 側

- `Frontend Bootstrap`: gateway を生成し、root view へ注入する frontend 側の手動 DI 入口
- `View`: Svelte component。表示と DOM event を扱う
- `UI Component`: `View` から呼ばれる Svelte component。表示、入力部品、表示規則を画面状態から分離する
- `ScreenController`: 画面操作の入口。screen local な依存を束ね、`Frontend UseCase` を起動する
- `Frontend UseCase`: 画面状態の更新手順を決め、`GatewayContract` と `Store` を使う
- `Presenter`: `Store` の状態を view model へ整形する
- `Store`: 画面状態を保持する
- `Gateway`: Wails Bind を呼ぶ frontend adapter
- `RuntimeEventAdapter`: Wails event を購読し、screen local な handler へ流す frontend adapter

### Backend 側

- `Backend Bootstrap`: production graph を手動 DI で組み立てる composition root
- `Controller`: backend の入口。Wails Bind の request / response DTO を内部境界へ写像する
- `Backend UseCase`: query / command を orchestrate する
- `Service`: 実処理を担い、concrete API を直接持たない
- `Repository`: 永続化 adapter
- `Adapter`: file、HTTP、driver などの concrete 実装

具体的な port 名、policy、notification 構造は backend greenfield reset 後に追加する。

## 3. 全体の依存方向（骨格）

- `Bootstrap` だけが concrete 実装を new する。
- 上位層は下位層を、`Bootstrap` で wire された interface 経由で参照する。
- adapter concrete を上位層から直接参照しない。
- frontend と backend は Wails 境界で接続する。

## 4. Wails 境界

- frontend の query / command は frontend の gateway から generated `wailsjs` を呼ぶ。
- backend の Bind 公開面は backend の controller 層が担当する。
- backend から frontend への push は実行側から transport adapter 経由で送る。
- runtime の concrete handle は bootstrap と adapter に閉じ込める。
- query / command の主経路は Bind call とし、event は push 通知専用に使う。

## 5. ディレクトリ正本（骨格）

具体的な path は backend greenfield reset 後に確定する。次は方向性として残す。

- frontend: `frontend/src/` 配下に View、UI Component、screen local な controller / usecase / presenter / store
- frontend: `frontend/src/controller/wails/` に gateway、DTO、generated binding wrapper
- frontend: `frontend/wailsjs/` に generated bindings（hand-edit しない）
- backend: `internal/bootstrap/` に composition root（greenfield 後に再構築）
- backend: `internal/` 配下に層別 package（greenfield 後に再構築）

## 6. 現在の状態（2026-06-06）

backend は `greenfield-reset` task で削減済み。残存:

- `internal/repository/provider_settings_keyring_secret_store.go`（汎用 keyring wrapper、再利用素材）
- `internal/repository/master_persona_keyring_secret_store.go`（同上、`keyringOpenFunc` type 定義を含む）

新 architecture は次 task で議論しながら確定し、本書に書き加える。
