# Scenario Design: body-translation-phase

- `skill`: scenario-design
- `status`: ready-human-review
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/body-translation-phase.md`
- `topic_abbrev`: `BTP`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - NPC ペルソナ生成フェーズ完了後だけ本文翻訳フェーズを開始できる。
  - Job Run で本文翻訳フェーズの current phase、progress、phase result、failure state を確認できる。
  - 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを同一 phase run の入力として参照できる。
  - 本文翻訳フェーズは `Job Setup` で設定した本文翻訳用の provider、model、execution mode を使う。
  - 完全一致した辞書 hit は provider request から除外し、部分一致は訳語固定制約として provider request に渡す。
  - 翻訳レコード種別と field type に応じた翻訳指示を構成できる。
  - 訳文、出力ステータス、保護要素検証結果を field 単位で保持できる。
  - 部分失敗では成功済み field result を保持し、phase 全体は回復可能な失敗として表示する。
  - 保護要素検証に失敗した訳文は保存前に拒否し、retry 対象にする。
  - 本文翻訳対象 0 件は正常完了にし、単語だけの plugin でも成果物出力へ進める。
  - 本文翻訳フェーズ完了時点で翻訳ジョブ全体を `Completed` とし、完了済み job から成果物を出力できる。
  - 取り消しは `Paused` からだけ可能にし、`Canceled` 後はフェーズ終端とし、途中成功結果は成果物出力に使わない。
  - provider 失敗、応答不正、保存失敗、保護要素検証失敗を successful Completed として扱わない。
  - secret、API key 平文、復号可能値を UI、error summary、structured log、debug log、fake transport log に出さない。
- `non_goals`:
  - 単語翻訳フェーズと NPC ペルソナ生成フェーズの再設計は扱わない。
  - translation-output-artifact の xTranslator row 生成規則は扱わない。
  - provider 実装方式、具体 API 名、DB migration、product code、product test、docs 正本、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-BTP-NNN` を一意 key として扱う。

- `candidate_count`: 64
- `adopted`: 11
- `merged`: 53
- `rejected`: 0
- `needs_human_decision`: 0
- `unresolved_conflicts`: 0
- `questionnaire`: `./scenario-design.questions.md`

candidate artifact 内の conflict hint は、最終 scenario の競合としては直接残さず、詳細要求の人間判断に集約した。
候補生成 agent は再起動していない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは sidecar JSON に分離した。
人間回答 10 件を反映済みである。
`needs_human_decision` は 0 件であり、この scenario matrix は人間設計レビューへ進められる。

### `REQ-BTP-001` NPC ペルソナ生成フェーズ完了後に本文翻訳フェーズを開始する

- `source_requirement`: `tasks/usecases/body-translation-phase.yaml` は precondition を NPC ペルソナ生成フェーズ完了とし、Job Run から本文翻訳フェーズを開始する。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: persona phase Completed、非 terminal job、active phase run なし、辞書と persona snapshot の参照成立を開始条件にする。

### `REQ-BTP-002` 入力 snapshot と provider request 入力を構成する

- `source_requirement`: 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータ、翻訳レコード種別を本文翻訳 request の入力として扱う。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: input summary は field 件数、辞書 snapshot digest、persona snapshot digest、metadata digest、prompt digest を持つ。provider、model、execution mode は `Job Setup` の本文翻訳用設定を使う。完全一致した辞書 hit は provider request から除外し、部分一致は訳語固定制約として provider request に渡す。

### `REQ-BTP-003` AI provider で本文翻訳を実行する

- `source_requirement`: 本文翻訳フェーズは AI provider を使い、訳文と provider execution summary を生成する。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし。
- `fixed_decisions`: paid real API は scenario validation の前提にしない。fake provider と fixed response で provider 境界を検証する。部分失敗では成功済み field result を保持し、phase 全体は回復可能な失敗として表示する。

### `REQ-BTP-004` 保護要素検証後に field result を保存する

- `source_requirement`: 訳文、出力ステータス、保護要素検証結果を `JOB_TRANSLATION_FIELD` と phase result から確認できるようにする。
- `requirement_kind`: persistence
- `needs_human_decision`: なし。
- `fixed_decisions`: 訳文と出力ステータスと保護要素検証結果は同じ field に対応付ける。保存失敗時は successful Completed にしない。保護要素検証に失敗した訳文は保存前に拒否し、retry 対象にする。出力ステータスは後続成果物出力に必要な語彙だけを細かく分ける。

