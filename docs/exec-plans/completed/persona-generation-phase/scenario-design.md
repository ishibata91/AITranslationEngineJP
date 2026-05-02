# Scenario Design: persona-generation-phase

- `skill`: scenario-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/persona-generation-phase.md`
- `topic_abbrev`: `PGP`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - 単語翻訳フェーズ完了後だけ NPC ペルソナ生成フェーズを開始できる。
  - Job Run で persona phase の current phase、progress、phase result を確認できる。
  - NPC 発話原文、NPC 属性メタデータ、会話文脈、共通ペルソナ参照から NPC ごとの生成対象を固定できる。
  - AI provider の persona 生成結果を、ジョブ内ペルソナまたは persona snapshot 参照へ写像できる。
  - 生成済み persona snapshot を本文翻訳フェーズの入力として参照できる。
  - provider 失敗、入力不備、保存失敗、partial state を successful Completed として扱わない。
  - secret、API key 平文、provider raw request / response、過剰な原文本文を UI、error summary、structured log に出さない。
- `non_goals`:
  - 本文翻訳フェーズの訳文生成は扱わない。
  - 共通ペルソナ構築フローの実行設計は扱わない。
  - ジョブ内ペルソナ flush の実行設計は扱わない。
  - product code、product test、docs 正本、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-PGP-NNN` を一意 key として扱う。

- `candidate_count`: 47
- `rejected`: 2
- `needs_human_decision`: 0
- `unresolved_conflicts`: 0
- `questionnaire`: `./scenario-design.questions.md`

candidate decision の未決はない。
候補生成 agent は再起動していない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは sidecar JSON に分離した。
人間質問票の回答を反映済みである。
この scenario matrix は人間レビュー待ちである。

### `REQ-PGP-001` 単語翻訳フェーズ完了後に NPC ペルソナ生成フェーズを開始する

- `source_requirement`: `tasks/usecases/persona-generation-phase.yaml` は precondition を単語翻訳フェーズ完了とし、Job Run から persona phase を開始する。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: term phase Completed、非 terminal job、active phase run なしを開始条件にする。Completed、Failed、Canceled は terminal state とし、RecoverableFailed は回復対象として扱う。

### `REQ-PGP-002` NPC ごとの生成対象と入力 snapshot を確定する

- `source_requirement`: persona phase は NPC 発話原文、NPC 属性メタデータ、会話文脈、共通ペルソナを input にする。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: 生成対象 summary は NPC、入力種類、対象件数、対象外理由を含む。共通ペルソナ hit 時は新規 `PERSONA` を作らず、job の persona snapshot 参照だけを固定する。生成対象 0 件は Completed とし、対象 0 件、provider 未実行、snapshot 空を result summary に出す。

### `REQ-PGP-003` AI provider で NPC ペルソナを生成する

- `source_requirement`: `docs/spec.md` は NPC の発言、種族、性別情報から AI に persona を生成させることを求める。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし。
- `fixed_decisions`: paid real API は検証前提にしない。fake provider と fixed response で検証する。Job Setup の persona 専用 provider、model、execution mode を継承する。1 NPC を 1 request unit とし、NPC 属性と会話文脈を同じ request で扱う。valid provider output は自動採用する。

### `REQ-PGP-004` ジョブ内ペルソナと persona snapshot を保存する

- `source_requirement`: task output はジョブ内ペルソナと翻訳時参照用 persona snapshot である。
- `requirement_kind`: persistence
- `needs_human_decision`: なし。
- `fixed_decisions`: `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA` の partial state は Completed にしない。共通ペルソナ hit 時は新規 `PERSONA` を作らず、job の persona snapshot 参照だけを固定する。

### `REQ-PGP-005` Job Run で phase、progress、生成結果、snapshot 参照状態を確認する

- `source_requirement`: Job Run manual check は progress、phase result、ジョブ内ペルソナの参照状態確認を含む。
- `requirement_kind`: display
- `needs_human_decision`: なし。
- `fixed_decisions`: UI は mock ではなく表示項目、主要操作、有効条件、状態差分の契約として固定する。persona phase は pause、resume、retry、cancel を body phase と同じように許可する。

### `REQ-PGP-006` persona phase 完了後だけ body phase readiness を成立させる

