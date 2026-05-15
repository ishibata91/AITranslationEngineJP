# Backend レビュー修正結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `agent`: `backend_implementer`
- `status`: `completed`
- `completed_at`: `2026-05-14`

## 変更ファイル

- `internal/service/body_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/term_translation_phase_service.go`

## 指摘別対応

- `behavior-001` / `contract-001`: 本文翻訳 summary は `recoverable_failed` で `CanResume=false`、`CanRetry=true` になる。
- `responsibility-boundary-001`: Service から `translationjobpolicy` は import せず、Service 側は read model 表示値と永続直前 guard に限定した。
- `behavior-002`: 単語翻訳 `PausePhase` は terminal job を拒否し、phase run と job を更新しない。
- `state-invariant-001`: body と persona は transaction 内で現在 job / phase run を再確認し、phase state は `UpdateJobPhaseRunWhenState` で条件付き更新する。

## 検証結果

- `gofmt -l internal/service`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。

## 重要エラー

初回 `backend-local` は本文翻訳 resume 拒否理由が期待値と異なり失敗した。
拒否理由を `phase_not_paused` へ揃えた後、再実行で通過した。
