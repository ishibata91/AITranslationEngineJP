# backend 実装入力: BE-02

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-02`
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `ready_wave`: `wave-2`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

単語翻訳フェーズの生成指示責務を分離する。
1 対象語を 1 実行単位にし、対象語、原文言語、訳文言語、応答対応識別子を同じ生成指示入力へ固定する。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md` に `implemented-with-validation-blocker` として記録済み。
- `internal/service/prompt_envelope.go`: 共有 `PromptEnvelope`、`PromptDigest`、要求形状識別子、安全要約境界を追加済み。
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/detail-specs/term-translation-phase.md`
- `internal/service/prompt_envelope.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/term_translation_phase_service.go`

## 変更許可範囲

- `internal/service/term_translation_provider_adapter.go`
- `internal/service/term_translation_phase_service.go`
- 必要な単語翻訳専用 prompt builder file。
- 上記範囲の backend test。

## 禁止範囲

- 共有 `PromptEnvelope` 契約の意味変更。必要な場合は実装を止めて報告する。
- NPC ペルソナ生成、本文翻訳の実装変更。
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
provider、model、execution mode、input count、output count、failure kind、`PromptDigest`。

出してはいけない値:
API key、raw prompt、source term 一覧の全文 dump、provider request / response body 全文、endpoint 実値。

## 初手

- path: `internal/service/term_translation_provider_adapter.go`
- 対象: `BuildTermTranslationPrompt`
- 変更種別: 単語翻訳専用 `PromptBuilder` への移動
- 対応する完了条件: completion_signal 1
- 理由: provider adapter から prompt 文言組み立て責務を外すため。

## 完了条件

1. 単語翻訳の生成指示は 1 対象語を 1 実行単位にし、対象語、原文言語、訳文言語、応答対応識別子を同じ入力として固定する。
2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、応答採用、辞書保存を持たない。
3. 応答欠落、余分な応答、空訳語、対象語不一致は対象語単位の失敗分類として返る。
4. 有効応答だけが確定訳語として対象ジョブ内辞書へ採用される。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-implementation-result.BE-02.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因