### `REQ-BTP-005` Job Run で phase result と後続出力 readiness を確認する

- `source_requirement`: Job Run で progress、failure state、訳文、出力ステータス、保護要素検証結果、後続出力可否を確認する。
- `requirement_kind`: display
- `needs_human_decision`: なし。
- `fixed_decisions`: UI は mock ではなく表示項目、主要操作、有効条件、状態差分の契約として固定する。本文翻訳フェーズ完了時点で翻訳ジョブ全体を `Completed` とする。結果確認から戻る導線とフィールド単体編集は本 task では扱わない。

### `REQ-BTP-006` pause、resume、retry、cancel を同じ phase run で扱う

- `source_requirement`: 本文翻訳フェーズは pause、resume、retry、cancel の可否を確認でき、recoverable failure から再試行できる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: retry 可能な失敗は同じ `JOB_PHASE_RUN` を継続する。新しい attempt 履歴テーブルは前提にしない。取り消しは `Paused` からだけ可能にし、`Canceled` 後はフェーズ終端とし、途中成功結果は成果物出力に使わない。

### `REQ-BTP-007` 監査要約と redaction を満たす

- `source_requirement`: AI 実行条件、入力 summary、error kind、件数、digest は確認できるが、secret と過剰本文は出さない。
- `requirement_kind`: security
- `needs_human_decision`: なし。
- `fixed_decisions`: 原文と訳文がローカルで見えること自体は問題にしない。API key、secret 本体、復号可能値は UI、error summary、structured log、debug log、fake transport log に出さない。

### `REQ-BTP-008` 翻訳対象 0 件または空 source の境界を扱う

- `source_requirement`: 本文翻訳フェーズの対象 field が 0 件、または source text が空の場合の終点を固定する必要がある。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: 本文翻訳対象 0 件は正常完了にする。provider 未実行、target count、skipped count、output readiness への影響を観測対象にし、単語だけの plugin でも成果物出力へ進める。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

人間回答 10 件を記録済みである。
未回答扱いの項目は残さない。

## Risks

- `implementation_risks`:
  - provider request の粒度と retry 単位を誤ると、成功済み訳文の重複保存や出力ステータス不整合が起きる。
  - debug log と provider request summary の扱いを誤ると、secret の露出につながる。
  - body phase Completed と job-level Completed の反映を分けると、translation-output-artifact の開始条件がぶれる。
- `test_data_risks`:
  - 保護要素欠落、重複、順序違い、余分追加の fixture を分ける必要がある。
  - provider 部分欠落、batch correlation error、save failure、late response の fixture を分ける必要がある。
  - secret 非露出は UI、structured log、debug log、fake transport log、phase result summary の複数観測点で確認する必要がある。

## Rules

- ケース ID は `SCN-BTP-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` は 0 件である。
- 未解決 conflict が残る間は scenario 完了にしない。
- paid real AI API を前提にしない。

## Scenario Matrix

以下は人間回答を反映済みの scenario matrix である。

### SCN-BTP-001 NPC ペルソナ生成フェーズ完了後に本文翻訳フェーズを開始する

- `分類`: 正常系 / 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Job Run で本文翻訳フェーズを開始し、current phase と progress を確認する。
- `受け入れ条件`: persona phase Completed、非 terminal job、active phase run なし、前段参照成立の場合だけ本文翻訳フェーズを開始できる。
- `事前条件`: persona phase Completed job、未完了 job、active phase job、terminal job fixture がある。
- `public_seam_or_api_boundary`: phase start boundary。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。開始条件、開始拒否、current phase、progress。
- `入力開始点`: Job Run UI。
- `主要 outcome`: 本文翻訳フェーズの current phase と progress が表示される。
- `開始操作`: Job Run を開く。
- `入力方法`: 対象 job を選ぶ。
- `主要操作列`: 本文翻訳フェーズ開始を実行し、開始結果を確認する。
- `期待結果`:
  1. 本文翻訳フェーズが current phase として表示される。
  2. progress と phase run 開始結果を確認できる。
  3. 開始条件を満たさない job では開始不可理由が表示される。
