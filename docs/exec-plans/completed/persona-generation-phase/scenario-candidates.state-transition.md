# Scenario Candidates: persona-generation-phase / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP-ST`

## Generator Scope

- `viewpoint`: state-transition
- `included_sources`: `plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `excluded_sources`: 他観点候補、最終シナリオ表、採否決定、実装範囲、プロダクトコード、プロダクトテスト、docs 正本化
- `generation_notes`: 単語翻訳フェーズ完了後から、NPC ペルソナ生成フェーズを経て、本文翻訳フェーズの readiness が成立するまでの状態遷移候補だけを扱う。

## Candidate Scenarios

### CAND-PGP-ST-001 term completed から persona phase を開始する

- `source requirement`: `plan.md:60-64`, `tasks/index.yaml:4-10`, `tasks/usecases/persona-generation-phase.yaml:19-27`, `docs/spec.md:100-114`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:295-318`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-001`
- `actor`: Job Run 利用者
- `pre-state`: 単語翻訳フェーズが Completed で、対象 job が terminal ではなく、同一 job に active phase run がない。
- `trigger`: Job Run から NPC ペルソナ生成フェーズ開始を実行する。
- `start condition`: 確定訳語とジョブ内辞書を後続入力として参照できる。
- `post-state`: NPC ペルソナ生成フェーズの `JOB_PHASE_RUN` が開始済みになり、current phase と progress が persona phase を示す。
- `expected outcome`: persona phase 開始結果、current phase、progress、開始不可理由を確認できる。
- `observable point`: Job Run UI、phase transition result、`JOB_PHASE_RUN`
- `acceptance relevance`: term phase 完了後だけ persona phase が開始できることを固定する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: term phase の禁止遷移シナリオと対になる正常遷移候補。
- `conflict hint`: lifecycle 観点の phase start 候補と統合対象になりうる。

### CAND-PGP-ST-002 term 未完了または active phase ありでは persona phase を開始しない

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:19-27`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:86-91`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:295-318`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-002`
- `actor`: Job Run 利用者
- `pre-state`: 単語翻訳フェーズが未開始、Running、Paused、RecoverableFailed、または辞書参照不能である。
- `trigger`: NPC ペルソナ生成フェーズ開始を試行する。
- `start condition`: 開始条件を満たさないため、遷移を拒否する。
- `post-state`: persona phase の `JOB_PHASE_RUN` は作成されず、既存 phase state は不変である。
- `expected outcome`: 拒否理由と後続 phase 開始不可状態を確認できる。
- `observable point`: phase transition result、phase run 件数、拒否理由
- `acceptance relevance`: term phase 未完了から persona phase へ進む禁止遷移を固定する。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: 禁止遷移として designer が term phase の完了条件と接続できる。
- `conflict hint`: failure 観点の recoverable failure 候補と拒否理由の粒度が重複しうる。

### CAND-PGP-ST-003 persona phase 開始時に生成対象 snapshot を固定する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9-18`, `tasks/usecases/persona-generation-phase.yaml:21-25`, `docs/spec.md:21-24`, `docs/spec.md:209-216`, `docs/er.md:32-46`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-003`
- `actor`: phase execution boundary
- `pre-state`: persona phase が開始済みで、NPC 発話原文、NPC 属性メタデータ、会話文脈、共通ペルソナを読める。
- `trigger`: persona phase の生成対象抽出を実行する。
- `start condition`: NPC プロファイルと翻訳フィールド参照から、生成対象 NPC を識別できる。
- `post-state`: phase run に紐づく生成対象 snapshot が固定され、対象 NPC ごとの処理状態を進捗へ反映できる。
- `expected outcome`: NPC ごとの生成対象数、対象外理由、progress が観測できる。
- `observable point`: Job Run UI、phase result、`PHASE_RUN_TRANSLATION_FIELD`
- `acceptance relevance`: 「NPC ごとの生成対象を確認できる」を状態遷移として検証可能にする。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: actor-goal 観点の「生成対象確認」と統合しやすい候補。
- `conflict hint`: snapshot の保存単位は詳細設計で未固定の可能性がある。

### CAND-PGP-ST-004 共通ペルソナ参照時のジョブ内 persona 遷移を決める

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:10-18`, `tasks/usecases/persona-generation-phase.yaml:21-25`, `docs/spec.md:23-25`, `docs/spec.md:34-35`, `docs/er.md:48-59`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-004`
- `actor`: phase execution boundary
- `pre-state`: 対象 NPC に共通ペルソナが存在し、persona phase の入力に共通ペルソナが含まれる。
- `trigger`: 対象 NPC の persona 生成または再利用判定を実行する。
- `start condition`: 共通ペルソナとジョブ内ペルソナの同時保持可否が判定できる。
- `post-state`: 共通ペルソナを参照した persona snapshot が body phase 入力へ渡せる状態になる。
- `expected outcome`: 共通ペルソナを再利用した対象と、新規生成した対象を phase result で区別できる。
- `observable point`: `PERSONA`、persona snapshot、phase result
- `acceptance relevance`: 共通ペルソナを入力に含めつつ、ジョブ内ペルソナを共通ペルソナと分離して扱う判断点を露出する。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: designer が共通ペルソナ再利用ルールを人間判断候補へ送る材料になる。
- `conflict hint`: task はジョブ内ペルソナ出力を求める一方、ER は同一 NPC プロファイルで共通とジョブ内を同時保持しないとしている。

### CAND-PGP-ST-005 persona 生成成功で job-scoped persona と phase link を確定する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:16-25`, `docs/spec.md:22-24`, `docs/spec.md:34-35`, `docs/spec.md:214-216`, `docs/er.md:50-56`, `docs/er.md:66-69`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-005`
- `actor`: phase execution boundary
- `pre-state`: persona phase が Running で、生成対象 NPC の provider 結果が valid である。
- `trigger`: 生成済み persona を保存し、phase run と紐づける。
- `start condition`: `translation_job_id` を持つジョブ内データとして保存できる。
- `post-state`: job-scoped `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA` が整合し、persona snapshot が参照可能になる。
- `expected outcome`: 生成済み persona 件数、phase link、本文翻訳入力 summary を確認できる。
- `observable point`: `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA`、phase result
- `acceptance relevance`: 生成済みペルソナを body phase 入力として参照できることを状態で証明する。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `state_requirement`
- `adoption hint`: body readiness 候補の前段になる主要成功遷移。
- `conflict hint`: external-integration 観点の provider valid response 候補と統合されうる。

### CAND-PGP-ST-006 partial persona 保存では Completed と body readiness を成立させない

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:21-25`, `tasks/usecases/body-translation-phase.yaml:10-22`, `docs/er.md:63-69`, `docs/architecture.md:117-120`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:214-238`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-006`
- `actor`: phase execution boundary
- `pre-state`: persona 保存中に `PERSONA`、evidence、phase link、snapshot の一部だけが成功した。
- `trigger`: persona phase を完了扱いにしようとする。
- `start condition`: 整合した phase output が揃っていない。
- `post-state`: persona phase は Completed にならず、body phase readiness は false のままである。
- `expected outcome`: partial state は成功扱いにならず、再試行可能性と不整合理由を確認できる。
- `observable point`: phase state、row count、phase result、拒否理由
- `acceptance relevance`: 不完全な persona を body phase 入力として渡さない禁止遷移を固定する。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `adoption hint`: failure 観点の保存失敗候補と接続する state invariant 候補。
- `conflict hint`: 失敗時に rollback するか recoverable partial として残すかは designer の統合時に確認が必要である。

