# Scenario Candidates: body-translation-phase / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`
- `task_artifact_dir`: `docs/exec-plans/completed/body-translation-phase/`
- `target_delta`: 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを参照して翻訳フィールド本文を翻訳し、訳文、出力ステータス、保護要素検証結果を確認できる本文翻訳フェーズを追加する。

## Generator Scope

- `viewpoint`: 失敗入力、参照不能、設定不整合、保存失敗、回復動作を扱う。
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/index.yaml`
  - `../../../../tasks/usecases/body-translation-phase.yaml`
  - `../../../../tasks/usecases/term-translation-phase.yaml`
  - `../../../../tasks/usecases/persona-generation-phase.yaml`
  - `../../../../tasks/usecases/translation-output-artifact.yaml`
  - `../../../spec.md`
  - `../../../er.md`
  - `../../../architecture.md`
  - `../term-translation-phase/scenario-design.md`
  - `../persona-generation-phase/scenario-design.md`
- `excluded_sources`:
  - プロダクトコード、プロダクトテスト、docs 正本更新、最終シナリオ表、採否判断。
- `generation_notes`:
  - 本文翻訳フェーズの failure 候補だけを列挙する。
  - 競合解消、採否、最終 scenario ID 付与は `designer` に残す。
  - paid real AI API を検証前提にしない。

## Candidate Scenarios

### CAND-BTP-001 persona phase 未完了または snapshot 参照不能で開始を拒否する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:21-30` は NPC ペルソナ生成フェーズ完了を precondition とし、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md:92-97` は snapshot 参照成立後だけ body readiness を成立させる。
- `viewpoint`: 参照不能 / 状態不整合
- `candidate scenario id`: `CAND-BTP-001`
- `actor`: Job Run を操作するユーザー。
- `trigger`: persona phase が `Completed` でない、または persona snapshot / `PHASE_RUN_PERSONA` を参照できない job で本文翻訳フェーズ開始を試行する。
- `expected outcome`: body phase run は作成されず、開始不可理由と persona snapshot 参照不能を Job Run で確認できる。
- `observable point`: Job Run UI、phase start result、`JOB_PHASE_RUN` 件数、body input summary。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: body phase の開始ガード候補。persona phase の詳細結果は再設計しない。
- `conflict hint`: lifecycle / state-transition 候補が body phase の開始条件を別に定義する場合、開始拒否理由と phase run 作成有無が競合しうる。

### CAND-BTP-002 確定訳語またはジョブ内辞書の参照不能で開始を拒否する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:10-16` は確定訳語とジョブ内辞書を input にし、`docs/exec-plans/completed/term-translation-phase/scenario-design.md:72-91` は辞書反映と参照不能時の後続 phase 禁止を固定している。
- `viewpoint`: 参照不能 / 設定不整合
- `candidate scenario id`: `CAND-BTP-002`
- `actor`: Job Run を操作するユーザー。
- `trigger`: term phase は完了扱いだが、確定訳語 summary、ジョブ内辞書 snapshot、または `PHASE_RUN_DICTIONARY_ENTRY` の参照が欠損している job で本文翻訳フェーズ開始を試行する。
- `expected outcome`: body phase run は作成されず、辞書参照不能、確定訳語参照不能、または term phase 結果不整合を確認できる。
- `observable point`: phase start result、body input summary、辞書 entry 件数、Job Run UI の開始不可理由。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 空の辞書と参照不能を分ける候補。空辞書が正常かどうかは `designer` が判断する。
- `conflict hint`: CAND-BTP-014 の空対象 / 空辞書扱いと競合しうる。

### CAND-BTP-003 翻訳フィールド識別子欠損で本文翻訳を拒否する

- `source requirement`: `docs/spec.md:41-43` は翻訳単位の `FormID`、`EditorID`、レコード種別、フィールド種別、原文、訳文、出力ステータスを lossless に保持することを求める。`docs/er.md:29-43` は翻訳レコードと翻訳フィールドの識別情報を定義する。
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-BTP-003`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: 翻訳対象 field に `FormID`、`EditorID`、record type、field type、source text のいずれかが欠損している。
- `expected outcome`: 欠損 field は provider request に入らず、訳文と成功ステータスは保存されない。field 欠損理由と識別可能な範囲の record 情報を phase result で確認できる。
- `observable point`: target extraction result、provider request payload、phase result、`JOB_TRANSLATION_FIELD` 件数。
- `related detail requirement type`: `boundary_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: xTranslator output に必要な識別情報を満たせない入力を扱う候補。
- `conflict hint`: actor-goal 候補が部分的に識別可能な field を翻訳継続とする場合、phase state と出力ステータスが競合しうる。

