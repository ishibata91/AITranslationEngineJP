# Scenario Candidates: term-translation-phase / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`
- `candidate_count`: `10`

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/usecases/term-translation-phase.yaml`
  - `../../../../tasks/index.yaml`
  - `../../../spec.md`
  - `../../../er.md`
  - `../../../architecture.md`
  - `../../../diagrams/er/combined-data-model-er.puml`
  - `../../completed/translation-job-setup/plan.md`
  - `../../completed/translation-job-setup/scenario-design.md`
  - `../../completed/translation-job-setup/implementation-scope.md`
- `excluded_sources`:
  - プロダクトコード
  - プロダクトテスト
  - docs 正本の変更
  - 他観点の `scenario-candidates.*.md`
- `generation_notes`:
  - 候補は単語翻訳フェーズの状態、禁止遷移、冪等再実行、回復時の状態不変に限定する。
  - 最終シナリオ表、採否、統合、競合解消は `designer` に残す。
  - `JOB_PHASE_RUN.state` の列挙値は正本で未固定のため、候補では状態名を概念名として扱う。

## Candidate Scenarios

### CAND-TTP-001 Ready ジョブから単語翻訳フェーズを開始する

- `source requirement`:
  - `tasks/usecases/term-translation-phase.yaml:21`
  - `tasks/usecases/term-translation-phase.yaml:32`
  - `docs/spec.md:147-148`
  - `docs/spec.md:226-248`
  - `docs/er.md:22-25`
  - `docs/er.md:63-69`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:241-264`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-001`
- `actor`: ユーザー
- `trigger`: Job Run で単語翻訳フェーズ開始を実行する。
- `pre_transition_state`: `TRANSLATION_JOB.state = Ready`。有効な翻訳ジョブが作成済みで、単語翻訳フェーズの active な `JOB_PHASE_RUN` はない。
- `start_condition`: Ready が成立している。terminal job ではなく、既存 active phase run がない。
- `post_transition_state`: `TRANSLATION_JOB.state = Running`。単語翻訳フェーズ用の `JOB_PHASE_RUN` が作成または開始され、`phase_type`、`state`、`progress_percent`、AI 実行設定を観測できる。
- `expected outcome`: 本文翻訳フェーズではなく、単語翻訳フェーズが current phase として進行する。
- `observable point`: Job Run の current phase、progress、`JOB_PHASE_RUN`、phase 対象の `PHASE_RUN_TRANSLATION_FIELD`。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: Ready 成立後だけ Running へ進む受け入れ条件に統合しやすい。
- `conflict hint`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の二重状態をどう表示するかは UI / lifecycle 候補と競合しうる。

### CAND-TTP-002 共通辞書の完全一致語を翻訳対象から除外する

- `source requirement`:
  - `tasks/usecases/term-translation-phase.yaml:13-18`
  - `tasks/usecases/term-translation-phase.yaml:23`
  - `docs/spec.md:27-33`
  - `docs/spec.md:218-219`
  - `docs/er.md:56`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-002`
- `actor`: フェーズ実行制御
- `trigger`: 単語翻訳フェーズが翻訳対象語と共通辞書を照合する。
- `pre_transition_state`: 単語翻訳フェーズが `Running` で、翻訳対象語の抽出結果と共通辞書が参照可能である。
- `start_condition`: source term が共通辞書の source term と完全一致する。
- `post_transition_state`: 完全一致語は AI 翻訳対象から除外され、置換対象の判定結果として保持される。対象語が全件除外された場合は、未翻訳語なしの phase 完了候補になる。
- `expected outcome`: 共通辞書に存在する語を重複翻訳せず、本文翻訳で置換に使える状態へ進む。
- `observable point`: phase result の除外件数、置換対象判定、`cached` に相当する内部観測情報、AI 実行対象語数。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`, `boundary_requirement`
- `adoption hint`: 共通辞書完全一致を境界条件として、AI 翻訳対象の状態遷移に統合できる。
- `conflict hint`: `cached` を `JOB_TRANSLATION_FIELD.output_status` に反映する時点は本文翻訳フェーズ候補と競合しうる。

### CAND-TTP-003 未辞書語の確定訳語をジョブ内辞書へ反映する

- `source requirement`:
  - `tasks/usecases/term-translation-phase.yaml:9-18`
  - `tasks/usecases/term-translation-phase.yaml:24-25`
  - `docs/spec.md:33`
  - `docs/spec.md:219-227`
  - `docs/er.md:24`
  - `docs/er.md:56`
  - `docs/diagrams/er/combined-data-model-er.puml:128-141`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-003`
