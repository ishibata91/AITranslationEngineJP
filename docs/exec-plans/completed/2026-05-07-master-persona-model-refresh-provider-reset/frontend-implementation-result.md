# frontend 実装結果

## 変更ファイル

- [master-persona-screen-controller.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts)
- [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)

## 実装内容

- モデル一覧更新操作の呼び先を、初期表示用の `loadAISettings()` から手動更新用の `refreshAISettings()` へ分離した。
- 手動更新では、現在 UI で選択している provider と実行方法を維持する。
- gateway が model を返した場合、返却 model を `aiSettings.model` と `modelOptions` へ反映する。
- 初期表示の保存済み AI 設定読込と、provider 変更時に model と `modelOptions` を空にする挙動は維持した。

## 検証結果

- 実行: `python3 scripts/harness/run.py --suite frontend-local`
- 結果: 失敗
- 失敗箇所: [translation-job-setup.usecase.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts)
- 行: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts:266`
- 内容: `getTranslationJobSetupOptions: vi.fn()` の型が `TranslationJobSetupGatewayContract` の戻り値型と一致しない。
- 補足: 失敗箇所は今回の承認済み実装範囲外である。
- 追加確認: `npm exec prettier -- --check src/application/usecase/master-persona/master-persona.usecase.ts src/controller/master-persona/master-persona-screen-controller.ts`
- 追加確認結果: 通過

## 残留リスク

- `frontend-local` は範囲外の型エラーで停止したため、harness 全体の通過は未確認である。
- 実画面で provider 維持と `fake-model` 選択可能状態の再確認は未実施である。

## follow-up: state-invariant-001 対応

## 変更ファイル

- [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)

## 実装内容

- `refreshAISettings()` 成功時の provider と実行方法は、`store.update` 時点の `draft.aiSettings` から維持する形へ変更した。
- refresh 開始後に provider が変わった場合、gateway が返した model と `modelOptions` を反映しない形へ変更した。
- refresh 失敗時は、開始時 snapshot の provider と `modelOptions` を復元しない形へ変更した。
- 初期表示の `loadAISettings()` と provider 変更時の model / `modelOptions` クリア処理は変更していない。

## 検証結果

- 実行: `npm --prefix frontend run test -- --run src/application/usecase/master-persona/master-persona.usecase.test.ts`
- 結果: 通過
- 詳細: `Test Files 1 passed (1)`、`Tests 28 passed (28)`
- 実行: `python3 scripts/harness/run.py --suite frontend-local`
- 結果: 失敗
- 失敗箇所: [translation-job-setup.usecase.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts)
- 行: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts:266`
- 内容: `getTranslationJobSetupOptions: vi.fn()` の型が `TranslationJobSetupGatewayContract` の戻り値型と一致しない。
- 補足: 失敗箇所は今回の承認済み実装範囲外である。
- 追加確認: `npm exec prettier -- --check src/application/usecase/master-persona/master-persona.usecase.ts`
- 追加確認結果: 通過

## 残留リスク

- `frontend-local` は範囲外の型エラーで停止したため、harness 全体の通過は未確認である。
- 実画面で refresh 中に provider を変更する競合操作は未確認である。
