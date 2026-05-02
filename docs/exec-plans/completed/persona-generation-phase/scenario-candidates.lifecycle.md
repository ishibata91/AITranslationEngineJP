# Scenario Candidates: persona-generation-phase / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP`

## Generator Scope

- `viewpoint`: 作成、生成、保存、参照、再開、完了、flush 前提の時間順 lifecycle。
- `included_sources`: `./plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `excluded_sources`: 他 generator の候補成果物、未承認 design bundle、implementation-scope、product code、product test、docs 正本更新。
- `generation_notes`: 採否、統合、最終シナリオ表の確定は扱わない。各候補は designer が判断できる lifecycle 候補として残す。

## Candidate Scenarios

### CAND-PGP-001 単語翻訳フェーズ完了後に NPC ペルソナ生成フェーズを開始する

- `source requirement`: `docs/exec-plans/active/persona-generation-phase/plan.md:60-64`, `tasks/usecases/persona-generation-phase.yaml:19-27`, `tasks/index.yaml:4-10`, `docs/spec.md:100-115`, `docs/spec.md:126-133`
- `viewpoint`: lifecycle / 開始。
- `candidate scenario id`: `CAND-PGP-001`
- `actor`: ユーザー、phase runner。
- `trigger`: Job Run で、単語翻訳フェーズ完了済みの翻訳ジョブから NPC ペルソナ生成フェーズ開始を実行する。
- `expected outcome`: NPC ペルソナ生成フェーズが current phase になり、progress と phase run 開始結果を確認できる。前段未完了の job では開始不可理由を確認できる。
- `observable point`: Job Run UI、phase start result、`JOB_PHASE_RUN`。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `data_requirement`
- `acceptance relevance`: phase 順序、current phase、progress の受け入れ条件に関係する。
- `adoption hint`: 正常 lifecycle の入口候補として検討できる。
- `conflict hint`: state-transition 観点で Ready / Running / phase state の扱いと競合しうる。開始可能な job state は designer 側で統合する。

### CAND-PGP-002 原文発話、属性、会話文脈、共通ペルソナから生成対象を確定する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9-18`, `tasks/usecases/persona-generation-phase.yaml:21-24`, `docs/spec.md:21-25`, `docs/spec.md:128-131`, `docs/spec.md:208-216`, `docs/er.md:32-33`, `docs/er.md:45-55`
- `viewpoint`: lifecycle / 生成対象確定。
- `candidate scenario id`: `CAND-PGP-002`
- `actor`: phase runner。
- `trigger`: NPC ペルソナ生成フェーズ開始後、対象 job の NPC 発話、NPC 属性メタデータ、会話文脈、共通ペルソナを読み取る。
- `expected outcome`: NPC ごとの生成対象、入力要約、共通ペルソナ参照有無、生成対象外理由を確認できる。
- `observable point`: phase target summary、`NPC_PROFILE`、`NPC_RECORD`、`TRANSLATION_FIELD_RECORD_REFERENCE`。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `testability_requirement`
- `acceptance relevance`: NPC ごとの生成対象確認と、本文翻訳に渡す persona snapshot の前提に関係する。
- `adoption hint`: 生成 lifecycle の前処理候補として検討できる。
- `conflict hint`: 共通ペルソナが存在する NPC を job 内生成対象から外すか、参照 snapshot に含めるかは未決になりうる。

### CAND-PGP-003 NPC ごとのジョブ内ペルソナを生成する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9-18`, `docs/spec.md:21-25`, `docs/spec.md:130-131`, `docs/spec.md:224-228`, `docs/er.md:48-56`
- `viewpoint`: lifecycle / 生成。
- `candidate scenario id`: `CAND-PGP-003`
- `actor`: phase runner、AI provider adapter。
- `trigger`: 生成対象 NPC ごとに、原文発話、NPC 属性メタデータ、会話文脈、共通ペルソナ参照を使って生成処理を実行する。
- `expected outcome`: NPC ごとのペルソナ生成結果を取得でき、生成根拠の翻訳フィールドと NPC profile を失わない。
- `observable point`: provider request summary、provider response contract、generated persona result、`PERSONA_FIELD_EVIDENCE`。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `acceptance relevance`: NPC ペルソナ生成フェーズの中心 outcome と、本文翻訳フェーズ入力の品質に関係する。
- `adoption hint`: 生成処理本体の正常 lifecycle 候補として検討できる。
- `conflict hint`: external-integration 観点で provider 選択、paid real API 不使用、fake provider 方針と統合が必要になる。

