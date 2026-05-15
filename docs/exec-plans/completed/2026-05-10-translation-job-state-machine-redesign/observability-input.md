# 観測ログ追加入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `observability_implementer`
- `skill`: `observability-implementer`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 完成済み実装成果物

- `backend-implementation-result.md`
- `unit-test-result.md`
- `scenario-test-result.md`

## 変更ファイル

- `.go-arch-lint.yml`
- `internal/usecase/translationjobpolicy/policy.go`
- `internal/usecase/phase_policy_helpers.go`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`

## 判定対象

UseCase の policy 拒否は、実行後に消える分岐理由である。
既存の phase command logging は service 実行後を主に記録しており、policy 拒否は service を呼ばない。
観測ログ追加が必要かを判定する。

## 追加する場合の制約

- `event`、`where`、`result` を共通 payload にする。
- 必要な場合だけ `id` と `reason` を追加する。
- `where` は UseCase 境界を示す。
- secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文は出さない。
- loop 内の 1 件ごとのログを追加しない。
- constructor 引数を広げない。
- context へ logger を埋め込まない。

## 禁止範囲

- プロダクトテストは変更しない。
- docs 正本、`.codex`、作業計画文書は変更しない。
- frontend、Wails DTO、DB schema、migration は変更しない。
- 新規機能実装、恒久修正は行わない。

## 返却内容

- 観測ログ追加の完了、追加不要、停止の判定。
- 追加ログまたは追加しない理由。
- 禁止ログ確認。
- 変更ファイル。
- 検証未実行理由または最終検証へ渡す理由。
