# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM`
- `candidate_count`: 10

## Generator Scope

- `viewpoint`: 利用者、運用確認者、実装者が、翻訳ジョブ状態機械に期待する目的、開始操作、成功結果を候補化する。
- `included_sources`: `plan.md`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/architecture.md`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、未承認の docs 正本化、最終シナリオ表、候補の採否、候補の統合判断。
- `generation_notes`: `actor-goal` 観点だけを扱う。状態遷移網羅、外部 provider 失敗、運用監査の網羅は別観点へ残す。
- `adopted_update`: 人間回答により、retry、resume、pause、cancel の可否は phase type で分けない。phase type で分ける対象は、start の開始前提、完了判定、呼び出す service method だけである。

## Candidate Scenarios

### CAND-TJSM-001 利用者が入力データから実行可能な Ready job を作成する

- `source requirement`: `translation-job-setup` は、登録済み入力データと 3 つの翻訳段階の AI 設定を固定し、`Ready` job を作成できることを求めている。`spec.md` は `Draft` から `Ready` への遷移をジョブ作成として定義している。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-001`
- `actor`: 翻訳入力データから翻訳 job を作成したい利用者。
- `trigger`: 利用者が Job Setup で入力データ、単語翻訳設定、NPC ペルソナ生成設定、本文翻訳設定を確認し、job 作成を実行する。
- `expected outcome`: 対象入力データに紐づく `Ready` job が作成される。Job Run 表示だけでは `Running` へ暗黙遷移しない。API key 本文、credential 参照実値、endpoint は表示されない。
- `observable point`: Job Management または Job Run で、job state、入力出自、AI 設定要約、credential 状態分類を確認できる。active な phase run は存在しない状態として開始可否を判定できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: Job Setup 作成後の job 初期状態と、Job Run 表示時の非遷移を固定する候補として扱える。
- `resolved conflict`: `Ready` job には phase run を事前作成しない。

### CAND-TJSM-002 利用者が Ready job から単語翻訳フェーズを開始する

- `source requirement`: `term-translation-phase` は、対象ジョブが `Ready` であり、active な単語翻訳 phase run がない時だけ開始できることを求めている。`spec.md` は `Ready` から `Running` への遷移を実行開始として定義している。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-002`
- `actor`: 翻訳 job を実行する利用者。
- `trigger`: 利用者が Job Run で `Ready` job の単語翻訳フェーズ開始を実行する。
- `expected outcome`: 単語翻訳フェーズが開始され、job は実行中として扱われる。開始不可条件がある場合は、開始されず、理由を確認できる。
- `observable point`: Job Run に current phase、phase state、progress、対象語件数、provider / model / execution mode / batch mode、credential 状態分類が表示される。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: `translationjobpolicy` が共通操作規則と start 開始前提を判定し、`JobIOService` が状態取得と保存を担う設計の代表候補として扱える。
- `conflict hint`: `docs/er.md` は job 状態を `JOB_PHASE_RUN` 群から集約すると定義する一方、`docs/spec.md` は job state を直接列挙している。

### CAND-TJSM-003 利用者が単語翻訳完了後に NPC ペルソナ生成フェーズへ進む

- `source requirement`: `persona-generation-phase` は、単語翻訳フェーズが `Completed`、job が terminal state ではない、active phase run がない場合だけ開始できることを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-003`
- `actor`: 翻訳 job を順番に進めたい利用者。
- `trigger`: 利用者が単語翻訳フェーズ完了後、Job Run で NPC ペルソナ生成フェーズ開始を実行する。
- `expected outcome`: NPC ペルソナ生成フェーズが開始される。単語翻訳フェーズ未完了、terminal job、active phase run ありの場合は開始されない。
- `observable point`: Job Run に current phase、phase state、progress、target count、generated count、failed count、skipped count、snapshot 参照状態、本文翻訳フェーズ readiness が表示される。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: phase 間の順序制約を利用者の次操作として確認する候補として扱える。
- `conflict hint`: terminal state の範囲が一覧表示、再開可否、後書き拒否で同じかは最終設計で確認が必要である。

### CAND-TJSM-004 利用者が NPC ペルソナ生成完了後に本文翻訳フェーズへ進む

