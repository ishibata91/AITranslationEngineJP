# backend 実装入力: BE-03

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-03`
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `ready_wave`: `wave-2`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

NPC ペルソナ生成フェーズの生成指示責務を分離する。
1 NPC を 1 実行単位にし、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ生成指示入力へ固定する。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md` に `implemented-with-validation-blocker` として記録済み。
- `internal/service/prompt_envelope.go`: 共有 `PromptEnvelope`、`PromptDigest`、要求形状識別子、安全要約境界を追加済み。
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/detail-specs/persona-generation-phase.md`
- `internal/service/prompt_envelope.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/persona_generation_phase_service.go`

## 変更許可範囲

- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/persona_generation_phase_service.go`
- 必要な NPC ペルソナ生成専用 prompt builder file。
- 上記範囲の backend test。

## 禁止範囲

- 共有 `PromptEnvelope` 契約の意味変更。必要な場合は実装を止めて報告する。
- 単語翻訳、本文翻訳の実装変更。
- frontend 実装。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文。
- `.codex/`。
- フェーズ追加、phase type 追加、ジョブ状態機械変更。
- 実 AI API を使う検証。
- 他 agent の変更の取り消し。

## secret 境界

表示、DTO、read model に出してよい値:
`CredentialRef`、provider、model、execution mode、request unit id、NPC 対応識別子、input count、output count、failure kind、`PromptDigest`。

出してはいけない値:
secret 本体、raw prompt、原文発話全文、会話文脈全文、provider request / response body 全文、endpoint 実値。

## 初手

- path: `internal/service/persona_generation_provider_adapter.go`
- 対象: `BuildPersonaGenerationPrompt`
- 変更種別: NPC ペルソナ生成専用 `PromptBuilder` への移動
- 対応する完了条件: completion_signal 1
- 理由: 原文発話と会話文脈の保護境界を provider adapter から分離するため。

## 完了条件

1. NPC ペルソナ生成の生成指示は 1 NPC を 1 実行単位にし、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ入力として固定する。
2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、応答採用、ペルソナ保存を持たない。
3. 応答欠落、余分な応答、NPC 対応識別子不一致、空のペルソナ本文は NPC 単位の失敗分類として返る。
4. 有効応答だけが翻訳ジョブ内ペルソナまたはペルソナ参照へ採用される。
5. debug log または audit summary は digest と redacted 値だけを出し、原文発話全文と会話文脈全文を出さない。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-implementation-result.BE-03.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因

