# 画面設計差分: `term-translation-phase`

- `skill`: design-bundle
- `status`: ready-for-human-review
- `screen_id`: `term-translation-phase`
- `source_screen_design`: `docs/screen-design/screens/term-translation-phase.md`
- `source_plan`: `./plan.md`

## 画面差分

### E-04 処理対象一覧

- `変更種別`: 追加

概要:
単語翻訳段階で処理する実体を確認する領域。

配置:
`job-run` の選択ジョブ概要の下、単語翻訳段階の状態概要の上に配置する。

表示内容:
- `処理対象`
- 現在段階名: `単語翻訳`
- 処理対象名: `共通辞書対象外の用語と固有名詞`
- 処理対象詳細: 共通辞書対象外の用語と固有名詞は、AIサービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存するもの。
- 処理対象ページング: 50 件程度を現在ページの表示範囲として扱う。

依存情報:
- 表示条件: 選択ジョブがあり、現在段階が単語翻訳である。
- 有効条件: 単語翻訳段階で処理する実体を表示できる。
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

- `docs/screen-design/screens/term-translation-phase.md` は、状態概要、進行状況兼操作、AI モデル選択を持つ。
- `docs/detail-specs/term-translation-phase.md` は、共通辞書対象外の用語と固有名詞を AI 翻訳対象として扱う。
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` は、単語翻訳段階の状態概要を画面先頭に表示している。