### CAND-BTP-004 翻訳指示構成に必要な field definition がない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:23-28` は翻訳レコード種別に応じた翻訳指示構成を completion criteria にする。`docs/er.md:42-46` は `TRANSLATION_FIELD_DEFINITION` が AI 向け説明と参照要件を保持すると定義する。
- `viewpoint`: 参照不能 / 設定不整合
- `candidate scenario id`: `CAND-BTP-004`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: 翻訳対象の `RecordType + SubrecordType` に対応する `TRANSLATION_FIELD_DEFINITION` がない、または翻訳対象フラグ / 参照要件が矛盾している。
- `expected outcome`: 対象 field は provider request に入らず、翻訳指示構成失敗として確認できる。成功訳文と出力ステータスは保存されない。
- `observable point`: prompt build result、target summary、phase result、Job Run error summary。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: provider 失敗ではなく、prompt / target 構成前の拒否候補。
- `conflict hint`: lifecycle 候補が unsupported field を skipped success として扱う場合、failure と skipped の境界が競合しうる。

### CAND-BTP-005 確定訳語と record type context が一致しない

- `source requirement`: `docs/exec-plans/completed/term-translation-phase/scenario-design.md:72-77` は同一 source term の重複判定 key を record type ごとに一意にする。`docs/spec.md:27-33` は辞書に基づく一貫した訳語と内部観測情報を求める。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-BTP-005`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: body phase の翻訳対象 record type と、参照する確定訳語 / ジョブ内辞書 entry の record type context が一致しない。
- `expected outcome`: 不一致 entry は翻訳指示へ混ぜず、辞書不整合として phase result に出る。誤った訳語を使った訳文は成功扱いにしない。
- `observable point`: dictionary lookup result、prompt input summary、phase result、Job Run error summary。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: 本文翻訳側の辞書参照 contract を固定するための失敗候補。
- `conflict hint`: actor-goal 候補が record type をまたいだ訳語再利用を許す場合、辞書適用範囲が競合しうる。

### CAND-BTP-006 provider 失敗または応答不正で訳文を成功保存しない

- `source requirement`: `docs/spec.md:49-58` は各フェーズの AI 選択、API key 保存、失敗時の暗黙 fallback 不要を定義する。`tasks/usecases/body-translation-phase.yaml:17-28` は訳文、出力ステータス、保護要素検証結果、recoverable failure 表示を output / criteria にする。
- `viewpoint`: 参照不能 / 回復動作
- `candidate scenario id`: `CAND-BTP-006`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: provider timeout、provider error、invalid response、response 欠落、request unit と response の対応不能が発生する。
- `expected outcome`: 失敗 field は成功訳文として保存されず、error kind、retryable flag、phase state、progress を確認できる。別 provider への暗黙 fallback は起きない。
- `observable point`: fake provider result、response validation result、phase result、Job Run error summary、`JOB_TRANSLATION_FIELD` 件数。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: external-integration 候補と統合されうる provider failure 候補。
- `conflict hint`: external-integration 候補が provider failure 粒度を request 単位にする場合、field 単位の失敗表示と競合しうる。

### CAND-BTP-007 保護要素検証失敗で訳文を成功扱いにしない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:17-28` は保護要素検証結果と埋め込み要素保持を求める。`docs/spec.md:41-43` は `<10gold>` などの構造や埋め込み要素を損なわない翻訳を求める。
- `viewpoint`: 失敗入力 / 整合性違反
- `candidate scenario id`: `CAND-BTP-007`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: provider 応答の訳文で、原文にあった埋め込み要素が欠落、改変、重複、順序違い、余分追加になっている。
- `expected outcome`: 訳文は成功ステータスにならず、保護要素検証結果、差分 summary、再試行可否を確認できる。
- `observable point`: protection validation result、phase result、Job Run UI、`JOB_TRANSLATION_FIELD` の訳文 / 出力ステータス。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 本文翻訳フェーズ固有の主要失敗候補。
- `conflict hint`: 正常系候補が provider 応答を即時保存する前提を置く場合、保存前検証の順序が競合しうる。

### CAND-BTP-008 訳文と出力ステータスの保存途中失敗を Completed にしない

