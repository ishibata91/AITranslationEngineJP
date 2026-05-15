# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM-ST`

## Generator Scope

- `viewpoint`: 状態、遷移、禁止遷移、再実行条件から、翻訳ジョブ状態機械の候補だけを生成する。
- `included_sources`: `plan.md`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/screen-design/README.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、最終シナリオ表、候補採否、候補統合。
- `generation_notes`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は画面責務で分ける。大枠画面は `TRANSLATION_JOB.state` を読み、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `adopted_update`: retry、resume、pause、cancel の可否は phase type で分けない。phase type で分ける対象は、start の開始前提、完了判定、呼び出す service method だけである。

## translationjobpolicy 入出力候補

- `input`: job id、現在 job state、phase run 集約結果、現在 phase state、操作種別、開始対象 phase type、terminal 判定、active phase run 有無、phase 別開始前提、retryable 判定、同一再送判定。
- `output`: allowed / rejected、共通操作規則の結果、start 専用の phase 別開始前提結果、遷移後 job state、遷移後 phase state、継続する `JOB_PHASE_RUN` id、作成してよい phase run type、呼び出してよい service method、UI 表示向け reason category、導出してよい summary 種別。
- `persistence boundary`: policy の判断結果、rule 名、判定履歴、`PolicyResult` は永続化しない。JobIOService は確定済み状態事実だけを保存する。
- `forbidden transitions`: terminal job への後書き、表示だけによる暗黙実行開始、active phase run 併存、前段未完了 phase の開始、Running job の削除、Running body phase の直接 cancel、参照不能 job の Job Run 表示、集約不能状態での危険操作。
- `rejection reason candidates`: `job_not_ready`, `terminal_job`, `active_phase_run_exists`, `previous_phase_not_completed`, `dictionary_reference_missing`, `persona_snapshot_missing`, `input_cache_missing`, `state_inconsistent`, `phase_progress_unaggregatable`, `not_retryable_failure`, `delete_running_job`, `delete_stop_pending_job`, `cancel_requires_paused`, `display_target_not_found`, `late_response_rejected`, `duplicate_phase_run_unresolved`, `secret_value_not_allowed_in_summary`。

## Candidate Scenarios

### CAND-TJSM-ST-001 Draft job を Ready job に作成する

- `source requirement`: `docs/spec.md` は `Draft --> Ready : ジョブ作成` を定義する。`translation-job-management.md` は同一入力データに対する job 一意制約を持たず、既存 job を削除または上書きしないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-001`
- `actor`: 利用者
- `trigger`: Job Setup で翻訳対象ファイルの設定が完了し、新しい翻訳 job を作成する。
- `Policy input`: 現在状態 `Draft`、入力データ参照、既存 job 一覧、作成操作。
- `expected outcome`: 新しい job は `Ready` になる。既存 job は状態を変更されない。
- `Policy output`: allowed、遷移後 job state `Ready`、作成 job id、既存 job 変更なし。
- `forbidden transition`: 同じ入力データの既存 job を上書きして `Ready` に戻す遷移を禁止する。
- `rejection reason candidates`: 入力データ参照不能なら `input_cache_missing` または `display_target_not_found`。
- `observable point`: 新しい job は未完了一覧に表示され、同じ入力由来の既存 job は job ID と作成日時で区別できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: job 作成遷移を translationjobpolicy の対象に含めるか、JobIOService の保存前検証に寄せるかを designer が判断する。
- `conflict hint`: `docs/er.md` は job 状態を `JOB_PHASE_RUN` 群から集約すると定義する一方、`TRANSLATION_JOB.state` も ER 図に存在する。

### CAND-TJSM-ST-002 Ready job の表示では Running へ暗黙遷移しない

- `source requirement`: `translation-job-management.md` は Job Run 表示だけでは Ready job を Running へ暗黙遷移させないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-002`
- `actor`: 利用者
- `trigger`: 未完了一覧で `Ready` job を選択し、Job Run の表示対象にする。
- `Policy input`: 現在状態 `Ready`、表示操作、active phase run なし。
- `expected outcome`: job state は `Ready` のまま維持される。単語翻訳フェーズ開始操作だけが実行入口として表示される。
- `Policy output`: allowed no-op、遷移後 job state `Ready`、phase state `idle_ready` 相当、実行開始未発火。
- `forbidden transition`: 表示操作を `Ready -> Running` として扱う遷移を禁止する。
- `rejection reason candidates`: 参照不能 job なら `display_target_not_found`。
- `observable point`: Job Run は read-only の実行入口として見え、開始ボタンを押すまで phase run を作らない。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 表示 no-op を translationjobpolicy の query 判定として明示すると、UI と backend の暗黙遷移を防げる。
- `conflict hint`: なし。

