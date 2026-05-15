# Backend 検証失敗修正入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 修正理由

レビュー修正後の `python3 scripts/harness/run.py --suite coverage` が失敗した。
失敗原因は Sonar maintainability high issue である。

## 失敗内容

- file: `internal/service/persona_generation_phase_service.go`
- line: `684`
- rule: `go:S3776`
- issue: `Refactor this method to reduce its Cognitive Complexity from 21 to the 15 allowed.`
- Sonar issue id: `AZ4mjQJAlhLVGXtYvDPR`

## 修正対象

- `internal/service/persona_generation_phase_service.go`

## 期待する修正

`mutatePhaseState` の transaction 内再確認処理を helper へ切り出し、cognitive complexity を閾値以下へ下げる。
挙動は変えない。

保持する条件:
- resume / retry / cancel は transaction 内で現在 job state と現在 phase state を再確認する。
- phase state 更新は `UpdateJobPhaseRunWhenState` で required state 不一致を更新しない。
- Service から `translationjobpolicy` を import しない。
- rejection payload、snapshot persistence、同じ `JOB_PHASE_RUN` 継続を変えない。

## 禁止範囲

- プロダクトテスト、検証データ、snapshot、test helper を変更しない。
- frontend、Wails DTO、DB schema、migration を変更しない。
- docs 正本、`.codex`、作業計画文書を変更しない。
- provider raw payload、secret、API key、credential 参照実値をログへ追加しない。

## 検証コマンド

- `gofmt -l internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`

## 返却内容

- backend 修正の完了、未完了、停止の判定。
- 変更ファイル。
- complexity を下げた方法。
- 実行した検証コマンドと結果。
- 未実行検証がある場合は未実行理由。
