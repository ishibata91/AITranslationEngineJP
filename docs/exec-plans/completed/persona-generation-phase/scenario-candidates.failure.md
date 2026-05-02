# Scenario Candidates: persona-generation-phase / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP`

## Generator Scope

- `viewpoint`: 失敗、入力不備、参照不能、整合性違反、保存失敗、回復、secret / redaction。
- `included_sources`: `plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`。
- `excluded_sources`: 他観点 candidate、最終 scenario 表、implementation-scope、プロダクトコード、プロダクトテスト、docs 正本更新。
- `generation_notes`: 採否、統合、競合解消は `designer` に残す。候補は failure 観点だけに限定する。

## Candidate Scenarios

### CAND-PGP-001 フェーズ開始時に必要入力が不足している

- `source requirement`: NPC ペルソナ生成フェーズは翻訳ジョブ、NPC 発話の原文、NPC 属性メタデータ、会話文脈、共通ペルソナを入力にする。単語翻訳フェーズ完了が前提である。`tasks/usecases/persona-generation-phase.yaml:9`, `tasks/usecases/persona-generation-phase.yaml:10`, `tasks/usecases/persona-generation-phase.yaml:20`, `docs/exec-plans/active/persona-generation-phase/plan.md:62`
- `viewpoint`: 失敗入力。
- `candidate scenario id`: `CAND-PGP-001`
- `actor`: Job Run から NPC ペルソナ生成フェーズを開始するユーザー。
- `trigger`: Ready 相当の job に見えるが、対象 NPC 発話、会話文脈、または単語翻訳フェーズ完了状態のいずれかが不足した状態で開始する。
- `expected outcome`: フェーズ実行を開始済み扱いにせず、不足入力名と開始不可理由を返す。`JOB_PHASE_RUN` や persona snapshot を成功状態で作らない。
- `observable point`: Job Run の開始結果、phase result、`JOB_PHASE_RUN` 件数。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `data_requirement`
- `acceptance relevance`: current phase、progress、phase result を確認できる要件と、本文翻訳フェーズの入力として参照できる要件の前提を守る。
- `adoption hint`: phase start boundary の拒否系として、UI と API の両方に関係する可能性がある。
- `conflict hint`: lifecycle 観点が「空対象なら Completed」とする場合、発話 0 件と入力欠落の扱いが競合する可能性がある。

### CAND-PGP-002 NPC 属性メタデータ欠落時にペルソナを捏造しない

- `source requirement`: NPC ペルソナは NPC の発言、種族、性別情報を元に生成する。NPC 属性メタデータは種族、性別、立場、性格傾向などを含む。`docs/spec.md:22`, `docs/spec.md:212`, `tasks/usecases/persona-generation-phase.yaml:12`, `tasks/usecases/persona-generation-phase.yaml:13`
- `viewpoint`: 失敗入力。
- `candidate scenario id`: `CAND-PGP-002`
- `actor`: NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: 生成対象 NPC に原文発話はあるが、種族、性別、または NPC_PROFILE 同一性に必要な属性が欠落している。
- `expected outcome`: 欠落属性を補完推測せず、対象 NPC を生成済みペルソナとして保存しない。生成不能な NPC と retryable かどうかを phase result に残す。
- `observable point`: provider request payload、phase result の対象 NPC 別 error、`PERSONA` row count。
- `related detail requirement type`: `boundary_requirement`, `failure_handling_requirement`, `testability_requirement`
- `acceptance relevance`: NPC ごとの生成対象確認と、生成済みペルソナを本文翻訳フェーズ入力にする要件の誤成立を防ぐ。
- `adoption hint`: per NPC の失敗扱いと phase 全体の失敗扱いを分ける候補として使える。
- `conflict hint`: 一部 NPC だけ失敗した時に phase 全体を RecoverableFailed にするか、部分完了として扱うかは人間判断候補になる。

### CAND-PGP-003 共通ペルソナとジョブ内ペルソナの整合性が崩れている

