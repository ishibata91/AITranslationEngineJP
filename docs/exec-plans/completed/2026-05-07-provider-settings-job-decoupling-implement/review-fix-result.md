# review 修正結果

## 対象

対象 task: `2026-05-07-provider-settings-job-decoupling-implement`

目的: 5 観点レビューの修正必須指摘を解消する。

## 解消内容

- `trust-boundary-001` / `contract-002`
  - Job Setup options response から `credentialRefs` を削除した。
  - Wails DTO、生成済み Wails 型、frontend gateway contract、fake API、presenter、usecase test fixture から `credentialRefs[].credentialRef` の公開を削除した。
  - provider model list の公開 request から `credentialRef` を削除した。
  - UI 表示面は API キー状態分類だけを使い、credential reference 実値と model list token を表示用 state へ出さない形へ揃えた。

- `contract-001`
  - create / validate request と phase runtime selection / summary の公開型から `credentialRef` と `modelListSourceToken` を削除した。
  - Wails DTO 変換だけで落とす形をやめ、Go usecase contract と frontend application gateway contract からも削除した。
  - controller、gateway、usecase、fake API、関連テストの payload 期待値を同じ公開型へ揃えた。

- `trust-boundary-002`
  - provider settings 対象 provider では、provider settings service が nil の場合に phase 開始を拒否するようにした。
  - term / persona / body の start と retry は、provider settings service から解決できた credential reference だけを provider adapter へ渡す。
  - provider settings 非対象 provider の既存経路は維持した。

- `responsibility-boundary-001`
  - `TranslationJobPhaseRuntimeSnapshot` と `TranslationJobPhaseRuntimeSnapshotDraft` から、DB が所有しない `CredentialRef`、`EndpointSummary`、`ModelListSourceToken` を削除した。
  - SQLite adapter、phase service の draft 生成、scenario test を非 secret 要約だけへ揃えた。

- `state-invariant-001`
  - retry の runtime snapshot は、phase が `Running` へ遷移した場合だけ永続化済みとして扱うようにした。
  - term / persona / body で、状態遷移拒否や失敗時に既存 snapshot を維持するようにした。
  - DB snapshot と in-memory `executionSnapshots` が同じ開始試行を指すよう、永続化成功後に in-memory snapshot を更新する形へ揃えた。

## 検証結果

- 成功: `GOCACHE=/tmp/aitranslationenginejp-go-cache go test ./internal/service ./internal/usecase ./internal/controller/wails ./internal/repository ./internal/apitest -run 'ProviderSettings|TranslationJobSetup|TermTranslation|PersonaGeneration|BodyTranslation|SCN_PSJD'`
- 成功: `npm --prefix frontend run test -- src/application src/controller/wails/translation-job-setup.gateway.test.ts src/ui/screens/translation-job-setup/JobSetupPage.test.ts src/controller/review-fake-api`
- 成功: `GOCACHE=/tmp/aitranslationenginejp-go-cache python3 scripts/harness/run.py --suite backend-local`
- 成功: `python3 scripts/harness/run.py --suite frontend-local`

## 重要エラー

- 初回の通常 `go test` は、Go の既定 cache 書き込み権限で失敗した。
- `GOCACHE=/tmp/aitranslationenginejp-go-cache` を指定して再実行し、対象検証は成功した。

## 未解消

- 未解消の review id はない。
