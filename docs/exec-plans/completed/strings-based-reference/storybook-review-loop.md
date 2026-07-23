# Storybook Review Loop: strings-based-reference

片側 Strings 欠け時の画面警告表示のレビュー記録。確定結果のみを残す。

## 対象 story（承認済み・通常分類へ復帰済み）

- `UI Components/MissingStringsWarning`
  - 日本語欠け / 英語欠け / 両方欠け / 揃い（非表示） / 未判定（非表示）
- `Screens/翻訳実行`: 「日本語 Strings 欠け」（画面文脈での確認用。既存 title へ追加）

## 確定した画面仕様

- 警告は見出し＋「状態と理由」「影響」「対処」の 3 段構成。daisyUI `alert alert-warning alert-soft`、`role="alert"`。
- 配置は `TranslationRunScreen` の「翻訳対象」セクション直下。`presence` 両方 true または未指定（未判定）では非表示。

## 関連資源

- コンポーネント: `frontend/src/ui/screens/translation-run/MissingStringsWarning.svelte`
- story: `frontend/src/ui/screens/translation-run/MissingStringsWarning.stories.ts`・`TranslationRunScreen.stories.ts`
- fixture: `frontend/src/ui/screens/translation-run/translation-run.fixtures.ts`（`MISSING_JAPANESE_STRINGS_STATE`）
- 型: `frontend/src/ui/screens/translation-run/translation-run-view.ts`（`StringsPresence`）
- 配置先: `TranslationRunScreen.svelte` の「翻訳対象」セクション直下（optional prop `stringsPresence`）

## 起動

- command: `npm --prefix frontend run storybook`
- URL: `http://localhost:6008/`

## 表示実装とロジック実装の境界

- 本モジュールの範囲: 表示コンポーネント・story・fixture・props 形（`presence?: { english: boolean; japanese: boolean }`）。
- implementation-module の範囲: 欠落状態の判定（backend が対象 plugin の Data フォルダ配下 `strings/` を言語別に判定）と `TranslationRunContainer` から `stringsPresence` への供給。未判定の間は prop を渡さず警告は出ない。

## 画面表示の根拠

- design.md「片側 Strings 欠け時の警告」: 片側欠けだと英日対を作れず、参照訳・固有名の確定訳語を再利用できず全文 AI 翻訳になる旨を利用者へ知らせる。

## 承認状態

- 承認済み（2026-07-24）。文言差し戻し 1 回（「欠け」だけでは伝わらない → 3 段構成へ修正）を経て承認。story は通常分類 `UI Components/MissingStringsWarning` へ復帰済み。