- `source requirement`: `body-translation-phase` は、NPC ペルソナ生成フェーズが `Completed`、job が terminal ではない、active phase run がない、辞書と persona snapshot の参照が成立している場合だけ開始できることを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-004`
- `actor`: Skyrim Mod 翻訳 job の本文翻訳を実行したい利用者。
- `trigger`: 利用者が Job Run で本文翻訳フェーズ開始を実行する。
- `expected outcome`: 本文翻訳フェーズが開始される。本文翻訳フェーズ完了時に、job 全体は `Completed` となり、output readiness を確認できる。
- `observable point`: Job Run に current phase、phase state、progress、field result summary、保護要素検証結果、output readiness、redacted error summary が表示される。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: job 完了状態と成果物出力準備をつなぐ利用者目的の候補として扱える。
- `conflict hint`: body phase `Completed` と job-level `Completed` のどちらが先に成立するかは、状態正本の設計で確認が必要である。

### CAND-TJSM-005 利用者が Paused または RecoverableFailed の job を再開またはリトライする

- `source requirement`: `spec.md` は `Paused` から `Running`、`RecoverableFailed` から `Running` または `Ready` への回復経路を定義している。各 phase 詳細仕様は、再送、再開、リトライでは同じ `JOB_PHASE_RUN` を継続することを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-005`
- `actor`: 中断または回復可能失敗から作業を続けたい利用者。
- `trigger`: 利用者が Job Management または Job Run で再開、リトライ、再実行準備の操作を実行する。
- `expected outcome`: 成功済みの辞書、persona、field result は重複作成されず、未処理対象だけが継続される。再開不可の場合は job 状態を変えず、理由を確認できる。
- `observable point`: 同じ `JOB_PHASE_RUN` の progress、retryable flag、error kind、再開不可理由、現在フェーズを確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `state_requirement`, `recovery_requirement`, `consistency_requirement`, `冪等性_requirement`
- `adoption hint`: 利用者にとっての回復成功と、実装者にとっての同一 phase run 継続条件を接続する候補として扱える。
- `resolved conflict`: retry と resume は同じ `JOB_PHASE_RUN` を継続する。start 再送は active phase run と同一再送判定で扱う。

### CAND-TJSM-006 利用者が停止またはキャンセルできる条件を確認する

- `source requirement`: `spec.md` は `Running` から `Paused`、`Ready` から `Canceled`、`Paused` から `Canceled` への操作系遷移を定義している。`body-translation-phase` は cancel を `Paused` からだけ可能にし、`Running` から直接 cancel しないことを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-006`
- `actor`: 実行中または停止中の job を安全に止めたい利用者。
- `trigger`: 利用者が Job Run または Job Management で pause、resume、cancel の操作可否を確認し、可能な操作を実行する。
- `expected outcome`: 許可された停止またはキャンセルだけが実行される。`Canceled` 後はフェーズ終端となり、途中成功結果は output readiness に使われない。
- `observable point`: 操作ボタンの有効状態、無効理由、phase state、job state、output readiness の無効理由を確認できる。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: 利用者操作から見た cancel 可能条件を候補化し、状態遷移候補との競合検出に渡せる。
- `resolved conflict`: Ready cancel は job-level、phase 開始後 cancel は Paused phase-level で固定する。

### CAND-TJSM-007 利用者が未完了一覧で job 状態と操作可否を比較する

- `source requirement`: `translation-job-management` は、Completed 以外の翻訳 job を一覧し、状態、操作可否、理由カテゴリ、Job Run 導線を確認できることを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-007`
- `actor`: 作成済みの未完了翻訳 job から次に扱う job を選びたい利用者。
- `trigger`: 利用者が翻訳管理を開き、未完了一覧を読み込む。
- `expected outcome`: `Ready`、`Running`、`Paused`、`RecoverableFailed`、`Failed`、`Canceled` の job が一覧対象になる。`Completed` job は未完了一覧に表示されない。一覧表示だけでは job 状態を変えない。
- `observable point`: 未完了一覧で job state、現在フェーズ、進捗、入力出自、操作可否、無効理由、AI 設定要約を比較できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: `translationjobpolicy` の結果を UI 表示値へ写す候補として扱える。
- `conflict hint`: `Failed` と `Canceled` が未完了一覧に残る扱いと、terminal state として後書きを拒否する扱いの関係は確認が必要である。

### CAND-TJSM-008 利用者が job 削除可否を状態から判断する

