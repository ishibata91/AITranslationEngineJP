# 画面設計差分: `translation-job-setup`

- `skill`: design-bundle
- `status`: approved
- `screen_id`: `translation-job-setup`
- `source_screen_design`: `docs/screen-design/screens/translation-job-setup.md`
- `source_plan`: `./plan.md`
- `source_storybook_review`: `./storybook-review-loop.md`

## 画面差分

### [2] 作成前設定領域

- `変更種別`: 変更

概要:
入力データの選択と作成確認を表示する領域にする。

表示内容:
- 入力データ領域
- ジョブ作成固定フッター

### [3] 入力データ領域

- `変更種別`: 変更

概要:
共通のファイル入力パネルでジョブ作成に使う入力データを表示する領域にする。

表示内容:
- `入力データ`
- 入力データ名
- 出自
- 翻訳レコード件数
- 登録日時
- 選択状態
- 既存 job 状態

依存部品:
- `FileImportPanel`: ファイル選択、選択ファイル名、補助情報を表示する。

### [7] ジョブ作成固定フッター

- `変更種別`: 変更

概要:
入力データの選択状態をもとに、翻訳ジョブ作成可否と作成操作を表示する固定フッターにする。

表示内容:
- `ジョブの作成確認`
- 入力データの確認説明
- 不足理由
- 作成に必要な確認状態
- `入力データの確認へ戻る`
- `単語翻訳へ進む`

操作:
- `入力データの確認へ戻る` を押す。
- `単語翻訳へ進む` を押す。

結果:
- 入力データの確認画面へ戻る。
- 翻訳ジョブを作成し、単語翻訳へ進む。

## 根拠

- `docs/screen-design/screens/translation-job-setup.md` は、作成前設定領域、入力データ領域、ジョブ作成固定フッターを持つ。
- `storybook-review-loop.md` の `入力ファイル UI 共通化` と `ジョブセットアップの JSON 入力単独化` は、人間レビュー承認済みの画面仕様を示す。
- `storybook-review-loop.md` の承認状態は `approved` である。

## 未決

- なし