- `actor`: フェーズ実行制御
- `trigger`: 未辞書語の訳語が確定する。
- `pre_transition_state`: 単語翻訳フェーズが `Running` で、未辞書語が AI 翻訳または確定処理の対象になっている。
- `start_condition`: source term、translated term、term kind、reusable 判定を確定できる。
- `post_transition_state`: `DICTIONARY_ENTRY` にジョブ内辞書項目が作成または更新され、`translation_job_id` と `PHASE_RUN_DICTIONARY_ENTRY` で phase run に紐づく。
- `expected outcome`: 確定訳語が本文翻訳フェーズの辞書入力として参照可能になる。
- `observable point`: ジョブ内辞書一覧、`DICTIONARY_ENTRY`、`PHASE_RUN_DICTIONARY_ENTRY`、phase result の確定訳語件数。
- `related detail requirement type`: `data_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: ジョブ内辞書への反映を、単語翻訳フェーズ成功の主要 state change として扱える。
- `conflict hint`: 確定訳語の上書き可否や同一 source term の重複扱いは failure / data invariant 候補と競合しうる。

### CAND-TTP-004 単語翻訳フェーズ完了で本文翻訳の参照入力を成立させる

- `source requirement`:
  - `tasks/usecases/term-translation-phase.yaml:22-25`
  - `tasks/usecases/term-translation-phase.yaml:33-34`
  - `docs/spec.md:33`
  - `docs/spec.md:100-103`
  - `docs/spec.md:129`
  - `docs/spec.md:226-248`
  - `docs/er.md:63-69`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-004`
- `actor`: フェーズ実行制御
- `trigger`: 単語翻訳フェーズの対象語がすべて除外済みまたは訳語確定済みになる。
- `pre_transition_state`: 単語翻訳フェーズの `JOB_PHASE_RUN` が `Running` で、未処理対象語が残っていない。
- `start_condition`: 共通辞書除外結果、ジョブ内辞書反映結果、置換対象判定を phase result として集約できる。
- `post_transition_state`: 単語翻訳フェーズの `JOB_PHASE_RUN` が完了状態になり、`progress_percent = 100` と `finished_at` を観測できる。本文翻訳フェーズは単語翻訳フェーズの結果を参照できる。
- `expected outcome`: current phase と progress を確認でき、確定訳語とジョブ内辞書を次フェーズ入力にできる。
- `observable point`: Job Run の phase result、progress 100、`JOB_PHASE_RUN.finished_at`、ジョブ内辞書、本文翻訳入力参照可否。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 本文翻訳フェーズ前提を固定する acceptance scenario に統合しやすい。
- `conflict hint`: 単語翻訳フェーズ完了後の `TRANSLATION_JOB.state` を `Running` のままにするか、次フェーズ待ち状態に戻すかは正本で未固定である。

### CAND-TTP-005 単語翻訳フェーズ未完了では本文翻訳フェーズ開始を拒否する

- `source requirement`:
  - `tasks/usecases/term-translation-phase.yaml:25`
  - `docs/spec.md:100-103`
  - `docs/spec.md:129`
  - `docs/spec.md:226-248`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:241-264`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-005`
- `actor`: 本文翻訳フェーズ起動元
- `trigger`: 単語翻訳フェーズ未完了のジョブで本文翻訳フェーズ開始を試行する。
- `pre_transition_state`: 単語翻訳フェーズの `JOB_PHASE_RUN` が未作成、`Running`、`Paused`、`RecoverableFailed` のいずれかである。
- `start_condition`: 本文翻訳フェーズが参照すべき確定訳語または置換対象判定が完了していない。
- `post_transition_state`: 本文翻訳フェーズの `JOB_PHASE_RUN` は作成されず、既存の job / phase 状態は不変である。
- `expected outcome`: 本文翻訳フェーズは単語翻訳フェーズ完了後だけ開始できる。
- `observable point`: phase transition result、本文翻訳 phase run 未作成、拒否理由、既存 phase run の state 不変。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `boundary_requirement`
- `adoption hint`: 禁止遷移として、後続フェーズ開始条件の acceptance scenario に統合できる。
- `conflict hint`: NPC ペルソナ生成フェーズを本文翻訳前に必須にする候補と、次フェーズ判定の順序が競合しうる。

