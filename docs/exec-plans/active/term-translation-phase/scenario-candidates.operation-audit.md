# Scenario Candidates: term-translation-phase / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`:
  - `./plan.md`
  - `tasks/usecases/term-translation-phase.yaml`
  - `tasks/index.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/architecture.md`
  - `docs/exec-plans/completed/translation-job-setup/plan.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
  - `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
- `excluded_sources`:
  - プロダクトコード
  - プロダクトテスト
  - docs 正本化
  - 他 generator の候補成果物
- `generation_notes`: 運用確認、監査ログ、履歴、再現材料、保存禁止の候補だけを作る。採否、統合、最終シナリオ表の確定は designer に残す。

## Candidate Scenarios

### CAND-TTP-001 単語翻訳フェーズ開始と対象件数を後追い確認する

- `source requirement`: `./plan.md:50-52`、`tasks/usecases/term-translation-phase.yaml:20-25`、`docs/spec.md:53-55`、`docs/er.md:63-69`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-001`
- `actor`: ユーザー、運用者
- `trigger`: `Ready` job の Job Run から単語翻訳フェーズを開始する。
- `audit event`: `term_translation_phase_started`
- `saved summary`: job ID、phase run ID、phase type、開始時刻、翻訳フィールド件数、抽出された翻訳対象語件数、共通辞書 hit 件数、選択済み AI runtime 要約。
- `redaction rule`: API key、credential 平文、provider request 原文、翻訳フィールド本文の全文は保存しない。本文由来の情報は件数、ID、digest だけにする。
- `expected outcome`: 現在 phase、progress、開始時の対象件数を Job Run と structured log から後追い確認できる。
- `observable point`: Job Run の current phase / progress、`JOB_PHASE_RUN`、structured log。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `data_requirement`
- `adoption hint`: designer は開始監査と progress 表示の受け入れ条件へ統合できる。
- `conflict hint`: lifecycle / state-transition 候補が phase 開始条件を別状態で固定する場合は統合時に前提をそろえる。

### CAND-TTP-002 共通辞書による除外と置換判定を監査する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:13-18`、`tasks/usecases/term-translation-phase.yaml:21-25`、`docs/spec.md:27-33`、`docs/er.md:56-59`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-002`
- `actor`: ユーザー、運用者
- `trigger`: 単語翻訳フェーズが翻訳対象語を共通辞書と照合する。
- `audit event`: `common_dictionary_exclusion_evaluated`
- `saved summary`: phase run ID、共通辞書 entry ID 群の参照、完全一致件数、除外件数、置換対象件数、未一致件数、内部出力ステータス `cached` の件数。
- `redaction rule`: structured log には全語彙リストや全訳語を重複保存しない。監査表示で語彙が必要な場合は `DICTIONARY_ENTRY` を参照して表示する。
- `expected outcome`: 共通辞書で除外された語と、用語翻訳へ送られた語の内訳を後から確認できる。
- `observable point`: phase result summary、Job Run result、`DICTIONARY_ENTRY`、structured log。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: designer は共通辞書 hit の確認ケース、または置換対象判定の監査観点へ統合できる。
- `conflict hint`: 実語と訳語を監査ログへ保存する方針は、保存禁止とデータ保持粒度に衝突する可能性がある。

### CAND-TTP-003 確定訳語のジョブ内辞書反映を後追い確認する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:15-18`、`tasks/usecases/term-translation-phase.yaml:23-25`、`docs/spec.md:33-35`、`docs/spec.md:217-220`、`docs/er.md:56-59`、`docs/er.md:69`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-003`
- `actor`: ユーザー、運用者
- `trigger`: 用語や固有名詞の訳語を確定し、ジョブ内辞書へ反映する。
- `audit event`: `job_dictionary_entries_confirmed`
- `saved summary`: job ID、phase run ID、作成または更新した job-scoped dictionary entry ID、確定訳語件数、dictionary scope、dictionary source、反映完了時刻。
- `redaction rule`: provider raw response、API key、翻訳フィールド本文の全文は保存しない。確定訳語そのものは業務データとして `DICTIONARY_ENTRY` に保持し、監査ログへ重複保存しない。
- `expected outcome`: どの phase run がどのジョブ内辞書 entry を作ったか、後から追跡できる。
- `observable point`: Job Run phase result、ジョブ内辞書確認画面、`DICTIONARY_ENTRY`、`PHASE_RUN_DICTIONARY_ENTRY`。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: designer は確定訳語確認、辞書反映確認、本文翻訳入力の前提確認へ統合できる。
- `conflict hint`: full text の監査ログ保存を求める候補がある場合は、業務データ保存と監査要約保存の境界を分ける必要がある。

### CAND-TTP-004 AI 実行設定と入力要約を再現材料として残す

- `source requirement`: `docs/spec.md:55-58`、`docs/er.md:25`、`docs/er.md:84`、`docs/architecture.md:124-140`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md:211-239`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-004`
- `actor`: 運用者、調査者
- `trigger`: 単語翻訳フェーズが AI provider に用語翻訳を依頼する。
- `audit event`: `term_translation_ai_execution_recorded`
- `saved summary`: phase run ID、provider、model、execution mode、instruction / prompt template version、候補抽出 rule version、入力語件数、出力語件数、credential 参照状態、外部 request 実行有無。
- `redaction rule`: API key、credential 平文、復号可能な secret、provider raw request / response、本文全量を保存しない。prompt は本文を含む可能性があるため、version と digest だけを保存候補にする。
- `expected outcome`: 障害調査時に、同じ provider / model / instruction version / 入力要約で実行条件を説明できる。
- `observable point`: `JOB_PHASE_RUN` の最終 AI 実行情報、AIProvider 境界の structured log、fake transport log。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: designer は paid API を使わない検証、fake transport 証跡、再現材料の候補へ統合できる。
- `conflict hint`: 完全再現のために full prompt や raw response を保存する案は、保存禁止と衝突する可能性がある。

### CAND-TTP-005 失敗・中断・再実行準備の監査要約を残す

- `source requirement`: `docs/spec.md:53-54`、`docs/spec.md:171-199`、`docs/er.md:63-67`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md:92-103`、`docs/exec-plans/completed/translation-job-setup/implementation-scope.md:24-28`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-005`
- `actor`: ユーザー、運用者
- `trigger`: AI 実行失敗、辞書反映失敗、保存失敗、ユーザー中断、再実行準備が発生する。
- `audit event`: `term_translation_phase_failed_or_paused`
- `saved summary`: phase run ID、状態、error kind、失敗カテゴリ、retryable flag、影響を受けた語件数、最後に成功した step、発生時刻。
- `redaction rule`: provider error 原文、外部応答原文、secret、翻訳フィールド本文の全文は保存しない。ユーザーに見せる error summary は分類と短い理由にする。
- `expected outcome`: RecoverableFailed / Failed / Paused の理由と再開可否を後から説明できる。
- `observable point`: Job Run error summary、`JOB_PHASE_RUN` state、structured log、永続化状態。
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `state_requirement`
- `adoption hint`: designer は failure / state-transition 候補と統合し、監査要約だけをこの候補から採れる。
- `conflict hint`: ER は Attempt 履歴テーブルを持たないため、再試行履歴を business table として保存する案とは衝突する。

### CAND-TTP-006 本文翻訳前に辞書 snapshot 参照を照合する

- `source requirement`: `tasks/index.yaml:4-11`、`tasks/usecases/term-translation-phase.yaml:21-25`、`docs/spec.md:100-115`、`docs/spec.md:128-133`、`docs/spec.md:224-248`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-006`
- `actor`: ユーザー、運用者
- `trigger`: 単語翻訳フェーズが完了し、本文翻訳フェーズへ進む前に確定訳語の参照状態を確認する。
- `audit event`: `term_dictionary_snapshot_ready_for_body_translation`
- `saved summary`: term phase run ID、ジョブ内辞書 entry ID set の digest、確定訳語件数、除外済み共通辞書件数、完了時刻、本文翻訳へ渡す参照 snapshot ID または digest。
- `redaction rule`: structured log には語彙全量を保存しない。本文翻訳側で必要な語彙はジョブ内辞書を参照する。
- `expected outcome`: 本文翻訳フェーズが参照する確定訳語 snapshot を、phase 結果から説明できる。
- `observable point`: Job Run phase result、ジョブ内辞書確認画面、後続 phase 開始前 summary。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: designer は本文翻訳への入力確認、phase 完了条件、Job Run result 表示へ統合できる。
- `conflict hint`: 本文翻訳フェーズ側の候補と責務が重なるため、最終シナリオでは cross-phase 参照だけを残すか判断が必要である。

