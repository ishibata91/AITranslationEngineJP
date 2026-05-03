# Scenario Design: translation-output-artifact

- `skill`: scenario-design
- `status`: ready-human-review
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-output-artifact.md`
- `topic_abbrev`: `TOA`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - Output Review は body phase Completed かつ job-level `Completed` の job だけを出力候補として扱う。
  - output readiness は訳文、出力ステータス、保護要素検証結果、field result 整合が成立する時だけ true にする。
  - Canceled、未完了、失敗中、不整合 job の途中成功結果を xTranslator 成果物へ使わない。
  - result summary は訳文件数、出力ステータス分布、入力ファイル出自、保護要素検証結果、出力可能性を確認できる。
  - diff preview は translation unit 単位で Source、Dest、Status、row 反映内容を確認できる。
  - xTranslator 互換 XML は `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を row 単位で再構成する。
  - 内部 `cached` は xTranslator `Status=1` へ写像し、辞書置換である事実は内部観測情報として分離する。
  - XML は UTF-8 で parse 可能であり、対象ゲームに対応した root element と `<String>` 子要素を持つ。
  - 再出力は同じ job に対して重複 artifact または重複 row を作らない。
  - 出力処理は AI provider、network、API key、secret store を必須経路にしない。
  - UI、error summary、structured log、debug log、runtime event へ secret、API key 平文、復号可能値、provider raw payload、過剰な本文全文を出さない。
