# backend 実装結果: BE-02

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-02`
- `implementation_skill`: `implement-backend`
- `status`: `implemented-with-validation-blocker`

## 実装結果

- `internal/service/term_translation_prompt_builder.go` を追加し、単語翻訳専用の `TermTranslationPromptInput` と `TermTranslationPromptBuilder` を分離した。
- 単語翻訳の生成指示は、1 対象語を 1 実行単位として、対象語、原文言語、訳文言語、応答対応識別子を同じ `PromptEnvelope` に固定した。
- `internal/service/term_translation_provider_adapter.go` は、単語翻訳専用 builder から受け取った raw prompt を provider client へ渡す処理に寄せた。
- provider adapter の audit summary は、`RequestShapeID` と `PromptDigest` を使い、raw prompt と secret を含まない形を維持した。
- `internal/service/term_translation_phase_service.go` は、対象語単位の応答対応識別子を provider request へ渡すようにした。
- 応答欠落、余分な応答、空訳語、対象語不一致は、既存の provider adapter と phase service の検査経路で `invalid_provider_response` として扱う。
- 有効応答だけが `ensureJobDictionaryEntry` と `linkPhaseRunDictionaryEntry` を通り、対象ジョブ内辞書へ採用される。

## 変更ファイル

- `internal/service/term_translation_prompt_builder.go`
- `internal/service/term_translation_prompt_builder_test.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/term_translation_phase_service_test.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-02.md`

## 完了条件確認

- completion_signal 1: 単語翻訳 builder が `request_unit_id`、`source_language`、`target_language`、`source_term` を 1 対象語の入力として固定した。
- completion_signal 2: provider adapter から prompt 文言組み立てを外し、AIサービス接続差異の吸収と provider response の形状検査に閉じた。
- completion_signal 3: 応答欠落、余分な応答、空訳語、対象語不一致は、対象語単位の `invalid_provider_response` へ分類される。
- completion_signal 4: 有効応答だけが対象ジョブ内辞書へ採用される既存境界を維持した。

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

## 残った失敗と原因

- backend-local は root package の `frontend/dist` 欠落で失敗した。
- `frontend/dist` 作成は frontend build artifact の扱いであり、BE-02 の変更許可範囲外である。
- `scripts/harness/__pycache__/harness_common.cpython-314.pyc` は検証実行で更新された可能性がある。