### CAND-TTP-007 監査表示とログで保存禁止情報を露出しない

- `source requirement`: `docs/spec.md:57-58`、`docs/er.md:84`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md:24-25`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md:211-239`、`docs/exec-plans/completed/translation-job-setup/implementation-scope.md:67-71`、`docs/exec-plans/completed/translation-job-setup/implementation-scope.md:167-168`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-007`
- `actor`: ユーザー、運用者、調査者
- `trigger`: Job Run の phase result、error summary、structured log、外部 request 証跡を表示または確認する。
- `audit event`: `term_translation_audit_redaction_verified`
- `saved summary`: redaction policy version、credential 参照状態、provider / model 名、表示対象の summary ID、保存禁止情報の非出力確認結果。
- `redaction rule`: API key、secret 本体、復号可能な値、provider raw response、本文全量、過剰な候補語一覧を UI、error summary、console、structured log に出さない。
- `expected outcome`: 監査確認に必要な情報を見せながら、secret と過剰本文が露出しない。
- `observable point`: Job Run result、error summary、console / structured log、fake secret store fixture。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `data_requirement`
- `adoption hint`: designer は secret 非露出の監査条件として、external-integration / trust-boundary 系候補と統合できる。
- `conflict hint`: 調査容易性のために raw response や full prompt の保存を求める候補と衝突する可能性がある。

### CAND-TTP-008 共通辞書参照 snapshot を監査する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:10-18`、`docs/spec.md:27-35`、`docs/er.md:24-25`、`docs/er.md:56-59`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md:92-98`、`docs/exec-plans/completed/translation-job-setup/implementation-scope.md:128`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TTP-008`
- `actor`: 運用者、調査者
- `trigger`: 単語翻訳フェーズが共通辞書を読み、除外対象や置換対象を決める。
- `audit event`: `shared_dictionary_snapshot_referenced`
- `saved summary`: phase run ID、共通辞書 snapshot ID または digest、参照 entry 件数、matched entry 件数、snapshot created_at / observed_at、基盤参照状態。
- `redaction rule`: structured log へ共通辞書の全 entry 内容を複製しない。必要な表示は共通辞書 entry を参照して行う。
- `expected outcome`: 後から共通辞書が更新されても、単語翻訳フェーズがどの辞書断面で判定したか説明できる。
- `observable point`: Job Run phase result、common dictionary reference summary、structured log、`DICTIONARY_ENTRY`。
- `related detail requirement type`: `observability_requirement`, `concurrency_requirement`, `consistency_requirement`
- `adoption hint`: designer は phase 実行中の共通基盤 lock / snapshot 未決の質問候補として扱える。
- `conflict hint`: lock、snapshot、変更検知 failure のどれを採るかは state-transition / consistency 候補と統合が必要である。

## Open Notes

- `human decision candidate`:
  - `CAND-TTP-002`: 監査表示で実語と訳語を表示する範囲。候補は structured log には ID と件数だけを保存する前提にしている。
  - `CAND-TTP-004`: full prompt / raw response を保存しない場合の再現性の十分性。候補は version と digest に寄せている。
  - `CAND-TTP-006`: 本文翻訳フェーズとの cross-phase 参照を、この task でどこまで固定するか。
  - `CAND-TTP-008`: 共通辞書参照中の lock、snapshot、変更検知 failure のどれを採るか。
- `merge candidate`:
  - `CAND-TTP-001` は lifecycle / state-transition の phase start 候補と統合候補である。
  - `CAND-TTP-005` は failure / state-transition の失敗・再開候補と統合候補である。
  - `CAND-TTP-007` は external-integration / trust-boundary 寄りの secret 非露出候補と統合候補である。
- `rejection candidate`:
  - なし。採否は designer が行う。
