# Frontend Implementation Result

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-FE-01`
- `implementation_skill`: `implement-frontend`
- `status`: `implemented-with-test-followup`

## 変更ファイル

- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`
- `frontend/src/controller/review-fake-api/default-review-fake-api-gateway-registry.ts`

## 実装結果

- Job Setup の legacy fallback から `credential reference` 選択表示を削除した。
- Job Setup の表示文言を AIサービス、モデル、実行方法、一括処理、APIキー状態分類へ寄せた。
- presenter の画面用 view model から `credentialRef` と `modelListSourceToken` を空値化した。
- APIキー未設定時はモデル一覧取得 gateway を呼ばず、`credential_missing` 状態へ閉じた。
- fakeAPI の作成後要約は作成時の phase 選択を返し、3 段階の要約を表示できるようにした。

## 検証結果

- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup src/application/usecase/translation-job-setup src/application/presenter/translation-job-setup src/application/store/translation-job-setup src/controller/review-fake-api`
  - 結果: 失敗。
  - 失敗: 5 件。
  - 理由: 既存 test が `credential reference` 表示、`登録済み` / `不要` 文言、APIキー未設定時の gateway 呼び出し、`credentialRef: "gemini-primary"` payload を期待しているため。
- `npm --prefix frontend run lint:types`
  - 結果: 通過。
- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: 失敗。
  - lint: 通過。
  - frontend test: 上記 5 件で失敗。

## fakeAPI 確認

review URL:

- `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- `http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#translation-management`
- `http://localhost:34115/?fakeApi=1&fakeScenario=error#translation-management`

確認状態:

- `success`: セットアップ表示で 3 段階の AIサービス、モデル、APIキー状態、モデル一覧状態、一括処理を確認した。
- `config-missing`: セットアップ表示で APIキー未設定、モデル一覧更新不可、作成不可を確認した。
- `error`: セットアップ表示で作成前確認失敗と作成不可を確認した。

未確認状態:

- 作成後要約への遷移は未確認。

未確認理由:

- agent-browser の `次へ` click 後も画面が要約へ遷移しなかったため。
- console error は出ていない。

## 許可範囲外の残件

- 既存 frontend test の期待値更新が必要である。
- Wails / DTO / backend 側で `credentialRef` と `modelListSourceToken` の公開契約を縮める作業は PSJD-INT-01 以降の範囲である。