- `観測点`: Job Run UI、phase start result、`JOB_PHASE_RUN`。
- `UI-visible outcome`: current phase、progress、開始不可理由。
- `fake_or_stub`: persona completed job fixture、not-ready job fixture、temp DB。

### SCN-BTP-002 本文翻訳入力 snapshot と request summary を固定する

- `分類`: 正常系 / 境界値
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 翻訳対象 field、確定訳語、ジョブ内辞書、persona snapshot、翻訳補助メタデータを同一 phase run の入力として固定する。
- `受け入れ条件`: input summary と digest から、同じ phase run の再開で同じ入力を参照できる。
- `事前条件`: 翻訳 field、job dictionary、persona snapshot、metadata fixture がある。
- `public_seam_or_api_boundary`: body input snapshot boundary。
- `contract_freeze`: あり。target count、dictionary digest、persona digest、metadata digest、prompt digest。
- `入力開始点`: body phase start 後の input collection。
- `主要 outcome`: provider request 前の入力 summary が成立する。
- `開始操作`: body input snapshot を作成する。
- `入力方法`: 前段 phase result と翻訳 field fixture を使う。
- `主要操作列`: target extraction、前段参照、summary 作成を確認する。
- `期待結果`:
  1. 対象 field 件数と対象外理由を確認できる。
  2. 辞書、persona、metadata の参照 summary を確認できる。
  3. 完全一致した辞書 hit は provider request から除外され、部分一致は訳語固定制約として request summary に出る。
  4. provider request へ渡す raw prompt 全文は summary 必須項目ではない。
- `観測点`: body input summary、snapshot digest、`PHASE_RUN_DICTIONARY_ENTRY`、`PHASE_RUN_PERSONA`、`PHASE_RUN_TRANSLATION_FIELD`。
- `UI-visible outcome`: 対象件数、辞書参照件数、persona 参照件数、対象外理由。
- `fake_or_stub`: temp DB、job dictionary fixture、persona snapshot fixture。

### SCN-BTP-003 翻訳レコード種別に応じた provider request を fake transport で実行する

- `分類`: 正常系 / 外部連携
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 翻訳フィールド本文と補助入力から provider request を作り、fake provider の valid response を field result へ写像する。
- `受け入れ条件`: record type、field type、field correlation key、保護要素 digest を失わず provider 境界を通過できる。
- `事前条件`: prompt builder fixture、fixed provider response、batch response fixture がある。
- `public_seam_or_api_boundary`: AIProvider / prompt builder / response adapter boundary。
- `contract_freeze`: あり。fake transport、request summary、response correlation、secret 非露出。
- `入力開始点`: body input snapshot fixture。
- `主要 outcome`: valid response が訳文候補と保護要素検証対象へ変換される。
- `開始操作`: body translation provider adapter を実行する。
- `入力方法`: provider setting、body input snapshot、fixed response を渡す。
- `主要操作列`: request mapping、fake provider execution、response adapter output を確認する。
- `期待結果`:
  1. paid real API を呼ばずに provider 実行を観測できる。
  2. response の field correlation key が `JOB_TRANSLATION_FIELD` と対応する。
  3. `Job Setup` の本文翻訳用 provider、model、execution mode を使った要約を確認できる。
  4. request unit count と output count の要約を確認できる。
- `観測点`: fake transport log、adapter output、provider execution summary。
- `UI-visible outcome`: なし。user-facing 表示は SCN-BTP-006 に統合する。
- `fake_or_stub`: fake provider、fixed response fixture、batch response fixture。

### SCN-BTP-004 保護要素検証後に訳文と出力ステータスを保存する