- `non_goals`:
  - 本文翻訳フェーズの再翻訳、field result 編集、provider 実行、Job Run UI の再設計は扱わない。
  - xTranslator 本体を自動操作する実装、real xTranslator import smoke の必須化は扱わない。
  - docs 正本、プロダクトコード、プロダクトテスト、`.codex`、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-...` を一意 key として扱う。

- `candidate_count`: 54
- `adopted`: 10
- `merged`: 44
- `rejected`: 0
- `needs_human_decision`: 0
- `unresolved_conflicts`: 0
- `questionnaire`: `./scenario-design.questions.md`

candidate artifact 内の human decision candidate は、正本から導出できるものを固定した。
正本だけでは値を固定できないものは、詳細要求の `deferred` と人間レビュー論点へ分けた。
候補生成 agent は再起動していない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは sidecar JSON に分離した。
`needs_human_decision` は 0 件であり、この scenario matrix は人間設計レビューへ進められる。
`deferred` は人間レビューまたは implementation-scope 作成時に再確認する。

### `REQ-TOA-001` Output Review で completed job と result summary を確認する

- `source_requirement`: `tasks/usecases/translation-output-artifact.yaml` は completed job 一覧、result summary、diff preview、output artifact 状態を completion criteria にする。
- `requirement_kind`: display
- `needs_human_decision`: なし。
- `fixed_decisions`: Output Review で completed job、readiness、拒否理由、訳文件数、出力ステータス分布、入力ファイル出自を確認する。

### `REQ-TOA-002` output readiness と禁止遷移を守る

- `source_requirement`: body phase 完了成果物は、body phase Completed、job-level `Completed`、field result と output status の整合時だけ output readiness を true にする。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: 未完了、失敗中、`Canceled`、field result 不整合、status 不整合では artifact 生成を開始しない。

### `REQ-TOA-003` xTranslator row と XML を生成する

- `source_requirement`: `docs/spec.md` は xTranslator 互換形式の `.xml` と row 必須列の再構成を求め、`docs/er.md` は `XTRANSLATOR_OUTPUT_ROW` を出力行モデルにする。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし。
- `fixed_decisions`: row 必須列、UTF-8 XML、Skyrim root element、local XML parser 検査、real xTranslator 非必須を固定する。

### `REQ-TOA-004` 内部 output status を xTranslator `Status` へ写像する

- `source_requirement`: `docs/spec.md` は内部 `cached` を xTranslator `Status=1` へ写像し、辞書置換情報を内部観測情報として保持するとする。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし。
- `fixed_decisions`: `cached` は `Status=1` にする。未定義 status は成功成果物へ混入しない。`cached` 以外の詳細写像表は review 後に再確認する。

### `REQ-TOA-005` diff preview と再出力状態を確認する

- `source_requirement`: `tasks/usecases/translation-output-artifact.yaml` は translation unit 単位の差分、再生成操作、再出力状態を completion criteria にする。
- `requirement_kind`: display
- `needs_human_decision`: なし。
- `fixed_decisions`: diff preview は現行 field result と生成済み row の差分確認を扱う。本文再翻訳は本 task の non-goal とし、再出力導線へ分離する。

### `REQ-TOA-006` output artifact の再出力と失敗回復を扱う

- `source_requirement`: `docs/er.md` と ER 図は `TRANSLATION_ARTIFACT.translation_job_id` を job 単位の成果物として扱い、`XTRANSLATOR_OUTPUT_ROW.job_translation_field_id` を一意にする。
- `requirement_kind`: persistence
- `needs_human_decision`: なし。
- `fixed_decisions`: 同じ job の再出力は現行 artifact を更新または置換する境界として扱い、重複 row を作らない。履歴の永続化拡張は deferred にする。

### `REQ-TOA-007` 監査要約と redaction を満たす

- `source_requirement`: `docs/spec.md` は API key の暗号化保存を求め、body phase 設計は secret と raw payload を UI、summary、log に出さない。
- `requirement_kind`: security
- `needs_human_decision`: なし。
- `fixed_decisions`: 出力処理は provider と secret store を呼ばない。監査要約は ID、件数、digest、error kind、status を中心にし、本文全文と secret を重複保存しない。

### `REQ-TOA-008` 本文翻訳対象 0 件の completed job を出力可能にする

- `source_requirement`: body phase completed 成果物は本文翻訳対象 0 件を正常完了にし、単語だけの plugin でも output readiness を成立させる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし。
- `fixed_decisions`: target count 0 は出力不可理由にしない。row count 0 の XML でも artifact summary と output status を確認できる。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答扱いの項目は残さない。
人間レビューで確認する論点は `Open Questions` に分離する。

## Risks

- `implementation_risks`:
  - output readiness の所有境界を body phase と output artifact で重複実装すると、開始可否と拒否理由が食い違う。
  - `TRANSLATION_ARTIFACT.translation_job_id` の一意性を無視すると、再出力で artifact と row が重複する。
  - `cached` と xTranslator `Status=1` を同一概念にすると、辞書置換の内部観測情報が失われる。
  - XML root、escaping、危険値分類を serializer だけへ押し込むと、UI summary と検証結果がずれる。
  - output task で provider 履歴を表示する時、secret や raw payload を誤って再露出する。
- `test_data_risks`:
  - completed、未完了、`Canceled`、field result 不整合、target count 0 の job fixture を分ける必要がある。
  - cached、translated、unknown status、欠損 row、重複 row、XML 特殊文字、長文、root element 未確定の fixture が必要になる。
  - readonly path、保存失敗、XML parse 失敗、再出力再送の fixture が必要になる。
  - UI、structured log、debug log、runtime event で redaction を別々に確認する必要がある。

## Rules

- ケース ID は `SCN-TOA-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` は 0 件である。
- 未解決 conflict は 0 件である。
- paid real AI API と real xTranslator 自動操作を前提にしない。

## Scenario Matrix

### SCN-TOA-001 Output Review で completed job と result summary を確認する

- `分類`: 正常系 / 表示
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 完了済み翻訳ジョブを選び、出力準備状態と結果要約を確認する。
- `受け入れ条件`: body phase Completed、job-level `Completed`、field result 整合、output status 整合の job だけが出力候補になる。
- `事前条件`: completed job、未完了 job、`Canceled` job、field result 不整合 job の fixture がある。
- `public_seam_or_api_boundary`: Output Review summary query。詳細 API 名は implementation-scope で固定する。
- `contract_freeze`: あり。completed job list、output readiness、拒否理由、result summary。
- `入力開始点`: Output Review UI。
- `主要 outcome`: completed job と result summary が確認できる。
- `開始操作`: Output Review を開く。
- `入力方法`: job list から対象 job を選ぶ。
- `主要操作列`: job を選択し、readiness、件数、出力ステータス分布、入力出自を確認する。
- `手順`:
  1. Output Review を開く。
  2. completed job を選ぶ。
  3. result summary を確認する。
- `期待結果`:
  1. completed job だけが出力開始可能として表示される。
  2. 訳文件数、出力ステータス分布、保護要素検証結果、入力ファイル出自を確認できる。
  3. 出力不可 job では拒否理由が表示される。
- `観測点`: Output Review UI、readiness result、result summary、`JOB_TRANSLATION_FIELD`。
- `UI-visible outcome`: completed job list、readiness、拒否理由、result summary。
- `fake_or_stub`: completed job fixture、not-ready job fixture、temp DB。
- `責務境界メモ`: job 完了判定は body phase の output readiness に依存し、Output Review は再翻訳を起動しない。

### SCN-TOA-002 未完了または不整合 job から artifact を作らない

- `分類`: 主要失敗系 / 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: output readiness が false の job で artifact 生成を拒否する。
- `受け入れ条件`: 未完了、失敗中、`Canceled`、field result 不整合、status 不整合では `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` を成功状態で作らない。
- `事前条件`: not completed、recoverable failed、failed、canceled、status mismatch の job fixture がある。
- `public_seam_or_api_boundary`: output artifact start boundary。
- `contract_freeze`: あり。readiness false、rejection reason、no artifact creation。
- `入力開始点`: output artifact 生成開始試行。
- `主要 outcome`: 出力開始が拒否され、理由を観測できる。
- `開始操作`: output artifact 生成を実行する。
- `入力方法`: 出力不可 job ID を渡す。
- `主要操作列`: readiness 判定、拒否理由作成、artifact 未生成確認を行う。
- `手順`:
  1. 出力不可 job で生成開始を試す。
  2. 拒否理由を確認する。
- `期待結果`:
  1. artifact 生成は開始されない。
  2. 途中成功結果は output readiness に使われない。
  3. rejection reason は本文全文や secret を含まない。
- `観測点`: readiness result、rejection reason、`TRANSLATION_ARTIFACT` 件数、`XTRANSLATOR_OUTPUT_ROW` 件数。
- `UI-visible outcome`: 出力不可理由、生成 button disabled。
- `fake_or_stub`: invalid job state fixture、temp DB。
- `責務境界メモ`: 状態不変条件は backend 側の開始境界で保証する。

### SCN-TOA-003 xTranslator row 必須列と status mapping を生成する

- `分類`: 正常系 / 互換性
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: output-ready field から xTranslator row 必須列を再構成し、内部 `cached` を `Status=1` へ写像する。
- `受け入れ条件`: 各 `XTRANSLATOR_OUTPUT_ROW` は 1 つの `JOB_TRANSLATION_FIELD` に対応し、`EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を持つ。
- `事前条件`: translated field、cached field、unknown status field、欠損識別子 field の fixture がある。
- `public_seam_or_api_boundary`: xTranslator row builder boundary。
- `contract_freeze`: あり。row required columns、one field one row、cached to `Status=1`。
- `入力開始点`: output-ready field set。
- `主要 outcome`: row set と status mapping summary が成立する。
- `開始操作`: row builder を実行する。
- `入力方法`: completed job の field result summary を渡す。
- `主要操作列`: field identity 読み取り、status mapping、row validation、summary 作成を確認する。
- `手順`:
  1. row builder を実行する。
  2. row 必須列と status mapping を確認する。