### CAND-TJSM-ST-003 Ready job から単語翻訳フェーズを開始する

- `source requirement`: `docs/spec.md` は `Ready --> Running : 実行開始` を定義する。`term-translation-phase.md` は Ready job かつ active な単語翻訳 phase run がない時だけ開始できると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-003`
- `actor`: 利用者
- `trigger`: Job Run で単語翻訳フェーズ開始を実行する。
- `Policy input`: job state `Ready`、開始対象 `term-translation`、active phase run なし、terminal false。
- `expected outcome`: job は `Running` になり、単語翻訳 `JOB_PHASE_RUN` は実行対象になる。
- `Policy output`: allowed、遷移後 job state `Running`、作成または継続する単語翻訳 phase run、phase state `Running`。
- `forbidden transition`: Ready 以外、terminal job、既存 active phase run ありの開始を禁止する。
- `rejection reason candidates`: `job_not_ready`, `terminal_job`, `active_phase_run_exists`。
- `observable point`: Job Run は current phase、phase state、progress、開始時刻を表示する。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: translationjobpolicy は provider 設定値を保存せず、開始可否と遷移結果だけを返す候補にする。
- `conflict hint`: `term-translation-phase.md` は job を Running のまま維持し phase state で完了や中断を区別すると定義する。job state 集約方式との整理が必要である。

### CAND-TJSM-ST-004 単語翻訳フェーズ完了後だけ NPC ペルソナ生成フェーズを開始する

- `source requirement`: `persona-generation-phase.md` は単語翻訳フェーズ Completed、非 terminal job、active phase run なしの場合だけ開始できると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-004`
- `actor`: 利用者
- `trigger`: Job Run で NPC ペルソナ生成フェーズ開始を実行する。
- `Policy input`: job state 非 terminal、単語翻訳 phase state `Completed`、辞書参照成立、active phase run なし。
- `expected outcome`: job は `Running` を維持し、NPC ペルソナ生成 phase run は実行対象になる。
- `Policy output`: allowed、遷移後 job state `Running`、作成または継続する persona phase run、phase state `Running`。
- `forbidden transition`: 単語翻訳未完了、単語翻訳失敗中、辞書参照不能、terminal job、active phase run ありの開始を禁止する。
- `rejection reason candidates`: `previous_phase_not_completed`, `dictionary_reference_missing`, `terminal_job`, `active_phase_run_exists`。
- `observable point`: Job Run は current phase として `NPC ペルソナ生成`、phase state、body phase readiness false を表示する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 前段 phase 完了と参照成立を Policy input に分けると、拒否理由を UI に近接表示できる。
- `conflict hint`: job state が `Running` のままでも、phase state が Completed の時に次 phase start を許可する集約規則が必要である。

### CAND-TJSM-ST-005 NPC ペルソナ生成フェーズ完了後だけ本文翻訳フェーズを開始する

