# Task Plan: observability-logging-foundation

- `workflow`: work
- `status`: planned
- `lane_owner`: designer
- `task_id`: observability-logging-foundation
- `task_mode`: implementation-planning
- `request_summary`: frontend は `pino` で console へ出し、backend は `log/slog` で観測ログの恒久基盤を作る。
- `goal`: frontend と backend の両方で、構造化ログ、trace ID、原因分離に必要な観測情報を残せる基盤を固定する。frontend のログは agent-browser で取得できる console 証跡を第一出力先にする。
- `constraints`: trace ID を全関数引数へ広げない。`console.log` 直書きを恒久実装にしない。secret、API key、provider raw payload、過剰本文をログへ出さない。ループ内で大量ログを出さない。
- `close_conditions`: plan と scenario-design が human review 可能であり、implementation-scope 作成前の判断材料が揃っている。

## Artifact Index

- `ux_task_frame`: `N/A`
- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `ux_review`: `N/A`
- `frontend_human_review`: `N/A`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`: `docs/tech-selection.md`, `docs/coding-guidelines-backend.md`, `docs/coding-guidelines-frontend.md`, `.codex/skills/observability-implementer/SKILL.md`
- `canonicalization_targets`: `docs/tech-selection.md`, `docs/coding-guidelines-backend.md`, `docs/coding-guidelines-frontend.md`
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`: `npm --prefix frontend run lint`, `go test ./...`

## HITL Status

- `functional_or_design_hitl`: `required-after-design-bundle`
- `ux_review`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `pending-after-design-bundle`

## Decision Record

### 決定

frontend の恒久ログ基盤は `pino` を第一候補にする。
frontend のログ出力先は browser console にする。
backend の恒久ログ基盤は既存方針どおり `log/slog` を使う。
backend の診断ログ保存先は `tmp/observability-logging-foundation/diagnostic-log.sqlite` にする。
trace ID は処理入口で生成し、logger scope または context に束ねる。

### 理由

構造化ログがない場合、実行時にしか確定しない値と実行後に消える分岐理由を後続調査で再構成できない。
trace ID を全関数引数へ広げる方式は、業務処理の署名変更が大きくなり、観測基盤のために実装責務が汚れる。
`pino` と `slog` は、key-value の構造化ログへ寄せやすい。
frontend は操作を前提に観測するため、agent-browser で取得できる console 証跡を第一出力先にすると、Wails command をログ輸送路にしなくて済む。

### 影響

frontend 実装者は `console.log` ではなく、repo 側の診断 logger 境界を使う。
frontend の診断 logger 境界は `pino` を通じて console へ構造化ログを出す。
backend 実装者は `context.Context` と `slog.Logger` の scope を使い、必要な境界だけで trace ID を付与する。
backend 診断ログは app data ではなく `tmp/` 配下の調査用 DB に残す。
観測ログ追加 agent は、完成済み成果物に対して既存基盤を使ったログ追加を判断できる。

### 未決事項

backend の app 起動時 logger 初期化場所と log level の設定方法は未決である。
trace ID の名前を `trace_id` にするか `traceId` にするかは、frontend と backend の検索性を見て固定する。

## Implementation Planning Notes

- backend: app 起動時に `slog.Logger` を初期化し、Wails controller 境界で trace ID を `context.Context` へ束ねる。
- backend: 診断 SQLite は `tmp/observability-logging-foundation/diagnostic-log.sqlite` に置く。
- backend: usecase、service、repository は必要な境界だけで logger scope を受け取り、secret と raw payload を出さない。
- frontend: `FrontendDiagnosticLogger` を repo 境界として定義し、実装に `pino` を使って console へ出す。
- frontend: 画面 controller、gateway、runtime event adapter の境界で logger scope を作る。
- frontend: Wails command へログを逐次送信しない。frontend SQL 永続化も初期対象にしない。
- common: ループや大量処理は件数、分類、代表 ID、最初の失敗、最後の失敗を優先する。

## CLI Check

- `sqlite3 --version`: `3.51.0`
- `sqlite3 tmp/observability-logging-foundation/diagnostic-log.sqlite`: create table、insert、select を確認済み。
- `select result`: `1|info|cli-check|trace-cli-check|diagnostic_sqlite_cli_check`

## Codex Implementation Result

- `completed_handoffs`: backend 観測ログ基盤、frontend 観測ログ基盤
- `touched_files`: `internal/infra/runtime/diagnostic_log.go`, `internal/infra/runtime/diagnostic_log_test.go`, `frontend/src/application/diagnostic/frontend-diagnostic-logger.ts`, `frontend/src/application/diagnostic/index.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts`, `frontend/package.json`, `frontend/package-lock.json`
- `implemented_scope`: backend は `slog.Handler` から `tmp/observability-logging-foundation/diagnostic-log.sqlite` へ構造化ログを書ける。frontend は `pino` を repo 側の診断 logger 境界から browser console へ出せる。
- `test_results`: `go test ./...` 通過。`npm --prefix frontend run lint` 通過。`npm --prefix frontend run test` 通過。`npm --prefix frontend run build` 通過。`python3 scripts/harness/run.py --suite frontend-local` 通過。`python3 scripts/harness/run.py --suite backend-local` 通過。`git diff --check` 通過。
- `implementation_investigation`: backend-local 初回失敗は backend lint で、export コメント不足、`QueryRowContext` 未使用、arch component 未所属が原因。診断ログ基盤を `internal/infra/runtime` へ移して解消した。
- `ui_evidence`: `N/A`
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: `N/A`
- `sonar_gate_result`: `N/A`
- `residual_risks`: 既存業務処理へのログ差し込みは未実施。backend app 起動時の logger 初期化場所と log level 設定は未固定。
- `docs_changes`: `N/A`

## Closeout Notes

- `canonicalized_artifacts`: `pending`
- `detail_spec_canonicalization`: `N/A`
- `follow_up`: human review 後に `implementation-scope.md` を作成する。

## Outcome

- planned