- `期待結果`:
  1. cached field は row `Status=1` になる。
  2. 辞書置換である事実は xTranslator `Status` とは別の内部 summary に残る。
  3. 未定義 status や必須列欠損は成功 row に混入しない。
- `観測点`: `XTRANSLATOR_OUTPUT_ROW`、row build result、status mapping summary。
- `UI-visible outcome`: row count、cached count、mapping error count。
- `fake_or_stub`: cached field fixture、unknown status fixture、missing identifier fixture。
- `責務境界メモ`: `cached` 以外の詳細 mapping は review 後に再確認するが、unknown status 成功混入は禁止する。

### SCN-TOA-004 UTF-8 の xTranslator XML を生成して parser で検査する

- `分類`: 正常系 / 外部連携
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: xTranslator 互換 XML を UTF-8 として出力し、root element と `<String>` row を検査する。
- `受け入れ条件`: XML は local parser で parse でき、対象ゲームに対応する root element、`<String>` 子要素、必須 tag を持つ。
- `事前条件`: Skyrim SE target fixture、Skyrim LE target fixture、XML escaping fixture、readonly path fixture がある。
- `public_seam_or_api_boundary`: XML serializer / filesystem adapter boundary。
- `contract_freeze`: あり。UTF-8、Skyrim root mapping、required tag、file write result。
- `入力開始点`: validated row set と output file path。
- `主要 outcome`: XML file と compatibility summary が生成される。
- `開始操作`: XML 出力を実行する。
- `入力方法`: row set、target game、出力 path を渡す。
- `主要操作列`: XML serialization、file write、parser check、artifact state 更新を確認する。
- `手順`:
  1. XML serializer を実行する。
  2. 生成 XML を local parser で検査する。
