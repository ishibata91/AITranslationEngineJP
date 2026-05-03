# Scenario Candidates: translation-job-setup-phase-provider-settings / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSPPS`

## Generator Scope

- `viewpoint`: 失敗。入力不備、参照不能、設定不整合、回復動作、secret 非露出を候補化する。
- `included_sources`:
  - `docs/exec-plans/active/translation-job-setup-phase-provider-settings/plan.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
  - `docs/exec-plans/completed/translation-job-setup/ui-design.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本変更、最終シナリオ表、採否、統合判断、他 generator 起動。
- `generation_notes`: 必須観点のうち、API key 未設定 provider、model list API 失敗、LM Studio、batch mode、master-persona 依存残存をそれぞれ候補として分離する。

## Candidate Scenarios

### CAND-TJSPPS-001 API key 未設定 provider の model list 取得を外部へ出さない

- `source requirement`: `plan.md:57-63`, `plan.md:91-100`, `scenario-design.md:24-25`, `ui-design.md:32-34`
- `viewpoint`: 参照不能、secret 非露出、外部連携抑止
- `candidate scenario id`: `CAND-TJSPPS-001`
- `actor`: Job Setup で phase 別 provider を選ぶ利用者
- `failure start condition`: API key が必要な provider を phase に選択し、対象 credential が未設定である。
- `rejected operation`: model list API を外部 provider へ送る操作。
- `trigger`: model selector を開く、または model list 再取得を実行する。
- `expected error`: `credential missing` 相当の blocking reason を表示し、API key 平文、secret 本体、復号可能値を表示しない。
- `expected outcome`: 外部 request は 0 件であり、create job は validation pass にならない。利用者は credential 設定導線または手動 model 選択可否の判断待ち状態へ戻れる。
- `observable point`: Job Setup runtime panel、validation summary、fake transport request log、secret redaction。
- `related detail requirement type`: external integration、display、trust boundary、validation failure。
- `adoption hint`: API key 未設定 provider を model list 取得前に止める候補として扱う。
- `conflict hint`: 手動 model 選択を API key 未設定時にも許可するかは designer の統合判断に残す。

### CAND-TJSPPS-002 model list API 失敗時に secret を出さず復旧操作へ戻せる

- `source requirement`: `plan.md:57-63`, `plan.md:91-100`, `scenario-design.md:62-69`, `ui-design.md:82-84`, `term-translation-phase.md:57-59`, `persona-generation-phase.md:64-70`, `body-translation-phase.md:42-43`
- `viewpoint`: 外部応答失敗、secret 非露出、回復動作
- `candidate scenario id`: `CAND-TJSPPS-002`
- `actor`: Job Setup で model 候補を取得する利用者
- `failure start condition`: credential は設定済みだが、model list API が network error、provider error、invalid response のいずれかで失敗する。
- `rejected operation`: 失敗した provider response、request header、API key、secret 値を UI、error summary、structured log、fake transport log へ出す操作。
- `trigger`: phase の provider 選択後に model list 取得を実行する。
- `expected error`: `model list fetch failed` 相当の短い理由を表示する。secret、API key 平文、provider raw request / response は表示しない。
- `expected outcome`: 利用者は model の手動選択または再取得へ戻れる。失敗状態のまま create job は validation pass にならない。
- `observable point`: model selector error state、retry action、manual selection state、request log、structured log、validation summary。
- `related detail requirement type`: external integration、display、redaction、retry。
- `adoption hint`: model list API の失敗と provider validation の失敗を混同しない候補として扱う。
- `conflict hint`: 既存 translation-job-setup は network reachability を blocking validation failure としているため、手動 model 選択後にどの段階で blocking に戻すかは designer が統合する。

### CAND-TJSPPS-003 LM Studio に API key 不足 failure を出さない

