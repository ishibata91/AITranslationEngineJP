# Scenario Design: translation-job-setup

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-job-setup.md`
- `topic_abbrev`: `TJS`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - `translation-input-intake` 完了後の 1 入力データから、1 件の翻訳ジョブを作成できる。
  - create job 前に、入力データ、共通辞書、共通ペルソナ、AI runtime、実行方式の validation 結果を確認できる。
  - validation pass なしで `Ready` job を作成しない。
  - 作成済み job は 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照し、入力出自を失わない。
  - AI 基盤設定は provider、model、credential 参照、実行方式を区別し、API key 平文を表示または保存要約に出さない。
  - paid な real AI API を scenario validation の前提にしない。
- `non_goals`:
  - 共通辞書と共通ペルソナの管理 UI は含めない。
  - 翻訳 phase 実行、訳文生成、`JOB_TRANSLATION_FIELD` の訳文更新、成果物出力は含めない。
  - job cancel、pause、resume、phase retry の実行操作は後続 job management / phase task へ送る。
  - docs 正本、product code、product test、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-TJS-NNN` を一意 key として扱う。

`needs_human_decision` は 0 件である。
未解決 conflict は 0 件である。
全 candidate は採用、統合、不採用、解決済み conflict のいずれかへ分類済みである。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは sidecar JSON に分離する。
質問票回答と human review 承認を反映済みである。

### `REQ-TJS-001` 1 入力データから 1 翻訳ジョブを作成する

- `source_requirement`: 1 入力データに対して 1 翻訳ジョブを作成し、実行前 validation を完了する。
- `requirement_kind`: operation
- `needs_human_decision`: なし
- `fixed_decisions`: Draft は UI 未保存状態である。DB job は validation pass 後の create で初めて作る。同一入力への 2 件目 job 作成は禁止するが、過去 job を廃棄できる手段は別途必要である。

### `REQ-TJS-002` create 前 validation を通す

- `source_requirement`: create job の前に validation を通せる。validation failure の理由を確認できる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 必須設定不足、参照不能、provider / mode 不整合、credential 参照不能は blocking validation failure にする。validation 実行履歴は business table ではなく structured log に残し、アプリ状態は直近結果、対象設定断面、失敗カテゴリ、job 作成時の pass 断面だけ保持する。

### `REQ-TJS-003` AI 基盤と実行方式を選ぶ

- `source_requirement`: 基盤参照、AI 基盤、実行方式を選択できる。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: credential 解決、provider capability、ネットワーク到達性はすべて blocking にする。provider list は real provider を扱い、fake provider を user-facing 選択肢にしない。test では外部 request / SDK transport だけを fake に差し替える。

### `REQ-TJS-004` job setup の状態と永続化を壊さない

- `source_requirement`: `TRANSLATION_JOB` は 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照し、ジョブ状態は `JOB_PHASE_RUN` 群から集約する。
- `requirement_kind`: persistence
- `needs_human_decision`: なし
- `fixed_decisions`: cache 欠落時は Job Setup をブロックし、Input Review の再構築導線へ戻す。Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集は許可しない。

### `REQ-TJS-005` Job Setup UI で確認と修正を完結する

- `source_requirement`: Job Setup を開き、入力データ、基盤参照、AI runtime を選び、validation pass 状態を確認して翻訳ジョブを作成する。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: UI は未保存 Draft、validation 結果、create 可否、作成後の read-only 要約を表示する。見た目 mock や product component 実装には踏み込まない。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

質問票回答は `Q-001` から `Q-008` まで反映済みである。
未回答質問はない。

## Risks

- `implementation_risks`:
  - 同一入力への 2 件目 job 作成は禁止するため、過去 job を廃棄する別手段がないとやり直し作業が詰まる。
  - provider network reachability を blocking にするため、外部要因で setup が止まる可能性がある。
  - 共通基盤の lock は phase 実行時へ deferred したため、phase 側 design で参照中更新を扱う必要がある。
  - Draft は UI 未保存状態のため、長い setup 作業を中断再開する UX は別途扱う必要がある。