- `source requirement`: `body-translation-phase.md` は NPC ペルソナ生成フェーズ Completed、非 terminal job、active phase run なし、辞書と persona snapshot の参照成立を開始条件にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-005`
- `actor`: 利用者
- `trigger`: Job Run で本文翻訳フェーズ開始を実行する。
- `Policy input`: job state 非 terminal、persona phase state `Completed`、辞書参照成立、persona snapshot 参照成立、active phase run なし。
- `expected outcome`: job は `Running` を維持し、本文翻訳 phase run は実行対象になる。
- `Policy output`: allowed、遷移後 job state `Running`、作成または継続する body phase run、phase state `Running`。
- `forbidden transition`: persona 未完了、snapshot 参照不能、辞書参照不能、terminal job、active phase run ありの開始を禁止する。
- `rejection reason candidates`: `previous_phase_not_completed`, `persona_snapshot_missing`, `dictionary_reference_missing`, `terminal_job`, `active_phase_run_exists`。
- `observable point`: Job Run は body phase readiness、対象 field 件数、provider / model 要約、credential 状態分類を表示する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: body phase start の拒否理由は、persona と辞書を別カテゴリに分ける候補にする。
- `conflict hint`: `docs/spec.md` の job state 図だけでは phase 連鎖の前提が表現されない。

### CAND-TJSM-ST-006 Running phase を Paused に中断する

- `source requirement`: `docs/spec.md` は `Running --> Paused : 中断` を定義する。persona と body の詳細仕様は pause を Running の時だけ有効にすると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-006`
- `actor`: 利用者
- `trigger`: Job Run で実行中 phase の pause を実行する。
- `Policy input`: job state `Running`、current phase state `Running`、active phase run あり。
- `expected outcome`: job または表示上の集約状態は `Paused` になる。current phase state は `Paused` になる。
- `Policy output`: allowed、遷移後 job state `Paused`、遷移後 phase state `Paused`、継続する `JOB_PHASE_RUN` id。
- `forbidden transition`: Running 以外の phase に pause を適用する遷移を禁止する。
- `rejection reason candidates`: `job_not_ready`, `state_inconsistent`, `terminal_job`。
- `observable point`: Job Run は phase state `Paused`、resume 入口、必要なら cancel 入口を表示する。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: pause は provider 停止方式ではなく状態遷移可否として候補化する。
- `conflict hint`: 実行中通信の止め方は `translation-job-management.md` の対象外である。

### CAND-TJSM-ST-007 Paused phase を同じ JOB_PHASE_RUN で再開する

- `source requirement`: `docs/spec.md` は `Paused --> Running : 再開` を定義する。各 phase 詳細仕様は再開で同じ `JOB_PHASE_RUN` を継続すると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-007`
- `actor`: 利用者
- `trigger`: Job Run で paused phase の resume を実行する。
- `Policy input`: job state `Paused`、current phase state `Paused`、同一 `JOB_PHASE_RUN` id、terminal false。
- `expected outcome`: job は `Running` に戻り、同じ `JOB_PHASE_RUN` が継続される。
- `Policy output`: allowed、遷移後 job state `Running`、遷移後 phase state `Running`、継続 phase run id。
- `forbidden transition`: 新しい `JOB_PHASE_RUN` 作成、terminal job の resume、参照不能 job の resume を禁止する。
- `rejection reason candidates`: `terminal_job`, `display_target_not_found`, `input_cache_missing`, `state_inconsistent`。
- `observable point`: Job Run は progress を継続値として表示し、重複 phase run を表示しない。
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: resume は新規開始ではなく継続遷移として固定する候補にする。
- `resolved conflict`: persona phase 固有の RecoverableFailed resume は採用しない。resume は Paused だけを許可する。

### CAND-TJSM-ST-008 Retryable failure を RecoverableFailed にし、retry で Running に戻す

- `source requirement`: `docs/spec.md` は `Running --> RecoverableFailed` と `RecoverableFailed --> Running` を定義する。各 phase 詳細仕様は retry で同じ `JOB_PHASE_RUN` を継続すると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-008`
- `actor`: 利用者
- `trigger`: provider 失敗、応答不正、保存失敗、部分失敗後に retry を実行する。
- `Policy input`: current phase state `RecoverableFailed`、retryable true、terminal false、同一 `JOB_PHASE_RUN` id。
- `expected outcome`: job は `Running` に戻り、同じ phase run で未処理対象だけ処理される。
- `Policy output`: allowed、遷移後 job state `Running`、遷移後 phase state `Running`、継続 phase run id。
- `forbidden transition`: retryable false、terminal job、新規 phase run 作成、成功済み結果の重複作成を禁止する。
- `rejection reason candidates`: `not_retryable_failure`, `terminal_job`, `duplicate_phase_run_unresolved`, `state_inconsistent`。
- `observable point`: Job Run は retryable flag、短い error kind、progress 継続、成功済み件数維持を表示する。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `冪等性_requirement`
- `adoption hint`: retry の戻り先は job `Running` と phase `Running` に揃える候補にする。
- `resolved conflict`: RecoverableFailed は retry だけを許可する。retry と resume の可否は phase type で分けない。

