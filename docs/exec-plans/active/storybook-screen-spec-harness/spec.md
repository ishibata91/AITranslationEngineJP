# Spec: storybook-screen-spec-harness

`spec.md` はこのtaskの確定仕様として、要求ごとの仕様を持つ。
要求は `plan.md`、設計理由・変更箇所・図は `design.md` が持つ。

---

## R-1 screens の画面仕様を Autodocs で読めるようにする

- R-1-1（正常系）: 画面状態の前提条件と、画面仕様ID付きの表示、状態バッジ、ボタンの文言、ボタンの操作可否を、`Screens`の画面単位のAutodocsで読めること
    - 前提条件: `frontend/src/ui/screens/` の画面状態を再現するfixtureとstoryがある。
    - 確かめ方: `Screens/翻訳対象プラグイン`、`Screens/プロンプトテンプレート`、`Screens/翻訳実行`のAutodocsに、画面状態の前提条件と画面仕様ID付きの画面仕様が表示されることを確認する。
    - 対応する実テスト: `npm --prefix frontend run build-storybook`と、`storybook-review-loop.md`に記録した三画面のAutodocs人間確認
- R-1-2（対象に入る側の境界）: 一つの画面状態に複数の画面仕様がある場合に、各画面仕様が異なる画面仕様IDでAutodocsに表示されること
    - 前提条件: 一つの画面状態が、状態バッジ、ボタンの文言、ボタンの操作可否を同時に持つ。
    - 確かめ方: 該当する画面状態のAutodocsに、複数の画面仕様と各画面仕様IDが表示されることを確認する。
    - 対応する実テスト: `npm --prefix frontend run build-storybook`と、`storybook-review-loop.md`に記録した43件の画面仕様IDの人間確認
- R-1-3（対象に入らない側の境界）: `UI Components` とナビゲーション確認用のstoryを画面仕様のAutodocsへ含めないこと
    - 前提条件: `frontend/src/ui/screens/`に、`UI Components`のstoryまたはナビゲーション確認用のstoryがある。
    - 確かめ方: `Screens`の三つの画面単位のAutodocsだけに画面仕様IDが表示されることを確認する。
    - 対応する実テスト: `storybook-review-loop.md`に記録した対象storyと分類の人間確認

---

## R-2 画面仕様を単体テストが消費していることをハーネスで確かめる

- R-2-1（正常系）: すべての画面仕様IDに対応する単体テストの検証関数が登録され、表示、状態バッジ、ボタンの文言、またはボタンの操作可否を確かめること
    - 前提条件: 画面仕様に、重複しない画面仕様IDが付いている。
    - 確かめ方: frontendのテスト出力に、各画面仕様IDを含む単体テストの成功が表示されることを確認する。
    - 対応する実テスト: `TargetPluginsScreen.spec.test.ts`、`TemplateEditorScreen.spec.test.ts`、`TranslationRunScreen.spec.test.ts`
- R-2-2（対象に入る側の境界）: 一つの画面状態に複数の画面仕様IDがある場合に、すべての画面仕様IDへ異なる単体テストの検証関数が登録されること
    - 前提条件: 一つの画面状態に複数の画面仕様IDがある。
    - 確かめ方: frontendのテスト出力に、同じ画面状態にある各画面仕様IDを含む単体テストの成功が一件ずつ表示されることを確認する。
    - 対応する実テスト: `TargetPluginsScreen.spec.test.ts`、`TemplateEditorScreen.spec.test.ts`、`TranslationRunScreen.spec.test.ts`の各画面仕様IDを含むtest case
- R-2-3（対象に入らない側の境界）: 単体テストの検証関数がない画面仕様ID、画面仕様にない画面仕様ID、または重複する画面仕様IDがある場合にハーネスが失敗すること
    - 前提条件: 画面仕様IDと単体テストの検証関数の対応に、不足、余分、または重複がある。
    - 確かめ方: ハーネス自身の単体テストの出力に、不足、余分、画面仕様側の重複、単体テスト側の重複をそれぞれ検出した結果が表示されることを確認する。
    - 対応する実テスト: `screen-spec-harness.test.ts`

**満たさない部分**: Storybookの `play`、backend、Wails境界、gateway、containerの処理が正しいことは、この画面仕様ハーネスでは検証しない。
