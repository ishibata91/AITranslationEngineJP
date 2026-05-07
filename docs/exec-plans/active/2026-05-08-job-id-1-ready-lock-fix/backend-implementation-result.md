# Backend 実装証跡

## 判断結果

backend 実装は完了した。
Job Setup 完了時の初期 `pending` phase run 作成を削除し、既存 DB 救済 migration を追加した。

## 変更ファイル

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:794)
- [016_remove_ready_job_initial_pending_phase_run.sql](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1)

## 実装内容

- `createTranslationJobInTransaction` から `createInitialTranslationPhaseRun` 呼び出しを削除した。
- `translationJobSetupJobLifecycleRepository` から Job Setup 用の `CreateJobPhaseRun` 要求を削除した。
- 未使用になった初期 phase run 作成関数と `pending` 用定数を削除した。
- `ready` job に紐づく未実行 placeholder だけを削除する migration を追加した。

## 検証結果

`python3 scripts/harness/run.py --suite backend-local` は通過した。

## 残留リスク

単体テストで、Job Setup が phase run を作らないことと migration の削除条件を明示的に証明する必要がある。
