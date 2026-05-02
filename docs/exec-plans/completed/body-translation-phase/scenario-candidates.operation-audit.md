# Scenario Candidates: body-translation-phase / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`:
  - `./plan.md`
  - `tasks/index.yaml`
  - `tasks/usecases/body-translation-phase.yaml`
  - `tasks/usecases/term-translation-phase.yaml`
  - `tasks/usecases/persona-generation-phase.yaml`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/architecture.md`
  - `docs/diagrams/er/combined-data-model-er.puml`
  - `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
  - `docs/exec-plans/completed/persona-generation-phase/scenario-design.md`
- `excluded_sources`: product code, product tests, docs 正本更新, 他 generator 成果物, 最終 scenario 採否, 統合判断
- `generation_notes`: 本文翻訳フェーズの運用確認、監査ログ、履歴、再現材料、保存禁止だけを候補化する。訳文と出力ステータスは業務データとしての保存要件がある一方、監査要約へ本文全量、raw payload、secret を重複保存する前提は置かない。

## Candidate Scenarios

### CAND-BTP-001 本文翻訳フェーズ開始時の入力 snapshot を後追い確認する

- `source requirement`: `./plan.md:13-16`、`./plan.md:71-78`、`tasks/index.yaml:4-11`、`tasks/usecases/body-translation-phase.yaml:9-23`、`docs/spec.md:100-115`、`docs/spec.md:128-133`、`docs/er.md:22-25`、`docs/er.md:63-69`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-001`
- `actor`: 利用者 / 運用確認者
- `trigger`: Job Run から本文翻訳フェーズを開始する。
- `audit event`: `body_translation_phase_started`
- `saved summary`: job ID、phase run ID、phase type、開始時刻、対象翻訳フィールド件数、term phase run ID、persona phase run ID、ジョブ内辞書 snapshot digest、persona snapshot digest、翻訳補助メタデータ snapshot digest、選択済み provider / model / execution mode。
- `redaction rule`: API key、secret 本体、provider raw request / response、full prompt、原文本文全文、訳文全文、ジョブ内辞書全文、persona 本文全文、翻訳補助メタデータ全文は監査要約へ保存しない。
- `expected outcome`: どの入力断面で本文翻訳フェーズを開始したかを、後から Job Run と phase run summary で確認できる。
- `observable point`: Job Run current phase / progress、`JOB_PHASE_RUN`、`PHASE_RUN_TRANSLATION_FIELD`、input snapshot summary。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: designer は phase 開始条件、入力 summary、cross-phase 参照の受け入れ条件へ統合できる。
- `conflict hint`: state-transition / lifecycle 候補が開始可能状態を固定する場合、operation-audit 側は入力断面の後追い確認だけを採る。

### CAND-BTP-002 翻訳指示の構成根拠を再現材料として確認する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:24-26`、`docs/spec.md:41-43`、`docs/spec.md:208-229`、`docs/er.md:42-46`、`docs/er.md:63-69`、`docs/architecture.md:124-140`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-002`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: 翻訳レコード種別、フィールド種別、辞書、ペルソナ、翻訳補助メタデータから本文翻訳用の翻訳指示を構成する。
- `audit event`: `body_translation_instruction_built`
- `saved summary`: phase run ID、instruction kind、prompt template version または digest、record type 件数、field type 件数、参照した dictionary snapshot digest、persona snapshot digest、metadata snapshot digest、構成成功 / 失敗件数。
- `redaction rule`: full prompt、原文本文全文、persona 本文全文、ジョブ内辞書全文、翻訳補助メタデータ全文、provider request body は保存しない。
- `expected outcome`: 障害調査時に、本文翻訳の指示構成がどの rule version と入力要約から作られたか説明できる。
- `observable point`: instruction summary、`JOB_PHASE_RUN.instruction_kind`、structured log、AIProvider 境界の fake transport log。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: external-integration 候補と統合する場合も、operation-audit では provider 呼び出し前の再現材料に限定する。
- `conflict hint`: 完全再現のために full prompt を保存する案は、保存禁止と衝突する可能性がある。

### CAND-BTP-003 翻訳フィールド単位の AI 実行結果を監査できる

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:10-20`、`tasks/usecases/body-translation-phase.yaml:24-28`、`docs/spec.md:11`、`docs/spec.md:41-43`、`docs/spec.md:53-58`、`docs/er.md:23-25`、`docs/diagrams/er/combined-data-model-er.puml:167-183`、`docs/architecture.md:129-140`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-003`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: 本文翻訳フェーズが AI provider へ翻訳フィールド本文の翻訳を依頼し、結果を受け取る。
- `audit event`: `body_translation_ai_execution_recorded`
- `saved summary`: phase run ID、provider、model、execution mode、credential ref、latest external run ID の参照、request unit count、success count、failure count、skipped count、output count、error kind。
- `redaction rule`: credential 平文、provider raw request / response、source text 全文、translated text 全文、full prompt は structured log と error summary へ保存しない。訳文は `JOB_TRANSLATION_FIELD` の業務データとして扱い、監査要約へ重複保存しない。
- `expected outcome`: どの AI 実行条件で何件処理され、どのカテゴリで失敗したかを確認できる。
- `observable point`: `JOB_PHASE_RUN` の AI 実行情報、AIProvider adapter summary、Job Run phase result、fake transport log。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: paid real API を検証前提にしない scenario へ統合し、fake transport で監査要約を検証できる。
- `conflict hint`: external-integration 候補が provider 接続や batch API 詳細を扱う場合、operation-audit 側は実行条件と redaction の後追い確認に絞る。

### CAND-BTP-004 訳文と出力ステータスの保存結果を後追い確認する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:17-20`、`tasks/usecases/body-translation-phase.yaml:33-37`、`tasks/usecases/translation-output-artifact.yaml:10-18`、`docs/spec.md:43`、`docs/spec.md:65-67`、`docs/er.md:23`、`docs/er.md:37-40`、`docs/diagrams/er/combined-data-model-er.puml:156-165`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-004`
- `actor`: 利用者 / 運用確認者
- `trigger`: AI 翻訳結果をジョブ内の翻訳フィールドへ反映する。
- `audit event`: `body_translation_field_output_persisted`
- `saved summary`: job ID、phase run ID、job translation field ID 群の digest、translated count、cached count、failed count、skipped count、output status 別件数、更新時刻、translation-output-artifact へ渡す result summary 参照。
- `redaction rule`: 監査ログには訳文全文、原文全文、provider raw response を重複保存しない。訳文と出力ステータスは `JOB_TRANSLATION_FIELD` の業務データとして保持し、監査要約は ID、件数、status、digest に限定する。
- `expected outcome`: 本文翻訳フェーズがどの翻訳フィールドを更新し、どの出力ステータスになったかを後から説明できる。
- `observable point`: Job Run result、`JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、output status summary。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: designer は訳文確認、出力ステータス確認、後続 output artifact の入力確認へ統合できる。
- `conflict hint`: 訳文を UI に表示することと、監査ログへ全文保存しないことを混同しないように統合が必要である。

