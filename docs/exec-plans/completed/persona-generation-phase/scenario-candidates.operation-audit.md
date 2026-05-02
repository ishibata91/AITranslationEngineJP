# Scenario Candidates: persona-generation-phase / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `excluded_sources`: product code, product tests, docs 正本更新, 他 generator 成果物, 最終 scenario 採否
- `generation_notes`: 生成対象、入力 snapshot、共通ペルソナ hit/miss、生成結果、本文翻訳参照、監査要約、redaction を operation-audit 観点だけで候補化する。

## Candidate Scenarios

### CAND-PGP-001 NPC ペルソナ生成対象を後追い確認できる

- `source requirement`: `./plan.md:13`, `./plan.md:61-63`, `tasks/usecases/persona-generation-phase.yaml:21-24`, `docs/er.md:22`, `docs/er.md:63`, `docs/er.md:69`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-001`
- `actor`: 利用者 / 運用確認者
- `trigger`: Job Run で NPC ペルソナ生成フェーズを開始する。
- `expected outcome`: phase run ごとに job ID、phase run ID、current phase、progress、生成対象 NPC 件数、対象 NPC の識別子を確認できる。
- `observable point`: Job Run summary、phase start result、`JOB_PHASE_RUN`、`PHASE_RUN_PERSONA`
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `testability_requirement`
- `acceptance relevance`: current phase と progress の確認、NPC ごとの生成対象確認、後続 designer の受け入れ条件固定に関係する。
- `adoption hint`: 生成対象一覧は監査対象である。本文や provider payload の保存要否はこの候補では確定しない。
- `conflict hint`: 生成対象の詳細を出しすぎると、原文発話や会話文脈の過剰保存と衝突する可能性がある。

### CAND-PGP-002 入力 snapshot と生成根拠を再現材料として確認できる

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:10-18`, `docs/spec.md:14-15`, `docs/spec.md:24`, `docs/spec.md:130-131`, `docs/spec.md:213`, `docs/er.md:20-33`, `docs/er.md:55`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-002`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: NPC ペルソナ生成フェーズが入力を固定して生成単位を作る。
- `expected outcome`: xEdit 抽出データ、NPC profile、NPC record、翻訳フィールド根拠、会話文脈、共通ペルソナ参照の snapshot ID または digest を確認できる。
- `observable point`: input snapshot summary、`NPC_PROFILE`、`NPC_RECORD`、`PERSONA_FIELD_EVIDENCE`
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `acceptance relevance`: persona snapshot の参照状態と、生成済みペルソナを本文翻訳フェーズの入力として追跡する条件に関係する。
- `adoption hint`: 再現材料は ID、version、digest、件数、根拠 field ID を中心にする候補として扱う。
- `conflict hint`: 原文発話全文や会話全文を snapshot として複製するかは未決である。redaction 候補と合わせて designer が判断する。

### CAND-PGP-003 共通ペルソナ hit/miss と重複回避理由を確認できる

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:15`, `tasks/usecases/persona-generation-phase.yaml:17-25`, `docs/spec.md:23-25`, `docs/spec.md:34`, `docs/spec.md:215-216`, `docs/er.md:50-55`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-003`
- `actor`: 利用者 / 運用確認者
- `trigger`: 生成対象 NPC ごとに共通ペルソナ参照を照合する。
- `expected outcome`: 対象 NPC ごとに common persona hit、common persona miss、生成 skipped、job persona generated の区別と理由を確認できる。
- `observable point`: hit/miss summary、persona selection summary、`PERSONA`、`PHASE_RUN_PERSONA`
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`
- `acceptance relevance`: ジョブ内ペルソナを共通ペルソナと分離して保持する条件、NPC ごとの生成対象確認に関係する。
- `adoption hint`: 共通ペルソナ再利用とジョブ内生成の差を、運用確認できる候補として残す。
- `conflict hint`: task はジョブ内ペルソナ生成を出力に持つが、ER は共通ペルソナ hit 時の同時保持を避ける前提を持つ。hit 時の最終扱いは designer の統合判断に残す。

### CAND-PGP-004 AI 生成結果と永続化結果を監査できる

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9`, `tasks/usecases/persona-generation-phase.yaml:17-24`, `docs/spec.md:22-23`, `docs/spec.md:55-58`, `docs/er.md:25`, `docs/er.md:50-55`, `docs/architecture.md:129`, `docs/architecture.md:138`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-004`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: AI provider 応答からペルソナ生成結果を作り、保存する。
- `expected outcome`: provider、model、execution mode、credential ref、prompt digest、input snapshot digest、生成件数、保存済み persona ID、失敗理由を確認できる。
- `observable point`: generation result summary、AIProvider adapter summary、`JOB_PHASE_RUN`、`PERSONA`
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `security_requirement`, `recovery_requirement`
- `acceptance relevance`: 生成済みペルソナ、phase result、失敗時の後追い確認に関係する。
- `adoption hint`: provider raw response ではなく、生成結果の業務要約と永続化結果を監査対象にする。
- `conflict hint`: 生成されたペルソナ本文は成果物として必要になる可能性がある。一方で raw request / response や full prompt は redaction 候補と衝突しうる。

