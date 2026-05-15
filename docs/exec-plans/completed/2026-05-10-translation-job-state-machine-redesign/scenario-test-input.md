# シナリオテスト引き継ぎ入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `implementation_scenario_tester`
- `skill`: `tests-scenario`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 依存完了情報

- backend 実装は `backend-implementation-result.md` の範囲で完了済み。
- 承認済みシナリオは `scenario-design.md` にある。
- 今回の対象は API テストである。

## 対象テスト範囲

- `internal/apitest/body_translation_recovery_terminal_readiness_test.go`
- 必要最小限の `internal/apitest/*` 補助。

`internal/service/` と `internal/usecase/` の単体テストは触らない。

## 証明対象

- `SCN-TJSM-003`: pause、resume、retry、cancel は phase type に依存せず、`JOB_PHASE_RUN.state` と job terminal 判定で許可または拒否する。
- `SCN-TJSM-003`: `retry` は `recoverable_failed` の同じ phase run を継続する。
- `SCN-TJSM-003`: `resume` は `paused` の同じ phase run を継続する。
- `SCN-TJSM-003`: `RecoverableFailed` の resume は拒否する。
- `SCN-TJSM-003`: `Running` の cancel は拒否し、状態と出力を変更しない。
- `SCN-TJSM-003`: active phase run がある時の `start` 再送は新規 phase run を作らない。
- `SCN-TJSM-008`: terminal job では状態変更操作を拒否する。

## 既知の旧期待値

- `TestSCN_TFN_011_RetryResumeAndStartResendReusePhaseRunWithoutDuplicateResults` は、`RecoverableFailed` の resume と active phase run 中の start 再送を成功側として扱っている可能性がある。
- `TestSCN_BTP_009_RunningCancelIsRejectedBeforeTerminalResultRewrite` は、Running cancel の拒否 error kind を旧仕様に合わせている可能性がある。

## 禁止範囲

- プロダクトコードは変更しない。
- `internal/service/` と `internal/usecase/` の単体テストは変更しない。
- docs 正本、`.codex`、作業計画文書は変更しない。
- UI、Wails DTO、DB schema、migration は変更しない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-local`

## 返却内容

- 変更ファイル。
- 証明したシナリオ ID。
- 公開接点、入力開始点、主要観測点、期待結果。
- 実行した検証コマンドと結果。
- 未証明シナリオと残留リスク。
