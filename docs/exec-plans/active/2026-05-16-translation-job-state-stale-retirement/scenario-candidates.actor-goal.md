# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSR`

## Generator Scope

- `viewpoint`: `actor-goal`
- `running task artifact location`: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/`
- `target difference`: 翻訳ジョブ状態関連の stale 廃止。対象は、正本 state にない `pending`、操作可否の重複 state 知識、`JobIOService` の空実体と正本参照の不一致、active observability task-local の旧名参照、`cancelled` spelling 差分である。
- `candidate artifact path`: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-candidates.actor-goal.md`
- `included_sources`: `./implement-lane-task-frame.md`, `./plan.md`, `./state-knowledge-investigation.md`, `./state-knowledge-investigation-lane-decision.md`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本本文の変更、`docs/exec-plans/completed/**`、UI 変更前提、`stale_selection`、`validation_stale`、`model_selection_stale` の削除判断
- `generation_notes`: 候補は scenario-design の入力である。最終シナリオ表、採否、統合、競合解消は `designer` に残す。
- `candidate_count`: 7

## Candidate Scenarios

### CAND-TJSR-001 Ready job の開始判断で `pending` を利用者へ漏らさない

- `source requirement`: `Ready` job には `JOB_PHASE_RUN` を事前作成せず、フェーズ開始が許可された時だけ対象フェーズの `JOB_PHASE_RUN` を作成する。`pending` は正本 state に存在しないが、3 phase service に内部 state として残っている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-001`
- `actor`: 翻訳 job を実行する利用者
- `trigger`: 利用者が `Ready` job を Job Run の表示対象にし、最初のフェーズ開始可否を確認する。
- `expected outcome`: 利用者は `Ready` job が実行入口であることを判断できる。利用者は正本 state にない `pending` を phase state として見ない。開始不可の場合は、`pending` ではなく既存の実行不可理由カテゴリで判断できる。
- `observable point`: Job Run の current phase、phase state、開始操作の有効状態、無効理由、`TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の永続状態
- `related detail requirement type`: `state_requirement`, `success_requirement`, `compatibility_requirement`
- `adoption hint`: `pending` を canonical state へ昇格する候補ではなく、利用者の開始判断に stale state が混ざらないことを確認する候補として扱える。
- `conflict hint`: lifecycle 観点では、開始時に既存 `pending` run を継続するのか、仕様外 state として隔離するのかが競合しうる。最終判断は `designer` に残す。

### CAND-TJSR-002 phase 共通操作可否を利用者が同じ規則で判断できる

- `source requirement`: pause、resume、retry、cancel の可否は `JOB_PHASE_RUN.state` と共通操作規則から決める。phase 固有の `canRetry`、`canResume`、`canPause`、`canCancel` は持たない。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-002`
- `actor`: 翻訳 job を進める利用者
- `trigger`: 利用者が単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかの Job Run で pause、resume、retry、cancel の操作可否を見る。
- `expected outcome`: 利用者は phase の種類に依存しない同じ規則で操作可否を判断できる。`Running` では pause、`Paused` では resume と cancel、`RecoverableFailed` では retry を確認できる。無効な操作は理由とともに確認できる。
- `observable point`: Job Run の action enablement、無効理由、phase state label、operation result summary
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: service 層の重複 helper 廃止の実装判断ではなく、利用者が見る操作可否の一貫性として scenario-design へ渡せる。
- `conflict hint`: state-transition 観点では、read model の操作可否を `TranslationJobPolicy` 由来にする方法と service 層で複製する方法が競合しうる。

### CAND-TJSR-003 状態不整合時に危険操作を利用者へ許可しない

- `source requirement`: 保存済み `TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違う場合、表示だけで状態を書き換えず、危険操作を無効にする。phase progress 集約不能は成功値として表示せず、危険操作を無効にする。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-003`
- `actor`: 未完了 job を確認する利用者
- `trigger`: 利用者が未完了 job 一覧または Job Run を開き、保存済み job state と phase state の食い違いがある job を確認する。
- `expected outcome`: 利用者は対象 job を成功状態として誤認しない。利用者は削除、開始、再開、リトライ、取り消しなどの危険操作が無効である理由を確認できる。表示確認だけでは state が書き換わらない。
- `observable point`: 未完了一覧の job state、操作可否、無効理由、Job Run 表示対象、永続 state の変化なし
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 操作可否の正本化と `pending` の扱いが未決でも、利用者保護の受け入れ条件として独立候補にできる。
- `conflict hint`: failure 観点の参照不能や集約不能の扱いと統合される可能性がある。

### CAND-TJSR-004 運用確認者が状態事実の保存境界を誤認しない

- `source requirement`: `JobIOService` は architecture 正本と lint component に残るが、実体 package は `doc.go` だけである。`JobIOService` は job と phase run の状態取得と保存だけを扱い、遷移可否や UI 表示文言を判断しない。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-004`
- `actor`: 状態不整合を調査する運用確認者
- `trigger`: 運用確認者が job state、phase run state、進捗、失敗 reason category の保存または取得経路を追う。
- `expected outcome`: 運用確認者は状態事実の取得と保存の責務境界を判断できる。運用確認者は実体のない `JobIOService` を実装済み境界として誤認しない。状態可否、terminal guard、provider 応答検証が状態保存境界に混ざらないことを確認できる。
- `observable point`: architecture 正本、arch-lint component、active task-local の境界名、状態事実を扱う実体の有無、状態事実の保存結果
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: `JobIOService` を外すか実体化するかは決めず、運用確認者が調査経路を誤らないという actor goal として渡せる。
- `conflict hint`: responsibility-boundary 観点では、`JobIOService` を architecture 正本から外す案と別 task で実体化する案が競合しうる。

