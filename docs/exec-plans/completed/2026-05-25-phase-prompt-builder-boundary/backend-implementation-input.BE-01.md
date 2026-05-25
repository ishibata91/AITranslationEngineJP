# backend 実装入力: BE-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `BE-01`
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `ready_wave`: `wave-1`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

3 フェーズ共通の生成指示受け渡し単位を固定する。
raw prompt、digest、要求形状識別子、利用者向け要約を混同しない。

## 依存完了

- `detail-spec-diff.md`: 承認済み。
- `design-diff.phase-prompt-builder-boundary.md`: 承認済み。
- `implementation-scope.md`: `ready-for-implementation`。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`

## 変更許可範囲

- `internal/service` の共有 prompt envelope / digest / redaction helper。
- 上記範囲の backend test。

## 禁止範囲

- frontend 実装。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文。
- `.codex/`。
- フェーズ追加、phase type 追加、ジョブ状態機械変更。
- 実 AI API を使う検証。

## secret 境界

表示、DTO、read model に出してよい値:
`credential_ref`、provider、model、execution mode、`PromptDigest`、入出力件数、失敗分類、保護対象を含まない要約。

出してはいけない値:
secret 本体、raw prompt、request body 全文、response body 全文、原文発話全文、会話文脈全文、API key、token、接続先実値。

## 初手

- path: `internal/service`
- 対象: `PromptEnvelope` 相当の型
- 変更種別: 共有受け渡し単位の追加
- 対応する完了条件: completion_signal 1
- 理由: 共有型が先にないと各フェーズが別々の raw prompt 受け渡しを作るため。

## 完了条件

1. 共有受け渡し単位が raw prompt、digest、要求形状識別子、利用者向け要約を区別している。
2. 共有 helper が secret 本体、raw prompt、request body 全文、response body 全文を log、DTO、read model へ出す形を持たない。
3. `PromptDigest` は復元不能な同一性情報として扱われ、生成規則の版選択値になっていない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-implementation-result.BE-01.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因
