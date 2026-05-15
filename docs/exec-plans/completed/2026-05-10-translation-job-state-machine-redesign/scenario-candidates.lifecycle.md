# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: product code, product test, docs canonical update, final scenario adoption, scenario merge decision
- `generation_notes`: 候補は lifecycle 観点に限定する。job state、phase run state、保存値、集約値、表示値の混同は競合候補として残す。最終シナリオ表、採否、統合、競合解消は designer が扱う。
- `adopted_update`: 人間回答により、retry、resume、pause、cancel の可否は phase type で分けない。phase type で分ける対象は、start の開始前提、完了判定、呼び出す service method だけである。

## State Vocabulary Notes

- `job state`: `docs/spec.md` は `Draft`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` を翻訳ジョブ状態として列挙している。
- `phase run state`: 各 phase 詳細仕様は、開始、pause、resume、retry、cancel、complete を `JOB_PHASE_RUN` の状態と操作可否で扱う。
- `保存値`: 大枠の一覧、導線、ジョブ全体の表示は `TRANSLATION_JOB.state` を正本にする。各フェーズ画面の操作可否は `JOB_PHASE_RUN.state` を正本にする。
- `集約値`: 未完了一覧、Job Run、Output Review は phase progress、active phase run、body phase Completed、field result 整合から job-level の状態または readiness を組み立てる。
- `表示値`: UI 詳細仕様は `idle_ready`、`running`、`empty_completed`、`recoverable_failed`、`blocked`、`not-ready`、`ready`、`starting` などの表示差分を持つ。表示値は保存値または集約値と同一視しない。

## Candidate Scenarios

### CAND-TJSM-001 翻訳 job 作成で Ready job を作る

- `source requirement`: `docs/spec.md` 1、6、7。`docs/detail-specs/translation-job-setup.md` の job 作成仕様。`docs/er.md` の 1 job 1 input と phase 設定保持。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-001`
- `actor`: 翻訳 job を作成する利用者
- `trigger`: 登録済み入力データと 3 つの翻訳段階の AI 設定がそろった状態で job 作成を実行する。
- `expected outcome`: Ready job が作成される。入力データと抽出 JSON 正本は保持される。job 作成だけでは phase run を実行中にしない。
- `observable point`: Job Setup の作成後 summary、未完了一覧の Ready job、Job Run の read-only 実行入口、保存された provider / model / execution mode / batch mode / credential 状態分類。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 作成直後の状態正本を確認する主候補になる。
- `resolved conflict`: `Ready` job には `JOB_PHASE_RUN` を事前作成しない。大枠画面は `TRANSLATION_JOB.state` を読む。

### CAND-TJSM-002 Ready job から単語翻訳フェーズを開始する

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/term-translation-phase.md` の Ready job 開始条件。`docs/detail-specs/translation-job-management.md` の Ready job は暗黙遷移しない仕様。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-002`
- `actor`: 翻訳 job を実行する利用者
- `trigger`: Ready job で active な単語翻訳 phase run が存在しない時に、Job Run から単語翻訳フェーズ開始を実行する。
- `expected outcome`: 単語翻訳用の `JOB_PHASE_RUN` が開始される。job-level は Running として観測できる。phase 開始時に最新 endpoint と credential 参照状態を再解決するが、secret 実値は表示または保存 summary に出さない。
- `observable point`: Job Run の current phase、phase state、progress、開始時刻、provider / model / execution mode / batch mode / credential 状態分類、開始ボタンの無効理由。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: Ready から Running へ進む最初の lifecycle 候補になる。
- `conflict hint`: `docs/spec.md` は job state の直接遷移を描くが、`docs/er.md` は job state を phase run 群から集約するとする。開始時の保存値と集約値の分離が必要である。

### CAND-TJSM-003 完了済み単語翻訳フェーズから NPC ペルソナ生成フェーズを開始する

