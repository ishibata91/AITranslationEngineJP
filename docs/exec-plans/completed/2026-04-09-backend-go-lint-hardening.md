# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: backend-go-lint-hardening
- task_id: 2026-04-09-backend-go-lint-hardening
- task_catalog_ref:
- parent_phase: backend-quality-foundation

## 要求要約

- Go lint を frontend 並みに強化し、AI 駆動の実装で取り込みたい静的診断を backend gate に追加する。

## 判断根拠

<!-- Decision Basis -->

- 現状の backend lint は `gofmt` と `go vet` だけであり、未使用、エラー処理、危険な記述、スタイル逸脱の多くは gate に乗っていない。
- `go vet ./...` と `go test ./...` は `frontend/node_modules` 配下の Go package まで拾っており、backend gate の対象範囲が不安定である。
- repo の source of truth docs は human 先行更新ルールがあるため、今回は `exec-plan` と実装側の command / config に閉じて導入する。

## 対象範囲

- `package.json`
- `scripts/harness/check_backend_lint.py`
- repo root の Go lint 設定 file
- repo root の Go lint helper script
- `docs/exec-plans/active/`

## 対象外

- `docs/tech-selection.md` と `docs/lint-policy.md` の正本更新
- CI workflow の追加
- backend 実装ロジック自体の拡張

## 依存関係・ブロッカー

- 導入ツールの設定は official docs に合わせる必要がある。
- 強めの linter 追加で既存コードに新規指摘が出る可能性がある。

## 並行安全メモ

- shared scope は `package.json`、root config、lint helper script、active exec-plan に限る。

## UI

- N/A

## Scenario

- 開発者は repo root で `npm run lint:backend` を実行し、backend package の format / vet / 高度 static analysis をまとめて検証できる。
- backend lint は `frontend/node_modules` 配下の Go package を対象に含めない。

## Logic

- backend lint の対象 package は `internal/...` に固定する。
- 高度 static analysis には official docs で設定方法を確認した `golangci-lint` を採用し、`staticcheck` を含む複数 linter を明示的に有効化する。
- 既存 harness は `lint:backend` を呼ぶ構造のため、その入口を強化しても harness 契約は維持する。

## 実装計画

<!-- Implementation Plan -->

- backend package 範囲を固定する helper を用意する。
- `golangci-lint` 用の repo root 設定を追加する。
- `lint:backend` を format / vet / advanced lint の分割実行へ更新する。
- 必要なら新規指摘に合わせて最小限のコード修正を入れる。
- active plan に導入済み範囲と未対応範囲を追記する。

## 受け入れ確認

- `npm run lint:backend` が成功する。
- `python3 scripts/harness/check_backend_lint.py` が成功する。
- `python3 scripts/harness/run.py --suite all` が成功する。
- active plan に導入済み範囲と未対応範囲が残る。

## 必要な証跡

<!-- Required Evidence -->

- `npm run lint:backend`
- `python3 scripts/harness/check_backend_lint.py`
- `python3 scripts/harness/run.py --suite all`

## HITL 状態

- N/A

## 承認記録

- user requested stronger Go lint rollout on 2026-04-09

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- docs 正本更新が必要になった場合は `updating-docs` へ分離する。

## 結果

<!-- Outcome -->

- 導入済み:
  - repo root に `.golangci.yml` を追加し、official docs に沿った `version: "2"` 設定で `golangci-lint` を導入した。
  - advanced lint では `staticcheck`、`errcheck`、`gosec`、`revive`、`unparam`、`wrapcheck` を含む複数 linter を明示的に有効化した。
  - `scripts/lint/run-go-backend-lint.sh` を追加し、backend package pattern を `./internal/...` に固定した。
  - `package.json` の `lint:backend` を `lint:backend:format`、`lint:backend:vet`、`lint:backend:static` の分割実行へ更新した。
  - `check_backend_lint.py` からは従来どおり `lint:backend` を呼び、強化後の backend lint を harness 契約そのままで実行できる。
  - `internal/*/doc.go` と `main.go` に package comment を追加し、`internal/controller/wails` には `doc.go` を新設した。
  - `internal/controller/wails/app_controller.go` に exported symbol comment を追加し、未使用 lifecycle parameter を `_` に置き換えた。
- 未対応:
  - `go test ./...` を使う execution harness は依然として `frontend/node_modules` 配下の Go package まで拾う。
  - backend lint の導入内容は `docs/tech-selection.md` と `docs/lint-policy.md` の正本へはまだ反映していない。
- 証跡:
  - `npm run lint:backend`
  - `python3 scripts/harness/check_backend_lint.py`
  - `python3 scripts/harness/run.py --suite all`