- `分類`: 正常系 / 永続化
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider 応答の訳文を保護要素検証し、訳文、出力ステータス、検証結果を field 単位で保存する。
- `受け入れ条件`: 保護要素検証を通過した field だけが output-ready 相当へ進み、保存結果を phase run から追跡できる。
- `事前条件`: valid translation result、protection pass / fail、save failure fixture がある。
- `public_seam_or_api_boundary`: body field result persistence boundary。
- `contract_freeze`: あり。field correlation、protection validation、field result persistence、partial state reject。
- `入力開始点`: provider response adapter output。
- `主要 outcome`: `JOB_TRANSLATION_FIELD` と `PHASE_RUN_TRANSLATION_FIELD` が整合する。
- `開始操作`: field result 保存を実行する。
- `入力方法`: 訳文、出力ステータス候補、保護要素検証結果、phase run ID を渡す。
- `主要操作列`: 検証、保存、phase link、result summary を確認する。
- `期待結果`:
  1. 訳文、出力ステータス、保護要素検証結果が同一 field に対応する。
  2. 保存失敗または検証失敗は Completed にならない。
  3. 後続 output artifact に必要な出力ステータス語彙だけを保持する。
  4. 後続 output artifact が参照できる field summary が成立する。
- `観測点`: `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase result、validation result。
- `UI-visible outcome`: translated count、validation pass / fail count、output status summary。
- `fake_or_stub`: temp DB、protection validator fixture、save failure injection。

### SCN-BTP-005 保護要素検証失敗を成功扱いにしない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider 応答が埋め込み要素を欠落、改変、重複、順序違い、余分追加した場合に成功保存しない。
- `受け入れ条件`: 不正な訳文は output-ready にならず、再試行可否と検証差分 summary を確認できる。
- `事前条件`: missing、modified、duplicate、reordered、extra の protection failure fixture がある。
- `public_seam_or_api_boundary`: protection validation boundary。
- `contract_freeze`: あり。protected element invariant、success save rejection、redacted diff summary。
- `入力開始点`: provider response fixture。
- `主要 outcome`: 保護要素不整合の field が成功扱いにならない。
- `開始操作`: protection validation を実行する。
- `入力方法`: 原文保護要素と provider 応答を渡す。
- `主要操作列`: validation result、field status、phase state を確認する。
- `期待結果`:
  1. 保護要素不一致は validation failed として観測できる。
  2. 失敗した訳文は保存前に拒否され、successful translated field として保存されない。
  3. retryable failure として再試行可否を確認できる。
- `観測点`: validation result、field state、phase result、Job Run UI。
- `UI-visible outcome`: 保護要素検証失敗、再試行可否、失敗件数。
- `fake_or_stub`: invalid response fixture、protection failure fixture、temp DB。

### SCN-BTP-006 Job Run で本文翻訳フェーズ result と操作可否を確認する

- `分類`: 表示系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: ユーザーが phase、progress、訳文、出力ステータス、保護要素検証結果、次操作可否を Job Run で確認する。
- `受け入れ条件`: running、paused、recoverable failed、completed、empty completed、validation failed、canceled の状態差分が見える。
- `事前条件`: phase summary、success result、failure result、validation failure result、empty result fixture がある。
- `public_seam_or_api_boundary`: Job Run body phase summary boundary。
- `contract_freeze`: あり。表示項目、button enablement、state variants。
- `入力開始点`: Job Run UI。
- `主要 outcome`: phase result と次操作可否を UI で判断できる。
- `開始操作`: Job Run を開く。
- `入力方法`: 対象 job を選ぶ。
- `主要操作列`: result summary、field summary、error summary、次操作可否を確認する。
- `期待結果`:
  1. current phase、phase state、progress、target count、translated count が表示される。
  2. 訳文、出力ステータス、保護要素検証結果を field 単位または summary で確認できる。
  3. 本文翻訳対象 0 件は completed として表示され、成果物出力へ進めることを確認できる。
  4. error summary は secret、API key 平文、復号可能値を含まない。
- `観測点`: Job Run UI、gateway response、button enablement。
- `UI-visible outcome`: current phase、progress、phase result、field result、次操作可否。
- `fake_or_stub`: mocked gateway result fixture。
- `責務境界メモ`: 詳細な UI 契約は `./ui-design.md` を正本にする。

### SCN-BTP-007 provider 失敗、応答不正、保存失敗を成功扱いにしない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider failure、invalid response、correlation error、save failure で訳文と出力ステータスを成功公開しない。
- `受け入れ条件`: error kind、retryable flag、phase state、progress、field count を確認できる。
- `事前条件`: provider failure、invalid response、batch partial response、save failure fixture がある。
- `public_seam_or_api_boundary`: body phase execution boundary。
- `contract_freeze`: あり。暗黙 fallback なし、invalid response reject、correlation reject、partial persistence reject。
- `入力開始点`: failure injection fixture。
- `主要 outcome`: failed target は successful translated field として保存されず、再試行可否を観測できる。
- `開始操作`: 本文翻訳フェーズを実行する。
- `入力方法`: failure fixture を使う。
- `主要操作列`: provider request、response validation、persistence failure、phase result を確認する。
- `期待結果`:
  1. 別 provider へ暗黙 fallback しない。
  2. invalid response は訳文として保存されない。
  3. provider response と field の対応不能は success にならない。
  4. 部分失敗では成功済み field result を保持し、phase 全体を回復可能な失敗として表示する。
  5. secret、API key 平文、復号可能値は表示またはログに出ない。
- `観測点`: phase result、error kind、row count、structured log。
- `UI-visible outcome`: 失敗理由、再試行可否、後続 output readiness 不可理由。
- `fake_or_stub`: fake transport、invalid response fixture、save failure injection、temp DB。

### SCN-BTP-008 retry、再開、開始再送で重複作成しない

- `分類`: 回復系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 同じ body phase run で開始再送、resume、retry を扱い、field result と phase link を重複作成しない。
- `受け入れ条件`: 同じ `JOB_PHASE_RUN` を継続し、成功済み field と未処理または失敗 field を区別する。
- `事前条件`: active phase run、paused phase run、recoverable failed phase run、partial success fixture がある。
- `public_seam_or_api_boundary`: phase resume / retry boundary。
- `contract_freeze`: あり。same phase run reuse、target snapshot stability、duplicate field guard。
- `入力開始点`: existing phase run fixture。
- `主要 outcome`: phase run と field link が二重作成されない。
- `開始操作`: 開始再送、再開、リトライを実行する。
- `入力方法`: 同一 job / same phase type fixture を使う。
- `主要操作列`: phase run ID、field count、未処理 count、latest error、progress を確認する。
- `期待結果`:
  1. phase run ID は同じである。
  2. `JOB_TRANSLATION_FIELD` と `PHASE_RUN_TRANSLATION_FIELD` は重複しない。
  3. latest error と progress が更新される。
- `観測点`: phase run ID、row count、progress、latest error。
- `UI-visible outcome`: 再開またはリトライ結果。
- `fake_or_stub`: phase run fixture、partial success fixture、temp DB。

### SCN-BTP-009 Paused からの cancel または terminal job には body translation result を後書きしない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: `Paused` から cancel した後、または Failed、Completed、Canceled の後に provider response、validation result、save request が到着しても結果を後書きしない。
- `受け入れ条件`: terminal job の `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase state は変化しない。
- `事前条件`: Running body phase、Paused body phase、terminal job fixture、late response fixture がある。
- `public_seam_or_api_boundary`: terminal job guard boundary。
- `contract_freeze`: あり。terminal state guard、late response rejection、state invariant。
- `入力開始点`: pause 後の cancel request、late provider response、body field save、body phase start。
- `主要 outcome`: terminal job の state 不変と拒否理由を確認できる。
- `開始操作`: terminal job で body phase 関連操作を試行する。
- `入力方法`: terminal job fixture を使う。
- `主要操作列`: 操作前後の phase run、field result、readiness を確認する。
- `期待結果`:
  1. Running から直接 cancel できず、Paused からだけ cancel できる。
  2. terminal job では body phase run が作成されない。
  3. 訳文と出力ステータスは後書きされない。
  4. Canceled 後の途中成功結果は output readiness に使われない。
  5. late response の拒否理由を確認できる。
