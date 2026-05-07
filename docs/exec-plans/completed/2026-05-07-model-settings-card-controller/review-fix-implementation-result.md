# レビュー指摘修正 実装結果

## 変更ファイル

- `frontend/src/application/gateway-contract/model-settings-card/*`: 共有カード状態から `requestToken` と `sourceToken` を削除。
- `frontend/src/application/gateway-contract/master-persona/*`: AI 設定読込応答と provider model list 応答を追加。
- `frontend/src/application/store/master-persona/master-persona.store.ts`: backend 由来の provider credential state を保持。
- `frontend/src/application/usecase/master-persona/master-persona.usecase.ts`: model list 更新を専用 gateway 経路へ接続し、遅延応答 token を usecase 内部へ閉じ込め。
- `frontend/src/application/presenter/master-persona/master-persona.presenter.ts`: providerOptions を backend 由来 state から表示へ反映。
- `frontend/src/controller/wails/master-persona.gateway.ts`: `MasterPersonaListProviderModels` binding を追加。
- `frontend/src/controller/wails/gateway-dto/master-persona/*`: master persona 用 DTO 型を追加。
- `internal/controller/wails/master_persona_controller.go`: AI 設定状態応答と model list DTO を追加。
- `internal/usecase/master_persona_usecase.go`: provider state と model list の usecase 境界を追加。
- `internal/service/master_persona_service.go`: provider settings 経由の credential state と model list 取得を追加。
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`: 共有カード呼び出しから token 受け渡しを削除。
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`: fallback 共有カード状態から token を削除。
- `frontend/src/controller/review-fake-api/default-review-fake-api-gateway-registry.ts`: 新しい master persona gateway contract に追従。
- `frontend/src/application/usecase/master-persona/master-persona.usecase.test.ts`: master persona の credential state と model list 更新テストを追従。
- `frontend/src/application/presenter/master-persona/master-persona.presenter.test.ts`: providerOptions 追従。
- `frontend/src/controller/wails/master-persona.gateway.test.ts`: 新 DTO と binding 追従。
- `internal/controller/wails/master_persona_controller_unit_test.go`: 新 DTO 追従。
- `internal/usecase/master_persona_usecase_test.go`: 新 usecase port 追従。
- `internal/bootstrap/app_controller_test.go`: 新 AI settings response 追従。

## 修正内容

- マスターペルソナ側の AI 設定読込は、`aiSettings`、`providerOptions`、`modelList` を Wails DTO から受け取る。
- マスターペルソナ側のモデル一覧更新は、保存済み AI 設定の再読込ではなく `MasterPersonaListProviderModels` を呼ぶ。
- provider credential state は provider 名から frontend で合成せず、backend の provider settings summary から受け取る。
- `requestToken` と `sourceToken` は共有カード state、view model、master persona DTO へ出さない。
- master persona の遅延応答破棄 token は `MasterPersonaUseCase` の private state に限定した。

## 解消したレビュー指摘

- `behavior-001`: master persona の更新ボタンは provider model list 取得経路へ接続済み。
- `contract-001`: credential state は Wails DTO 由来の `providerOptions` と `modelList` から共有カードへ渡る。
- `trust-boundary-001`: 共有カード契約と master persona 画面状態、DTO、view model から内部 token を削除済み。

## 検証結果

- `npm --prefix frontend run check`: 通過。
- `npm --prefix frontend run test -- src/application/usecase/master-persona/master-persona.usecase.test.ts src/application/presenter/master-persona/master-persona.presenter.test.ts src/controller/wails/master-persona.gateway.test.ts`: 3 files / 56 tests 通過。
- `go test ./internal/controller/wails ./internal/usecase ./internal/service ./internal/infra/ai -run 'ProviderSettings|Model|MasterPersona|Fake'`: 通過。
- `python3 scripts/harness/run.py --suite frontend-local`: 通過。
- `python3 scripts/harness/run.py --suite backend-local`: 通過。

## 残留リスク

- 実画面操作は未実行。レーン内の自動検証で contract と回帰を確認した。
- Translation Job Setup の provider model list DTO は既存契約維持のため残した。共有カード state への token 露出は削除済み。
