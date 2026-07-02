# Storybook レビューループ記録: 汎用台詞・PC 発話への口調付与

履歴は残さず、確定結果だけを記録する。

## 起動

- URL: `http://localhost:6008/`
- 起動 command: `npm --prefix frontend run storybook`

## 対象 story（作業中分類 → 通常分類）

### 結果行コンポーネント（TranslationResultRow）

- 作業中分類: `Review/Changed Components/TranslationResultRow`
- 通常分類（承認後に戻す）: `UI Components/TranslationResultRow`
- 追加・変更した表示状態:
  - `展開（汎用既定・口調メタ）`: 見出し「汎用台詞」、根拠「感情 中 ・ 性別 男性」。対人段階・セル・印は出さない。
  - `展開（PC既定・口調メタ）`: 見出し「PC発話」、根拠「感情 抑制 ・ 性別 男性」。対人段階・セル・印は出さない。
  - `畳む（汎用台詞）`: 話者を特定できなくても口調チップが付く。
  - 既存の本文・声質・保留の口調メタ表示は非劣化（対人段階・印を出す経路のまま）。

### テンプレート編集画面（TemplateEditorScreen）

- 作業中分類: `Review/Changed Screens/プロンプトテンプレート`
- 通常分類（承認後に戻す）: `Screens/プロンプトテンプレート`
- 追加・変更した表示状態:
  - `レコード別タブ`: 「話者なし台詞の口調」節を追加（汎用・PC の自由記述欄 2 つ＋PC 性別の選択）。
  - `レコード別タブ（口調と PC 性別を変更）`: 汎用口調を書き換え、PC 性別を男性に選んだ表示。

## fixture・関連資源

- 結果行: `translation-run.fixtures.ts`（汎用台詞・PC 発話の行を追加）、story 内 `GENERIC_ROW`・`PC_ROW`。
- 結果行 presentation: `translation-run-presentation.ts`（`DECISION_PATH_LABEL`・`DECISION_PATH_HINT` に汎用・PC、口調メタ根拠の組み立て `personaMetaParts`）。
- テンプレート編集: `template-editor.fixtures.ts`（`DEFAULT_TEMPLATE_FORM` に自由記述 2 つ＋PC 性別、`TONE_DEFAULT_EDITED_FORM`）、`template-editor-presentation.ts`（`TONE_TEXT_FIELDS`・`PC_SEX_OPTIONS`）、コンポーネント `ToneDefaultPane.svelte`。

## frontend 表示実装境界（本モジュールで扱った範囲）

- 扱った: svelte 表示コンポーネント、props、story、fixture、表示文言、自由記述欄・性別選択の表示。
- 扱わない（implementation-module へ）: 汎用口調・PC 口調・PC 性別の state 保持・gateway・Wails bridge・保存配線。`PromptTemplateForm` の新フィールド（`genericToneText`・`pcToneText`・`pcSex`）と `PersonaMeta` の汎用・PC 用フィールドは配線前でも壊れないよう任意にした。`onFieldInput` を再利用し、新しい callback は足さない。

## 画面表示の根拠

- design.md・implementation-scope.md の確定判断（#2 既定対人段階の編集面、#4 公開境界の DecisionPath 新値・性別根拠表示、#5 PC 判定）。

## 承認状態

- 承認済み（2026-06-30）。story を通常分類へ戻した。
  - 結果行: `UI Components/TranslationResultRow`。
  - 編集画面: `Screens/プロンプトテンプレート`。
