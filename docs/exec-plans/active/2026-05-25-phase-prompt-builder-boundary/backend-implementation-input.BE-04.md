# backend 実装入力: BE-04

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-04`
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `ready_wave`: `wave-2`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

本文翻訳フェーズの既存 prompt builder を共有境界へ揃える。
1 翻訳項目を 1 実行単位にし、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、辞書制約、保持要素を同じ生成指示入力へ固定する。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md` に `implemented-with-validation-blocker` として記録済み。
- `internal/service/prompt_envelope.go`: 共有 `PromptEnvelope`、`PromptDigest`、要求形状識別子、安全要約境界を追加済み。
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/detail-specs/body-translation-phase.md`
- `internal/service/prompt_envelope.go`
- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`
- `internal/service/body_translation_response_adapter.go`
- `internal/service/body_translation_phase_service.go`

## 変更許可範囲

- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`
- `internal/service/body_translation_response_adapter.go`
- `internal/service/body_translation_phase_service.go`
- 上記範囲の backend test。

## 禁止範囲

- 共有 `PromptEnvelope` 契約の意味変更。必要な場合は実装を止めて報告する。
- 単語翻訳、NPC ペルソナ生成の実装変更。
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
`CredentialRef`、provider、model、execution mode、request unit id、field correlation key、input count、output count、protected element count、failure kind、`PromptDigest`。

出してはいけない値:
secret 本体、raw prompt、provider request / response body 全文、生成指示の原文全文、外部サービスとの生データ、endpoint 実値。

## 初手

- path: `internal/service/body_translation_prompt_builder.go`
- 対象: `BuildBodyTranslationPromptEnvelope`
- 変更種別: 本文翻訳専用 `PromptInput` と共有 `PromptEnvelope` の整合
- 対応する完了条件: completion_signal 1
- 理由: 本文翻訳だけ既存 builder があるため、共有境界への整合を先に閉じるため。

## 完了条件

1. 本文翻訳の生成指示は 1 翻訳項目を 1 実行単位にし、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、辞書制約、保持要素を同じ入力として固定する。
2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、採用判断、保存を持たない。
3. 応答欠落、余分な応答、翻訳項目識別子不一致、空訳文、保持要素不整合は翻訳項目単位の失敗分類として返る。
4. 有効応答だけが翻訳項目単位の訳文候補と保護要素検証対象へ採用される。
5. request summary は件数、digest、識別子、保持要素件数に限定し、raw prompt と外部サービス生データを出さない。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-implementation-result.BE-04.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因