### CAND-BTP-005 保護要素検証結果を監査できる

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:17-20`、`tasks/usecases/body-translation-phase.yaml:24-26`、`tasks/usecases/body-translation-phase.yaml:33-37`、`docs/spec.md:41-43`、`docs/spec.md:230-231`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-005`
- `actor`: 利用者 / 運用確認者 / 障害調査者
- `trigger`: 訳文生成後に、原文の埋め込み要素や保持すべき構造が損なわれていないか検証する。
- `audit event`: `body_translation_protected_element_validation_completed`
- `saved summary`: phase run ID、job translation field ID 群の digest、検出した保護要素件数、検証 pass count、mismatch count、missing count、extra count、validation rule version、failure field count。
- `redaction rule`: 原文本文全文、訳文全文、保護要素を含む本文断片を監査要約へ保存しない。保護要素値の保存が必要な場合でも、候補段階では digest、位置種別、件数に寄せる。
- `expected outcome`: 保護要素検証の合否と失敗カテゴリを Job Run と監査要約から確認できる。
- `observable point`: Job Run の保護要素検証結果、validation result summary、structured log、`JOB_TRANSLATION_FIELD`。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: designer は出力ステータスと保護要素検証を同じ結果確認シナリオへ統合できる。
- `conflict hint`: 保護要素の実値を残す方が調査しやすい一方、本文断片の過剰保存と衝突する可能性がある。

