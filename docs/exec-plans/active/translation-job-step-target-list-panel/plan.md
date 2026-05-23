# 翻訳ジョブステップ処理対象一覧表示パネル

## 目的

翻訳ジョブ実行中に、利用者が各段階の処理対象を同じ場所で確認できるようにする。
表示部品は各フェーズ画面へ個別実装せず、`job-run` 共通部品として再利用する。

## 成果物依存表

| 成果物 | 状態 | 根拠 |
| --- | --- | --- |
| `task 枠` | 完了 | `task-frame.md` |
| `branch 準備` | 完了 | `codex/translation-job-step-target-list-panel` を作成済み |
| `詳細仕様差分` | 完了 | `detail-spec-diff.md` |
| `画面設計差分` | 完了 | `screen-design-diff.job-run.md`, `screen-design-diff.term-translation-phase.md`, `screen-design-diff.persona-generation-phase.md`, `screen-design-diff.body-translation-phase.md`, `screen-design-diff.translation-complete.md` |
| `設計差分図` | 完了 | `design-diff.translation-job-step-target-list-panel.md` |
| `人間設計レビュー` | 完了 | 2026-05-23 に人間が `approved` として承認 |
| `実装範囲` | 完了 | `implementation-scope.md` |
| `実装引き継ぎ入力` | 完了 | `implementation-handoff.frontend-processing-target-panel.md` |
| `frontend 実装` | 完了 | `ProcessingTargetListPanel` と Storybook review story を追加 |
| `Storybook人間レビュー依頼` | 完了 | この `plan.md` の `Storybook 人間レビュー依頼` |
| `frontend 実装後人間レビュー` | 停止中 | 人間の Storybook 確認待ち |
| `観測ログ追加` | 完了 | 追加不要。runtime 分岐と外部境界がないため |
| `最終検証` | 完了 | frontend-local と Storybook build が通過 |
| `実装後ブラウザ確認` | 完了 | Storybook 初期表示と Playwright ページ操作確認が通過 |

## 判断

- この task は UI を追加するため、frontend 実装と Storybook 人間レビューが必要である。
- 処理対象一覧表示パネルは、単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果確認の上に共通表示する。
- 処理対象一覧表示パネルは、選択中ジョブの入力、段階、処理対象名、処理対象詳細を示す。
- 処理対象一覧表示パネルは、50 件程度を既定ページサイズとして扱う。
- 処理対象一覧表示パネルは、数万件レベルの処理対象でも現在ページの表示範囲だけを画面要素にする。
- この task の表示パネルは、処理対象の確認位置を統一するために追加する。
- `screen-design-diff.job-run.md` は共通配置を扱う。
- フェーズ別の画面設計差分は各 `screen-design-diff.<screen-id>.md` が扱う。

## Storybook 人間レビュー依頼

- レビュー分類: `Review/Changed Screens/Job Run/ProcessingTargetListPanel`
- 通常分類: `Screens/Job Run/ProcessingTargetListPanel`
- 現在分類: `Review/Changed Screens/Job Run/ProcessingTargetListPanel`
- Storybook URL: `http://localhost:6008/?path=/story/review-changed-screens-job-run-processingtargetlistpanel--term-translation`
- Storybook 起動 command: `npm --prefix frontend run storybook`
- 変更部品: `ProcessingTargetListPanel`
- 追加状態: 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果の確認、1 ページ目、最終ページ、長い表示文言。
- fixture: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts` の `processingTargetListPanelFixtures`
- 関連資源: `frontend/src/ui/screens/job-run/job-run-shell-props.ts`, `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.svelte`, `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- Storybook 検証結果: `npm --prefix frontend run build-storybook` は通過。Vite chunk size warning は既存系の警告として扱う。
- Codex 内蔵ブラウザのコメント受付条件: コメント本文、対象 story、対象 selector、frame URL、marker screenshot を 1 件ずつ記録する。

## frontend 実装結果

- 追加: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.svelte`
- 変更: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- 変更: `frontend/src/ui/screens/job-run/job-run-shell-props.ts`
- 変更: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
- 追加: `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- 追加: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`
- 保護前状態: frontend 実装後人間レビューが未完了のため、合意済み frontend 保護対象は未固定。

## 観測ログ追加

- 結果: 追加不要。
- 理由: 変更は Storybook fixture と UI 表示部品に閉じており、実行後に消える runtime 分岐、外部境界の失敗分類、永続化、Wails bridge を扱わない。
- 変更ファイル: なし。

## 最終検証

- `npm --prefix frontend run test -- ProcessingTargetListPanel`: pass。1 file、4 tests passed。
- `npm --prefix frontend run build-storybook`: pass。Storybook build completed successfully。Vite chunk size warning あり。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。56 files、515 tests passed。
- coverage 値: 未測定。
- issue 数: 未測定。
- system test 件数: frontend-local は system test ではないため該当なし。

## 実装後ブラウザ確認

- 確認 URL: `http://localhost:6008/iframe.html?id=review-changed-screens-job-run-processingtargetlistpanel--term-translation&viewMode=story`
- 起動状態: `npm --prefix frontend run storybook` で `http://localhost:6008/` が起動済み。
- 初期表示: `処理対象一覧`、`単語翻訳`、`用語候補 001`、`用語候補 050`、`1-50 / 125 件` を確認。
- 非表示確認: 初期表示で `用語候補 051` は表示されないことを確認。
- 操作確認: Playwright で `次へ` 押下後に `51-100 / 125 件`、`用語候補 051`、`用語候補 100` を確認し、`用語候補 001` が表示されないことを確認。
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel/playwright-after-next.png`
- agent-browser 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel/` と `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel/retry-next/`
- 注意: agent-browser では `次へ` 後の DOM 変化を観測できなかった。Playwright の直接確認と単体テストではページ切替が通過した。
- 異常: Storybook 起動時に `EMFILE` watcher warning と `.storybook/settings.json` の `EPERM` warning が出た。Storybook は起動済み。

## 人間設計レビュー

- 結果: 承認
- 記録日: 2026-05-23
- 入力: `approved`
- 対象: `detail-spec-diff.md`, `screen-design-diff.*.md`, `design-diff.translation-job-step-target-list-panel.md`

## 実装引き継ぎ入力

- 対象 agent: `frontend_implementer`
- 対象 skill: `implement-frontend`
- 入力: `implementation-handoff.frontend-processing-target-panel.md`
- 依存対象: `implementation-scope.md`
- Storybook レビュー状態: 人間レビュー前

## 停止理由

frontend 実装後人間レビューが未完了である。
`implement-lane` の規約により、UI がある task では人間の Storybook 確認と承認なしに合意済み frontend 保護、後続完了、作業 commit へ進めない。
