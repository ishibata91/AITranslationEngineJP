# observability 入力: OBS-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `OBS-01`
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `ready_wave`: `wave-5`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

3 フェーズの provider settings、provider execute、bulk summary の恒久ログを安全要約へ揃える。
追加が不要なログは、追加しない理由を根拠参照付きで返す。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md`
- `BE-02`: `./backend-implementation-result.BE-02.md`
- `BE-03`: `./backend-implementation-result.BE-03.md`
- `BE-04`: `./backend-implementation-result.BE-04.md`
- `TU-01`: `./unit-test-result.TU-01.md`
- `SCN-01`: `./scenario-test-result.SCN-01.md`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/unit-test-result.TU-01.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/scenario-test-result.SCN-01.md`
- `docs/observability-logging.md`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/body_translation_provider_adapter.go`

## 完成済み実装成果物

- 共有 `PromptEnvelope`、`PromptDigest`、要求形状識別子、安全要約境界。
- 単語翻訳の専用 `PromptInput` / `PromptBuilder` と provider adapter 接続。
- NPC ペルソナ生成の専用 `PromptInput` / `PromptBuilder` と provider adapter 接続。
- 本文翻訳の専用 `PromptInput` / `PromptBuilder` と provider adapter 接続。
- 単体テストと backend API シナリオテストの安全要約確認。

## 変更許可範囲

- backend 実装 handoff で変更済みの provider 境界ログだけ。
- 具体候補:
  - `internal/service/term_translation_phase_service.go`
  - `internal/service/persona_generation_phase_service.go`
  - `internal/service/body_translation_phase_service.go`
  - `internal/service/term_translation_provider_adapter.go`
  - `internal/service/persona_generation_provider_adapter.go`
  - `internal/service/body_translation_provider_adapter.go`

## 禁止範囲

- 新規機能。
- プロダクトテスト追加。
- frontend 実装。
- UI 表示、画面、部品、文言、style。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文。
- `.codex/`。
- secret、raw prompt、原文発話全文、会話文脈全文、外部サービス生データをログへ出す変更。
- 他 agent の変更の取り消し。

## secret 境界

ログに出してよい値:
event、where、result、provider、件数、failure kind、redacted reason。

ログに出してはいけない値:
secret 本体、raw prompt、request body 全文、response body 全文、endpoint 実値、原文発話全文、会話文脈全文。

## 初手

- path: `internal/service/persona_generation_phase_service.go`
- 対象: provider execute 失敗ログ
- 変更種別: 変更または追加しない理由の固定
- 対応する完了条件: completion_signal 1
- 理由: NPC ペルソナ生成は保護対象が多いため、禁止ログを最初に確認する。

## 完了条件

1. 各フェーズの provider 境界ログは event、where、result、provider、件数、失敗分類だけで原因候補を分離できる。
2. 同種の大量ログを増やさず、bulk summary は集約件数を優先している。
3. secret、raw prompt、原文発話全文、会話文脈全文、外部サービス生データをログに出さない確認結果が返る。
4. 既存ログで十分な対象は、追加しない理由を根拠参照付きで返る。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`

## 既知 blocker

- `backend-local` は root package の `frontend/dist` 欠落で失敗する可能性がある。
- coverage は frontend import resolution と Sonar gate で失敗する可能性がある。
- 上記は OBS-01 の変更許可範囲外である。

## 期待出力

- `observability-result.OBS-01.md`
- 追加ログまたは追加しない理由
- 禁止ログ確認
- 変更ファイル一覧
- 検証結果
- 最終検証へ進む判断材料