- `source requirement`: `translation-job-management` は、`Running` job を削除できず、非実行中 job を削除しても input data と抽出 JSON 正本を残すことを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-008`
- `actor`: 不要な未完了 job を整理したい利用者。
- `trigger`: 利用者が Job Management で対象 job の削除を実行する、または削除不可理由を確認する。
- `expected outcome`: `Running` job の削除は拒否され、停止入口が表示される。非実行中 job の削除が成功した場合、対象 job は未完了一覧から外れ、入力データは保持される。
- `observable point`: 削除ボタンの有効状態、削除拒否理由、停止入口、削除後の一覧状態、入力データの保持を確認できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: `translationjobpolicy` が危険操作を拒否する利用者目的の候補として扱える。
- `conflict hint`: 停止要求中の削除可否は、job state と phase state の集約規則が不明だと判定が揺れる可能性がある。

### CAND-TJSM-009 利用者が Completed job から翻訳成果物を出力する

- `source requirement`: `translation-output-artifact` は、body phase が `Completed` であり、job-level 状態も `Completed` である翻訳 job だけを出力候補にすることを求めている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-009`
- `actor`: Skyrim Mod 翻訳成果物を確認して xTranslator 互換 XML を出力したい利用者。
- `trigger`: 利用者が Output Review で completed job を選択し、出力準備状態を確認する。
- `expected outcome`: output readiness、result summary、output status distribution、diff preview、artifact status を確認できる。未完了、失敗中、`Canceled`、field result 不整合、status 不整合の job では成果物生成を開始できない。
- `observable point`: completed job list、selected job summary、拒否理由、row count、file path、generated_at、compatibility summary を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: job-level `Completed` と body phase `Completed` の整合を確認する利用者目的の候補として扱える。
- `conflict hint`: `Completed` job は未完了一覧に表示しないため、Output Review 側の一覧と Job Management 側の一覧は別の状態解釈を持つ。

### CAND-TJSM-010 実装者と運用確認者が状態判定の責務境界を確認する

- `source requirement`: `architecture.md` は既存名 `StateMachine` が状態遷移規則だけを保持し、`JobIOService` が job 状態の取得と保存だけを扱うと定義している。今回の plan は既存名を `translationjobpolicy` に置き換える。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSM-010`
- `actor`: 状態判定の責務境界を守りたい実装者と、状態不整合を安全側に確認したい運用確認者。
- `trigger`: 実装者が job 操作の usecase を設計または確認し、運用確認者が集約不能または参照不能 job の表示を確認する。
- `expected outcome`: 状態遷移の許可、拒否、拒否理由は状態規則として一貫する。状態取得と保存は I/O 責務として分離される。集約不能または参照不能の場合は成功状態として表示されず、危険操作が無効になる。
- `observable point`: UseCase から見て、状態判定結果、永続化対象、表示用理由カテゴリを分けて確認できる。Job Management では stale selection、参照不能、phase progress 集約不能を空状態や成功状態と区別できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`, `failure_handling_requirement`
- `adoption hint`: 実装者を actor とする責務分離の候補であり、最終シナリオでは lower-level only または API 境界の検証候補へ変換される可能性がある。
- `resolved conflict`: `translationjobpolicy` は操作可否、拒否理由、状態作用、service command を返す。`PolicyResult` は保存せず、UseCase が確定済み状態事実へ変換してから `JobIOService` へ渡す。

## Open Notes

- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `resolved decision`: `Ready` job に phase run を事前作成しない。
- `resolved decision`: `pending` は公開仕様上の状態に含めない。
- `resolved decision`: retry と resume は同じ `JOB_PHASE_RUN` を継続する。
- `resolved decision`: cancel は `Ready` job-level または `Paused` phase-level で扱う。
- `resolved decision`: `translationjobpolicy` は操作可否、拒否理由、状態作用を返す。
- `merge candidate`: CAND-TJSM-002、CAND-TJSM-003、CAND-TJSM-004 は、designer がフェーズ順序の最終シナリオへ統合する可能性がある。
- `merge candidate`: CAND-TJSM-005 と CAND-TJSM-006 は、designer が回復操作と停止操作を分けるか、同じ操作系シナリオへ統合するかを判断する可能性がある。
- `merge candidate`: CAND-TJSM-007、CAND-TJSM-008、CAND-TJSM-009 は、画面別の状態表示シナリオとして分けるか、状態表示の acceptance group として統合する可能性がある。
- `rejection candidate`: 実装者 actor の CAND-TJSM-010 は、最終利用者シナリオだけに絞る方針になった場合、designer が implementation-scope 側の責務境界確認へ移す可能性がある。
