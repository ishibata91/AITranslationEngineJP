# 合意済み frontend 保護

- `status`: approved
- `source`: frontend 実装後人間レビュー approve
- `ux_review`: `./ux-review.yaml`
- `implementation_scope`: `./implementation-scope.md`

## 承認済み画面

- 翻訳管理初期画面は、未完了のジョブ一覧を表示する。
- 未完了のジョブ一覧だけが、翻訳セクションから直接開ける。
- データロード、翻訳設定、各翻訳段階、翻訳結果の確認、出力管理は翻訳セクションでは参照表示にする。
- ジョブ選択後だけ、翻訳段階ページを表示する。
- データロード、翻訳設定、各翻訳段階、翻訳結果の確認は sticky footer で次の移動を扱う。
- 出力管理へ移動しても、出力対象ジョブは自動選択しない。

## 承認済み表示規則

- 利用者向け表示では `job` / `Job` を `ジョブ` と書く。
- 利用者向け表示では `phase` / `フェーズ` を原則 `翻訳段階` と書く。
- footer は移動先と移動できない理由だけを伝える。
- 旧 `Job Run` 直リンク、ジョブ ID 手入力、summary 取得入口は表示しない。
- secret、API key 平文、provider raw payload、過剰本文を UI、DTO、console、error summary に出さない。

## 確認済み実画面

- review URL: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management`
- UX 事前確認: `./ux-review.yaml`
- 検証: `python3 scripts/harness/run.py --suite frontend-local`
- 検証結果: pass

## 変更禁止範囲

後続 agent は、backend 実装、統合境界実装、テスト実装で次の frontend 表示と導線を変更しない。

- `frontend/src/ui/stores/shell-state.ts`
- `frontend/src/ui/views/AppShell.svelte`
- `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
- `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte`
- `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte`
- `frontend/src/ui/components/StickyActionFooter.svelte`
- `frontend/src/ui/screens/translation-output-artifact/TranslationOutputArtifactPage.svelte`

## 残留事項

- `UX-TFN-004` は minor として残る。
- 内容: 未完了のジョブ一覧で削除操作の視覚的強さを人間レビューで見る。
- blocker / major はない。

