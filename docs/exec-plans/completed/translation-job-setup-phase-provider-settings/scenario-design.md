# Scenario Design: translation-job-setup-phase-provider-settings

- `skill`: scenario-design
- `status`: pending-human-review
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-job-setup-phase-provider-settings.md`
- `topic_abbrev`: `TJSPPS`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - Job Setup は master-persona の provider / model / credential / execution mode を既定値、保存元、validation 対象、secret 解決元として使わない。
  - Job Setup は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 phase ごとに provider、model、credential 参照、execution mode、batch mode を保持する。
  - phase 実行時は対象 phase の Job Setup 保存設定だけを参照し、他 phase と master-persona の provider 設定を参照しない。
  - model 候補は provider 別 getModels 系 API から取得する。API key 必須 provider は credential 参照が解決できる場合だけ外部取得する。
  - model list API が失敗した場合、または API key 必須 provider で API key が未設定の場合、手動 model 入力は許可しない。
  - API key 未設定、credential 参照不能、credential 失効では外部 getModels と provider validation request を送らない。
  - LM Studio は API key 不要 provider として扱い、API key 入力、credential select、API key 未設定 warning、credential missing failure を出さない。
  - batch mode は Gemini と xAI だけで選択可能にし、他 provider では選択不能または blocking validation failure にする。
  - batch mode 操作は checkbox として固定する。理由は有効 / 無効の 2 値であり、UI 設計規約上の binary control に該当するためである。
  - 翻訳段階ごとの明示的な確認ボタンは必須にしない。3 phase すべてで API key 必須 provider の API key が設定済みで、モデル一覧から model が選ばれている場合、create job を実行できる。
  - API key 平文、secret 本体、復号可能値、provider raw request / response、raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
- `non_goals`:
  - product code、product test、docs 正本、implementation-scope は扱わない。
  - master-persona provider 設定の保存 UI と secret 保存 UI は扱わない。
  - LM Studio の base URL 設定画面、timeout、retry 回数の恒久仕様は扱わない。
  - 後続 phase の翻訳実行、pause、resume、retry、cancel の操作 UI は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-*` を一意 key として扱う。

`needs_human_decision` は 0 件である。
未解決 conflict は 0 件である。
人間回答 `Q-TJSPPS-001` により、手動 model 入力は許可しない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは sidecar JSON に分離する。
人間回答 `Q-TJSPPS-001` により、詳細要求タイプの人間判断待ちは残っていない。

### `REQ-TJSPPS-001` phase 別 provider 設定を作成する

- `source_requirement`: Job Setup は単語翻訳、NPC ペルソナ生成、本文翻訳の各 phase に独立した provider / model / credential / execution mode / batch mode を持つ。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: Draft は UI 状態として作り、create job 時に 3 phase の設定断面を保存する。master-persona provider 設定は初期値にも fallback にも使わない。

### `REQ-TJSPPS-002` model 候補を provider 別 API から取得する

- `source_requirement`: model 候補は provider 別 getModels 系 API から取得し、API key が設定済みの場合だけ外部取得する。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: API key 未設定または credential 参照不能では外部取得しない。取得失敗は secret 非露出の retry state として扱う。手動 model 入力欄は表示しない。

### `REQ-TJSPPS-003` LM Studio を credential 不要 provider として扱う

- `source_requirement`: LM Studio は API key を要求しないため、API key 入力、credential select、API key 未設定 warning を出さない。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: LM Studio の credential 状態は not_required とし、credential missing は validation failure にしない。endpoint 不通や model list 失敗は別カテゴリで表示する。

### `REQ-TJSPPS-004` batch mode を Gemini / xAI だけに限定する

- `source_requirement`: batch mode は暗黙推定にせず、Gemini と xAI だけで明示的に切り替える。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: batch mode 操作は checkbox として固定する。対象外 provider では batch checkbox を非表示または disabled にし、stale batch 値を保存対象にしない。

### `REQ-TJSPPS-005` create 可否と遅延 model list 結果の扱いを固定する

- `source_requirement`: 翻訳段階ごとの明示的な確認ボタンは不要であり、未設定がなければ create job を実行できる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: create job は 3 phase すべてで API key 必須 provider の API key が設定済みで、モデル一覧から model が選ばれている場合だけ許可する。遅延 model list 結果は現在の provider / APIキー状態へ混入させない。

### `REQ-TJSPPS-006` secret 非露出と監査要約を固定する

- `source_requirement`: credential 参照状態、provider、model、execution mode、batch mode は確認可能にし、API key 平文や raw payload は出さない。
- `requirement_kind`: security
- `needs_human_decision`: なし
- `fixed_decisions`: credential は存在状態または参照状態だけを表示する。business state に model list 取得履歴の全履歴は持たず、直近 UI 状態と structured log だけを観測点にする。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答質問は 0 件である。
`Q-TJSPPS-001` の承認記録により、人間レビュー待ちの design artifact へ進められる。

## Risks

