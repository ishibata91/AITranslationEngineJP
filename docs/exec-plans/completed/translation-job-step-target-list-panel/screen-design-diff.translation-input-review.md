# 画面設計差分: `translation-input-review`

- `skill`: design-bundle
- `status`: approved
- `screen_id`: `translation-input-review`
- `source_screen_design`: `docs/screen-design/screens/translation-input-review.md`
- `source_plan`: `./plan.md`
- `source_storybook_review`: `./storybook-review-loop.md`

## 画面差分

### [2] ロード準備領域

- `変更種別`: 変更

概要:
共通のファイル入力パネルで翻訳入力 JSON の選択と登録を表示する領域にする。

表示内容:
- `ロード準備`
- `ロード対象を選ぶ`
- 選択 JSON 名
- `この JSON を登録`
- `選び直す`
- 登録状態

依存部品:
- `FileImportPanel`: ファイル選択、選択ファイル名、補助情報を表示する。

### [3] 読み込み済みデータ一覧

- `変更種別`: 変更

概要:
読み込み済みデータの登録結果と次の作業に必要な情報を一覧で確認する領域にする。

表示内容:
- `読み込み済みデータ`
- 空状態文
- ファイル名
- 登録状態
- 登録結果
- 読み込み日時
- 問題区分
- 選択状態

操作:
- 読み込み済みデータの行を選択する。

結果:
- 選択した読み込み済みデータを次の作業フッターの対象にする。

### [7] 次の作業フッター

- `変更種別`: 変更

概要:
選択済みの読み込み済みデータから翻訳設定へ進む固定フッターにする。

表示内容:
- `次の作業`
- 選択済み入力データの説明
- `翻訳設定へ進む`

依存情報:
- 表示条件: 選択データの状態が `登録済み` または `警告あり` である。
- 有効条件: 表示されている場合に操作できる。

操作:
- `翻訳設定へ進む` を押す。

結果:
- 選択した入力データで翻訳設定の確認へ進む。

## 根拠

- `docs/screen-design/screens/translation-input-review.md` は、ロード準備領域、読み込み済みデータ一覧、選択データ領域、次の作業フッターを持つ。
- `storybook-review-loop.md` の `入力ファイル UI 共通化`、`翻訳入力ロード準備の共通化`、`翻訳入力の選択データ詳細パネル廃止`、`翻訳入力の読み込み済みデータ件数表示廃止` は、人間レビュー承認済みの画面仕様を示す。
- `storybook-review-loop.md` の承認状態は `approved` である。

## 未決

- なし
