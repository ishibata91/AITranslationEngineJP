# 画面設計差分: `translation-job-management`

- `skill`: design-bundle
- `status`: approved
- `screen_id`: `translation-job-management`
- `source_screen_design`: `docs/screen-design/screens/translation-job-management.md`
- `source_plan`: `./plan.md`
- `source_storybook_review`: `./storybook-review-loop.md`

## 画面差分

### [13] 全体進捗パネル

- `変更種別`: 追加

概要:
翻訳管理内の作業位置を細く表示する領域。

表示内容:
- 番号
- 現在状態
- 作業名

依存情報:
- 表示条件: 翻訳管理内の作業順を表示する。
- 有効条件: なし。
- データ種別: 作業番号、作業状態、作業名。

操作:
- なし。

結果:
- なし。

依存部品:
- `TranslationManagementStepper`: 翻訳管理内の作業順を表示する。

## 根拠

- `docs/screen-design/screens/translation-job-management.md` は、未完了ジョブ一覧とジョブ選択を扱う。
- `storybook-review-loop.md` の `翻訳管理の細い全体進捗パネル` は、人間レビュー承認済みの画面仕様を示す。
- `storybook-review-loop.md` の承認状態は `approved` である。

## 未決

- なし