### CAND-TJSM-ST-009 回復不能失敗では Failed にし、再開入口を出さない

- `source requirement`: `docs/spec.md` は `Running --> Failed : 回復不能な失敗` を定義する。detail specs は provider 失敗、応答不正、保存失敗を successful Completed として扱わないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-009`
- `actor`: システム
- `trigger`: retryable ではない failure が発生する。
- `Policy input`: current phase state failure、retryable false、回復不能分類、terminal 判定候補。
- `expected outcome`: job は `Failed` になる候補を持ち、resume / retry 入口は無効になる。
- `Policy output`: allowed failure transition、遷移後 job state `Failed`、危険操作 disabled、理由カテゴリ表示。
- `forbidden transition`: 回復不能失敗を `Completed` として扱う遷移、retryable false の retry を禁止する。
- `rejection reason candidates`: retry 操作に対して `not_retryable_failure`。
- `observable point`: 未完了一覧または Job Run は failure state、error kind、retryable false、再開不可理由を表示する。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `adoption hint`: Failed が terminal job に含まれるかを明示する必要がある。
- `conflict hint`: `translation-job-management.md` は Failed job を未完了一覧の対象にする。terminal state と未完了一覧対象の関係を designer が整理する必要がある。

### CAND-TJSM-ST-010 本文翻訳フェーズ完了で job を Completed にする

- `source requirement`: `docs/spec.md` は `Running --> Completed : 翻訳完了` を定義する。`body-translation-phase.md` は本文翻訳フェーズ完了時点で翻訳 job 全体が `Completed` になると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-010`
- `actor`: システム
- `trigger`: body phase が Completed になり、field result 整合、output status 整合、output readiness が成立する。
- `Policy input`: body phase state `Completed`、field result 整合 true、output status 整合 true、output readiness true。
- `expected outcome`: job は `Completed` になる。Completed job は未完了一覧から外れる。
- `Policy output`: allowed、遷移後 job state `Completed`、terminal candidate true、未完了一覧対象 false。
- `forbidden transition`: 保存失敗、検証失敗、partial state、output readiness false の `Completed` 化を禁止する。
- `rejection reason candidates`: `state_inconsistent`, `previous_phase_not_completed`, `phase_progress_unaggregatable`。
- `observable point`: Job Run は Completed、output readiness、成果物出力へ渡せる summary を表示する。未完了一覧には表示しない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: job Completed は body phase だけから導出する候補にする。
- `conflict hint`: 単語翻訳対象 0 件や persona 生成対象 0 件は phase Completed だが job Completed ではない。

### CAND-TJSM-ST-011 Ready または Paused から cancel し、Running からの直接 cancel は拒否する

- `source requirement`: `docs/spec.md` は `Ready --> Canceled` と `Paused --> Canceled` を定義する。`body-translation-phase.md` は cancel を `Paused` からだけ可能にし、Running から直接 cancel しないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-011`
- `actor`: 利用者
- `trigger`: Job Run で cancel を実行する。
- `Policy input`: job state `Ready` または `Paused`、current phase state、cancel 操作。
- `expected outcome`: 許可状態では job または phase は `Canceled` になる。Running からの cancel は拒否される。
- `Policy output`: allowed なら遷移後 job state `Canceled`、rejected なら状態不変と reason category。
- `forbidden transition`: Running body phase から直接 `Canceled` へ遷移することを禁止する。
- `rejection reason candidates`: `cancel_requires_paused`, `terminal_job`, `state_inconsistent`。
- `observable point`: Job Run は cancel 無効理由を近接表示し、Canceled 後は output readiness を使わない。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: Ready cancel を job-level cancel として残すか、phase-level cancel と分けるかを designer が判断する。
- `conflict hint`: spec は Ready cancel を許可するが、phase 詳細仕様は主に Paused cancel を扱う。

### CAND-TJSM-ST-012 Running job と停止要求中 job の削除を拒否する

- `source requirement`: `translation-job-management.md` は Running job を削除できず、停止要求中も削除できないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-012`
- `actor`: 利用者
- `trigger`: 未完了一覧または Job Run で job 削除を実行する。
- `Policy input`: job state、停止要求中 flag、active phase run 有無、削除操作。
- `expected outcome`: Running または停止要求中の job は削除拒否になる。非実行中 job の削除は input data と抽出 JSON 正本を残す。
- `Policy output`: rejected なら状態不変と reason category、allowed なら job 削除可と入力データ保持。
- `forbidden transition`: Running job の削除、停止要求中 job の削除、input data の連鎖削除を禁止する。
- `rejection reason candidates`: `delete_running_job`, `delete_stop_pending_job`, `state_inconsistent`。
- `observable point`: UI は削除拒否理由と停止入口を表示する。削除成功後、対象 job は未完了一覧から外れる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `authorization_requirement`
- `adoption hint`: 削除は job state 遷移ではないが、translationjobpolicy が operation availability を返す候補に含める。
- `conflict hint`: 削除後の履歴保持と監査表示は対象外である。

