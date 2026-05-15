# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / state-transition

- `generator`: `state-transition`
- `source_plan`: `./implement-lane-task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSR-ST`

## Generator Scope

- `viewpoint`: 状態遷移。`TRANSLATION_JOB.state`、`JOB_PHASE_RUN.state`、内部一時 state、read model state の混在を検出し、許可遷移、禁止遷移、冪等再実行の候補を出す。
- `included_sources`: `./implement-lane-task-frame.md`, `./state-knowledge-investigation.md`, `./state-knowledge-investigation-lane-decision.md`, `../../../spec.md`, `../../../architecture.md`, `../../../detail-specs/translation-job-management.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`, `../../../../internal/usecase/translationjobpolicy/policy.go`, `../../../../internal/service/phase_action_enablement_helpers.go`
- `excluded_sources`: product code 変更、product test 変更、docs 正本本文の変更、`docs/exec-plans/completed/**`、UI 変更、`stale_selection`、`validation_stale`、`model_selection_stale` の削除判断。
- `generation_notes`: 正本 state は `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` に分ける。`pending` は正本 state として扱わない。`idle_ready`、`blocked`、`not-ready`、`ready`、`starting`、`validation failed`、`empty completed`、`not_started`、`rejected`、`snapshot_missing` は read model state または表示差分として扱い、永続 state 候補へ混ぜない。

## Candidate Scenarios

### CAND-TJSR-ST-001 Ready job は phase start 許可時だけ Running へ進む

- `source requirement`: `docs/spec.md:141-145`, `docs/spec.md:181-182`, `docs/detail-specs/translation-job-management.md:28-32`, `docs/detail-specs/term-translation-phase.md:28-35`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-001`
- `actor`: 翻訳 job 実行 usecase
- `trigger`: 対象 job が `Ready` で、active な `JOB_PHASE_RUN` がなく、単語翻訳フェーズの開始前提が成立した状態で phase start を要求する。
- `expected outcome`: `TRANSLATION_JOB.state` は `Ready` から `Running` へ遷移する。対象フェーズの `JOB_PHASE_RUN.state` は作成時点から `Running` になる。`pending` の永続 state は作らない。
- `observable point`: job snapshot は `Running` を持つ。phase run snapshot は `Running` を持つ。Ready job 表示だけでは `TRANSLATION_JOB.state` も `JOB_PHASE_RUN` も変わらない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: 正本遷移の start 成功候補として扱える。
- `conflict hint`: `pending` を canonical state へ昇格する判断と競合する。

### CAND-TJSR-ST-002 pending は正本 state ではなく危険操作を許可しない

- `source requirement`: `docs/spec.md:170-191`, `state-knowledge-investigation.md:24-27`, `state-knowledge-investigation.md:35-40`, `docs/detail-specs/translation-job-management.md:45-47`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-002`
- `actor`: 翻訳 job 実行 usecase
- `trigger`: 既存データまたは内部処理の途中値として `JOB_PHASE_RUN.state = pending` が観測される。
- `expected outcome`: `pending` は `JOB_PHASE_RUN.state` の正本遷移として扱わない。pause、resume、retry、cancel は許可しない。表示や集約だけで `pending` を `Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` へ書き換えない。
- `observable point`: 操作可否は危険操作無効になる。無効理由は状態不整合または phase state 不一致として観測できる。正本 state、内部一時 state、read model state が同じ項目へ混在しない。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: `pending` を内部一時 state として隔離する場合の候補として扱える。
- `conflict hint`: `pending` を canonical state に追加する場合、候補の期待結果を変更する必要がある。

### CAND-TJSR-ST-003 read model state は永続 state を更新しない