### CAND-TTP-006 同一ジョブの単語翻訳フェーズ開始を再送しても重複作成しない

- `source requirement`:
  - `docs/er.md:63-69`
  - `docs/er.md:66`
  - `docs/diagrams/er/combined-data-model-er.puml:167-183`
  - `docs/exec-plans/completed/translation-job-setup/implementation-scope.md:117-119`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-006`
- `actor`: ユーザーまたは再送処理
- `trigger`: 同一ジョブで単語翻訳フェーズ開始を二重送信する。
- `pre_transition_state`: 同一 job に単語翻訳フェーズの active または completed `JOB_PHASE_RUN` が存在する。
- `start_condition`: 同一 `translation_job_id` と同一 `phase_type` の phase run がすでに存在する。
- `post_transition_state`: 新しい `JOB_PHASE_RUN` は増えず、既存 phase run が返るか、重複開始として拒否される。ジョブ内辞書の重複項目も作られない。
- `expected outcome`: 再送で phase run、phase 対象、ジョブ内辞書が二重作成されない。
- `observable point`: `JOB_PHASE_RUN` 件数、`PHASE_RUN_DICTIONARY_ENTRY` 件数、`DICTIONARY_ENTRY` 件数、重複開始結果。
- `related detail requirement type`: `冪等性_requirement`, `concurrency_requirement`, `consistency_requirement`
- `adoption hint`: UI double click、通信再送、再実行ボタンの idempotency scenario に統合できる。
- `conflict hint`: 重複開始時に既存 run を返すか拒否するかは user-facing behavior 候補と競合しうる。

### CAND-TTP-007 中断と再開は同じ単語翻訳フェーズ実行へ戻る

- `source requirement`:
  - `docs/spec.md:162-165`
  - `docs/spec.md:193-196`
  - `docs/spec.md:196`
  - `docs/er.md:63-69`
  - `docs/er.md:66`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-007`
- `actor`: ユーザー
- `trigger`: 単語翻訳フェーズ実行中に中断し、その後に再開する。
- `pre_transition_state`: `TRANSLATION_JOB.state = Running` で、単語翻訳フェーズの `JOB_PHASE_RUN` が進行中である。
- `start_condition`: 中断可能な実行中 job で、terminal state ではない。
- `post_transition_state`: 中断で `TRANSLATION_JOB.state = Paused` になり、再開で `Running` へ戻る。同じ `JOB_PHASE_RUN` が継続し、progress と既存の phase link は保持される。
- `expected outcome`: 中断、再開によって phase run やジョブ内辞書が作り直されない。
- `observable point`: job state、phase run ID、progress、`PHASE_RUN_TRANSLATION_FIELD`、`PHASE_RUN_DICTIONARY_ENTRY`。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: job 操作系の状態遷移 scenario に統合しやすい。
- `conflict hint`: job state `Paused` と phase run state の対応関係は正本で未固定である。

### CAND-TTP-008 回復可能失敗のリトライは同じ JOB_PHASE_RUN を戻す

- `source requirement`:
  - `docs/spec.md:178-180`
  - `docs/spec.md:196`
  - `docs/er.md:63-69`
  - `docs/er.md:66`
  - `docs/diagrams/er/combined-data-model-er.puml:167-183`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-008`
- `actor`: ユーザーまたはリトライ制御
- `trigger`: 単語翻訳フェーズが回復可能失敗になった後、リトライまたは再開を実行する。
- `pre_transition_state`: `TRANSLATION_JOB.state = RecoverableFailed` で、単語翻訳フェーズの `JOB_PHASE_RUN.latest_error` を観測できる。
- `start_condition`: 失敗が回復可能で、同じ phase run を再実行できる。
- `post_transition_state`: `TRANSLATION_JOB.state = Running` へ戻る。同じ `JOB_PHASE_RUN` の状態が戻り、attempt 履歴 table は増えない。
- `expected outcome`: リトライで phase run を新規作成せず、最新エラーと進捗を更新しながら再実行できる。
- `observable point`: job state、phase run ID、`latest_error`、progress、phase result。
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `冪等性_requirement`
- `adoption hint`: 外部 AI 失敗や一時的な保存失敗からの回復 scenario に統合できる。
- `conflict hint`: 失敗前に作成済みのジョブ内辞書項目を保持、無効化、削除のどれにするかは human decision 候補である。

### CAND-TTP-009 回復可能失敗から再実行準備へ戻しても本文入力を成立させない

- `source requirement`:
  - `docs/spec.md:178-180`
  - `docs/spec.md:193-196`
  - `docs/spec.md:226-248`
  - `docs/er.md:66`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-009`
