# Work Plan

- workflow: work
- status: completed
- lane_owner: orchestrate
- scope: master-dictionary-sqlite-persistence
- task_id: master-dictionary-sqlite-persistence
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml
- parent_phase: live-workflow

## Request Summary

- マスター辞書登録を `sqlx` による SQLite 永続化へ切り替える。
- SQLite の DB ファイルは repo root `db/master-dictionary.sqlite3` へ固定する。
- `sqlc` と query builder は採用せず、既存の一覧契約と CRUD / import semantics を維持する。

## Decision Basis

- user request は既存の master dictionary 永続化試行を `sqlc` から `sqlx` へ切り替える direction change である。
- `docs/spec.md` はマスター辞書を事前に構築し、UI から観測可能であることを要求している。
- `docs/architecture.md` は backend の永続化を `Repository` adapter に閉じ込めること、concrete wiring を `internal/bootstrap/` で行うことを要求している。
- 既存 completed plan `2026-04-11-master-dictionary-management.md` は list/search/category filter/pagination/UpdatedAt ordering と CRUD / import 導線を public contract として定義済みである。
- `docs/tech-selection.md` は現時点では `sqlc` 前提のままであり、approved direction と source-of-truth の同期が close 前に必要である。
- 今回の変更は非自明であり、active plan と implementation-scope artifact の再固定、および design-review の再実施が必要である。

## Task Mode

- `task_mode`: implement
- `goal`: master dictionary registration と参照更新系を `sqlx` + SQLite 永続化へ移し、migration を起動時に適用し、DB ファイル保存先を repo root `db/master-dictionary.sqlite3` に固定する。
- `constraints`:
  - orchestrate 自身は実装しない。
  - `docs/` 正本は human 先行変更がない限り更新しない。
  - library と D2 の書き方は必要時に Context7 で確認する。
  - backend を含む task として Sonar gate と implementation review pass を close 条件に含める。
  - bootstrap は wiring のみを担い、repository concrete や DB 初期化責務を論理所有しない。
  - `sqlx` は `database/sql` 拡張として使い、query builder は導入しない。
  - approved direction が `sqlc` から `sqlx` へ変わったため、design-review の再reviewが終わるまで実装 close は block する。
- `close_conditions`:
  - sqlx direction を反映した design-review が pass を返すこと。
  - implementation-review が pass を返すこと。
  - `docs/tech-selection.md` への human-approved source-of-truth sync が close 前に完了すること。
  - backend 実装と review で Sonar 件数ゲートを確認すること。
  - `HIGH` / `BLOCKER` の open issue が 0 件であること。
  - open reliability issue が 0 件であること。
  - open security issue が 0 件であること。
  - `python3 scripts/harness/run.py --suite all` が通ること。

## Facts

- 現在の production wiring は `internal/bootstrap/app_controller.go` で SQLite-backed adapter を組み立て、repo root `db/master-dictionary.sqlite3` を使う前提の test が `internal/bootstrap/app_controller_test.go` に存在する。
- 現在の SQLite repository 試行は `internal/repository/master_dictionary_sqlite_repository.go` と `internal/infra/sqlite/sqlite.go` で generated `sqlc` package を import している。
- repo root には `sqlc.yaml` があり、`internal/repository/sqlc/` と `tmp/disabled-repository-sqlc/` が failed attempt の artifact として残っている。
- `.go-arch-lint.yml` には `repository_sqlc` ルールが残っており、`sqlc` cleanup と同時に整合を取り直す必要がある。
- migration は app 起動時に再適用安全である必要があり、空 DB に対する seed は 1 回だけでなければならない。
- 既存の public contract として list/search/category filter/pagination/`UpdatedAt` ordering、および CRUD / import semantics が維持対象である。
- structure harness は 2026-04-15 時点で pass している。

## Functional Requirements

- `summary`:
  - マスター辞書の登録処理を `sqlx` による SQLite 永続化へ接続する。
  - query builder は使わず、repo-owned SQL と `sqlx` helper で query / command を実装する。
  - SQLite DB は repo root `db/master-dictionary.sqlite3` に保存する。
  - 既存の list/search/category filter/pagination/`UpdatedAt` ordering と CRUD / import semantics を維持する。
  - safe に除去できる `sqlc` artifact は owned scope に含めて cleanup する。
- `in_scope`:
  - master dictionary 永続化の現状整理。
  - `sqlx` 導入と既存 SQLite foundation の再固定。
  - startup migration と one-time seed の責務分離。
  - master dictionary repository / adapter / bootstrap wiring の成立。
  - repo root `db/master-dictionary.sqlite3` path 制御。
  - `sqlc.yaml`、`internal/repository/sqlc/`、`tmp/disabled-repository-sqlc/`、関連 lint 設定の safe cleanup。
  - 一覧の list/search/category filter/pagination/`UpdatedAt` ordering と CRUD / import semantics の維持確認。
  - 既存導線に対する必要な tests / validation の追加更新。
  - close 前の `docs/tech-selection.md` source-of-truth sync を canonicalization target として固定すること。