- `source requirement`: `docs/spec.md:137-140`, `docs/detail-specs/term-translation-phase.md:83-85`, `docs/detail-specs/body-translation-phase.md:78-84`, `state-knowledge-investigation.md:26`, `docs/detail-specs/translation-job-management.md:45-50`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-003`
- `actor`: Job Run 読み取り usecase
- `trigger`: phase summary が `idle_ready`、`blocked`、`not-ready`、`ready`、`starting`、`validation failed`、`empty completed` などの表示用状態差分を導出する。
- `expected outcome`: read model state は表示差分として返す。`TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は更新しない。
- `observable point`: 一覧表示、Job Run 表示、実行不可理由表示の読み取りだけでは DB の job state と phase state が変わらない。集約不能は成功状態や空状態と区別できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 正本 state と表示用 state の分離を確認する候補として扱える。
- `conflict hint`: UI 変更を前提にした候補ではない。表示文言の確定は別成果物に残す。

### CAND-TJSR-ST-004 common operation は TranslationJobPolicy と同じ状態事実で決まる

- `source requirement`: `docs/spec.md:196-203`, `docs/architecture.md:169-180`, `internal/usecase/translationjobpolicy/policy.go:67-123`, `internal/service/phase_action_enablement_helpers.go:25-61`, `state-knowledge-investigation.md:27-28`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-004`
- `actor`: 翻訳 job 操作 usecase
- `trigger`: `pause`、`resume`、`retry`、`cancel` の操作可否を評価する。
- `expected outcome`: `Running` の phase run だけ pause を許可する。`Paused` の phase run だけ resume と cancel を許可する。`RecoverableFailed` の phase run だけ retry を許可する。terminal job ではすべて拒否する。
- `observable point`: 状態ごとの操作可否と拒否理由は、phase type に依存しない。read model の `CanPause`、`CanResume`、`CanRetry`、`CanCancel` は共通操作規則と矛盾しない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: `commonPhaseActionAvailability` の重複知識を検出する候補として扱える。
- `conflict hint`: `TranslationJobPolicy` は UseCase だけが呼ぶという architecture 制約と、read model 操作可否で同じ規則を再利用する要求が競合しうる。

### CAND-TJSR-ST-005 Paused cancel は canceled spelling の terminal state にそろう

- `source requirement`: `docs/spec.md:161-165`, `docs/spec.md:188-191`, `docs/detail-specs/body-translation-phase.md:43-44`, `state-knowledge-investigation.md:33`, `state-knowledge-investigation-lane-decision.md:18`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-005`
- `actor`: 翻訳 job 操作 usecase
- `trigger`: `TRANSLATION_JOB.state = Paused` かつ対象 `JOB_PHASE_RUN.state = Paused` の状態で cancel を要求する。
- `expected outcome`: `TRANSLATION_JOB.state` と対象 `JOB_PHASE_RUN.state` は `Canceled` へ遷移する。永続値と API 由来の state spelling は `canceled` にそろえる。`cancelled` は正本 spelling として扱わない。
- `observable point`: cancel 後の job は terminal job として扱われる。phase result の途中成功結果は output readiness に使われない。`cancelled` が観測された場合は正本 state として一致しない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: cancel spelling の候補として扱える。
- `conflict hint`: `cancelled` fixture spelling を今回の stale 廃止に含めるかは人間判断候補である。

### CAND-TJSR-ST-006 terminal job は phase 作成、保存、readiness 更新、late response 後書きを拒否する

- `source requirement`: `docs/spec.md:196-203`, `docs/detail-specs/term-translation-phase.md:59-61`, `docs/detail-specs/persona-generation-phase.md:43-48`, `docs/detail-specs/body-translation-phase.md:43-46`, `internal/usecase/translationjobpolicy/policy.go:116-123`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-006`
- `actor`: phase 実行 usecase
- `trigger`: `TRANSLATION_JOB.state` が `Completed`、`Failed`、`Canceled` のいずれかで、phase start、結果保存、readiness 更新、late response 後書きが発生する。
- `expected outcome`: 状態変更を拒否する。既存の `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は変わらない。
- `observable point`: phase run は新規作成されない。成功結果、persona snapshot、body readiness、field result は後書きされない。拒否理由は terminal job として観測できる。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: terminal guard の禁止遷移候補として扱える。
- `conflict hint`: observability 側の late response 記録対象と統合時に観測対象が競合しうる。

