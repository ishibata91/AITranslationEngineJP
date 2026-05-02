# Scenario Candidates: term-translation-phase / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`
- `viewpoint`: `lifecycle`
- `task_artifact_dir`: `docs/exec-plans/active/term-translation-phase/`
- `target_diff`: 本文翻訳フェーズの前に、用語や固有名詞の訳語を確定し、ジョブ内辞書へ反映する単語翻訳フェーズを追加する。

## Generator Scope

- `viewpoint`: 作成、開始、対象語確定、実行、保存、完了、再開、再実行の時間順だけを扱う。
- `included_sources`: `./plan.md`, `tasks/usecases/term-translation-phase.yaml`, `tasks/index.yaml`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/architecture.md`, `docs/exec-plans/completed/translation-job-setup/plan.md`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md`, `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本化、他観点 generator の候補、最終シナリオ表。
- `generation_notes`: 候補は designer 統合前の母集団である。採否、統合、競合解消は行わない。

## Candidate Scenarios

### CAND-TTP-001 単語翻訳フェーズを開始して current phase と progress を観測する

- `source requirement`: `docs/exec-plans/active/term-translation-phase/plan.md:49-54`, `tasks/usecases/term-translation-phase.yaml:20-25`, `docs/spec.md:100-114`, `docs/spec.md:128-133`, `docs/er.md:61-69`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md:122-149`, `docs/exec-plans/completed/translation-job-setup/implementation-scope.md:21-29`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-001`
- `lifecycle stage`: `開始`
- `start condition`: `translation-job-setup` が完了し、翻訳ジョブが `Ready` として観測できる。
- `actor`: ユーザー
- `trigger`: Job Run で単語翻訳フェーズを開始する。
- `expected outcome`: 単語翻訳フェーズの `JOB_PHASE_RUN` が開始され、current phase と progress を確認できる。
- `observable point`: Job Run の current phase、progress、phase result、`JOB_PHASE_RUN.state`、`JOB_PHASE_RUN.progress_percent`
- `related detail requirement type`: `state_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: phase 開始の正常系候補として扱える。
- `conflict hint`: `state-transition` 観点が `Ready -> Running` の許可条件を別候補として切る可能性がある。

### CAND-TTP-002 共通辞書一致語を翻訳対象から除外する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:13-18`, `tasks/usecases/term-translation-phase.yaml:21-25`, `docs/spec.md:27-33`, `docs/spec.md:217-220`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-002`
- `lifecycle stage`: `対象語確定`
- `start condition`: 単語翻訳フェーズが開始され、翻訳対象語と共通辞書を参照できる。
- `actor`: システム
- `trigger`: 単語翻訳フェーズが対象語一覧を作る。
- `expected outcome`: 共通辞書に完全一致する語は単語翻訳対象から除外され、置換対象の判定結果として観測できる。
- `observable point`: phase result の除外語一覧、置換対象判定、内部 `cached` 相当の観測情報、対象語件数
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 共通辞書を先に適用する lifecycle 候補として扱える。
- `conflict hint`: `failure` 観点が共通辞書の重複、空訳語、不一致を別候補として切る可能性がある。

### CAND-TTP-003 用語や固有名詞の訳語を確定する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:9-18`, `tasks/usecases/term-translation-phase.yaml:21-25`, `docs/spec.md:33`, `docs/spec.md:224-228`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-003`
- `lifecycle stage`: `実行`
- `start condition`: 共通辞書除外後の翻訳対象語があり、単語翻訳フェーズが実行中である。
- `actor`: システム
- `trigger`: 単語翻訳フェーズが用語や固有名詞の訳語生成を実行する。
- `expected outcome`: 翻訳対象語ごとの確定訳語を phase result として確認できる。
- `observable point`: 確定訳語一覧、翻訳対象語件数、phase progress、AI 実行結果または fake 実行結果
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 訳語確定の主経路候補として扱える。
- `conflict hint`: `external-integration` 観点が AI provider、Batch API、fake transport を別候補として切る可能性がある。

### CAND-TTP-004 確定訳語をジョブ内辞書へ反映する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:15-18`, `tasks/usecases/term-translation-phase.yaml:23-25`, `docs/spec.md:217-220`, `docs/spec.md:226-228`, `docs/er.md:48-59`, `docs/diagrams/er/combined-data-model-er.puml:128-141`, `docs/diagrams/er/combined-data-model-er.puml:199-204`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-004`
- `lifecycle stage`: `保存`
- `start condition`: 訳語が確定し、ジョブ内辞書へ反映できる状態である。
- `actor`: システム
- `trigger`: 単語翻訳フェーズが確定訳語の保存を実行する。
- `expected outcome`: ジョブ内辞書として `DICTIONARY_ENTRY` が作られ、単語翻訳フェーズとの関連を確認できる。
- `observable point`: `DICTIONARY_ENTRY.translation_job_id`、`dictionary_lifecycle`、`dictionary_scope`、`dictionary_source`、`source_term`、`translated_term`、`reusable`、`PHASE_RUN_DICTIONARY_ENTRY`
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `observability_requirement`
- `adoption hint`: 確定訳語を本文翻訳へ渡す前の保存候補として扱える。
- `conflict hint`: `operation-audit` 観点が保存履歴、監査ログ、再現性を別候補として切る可能性がある。