- `non_functional_requirements`:
  - 既存のアーキテクチャ依存方向を破らない。
  - bootstrap は wiring だけを担い、repository concrete、migration、seed、SQL 所有を持たない。
  - `sqlx` は `database/sql` の補助に限定し、query builder や codegen を持ち込まない。
  - 永続化層は再現可能で、ローカル実行時に repo root `db/` 配下を利用する。
  - 既存 DB 起動時の migration 再適用安全性と seed 非再挿入を明示的に検証する。
- `out_of_scope`:
  - master dictionary 以外の基盤データ永続化。
  - human 承認なしの `docs/` 正本更新。
  - unrelated frontend redesign。
  - query builder や ORM の導入。
- `open_questions`:
  - なし。
- `required_reading`:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-master-dictionary-management.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml

## Artifacts

- `ui_artifact_path`:
- `final_mock_path`:
- `scenario_artifact_path`:
- `final_scenario_path`:
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/master-dictionary-sqlite-persistence.implementation-scope.md
- `review_diff_diagrams`:
- `source_diagram_targets`: なし
- `canonicalization_targets`:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md

## Work Brief

- `implementation_target`: master dictionary backend persistence path
- `accepted_scope`:
  - distill で facts / constraints / related_code_pointers を抽出する。
  - design で `sqlx` direction を反映した implementation-scope を再固定する。
  - review で design-review を再実施し、pass 前は close を block する。
  - implement は `sqlx` ベースの master dictionary 永続化と repo root `db/master-dictionary.sqlite3` 保存先、startup migration、safe な `sqlc` artifact cleanup に限定する。
  - bootstrap は wiring のみを担い、repository concrete と DB 初期化は `internal/repository/` または `internal/infra/` に置く。
  - tests は migration 再適用安全性、one-time seed、list/search/category filter/pagination/`UpdatedAt` ordering、不変な CRUD / import semantics を証明する。
  - close 前に `docs/tech-selection.md` へ human-approved source-of-truth sync を適用する。
  - review は implementation-review を実施する。
- `parallel_task_groups`:
  - `group_id`: `distill-entry`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: なし
    - `completion_signal`: 現行コード上の関連境界、`sqlx` 導入点、残存 `sqlc` artifact、必要 validation が明確になる。
  - `group_id`: `design-scope-freeze`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `distill-entry`
    - `completion_signal`: 実装対象ファイル、cleanup 対象、close gate が narrow scope で固定される。
  - `group_id`: `sqlite-foundation-and-sqlx-cleanup`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `design-scope-freeze`
    - `completion_signal`: `sqlx` を使う SQLite foundation が整い、startup migration と one-time seed の境界が成立し、不要な `sqlc` artifact を safe に除去できる。
  - `group_id`: `master-dictionary-repository-wiring`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `sqlite-foundation-and-sqlx-cleanup`
    - `completion_signal`: bootstrap wiring と repository adapter 差し替えで永続化が成立し、一覧契約と CRUD / import semantics が維持される。
  - `group_id`: `persistence-proof-tests`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `master-dictionary-repository-wiring`
    - `completion_signal`: 永続化成立、一覧契約維持、migration / seed 不変条件が test で証明される。
  - `group_id`: `final-review-and-sync-gate`
    - `can_run_in_parallel_with`: なし
    - `blocked_by`: `master-dictionary-repository-wiring`, `persistence-proof-tests`
    - `completion_signal`: design-review 再review、implementation-review、human-approved docs sync、close gate 確認が完了する。
