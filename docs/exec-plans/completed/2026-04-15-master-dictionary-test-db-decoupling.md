# Fix Plan

- workflow: fix
- status: in_review_blocked_by_global_gate
- lane_owner: orchestrate
- scope: master-dictionary-test-db-decoupling
- task_id: 2026-04-15-master-dictionary-test-db-decoupling
- task_catalog_ref:
- parent_phase:

## Request Summary

- マスター辞書のテストが本体用 DB の初期データカテゴリに依存している。
- テスト実行時に本体用 DB 依存を切り離し、テスト用 DB へ差し替えて検証できるようにしたい。
- 対象は master dictionary 周辺の test setup、DB 初期化経路、関連 validation である。

## Decision Basis

- 既存テストの前提が production 向け初期データカテゴリに結びついているため、既存挙動に対する不適切な環境依存とみなせる。
- 要件追加ではなく既存テスト基盤の不具合是正が主目的であるため、`task_mode: fix` とする。
- 依存箇所が test setup、DB seed、master dictionary 実装の複数層にまたがる可能性があるため、`distill` と `investigate` で入口整理と原因特定を先行させる。

## Task Mode

- `task_mode`: `fix`
- `goal`: master dictionary の test が本体用 DB 初期カテゴリへ依存せず、test 用 DB で独立実行できる状態にする。
- `constraints`: `docs/` 正本は更新しない。fix workflow に従い evidence を先に揃える。最小差分で修正する。
- `close_conditions`: master dictionary 関連 test が test 用 DB で通る。`review_mode: implementation-review` が pass を返す。backend を含むため Sonar gate と open issue gate を確認する。

## Facts

- user は master dictionary test が本体用 DB の初期データカテゴリに依存していると報告している。
- user は test 用 DB に差し替えて test 可能にしたいと求めている。
- structure harness は 2026-04-15 時点で pass した。
- `NewAppController()` は `repository.DefaultMasterDictionarySeed(now())` を常時注入している。
- SQLite 初期化は DB が空の時に受け取った seed を投入する。
- repository 層の SQLite test は test DB path と test seed を直接注入できる。
- bootstrap 経由の test では DB path だけ test 用へ差し替わり、seed は production seed のままである。
- bootstrap test の一部は `Whiterun` と `地名` を前提にしている。

## Functional Requirements

- `summary`: master dictionary test の DB 依存を分離し、test fixture または test seed で完結させる。
- `in_scope`: master dictionary 関連 test、test setup、DB 初期化経路、必要最小限の実装修正。
- `non_functional_requirements`: production DB の初期カテゴリに依存しない。再現性がある。validation で確認できる。
- `out_of_scope`: docs 正本更新、master dictionary の新機能追加、無関係な DB リファクタ。
- `open_questions`: 失敗している具体的な test command は何か。fix 対象が bootstrap test だけか。seed 注入口を本体 constructor に持たせるか test 専用 builder に分けるか。
- `required_reading`: `docs/detail-specs/master-dictionary.md`、`internal/bootstrap/app_controller.go`、`internal/bootstrap/app_controller_test.go`、`internal/infra/sqlite/sqlite.go`、`internal/repository/master_dictionary_repository.go`、`internal/repository/master_dictionary_sqlite_repository_test.go`

## Artifacts

- `ui_artifact_path`:
- `final_mock_path`:
- `scenario_artifact_path`:
- `final_scenario_path`:
- `implementation_scope_artifact_path`:
- `review_diff_diagrams`:
- `source_diagram_targets`:
- `canonicalization_targets`:

## Work Brief

- `implementation_target`: master dictionary test DB decoupling
- `accepted_scope`: `NewAppController()` の既定挙動は維持したまま、bootstrap test が production seed ではなく test seed を注入できる constructor または builder を追加する。直接依存している bootstrap test を test seed 前提へ置き換える。repository 層の既存 seed 注入境界は流用し、無関係な DB リファクタは行わない。
- `parallel_task_groups`: なし
- `tasks`: distill で入口整理。investigate で再現と原因特定。implement で bootstrap seed 注入境界と関連 test を最小修正。review で implementation-review。
- `validation_commands`: `python3 scripts/harness/run.py --suite all`、対象 test command、Sonar MCP gate 確認

## Investigation

- `reproduction_status`: code-confirmed / execution-run-pass-with-coupling
- `trace_hypotheses`: test seed が production seed を共有している。master dictionary test がカテゴリ初期値を固定前提で参照している。DB 切替の injection 境界が不足している。
- `observation_points`: `internal/bootstrap/app_controller.go` の seed 注入、`internal/infra/sqlite/sqlite.go` の空 DB seed 投入、`internal/bootstrap/app_controller_test.go` の test DB helper と seed 前提 assertion
- `residual_risks`: test 専用 DB 切替が他 test に影響する可能性がある。`NewAppController()` を使う他 test も暗黙に production seed を読んでいる可能性がある。

## Acceptance Checks

- master dictionary 関連 test が production 用初期カテゴリを前提にせず成功する。
- test 実行ごとに独立した test DB または test seed を使う。
- 既存の master dictionary 振る舞いに回帰を入れない。

## Required Evidence

- 依存箇所の特定結果。
- coupling を示す最小再現コマンドと実行結果。
- test 用 DB へ切替えた実装差分。
- 関連 test と full harness の実行結果。
- Sonar gate と open issue gate の確認結果。

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: fix workflow により not-required

## Closeout Notes

- `canonicalized_artifacts`: none
- bootstrap 側に test seed を注入できる constructor 境界を追加した。
- bootstrap test は test seed helper を使うよう更新し、production seed の `Whiterun` / `地名` 前提を外した。
- `NewAppController()` の public constructor が production seed を使う既定挙動を直接監視する回帰 test を追加した。
- `go test -count=1 ./internal/bootstrap ./internal/repository` は pass した。
- `python3 scripts/harness/run.py --suite all` は frontend 既存 TypeScript 型不整合で失敗した。
- Sonar project `ishibata91_AITranslationEngineJP` では open `HIGH/BLOCKER=2`、open reliability `0`、open security `0` を確認した。

## Outcome

- scope 内の fix は実装済み。
- review は scope 内 defect なしだが、global Sonar gate 未通過のため `reroute`。
- 次手は `internal/service/master_dictionary_import_service_test.go` の既存 Sonar issue 解消、または gate 例外の human 判断。
