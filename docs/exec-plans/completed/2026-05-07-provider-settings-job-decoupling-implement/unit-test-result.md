# unit test result

- task_id: `2026-05-07-provider-settings-job-decoupling-implement`
- handoff_id: `PSJD-UT-01`
- implementation_skill: `tests-unit`
- status: `completed`

## 実装判断

既存の unit test 群が、入力で指定された 4 観点をすでに証明していた。
product code 不足は検出されなかったため、product code 変更は実施していない。

## 証明観点と根拠

- Job Setup の選択値永続化
  - `internal/service/translation_job_setup_service_test.go`
  - `TestTranslationJobSetupServiceReadSummaryReturnsPersistedPhaseRuntimeSnapshots`
- phase 開始時の provider settings 再解決（3 phase）
  - `internal/service/term_translation_phase_service_test.go`
  - `internal/service/persona_generation_phase_service_test.go`
  - `internal/service/body_translation_phase_service_test.go`
- Running 非 secret 要約の保存禁止値
  - `internal/service/term_translation_phase_service_test.go`
  - `internal/service/persona_generation_phase_service_test.go`
  - `internal/service/body_translation_phase_service_test.go`
- frontend presenter/usecase/view の credential 参照値非表示
  - `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts`
  - `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`

## 変更ファイル

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/unit-test-result.md`

## 検証結果

- `go test ./internal/...`: pass
- `npm --prefix frontend run test -- src/application src/ui/screens/translation-job-setup src/controller`: pass
- `python3 scripts/harness/run.py --suite backend-local`: pass
- `python3 scripts/harness/run.py --suite frontend-local`: pass

## 残った失敗と原因

残った失敗はない。