### CAND-PGP-004 生成済みペルソナをジョブ内ペルソナとして保存する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:16-25`, `docs/spec.md:34-35`, `docs/spec.md:215-216`, `docs/er.md:48-59`, `docs/er.md:61-69`
- `viewpoint`: lifecycle / 保存。
- `candidate scenario id`: `CAND-PGP-004`
- `actor`: phase runner、repository。
- `trigger`: NPC ごとの valid persona result を受け取った後、対象 job のジョブ内ペルソナとして保存する。
- `expected outcome`: `PERSONA.translation_job_id` が対象 job に紐づき、共通ペルソナとは分離される。phase run と対象 persona の関連、生成根拠、persona snapshot が追跡できる。
- `observable point`: `PERSONA`, `PERSONA_FIELD_EVIDENCE`, `PHASE_RUN_PERSONA`, phase result。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `state_requirement`
- `acceptance relevance`: ジョブ内ペルソナ保持、共通ペルソナ分離、本文翻訳入力参照の受け入れ条件に関係する。
- `adoption hint`: 保存 lifecycle と永続化境界の候補として検討できる。
- `conflict hint`: 同一 `NPC_PROFILE` に共通 persona と job 内 persona を同時保持しない ER 方針と、spec の共通 / ジョブ内分離要件の統合判断が必要になる可能性がある。

### CAND-PGP-005 保存済み persona snapshot を本文翻訳フェーズ入力として参照する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:16-25`, `tasks/usecases/body-translation-phase.yaml:9-22`, `docs/spec.md:21-24`, `docs/spec.md:33-35`, `docs/spec.md:113-115`, `docs/spec.md:246-248`
- `viewpoint`: lifecycle / 参照。
- `candidate scenario id`: `CAND-PGP-005`
- `actor`: ユーザー、body translation phase runner。
- `trigger`: NPC ペルソナ生成フェーズ完了後、本文翻訳フェーズの開始または入力 summary 確認を行う。
- `expected outcome`: 本文翻訳フェーズが、対象 job のジョブ内ペルソナと persona snapshot を入力として参照できる。参照不能な場合は本文翻訳フェーズ開始不可理由を確認できる。
- `observable point`: body phase input summary、phase transition result、Job Run UI、`PHASE_RUN_PERSONA`。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `acceptance relevance`: 本文翻訳フェーズの precondition と、ペルソナを翻訳補助メタデータとして提供する要件に関係する。
- `adoption hint`: フェーズ間参照 lifecycle の候補として検討できる。
- `conflict hint`: body-translation-phase 側の候補と、参照単位を persona record、snapshot、summary のどれで固定するかが競合しうる。

### CAND-PGP-006 Job Run を開き直して phase result と persona 参照状態を確認する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:21-24`, `tasks/usecases/persona-generation-phase.yaml:30-34`, `docs/spec.md:35`, `docs/spec.md:53-55`, `docs/architecture.md:92-98`, `docs/architecture.md:173-180`
- `viewpoint`: lifecycle / 参照状態確認。
- `candidate scenario id`: `CAND-PGP-006`
- `actor`: ユーザー。
- `trigger`: NPC ペルソナ生成フェーズの実行中、完了後、または画面再表示後に Job Run を開く。
- `expected outcome`: current phase、progress、phase result、NPC ごとの生成結果、ジョブ内ペルソナ参照状態を UI で確認できる。
- `observable point`: Job Run UI、gateway query response、runtime event 後の reload result。
- `related detail requirement type`: `success_requirement`, `observability_requirement`, `testability_requirement`
- `acceptance relevance`: UI から実行前後の翻訳補助メタデータを観測可能にする要件に関係する。
- `adoption hint`: lifecycle 終了後利用と画面再表示の候補として検討できる。
- `conflict hint`: operation-audit 観点の監査保存対象と重なりうる。UI 表示項目と監査ログ項目は designer 側で分ける。

### CAND-PGP-007 中断、再開、リトライで同じ phase run を継続する

