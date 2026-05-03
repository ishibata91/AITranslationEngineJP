# Implementation Handoff: normal-flow continuation

- `skill`: implement-integration
- `status`: ready
- `source_plan`: `./plan.md`
- `source_exploration_plan`: `./exploration-test-plan.md`
- `source_findings`: `./exploration-test-findings.md`
- `return_to`: exploration_test_lane

## Goal

`Job Setup` で secret なしの `LMStudio` 通常フローを validation と ready job 作成へ進められるようにする。
`LMStudio` は provider client 側で API key なしを許可しているため、Job Setup の credential 参照状態と validation が同じ扱いになることを確認する。

## Approved Scope

- `internal/service/translation_job_setup_service.go`
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`
- `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts`
- 既存 DTO / Wails binding の互換性確認

## Required Behavior

- `lm_studio` の credential 参照は、secret 本体が空でも `configured` として扱える。
- `lm_studio` の validation は secret 未保存だけを理由に `credential_missing` へしない。
- `openai`、`gemini`、`xai` は secret 未保存を引き続き `credential_missing` とする。
- UI / DTO / error summary / log に secret 本体を出さない。
- `provider / model / execution mode` と `credential reference` の選択状態は既存の表示契約を維持する。

## Evidence

- 停止証跡: `tmp/agent-browser/20260503-section3-job-setup-blocked.png`
- 停止表示: `credential 参照を選択してください。`
- 後続停止: `tmp/agent-browser/20260503-section4-job-run-blocked.png`
- 後続停止: `tmp/agent-browser/20260503-section5-output-management.png`

## Validation

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

## Forbidden

- プロダクトテスト、docs 正本本文、`.codex/` は変更しない。
- fake provider を user-facing provider list に追加しない。
- 実 API key、secret 平文、復号可能値を UI / DTO / log へ出さない。