- `観測点`: phase transition result、row count、拒否理由、state snapshot。
- `UI-visible outcome`: cancel 可否、terminal job の開始不可理由、output readiness 不可理由。
- `fake_or_stub`: terminal job fixture、late response fixture、temp DB。

### SCN-BTP-010 本文翻訳結果 summary を後続成果物出力へ渡せる

- `分類`: 後続フェーズ境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 本文翻訳フェーズ完了後に、後続の translation-output-artifact が参照できる readiness summary を作る。
- `受け入れ条件`: body phase Completed で翻訳ジョブ全体が `Completed` になり、訳文、出力ステータス、保護要素検証結果が field 単位で整合している時だけ readiness が成立する。
- `事前条件`: completed body phase result、failed body phase result、status inconsistency fixture がある。
- `public_seam_or_api_boundary`: output artifact readiness boundary。
- `contract_freeze`: あり。body phase completion requirement、field status consistency、readiness rejection。
- `入力開始点`: output readiness query または output phase start 試行。
- `主要 outcome`: 後続 phase の開始可否と拒否理由を確認できる。
- `開始操作`: output readiness を確認する。
- `入力方法`: body phase state fixture と field result fixture を使う。
- `主要操作列`: readiness query、開始拒否理由、field status summary を確認する。
- `期待結果`:
  1. body phase Completed かつ field result 整合時に job-level `Completed` になり、readiness が true になる。
  2. 未完了、失敗中、status 不整合では output artifact は開始できない。
  3. 本文翻訳対象 0 件でも body phase と job-level `Completed` になり、単語翻訳結果を成果物出力へ渡せる。
