# Scenario Candidates: body-translation-phase / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`:
  - `./plan.md`
  - `tasks/index.yaml`
  - `tasks/usecases/body-translation-phase.yaml`
  - `tasks/usecases/term-translation-phase.yaml`
  - `tasks/usecases/persona-generation-phase.yaml`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/architecture.md`
  - `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
  - `docs/exec-plans/completed/persona-generation-phase/scenario-design.md`
- `excluded_sources`:
  - product code
  - product test
  - docs canonicalization
  - final scenario adoption decision
  - other scenario candidate generator output
- `generation_notes`: 実行者の目的、開始経路、成功結果、観測点だけを候補化する。状態遷移網羅、失敗詳細、外部 provider 詳細、監査保存詳細は他観点または designer の統合対象へ残す。

## Candidate Scenarios

### CAND-BTP-001 Job Run から本文翻訳フェーズを開始して進行を確認する

- `source requirement`: `body-translation-phase.yaml` の precondition、manual_check_steps、completion_criteria。`persona-generation-phase` の body phase readiness。`docs/spec.md` の翻訳フローとジョブ進捗確認要件。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-001`
- `actor`: 翻訳ジョブを実行するユーザー
- `trigger`: Job Run で NPC ペルソナ生成フェーズ完了済みの job を選び、本文翻訳フェーズ開始を実行する。
- `expected outcome`: 本文翻訳フェーズが開始され、current phase と progress が Job Run で確認できる。
- `observable point`: Job Run UI、phase start result、`JOB_PHASE_RUN`
- `related detail requirement type`: workflow / display
- `adoption hint`: 後続 designer は persona phase Completed、非 terminal job、active phase run なしを開始条件候補として扱う。
- `conflict hint`: state-transition 観点の開始可否、terminal job guard、active phase run 重複防止と統合が必要になる。

### CAND-BTP-002 確定訳語とジョブ内辞書を参照して本文訳語を一貫させる

- `source requirement`: `body-translation-phase.yaml` の inputs と goal。`term-translation-phase.yaml` の outputs と completion_criteria。`term-translation-phase/scenario-design.md` の確定訳語、ジョブ内辞書、後続フェーズ参照要件。`docs/spec.md` の共通辞書、単語翻訳フェーズ、本文翻訳フェーズ要件。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-002`
- `actor`: 翻訳品質を確認するユーザー
- `trigger`: 確定訳語とジョブ内辞書を持つ job で本文翻訳フェーズを実行する。
- `expected outcome`: 翻訳フィールド本文の訳文が、確定訳語とジョブ内辞書を参照した内容として生成される。
- `observable point`: provider request input summary、訳文、出力ステータス、ジョブ内辞書参照 summary、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: workflow / persistence / display
- `adoption hint`: 参照された辞書件数、適用された確定訳語件数、未適用理由を Job Run の result summary 候補へ含めると designer が表示要件へ接続しやすい。
- `conflict hint`: term phase 側の辞書 snapshot 固定時点、record type ごとの重複 key、cached 相当の内部観測情報との整合が必要になる。

### CAND-BTP-003 ジョブ内ペルソナを参照して NPC 発話の口調を反映する

