# 修正前調査

## 人間観測の補正

- 補正日時: 2026-05-09
- 補正内容: 「何も表示されない」という表現は強すぎる。
- 補正後の観測: 初回操作では、少なくとも単語翻訳フェーズ UI が表示されない。
- 影響: この調査に含まれる `空表示` という表現は、補足調査前の表現として扱う。
- 影響: 原因候補は、画面全体の空表示ではなく、単語翻訳フェーズ UI の表示条件を中心に再確認する必要がある。

## 判断結果

- 判定: 完了
- 調査 mode: `修正前調査`, `UI 根拠`, `trace`
- 引き継ぎ先: `fix_lane`
- 判断: 初回だけ単語翻訳フェーズ UI が表示されない症状について、frontend 側の状態同期競合を原因候補として切り分け可能な観測事実を揃えた。
- 判断: `agent-browser` では人間観測どおりの遷移完了までは再現できなかった。
- 判断: backend 側の明示エラーは今回の症状へ直結する形では観測できなかった。

## 根拠参照

- 人間観測: `./human-observation.md`
- 起動入力: `./investigation-input.md`
- UI 証跡:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/list-before-first-open.png`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/before-first-open.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/link-inspect.json`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/app-text.txt`
- frontend log:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/before-first-open.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/console.filtered.txt`
- backend log:
  - `tmp/logs/wails-dev.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.before-first-open.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.backend.filtered.log`
- 参照実装:
  - `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`
  - `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## 観測事実

- 観測事実: `#translation-management` の一覧には `ジョブ #1` が 1 件表示され、`現在の翻訳段階へ進む` は enabled 表示だった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/before-first-open.snapshot.txt`, `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/app-text.txt`
- 観測事実: DOM 上のジョブリンクは `href="#translation-management/job-run"` を持っていた。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/link-inspect.json`
- 観測事実: `agent-browser click @e18` と `agent-browser click @e17` の後も URL は `http://127.0.0.1:34115/#translation-management` のままだった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.url.txt`
- 観測事実: click 後 snapshot でも `未完了ジョブ一覧` が残り、`job-run` 側の見出しや `ジョブ未選択` 表示は観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.snapshot.txt`
- 観測事実: `TranslationJobManagementPage` は click 時に、一覧カード上の `job.jobRunTarget` を使って即座に `onOpenJobRun(job.jobRunTarget)` を呼び、その後で `await controller.selectJob(job.jobId)` を実行する。根拠: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:103-111`
- 観測事実: 同じ画面は subscription のたびに `onJobRunTargetChange(nextViewModel.jobRunTarget)` を呼ぶ。根拠: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:46-49`
- 観測事実: `selectJob` の loading 更新では `selectedJobId` と `detailPhase` だけを先に更新し、その時点では `selectedJobDetail` を埋めない。根拠: `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts:119-129`
- 観測事実: presenter は一覧カードごとの `jobRunTarget` を summary から生成するが、画面全体の `viewModel.jobRunTarget` は `selectedJobDetail` からだけ生成する。根拠: `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:264-315`, `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:318-342`, `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts:442`
- 観測事実: `AppShell` は `openTranslationJobRun` で `selectedJobRunTarget` が non-null なら `selectedTranslationManagementViewId = "job-run"` と `#translation-management/job-run` を設定する。根拠: `frontend/src/ui/views/AppShell.svelte:177-193`
- 観測事実: `JobRunPage` は `selectedJobTarget` が null になると phase controller 全部へ `setJobId(null)` を送り、`ジョブ未選択` 分岐を描画する。根拠: `frontend/src/ui/screens/job-run/JobRunPage.svelte:116-133`, `frontend/src/ui/screens/job-run/JobRunPage.svelte:397-410`

## UI 証跡

### 一覧画面

- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/list-before-first-open.png`
  - `ジョブ #1`
  - `現在の翻訳段階へ進む`
  - 一覧 1 件
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/before-first-open.snapshot.txt`
  - `link "ジョブ 1 を選択して現在の翻訳段階へ進む"`
  - `button "現在の翻訳段階へ進む"`

### click 後

- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.url.txt`
  - hash は `#translation-management`
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/after-click.snapshot.txt`
  - 一覧画面のまま
- `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/link-inspect.json`
  - DOM 上の link target は `#translation-management/job-run`

## ログ証跡

### frontend log

- 観測事実: filter 後 console では `GetJobDetail`, `ListIncompleteJobs`, `translation-management/job-run` を含む行を確認できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/console.filtered.txt`
- 観測事実: filter 後 console では `frontend.runtime.master_dictionary` の既存 progress 行が大量に残っていた。今回の job open 操作に対応する新規行は切り出せなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/console.filtered.txt`

