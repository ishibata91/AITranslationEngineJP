# observability 結果: OBS-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `OBS-01`
- `implementation_skill`: `observability-implementer`
- `status`: `completed-with-known-backend-local-blocker`

## 判断結果

provider 境界ログの観測追加は完了。

3 フェーズの bulk summary ログへ `provider` を追加した。
provider settings と provider execute の失敗ログは既存の `logProviderBoundaryFailure` / `logProviderBoundarySkipped` で十分なため追加しない。

## 根拠参照

- `observability-input.OBS-01.md:48-57`: 変更許可範囲は backend 実装 handoff で変更済みの provider 境界ログだけ。
- `observability-input.OBS-01.md:72-78`: ログに出してよい値は `event`、`where`、`result`、`provider`、件数、failure kind、redacted reason。
- `observability-input.OBS-01.md:88-93`: 完了条件は provider 境界ログの安全要約、大量ログ抑制、禁止ログ確認、追加しない理由。
- `docs/observability-logging.md:8-14`: backend log は `slog` の JSON log を `stderr` へ出す。
- `docs/observability-logging.md:16-32`: 共通 payload は `event`、`where`、`result` で、必要な場合だけ件数や理由を追加する。
- `docs/observability-logging.md:38-50`: provider 境界の失敗分類と大量処理の集約値は追加対象で、secret、provider raw payload、prompt 全文は出さない。

## 追加ログ

- `internal/service/term_translation_phase_service.go:1051`: 単語翻訳の provider bulk summary 呼び出しへ `plan.Execution.Provider` を渡す。
- `internal/service/term_translation_phase_service.go:1072-1084`: 単語翻訳の bulk summary は `event`、`where`、`result`、`provider`、件数、失敗分類だけを出す。
- `internal/service/persona_generation_phase_service.go:1405`: NPC ペルソナ生成の provider bulk summary 呼び出しへ `run.AIProvider` を渡す。
- `internal/service/persona_generation_phase_service.go:1421-1433`: NPC ペルソナ生成の bulk summary は `event`、`where`、`result`、`provider`、件数、失敗分類だけを出す。
- `internal/service/body_translation_phase_service.go:1599`: 本文翻訳の provider bulk summary 呼び出しへ `loaded.execution.Provider` を渡す。
- `internal/service/body_translation_phase_service.go:1670-1682`: 本文翻訳の bulk summary は `event`、`where`、`result`、`provider`、件数、失敗分類だけを出す。

## 追加しない理由

- provider settings 失敗ログは、既存 helper が `event`、`where`、`result=failed`、`provider`、redacted `reason` を出すため追加しない。
- provider settings skip ログは、既存 helper が `event`、`where`、`result=skipped`、`provider`、redacted `reason`、`count` を出すため追加しない。
- provider execute 失敗ログは、各フェーズで既存 helper を呼び、`provider_unavailable`、`provider_failure`、`invalid_provider_response`、`correlation_error` などの失敗分類だけを出すため追加しない。
- provider adapter 内部ログは追加しない。adapter は raw prompt を provider client へ渡す境界を持つため、ログ追加で raw prompt や provider raw response の誤出力リスクが増える。

## 禁止ログ確認

- secret 本体、API key、credential token はログ payload に追加していない。
- raw prompt、request body 全文、response body 全文はログ payload に追加していない。
- 原文発話全文、会話文脈全文、翻訳本文全文、ペルソナ本文全文はログ payload に追加していない。
- endpoint 実値と外部サービス生データはログ payload に追加していない。
- ループ内の 1 件ごとのログは増やしていない。bulk summary の 1 回ログに provider と集約件数を残す。
- trace ID、constructor 引数、context logger、frontend log 転送は追加していない。

## 変更ファイル

- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/observability-result.OBS-01.md`

## 検証結果

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: failed with known blocker
  - backend lint: pass
  - `internal/...` package tests: pass
  - root package `aitranslationenginejp`: fail
  - failure: `main.go:18:12: pattern all:frontend/dist: no matching files found`

## 影響範囲修正

最初の実装では bulk summary helper の引数順を変更し、既存テストの private helper 呼び出しを壊した。
プロダクトテストは変更せず、helper の既存引数を維持し、provider を任意引数として受ける形へ修正した。

## 最終検証へ進む判断材料

- OBS-01 の追加ログは provider 境界の安全要約に限定している。
- 禁止ログに該当する secret、raw prompt、原文発話全文、会話文脈全文、外部サービス生データは追加していない。
- 内部 package の指定テストは通過している。
- `backend-local` の失敗は入力で記録済みの `frontend/dist` 欠落であり、OBS-01 の変更許可範囲外である。