- `source requirement`: `body-translation-phase.yaml` の inputs と goal。`persona-generation-phase.yaml` の outputs。`persona-generation-phase/scenario-design.md` の persona snapshot 参照状態と body readiness。`docs/spec.md` の NPC ペルソナと本文翻訳フェーズ要件。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-003`
- `actor`: NPC 会話文を確認するユーザー
- `trigger`: ジョブ内ペルソナまたは persona snapshot 参照がある dialogue field を本文翻訳フェーズで処理する。
- `expected outcome`: 訳文が対象 NPC の persona snapshot を参照した結果として生成され、参照状態を確認できる。
- `observable point`: persona snapshot summary、訳文、field result、`PHASE_RUN_PERSONA`、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: workflow / persistence / display
- `adoption hint`: persona 参照あり、共通 persona hit、空 snapshot の各成功パターンを別候補へ分けるかは designer が判断する。
- `conflict hint`: persona snapshot が空の Completed job を成功として扱うか、対象 field で persona 参照不要とするかは persona phase の固定判断と統合が必要になる。

### CAND-BTP-004 翻訳補助メタデータからレコード種別に応じた翻訳指示を構成する

- `source requirement`: `body-translation-phase.yaml` の completion_criteria。`docs/spec.md` の翻訳レコード種別に応じた翻訳指示、ダイアログ補助メタデータ、クエスト補助メタデータ。`docs/er.md` の `TRANSLATION_FIELD_DEFINITION` と翻訳フィールド参照。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-004`
- `actor`: 本文翻訳フェーズ実行処理
- `trigger`: 翻訳フィールドの record type、field type、参照関係、翻訳補助メタデータを読み、provider request 用の翻訳指示を作る。
- `expected outcome`: dialogue、quest、item などの翻訳フィールド本文に対して、種別に応じた翻訳指示が構成される。
- `observable point`: provider request input summary、prompt version または digest、record type / field type summary、phase result
- `related detail requirement type`: workflow / external_integration
- `adoption hint`: actor は内部処理であり、UI 表示は digest と summary に留める候補として扱う。
- `conflict hint`: raw prompt や本文全文を保存または表示するかは security / operation-audit 観点と競合しうる。

### CAND-BTP-005 埋め込み要素を損なわずに訳文と検証結果を確認する

- `source requirement`: `body-translation-phase.yaml` の outputs と completion_criteria。`docs/spec.md` の `<10gold>` などの埋め込み要素保持要件。`docs/er.md` の翻訳結果保持方針。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-005`
- `actor`: 翻訳結果を確認するユーザー
- `trigger`: 埋め込み要素を含む翻訳フィールド本文を本文翻訳フェーズで処理する。
- `expected outcome`: 訳文が生成され、保護要素検証結果が成功または確認可能な注意状態として表示される。
- `observable point`: 訳文、保護要素検証結果、出力ステータス、Job Run result summary、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: workflow / display / persistence
- `adoption hint`: 成功表示だけでなく、保護要素検証が注意状態のときに output status をどう扱うかを designer の人間判断候補に残す。
- `conflict hint`: 検証失敗を recoverable failure、field-level warning、または phase failure のどれにするかは failure / state-transition 観点と競合する。

### CAND-BTP-006 翻訳済みフィールドの訳文と出力ステータスを確認する

- `source requirement`: `body-translation-phase.yaml` の outputs と manual_check_steps。`translation-output-artifact.yaml` の inputs。`docs/spec.md` の `FormID`、`EditorID`、レコード種別、フィールド種別、原文、訳文、出力ステータスの lossless 保持要件。`docs/er.md` の `JOB_TRANSLATION_FIELD`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-006`
- `actor`: 翻訳成果物出力へ進めたいユーザー
- `trigger`: 本文翻訳フェーズ完了後に Job Run で翻訳フィールド結果を確認する。
- `expected outcome`: 翻訳フィールドごとの訳文、出力ステータス、保護要素検証結果を確認でき、後続の翻訳成果物出力の入力になる。
- `observable point`: Job Run UI、field result list、`JOB_TRANSLATION_FIELD`、translation output readiness summary
- `related detail requirement type`: display / persistence / workflow
- `adoption hint`: output artifact phase へ渡す readiness は designer 統合時の接続候補に留める。
- `conflict hint`: xTranslator `Status` への写像や output artifact の diff preview は後続 task の範囲であり、本候補では確定しない。

### CAND-BTP-007 本文翻訳フェーズ中の操作可否を確認する

