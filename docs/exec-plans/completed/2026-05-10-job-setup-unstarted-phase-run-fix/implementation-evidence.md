# ジョブセットアップ未開始 phase run 修正 実装証跡

## 判断結果

- backend 実装は完了した。
- 担当 agent は `backend_implementer` である。
- 使用 skill は `implement-backend` である。
- 初回 `backend-local` は旧期待値 test 1 件で失敗した。
- 回帰テスト更新後の `backend-local` は成功した。

## 変更ファイル

- `internal/service/translation_job_setup_service.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/translation_job_management_service.go`
- `internal/repository/job_lifecycle_repository.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`

## 変更した symbol

- `TranslationJobSetupService.createUnstartedPhaseRuns`
- `translationJobSetupUnstartedPhaseRunDraft`
- `TermTranslationPhaseService.ensureExecutionPlanRun`
- `termTranslationExecutionBasePhase`
- `PersonaGenerationPhaseService.StartPhase`
- `PersonaGenerationPhaseService.startPhaseRunTransaction`
- `BodyTranslationPhaseService.StartPhase`
- `BodyTranslationPhaseService.createBodyTranslationRunRecord`
- `findTranslationJobManagementCurrentRun`
- `translationJobManagementPhaseRunCompleted`
- `repository.JobPhaseRunUpdateDraft`
- `SQLiteJobLifecycleRepository.UpdateJobPhaseRun`

## 変更内容

- setup 完了時に未開始 `JOB_PHASE_RUN` を 4 件作成する。
- `translation` と `term_translation` と `body_translation` は `idle_ready` として作る。
- `persona_generation` は `not_started` として作る。
- phase start は先置き済み row を `running` へ更新する。
- 未開始 row だけの job を job management が後続 phase 進行中と誤判定しないようにする。
- `UpdateJobPhaseRun` で実行設定と snapshot digest を更新できるようにする。

## 検証結果

- 実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- 初回結果: 失敗
- lint: 通過
- backend test: `TestTJSPPS007TranslationJobSetupServiceCreateCapturesOnlyTargetPhaseRuntimeSettings` が失敗
- 失敗理由: 旧 test は setup 時に `JOB_PHASE_RUN` を作らない期待を持つため、今回方針と衝突した。
- 回帰テスト更新後の実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- 回帰テスト更新後の結果: 成功
- coverage 実行コマンド: `python3 scripts/harness/run.py --suite coverage`
- coverage 結果: 成功
- coverage summary: `70.8%`

## 次判断材料

- レビューでは、未開始 row の state、start 遷移、delete guard 境界、旧 DB fallback 維持を確認する。

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/fix-execution-input.md`
- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/human-observation.md`
- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/investigation.md`
- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/cause-sequence.puml`
