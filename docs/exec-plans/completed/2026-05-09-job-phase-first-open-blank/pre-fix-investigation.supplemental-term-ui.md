# 補足修正前調査: 単語翻訳フェーズ UI 初回未表示

## 判断結果

- 判定: 補足調査は完了
- 調査 mode: `修正前調査`, `UI 根拠`, `trace`
- 引き継ぎ先: `fix_lane`
- 判断: `agent-browser` では、人間観測どおりの「初回押下後に job-run へ入る遷移」は再現できなかった。
- 判断: ただし、単語翻訳フェーズ UI の表示条件、`job-run` route の成立条件、`selectedJobRunTarget` と `selectedTranslationManagementViewId` の更新経路は、UI 証跡とコード参照で補足できた。
- 判断: 人間観測どおりの初回遷移直後 DOM は未取得のため、原因箇所の確定には足りない。

## 根拠参照

- 人間観測: `./human-observation.md`
- 起動入力: `./investigation-input.md`
- 既存調査: `./pre-fix-investigation.md`
- 補足 UI 証跡:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.png`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.eval.json`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.png`
- frontend log:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.console.filtered.txt`
- backend log:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.supplemental.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.supplemental.backend.filtered.log`
- 参照実装:
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`
  - `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## 観測点

- 入口: `http://127.0.0.1:34115/#translation-management`
- 対象 job: `jobID1`, 画面表示では `ジョブ #1`
- 操作:
  - 一覧で `現在の翻訳段階へ進む` を role 指定で押下
  - `#translation-management/job-run` を直接開く
- 確認対象:
  - 初回押下後 route
  - 初回押下後 DOM
  - 単語翻訳フェーズ UI の表示条件
  - `selectedJobRunTarget` と `selectedTranslationManagementViewId` の更新経路
  - job detail 読込順

## 観測事実

### UI 証跡

- 観測事実: 一覧画面には `ジョブ #1` と `現在の翻訳段階へ進む` が表示されていた。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.snapshot.txt`
- 観測事実: `現在の翻訳段階へ進む` を role 指定で押下した後も URL は `http://127.0.0.1:34115/#translation-management` のままだった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.url.txt`
- 観測事実: 押下後 snapshot でも `未完了ジョブ一覧` が残り、`選択中のジョブ`、`ジョブ未選択`、`単語翻訳の次の作業` は観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.snapshot.txt`
- 観測事実: `#translation-management/job-run` を直接開いても、表示 URL は `http://127.0.0.1:34115/#translation-management` だった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.url.txt`
- 観測事実: 直接 open 後の DOM でも一覧画面の本文が残り、`job-run` 画面の見出しは観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.snapshot.txt`, `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.eval.json`

### route と表示条件のコード観測

- 観測事実: `AppShell` は `selectedTranslationManagementViewId === "job-run"` の時だけ `JobRunPage` を描画し、`selectedJobTarget={selectedJobRunTarget}` を渡す。根拠: `frontend/src/ui/views/AppShell.svelte:97-103`, `frontend/src/ui/views/AppShell.svelte:387-396`
- 観測事実: `openTranslationJobRun` は、渡された target が non-null の時だけ `selectedJobRunTarget` を設定する。`selectedJobRunTarget` が null のままなら `selectedTranslationManagementViewId = "job-management"` で return する。non-null の時だけ `selectedTranslationManagementViewId = "job-run"` と `#translation-management/job-run` を設定する。根拠: `frontend/src/ui/views/AppShell.svelte:177-193`
- 観測事実: `syncRouteFromHash` は `#translation-management/job-run` を検出すると、`selectedJobRunTarget = null`、`selectedTranslationManagementViewId = "job-management"` を設定し、hash を `#translation-management` へ置き換える。根拠: `frontend/src/ui/views/AppShell.svelte:119-145`
- 観測事実: `JobRunPage` は `selectedJobTarget` が falsy の時、phase controller 全部へ `setJobId(null)` を送り、`ジョブ未選択` 分岐だけを描画する。`selectedJobTarget` が truthy の時だけ `TermTranslationPhasePanel` を含む phase UI を描画できる。根拠: `frontend/src/ui/screens/job-run/JobRunPage.svelte:116-133`, `frontend/src/ui/screens/job-run/JobRunPage.svelte:300-410`
- 観測事実: `handleOpenJob` は一覧カード上の `job.jobRunTarget` を使って、先に `onJobRunTargetChange(job.jobRunTarget)` と `onOpenJobRun(job.jobRunTarget)` を呼び、その後で `await controller.selectJob(job.jobId)` を実行する。根拠: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:103-112`
- 観測事実: 一覧画面の subscription は、viewModel 更新のたびに `onJobRunTargetChange(nextViewModel.jobRunTarget)` を呼ぶ。根拠: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:46-49`
- 観測事実: `selectJob` は detail 読込前に `selectedJobId` と `detailPhase = "loading"` を更新するが、その時点では `selectedJobDetail` を設定しない。detail 取得成功後にだけ `selectedJobDetail = detail` と `detailPhase = "ready"` を設定する。根拠: `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts:106-137`
- 観測事実: presenter は一覧カードごとの `jobRunTarget` を `jobs` summary から作るが、画面全体の `viewModel.jobRunTarget` は `selectedJobDetail` からだけ作る。`selectedJobDetail` が null の間、`viewModel.jobRunTarget` は null になる。根拠: `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:286-301`, `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:318-342`, `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:442`

