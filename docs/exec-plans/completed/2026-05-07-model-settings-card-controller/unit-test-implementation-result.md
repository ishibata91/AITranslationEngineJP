# unit-test-implementation-result

## 判断結果
単体テストを実装した。
対象は frontend usecase 単体テストのみで、プロダクトコードは未変更である。

## 根拠参照
- 単一引き継ぎ入力: `docs/exec-plans/active/2026-05-07-model-settings-card-controller/unit-test-implementation-input.md`
- 実装済み対象: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- 変更ファイル: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`

## 変更ファイル
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- `docs/exec-plans/active/2026-05-07-model-settings-card-controller/unit-test-implementation-result.md`

## 証明した公開振る舞い
- provider 切替後、`word_translation` の `model` が空へ戻る。
- provider 切替後、`providerModelLists` が切替先 provider の空一覧へ初期化される。
- 上記状態は `modelSettingsCards` に同期され、旧 provider の model list / model が混入しない。
- 不足状態のまま `createJob` を呼ぶと保存拒否され、gateway `createTranslationJob` は呼ばれない。

## 証明した分岐
- 遅延 success 応答が旧 provider 由来なら破棄され、最新 provider 側 state を維持する分岐。
- 遅延 failed 応答が旧 provider 由来なら破棄され、最新 provider 側 state を維持する分岐。

## 証明したエラー経路
- provider 切替後に model 未選択状態が残る場合、`createJob` は実行拒否され `errorMessage` を返す経路。

## 検証結果
- `npm --prefix frontend run test -- --run 'model|provider|master-persona|translation-job-setup'`
  - 結果: 失敗（Vitest の `--run` フィルタで対象 0 件）
- 代替: `npm --prefix frontend run test -- src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts src/application/usecase/provider-settings/provider-settings.usecase.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts`
  - 結果: 成功（3 files / 48 tests passed）
- `go test ./internal/usecase ./internal/service ./internal/repository ./internal/infra/ai -run 'Model|ProviderSettings|MasterPersona|TranslationJobSetup|Redaction'`
  - 結果: 成功
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: 初回 lint 失敗を修正後に再実行し成功

## 網羅率検証結果
- `python3 scripts/harness/run.py --suite coverage`
  - 結果: 成功
  - harness summary: coverage 70.7%（line 71.8%, branch 62.9%）
  - 備考: suite 内で `scan:sonar` が自動実行された。

## 未証明小範囲
- backend 側 `Model|ProviderSettings|MasterPersona|TranslationJobSetup|Redaction` の追加ケース拡張は未実施（今回は frontend usecase の指定初手を優先）。
- presenter/store 層の追加分岐拡張は未実施（既存テスト範囲で変更なし）。

## 残留リスク
- Vitest `--run` フィルタ挙動が環境差で不安定な可能性がある。
- `coverage` suite が Sonar 実行を内包するため、レーン制約運用とハーネス実装の乖離が残る。