- `source requirement`: `plan.md:57-63`, `plan.md:91-100`, `ui-design.md:30-34`, `persona-generation-phase.md:31-33`, `body-translation-phase.md:27-32`
- `viewpoint`: 設定不整合、誤判定防止
- `candidate scenario id`: `CAND-TJSPPS-003`
- `actor`: LM Studio を phase provider に選ぶ利用者
- `failure start condition`: LM Studio が選択され、credential 参照が空である。
- `rejected operation`: API key 入力、API key 未設定 warning、credential missing validation failure を LM Studio に出す操作。
- `trigger`: LM Studio の model selector を開く、または Job Setup validation を実行する。
- `expected error`: API key 不足 error は出ない。LM Studio endpoint 不通、model list 取得失敗、model 未選択など別カテゴリの failure だけを区別して表示する。
- `expected outcome`: credential 欠落を理由に create job を止めない。別の blocking failure がある場合だけ、別カテゴリとして validation failure にする。
- `observable point`: runtime panel、credential field visibility、validation summary、error kind、fake transport request log。
- `related detail requirement type`: provider capability、display、validation failure。
- `adoption hint`: provider capability によって credential 必須性が変わる候補として扱う。
- `conflict hint`: 既存 translation-job-setup の credential missing blocking 条件を provider capability で例外化する必要がある。

### CAND-TJSPPS-004 Gemini / xAI 以外の batch mode 指定を validation failure にする

- `source requirement`: `plan.md:57-64`, `plan.md:91-100`, `scenario-design.md:62-69`, `scenario-design.md:157-178`, `ui-design.md:30-32`, `body-translation-phase.md:27-32`
- `viewpoint`: 設定不整合、provider capability
- `candidate scenario id`: `CAND-TJSPPS-004`
- `actor`: phase の execution mode を設定する利用者
- `failure start condition`: provider が Gemini または xAI 以外であり、batch mode が指定されている。
- `rejected operation`: provider capability に反する batch mode のまま validation pass または create job を許可する操作。
- `trigger`: batch mode を選択した状態で validation を実行する。
- `expected error`: `unsupported provider / mode` 相当の blocking validation failure を表示する。
- `expected outcome`: create job は無効または拒否される。利用者は execution mode を通常 mode へ戻して再 validation できる。
- `observable point`: phase runtime setting、validation summary、create button state、external request 未実行証跡。
- `related detail requirement type`: validation failure、provider capability、display。
- `adoption hint`: batch mode を暗黙推定にしない条件の失敗候補として扱う。
- `conflict hint`: Gemini と xAI の batch capability の詳細範囲は designer が最終統合で固定する。

### CAND-TJSPPS-005 master-persona provider 設定への依存残存を失敗として検出する

- `source requirement`: `plan.md:8-10`, `plan.md:51-64`, `scenario-design.md:24-25`, `scenario-design.md:217-239`, `persona-generation-phase.md:31-33`, `body-translation-phase.md:27-32`
- `viewpoint`: 参照境界違反、設定不整合、回復不能な stale dependency
- `candidate scenario id`: `CAND-TJSPPS-005`
- `actor`: Job Setup で phase 別 provider 設定を保存または検証する利用者
- `failure start condition`: master-persona provider / model / secret 参照が保存済みであり、Job Setup の phase 別設定とは異なる。
- `rejected operation`: master-persona の provider、model、secret key を Job Setup の既定値、保存元、validation 対象、phase 実行設定として使う操作。
- `trigger`: Job Setup を開く、phase 別 validation を実行する、または create job を実行する。
- `expected error`: `phase runtime setting missing` または `invalid phase runtime source` 相当の blocking validation failure を表示する。
- `expected outcome`: phase 別 provider 設定が未設定なら create job はできない。master-persona 設定を代替として使った validation pass は許可しない。
- `observable point`: runtime option source、validation target snapshot、secret key namespace、create job payload summary、phase execution summary。
- `related detail requirement type`: responsibility boundary、external integration、validation failure、trust boundary。
- `adoption hint`: master-persona 依存の残存を検出する回帰候補として扱う。
- `conflict hint`: 既存 translation-job-setup の保存済み AI 設定復元シナリオと、今回の phase 別設定独立要件を designer が分離する必要がある。

### CAND-TJSPPS-006 phase 別 provider 設定の不足を phase 単位で validation failure にする

