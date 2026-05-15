# 観測ログ追加結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `agent`: `observability_implementer`
- `status`: `completed`
- `completed_at`: `2026-05-14`

## 判定

観測ログ追加は完了した。
UseCase の `translationjobpolicy` 拒否は service 呼び出し前に戻るため、既存の service 実行後ログでは観測できない経路だった。

## 追加ログ

- `internal/usecase/phase_policy_helpers.go`: policy 拒否専用の `slog.WarnContext` を追加した。
- `internal/usecase/term_translation_phase_usecase.go`: 単語翻訳 phase の policy 拒否を記録する。
- `internal/usecase/persona_generation_phase_usecase.go`: NPC ペルソナ生成 phase の policy 拒否を記録する。
- `internal/usecase/body_translation_phase_usecase.go`: 本文翻訳 phase の policy 拒否を記録する。

## payload

- `event`
- `where`
- `result`
- `id`
- `reason`

## 禁止ログ確認

secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文は含めていない。
loop 内ログ、trace ID、constructor 引数追加、context への logger 埋め込みも追加していない。

## 検証結果

- `git diff --check -- internal/usecase/...`: pass。
- `go test ./internal/usecase/...`: pass。

## 重要エラー

初回 `go test ./internal/usecase/...` は存在しない `ReasonInvalidPhaseState` 参照で失敗した。
`ReasonPhaseRunRequired` に修正済み。
