# frontend 情報保持再実装入力

## 対象成果物

`frontend 実装`

## 満たされた依存対象

- `task 枠`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-list-ux-refactor/task-frame.md)
- `人間UIレビュー`: [human-ui-review.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-list-ux-refactor/human-ui-review.md)

## 実装要求

- `生成準備 / マスターペルソナ作成` の画面内カードは戻さない。
- 削除された `viewModel.runStatus.runState` は、カードではない軽量な状態表示として復活させる。
- 上部ヒーロー説明文は `shell-state.ts` の route lead に保持する。
- ペルソナ一覧の件数、ページ、現在範囲など、一覧利用に必要な情報は削りすぎない。
- `Test NPC A を選択中です。` と `更新と削除を行えます` の可視表示は戻さない。
- 情報保持のために余白を増やしすぎない。

## 起動状態

- frontend agent `019df61d-f8f5-7e70-875b-45cbf5df9043` を再開した。
- 人間UIレビューが通るまで frontend agent を close しない。

