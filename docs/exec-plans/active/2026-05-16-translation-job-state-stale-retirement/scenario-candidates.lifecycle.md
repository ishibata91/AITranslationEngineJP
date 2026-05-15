# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSR`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./implement-lane-task-frame.md`, `./plan.md`, `./backend-implementation-result.md`, `./state-knowledge-investigation.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本本文の変更、`docs/exec-plans/completed/**`、UI 変更前提、`stale_selection`、`validation_stale`、`model_selection_stale` の削除
- `generation_notes`: lifecycle 段階ごとに候補を分ける。最終採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-TJSR-001 Ready job 作成時に phase run を事前作成しない

- `source requirement`: `docs/spec.md` は、`Ready` job には `JOB_PHASE_RUN` を事前作成せず、フェーズ開始が許可された時だけ対象フェーズの `JOB_PHASE_RUN` を作成すると定義している。`docs/detail-specs/translation-job-management.md` も同じ前提を持つ。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-001`
- `actor`: 翻訳 job を作成する利用者
- `trigger`: 入力データから新しい翻訳 job を作成する。
- `expected outcome`: `TRANSLATION_JOB.state` は `Ready` になる。active な `JOB_PHASE_RUN` は存在しない。Job Run 表示だけでは `Ready` job を `Running` へ暗黙遷移させない。
- `observable point`: job 一覧または Job Run の read model で、job state は `Ready`、現在フェーズは未開始、危険操作は安全側の可否になる。DB 観測では対象 job に事前作成された phase run がない。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: `pending` を canonical state にしない案でも、内部一時 state として隔離する案でも、作成直後の `Ready` job に `JOB_PHASE_RUN` を作らない受け入れ条件として残せる。
- `conflict hint`: 既存実装やテストが作成直後の `pending` phase run を期待する場合、正本仕様と衝突する。

### CAND-TJSR-002 フェーズ開始許可時だけ phase run を作成して Running へ進める

- `source requirement`: `docs/spec.md` は、開始許可時に `JOB_PHASE_RUN` を作成し、phase state は `Running` から始まると定義している。各 phase 詳細仕様は、開始条件成立時だけ対象 phase run を作成すると定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-002`
- `actor`: 翻訳フェーズを開始する利用者
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかで開始条件を満たした状態から start を実行する。
- `expected outcome`: 対象 phase の `JOB_PHASE_RUN` が作成または継続され、利用者に観測される phase state は `Running` になる。job 全体は phase 実行中として扱われる。phase 固有の開始前提だけが phase type ごとの差分になる。
- `observable point`: command 結果、read model、DB 状態で、開始対象 phase、job state、phase state、active phase run、開始時刻、進捗初期値を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: service 層に残る `pending` を廃止する場合、start 後の最初の永続 state が `Running` であることを検証する候補になる。
- `conflict hint`: `pending` を内部一時 state として残す場合、外部観測点では `pending` が見えないことを別候補で保証する必要がある。

### CAND-TJSR-003 `pending` を canonical state にしない受け入れ観点を固定する

- `source requirement`: `docs/spec.md` の job state と phase state に `pending` は存在しない。`state-knowledge-investigation.md` は、3 phase service に `pending` が内部 state として残っていると観測している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-003`
- `actor`: lifecycle 仕様を確認する設計者
- `trigger`: `pending` を正本 state へ昇格しない方針で lifecycle を設計する。
- `expected outcome`: `pending` は `TRANSLATION_JOB.state`、`JOB_PHASE_RUN.state`、read model の phase state、operation summary、Wails DTO、詳細仕様の state 一覧に出ない。start、pause、resume、retry、cancel、完了、失敗の受け入れ条件は `Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` だけで記述される。
- `observable point`: state 定義、scenario、read model、operation result、DB 検証で `pending` が canonical state として観測されない。`stale_selection`、`validation_stale`、`model_selection_stale` は削除対象から除外されたまま残る。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`, `consistency_requirement`
- `adoption hint`: stale 廃止の中心候補として使える。既存の正本仕様を広げず、状態語彙の増加を避ける。
- `conflict hint`: service 内の既存 `pending` 分岐を残す場合、この候補だけでは実装との整合を証明できない。

