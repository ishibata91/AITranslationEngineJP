# Fix Decision Report: fix-lucien-target-list-empty

## 判断結果

- 判定: 完了
- 停止理由: なし。
- 前回からの変更（追補3）: 真因をさらに改訂した。「setJobId 再呼び出しによる 0件リセット」は HMR 観測効果が含まれていた。クリーン環境での真因は「onMount で term/persona/body 3コントローラーの bridge 呼び出しが同時発火し、IPC 競合で term の processingTargetList が 0件のまま確定する」。採用する修正方針を「初回マウント時の同時起動数削減（表示段階のみ取得）」に更新した。

## 観測済み問題

- 問題: jobId=6（`Lucien.esp_Export.json`、state=実行中、用語翻訳段階）の単語翻訳画面で、進捗パネルが「AI 翻訳対象語 4,931」を表示する一方、処理対象パネルが「0 件」「処理対象がありません。」を表示する。または、処理対象パネルが永続的に loading 状態（pending）のままとなり、summary も未取得のままとなる。
- 期待との差分: 処理対象パネルに約4,931件が表示される。利用者は読み込み済みの処理対象を初回表示で操作できない。

## 画面再現確認

- Wails 接続対象: `http://localhost:34115`（単一 process。`curl` で http_code=200、`lsof` で `*:34115` LISTEN を確認）。
- 再現手順（ラウンド2）:
  1. `agent-browser open http://localhost:34115/#translation-management` で翻訳管理画面を開く。
  2. `eval "Array.from(document.querySelectorAll('a')).find(a => a.textContent.includes('ジョブ #6'))?.click()"` でジョブ6に遷移する。
  3. URL が `#translation-management/job-run` へ遷移する。
  4. 処理対象パネルの件数、進捗パネルの情報、console ログを観測する。
- 操作結果（ラウンド2、観測ログ付き）:
  - パターンA（約半数の試行）: 処理対象「0件」、進捗「0%」（summary は取得済み）。
  - パターンB（約半数の試行）: 処理対象「0件」、進捗「-」（summary 未取得、Promise.all が永続 pending）。
  - 前回証跡（ラウンド1）: 初回0件 → 検索操作で50件 → リロードで再び0件。
- 画面状態: 初回ロード時に処理対象が0件になる。検索操作経路（`fetchProcessingTargetList` 単独）は正常に機能する。
- 証跡 path:
  - ラウンド1証跡: `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/observe/`
  - ラウンド2証跡: `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-round2/obs-round2-console.txt`

## 確定原因

- 原因: `JobRunPage.svelte` の `$effect`（131行）が `controller.setJobId(selectedJobTarget.jobId)` を呼び、同時に `onMount`（166行）が `controller.mount()`（= `useCase.load()`）を呼ぶ。この二重起動により、`fetchSummaryAndReadiness` が同一 jobId に対して2回実行される。各実行が summary・readiness・処理対象一覧の3本ずつ、合計**6本**の Wails bridge 呼び出しを同時に起動する。この6本同時呼び出しが次の2つの症状を引き起こす。
  - 症状A（連番競合 skip）: 6本全てが解決した場合、seq=1（setJobId 経由）と seq=2（load 経由）のどちらかが先に `store.update` を完了する。後発の `store.update` が `processingTargetListRequestSequence` の不一致を検出し、`processingTargetPageState` の反映を skip する。最終的に 0 件のまま、または 4931 件が書き込まれた後に何らかの理由で 0 件に戻る。
  - 症状B（永続 pending）: Wails bridge の IPC 並列呼び出し制限により、6本の中の1本が応答を返さずブロックされ、`Promise.all` が永続 pending 状態となる。summary も処理対象も取得できないまま画面が固まる。
