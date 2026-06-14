# Storybook レビューループ記録（results-paging-bulk-persona）

確定結果のみを記録する。人間コメントの履歴は残さない。

## 起動

- 固定 URL: `http://localhost:6008/`
- 起動 command: `npm --prefix frontend run storybook`
- build 確認 command: `npm --prefix frontend run build-storybook`

## 対象 story（作業中分類）

| 確認資源 | 作業中分類 story id | 状態名 |
| --- | --- | --- |
| 追加コンポーネント | `review-changed-components-resultspager--first-page` | 先頭ページ（前へ無効） |
| 追加コンポーネント | `review-changed-components-resultspager--middle-page` | 中間ページ（両方有効） |
| 追加コンポーネント | `review-changed-components-resultspager--last-page` | 末尾ページ（次へ無効） |
| 追加コンポーネント | `review-changed-components-resultspager--single-page` | 単一ページ（両方無効） |
| 変更コンポーネント | `review-changed-components-resultspanel--paging-first` | 複数ページ・先頭 |
| 変更コンポーネント | `review-changed-components-resultspanel--paging-middle` | 複数ページ・中間 |
| 変更コンポーネント | `review-changed-components-resultspanel--paging-last` | 複数ページ・末尾 |

## 通常分類（承認後の戻し先）

- 追加コンポーネント `ResultsPager` → `UI Components/ResultsPager/<状態名>`（`ResultsPager.stories.ts` の title を変更）
- 変更コンポーネント `ResultsPanel` の paging 状態 → `UI Components/ResultsPanel/<状態名>`（既存 `ResultsPanel.stories.ts` へ統合、一時 `ResultsPanelPaging.stories.ts` は削除）
- 変更画面 `翻訳実行` の複数ページ状態 → `Screens/翻訳実行/<状態名>`（承認後に fixture へ paging を足して既存 story へ統合）

## fixture / 関連資源

- 一時 fixture: `ResultsPanelPaging.stories.ts` 内の `PAGE_ROWS`（総件数 121・ページサイズ 50 を想定した 1 ページ分）。
- 関連資源: `translation-run-view.ts`（`ResultsPaging` 型）。

## frontend 表示実装境界（storybook-module で扱う範囲）

- 扱う: `ResultsPager.svelte`（新規ページャ）、`ResultsPanel.svelte`（総件数バッジ・ページャ配置・paging props）、`TranslationRunScreen.svelte`（paging props の受け渡し）、story、fixture、表示用 view 型（`ResultsPaging`）。
- 扱わない（implementation-module へ）: keyset cursor 履歴 state、`listResultsPage` の gateway 経路、container のページ state とページ送り操作、`RunExtractAndTranslate`／`ListResultsPage` 呼び出し。

## 画面表示の根拠

- `plan.md` の実装範囲（keyset cursor ページング、ページャ表示）と設計判断（判断 1=keyset cursor・ページサイズ 50、判断 3=連結列の順次送り）。

## 確定結果

- 承認状態: 承認済み。
- 確定した変更画面仕様: 結果一覧へ keyset ページャを追加。順次送り（← 前へ ｜ ページ N ｜ 次へ →）で番号ジャンプは持たない。先頭で「前へ」、末尾で「次へ」を無効化（淡色）。件数バッジを現在ページ件数から総件数（例: 121 件）へ変更。結果行（`TranslationResultRow`）の表示形（コンパクト 1 行・口調チップ・展開）は不変。
- 現在分類: 全 story を通常分類へ復帰済み。`UI Components/ResultsPager`（新規）、`UI Components/ResultsPanel`（paging story を統合）、`Screens/翻訳実行`（DONE 系 fixture に paging を付与）。作業中分類（`Review/Changed/...`）と一時 `ResultsPanelPaging.stories.ts` は削除済み。
- 反映先: `ResultsPager.svelte`（追加）、`ResultsPanel.svelte`・`TranslationRunScreen.svelte`（変更）、`translation-run-view.ts`（`ResultsPaging` 型）、`ResultsPager.stories.ts`・`ResultsPanel.stories.ts`・`TranslationRunScreen.stories.ts`・`translation-run.fixtures.ts`。
- 検証: `build-storybook` 成功、`frontend-local` suite 通過。console エラーなし。実画面（localhost:6008）で中間ページ（両端有効）・先頭ページ（前へ無効）を確認。
