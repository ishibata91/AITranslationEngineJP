# scenario test 入力: SCN-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `SCN-01`
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `ready_wave`: `wave-4`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

backend 公開接点から、生成指示全文ではなく digest、件数、失敗分類、保護対象を含まない要約だけを観測できることを API テストで証明する。
invalid provider response はフェーズ別の失敗分類として返し、raw provider response を公開しないことを確認する。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md`
- `BE-02`: `./backend-implementation-result.BE-02.md`
- `BE-03`: `./backend-implementation-result.BE-03.md`
- `BE-04`: `./backend-implementation-result.BE-04.md`
- `TU-01`: `./unit-test-result.TU-01.md`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/unit-test-result.TU-01.md`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/coding-guidelines-tests.md`

## 実装済み対象

- `internal/service/prompt_envelope.go`
- `internal/service/term_translation_prompt_builder.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/persona_generation_prompt_builder.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/controller/wails/persona_generation_phase_controller_unit_test.go`
- `internal/controller/wails/body_translation_phase_controller_unit_test.go`

## 変更許可範囲

- `internal/usecase/*_scenario_test.go`
- `internal/controller/wails/*_unit_test.go` のうち backend 公開接点の安全要約を証明する範囲。
- 必要最小限のテスト補助。

## 禁止範囲

- プロダクトコード変更。
- 単体分岐だけの補強。
- frontend 実装。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文。
- `.codex/`。
- 実 AI API 呼び出し。
- 他 agent の変更の取り消し。

## secret 境界

表示、DTO、read model に出してよい値:
provider、model、execution mode、input count、output count、failure kind、`PromptDigest`。

出してはいけない値:
fake secret、raw prompt、request body 全文、response body 全文、原文発話全文、会話文脈全文。

## 証明対象

1. backend 公開接点から、生成指示全文ではなく digest、件数、失敗分類、保護対象を含まない要約だけを観測できる。
2. invalid provider response がフェーズ別の失敗分類として返り、raw provider response は公開されない。
3. 実 AI API を呼ばず、provider stub または contract stub で API テストとして成立している。

## 初手

- path: `internal/usecase/body_translation_phase_scenario_test.go`
- 対象: 本文翻訳開始結果または公開 summary の safe summary
- 変更種別: backend API シナリオテスト追加または既存テスト補強
- 対応する完了条件: completion_signal 1
- 理由: 本文翻訳は保持要素と要約境界の公開確認が最も広い。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/usecase ./internal/controller/wails`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`

## 既知 blocker

- `backend-local` は root package の `frontend/dist` 欠落で失敗する可能性がある。
- この blocker は `backend-implementation-result.BE-01.md` から `BE-04.md`、`unit-test-result.TU-01.md` に記録済みである。
- テスト担当は、backend lint と `internal/...` package tests が通過するかを分けて記録する。

## 期待出力

- `scenario-test-result.SCN-01.md`
- 変更ファイル一覧
- 証明済みシナリオ結果
- 検証結果
- 未証明小範囲
- 残った失敗と原因
