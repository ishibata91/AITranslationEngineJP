# Investigation 概要

## 問題

`JobRunPage.svelte` の `$effect` が `currentPhasePage` を依存として登録してしまい、ユーザー操作で phase を進めても `$effect` 再実行が初期 phase に戻す。結果、E2E-UC-048（term→persona）と E2E-UC-049（body→complete）が fail。

原因箇所: `setCurrentPhasePage` 内の `const previous = currentPhasePage`（line 113）が `$state` を読むため、`$effect` から呼ばれた時に依存登録される。

## 修正

`frontend/src/ui/screens/job-run/JobRunPage.svelte`:

- `import { untrack } from "svelte"` 追加
- line 113 を `const previous = untrack(() => currentPhasePage)` に変更

## テスト追加

- E2E-UC-048-B1（境界）: persona 遷移後に term に戻らないことを確認（500ms 待機）
- 単体テスト: `setCurrentPhasePage` 呼び出し後 `$effect` が再実行されないこと

詳細は同フォルダの `fix-decision.md` / `missing-usecases.md` / `missing-tests.md` / `missing-selectors.md` / `test-design.csv` を参照。
