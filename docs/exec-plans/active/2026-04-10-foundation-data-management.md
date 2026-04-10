# 実装計画

- workflow: impl
- status: planned
- lane_owner: codex
- scope: tasks/usecases/foundation-data-management.yaml, docs/exec-plans, docs/screen-design, frontend, internal
- task_id: 2026-04-10-foundation-data-management
- task_catalog_ref: tasks/usecases/foundation-data-management.yaml
- parent_phase: phase-1-distill

## 要求要約

- Foundation Data 画面で、マスターペルソナ一覧とマスター辞書一覧を観測できるようにする
- 選択中の基盤エントリ詳細を確認できるようにする
- 基盤データの編集導線と Rebuild 導線を確認できるようにする

## 判断根拠

- usecase `foundation-data-management` は、Foundation Data 画面から Persona / Dictionary の切り替え、詳細確認、編集導線、Rebuild 導線の確認を完了条件に置いている
- `docs/spec.md` は、基盤データを `マスターペルソナ` と `マスター辞書` と定義し、UI から観測可能であることを恒久要件に置いている
- `docs/architecture.md` は、`Wails + Go + Svelte` と `Controller -> UseCase -> Service -> Repository` の責務境界を正本としている
- `docs/screen-design/` 直下に `app-shell.md` と `foundation-data.md` が現時点で存在せず、画面詳細は第2段階で task-local artifact として補う必要がある
- 現行 frontend 実装として確認できた実体は `frontend/src/ui/App.svelte`、`frontend/src/ui/views/AppShell.svelte`、`frontend/src/ui/stores/shell-state.ts` であり、Foundation Data 専用画面は未実装である
- 現行 backend の Wails controller は `internal/controller/wails/app_controller.go` の `Health()` のみであり、Foundation Data 用 API は未確認である

## 対象範囲

- `tasks/usecases/foundation-data-management.yaml`
- `docs/exec-plans/active/2026-04-10-foundation-data-management.md`
- `docs/exec-plans/active/2026-04-10-foundation-data-management.ui.html`
- `docs/exec-plans/active/2026-04-10-foundation-data-management.scenario.md`
- `frontend/`
- `internal/`

## 対象外

- `docs/spec.md` など `docs/` 正本の恒久仕様変更
- Foundation Data 以外の画面実装
- 基盤データ以外の翻訳フロー実装

## 依存関係・ブロッカー

- 第2段階で Foundation Data 画面の task-local design を固定する必要がある
- `app-shell.md` と `foundation-data.md` が未作成のため、既存 screen-design 正本だけでは画面詳細を確定できない
- Foundation Data に対応する frontend component、UseCase、Repository、DTO の有無は第2段階以降の詳細設計で切り分ける必要がある
- 編集導線と Rebuild 導線が既存 API の接続で足りるか、新規 API が必要かは未確定である
- design review reroute: UI モックに未選択時プレースホルダ、Dictionary 側 action rail、取得失敗状態の固定が不足している
- design review reroute: mock に検索入力と filter chip が含まれており、usecase / scenario にない挙動を先行露出している
- 人間承認を `承認記録` と `HITL 状態` に固定するまで実装へ進めない

## 並行安全メモ

- `frontend/` と `internal/` の両方にまたがる可能性があるため、第2段階で task group と依存関係を明示してから並列 handoff する
- task-local artifact は active exec-plan 配下に限定し、`docs/` 正本へは反映しない

## UI モック

