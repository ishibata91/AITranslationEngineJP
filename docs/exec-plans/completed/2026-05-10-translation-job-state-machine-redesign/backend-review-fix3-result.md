# Backend レビュー修正 3 結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `completed`
- `completed_at`: `2026-05-14`

## 修正対象

`state-invariant-002` を修正した。
単語翻訳 phase の resume / retry は、同じ `JOB_PHASE_RUN` を required state 条件付きで `running` に更新する。

## 変更ファイル

- `internal/service/term_translation_phase_service.go`
- `internal/service/term_translation_phase_service_test.go`

## 修正内容

- `resume` は `UpdateJobPhaseRunWhenState` に expected state `paused` を渡す。
- `retry` は `UpdateJobPhaseRunWhenState` に expected state `recoverable_failed` を渡す。
- `updateExecutionPlanRunState` に expected state 引数を追加した。
- terminal job 拒否、phase run id 一致確認、同じ phase run 継続、runtime snapshot 保存順序は維持した。

## 追加テスト

- resume 成功時に expected state `paused` が渡ることを確認した。
- retry 成功時に expected state `recoverable_failed` が渡ることを確認した。
- expected state 不一致相当の `repository.ErrConflict` では provider request が発生せず、phase run が `running` にならないことを確認した。

## backend_implementer 検証結果

- `gofmt -l internal/service`: 通過。
- `python3 scripts/harness/run.py --suite backend-local`: 通過。
- `python3 scripts/harness/run.py --suite coverage`: 通過。

## 注意

backend_implementer は追加テストの cognitive complexity による Sonar maintainability high を一度検出した。
helper 分割後、coverage は通過した。

## 残事項

親側で backend-local と coverage を再実行し、5 観点レビューを再実行する。
