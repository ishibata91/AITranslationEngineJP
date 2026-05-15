# Backend レビュー修正 3 入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 修正理由

5 観点レビュー再実行 2 で `state-invariant-002` が open になった。
単語翻訳 phase の resume / retry が `UpdateJobPhaseRunWhenState` による required state 条件付き更新を使っていない。

## 対象指摘

- file: `reviewback.state-invariant.yaml`
- issue id: `state-invariant-002`
- level: `major`
- status: `open`

## 修正対象

product code:
- `internal/service/term_translation_phase_service.go`

product test:
- `internal/service/term_translation_phase_service_test.go`

## 期待する修正

単語翻訳 phase の resume / retry で、同じ `JOB_PHASE_RUN` を required state 条件付きで `running` に更新する。

期待する required state:
- resume: `paused`
- retry: `recoverable_failed`

保持する条件:
- terminal job 拒否を維持する。
- phase run id 一致確認を維持する。
- 同じ `JOB_PHASE_RUN` 継続を維持する。
- runtime snapshot の保存順序を壊さない。
- `RecoverableFailed` の resume は拒否し、retry だけ許可する。
- Service から `translationjobpolicy` を import しない。
- provider raw payload、secret、API key、prompt 全文、翻訳本文全文をログへ追加しない。

## 期待するテスト

`internal/service/term_translation_phase_service_test.go` に次の確認を追加または既存テストへ追加する。

- resume が `UpdateJobPhaseRunWhenState` に expected state `paused` を渡す。
- retry が `UpdateJobPhaseRunWhenState` に expected state `recoverable_failed` を渡す。
- expected state 不一致相当の `repository.ErrConflict` では provider request を発生させず、phase run を running にしない。

## 禁止範囲

- frontend、Wails DTO、DB schema、migration を変更しない。
- docs 正本、`.codex`、作業計画文書を変更しない。
- 別 task `2026-05-13-notification-module-dependency-separation` を変更しない。
- 初回レビューと再レビューで resolved になった修正を戻さない。

## 検証コマンド

- `gofmt -l internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`

## 返却内容

- backend 修正の完了、未完了、停止の判定。
- 変更ファイル。
- resume / retry の expected state を固定した方法。
- 追加または更新したテスト。
- 実行した検証コマンドと結果。
- 未実行検証がある場合は未実行理由。
