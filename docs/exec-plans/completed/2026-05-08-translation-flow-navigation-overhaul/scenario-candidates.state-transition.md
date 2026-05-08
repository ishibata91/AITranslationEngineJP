# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFNO-ST`

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`: `plan.md`, `navigation-state-machine.puml`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、`.codex` 変更、ツール権限、実装指示
- `generation_notes`: state-transition 観点だけを候補化する。最終シナリオ表、採否、統合、競合解消は扱わない。

## Candidate Scenarios

### CAND-TFNO-ST-001 新規 job 作成後に単語翻訳ページへ進む

- `source requirement`: `plan.md:16`, `plan.md:69-70`, `navigation-state-machine.puml:27-30`, `navigation-state-machine.puml:70-71`, `translation-job-setup.md:20-21`, `translation-job-setup.md:45`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-001`
- `actor`: 利用者
- `trigger`: `Job Setup` ページで登録済み入力データと 3 つの翻訳段階の設定が作成条件を満たし、job 作成を完了する。
- `transition before state`: `JobSetupPage`。job は未作成で、選択 input と各翻訳段階の設定が固定済み。
- `start condition`: 3 つの翻訳段階で API キー不足と model 未選択がない。
- `transition after state`: Ready job と初期 phase 状態が作成され、表示状態は `TermPhasePage` になる。
- `expected outcome`: 作成直後の job は旧 `Job Run` のセッション取得を経由せず、単語翻訳ページで対象 job と単語翻訳 summary を読める。
- `observable point`: `TermPhasePage` に jobId、単語翻訳 summary、開始可否、次へ進めない理由が表示される。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 新規作成直後の許可遷移として扱える。
- `conflict hint`: 旧 `Job Run` 表示対象への移動を残す候補と競合する。

### CAND-TFNO-ST-002 未完了 job 一覧から再開対象を固定する

- `source requirement`: `plan.md:17`, `plan.md:43`, `navigation-state-machine.puml:32-35`, `navigation-state-machine.puml:74-76`, `translation-job-management.md:13-14`, `translation-job-management.md:26-35`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-002`
- `actor`: 利用者
- `trigger`: 未完了 job 一覧で Ready、Running、Paused、RecoverableFailed、Failed、Canceled のいずれかの job を選択する。
- `transition before state`: `IncompleteJobListPage`。一覧は Completed 以外の job を表示し、一覧表示だけでは job 状態を変えない。
- `start condition`: 選択 job が参照可能であり、phase progress を安全側に集約できる。
- `transition after state`: 選択 job の jobId と表示フェーズが固定され、該当する単語翻訳、NPC ペルソナ生成、本文翻訳ページのいずれかへ移動する。
- `expected outcome`: 途中再開は一覧選択だけを入口にし、選択していない job のフェーズページを表示しない。
- `observable point`: 遷移先ページで current phase、phase state、progress、操作可否、再開不可理由を確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 途中再開の許可遷移として扱える。
- `conflict hint`: 参照不能 job を表示対象にする候補、または session fetch で対象 job を後から取得する候補と競合する。

### CAND-TFNO-ST-003 Completed job を翻訳管理の未完了一覧から除外する

- `source requirement`: `plan.md:91-98`, `plan.md:134-138`, `translation-job-management.md:26-27`, `translation-job-management.md:82`, `translation-output-artifact.md:20`, `translation-output-artifact.md:27`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-003`
- `actor`: 利用者
- `trigger`: 翻訳管理の未完了 job 一覧を開く。
- `transition before state`: job が Completed、または Completed 以外の状態で存在する。
- `start condition`: job 一覧が読み込める。
- `transition after state`: Completed job は未完了一覧へ入らず、Completed job は成果物出力セクションの completed job 一覧で扱われる。
- `expected outcome`: 翻訳管理は job 作成、未完了一覧、フェーズ実行、再開までを扱い、成果物出力処理を直接開始しない。
- `observable point`: 未完了一覧に Completed job が表示されず、成果物出力側だけが completed job list を表示する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 翻訳管理と成果物出力の分離不変条件として扱える。
- `conflict hint`: 完了済み job を翻訳管理一覧に残す候補と競合する。

### CAND-TFNO-ST-004 フェーズページへの直リンクを未完了 job 一覧へ戻す

