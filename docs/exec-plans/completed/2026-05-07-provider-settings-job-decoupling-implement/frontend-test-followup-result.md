# Frontend Test Follow-up Result

- task: `2026-05-07-provider-settings-job-decoupling-implement`
- handoff: `PSJD-FE-01`
- role: `implementation_unit_tester`
- status: `completed`

## 変更ファイル

- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`

## 追従内容

- usecase test の `credentialStatus: missing` 分岐を新契約に合わせ、gateway 未呼び出しと `credential_missing` 維持を検証する期待へ更新した。
- usecase test の validation payload 期待を更新し、`credentialRef` は空文字を検証するようにした。
- JobSetupPage test の表示期待を UI 契約へ更新し、`設定済み` / `APIキー不要` / `要確認` / `再確認が必要` を検証するようにした。
- JobSetupPage test で `credential reference` 操作期待を削除し、`AIサービス / モデル / 実行方法` と `確認を実行` 操作委譲へ更新した。
- JobSetupPage test で禁止表示を検証し、`credential reference` と `modelListSourceToken` が画面表示されないことを確認した。

## 検証結果

- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup src/application/usecase/translation-job-setup src/application/presenter/translation-job-setup src/application/store/translation-job-setup src/controller/review-fake-api`
  - 結果: pass
  - Test Files: 6 passed
  - Tests: 41 passed
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: pass
  - lint: pass
  - frontend test: 57 passed / 494 passed

## 残失敗

- なし
