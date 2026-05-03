# Scenario Candidates: translation-job-setup-phase-provider-settings / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS-PPS`
- `candidate_artifact`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-candidates.operation-audit.md`

## Generator Scope

- `viewpoint`: operation-audit
- `running_task_artifact_location`: `docs/exec-plans/active/translation-job-setup-phase-provider-settings/`
- `target_diff`: Job Setup の AI runtime 設定を単一 master-persona 参照から phase 別 provider / model / credential / execution mode / batch mode へ分離する。
- `included_sources`:
  - `docs/exec-plans/active/translation-job-setup-phase-provider-settings/plan.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本、他 generator 成果物、final scenario matrix。
- `generation_notes`: 採否、統合、競合解消は designer へ残す。operation-audit 観点の候補だけを列挙する。

## Candidate Scenarios

### CAND-TJS-PPS-001 job 作成後に phase 別 runtime 要約を確認する

- `source requirement`:
  - `plan.md:8-9` phase ごとの provider / model / credential / batch mode を設定し、master-persona provider 設定から切り離す。
  - `plan.md:59-60` Job Setup は master-persona の provider 設定を既定値または保存元にせず、phase ごとに provider、model、credential 参照、execution mode を持つ。
  - `scenario-design.md:24` AI 基盤設定は provider、model、credential 参照、実行方式を区別し、API key 平文を表示または保存要約に出さない。
  - `term-translation-phase.md:24` 単語翻訳 phase の主要データは provider / model / execution mode の要約を含む。
  - `persona-generation-phase.md:31` NPC ペルソナ生成 phase は Job Setup の persona 専用設定を継承する。
  - `body-translation-phase.md:27` 本文翻訳 phase は Job Setup の本文翻訳用 provider、model、execution mode を使う。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-001`
- `actor`: Job Setup を完了した利用者、運用確認者。
- `trigger`: validation pass 後に job を作成し、作成済み job の read-only 要約を再表示する。
- `audit event`: job 作成完了時に、単語翻訳、NPC ペルソナ生成、本文翻訳の各 phase の実行設定断面を確認する。
- `saved summary`: phase 名、provider、model、execution mode、credential 参照状態、batch mode 選択、validation pass 断面、作成時刻。
- `redaction rule`: API key 平文、secret 本体、復号可能な値、provider raw request / response は保存要約と監査表示へ出さない。
- `expected outcome`: 作成後の job 要約で、3 phase の provider / model / execution mode が phase 別に追える。
- `observable point`: create result、job detail の read-only 要約、phase 実行開始前の設定断面。
- `related detail requirement type`: operation、display、persistence、redaction。
- `adoption hint`: 既存 `SCN-TJS-001` と `SCN-TJS-004` の作成後要約を phase 別設定へ拡張する候補。
- `conflict hint`: 保存単位が単一 runtime 前提の既存 Job Setup と衝突する。保存 schema や DTO の確定は implementation-scope へ残す。

### CAND-TJS-PPS-002 credential 参照状態を表示し secret を露出しない

- `source requirement`:
  - `plan.md:60` Job Setup は phase ごとに credential 参照を持つ。
  - `scenario-design.md:218` API key 平文は UI、DB、validation summary、監査表示に出ない。
  - `ui-design.md:13` AI runtime selector は credential 参照状態を表示する。
  - `ui-design.md:33` API key 平文、secret 本体、復号可能な値を表示しない。
  - `term-translation-phase.md:57-59` secret と raw request / response を出さず、監査要約には provider、model、execution mode などを残す。
  - `persona-generation-phase.md:66-70` 障害調査用の要約では credential ref と error kind を確認できる。
  - `body-translation-phase.md:42` secret、API key 平文、復号可能値、provider raw request / response、raw prompt は UI、DTO、error summary、structured log、debug log、fake transport log に出さない。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-002`
- `actor`: Job Setup を確認する利用者、障害調査を行う運用確認者。
- `trigger`: provider と model を phase 別に選び、credential が設定済みまたは未設定の状態で validation summary を確認する。
- `audit event`: phase 別 credential 参照状態と secret 非露出を確認する。
- `saved summary`: credential の有無、credential ref の識別子または digest、解決状態、失敗カテゴリ、確認時刻。
- `redaction rule`: API key 平文、secret 本体、復号可能な値、認証 header、外部 provider 応答原文を UI、DB summary、structured log、error summary に残さない。
- `expected outcome`: credential 参照状態は確認できるが、secret 値は画面、要約、ログ候補へ出ない。
- `observable point`: runtime selector、validation summary、error summary、structured log の redacted summary。
- `related detail requirement type`: display、external_integration、redaction、trust_boundary。
- `adoption hint`: 既存 `SCN-TJS-004` の secret 非露出確認を phase 別 credential 参照へ拡張する候補。
- `conflict hint`: credential ref の表示粒度、digest 化、保持期間は security / data requirement と衝突する可能性がある。

