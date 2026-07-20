# Storybook レビューループ画面仕様

## 対象

- 確定した story: `TranslationRunScreen.stories.ts`（`Screens/翻訳実行`）
- 確定した `fixture`: `translation-run.fixtures.ts`（`EMPTY_FORM.tokenBudget=""`、`FILLED_FORM.tokenBudget="2"`）
- 関連資源: `TranslationRunScreen.svelte`、`translation-run-view.ts`、`translation-run-presentation.ts`、`TextField.svelte`、`Field.svelte`
- 作業中分類: `Review/Changed Screens/翻訳実行`
- 通常分類: `Screens/翻訳実行`
- 現在分類: `Screens/翻訳実行`（承認後に通常分類へ戻し済み）

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| 翻訳実行画面「AI サービス」カード | モデル選択の隣に「最大トークン」欄を追加。トークン予算を千トークン（k）単位で入力。入力ボックスは約 1/4 幅（`w-24`）で、数値の右に単位「k」を表示 | `TranslationRunScreen.svelte`、`translation-run-presentation.ts`（`TOKEN_BUDGET_FIELD`） | なし |
| `TextField` 共有部品 | 任意 `suffix?`（単位表示）と任意 `inputWidthClass?`（既定 `w-full`）を追加。既存の利用箇所は見た目不変 | `TextField.svelte` | なし |
| `Field` 共有部品 | hint を折り返し表示（`whitespace-normal`）にし、長い hint がカードをはみ出さないようにした | `Field.svelte` | なし |

## 現在状態

- 変更ファイル:
  - `frontend/src/ui/screens/translation-run/translation-run-view.ts`
  - `frontend/src/ui/screens/translation-run/translation-run-presentation.ts`
  - `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`
  - `frontend/src/ui/screens/translation-run/translation-run.fixtures.ts`
  - `frontend/src/ui/components/TextField.svelte`
  - `frontend/src/ui/components/Field.svelte`
- 通常分類へ戻した story: `Screens/翻訳実行`
- 承認状態: Storybook 人間レビュー承認済み
