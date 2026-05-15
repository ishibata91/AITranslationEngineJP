# Backend 実装結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `agent`: `backend_implementer`
- `status`: `implemented`
- `completed_at`: `2026-05-14`

## 実装結果

`translationjobpolicy` を UseCase 専用の pure rule として追加した。
各 phase UseCase は、service 呼び出し前に共通操作規則を評価する。
各 phase Service は、resume / retry / cancel の表示規則と実行 guard を共通操作規則へ寄せた。

## 変更ファイル

- `.go-arch-lint.yml`: `usecase` から `translationjobpolicy` への依存を許可した。
- `internal/usecase/translationjobpolicy/policy.go`: 共通操作規則を追加した。
- `internal/usecase/phase_policy_helpers.go`: phase policy 評価用 helper を追加した。
- `internal/usecase/term_translation_phase_usecase.go`: 単語翻訳 phase 操作前に policy を評価する。
- `internal/usecase/persona_generation_phase_usecase.go`: NPC ペルソナ生成 phase 操作前に policy を評価する。
- `internal/usecase/body_translation_phase_usecase.go`: 本文翻訳 phase 操作前に policy を評価する。
- `internal/service/term_translation_phase_service.go`: recoverable failed の resume を拒否し、retry だけを許可する。
- `internal/service/persona_generation_phase_service.go`: resume / retry / cancel の条件を共通操作規則へ揃える。
- `internal/service/body_translation_phase_service.go`: resume / retry / cancel の条件を共通操作規則へ揃える。

## 実装済み操作規則

- terminal job は状態変更操作を拒否する。
- `start` は active phase run がある場合に拒否する。
- `start` は service summary 上の開始前提が満たされた場合だけ許可する。
- `pause` は `running` の phase run だけ許可する。
- `resume` は `paused` の phase run だけ許可する。
- `retry` は `recoverable_failed` の phase run だけ許可する。
- `cancel` は `paused` の phase run だけ許可する。

## 検証結果

- `gofmt -l internal/usecase internal/service`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: fail。

## backend-local 失敗

既存 test が旧仕様の期待値を持っている。
失敗した主な test は次である。

- `TestSCN_BTP_009_RunningCancelIsRejectedBeforeTerminalResultRewrite`
- `TestPersonaGenerationPhaseService_CommandMutationsAndReadinessBranches`
- `TestPersonaGenerationPhaseService_RetryRejectsNonRetryableStatesAndErrorsWithoutMutation`
- `TestTermTranslationPhaseServiceReadSummaryAllowsResumeForRecoverableFailedRun`

## 後続判断

単体テストは、`translationjobpolicy` と service 局所期待値を新しい共通操作規則へ更新する。
シナリオテストは、承認済みシナリオの API 受け入れ経路を新しい共通操作規則へ更新する。
