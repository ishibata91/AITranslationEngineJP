# fix-job-run-shell-effect-untrack

## 依頼要約

`frontend/src/ui/screens/job-run/JobRunPage.svelte` の `$effect` が `setCurrentPhasePage` 関数経由で `currentPhasePage`（`$state`）を読み取り、Svelte 5 の reactivity が依存に登録されている。その結果、ユーザー操作で `setCurrentPhasePage("persona")` または `setCurrentPhasePage("complete")` を呼ぶと、`$effect` 再実行が `selectedJobTarget.currentPhase` 起点の初期 phase（"term" / "body"）へ戻してしまい、画面遷移が機能しない。`E2E-UC-048`（term → persona）と `E2E-UC-049`（body → complete）が pass しない。

## 分岐元

- 分岐元 task: `fix-phase-ai-settings-pill-update-after-model-select`（agent が発見、本 task 主旨と独立のため別出し）
- 分岐元 branch（決定保留）: `master`（前 task の merge 後を分岐元にする想定）

## 観測根拠（前 task の implementation_tester 起動で agent が確認）

- `tests/system/job-run-shell.spec.ts:254` E2E-UC-048 fail: `clickNext()` 後に `persona-generation-phase-screen` が visible にならず timeout
- `tests/system/job-run-shell.spec.ts:273` E2E-UC-049 fail: `clickBodyCompleteNext()` 後に `translation-complete-translation-complete-screen` が visible にならず timeout
- 循環フロー:
  1. `clickNext()` → `setCurrentPhasePage("persona")` → `currentPhasePage = "persona"`
  2. `currentPhasePage` の変化が `$effect` 再実行をトリガー
  3. `$effect` 内で `resolveInitialPhasePage(selectedJobTarget)` が "term"（selectedJobTarget.currentPhase が `"term_translation"` のため）
  4. `setCurrentPhasePage("term")` で persona 画面が消える

## 修正方針候補

- `setCurrentPhasePage` 内の `currentPhasePage` 読み取りを Svelte 5 `untrack()` でラップ
- または `$effect` の依存から `currentPhasePage` を除外する設計変更
- いずれも frontend ロジック層の変更で、storybook-module 範囲ではない（svelte 表示構造変更ではない）

## 影響範囲

- `frontend/src/ui/screens/job-run/JobRunPage.svelte` の `$effect`（line 156-177 付近）と `setCurrentPhasePage`（line 112-139 付近）
- 単体テスト: `JobRunPage` の state 遷移を Svelte testing で証明できるか別途検討
- システムテスト: `E2E-UC-048/049` を pass にする

## 後続フロー

- 入口: `preparation-module` → `investigation-module` → `implementation-module` → `finalization-module`
- 主因確定済みのため、investigation-module は plan.md 内の観測記録と修正方針候補を fix_decider に渡せばよい
