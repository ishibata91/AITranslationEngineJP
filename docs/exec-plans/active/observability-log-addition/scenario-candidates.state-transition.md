# Scenario Candidates: observability-log-addition / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OBS-ST`

## Generator Scope

- `viewpoint`: 状態遷移
- `included_sources`: `plan.md`, `docs/observability-logging.md`, `docs/architecture.md`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、`.codex/` 変更、他 agent の候補成果物
- `generation_notes`: 変更前、変更後、拒否理由、遷移理由を同じ境界で取れる候補だけを残す。最終採否、統合、競合解消は `designer` が扱う。
- `candidate_count`: 8

## Candidate Scenarios

### CAND-OBS-ST-001 翻訳ジョブ削除の許可と拒否を観測する

- `source requirement`: 観測ログ仕様 4 は、状態遷移の変更前、変更後、拒否理由が同じ場所で取れる箇所を追加対象にする。要件一覧 7 は、`Ready`、`Paused`、`RecoverableFailed`、`Completed`、`Canceled` などのジョブ状態を定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-001`
- `actor`: backend service
- `trigger`: 利用者が未完了ジョブ一覧からジョブ削除を要求する。
- `previous state`: 対象ジョブは `ready`、`paused`、`recoverable_failed`、`running`、`completed` のいずれかで永続化されている可能性がある。
- `start condition`: 削除候補を読み、削除可否を判定できる。
- `transition reason`: 利用者の削除要求である。
- `expected outcome`: 削除可能なジョブは削除済みになる。`running` は削除されず、`running_delete_blocked` を返す。`completed` または存在しないジョブは `stale_selection` を返す。
- `observable point`: `event=job_delete_decision`, `where=backend.service.translation_job_management`, `result=deleted|rejected`, `id=job:<id>`, `reason=running_delete_blocked|stale_selection|delete_failed`, `before_state`, `after_state` を残す候補にする。
- `observed disappearing information`: 削除判定後に消える削除前ジョブ状態、削除拒否の理由分類、一覧再読込前の選択 stale 判定。
- `forbidden log`: 入力ファイル全文、抽出 JSON 全文、source path の全文、関連 phase run の全 payload は出さない。
- `related detail requirement type`: `state_requirement`, `observability_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 削除処理は拒否理由と削除結果を同じ処理境界で持つため、状態遷移ログの優先候補にする。
- `conflict hint`: 削除ログを監査履歴として永続保存するかどうかは未決である。ログ仕様は stderr / browser console を出力先にしている。

### CAND-OBS-ST-002 翻訳ジョブ停止と再開の可否を観測する

