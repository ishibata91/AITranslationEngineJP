# 5 観点レビュー起動入力

## 対象

ジョブID 1 が作成直後に実行も削除もできなくなる不具合の恒久修正をレビューする。

## 入力成果物

- [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/human-observation.md:1)
- [investigation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/investigation.md:1)
- [cause-sequence.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/cause-sequence.puml:1)
- [backend-implementation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/backend-implementation-input.md:1)
- [backend-implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/backend-implementation-result.md:1)
- [regression-test-evidence.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/regression-test-evidence.md:1)
- [browser-confirmation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/browser-confirmation-result.md:1)

## 変更ファイル

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:794)
- [016_remove_ready_job_initial_pending_phase_run.sql](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1)
- [translation_job_setup_service_test.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service_test.go:445)
- [migration_test.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migration_test.go:780)

## 検証

- `python3 scripts/harness/run.py --suite backend-local`: PASS。
- `python3 scripts/harness/run.py --suite coverage`: PASS。
- ブラウザ確認: job ID 1 の状態不整合削除拒否は表示されない。

## reviewback 出力先

- `docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/reviewback.behavior.yaml`
- `docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/reviewback.contract.yaml`
- `docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/reviewback.trust-boundary.yaml`
- `docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/reviewback.state-invariant.yaml`
- `docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/reviewback.responsibility-boundary.yaml`