- `source_requirement`: body phase は precondition として persona phase 完了を要求する。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: body phase の訳文生成は扱わず、readiness と input summary までを扱う。一部 NPC 失敗時は成功分を維持し、phase は RecoverableFailed として未処理 NPC だけ retry する。全対象が完了するまで body readiness は成立させない。生成対象 0 件は Completed とし、空 snapshot summary を body phase 入力として扱う。

### `REQ-PGP-007` 失敗、再開、リトライで partial state を成功扱いにしない

- `source_requirement`: job は中断、再開、失敗回復の対象であり、failure candidates は provider 失敗、入力不備、保存失敗を扱う。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: provider failure、invalid response、save failure は success として保存しない。一部 NPC 失敗時は成功分を維持し、phase は RecoverableFailed として未処理 NPC だけ retry する。persona phase は pause、resume、retry、cancel を body phase と同じように許可する。

### `REQ-PGP-008` 監査要約と redaction を満たす

- `source_requirement`: APIKey は secret store 参照で扱い、UI と log は redacted summary を表示する。
- `requirement_kind`: security
- `needs_human_decision`: なし。
- `fixed_decisions`: UI と DB summary には ID、digest、件数、evidence ref、redacted phase result summary だけを出し、全文と raw prompt は出さない。ログまたは debug 用には prompt / request body を確認できる導線を持つ。secret と API key はどこにも出さない。Job Run 再表示用の redacted phase result summary は保持するが、直接 DB 保存に限定せず、進行中の job state から復元できる形を許容する。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

Q-001 から Q-010 まで回答済みである。
回答内容は `./scenario-design.questions.md` と coverage JSON に保持する。

## Risks

- `implementation_risks`:
  - debug log に prompt / request body を出すため、secret と API key の redaction assertion が必須になる。
- `test_data_risks`:
  - common persona hit / miss、orphan NPC reference、partial save failure、invalid provider response の fixture が必要になる。
  - redaction は UI、structured log、fake transport log、phase result summary の複数観測点で確認する必要がある。

## Rules