- `source requirement`: 要件一覧 4 は、翻訳ジョブの中断、再開、失敗回復を継続的に行えることを求める。要件一覧 7 は、`Running` から `Paused`、`Paused` から `Running`、`RecoverableFailed` から `Running` への遷移を定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-002`
- `actor`: backend service
- `trigger`: 利用者が未完了ジョブ一覧またはジョブ詳細から停止または再開を要求する。
- `previous state`: 対象ジョブは `running`、`paused`、`recoverable_failed`、`failed`、`canceled` のいずれかで永続化されている可能性がある。
- `start condition`: 停止または再開の対象ジョブと現在 phase run を読み込める。
- `transition reason`: 利用者の停止要求または再開要求である。
- `expected outcome`: 許可状態ではジョブ状態または phase run 状態が遷移する。終端状態や不整合状態では状態を変えず、`terminal_state`、`resume_failed`、`stop_failed`、`state_projection_inconsistent` などの理由を返す。
- `observable point`: `event=job_operation_transition`, `where=backend.service.translation_job_management`, `result=transitioned|rejected|failed`, `id=job:<id>`, `reason`, `before_state`, `after_state`, `operation=stop|resume` を残す候補にする。
- `observed disappearing information`: UI 再読込後に消える操作時点のジョブ状態、拒否分類、状態投影の不整合分類。
- `forbidden log`: phase run の `latest_error` 全文、provider raw payload、翻訳本文全文、prompt 全文は出さない。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `observability_requirement`, `consistency_requirement`
- `adoption hint`: 停止と再開は人間操作の状態遷移であり、拒否理由が利用者の次操作に直結する。
- `conflict hint`: job state と phase run state のどちらを `after_state` の主語にするかは、designer が統合時に決める必要がある。

### CAND-OBS-ST-003 フェーズ開始の許可遷移を観測する

- `source requirement`: ER 正本は `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` を保持する。観測ログ仕様 4 は、変更前と変更後を同じ場所で取れる状態遷移を追加対象にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-003`
- `actor`: backend service
- `trigger`: 用語翻訳、NPC ペルソナ生成、本文翻訳のいずれかのフェーズ開始が要求される。
- `previous state`: ジョブは開始可能状態で、対象フェーズの active run が存在しない。前段フェーズが必要な場合は完了済みである。
- `start condition`: job、前段 phase run、runtime 設定、対象件数を読み込める。
- `transition reason`: フェーズ開始要求である。対象件数が 0 件の場合は即時完了の遷移理由も持つ。
- `expected outcome`: `JOB_PHASE_RUN` が `running` または `completed` で作られる。対象件数が 0 件なら phase run と job が完了側へ進む可能性がある。
- `observable point`: `event=phase_start_transition`, `where=backend.service.<phase>`, `result=transitioned`, `id=job:<id>`, `phase`, `before_state`, `after_state`, `reason=start_requested|zero_target_completed`, `count` を残す候補にする。
- `observed disappearing information`: phase run 作成前の「active run が無い」状態、0 件完了にした理由、開始時点の対象件数。
- `forbidden log`: API key、credential_ref の実値、prompt 全文、翻訳本文全文、XML 全文、provider raw payload は出さない。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: フェーズ開始は job と phase run の二重状態を同時に変えるため、原因分離価値が高い。
- `conflict hint`: phase ごとに `phase` の固定値を統一する必要がある。候補は命名を確定しない。

### CAND-OBS-ST-004 フェーズ開始拒否を観測する

- `source requirement`: 観測ログ仕様 3 は `reason` を拒否、破棄、失敗分類に使える任意 payload とする。用語翻訳、NPC ペルソナ生成、本文翻訳の開始処理は、前提状態や active run の存在で拒否する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-004`
- `actor`: backend service
- `trigger`: 開始条件を満たさない状態でフェーズ開始が要求される。
- `previous state`: ジョブが `ready` 以外、前段 phase run が未完了、対象 phase run が active、または runtime 設定が不足している。
- `start condition`: 拒否理由を判定するための job state、phase run state、runtime 設定状態を読める。
- `transition reason`: 開始要求は受けたが、開始前提を満たさない。
- `expected outcome`: 永続化状態を変更しない。結果は `rejected` になり、理由は `ready_required`、`active_phase_exists`、`term_phase_incomplete`、`persona_phase_incomplete`、`terminal_job`、`runtime_snapshot_missing` などに分類される。
- `observable point`: `event=phase_start_rejected`, `where=backend.service.<phase>`, `result=rejected`, `id=job:<id>`, `phase`, `before_state`, `after_state=unchanged`, `reason` を残す候補にする。
- `observed disappearing information`: response を返した後に消える拒否判定の局所理由、前段 phase run の状態、active run の存在。
- `forbidden log`: secret、API key、credential_ref の実値、source_utterance_full_text、provider error raw body は出さない。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: 拒否理由は画面操作後に UI 文言へ畳まれるため、backend の分類ログとして残す価値が高い。
- `conflict hint`: public error kind と内部 reason のどちらをログに残すかは、セキュリティ観点の候補と統合が必要である。

### CAND-OBS-ST-005 フェーズ一時停止、再開、リトライ、キャンセルを観測する

- `source requirement`: 要件一覧 7 は `Running`、`Paused`、`RecoverableFailed`、`Canceled` の操作系と異常系遷移を定義する。architecture 正本は状態遷移規則を `StateMachine` または service 境界で扱う。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-005`
- `actor`: backend service
- `trigger`: 利用者が phase run に対して pause、resume、retry、cancel を要求する。
- `previous state`: phase run は `running`、`paused`、`recoverable_failed`、`completed`、`canceled` のいずれかである。
- `start condition`: 対象 phase run ID と job ID が一致し、現在 state を読める。
- `transition reason`: 利用者の phase 操作要求である。
- `expected outcome`: 許可状態では `running -> paused`、`paused -> running`、`recoverable_failed -> running`、cancel 可能状態から `canceled` へ遷移する。非許可状態では状態を変えず、`phase is not running`、`phase is not resumable`、`phase is not retryable`、`phase is not cancelable` などを返す。
- `observable point`: `event=phase_operation_transition`, `where=backend.service.<phase>`, `result=transitioned|rejected`, `id=phase_run:<id>`, `phase`, `operation`, `before_state`, `after_state`, `reason` を残す候補にする。
- `observed disappearing information`: 操作直前の phase run state、操作別の拒否理由、job state が終端かどうかの判定。
- `forbidden log`: phase 対象レコードごとの 1 件ログ、翻訳本文、prompt、latest_error 全文、provider raw payload は出さない。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `observability_requirement`, `冪等性_requirement`
- `adoption hint`: phase 操作は同じ ID に対する再送があり得るため、重複実行ではなく状態不変の拒否として観測できるとよい。
- `conflict hint`: cancel の job state 連動範囲は phase ごとに異なる可能性があるため、最終シナリオで分割が必要になる可能性がある。

