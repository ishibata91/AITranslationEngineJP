# ジョブセットアップ未開始 phase run 修正 回帰テスト証跡

## 判断結果

- 単体テストによる回帰確認は完了した。
- 担当 agent は `implementation_unit_tester` である。
- 使用 skill は `tests-unit` である。

## 変更ファイル

- `internal/service/translation_job_setup_service_test.go`
- `internal/service/term_translation_phase_service_test.go`
- `internal/service/persona_generation_phase_service_test.go`
- `internal/service/body_translation_phase_service_test.go`
- `internal/repository/job_lifecycle_sqlite_repository_test.go`

## 証明済み完了条件

- setup 完了時に未開始 `JOB_PHASE_RUN` が 4 件作成される。
- 未開始 row は `pending` ではない。
- `translation`、`term_translation`、`body_translation` は `idle_ready` で作成される。
- `persona_generation` は `not_started` で作成される。
- 単語翻訳 start は、先置き済み row を `running` へ更新する。
- ペルソナ生成 start は、先置き済み row を `running` へ更新する。
- 本文翻訳 start は、先置き済み row を `running` へ更新する。
- phase run が無い ready job の互換 fallback は維持される。
- SQLite repository は `pending` を delete unsafe と扱う。
- SQLite repository は `idle_ready` を delete safe と扱う。

## 検証結果

- 実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- 結果: 成功
- 実行コマンド: `python3 scripts/harness/run.py --suite coverage`
- 結果: 成功
- coverage summary: `70.8%`
- security、reliability、maintainability は基準内で成功した。

## 未証明小範囲

- `translation_job_management_service_test.go` への直接追記はしていない。
- 未開始 row の投影挙動は `backend-local` 全体通過で回帰なしを確認した。

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/implementation-evidence.md`
- `internal/service/translation_job_setup_service_test.go`
- `internal/service/term_translation_phase_service_test.go`
- `internal/service/persona_generation_phase_service_test.go`
- `internal/service/body_translation_phase_service_test.go`
- `internal/repository/job_lifecycle_sqlite_repository_test.go`
