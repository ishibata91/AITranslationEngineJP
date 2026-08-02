# Storybook Review Loop: storybook-screen-spec-harness

## 対象 story

| 画面 | 対象 story | fixture | 画面仕様ID数 |
| --- | --- | --- | ---: |
| 翻訳対象プラグイン | `Empty`、`Loading`、`List`、`Selected`、`ConfirmDelete`、`Deleting`、`Errored` | `target-plugins.fixtures.ts` | 12 |
| プロンプトテンプレート | `BaseTab`、`RecordTab`、`RecordTabToneDefaultEdited`、`RecordTabDirty` | `template-editor.fixtures.ts` | 10 |
| 翻訳実行 | `NotStarted`、`Running`、`Paused`、`Done`、`DoneWithUntranslated`、`Failed` | `translation-run.fixtures.ts` | 21 |

## 関連資源

- 画面仕様と画面状態は、各画面の `*-screen-specs.ts` に置く。
- 前提条件と仕様ID付きの仕様文は、`screenStateDescription` が Autodocs の story 説明へ変換する。
- 画面表示は既存の `*Screen.svelte` と既存 fixture を使う。

## 起動

- URL: `http://localhost:6008/`
- command: `npm --prefix frontend run storybook`

## 分類

- 人間レビュー中: `Review/Changed Screens/<画面名>`
- 承認後: `Screens/<画面名>`

## 表示実装とロジック実装の境界

- 今回の Storybook レビューは、既存画面状態、Autodocs の前提条件、43件の画面仕様ID、仕様文の対応を扱う。
- 画面の表示、文言、操作、fixture値、state、API、Wails境界、ルーティング、副作用は変更しない。
- 単体テストが43件の画面仕様IDを消費する処理は、人間レビュー承認後に `implementation-module` で扱う。

## 画面表示の根拠

- `design.md` の R-1 にある三つの対応表を根拠とする。
- 対応表にない画面状態、画面仕様ID、画面仕様は追加しない。

## 承認状態

- 承認済み。
- 人間レビューで、Autodocs のプロパティ表を不要と判断した。画面用の Docs ページはプロパティ表を表示しない。
- 人間レビューで、画面用の Autodocs をダークテーマにすることを確定した。
- 承認後の分類は `Screens/<画面名>` とする。