- `test_data_risks`:
  - paid な real AI API を使わず、provider capability と validation 結果を fixture / fake transport で観測する必要がある。
  - terminal job、stale foundation ref、cache missing、partial create failure は fixture を分ける必要がある。
  - validation 実行履歴は structured log で検証し、business table の履歴として期待しない。

## Rules

- ケース ID は `SCN-TJS-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | final validation` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が残る場合は scenario 完了にしない。
- 未解決 conflict が残る場合は scenario 完了にしない。
- `not_applicable` と `deferred` は理由なしで通さない。
- paid な real AI API を前提にしない。

## Scenario Matrix

質問票回答を反映済みである。
この scenario matrix は human review 承認済みである。

### SCN-TJS-001 validation pass 済み setup から Ready job を作成する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 1 入力データに対して validation pass 後に 1 件の job を作成する。
- `受け入れ条件`: 入力データ、基盤参照、AI runtime、実行方式が選択され、validation pass が表示された後だけ create できる。
- `事前条件`: `translation-input-intake` が完了し、取り込み済み入力データが 1 件以上ある。
- `public_seam_or_api_boundary`: Job Setup の create job boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。1 input = 1 job、input 出自保持、実行設定保持。
- `入力開始点`: Job Setup UI。
- `主要 outcome`: 作成済み job が `Ready` として観測できる。
- `開始操作`: Job Setup を開く。
- `入力方法`: UI で入力データ、共通辞書、共通ペルソナ、AI runtime、実行方式を選ぶ。
- `主要操作列`: validation を実行し、pass 表示を確認し、create job を実行する。
- `手順`:
  1. Job Setup で取り込み済み入力データを選ぶ。
  2. 共通辞書、共通ペルソナ、AI runtime、実行方式を選ぶ。
  3. validation pass を確認して create job を実行する。
- `期待結果`:
  1. 選択入力に紐づく job が 1 件作成される。
  2. 作成後の表示で input data ID、入力出自、execution setting、validation 結果を確認できる。
  3. API key 平文は表示されない。
- `観測点`: Job Setup UI、job detail / job list、永続化済み input data 参照。
- `UI-visible outcome`: 作成完了、`Ready`、入力出自、設定要約が表示される。
- `fake_or_stub`: fixed input fixture、foundation data fixture、fake transport。
- `責務境界メモ`: Draft は UI 未保存状態で、DB job は create 成功時に初めて作る。Ready job は再表示だけ許可し、再編集は許可しない。

### SCN-TJS-002 validation failure の理由を確認して create を止める

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: create 前 validation failure を UI で確認し、無効な job 作成を防ぐ。
- `受け入れ条件`: 必須設定不足、参照不能、provider / mode 不整合、credential 参照不能、cache 欠落がある場合は create job ができない。
- `事前条件`: 必須設定不足、不整合 runtime、参照不能 foundation のいずれかを作れる fixture がある。
- `public_seam_or_api_boundary`: validation boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。blocking failure は create job を禁止する。
- `入力開始点`: Job Setup UI。
- `主要 outcome`: `Ready` job は作成されず、failure reason が表示される。
- `開始操作`: validation を実行する。
- `入力方法`: 不足または不整合を含む設定を選ぶ。
- `主要操作列`: validation failure を確認し、設定を直し、再 validation できることを確認する。
- `手順`:
  1. 無効な設定で validation を実行する。
  2. failure reason と create 可否を確認する。
  3. 設定を修正して再 validation を実行する。
- `期待結果`:
  1. failure reason が UI に表示される。
  2. blocking failure がある間は create job が無効または拒否される。
  3. 修正後に validation 結果が最新設定へ更新される。
  4. validation 実行履歴は structured log に残り、アプリ状態は直近結果、対象設定断面、失敗カテゴリだけを保持する。
- `観測点`: validation summary、create button state、job 未作成の永続化結果。
- `UI-visible outcome`: 失敗理由、修正対象、再 validation の状態が見える。
- `fake_or_stub`: invalid setting fixture、stale foundation ref fixture、fake transport。
- `責務境界メモ`: 入力 cache 欠落時は Job Setup 内で再構築せず、Input Review の再構築導線へ戻す。