- `artifact_path`: `docs/exec-plans/active/2026-04-10-foundation-data-management.ui.html`
- `summary`: Foundation Data 画面の一覧、詳細、編集導線、Rebuild 導線を task-local wireframe として固定する
- `note`: app shell 上端 nav から Foundation Data を開き、collection switch と detail pane を同一画面で観測する
- `note`: 一覧は `マスターペルソナ` と `マスター辞書` を segmented control で切り替え、検索や filter など usecase 外の導線は置かない
- `note`: 未選択初期状態では detail placeholder を表示し、編集導線は disabled または非表示相当で誤操作不能にする
- `note`: Persona と Dictionary の両方で action rail を固定し、Rebuild target 表示が選択中 collection に一致するようにする
- `note`: 初期ロード失敗時は成功状態と誤認しない failure surface を同一 artifact 内で示す

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/2026-04-10-foundation-data-management.scenario.md`
- `summary`:
  - app-shell から Foundation Data への遷移、Persona / Dictionary 切替、詳細追従、編集導線、Rebuild 導線を SCN-FDM-001..006 で固定
  - 主要例外として基盤データ取得失敗時の観測可能性を SCN-FDM-007 で固定
  - 責務境界として UI から backend 参照が Wails 境界経由であることを SCN-FDM-008 で固定

## 実装計画

<!-- Implementation Plan -->

- `ordered_scope`:
  1. `foundation-data-backend-contract`: Foundation Data 一覧取得、詳細取得、編集導線、Rebuild 導線を返す backend bind / DTO / usecase / service / repository の責務境界を固定する
  2. `foundation-data-frontend-application`: backend contract を受ける frontend gateway contract / wails adapter / screen state を固定する
  3. `foundation-data-ui-shell`: App Shell から Foundation Data へ遷移し、Persona / Dictionary 切り替え、詳細表示、編集導線、Rebuild 導線を観測できる UI を固定する
- `parallel_task_groups`:
  - `group_id`: foundation-data-contract
  - `scope`: `internal/controller/wails/`, `internal/usecase/`, `internal/service/`, `internal/repository/`
  - `can_run_in_parallel_with`: []
  - `blocked_by`: [`phase-1-distill`]
  - `completion_signal`: Foundation Data query / command DTO と Wails bind method 名、UseCase 呼び出し境界、Repository 依存先が固定される
  - `tasks`:
    - `task_id`: backend-query-command-contract
    - `owned_scope`: `internal/controller/wails/`, `internal/usecase/`, `internal/service/`, `internal/repository/`
    - `focus`: Persona / Dictionary collection 切り替え、詳細取得、編集導線、Rebuild 導線の backend contract を定義する
    - `required_reading`: `tasks/usecases/foundation-data-management.yaml`, `docs/spec.md`, `docs/architecture.md`, `internal/controller/wails/app_controller.go`
    - `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `go test ./...`
  - `group_id`: foundation-data-frontend-core
  - `scope`: `frontend/src/application/`, `frontend/src/controller/wails/`
  - `can_run_in_parallel_with`: [`foundation-data-ui-shell`]
  - `blocked_by`: [`foundation-data-contract`]
  - `completion_signal`: frontend gateway DTO mapping、screen state、action entrypoint が backend contract に追従して固定される
  - `tasks`:
    - `task_id`: frontend-gateway-and-screen-state
    - `owned_scope`: `frontend/src/application/`, `frontend/src/controller/wails/`
    - `focus`: Foundation Data response を UI 用 state へ写像し、一覧読込、選択切替、編集導線、Rebuild 導線の操作入口を提供する
    - `required_reading`: `docs/spec.md`, `docs/architecture.md`, `frontend/src/application/README.md`, `frontend/src/controller/wails/README.md`, `internal/controller/wails/app_controller.go`
    - `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `npm --prefix frontend run check`, `npm --prefix frontend run test`
  - `group_id`: foundation-data-ui-shell
  - `scope`: `frontend/src/ui/`
  - `can_run_in_parallel_with`: [`foundation-data-frontend-core`]
  - `blocked_by`: [`foundation-data-contract`]
  - `completion_signal`: App Shell navigation と Foundation Data view が screen state と callback prop で接続される
  - `tasks`:
    - `task_id`: ui-shell-and-foundation-data-view
    - `owned_scope`: `frontend/src/ui/`
    - `focus`: App Shell navigation、collection toggle、detail pane、edit / rebuild affordance を UI に反映する
    - `required_reading`: `frontend/src/ui/App.svelte`, `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/stores/shell-state.ts`, `docs/exec-plans/active/2026-04-10-foundation-data-management.ui.html`, `docs/exec-plans/active/2026-04-10-foundation-data-management.scenario.md`
    - `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `npm --prefix frontend run check`, `npm --prefix frontend run test`
- `task_dependencies`:
  - `backend-query-command-contract` -> []
  - `frontend-gateway-and-screen-state` -> [`backend-query-command-contract`]
  - `ui-shell-and-foundation-data-view` -> [`backend-query-command-contract`]
  - `ui-shell-and-foundation-data-view` -> [`frontend-gateway-and-screen-state`] for live data binding
- `owned_scope`: backend contract first, then frontend application and UI split. backend は Foundation Data bind / DTO / orchestration を持ち、frontend は gateway mapping / screen state / shell view を持つ。
- `required_reading`:
  - `tasks/usecases/foundation-data-management.yaml`
  - `docs/exec-plans/active/2026-04-10-foundation-data-management.md`
  - `docs/exec-plans/active/2026-04-10-foundation-data-management.ui.html`
  - `docs/exec-plans/active/2026-04-10-foundation-data-management.scenario.md`
  - `docs/spec.md`
  - `docs/architecture.md`
  - `docs/tech-selection.md`
  - `docs/coding-guidelines.md`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/stores/shell-state.ts`
  - `frontend/src/application/README.md`
  - `frontend/src/controller/wails/README.md`
  - `internal/controller/wails/app_controller.go`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `go test ./...`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test`

## 受け入れ確認

- app-shell から Foundation Data を開ける
- Persona と Dictionary を切り替えて一覧を観測できる
- 選択中の基盤エントリ詳細を確認できる
- 編集導線と Rebuild 導線を確認できる

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`
- 第2段階の UI モック artifact
- 第2段階の Scenario artifact
- 人間承認記録

## HITL 状態

- pending: design review reroute に伴う論点整理中

## 承認記録

- pending
- design review 前提未充足のため未承認

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- `docs/screen-design/` に関連 screen 正本が未作成のため、このタスクでは active exec-plan 配下の artifact を正本代替として使う
- review 用差分図が必要になった時だけ作成し、close 前に削除する
- 第2段階の task-local reading は `tasks/usecases/foundation-data-management.yaml`、active exec-plan、`docs/spec.md`、`docs/architecture.md`、現行 `frontend/src/ui/` と `internal/controller/wails/` を起点にする
- design review reroute により、UI モックは未選択初期状態、Dictionary 側 action rail、取得失敗状態、対象外導線の除去を反映してから次段階へ進める
- 人間判断が必要な論点は、Dictionary 側 action rail の文言差分、初期表示の自動選択有無、取得失敗表示の粒度である

## 結果

<!-- Outcome -->

- pending