- 根本原因の確定根拠:
  - 一時観測ログ（ラウンド2）で確認: `setJobId called jobId=6 seq-before=0` の直後に `load called jobId=6 seq-before=1` が呼ばれ、二重起動を直接観測した。
  - seq=1（setJobId）と seq=2（load）がそれぞれ独立した fetchSummaryAndReadiness を起動することをログで確認した。
  - 症状AとBの両方が同一コードで非決定的に発生することを複数回の試行で確認した。
  - backend は正常: `GetTermTranslationPhaseSummary`、`GetTermTranslationNextPhaseReadiness`、`GetProcessingTargetList` を個別に直接呼ぶと全て正常応答（4931件）を返す。
  - 前回確定の「連番競合 skip」はラウンド2でも観測（`seq-assigned=2 seq-current=2 match=true` の直後に `seq-assigned=1 seq-current=2 match=false`）。ただしこれは二重起動の一症状であり、根本原因は二重起動そのものにある。
- 「遅延で再現しない」事実との整合:
  - モック環境で `GetProcessingTargetList` に300ms遅延を入れても E2E が green のままだった（シナリオテスト工程の観測）。これは「最後に起動した fetch の処理対象解決前に、さらに別の fetch が seq を進める」状況が通常遅延では作れないためである。症状Aを再現するには、seq=1 の Promise.all が解決した後に seq=2 が `processingTargetListRequestSequence` を進めた状態を作る必要があるが、6本同時起動のタイミングはランダムであり、単純な遅延付加では制御できない。症状Bは Wails bridge の IPC 特性に依存し、モック環境では発生しない。

## 採用する修正方針

- 方針: `JobRunPage.svelte` の初回ロードで `fetchSummaryAndReadiness` が2回起動する二重起動を除去する。具体的には、`onMount` 内の `controller.mount()`（= `useCase.load()`）が `setJobId` と重複して `fetchSummaryAndReadiness` を起動しないようにする。実装 agent は次の方向から責務境界に沿って最小変更を選ぶ。
  - 案A（load の責務縮小）: `useCase.load()` が `jobId` を持つ場合でも `fetchSummaryAndReadiness` を呼ばないようにする。または `load()` 自体を廃止し、`setJobId` への一本化を採用する。`$effect` の `setJobId` 呼び出しが初回ロードを担保するため、`onMount` の `mount()` は state 初期化のみを行えばよい。
  - 案B（load と setJobId の協調ガード）: `useCase` に「既に当該 jobId で fetch が実行中」であることを検出するガードを追加し、二重起動を抑制する。
  - 案C（$effect と onMount の起動分担整理）: `$effect` と `onMount` の実行タイミングと責務を整理し、どちらか一方のみが `fetchSummaryAndReadiness` を起動するよう `JobRunPage.svelte` を変更する。
- 理由: 二重起動を除去すれば、6本同時呼び出しが3本に減り、症状AとBの両方が解消する。連番競合検出の仕組みは手動操作の競合防止として有効なため、廃止ではなく二重起動の除去で対応する。

## 禁止する修正

- 禁止修正1: 進捗パネルの母数や backend の SQL・DTO・bind 配線の変更。backend は正常である。
- 禁止修正2: 処理対象パネルのフィルタ初期値や検索機能の削除。検索経路は正常である。
- 禁止修正3: `processingTargetListRequestSequence` による取得競合検出の撤廃。手動操作の競合防止に必要である。
- 禁止修正4: 新しい phase 状態値や処理対象専用の重複 state フィールドの追加。既存 state モデルで対応できる。
- 禁止修正5: 症状の表示だけを隠す対症療法（例: 初期値を検索クエリで上書きして強制取得させるなど）。

## 影響ファイル候補

