# frontend 再実装起動入力

## 対象成果物

`frontend 実装`

## 満たされた依存対象

- `task 枠`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-list-ux-refactor/task-frame.md)
- `人間UIレビュー`: [human-ui-review.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-list-ux-refactor/human-ui-review.md)

## 実装 skill

- [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md)

## 変更対象

- [shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)
- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)

## 実装要求

- `shell-state.ts` の `master-persona` route `lead` を、画面内カードにあった説明文へ置き換える。
- `MasterPersonaPage.svelte` の `生成準備 / マスターペルソナ作成` カードを削除する。
- 削除したカード用の不要CSSも削除する。
- `PersonaReviewPanel.svelte` から `selectedSummary` 表示を削除する。
- `PersonaReviewPanel.svelte` から `detailLockText` と `detailStatusText` の表示を削除する。
- ペルソナ一覧行をさらに細くする。目安は `min-height` を 34px 前後、上下 padding を 6px 前後、gap を 4px 前後にする。
- panel と一覧周辺の余白を詰める。説明削除後に header、filter、list の間へ広い空白を残さない。

## 停止条件

- backend、Wails gateway、controller、usecase、store の変更が必要になった場合は停止する。
- プロダクトテスト変更が必要になった場合は停止する。
- 表示項目追加または操作追加が必要になった場合は停止する。

## 検証

- `python3 scripts/harness/run.py --suite frontend-local`
- 可能なら `agent-browser` で `http://localhost:34115` のマスターペルソナ画面を確認する。

