# 画面設計差分: `persona-generation-phase`

- `skill`: design-bundle
- `status`: ready-for-human-review
- `screen_id`: `persona-generation-phase`
- `source_screen_design`: `docs/screen-design/screens/persona-generation-phase.md`
- `source_plan`: `./plan.md`

## 画面差分

### E-04 処理対象一覧

- `変更種別`: 追加

概要:
NPC ペルソナ生成段階で生成する実体の入力を確認する領域。

配置:
`job-run` の選択ジョブ概要の下、NPC ペルソナ生成段階の状態概要の上に配置する。

表示内容:
- `処理対象`
- 現在段階名: `NPC ペルソナ生成`
- 処理対象名: `NPC ごとのペルソナ生成入力`
- 処理対象詳細: NPC ごとのペルソナ生成入力は、NPC の原文発話、NPC 属性、会話文脈、共通ペルソナ参照からペルソナ参照情報を作るもの。
- 処理対象ページング: 50 件程度を現在ページの表示範囲として扱う。

依存情報:
- 表示条件: 選択ジョブがあり、現在段階が NPC ペルソナ生成である。
- 有効条件: NPC ペルソナ生成段階で生成に使う実体を表示できる。
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

- `docs/screen-design/screens/persona-generation-phase.md` は、状態概要、進行状況兼操作、AI モデル選択を持つ。
- `docs/detail-specs/persona-generation-phase.md` は、生成対象を NPC レコード、翻訳対象項目、会話文脈、共通ペルソナ参照、ペルソナ参照情報で構成する。
- `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte` は、NPC ペルソナ生成段階の状態概要を画面先頭に表示している。
