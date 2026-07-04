# Storybook レビューループ記録: xtranslator-export

## 対象

- 確認対象: 結果一覧パネル（`ResultsPanel`）へ追加した「xTranslator へ書き出し」ボタンの表示。
- 変更コンポーネント: `frontend/src/ui/screens/translation-run/ResultsPanel.svelte`（ヘッダ右に書き出しボタンを追加）。
- 配線のみ変更: `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`（`onExportXml`・`exporting` を受け `ResultsPanel` へ流す）。

## story

- story ファイル: `frontend/src/ui/screens/translation-run/ResultsPanel.stories.ts`
- 追加・変更した表示状態:
    - `結果あり（書き出し可能）`: 結果があり、書き出しボタンが有効。
    - `書き出し中`: `exporting=true`。ボタン無効化＋スピナー＋「書き出し中…」。
- 変更なしで書き出しボタンが出ない状態: `空（未実行）`・`空（実行中）`（結果 0 件のためボタン非表示）。

## fixture・関連資源

- props fixture: `ROWS`（`ResultsPanel.stories.ts` 内、結果 2 件）。
- 書き出しボタンは `results.length > 0` のときだけ表示する（空状態では出さない）。

## 分類

- 作業中分類（レビュー中に使用）: `Review/Changed Components/ResultsPanel`
- 通常分類（承認後に戻した）: `UI Components/ResultsPanel`
- 現在分類: 通常分類（承認済み）。

## 承認・検証結果

- 人間 Storybook レビュー: 承認済み（配置・表示条件・書き出し中の見た目・文言「xTranslator へ書き出し」）。
- 確定した story: `結果あり（書き出し可能）`、`書き出し中`（`ResultsPanel.stories.ts`）。
- 検証:
    - `npm --prefix frontend run build-storybook` 通過（exit 0）。
    - `npm --prefix frontend run test`（vitest）通過（exit 0）。
    - `npm --prefix frontend run check`（svelte-check）は既存の型宣言エラー 1 件のみ（`node_modules/@storybook/svelte/dist/index.d.ts` の宣言不足）。変更前後で同一（stash で確認済み）であり、本 task の 3 ファイル変更が原因ではない。
- skill 記載の `python3 scripts/harness/run.py --suite frontend-local` は harness に実体が無い（`scripts/harness/run.py` 不在）。frontend の npm scripts（build-storybook・test・check）で代替した。

## 起動

- 起動 command: `npm --prefix frontend run storybook`
- 固定 URL: `http://localhost:6008/`

## frontend 表示実装境界

- 本モジュールで扱うのは表示のみ（ボタンの layout・文言・状態別 style、story、fixture）。
- 扱わない（implementation-module へ）: 書き出しの実処理、`TranslationRunContainer` の `exporting` 状態と `onExportXml` ハンドラ、gateway ラッパ、Wails bindings。
