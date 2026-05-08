# 単体テスト入力

## 対象成果物

`implementation_unit_tester` は `tests-unit` として回帰テストを追加または更新する。

## 根拠

- [backend-implementation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/backend-implementation-input.md:1)
- [backend-implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/backend-implementation-result.md:1)

## 実装済み対象

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:794)
- [016_remove_ready_job_initial_pending_phase_run.sql](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1)

## 対象テスト範囲

- [translation_job_setup_service_test.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service_test.go:1)
- [migration_test.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migration_test.go:1)

## 証明対象

- Job Setup 完了後に `CreateJobPhaseRun` が呼ばれない。
- Job Setup 完了後も runtime snapshot は 3 phase 分保存される。
- migration は `ready` job の未実行 placeholder だけを削除する。
- migration は `running`、進捗あり、開始時刻あり、`ready` 以外の job の phase run を削除しない。

## 検証コマンド

`python3 scripts/harness/run.py --suite backend-local`
