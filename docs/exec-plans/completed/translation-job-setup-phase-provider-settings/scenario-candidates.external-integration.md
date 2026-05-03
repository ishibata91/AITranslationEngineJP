# Scenario Candidates: translation-job-setup-phase-provider-settings / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSPPS`
- `candidate_count`: 9

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./plan.md`, `../../completed/translation-job-setup/scenario-design.md`, `../../completed/translation-job-setup/ui-design.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本変更、最終 scenario matrix、候補の採否、候補の統合判断
- `generation_notes`: provider、secret、adapter、network、validation の外部境界だけを候補化する。API key 平文、secret 本体、provider raw request / response は出力しない。

## Candidate Scenarios

### CAND-TJSPPS-001 provider 別 getModels は secret 参照後だけ外部取得する

- `source requirement`: `plan.md:57-64` は provider ごとの model list API、API key 設定済み時だけの外部取得、LM Studio の API key 非要求、Gemini / xAI batch mode 明示を要求している。`scenario-design.md:64-70` は credential 解決、provider capability、fake 非表示を固定している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-001`
- `external boundary`: provider 別 getModels API と secret store 参照の境界。
- `actor`: Job Setup 利用者。
- `trigger`: API key が必要な provider を phase 設定で選び、model 候補取得を実行する。
- `start condition`: 対象 provider の secret store 参照が設定済みである。
- `expected outcome`: secret store から取得可否だけを判定した後、provider 別 getModels を呼び、model 候補を返す。UI、DTO、ログ、error summary に API key 平文は出ない。
- `observable point`: provider 選択状態、model list loading / success、secret redaction、外部 request 証跡。
- `fake_or_stub`: fake secret store、fake transport、fixed model list response。
- `related detail requirement type`: external_integration、secret 境界、provider 境界、network 境界。
- `adoption hint`: provider ごとの model 候補取得を acceptance に入れる場合の中心候補になる。
- `conflict hint`: API key 未設定時の候補取得可否は `CAND-TJSPPS-002` と組み合わせて designer が整理する。

### CAND-TJSPPS-002 API key 未設定 provider は外部取得せず状態だけ返す

- `source requirement`: `plan.md:61-62` は API key 設定済みの場合だけ外部取得し、LM Studio の API key warning を出さないことを要求している。`ui-design.md:28-34` は credential missing と secret 非露出を UI 状態として区別する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-002`
- `external boundary`: secret store 未設定判定と provider request 抑止の境界。
- `actor`: Job Setup 利用者。
- `trigger`: API key が必要な provider を選ぶが、secret store に参照がない状態で model 候補取得または validation を実行する。
- `start condition`: provider は API key 必須であり、credential 参照は未設定または参照不能である。
- `expected outcome`: 外部 getModels と validation 用 provider request は実行されない。画面には credential missing または参照不能の状態だけが返り、secret 値は表示されない。
- `observable point`: request log に外部 request がないこと、validation failure category、UI の credential 状態、secret redaction。
- `fake_or_stub`: fake secret store、request spy、invalid credential fixture。
- `related detail requirement type`: external_integration、secret 境界、validation 境界。
- `adoption hint`: API key 未設定時の安全な failure path として採用候補になる。
- `conflict hint`: failure 観点では必須設定不足として扱えるため、外部 request 抑止の観測点を残す必要がある。

### CAND-TJSPPS-003 LM Studio は API key なし provider として model 候補を扱う

