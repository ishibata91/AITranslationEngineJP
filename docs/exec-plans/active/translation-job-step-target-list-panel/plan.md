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
| `画面設計差分` | 完了 | `screen-design-diff.*.md` |
| `設計差分図` | 完了 | `design-diff.translation-job-step-target-list-panel.md` |
| `人間設計レビュー` | 完了 | 2026-05-23 に人間が `approved` として承認 |
| `実装範囲` | 完了 | `implementation-scope.md` |
| `実装引き継ぎ入力` | 完了 | `implementation-handoff.frontend-processing-target-panel.md` |
| `frontend 実装` | 完了 | `ProcessingTargetListPanel` と Storybook review story を追加 |
| `Storybookレビューループ入力確認` | 完了 | この `plan.md` の `Storybook 人間レビュー依頼` |
| `Storybookレビューループ完了証跡` | 完了 | `storybook-review-loop.md`。承認状態は `approved` |
| `frontend 実装後人間レビュー` | 完了 | 2026-05-24 に人間が Storybook フロント実装を承認 |
| `Storybook後画面設計差分整合` | 完了 | `designer` が `screen-design-diff.*.md` へ反映。未決 0 件 |
| `合意済みfrontend保護` | 完了 | この `plan.md` の `合意済み frontend 保護` |
| `観測ログ追加` | 完了 | 追加不要。runtime 分岐と外部境界がないため |
| `最終検証` | 完了 | frontend-local と Storybook build が通過 |
| `実装後ブラウザ確認` | 完了 | Storybook 初期表示と Playwright ページ操作確認が通過 |
| `正本化判断` | 完了 | 人間承認済み詳細仕様差分と画面設計差分の docs 正本化が必要 |
| `詳細仕様正本反映` | 完了 | `docs_updater` が docs 正本へ反映。残留不足なし |
| `作業 commit` | 完了予定 | local commit 作成後の hash は `git log -1` で確認する |
| `マージ準備入力` | 完了 | この `plan.md` の `マージ準備入力` |

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

## Storybook レビューループ完了証跡

- 証跡: `storybook-review-loop.md`
- 承認状態: `approved`
- 人間承認入力: 2026-05-24 の `storybookフロント実装承認`
- 確定した story: `Screens/Job Run/JobRunPage`, `Screens/Job Run/TranslationCompletePage`, `Screens/Master Dictionary/MasterDictionaryPage`, `Screens/Master Persona/MasterPersonaPage`, `Screens/Translation Input/InputReviewPage`, `UI Components/ProcessingTargetListPanel`
- 通常分類へ戻した story: `Screens/Job Run/JobRunPage`, `Screens/Job Run/TranslationCompletePage`, `Screens/Master Dictionary/MasterDictionaryPage`, `Screens/Master Persona/MasterPersonaPage`, `Screens/Translation Input/InputReviewPage`, `UI Components/ProcessingTargetListPanel`
- 設計整合入力: `storybook-review-loop.md` の `変更された画面仕様` は、`screen-design-diff.*.md` へ反映済みである。

## frontend 実装後人間レビュー

- 結果: 承認
- 記録日: 2026-05-24
- 入力: `storybookフロント実装承認`
- 根拠: `storybook-review-loop.md` の承認状態 `approved`
- 注意: 承認は frontend 表示成果物への承認である。

## Storybook後画面設計差分整合

- 結果: 完了
- 担当: `designer`
- 更新成果物: `screen-design-diff.job-run.md`, `screen-design-diff.term-translation-phase.md`, `screen-design-diff.persona-generation-phase.md`, `screen-design-diff.body-translation-phase.md`, `screen-design-diff.translation-complete.md`
- 更新成果物: `screen-design-diff.master-dictionary.md`, `screen-design-diff.master-persona.md`, `screen-design-diff.translation-input-review.md`, `screen-design-diff.translation-job-setup.md`, `screen-design-diff.translation-job-management.md`
- 反映根拠: `storybook-review-loop.md` の `変更された画面仕様`
- 未反映項目: Storybook 分類、fixture 表示パターン、変更ファイル一覧、承認状態は運用記録であり、画面内容の差分本文へ入れない。
- 未決事項: 0 件
- 判断: `implement_lane` は `合意済みfrontend保護` へ進める。

## 合意済み frontend 保護

- 状態: 完了
- 承認済み画面: `JobRunPage`, `TranslationCompletePage`, `MasterDictionaryPage`, `MasterPersonaPage`, `InputReviewPage`
- 承認済み部品: `ProcessingTargetListPanel`, `ProcessingTargetListWrapper`, `FileImportPanel`, `PhaseStatusPanel`, `PhaseProgressPanel`, `AIModelSelectionCard`, `TranslationManagementStepper`
- 承認済み表示規則: `screen-design-diff.*.md` に反映した画面配置、一覧、検索、ページ操作、展開行、入力ファイル、進行状況、Storybook 通常分類。
- 確認済み Storybook 状態: `storybook-review-loop.md` の確定 story と通常分類。
- Storybook 確認資源: `storybook-review-loop.md` の `対象` と `現在状態` に記録した story、fixture、関連資源。
- Storybook 画面仕様: `screen-design-diff.*.md`
- 変更禁止範囲: `storybook-review-loop.md` の `関連資源` と `変更ファイル` に含まれる frontend 表示、文言、layout、style、story、fixture。

## frontend 実装結果