### SCN-TJS-003 複数入力を混線させず job を準備する

- `分類`: 境界条件
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 複数入力の job 作成で入力データと execution setting を混線させない。
- `受け入れ条件`: 各 job はちょうど 1 つの入力データだけを参照し、同一入力への 2 件目 job 作成を禁止する。
- `事前条件`: 取り込み済み入力データが複数ある。
- `public_seam_or_api_boundary`: job create boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。`TRANSLATION_JOB -> X_EDIT_EXTRACTED_DATA` は 1:1 参照にする。
- `入力開始点`: 複数入力 fixture。
- `主要 outcome`: job と input の対応が一意に観測できる。
- `開始操作`: それぞれの入力で job create を実行する。
- `入力方法`: input A と input B を切り替えて setup を作る。
- `主要操作列`: input A の job を作成し、input B の job を作成し、参照先を確認する。
- `手順`:
  1. input A で job を作成する。
  2. input B で job を作成する。
  3. job detail または repository query で input 参照を確認する。
- `期待結果`:
  1. job A は input A だけを参照する。
  2. job B は input B だけを参照する。
  3. foundation 参照と validation 結果が別入力へ混線しない。
  4. 同一入力に既存 job がある場合、状態に関係なく create は拒否される。
- `観測点`: job list、job detail、repository query。
- `UI-visible outcome`: 各 job の入力名、出自、状態が分離表示される。
- `fake_or_stub`: two-input fixture、temp DB。
- `責務境界メモ`: 過去 job を廃棄できる手段は別途必要である。

### SCN-TJS-004 AI 基盤設定を復元し secret を露出しない

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 保存済み provider / model / credential 参照を使い、Job Setup で再入力なしに validation へ進む。
- `受け入れ条件`: API key 平文は UI、DB、validation summary、監査表示に出ない。
- `事前条件`: 保存済み AI 設定と secret store 参照がある。
- `public_seam_or_api_boundary`: AI settings read boundary、validation boundary。
- `contract_freeze`: あり。credential 解決、provider capability、ネットワーク到達性は blocking validation 条件にする。
- `入力開始点`: 保存済み AI 設定 fixture。
- `主要 outcome`: provider、model、credential 参照状態、execution mode を確認できる。
- `開始操作`: Job Setup を開く。
- `入力方法`: 保存済み設定を確認し、必要時だけ provider / mode を変更する。
- `主要操作列`: 設定復元、validation、secret 非露出確認。
- `手順`:
  1. 保存済み AI 設定がある状態で Job Setup を開く。
  2. provider、model、credential 参照状態を確認する。
  3. validation を実行する。
- `期待結果`:
  1. provider と model が復元される。
  2. API key 平文は表示されない。
  3. fake transport で validation 経路を検証できる。
  4. credential 解決、provider capability、ネットワーク到達性の失敗は blocking として表示される。
- `観測点`: UI 表示、validation result、secret redaction、external request 未実行証跡。
- `UI-visible outcome`: 保存済み credential があるかだけが表示される。
- `fake_or_stub`: fake secret store、fake transport。
- `責務境界メモ`: API key 保存そのものの UI は job setup では扱わない。

### SCN-TJS-005 Ready 未成立では Running へ進めない

- `分類`: 責務境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: validation 未通過または job 未作成の状態から翻訳 phase 実行を開始しない。
- `受け入れ条件`: Running への遷移は Ready job の成立後だけ許可される。
- `事前条件`: validation failure の setup、Ready job の fixture がある。
- `public_seam_or_api_boundary`: job state transition boundary。
- `contract_freeze`: あり。Draft / validation failure から Running へ直接進まない。
- `入力開始点`: job state fixture。
- `主要 outcome`: Ready 以外では phase run が開始されない。
- `開始操作`: 翻訳実行開始を試行する。
- `入力方法`: Draft または validation failure の setup を使う。
- `主要操作列`: 実行開始を試行し、拒否と理由を確認する。
- `手順`:
  1. validation 未通過の setup を用意する。
  2. 実行開始相当の boundary を呼ぶ。
  3. 状態と error kind を確認する。
