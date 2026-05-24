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
| `実装範囲` | 完了 | `implementation-scope.md`。2026-05-24 の人間指摘後に backend / integration 必須として再作成済み |
| `実装引き継ぎ入力` | 完了 | `implementation-scope.md` の `Ready Waves` と `Handoffs` |
| `frontend 実装` | 完了 | controlled UI と production path 接続を完了 |
| `Storybookレビューループ入力確認` | 完了 | この `plan.md` の `Storybook 人間レビュー依頼` |
| `Storybookレビューループ完了証跡` | 完了 | `storybook-review-loop.md`。承認状態は `approved` |
| `frontend 実装後人間レビュー` | 完了 | 2026-05-24 に人間が Storybook フロント実装を承認 |
| `Storybook後画面設計差分整合` | 完了 | `designer` が `screen-design-diff.*.md` へ反映。未決 0 件 |
| `合意済みfrontend保護` | 完了 | この `plan.md` の `合意済み frontend 保護` |
| `backend 実装` | 完了 | `backend-processing-target-list-read-model` を実装 |
| `統合境界実装` | 完了 | `integration-processing-target-list-public-seam` を実装 |
| `単体テスト` | 完了 | 承認済み UI 追従 test と phase targeted tests が通過 |
| `シナリオテスト` | 完了 | `frontend-local`、`backend-local`、Storybook build が通過 |
| `観測ログ追加` | 完了 | backend の Wails boundary と repository failure へ `slog` を追加 |
| `最終検証` | 完了 | `frontend-local`、`backend-local`、Storybook build、`git diff --check` が通過 |
| `実装後ブラウザ確認` | 完了 | サンドボックス外 Wails dev で production `AppShell` と job-run 空状態を確認 |
| `正本化判断` | 完了 | 追加 docs 正本化は不要 |
| `詳細仕様正本反映` | 旧範囲完了 | scope 再作成前の docs 正本化は完了。追加反映なし |
| `作業 commit` | 完了 | `768d09d` |
| `マージ準備入力` | 完了 | `merge-ready.md` |

## 2026-05-24 scope 再開記録

- 人間指摘: `バックエンド整合，テスト修正やってなくない？`、`本当にバックエンド不要？`
- 再確認結果: `implementation-scope.md` は backend / integration / frontend production path を必要と判定した。
- 旧完了証跡: `最終検証`、`実装後ブラウザ確認`、`マージ準備入力` は frontend-only scope の証跡として残す。
- 現在の進行単位: `Ready Waves` に従い、`wave-1` から `wave-5` へ順に進める。
- 次の着手可能成果物: `backend-processing-target-list-read-model`。

## wave-1 実装結果

- `frontend-processing-target-list-controlled-ui`: controlled page state、検索語、件数、ページ操作 callback を frontend 部品と phase panel へ追加した。
- `unit-tests-translation-input-review-approved-ui`: 翻訳入力画面の unit test を承認済み UI へ追従した。
- `unit-tests-job-setup-approved-ui`: ジョブセットアップ画面の unit test を承認済み UI へ追従した。
- 検証: `python3 scripts/harness/run.py --suite frontend-local` は pass。56 files、523 tests passed。
- 残留確認: `screen-design-diff.translation-input-review.md` 由来の文言 `翻訳設定へ進む` と、現在 product / test の `単語翻訳へ進む` は差分確認が必要である。

## wave-2 backend 実装結果

- `backend-processing-target-list-read-model`: 処理対象一覧の read model 用 repository、service、usecase contract、usecase を追加した。
- request: `jobId`、`phase`、`page`、`pageSize`、`searchQuery` を受け取る。
- response: `items`、`page`、`pageSize`、`totalCount`、`searchQuery` を返す。
- phase: `term_translation`、`persona_generation`、`body_translation`、`translation_complete` を扱う。
- 検証: `python3 scripts/harness/run.py --suite backend-local` は pass。

## wave-3 統合境界実装結果

- `integration-processing-target-list-public-seam`: phase gateway contract へ `getProcessingTargetList` を追加した。
- Wails: `ProcessingTargetController` と `GetProcessingTargetList` を追加した。
- DTO: frontend gateway DTO と backend Wails DTO を追加した。
- 保持方針: `TranslationJobManagementJobRunTarget` は軽量のまま維持した。
- 検証: `python3 scripts/harness/run.py --suite backend-local` と frontend phase targeted tests は pass。

