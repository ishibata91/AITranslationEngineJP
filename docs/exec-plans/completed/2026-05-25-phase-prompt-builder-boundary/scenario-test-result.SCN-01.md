# scenario test 結果: SCN-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `SCN-01`
- `implementation_skill`: `tests-scenario`
- `status`: `completed-with-known-backend-local-blocker`

## 証明済みシナリオ結果

1. backend 公開接点の開始結果は、生成指示全文ではなく digest と件数を返す。
- 証明テスト: `TestSCN_BTP_008_StartPublicResultExposesOnlySafePromptSummary`
- 公開接点: `BodyTranslationPhaseUsecase.StartBodyTranslationPhase`
- 入力開始点: `StartBodyTranslationPhaseRequest{JobID: 1005}`
- 主要観測点: `InputSummary.PromptDigest`, `RequestSummary.ProviderTargetCount`, `Execution.RequestUnitCount`
- 保護確認: JSON 化した結果に fake secret、raw prompt、provider raw response、原文発話全文、会話文脈全文が含まれない。

2. invalid provider response は本文翻訳フェーズの失敗分類として返る。
- 証明テスト: `TestSCN_BTP_008_StartInvalidProviderResponseReturnsRedactedFailureKind`
- 公開接点: `BodyTranslationPhaseUsecase.StartBodyTranslationPhase`
- 入力開始点: provider stub 相当の fake service が invalid provider response を返す開始要求。
- 主要観測点: `ErrorSummary.ErrorKind == invalid_provider_response`, `ErrorSummary.IsRedacted == true`, `ErrorSummary.Retryable == true`
- 保護確認: JSON 化した結果に raw provider response が含まれない。

3. Wails 公開接点は invalid provider response の raw error を公開しない。
- 証明テスト: `TestBodyTranslationPhaseControllerStartReturnsOnlySafeSummaryForInvalidProviderResponse`
- 公開接点: `BodyTranslationPhaseController.StartBodyTranslationPhase`
- 入力開始点: `StartBodyTranslationPhaseRequestDTO{JobID: 778}`
- 主要観測点: `error == nil`, `errorSummary.errorKind == invalid_provider_response`, `inputSummary.promptDigest`, `requestSummary.providerTargetCount`, `execution.requestUnitCount`
- 保護確認: usecase stub の error 文字列に raw provider response 相当文字列を含めても、Wails response JSON には含まれない。

4. 実 AI API 呼び出しは発生していない。
- 証明方法: usecase scenario test は `fakeBodyTranslationScenarioService` を使う。
- 証明方法: Wails controller unit test は `fakeBodyTranslationPhaseUsecase` を使う。

## 変更ファイル

- `internal/usecase/body_translation_phase_scenario_test.go`
- `internal/controller/wails/body_translation_phase_controller_unit_test.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/scenario-test-result.SCN-01.md`

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/usecase ./internal/controller/wails`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed
  - backend lint: pass
  - `internal/...` package tests: pass
  - root package `aitranslationenginejp`: fail

## 未証明小範囲

- 単語翻訳フェーズと NPC ペルソナ生成フェーズの Wails 公開接点は、SCN-01 の初手と証明範囲では補強していない。
- 本文翻訳フェーズの実 repository 接続を伴う provider 実行は、実 AI API 呼び出し禁止と provider stub 条件により対象外である。

## 残った失敗と原因

1. `backend-local` の root package test が失敗している。
- 原因: `main.go:18:12: pattern all:frontend/dist: no matching files found`
- 判断: `scenario-test-input.SCN-01.md` に記録済みの既知 blocker である。
- 影響: backend lint と `internal/...` package tests は通過しているため、SCN-01 追加テストの失敗ではない。