- `source requirement`: `docs/er.md:22-25` は翻訳結果と出力ステータスを `JOB_TRANSLATION_FIELD` に保持し、最終 AI 実行情報を `JOB_PHASE_RUN` に保持すると定義する。`docs/er.md:66-69` は phase 再実行と `PHASE_RUN_TRANSLATION_FIELD` を定義する。
- `viewpoint`: 保存失敗
- `candidate scenario id`: `CAND-BTP-008`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: 訳文、出力ステータス、保護要素検証結果、phase link のいずれかの保存で失敗する。
- `expected outcome`: partial state は successful Completed にならず、保存済み / 未保存の件数、retryable flag、再試行対象を確認できる。
- `observable point`: transaction result、`JOB_TRANSLATION_FIELD` 件数、`PHASE_RUN_TRANSLATION_FIELD` 件数、phase result。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 永続化境界の atomicity 候補。実装方式は固定しない。
- `conflict hint`: lifecycle 候補が成功済み field の保持を前提にする場合、all-or-nothing と部分成功維持の扱いが競合しうる。

### CAND-BTP-009 pause、resume、retry、cancel の不許可状態を確認する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:23-28` は pause、resume、retry、cancel の可否と recoverable failure 情報を確認できることを求める。`docs/spec.md:155-199` は Paused、RecoverableFailed、Failed、Canceled の状態を定義する。
- `viewpoint`: 回復動作 / 状態不整合
- `candidate scenario id`: `CAND-BTP-009`
- `actor`: Job Run を操作するユーザー。
- `trigger`: Running 以外で pause、Paused 以外で resume、RecoverableFailed 以外で retry、terminal job で cancel / retry を実行しようとする。
- `expected outcome`: 不許可操作は実行されず、操作不可理由と現在 state を確認できる。phase result と保存済み訳文は変化しない。
- `observable point`: Job Run button enablement、operation result、state snapshot、phase result。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `adoption hint`: UI-visible な回復操作の失敗候補。
- `conflict hint`: state-transition 候補が操作可能 state を別に定義する場合、button enablement と operation result が競合しうる。

### CAND-BTP-010 RecoverableFailed からの retry で成功済み訳文を重複作成しない

- `source requirement`: `docs/er.md:66-69` は同じ `JOB_PHASE_RUN` の状態を戻す再実行と phase target link を定義する。`docs/spec.md:178-196` は RecoverableFailed から Running への再開 / リトライを定義する。
- `viewpoint`: 回復動作 / 保存失敗後の再試行
- `candidate scenario id`: `CAND-BTP-010`
- `actor`: Job Run を操作するユーザー。
- `trigger`: 一部 field が成功し、別 field が provider 失敗または保護要素検証失敗で RecoverableFailed になった後に retry する。
- `expected outcome`: 同じ phase run を継続し、成功済み訳文と保護要素検証結果を重複保存しない。未処理または失敗 field だけが retry 対象になる。
- `observable point`: phase run ID、`JOB_TRANSLATION_FIELD` 件数、`PHASE_RUN_TRANSLATION_FIELD` 件数、progress、latest error。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: 本文翻訳フェーズの部分成功維持と retry の候補。
- `conflict hint`: CAND-BTP-008 の保存失敗 atomicity と統合時に、成功済み field を維持する範囲を決める必要がある。

### CAND-BTP-011 terminal job へ late provider response を後書きしない

- `source requirement`: `docs/spec.md:155-199` は Completed、Failed、Canceled を終了状態として整理する。`docs/er.md:22-25` は job state を phase run から集約し、翻訳結果を job ごとに保持すると定義する。
- `viewpoint`: 設定不整合 / 保存失敗
- `candidate scenario id`: `CAND-BTP-011`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: cancel、Failed、Completed の後に provider response、validation result、または save request が到着する。
- `expected outcome`: terminal job の `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase state は更新されず、late result の拒否理由を確認できる。
- `observable point`: state snapshot、save result、row count、structured log summary。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: cancel / terminal state の後書き防止候補。
- `conflict hint`: operation-audit 候補が late response を監査対象にする場合、保存対象と redaction 範囲が競合しうる。

### CAND-BTP-012 出力ステータス不整合で output artifact readiness を成立させない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:17-20` は訳文と出力ステータスを output にする。`tasks/usecases/translation-output-artifact.yaml:10-21` は本文翻訳フェーズ完了後の訳文と出力ステータスを input にする。`docs/spec.md:65-66` は xTranslator 互換出力行の `Status` 再構成を求める。
- `viewpoint`: 整合性違反 / 参照不能
- `candidate scenario id`: `CAND-BTP-012`
- `actor`: Job Run を操作するユーザー。
- `trigger`: 訳文はあるが出力ステータスがない、保護要素検証結果が failed なのに translated status になっている、または cached / translated / failed の内部状態が output artifact input と矛盾している。
- `expected outcome`: body phase 完了または output artifact readiness は成立せず、不整合 field と不整合種別を確認できる。
- `observable point`: body phase result、output readiness result、`JOB_TRANSLATION_FIELD`、Job Run UI。
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `adoption hint`: 後続 output artifact との境界候補。xTranslator mapping 自体は後続 task に残す。
- `conflict hint`: translation-output-artifact 側が不整合 row を export 時に除外する候補を持つ場合、body phase 側で拒否する範囲が競合しうる。

