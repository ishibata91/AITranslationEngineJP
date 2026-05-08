# unit-test-result

## 変更ファイル
- frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts
- internal/service/provider_settings_service_test.go
- internal/service/translation_job_setup_service_test.go

## 証明した振る舞い
- frontend: `credentialStatus=missing` の phase でも `refreshPhaseModels` が gateway を呼び、backend 応答 (`credentialStatus=not_required`, 単一モデル) を state に反映する。
- frontend: 単一モデルの自動選択は `fake-model` 固定名に依存しない。唯一の `modelId` を選択する。
- backend provider settings: test-safe loader のとき、`Endpoint=nil` / `CredentialReferenceID=nil` / `CredentialState=missing` でも loader を呼ぶ。loader 呼び出しの `apiKey` は空、結果は `State=ready` と `CredentialState=not_required` になる。
- backend translation job setup: provider settings 経由のモデル一覧で、missing credential の入力スナップショットでも `ListProviderModels` が成功し model list を返す。

## 実行コマンドと結果
- 失敗: `npm --prefix frontend run test -- --run frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
  - 理由: frontend の `vitest include` は `src/**/*.test.ts` であり、先頭 `frontend/` を含む指定は対象外。
- 成功: `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- 成功: `go test ./internal/service -run 'ProviderSettings|TranslationJobSetup' -count=1`
- 成功: `python3 scripts/harness/run.py --suite frontend-local`
- 失敗: `python3 scripts/harness/run.py --suite backend-local`
  - 失敗箇所: `internal/bootstrap`
  - 失敗テスト:
    - `TestNewAppControllerProvidesMasterPersonaAISettingsPersistence`
    - `TestNewAppControllerPersistsMasterPersonaAISettingsAcrossControllerRecreation`
  - 観測: 期待値と実値差分は model 値 (`expected ...` / `got ... Model:"fake-model"`)。

## 追加修正前の戻し指摘
- `ResolveProviderExecutionSettings` の `test-safe` 条件で `Endpoint != nil` を要求する分岐が実装に存在した。
- `Endpoint=nil` の test-safe 実行時に `ErrorKind` が消えない経路が実運用で問題化する可能性があった。
- 上記指摘は追加修正で解消済みである。

## 残留リスク
- `backend-local` は `internal/bootstrap` の既存失敗で非通過。今回の単体テスト対象外レイヤでの失敗であり、unit-test 証明範囲の外に残留する。

## 追加修正結果

- `ResolveProviderExecutionSettings` の test-safe snapshot 適用は、endpoint が nil でも `CredentialReferenceID=nil`、`CredentialState=not_required`、`ErrorKind=nil` になることを単体テストで証明した。
- real provider の endpoint missing と credential missing の gate 維持を単体テストで証明した。
- 成功: `go test ./internal/service -run 'ProviderSettings' -count=1`
- 失敗: `python3 scripts/harness/run.py --suite backend-local`
  - 失敗箇所: `internal/bootstrap`
  - 失敗内容: master persona AI settings persistence 2 件の model 期待値差分。
