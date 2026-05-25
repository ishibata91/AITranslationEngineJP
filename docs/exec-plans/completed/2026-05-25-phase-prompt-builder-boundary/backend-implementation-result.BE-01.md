# backend 実装結果: BE-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-01`
- `implementation_skill`: `implement-backend`
- `status`: implemented-with-validation-blocker

## 実装結果

- `internal/service/prompt_envelope.go` を追加し、`PromptEnvelope`、`PromptDigest`、`RequestShapeID`、`PromptSafeSummary` を分離した。
- 3 フェーズの要求形状識別子を `TERM_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1`、`BODY_TRANSLATION_REQUEST_V1` として共有定数へ固定した。
- 単語翻訳、NPC ペルソナ生成、本文翻訳の prompt builder に `PromptEnvelope` 生成関数を追加した。
- provider adapter は raw prompt を provider client 呼び出しにだけ渡し、audit summary には digest と要求形状識別子だけを渡す形へ寄せた。
- NPC ペルソナ生成の debug log redaction は raw prompt と request body 全文を返さず、復元不能 digest だけを返す helper へ寄せた。
- `internal/service/prompt_envelope_test.go` を追加し、raw prompt、digest、要求形状識別子、利用者向け安全要約の分離を検証した。

## 変更ファイル

- `internal/service/prompt_envelope.go`
- `internal/service/prompt_envelope_test.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`

## 完了条件確認

- 共有受け渡し単位は `RawPrompt`、`Digest`、`RequestShapeID`、`Summary` を別 field として持つ。
- 共有 helper は `PromptSafeSummary` と `RedactedPromptDiagnostic` を返し、secret 本体、raw prompt、request body 全文、response body 全文を公開要約へ出さない。
- `PromptDigest` は raw prompt から SHA-256 digest を作るだけで、生成規則の版選択値として使っていない。

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed

backend-local の内訳:

- backend lint: pass
- `internal/...` package tests: pass
- root package `aitranslationenginejp`: fail

失敗内容:

```text
main.go:18:12: pattern all:frontend/dist: no matching files found
```

## 残留リスク

- backend-local は root package の frontend build artifact 欠落で失敗した。
- `frontend/dist` 作成は frontend build artifact の扱いであり、BE-01 の変更許可範囲外である。
- 検証実行により `scripts/harness/__pycache__/harness_common.cpython-314.pyc` が更新された可能性がある。