- `観測点`: readiness result、`JOB_TRANSLATION_FIELD`、phase result、output status summary。
- `UI-visible outcome`: output readiness、拒否理由、result summary。
- `fake_or_stub`: phase state fixture、field result fixture、temp DB。

### SCN-BTP-011 監査要約と secret 非露出を確認する

- `分類`: セキュリティ / 観測性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: UI、error summary、structured log、debug log、fake transport log が secret を露出しない。
- `受け入れ条件`: provider、model、execution mode、credential ref、input count、output count、digest、error kind は観測できる。原文と訳文はローカル表示できるが、secret、API key 平文、復号可能値は出ない。
- `事前条件`: fake secret store、fake transport、structured log capture、debug log capture、redaction assertion fixture がある。
- `public_seam_or_api_boundary`: runtime log / audit summary boundary。
- `contract_freeze`: あり。secret 非露出、credential_ref、原文と訳文のローカル表示許容。
- `入力開始点`: phase execution fixture。
- `主要 outcome`: 障害調査に必要な要約を確認でき、保存禁止情報は出ない。
- `開始操作`: 本文翻訳フェーズを実行する。
- `入力方法`: provider / model / credential ref / prompt digest fixture を使う。
- `主要操作列`: success、failure、redaction assertion を確認する。
- `期待結果`:
  1. API key、secret 本体、復号可能値は表示またはログに出ない。
  2. 原文と訳文がローカルで見えること自体は失敗扱いにしない。
  3. provider、model、execution mode、input count、output count、prompt digest を確認できる。
- `観測点`: structured log、debug log capture、Job Run summary、fake secret store assertion、fake transport log。
- `UI-visible outcome`: provider / model / credential 参照状態、secret を含まない error summary。
- `fake_or_stub`: fake secret store、fake transport、log capture。

## Acceptance Checks

- `REQ-BTP-001`: `SCN-BTP-001`, `SCN-BTP-009`
- `REQ-BTP-002`: `SCN-BTP-002`, `SCN-BTP-003`, `SCN-BTP-011`
- `REQ-BTP-003`: `SCN-BTP-003`, `SCN-BTP-007`, `SCN-BTP-008`
- `REQ-BTP-004`: `SCN-BTP-004`, `SCN-BTP-005`, `SCN-BTP-010`
- `REQ-BTP-005`: `SCN-BTP-006`, `SCN-BTP-010`
- `REQ-BTP-006`: `SCN-BTP-006`, `SCN-BTP-008`, `SCN-BTP-009`
- `REQ-BTP-007`: `SCN-BTP-007`, `SCN-BTP-011`
- `REQ-BTP-008`: `SCN-BTP-002`, `SCN-BTP-006`, `SCN-BTP-010`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/body-translation-phase/scenario-design.md --coverage docs/exec-plans/completed/body-translation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/body-translation-phase/scenario-design.candidate-coverage.json --report-out docs/exec-plans/completed/body-translation-phase/scenario-design.requirement-gate.md --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Decisions

- なし。
- `scenario-design.questions.md` に人間回答 10 件を記録済みである。
- `implementation-scope.md` は人間設計レビュー未承認のため未作成である。
