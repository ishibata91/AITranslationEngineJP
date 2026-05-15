# シナリオテスト結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `agent`: `implementation_scenario_tester`
- `status`: `passed`
- `completed_at`: `2026-05-14`

## 変更ファイル

- `internal/apitest/body_translation_recovery_terminal_readiness_test.go`: `SCN-TJSM-003` と `SCN-TJSM-008` の API テスト期待値を更新した。

## 証明済み

- `SCN-TJSM-003`: `recoverable_failed` の `retry` は同じ `JOB_PHASE_RUN` を継続し、既存 field result を重複作成しない。
- `SCN-TJSM-003`: `paused` の `resume` は同じ `JOB_PHASE_RUN` を `running` として継続する。
- `SCN-TJSM-003`: `recoverable_failed` の `resume` は拒否され、行を変更しない。
- `SCN-TJSM-003`: `running` の `cancel` は拒否され、状態と出力を変更しない。
- `SCN-TJSM-003`: active phase run 中の `start` 再送は `active_phase_exists` で拒否され、新規 phase run を作らない。
- `SCN-TJSM-008`: terminal job の `start` と `cancel` は `terminal_job` で拒否され、状態と行を変更しない。

## 公開接点

- `BodyTranslationPhaseController.StartBodyTranslationPhase`
- `BodyTranslationPhaseController.ResumeBodyTranslationPhase`
- `BodyTranslationPhaseController.RetryBodyTranslationPhase`
- `BodyTranslationPhaseController.CancelBodyTranslationPhase`

## 検証結果

- `go test ./internal/apitest -run 'TestSCN_TJSM_003|TestSCN_TJSM_008|TestSCN_BTP_009_RunningCancel'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。

## 未証明範囲

`SCN-TJSM-008` の terminal job による保存拒否、readiness 更新拒否、late response 後書き拒否は、既存 late response test があるが、terminal job 条件の専用追加証明ではない。