- `implementation_risks`:
  - 既存 Job Setup は単一 runtime option 前提であるため、3 phase 分の設定断面へ public seam を拡張する必要がある。
  - model list API と validation が別タイミングで動くため、遅延 response が現在の provider / credential へ混入しない状態管理が必要である。
  - master-persona の secret namespace を再利用すると、今回の非依存要件に違反する。
  - 手動 model 入力を許可しないため、取得失敗時の retry 導線と credential 設定導線が弱いと利用者が復旧できない。
- `test_data_risks`:
  - paid real API を使わず、fake secret store と fake transport で getModels、validation、redaction を観測する必要がある。
  - Gemini、xAI、LM Studio、API key 必須 provider の credential missing を fixture で分ける必要がある。

## Rules

- ケース ID は `SCN-TJSPPS-NNN` 形式にする。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- paid な real AI API を前提にしない。
- 人間判断待ちと未解決競合は残さない。

## Scenario Matrix

人間回答 `Q-TJSPPS-001` により、model selector は取得済み候補からの選択だけを許可する。
model list API 失敗時と API key 未設定時は、手動 model 入力欄を表示しない。

### SCN-TJSPPS-001 phase 別 provider 設定で Ready job を作成する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Job Setup で 3 phase の provider / model / credential / execution mode / batch mode を個別に設定する。
- `受け入れ条件`: 3 phase すべてで API key 必須 provider の API key が設定済みで、モデル一覧から model が選ばれている時に create job を実行できる。
- `入力開始点`: Job Setup UI。
- `主要 outcome`: `Ready` job と phase 別 runtime 要約が作成される。
- `主要操作列`: Job Setup を開き、3 phase の runtime 設定を選び、未設定がないことを確認して create job を実行する。
- `期待結果`:
  1. 単語翻訳、NPC ペルソナ生成、本文翻訳の runtime 設定が別々に表示される。
  2. create result と read-only 要約に phase 別 provider、model、execution mode、batch mode、credential 参照状態が出る。
  3. API key 平文は表示されない。
- `観測点`: Job Setup UI、create response、read-only job summary。
- `fake_or_stub`: fixed input fixture、fake secret store、fake transport。

### SCN-TJSPPS-002 master-persona provider 設定を参照しない

- `分類`: 責務境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: Job Setup の runtime 設定が master-persona provider 設定から独立していることを確認する。
- `受け入れ条件`: master-persona 側に provider / model / secret が存在しても、Job Setup の未設定 phase は作成可能にならない。
- `入力開始点`: master-persona provider 設定あり、Job Setup phase 設定なしの fixture。
- `主要 outcome`: master-persona 設定を既定値、保存元、secret 解決元として使わない。
- `主要操作列`: Job Setup options 取得、validation、create を試行する。
- `期待結果`:
  1. phase 別設定が空の場合は phase runtime missing になる。
  2. `master-persona:<provider>` 相当の secret key 参照では作成可能にならない。
  3. 作成後要約の setting source は Job Setup phase 設定である。
- `観測点`: validation target snapshot、secret key namespace、create payload summary。
- `fake_or_stub`: master-persona settings fixture、fake secret store。

### SCN-TJSPPS-003 provider 別 getModels を secret gate 後だけ実行する

- `分類`: 外部連携
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: API key 必須 provider の model 候補取得を credential 解決後だけ実行する。
- `受け入れ条件`: credential 参照が解決できる場合だけ getModels を呼び、credential missing または失効では外部 request を送らない。
- `入力開始点`: credential 有無を切り替えられる phase runtime fixture。
- `主要 outcome`: model list の loading、success、failure、skipped を phase 別に観測できる。
- `主要操作列`: provider 選択、model list 取得、credential 失効、再取得を実行する。
- `期待結果`:
  1. credential ありでは fake getModels response から model 候補が表示される。
  2. credential missing では外部 request が 0 件で、credential missing が表示される。
  3. getModels 失敗では secret 非露出の retry state になる。
  4. credential missing または getModels 失敗では手動 model 入力欄が表示されず、create job は実行できない。
- `観測点`: model selector、request spy、validation summary、redacted log。
- `fake_or_stub`: fake secret store、phase-tagged fake transport、fixed model response、failure response。

### SCN-TJSPPS-004 LM Studio を API key 不要 provider として設定する

- `分類`: 代替正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: LM Studio 選択時に credential UI と API key warning を出さない。
- `受け入れ条件`: LM Studio の credential 状態は not_required であり、credential missing は validation failure にならない。
- `入力開始点`: Job Setup UI。
- `主要 outcome`: LM Studio の provider、model、通常 execution mode を設定できる。
- `主要操作列`: 任意の phase で LM Studio を選び、model 候補状態と validation summary を確認する。
- `期待結果`:
  1. API key 入力、credential select、API key 未設定 warning は表示されない。
  2. validation summary に credential missing は出ない。
  3. endpoint 不通や model list failure は API key missing とは別カテゴリで表示される。
