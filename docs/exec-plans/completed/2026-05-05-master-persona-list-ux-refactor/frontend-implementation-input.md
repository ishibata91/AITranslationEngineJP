# frontend 実装起動入力

## 対象成果物

`frontend 実装`

## 満たされた依存対象

- `task 枠`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-list-ux-refactor/task-frame.md)

## 実装 skill

- [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md)

## 変更対象

- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)

## 実装要求

- `MasterPersonaPage.svelte` の上側ヒーロー説明文へ、一覧と詳細を同じ画面から確認できる意味のページ説明文を統合する。
- 下側ヒーローに相当する重複説明ブロックが `PersonaReviewPanel.svelte` 側にある場合は削る。
- ペルソナ一覧の `.list-row` の行間、gap、padding、min-height を詰める。
- `.list-row strong` の文字色を `var(--text)` 系へ寄せ、黒表示を避ける。
- 検索窓とプラグインフィルタは同じ高さ、同じ入力面の見え方、同じ行内整列にする。

## 停止条件

- backend、Wails gateway、controller、usecase、store の変更が必要になった場合は停止する。
- プロダクトテスト変更が必要になった場合は停止する。
- 表示項目追加または操作追加が必要になった場合は停止する。
- 実画面確認に必要な環境が起動できない場合は、未取得理由を返す。

## 検証

- `python3 scripts/harness/run.py --suite frontend-local`
- 可能なら `agent-browser` で `http://localhost:34115` のマスターペルソナ画面を確認する。

