# Scenario Candidates: body-translation-phase / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`

## Generator Scope

- `viewpoint`: lifecycle
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
- `excluded_sources`: product code、product test、docs 正本化、最終シナリオ表、採否、統合判断、他 agent 出力の生成。
- `generation_notes`: 本文翻訳フェーズが作られてから後続の翻訳成果物出力へ渡るまでの時間順だけを候補化する。

## Candidate Scenarios

### CAND-BTP-001 persona phase 完了後に本文翻訳フェーズを開始する

- `source requirement`: `./plan.md:74-77`, `../../../../tasks/usecases/body-translation-phase.yaml:21-30`, `../persona-generation-phase/scenario-design.md:92-97`, `../../../spec.md:129-133`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-001`
- `lifecycle stage`: 開始
- `start condition`: persona phase が Completed で、persona snapshot 参照が成立し、terminal job ではなく、同時 active phase run がない。
- `actor`: ユーザー
- `trigger`: Job Run で本文翻訳フェーズ開始を実行する。
- `expected outcome`: body phase の `JOB_PHASE_RUN` が開始され、current phase と progress を確認できる。開始不可の場合は phase run を作らず、拒否理由を確認できる。
- `observable point`: Job Run UI、phase start result、`JOB_PHASE_RUN`、body input summary。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `compatibility_requirement`
- `adoption hint`: body phase の主開始シナリオ候補として、state-transition 候補の開始条件と突き合わせる。
- `conflict hint`: persona snapshot 参照不能時の扱いは failure 観点と競合しうる。terminal job と active phase run の厳密条件は state-transition 観点に残す。

### CAND-BTP-002 本文翻訳入力 snapshot を固定する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:11-17`, `../../../../tasks/usecases/body-translation-phase.yaml:23-26`, `../../../spec.md:33-35`, `../../../spec.md:41-43`, `../../../spec.md:227-229`, `../term-translation-phase/scenario-design.md:72-77`, `../persona-generation-phase/scenario-design.md:78-83`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-002`
- `lifecycle stage`: 作成後の入力固定
- `start condition`: body phase run が開始済みで、確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを読める。
- `actor`: システム
- `trigger`: body phase 実行の最初に翻訳対象フィールドと参照データを収集する。
- `expected outcome`: 翻訳対象フィールド、確定訳語、ジョブ内辞書、persona snapshot、翻訳補助メタデータ、翻訳レコード種別別の翻訳指示が同一 phase run の入力 summary として固定される。
- `observable point`: input summary、input snapshot digest、対象フィールド件数、対象外理由、翻訳指示構成 summary。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: provider request 前の lifecycle 候補として、external-integration 候補の request 構成と接続する。
- `conflict hint`: snapshot を phase 開始時に固定するか、再開時に再読込するかは state-transition / operation-audit 観点と競合しうる。

### CAND-BTP-003 翻訳フィールド本文を実行中 progress として進める

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:23-28`, `../../../spec.md:11`, `../../../spec.md:41-43`, `../../../spec.md:53-57`, `../../../architecture.md:115-129`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-003`
- `lifecycle stage`: 実行
- `start condition`: body input snapshot が固定済みで、対象フィールドが 1 件以上あり、AI 実行に必要な phase 設定を参照できる。
- `actor`: システム
- `trigger`: body phase executor が翻訳フィールド本文の処理を開始する。
- `expected outcome`: 翻訳フィールド本文が翻訳指示、確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを使って処理され、phase progress が更新される。
- `observable point`: Job Run progress、phase execution summary、provider execution summary、処理済み件数、未処理件数。
- `related detail requirement type`: `success_requirement`, `performance_requirement`, `observability_requirement`
- `adoption hint`: 実行中 progress の lifecycle 候補として、actor-goal 候補のユーザー確認観点と統合できる。
- `conflict hint`: request unit の粒度、Batch API 利用時の progress 単位、provider 設定継承は external-integration 観点と競合しうる。

### CAND-BTP-004 保護要素検証を field success の前に完了する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:18-20`, `../../../../tasks/usecases/body-translation-phase.yaml:25-28`, `../../../spec.md:42-43`, `../../../spec.md:231`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-004`
- `lifecycle stage`: 段階内検証
- `start condition`: 翻訳候補の訳文が field 単位で返っている。
- `actor`: システム
- `trigger`: 訳文を保存可能結果へ進める前に埋め込み要素と保護要素を照合する。
- `expected outcome`: 埋め込み要素を損なわない訳文だけが成功候補になり、保護要素検証結果を field result として確認できる。検証失敗は successful Completed に進めない。
- `observable point`: 保護要素検証結果、field result、phase result、error summary、Job Run UI。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: body phase 固有の lifecycle 検証候補として、failure 候補の validation failure と接続する。
- `conflict hint`: 保護要素検証失敗を field failed、phase RecoverableFailed、または retryable validation error のどれで扱うかは人間判断候補になりうる。

### CAND-BTP-005 訳文、出力ステータス、検証結果を保存して phase を完了する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:18-20`, `../../../spec.md:43`, `../../../er.md:22-25`, `../../../er.md:40`, `../../../er.md:63-69`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-005`
- `lifecycle stage`: 保存と完了
- `start condition`: 対象 field の訳文と保護要素検証結果が成功扱いにできる。
- `actor`: システム
- `trigger`: field result を job-scoped result として永続化する。
- `expected outcome`: `JOB_TRANSLATION_FIELD` に訳文と出力ステータスが保持され、`PHASE_RUN_TRANSLATION_FIELD` から body phase が対象にした field を追跡できる。全対象 field が成功または明示的な terminal result になった時だけ body phase が Completed になる。
- `observable point`: `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase result、Job Run result summary。
- `related detail requirement type`: `data_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 完了条件候補として、state-transition 候補の Completed 遷移と統合する。
- `conflict hint`: field 単位の terminal result 種別と phase Completed 条件は failure / state-transition 観点と競合しうる。

### CAND-BTP-006 翻訳対象 0 件または provider 未実行でも終点を確認する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:23-28`, `../../../spec.md:30-35`, `../term-translation-phase/scenario-design.md:58-63`, `../persona-generation-phase/scenario-design.md:64-69`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-006`
- `lifecycle stage`: 空対象または翻訳不要完了
- `start condition`: body input snapshot の結果、本文翻訳対象 field が 0 件、または provider 実行対象が 0 件になる。
- `actor`: システム
- `trigger`: body phase が provider request を作る前に対象件数を確定する。
- `expected outcome`: provider を呼ばずに phase result が作られ、対象 0 件、provider 未実行、後続出力可否を確認できる。
- `observable point`: phase result、target count、provider skipped reason、Job Run UI、output readiness。
- `related detail requirement type`: `boundary_requirement`, `state_requirement`, `observability_requirement`
- `adoption hint`: 空対象 lifecycle 候補として、designer が body phase の Completed 条件に含めるか判断する。
- `conflict hint`: 対象 0 件を Completed とするか human decision とするかは、本文翻訳フェーズ固有の業務判断と競合しうる。

### CAND-BTP-007 pause、resume、retry で同じ phase run を継続する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:27-28`, `../../../spec.md:53`, `../../../spec.md:162-180`, `../../../er.md:63-69`, `../term-translation-phase/scenario-design.md:320-327`, `../persona-generation-phase/scenario-design.md:314-321`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-007`
- `lifecycle stage`: 中断、再開、リトライ
- `start condition`: body phase run が Running、Paused、RecoverableFailed のいずれかである。
- `actor`: ユーザーまたはシステム
- `trigger`: pause、resume、retry、開始再送を実行する。
- `expected outcome`: 同じ `JOB_PHASE_RUN` を継続し、成功済み field result と未処理 field を区別する。再開またはリトライでは重複 field result を作らず、progress と latest error が更新される。
- `observable point`: phase run ID、progress、latest error、`JOB_TRANSLATION_FIELD` 件数、`PHASE_RUN_TRANSLATION_FIELD` 件数、Job Run UI。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `concurrency_requirement`
- `adoption hint`: 回復 lifecycle 候補として、state-transition 候補と failure 候補の接続点にする。
- `conflict hint`: 成功済み訳文を retry で維持するか再翻訳するかは人間判断候補になりうる。

### CAND-BTP-008 recoverable failure では後続出力へ進めない

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:27-28`, `../../../../tasks/usecases/translation-output-artifact.yaml:20-28`, `../../../spec.md:178-182`, `../../../spec.md:189-199`, `../persona-generation-phase/scenario-design.md:99-104`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-008`
- `lifecycle stage`: 失敗回復
- `start condition`: provider 失敗、応答不正、保存失敗、保護要素検証失敗、入力参照不能のいずれかが発生する。
- `actor`: システム
- `trigger`: body phase 実行中に recovery 可能な失敗を記録する。
- `expected outcome`: phase は successful Completed にならず、recoverable failure 情報と再試行可否を確認できる。本文翻訳フェーズが完了するまで translation-output-artifact は開始可能にならない。
- `observable point`: phase state、error kind、retryable flag、output readiness、Job Run failure state。
- `related detail requirement type`: `recovery_requirement`, `failure_handling_requirement`, `state_requirement`
- `adoption hint`: lifecycle 側では後続出力への遮断だけを候補化し、失敗種別の詳細は failure 観点へ残す。
- `conflict hint`: failure 種別ごとの RecoverableFailed / Failed の分類は failure 観点と state-transition 観点に委ねる。

### CAND-BTP-009 cancel で body phase を終了し、後続出力を遮断する

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:27-28`, `../../../spec.md:162-166`, `../../../spec.md:189-199`, `../persona-generation-phase/scenario-design.md:364-383`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-009`
- `lifecycle stage`: 取消と終了
- `start condition`: body phase が cancel 可能な状態である。
- `actor`: ユーザー
- `trigger`: Job Run で cancel を実行する。
- `expected outcome`: cancel 後は body phase を Completed として扱わず、訳文、出力ステータス、保護要素検証結果の後書きを拒否する。translation-output-artifact の readiness は成立しない。
- `observable point`: phase state、cancel result、後書き拒否理由、output readiness、Job Run UI。
- `related detail requirement type`: `state_requirement`, `authorization_requirement`, `consistency_requirement`
- `adoption hint`: lifecycle の終端候補として、state-transition 候補の terminal guard と統合する。
- `conflict hint`: Running から直接 cancel できるか、Paused を経由するかは body phase 固有の人間判断候補になりうる。

### CAND-BTP-010 Completed body phase を結果確認と出力フェーズへ渡す

- `source requirement`: `../../../../tasks/usecases/body-translation-phase.yaml:33-37`, `../../../../tasks/usecases/translation-output-artifact.yaml:9-28`, `../../../spec.md:100-106`, `../../../spec.md:129-133`, `../../../er.md:77`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-BTP-010`
- `lifecycle stage`: 終了後利用
- `start condition`: body phase が Completed で、訳文、出力ステータス、保護要素検証結果が参照可能である。
- `actor`: ユーザー
- `trigger`: Job Run で本文翻訳結果を確認し、後続の Output Review または translation-output-artifact へ進む。
- `expected outcome`: 訳文、出力ステータス、保護要素検証結果を確認でき、completed job として translation-output-artifact の入力にできる。
- `observable point`: Job Run result summary、`JOB_TRANSLATION_FIELD`、output readiness、Output Review 入力 summary。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 終了後利用候補として、translation-output-artifact 側の scenario と接続する。
- `conflict hint`: 結果確認で修正ありの場合に body phase を再実行するか、別 job / review workflow に戻すかは designer の統合判断に残す。

## Open Notes

- `human decision candidate`:
  - `CAND-BTP-004`: 保護要素検証失敗を field failed、phase RecoverableFailed、または retryable validation error のどれで扱うか。
  - `CAND-BTP-006`: 翻訳対象 0 件または provider 未実行を body phase Completed としてよいか。
  - `CAND-BTP-007`: retry 時に成功済み訳文を維持して未処理だけ進めるか、成功済み field も再翻訳可能にするか。
  - `CAND-BTP-009`: Running から直接 cancel できるか、Paused 経由だけにするか。
  - `CAND-BTP-010`: 結果確認で修正ありの場合の戻り先を body phase retry、ジョブ再実行準備、review workflow のどれにするか。
- `merge candidate`:
  - `CAND-BTP-001` は state-transition の開始可否候補と統合しうる。
  - `CAND-BTP-002` は external-integration の request 構成候補と統合しうる。
  - `CAND-BTP-004` は failure の保護要素検証失敗候補と統合しうる。
  - `CAND-BTP-007` と `CAND-BTP-008` は recovery / retry の最終シナリオで統合しうる。
  - `CAND-BTP-010` は translation-output-artifact の開始候補と接続しうる。
- `rejection candidate`: なし。lifecycle 観点では採否を判断しない。
