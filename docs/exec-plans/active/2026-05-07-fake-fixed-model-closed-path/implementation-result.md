# Implementation Result

## 変更ファイル

- `internal/infra/ai/provider_models.go`
- `internal/infra/ai/transport.go`
- `internal/infra/ai/provider_client.go`
- `internal/bootstrap/app_controller.go`
- `internal/service/provider_settings_service.go`
- `internal/service/translation_job_setup_service.go`
- `internal/service/master_persona_service.go`
- `internal/service/master_persona_provider_transport.go`

## 実装内容

- test-safe transport の model list 契約で `fake-model` 1 件を返すようにした。
- fake mode の model list loader を test-safe transport に差し替え、外部 HTTP へ出ないようにした。
- provider settings と Job Setup の model list / validation で、test-safe loader の場合は外部 secret を要求しないようにした。
- マスターペルソナの AI 設定読込で、fake mode かつ model 未保存の場合だけ user-facing provider と `fake-model` を返すようにした。
- `ProviderFake` を AIサービス設定 provider 一覧、Job Setup provider 一覧、provider catalog、frontend 分岐へ追加していない。

## 検証結果

- `go test ./internal/infra/ai ./internal/bootstrap ./internal/service ./internal/controller/wails`: 通過。
- `npm --prefix frontend run test -- --run src/application/usecase/master-persona/master-persona.usecase.test.ts src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts`: 通過。
- `python3 scripts/harness/run.py --suite backend-local`: 初回は公開 method コメント不足で失敗。コメント追加後に通過。
- `python3 scripts/harness/run.py --suite frontend-local`: 通過。

## 未実行理由

- なし。

## 残留リスク

- 実画面確認は未実施。今回の承認済み範囲は既存 Wails method と既存 model options 生成経路の backend 契約修正に限定した。
