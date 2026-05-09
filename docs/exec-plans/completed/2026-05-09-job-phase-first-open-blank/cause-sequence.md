# 原因箇所シーケンス図

## 判断結果

- 判定: 完了
- 図成果物種別: `原因箇所シーケンス図`
- source: `./cause-sequence.puml`
- 描画結果: `./cause-sequence.svg`

## 図説明

- 対象: 初回操作では route が `job-run` へ進んだ後に、`JobRunPage` へ渡す job target が消え、下部 panel が `ジョブ未選択` になる呼び出し順序を固定する。
- 対象: 再実行では同じ route で `ジョブ #1` と `単語翻訳` UI が表示される差分を、target 保持の有無に限定して示す。
- 範囲: `TranslationJobManagementPage.handleOpenJob`、`selectJob` の loading 更新、presenter の `viewModel.jobRunTarget`、`AppShell.selectedJobRunTarget`、`JobRunPage` の null 分岐だけを扱う。

## 問題点

- 問題点: 初回操作では、一覧 card 由来の `job.jobRunTarget` で `job-run` へ進んだ後に、loading 中の `viewModel.jobRunTarget = null` が親の `selectedJobRunTarget` を上書きする。
- 問題点: `JobRunPage` は `selectedJobTarget = null` を受けると、phase controller へ `setJobId(null)` を送り、下部 panel を `ジョブ未選択` 分岐へ切り替える。
- 問題点: 人間観測と手動再現証跡の差分は route 遷移ではなく、`JobRunPage` に target が残るかどうかにある。

## 修正方針

- 修正方針: detail loading 中は、`selectedJobId` に対応する一覧 summary から `viewModel.jobRunTarget` を生成する。
- 修正方針: `selectedJobDetail` が存在する時は従来どおり detail を優先し、stale 状態では summary から target を復元しない。

## 根拠参照

- 人間観測: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/human-observation.md)
- 手動再現証跡: [pre-fix-investigation.manual-reproduction.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/pre-fix-investigation.manual-reproduction.md)
- 調査まとめ: [pre-fix-investigation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/pre-fix-investigation.md)
- 補足 UI 調査: [pre-fix-investigation.supplemental-term-ui.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/pre-fix-investigation.supplemental-term-ui.md)
- 実装根拠:
  - [TranslationJobManagementPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:46)
  - [translation-job-management.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts:106)
  - [translation-job-management.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:318)
  - [AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte:177)
  - [JobRunPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/job-run/JobRunPage.svelte:116)

## 実装根拠の要点

- `TranslationJobManagementPage` は subscription のたびに `onJobRunTargetChange(nextViewModel.jobRunTarget)` を呼ぶ。根拠: [TranslationJobManagementPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:46)
- `handleOpenJob` は `onOpenJobRun(job.jobRunTarget)` を先に呼び、その後で `await controller.selectJob(job.jobId)` を実行する。根拠: [TranslationJobManagementPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:103)
- `selectJob` の loading 更新は `selectedJobId` と `detailPhase` だけを更新し、その時点では `selectedJobDetail` を設定しない。根拠: [translation-job-management.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts:119)
- 修正前の presenter は、一覧 card ごとの `jobRunTarget` を summary から作る一方で、画面全体の `viewModel.jobRunTarget` は `selectedJobDetail` からだけ作っていた。`selectedJobDetail` が null の間は `viewModel.jobRunTarget = null` になっていた。根拠: [translation-job-management.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:290), [translation-job-management.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:318), [translation-job-management.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:442)
- `AppShell` は `selectedJobRunTarget` が non-null の時だけ `job-run` を維持し、`JobRunPage` へ `selectedJobTarget={selectedJobRunTarget}` を渡す。根拠: [AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte:177), [AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte:387)
- `JobRunPage` は `selectedJobTarget` が null の時に `setJobId(null)` を送り、`ジョブ未選択` 分岐を描画する。根拠: [JobRunPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/job-run/JobRunPage.svelte:116), [JobRunPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/job-run/JobRunPage.svelte:397)

## 観測根拠の要点

- 人間観測では、初回操作後に `ジョブの進み方` と `単語翻訳` の強調は見えるが、下部 panel は `ジョブ未選択` を表示した。根拠: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/human-observation.md)
- 手動再現証跡では、初回操作後の URL は `#translation-management/job-run` で、下部 panel は `未完了ジョブ一覧でジョブを選んでください` を表示した。再実行後は同じ route で `ジョブ #1` と `単語翻訳` UI が表示された。根拠: [pre-fix-investigation.manual-reproduction.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/pre-fix-investigation.manual-reproduction.md)

## 検証結果

- 検証結果: `plantuml -tsvg docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/cause-sequence.puml` は成功した。
- 検証結果: `docs/exec-plans/active/2026-05-09-job-phase-first-open-blank/cause-sequence.svg` の生成を確認した。
- 追加検証結果: `npm --prefix frontend run test -- translation-job-management.presenter.test.ts AppShell.test.ts` は成功した。