### CAND-PGP-005 本文翻訳フェーズが参照する persona snapshot を確認できる

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:18`, `tasks/usecases/persona-generation-phase.yaml:24`, `tasks/usecases/body-translation-phase.yaml:9-16`, `tasks/usecases/body-translation-phase.yaml:21-30`, `docs/spec.md:101-102`, `docs/spec.md:130-131`, `docs/spec.md:225-227`, `docs/er.md:69`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-005`
- `actor`: 利用者 / 運用確認者
- `trigger`: NPC ペルソナ生成フェーズが完了し、本文翻訳フェーズの入力 summary を作る。
- `expected outcome`: 本文翻訳フェーズが参照する persona snapshot ID、persona count、missing count、参照不可理由、snapshot digest を確認できる。
- `observable point`: phase result、body translation input summary、Job Run summary、`PHASE_RUN_PERSONA`
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `state_requirement`, `testability_requirement`
- `acceptance relevance`: 生成済みペルソナを本文翻訳フェーズの入力として参照できる条件に直接関係する。
- `adoption hint`: 本文翻訳の実行や訳文生成は扱わず、参照可能性の監査だけを候補にする。
- `conflict hint`: snapshot が共通ペルソナ参照かジョブ内ペルソナ参照かは、共通ペルソナ hit/miss の統合判断に依存する。

### CAND-PGP-006 phase 監査要約で成功、失敗、再開を確認できる

- `source requirement`: `./plan.md:22`, `tasks/usecases/persona-generation-phase.yaml:21-25`, `tasks/usecases/body-translation-phase.yaml:27-28`, `docs/spec.md:53-54`, `docs/spec.md:178-196`, `docs/er.md:63-67`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:93-98`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:346-372`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-006`
- `actor`: 利用者 / 運用確認者 / 障害調査者
- `trigger`: NPC ペルソナ生成フェーズが success、failure、pause、resume、retry のいずれかで状態更新する。
- `expected outcome`: phase run ID、state、progress、target count、hit count、miss count、generated count、error kind、retryable flag、latest error summary を確認できる。
- `observable point`: structured audit summary、Job Run phase result、runtime event summary、`JOB_PHASE_RUN`
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `state_requirement`, `testability_requirement`
- `acceptance relevance`: current phase、progress、phase result、recoverable failure の確認に関係する。
- `adoption hint`: 単語翻訳フェーズの監査要約と redaction の先行判断を、ペルソナ生成フェーズの監査候補として継承する。
- `conflict hint`: 監査要約を UI、structured log、DB のどこまでに保存するかは未決である。保持期間と粒度は designer の質問候補になりうる。

### CAND-PGP-007 redaction で secret、raw payload、過剰本文を保存しない

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:12-15`, `docs/spec.md:57-58`, `docs/er.md:84`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:95-98`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:352-365`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-007`
- `actor`: 運用確認者 / セキュリティ確認者
- `trigger`: UI summary、error summary、structured log、phase result、永続化 summary を作る。
- `expected outcome`: API key、secret 本体、復号可能な値、provider raw request / response、full prompt、原文発話全文、会話文脈全文を保存または表示しないことを確認できる。
- `observable point`: redaction assertion、Job Run summary、structured log capture、fake secret store assertion
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `data_requirement`, `testability_requirement`
- `acceptance relevance`: 監査要約と障害調査材料を残しながら、保存禁止情報を露出しない条件に関係する。
- `adoption hint`: 保存する情報は ID、digest、version、件数、短い error kind、credential ref の参照状態に限定する候補として扱う。
- `conflict hint`: 再現性を高めるための入力 snapshot と、保存禁止情報の最小化が衝突しうる。全文保存の可否は AI が確定しない。

### CAND-PGP-008 ジョブ内ペルソナ flush 後も履歴を確認できる

- `source requirement`: `docs/spec.md:24`, `docs/spec.md:34`, `docs/spec.md:59`, `docs/er.md:53-59`, `docs/spec.md:215-216`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-PGP-008`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: 翻訳完了後、ジョブ内ペルソナを flush 対象にする。
- `expected outcome`: flush 対象の job persona count、共通ペルソナ非削除、本文翻訳参照済み snapshot digest、削除後に残す監査要約を確認できる。
- `observable point`: flush summary、persona lifecycle summary、`PERSONA`、`PHASE_RUN_PERSONA`
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `recovery_requirement`, `security_requirement`
- `acceptance relevance`: ジョブ内ペルソナを共通ペルソナと分離し、翻訳完了後に flush 可能にする条件に関係する。
- `adoption hint`: persona-generation-phase の直接完了後ではなく、後続 close / cleanup 時の監査候補として残す。
- `conflict hint`: flush 後に persona 本文、snapshot digest、監査要約のどれを残すかは保持方針と redaction 方針に依存する。

## Open Notes

- `human decision candidate`: 入力 snapshot に原文発話または会話文脈の全文を含めるか、ID / digest / 短い要約だけにするかは未決である。
- `human decision candidate`: 共通ペルソナ hit 時に job persona を作らないのか、body translation 用 snapshot だけを job に固定するのかは未決である。
- `human decision candidate`: 監査要約の保存場所、保持期間、UI 表示粒度は未決である。
- `merge candidate`: CAND-PGP-002、CAND-PGP-004、CAND-PGP-007 は redaction と再現材料の同一シナリオへ統合される可能性がある。
- `merge candidate`: CAND-PGP-003 と CAND-PGP-005 は persona snapshot 参照可否の同一シナリオへ統合される可能性がある。
- `rejection candidate`: CAND-PGP-008 は persona-generation-phase の直接受け入れ条件外として、後続 cleanup / output artifact 側へ移される可能性がある。
