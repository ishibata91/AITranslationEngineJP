# unit test 入力: TU-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `TU-01`
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `ready_wave`: `wave-3`
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

3 フェーズの prompt 境界、provider adapter、response validation、redaction を単体テストで証明する。
`PromptDigest` と `REQUEST_V1` 系の値が、生成規則版ではなく内部同一性情報と要求形状識別子に留まることをテスト名または期待値で示す。

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md`
- `BE-02`: `./backend-implementation-result.BE-02.md`
- `BE-03`: `./backend-implementation-result.BE-03.md`
- `BE-04`: `./backend-implementation-result.BE-04.md`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`: pass。
- `git diff --check`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-02.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-03.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-04.md`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/coding-guidelines-tests.md`

## 実装済み対象

- `internal/service/prompt_envelope.go`
- `internal/service/term_translation_prompt_builder.go`
- `internal/service/term_translation_provider_adapter.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_prompt_builder.go`
- `internal/service/persona_generation_provider_adapter.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_prompt_builder.go`
- `internal/service/body_translation_provider_adapter.go`
- `internal/service/body_translation_response_adapter.go`
- `internal/service/body_translation_phase_service.go`

## 変更許可範囲

- `internal/service/prompt_envelope_test.go`
- `internal/service/term_translation_prompt_builder_test.go`
- `internal/service/term_translation_provider_adapter_test.go`
- `internal/service/term_translation_phase_service_test.go`
- `internal/service/persona_generation_provider_adapter_test.go`
- `internal/service/persona_generation_phase_service_test.go`
- `internal/service/body_translation_provider_adapter_test.go`
- `internal/service/body_translation_phase_service_test.go`
- `internal/service/body_translation_field_result_test.go`
- 必要最小限の `internal/service` test helper。

## 禁止範囲

- プロダクトコード変更。
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
digest、件数、redacted summary、failure kind。

出してはいけない値:
fake secret が DTO、read model、log capture、error summary へ出る期待値。
raw prompt、provider request / response body 全文、原文発話全文、会話文脈全文。

## 初手

- path: `internal/service/term_translation_provider_adapter_test.go`
- 対象: 単語翻訳 builder の要求形状識別子と raw prompt 非公開
- 変更種別: 単体テスト追加または既存テスト補強
- 対応する完了条件: completion_signal 1
- 理由: 単語翻訳は最小入力で共有境界を確認しやすい。

## 完了条件

1. 単語翻訳は対象語不一致、空訳語、余分な応答、欠落応答を対象語単位の invalid response として証明している。
2. NPC ペルソナ生成は対応識別子不一致、空ペルソナ、debug log redaction、credential 非公開を証明している。
3. 本文翻訳は翻訳項目識別子不一致、空訳文、保持要素不整合、request summary の raw prompt 非公開を証明している。
4. `PromptDigest` と `REQUEST_V1` 系の値が生成規則版ではなく、要求形状識別子と内部同一性情報である期待をテスト名または期待値で示している。

## 検証コマンド

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite coverage`

## 既知 blocker

- `backend-local` は root package の `frontend/dist` 欠落で失敗する可能性がある。
- この blocker は `backend-implementation-result.BE-01.md` から `BE-04.md` までに記録済みである。
- テスト担当は、backend lint と `internal/...` package tests が通過するかを分けて記録する。

## 期待出力

- `unit-test-result.TU-01.md`
- 変更ファイル一覧
- 検証結果
- coverage 結果または未実行理由
- 未証明小範囲
- 残った失敗と原因