### frontend log

- 観測事実: 補足採取した filtered console には `translation-management/job-run`、`GetJobDetail`、`ListIncompleteJobs` は残っていなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.console.filtered.txt`
- 観測事実: filtered console には `Connected to backend`、`Disconnected from backend`、`runtime:ready` の行が残っていた。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.console.filtered.txt`

### backend log

- 観測事実: 補足採取した backend filtered log には `runtime:ready -> Unknown message from front end: runtime:ready` が複数行あった。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.supplemental.backend.filtered.log`
- 観測事実: 補足採取した backend filtered log には `GetJobDetail`、`ListIncompleteJobs`、panic は確認できなかった。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.supplemental.backend.filtered.log`

## UI 証跡

- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.snapshot.txt`
  - 一覧画面のまま
  - `未完了ジョブ一覧`
  - `現在の翻訳段階へ進む`
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.url.txt`
  - hash は `#translation-management`
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.url.txt`
  - 直接 open 後も hash は `#translation-management`
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.direct-job-run-open.eval.json`
  - bodyText は一覧画面本文

## ログ証跡

### frontend log

- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/supplemental.after-role-click.console.filtered.txt`
  - `Connected to backend`
  - `Disconnected from backend`
  - `runtime:ready`
  - `GetJobDetail` 行なし
  - `translation-management/job-run` 行なし

### backend log

- `tmp/logs/2026-05-09-job-phase-first-open-blank.supplemental.backend.filtered.log`
  - `runtime:ready -> Unknown message from front end: runtime:ready`
  - `GetJobDetail` 行なし
  - panic 行なし

## 仮説

- 仮説: `handleOpenJob` の先行 `onOpenJobRun(job.jobRunTarget)` の後で、subscription による `onJobRunTargetChange(nextViewModel.jobRunTarget)` が `selectedJobDetail` 未読込の null を反映し、`selectedJobRunTarget` を空へ戻す可能性がある。
- 仮説: `selectedJobRunTarget` が空へ戻ると、`JobRunPage` は単語翻訳フェーズ UI を維持できず、`selectedTranslationManagementViewId` との同期順によって初回だけ UI が不安定になる可能性がある。
- 注意: 上記はコード経路に基づく原因候補であり、人間観測どおりの初回遷移直後 DOM では未確認である。

## 再実行時の表示条件に関する補足

- 観測事実: `JobRunPage` が単語翻訳フェーズ UI を描画する条件は `selectedJobTarget` が truthy で、かつ `currentPhasePage === "term"` の場合である。根拠: `frontend/src/ui/screens/job-run/JobRunPage.svelte:127-133`, `frontend/src/ui/screens/job-run/JobRunPage.svelte:331-347`
- 観測事実: `currentPhasePage` は `selectedJobTarget.currentPhase` が `persona_generation` と `body_translation` 以外なら `term` に解決される。根拠: `frontend/src/ui/screens/job-run/JobRunPage.svelte:75-83`
- 未確認: 一覧へ戻って同じ操作を再実行した時に、`selectedJobDetail` が保持済みで `viewModel.jobRunTarget` が null を経由しないかどうか。理由: 人間観測どおりの再実行経路を `agent-browser` で再現できていない。

## 影響ファイル候補

- `frontend/src/ui/views/AppShell.svelte`
  - 理由: `selectedTranslationManagementViewId`、`selectedJobRunTarget`、hash 正規化を持つ。
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - 理由: open と select の呼び順、subscription 同期を持つ。
- `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`
  - 理由: `selectedJobDetail` の読込順を持つ。
- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - 理由: 一覧 card target と画面全体 target の生成元が異なる。
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - 理由: target null 時に単語翻訳フェーズ UI を描画しない。

## 残り不足

- 未確認: 人間観測どおりの初回押下直後に `#translation-management/job-run` が一度でも成立したか。理由: `agent-browser` 押下では hash 変化を採取できなかった。
- 未確認: 初回押下直後 DOM に `JobRunPage` の `ジョブ未選択` 分岐が一瞬でも出たか。理由: `agent-browser` では一覧画面のままだった。
- 未確認: 再実行時に単語翻訳フェーズ UI が表示された直前の `selectedJobRunTarget` と `selectedTranslationManagementViewId` の値。理由: state を直接記録する観測点がない。
- 未確認: `GetJobDetail` 呼び出しが初回押下時に実際に発火したか。理由: frontend console と backend log の今回採取分では該当行を切り出せなかった。

## 残留リスク

- リスク: 補足調査の UI 証跡は、人間観測どおりの初回遷移完了状態を含まない。
- リスク: `runtime:ready` の接続揺れが補足採取中に存在し、押下不成立と症状本体が混在している可能性がある。
- リスク: 仮説どおりなら frontend state 同期順の問題だが、現在は観測事実だけで原因箇所を確定できない。

## 推奨 next step

- 推奨: `fix_lane` は、現時点では `原因箇所シーケンス図` へ進まず停止する。
- 推奨: 進める条件は、人間観測どおりの初回遷移直後について、少なくとも route、DOM、`selectedJobRunTarget`、`selectedTranslationManagementViewId` の時系列が追加観測できた後とする。
