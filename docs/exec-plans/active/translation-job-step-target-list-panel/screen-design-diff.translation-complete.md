# 画面設計差分: `translation-complete`

- `skill`: design-bundle
- `status`: approved
- `screen_id`: `translation-complete`
- `source_screen_design`: `docs/screen-design/screens/translation-complete.md`
- `source_plan`: `./plan.md`
- `source_storybook_review`: `./storybook-review-loop.md`

## 画面差分

### [3] 処理対象一覧領域

- `変更種別`: 変更

概要:
翻訳完了確認で確認する実体を、処理対象一覧として表示する領域にする。

表示内容:
- `処理対象`
- 現在段階名: `翻訳結果の確認`
- 処理対象名: 翻訳項目単位の訳文
- 処理対象詳細: 本文翻訳で保持された訳文として出力管理へ進む前に確認するもの。
- テーブル見出し: `原文`、`訳文`
- 原文セル: 原文の内容
- 訳文セル: 訳文の内容
- 件数
- ページ操作

依存情報:
- 表示条件: 選択ジョブがあり、現在段階が翻訳完了確認である。
- 有効条件: 翻訳完了確認で確認する実体を表示できる。
- データ種別: 表示用の処理対象情報、ページング状態。
- 画面非機能要件: 数万件レベルの処理対象でも、画面要素は現在ページの表示範囲に限定し、ページ切替操作を維持する。

操作:
- 処理対象行を開く。
- ページを切り替える。

結果:
- 処理対象行の metadata を表示する。
- 表示する処理対象の範囲が切り替わる。

セレクタ属性:
- `aria-label`: `処理対象一覧`

依存部品:
- `ProcessingTargetListPanel`: 翻訳完了確認の処理対象一覧を表示する。

## 根拠

- `docs/screen-design/screens/translation-complete.md` は、完了概要、原文と訳文、翻訳結果ページング、翻訳完了後の次の作業を持つ。
- `docs/detail-specs/body-translation-phase.md` は、本文翻訳フェーズ完了時点で翻訳ジョブ全体が `Completed` になり、成果物出力条件が成立することを扱う。
- `storybook-review-loop.md` の `翻訳完了画面の一覧置き換え` と `処理対象一覧の 2 列表示` は、人間レビュー承認済みの画面仕様を示す。
- `storybook-review-loop.md` の承認状態は `approved` である。

## 未決

- なし
