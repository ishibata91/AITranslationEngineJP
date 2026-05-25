# backend 実装結果: BE-04

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-04`
- `implementation_skill`: `implement-backend`
- `status`: implemented-with-validation-blocker

## 実装結果

- `internal/service/body_translation_prompt_builder.go` は、本文翻訳 prompt の要求形状識別子を共有 `BodyTranslationRequestShapeV1` に揃えている。
- `internal/service/body_translation_prompt_builder.go` は、provider 接続情報を持たない `BodyTranslationPromptInput` を追加した。
- `internal/service/body_translation_prompt_builder.go` は、`BodyTranslationPromptBuilder` と default 実装を追加し、本文翻訳 prompt 生成を専用 input から `PromptEnvelope` を作る境界へ分けた。
- `internal/service/body_translation_prompt_builder.go` は、本文翻訳の `PromptEnvelope` に 1 翻訳項目、翻訳項目識別子、field correlation key、保持要素件数、辞書制約件数を固定している。
- `internal/service/body_translation_prompt_builder.go` は、`BuildBodyTranslationPromptEnvelope` と `BuildBodyTranslationPrompt` を provider request から専用 input へ詰め替える互換 wrapper として残した。
- `internal/service/body_translation_prompt_builder.go` は、provider audit 用 request summary から辞書制約の raw 文字列と保持要素 raw 一覧を除外した。
- `internal/service/body_translation_provider_adapter_test.go` は、専用 input が provider 接続情報を持たないこと、builder が 1 翻訳項目単位の envelope と安全要約を作ること、audit summary には識別子と件数だけが残ることを検証している。

## 変更ファイル

- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter_test.go`

## 完了条件確認

- 本文翻訳の生成指示は、1 翻訳項目を 1 実行単位として `BuildBodyTranslationPromptEnvelope` へ固定している。
- 本文翻訳の生成指示入力は `BodyTranslationPromptInput` に分離し、provider、model、credential、endpoint を含めていない。
- `BodyTranslationPromptBuilder` は `Build(input BodyTranslationPromptInput) (PromptEnvelope, error)` を持つ。
- provider adapter は `PromptEnvelope.RawPrompt` を provider client 呼び出しへ渡すだけで、保存判断や採用判断を持たない。
- 応答欠落、余分な応答、翻訳項目識別子不一致、field correlation key 不一致、空訳文は `mapBodyTranslationProviderResponse` で invalid provider response として返る。
- 保持要素不整合は `ValidateBodyTranslationProtection` 経由で翻訳項目単位の `protection_validation_failed` として返る。
- request summary は、request unit id、field correlation key、record type、field type、protection source digest、保持要素件数、辞書制約件数に限定した。

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed

backend-local の失敗内容:

```text
FAIL	aitranslationenginejp [setup failed]
# aitranslationenginejp
main.go:18:12: pattern all:frontend/dist: no matching files found
```

## 残った失敗と原因

- backend-local は root package の `frontend/dist` 欠落で失敗した。
- `frontend/dist` 作成は frontend build artifact の扱いであり、BE-04 の変更許可範囲外である。
- backend-local の backend lint と `internal/service` test は通過している。

## 根拠参照

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-input.BE-04.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/detail-specs/body-translation-phase.md`
- `internal/service/prompt_envelope.go`