- `source requirement`: `plan.md:56-65`, `plan.md:119`, `navigation-state-machine.puml:68`, `navigation-state-machine.puml:130-132`, `navigation-state-machine.puml:147`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-004`
- `actor`: 利用者
- `trigger`: route state または復元状態が不整合なまま、単語翻訳、NPC ペルソナ生成、本文翻訳ページへ入る。
- `transition before state`: 対象 job が未確定で、フェーズページが要求する jobId と summary の前提が成立していない。
- `start condition`: フェーズページに入ろうとした時点で対象 job を一意に確定できない。
- `transition after state`: job 状態を変更せず、`IncompleteJobListPage` へ戻す。
- `expected outcome`: フェーズページは対象 job が確定している場合だけ表示される。
- `observable point`: 未完了 job 一覧が表示され、利用者は再開対象 job を選び直せる。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 直リンク防止の禁止遷移候補として扱える。
- `conflict hint`: 直リンク時に最後の選択 job を推測復元する候補と競合する可能性がある。

### CAND-TFNO-ST-005 フェーズページから入力データページまたは Job Setup へ戻らない

- `source requirement`: `plan.md:34-43`, `navigation-state-machine.puml:124-129`, `translation-job-setup.md:29`, `translation-job-management.md:30-35`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-005`
- `actor`: 利用者
- `trigger`: 作成済み job のフェーズページから入力データページまたは `Job Setup` ページへ戻ろうとする。
- `transition before state`: 対象 job が作成済みで、入力と設定は job の前提として固定済み。
- `start condition`: 現在ページが単語翻訳、NPC ペルソナ生成、本文翻訳ページのいずれかである。
- `transition after state`: 戻り遷移を表示しない。誤った復元状態では job 状態を変えず、未完了 job 一覧で再選択させる。
- `expected outcome`: 既存 job を再編集できるように見せない。入力や設定を変える場合は新しい job 作成で扱う。
- `observable point`: `sticky footer` には `次へ進む` と `一覧へ戻る` を中心にした移動だけが表示される。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 作成済み job の状態不変条件として扱える。
- `conflict hint`: フェーズページから `Job Setup` へ戻って設定変更できる候補と競合する。

### CAND-TFNO-ST-006 単語翻訳完了と辞書参照成立後だけ NPC ペルソナ生成へ進む

- `source requirement`: `navigation-state-machine.puml:37-42`, `navigation-state-machine.puml:78-80`, `term-translation-phase.md:22-23`, `term-translation-phase.md:55-56`, `term-translation-phase.md:71-72`, `docs/spec.md:129`, `docs/spec.md:226-228`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-006`
- `actor`: 利用者
- `trigger`: 単語翻訳ページで `次へ進む` を選ぶ。
- `transition before state`: `TermPhasePage`。単語翻訳 phase state とジョブ内辞書参照状態が表示されている。
- `start condition`: 単語翻訳フェーズが Completed であり、ジョブ内辞書参照が成立している。
- `transition after state`: `PersonaPhasePage` へ移動する。
- `expected outcome`: 単語翻訳未完了、失敗中、辞書参照不能の場合は後続 phase run を作らず、ページ状態を維持する。
- `observable point`: `sticky footer` に次へ進めない理由が表示され、許可時だけ NPC ペルソナ生成ページへ進む。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 単語翻訳から NPC ペルソナ生成への許可遷移と禁止遷移を同時に扱える。
- `conflict hint`: 単語翻訳未完了で NPC ペルソナ生成を開始する候補と競合する。

### CAND-TFNO-ST-007 NPC ペルソナ生成完了と snapshot 参照成立後だけ本文翻訳へ進む

- `source requirement`: `navigation-state-machine.puml:44-49`, `navigation-state-machine.puml:82-84`, `persona-generation-phase.md:13-23`, `persona-generation-phase.md:44-45`, `persona-generation-phase.md:57-60`, `docs/spec.md:130`, `docs/spec.md:225-228`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-007`
- `actor`: 利用者
- `trigger`: NPC ペルソナ生成ページで `次へ進む` を選ぶ。
- `transition before state`: `PersonaPhasePage`。persona phase state と body readiness が表示されている。
- `start condition`: NPC ペルソナ生成フェーズが Completed であり、persona snapshot 参照が成立している。
- `transition after state`: `BodyPhasePage` へ移動する。
- `expected outcome`: persona 未完了、失敗、snapshot 参照不能では本文翻訳 phase run を作らず、ページ状態を維持する。
- `observable point`: body phase readiness、snapshot 参照状態、次へ進めない理由が表示される。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: NPC ペルソナ生成から本文翻訳への許可遷移と禁止遷移を同時に扱える。
- `conflict hint`: snapshot 参照不能でも本文翻訳へ進める候補と競合する。

