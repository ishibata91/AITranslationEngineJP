# task 枠

## 人間依頼

- 対象画面は翻訳管理の `Input Review` 画面である。
- 画面概念は `input review` ではなくデータロード画面へ変える。
- UX観点を適用し、ロード対象、状態、結果、再構築判断を判断しやすくする。
- コンポーネント化を考慮して作り直す。

## 既存画面根拠

- [InputReviewPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.svelte)
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:115`
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:137`
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:192`
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:257`
- [shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)
  - `frontend/src/ui/stores/shell-state.ts:77`
- [UX-standard.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/UX-standard.md)
  - `docs/UX-standard.md:12`
  - `docs/UX-standard.md:14`
  - `docs/UX-standard.md:16`
  - `docs/UX-standard.md:40`

## 変更許可範囲

- [InputReviewPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.svelte)
- `frontend/src/ui/screens/translation-input/` 配下の画面専用 Svelte component
- [shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)
- [AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte)
- 変更内容は表示文言、表示構造、CSS、画面専用部品分割に限定する。

## 禁止範囲

- backend、Wails gateway、controller、usecase、store、gateway contract は変更しない。
- 保存先、ログ出力、secret、外部入力の扱いは変更しない。
- プロダクトテスト、検証データ、snapshot、test helper は原則変更しない。
- 人間UIレビュー後の `テスト修正証跡` に限り、旧 `Input Review` 表示名を期待する frontend 単体テストだけを変更してよい。
- 新しい入力形式、登録仕様、再構築仕様、ジョブ作成導線は追加しない。

## テスト修正許可範囲

- [InputReviewPage.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-input/InputReviewPage.test.ts)
- [translation-input-app.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/translation-input-app.test.ts)
- 変更内容は `Input Review` 旧見出し期待をデータロード画面の表示へ追従する範囲に限定する。

## UX適用方針

- 画面目的は JSON を読み込み、登録結果と再構築可否を確認するデータロードに固定する。
- 主要CTAはロード対象の選択と登録に絞り、再構築は選択済みデータの補助操作として分離する。
- 情報階層は `ロード準備`、`読み込み済みデータ`、`選択データの内容`、`問題と再構築` の順にする。
- UI文言は内部状態名を避け、利用者が次に行う作業で読める日本語へ寄せる。

## コンポーネント化方針

- 画面専用部品は `frontend/src/ui/screens/translation-input/` 配下に置く。
- `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail` 程度の意味単位を優先する。
- UI Component は `ViewModel` の表示値と callback props だけを受ける。
- component から controller、store、gateway、generated binding を直接参照しない。

## 人間UIレビュー観点

- 翻訳管理タブと画面見出しが `Input Review` ではなくデータロードの概念で読める。
- 画面上部で接続状態、作業状態、エラー、次操作が判断できる。
- JSON 選択、登録、選び直しが主作業としてまとまって見える。
- 一覧と詳細が、読み込み済みデータの確認画面として自然に読める。
- cache 再構築は主作業ではなく、選択済みデータの補助操作として分かる。