### CAND-TJS-PPS-003 model list API の取得結果と失敗理由を後追い確認する

- `source requirement`:
  - `plan.md:8` model 候補は provider の model list API から取得する。
  - `plan.md:61` API key が設定済みの場合だけ外部取得を試みる。
  - `plan.md:100` designer must include は model list 取得失敗時の UI 状態と API key 未設定 provider の候補表示可否を含む。
  - `scenario-design.md:62` validation 実行履歴は structured log に残し、アプリ状態は直近結果、対象設定断面、失敗カテゴリ、job 作成時の pass 断面だけ保持する。
  - `ui-design.md:14` validation summary は失敗理由、最終検証時刻または検証断面を表示する。
  - `ui-design.md:102-103` runtime と validation summary を同時に読み、設定変更時の validation 失効を確認する。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-003`
- `actor`: model 候補の取得状態を確認する利用者、外部連携失敗を調査する運用確認者。
- `trigger`: provider 選択後、model list API 取得が成功、失敗、または credential 未設定で未実行になる。
- `audit event`: provider ごとの model list 取得結果、失敗理由、取得時刻、取得対象 credential 参照状態を確認する。
- `saved summary`: provider、取得状態、model count、failure category、failure reason summary、retrieved_at、credential 参照状態、取得に使った provider capability version または digest。
- `redaction rule`: model list API の raw response、認証情報、provider error raw body、request header は保存要約に出さない。model 名は候補表示に必要な範囲だけ残す。
- `expected outcome`: 取得成功時は候補数と取得時刻を確認できる。失敗時は失敗理由と再取得可否を確認できる。credential 未設定時は外部取得未実行として区別できる。
- `observable point`: model selector、model list status、retrieved_at 表示、failure reason summary、dirty-validation 表示。
- `related detail requirement type`: external_integration、display、operation_audit、redaction。
- `adoption hint`: external-integration 候補と統合し、model list 取得を create 前 validation の前段観測点として扱う候補。
- `conflict hint`: model list 取得履歴を business state に残すか structured log に残すかは designer の統合判断が必要である。

### CAND-TJS-PPS-004 LM Studio では API key 欠落 warning を監査要約に残さない

- `source requirement`:
  - `plan.md:62` LM Studio は API key を要求しないため、API key 入力、API key 未設定 warning、credential select に出さない。
  - `plan.md:91` must include は LM Studio の API key 非表示を含む。
  - `scenario-design.md:62` 必須設定不足と credential 参照不能は blocking validation failure にする。
  - `ui-design.md:30-32` credential missing と provider / mode 不整合は区別し、credential 解決は blocking validation failure にする。
  - `scenario-design.md:218` API key 平文は UI、DB、validation summary、監査表示に出ない。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-004`
- `actor`: LM Studio を使う利用者、validation warning を確認する運用確認者。
- `trigger`: phase 別 provider のいずれかに LM Studio を選び、credential 未設定のまま validation と job 作成後要約を確認する。
- `audit event`: LM Studio の credential 不要状態が、warning や blocking failure として残っていないことを確認する。
- `saved summary`: provider は LM Studio、credential policy は not_required、credential warning count は 0、validation status、確認時刻。
- `redaction rule`: API key 欄、credential ref、secret 値、API key 欠落 warning を LM Studio の要約に出さない。
- `expected outcome`: LM Studio では API key 未設定 warning が UI、validation summary、job 作成後監査要約に残らない。
- `observable point`: runtime selector、validation summary、create result、job detail の phase 別要約。
- `related detail requirement type`: display、external_integration、operation_audit、redaction。
- `adoption hint`: failure 候補の credential missing と統合し、provider ごとの credential policy 例外として扱う候補。
- `conflict hint`: generic credential missing rule と LM Studio の not_required policy が衝突するため、provider capability の正本化は designer が整理する。

### CAND-TJS-PPS-005 batch mode の明示選択を監査要約に残す