### CAND-OBS-ST-006 フェーズ実行結果の完了と失敗を観測する

- `source requirement`: 要件一覧 4 は、翻訳ジョブの失敗回復と進捗確認を求める。観測ログ仕様 4 は、大量処理の件数、分類、最初の失敗、最後の失敗が集約済みの箇所を追加対象にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-006`
- `actor`: backend service
- `trigger`: provider 実行、辞書ヒット適用、ペルソナ生成、本文翻訳の処理結果が確定する。
- `previous state`: phase run は `running` である。
- `start condition`: 実行対象件数、成功件数、skip 件数、失敗分類、保存結果を集約できる。
- `transition reason`: 実行結果の集約により、phase run の完了または失敗が決まる。
- `expected outcome`: 全対象が処理済みなら `completed` へ遷移する。回復可能な失敗なら `recoverable_failed` へ遷移する。回復不能な失敗なら `failed` 系の結果になる。
- `observable point`: `event=phase_execution_settled`, `where=backend.service.<phase>`, `result=completed|recoverable_failed|failed`, `id=phase_run:<id>`, `phase`, `before_state=running`, `after_state`, `reason`, `count` を残す候補にする。
- `observed disappearing information`: 集約後に畳まれる成功件数、skip 件数、最初の失敗分類、最後の失敗分類、完了判定理由。
- `forbidden log`: loop 内 1 件ごとのログ、provider raw payload、prompt 全文、翻訳本文全文、XML 全文、secret は出さない。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `observability_requirement`, `performance_requirement`
- `adoption hint`: 実行結果の確定点は後続調査で最も参照されるため、集約ログだけを候補にする。
- `conflict hint`: 失敗分類の粒度は failure 観点候補と重複する可能性がある。

### CAND-OBS-ST-007 設定検証の stale 遷移を観測する

- `source requirement`: 観測ログ仕様 4 は、frontend runtime event の破棄理由と拒否理由を残す対象にする。provider 設定検証と job setup 検証は `validation_stale` と `model_selection_stale` を返す。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-007`
- `actor`: backend service
- `trigger`: provider 設定検証または翻訳ジョブ作成前検証が request token や model list source token の不一致を検出する。
- `previous state`: 設定検証は `pending`、`not_validated`、`validated`、または画面側の `fresh|stale` 相当状態である。
- `start condition`: request token、model list source token、検証対象の現在 snapshot を比較できる。
- `transition reason`: 利用者が画面上で設定を変更した後、古い検証結果または古い model 選択が戻った。
- `expected outcome`: 現在状態を古い結果で上書きしない。結果は `rejected` または `not_validated` になり、`validation_stale` または `model_selection_stale` を返す。
- `observable point`: `event=settings_validation_stale`, `where=backend.service.provider_settings|backend.service.translation_job_setup`, `result=skipped|rejected`, `id=provider:<id>`, `before_state`, `after_state=unchanged`, `reason=validation_stale|model_selection_stale` を残す候補にする。
- `observed disappearing information`: 古い request token が破棄された理由、古い model list と現在選択の不一致、検証結果を上書きしなかった理由。
- `forbidden log`: API key、credential_ref の実値、endpoint URL の secret 部分、provider response body は出さない。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: stale 判定は UI 上では単なる再検証要求に見えるため、原因分離ログが有効である。
- `conflict hint`: provider ID を `id` として残す場合、secret ではないことを security 観点と合わせて確認する必要がある。

