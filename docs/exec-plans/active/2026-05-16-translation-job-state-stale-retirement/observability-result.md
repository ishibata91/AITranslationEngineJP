# Observability Result: 2026-05-16-translation-job-state-stale-retirement

- `skill`: `observability-implementer`
- `status`: `completed`
- `return_to`: `implement_lane`

## 判断結果

観測ログ追加は完了した。
本文翻訳段階の start-on-demand 分岐だけに恒久ログを追加した。

## 根拠参照

- `scenario-design.md`: `Ready` job には `JOB_PHASE_RUN` を事前作成せず、phase start 許可時だけ作成する。
- `implementation-scope.md`: body phase run 未存在時の start-on-demand は、承認済みシナリオ `SCN-TJSR-001` の実装範囲である。
- `implementation-wave-result.md`: 初回シナリオテスト失敗後、`internal/service/body_translation_phase_service.go` で body phase run 未存在時の start-on-demand を実装した。
- `docs/observability-logging.md`: backend は `slog` を使い、`event`、`where`、`result` を基本 payload とする。

## 追加ログ

- 対象: body phase run が未存在の `Ready` job を start した時に、`CreateJobPhaseRun` と初期化が完了した分岐。
- 目的: 実行後に消える「phase run が事前作成済みだったか、start 時に作成されたか」という原因候補を分離する。
- 変更ファイル: `internal/service/body_translation_phase_service.go`
- payload:
  - `event`: `body_translation_phase_run_created_on_start`
  - `where`: `backend.service.body_translation_phase.start`
  - `result`: `created`
  - `id`: `job:<id> phase_run:<id>`
  - `reason`: `missing_precreated_phase_run`

## 追加しない理由

- `JobIOService` 削除と architecture lint 定義削除は、実行時に消える状態や分岐理由を生まないため、ログを追加しない。
- `cancelled` fixture spelling 統一は、test fixture の期待値差分であり、実行時の原因分離情報を生まないため、ログを追加しない。
- 単体テストとシナリオテストは product runtime の観測点ではないため、ログを追加しない。

## 禁止ログ確認

- provider raw payload はログへ出していない。
- prompt 全文と翻訳本文全文はログへ出していない。
- credential 実値、API key、endpoint 実値はログへ出していない。
- 大量処理の loop 内 1 件ごとのログは追加していない。
- trace ID、全 command の start / finish log、frontend から backend へのログ送信は追加していない。

## 変更ファイル

- `internal/service/body_translation_phase_service.go`: start-on-demand で body phase run を作成した分岐を `slog` で記録した。
- `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/observability-result.md`: 観測ログ追加結果を記録した。

## 検証結果

- `go test ./internal/service`: pass

## 検証未実行理由

追加ログは `internal/service/body_translation_phase_service.go` の観測 payload 追加だけである。
API 境界、integration 境界、scenario gate は完成済みテスト成果物で通過済みのため、観測ログ追加では再実行していない。

## 次判断材料

`implement_lane` は最終検証へ進める。
追加ログは承認済み実装範囲内にあり、UI、docs 正本本文、DB schema、Wails DTO、secret / trust boundary は変更していない。
