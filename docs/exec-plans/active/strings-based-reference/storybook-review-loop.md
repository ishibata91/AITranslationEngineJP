# Storybook レビュー記録: strings-based-reference

Storybook 人間レビューの確定結果を残す。人間コメントの履歴は残さず、確定した story・画面仕様・反映先・分類・承認状態だけを記録する。

## 確定した story

- 対象 story: `StringsOneSidedWarning`（表示名「Strings 片側警告」）。
- 反映先ファイル:
  - `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`
  - `frontend/src/ui/screens/translation-run/translation-run.fixtures.ts`
  - `frontend/src/ui/screens/translation-run/TranslationRunScreen.stories.ts`
- fixture: `STRINGS_ONE_SIDED_WARNING_STATE`。
- 起動 command: `npm --prefix frontend run storybook`（`http://localhost:6008/`）。

## 変更された画面仕様

- 翻訳実行画面（`TranslationRunScreen`）に警告表示を追加した。
- 追加 prop: `warning?: string`（空なら出さない）。
- 表示: `role="alert"` の `alert alert-warning alert-soft`（amber）。既存の `notice`（info）と `errorMessage`（error）の並びに合わせ、実行アクション直前に置く。provider（sync / xai）で出し分けない。
- 表示条件（画面表示の根拠）: localized plugin で english / japanese の Strings が片方しか無い場合に出す。両方揃えば警告を出さず既存訳を再利用する。非 localized（Strings 不使用）は出さない。
- 確定文言: 「英語と日本語の Strings の両方が揃っていないため、既存訳を再利用できません。全文を AI 翻訳します。」

## 分類

- 現在分類: 通常分類 `Screens`（`title: "Screens/翻訳実行"`）。追加した variant も同 title 配下。

## 承認状態

- Storybook 人間レビュー: 承認済み。

## 検証

- `npm run test:frontend`: 通過（1 file / 2 tests）。
- `npm --prefix frontend run build-storybook`: 通過。

## implementation-module へ渡す残課題（表示範囲外）

- `warning` prop への実データ配線（Container・state・API）。片方しか無いの判定は backend（C# 抽出器の Strings 読み）で行い、その結果を `warning` へ渡す。
- 既存の lint 赤（`TranslationRunScreen.svelte` の未使用 props `onRefresh` / `canRefresh` / `refreshing` / `canApply`）。変更前から存在する Container 互換の宣言由来で、除去は Container 配線に影響しうるため storybook-module では触らない。Container 配線を扱う implementation-module で解消するか、残す理由を確定する。