### CAND-TJSR-004 `pending` を内部一時 state として隔離する受け入れ観点を固定する

- `source requirement`: `state-knowledge-investigation.md` は、`pending` が phase run 作成直後の一時永続 state として温存されている可能性を挙げている。`docs/architecture.md` は、policy 判定結果や rule 名を DB、DTO、repository 永続契約へ出さないと定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-004`
- `actor`: lifecycle 仕様を確認する設計者
- `trigger`: `pending` を canonical state へ昇格せず、内部一時 state として隔離する方針で lifecycle を設計する。
- `expected outcome`: `pending` は利用者向け state、公開 DTO、詳細仕様の state 一覧、共通操作規則の判定入力に出ない。内部で使う場合も、start 処理の境界内で `Running`、失敗、またはロールバックへ解決される。`Ready` job の事前 phase run 作成には使わない。
- `observable point`: command 完了後の read model と DB の安定状態で `pending` が残らない。内部一時 state を許容する場合は、テストが一時 state の露出禁止、terminal job への後書き拒否、危険操作無効化を確認する。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `concurrency_requirement`, `testability_requirement`
- `adoption hint`: 既存実装の中間処理をすぐ消さない設計にする場合の候補になる。
- `conflict hint`: `pending` が DB に安定状態として残る、read model に出る、操作可否の入力になる、といういずれかが起きる場合、正本仕様と衝突する。

### CAND-TJSR-005 Running phase を pause して Paused へ移す

- `source requirement`: `docs/spec.md` は、`Running` の `JOB_PHASE_RUN` だけを pause できると定義している。各 phase 詳細仕様も pause は `Running` の時だけ有効と定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-005`
- `actor`: 実行中 phase を中断する利用者
- `trigger`: 対象 phase が `Running` の状態で pause を実行する。
- `expected outcome`: 対象 phase run は `Paused` になる。job は利用者が resume または cancel を判断できる状態になる。phase 固有の `canPause` 分岐は持たず、共通操作規則で判断できる。
- `observable point`: command 結果と read model で、phase state、progress、pause 後の resume 可否、cancel 可否、pause 不可理由を確認できる。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: `commonPhaseActionAvailability` と `TranslationJobPolicy` の重複判断を designer が扱う時の lifecycle 候補になる。
- `conflict hint`: service helper が policy と異なる理由分類を返す場合、contract 観点または state-transition 観点との競合候補になる。

### CAND-TJSR-006 Paused phase を resume して同じ phase run を継続する

- `source requirement`: `docs/spec.md` は、`Paused` の `JOB_PHASE_RUN` だけを resume でき、retry、resume、開始再送は同じ `JOB_PHASE_RUN` を継続すると定義している。各 phase 詳細仕様も同じ phase run の継続と重複作成禁止を定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-006`
- `actor`: 中断中 phase を再開する利用者
- `trigger`: 対象 phase が `Paused` の状態で resume を実行する。
- `expected outcome`: 既存の `JOB_PHASE_RUN` が `Running` へ戻る。既存 progress、成功済み結果、phase run id は継続する。新しい phase run や重複 entry は作らない。
- `observable point`: command 結果、DB、read model で、phase run id の継続、state の `Paused` から `Running` への変化、重複作成なし、resume 不可理由を確認できる。
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `data_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: 「resume」と「開始再送」の lifecycle 差分を designer が整理する材料になる。
- `conflict hint`: `pending` を内部一時 state として使う場合、resume は `pending` を経由するかどうかを別途決める必要がある。

### CAND-TJSR-007 RecoverableFailed phase を retry して同じ phase run を継続する