### CAND-BTP-013 failure summary と log に secret、raw response、過剰本文を出さない

- `source requirement`: `docs/spec.md:52-58` は暗黙 fallback 不要、APIKey 保存、暗号化保存を定義する。`docs/exec-plans/completed/term-translation-phase/scenario-design.md:93-98` と `docs/exec-plans/completed/persona-generation-phase/scenario-design.md:106-111` は UI / log の redaction 方針を固定している。
- `viewpoint`: 失敗時表示 / セキュリティ
- `candidate scenario id`: `CAND-BTP-013`
- `actor`: Job Run を操作するユーザー。
- `trigger`: provider 失敗、invalid response、保護要素検証失敗、保存失敗が発生する。
- `expected outcome`: UI、error summary、structured log には secret、API key、provider raw request / response、翻訳フィールド本文の全文を出さない。error kind、件数、field identifier、digest、retryable flag は確認できる。
- `observable point`: Job Run UI、error summary、structured log、fake secret store assertion、fake transport log。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `failure_handling_requirement`
- `adoption hint`: failure 観点で redaction を確認する候補。operation-audit 候補と統合されうる。
- `conflict hint`: debug log で prompt / request body を許す候補が出る場合、本文全文と raw response の保存可否が競合しうる。

### CAND-BTP-014 翻訳対象 field 0 件または空 source の扱いを人間判断へ残す

- `source requirement`: `tasks/usecases/body-translation-phase.yaml:9-20` は翻訳フィールド本文から訳文、出力ステータス、保護要素検証結果を作る。`docs/exec-plans/completed/term-translation-phase/scenario-design.md:58-63` と `docs/exec-plans/completed/persona-generation-phase/scenario-design.md:64-69` は 0 件時 Completed を採用しているが、本文翻訳フェーズの 0 件 / 空 source は明示されていない。
- `viewpoint`: 失敗入力 / 境界値
- `candidate scenario id`: `CAND-BTP-014`
- `actor`: 本文翻訳フェーズ実行処理。
- `trigger`: body phase の翻訳対象 field が 0 件、または source text が空の field だけで構成される。
- `expected outcome`: AI だけでは Completed、Blocked、Skipped、Failed のどれにするかを確定しない。候補として、provider 未実行、0 件 summary、output artifact readiness への影響を観測対象に残す。
- `observable point`: target extraction result、phase result、Job Run UI、output readiness result。
- `related detail requirement type`: `boundary_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: human decision candidate。本文翻訳の空対象を正常完了に寄せるか、入力不備に寄せるかを `designer` が質問票へ送る候補。
- `conflict hint`: lifecycle 候補が 0 件 Completed を採用する場合、failure 候補としては採否または merge 判断が必要になる。

## Open Notes

- `human decision candidate`:
  - `CAND-BTP-014`: 翻訳対象 field 0 件または空 source だけの job を Completed、Blocked、Skipped、Failed のどれにするか。
  - `CAND-BTP-008` / `CAND-BTP-010`: 保存失敗や provider 部分失敗で成功済み field を維持する範囲をどこまで許すか。
  - `CAND-BTP-007`: 保護要素検証失敗時に、訳文を一時保持して再試行材料にするか、成功保存せず破棄するか。
- `merge candidate`:
  - `CAND-BTP-006` と external-integration 観点の provider failure 候補。
  - `CAND-BTP-009`、`CAND-BTP-010`、`CAND-BTP-011` と state-transition / lifecycle 観点の操作状態候補。
  - `CAND-BTP-013` と operation-audit 観点の redaction / audit summary 候補。
- `rejection candidate`:
  - 正常系の裏返しだけで観測点がない候補は除外した。
  - 実装方式、トランザクション方式、具体 API 名だけを固定する候補は除外した。