### CAND-TJSR-ST-007 RecoverableFailed retry は同じ JOB_PHASE_RUN を継続する

- `source requirement`: `docs/spec.md:157-158`, `docs/spec.md:173-185`, `docs/detail-specs/term-translation-phase.md:50-56`, `docs/detail-specs/persona-generation-phase.md:40-45`, `docs/detail-specs/body-translation-phase.md:40-42`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-007`
- `actor`: 翻訳 job 操作 usecase
- `trigger`: 対象 `JOB_PHASE_RUN.state = RecoverableFailed` で retry を要求する。
- `expected outcome`: `JOB_PHASE_RUN.state` は `RecoverableFailed` から `Running` へ遷移する。`TRANSLATION_JOB.state` は retry 実行中の状態と矛盾しない。新しい `JOB_PHASE_RUN` は作らない。
- `observable point`: phase run id は retry 前後で同一である。成功済み辞書 entry、persona、field result は重複作成されない。未処理対象だけが retry 対象になる。
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: 冪等再実行と状態遷移を結びつける候補として扱える。
- `conflict hint`: retryable failure の判定材料は phase 別詳細仕様と統合する必要がある。

### CAND-TJSR-ST-008 body phase Completed 後だけ job Completed へ進む

- `source requirement`: `docs/spec.md:159-160`, `docs/detail-specs/body-translation-phase.md:20-23`, `docs/detail-specs/body-translation-phase.md:40-46`, `docs/detail-specs/translation-job-management.md:26-47`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-TJSR-ST-008`
- `actor`: 本文翻訳 phase 実行 usecase
- `trigger`: 本文翻訳 phase の保存、保護要素検証、field result 整合、output readiness 判定が完了する。
- `expected outcome`: body phase が `Completed` で、field result 整合と output readiness が成立した場合だけ、job は `Running` から `Completed` へ遷移する。保存失敗、検証失敗、partial state では `Completed` にしない。
- `observable point`: Completed job は未完了一覧から外れる。失敗または検証失敗では failure state と retryable flag を確認できる。表示だけで job state は変わらない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: job-level 完了遷移の候補として扱える。
- `conflict hint`: output readiness の詳細条件は body phase 側の候補と統合する必要がある。

## Open Notes

- `human decision candidate`: `pending` を canonical state へ昇格するか、内部一時 state として隔離するかは未確定である。
- `human decision candidate`: `TranslationJobPolicy` の共通操作規則を read model の操作可否へどう再利用するかは未確定である。
- `human decision candidate`: `cancelled` fixture spelling を今回の stale 廃止に含めるかは未確定である。
- `merge candidate`: `CAND-TJSR-ST-001`、`CAND-TJSR-ST-002`、`CAND-TJSR-ST-003` は、正本 state、内部一時 state、read model state の分離として統合できる可能性がある。
- `merge candidate`: `CAND-TJSR-ST-004`、`CAND-TJSR-ST-005`、`CAND-TJSR-ST-006` は、共通操作規則と terminal guard の候補として統合できる可能性がある。
- `merge candidate`: `CAND-TJSR-ST-007`、`CAND-TJSR-ST-008` は、phase 完了、失敗、再実行の一貫性候補として統合できる可能性がある。
- `rejection candidate`: UI 変更を前提にする候補は今回の state-transition 候補から除外する。
- `rejection candidate`: `stale_selection`、`validation_stale`、`model_selection_stale` を削除する候補は今回の state-transition 候補から除外する。

## Completion Material

- `candidate_count`: 8
- `candidate_artifact_path`: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-candidates.state-transition.md`
- `viewpoint`: `state-transition`
- `remaining_risk`: `pending` の扱い、read model 操作可否への共通操作規則再利用、`cancelled` fixture spelling の今回範囲は人間判断が残る。