- `tasks`:
  - `task_id`: `distill-entry`
    - `owned_scope`: active plan と既存 docs から facts / constraints / related code を抽出し、次 skill へ渡す。
    - `depends_on`: なし
    - `parallel_group`: `distill-entry`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
  - `task_id`: `design-scope-freeze`
    - `owned_scope`: `sqlx` + SQLite 永続化、`sqlc` cleanup、close 前 docs sync を含む narrow implementation scope と required validation を固定する。
    - `depends_on`: `distill-entry`
    - `parallel_group`: `design-scope-freeze`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite structure`
  - `task_id`: `sqlite-foundation-and-sqlx-cleanup`
    - `owned_scope`: `go.mod`、`.gitignore`、`.go-arch-lint.yml`、`internal/infra/sqlite/`、`sqlc.yaml`、`internal/repository/sqlc/`、`tmp/disabled-repository-sqlc/` を対象に、`sqlx` foundation への差し替えと safe cleanup を行う。
    - `depends_on`: `design-scope-freeze`
    - `parallel_group`: `sqlite-foundation-and-sqlx-cleanup`
    - `validation_commands`:
      - `go test ./internal/infra/...`
      - `python3 scripts/harness/run.py --suite structure`
    - `notes`: `sqlx` は `Get` / `Select` / `NamedExec` / transaction helper の範囲で使い、query builder は導入しない。
  - `task_id`: `master-dictionary-repository-wiring`
    - `owned_scope`: `internal/repository/` 配下の SQLite repository concrete、`internal/service/master_dictionary_sqlite_repository_adapter.go`、`internal/bootstrap/app_controller.go` の wiring、repo root `db/master-dictionary.sqlite3` path 制御、一覧契約を壊さない master dictionary backend 永続化実装。
    - `depends_on`: `sqlite-foundation-and-sqlx-cleanup`
    - `parallel_group`: `master-dictionary-repository-wiring`
    - `validation_commands`:
      - `go test ./internal/repository ./internal/service ./internal/bootstrap ./internal/usecase ./internal/controller/wails`
      - `python3 scripts/harness/run.py --suite execution`
  - `task_id`: `persistence-proof-tests`
    - `owned_scope`: empty DB 初回 seed、既存 DB 起動時の migration 再適用安全性、seed 非再挿入、CRUD / import 永続化、list/search/category filter/pagination/`UpdatedAt` ordering、controller 再生成後の再読込を証明する backend tests と validation 更新。
    - `depends_on`: `master-dictionary-repository-wiring`
    - `parallel_group`: `persistence-proof-tests`
    - `validation_commands`:
      - `go test ./internal/...`
      - `python3 scripts/harness/run.py --suite all`
  - `task_id`: `final-review-and-sync-gate`
    - `owned_scope`: design-review 再reviewの通過、implementation-review、Sonar gate、`docs/tech-selection.md` の human-approved source-of-truth sync、full harness、残課題整理。
    - `depends_on`: `master-dictionary-repository-wiring`, `persistence-proof-tests`
    - `parallel_group`: `final-review-and-sync-gate`
    - `validation_commands`:
      - `python3 scripts/harness/run.py --suite all`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite execution`
  - `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: not_applicable
- `trace_hypotheses`:
  - `sqlc` 前提の current SQLite attempt を `sqlx` へ差し替える際に、repository / infra / lint config の cleanup が同時に必要になる。
  - startup migration と one-time seed の責務が bootstrap へ漏れると architecture violation になりやすい。
- `observation_points`:
  - master dictionary controller / usecase / service / repository 境界
  - `internal/infra/sqlite/` の DB open / migration / seed 責務
  - 残存 `sqlc` artifact と関連 lint 設定
  - 既存一覧の list/search/category filter/pagination/`UpdatedAt` ordering
- `residual_risks`:
  - docs 正本が `sqlc` 前提のまま close されると source-of-truth と実装が乖離する。
  - migration / seed の責務分離が不十分だと、再起動時の seed 重複や bootstrap 過責務が再発しやすい。

## Acceptance Checks

- マスター辞書登録後に SQLite へ保存される。
- アプリ再起動後もマスター辞書データが保持される。
- DB ファイルが repo root `db/master-dictionary.sqlite3` に生成される。
- 一覧の list/search/category filter/pagination/`UpdatedAt` ordering が永続化後も維持される。
- CRUD / import semantics が `sqlx` 差し替え後も維持される。
- 既存 DB での起動時に migration 再適用で失敗せず、seed が再挿入されない。
- `sqlc` 依存 artifact が safe に cleanup されるか、残す場合は明示的な理由が evidence に残る。

## Required Evidence

- distill の facts / constraints / related_code_pointers
- design の implementation scope artifact
- `sqlx` direction に対する design-review 再review結果
- 既存 DB 起動時の migration 再適用安全性と seed 非再挿入を示す test / validation 結果
- implementation-review の結果
- `docs/tech-selection.md` source-of-truth sync 記録
- Sonar gate 確認結果
- `python3 scripts/harness/run.py --suite all` の結果

## HITL Status

- `functional_or_design_hitl`: approved
- `approval_record`:
  - 2026-04-15 orchestrate: active plan を作成し、distill へ handoff 準備。
  - 2026-04-15 design-review: reroute。bootstrap ownership、一覧契約不変条件、既存 DB 起動時 validation、dependency alignment の補強が必要と判断された。
  - 2026-04-15 design: reroute 指摘を active plan と implementation-scope artifact へ反映した。実装は次の design-review 結果まで継続して block する。
  - 2026-04-15 human direction change: approved direction を `sqlc` から `sqlx` へ変更し、query builder を使わない方針を確定した。design-review はこの変更を前提に再reviewが必要である。
  - 2026-04-15 design-review: sqlx direction と cleanup scope を確認し、pass。
  - 2026-04-15 implementation-review: pass。指摘なし。

## Closeout Notes

- `canonicalized_artifacts`:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md
- `implementation_summary`:
  - `sqlc` 依存を `sqlx` に置き換え、startup migration と one-time seed を SQLite 初期化へ閉じ込めた。
  - master dictionary repository / adapter / wiring を `sqlx` 前提へ差し替え、repo root `db/master-dictionary.sqlite3` 永続化を成立させた。
  - stale な `sqlc` dependency rule を cleanup し、失敗試行 artifact は archive へ移動した。
- `validation_summary`:
  - `go test ./internal/...`: pass
  - `python3 scripts/harness/run.py --suite all`: pass
  - Sonar open `HIGH` / `BLOCKER`: 0
  - Sonar open reliability: 0
  - Sonar open security: 0

## Outcome

- completed