### CAND-BTP-006 pause、resume、retry、cancel と recoverable failure の監査要約を確認する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:23-29`、`tasks/usecases/body-translation-phase.yaml:33-37`、`docs/spec.md:53-54`、`docs/spec.md:135-199`、`docs/er.md:63-67`、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md:99-111`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-006`
- `actor`: 利用者 / 運用確認者 / 障害調査者
- `trigger`: 本文翻訳フェーズが pause、resume、retry、cancel、recoverable failure、failed のいずれかで状態更新する。
- `audit event`: `body_translation_phase_interrupted_or_recovered`
- `saved summary`: phase run ID、state、progress、retryable flag、affected field count、last successful step、latest error kind、発生時刻、再開可能な未処理 field count。
- `redaction rule`: provider error 原文、raw response、原文本文全文、訳文全文、secret は保存しない。error summary は短い分類、対象件数、再試行可否に限定する。
- `expected outcome`: 中断、再開、リトライ、キャンセル、回復可能失敗の理由と再開可否を後から説明できる。
- `observable point`: Job Run failure state、`JOB_PHASE_RUN.state`、`JOB_PHASE_RUN.latest_error`、structured log、runtime event summary。
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `state_requirement`, `security_requirement`
- `adoption hint`: failure / state-transition 候補と統合し、operation-audit からは監査要約と redaction 条件を採れる。
- `conflict hint`: ER は attempt 履歴テーブルを持たないため、再試行ごとの business history table を前提にする案とは衝突する。

### CAND-BTP-007 前段の辞書・ペルソナ参照と本文翻訳結果を突合できる

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:21-25`、`tasks/usecases/persona-generation-phase.yaml:21-25`、`tasks/usecases/body-translation-phase.yaml:9-20`、`docs/spec.md:29-35`、`docs/spec.md:129-130`、`docs/spec.md:246-248`、`docs/exec-plans/completed/term-translation-phase/scenario-design.md:72-98`、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md:78-111`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-007`
- `actor`: 運用確認者 / 障害調査者
- `trigger`: 本文翻訳フェーズの結果を、前段の確定訳語、ジョブ内辞書、persona snapshot と突合して確認する。
- `audit event`: `body_translation_cross_phase_reference_verified`
- `saved summary`: body phase run ID、term phase run ID、persona phase run ID、dictionary entry set digest、persona snapshot digest、metadata snapshot digest、translated field count、dictionary-applied count、persona-applied count、missing reference count。
- `redaction rule`: 辞書全量、persona 本文全文、翻訳補助メタデータ全文、原文本文全文、訳文全文を監査ログへ複製しない。必要な表示は業務データを参照して行う。
- `expected outcome`: 本文翻訳フェーズが、どの前段成果を参照して訳文を生成したかを説明できる。
- `observable point`: Job Run input / result summary、ジョブ内辞書確認、persona snapshot summary、`PHASE_RUN_DICTIONARY_ENTRY`、`PHASE_RUN_PERSONA`、`PHASE_RUN_TRANSLATION_FIELD`。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: designer は body phase readiness、結果確認、後続 output artifact の traceability へ統合できる。
- `conflict hint`: cross-phase 参照を body phase 側でどこまで固定するかは、term / persona 完了済み scenario と統合が必要である。

### CAND-BTP-008 監査表示とログで保存禁止情報を露出しない

- `source requirement`: `docs/spec.md:57-58`、`docs/er.md:84`、`docs/exec-plans/completed/term-translation-phase/scenario-design.md:93-98`、`docs/exec-plans/completed/term-translation-phase/scenario-design.md:346-372`、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md:106-111`、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md:339-363`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-008`
- `actor`: 運用確認者 / セキュリティ確認者
- `trigger`: Job Run summary、error summary、structured log、debug log、fake transport log、永続化 summary を確認する。
- `audit event`: `body_translation_audit_redaction_verified`
- `saved summary`: redaction policy version、provider / model / execution mode、credential ref、summary ID、検証対象 log 種別、保存禁止情報の非出力確認結果。
- `redaction rule`: API key、secret 本体、復号可能な値、provider raw request / response、full prompt、原文本文全文、訳文全文、ジョブ内辞書全文、persona 本文全文、翻訳補助メタデータ全文を UI、error summary、structured log、debug log、fake transport log に出さない。
- `expected outcome`: 監査確認に必要な provider、model、件数、digest、error kind は見えるが、secret と過剰本文は露出しない。
- `observable point`: Job Run summary、error summary、structured log capture、debug log capture、fake secret store assertion、fake transport log。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: designer は redaction の横断受け入れ条件として、external-integration / trust-boundary 寄り候補と統合できる。
- `conflict hint`: persona-generation-phase の先行 scenario は debug log で prompt / request body を確認できる余地を残している。本文翻訳では原文・訳文の過剰本文リスクが高いため、debug log 粒度は人間判断候補に残る。

