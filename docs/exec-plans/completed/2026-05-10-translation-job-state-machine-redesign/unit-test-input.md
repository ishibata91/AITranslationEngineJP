# 単体テスト引き継ぎ入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `implementation_unit_tester`
- `skill`: `tests-unit`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 依存完了情報

- backend 実装は `backend-implementation-result.md` の範囲で完了済み。
- `backend-local` は旧仕様の test 期待値で失敗している。
- 単体テストの目的は、実装済み backend 責務を新しい共通操作規則で証明することである。

## 対象テスト範囲

- `internal/usecase/translationjobpolicy/policy_test.go`
- `internal/usecase/phase_policy_helpers_test.go`
- `internal/service/term_translation_phase_service_test.go`
- `internal/service/persona_generation_phase_service_test.go`
- `internal/service/body_translation_phase_service_test.go`

`internal/apitest/` は触らない。

## 実装済み対象

- `internal/usecase/translationjobpolicy/policy.go`
- `internal/usecase/phase_policy_helpers.go`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`

## 証明対象

- terminal job では、`start`、`pause`、`resume`、`retry`、`cancel` を拒否する。
- active phase run がある場合、`start` を拒否する。
- `start` は phase 別開始前提が満たされた場合だけ許可する。
- `pause` は `running` の phase run だけ許可する。
- `resume` は `paused` の phase run だけ許可する。
- `retry` は `recoverable_failed` の phase run だけ許可する。
- `cancel` は `paused` の phase run だけ許可する。
- `retry`、`resume`、`pause`、`cancel` は phase type で条件を変えない。
- `RecoverableFailed` の resume は拒否し、retry は許可する。
- persona の `latestError` は retry 可否の条件にしない。

## 既知の旧期待値

- `TestTermTranslationPhaseServiceReadSummaryAllowsResumeForRecoverableFailedRun` は、`RecoverableFailed` を resume 可能と期待している。
- `TestPersonaGenerationPhaseService_CommandMutationsAndReadinessBranches` は、pause 後の retry / cancel の扱いが旧仕様に寄っている。
- `TestPersonaGenerationPhaseService_RetryRejectsNonRetryableStatesAndErrorsWithoutMutation` は、`latestError` による retry 拒否を期待している。

## 禁止範囲

- プロダクトコードは変更しない。
- `internal/apitest/` は変更しない。
- docs 正本、`.codex`、作業計画文書は変更しない。
- UI、Wails DTO、DB schema、migration は変更しない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`

## 返却内容

- 変更ファイル。
- 証明した公開振る舞い、分岐、エラー経路。
- 実行した検証コマンドと結果。
- coverage の結果または未実行理由。
- 未証明範囲と残留リスク。
