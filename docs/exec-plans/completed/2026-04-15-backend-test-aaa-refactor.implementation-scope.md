# 実装スコープ固定

- `task_id`: `backend-test-aaa-refactor`
- `task_mode`: `refactor`
- `design_review_status`: `not_run`
- `hitl_status`: `approved`
- `summary`: backend 既存テストの AAA / single-intent 未準拠最小集合を 11 file に固定し、package 単位の 3 handoff で振る舞い不変 refactor を進める。

## 共通ルール

- product code は変更しない。
- 実装 scope は distill で確定した 11 file から広げない。
- 1 test は 1 intent を原則とし、Arrange / Act / Assert が body 構造から読める形へ寄せる。
- 複数観測点が別意図になる場合は test を分割する。Act は主操作を 1 つに絞る。
- helper 抽出は同一 package 内の重複削減に限る。新しい cross-package fixture や production seam は追加しない。
- validation は package 単位の targeted `go test` を通した後に `python3 scripts/harness/run.py --suite all` を通す。

## 実装 handoff 一覧

### `backend-repository-and-sqlite-test-aaa`

- `implementation_target`: `backend`
- `handoff_skill`: `tests`
- `owned_scope`:
  - `internal/repository/master_dictionary_repository_test.go`
  - `internal/repository/master_dictionary_sqlite_repository_test.go`
  - `internal/infra/sqlite/sqlite_test.go`
- `depends_on`: `none`
- `validation_commands`:
  - `go test ./internal/repository ./internal/infra/sqlite`
- `completion_signal`: repository / sqlite 系 3 file の test 名と body から AAA と single-intent を読める。CRUD 観点と error 観点が混在せず、既存 package test が pass する。
- `notes`:
  - repository helper の整理は同一 package 内へ閉じる。
  - DB 初期化や transaction 観点は成功系と異常系を分離し、assertion bundle を残さない。

### `backend-bootstrap-and-wails-controller-test-aaa`

- `implementation_target`: `backend`
- `handoff_skill`: `tests`
- `owned_scope`:
  - `internal/bootstrap/app_controller_test.go`
  - `internal/controller/wails/master_dictionary_controller_unit_test.go`
  - `internal/controller/wails/app_controller_test.go`
- `depends_on`: `backend-repository-and-sqlite-test-aaa`
- `validation_commands`:
  - `go test ./internal/bootstrap ./internal/controller/wails`
- `completion_signal`: bootstrap / Wails controller 系 3 file で controller 生成、依存注入、command 呼び出し、error 伝播の観点が test ごとに分離され、既存 package test が pass する。
- `notes`:
  - controller mock expectation と response assertion は 1 test 1 intent に分解する。
  - Wails bridge の正常系と異常系を同一 test に混在させない。

### `backend-service-and-usecase-test-aaa`

- `implementation_target`: `backend`
- `handoff_skill`: `tests`
- `owned_scope`:
  - `internal/service/master_dictionary_command_service_test.go`
  - `internal/service/master_dictionary_import_service_test.go`
  - `internal/service/master_dictionary_query_service_test.go`
  - `internal/service/master_dictionary_xml_adapter_test.go`
  - `internal/usecase/master_dictionary_usecase_test.go`
- `depends_on`: `backend-bootstrap-and-wails-controller-test-aaa`
- `validation_commands`:
  - `go test ./internal/service ./internal/usecase`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: service / usecase 系 5 file で command、import、query、XML adapter、usecase 境界の観点が AAA で読める。multi-intent test が残らず、既存 package test と full harness が pass する。
- `notes`:
  - import progress、parse、repository interaction、usecase orchestration は別意図へ分割する。
  - XML adapter の入力 variation と parse error は独立 test に寄せる。

## Canonicalization

- `backend_repository_targets`:
  - `internal/repository/master_dictionary_repository_test.go`
  - `internal/repository/master_dictionary_sqlite_repository_test.go`
  - `internal/infra/sqlite/sqlite_test.go`
- `backend_controller_targets`:
  - `internal/bootstrap/app_controller_test.go`
  - `internal/controller/wails/master_dictionary_controller_unit_test.go`
  - `internal/controller/wails/app_controller_test.go`
- `backend_service_targets`:
  - `internal/service/master_dictionary_command_service_test.go`
  - `internal/service/master_dictionary_import_service_test.go`
  - `internal/service/master_dictionary_query_service_test.go`
  - `internal/service/master_dictionary_xml_adapter_test.go`
  - `internal/usecase/master_dictionary_usecase_test.go`

## 完了シグナル

- 3 handoff の owned scope がすべて完了し、scope 拡張が発生していない。
- `go test ./internal/repository ./internal/infra/sqlite` が pass する。
- `go test ./internal/bootstrap ./internal/controller/wails` が pass する。
- `go test ./internal/service ./internal/usecase` が pass する。
- `python3 scripts/harness/run.py --suite all` が pass する。

## Open Questions

- なし
