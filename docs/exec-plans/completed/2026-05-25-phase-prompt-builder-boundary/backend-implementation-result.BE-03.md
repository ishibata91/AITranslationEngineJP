# backend 実装結果: BE-03

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-03`
- `implementation_skill`: `implement-backend`
- `status`: `implemented-with-validation-blocker`

## 実装結果

- `internal/service/persona_generation_prompt_builder.go` を追加し、NPC ペルソナ生成の prompt 文言組み立てと `PromptEnvelope` 生成を provider adapter から分離した。
- `PersonaGenerationPromptInput` を追加し、prompt 生成に必要な NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約だけを専用 input に固定した。
- `PersonaGenerationPromptBuilder` interface と `DefaultPersonaGenerationPromptBuilder` を追加し、`Build(input PersonaGenerationPromptInput) (PromptEnvelope, error)` を NPC ペルソナ生成フェーズ内の差し替え境界にした。
- NPC ペルソナ生成 prompt は 1 NPC を 1 実行単位にし、`RequestUnitID`、`NPCCorrelationID`、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ入力へ固定した。
- `internal/service/persona_generation_provider_adapter.go` は `PromptEnvelope` から raw prompt だけを provider client 呼び出しへ渡し、audit summary は digest と要求形状識別子だけを持つ形へ寄せた。
- provider debug log は `RedactedPromptDiagnostic` を使い、prompt と request body の全文を返さず digest 化した診断値だけを返す。
- `internal/service/persona_generation_phase_service.go` の集約 prompt digest 算出を `BuildPersonaGenerationPromptEnvelope` 経由にし、service 側も共有 prompt 境界に揃えた。
- `BuildPersonaGenerationPromptEnvelope` は provider request から専用 input へ詰め替える互換 wrapper として残した。
- `internal/service/persona_generation_provider_adapter_test.go` に、専用 input / builder の 1 NPC prompt envelope、安全要約、provider adapter の digest / redaction 検証を追加した。

## 変更ファイル

- `internal/service/persona_generation_prompt_builder.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/persona_generation_provider_adapter_test.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-03.md`

## 完了条件確認

- NPC ペルソナ生成の生成指示は 1 NPC を 1 実行単位にし、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ prompt builder 入力として扱う。
- prompt 生成責務は `PersonaGenerationPromptInput` と `PersonaGenerationPromptBuilder` に閉じ、provider request は互換 wrapper で専用 input へ変換する。
- provider adapter は prompt 文言組み立てを持たず、provider client 呼び出し、provider response shape の解釈、NPC 単位の invalid response 判定に閉じる。
- 応答欠落、余分な応答、NPC 対応識別子不一致、空のペルソナ本文は、既存の `invalid_provider_response` failure として NPC 単位に返る。
- 有効応答だけが `PersonaBody` として service へ返り、service 側の `persistGeneratedPersonaTarget` 経由で翻訳ジョブ内ペルソナへ採用される。
- debug log と audit summary は digest、要求形状識別子、redacted header、件数、識別子だけを返し、原文発話全文と会話文脈全文を返さない。

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed

追加戻し後の再検証:

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

- `backend-local` は root package の `frontend/dist` 欠落で失敗した。
- 同じ失敗は `backend-implementation-result.BE-01.md` に既知 blocker として記録されている。
- `frontend/dist` の作成は frontend build artifact の扱いであり、BE-03 の変更許可範囲外である。
- 検証実行により `scripts/harness/__pycache__/harness_common.cpython-314.pyc` が更新された可能性がある。

## 注意

- 作業開始時点で BE-02、BE-04、BE-01 と見られる未コミット変更が同じ worktree に存在した。
- BE-03 は他 agent の変更を戻さず、NPC ペルソナ生成フェーズの対象ファイルだけへ変更を重ねた。