- `source requirement`: `docs/spec.md` は、`RecoverableFailed` の `JOB_PHASE_RUN` だけを retry でき、`RecoverableFailed` から `Ready` へ戻す経路を作らないと定義している。各 phase 詳細仕様は retry で同じ phase run を継続し、未処理対象だけを進めると定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-007`
- `actor`: 回復可能失敗から再試行する利用者
- `trigger`: provider 失敗、応答不正、保存失敗、保護要素検証失敗などで phase が `RecoverableFailed` になった後、retry を実行する。
- `expected outcome`: 既存の `JOB_PHASE_RUN` が `Running` へ戻る。成功済み結果は維持され、未処理対象だけが retry 対象になる。`Ready` への巻き戻しや新規 phase run 作成は起きない。
- `observable point`: command 結果、DB、read model で、phase run id の継続、成功済み結果の維持、未処理対象数、retryable flag、retry 不可理由を確認できる。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `state_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: stale 廃止後も、旧 state 名を経由せず `RecoverableFailed` を起点に回復できることを確認する候補になる。
- `conflict hint`: `pending` または read model 用派生 state が retry の入力に混ざる場合、canonical state と内部表示差分の境界が競合しうる。

### CAND-TJSR-008 Paused phase を cancel して終端化する

- `source requirement`: `docs/spec.md` は、phase 開始後の cancel は `Paused` の対象フェーズからだけ許可すると定義している。本文翻訳詳細仕様は、`Canceled` 後はフェーズ終端とし、途中成功結果は output readiness に使わないと定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-008`
- `actor`: 中断中 phase を取り消す利用者
- `trigger`: 対象 phase が `Paused` の状態で cancel を実行する。
- `expected outcome`: 対象 phase run は `Canceled` になる。以後の resume、retry、late response 後書き、output readiness 更新は拒否される。Running から直接 cancel はできない。
- `observable point`: command 結果、read model、DB で、phase state の `Canceled`、terminal guard、危険操作の無効化、途中成功結果が後続 readiness に使われないことを確認できる。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `security_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: `cancelled` fixture spelling を `canceled` へ寄せるかどうかを designer が判断する材料になる。
- `conflict hint`: `Canceled` と `cancelled` が混在する場合、検索漏れと state 比較失敗の競合候補になる。

### CAND-TJSR-009 body phase 完了で job 全体を Completed にする

- `source requirement`: `docs/spec.md` は、`Running` から `Completed` への job 遷移を本文翻訳完了で行うと定義している。本文翻訳詳細仕様は、body phase Completed、field result 整合、output status 整合を満たす時だけ output readiness を true にすると定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-009`
- `actor`: 本文翻訳を完了させる利用者
- `trigger`: 本文翻訳フェーズの全対象が成功し、保存と整合性確認が完了する。
- `expected outcome`: body phase は `Completed` になる。job 全体は `Completed` になる。Completed job は未完了一覧から外れ、成果物出力へ進める状態になる。
- `observable point`: read model、DB、operation result で、body phase state、job state、field result 整合、output readiness、未完了一覧からの除外を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: lifecycle の終点を `Completed` として固定し、stale 廃止が完了条件を変えていないことを示す候補になる。
- `conflict hint`: detail-spec の read model 用派生 state `empty completed` と、永続 state `Completed` の境界を designer が分ける必要がある。

### CAND-TJSR-010 回復不能失敗で Failed へ終端化する

- `source requirement`: `docs/spec.md` は、`Running` から `Failed` への回復不能失敗遷移と、`Failed` を terminal state として定義している。各 phase 詳細仕様は、provider 失敗、応答不正、保存失敗、検証失敗を successful Completed として扱わないと定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-010`
- `actor`: 失敗状態を確認する利用者
- `trigger`: 回復不能な失敗、または retry できない失敗が phase 実行中に発生する。
- `expected outcome`: 対象 phase run または job は `Failed` へ終端化する。terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。失敗理由は理由カテゴリとして観測できる。
- `observable point`: command 結果、read model、DB、structured summary で、`Failed`、terminal guard、retry 不可理由、成功扱いされない partial state を確認できる。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: stale 廃止が失敗終端の意味を変えないことを確認する候補になる。
- `conflict hint`: 回復可能失敗を `Failed` と同一視する実装がある場合、retry lifecycle と衝突する。

