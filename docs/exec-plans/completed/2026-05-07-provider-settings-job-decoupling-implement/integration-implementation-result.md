# 統合境界実装結果

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-INT-01`
- `implementation_artifact`: `統合境界実装`
- `status`: `completed`

## 実装結果

- create / validate / summary の Wails DTO と frontend gateway DTO から、`credentialRef` と `modelListSourceToken` を送受信しない形へ縮めた。
- Wails controller は create / validate / summary で provider、model、execution mode、batch mode、credential 状態分類だけを usecase へ渡す。
- frontend Wails gateway は application 側の既存 payload から、Wails binding へ渡す前に禁止 field を落とす。
- generated binding は `wails dev` により確認した。対象 DTO では `credentialRef` と `modelListSourceToken` が出力されていない。

## 変更ファイル

- `internal/usecase/translation_job_setup_contract.go`
- `internal/usecase/translation_job_setup_usecase.go`
- `internal/controller/wails/translation_job_setup_controller.go`
- `frontend/src/application/gateway-contract/translation-job-setup/index.ts`
- `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/index.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/translation-job-setup-gateway-dto.ts`
- `frontend/src/controller/wails/translation-job-setup.gateway.ts`
- `frontend/src/controller/wails/translation-job-setup.gateway.test.ts`

## 検証結果

- `go test ./internal/controller/wails ./internal/usecase -run 'TranslationJobSetup|ProviderSettings'`: pass。
- `npm --prefix frontend run test -- src/controller/wails/translation-job-setup.gateway.test.ts src/controller/wails/gateway-dto/translation-job-setup src/application/gateway-contract/translation-job-setup`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。

## 実画面確認

- 起動: `npm run dev:wails:agent-browser`
- 成功状態: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- 確認結果: Job Setup で 3 phase が `設定済み`、モデル `fake-model`、作成前確認 `不足はありません。` を表示した。
- 不足状態: `http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#translation-management`
- 確認結果: Job Setup で 3 phase が `APIキー未設定`、モデル一覧更新不可、作成操作 disabled を表示した。
- `agent-browser errors`: 出力なし。

## 残留事項

- Job Setup options と provider model list 接点には既存の credential reference / source token 系 field が残る。
- 対象理由: `PSJD-INT-01` の必須対象は create / validate / summary の公開接点であり、options と provider model list の契約再設計は今回の完了条件外である。
- 許可範囲外で必要になった残件はない。