- `期待結果`:
  1. XML は UTF-8 として parse できる。
  2. Skyrim SE は `SSETranslator`、Skyrim LE は `TESVTranslator` を root element にできる。
  3. XML 特殊文字と日本語 text は壊れない。
  4. real xTranslator 起動は必須にしない。
- `観測点`: generated XML、XML parser result、file write result、artifact state。
- `UI-visible outcome`: XML 出力結果、file path、compatibility summary。
- `fake_or_stub`: temp filesystem、XML parser、target game fixture。
- `責務境界メモ`: filesystem concrete は adapter 境界に閉じる。

### SCN-TOA-005 diff preview と再出力必要状態を確認する

- `分類`: 正常系 / 代替成功
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 現行 field result と生成済み artifact row の差分を translation unit 単位で確認する。
- `受け入れ条件`: diff preview は Source、Dest、Status、row 反映内容、stale reason、再出力可否を表示する。
- `事前条件`: 生成済み artifact、変更なし field、変更あり field、欠落 row の fixture がある。
- `public_seam_or_api_boundary`: diff preview query。
- `contract_freeze`: あり。diff row identity、stale reason、re-output enablement。
- `入力開始点`: Output Review UI。
- `主要 outcome`: 差分と再出力導線を判断できる。
- `開始操作`: diff preview を開く。
- `入力方法`: selected job と selected artifact を選ぶ。
- `主要操作列`: diff preview を開き、unit を選び、再出力状態を確認する。
- `手順`:
  1. 生成済み artifact を持つ job を開く。
  2. diff preview を開く。
  3. 再出力可否を確認する。
- `期待結果`:
  1. translation unit 単位の差分を確認できる。
  2. 参照不能または古い row は正しい preview として表示されない。
  3. 本文再翻訳ではなく artifact 再出力導線として扱われる。
- `観測点`: diff preview UI、row digest、field digest、re-output state。
- `UI-visible outcome`: diff row、stale badge、re-output button enablement。
- `fake_or_stub`: stale artifact fixture、missing row fixture、temp DB。
- `責務境界メモ`: translation unit の再翻訳は本 task の非対象とする。

### SCN-TOA-006 同じ job の再出力で artifact と row を重複作成しない

- `分類`: 正常系 / 冪等性
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 同じ completed job と同じ field set で再出力しても、現行 artifact と row count が破損しない。
- `受け入れ条件`: `TRANSLATION_ARTIFACT.translation_job_id` は job 単位の現行 artifact を指し、`XTRANSLATOR_OUTPUT_ROW.job_translation_field_id` は同じ field に重複 row を作らない。
- `事前条件`: 生成済み artifact、同一 field set、変更あり field set、失敗 artifact summary の fixture がある。
- `public_seam_or_api_boundary`: re-output command boundary。
- `contract_freeze`: あり。same job re-output、unique artifact per job、unique row per field。
- `入力開始点`: re-output command。
- `主要 outcome`: 再出力後も artifact summary と row count が一貫する。
- `開始操作`: 再出力を実行する。
- `入力方法`: completed job ID と出力条件を渡す。
- `主要操作列`: existing artifact 読み取り、row rebuild、replace/update、summary 更新を確認する。
- `手順`:
  1. 生成済み artifact がある job で再出力を実行する。
  2. artifact と row 件数を確認する。
