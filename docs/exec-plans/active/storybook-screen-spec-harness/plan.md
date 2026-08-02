# Task Plan: storybook-screen-spec-harness

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

- `frontend/src/ui/screens/` の画面仕様は、画面の story、fixture、Svelte コンポーネント、単体テストへ分かれており、人間が画面単位で読める仕様ID付きの一覧がない。
- 画面の story と単体テストが同じ画面仕様を扱っているかを機械的に確かめる手段がない。

## 要求

- **R-1 screens の画面仕様を Autodocs で読めるようにする**: 既存の`*Screen.stories.ts`が持つ画面状態を対象に、story間で異なるfixtureの値が`*Screen.svelte`の直接の分岐で変える表示、状態バッジ、ボタンの文言、ボタンの操作可否と、storyの説明文が明記する非表示を画面仕様ID付きで抽出し、fixtureとstoryを使って、`Screens`の画面単位のStorybook Autodocsに表示する。状態間で変わらずstoryの説明文にもない表示、`UI Components`の内部表示、ナビゲーション確認用のstoryは対象に含めない。
- **R-2 画面仕様を単体テストが消費していることをハーネスで確かめる**: R-1の画面状態、表示、状態バッジ、ボタンの文言、ボタンの操作可否を持つ画面仕様IDと単体テストの検証関数を対応させ、画面仕様IDの不足、余分、重複をハーネス自身の単体テストとfrontendのテスト出力で検出する。

## branch 情報

- `execution_branch`: `codex/storybook-screen-spec-harness`
- `target_branch`: `master`
- `source_commit`: `a6235add7b777248a87a1b45e7a989d5c296fb61`

## やらないこと

- `frontend/src/ui/screens/` の画面表示、文言、操作、状態遷移を変更しない。
- `UI Components` の story とナビゲーション確認用の storyを画面仕様の対象に含めない。
- Storybook の `play` を使わない。
- ボタンを押した時にcallbackが呼ばれることを画面仕様の対象に含めない。
- backend、Wails境界、gateway、containerの処理を画面仕様ハーネスで検証しない。
- backendの `test-oracle/` と既存の単体テスト分類を変更しない。