- `source requirement`: `docs/spec.md:53-55`, `docs/spec.md:133-199`, `docs/er.md:61-69`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:320-344`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:389-399`
- `viewpoint`: lifecycle / 再開。
- `candidate scenario id`: `CAND-PGP-007`
- `actor`: ユーザー、phase runner。
- `trigger`: NPC ペルソナ生成フェーズが Paused または RecoverableFailed の状態で、再開またはリトライを実行する。
- `expected outcome`: 同じ `JOB_PHASE_RUN` が継続され、保存済み persona と phase link は重複作成されない。未処理 NPC だけが生成対象に戻り、progress と最新 error が更新される。
- `observable point`: phase run ID、`PERSONA` 件数、`PHASE_RUN_PERSONA` 件数、progress、latest error。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `consistency_requirement`
- `acceptance relevance`: 翻訳ジョブの中断、再開、失敗回復の要件に関係する。
- `adoption hint`: 回復 lifecycle の候補として検討できる。
- `conflict hint`: failure 観点の recoverable failure と、state-transition 観点の禁止遷移に統合が必要になる。

### CAND-PGP-008 全対象の保存後に NPC ペルソナ生成フェーズを完了する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:21-25`, `tasks/usecases/body-translation-phase.yaml:21-30`, `docs/spec.md:100-115`, `docs/spec.md:187-199`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:86-91`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:295-318`
- `viewpoint`: lifecycle / 完了。
- `candidate scenario id`: `CAND-PGP-008`
- `actor`: phase runner。
- `trigger`: 生成対象 NPC の処理とジョブ内ペルソナ保存がすべて成功する。
- `expected outcome`: NPC ペルソナ生成フェーズの phase result が Completed になり、本文翻訳フェーズの開始条件が満たされる。未完了、失敗、persona 参照不能では後続 phase run を作成しない。
- `observable point`: phase result、progress、body phase start enablement、phase transition result。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `acceptance relevance`: completion criteria と後続フェーズ precondition に関係する。
- `adoption hint`: 正常 lifecycle の完了候補として検討できる。
- `conflict hint`: 一部 NPC だけ失敗した場合に phase 全体を failed、recoverable、partial completed のどれで扱うかは人間判断候補になりうる。

### CAND-PGP-009 翻訳完了後にジョブ内ペルソナを flush 可能にする

- `source requirement`: `docs/spec.md:34`, `docs/spec.md:59`, `docs/spec.md:215-216`, `docs/er.md:48-59`, `tasks/usecases/body-translation-phase.yaml:17-30`
- `viewpoint`: lifecycle / flush 前提。
- `candidate scenario id`: `CAND-PGP-009`
- `actor`: job lifecycle manager。
- `trigger`: 本文翻訳フェーズと後続の出力利用が終わり、未完了ジョブが対象ジョブ内ペルソナを参照していないことを確認する。
- `expected outcome`: 共通ペルソナは残し、対象 job のジョブ内ペルソナだけを flush 対象として識別できる。flush 後も出力済み成果物と共通ペルソナ再利用の前提を壊さない。
- `observable point`: flush target summary、`PERSONA.persona_scope`, `PERSONA.translation_job_id`, dependent phase references。
- `related detail requirement type`: `data_requirement`, `recovery_requirement`, `compatibility_requirement`
- `acceptance relevance`: ジョブ内生成物の削除、共通ペルソナ再利用、入力キャッシュ再構築可能性に関係する。
- `adoption hint`: lifecycle 終了後の保存期間と flush 前提候補として検討できる。
- `conflict hint`: flush の実行タイミングは、本文翻訳完了後、成果物出力後、ジョブ Completed 後のどれかで競合しうる。designer 側で人間判断に回す可能性がある。

## Open Notes

- `human decision candidate`: 共通 persona が存在する NPC の job 内 persona 生成有無。partial success の phase state。persona snapshot の保存単位。flush 実行タイミング。
- `merge candidate`: `CAND-PGP-001` と state-transition の開始候補。`CAND-PGP-007` と failure / state-transition の回復候補。`CAND-PGP-006` と operation-audit の観測候補。
- `rejection candidate`: lifecycle 観点では確定しない。designer が他観点と統合して判断する。