- 追加: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.svelte`
- 変更: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- 変更: `frontend/src/ui/screens/job-run/job-run-shell-props.ts`
- 変更: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
- 追加: `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- 追加: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`
- 保護状態: frontend 実装後人間レビューと Storybook 後画面設計差分整合が完了したため、合意済み frontend 保護対象は固定済みである。

## 正本化判断

- 結果: docs 正本化が必要。
- 理由: `detail-spec-diff.md` と `screen-design-diff.*.md` は、人間承認済みの恒久仕様を含む。
- 承認記録: 2026-05-23 の人間設計レビュー承認。
- 承認記録: 2026-05-24 の Storybook フロント実装承認。
- 正本化対象成果物: `detail-spec-diff.md`
- 正本化対象成果物: `screen-design-diff.job-run.md`, `screen-design-diff.term-translation-phase.md`, `screen-design-diff.persona-generation-phase.md`, `screen-design-diff.body-translation-phase.md`, `screen-design-diff.translation-complete.md`
- 正本化対象成果物: `screen-design-diff.master-dictionary.md`, `screen-design-diff.master-persona.md`, `screen-design-diff.translation-input-review.md`, `screen-design-diff.translation-job-setup.md`, `screen-design-diff.translation-job-management.md`
- 正本化先: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- 正本化先: `docs/screen-design/screens/job-run.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`, `docs/screen-design/screens/translation-complete.md`
- 正本化先: `docs/screen-design/screens/master-dictionary.md`, `docs/screen-design/screens/master-persona.md`, `docs/screen-design/screens/translation-input-review.md`, `docs/screen-design/screens/translation-job-setup.md`, `docs/screen-design/screens/translation-job-management.md`
- 起動先: `docs_updater` / `updating-docs`
- 未決事項: 0 件

## 詳細仕様正本反映

- 結果: 完了
- 担当: `docs_updater`
- 反映元: `detail-spec-diff.md`
- 反映元: `screen-design-diff.job-run.md`, `screen-design-diff.term-translation-phase.md`, `screen-design-diff.persona-generation-phase.md`, `screen-design-diff.body-translation-phase.md`, `screen-design-diff.translation-complete.md`
- 反映元: `screen-design-diff.master-dictionary.md`, `screen-design-diff.master-persona.md`, `screen-design-diff.translation-input-review.md`, `screen-design-diff.translation-job-setup.md`, `screen-design-diff.translation-job-management.md`
- 更新 docs: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- 更新 docs: `docs/screen-design/screens/job-run.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`, `docs/screen-design/screens/translation-complete.md`
- 更新 docs: `docs/screen-design/screens/master-dictionary.md`, `docs/screen-design/screens/master-persona.md`, `docs/screen-design/screens/translation-input-review.md`, `docs/screen-design/screens/translation-job-setup.md`, `docs/screen-design/screens/translation-job-management.md`
- 未反映項目: Storybook 分類、fixture、変更ファイル一覧、承認状態は運用記録のため docs 正本へ入れない。
- 未反映項目: `implementation-scope.md`、実装手順、テスト手順、agent handoff は docs 正本化対象外である。
- 検証: `git diff --check -- docs/detail-specs docs/screen-design/screens` は pass。
- 検証: `python3 scripts/harness/run.py --suite structure` は pass。
- 残留不足: なし。

## 作業 commit

- 対象 branch: `codex/translation-job-step-target-list-panel`
- commit 対象差分: `docs/exec-plans/active/translation-job-step-target-list-panel/`
- commit 対象差分: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- commit 対象差分: `docs/screen-design/screens/job-run.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`, `docs/screen-design/screens/translation-complete.md`
- commit 対象差分: `docs/screen-design/screens/master-dictionary.md`, `docs/screen-design/screens/master-persona.md`, `docs/screen-design/screens/translation-input-review.md`, `docs/screen-design/screens/translation-job-setup.md`, `docs/screen-design/screens/translation-job-management.md`
- commit hash: local commit 作成後に `git log -1` で確認する。

## マージ準備入力

- active plan folder: `docs/exec-plans/active/translation-job-step-target-list-panel/`
- source branch: `codex/translation-job-step-target-list-panel`
- target branch: `master`
- commit hash: local commit 作成後に `git log -1` で確認する。
- 検証結果: `npm --prefix frontend run test -- ProcessingTargetListPanel` は pass。1 file、4 tests passed。
- 検証結果: `npm --prefix frontend run build-storybook` は pass。Storybook build completed successfully。Vite chunk size warning あり。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は pass。56 files、515 tests passed。
- 検証結果: `git diff --check -- docs/detail-specs docs/screen-design/screens docs/exec-plans/active/translation-job-step-target-list-panel` は pass。
- 検証結果: `python3 scripts/harness/run.py --suite structure` は pass。
- 実装後ブラウザ確認結果: Storybook 初期表示と Playwright ページ操作確認が通過。
- 残留リスク: agent-browser では `次へ` 後の DOM 変化を観測できなかった。Playwright の直接確認と単体テストではページ切替が通過している。
- 残留リスク: Storybook 起動時の `EMFILE` watcher warning と `.storybook/settings.json` の `EPERM` warning は、Storybook 起動済みの状態で発生した環境警告として扱う。
- 除外差分: 作業開始時点から存在した `.codex/README.md`, `.codex/agents/implement_lane.toml`, `.codex/skills/implement-lane/SKILL.md`, `docs/exec-plans/templates/task-folder/plan.md` は今回の commit 対象に含めない。

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

なし。
