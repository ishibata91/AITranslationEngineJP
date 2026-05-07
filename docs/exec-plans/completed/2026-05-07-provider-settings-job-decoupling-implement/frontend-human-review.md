# frontend 実装後人間レビュー入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-FE-01`
- `status`: `approved`
- `next_artifact_after_approval`: `backend 実装`

## レビュー対象

- `frontend-implementation-result.md`
- `frontend-test-followup-result.md`
- `ui-design.md`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- `frontend/src/controller/review-fake-api/default-review-fake-api-gateway-registry.ts`

## 確認してほしいこと

- Job Setup に endpoint 原文、`credential_ref` 実値、secret store key 名が表示されない。
- Job Setup は AIサービス、モデル、実行方法、一括処理、APIキー状態分類だけを表示する。
- APIキー未設定、モデル未選択、モデル一覧未更新、モデル一覧取得失敗が別状態で見える。
- 作成後要約に endpoint、`credential_ref`、secret store 参照実値が表示されない。

## review URL

- success: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- config-missing: `http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#translation-management`
- error: `http://localhost:34115/?fakeApi=1&fakeScenario=error#translation-management`

## 確認済み状態

- `success`: セットアップ表示で 3 段階の AIサービス、モデル、APIキー状態、モデル一覧状態、一括処理を確認済み。
- `config-missing`: セットアップ表示で APIキー未設定、モデル一覧更新不可、作成不可を確認済み。
- `error`: セットアップ表示で作成前確認失敗と作成不可を確認済み。

## 未確認状態

- 作成後要約への遷移は未確認。

未確認理由:
agent-browser の `次へ` click 後も画面が要約へ遷移しなかったため。
console error は出ていない。

## 検証済み

- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup src/application/usecase/translation-job-setup src/application/presenter/translation-job-setup src/application/store/translation-job-setup src/controller/review-fake-api`
  - 結果: pass
  - Test Files: 6 passed
  - Tests: 41 passed
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: pass
  - frontend lint: pass
  - frontend test: 57 files / 494 tests pass

## 次に進む条件

人間が frontend 実装後レビューを承認した場合だけ、`PSJD-BE-01` の backend 実装へ進む。
差し戻しの場合は、差し戻し内容を `frontend 実装` へ戻す。

## 承認記録

- `approved_at`: `2026-05-07`
- `human_response`: `フロントで見たいものないから先進めていい`
- `review_decision`: `approved_without_additional_visual_check`