- `観測点`: runtime panel、credential field visibility、validation summary、request spy。
- `fake_or_stub`: fake LM Studio endpoint、local model response、endpoint failure fixture。

### SCN-TJSPPS-005 Gemini / xAI だけ batch mode checkbox を使える

- `分類`: 境界条件
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: batch mode を Gemini / xAI に限定し、利用者が checkbox で明示する。
- `受け入れ条件`: Gemini / xAI では batch checkbox を表示する。その他 provider では batch checkbox を非表示または disabled にし、stale batch 値を保存対象にしない。
- `入力開始点`: Job Setup UI。
- `主要 outcome`: execution mode は暗黙推定ではなく phase 別の明示値として保存される。
- `主要操作列`: Gemini、xAI、LM Studio を切り替え、batch checkbox と validation summary を確認する。
- `期待結果`:
  1. Gemini / xAI では checkbox によって batch mode を on / off できる。
  2. provider 変更後は対象 provider に合わせて batch checkbox の表示が更新される。
  3. 対象外 provider では batch mode が選択できない。
- `観測点`: batch checkbox、create button state、phase runtime summary。
- `fake_or_stub`: provider capability fixture、fake transport。

### SCN-TJSPPS-006 設定変更と遅延結果で model 選択を保護する

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 古い model list 結果で、現在の AIサービスと APIキー状態に合わない model が選択済みにならないことを確認する。
- `受け入れ条件`: provider または APIキー状態が変わると model 選択は未選択に戻り、現在の model list 取得が成功するまで create job は有効にならない。
- `入力開始点`: model 選択済み Draft。
- `主要 outcome`: 不足理由として対象 phase の model 未選択または APIキー未設定が見える。
- `主要操作列`: model 選択後に provider または APIキー状態を変更し、変更前の遅延 getModels result を返す。
- `期待結果`:
  1. create job は現在の model が未選択の間 disabled になる。
  2. 遅延 getModels result は現在 phase の model 候補へ混入しない。
  3. 作成不可理由は、該当 phase の APIキー未設定または model 未選択だけを示す。
- `観測点`: model selector、model list source provider、create button state、作成不可理由。
- `fake_or_stub`: delayed model response fixture、fake transport。

### SCN-TJSPPS-007 phase 実行時に対象 phase 専用設定を参照する

- `分類`: 責務境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: phase run が Job Setup で確定した対象 phase 専用設定を読む。
- `受け入れ条件`: 単語翻訳、NPC ペルソナ生成、本文翻訳はそれぞれの provider / model / credential / execution mode を参照し、他 phase と master-persona 設定を参照しない。
- `入力開始点`: phase 別に異なる runtime 設定を持つ Ready job fixture。
- `主要 outcome`: phase run summary に対象 phase の設定だけが表示される。
- `主要操作列`: 各 phase 開始境界を呼び、provider summary を確認する。
- `期待結果`:
  1. 単語翻訳 phase は単語翻訳用設定を読む。
  2. NPC ペルソナ生成 phase は persona 用設定を読む。
  3. 本文翻訳 phase は本文翻訳用設定を読む。
  4. 開始時に provider 再選択 UI は出ない。
- `観測点`: phase run summary、credential 参照状態、provider request unit。
- `fake_or_stub`: Ready job fixture、fake transport。

### SCN-TJSPPS-008 secret 非露出と paid API 非依存を確認する

- `分類`: セキュリティ / テスト容易性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: provider list は real provider のまま、外部 request だけ fake に差し替えて検証する。
- `受け入れ条件`: API key 平文、secret 本体、provider raw request / response、raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出ない。
- `入力開始点`: Gemini、xAI、LM Studio を含む provider fixture。
- `主要 outcome`: redacted summary と request 未実行証跡を確認できる。
- `主要操作列`: model list、validation、create summary を fake transport で実行する。
- `期待結果`:
  1. user-facing provider list に fake provider は出ない。
  2. paid real API は呼ばれない。
  3. credential 参照状態、provider、model、execution mode、batch mode だけが要約に出る。
- `観測点`: provider list、request spy、redaction assertion、structured log。
- `fake_or_stub`: fake transport、redaction assertion fixture。

## Acceptance Checks

- `REQ-TJSPPS-001`: `SCN-TJSPPS-001`, `SCN-TJSPPS-002`, `SCN-TJSPPS-007`
- `REQ-TJSPPS-002`: `SCN-TJSPPS-003`, `SCN-TJSPPS-006`, `SCN-TJSPPS-008`
- `REQ-TJSPPS-003`: `SCN-TJSPPS-004`, `SCN-TJSPPS-008`
- `REQ-TJSPPS-004`: `SCN-TJSPPS-005`, `SCN-TJSPPS-006`
- `REQ-TJSPPS-005`: `SCN-TJSPPS-001`, `SCN-TJSPPS-006`
- `REQ-TJSPPS-006`: `SCN-TJSPPS-002`, `SCN-TJSPPS-003`, `SCN-TJSPPS-008`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup-phase-provider-settings/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Questions

- none