## wave-4 frontend production path 実装結果

- `frontend-processing-target-list-production-path`: production gateway から処理対象一覧を取得して `JobRunPage` へ渡す経路を追加した。
- Storybook: `processingTargetItemsByPhase` がある場合は fixture 表示を優先する。
- 状態分離: `body_translation` と `translation_complete` の page state は phase key ごとに分離した。
- 保護: 承認済み画面の layout、文言、style は変更していない。
- 検証: `python3 scripts/harness/run.py --suite frontend-local` は pass。56 files、523 tests passed。

## 観測ログ追加 再判定

- 結果: 追加済み。
- backend Wails boundary: usecase 未設定と usecase 失敗を `processing_target_list_boundary_failed` として記録する。
- backend repository: unsupported phase、count failed、list failed を `processing_target_list_repository_failed` として記録する。
- 保護: 検索語、DTO 全体、secret、大量本文はログへ出さない。
- 検証: `python3 scripts/harness/run.py --suite backend-local` と `git diff --check` は pass。

## wave-5 最終検証

- `python3 scripts/harness/run.py --suite frontend-local`: pass。56 files、523 tests passed。
- `npm --prefix frontend run build-storybook`: pass。Vite chunk size warning あり。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `git diff --check`: pass。
- `wails build -clean`: fail。Wails CLI は `exit status 1` を返す。
- 追加確認: Wails が出した Go build 相当 command は直接実行で pass。

## production path 実装後ブラウザ確認

- 確認 URL: `http://localhost:34115`
- 結果: 完了。
- 停止理由: `npm run dev:wails:agent-browser` は `tmp/logs/wails-dev.log` に `Build error - exit status 1` を出し、`34115` が listen しない。
- 人間確認: 手動の Wails dev 起動は通るため、sandbox 起因の可能性が高い。
- Codex 再確認: `lsof -nP -iTCP:34115 -sTCP:LISTEN` は listen なし、`curl -I --max-time 5 http://localhost:34115` は connection refused。
- サンドボックス外起動: `npm run dev:wails:agent-browser` は `34115` listen と HTTP 到達を確認済み。
- 確認結果: `http://localhost:34115/#translation-management/job-run` へ到達し、`ジョブ #7`、`単語翻訳`、`処理対象`、検索欄、ページング、`処理対象がありません` を確認。
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/snapshot.txt`
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/job-run-errors.txt`
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/job-run-confirmed.png`
- 異常: `agent-browser errors` では新しい console error は確認されなかった。
- 未確認理由: なし。
- 切り分け: 通常 sandbox 起動だけ失敗し、サンドボックス外起動では production 画面確認が通過した。

## 追加正本化判断

- 結果: 追加 docs 正本化は不要。
- 理由: 今回の再開後実装は、承認済み `implementation-scope.md` の backend、integration、production path 接続を実装した。
- 理由: 恒久仕様の追加は、既に正本化済みの `detail-spec-diff.md` と `screen-design-diff.*.md` の範囲を超えていない。
- 対象外: Wails dev 起動不能は環境確認結果であり、詳細仕様正本へ入れない。

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
- 状態: 完了。
- commit 対象差分: `docs/exec-plans/active/translation-job-step-target-list-panel/`
- commit 対象差分: 処理対象一覧 read model の backend 実装。
- commit 対象差分: 処理対象一覧 Wails / gateway 境界の integration 実装。
- commit 対象差分: 処理対象一覧 production path の frontend 実装と追従 test。
- 除外差分: 作業開始時点または別作業の `.codex`、`AGENTS.md`、template 差分は今回 commit 対象に含めない。
- commit hash: `768d09d`
- commit message: `Implement processing target production path`

## マージ準備入力

- active plan folder: `docs/exec-plans/active/translation-job-step-target-list-panel/`
- source branch: `codex/translation-job-step-target-list-panel`
- target branch: `master`
- 状態: 完了。
- merge ready file: `merge-ready.md`
- commit hash: `768d09d`
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は pass。56 files、523 tests passed。
- 検証結果: `npm --prefix frontend run build-storybook` は pass。Vite chunk size warning あり。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は pass。
- 検証結果: `git diff --check` は pass。
- 実装後ブラウザ確認結果: pass。サンドボックス外 Wails dev で `#translation-management/job-run` へ到達した。
- 実装後ブラウザ確認結果: `ジョブ #7`、`単語翻訳`、`処理対象`、検索欄、ページング、`処理対象がありません` を確認した。
- 残留リスク: 通常 sandbox の `npm run dev:wails:agent-browser` は `Build error - exit status 1` を返す。
- 切り分け: サンドボックス外 Wails dev では `34115` listen と HTTP 到達が通過した。