### backend log

- 観測事実: `tmp/logs/wails-dev.log` には起動時の `runtime:ready` 未知 message が 3 回記録されていた。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.backend.filtered.log`
- 観測事実: backend log からは `jobID1` の初回 open 失敗に対応する job detail 取得失敗、phase 読込失敗、panic は確認できなかった。根拠: `tmp/logs/wails-dev.log`

## 仮説

- 仮説: 初回 open では、一覧カード由来の `job.jobRunTarget` で `job-run` へ入った直後に、loading 中の subscription が `nextViewModel.jobRunTarget = null` を返し、`AppShell` の `selectedJobRunTarget` を一度空に戻す可能性がある。
- 仮説: `selectedJobRunTarget` が一度 null へ戻ると、`JobRunPage` の `$effect` が `setJobId(null)` を各 phase controller へ送り、初回だけ target 未固定状態を作る可能性がある。
- 仮説: 2 回目 open では `selectedJobDetail` がすでに準備済みのため、`viewModel.jobRunTarget` が null を経由せず、画面が安定する可能性がある。
- 注意: 上記 3 件はコード上の整合仮説であり、今回の `agent-browser` 実行では実際の phase 画面遷移まで確認できていない。

## 観測点

- 入口: `http://127.0.0.1:34115/#translation-management`
- 対象 job: `ジョブ #1`
- 操作: 一覧カード link 押下、一覧カード action button 押下
- 参照境界:
  - 一覧カード click 処理
  - job detail 取得の loading 更新
  - presenter の `jobRunTarget` 生成元
  - `AppShell` の route 切替
  - `JobRunPage` の null target 分岐

## 影響ファイル候補

- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - 理由: open と select の順序、subscription での target 同期を持つ。
- `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`
  - 理由: 初回 `selectJob` loading 更新で `selectedJobDetail` を埋めない。
- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - 理由: 一覧カード target と selected detail target の生成元が分かれている。
- `frontend/src/ui/views/AppShell.svelte`
  - 理由: `selectedJobRunTarget` を job-run 遷移の成立条件として持つ。
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - 理由: null target 時に controller を null jobId へ戻す。

## 残り不足

- 未確認: 人間観測どおりに初回 click で `job-run` へ遷移した直後の DOM と単語翻訳フェーズ UI の有無。理由: `agent-browser` click では hash 変化まで確認できなかった。
- 未確認: 初回 open 時の `selectedJobRunTarget` の実時間変化。理由: frontend state を直接記録する観測ログが現状ない。
- 未確認: `GetJobDetail` 呼び出しの直前直後に console へ何が出るか。理由: console に過去セッションの大量行が残り、今回操作の切り出しに失敗した。
- 未確認: `jobID1` 固有か、初回 open 一般か。理由: 今回は人間観測対象の 1 件だけを扱った。

## 残留リスク

- リスク: UI 再現が `agent-browser` 上で最後まで成立していないため、単語翻訳フェーズ UI が表示されない位置は未確定。
- リスク: backend `runtime:ready` 行は存在するが、今回の症状との因果は未確認。
- リスク: 仮説どおりなら frontend の同期順序問題だが、phase controller 側の初期 mount と併発している可能性も残る。

## 原因未確認による停止理由

- 停止理由: `原因箇所シーケンス図` に必要な `原因箇所`、`問題点`、`修正方針` を確認済みとは言えない。
- 根拠: 現時点の `原因箇所` は `仮説` 節に分離しており、観測事実だけでは確定していない。
- 根拠: `agent-browser` では人間観測どおりの初回遷移完了を再現できず、単語翻訳フェーズ UI が表示されない時点の DOM と route 状態を取得できていない。
- 根拠: frontend log から `GetJobDetail` 呼び出しと `selectedJobRunTarget` の時系列変化を切り出せていない。
- 根拠: backend log からは今回症状に対応する失敗ログを確認できていない。
- 判断: `diagramming` skill の停止条件である「未確認の原因または未確認の修正案を含める場合」に該当するため、`fix_lane` は `原因箇所シーケンス図` へ進まず停止すべきである。

## 推奨 next step

- 推奨: `fix_lane` は `原因箇所シーケンス図` を起動せず停止し、追加観測の要否を人間判断へ戻す。
- 推奨: 追加観測を行う場合は、frontend 一時観測点で `selectedJobRunTarget` と `selectedTranslationManagementViewId` の遷移だけを取る。
- 推奨: 追加観測で初回遷移直後の DOM、route、job detail 読込順、単語翻訳フェーズ UI の表示条件が確認できた後にだけ、`原因箇所シーケンス図` へ進む。