### CAND-TFNO-ST-008 本文翻訳完了時だけ output readiness を成立させる

- `source requirement`: `navigation-state-machine.puml:51-56`, `body-translation-phase.md:13-22`, `body-translation-phase.md:34-43`, `body-translation-phase.md:70-75`, `translation-output-artifact.md:20`, `translation-output-artifact.md:27-28`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-008`
- `actor`: 利用者
- `trigger`: 本文翻訳ページで本文翻訳フェーズを完了、失敗、取消、または再試行する。
- `transition before state`: `BodyPhasePage`。本文翻訳 phase state、field result 整合、output status 整合が表示されている。
- `start condition`: body phase が Completed であり、field result 整合と output status 整合を満たす。
- `transition after state`: job-level 状態が Completed になり、output readiness が true になる。
- `expected outcome`: Failed、Canceled、保存失敗、検証失敗、field result 不整合では output readiness を true にしない。
- `observable point`: output readiness、failure state、error kind、retryable flag、field result summary が状態に応じて表示される。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 本文翻訳から成果物出力可能状態への不変条件として扱える。
- `conflict hint`: `navigation-state-machine.puml:87` は `job Completed / Canceled / Failed` で翻訳完了ページへ進む表記を持つ。Completed 以外を翻訳完了ページで扱うかは競合候補になる。

### CAND-TFNO-ST-009 本文翻訳の再送、再開、リトライで重複作成しない

- `source requirement`: `body-translation-phase.md:39-41`, `body-translation-phase.md:56-57`, `term-translation-phase.md:43-46`, `persona-generation-phase.md:38-42`, `translation-output-artifact.md:31`, `translation-output-artifact.md:37`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-009`
- `actor`: 利用者
- `trigger`: phase の開始再送、再開、リトライ、または late response が発生する。
- `transition before state`: 対象 job に既存 `JOB_PHASE_RUN`、成功済み result、または terminal job 状態が存在する。
- `start condition`: 再送、再開、リトライで同じ phase run を継続できる。terminal job では後書き対象ではない。
- `transition after state`: 同じ `JOB_PHASE_RUN` を継続し、成功済み result と対応行を重複作成しない。terminal job では状態を変えない。
- `expected outcome`: job、phase run、辞書 entry、persona、translation field、artifact row の重複作成を防ぐ。
- `observable point`: retryable failure、late response rejected、重複なしの progress、既存 result summary が確認できる。
- `related detail requirement type`: `冪等性_requirement`, `concurrency_requirement`, `state_requirement`
- `adoption hint`: 状態遷移の冪等再実行候補として扱える。
- `conflict hint`: phase ごとの重複防止を個別候補へ分割する設計と統合粒度が競合する可能性がある。

### CAND-TFNO-ST-010 翻訳完了ページは確認と出力管理への移動だけを扱う