- `期待結果`:
  1. Running 状態は作成されない。
  2. Ready 未成立理由が返る。
  3. job setup の責務は create 前 validation と Ready 作成までに留まる。
- `観測点`: state transition result、phase run 未作成、error kind。
- `UI-visible outcome`: 実行開始不可理由が確認できる。
- `fake_or_stub`: state fixture、temp DB。
- `責務境界メモ`: 実行開始 UI の詳細は後続 phase task へ送る。

### SCN-TJS-006 create 途中の保存失敗で partial state を残さない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: job と execution setting の保存が途中で失敗しても不完全な job を残さない。
- `受け入れ条件`: create 全体が atomic に失敗し、再試行可否を観測できる。
- `事前条件`: DB 書き込み失敗または整合性違反を起こす fixture がある。
- `public_seam_or_api_boundary`: job create transaction boundary。
- `contract_freeze`: あり。partial `TRANSLATION_JOB` や欠けた `JOB_PHASE_RUN` を残さない。
- `入力開始点`: valid setup fixture と failure injection。
- `主要 outcome`: 失敗後も永続化状態が整合している。
- `開始操作`: create job を実行する。
- `入力方法`: 保存失敗を起こす fixture を使う。
- `主要操作列`: create を実行し、失敗後の row count と error kind を確認する。
- `手順`:
  1. validation pass 済み setup を用意する。
  2. create 永続化途中で失敗させる。
  3. job と関連設定の永続化状態を確認する。
- `期待結果`:
  1. create は失敗として返る。
  2. `TRANSLATION_JOB` だけが残る partial state はない。
  3. UI または API response で再試行可否を確認できる。
- `観測点`: transaction result、row count、error kind。
- `UI-visible outcome`: 作成失敗と再試行可否が表示される。
- `fake_or_stub`: temp DB、failure injection。
- `責務境界メモ`: attempt 履歴 table は前提にしない。

### SCN-TJS-007 paid API なしで provider validation 経路を検証する

- `分類`: テスト容易性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: real provider list を保ちつつ、外部 request / SDK transport だけを fake に差し替える。
- `受け入れ条件`: paid な real AI API を呼ばずに validation result を観測できる。
- `事前条件`: test mode または DI 可能な AIProvider boundary がある。
- `public_seam_or_api_boundary`: AIProvider / transport boundary。
- `contract_freeze`: あり。fake provider を user-facing provider list に出さない。
- `入力開始点`: provider validation fixture。
- `主要 outcome`: validation 経路は共通で、外部通信だけが fake になる。
- `開始操作`: provider validation を実行する。
- `入力方法`: provider / model / execution mode fixture を渡す。
- `主要操作列`: fake transport を注入し、validation を実行し、外部 request 未実行を確認する。
- `手順`:
  1. real provider list を使う test fixture を用意する。
  2. transport を fake に差し替える。
  3. validation result と外部 request 証跡を確認する。
- `期待結果`:
  1. provider list に fake は表示されない。
  2. paid API は呼ばれない。
  3. validation result は Job Setup と同じ経路で返る。
- `観測点`: provider list、transport fake、request log、validation result。
- `UI-visible outcome`: なし。user-facing UI は real provider list のまま。
- `fake_or_stub`: fake transport、fixed provider response。
- `責務境界メモ`: user-facing scenario ではなく、scenario acceptance を支える lower-level 条件である。

## Acceptance Checks

- `REQ-TJS-001`: `SCN-TJS-001`, `SCN-TJS-003`, `SCN-TJS-005`
- `REQ-TJS-002`: `SCN-TJS-001`, `SCN-TJS-002`, `SCN-TJS-006`
- `REQ-TJS-003`: `SCN-TJS-004`, `SCN-TJS-007`
- `REQ-TJS-004`: `SCN-TJS-003`, `SCN-TJS-005`, `SCN-TJS-006`
- `REQ-TJS-005`: `SCN-TJS-001`, `SCN-TJS-002`, `SCN-TJS-004`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-setup/scenario-design.md --coverage docs/exec-plans/active/translation-job-setup/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-setup/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Questions

- なし。