### CAND-BTP-009 後続成果物出力へ渡す本文翻訳結果 summary を確認する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:17-20`、`tasks/usecases/body-translation-phase.yaml:33-37`、`tasks/usecases/translation-output-artifact.yaml:9-28`、`docs/spec.md:61-67`、`docs/er.md:71-77`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-BTP-009`
- `actor`: 利用者 / 運用確認者
- `trigger`: 本文翻訳フェーズが完了し、翻訳成果物出力フェーズへ進めるかを確認する。
- `audit event`: `body_translation_output_readiness_summarized`
- `saved summary`: body phase run ID、completed job ID、translated count、output status 別件数、protected validation pass / fail count、translation artifact 前提 summary ID、完了時刻。
- `redaction rule`: result summary には訳文全文と原文全文を監査ログとして複製しない。diff preview や output artifact が必要とする本文表示は後続 task の業務表示として扱い、この候補では保存しない。
- `expected outcome`: 後続の translation-output-artifact が本文翻訳結果を参照できる状態か、件数と失敗カテゴリで確認できる。
- `observable point`: Job Run result summary、Output Review への入力 summary、`JOB_TRANSLATION_FIELD`、`XTRANSLATOR_OUTPUT_ROW` への対応前提。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: designer は body phase 完了条件と output artifact precondition の橋渡し候補として扱える。
- `conflict hint`: translation-output-artifact task の候補と責務が重なるため、最終 scenario では「本文翻訳フェーズが渡す readiness summary」までに限定するか判断が必要である。

## Open Notes

- `human decision candidate`:
  - `CAND-BTP-002`: full prompt を保存しない場合の再現性を、prompt template version / digest / 入力 snapshot digest だけで十分とするか。
  - `CAND-BTP-005`: 保護要素の実値を監査要約へ保存するか、digest / 件数 / 位置種別だけにするか。
  - `CAND-BTP-008`: 本文翻訳フェーズの debug log で prompt / request body を確認可能にするか、過剰本文リスクを優先して digest と件数だけにするか。
  - `CAND-BTP-009`: 後続 output artifact との境界として、本文翻訳フェーズ側に result summary をどこまで保持するか。
- `merge candidate`:
  - `CAND-BTP-001` と `CAND-BTP-007` は cross-phase 入力確認の同一シナリオへ統合される可能性がある。
  - `CAND-BTP-003`、`CAND-BTP-006`、`CAND-BTP-008` は AI 実行監査と redaction の同一シナリオへ統合される可能性がある。
  - `CAND-BTP-004`、`CAND-BTP-005`、`CAND-BTP-009` は Job Run result と後続 output readiness の同一シナリオへ統合される可能性がある。
- `rejection candidate`:
  - provider raw request / response、full prompt、原文本文全文、訳文全文、辞書全文、persona 本文全文を監査ログまたは structured log に保存する案。
  - attempt 履歴テーブルや固定ログ形式をこの候補成果物で確定する案。
