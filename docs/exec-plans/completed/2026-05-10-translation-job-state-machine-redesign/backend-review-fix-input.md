# Backend レビュー修正入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 修正理由

5 観点レビューで未解決 `major` が残った。
`review-aggregation.md` により `implementation_action=fix` とする。

## 修正対象レビュー指摘

- `behavior-001`: 本文翻訳 read model が `recoverable_failed` の `resume` を許可表示する。
- `contract-001`: 本文翻訳 summary の resume 可否が command 拒否契約と一致しない。
- `responsibility-boundary-001`: 本文翻訳 Service の read model が resume 共通操作規則を独自保持している。
- `behavior-002`: 単語翻訳 Service の `PausePhase` が terminal job を拒否しない。
- `state-invariant-001`: body と persona の状態変更 guard が永続更新と原子的でない。

## 修正対象ファイル

- `internal/service/body_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/term_translation_phase_service.go`

必要な場合だけ、同じ責務内の service fake ではなく product service code を変更する。
プロダクトテストは変更しない。

## 期待する修正

- 本文翻訳 summary は `recoverable_failed` で `CanResume=false`、`CanRetry=true` を返す。
- 本文翻訳 summary の `ResumeBlockedReason` は、`recoverable_failed` の時に resume 不可理由を返す。
- 単語翻訳 `PausePhase` は terminal job を拒否し、phase run と job を更新しない。
- body の resume / retry / cancel は、永続更新直前または transaction 内で現在 job state と現在 phase state を再確認する。
- persona の resume / retry / cancel は、永続更新直前または transaction 内で現在 job state と現在 phase state を再確認する。
- phase state 更新では、可能な範囲で既存の `UpdateJobPhaseRunWhenState` を使い、required state と異なる現在状態を更新しない。
- Service から `translationjobpolicy` を import しない。

## 禁止範囲

- プロダクトテスト、検証データ、snapshot、test helper を変更しない。
- frontend、Wails DTO、DB schema、migration を変更しない。
- docs 正本、`.codex`、作業計画文書を変更しない。
- repository schema の意味拡張をしない。
- provider raw payload、secret、API key、credential 参照実値をログへ追加しない。

## 参照

- `review-aggregation.md`
- `reviewback.behavior.yaml`
- `reviewback.contract.yaml`
- `reviewback.responsibility-boundary.yaml`
- `reviewback.state-invariant.yaml`
- `reviewback.trust-boundary.yaml`
- `implementation-scope.md`
- `backend-implementation-result.md`
- `final-validation.md`

## 検証コマンド

- `gofmt -l internal/service`
- `python3 scripts/harness/run.py --suite backend-local`

## 返却内容

- backend 修正の完了、未完了、停止の判定。
- 変更ファイル。
- 各レビュー指摘に対する修正内容。
- 実行した検証コマンドと結果。
- 未実行検証がある場合は未実行理由。
- 後続の unit / scenario tester へ渡す事項。