- `source requirement`: `docs/spec.md` 6。`docs/detail-specs/term-translation-phase.md` の後続 phase 条件。`docs/detail-specs/persona-generation-phase.md` の開始条件。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-003`
- `actor`: 翻訳 job を進める利用者
- `trigger`: 単語翻訳フェーズが Completed であり、ジョブ内辞書参照が成立し、terminal job ではなく、active phase run が存在しない時に NPC ペルソナ生成を開始する。
- `expected outcome`: NPC ペルソナ生成用の `JOB_PHASE_RUN` が開始される。単語翻訳の完了結果は後続 phase の入力 summary として使われる。job-level の進行は Running として観測できる。
- `observable point`: current phase、phase state、単語翻訳結果 summary、辞書参照成立、active phase run 不在、NPC ペルソナ生成の readiness。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: phase 間遷移の lifecycle 候補になる。
- `conflict hint`: 単語翻訳 phase state の `Completed` と、job-level 表示の `Running` を同一状態として扱わない必要がある。

### CAND-TJSM-004 完了済み NPC ペルソナ生成フェーズから本文翻訳フェーズを開始する

- `source requirement`: `docs/spec.md` 6。`docs/detail-specs/persona-generation-phase.md` の body phase readiness。`docs/detail-specs/body-translation-phase.md` の開始条件。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-004`
- `actor`: 本文翻訳を開始する利用者
- `trigger`: NPC ペルソナ生成フェーズが Completed であり、persona snapshot 参照、辞書参照、metadata 参照、active phase run 不在が成立した時に本文翻訳フェーズを開始する。
- `expected outcome`: 本文翻訳用の `JOB_PHASE_RUN` が開始される。本文翻訳入力 snapshot が固定される。job-level は terminal ではない Running として扱われる。
- `observable point`: body phase readiness、辞書 snapshot digest、persona snapshot digest、metadata digest、prompt digest、current phase、開始ボタンの有効条件。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 最終 phase 開始までの lifecycle をつなぐ候補になる。
- `conflict hint`: readiness は表示値または集約値であり、保存済み phase state と混同しない。

### CAND-TJSM-005 Running phase を pause して再開可能な停止状態にする

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/persona-generation-phase.md` と `docs/detail-specs/body-translation-phase.md` の pause 条件。`docs/detail-specs/translation-job-management.md` の停止入口。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-005`
- `actor`: 実行中 job を一時停止する利用者
- `trigger`: 対象 phase run が Running の時に pause を実行する。
- `expected outcome`: 対象 phase run は Paused として観測できる。job-level は Paused として集約または表示される。途中まで成功した結果は phase 詳細仕様の範囲で維持される。
- `observable point`: phase state、job state 表示、progress、再開入口、削除可否の再判定、停止要求中の削除不可理由。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `consistency_requirement`, `concurrency_requirement`
- `adoption hint`: pause と delete 可否の lifecycle 境界を確認する候補になる。
- `conflict hint`: Running job は削除不可である。停止要求中、Paused 到達後、削除可能判定の境界が未確定である可能性がある。

### CAND-TJSM-006 Paused phase を resume する

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/persona-generation-phase.md` の resume 条件。`docs/detail-specs/body-translation-phase.md` の resume 条件。`docs/detail-specs/translation-job-management.md` の再開入口。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-006`
- `actor`: 中断または回復可能失敗から処理を戻す利用者
- `trigger`: Paused の phase run で resume を実行する。
- `expected outcome`: 同じ `JOB_PHASE_RUN` が Running へ戻る。成功済みの辞書、persona、field result は重複作成されない。入力キャッシュ欠落、terminal state、状態不整合では再開不可理由だけを表示し、状態を変えない。
- `observable point`: phase state、retryable flag、再開不可理由カテゴリ、既存 result 数、重複作成なし、redacted phase result summary。
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: resume は phase type ではなく `JOB_PHASE_RUN.state=Paused` で許可する共通操作規則の候補になる。
- `resolved conflict`: persona phase 固有の RecoverableFailed resume は採用しない。RecoverableFailed は retry だけを許可する。

### CAND-TJSM-007 RecoverableFailed phase を retry する

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/term-translation-phase.md`、`docs/detail-specs/persona-generation-phase.md`、`docs/detail-specs/body-translation-phase.md` の retry 条件。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-007`
- `actor`: 回復可能失敗を再試行する利用者
- `trigger`: retryable failure が存在する RecoverableFailed phase run で retry を実行する。
- `expected outcome`: 同じ `JOB_PHASE_RUN` を継続する。未処理 term、未処理 NPC、失敗 field だけを再処理対象にする。phase 開始時と同じく最新 endpoint と credential 参照状態を再解決する。
- `observable point`: retryable flag、error kind、progress、既存成功結果の維持、重複 entry / persona / field result なし、provider / model / credential 状態分類。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `state_requirement`, `data_requirement`, `security_requirement`
- `adoption hint`: retry と開始再送の同一 phase run 継続を検証する候補になる。
- `conflict hint`: retry が Running へ戻す transition 結果なのか、phase run 内の attempt 更新なのかは designer の整理対象である。

