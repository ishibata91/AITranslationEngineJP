# Storybook レビューループ記録: 翻訳実行画面

storybook-module の確定結果だけを残す。人間コメントの履歴は残さない。

## 承認状態

- Storybook 人間レビュー: 承認済み（「UI は OK」）。
- 現在分類: 通常分類へ復帰済み（`Screens/翻訳実行`、`UI Components/*`）。

## 確定した画面仕様

- 画面目的: plugin を選んで叙述文を抽出し、AI 翻訳を実行して原文・訳文を 1 画面で確認する。
- 入力（業務単位で 2 群）:
  - 翻訳対象: plugin の**ファイル選択**（フルパス）。Data フォルダは plugin の親ディレクトリとして派生し、画面入力には持たない。
  - AI サービス: エンドポイント（テキスト）、API キー（マスク・非永続）、モデル（**取得した一覧からの選択**。`getModels` を呼ぶ「取得」ボタン付き）。
- 主操作: 「実行」1 つ。状態バッジ（未実行・実行中・完了・失敗）を併置。
- 結果一覧: 見開き対訳（左 原文 / 右 訳文）。未訳は「（訳文なし）」、状態は意味トーンのバッジ。
- 状態網羅: 空状態・モデル取得中・入力済み・実行中・完了・エラー。

## 確定した story（通常分類）

- 画面: `Screens/翻訳実行` — 空状態 / モデル取得中 / 入力済み / 実行中 / 完了 / エラー（6）。
- 部品: `UI Components/StatusBadge`（6）, `UI Components/TextField`（3）, `UI Components/SelectField`（3）, `UI Components/FileSelectField`（2）, `UI Components/TranslationResultRow`（2）, `UI Components/ResultsPanel`（3）。

## 反映先 frontend ファイル

- テーマ: `frontend/src/ui/styles/app.css`（Tailwind v4 + daisyUI v5、独自テーマ `dovahkael`）。
- 共有部品: `frontend/src/ui/components/{Field,TextField,SelectField,FileSelectField,StatusBadge}.svelte`、`status-badge.ts`。
- 画面・ドメイン部品: `frontend/src/ui/screens/translation-run/{TranslationRunScreen,ResultsPanel,TranslationResultRow}.svelte`、`translation-run-view.ts`、`translation-run.fixtures.ts`。
- story: 各 `*.stories.ts`。
- 配線: `frontend/vite.config.ts`、`frontend/.storybook/{main.ts,preview.ts}`（Tailwind plugin 追加・app.css 読込）。旧 `preview.css`（amber）は削除。

## 検証結果

- `npm --prefix frontend run build-storybook`: 成功。
- frontend lint: `eslint` 通過、`tsc --noEmit` 通過。
- `knip`（`lint:exports`, `--production`）: 未達。production entry（`main.ts`）が未配線のため、本画面・部品・既存の未配線ファイル（`shell-state.ts`・diagnostic・`pino`）が production グラフから到達不能。greenfield reset の構造状態であり、implementation-module の production 配線で解消する。
- `lint:boundaries`: boundary テストファイルが reset で不在のため no-op（同上、構造状態）。

## 表示範囲外の残課題（implementation-module へ）

- ファイル選択ダイアログ（Wails `OpenFileDialog`）の配線と、plugin パス → Data フォルダ派生。
- `getModels`（OpenAI 互換 `/v1/models`）の gateway 配線と `models`/`modelsLoading` の state 化。
- フォーム validation（必須充足 → `canRun`）と実行手続き（抽出 → 翻訳 → 結果反映）の配線。
- production entry（`main.ts`）作成と本画面の wiring（knip/boundaries の構造未達もこれで解消）。
- フォントの self-host（現在は Google Fonts CDN）。