- `source requirement`: ジョブ内ペルソナは共通ペルソナと分離して保持する。`PERSONA` は `NPC_PROFILE` と 1:1 で紐づき、同一 NPC プロファイルに共通とジョブ内を同時保持しない。`tasks/usecases/persona-generation-phase.yaml:15`, `tasks/usecases/persona-generation-phase.yaml:25`, `docs/spec.md:34`, `docs/er.md:50`, `docs/er.md:51`
- `viewpoint`: 設定不整合。
- `candidate scenario id`: `CAND-PGP-003`
- `actor`: NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: 同一 `NPC_PROFILE` に共通ペルソナとジョブ内ペルソナが同時に見える、または共通ペルソナの `scope` / `source` / `translation_job_id` が矛盾している。
- `expected outcome`: 矛盾した共通ペルソナを正常な既存ペルソナとして扱わず、ジョブ内ペルソナの重複作成もしない。整合性違反を phase result に残す。
- `observable point`: `PERSONA`、`PHASE_RUN_PERSONA`、phase result、本文翻訳フェーズ入力 summary。
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `failure_handling_requirement`
- `acceptance relevance`: ジョブ内ペルソナを共通ペルソナと分離して保持し、本文翻訳フェーズへ正しい persona snapshot を渡す要件に関係する。
- `adoption hint`: 共通ペルソナが存在する時の skip / reuse / error の境界を `designer` が固定する材料になる。
- `conflict hint`: actor-goal 観点が共通ペルソナ再利用を正常系にする場合、同一 profile の job 内生成可否と競合する可能性がある。

### CAND-PGP-004 provider 失敗または不正応答でペルソナを保存しない

- `source requirement`: 各フェーズでは provider / model を選択できる。失敗しても別 provider fallback は不要で、ジョブは失敗回復対象になる。過去の単語翻訳フェーズでは provider 失敗、応答不正、保存失敗を成功扱いにしない方針が固定されている。`docs/spec.md:52`, `docs/spec.md:53`, `docs/spec.md:55`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:25`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:265`
- `viewpoint`: 参照不能 / 回復動作。
- `candidate scenario id`: `CAND-PGP-004`
- `actor`: NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: AI provider が timeout、network error、invalid response、対象 NPC 欠落 response のいずれかを返す。
- `expected outcome`: 暗黙 fallback を行わず、invalid response をペルソナとして保存しない。対象 NPC 単位または phase 単位の error kind、retryable flag、progress を観測できる。
- `observable point`: AIProvider adapter output、fake transport log、phase result、`PERSONA` row count。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `testability_requirement`
- `acceptance relevance`: 生成済みペルソナを本文翻訳フェーズ入力にする要件と、失敗回復を継続できる要件に関係する。
- `adoption hint`: external-integration 観点の provider 境界候補と統合される可能性が高い。
- `conflict hint`: external-integration 観点が provider 別の細分ケースを出す場合、本候補は umbrella として merge 対象になる。

### CAND-PGP-005 保存失敗で partial persona state を成功扱いにしない

- `source requirement`: `PERSONA_FIELD_EVIDENCE` は生成根拠の翻訳フィールドを保持し、`PHASE_RUN_PERSONA` はフェーズ対象のペルソナを表す。ジョブ内ペルソナは本文翻訳フェーズの入力として参照される。`docs/er.md:55`, `docs/er.md:69`, `tasks/usecases/persona-generation-phase.yaml:17`, `tasks/usecases/persona-generation-phase.yaml:18`, `tasks/usecases/body-translation-phase.yaml:15`
- `viewpoint`: 保存失敗。
- `candidate scenario id`: `CAND-PGP-005`
- `actor`: NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: provider response は valid だが、`PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA`、persona snapshot のいずれかの保存に失敗する。
- `expected outcome`: 一部保存済みの状態を Completed としない。本文翻訳フェーズから参照可能な persona snapshot として公開せず、再試行可能な保存失敗として観測する。
- `observable point`: temp DB の row count、transaction 結果、phase result、本文翻訳開始可否。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `acceptance relevance`: 生成済みペルソナを本文翻訳フェーズの入力として参照できる要件の整合性を守る。
- `adoption hint`: persistence boundary の API テスト候補として使える。
- `conflict hint`: 回復系候補が「既存成功分を維持して未処理だけ進める」とする場合、保存単位と rollback 単位を揃える必要がある。

