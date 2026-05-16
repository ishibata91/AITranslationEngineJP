# 状態知識追加調査

## 判断結果

- 判定: 完了
- 結論: 現在の削減範囲だけでは stale 廃止は閉じない。`StateMachine` 旧名、`JobIOService` 旧設計名、`pending` を含む phase state 派生知識、read model 用の派生 state 名が複数箇所へ残っている。
- 引き継ぎ先: `designer`

## 調査 mode

- `リスク報告`

## 観測点

- `docs/spec.md` の `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` 正本定義を確認した。参照: `docs/spec.md:137-218`
- `docs/architecture.md`、`docs/diagrams/backend/backend-architecture.puml`、`.go-arch-lint.yml` の構造主語と依存主語を確認した。参照: `docs/architecture.md:30-31`, `docs/architecture.md:53`, `docs/architecture.md:160-186`, `docs/architecture.md:235-256`, `docs/diagrams/backend/backend-architecture.puml:56-69`, `docs/diagrams/backend/backend-architecture.puml:146-150`, `.go-arch-lint.yml:11-16`, `.go-arch-lint.yml:48-59`
- `internal/usecase/translationjobpolicy/policy.go`、`internal/usecase/phase_policy_helpers.go`、`internal/service/phase_action_enablement_helpers.go` を確認した。参照: `internal/usecase/translationjobpolicy/policy.go:22-27`, `internal/usecase/translationjobpolicy/policy.go:67-123`, `internal/usecase/phase_policy_helpers.go:10-39`, `internal/service/phase_action_enablement_helpers.go:5-77`
- phase service の state 定義と start 直後処理を確認した。参照: `internal/service/term_translation_phase_service.go:17-34`, `internal/service/term_translation_phase_service.go:697-710`, `internal/service/body_translation_phase_service.go:19-29`, `internal/service/body_translation_phase_service.go:367-394`, `internal/service/body_translation_phase_service.go:487-492`, `internal/service/persona_generation_phase_service.go:22-36`, `internal/service/persona_generation_phase_service.go:375-386`, `internal/service/persona_generation_phase_service.go:1476-1482`
- detail-spec と active task-local の state 語彙を確認した。参照: `docs/detail-specs/translation-job-management.md:26-57`, `docs/detail-specs/term-translation-phase.md:28-35`, `docs/detail-specs/term-translation-phase.md:76-85`, `docs/detail-specs/persona-generation-phase.md:29-48`, `docs/detail-specs/persona-generation-phase.md:59-69`, `docs/detail-specs/body-translation-phase.md:28-46`, `docs/detail-specs/body-translation-phase.md:74-84`, `docs/exec-plans/active/observability-log-addition/scenario-design.md:107-110`, `docs/exec-plans/active/observability-log-addition/scenario-candidates.operation-audit.md:19-25`
- `internal/jobio/` と `internal/statemachine/` の filesystem 状態を確認した。2026-05-16 の listing では `internal/jobio/` は `doc.go` だけ、`internal/statemachine/` は空 directory だった。

## 観測事実