### CAND-TJSM-008 Paused phase を cancel して終端状態にする

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/body-translation-phase.md` の cancel 条件。`docs/detail-specs/translation-job-management.md` の terminal state と削除可否。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-008`
- `actor`: 処理継続をやめる利用者
- `trigger`: Paused の phase run で cancel を実行する。
- `expected outcome`: phase run または job-level は Canceled として終端になる。terminal job では後続 phase run 作成、field save、readiness update、late response 後書きを拒否する。途中成功結果は output readiness に使わない。
- `observable point`: Canceled 表示、後続操作の無効理由、late response rejected、output readiness false、未完了一覧での表示有無。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: cancel 後の終端性と late response 拒否を確認する候補になる。
- `conflict hint`: `docs/spec.md` は Ready から Canceled への cancel を含むが、body phase 詳細仕様は Paused からだけ cancel 可能にする。job-level cancel と phase-level cancel の境界は未解決である。

### CAND-TJSM-009 本文翻訳フェーズ完了で job を Completed にする

- `source requirement`: `docs/spec.md` 7。`docs/detail-specs/body-translation-phase.md` の job 完了条件。`docs/detail-specs/translation-output-artifact.md` の出力開始条件。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-009`
- `actor`: 翻訳結果を確認して成果物出力へ進みたい利用者
- `trigger`: 本文翻訳フェーズが Completed になり、field result 整合、output status 整合、output readiness が成立する。
- `expected outcome`: job-level 状態は Completed として観測できる。Completed job は未完了一覧に表示されない。Output Review の completed job list と artifact 生成候補に出る。
- `observable point`: body phase Completed、job-level Completed、output readiness true、completed job list、未完了一覧からの除外、field result summary。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: job lifecycle の正常終点候補になる。
- `conflict hint`: body phase state、job-level 集約値、Output Review の表示値を同一保存値として扱わない必要がある。

### CAND-TJSM-010 非実行中 job を delete して input data を保持する

- `source requirement`: `docs/spec.md` 4。`docs/detail-specs/translation-job-management.md` の削除仕様。`docs/detail-specs/translation-job-setup.md` の input 削除との分離。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSM-010`
- `actor`: 不要な未完了 job を整理する利用者
- `trigger`: Running ではなく、停止要求中でもない job に対して Job Management から削除を実行する。
- `expected outcome`: 対象 job は未完了一覧から外れる。input data と抽出 JSON 正本は残る。Job Setup の input 削除とは別操作として扱われる。
- `observable point`: 削除ボタンの有効条件、削除拒否理由、削除後の未完了一覧、input data の残存、Job Setup 側の入力候補または既存 job summary。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: lifecycle の終了後整理候補になる。
- `conflict hint`: Canceled、Failed、RecoverableFailed、Paused、Ready の削除可否を同一規則にするか、状態別に理由を分けるかは最終設計で確認が必要である。

## Conflict Candidates

- `CONF-TJSM-001`: 解消済み。大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `CONF-TJSM-002`: 解消済み。`Ready` job は active ではない phase run を持たない。
- `CONF-TJSM-003`: 解消済み。`pending` は公開状態にしない。
- `CONF-TJSM-004`: 解消済み。resume は Paused だけ、retry は RecoverableFailed だけにする。
- `CONF-TJSM-005`: `docs/spec.md` は Ready から Canceled への cancel を含むが、body phase は Paused からだけ cancel 可能にする。
- `CONF-TJSM-006`: Completed job は未完了一覧から除外されるが、Output Review では completed job list の対象になる。表示場所ごとの表示値を分ける必要がある。

## Open Notes

- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `resolved decision`: `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
- `resolved decision`: `pending` は公開状態にしない。
- `resolved decision`: `Ready` cancel は job-level、phase 開始後 cancel は Paused phase からだけ許可する。
- `resolved decision`: resume と retry は phase 共通規則にする。
- `merge candidate`: CAND-TJSM-003 と CAND-TJSM-004 は phase 間開始 lifecycle として統合できる可能性がある。
- `merge candidate`: CAND-TJSM-006 と CAND-TJSM-007 は回復 lifecycle として統合できる可能性がある。
- `rejection candidate`: なし。designer が最終採否を判断する。
