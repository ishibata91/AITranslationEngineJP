# テスト修正証跡

## 状態

- `artifact`: `テスト修正証跡`
- `status`: `completed`
- `agent`: `implementation_unit_tester`

## 対象テスト範囲

- [InputReviewPage.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.test.ts)
- [translation-input-app.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/translation-input-app.test.ts)

## 検証目的

- データロード画面への表示名変更後も、画面本体が描画されることを証明する。
- route を離れて戻っても、読み込み済みデータ一覧、選択状態、再構築入口が維持されることを証明する。
- 旧 `Input Review` 見出し期待を残さない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite frontend-local`

## 変更内容

- [InputReviewPage.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.test.ts): `h2` 見出し期待を `データロード` へ更新し、旧 `Input Review` が存在しないことを追加で検証した。
- [translation-input-app.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/translation-input-app.test.ts): route 初回描画と復帰後描画の見出し期待を `データロード` へ更新し、旧 `Input Review` 非表示を検証した。

## 検証結果

- `python3 scripts/harness/run.py --suite frontend-local`: pass
- `check_frontend_lint.py`: pass
- `check_frontend_test.py`: pass (48 files / 441 tests)
- 補足: `vite-plugin-svelte` の `optimizeDeps.esbuildOptions` deprecation warning は既存の注意であり、失敗原因ではない。

## 証明済み完了条件

- データロード画面への表示名変更後も、画面本体が描画されることを単体テストで証明した。
- route を離れて戻っても、読み込み済みデータ一覧、選択状態、再構築入口が維持されることを単体テストで証明した。
- 旧 `Input Review` 見出し期待をテストから除去し、非表示検証へ置換した。
