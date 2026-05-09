# Null Source Investigation

## 判断結果

- 判定: 追加調査完了
- 対象: `viewModel.jobRunTarget = null` の発生源
- 結論: `null` は backend response ではなく、frontend presenter が detail loading 中の state から生成していた。

## 発生順序

- `TranslationJobManagementPage.handleOpenJob` は、一覧 card の `job.jobRunTarget` を使って `job-run` を開く。
- その直後に `controller.selectJob(job.jobId)` が呼ばれる。
- `selectJob` は loading 更新で `selectedJobId` と `detailPhase` を更新する。
- loading 更新時点では `selectedJobDetail` はまだ `null` である。
- presenter は画面全体の `jobRunTarget` を `selectedJobDetail` だけから作っていた。
- そのため、loading 更新の view model では `jobRunTarget` が `null` になっていた。

## 根拠

- `TranslationJobManagementPage.svelte`: `handleOpenJob` は `onOpenJobRun(job.jobRunTarget)` の後に `controller.selectJob(job.jobId)` を呼ぶ。
- `translation-job-management.usecase.ts`: `selectJob` の loading 更新は `selectedJobDetail` を設定しない。
- `translation-job-management.presenter.ts`: 修正前の画面全体 `jobRunTarget` は `toJobRunTarget(state.selectedJobDetail)` から生成されていた。
- `TranslationJobManagementPage.svelte`: subscription は `onJobRunTargetChange(nextViewModel.jobRunTarget)` を親へ通知する。

## 修正方針

- `detailPhase` が `loading` で `selectedJobId` が存在する場合だけ、一覧 summary から `jobRunTarget` を生成する。
- `selectedJobDetail` が存在する場合は、従来どおり detail を優先する。
- `detailPhase` が `stale`、`idle`、または選択 job が一覧に存在しない場合は、`jobRunTarget` を `null` のままにする。

## 影響

- 初回操作の loading 中でも、`job-run` へ渡す selected job target が維持される。
- detail 取得失敗時の stale 状態では、summary から target を復元しない。
- `AppShell` 側で `null` を握りつぶす必要はない。

## 追加検証

- `npm --prefix frontend run test -- translation-job-management.presenter.test.ts AppShell.test.ts`: 成功。
- `python3 scripts/harness/run.py --suite frontend-local`: 成功。
- `agent-browser open http://127.0.0.1:34115/#translation-management`: 成功。
- `agent-browser click @e18`: 成功。初回操作後に `ジョブ #1` と `単語翻訳` UI を確認した。
- screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/null-source-fix-after-first-open.png`
- `agent-browser errors`: なし。