- `source requirement`: `plan.md:57-64`, `scenario-design.md:57-83`, `term-translation-phase.md:21-24`, `persona-generation-phase.md:21-33`, `body-translation-phase.md:20-32`
- `viewpoint`: 失敗入力、参照不能、phase runtime completeness
- `candidate scenario id`: `CAND-TJSPPS-006`
- `actor`: 3 phase の runtime 設定を作る利用者
- `failure start condition`: 単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかで provider、model、credential 参照、execution mode の必要項目が欠けている。
- `rejected operation`: 欠損 phase を既定値や他 phase の値で補完し、validation pass または create job を許可する操作。
- `trigger`: 欠損を含む phase 別設定で validation を実行する。
- `expected error`: 欠損 phase 名と欠損項目を含む blocking validation failure を表示する。secret 値は表示しない。
- `expected outcome`: create job は無効または拒否される。利用者は該当 phase の設定だけを修正して再 validation できる。
- `observable point`: phase runtime panel、validation summary、create button state、validation target snapshot。
- `related detail requirement type`: validation failure、display、phase execution setting。
- `adoption hint`: 3 phase の独立設定をまとめて validation する候補として扱う。
- `conflict hint`: credential が不要な LM Studio は CAND-TJSPPS-003 と統合時に例外化する。

### CAND-TJSPPS-007 model list 取得後の provider 変更で stale model を validation failure にする

- `source requirement`: `plan.md:57-64`, `plan.md:91-100`, `scenario-design.md:151-179`, `ui-design.md:22-33`, `body-translation-phase.md:27-32`
- `viewpoint`: 設定不整合、validation stale、回復動作
- `candidate scenario id`: `CAND-TJSPPS-007`
- `actor`: model list 取得後に provider または execution mode を変更する利用者
- `failure start condition`: provider A で取得した model を保持したまま、同じ phase の provider を provider B へ変更する。
- `rejected operation`: provider A の model を provider B の model として validation pass または create job に使う操作。
- `trigger`: provider 変更後に validation を実行する、または create job を押す。
- `expected error`: `stale model selection` または `provider / model mismatch` 相当の blocking validation failure を表示する。
- `expected outcome`: create job は無効になる。利用者は model list 再取得または手動 model 選択で対象 phase を更新できる。
- `observable point`: dirty-validation state、model selector source provider、validation summary、create button state。
- `related detail requirement type`: validation stale、provider capability、display。
- `adoption hint`: validation stale 条件を phase 別 provider 設定に適用する候補として扱う。
- `conflict hint`: 手動 model 名の provider 所属確認をどこまで行うかは designer が統合する。

### CAND-TJSPPS-008 credential 参照失効後の model list と validation を secret 非露出で止める

- `source requirement`: `plan.md:57-64`, `scenario-design.md:57-69`, `scenario-design.md:211-239`, `ui-design.md:28-34`, `term-translation-phase.md:57-59`, `persona-generation-phase.md:64-70`, `body-translation-phase.md:42-43`
- `viewpoint`: 参照不能、secret 非露出、回復動作
- `candidate scenario id`: `CAND-TJSPPS-008`
- `actor`: 保存済み credential 参照を使う利用者
- `failure start condition`: model list 取得後、または validation 前に credential 参照が削除、失効、解決不能になっている。
- `rejected operation`: stale credential を使って外部 request を送る操作、または secret 解決失敗の詳細値を表示する操作。
- `trigger`: model list 再取得、validation、create job のいずれかを実行する。
- `expected error`: `credential reference unavailable` 相当の blocking validation failure を表示する。secret、API key 平文、復号可能値は表示しない。
- `expected outcome`: 外部 request は送られない。利用者は credential を再選択し、model list 再取得または validation 再実行へ戻れる。
- `observable point`: credential ref state、request log、validation summary、structured log、create button state。
- `related detail requirement type`: credential resolution、trust boundary、validation failure、retry。
- `adoption hint`: model list と create validation の間で credential が失効する候補として扱う。
- `conflict hint`: LM Studio は credential 不要であるため、provider capability に基づく除外が必要である。

## Open Notes

- `human decision candidate`: API key 未設定 provider でも手動 model 選択を許可するか。
- `human decision candidate`: model list API 失敗後の手動 model 選択を validation pass まで進める条件。
- `human decision candidate`: LM Studio の model list 取得先が未設定または到達不能の時の failure category。
- `merge candidate`: CAND-TJSPPS-001、CAND-TJSPPS-002、CAND-TJSPPS-008 は secret 非露出と外部 request 抑止で統合候補になる。
- `merge candidate`: CAND-TJSPPS-004、CAND-TJSPPS-006、CAND-TJSPPS-007 は phase runtime validation の統合候補になる。
- `rejection candidate`: 正常系の provider / model 選択成功、create job 成功、phase 実行開始成功は failure generator の対象外である。
