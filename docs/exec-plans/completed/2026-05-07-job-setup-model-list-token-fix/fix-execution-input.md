# 修正実行入力

## 根拠

- 人間観測: [human-observation.md](./human-observation.md)
- 原因候補確認: `internal/service/translation_job_setup_service.go` の `listProviderModelsViaProviderSettings`
- 比較経路: `internal/service/master_persona_service.go` の `ListProviderModels`

## 影響ファイル候補

- `internal/service/translation_job_setup_service.go`
- `internal/service/translation_job_setup_service_test.go`

## 禁止変更範囲

- frontend の画面状態、表示部品、Wails binding は変更しない。
- provider 設定の公開契約、secret 保存、credential 解決境界は変更しない。
- docs 正本本文は変更しない。

## 実装 skill

- `implement-backend`

## 修正方針

- Job Setup の画面操作用 request token は、frontend の遅延応答破棄用として response に残す。
- provider 設定側の `ListProviderModels` へ渡す `RequestToken` は、provider 設定 summary の request token にする。
- 回帰テストで、provider 設定側へ渡す token が provider 設定 snapshot token であることを証明する。

## 回帰確認観点

- Job Setup の provider 設定経由モデル一覧取得が `success` になる。
- provider 設定 service の snapshot 照合に画面操作用 token を渡さない。
- response の `RequestToken` は画面操作用 token のまま維持される。