### CAND-OBS-ST-008 runtime event 購読と破棄を観測する

- `source requirement`: architecture 正本は `RuntimeEventAdapter` を Wails event 購読から screen local handler へ写像する境界とする。観測ログ仕様 4 は、frontend runtime event の破棄理由が画面操作後に消える箇所を追加対象にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-OBS-ST-008`
- `actor`: frontend runtime adapter
- `trigger`: 画面 mount、画面 dispose、Wails runtime event 到着、runtime 不在、payload parse 失敗が発生する。
- `previous state`: runtime event adapter は未購読、購読中、解除済みのいずれかである。
- `start condition`: runtime bridge の有無、detach callback の有無、payload parse 結果を判定できる。
- `transition reason`: 画面の mount / dispose、または runtime event 受信である。
- `expected outcome`: runtime がある場合は購読中へ遷移する。runtime がない場合は未購読のまま `runtime_unavailable` として扱う。dispose では解除済みへ遷移する。payload が object でない場合は event を空 payload へ畳むか破棄する。
- `observable point`: `event=runtime_event_adapter_transition`, `where=frontend.runtime.<screen>`, `result=subscribed|detached|skipped|dropped`, `reason=runtime_unavailable|payload_invalid|dispose|resubscribe`, `before_state`, `after_state` を残す候補にする。
- `observed disappearing information`: runtime 不在で購読しなかった事実、再購読時に前の購読を解除した事実、payload parse 失敗理由。
- `forbidden log`: event payload 全体、翻訳本文、XML 全文、provider raw payload、backend への frontend log 転送は出さない。
- `related detail requirement type`: `state_requirement`, `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: frontend runtime event は画面状態へ吸収されると破棄理由が消えるため、console 側の軽量ログ候補にする。
- `conflict hint`: frontend log は backend file へ集約しないという観測ログ仕様と衝突しないようにする。

## Open Notes

- `human decision candidate`: job state と phase run state の両方が変わる遷移で、ログの主語を job にするか phase run にするかは人間判断が必要である。
- `human decision candidate`: `credential_ref` は実値を出さない方針だが、redacted ID もログ禁止に含めるかどうかは人間判断が必要である。
- `human decision candidate`: runtime event の payload parse 失敗を `dropped` とするか、空 payload に畳んだ `skipped` とするかは人間判断が必要である。
- `merge candidate`: CAND-OBS-ST-003、CAND-OBS-ST-004、CAND-OBS-ST-005、CAND-OBS-ST-006 は phase 共通の状態遷移シナリオとして統合できる可能性がある。
- `merge candidate`: CAND-OBS-ST-007 は external-integration 観点の provider 境界候補と統合される可能性がある。
- `rejection candidate`: 全 command の start / finish log を一律に追加する候補は、観測ログ仕様の禁止事項に反するため候補から除外する。
- `rejection candidate`: loop 内 1 件ごとの phase target log は、観測ログ仕様の禁止事項に反するため候補から除外する。
- `rejection candidate`: trace ID 導入、frontend log の backend 集約、constructor 引数拡張、context logger 埋め込みは、観測ログ仕様の禁止事項に反するため候補から除外する。