### CAND-PGP-ST-007 persona phase Completed で body phase readiness が成立する

- `source requirement`: `tasks/index.yaml:4-10`, `tasks/usecases/persona-generation-phase.yaml:16-27`, `tasks/usecases/body-translation-phase.yaml:9-22`, `docs/spec.md:100-114`, `docs/spec.md:225-248`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-007`
- `actor`: Job Run 利用者
- `pre-state`: persona phase が Completed で、ジョブ内ペルソナまたは persona snapshot が body phase 入力として参照可能である。
- `trigger`: Job Run で次 phase 操作または body phase input summary を確認する。
- `start condition`: 本文翻訳フェーズの precondition が満たされている。
- `post-state`: body phase readiness が true になり、本文翻訳フェーズの開始可否を判断できる。
- `expected outcome`: 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータの input summary が確認できる。
- `observable point`: Job Run UI、phase result、button enablement、input summary
- `acceptance relevance`: persona phase 完了が body phase readiness の必要条件であることを固定する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: body phase へ渡す readiness 境界の主候補。
- `conflict hint`: body phase 開始そのものを含めるか、readiness 表示までに留めるかは designer 統合で切り分ける。

### CAND-PGP-ST-008 persona 未完了または persona 参照不能では body phase readiness を成立させない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:10-22`, `docs/spec.md:130-131`, `docs/spec.md:225-248`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:295-318`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-008`
- `actor`: Job Run 利用者
- `pre-state`: persona phase が未開始、Running、Paused、RecoverableFailed、または persona snapshot 参照不能である。
- `trigger`: body phase readiness 確認、または本文翻訳フェーズ開始を試行する。
- `start condition`: persona phase completion requirement を満たさない。
- `post-state`: body phase readiness は false のままで、body phase の `JOB_PHASE_RUN` は作成されない。
- `expected outcome`: body phase 開始不可理由と persona phase の未完了状態を確認できる。
- `observable point`: phase transition result、phase run 件数、button enablement、拒否理由
- `acceptance relevance`: persona phase 未完了から body phase へ進む禁止遷移を固定する。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `adoption hint`: term phase の後続禁止遷移と同じ構造で採否判断しやすい。
- `conflict hint`: body phase 側の failure 観点候補と拒否理由の粒度が重複しうる。

