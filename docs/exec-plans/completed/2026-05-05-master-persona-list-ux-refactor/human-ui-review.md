# 人間UIレビュー

## 判定

- 結果: 差し戻し
- 対象成果物: `frontend 実装`
- 戻し先: `implementation_implementer`

## 差し戻し内容

- `生成準備 / マスターペルソナ作成` のカードは残さず削除する。
- 削除するカードの説明文は、上部ヒーローの説明文と入れ替える。
- 詳細カードの `Test NPC A を選択中です。` は不要である。
- 詳細カードの `更新と削除を行えます` が重複しているため不要である。
- ペルソナ一覧の行がまだ太すぎる。
- カード説明と一覧の間に不要な余白がある。

## 固定する修正方針

- 上部ヒーローの説明文は `frontend/src/ui/stores/shell-state.ts` の `master-persona` route lead を更新する。
- `MasterPersonaPage.svelte` の画面内ヒーローカードは削除する。
- `PersonaReviewPanel.svelte` の詳細見出し下の選択中説明と `detailLockText` / `detailStatusText` 表示を削除する。
- `PersonaReviewPanel.svelte` の一覧行、panel、section、filter、list の余白をさらに詰める。

## 二回目差し戻し

- 画面内ヒーローカード削除時に、必要な情報まで削れている。
- 削除対象は重複カードであり、状態表示や一覧利用に必要な情報は残す。
- `viewModel.runStatus.runState` はカードではない軽量な状態表示として復活させる。
- ペルソナ一覧の件数、ページ、現在範囲など、一覧利用に必要な情報は削りすぎない。
- `Test NPC A を選択中です。` と `更新と削除を行えます` の可視表示は戻さない。
- frontend agent は人間UIレビューが通るまで継続し、close しない。
