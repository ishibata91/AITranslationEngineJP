# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: backend-boundary-lint
- task_id: 2026-04-09-backend-boundary-lint
- task_catalog_ref:
- parent_phase: backend-quality-foundation

## 要求要約

- Go backend lint に内部境界違反の検出を追加する。

## 判断根拠

<!-- Decision Basis -->

- `docs/architecture.md` には backend の依存方向が定義されている。
- 現状の backend lint は品質系診断は持つが、層間依存の逸脱は機械的に検出していない。
- 現在の package 構成は薄いため、今のうちに deny rule を入れておくと将来の AI 実装差分を早く止められる。

## 対象範囲

- `.golangci.yml`
- `.go-arch-lint.yml`
- `.gomodguard.yaml`
- `scripts/lint/run-go-backend-lint.sh`
- `package.json`
- `docs/exec-plans/active/`

## 対象外

- `docs/architecture.md` の正本更新
- package 再編
- custom linter 実装

## 依存関係・ブロッカー

- `depguard` の設定形式は `golangci-lint v2` の書式に合わせる必要がある。
- 既存 package は現在ほぼ空なので、まずは architecture で定義済みの禁止方向だけを deny する。

## 並行安全メモ

- shared scope は `.golangci.yml` と exec-plan に限る。

## UI

- N/A

## Scenario

- 開発者は `npm run lint:backend` を実行した時に、backend package が architecture で禁止された internal package を import すると lint failure になる。

## Logic

- `depguard` を有効化し、`internal` 配下の層ごとに禁止 import を deny list として定義する。
- deny rule は `docs/architecture.md` の依存方向だけを対象にし、同層内と許可された下位依存は妨げない。

## 実装計画

<!-- Implementation Plan -->

- `depguard` を backend lint の有効 linter に追加し、production / test 依存ルールを分ける。
- `go-arch-lint` の architecture spec を追加する。
- `gomodguard` の module policy を追加する。
- `npm run lint:backend` と backend harness を実行する。
- plan を completed へ移し、導入済み内容を記録する。

## 受け入れ確認

- `npm run lint:backend` が成功する。
- `python3 scripts/harness/check_backend_lint.py` が成功する。
- deny rule が architecture の依存方向を表現している。

## 必要な証跡

<!-- Required Evidence -->

- `npm run lint:backend`
- `python3 scripts/harness/check_backend_lint.py`
- `python3 scripts/harness/run.py --suite all`

## HITL 状態

- N/A

## 承認記録

- user requested backend boundary lint on 2026-04-09

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- docs 正本への恒久反映が必要なら `updating-docs` へ分離する。

## 結果

<!-- Outcome -->

- 導入済み:
  - `.golangci.yml` で `depguard` を backend lint の有効 linter に追加した。
  - `internal/controller`、`usecase`、`service`、`statemachine`、`jobio`、`repository`、`aiprovider`、`infra/ai`、`infra/runtime` ごとに file rule を定義した。
  - deny list は `docs/architecture.md` の依存方向に合わせ、禁止された internal package への import を失敗にする。
  - production code では `github.com/stretchr/testify`、`github.com/google/go-cmp`、`io/ioutil`、`github.com/pkg/errors` を禁止し、test file には test library だけを許可する rule を追加した。
  - `.go-arch-lint.yml` を追加し、`controller -> usecase -> service -> repository` 系の依存方向と `jobio -> repository`、`infra/ai -> aiprovider, infra/runtime` を architecture gate として固定した。
  - `.gomodguard.yaml` を追加し、`github.com/pkg/errors` の混入禁止、`github.com/mitchellh/go-homedir <= 1.1.0` の version block、local replace directive の禁止を追加した。
  - `package.json` と `scripts/lint/run-go-backend-lint.sh` を更新し、`lint:backend:arch` と `lint:backend:module` を backend lint に統合した。
  - `npm run lint:backend`、`python3 scripts/harness/check_backend_lint.py`、`python3 scripts/harness/run.py --suite all` は成功した。
- 未対応:
  - docs 正本の lint policy / architecture への反映はまだ行っていない。
  - package が増えた時は `depguard` と `go-arch-lint` の rule を追従更新する必要がある。
  - `go test ./...` を使う execution harness は依然として `frontend/node_modules` 配下の Go package まで拾う。
- 証跡:
  - `npm run lint:backend`
  - `python3 scripts/harness/check_backend_lint.py`
  - `python3 scripts/harness/run.py --suite all`