### CAND-PGP-ST-009 persona phase の再送、再開、リトライで重複作成しない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:23-30`, `docs/spec.md:95-102`, `docs/spec.md:228-229`, `docs/er.md:61-69`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:320-343`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-009`
- `actor`: Job Run 利用者
- `pre-state`: 同一 job の persona phase run が active、paused、recoverable failed、または retryable である。
- `trigger`: persona phase の開始再送、再開、リトライを実行する。
- `start condition`: 同一 phase type の既存 `JOB_PHASE_RUN` が存在する。
- `post-state`: 同じ `JOB_PHASE_RUN` の状態が戻り、persona、phase link、snapshot が二重作成されない。
- `expected outcome`: phase run ID、生成済み persona 件数、未処理対象数、latest error、progress を確認できる。
- `observable point`: phase run ID、row count、progress、latest error
- `acceptance relevance`: 再送やリトライで body phase readiness が重複または破損しないことを固定する。
- `related detail requirement type`: `冪等性_requirement`, `concurrency_requirement`, `recovery_requirement`
- `adoption hint`: term phase の再実行方針を persona phase へ展開する候補。
- `conflict hint`: body phase usecase は pause/resume/retry/cancel を明記するが、persona phase の操作可否は未確定である。

### CAND-PGP-ST-010 terminal job には persona と body readiness を後書きしない

- `source requirement`: `docs/spec.md:95-102`, `docs/spec.md:228-229`, `docs/er.md:20-25`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:295-318`, `docs/architecture.md:117-120`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-PGP-ST-010`
- `actor`: Job Run 利用者
- `pre-state`: 対象 job が terminal state である。
- `trigger`: persona phase 開始、persona 保存、body phase readiness 更新のいずれかを試行する。
- `start condition`: terminal job への後書きは禁止される。
- `post-state`: `JOB_PHASE_RUN`、`PERSONA`、`PHASE_RUN_PERSONA`、body readiness は変更されない。
- `expected outcome`: terminal job の拒否理由と state 不変を確認できる。
- `observable point`: phase transition result、row count、拒否理由、state snapshot
- `acceptance relevance`: 完了済みまたは終了済み job の state invariant を保つ。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `authorization_requirement`
- `adoption hint`: 禁止遷移候補として phase 共通 state machine へ統合されうる。
- `conflict hint`: terminal state の列挙は現時点の入力だけでは確定できない。

## Open Notes

- `human decision candidate`: 共通ペルソナが存在する NPC で、ジョブ内 persona を作るのか、共通 persona を snapshot 参照するだけにするのかは未決である。
- `human decision candidate`: persona phase の pause、resume、retry、cancel の許可状態は、body phase では明記されているが persona phase では未確定である。
- `merge candidate`: `CAND-PGP-ST-001` と `CAND-PGP-ST-002` は phase start boundary の正常系 / 禁止遷移として統合されうる。
- `merge candidate`: `CAND-PGP-ST-007` と `CAND-PGP-ST-008` は body phase readiness boundary の正常系 / 禁止遷移として統合されうる。
- `rejection candidate`: body phase の実翻訳開始、訳文生成、保護要素検証は state-transition 候補の範囲外である。
