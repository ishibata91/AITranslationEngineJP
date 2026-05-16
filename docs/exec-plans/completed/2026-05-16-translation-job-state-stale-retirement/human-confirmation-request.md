# Human Confirmation Request

- `caller`: `light_change_lane`
- `status`: `waiting_human`
- `request`: backend 実装後の削減範囲を確認する。

## 確認してほしい結果

- `internal/statemachine/` は旧設計 package として削除されている。
- `.go-arch-lint.yml` から `statemachine` component と許可依存が削除されている。
- UseCase 層の phase 別 policy input 生成は `phasePolicyInput` へ寄っている。
- Service 層の phase 別 action enablement 分岐は `commonPhaseActionAvailability` へ寄っている。
- `internal/jobio/` は未決事項として残っている。
- docs 正本本文、docs 正本図、DB schema、Wails DTO、frontend UI は変更されていない。

## 確認根拠

- `backend-implementation-result.md`
- `design-diff.md`
- `git diff --stat`

## 検証済み

- `gofmt -l internal/usecase internal/service`
- `go test ./internal/usecase ./internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite coverage`
- `git diff --check`

## 人間の返答でほしい判断

- `approved`: 削減範囲を承認し、次のレビュー準備へ進める。
- `changes_requested`: 削減範囲の差し戻し内容を指定する。
- `question`: 追加で確認したい点を指定する。

## 未決事項

- `internal/jobio/` を削除するか、architecture 正本に合わせて別 task で実装するか。
- `docs/exec-plans/active/observability-log-addition/` に残る `StateMachine` / `JobIOService` 旧名参照を、この task 外でどう扱うか。
