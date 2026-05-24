# Frontend Implementation Handoff: frontend-processing-target-panel

- `skill`: implement-lane
- `target_agent`: frontend_implementer
- `target_skill`: implement-frontend
- `task_id`: translation-job-step-target-list-panel
- `source_scope`: `./implementation-scope.md`
- `storybook_review_state`: 人間レビュー前

## 満たされた依存対象

- `task 枠`: `task-frame.md`
- `詳細仕様差分`: `detail-spec-diff.md`
- `画面設計差分`: `screen-design-diff.*.md`
- `設計差分図`: `design-diff.translation-job-step-target-list-panel.md`
- `人間設計レビュー`: 2026-05-23 に人間が `approved` として承認
- `実装範囲`: `implementation-scope.md`

## 読むファイル

- `.codex/skills/implement-frontend/SKILL.md`
- `docs/coding-guidelines-frontend.md`
- `docs/lint-policy.md`
- `docs/UX-standard.md`
- `docs/architecture.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/detail-spec-diff.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/screen-design-diff.job-run.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/screen-design-diff.term-translation-phase.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/screen-design-diff.persona-generation-phase.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/screen-design-diff.body-translation-phase.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/screen-design-diff.translation-complete.md`
- `docs/exec-plans/active/translation-job-step-target-list-panel/implementation-scope.md`

## 実装対象

- `ProcessingTargetListPanel` を `job-run` 共通部品として追加する。
- `JobRunPage.svelte` で、選択ジョブ概要の下、`PhaseHost` の上に共通パネルを配置する。
- 現在段階に応じた段階名、処理対象名、処理対象詳細を表示する。
- 50 件程度を既定ページサイズとして扱い、現在ページの表示範囲だけを描画する。
- `前へ` と `次へ` でページを切り替え、1 ページ目と最終ページでは該当操作を無効にする。

## 変更してよい frontend ファイル

- `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.svelte`
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- `frontend/src/ui/screens/job-run/job-run-shell-props.ts`
- `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
- `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`

## Storybook 確認対象

- レビュー分類: `Review/Changed Screens/Job Run/ProcessingTargetListPanel`
- 通常分類: `Screens/Job Run/ProcessingTargetListPanel`
- 変更部品: `ProcessingTargetListPanel`
- 追加状態: 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果の確認、1 ページ目、最終ページ、長い表示文言
- `fixture`: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
- 関連資源: `job-run-shell-props.ts`, `ProcessingTargetListPanel.svelte`, `ProcessingTargetListPanel.stories.ts`

## 禁止事項

- backend、Wails bridge、DTO、gateway、repository、SQLite schema を変更しない。
- `docs/` 正本、`.codex/`、`.codex/skills`、`.codex/agents`、`plan.md` を変更しない。
- 翻訳ジョブ実行、開始、一時停止、再開、再試行、取り消し、出力管理への導線の既存動作を変更しない。
- 処理対象一覧のために新しい永続化、API、実データ取得経路を追加しない。
- 未承認のフィルタ、検索、並べ替え、一括操作、空状態、エラー状態を追加しない。
- Storybook review のための fake gateway や長寿命 mock API を追加しない。

## 検証コマンド

- `npm --prefix frontend run test -- ProcessingTargetListPanel`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

## 返却に含めること

- 変更ファイル
- 追加部品と追加状態
- Storybook のレビュー分類、通常分類、現在分類
- `fixture` と関連資源
- 検証結果と未実行理由
- 実画面確認の根拠または未取得理由
- 残留リスク