- `source requirement`: `body-translation-phase.yaml` の completion_criteria。`docs/spec.md` の中断、再開、失敗回復、ジョブ進捗確認要件。`persona-generation-phase/scenario-design.md` の pause、resume、retry、cancel を body phase と同じように許可する固定判断。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-007`
- `actor`: 長時間実行を管理するユーザー
- `trigger`: 本文翻訳フェーズが running、paused、recoverable failed のいずれかになった Job Run を開く。
- `expected outcome`: ユーザーが pause、resume、retry、cancel の可否と次に可能な操作を確認できる。
- `observable point`: Job Run UI、button enablement、phase state summary、latest error summary
- `related detail requirement type`: display / workflow
- `adoption hint`: actor の成功体験は操作可否の確認であり、状態遷移網羅は state-transition 観点へ委ねる。
- `conflict hint`: cancel 後の出力ステータス、retry 時の部分成功保持、active phase run 再利用は lifecycle / state-transition / failure 観点と統合が必要になる。

### CAND-BTP-008 recoverable failure の情報を見て再実行判断を行う

- `source requirement`: `body-translation-phase.yaml` の completion_criteria と manual_check_steps。`docs/spec.md` の RecoverableFailed、再開、リトライ、失敗回復要件。`term-translation-phase/scenario-design.md` と `persona-generation-phase/scenario-design.md` の secret / raw response 非露出判断。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-008`
- `actor`: 失敗した本文翻訳フェーズを回復したいユーザー
- `trigger`: 本文翻訳フェーズが recoverable failure になった Job Run を確認する。
- `expected outcome`: 失敗理由、影響範囲、再試行可否を確認でき、secret、raw request、raw response、過剰本文は表示されない。
- `observable point`: Job Run error summary、retry enablement、phase result、redacted provider summary
- `related detail requirement type`: display / workflow / security
- `adoption hint`: ユーザーの目的は回復判断であり、failure の種類や再試行単位は failure 観点へ残す。
- `conflict hint`: 失敗時に field-level result を公開する範囲、partial translated text を保持するかは failure / state-invariant 統合で判断する必要がある。

### CAND-BTP-009 後続の翻訳成果物出力へ進める完了結果を確認する

- `source requirement`: `body-translation-phase.yaml` の outputs。`translation-output-artifact.yaml` の preconditions と inputs。`docs/spec.md` の結果確認から翻訳成果物出力へ進む業務フロー。`docs/er.md` の `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-BTP-009`
- `actor`: xTranslator 互換成果物を作りたいユーザー
- `trigger`: 本文翻訳フェーズが完了した Job Run で出力前確認を行う。
- `expected outcome`: 本文翻訳フェーズの完了、訳文件数、出力ステータス、保護要素検証結果の summary を確認し、翻訳成果物出力へ進めるか判断できる。
- `observable point`: Job Run result summary、output phase readiness、field count summary、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: workflow / display
- `adoption hint`: 後続 task の開始条件候補として残し、xTranslator output row の生成ルールは translation-output-artifact 側で扱う。
- `conflict hint`: body phase が warning を含む Completed のとき output artifact へ進めるかは人間判断候補になりうる。

## Open Notes

- `human decision candidate`:
  - 保護要素検証が warning の場合に、出力ステータスを success 相当にするか、field-level warning にするか、recoverable failure にするか。
  - persona snapshot が空の Completed job で、dialogue 以外の field を含む本文翻訳フェーズを成功として扱うか。
  - 本文翻訳フェーズ result summary に prompt digest、provider request input summary、辞書 / persona 参照 summary のどこまでを UI 表示するか。
  - partial translated text を recoverable failure 時にユーザーへ見せるか、Completed まで非公開にするか。
- `merge candidate`:
  - `CAND-BTP-001` と `CAND-BTP-007` は Job Run の phase 操作確認として統合される可能性がある。
  - `CAND-BTP-002` と `CAND-BTP-003` と `CAND-BTP-004` は provider request input summary の観測候補として統合される可能性がある。
  - `CAND-BTP-006` と `CAND-BTP-009` は body phase result と output readiness の表示候補として統合される可能性がある。
- `rejection candidate`:
  - `CAND-BTP-008` は failure 観点の候補と重複する場合、actor-goal 側では UI の回復判断表示だけへ縮小できる。
  - `CAND-BTP-004` は external-integration 観点の provider request 候補と重複する場合、actor-goal 側では record type 別 translation instruction summary の確認だけへ縮小できる。
