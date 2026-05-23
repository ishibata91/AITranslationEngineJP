# 画面設計差分: `translation-complete`

- `skill`: design-bundle
- `status`: ready-for-human-review
- `screen_id`: `translation-complete`
- `source_screen_design`: `docs/screen-design/screens/translation-complete.md`
- `source_plan`: `./plan.md`

## 画面差分

### [8] 処理対象一覧

- `変更種別`: 追加

概要:
翻訳完了確認で確認する実体を確認する領域。

配置:
`job-run` の選択ジョブ概要の下、翻訳完了確認画面の完了概要の上に配置する。

表示内容:
- `処理対象`
- 現在段階名: `翻訳結果の確認`
- 処理対象名: `翻訳項目単位の訳文`
- 処理対象詳細: 翻訳項目単位の訳文は、本文翻訳で保持された訳文として出力管理へ進む前に確認するもの。
- 処理対象ページング: 50 件程度を現在ページの表示範囲として扱う。

依存情報:
- 表示条件: 選択ジョブがあり、現在段階が翻訳完了確認である。
- 有効条件: 翻訳完了確認で確認する実体を表示できる。
- データ種別: 表示用の処理対象情報、ページング状態。
- 画面非機能要件: 数万件レベルの処理対象でも、画面要素は現在ページの表示範囲に限定し、ページ切替操作を維持する。

操作:
- なし。

結果:
- なし。

セレクタ属性:
- `aria-label`: `処理対象一覧`

依存部品:
- `ProcessingTargetListPanel`: 各段階で同じ構造の処理対象一覧を表示する。

## 根拠

- `docs/screen-design/screens/translation-complete.md` は、完了概要、原文と訳文、翻訳結果ページング、翻訳完了後の次の作業を持つ。
- `docs/detail-specs/body-translation-phase.md` は、本文翻訳フェーズ完了時点で翻訳ジョブ全体が `Completed` になり、成果物出力条件が成立することを扱う。
- `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte` は、完了概要と翻訳結果一覧を表示している。