- `期待結果`:
  1. 同一 field の row が重複しない。
  2. 成功済み artifact と失敗 summary を区別できる。
  3. result summary は再出力前後で破損しない。
- `観測点`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、row count、artifact status。
- `UI-visible outcome`: 再出力結果、row count、generated_at、status。
- `fake_or_stub`: existing artifact fixture、retry fixture、temp DB。
- `責務境界メモ`: 複数世代の artifact 履歴は現行 ER からは導出せず、将来拡張として deferred にする。

### SCN-TOA-007 生成失敗を成功成果物として扱わず回復可能にする

- `分類`: 主要失敗系 / 回復
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: row validation、XML serialization、file write、artifact 保存失敗を success として公開しない。
- `受け入れ条件`: 生成完了条件を満たさない場合、artifact は成功状態にならず、失敗理由と再出力可否を返す。
- `事前条件`: missing row、invalid status、XML serialization failure、readonly path、DB save failure の fixture がある。
- `public_seam_or_api_boundary`: artifact generation transaction boundary。
- `contract_freeze`: あり。failure kind、no success publish、retryable flag。
- `入力開始点`: artifact generation command。
- `主要 outcome`: 失敗 summary と再出力可能性を確認できる。
- `開始操作`: 失敗 fixture で artifact 生成を実行する。
- `入力方法`: invalid row set または save failure injection を渡す。
- `主要操作列`: validation、serialization、save、failure mapping を確認する。
- `手順`:
  1. 失敗 fixture で生成を試す。
  2. artifact state と error summary を確認する。
- `期待結果`:
  1. 不完全な XML または row は成功 artifact にならない。
  2. 保存失敗後に成功 XML として公開されない。
  3. retryable flag と failed stage を確認できる。
- `観測点`: artifact status、row count、error summary、structured log。
- `UI-visible outcome`: 失敗理由、再出力可否、保存済み件数。
- `fake_or_stub`: save failure injection、readonly path fixture、invalid XML fixture。
- `責務境界メモ`: 部分保存を残すか破棄するかは実装時 transaction 方針で確認するが、success 公開は禁止する。

### SCN-TOA-008 本文翻訳対象 0 件の completed job でも出力結果を確認する

- `分類`: 境界値 / 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 本文翻訳対象 0 件で Completed になった job を output artifact へ進める。
- `受け入れ条件`: target count 0 は output readiness false の理由にならず、row count 0 の artifact summary を確認できる。
- `事前条件`: body translation target count 0、term result あり、field result 0 の completed job fixture がある。
- `public_seam_or_api_boundary`: output readiness / row builder boundary。
- `contract_freeze`: あり。target count 0、row count 0、artifact summary。
- `入力開始点`: output artifact generation command。
- `主要 outcome`: row count 0 でも成果物出力処理が正常に終了する。
- `開始操作`: target count 0 job で出力を実行する。
- `入力方法`: target count 0 の completed job ID を渡す。
- `主要操作列`: readiness、row builder、XML summary、artifact state を確認する。
- `手順`:
  1. target count 0 job を選ぶ。
  2. artifact 出力を実行する。
- `期待結果`:
  1. target count 0 は出力不可理由にならない。
  2. row count 0 と skipped count を result summary で確認できる。
  3. provider 未実行の事実は監査要約で確認できる。
- `観測点`: readiness result、row count、artifact status、result summary。
- `UI-visible outcome`: target count 0、row count 0、出力状態。
- `fake_or_stub`: zero target completed job fixture、temp DB。
- `責務境界メモ`: 単語翻訳フェーズ自体の再設計は扱わない。

### SCN-TOA-009 xTranslator 互換上の危険値を検出して summary に出す

- `分類`: 境界値 / 互換性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: xTranslator 互換上の危険値を検出し、無警告で成功確認にしない。
- `受け入れ条件`: 文字列サイズ上限超過、翻訳非推奨 field、RACE 先頭スペース、末尾スペースなどを compatibility summary に出す。
- `事前条件`: large text、WOOP field、RACE leading space、trailing space の fixture がある。
- `public_seam_or_api_boundary`: compatibility validator boundary。
- `contract_freeze`: あり。detection result、severity、summary。
- `入力開始点`: row validation または diff preview。
- `主要 outcome`: 危険値の種別、対象 row、互換性への影響を確認できる。
- `開始操作`: compatibility validator を実行する。
- `入力方法`: boundary value row set を渡す。
- `主要操作列`: detection、summary、success blocking または warning classification を確認する。
- `手順`:
  1. 危険値 fixture を検証する。
  2. compatibility summary を確認する。