- `source requirement`: `plan.md:78-88`, `plan.md:97-98`, `navigation-state-machine.puml:58-61`, `navigation-state-machine.puml:90-92`, `translation-output-artifact.md:13-15`, `translation-output-artifact.md:63-68`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-010`
- `actor`: 利用者
- `trigger`: 翻訳完了ページで出力管理への移動を選ぶ。
- `transition before state`: `TranslationCompletePage`。原文と訳文をページング表示している。
- `start condition`: 翻訳完了ページが出力処理そのものを扱わない。
- `transition after state`: `OutputArtifactSection` へ移動する。artifact 生成状態はまだ開始しない。
- `expected outcome`: 翻訳完了ページのボタンは出力管理への導線に限定され、XML 出力、preview、再出力、互換性確認は成果物出力側で扱う。
- `observable point`: 出力管理側で completed job list、selected job summary、output readiness、diff preview、出力 action 可否を確認できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 翻訳管理から成果物出力へのセクション遷移候補として扱える。
- `conflict hint`: 出力対象 job を自動選択するかどうかは `plan.md:87` と `plan.md:153` で未決である。

### CAND-TFNO-ST-011 成果物出力は Completed job だけを出力候補にする

- `source requirement`: `plan.md:91-98`, `plan.md:134-138`, `translation-output-artifact.md:20`, `translation-output-artifact.md:27-28`, `translation-output-artifact.md:65-68`, `body-translation-phase.md:16`, `body-translation-phase.md:43`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-011`
- `actor`: 利用者
- `trigger`: 成果物出力セクションで completed job を選択し、出力 action または再出力 action を実行しようとする。
- `transition before state`: `CompletedJobListPage` または `OutputReviewPage`。job-level 状態、body phase、field result 整合、output status 整合が表示されている。
- `start condition`: body phase Completed、job-level Completed、field result 整合、output status 整合、row validation pass、出力先 path valid を満たす。
- `transition after state`: `OutputReviewPage` で preview-ready、generating、success、failed、stale、re-output-needed の出力状態へ遷移する。
- `expected outcome`: 未完了、失敗中、Canceled、不整合 job では artifact 生成を開始せず、成功状態の artifact と row を作らない。
- `observable point`: disabled action の理由、output readiness、row validation、artifact status、re-output state が表示される。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 成果物出力側の許可遷移と禁止遷移をまとめた候補として扱える。
- `conflict hint`: 翻訳管理側から成果物生成を直接開始する候補と競合する。

### CAND-TFNO-ST-012 旧 Job Run のセッション取得で対象 job を探さない

- `source requirement`: `plan.md:67-76`, `plan.md:130`, `navigation-state-machine.puml:135-137`, `translation-job-management.md:33-35`, `translation-job-management.md:44`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TFNO-ST-012`
- `actor`: 利用者
- `trigger`: フェーズページで対象 job 未確定のまま表示更新または対象取得を行おうとする。
- `transition before state`: フェーズページが job を持たない、または参照不能 job が指定されている。
- `start condition`: job は `Job Setup` 完了直後、または未完了 job 一覧の選択結果から受け取る必要がある。
- `transition after state`: セッション取得による job 探索は行わず、参照不能 job は表示対象にしない。必要なら未完了 job 一覧へ戻す。
- `expected outcome`: 一覧選択と別の対象取得経路が併存しない。
- `observable point`: フェーズページ上にセッション取得操作がなく、選択済み jobId または未完了 job 一覧への復帰だけが観測できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 旧 `Job Run` 廃止に伴う禁止遷移候補として扱える。
- `conflict hint`: session fetch を fallback として残す候補と競合する。

## Open Notes

- `candidate_count`: 12
- `human decision candidate`: [CAND-TFNO-ST-008] `navigation-state-machine.puml:87` は `job Completed / Canceled / Failed` で翻訳完了ページへ進む表記を持つ。翻訳完了ページを Completed 専用にするか、Failed / Canceled の確認ページも同じページで扱うかは designer の統合判断に残す。
- `human decision candidate`: [CAND-TFNO-ST-010] 出力管理へ移動した後に出力対象 job を自動選択するか、出力管理側で選ばせるかは `plan.md:87` と `plan.md:153` で未決である。
- `merge candidate`: [CAND-TFNO-ST-006] と [CAND-TFNO-ST-007] は、前フェーズ完了と参照成立後だけ次フェーズへ進む同型候補である。
- `merge candidate`: [CAND-TFNO-ST-003]、[CAND-TFNO-ST-010]、[CAND-TFNO-ST-011] は、翻訳管理と成果物出力の分離として統合候補になる。
- `rejection candidate`: なし。採否は designer に残す。

## Evidence Summary

- `task_artifact_root`: `docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/`
- `candidate_artifact`: `docs/exec-plans/active/2026-05-08-translation-flow-navigation-overhaul/scenario-candidates.state-transition.md`
- `target_diff`: 翻訳セクション間の移動、フェーズページ直移動禁止、新規 job 作成直後の単語翻訳ページ遷移、未完了 job 一覧からの途中再開、成果物出力分離、旧 `Job Run` セッション取得廃止
- `viewpoint`: `state-transition`
