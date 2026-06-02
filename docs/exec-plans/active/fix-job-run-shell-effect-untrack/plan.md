# fix-job-run-shell-effect-untrack

## 依頼要約

`frontend/src/ui/screens/job-run/JobRunPage.svelte` の `$effect` が `setCurrentPhasePage` 関数経由で `currentPhasePage`（`$state`）を読み取り、Svelte 5 の reactivity が依存に登録されている。その結果、ユーザー操作で `setCurrentPhasePage("persona")` または `setCurrentPhasePage("complete")` を呼ぶと、`$effect` 再実行が `selectedJobTarget.currentPhase` 起点の初期 phase（"term" / "body"）へ戻してしまい、画面遷移が機能しない。`E2E-UC-048`（term → persona）と `E2E-UC-049`（body → complete）が pass しない。

## 分岐元

- 分岐元 task: `fix-phase-ai-settings-pill-update-after-model-select`（agent が発見、本 task 主旨と独立のため別出し）
- 分岐元 branch: `master`
- 分岐元 commit: `a5d0dbe4`
- 作業 branch: `claude/fix-job-run-shell-effect-untrack`

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

## 想定 Y/N 評価（investigation-module 入口）

- 仕様変更または仕様追加がある: N。term→persona、body→complete の既存遷移仕様が動かないのを直す（`tests/system/job-run-shell.spec.ts` の E2E-UC-048/049 が既存期待）。
- 画面変更がある: N。phase 画面構造は変えず、遷移ロジックだけを直す（`frontend/src/ui/screens/job-run/JobRunPage.svelte`）。
- 内部構造変更がある: Y。`$effect` と `setCurrentPhasePage` の Svelte 5 reactivity 境界を変える（`JobRunPage.svelte:112-177`）。
- 画面の表示変更がある: N。layout、文言、style、表示構造はいずれも変えない。
- frontend ロジック変更がある: Y。state（`currentPhasePage`）と副作用（`$effect`）のロジック層を変える（`implementation-module` 範囲、`storybook-module` 範囲外）。
- backend 変更がある: N。修正対象は frontend のみ。
- frontend と backend を接続する: N。Wails bridge、API、DTO の境界変更はない。
- 実装済み責務を独立に証明したい: Y。`JobRunPage` の phase 遷移 state ロジックを Svelte testing で証明したい（plan の影響範囲節に既記）。
- 実行時にしか確定しない値または原因分離が要る分岐がある: Y。Svelte 5 の依存追跡（`untrack` 適用範囲、`$effect` 再実行条件）は実行時挙動で、`chrome-devtools` 再現が要る。

→ 「仕様変更または仕様追加がある」が N のため investigation-module 続行。

## 修正実行入力（investigation-module 出口、人間承認済み）

承認日: 2026-06-02。詳細は `investigation-summary.md` および同フォルダの個別成果物参照。

- **承認済み修正方針**: `frontend/src/ui/screens/job-run/JobRunPage.svelte` line 113 の `const previous = currentPhasePage` を `const previous = untrack(() => currentPhasePage)` に変更し、`import { untrack } from "svelte"` を追加する。
- **禁止する修正**:
  1. `$effect` 内 `setCurrentPhasePage` 削除して直接代入する変更（初期 fetch トリガー喪失）
  2. `isUserNavigated` 等の新 `$state` 追加で `$effect` 再実行をスキップする対症療法
  3. `$effect.pre` / `afterUpdate` 相当でタイミングをずらす変更
  4. `resolveInitialPhasePage` の戻り値キャッシュ比較
- **影響ファイル候補**:
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`（`untrack` import 追加と line 113 の 1 行変更）
  - `frontend/src/ui/screens/job-run/` 配下の新規単体テストファイル
- **承認済み UC 差分**: 差分なし（仕様正本への追記なし）
- **承認済み E2E テスト観点差分**:
  - 既存 E2E-UC-048 / 049 / 050: 差分なし
  - 追加: `E2E-UC-048-B1`（境界。persona 遷移後 500ms 待機で term に差し戻されないことを確認。selector は既存）
- **画面再現確認の再現手順と修正後の期待状態**: `npm run dev:wails:run` で `http://localhost:34115` を起動し、jobId=3 で `setCurrentPhasePage("persona")` を経由するユーザー操作後、`persona-generation-phase-screen` が visible のまま維持され、`$effect` 再実行で term に戻らないこと。

## implementation-module decision table 結果

想定 Y/N 評価（investigation-module 出口）に基づく：

| 想定 | Y/N | 起動 artifact |
| --- | --- | --- |
| frontend ロジック変更がある | Y | frontend ロジック実装 |
| backend 変更がある | N | - |
| frontend と backend を接続する | N | - |
| 実装済み責務を独立に証明したい | Y | 単体テスト |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | 観測ログ追加 判断 |

加えて修正系 task 経路の fail-test 要件として、E2E-UC-048-B1 追加と既存 E2E-UC-048/049 の pass 確認のため `シナリオテスト` も起動する。

## 実装結果

- frontend ロジック実装: `frontend/src/ui/screens/job-run/JobRunPage.svelte` に `untrack` import 追加と line 113 ラップ適用済み
- 単体テスト追加: `frontend/src/ui/screens/job-run/__tests__/untrack-effect-dependency.test.ts`（3 テスト）と補助 `UntrackEffectProbe.svelte`、`frontend/package.json` の knip ignore エントリ追加
- シナリオテスト追加: `tests/system/job-run-shell.spec.ts` に E2E-UC-048-B1 追加。E2E-UC-048 / 049 / 048-B1 すべて pass 確認済み

## 観測ログ追加判断

- 判断: 省略
- 理由: investigation-module で確定原因が `Svelte 5 reactivity の依存追跡` という静的な構文起因に固定されており、実行時データ依存の分岐ではない。一時観測ログは fix_decider が削除済み。回帰は E2E-UC-048-B1 の境界観点で検出可能なため、恒久観測ログは不要と判断した。

## 最終検証

- 実行 suite: `python3 scripts/harness/run.py --suite frontend-local`
- 結果: PASS
  - lint: PASS（ESLint、tsc、knip、boundary）
  - test: 582 tests / 54 files 全 PASS
- E2E playwright 個別実行: E2E-UC-048 / 049 / 048-B1 すべて pass

## 正本化判断（finalization-module）

- 仕様変更または仕様追加: なし（investigation-module の想定 Y/N 評価で固定、UC 差分候補も「差分なし」）
- 人間承認済みの恒久仕様: なし
- 詳細仕様正本反映: 不要

## 実画面確認

- 修正前: `fix_decider` が `http://localhost:34115` で jobId=3 を選択し、`setCurrentPhasePage("persona")` 経路を踏んで `$effect` 再実行による term への差し戻しをコンソールログで観測（`fix-decision.md` 画面再現確認節）。
- 修正後: dev seed には完了済み term ジョブが存在せず、実 app での「次へ進む」クリック経路は到達できないため、E2E playwright（headless Chromium、fixture `system-test-completed-term`）の E2E-UC-048 / 049 / 048-B1 pass を実画面動作確認の代替とする。
