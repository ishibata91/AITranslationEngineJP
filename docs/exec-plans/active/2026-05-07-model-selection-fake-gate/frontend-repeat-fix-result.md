# frontend 修正結果

## 変更ファイル

- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)

## 変更内容

- `TranslationJobSetupUseCase.refreshPhaseModels()` から `credentialStatus === "missing"` の事前 return を削除した。
- `missing` の phase でも `listTranslationJobSetupProviderModels()` へ進むようにした。
- backend 応答が利用不可の場合の失敗反映と、単一モデル一覧の自動選択は既存処理を維持した。
- frontend に fake provider ID、`fake-model` 固有分岐、provider catalog 追加は入れていない。

## 検証結果

- `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
  - 結果: 失敗。
  - 理由: 既存テスト `credentialStatus が missing の phase は refreshPhaseModels で gateway を呼ばず missing 状態を維持する` が、旧仕様どおり gateway 呼び出し 0 回を期待している。
  - 実際: gateway 呼び出しは 1 回になった。修正入力の期待どおり。
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: 失敗。
  - 通過: frontend lint harness。
  - 失敗: frontend test harness。
  - 理由: 上記と同じ既存テスト 1 件の旧期待値。

## 残留リスク

- プロダクトテスト変更は禁止範囲のため、旧期待値のテスト更新は未実施。
- UI 実画面確認は未実施。今回の変更は usecase の gateway 呼び出し条件だけで、UI 文言と layout は変更していない。