- `source requirement`:
  - `plan.md:8` Job Setup は phase ごとの batch mode を設定できる。
  - `plan.md:63` batch mode は暗黙推定にせず、対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
  - `plan.md:91` must include は Gemini / xAI batch mode 明示切替を含む。
  - `term-translation-phase.md:36` Batch API を使う場合も batch item は 1 対象語単位にする。
  - `body-translation-phase.md:64` 本文翻訳 phase の表示項目は provider / model / execution mode 要約、request unit count、output count を含む。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-005`
- `actor`: batch mode を選ぶ利用者、処理結果を後追い確認する運用確認者。
- `trigger`: Gemini または xAI を選択し、batch mode を有効または無効として明示選択する。
- `audit event`: phase 別の batch mode 明示選択が job 作成時の設定断面に残ることを確認する。
- `saved summary`: phase 名、provider、model、execution mode、batch mode selected、batch 対象 provider 判定、request unit policy、作成時刻。
- `redaction rule`: batch request payload、provider raw request / response、prompt raw body は残さない。
- `expected outcome`: batch mode が暗黙推定ではなく、選択値として job 要約と phase 実行結果要約から確認できる。
- `observable point`: runtime selector の batch mode control、validation summary、create result、phase result summary。
- `related detail requirement type`: operation、external_integration、display、operation_audit。
- `adoption hint`: state-transition 候補の provider / mode 整合性と統合し、Gemini / xAI だけの明示選択として扱う候補。
- `conflict hint`: execution mode と batch mode の用語、保存単位、provider capability との関係は designer が整理する必要がある。

### CAND-TJS-PPS-006 master-persona provider 設定を参照していないことを監査する

- `source requirement`:
  - `plan.md:9` Job Setup を master-persona の provider 設定から切り離す。
  - `plan.md:52-53` 現状は master-persona provider / model から runtime option を 1 件作り、secret key を `master-persona:<provider>` として解決する。
  - `plan.md:59` Job Setup は master-persona の provider 設定を既定値または保存元として扱わない。
  - `scenario-design.md:24` AI 基盤設定は provider、model、credential 参照、実行方式を区別する。
  - `persona-generation-phase.md:31` NPC ペルソナ生成 phase は Job Setup の persona 専用設定を継承する。
  - `body-translation-phase.md:27` 本文翻訳 phase は Job Setup の本文翻訳用 provider、model、execution mode を使う。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJS-PPS-006`
- `actor`: Job Setup の設定出所を確認する運用確認者。
- `trigger`: master-persona provider 設定が存在する環境で、Job Setup の phase 別 provider 設定を作成し、job 作成後要約を確認する。
- `audit event`: Job Setup の phase 別 runtime が master-persona provider / model / secret key を参照していないことを確認する。
- `saved summary`: phase 名、provider、model、credential ref、setting_source、master_persona_runtime_ref_used=false、作成時刻。
- `redaction rule`: master-persona 側 secret key、API key 平文、復号可能な値は比較結果や error summary に出さない。
- `expected outcome`: job 作成後の設定断面は Job Setup の phase 別入力から構成され、master-persona 設定を既定値、保存元、secret 解決元として扱わない。
- `observable point`: create result、job detail の setting_source、validation summary、credential ref resolution summary。
- `related detail requirement type`: responsibility_boundary、external_integration、operation_audit、redaction。
- `adoption hint`: actor-goal または lifecycle 候補の phase 別設定保存と統合し、設定出所の監査確認点として扱う候補。
- `conflict hint`: 共通ペルソナ参照と master-persona provider 設定を混同する可能性がある。名称と責務境界の整理は designer へ残す。

## Open Notes

- `human decision candidate`:
  - model list API の取得結果を business state に保持するか、structured log と直近 UI 状態に限定するか。
  - credential ref の表示粒度を ID、digest、存在状態のどれにするか。
  - `setting_source` と `master_persona_runtime_ref_used=false` を保存要約に持つか、検証用の観測点に限定するか。
- `merge candidate`:
  - `CAND-TJS-PPS-001` は lifecycle の job 作成後要約候補と統合候補である。
  - `CAND-TJS-PPS-003` は external-integration の model list API 候補と統合候補である。
  - `CAND-TJS-PPS-004` と `CAND-TJS-PPS-005` は failure または state-transition の provider capability 候補と統合候補である。
- `rejection candidate`:
  - operation-audit 観点では不採用確定なし。採否は designer が決める。
- `conflict candidate`:
  - credential missing の一般 rule と LM Studio の credential not_required rule の衝突。
  - model list API 取得履歴の保存場所と保持粒度。
  - batch mode と execution mode の保存単位。
  - 共通ペルソナ参照と master-persona provider 設定の名称衝突。
