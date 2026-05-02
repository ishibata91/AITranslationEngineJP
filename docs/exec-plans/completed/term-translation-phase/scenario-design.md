# Scenario Design: term-translation-phase

- `skill`: scenario-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/term-translation-phase.md`
- `topic_abbrev`: `TTP`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - Ready の翻訳ジョブだけが単語翻訳フェーズを開始できる。
  - 単語翻訳フェーズの current phase、progress、phase result を Job Run で確認できる。
  - 共通辞書に完全一致する語は provider request へ送らず、置換対象判定として観測できる。
  - 共通辞書にない用語や固有名詞は、確定訳語としてジョブ内辞書へ反映できる。
  - ジョブ内辞書は対象 job だけに紐づき、本文翻訳フェーズ前の入力として参照できる。
  - provider 失敗、応答不正、保存失敗では secret と raw response を露出せず、辞書 partial state を成功扱いにしない。
- `non_goals`:
  - 本文翻訳フェーズの訳文生成は扱わない。
  - NPC ペルソナ生成フェーズの設計は扱わない。
  - 共通辞書管理 UI、xTranslator import、xTranslator export は扱わない。
  - product code、product test、docs 正本、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-TTP-NNN` を一意 key として扱う。

candidate decision の `needs_human_decision` は 0 件である。
未解決 conflict は 0 件である。
候補生成 agent は再起動していない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは sidecar JSON に分離する。
人間質問票の回答を反映済みである。
この scenario matrix は人間レビュー待ちである。

### `REQ-TTP-001` Ready job から単語翻訳フェーズを開始する

- `source_requirement`: Ready job の Job Run から単語翻訳フェーズを開始し、current phase と progress を観測する。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: Ready 以外、terminal job、既存 active phase run がある job は開始を拒否する。

### `REQ-TTP-002` 共通辞書の完全一致語を除外する

- `source_requirement`: 共通辞書として登録済みの完全一致語は翻訳対象にせず、置換対象判定として保持する。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 完全一致ではない hit は `cached` または置換済みとして保存しない。共通辞書は phase 開始時 snapshot で固定する。共通辞書除外後に対象語が 0 件でも phase result は Completed とし、provider 未実行を result summary に出す。

### `REQ-TTP-003` 用語や固有名詞の訳語を確定する

- `source_requirement`: 共通辞書にない用語や固有名詞を単語翻訳フェーズで訳し、確定訳語として扱う。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: paid な real AI API は scenario validation の前提にしない。fake transport で provider 応答を検証する。AI provider の用語翻訳結果は自動で確定訳語にする。初期実装は 1 対象語 1 request unit とし、Batch API を使う場合も batch item は 1 対象語単位にする。単語翻訳フェーズは Job Setup の単語翻訳用 provider、model、execution mode を継承する。

### `REQ-TTP-004` 確定訳語をジョブ内辞書へ反映する

- `source_requirement`: 確定訳語を対象 job のジョブ内辞書として保存し、後続フェーズが参照できるようにする。
- `requirement_kind`: persistence
- `needs_human_decision`: なし
- `fixed_decisions`: `DICTIONARY_ENTRY` と `PHASE_RUN_DICTIONARY_ENTRY` の整合が取れない partial state は成功扱いにしない。再開、リトライ、再実行では既存 entry を維持し、未処理だけ進める。共通辞書完全一致は API を投げずに確定済みとして適用する。同一 source term の重複判定 key は record type ごとに一意にする。

### `REQ-TTP-005` Job Run で phase result を確認する

- `source_requirement`: progress、phase result、確定訳語、ジョブ内辞書反映、エラー理由を Job Run で確認する。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: UI は mock ではなく、表示項目、操作、有効条件、状態差分の契約として固定する。

### `REQ-TTP-006` 未完了または失敗時に後続フェーズへ進めない

- `source_requirement`: 単語翻訳フェーズ未完了、失敗、ジョブ内辞書参照不能では本文翻訳フェーズへ進めない。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: term phase 未完了時は後続 phase run を作成しない。job は Running のまま phase state で完了、中断、回復可能失敗、再実行準備を区別する。

### `REQ-TTP-007` 監査と redaction を満たす