- `source requirement`: `plan.md:8` と `plan.md:62` は LM Studio に API key 入力を出さないことを要求している。`plan.md:91` は LM Studio の API key 非表示を必須観点にしている。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-003`
- `external boundary`: API key なし provider の model list 取得境界。
- `actor`: Job Setup 利用者。
- `trigger`: LM Studio を phase 設定の provider として選び、model 候補を取得する。
- `start condition`: LM Studio は provider として選択済みであり、API key 入力、credential select、API key 未設定 warning は表示されていない。
- `expected outcome`: secret store 参照なしで LM Studio の model list 取得を試みる。失敗時も API key 未設定とは別の network / provider failure として返る。
- `observable point`: LM Studio 選択時の credential UI 非表示、secret store 未参照、model list request、network failure category。
- `fake_or_stub`: fake LM Studio endpoint、request spy、fixed local model response。
- `related detail requirement type`: external_integration、provider 境界、network 境界、display。
- `adoption hint`: LM Studio 固有の API key 非要求を scenario に明示する候補になる。
- `conflict hint`: network failure と API key missing を同じ failure 表示に統合すると、LM Studio の要件を失う可能性がある。

### CAND-TJSPPS-004 phase ごとの provider / model / credential 参照を独立して検証する

- `source requirement`: `plan.md:59-64` は Job Setup が phase ごとの provider、model、credential 参照、execution mode を持つことを要求している。`term-translation-phase.md:35-37`、`persona-generation-phase.md:30-32`、`body-translation-phase.md:27-32` は各 phase が provider 境界を使うことを示している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-004`
- `external boundary`: phase 別 runtime 設定と secret / provider validation の境界。
- `actor`: Job Setup 利用者。
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳に別々の provider / model / credential 参照を選び、validation を実行する。
- `start condition`: 3 phase の runtime 設定が独立して入力されている。
- `expected outcome`: validation は phase ごとの credential 解決、provider capability、model 参照を個別に判定する。1 phase の credential 不備は、その phase の failure reason として返る。
- `observable point`: phase 別 validation result、phase 別 credential state、phase 別 model list source、外部 request の provider / phase 対応。
- `fake_or_stub`: fake secret store、phase-tagged fake transport、provider capability fixture。
- `related detail requirement type`: external_integration、secret 境界、adapter 境界、validation 境界。
- `adoption hint`: master-persona 設定からの切り離しを外部連携面で確認する候補になる。
- `conflict hint`: state-transition 観点では validation stale 条件と統合される可能性がある。

### CAND-TJSPPS-005 fake provider は user-facing list に出さず transport だけ差し替える

- `source requirement`: `scenario-design.md:64-70` と `scenario-design.md:299-326` は fake provider を user-facing 選択肢に出さず、外部 request / SDK transport だけを fake に差し替えることを固定している。`body-translation-phase.md:31` は paid real API を検証前提にしないことを要求している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-005`
- `external boundary`: user-facing provider list と SDK transport seam の境界。
- `actor`: Job Setup 利用者、scenario acceptance を支える検証実行者。
- `trigger`: provider list 表示、model list 取得、validation を fake transport で実行する。
- `start condition`: provider list は real provider だけで構成され、transport は fake に差し替えられている。
- `expected outcome`: UI の provider list に fake provider は出ない。model list 取得と validation は同じ application 経路を通り、外部 request だけ fake transport に置き換わる。
- `observable point`: provider list、transport selection、request log、validation result、paid API 未実行証跡。
- `fake_or_stub`: fake transport、fixed provider response、request spy。
- `related detail requirement type`: external_integration、adapter 境界、network 境界、testability。
- `adoption hint`: 有料 real API を前提にしない acceptance 条件の下支え候補になる。
- `conflict hint`: fake provider を fixture 名として UI に出す案は source requirement と衝突する。

### CAND-TJSPPS-006 Gemini / xAI だけ batch capability を有効候補として検証する

- `source requirement`: `plan.md:63-64` は batch mode を暗黙推定にせず、対象 provider を Gemini と xAI だけに限定し、checkbox または select で明示することを要求している。`plan.md:100` は batch mode の対象 provider を designer の必須観点にしている。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-006`
- `external boundary`: provider capability と validation の境界。
- `actor`: Job Setup 利用者。
- `trigger`: Gemini、xAI、その他 provider で batch mode を切り替え、validation を実行する。
- `start condition`: provider / model / execution mode が phase ごとに選択済みである。
- `expected outcome`: Gemini と xAI では batch mode が明示選択できる。その他 provider で batch mode が選ばれた場合は provider / mode 不整合の blocking validation failure になる。
- `observable point`: batch control state、provider capability result、validation failure category、phase 別 execution mode summary。
- `fake_or_stub`: provider capability fixture、fake transport、validation fixture。
- `related detail requirement type`: external_integration、provider 境界、validation 境界。
- `adoption hint`: batch capability の対象 provider と validation 境界を固定する候補になる。
- `conflict hint`: UI 観点では control 表示条件、failure 観点では unsupported provider / mode と統合される可能性がある。

### CAND-TJSPPS-007 model list 取得失敗は redacted な retry 状態として返す