### CAND-TJSR-011 terminal job への後書きを拒否する

- `source requirement`: `docs/spec.md` は、terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否すると定義している。本文翻訳詳細仕様と NPC ペルソナ生成詳細仕様も terminal job への後書き拒否を定義している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-011`
- `actor`: 実行結果の完了後応答を処理するシステム
- `trigger`: job または phase が `Completed`、`Failed`、`Canceled` になった後に、遅れて provider 応答や保存要求が届く。
- `expected outcome`: terminal state は維持される。late response は保存、readiness 更新、phase run 作成に使われない。利用者向け summary には secret、API key、provider raw payload を含めない。
- `observable point`: DB、operation result、structured summary で、terminal state の不変、後書き拒否、危険操作無効化、redaction を確認できる。
- `related detail requirement type`: `state_requirement`, `security_requirement`, `consistency_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: `pending` を内部一時 state として残す場合でも、terminal 後の内部 state 後書きがないことを保証する候補になる。
- `conflict hint`: external-integration 観点では provider の late response 扱いと競合または統合対象になる可能性がある。

### CAND-TJSR-012 状態不整合の再表示で状態を書き換えず危険操作を無効化する

- `source requirement`: `docs/detail-specs/translation-job-management.md` は、保存済み `TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違う場合、表示だけで状態を書き換えず、危険操作を無効にすると定義している。`state-knowledge-investigation.md` は、`pending` と read model 用派生 state の混在を観測している。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TJSR-012`
- `actor`: Job Run または未完了一覧を再表示する利用者
- `trigger`: 保存済み job state、phase state、read model 集約結果に食い違いがある状態で一覧または Job Run を読み込む。
- `expected outcome`: 読み込みだけでは状態を修復または変更しない。危険操作は無効になる。実行不可理由は理由カテゴリとして表示できる。再開、retry、cancel の入口は正本 state と整合する場合だけ有効になる。
- `observable point`: read model、operation availability、DB で、読み込み前後の永続 state 不変、危険操作無効化、理由カテゴリ、`pending` または派生 state の露出有無を確認できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `recovery_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 「再開」を実行操作だけでなく、既存 job の再表示と操作再判定の lifecycle として扱う候補になる。
- `conflict hint`: UI 変更を前提にせず検証する必要がある。UI 変更が必要になる場合は、この task の禁止範囲と衝突する。

## Open Notes

- `human decision candidate`: `pending` を canonical state へ昇格するか、内部一時 state として隔離するかは人間判断候補である。候補 `CAND-TJSR-003` と `CAND-TJSR-004` は両案の受け入れ観点として残す。
- `human decision candidate`: `JobIOService` を architecture 正本から外すか、別 task で実体化するかは lifecycle 候補だけでは確定しない。
- `human decision candidate`: active task-local に残る `StateMachine` / `JobIOService` 旧名参照を今回の stale 廃止に含めるかは、人間判断候補である。
- `merge candidate`: `CAND-TJSR-005`、`CAND-TJSR-006`、`CAND-TJSR-007`、`CAND-TJSR-008` は、共通操作規則の受け入れシナリオとして統合される可能性がある。
- `merge candidate`: `CAND-TJSR-010` と `CAND-TJSR-011` は、terminal guard の受け入れシナリオとして統合される可能性がある。
- `rejection candidate`: `pending` を正本 state として追加する候補は、現行 `docs/spec.md` と衝突するため、この lifecycle 候補では採用前提にしない。

