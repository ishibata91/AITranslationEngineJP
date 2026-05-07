# backend 実装入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-BE-01`
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `ready_wave`: `wave-2`
- `source_scope`: `implementation-scope.md`

## 目的

Job Setup の永続化境界を変更する。
Job 側 DB は endpoint、endpoint summary、`credential_ref` 実値、secret store 参照実値、`modelListSourceToken` を所有しない。

## 依存完了

- `frontend 実装`: 完了。
- `frontend test 追従`: 完了。
- `frontend 実装後人間レビュー`: 承認済み。

## 読むファイル

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/diagramming-result.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/implementation-scope.md`
- `internal/infra/sqlite/dbinit/migrations/003_canonical_er_v1_tables.sql`
- `internal/infra/sqlite/dbinit/migrations/009_translation_job_phase_runtime_snapshots.sql`
- `internal/repository/job_lifecycle_repository.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`
- `internal/service/translation_job_setup_service.go`
- `internal/usecase/translation_job_setup_usecase.go`

## 変更許可範囲

- `internal/infra/sqlite/dbinit/migrations/`
- `internal/repository/job_lifecycle_repository.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`
- `internal/service/translation_job_setup_service.go`
- `internal/usecase/translation_job_setup_usecase.go`
- 上記範囲の backend test

## 禁止範囲

- Wails controller / DTO / generated binding の公開契約変更。
- frontend 実装。
- docs 正本本文。
- `.codex/`
- secret 本体を保存または出力する変更。

## secret 境界

保存してよい値:
provider、model、execution mode、batch mode、credential 状態分類、接続確認状態、再解決結果分類、再解決時刻。

保存禁止:
`credential_ref` 実値、endpoint 原文、endpoint summary、secret store 参照実値、`modelListSourceToken`、APIキー本体、raw request、raw response、raw prompt。

## 初手

- path: `internal/repository/job_lifecycle_repository.go`
- 対象: `TranslationJobPhaseRuntimeSnapshot`
- 変更種別: Job 側 snapshot の所有 field 削減

## 完了条件

- `JOB_PHASE_RUN` と `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は Job Setup 作成時に `credential_ref` 実値と endpoint 系値を所有しない。
- Ready job 作成時に保存する phase 値は provider、model、execution mode、batch mode、credential 状態分類だけである。
- `modelListSourceToken` は Job 側 DB、Job Setup summary、利用者向け表示に残さない。
- provider settings revision と更新履歴を Job 側へ保存しない。
- 既存の Job Setup 作成経路は provider settings を fallback として Job にコピーしない。

## 検証コマンド

- `go test ./internal/repository ./internal/infra/sqlite/dbinit ./internal/service ./internal/usecase -run 'JobLifecycle|TranslationJobSetup|ProviderSettings|Migration'`
- `python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-implementation-result.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因