- `source requirement`: `plan.md:100` は model list 取得失敗時の UI 状態を designer の必須観点にしている。`ui-design.md:82-84` は provider / mode 不整合、credential 不備、create 失敗を区別し、secret 値を error message に含めないことを要求している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-007`
- `external boundary`: provider getModels の network failure と UI 状態返却の境界。
- `actor`: Job Setup 利用者。
- `trigger`: secret 設定済み provider で model list API が timeout、network failure、invalid status のいずれかになる。
- `start condition`: API key が必要な provider では secret store 参照が設定済みである。LM Studio では secret store 参照なしで provider endpoint が参照される。
- `expected outcome`: model 候補は成功扱いにならない。failure reason は redacted に返り、retry 可能性と provider / phase 対応を確認できる。
- `observable point`: model list error state、retry action availability、request log、secret redaction、provider raw response 非表示。
- `fake_or_stub`: timeout fake transport、invalid status fixture、request spy。
- `related detail requirement type`: external_integration、network 境界、display。
- `adoption hint`: model 候補取得失敗時の画面状態と再試行可否を扱う候補になる。
- `conflict hint`: validation failure と model list failure のどちらで blocking にするかは designer 側の統合判断に残す。

### CAND-TJSPPS-008 provider 固有 model response を内部 model 候補へ変換する

- `source requirement`: `plan.md:61` は provider ごとの model list API から model 候補を取得することを要求している。`scenario-design.md:299-326` は fixed provider response と fake transport で外部境界を検証できることを示している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-008`
- `external boundary`: provider response adapter と内部 model candidate contract の境界。
- `actor`: Job Setup 利用者。
- `trigger`: provider ごとの model list response を受け取り、Job Setup の model 候補へ変換する。
- `start condition`: provider response は fake transport または real transport から返る。
- `expected outcome`: provider 固有 response は内部 model 候補に正規化される。空 response、重複 model、必須項目欠落、想定外 schema は成功扱いにならない。
- `observable point`: normalized model candidate list、adapter error kind、raw response 非露出、provider / model 名の表示。
- `fake_or_stub`: fixed provider responses、invalid schema fixture、adapter-level fake。
- `related detail requirement type`: external_integration、adapter 境界、provider 境界。
- `adoption hint`: provider 追加時の adapter 境界を acceptance 候補として残せる。
- `conflict hint`: product 実装方式の固定に踏み込まないよう、具体 class 名や SDK 名は designer 後段へ送る。

### CAND-TJSPPS-009 phase detail の provider summary は secret と raw payload を露出しない

- `source requirement`: `term-translation-phase.md:57-59`、`persona-generation-phase.md:64-70`、`body-translation-phase.md:42-43` は secret、API key 平文、provider raw request / response、raw prompt の非露出を要求している。`ui-design.md:99-105` は API key 平文が画面、console、error summary に出ないことを確認点にしている。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSPPS-009`
- `external boundary`: Job Setup の phase runtime summary と後続 phase provider summary の redaction 境界。
- `actor`: Job Setup 利用者、運用確認者。
- `trigger`: phase 別 provider / model / execution mode / credential 参照を保存または validation summary に表示する。
- `start condition`: phase 別 runtime 設定が選択済みで、secret store 参照がある provider と API key なし provider が混在している。
- `expected outcome`: 表示と summary には provider、model、execution mode、credential 参照状態だけが出る。secret、復号可能値、provider raw request / response、raw prompt は出ない。
- `observable point`: UI summary、validation summary、structured log、fake transport log、phase result summary。
- `fake_or_stub`: fake secret store、redaction assertion fixture、fake transport log。
- `related detail requirement type`: external_integration、secret 境界、operation-audit 競合候補。
- `adoption hint`: Job Setup と各 phase detail-spec の redaction 整合を確認する候補になる。
- `conflict hint`: operation-audit 観点の監査要約候補と重複しうるため、外部連携側は secret / raw payload 非露出に絞る。

## Open Notes

- `human decision candidate`: LM Studio の接続先、到達性確認の timeout、retry 回数、base URL の入力場所は、この候補生成では確定しない。
- `human decision candidate`: Gemini / xAI の batch capability の参照元を固定するか、provider adapter fixture だけで判定するかは designer に残す。
- `merge candidate`: `CAND-TJSPPS-001` と `CAND-TJSPPS-002` は getModels の secret gate として統合可能である。
- `merge candidate`: `CAND-TJSPPS-006` は validation failure 観点の unsupported provider / mode と統合可能である。
- `rejection candidate`: fake provider を user-facing provider list に出す案は、既存 scenario design の固定要件と衝突するため不採用候補である。
- `conflict candidate`: model list failure を create 前 validation の blocking failure に含めるか、model 選択前の retry state に留めるかは最終 scenario 統合で判断する。