### CAND-TTP-005 単語翻訳フェーズを完了し本文翻訳フェーズの入力として参照可能にする

- `source requirement`: `tasks/index.yaml:7-10`, `tasks/usecases/term-translation-phase.yaml:21-25`, `docs/spec.md:100-115`, `docs/spec.md:128-130`, `docs/spec.md:246-248`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-005`
- `lifecycle stage`: `完了`
- `start condition`: 確定訳語がジョブ内辞書へ反映され、単語翻訳フェーズの結果を表示できる。
- `actor`: システム
- `trigger`: 単語翻訳フェーズの完了処理を行う。
- `expected outcome`: 単語翻訳フェーズが完了し、本文翻訳フェーズが確定訳語を参照できる。
- `observable point`: phase result、完了 progress、ジョブ内辞書参照、後続フェーズの入力参照
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: フェーズ完了と後続フェーズ連携の候補として扱える。
- `conflict hint`: `state-transition` 観点が `term-translation-phase -> persona-generation-phase -> body-translation-phase` の遷移条件を別候補として切る可能性がある。

### CAND-TTP-006 翻訳対象語が残らない場合でも phase result を終了状態にする

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:21-25`, `docs/spec.md:29-33`, `docs/spec.md:128-130`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-006`
- `lifecycle stage`: `空対象完了`
- `start condition`: 共通辞書の完全一致除外後、単語翻訳対象語が 0 件になる。
- `actor`: システム
- `trigger`: 対象語確定処理が 0 件の翻訳対象を返す。
- `expected outcome`: AI による単語翻訳を実行せず、対象語なしの phase result と完了 progress を確認できる。
- `observable point`: 対象語 0 件、AI 実行なし、phase result、完了 progress、本文翻訳フェーズから参照できる共通辞書置換結果
- `related detail requirement type`: `boundary_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 境界値 lifecycle 候補として扱える。
- `conflict hint`: 対象語 0 件を完了扱いにするか、スキップ扱いにするかは designer 統合時に確認が必要である。

### CAND-TTP-007 単語翻訳フェーズを中断後に再開し、既存結果を重複させない

- `source requirement`: `docs/spec.md:53-59`, `docs/spec.md:133`, `docs/spec.md:155-184`, `docs/er.md:61-69`, `docs/diagrams/er/combined-data-model-er.puml:167-204`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-007`
- `lifecycle stage`: `再開`
- `start condition`: 単語翻訳フェーズが中断または回復可能失敗になり、同じジョブで再開できる。
- `actor`: ユーザー
- `trigger`: Job Run で単語翻訳フェーズを再開する。
- `expected outcome`: 同じ phase run の状態が戻り、既に反映済みの確定訳語を重複作成せず progress が継続する。
- `observable point`: `JOB_PHASE_RUN.state`、`progress_percent`、既存 `DICTIONARY_ENTRY` 件数、再開後の phase result
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `冪等性_requirement`
- `adoption hint`: 中断再開の lifecycle 候補として扱える。
- `conflict hint`: `state-transition` 観点が pause/resume 遷移を、`failure` 観点が recoverable failed retry を別候補として切る可能性がある。

### CAND-TTP-008 単語翻訳フェーズを再実行し、既存ジョブ内辞書の扱いを観測する

- `source requirement`: `docs/er.md:61-69`, `docs/spec.md:53-59`, `docs/spec.md:217-220`, `docs/diagrams/er/combined-data-model-er.puml:128-141`, `docs/diagrams/er/combined-data-model-er.puml:167-204`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TTP-008`
- `lifecycle stage`: `再実行`
- `start condition`: 単語翻訳フェーズが完了済み、または回復可能失敗から再実行準備に戻せる。
- `actor`: ユーザー
- `trigger`: Job Run で単語翻訳フェーズを再実行する。
- `expected outcome`: 既存ジョブ内辞書を維持、上書き、追加のどれで扱うかを phase result と辞書状態から確認できる。
- `observable point`: 再実行前後の `DICTIONARY_ENTRY`、`PHASE_RUN_DICTIONARY_ENTRY`、phase result、対象語件数、progress
- `related detail requirement type`: `冪等性_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 再実行と終了後利用の lifecycle 候補として扱える。
- `conflict hint`: 既存ジョブ内辞書の維持、上書き、追加の業務ルールは外部正本だけでは確定できない。

## Open Notes

- `human decision candidate`: `CAND-TTP-003` の「確定」は AI 出力の自動確定か、人間確認後の確定か。
- `human decision candidate`: `CAND-TTP-006` の対象語 0 件を `completed` として扱うか、`skipped` 相当の別結果にするか。
- `human decision candidate`: `CAND-TTP-007` の部分反映済み訳語を再開時に再利用するか、再検証対象にするか。
- `human decision candidate`: `CAND-TTP-008` の再実行時に既存ジョブ内辞書を維持、上書き、追加のどれで扱うか。
- `merge candidate`: `CAND-TTP-001` と state-transition 候補の `Ready -> Running` 条件。
- `merge candidate`: `CAND-TTP-005` と actor-goal 候補の本文翻訳フェーズ入力参照。
- `merge candidate`: `CAND-TTP-007` と failure / state-transition 候補の中断、回復可能失敗、再開。
- `rejection candidate`: なし。designer 統合時に最終判断する。
