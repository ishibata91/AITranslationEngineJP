# 実装スコープ固定

- `task_id`: `master-dictionary-sqlite-persistence`
- `task_mode`: `implement`
- `design_review_status`: `pass`
- `hitl_status`: `approved`
- `summary`: master dictionary の production persistence path だけを対象に、current SQLite attempt の `sqlc` 依存を `sqlx` へ置き換える。bootstrap は wiring のみを担い、repository concrete と SQLite 初期化責務は `internal/repository/` または `internal/infra/` に固定する。既存の list/search/category filter/pagination/`UpdatedAt` ordering と CRUD / import semantics は維持する。`sqlc` query/codegen folder は owned scope から外し、safe に除去できる partial artifact は cleanup 対象へ含める。

## 共通ルール

- SQLite driver は現行実装の `sqlite` driver 名を維持しつつ、DB access は `sqlx` を `database/sql` の補助として使う。
- runtime DB file は repo root `db/master-dictionary.sqlite3` に固定し、migration は app startup で実行する。
- query builder は導入しない。query / command は repo-owned SQL を handwritten で保持し、`sqlx` の `Get` / `Select` / `NamedExec` / transaction helper の範囲で扱う。
- bootstrap は composition と wiring だけを担う。repository concrete、DB open、migration 実行、seed 判定、SQL 所有は `internal/repository/` または `internal/infra/` へ置き、bootstrap に残さない。
- `sqlc.yaml`、`internal/repository/sqlc/`、`tmp/disabled-repository-sqlc/` は source of truth として扱わない。`sqlx` 置換後に参照が残らないなら cleanup 対象に含める。
- `.go-arch-lint.yml` などの周辺設定は、`sqlc` artifact cleanup と矛盾しない範囲だけ更新対象に含める。
- 初回起動で空 table の場合だけ `DefaultMasterDictionarySeed(...)` を投入し、既存 DB には再 seed しない。既存 DB 起動時は migration の再適用が安全に完了し、seed が再挿入されないことを必ず検証する。
- public contract invariants として、list/search/category filter/pagination/`UpdatedAt` ordering は create、update、delete、XML import、再起動の前後で維持する。
- docs 正本更新、frontend UI 変更、master dictionary 以外の永続化導入は scope 外とする。
- `docs/tech-selection.md` は human-approved source-of-truth sync を close 前に適用する canonicalization target として扱う。今回の turn では planning artifact だけを更新する。

## 実装 handoff 一覧

### `sqlite-foundation-and-sqlx-cleanup`

- `implementation_target`: `backend`
- `owned_scope`:
  - `go.mod`
  - `.gitignore`
  - `.go-arch-lint.yml`
  - `internal/infra/sqlite/`
  - `sqlc.yaml`
  - `internal/repository/sqlc/`
  - `tmp/disabled-repository-sqlc/`
- `depends_on`: `none`
- `validation_commands`:
  - `go test ./internal/infra/...`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: `sqlx` を使う SQLite foundation が成立し、repo root `db/master-dictionary.sqlite3` を開いて migration を適用できる。既存 DB に対する migration 再適用が安全で、seed 判定を infra 側で実行できる。不要な `sqlc` artifact は safe に除去され、周辺設定に stale な `sqlc` 参照が残らない。
- `notes`:
  - `sqlx` は `database/sql` 互換の薄い拡張として使う。query builder と codegen は導入しない。
  - handwritten SQL は repository / infra 内へ閉じ込め、migration SQL と混在させない。
  - cleanup は compile / lint / test を壊さない範囲で行い、残置が必要な artifact がある場合は理由を evidence に残す。

### `master-dictionary-repository-wiring`

- `implementation_target`: `backend`
- `owned_scope`:
  - `internal/repository/master_dictionary_repository.go`
  - `internal/repository/master_dictionary_sqlite_repository.go`
  - `internal/service/master_dictionary_sqlite_repository_adapter.go`
  - `internal/bootstrap/app_controller.go`
- `depends_on`: `sqlite-foundation-and-sqlx-cleanup`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/bootstrap ./internal/usecase ./internal/controller/wails`
  - `python3 scripts/harness/run.py --suite execution`
- `completion_signal`: production wiring が `sqlx` backed repository adapter を使う構成に揃い、SQLite concrete は `internal/repository/` または `internal/infra/` にだけ残る。CRUD と XML import が同じ SQLite file に対して継続して読書きでき、既存の list/search/category filter/pagination/`UpdatedAt` ordering と CRUD / import semantics が維持される。
- `notes`:
  - bootstrap は concrete adapter を論理所有しない。constructor 呼び出しによる wiring だけに限定する。
  - query / command / import service の public contract は変えず、repository concrete の差し替えだけで成立させる。
  - migration 実行と one-time seed の責務は bootstrap から分離し、`internal/infra/sqlite/` 側へ閉じ込める。

### `persistence-proof-tests`

- `implementation_target`: `backend`
- `owned_scope`:
  - `internal/repository/master_dictionary_repository_test.go`
  - `internal/repository/master_dictionary_sqlite_repository_test.go`
  - `internal/bootstrap/app_controller_test.go`
- `depends_on`: `master-dictionary-repository-wiring`
- `validation_commands`:
  - `go test ./internal/...`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: file 生成、初回 seed、既存 DB 起動時の migration 再適用安全性、seed 非再挿入、create/update/delete、XML import、controller 再生成後の再読込、list/search/category filter/pagination/`UpdatedAt` ordering が backend test で証明される。
- `notes`:
  - persistence proof は repository wiring 完了後に着手する。永続化振る舞いの証明は wiring 前提のため、この split を先行させない。
  - UI / Wails DTO は不変なので、system test 追加より backend integration test を優先する。
  - test は実 file を使う場合でも case ごとに path を分離し、共有状態を残さない。

### `final-review-and-sync-gate`

- `implementation_target`: `review`
- `owned_scope`:
  - `docs/tech-selection.md`
  - design-review 再review記録
  - implementation-review 記録
- `depends_on`: `master-dictionary-repository-wiring`, `persistence-proof-tests`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: `sqlx` direction を反映した design-review 再reviewが通過し、implementation-review と Sonar gate が完了し、`docs/tech-selection.md` への human-approved source-of-truth sync が close 前に適用されている。
- `notes`:
  - `docs/tech-selection.md` の更新自体は human-approved docs sync として扱う。この artifact は close gate と canonicalization target を固定するためのものとする。
  - source-of-truth sync が未完了の間は task close を認めない。