## 観測ログ追加

- 結果: 追加済み。
- 追加先: backend Wails boundary。
- 追加先: processing target read model repository。
- event: `processing_target_list_boundary_failed`
- event: `processing_target_list_repository_failed`
- 保護: 検索語、DTO 全体、secret、大量本文はログへ出さない。

## 最終検証

- `python3 scripts/harness/run.py --suite frontend-local`: pass。56 files、523 tests passed。
- `npm --prefix frontend run build-storybook`: pass。Vite chunk size warning あり。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `git diff --check`: pass。
- `wails build -clean`: fail。Wails CLI は `exit status 1` を返す。
- 直接 Go build: pass。Wails dev と同等 tag の backend binary は生成できる。
- coverage 値: 未測定。
- issue 数: 未測定。
- system test 件数: frontend-local は system test ではないため該当なし。

## 実装後ブラウザ確認

- 確認 URL: `http://localhost:34115`
- 結果: 完了。
- 起動状態: サンドボックス外 `npm run dev:wails:agent-browser` で `34115` listen と HTTP 到達を確認。
- 確認結果: `http://localhost:34115/#translation-management/job-run` へ到達した。
- 確認結果: `ジョブ #7`、`単語翻訳`、`処理対象`、検索欄、ページング、`処理対象がありません` を確認した。
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/snapshot.txt`
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/job-run-errors.txt`
- 証跡: `frontend/test-results/browser-confirmation/translation-job-step-target-list-panel-production-path/job-run-confirmed.png`

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

## Merge Result

- `merge_status`: `merged-to-master`
- `merge_command`: `git merge --no-ff --no-commit codex/translation-job-step-target-list-panel`
- `source_branch`: `codex/translation-job-step-target-list-panel`
- `target_branch`: `master`
- `source_branch_head`: `a3f8dd4`
- `work_commit_hash`: `768d09d`
- `conflict_resolution`: conflict なし。
- `post_merge_validation`: `git diff --check` pass。
- `post_merge_validation`: `python3 scripts/harness/run.py --suite backend-local` pass。
- `post_merge_validation`: `python3 scripts/harness/run.py --suite frontend-local` pass。56 files、523 tests passed。
- `post_merge_validation`: `npm --prefix frontend run build-storybook` pass。Vite chunk size warning あり。
- `post_merge_validation`: `python3 scripts/harness/run.py --suite structure` pass。
- `validation_note`: 初回 `frontend-local` は temp worktree の依存不足により `eslint: command not found` で fail。`npm --prefix frontend install` 後に pass。
- `validation_note`: 初回 `backend-local` は temp worktree の `frontend/dist` 不足により fail。`npm --prefix frontend run build` 後に pass。
- `completed_move`: `docs/exec-plans/active/translation-job-step-target-list-panel/` から `docs/exec-plans/completed/translation-job-step-target-list-panel/` へ移動済み。
- `merge_commit_hash`: local commit 作成後に merge lane の返却で記録する。
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: `detail-spec-diff.md` と `screen-design-diff.*.md` は既に docs 正本へ反映済み。
- `detail_spec_canonicalization`: 追加 docs 正本化は不要。
- `residual_risk`: 通常 sandbox の `npm run dev:wails:agent-browser` は `Build error - exit status 1` を返す。
- `residual_risk`: サンドボックス外 Wails dev では `34115` listen と HTTP 到達が通過した。
- `residual_risk`: browser confirmation は production job-run の空状態を確認した。処理対象 item が存在する job での非空一覧は未確認である。
- `follow_up`: なし。

## 停止理由

なし。
