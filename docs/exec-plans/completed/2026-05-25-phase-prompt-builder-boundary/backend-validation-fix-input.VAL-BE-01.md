# backend validation fix 入力: VAL-BE-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `VAL-BE-01`
- `implementation_artifact`: backend 実装の検証戻し
- `implementation_skill`: implement-backend
- `source_scope`: `./implementation-scope.md`
- `human_review`: approved by human reply `approve` on 2026-05-25

## 目的

最終検証の `coverage` で検出された、今回追加した backend 実装範囲内の Sonar HIGH maintainability issue を解消する。

## 失敗内容

- command: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite coverage`
- issue: `internal/service/prompt_envelope.go:62 go:S1192`
- message: `Define a constant instead of duplicating this literal "sha256:" 3 times.`
- Sonar issue id: `AZ5fEvP5P5vDm0ouqcIm`

## 依存完了

- `BE-01`: `./backend-implementation-result.BE-01.md`
- `OBS-01`: `./observability-result.OBS-01.md`
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/backend-implementation-result.BE-01.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/observability-result.OBS-01.md`
- `internal/service/prompt_envelope.go`
- `internal/service/prompt_envelope_test.go`

## 変更許可範囲

- `internal/service/prompt_envelope.go`
- 必要な場合だけ `internal/service/prompt_envelope_test.go`

## 禁止範囲

- prompt envelope の意味変更。
- 3 フェーズの provider adapter / phase service の挙動変更。
- frontend 実装。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文。
- `.codex/`。
- 他 agent の変更の取り消し。

## 完了条件

1. `sha256:` literal の重複を定数化し、`PromptDigest`、`PromptDigestString`、`RedactedPromptDiagnostic` の意味を変えない。
2. raw prompt、secret、request body 全文、response body 全文を公開する経路を増やさない。
3. `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service` が pass する。

## 期待出力

- `backend-validation-fix-result.VAL-BE-01.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因
