# frontend 実装入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-FE-01`
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `ready_wave`: `wave-1`
- `source_scope`: `implementation-scope.md`

## 目的

Job Setup から endpoint、`credential_ref` 実値、secret store 参照実値を表示しない。
Job Setup は provider、model、execution mode、batch mode、APIキー状態分類だけを扱う画面契約へ寄せる。

## 読むファイル

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/ui-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/implementation-scope.md`
- `docs/frontend-fake-api.md`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`
- `frontend/src/application/store/translation-job-setup/translation-job-setup.store.ts`

## 変更許可範囲

- `frontend/src/ui/screens/translation-job-setup/`
- `frontend/src/application/usecase/translation-job-setup/`
- `frontend/src/application/presenter/translation-job-setup/`
- `frontend/src/application/store/translation-job-setup/`
- `frontend/src/controller/review-fake-api/`
- 上記範囲の frontend test

## 禁止範囲

- `frontend/src/controller/wails/`
- `frontend/wailsjs/`
- `internal/`
- docs 正本本文
- `.codex/`

## secret 境界

表示してよい値:
provider、model、execution mode、batch mode、APIキー状態分類。

表示禁止:
endpoint 原文、`credential_ref` 実値、secret store key 名、APIキー本体、復号可能値、raw request、raw response、raw prompt、`modelListSourceToken`。

## 初手

- path: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- 対象: legacy fallback の `credential reference` select
- 変更種別: 表示削除

## 完了条件

- Job Setup は 3 つの翻訳段階で AIサービス、モデル、実行方法、一括処理、APIキー状態分類だけを表示する。
- 作成後要約は endpoint、`credential_ref`、secret store 参照実値を表示しない。
- APIキー未設定、モデル未選択、モデル一覧未更新、モデル一覧取得失敗を別状態で表示する。
- fakeAPI で `success`、`config-missing`、`error` を確認できる。
- 人間レビュー入力として review URL、確認状態、未確認状態、未確認理由を task 成果物へ記録する。

## 検証コマンド

- `npm --prefix frontend run test -- src/ui/screens/translation-job-setup src/application/usecase/translation-job-setup src/application/presenter/translation-job-setup src/application/store/translation-job-setup src/controller/review-fake-api`
- `python3 scripts/harness/run.py --suite frontend-local`

## 期待出力

- `frontend-implementation-result.md`
- 変更ファイル一覧
- 検証結果
- fakeAPI review URL と確認状態
- 未確認状態と未確認理由