- `期待結果`:
  1. 危険値は検出される。
  2. 致命的な構造違反は success artifact にならない。
  3. warning と reject の詳細分類は人間レビュー後に再確認できる。
- `観測点`: compatibility validation result、result summary、diff preview。
- `UI-visible outcome`: compatibility warning、reject reason、affected row count。
- `fake_or_stub`: boundary value fixture、XML parser。
- `責務境界メモ`: 危険値分類の閾値は `deferred` とし、検出と表示を先に固定する。

### SCN-TOA-010 出力処理で provider、network、secret を使わず redaction を守る

- `分類`: セキュリティ / 監査
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: completed result からの成果物出力で AI provider、network、secret store を呼ばず、保存禁止情報を露出しない。
- `受け入れ条件`: output artifact の success、failure、re-output で provider call count は 0 になり、UI、summary、log に secret と raw payload が出ない。
- `事前条件`: fake provider fail-on-call、fake secret store fail-on-call、log capture、completed result fixture がある。
- `public_seam_or_api_boundary`: output artifact service / XML adapter boundary。
- `contract_freeze`: あり。no provider call、no secret resolution、redacted summary。
- `入力開始点`: output artifact generation command。
- `主要 outcome`: provider 非接続と redaction を確認できる。
- `開始操作`: 出力処理を実行する。
- `入力方法`: provider 履歴を持つ completed job fixture を渡す。
- `主要操作列`: provider call count、secret store call count、summary、log capture を確認する。
- `手順`:
  1. fake provider と fake secret store を fail-on-call にする。
  2. 出力処理を実行する。
  3. UI summary と log capture を確認する。
- `期待結果`:
  1. provider、network、secret store は呼ばれない。
  2. secret、API key、復号可能値、provider raw payload は出力されない。
  3. 監査要約は artifact id、operation kind、row count、digest、error kind、status を確認できる。
- `観測点`: fake provider call count、fake secret store assertion、UI error summary、structured log、debug log。
- `UI-visible outcome`: redacted error summary、operation summary。
- `fake_or_stub`: fake provider fail-on-call、fake secret store、log capture。
- `責務境界メモ`: output task は provider 実行履歴を読む場合も secret 本体を再解決しない。

## Acceptance Checks

- `REQ-TOA-001`: `SCN-TOA-001`, `SCN-TOA-005`
- `REQ-TOA-002`: `SCN-TOA-001`, `SCN-TOA-002`
- `REQ-TOA-003`: `SCN-TOA-003`, `SCN-TOA-004`, `SCN-TOA-009`
- `REQ-TOA-004`: `SCN-TOA-003`, `SCN-TOA-007`
- `REQ-TOA-005`: `SCN-TOA-005`, `SCN-TOA-006`
- `REQ-TOA-006`: `SCN-TOA-006`, `SCN-TOA-007`
- `REQ-TOA-007`: `SCN-TOA-010`
- `REQ-TOA-008`: `SCN-TOA-008`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-output-artifact/scenario-design.md --coverage docs/exec-plans/active/translation-output-artifact/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-output-artifact/scenario-design.candidate-coverage.json --report-out docs/exec-plans/active/translation-output-artifact/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/translation-output-artifact/scenario-design.questions.md`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Questions

- 人間レビュー論点: `cached` 以外の内部 output status を xTranslator `Status=0..4` のどれへ写像するか。
- 人間レビュー論点: compatibility validator が検出する危険値を reject、warning、許容のどれへ分類するか。
- 人間レビュー論点: 再出力履歴を現行 ER の 1 job 1 artifact に留めるか、将来 revision 履歴へ拡張するか。
- 人間レビュー論点: source file path、`Source`、`Dest`、diff preview 本文の UI 表示とログ保存の伏せ字範囲をどこまでにするか。
- 人間レビュー論点: real xTranslator import smoke を最終検証の任意手動確認に含めるか。