### CAND-TJSR-005 active observability task 再開時に旧名で状態境界を参照しない

- `source requirement`: `StateMachine` 旧名は product code から外れたが、active observability task-local に残っている。`JobIOService` 旧設計名も active observability task-local と architecture 正本に残っている。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-005`
- `actor`: active observability task を再開する運用確認者
- `trigger`: 運用確認者が `observability-log-addition` の active task-local を参照し、状態遷移、状態事実、通知観測点を確認する。
- `expected outcome`: 運用確認者は旧名の `StateMachine` を現在の状態境界として扱わない。運用確認者は `JobIOService` が実装済みの観測境界か、未決の設計主語かを区別できる。古い責務境界を次の観測ログ設計へ再注入しない。
- `observable point`: active task-local の scenario-design、scenario candidates、設計差分図、状態境界名、観測ログ追加時の対象境界
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: product code 変更ではなく、active task-local の旧名参照が次作業へ与える影響を scenario-design の入力にできる。
- `conflict hint`: operation-audit 観点では、旧名参照を今回更新するか、observability task 再開時に更新するかが競合しうる。

### CAND-TJSR-006 取消状態の spelling 差分で検索と検証を失敗させない

- `source requirement`: 正本仕様と service 実装は `Canceled` / `canceled` で統一されている。`PersonaGenerationPhaseContractStub` の cancel fixture は `cancelled` を返す。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-006`
- `actor`: 取消状態を確認する運用確認者
- `trigger`: 運用確認者が NPC ペルソナ生成フェーズの cancel 結果、fixture 応答、検索結果を確認する。
- `expected outcome`: 運用確認者は正本 spelling の `Canceled` / `canceled` だけを状態値として追跡できる。`cancelled` 由来の検索漏れや fixture だけの状態差分で、取消済み phase の原因調査を誤らない。
- `observable point`: stub 応答、contract fixture、phase state 表示、検索結果、取消後の terminal state
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `state_requirement`
- `adoption hint`: 利用者向け UI 変更ではなく、取消状態の検証と運用検索の一貫性として扱える。
- `conflict hint`: failure 観点では、取消状態の spelling 差分をテスト fixture 修正だけで閉じるか、仕様用語の確認まで広げるかが競合しうる。

### CAND-TJSR-007 利用者向け stale reason を stale 廃止対象と混同しない

- `source requirement`: この task の `stale` は古い設計名、使われていない package、同じ規則を繰り返す wrapper、完了済み判断とずれた active task-local 参照を意味する。`stale_selection`、`validation_stale`、`model_selection_stale` は利用者向けまたは API 向け理由分類として廃止対象にしない。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJSR-007`
- `actor`: 未完了 job 一覧と Job Run を確認する利用者
- `trigger`: 利用者が読み込み失敗、stale selection、設定検証 stale、model selection stale に関係する理由表示または API 結果を見る。
- `expected outcome`: 利用者は stale 廃止後も、利用者向け stale reason を確認できる。利用者向け stale reason は空状態や成功状態と区別される。状態関連の stale 廃止が、既存の理由分類を消さない。
- `observable point`: 未完了一覧の理由カテゴリ、Job Run の無効理由、API response の reason category、状態変更なし
- `related detail requirement type`: `compatibility_requirement`, `failure_handling_requirement`, `state_requirement`
- `adoption hint`: 廃止禁止対象の保護シナリオとして、今回の stale 廃止が利用者向け理由分類を壊さないことを確認できる。
- `conflict hint`: failure 観点や operation-audit 観点の stale 判定候補と統合される可能性がある。

## Open Notes

- `human decision candidate`: `pending` を canonical state へ昇格するか、内部一時 state として隔離するか。
- `human decision candidate`: `TranslationJobPolicy` の共通操作規則を read model の操作可否へどう再利用するか。
- `human decision candidate`: `JobIOService` を architecture 正本から外すか、別 task で実体化するか。
- `human decision candidate`: `observability-log-addition` の `StateMachine` / `JobIOService` 旧名参照を今回の active task-local 更新対象に含めるか。
- `human decision candidate`: `cancelled` fixture spelling を今回の stale 廃止に含めるか。
- `merge candidate`: `CAND-TJSR-001` と lifecycle 観点の phase start 候補は、開始時の `JOB_PHASE_RUN` 作成条件として統合される可能性がある。
- `merge candidate`: `CAND-TJSR-002` と state-transition 観点の pause、resume、retry、cancel 候補は、共通操作規則の受け入れ条件として統合される可能性がある。
- `merge candidate`: `CAND-TJSR-004` と `CAND-TJSR-005` は、運用確認者が古い責務境界を参照しない候補として統合される可能性がある。
- `rejection candidate`: `CAND-TJSR-007` は、他観点で既存 reason category 保護が十分に扱われる場合、重複候補として不採用になる可能性がある。
