# Backend 修正実行入力

## 対象成果物

`implementation_implementer` は `implement-backend` として実装する。
実装対象は、作成直後の `ready` job に未実行の `pending` phase run が残らない状態へ戻す backend 修正である。

## 根拠

- [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/human-observation.md:1)
- [investigation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/investigation.md:1)
- [cause-sequence.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-08-job-id-1-ready-lock-fix/cause-sequence.puml:1)

## 変更してよい範囲

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:822)
- [dbinit migrations](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/sqlite/dbinit/migrations/015_translation_job_phase_runtime_snapshot_non_secret_boundary.sql:1)

## 修正方針

- Job Setup は作成直後の `JOB_PHASE_RUN` を作成しない。
- Job Setup は `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` を 3 phase 分保存する既存挙動を維持する。
- 実行開始時の phase run 作成は、`term_translation` の `running` を作る既存 flow に寄せる。
- 既存DBの救済として、`ready` job に紐づく未実行 placeholder だけを migration で削除する。

## placeholder 削除条件

削除対象は次の全条件を満たす `JOB_PHASE_RUN` だけに限定する。

- `JOB_PHASE_RUN.phase_type = 'translation'`
- `JOB_PHASE_RUN.state = 'pending'`
- `JOB_PHASE_RUN.progress_percent = 0`
- `JOB_PHASE_RUN.started_at IS NULL`
- `JOB_PHASE_RUN.finished_at IS NULL`
- 親 `TRANSLATION_JOB.state = 'ready'`

## 禁止変更範囲

- `docs/` 正本本文は変更しない。
- `.codex/` と `.codex/skills` は変更しない。
- frontend は変更しない。
- 削除 guard の安全判定はこの実装では緩めない。
- 新しい phase 状態や新しい表示状態は追加しない。

## 回帰確認観点

- 新規 Job Setup 完了後に `CreateJobPhaseRun` が呼ばれない。
- 新規 Job Setup 完了後も runtime snapshot は 3 phase 分保存される。
- migration は未実行 placeholder だけを削除する。
- migration は `running`、進捗あり、開始時刻あり、`ready` 以外の job の phase run を削除しない。

## 検証コマンド

`python3 scripts/harness/run.py --suite backend-local`