### CAND-PGP-006 参照不能な NPC / 会話文脈を provider request に入れない

- `source requirement`: 翻訳フィールドは発話者 NPC、親クエスト、会話トピックなどを参照する。NPC ペルソナ生成フェーズは NPC の原文発話、属性メタデータ、会話文脈からジョブ内ペルソナを生成する。`docs/er.md:45`, `docs/er.md:46`, `tasks/usecases/persona-generation-phase.yaml:12`, `tasks/usecases/persona-generation-phase.yaml:14`, `docs/exec-plans/active/persona-generation-phase/plan.md:13`
- `viewpoint`: 参照不能。
- `candidate scenario id`: `CAND-PGP-006`
- `actor`: NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: 翻訳フィールドから発話者 NPC、親クエスト、会話トピック、または NPC_PROFILE へ解決できない参照がある。
- `expected outcome`: orphan な会話文脈を provider request に混ぜない。参照不能な対象と理由を phase result に残し、生成済み扱いの persona snapshot を作らない。
- `observable point`: reference resolution result、provider request payload、phase result、`PERSONA_FIELD_EVIDENCE`。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `observability_requirement`
- `acceptance relevance`: NPC ごとの生成対象確認と、本文翻訳フェーズで参照する persona snapshot の信頼性に関係する。
- `adoption hint`: 入力解決 boundary と provider request composition boundary の両方にまたがる候補として扱う。
- `conflict hint`: operation-audit 観点が参照不能の監査要約を出す場合、保存対象の粒度が競合する可能性がある。

### CAND-PGP-007 secret と provider raw payload を露出しない

- `source requirement`: APIKey は暗号化保存し、`credential_ref` は secret store への参照だけを保持する。過去の単語翻訳フェーズでは secret、provider raw request / response、本文全文を UI、error summary、structured log に出さない方針が固定されている。`docs/spec.md:57`, `docs/spec.md:58`, `docs/er.md:84`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:93`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:98`
- `viewpoint`: 失敗入力 / 設定不整合 / 回復動作。
- `candidate scenario id`: `CAND-PGP-007`
- `actor`: provider 設定済みの Job Run で NPC ペルソナ生成フェーズを実行するユーザー。
- `trigger`: provider failure、invalid response、保存失敗、設定不整合の error summary または structured log が生成される。
- `expected outcome`: API key、secret 本体、復号可能な値、provider raw request / response、NPC 発話全文を UI、error summary、structured log に出さない。provider、model、credential ref、input count、error kind などの要約だけを確認できる。
- `observable point`: Job Run error summary、structured log、fake secret store assertion、fake transport log。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `failure_handling_requirement`
- `acceptance relevance`: 失敗時の観測可能性を保ちながら、APIKey 保存要件と過剰本文非露出の回帰を防ぐ。
- `adoption hint`: security / observability の横断候補として、provider failure 候補と別に残す価値がある。
- `conflict hint`: operation-audit 観点が障害再現材料として raw payload 保存を要求する場合、redaction 方針と衝突する。

## Open Notes

- `human decision candidate`: 一部 NPC だけ失敗した時に phase 全体を RecoverableFailed にするか、成功分を維持した partial state とするかは未確定である。
- `human decision candidate`: 共通ペルソナがある NPC をジョブ内生成で skip するか、明示的な参照 link を作るかは未確定である。
- `merge candidate`: `CAND-PGP-004` は external-integration 観点の provider 失敗候補と統合される可能性がある。
- `merge candidate`: `CAND-PGP-007` は operation-audit 観点の redaction / audit 候補と統合される可能性がある。
- `rejection candidate`: 正常系の「生成対象 0 件」の扱いは failure 観点だけでは確定しないため、本ファイルでは採否判断しない。