### CAND-TJSM-ST-013 terminal job への phase run 作成、保存、readiness 更新、後書きを拒否する

- `source requirement`: term、persona、body の詳細仕様は terminal job への後書きや phase run 作成を拒否し、既存 state を変更しないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-013`
- `actor`: システム
- `trigger`: terminal job に対して phase start、save、readiness update、late response write が要求される。
- `Policy input`: terminal true、操作種別、現在 job state、現在 phase state。
- `expected outcome`: 操作は拒否され、job state と phase state は変更されない。
- `Policy output`: rejected、状態不変、reason category `terminal_job` または `late_response_rejected`。
- `forbidden transition`: `Completed`、`Failed`、`Canceled` など terminal candidate から `Running` や phase state 更新へ戻す遷移を禁止する。
- `rejection reason candidates`: `terminal_job`, `late_response_rejected`。
- `observable point`: late response は保存されず、UI と summary は terminal 後の結果を成功値として表示しない。
- `related detail requirement type`: `state_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: terminal 判定は Policy input に明示する候補にする。
- `conflict hint`: Failed と Canceled は未完了一覧対象でもあるため、terminal 判定と一覧表示対象は別概念にする必要がある。

### CAND-TJSM-ST-014 開始再送、再開、リトライでは同じ JOB_PHASE_RUN を継続する

- `source requirement`: term、persona、body の詳細仕様は再送、再開、リトライで同じ `JOB_PHASE_RUN` を継続し、重複作成しないと定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-014`
- `actor`: 利用者またはシステム
- `trigger`: 同じ phase 開始、resume、retry が再送される。
- `Policy input`: 対象 phase type、既存 `JOB_PHASE_RUN` id、既存 phase state、同一再送判定、active phase run 有無。
- `expected outcome`: 同一再送と判定できる場合は既存 phase run を継続する。判定不能または別 active phase がある場合は新規作成を拒否する。
- `Policy output`: allowed continue または rejected、継続 phase run id、重複作成なし。
- `forbidden transition`: 同じ phase の重複 `JOB_PHASE_RUN` 作成、成功済み `JOB_TRANSLATION_FIELD` や `PHASE_RUN_*` の重複作成を禁止する。
- `rejection reason candidates`: `active_phase_run_exists`, `duplicate_phase_run_unresolved`, `state_inconsistent`。
- `observable point`: progress、成功済み件数、failed count は継続される。DB 上の phase run は増えない。
- `related detail requirement type`: `冪等性_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: 同一再送判定は translationjobpolicy の入力にする。JobIOService はその判定結果ではなく、確定済み状態事実だけを保存する。
- `conflict hint`: 単語翻訳フェーズ開始条件は active phase run なしを要求するため、開始再送を継続扱いにする条件を明示する必要がある。

### CAND-TJSM-ST-015 集約不能または状態不整合では危険操作を無効化し、状態を変えない

