# unit test 結果: TU-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `TU-01`
- `implementation_skill`: `tests-unit`
- `status`: `completed-with-external-blockers`

## 証明済み完了条件

1. 単語翻訳の invalid response 分岐を証明した。
- 対象語不一致: `TestTermTranslationProviderAdapterReturnsInvalidResponseWhenSourceTermMismatches`
- 空訳語: `TestTermTranslationProviderAdapterReturnsInvalidResponseWhenTranslatedTermIsEmpty`
- 余分な応答: `TestTermTranslationProviderAdapterReturnsInvalidResponseWhenProviderReturnsExtraItem`
- 欠落応答: `TestTermTranslationProviderAdapterReturnsInvalidResponseWhenProviderReturnsMissingItem`

2. NPC ペルソナ生成の invalid response と境界を証明した。
- 対応識別子不一致: `TestPersonaGenerationProviderAdapterRejectsMismatchedCorrelation`
- 空ペルソナ: `TestPersonaGenerationProviderAdapterRejectsEmptyPersonaBody`
- debug log redaction: `TestPersonaGenerationProviderAdapterMapsValidResponse`
- credential 非公開: `TestPersonaGenerationProviderAdapterRedactsProviderCredentialFailure`

3. 本文翻訳の invalid response と境界を証明した。
- 翻訳項目識別子不一致: `TestBodyTranslationProviderAdapterRejectsMismatchedFieldCorrelationKey`
- 空訳文: `TestBodyTranslationProviderAdapterRejectsEmptyTranslatedText`
- 保持要素不整合: `TestValidateBodyTranslationProviderResultsReturnsProtectionValidationFailure`
- request summary の raw prompt 非公開: `TestBodyTranslationProviderAdapterAuditSummaryDoesNotExposeRawPromptInputs`

4. `PromptDigest` と `REQUEST_V1` 系の意味境界を証明した。
- `PromptDigest` は内部同一性情報: `TestTermTranslationProviderAdapterReturnsRequestShapeIdentifierAndDigestAsInternalIdentity`
- `REQUEST_V1` は要求形状識別子: `TestPromptEnvelopeSeparatesRawPromptDigestShapeAndSafeSummary` と `TestPhasePromptEnvelopesUseRequestShapeIdentifiers`

## 変更ファイル

- `internal/service/term_translation_provider_adapter_test.go` (新規)
- `internal/service/persona_generation_provider_adapter_test.go`
- `internal/service/body_translation_provider_adapter_test.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/unit-test-result.TU-01.md` (新規)

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed
  - backend lint: pass
  - `internal/...` package tests: pass
  - root package `aitranslationenginejp`: fail (`main.go:18:12: pattern all:frontend/dist: no matching files found`)

## coverage 結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite coverage`: failed
- 失敗理由:
  - frontend coverage で `wailsjs/go/wails/AppController.js` 解決失敗と gateway import parse error が発生した。
  - Sonar coverage summary が `62.5% < 70.0%` で gate 未達になった。
  - Sonar maintainability HIGH issue 1 件 (`internal/service/prompt_envelope.go` の重複 literal) が検出された。
- 備考: 上記は TU-01 の許可範囲外ファイルまたは既存実装由来であり、今回の単体テスト変更だけでは解消できない。

## 未証明小範囲

- なし（TU-01 の完了条件 1〜4 に対応する公開振る舞い、分岐、エラー経路は対象テストで証明済み）。

## 残った失敗と原因

1. `backend-local` 失敗
- 原因: root package テストに必要な `frontend/dist` が欠落している。
- 範囲: frontend build artifact 準備であり、TU-01 の変更許可範囲外。

2. `coverage` 失敗
- 原因: frontend 側の import 解決と Sonar gate 未達、HIGH issue 残存。
- 範囲: TU-01 の変更許可範囲外。