- 正本仕様は job state を `Draft`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` に限定し、phase state を `Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` に限定している。`pending` は正本図に存在しない。参照: `docs/spec.md:141-218`
- 実装側 service は 3 phase とも `pending` を内部 state として保持している。単語翻訳は `pending -> running` 更新を持つ。本文翻訳は `bodyRun == nil || pending` を start 分岐に含め、既存 `pending` run を更新する。NPC ペルソナ生成も `pending` run を再利用し、`UpdateJobPhaseRunWhenState(..., pending, ...)` を持つ。参照: `internal/service/term_translation_phase_service.go:21-27`, `internal/service/term_translation_phase_service.go:697-710`, `internal/service/body_translation_phase_service.go:19-25`, `internal/service/body_translation_phase_service.go:367-394`, `internal/service/body_translation_phase_service.go:487-492`, `internal/service/persona_generation_phase_service.go:23-32`, `internal/service/persona_generation_phase_service.go:375-386`, `internal/service/persona_generation_phase_service.go:1476-1482`
- read model 用の派生 state 名が phase ごとに増えている。単語翻訳 detail-spec は `idle_ready`、`empty_completed`、`blocked` を `phase state` と同じ節で列挙している。本文翻訳 detail-spec は `not-ready`、`ready`、`starting`、`validation failed`、`empty completed` を列挙している。NPC ペルソナ生成 service は `not_started`、`rejected`、`snapshot_missing`、`empty_completed` を実装している。正本仕様の `JOB_PHASE_RUN.state` 節にはこれらが存在しない。参照: `docs/detail-specs/term-translation-phase.md:83-85`, `docs/detail-specs/body-translation-phase.md:78-84`, `internal/service/persona_generation_phase_service.go:22-32`, `internal/service/persona_generation_phase_service.go:565-620`, `internal/service/persona_generation_phase_service.go:2403-2411`
- 共通操作規則は `TranslationJobPolicy` に正本化されているが、service 層にも同じ state 規則が別実装で残っている。UseCase 側 policy は terminal job 判定と `running|paused|recoverable_failed` の可否を評価する。service helper も terminal job 判定と同じ 3 state の `pause|resume|retry|cancel` 可否を文字列で再実装している。参照: `internal/usecase/translationjobpolicy/policy.go:22-27`, `internal/usecase/translationjobpolicy/policy.go:67-123`, `internal/usecase/phase_policy_helpers.go:17-39`, `internal/service/phase_action_enablement_helpers.go:5-77`
- この重複は今回の plan が目指した「phase service は開始前提、read model 集約、provider 実行だけを残す」と完全には一致していない。plan は `pause`、`resume`、`retry`、`cancel` の可否を `TranslationJobPolicy` の共通操作規則から導出すると書いているが、実装結果では service helper が `TranslationJobPolicy` package を直接 import していないと明記されている。参照: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:67-69`, `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/backend-implementation-result.md:51-53`
- `StateMachine` 旧名は product code と `.go-arch-lint.yml` からは外れているが、active task-local には残っている。observability task は `StateMachine` を状態遷移境界としてまだ参照している。参照: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/backend-implementation-result.md:43`, `docs/exec-plans/active/observability-log-addition/scenario-design.md:110`, `docs/exec-plans/active/observability-log-addition/scenario-candidates.failure.md:95`, `docs/exec-plans/active/observability-log-addition/scenario-candidates.operation-audit.md:19-25`, `docs/exec-plans/active/observability-log-addition/design-diff.component.puml:45`, `docs/exec-plans/active/observability-log-addition/design-diff.sequence.puml:21`, `docs/exec-plans/active/observability-log-addition/design-diff.sequence.puml:40`
- `JobIOService` は architecture 正本、backend architecture 図、arch-lint component、active observability task に残っている一方で、実体 package は `internal/jobio/doc.go` だけである。plan もこの衝突を未決事項として残している。参照: `docs/architecture.md:31`, `docs/architecture.md:53`, `docs/architecture.md:160-186`, `docs/architecture.md:235-256`, `docs/diagrams/backend/backend-architecture.puml:62-69`, `docs/diagrams/backend/backend-architecture.puml:148-149`, `.go-arch-lint.yml:15-16`, `.go-arch-lint.yml:52-59`, `internal/jobio/doc.go:1-2`, `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:36-38`, `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:111-112`
- `stale_*` の意味は現 plan 上では分離されている。今回の `stale` は古い設計名、未使用 package、重複 wrapper、古い task-local 参照を指すと明記されている。一方、`stale_selection`、`validation_stale`、`model_selection_stale` は利用者向けまたは API 向け理由分類として別扱いになっている。observability task でも runtime event stale と設定検証 stale を別候補として扱っている。参照: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:15-16`, `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:46`, `internal/usecase/translation_job_management_contract.go:3-19`, `docs/exec-plans/active/observability-log-addition/scenario-candidates.state-transition.md:128-138`
- ただし、同じ `stale` という語が task 名、runtime event、validation error、artifact stale reason で併存している。detail-spec と active task-local の中では用途別の修飾語が付くが、state 廃止 task の文脈だけで `stale` と書く箇所も残る。参照: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/plan.md:15-16`, `docs/detail-specs/translation-job-management.md:83`, `docs/detail-specs/translation-job-setup.md:41`, `docs/detail-specs/translation-output-artifact.md:35`, `docs/detail-specs/translation-output-artifact.md:67-69`
- 旧 state spelling の残りがある。`PersonaGenerationPhaseContractStub` の cancel fixture は `canceled` ではなく `cancelled` を返す。正本仕様と service 実装は `canceled` で統一されている。参照: `internal/usecase/persona_generation_phase_contract.go:353-360`, `docs/spec.md:161-165`, `internal/service/persona_generation_phase_service.go:29`

## 仮説

- `pending` は phase run 作成直後の一時永続 state として残っており、2026-05-10 の `Ready job には JOB_PHASE_RUN を事前作成しない` 修復とは別の経路で温存されている可能性がある。
- service 層の共通操作規則重複は、`TranslationJobPolicy` を UseCase 専用に固定した architecture 制約を避けるために、read model 用判定を service 内へ複製した可能性がある。
- detail-spec の `phase state` 表現には、永続 state、read model state、UI 表示 state が混在している可能性がある。特に本文翻訳 detail-spec と NPC ペルソナ生成 service はその混在が強い。

## 影響ファイル候補

- `docs/spec.md`
- `docs/architecture.md`
- `docs/diagrams/backend/backend-architecture.puml`
- `.go-arch-lint.yml`
- `internal/jobio/doc.go`
- `internal/service/phase_action_enablement_helpers.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/usecase/persona_generation_phase_contract.go`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/exec-plans/active/observability-log-addition/scenario-design.md`
- `docs/exec-plans/active/observability-log-addition/scenario-candidates.operation-audit.md`
- `docs/exec-plans/active/observability-log-addition/scenario-candidates.failure.md`
- `docs/exec-plans/active/observability-log-addition/scenario-candidates.state-transition.md`

## 残り不足

- `pending` が DB 永続 state として実際にどの経路で作られるかは、repository write path の追加追跡が未実施である。
- `JobIOService` を architecture 正本から外すのか、実装するのかの人間判断が未確定である。
- detail-spec の `phase state` が永続 state なのか、read model state なのか、UI 差分名なのかを task-local だけでは確定できない。

## 残留リスク

- `pending` を含む phase state が仕様外のまま残ると、次の stale 廃止で再び state 語彙の棚卸しが必要になる。
- `JobIOService` を architecture 正本と arch-lint に残したままにすると、未使用 package と設計主語の不一致が継続する。
- observability task が `StateMachine` 前提のまま進むと、次の設計や実装で古い責務境界を再注入する可能性がある。
- `cancelled` と `canceled` の spelling 差異は、fixture 由来でも state 知識の検索漏れを起こす可能性がある。

## 推奨 next step

- `designer` は、今回の stale 廃止を close する前に、少なくとも次の 3 点の判断を分けて持つのが妥当である。
- `pending` を canonical state へ昇格するのか、内部一時 state として隔離するのかを明示する。
- `JobIOService` を architecture 正本から外すのか、別 task で実体化するのかを決める。
- `observability-log-addition` の `StateMachine` / `JobIOService` 旧名参照を、active task-local 更新対象へ含めるかを決める。