- `source_requirement`: phase 開始、除外判定、AI 実行、辞書反映、失敗、中断を後追い確認でき、secret と過剰本文を保存しない。
- `requirement_kind`: security
- `needs_human_decision`: なし
- `fixed_decisions`: API key、secret 本体、復号可能な値、provider raw request / response、翻訳フィールド本文の全文は UI、error summary、structured log に出さない。共通辞書 snapshot の digest または version を監査要約へ残す。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

Q-001 から Q-009 まで回答済みである。
回答内容は `./scenario-design.questions.md` と coverage JSON の `human_answer` に保持する。

## Risks

- `implementation_risks`:
  - Job Setup に単語翻訳用 provider、model、execution mode の設定が必要になる。
  - 1 対象語 1 request unit にするため、大量対象語では request unit 数と progress 表示の整合が重要になる。
  - record type ごとの重複 key により、本文翻訳側の辞書参照でも record type context が必要になる。
- `test_data_risks`:
  - 1 対象語 request unit の response 欠落、invalid response、save failure の fixture を分ける必要がある。
  - secret redaction は UI、structured log、fake transport log の複数観測点で確認する必要がある。

## Rules

- ケース ID は `SCN-TTP-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が再発した場合は scenario 完了にしない。
- 未解決 conflict が再発した場合は scenario 完了にしない。
- `not_applicable` と `deferred` は理由なしで通さない。
- paid な real AI API を前提にしない。

## Draft Scenario Matrix

人間質問票の回答を反映済みである。
以下を人間レビュー対象の scenario matrix とする。

### SCN-TTP-001 Ready job から単語翻訳フェーズを開始する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Ready job の Job Run で単語翻訳フェーズを開始し、current phase と progress を確認する。
- `受け入れ条件`: Ready job だけが単語翻訳フェーズを開始できる。
- `事前条件`: `translation-job-setup` が完了し、Ready job が存在する。
- `public_seam_or_api_boundary`: phase start boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。Ready 以外からの開始拒否、active phase run 重複防止、progress 表示。
- `入力開始点`: Job Run UI。
- `主要 outcome`: 単語翻訳フェーズの current phase と progress が表示される。
- `開始操作`: Job Run を開く。
- `入力方法`: Ready job を選ぶ。
- `主要操作列`: 単語翻訳フェーズ開始を実行し、current phase と progress を確認する。
- `期待結果`:
  1. 単語翻訳フェーズが current phase として表示される。
  2. progress と phase run 開始結果を確認できる。
  3. Ready 以外では開始不可理由が表示される。
- `観測点`: Job Run UI、phase start result、`JOB_PHASE_RUN`。
- `UI-visible outcome`: current phase、progress、開始不可理由。
- `fake_or_stub`: Ready job fixture、not-ready job fixture、temp DB。
- `責務境界メモ`: job は Running のまま、phase state で単語翻訳フェーズの完了、中断、回復可能失敗、再実行準備を区別する。

### SCN-TTP-002 共通辞書の完全一致語を provider request から除外する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 共通辞書完全一致語を provider request へ含めず、置換対象判定として観測する。
- `受け入れ条件`: 完全一致 hit だけが除外され、非完全一致 hit は除外されない。
- `事前条件`: 共通辞書 fixture と翻訳対象語 fixture がある。
- `public_seam_or_api_boundary`: term extraction / provider request filtering boundary。
- `contract_freeze`: あり。完全一致、provider request filtering、`cached` 相当の内部観測情報。
- `入力開始点`: Ready job fixture。
- `主要 outcome`: provider request には未一致語だけが含まれる。
- `開始操作`: 単語翻訳フェーズを実行する。
- `入力方法`: 共通辞書 hit と未一致語を含む fixture を使う。
- `主要操作列`: 対象語抽出、共通辞書照合、provider request payload 確認を行う。
- `期待結果`:
  1. 完全一致語は provider request から除外される。
  2. 除外件数と置換対象判定を phase result で確認できる。
  3. 非完全一致語は `cached` として保存されない。
- `観測点`: provider request payload、phase result、置換対象判定。
- `UI-visible outcome`: 除外件数、置換対象件数、未一致件数。
- `fake_or_stub`: common dictionary fixture、fake AI transport。
- `責務境界メモ`: 共通辞書は phase 開始時 snapshot で固定する。対象語 0 件でも phase result は Completed とし、provider 未実行を result summary に出す。

### SCN-TTP-003 provider 応答を確定訳語候補として扱う

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 未登録の用語や固有名詞を provider へ 1 対象語 1 request unit で渡し、応答を確定訳語へ写像できる。
- `受け入れ条件`: provider 応答の原語、訳語、対応 key を失わずに内部 contract へ変換し、自動で確定訳語にできる。
- `事前条件`: 未登録語 fixture と fake provider response がある。
- `public_seam_or_api_boundary`: AIProvider / response adapter boundary。
- `contract_freeze`: あり。fake transport、response correlation、secret 非露出。
- `入力開始点`: provider response fixture。
- `主要 outcome`: provider 応答が job dictionary entry 候補へ変換される。
- `開始操作`: 単語翻訳 provider adapter を実行する。
- `入力方法`: fixed response、1 対象語 request unit の batch item response、invalid response fixture を使う。
- `主要操作列`: provider request、response parse、adapter output を確認する。
- `期待結果`:
  1. valid response は source term と translated term の対応を保持する。
  2. valid response は自動で確定訳語として扱われる。
  3. secret は request log、error summary、UI に出ない。
  4. paid API を呼ばずに検証できる。
- `観測点`: adapter output、fake transport log、redaction assertion。
- `UI-visible outcome`: なし。user-facing 表示は SCN-TTP-005 に統合する。
- `fake_or_stub`: fake transport、fixed provider response、1 対象語 request unit response fixture。
- `参考資料`:
  - [`xai_openapi_full.json`](../../../references/vendor-api/xai_openapi_full.json)
  - [`gemini batch ref.md`](../../../references/vendor-api/gemini%20batch%20ref.md)
- `責務境界メモ`: Job Setup の単語翻訳用 provider、model、execution mode を継承する。

### SCN-TTP-004 確定訳語をジョブ内辞書へ反映する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 確定訳語を対象 job のジョブ内辞書として保存し、phase run と紐づける。
- `受け入れ条件`: `DICTIONARY_ENTRY.translation_job_id` と `PHASE_RUN_DICTIONARY_ENTRY` から、対象 job と phase run を追跡できる。重複判定 key は record type ごとに一意である。
- `事前条件`: 確定訳語 fixture がある。
- `public_seam_or_api_boundary`: job dictionary persistence boundary。
- `contract_freeze`: あり。job-scoped dictionary、phase run link、partial state rollback。
- `入力開始点`: confirmed term fixture。
- `主要 outcome`: 本文翻訳フェーズが参照できる job dictionary entry が作成される。
- `開始操作`: 確定訳語保存を実行する。
- `入力方法`: source term、translated term、dictionary source、reusable 判定を渡す。
- `主要操作列`: 保存、phase link 作成、本文翻訳入力 summary 確認を行う。
- `期待結果`:
  1. ジョブ内辞書 entry が対象 job にだけ紐づく。
  2. 同一 job、同一 record type、同一 source term の重複 entry を作らない。
  3. phase run と辞書 entry の関連を確認できる。
  4. 保存途中失敗では partial state を成功扱いにしない。
- `観測点`: `DICTIONARY_ENTRY`、`PHASE_RUN_DICTIONARY_ENTRY`、phase result。
- `UI-visible outcome`: 確定訳語件数、ジョブ内辞書反映件数。
- `fake_or_stub`: temp DB、save failure injection。
- `責務境界メモ`: 共通辞書完全一致は API を投げずに確定済みとして扱う。再開、リトライ、再実行では既存 entry を維持し、未処理だけ進める。

### SCN-TTP-005 Job Run で単語翻訳フェーズ結果を確認する

- `分類`: 表示系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: ユーザーが phase result、progress、確定訳語、ジョブ内辞書反映を Job Run で確認する。
- `受け入れ条件`: 成功、空対象、失敗、再試行可能の各状態で必要な結果と次操作が見える。
- `事前条件`: 成功 result、error result、empty result の fixture がある。
- `public_seam_or_api_boundary`: Job Run summary boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。表示項目、状態差分、button enablement。
- `入力開始点`: Job Run UI。
- `主要 outcome`: phase result と次操作可否が UI で判断できる。
- `開始操作`: Job Run を開く。
- `入力方法`: 対象 job を選ぶ。
- `主要操作列`: result summary、辞書 summary、error summary、次 phase 操作を確認する。
- `期待結果`:
  1. phase、progress、確定訳語件数、辞書反映件数が表示される。
  2. error summary は secret と raw response を含まない。
  3. 未完了または失敗時は後続 phase 開始不可理由が表示される。
- `観測点`: Job Run UI、error summary、button enablement。
- `UI-visible outcome`: current phase、progress、phase result、次操作可否。
- `fake_or_stub`: mocked gateway result fixture。
- `責務境界メモ`: 詳細な UI 契約は `./ui-design.md` を正本にする。

### SCN-TTP-006 provider 失敗や応答不正で辞書を汚さない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider 失敗、応答不正、保存失敗を成功扱いにせず、辞書 partial state を残さない。
- `受け入れ条件`: 失敗時は error kind、retryable flag、phase state、辞書未更新を確認できる。
- `事前条件`: provider failure、invalid response、save failure fixture がある。
- `public_seam_or_api_boundary`: term phase execution boundary。
- `contract_freeze`: あり。暗黙 fallback なし、invalid response reject、atomic persistence。
- `入力開始点`: failure injection fixture。
- `主要 outcome`: 確定訳語は保存されず、再試行可否を観測できる。
- `開始操作`: 単語翻訳フェーズを実行する。
- `入力方法`: failure fixture を使う。
- `主要操作列`: provider request、response validation、persistence failure を観測する。
- `期待結果`:
  1. 別 provider へ暗黙 fallback しない。
  2. invalid response は確定訳語として保存されない。
  3. 1 対象語 request unit の response 欠落は、その対象語の failed / retryable として扱う。
  4. partial dictionary state は残らない。
  5. secret と raw response は表示またはログに出ない。
- `観測点`: phase result、error kind、row count、structured log。
- `UI-visible outcome`: 失敗理由、再試行可否、後続 phase 不可理由。
- `fake_or_stub`: fake transport、invalid response fixture、temp DB。
- `参考資料`:
  - [`xai_openapi_full.json`](../../../references/vendor-api/xai_openapi_full.json)
  - [`gemini batch ref.md`](../../../references/vendor-api/gemini%20batch%20ref.md)
- `責務境界メモ`: 複数対象語を 1 応答に詰める前提は置かない。Batch API を使う場合も batch item は 1 対象語単位にする。

### SCN-TTP-007 単語翻訳フェーズ未完了では後続フェーズへ進めない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 単語翻訳フェーズ未完了、失敗、辞書参照不能では後続フェーズを開始しない。
- `受け入れ条件`: 後続 phase run は作成されず、拒否理由を確認できる。
- `事前条件`: 未開始、Running、Paused、RecoverableFailed、Completed の phase fixture がある。
- `public_seam_or_api_boundary`: next phase start boundary。
- `contract_freeze`: あり。term phase completion requirement。
- `入力開始点`: job phase state fixture。
- `主要 outcome`: term phase 完了後だけ後続 phase が開始可能になる。
- `開始操作`: 後続 phase 開始を試行する。
- `入力方法`: 未完了または失敗中の job fixture を使う。
- `主要操作列`: 後続 phase 開始を試行し、拒否理由と state 不変を確認する。
- `期待結果`:
  1. 未完了または失敗中は後続 phase run が作成されない。
  2. 完了時だけ後続 phase 入力 summary が成立する。
  3. terminal job への後書きは拒否される。
- `観測点`: phase transition result、phase run 件数、拒否理由。
- `UI-visible outcome`: 後続 phase 開始不可理由。
- `fake_or_stub`: phase state fixture、temp DB。
- `責務境界メモ`: job は Running のまま、phase state が後続 phase 開始可否を決める。

### SCN-TTP-008 phase 再送、再開、リトライで重複作成しない

- `分類`: 回復系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 同一 job の単語翻訳 phase run、phase link、辞書 entry を重複作成しない。
- `受け入れ条件`: 再送、再開、リトライは同じ `JOB_PHASE_RUN` の状態を戻し、既存 entry を維持して未処理だけ進める。
- `事前条件`: active phase run、paused phase run、recoverable failed phase run fixture がある。
- `public_seam_or_api_boundary`: phase resume / retry boundary。
- `contract_freeze`: あり。same phase run reuse、attempt table なし。
- `入力開始点`: existing phase run fixture。
- `主要 outcome`: phase run と辞書 entry が二重作成されない。
- `開始操作`: 開始再送、再開、リトライを実行する。
- `入力方法`: 同一 job / same phase type fixture を使う。
- `主要操作列`: 操作前後の phase run ID、辞書 entry 件数、progress を確認する。
- `期待結果`:
  1. phase run ID は同じである。
  2. `DICTIONARY_ENTRY` と `PHASE_RUN_DICTIONARY_ENTRY` が重複しない。
  3. 既存 entry は再作成せず、未処理 term だけ provider request へ進む。
  4. retryable failure では最新 error と progress が更新される。
- `観測点`: phase run ID、row count、progress、latest error。
- `UI-visible outcome`: 再開またはリトライ結果。
- `fake_or_stub`: phase run fixture、temp DB。
- `責務境界メモ`: API コスト節約のため、再実行では既存 entry を維持する。

### SCN-TTP-009 監査要約と redaction を確認する

- `分類`: セキュリティ / 観測性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: phase 実行の再現材料を残し、secret と過剰本文を露出しない。
- `受け入れ条件`: structured log と UI summary は ID、件数、digest、version、error kind だけを持ち、secret と raw payload を持たない。
- `事前条件`: fake secret store、fake transport、structured log capture がある。
- `public_seam_or_api_boundary`: runtime log / audit summary boundary。
- `contract_freeze`: あり。redaction、credential_ref、raw request / response 保存禁止。
- `入力開始点`: phase execution fixture。
- `主要 outcome`: 障害調査に必要な要約を確認でき、保存禁止情報は出ない。
- `開始操作`: 単語翻訳フェーズを実行する。
- `入力方法`: provider / model / credential ref / prompt version fixture を使う。
- `主要操作列`: success、failure、redaction assertion を確認する。
- `期待結果`:
  1. API key と secret 本体は表示またはログに出ない。
  2. provider raw request / response と本文全量は保存されない。
  3. provider、model、execution mode、input count、output count、prompt version または digest を確認できる。
- `観測点`: structured log、Job Run summary、fake secret store assertion。
- `UI-visible outcome`: provider / model / credential 参照状態、短い error summary。
- `fake_or_stub`: fake secret store、fake transport、log capture。
- `参考資料`:
  - [`xai_openapi_full.json`](../../../references/vendor-api/xai_openapi_full.json)
  - [`gemini batch ref.md`](../../../references/vendor-api/gemini%20batch%20ref.md)
- `責務境界メモ`: 共通辞書 snapshot の digest または version を残し、full prompt と raw response は保存しない。

## Acceptance Checks

- `REQ-TTP-001`: `SCN-TTP-001`, `SCN-TTP-008`
- `REQ-TTP-002`: `SCN-TTP-002`, `SCN-TTP-005`
- `REQ-TTP-003`: `SCN-TTP-003`, `SCN-TTP-006`
- `REQ-TTP-004`: `SCN-TTP-004`, `SCN-TTP-008`
- `REQ-TTP-005`: `SCN-TTP-001`, `SCN-TTP-005`
- `REQ-TTP-006`: `SCN-TTP-007`
- `REQ-TTP-007`: `SCN-TTP-006`, `SCN-TTP-009`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/term-translation-phase/scenario-design.md --coverage docs/exec-plans/active/term-translation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/term-translation-phase/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Answered Decisions

- `Q-001`: AI provider の用語翻訳結果を自動で確定訳語にする。
- `Q-002`: 共通辞書参照は phase 開始時 snapshot で固定する。
- `Q-003`: 対象語 0 件でも phase result は Completed とする。
- `Q-004`: 再開、リトライ、再実行では既存 entry を維持し、未処理だけ進める。
- `Q-005`: 初期実装は 1 対象語 1 request unit とし、欠落は対象語単位の失敗として扱う。
- `Q-006`: 共通辞書完全一致は API を投げずに確定済みとして適用する。
- `Q-007`: 同一 source term の重複判定 key は record type ごとに一意にする。
- `Q-008`: 単語翻訳フェーズは Job Setup の単語翻訳用 provider、model、execution mode を継承する。
- `Q-009`: job は Running のまま phase state で単語翻訳フェーズの状態を区別する。