- `actor`: ユーザー
- `trigger`: 回復可能失敗した単語翻訳フェーズを再実行準備へ戻す。
- `pre_transition_state`: `TRANSLATION_JOB.state = RecoverableFailed` で、単語翻訳フェーズは完了していない。
- `start_condition`: 再実行準備へ戻せるが、確定訳語と置換対象判定が完了状態ではない。
- `post_transition_state`: `TRANSLATION_JOB.state = Ready` へ戻る。単語翻訳フェーズは未完了扱いのままで、本文翻訳フェーズ開始条件は成立しない。
- `expected outcome`: 再実行準備後に、単語翻訳フェーズをやり直すまで本文翻訳へ進めない。
- `observable point`: job state、単語翻訳 phase result、本文翻訳 phase run 未作成、再実行可能表示。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: RecoverableFailed から Ready へ戻る既存 job 状態遷移を、単語翻訳フェーズの前提不成立と結びつけられる。
- `conflict hint`: Ready へ戻した時に既存 phase run をどの state に置くかは正本で未固定である。

### CAND-TTP-010 terminal job では単語翻訳フェーズの状態変更を拒否する

- `source requirement`:
  - `docs/spec.md:149-150`
  - `docs/spec.md:164-165`
  - `docs/spec.md:181`
  - `docs/spec.md:197-199`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md:241-264`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TTP-010`
- `actor`: ユーザーまたはフェーズ実行制御
- `trigger`: `Completed`、`Failed`、`Canceled` の job で単語翻訳フェーズ開始、リトライ、辞書反映を試行する。
- `pre_transition_state`: `TRANSLATION_JOB.state` が terminal state である。
- `start_condition`: terminal job で phase state やジョブ内辞書を書き換えようとしている。
- `post_transition_state`: job state、phase run、ジョブ内辞書は不変で、禁止遷移として拒否理由を返す。
- `expected outcome`: 終了済み job に対して、単語翻訳フェーズの再開始や辞書の後書きが起きない。
- `observable point`: state transition result、job state 不変、`JOB_PHASE_RUN` 件数不変、`DICTIONARY_ENTRY` 件数不変。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: terminal state の禁止遷移 scenario に統合できる。
- `conflict hint`: 完了後の再翻訳や修正導線が別 task で追加される場合、禁止条件の範囲が競合しうる。

## Open Notes

- `human decision candidate`:
  - `JOB_PHASE_RUN.state` の列挙値と、job state との対応関係を確定する必要がある。
  - 単語翻訳フェーズ完了後の `TRANSLATION_JOB.state` を `Running` のままにするか、次フェーズ待ちの状態へ戻すかは未固定である。
  - 回復可能失敗前に作成されたジョブ内辞書項目を保持、無効化、削除のどれにするかは未固定である。
  - 共通辞書一致時の `cached` 反映時点を、単語翻訳フェーズ完了時にするか本文翻訳時にするかは未固定である。
  - 対象語が全件共通辞書一致した場合に、`JOB_PHASE_RUN` を作成して completed にするか、phase run 作成自体を省略するかは未固定である。
- `merge candidate`:
  - `CAND-TTP-001` と `CAND-TTP-004` は lifecycle 観点の phase 開始 / 完了候補と統合されうる。
  - `CAND-TTP-005` は本文翻訳フェーズ側の前提条件候補と統合されうる。
  - `CAND-TTP-007` から `CAND-TTP-009` は failure 観点の回復候補と統合されうる。
- `conflict candidate`:
  - lifecycle 観点が phase 順序を `単語翻訳 -> NPC ペルソナ生成 -> 本文翻訳` と固定する場合、本文翻訳開始前提はペルソナ生成完了も含める必要がある。
  - external-integration 観点が AI 実行単位を固定する場合、全件共通辞書一致時の AI 未実行完了と競合しうる。
  - operation-audit 観点が failed phase の辞書項目保存を監査対象にする場合、回復時の保持 / 削除判断と競合しうる。
- `rejection candidate`:
  - 採否判断は行わない。純粋な UI 表示だけの進捗確認は state-transition 候補から外し、designer の統合判断に残す。