- ケース ID は `SCN-PGP-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が残る間は scenario 完了にしない。
- 未解決 conflict が残る間は scenario 完了にしない。
- paid real AI API を前提にしない。

## Draft Scenario Matrix

以下を人間レビュー対象の scenario matrix とする。

### SCN-PGP-001 単語翻訳フェーズ完了後に persona phase を開始する

- `分類`: 正常系 / 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Job Run で persona phase を開始し、current phase と progress を確認する。
- `受け入れ条件`: term phase Completed、非 terminal job、active phase run なしの場合だけ persona phase を開始できる。
- `事前条件`: term phase Completed の job と、term 未完了 / active phase / terminal job fixture がある。
- `public_seam_or_api_boundary`: phase start boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。開始条件、開始拒否、current phase、progress。
- `入力開始点`: Job Run UI。
- `主要 outcome`: persona phase の current phase と progress が表示される。
- `開始操作`: Job Run を開く。
- `入力方法`: 対象 job を選ぶ。
- `主要操作列`: NPC ペルソナ生成フェーズ開始を実行し、開始結果を確認する。
- `期待結果`:
  1. persona phase が current phase として表示される。
  2. progress と phase run 開始結果を確認できる。
  3. term 未完了、active phase あり、terminal job では開始不可理由が表示される。
- `観測点`: Job Run UI、phase start result、`JOB_PHASE_RUN`。
- `UI-visible outcome`: current phase、progress、開始不可理由。
- `fake_or_stub`: term completed job fixture、not-ready job fixture、active phase fixture、temp DB。

### SCN-PGP-002 NPC ごとの生成対象と入力 snapshot を確認する

- `分類`: 正常系 / 境界値
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: NPC 発話原文、属性、会話文脈、共通ペルソナから生成対象 snapshot を作る。
- `受け入れ条件`: 生成対象 NPC、入力種類、対象件数、common persona hit/miss、対象外理由を確認できる。
- `事前条件`: NPC record、translation field reference、common persona hit / miss、orphan reference fixture がある。
- `public_seam_or_api_boundary`: target extraction / input snapshot boundary。
- `contract_freeze`: あり。target count、input snapshot digest、common persona status、対象外理由。
- `入力開始点`: persona phase start 後の target extraction。
- `主要 outcome`: provider input へ進める対象 NPC と進めない対象 NPC が区別される。
- `開始操作`: 生成対象抽出を実行する。
- `入力方法`: NPC record と会話文脈 fixture を使う。
- `主要操作列`: NPC 解決、共通ペルソナ照合、target summary 作成を確認する。
- `期待結果`:
  1. 対象 NPC 件数と input 種類を確認できる。
  2. 参照不能な NPC / 会話文脈は provider request に入らない。
  3. common persona hit / miss と対象外理由が summary に出る。
- `観測点`: target summary、input snapshot digest、`NPC_PROFILE`、`NPC_RECORD`。
- `UI-visible outcome`: 対象件数、common persona hit / miss、対象外理由。
- `fake_or_stub`: NPC record fixture、common persona fixture、orphan reference fixture。

### SCN-PGP-003 NPC persona を AI provider で生成する

- `分類`: 正常系 / 外部連携
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: NPC ごとの入力 snapshot から provider request を作り、valid response を persona result へ変換する。
- `受け入れ条件`: fake provider で persona result を取得でき、secret と raw payload を露出しない。
- `事前条件`: fixed provider response、provider failure、invalid response、common persona fixture がある。
- `public_seam_or_api_boundary`: AIProvider / prompt builder / response adapter boundary。
- `contract_freeze`: あり。provider setting、request unit、response validation、redaction。
- `入力開始点`: target snapshot fixture。
- `主要 outcome`: persona result が adapter output として得られる。
- `開始操作`: persona provider adapter を実行する。
- `入力方法`: target snapshot、provider setting、fixed response を渡す。
- `主要操作列`: prompt input mapping、fake provider execution、response validation を確認する。
- `期待結果`:
  1. provider、model、execution mode の要約を確認できる。
  2. valid response は persona result へ写像される。
  3. invalid response は persona として保存されない。
  4. paid real API を呼ばずに検証できる。
- `観測点`: adapter output、fake transport log、provider execution summary。
- `UI-visible outcome`: なし。表示は SCN-PGP-005 に統合する。
- `fake_or_stub`: fake provider、fixed response fixture、invalid response fixture。

### SCN-PGP-004 job-scoped persona と phase link を保存する

- `分類`: 正常系 / 永続化
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: valid persona result を job-scoped persona、evidence、phase link、persona snapshot summary へ保存する。
- `受け入れ条件`: body phase が参照できる persona snapshot と `PHASE_RUN_PERSONA` が整合する。
- `事前条件`: valid persona result、common persona hit / miss、save failure fixture がある。
- `public_seam_or_api_boundary`: persona persistence boundary。
- `contract_freeze`: あり。job-scoped persona、evidence、phase link、snapshot summary、partial state reject。
- `入力開始点`: persona result fixture。
- `主要 outcome`: job 内 persona または snapshot ref が body phase input summary から参照できる。
- `開始操作`: persona result 保存を実行する。
- `入力方法`: persona result、evidence ref、phase run ID を渡す。
- `主要操作列`: 保存、phase link 作成、snapshot summary 更新を確認する。
- `期待結果`:
  1. `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA` が整合する。
  2. partial save failure は Completed にならない。
  3. common persona hit 時の保存対象は人間回答に従う。
- `観測点`: `PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA`、phase result。
- `UI-visible outcome`: generated count、snapshot reference status。
- `fake_or_stub`: temp DB、save failure injection、common persona fixture。

### SCN-PGP-005 Job Run で persona phase result を確認する

- `分類`: 表示系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Job Run で phase、progress、生成対象、生成結果、persona snapshot 参照状態を見る。
- `受け入れ条件`: success、running、paused、recoverable failed、blocked、empty completed の状態差分が見える。
- `事前条件`: phase summary fixture、success result、failure result、empty result、snapshot missing result がある。
- `public_seam_or_api_boundary`: Job Run persona phase summary boundary。
- `contract_freeze`: あり。表示項目、button enablement、state variants。
- `入力開始点`: Job Run UI。
- `主要 outcome`: phase result と次操作可否を UI で判断できる。
- `開始操作`: Job Run を開く。
- `入力方法`: 対象 job を選ぶ。
- `主要操作列`: result summary、persona snapshot summary、error summary、次操作可否を確認する。
- `期待結果`:
  1. current phase、phase state、progress、target count、generated count が表示される。
  2. persona snapshot ID または参照状態、missing count、body readiness が表示される。
  3. error summary は secret、raw payload、過剰本文を含まない。
- `観測点`: Job Run UI、gateway response、button enablement。
- `UI-visible outcome`: current phase、progress、phase result、snapshot reference status、次操作可否。
- `fake_or_stub`: mocked gateway result fixture。

### SCN-PGP-006 persona phase 完了後だけ body readiness を成立させる

- `分類`: 後続フェーズ境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: persona phase 完了と snapshot 参照成立後だけ body phase readiness を true にする。
- `受け入れ条件`: persona 未完了、失敗、snapshot 参照不能では body phase run を作らない。
- `事前条件`: Completed、Running、Paused、RecoverableFailed、snapshot missing fixture がある。
- `public_seam_or_api_boundary`: next phase readiness boundary。
- `contract_freeze`: あり。persona phase completion requirement、snapshot reference requirement。
- `入力開始点`: body phase readiness query または body phase start 試行。
- `主要 outcome`: body phase の開始可否と拒否理由を確認できる。
- `開始操作`: body phase readiness を確認する。
- `入力方法`: phase state fixture と snapshot fixture を使う。
- `主要操作列`: readiness query、開始拒否理由、phase run 件数を確認する。
- `期待結果`:
  1. persona phase Completed かつ snapshot 参照可能時だけ body readiness が true になる。
  2. 未完了または参照不能では body phase run が作成されない。
  3. input summary に persona count、missing count、snapshot digest が出る。
- `観測点`: phase transition result、body input summary、`PHASE_RUN_PERSONA`。
- `UI-visible outcome`: body phase 開始可否、開始不可理由。
- `fake_or_stub`: phase state fixture、temp DB。

### SCN-PGP-007 provider 失敗、入力不備、保存失敗を成功扱いにしない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider failure、invalid response、input missing、save failure で persona snapshot を成功公開しない。
- `受け入れ条件`: error kind、retryable flag、phase state、progress、body readiness false を確認できる。
- `事前条件`: provider failure、invalid response、orphan reference、save failure fixture がある。
- `public_seam_or_api_boundary`: persona phase execution boundary。
- `contract_freeze`: あり。暗黙 fallback なし、invalid response reject、partial persistence reject。
- `入力開始点`: failure injection fixture。
- `主要 outcome`: failed target は persona として保存されず、再試行可否を観測できる。
- `開始操作`: persona phase を実行する。
- `入力方法`: failure fixture を使う。
- `主要操作列`: provider request、response validation、persistence failure、phase result を確認する。
- `期待結果`:
  1. 別 provider へ暗黙 fallback しない。
  2. invalid response は persona として保存されない。
  3. partial state は Completed にならない。
  4. body readiness は false のままである。
- `観測点`: phase result、error kind、row count、body readiness。
- `UI-visible outcome`: 失敗理由、再試行可否、後続 phase 不可理由。
- `fake_or_stub`: fake transport、invalid response fixture、save failure injection、temp DB。

### SCN-PGP-008 再送、再開、リトライで persona を重複作成しない

- `分類`: 回復系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 同じ persona phase run で開始再送、再開、リトライを扱い、persona と phase link を重複作成しない。
- `受け入れ条件`: 同じ `JOB_PHASE_RUN` を継続し、成功済み persona と未処理 target を区別する。
- `事前条件`: active phase run、paused phase run、recoverable failed phase run、partial success fixture がある。
- `public_seam_or_api_boundary`: phase resume / retry boundary。
- `contract_freeze`: あり。same phase run reuse、target snapshot stability、duplicate persona guard。
- `入力開始点`: existing phase run fixture。
- `主要 outcome`: phase run と persona link が二重作成されない。
- `開始操作`: 開始再送、再開、リトライを実行する。
- `入力方法`: 同一 job / same phase type fixture を使う。
- `主要操作列`: phase run ID、persona count、未処理 count、latest error、progress を確認する。
- `期待結果`:
  1. phase run ID は同じである。
  2. 成功済み persona と `PHASE_RUN_PERSONA` は重複しない。
  3. 未処理 NPC だけ provider request 対象へ戻る。
  4. latest error と progress が更新される。
- `観測点`: phase run ID、row count、progress、latest error。
- `UI-visible outcome`: 再開またはリトライ結果。
- `fake_or_stub`: phase run fixture、partial success fixture、temp DB。

### SCN-PGP-009 redacted summary と debug log を確認する

- `分類`: セキュリティ / 観測性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: UI と DB summary は redacted に保ち、debug log では prompt / request body を調整用に確認できる。
- `受け入れ条件`: provider、model、execution mode、credential ref、input count、output count、digest、error kind は観測できる。UI と DB summary には保存禁止情報が出ない。debug log に prompt / request body が出る場合も secret と API key は出ない。
- `事前条件`: fake secret store、fake transport、structured log capture、redaction assertion fixture がある。
- `public_seam_or_api_boundary`: runtime log / audit summary boundary。
- `contract_freeze`: あり。UI / DB summary の redaction、credential_ref、debug log の prompt / request body、secret 非露出。
- `入力開始点`: phase execution fixture。
- `主要 outcome`: 障害調査に必要な要約を確認でき、保存禁止情報は出ない。
- `開始操作`: persona phase を実行する。
- `入力方法`: provider / model / credential ref / prompt digest fixture を使う。
- `主要操作列`: success、failure、redaction assertion を確認する。
- `期待結果`:
  1. API key と secret 本体は表示またはログに出ない。
  2. UI と DB summary には provider raw request / response、full prompt、原文発話全文が出ない。
  3. provider、model、execution mode、input count、output count、prompt digest を確認できる。
  4. debug log では prompt / request body を確認でき、secret と API key は redacted される。
- `観測点`: structured log、debug log capture、Job Run summary、fake secret store assertion、fake transport log。
- `UI-visible outcome`: provider / model / credential 参照状態、短い error summary。
- `fake_or_stub`: fake secret store、fake transport、log capture。

### SCN-PGP-010 terminal job には persona と body readiness を後書きしない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: terminal job で persona phase 開始、persona 保存、body readiness 更新を拒否する。
- `受け入れ条件`: terminal job の `JOB_PHASE_RUN`、`PERSONA`、`PHASE_RUN_PERSONA`、body readiness は変更されない。
- `事前条件`: terminal job fixture がある。
- `public_seam_or_api_boundary`: terminal job guard boundary。
- `contract_freeze`: あり。terminal state guard、state invariant。
- `入力開始点`: persona phase start、persona save、body readiness update。
- `主要 outcome`: terminal job の state 不変と拒否理由を確認できる。
- `開始操作`: terminal job で persona phase 関連操作を試行する。
- `入力方法`: terminal job fixture を使う。
- `主要操作列`: 操作前後の phase run、persona、readiness を確認する。
- `期待結果`:
  1. terminal job では persona phase run が作成されない。
  2. persona と phase link は後書きされない。
  3. body readiness は変更されない。
- `観測点`: phase transition result、row count、拒否理由、state snapshot。
- `UI-visible outcome`: terminal job の開始不可理由。
- `fake_or_stub`: terminal job fixture、temp DB。

## Acceptance Checks

- `REQ-PGP-001`: `SCN-PGP-001`, `SCN-PGP-010`
- `REQ-PGP-002`: `SCN-PGP-002`, `SCN-PGP-003`
- `REQ-PGP-003`: `SCN-PGP-003`, `SCN-PGP-007`, `SCN-PGP-009`
- `REQ-PGP-004`: `SCN-PGP-004`, `SCN-PGP-008`
- `REQ-PGP-005`: `SCN-PGP-001`, `SCN-PGP-005`
- `REQ-PGP-006`: `SCN-PGP-004`, `SCN-PGP-006`
- `REQ-PGP-007`: `SCN-PGP-007`, `SCN-PGP-008`
- `REQ-PGP-008`: `SCN-PGP-005`, `SCN-PGP-009`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/persona-generation-phase/scenario-design.md --coverage docs/exec-plans/active/persona-generation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/persona-generation-phase/scenario-design.candidate-coverage.json --report-out docs/exec-plans/active/persona-generation-phase/scenario-design.requirement-gate.md --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Decisions

- なし。