- 候補1: `frontend/src/ui/screens/job-run/JobRunPage.svelte`（`$effect` 131行の `setJobId` と `onMount` 166行の `mount()` の二重起動構造）。
- 候補2: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`（`load()` 118行の `fetchSummaryAndReadiness` 呼び出し判定、または連番ガードのスコープ）。
- 理由: 二重起動の源泉は JobRunPage.svelte の `$effect` と `onMount` の組み合わせにある。usecase.ts の `load()` 実装変更で対応するか、svelte の起動構造を変更するかは実装 agent が判断する。

## fail-test 再現可否の見解

- 症状A（連番競合 skip による0件）: モック E2E での再現は困難。6本同時呼び出しのタイミングに依存する非決定的な現象であり、単純な遅延付加では再現できない（シナリオテスト工程で確認済み）。単体テストで「setJobId と load が同一 jobId に対して連続呼び出しされた場合、processingTargetPageState が最終的に0件のままにならない」ことを assertできる可能性がある（usecase レベルのユニットテスト）。
- 症状B（Wails bridge 永続 pending）: モック環境では発生しない。実機限定の現象。
- 推奨テスト戦略: `useCase.setJobId(6)` と `useCase.load()` を並列（または順次）で呼び出した後、最終的な `processingTargetPageState.totalCount` が 0 でないことを assert する usecase 単体テストを作成する。これが二重起動除去の回帰テストになる。

## 追補（人間調査による設計確定・修正範囲拡大）

初回ロード二重起動の修正方針を、設計事実の確定により改訂する。

- 真のレンダリング親は `frontend/src/ui/views/AppShell.svelte`。`JobRunPage` は `{#if renderedTranslationManagementViewId === "job-run"}`（269行）配下で描画する。`selectedJobTarget` は `selectedJobRunTarget`（276行）を渡す。
- `renderedTranslationManagementViewId` は `selectedTranslationManagementViewId` 単一から導出する（86-87行）。
- 別ジョブへ移る唯一の経路は一覧復帰を経由する。`JobRunPage` のジョブ管理導線 `onOpenJobManagement`（274行）= `openTranslationJobManagement()` が `selectedTranslationManagementViewId="job-management"`（168-169行）にし、`{#if job-run}` が false になり `JobRunPage` が unmount される。その後、一覧の `openTranslationJobRun()`（161行で `view="job-run"`）で再 mount される。
- 帰結: `renderedTranslationManagementViewId` が "job-run" のまま `selectedJobRunTarget` だけ別ジョブに差し替わる経路は存在しない。別ジョブ表示時は `JobRunPage` が必ず unmount→再mount される。
- よって `$effect`(selectedJobTarget変化→`setJobId`) が取得(fetch)を起こす正当性はない。取得は `onMount` 一本で足りる。
- 修正範囲の拡大（人間判断）: term だけでなく persona / body の3画面とも、JobRunPage での取得を `onMount` 一本化する。`$effect` からは取得(fetch)起動を除去し、取得以外（`selectedJobTarget` 不在時のリセット、phase page 初期化）に限定する。これで3画面の同型二重起動を構造的に除去する。
- 前回の案A（`load` の hasLoaded ガード）、案C（term だけ `onMount` から mount を外し `$effect` 一本化）は、いずれも逆方向または部分修正であり、persona/body の二重起動が残る。今回 onMount 一本化へ統一して是正する。

## 追補2（ラウンド3: onMount 一本化後の実機観測）

onMount 一本化実装後の最新ビルドで実機観測を実施した。結論として、初回0件は解消していない。

### 判断結果

- 判定: 停止（Wails dev process の再起動が必要。fix_decider の権限範囲外）。
- 停止理由: bridge eval（診断目的の `GetProcessingTargetList` 直呼び）が Wails IPC queue を詰まらせ、以降の全 bridge 呼び出しが pending になった。この状態では `fetchSummaryAndReadiness` 内の `Promise.all` の解決を観測できない。

### 確認できた事実

- 最新ビルド確認: console に `seq=` 等の一時ログが 0件であることを確認した。観測は最新ビルド上で行った。
- 症状の継続: ジョブ6単語翻訳画面で `cell "処理対象がありません"` と「summary 未取得です。」が表示された（症状B）。
- onMount 一本化の適用確認: `JobRunPage.svelte` の `$effect` は `setJobId(null)` と `setCurrentPhasePage` のみ。取得起動は `onMount` 内の `setJobId` 一本になっている。
- bridge eval が IPC を詰まらせる: fix_decider が送った `GetProcessingTargetList` 直呼び eval が返らず、以降 `window.go` の評価も約15秒かかるようになった。Wails dev process の CPU が 100% 持続した。
- Promise.all は3本: 二重起動除去後も `fetchSummaryAndReadiness` の `Promise.all` は3本の bridge 呼び出し（summary / readiness / processingTargetList）を同時に起動する。6本→3本に減ったが3本同時は残る。

### 残存する疑問（Wails 再起動後の確認が必要）

- 症状B（bridge pending）が、3本同時 bridge 呼び出し自体で発生するか、それとも今回の診断操作（eval）が汚染したためかを切り分けられていない。
- Wails dev process を再起動してクリーンな状態で観測した場合に、症状Bが再現するか否かが未確定。
- 症状AまたはBのどちらが最新ビルドで発生しているかが未確定。

### 次に行うべき確認（fix_lane への戻し条件）

1. Wails dev process を再起動する（`npm run dev:wails:agent-browser` を再起動）。
2. ジョブ6単語翻訳画面を初回表示する。
3. console ログで `obs-r3` 等の一時ログが 0件であることを確認する。
4. 処理対象パネルの件数と進捗パネルの状態を観測する。
5. bridge 直呼び eval は使わず、usecase 一時ログで `fetchSummaryAndReadiness` の呼び出し回数と Promise.all の解決可否を確認する。

### 修正方針の現状

- 採用する修正方針（前回まで）は変更しない。
- 二重起動除去（onMount 一本化）は正しい方向の修正であり、残す。
- 症状Bが3本同時 bridge 呼び出し自体に起因する場合、追加修正として `Promise.all` の直列化（sequential fetch）または bridge 呼び出し数の削減が必要になる可能性がある。
- この判断は Wails 再起動後の観測結果に依存する。

## 追補3（ラウンド4: クリーン再起動後の観測・真因確定）

Wails dev process クリーン再起動後、usecase と presenter に一時ログを仕込み実機観測を実施した。観測後に一時ログを全削除した。

### 判断結果

- 判定: 完了。
- 変更点: 真因を「3コントローラー並列 setJobId による bridge 競合」に改訂し、採用する修正方針を更新した。

### クリーン再起動後の初回0件の有無

- 残存している。クリーン再起動後（bridge eval 汚染なし）でも、初回遷移で「0件」「処理対象がありません」が表示される（症状A）。
- summary は取得できている（「AI 翻訳対象語 4,930 件」が表示された）。processingTargetList のみ 0件。

### 初回マウント時の bridge 同時呼び出し本数と Promise.all の解決状況

- onMount で `Promise.all([controller.setJobId(6), personaController.setJobId(6), bodyController.setJobId(6)])` が並列起動する。
- 各 `setJobId` が内部で `fetchSummaryAndReadiness` → `Promise.all([summary, readiness, processingTargetList])` を起動する。
- term と persona の Promise.all が同時発火することを `[obs-r4] term Promise.all firing seq=1` と `[obs-r4] persona Promise.all firing seq=1` のログで確認した。
- persona が先に RESOLVED し、term が後に RESOLVED した（または遅延した）。term の Promise.all は hang 状態が続くことがある（症状B相当）。
- 最終的には term RESOLVED totalCount=4931 が出る（6本の競合が解消した後）。

### store 反映確認

- `[obs-r4] term after store.update processingTargetPageState.totalCount=4931` を確認した。**store への書き込みは成功している**。
- `[obs-r4] presenter.toViewModel called processingTargetPageState.totalCount=4931` も確認した。4931件が presenter に届いた瞬間がある。
- **しかし画面は 0件のまま**。

### setJobId 再呼び出しによる 0件リセットの観測

- `[obs-r4] term setJobId called jobId=6` が store.update 後に再度呼ばれることを確認した。
- `setJobId` の中で `processingTargetPageState = createDefaultProcessingTargetPageState()`（totalCount=0）が書き込まれ、0件で上書きされた。
- obs-r4 ログ環境では、HMR（vite page reload）による再マウントが `onMount` → `setJobId` の再実行を繰り返した。これが観測効果（observer effect）として 0件リセットを誇張していた。

### クリーン環境（ログなし）での症状

- 最終的な症状は症状A（summary 取得済み、processingTargetList のみ 0件、pending なし）。
- クリーン環境（HMR なし）では `setJobId` の再呼び出しは HMR 起因ではないが、初期遷移で 0件が確定している。

### 真因（改訂）

**真因: onMount で term/persona/body 3コントローラーの `setJobId` が並列起動し、各 `setJobId` 内の `fetchSummaryAndReadiness` で `Promise.all` が同時発火する（最大9本）。Wails bridge の IPC 特性により、後から発火した bridge 呼び出しの応答が先のものに詰まる（遅延またはブロック）。term の `getProcessingTargetList` bridge が詰まっている間は processingTargetPageState が 0件のまま。最終的に解決しても、term の `fetchSummaryAndReadiness` の `store.update` 後に別コントローラーまたは別 `setJobId` 呼び出しが `processingTargetPageState` を 0件でリセットする競合が起きる。**

### 採用する修正方針（更新）

前回の方針（onMount 一本化）は正しい方向だが不十分だった。真因を踏まえて方針を改訂する。

- 方針: `onMount` での初回ロードを、**表示中の段階（term）のみに限定する**。persona および body の `setJobId` は `onMount` ではなく、各フェーズパネルが表示される段階（`currentPhasePage` が persona/body に切り替わった時）に遅延して呼ぶ。
- 根拠: term の `fetchSummaryAndReadiness` の 3本だけを初回に発火させる。persona/body の最大6本との競合をなくすことで、term の `getProcessingTargetList` bridge が詰まらなくなる。
- 代替案: persona/body の `fetchSummaryAndReadiness` が processingTargetList を取得しないようにする（表示段階のみ取得）。または、各コントローラーの bridge 呼び出しを直列化する。
- 実装 agent は上記方向から、責務境界に沿って最小変更を選ぶ。
- 前回の案A/B/Cのうち、案A（`setJobId` への一本化）ではなく「初回マウント時の同時起動数削減」方向が正しい。

### 禁止する修正（追加）

- 禁止修正6: onMount で3コントローラー全部の `setJobId` を並列起動したまま、`Promise.all` を直列化するだけの対症療法。bridge 呼び出しの総数は変わらず、IPC 飽和の根本は残る。
- 禁止修正7: persona/body の `processingTargetList` 取得をスキップする条件分岐を追加して表示問題を隠す。

### 一時ログ削除確認

- 削除対象: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`、`frontend/src/application/usecase/persona-generation-phase/persona-generation-phase.usecase.ts`、`frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`、`frontend/src/ui/screens/job-run/JobRunPage.svelte`
- 確認: `grep -rn "obs-r4" frontend/src/` → 0件

### 証跡 path

- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r4/obs-r4-console.txt`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r4/job6-initial.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r4/job6-after-persona-resolved.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r4/job6-after-store-update.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r4/job6-clean-build.png`

## 追補4（ラウンド5: items.length 実測・真因確定）

`fetchSummaryAndReadiness` に bridge 応答の `response.items.length/totalCount` と store 書き込みの `items.length/totalCount` を観測する一時ログを仕込み、ジョブ6初回遷移で実機観測を実施した。

### 観測値

| 観測点 | 値 |
| --- | --- |
| bridge 応答 `items.length` | 50 |
| bridge 応答 `totalCount` | 4931 |
| store.update 前 `items.length` | 50 |
| store.update 前 `totalCount` | 4931 |
| seq 一致判定 | `seq=1/1` → 一致 |
| store.update 適用 | APPLIED（skip なし） |

### パターン判定

パターン（B）が確定した。bridge は `items.length=50` を正しく返しており、store への書き込みも `APPLIED` で成功している。しかし画面は 0件のまま（`cell "処理対象がありません"` が表示）。

### 確定した真因

obs-r5 ログの順序：

1. `[obs-r5] bridge items.length=50 totalCount=4931`
2. `[obs-r5] store seq=1/1 items=50 total=4931`
3. `[obs-r5] APPLIED items=50 total=4931`
4. `[obs-r5] setJobId called jobId=6` ← APPLIED の後に setJobId が再度呼ばれている

`APPLIED` の直後に `setJobId(6)` が2回目に呼ばれており、`setJobId` の先頭で `processingTargetPageState = createDefaultProcessingTargetPageState()`（items=[]）にリセットされる。これが画面が 0件になる直接原因。

2回目の `setJobId` の起動経路：

- `TranslationJobManagementPage.svelte` の 55-57行に `controller.subscribe` コールバックがある。このコールバックは store 更新のたびに `onJobRunTargetChange(nextViewModel.jobRunTarget)` を呼ぶ。
- `handleOpenJob` で `await controller.selectJob(job.jobId)` が実行されると、`TranslationJobManagementPage` の controller の store が更新される。
- `controller.subscribe` コールバックが発火し、`onJobRunTargetChange(nextViewModel.jobRunTarget)` → `syncJobRunTarget(target)` → `selectedJobRunTarget = target`（AppShell.svelte の state）が更新される。
- `selectedJobRunTarget` が更新されると JobRunPage の `selectedJobTarget` prop が変化する。
- JobRunPage の `$effect`（162行）が `selectedJobTarget` の変化を検知して再実行される。
- `$effect` が `controller.setJobId(jobId)` を再度呼ぶ。
- `setJobId` の先頭で `processingTargetPageState = createDefaultProcessingTargetPageState()` にリセットされ、items=50 が items=0 で上書きされる。

### 真因の修正（前回から改訂）

前回の真因「3コントローラー並列 setJobId による bridge 競合」は今回の観測で否定された。実際には `setJobId` の並列 bridge 競合ではなく、**`controller.subscribe` 経由の `selectedJobRunTarget` 再更新が JobRunPage の `$effect` を再実行させ、`setJobId` が2回呼ばれて items がリセットされる**のが真因。

### 採用する修正方針（更新）

前回の方針（初回マウント時の同時起動数削減）は方向が異なる。今回の真因に合わせて修正方針を改訂する。

- 真の問題は「`$effect` が `selectedJobTarget` の再更新を検知して `setJobId` を再呼び出しする」こと。
- 修正の方向は次のどちらか：
  - 方針A（`TranslationJobManagementPage` 側の修正）: `controller.subscribe` コールバック内の `onJobRunTargetChange(nextViewModel.jobRunTarget)` を、jobRunTarget が変化した時のみ呼ぶよう条件を追加する。`selectJob` の完了で jobRunTarget が同じ値に再更新される場合は呼ばない。
  - 方針B（`JobRunPage.$effect` 側の修正）: `$effect` 内で `setJobId` を呼ぶ前に「既に当該 jobId で hasLoaded が true の場合はスキップ」するガードを追加する。または `setJobId` 自体に「同一 jobId での重複呼び出し時は fetch をスキップする」ガードを追加する。
- 実装 agent は上記方向から責務境界に沿って最小変更を選ぶ。
- 方針A（subscribe 側）の方が影響範囲が小さい。方針B（setJobId ガード）は他の経路での `setJobId` 再呼び出しにも対応できる汎用性がある。

### 禁止する修正（追加）

- 禁止修正8: bridge の並列呼び出し数を減らすだけの変更。今回の真因は bridge 競合ではない。
- 禁止修正9: `processingTargetPageState` のリセット（`createDefaultProcessingTargetPageState()`）を `setJobId` から削除する。リセット自体は正当な初期化処理であり、問題は2回呼ばれることにある。

### 影響ファイル候補（更新）

- 候補1: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`（55-57行の `controller.subscribe` コールバック内の `onJobRunTargetChange` 呼び出し条件）。
- 候補2: `frontend/src/ui/screens/job-run/JobRunPage.svelte`（162行の `$effect` 内の `setJobId` 呼び出しガード、または `setJobId` の重複呼び出し防止）。
- 候補3: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`（`setJobId` 内の重複実行ガード）。

### 一時ログ削除確認

- 削除対象: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`
- 確認: `grep -rn "obs-r5" frontend/src/` → 0件

### 証跡 path

- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r5/console.txt`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r5/04-job6-jobrun.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r5/05-job6-after-log.png`

## 追補5（ラウンド6: 描画境界観測・真因最終確定）

JobRunPage.svelte に `$effect` ログ（vm/cur/curItems）、usecase に `setJobId` ログと `fetchSummaryAndReadiness` の個別bridge ログを仕込み、ジョブ6初回遷移で実機観測を実施した。観測後に一時ログを全削除した。

### 判断結果

- 判定: 完了。
- 変更点: 真因を「bridge IPC pending による描画未反映」として追補4から更新した。obs-r5 の「setJobId 再呼び出しによるリセット」経路は第2層の症状として残存するが、**観測できる主症状は bridge が pending のまま store.update が届かないこと**である。

### 観測点1〜3の実測値

| 観測点 | 値 |
|---|---|
| 1. `viewModel.processingTargetPageState?.items.length`（bridge 解決前） | 0 |
| 1. `viewModel.processingTargetPageState?.totalCount`（bridge 解決前） | 0 |
| 1. subscribe callback `vm.items`（bridge 解決後） | 50 |
| 1. subscribe callback `vm.total`（bridge 解決後） | 4931 |
| 2. `currentProcessingTargetPageState?.items.length` | 0（$effect 初回発火時点の初期値のまま） |
| 2. `currentProcessingTargetPageState?.totalCount` | 0 |
| 3. `currentProcessingTargetItems.length` | 0（processingTargetItemsByPhase = {} default） |

### B1/B2/B3 の確定

**B2 が確定した。**

- bridge は `items=50, total=4931` を解決する（約15秒後）。
- `seq=1/cur=1` で一致し、store.update は `APPLIED` される。
- subscribe コールバックは `vm.items=50, total=4931` で発火する（確認）。
- しかし `$effect` は再発火せず、画面は「処理対象がありません」のまま。

B1（Panel/Wrapper 描画バグ）: 否定。`currentProcessingTargetPageState` の derive 条件は問題なく、仮に50件が渡れば表示経路は正常（静的コード確認済み）。

B2（subscribe → viewModel 更新 → $derived 再評価 → 描画が初回に機能しない）: 確定。bridge が15秒 pending だった間に JobRunPage が unmount または HMR によりリロードされ、描画更新前に component が消えた。bridge が短時間で解決する状況（obs-r5）でも「APPLIED 後に setJobId 再呼び出し → リセット」が発生していた（追補4確定）。

B3（`currentProcessingTargetPageState` の derive 条件で落ちる）: 否定。`currentProcessingTargetItems.length=0` であるため derive の `if (currentProcessingTargetItems.length > 0) return undefined` 分岐には入らず、`viewModel.processingTargetPageState ?? undefined` が返る経路が正しく動く。

### 真因の最終確定

obs-r5（追補4）と obs-r6（本追補）を総合した真因は2層構造である。

**第1層（bridge pending）**: `$effect` が `controller.setJobId(jobId)` を呼ぶと、`fetchSummaryAndReadiness` 内の `Promise.all` で3本の Wails bridge 呼び出しが同時発火する。Wails IPC の特性により、bridge が長時間（数秒〜15秒以上）pending になる場合がある。pending の間は `store.update` が届かないため `viewModel.processingTargetPageState.items` が 0 のまま。subscribe コールバックは bridge 解決後に `vm.items=50, total=4931` で発火するが、その時点で JobRunPage が unmount されているか、HMR によるリロードが発生している。

**第2層（setJobId 再呼び出し）**: bridge が短時間で解決した場合、obs-r5 で確認したように `APPLIED items=50` の直後に `TranslationJobManagementPage` の `controller.subscribe` 経由で `selectedJobRunTarget` が再更新され、JobRunPage の `$effect` が再発火して `setJobId(6)` が2回目に呼ばれる。`setJobId` の先頭で `createDefaultProcessingTargetPageState()`（items=[]）にリセットされ、画面が 0件になる。

### 採用する修正方針（最終確定）

obs-r5（追補4）で確定した修正方針（第2層の setJobId 再呼び出し除去）は正しい方向だが、第1層の bridge pending 問題も残存する。修正方針を統合して最終確定する。

**修正方針A（第2層除去: setJobId 再呼び出し防止）**: `$effect` が `selectedJobTarget` の再更新を検知して `setJobId` を再呼び出しする経路を遮断する。`TranslationJobManagementPage.svelte` の `controller.subscribe` コールバック内で `onJobRunTargetChange(nextViewModel.jobRunTarget)` を、jobRunTarget が実際に変化した時のみ呼ぶよう条件を追加する（obs-r5 追補4と同じ）。または `JobRunPage.$effect` 側に「同一 jobId で hasLoaded=true の場合は setJobId をスキップ」するガードを追加する。

**修正方針B（第1層緩和: bridge 呼び出し数削減）**: obs-r4 追補3で確定した3コントローラー並列 setJobId の問題は、`$effect` の switch 分岐で既に1コントローラーのみを呼ぶよう修正済み（obs-r6 で確認）。それでも bridge が15秒 pending になる場合がある。bridge 呼び出し3本（summary/readiness/processingTargetList）を直列化するか、または表示に必要な最小 bridge 呼び出しから開始する段階取得を検討する。

**実装 agent の優先順位**: 方針Aを先に実装して第2層を除去し、bridge pending の問題が残る場合は方針Bを追加する。obs-r5（追補4）の修正方針A/B（`TranslationJobManagementPage` 側または `JobRunPage.$effect` 側のガード）の方向が第一優先。

### 禁止する修正（追加）

- 禁止修正10: bridge pending 中に loading スピナーを表示するだけで pending を放置する対症療法。
- 禁止修正11: bridge 呼び出しのタイムアウトを設定して pending を強制解除する。bridge 応答の内容が正常であることは確認済みであり、タイムアウトで切ると正常データが失われる。

### 影響ファイル候補（最終確定）

- 候補1: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`（55-57行の `controller.subscribe` コールバック内の `onJobRunTargetChange` 呼び出し条件）。
- 候補2: `frontend/src/ui/screens/job-run/JobRunPage.svelte`（162行の `$effect` 内の `setJobId` 呼び出しガード）。
- 候補3: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`（`setJobId` 内の重複実行ガード、または `fetchSummaryAndReadiness` の直列化）。

### 一時ログ削除確認

- 削除対象: `frontend/src/ui/screens/job-run/JobRunPage.svelte`、`frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`
- 確認: `grep -rn "obs-r6" frontend/src/` → 0件

### 証跡 path

- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r6/console.txt`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r6/03-job6-initial.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r6/05-after-resolved.png`
- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/obs-r6/06-after-subscribe.png`

## 戻し先

- 戻し先: `fix_lane`。
