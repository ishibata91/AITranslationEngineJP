# Storybook レビュー記録: record-type-translation-expansion

- task-id: 2026-06-23-record-type-translation-expansion
- 由来: storybook-module（design-module 承認後、画面が動くため）。
- 状態: Storybook 人間レビュー承認済み（2026-06-25）。確定結果のみを記録する。

## 起動

- 起動 command: `npm --prefix frontend run storybook`
- 起動 URL: `http://localhost:6008/`

## 確定した story（通常分類へ復帰済み）

| story | 通常分類（title） | 状態 |
| --- | --- | --- |
| プロンプトテンプレート画面 | `Screens/プロンプトテンプレート` | ベースタブ・レコード別タブ・未保存 |
| 指示文エディタ | `UI Components/DirectiveEditor` | 指示文一覧（口調・文体・固有名・定型句） |
| 結果行 | `UI Components/TranslationResultRow` | 種別バッジ（台詞・叙述文・定型句・固有名） |

## 確定した画面仕様

- プロンプトテンプレート画面にサブタブ [ベース][レコード別] を持つ。
  - ベース: Base 翻訳指示（全種別共通の前提）。
  - レコード別: 種別ごとの指示文（directive）を一律に並べる。各 directive は 本文（編集可）＋変数（あれば）＋対象 REC:FIELD（読み取り専用）。口調は `{traits}` 変数を持つ。固有名・定型句も指示文を編集できる。
- プロンプト = Base ＋ その REC:FIELD の directive（変数を実行時に埋める）。REC:FIELD → directive の割り当ては固定。翻訳しない種別は載せない。
- 結果行の先頭に元レコード種別バッジ（箱 ・ REC:FIELD）を出す。

## 反映先 frontend

- `frontend/src/ui/screens/template-editor/`: `TemplateEditorScreen.svelte`・`TemplateBasePane.svelte`・`DirectiveEditor.svelte`・`directive-view.ts`・`directive-presentation.ts`・`template-editor.fixtures.ts`・`template-editor-view.ts`・`TemplateEditorContainer.svelte`（placeholders 撤去）。
- `frontend/src/ui/screens/translation-run/`: `TranslationResultRow.svelte`・`translation-run-view.ts`。
- 削除: `frontend/src/ui/screens/record-type-master/` 一式。

## frontend 表示実装境界（implementation-module へ渡す残り）

- state・gateway・Wails bridge・タブ状態・directive データの読み出しと保存の配線。
- `TemplateEditorContainer` を新 props（activeTab・directives・assignments・onInstructionInput・onTabChange）へ配線し直す。
- 結果行の `recordType` を実データから埋める配線。

## 検証

- `npm --prefix frontend run check`: 通過（残る 1 error は node_modules の @storybook/svelte 型宣言で既存・本変更外）。
- `frontend-local` suite（eslint・tsc・knip・boundaries・vitest）: 通過。
- `npm --prefix frontend run build-storybook`: 通過。
