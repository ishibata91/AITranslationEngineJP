# unit repeat test result

## 対象

- task: `2026-05-07-model-selection-fake-gate`
- test scope: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- 実装変更は対象外。単体テストのみ更新。

## 変更内容

- `credentialStatus が missing の phase` テストを新仕様へ更新した。
- `refreshPhaseModels("word_translation")` で gateway 呼び出しが 1 回であることを検証した。
- backend 応答 `credentialStatus: not_required` かつ単一モデル時に `model` が返却 `modelId` へ自動選択されることを検証した。

## 検証結果

- `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
  - 結果: pass (1 file, 18 tests)
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: pass
  - frontend lint harness: pass
  - frontend test harness: pass
- `python3 scripts/harness/run.py --suite coverage`
  - 結果: pass
  - Sonar coverage summary: `coverage=70.6%`（70.0% 以上）

## 残留リスク

- `refreshPhaseModels` の backend エラー分岐（通信失敗時の挙動）は今回の更新対象外。
- Sonar の SCM 警告（blame 情報不足）が 1 件出るが、テスト成否には影響しない。