- `source requirement`: `translation-job-management.md` は phase progress 集約不能を成功値として表示せず、危険操作を無効にすると定義する。入力キャッシュ欠落、terminal state、状態不整合では再開不可理由を表示すると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-015`
- `actor`: 利用者
- `trigger`: 未完了一覧または Job Run で、集約不能 job の操作可否を確認する。
- `Policy input`: phase progress 集約結果 failure、入力キャッシュ状態、job state、phase run 集約状態。
- `expected outcome`: 再開、削除、開始などの危険操作は安全側に disabled になる。理由表示だけでは job state を変えない。
- `Policy output`: rejected or disabled、状態不変、reason category。
- `forbidden transition`: 集約不能を Completed や Running の成功値として扱う遷移を禁止する。
- `rejection reason candidates`: `phase_progress_unaggregatable`, `state_inconsistent`, `input_cache_missing`。
- `observable point`: UI は空一覧や成功状態と区別して、再開不可理由または無効理由を表示する。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: translationjobpolicy は状態を変えず、操作可否と reason category を返す候補にする。
- `conflict hint`: なし。

### CAND-TJSM-ST-016 対象 0 件 phase は Completed にするが job 完了条件は phase ごとに分ける

- `source requirement`: term、persona、body の詳細仕様は対象 0 件を Completed として扱う。body phase の 0 件 Completed では成果物出力へ進めると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-016`
- `actor`: システム
- `trigger`: phase 開始後、provider 実行対象が 0 件と判定される。
- `Policy input`: 対象 phase type、target count 0、provider 未実行、前段参照成立。
- `expected outcome`: current phase は `Completed` になる。term と persona では次 phase readiness を評価する。body では job `Completed` と output readiness を評価する。
- `Policy output`: allowed、phase state `Completed`、phase type ごとの次状態候補、provider 未実行 summary。
- `forbidden transition`: term または persona の 0 件 Completed だけで job 全体を `Completed` にする遷移を禁止する。
- `rejection reason candidates`: `previous_phase_not_completed`, `dictionary_reference_missing`, `persona_snapshot_missing`, `state_inconsistent`。
- `observable point`: Job Run は `empty_completed` 相当、provider 未実行、次 phase readiness または output readiness を表示する。
- `related detail requirement type`: `boundary_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 0 件完了は phase state と job state の集約規則を分ける候補にする。
- `conflict hint`: job state だけの状態図では、phase ごとの 0 件完了後の分岐を表現できない。

### CAND-TJSM-ST-017 RecoverableFailed から Ready への再実行準備は不採用にする

- `source requirement`: `docs/spec.md` は `RecoverableFailed --> Ready : 再実行準備` を定義する。一方、phase 詳細仕様は retry、resume、開始再送で同じ `JOB_PHASE_RUN` を継続すると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSM-ST-017`
- `actor`: 利用者
- `trigger`: RecoverableFailed job に対して再実行準備を行う。
- `Policy input`: job state `RecoverableFailed`、current phase state `RecoverableFailed`、再実行準備操作。
- `expected outcome`: `Ready` へ戻す経路は採用しない。同じ `JOB_PHASE_RUN` を retry で `Running` へ戻す。
- `Policy output`: rejected、状態不変、reason category `recoverable_failed_ready_removed`。
- `forbidden transition`: `RecoverableFailed -> Ready` を確定仕様として採用することを禁止する。
- `rejection reason candidates`: `state_inconsistent`, `duplicate_phase_run_unresolved`。
- `observable point`: Job Run では retry 入口だけを表示し、同じ phase run と progress が残る。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `冪等性_requirement`
- `adoption hint`: `docs/spec.md` の再実行準備を docs 更新対象として削除する候補にする。
- `resolved conflict`: `RecoverableFailed -> Ready` は詳細仕様の同一 `JOB_PHASE_RUN` 継続と衝突するため廃止する。

## Open Notes

- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `resolved decision`: terminal guard と未完了一覧対象は別概念として扱う。
- `resolved decision`: `RecoverableFailed -> Ready` は廃止し、retry / resume は同一 phase run を継続する。
- `resolved decision`: Policy output は operation availability、拒否理由、状態作用を含める。
- `resolved decision`: Ready job の cancel は job-level 操作として維持し、phase 開始後の cancel は Paused phase だけに限定する。
- `resolved decision`: retry、resume、pause、cancel の可否は phase type で分けない。
- `merge candidate`: CAND-TJSM-ST-004 と CAND-TJSM-ST-005 は、前段 phase 完了と参照成立による次 phase start guard として統合できる。
- `merge candidate`: CAND-TJSM-ST-007、CAND-TJSM-ST-008、CAND-TJSM-ST-014 は、同一 `JOB_PHASE_RUN` 継続と重複作成防止として統合できる。
- `rejection candidate`: provider 実装方式、credential 実値、API key、provider raw response を Policy input / output に含める候補は除外する。
- `candidate count`: 17
