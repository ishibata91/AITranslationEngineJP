# Implementation Scope: term-target-rec-config

- `skill`: implementation-scope
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `human_review_status`: 承認済み詳細仕様差分と承認済み設計差分図を入力に固定する。本成果物自体は人間レビュー待ち。
- `approval_record`: 詳細仕様差分 `./detail-spec-diff.md`（status: ready-for-human-review、回答済み、未決 0 件）、設計差分図 `./design-diff.md`、人間レビュー記録（2026-05-31 確定事項）
- `module_entry`: `.claude/skills/implementation-module/SKILL.md`
- `handoff_runtime`: `claude-module`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`: `N/A`（画面変更なし、plan.md の想定 Y/N で N 確定）
- `design_diff`: `./design-diff.md`

## Fixed Decisions

- 単語翻訳対象 REC と XML 辞書取り込み対象 REC は同一の 13 種別集合とする。
- 共通判定関数は `IsTermTarget(rec string) bool` 1 つだけを公開し、`XMLImportAllowed` のような別関数や別集合は作らない。
- 共通 config の置き場所は新規 package `internal/recclassification` とする。package には 13 種別集合定数、`IsTermTarget`、SQL 埋め込み用 REC リスト生成手段を含める。
- `dictionary_scope` の保存形式と比較形式を `RECORD:FIELD` 形式（例 `NPC_:FULL`）に変える。既存翻訳データは reset 前提とし、互換 migration は作らない。
- SQL の REC リストは Go 側共通 config から生成して埋め込む。SQL 直書きの REC リテラルを作らない。
- frontend、UI、Wails bridge、DTO は変更しない。統合境界 handoff は不要。
- `unanswered_questions`: `0`
- secret 取り扱いなし。`secret_boundary.status = not_required` を全 handoff に書く。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `B1-recclassification-config` | `なし` | `なし` | `なし` |
| `wave-2` | `B2-term-phase-candidate-filter`, `B3-processing-target-sql-filter`, `B4-master-dictionary-import-gate` | `B1-recclassification-config` | `B2 <-> B3`, `B2 <-> B4`, `B3 <-> B4` | `なし` |
| `wave-3` | `UT1-recclassification`, `UT2-term-phase-candidate`, `UT3-processing-target-sql`, `UT4-master-dictionary-import`, `ST1-processing-target-and-execution-rec-parity` | `B1`, `B2`, `B3`, `B4` 全完了 | `UT1 <-> UT2 <-> UT3 <-> UT4`、`ST1` は並列可 | `なし` |

注: wave-2 は変更ファイルが完全に分離する（`internal/service/term_translation_phase_service.go` / `internal/repository/processing_target_sqlite_repository.go` / `internal/service/master_dictionary_service.go`）。共有変更は wave-1 の共通 config に集約しており、wave-2 では package 公開 API（`IsTermTarget` と REC リスト生成）の呼び出しだけ行う。

## Handoffs

### `B1-recclassification-config`:

- `implementation_target`: 単語翻訳対象 REC 共通判定 package を新規追加する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-008`、`master-dictionary-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規ファイル: `internal/recclassification/term_target.go`
  - 公開: `TermTargets`（13 種別集合、定義順は `BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL`）
  - 公開: `IsTermTarget(rec string) bool`
  - 公開: SQL 埋め込み用 REC リスト取得手段（例 `TermTargetRECList() []string`）。SQL 側 IN 句の二重管理を防ぐため共通 config が単一情報源になる。
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/recclassification/term_target.go`
  - symbol: `TermTargets` と `IsTermTarget`
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「`IsTermTarget` が 13 種別だけに true を返し、それ以外で false を返す」
  - 1 手目にする理由: wave-2 の 3 引き継ぎがすべてこの公開 API に依存するため、最小公開面を最初に閉じる。
- `validation_commands`:
  - `go build ./internal/recclassification/...`
  - `go vet ./internal/recclassification/...`
- `completion_signal`:
  - 13 種別集合が定数として package 内に定義されている。
  - `IsTermTarget` が 13 種別だけに true を返し、それ以外で false を返す。
  - SQL 埋め込み用 REC リスト取得手段が、SQL からの二重管理を防ぐ単一情報源として公開されている。
  - 13 種別集合に `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` を含まない。
  - `NPC_:FULL` と `NPC_:SHRT` を別 REC として両方含む。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff は config 単独で、外部 IO や状態を持たない。単体検証だけで十分な担当者にする。

### `B2-term-phase-candidate-filter`:

- `implementation_target`: `collectCandidates` の REC 絞り込みと candidate key を `RECORD:FIELD` 形式へ変える。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-003`、`term-translation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 変更ファイル: `internal/service/term_translation_phase_service.go`
  - 変更対象: `collectCandidates`（`internal/service/term_translation_phase_service.go:1667` 付近）
  - candidate key を `record.RecordType + ":" + field.SubrecordType` で組み立てる `RECORD:FIELD` 形式に変える。
  - 13 種別判定は `recclassification.IsTermTarget` を呼ぶ。
  - 絞り込みに該当しない候補は集合外として除外する。
  - `confirmedTerms` の key 生成、確定訳語の保存単位も `RECORD:FIELD` 形式の REC 単位に合わせる。
- `depends_on`: `B1-recclassification-config`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `B3-processing-target-sql-filter`, `B4-master-dictionary-import-gate`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/term_translation_phase_service.go`
  - symbol: `collectCandidates`
  - 変更種別: 既存関数の変更
  - 対応する `completion_signal` clause: 「`collectCandidates` は 13 種別に該当する候補だけを返す」
  - 1 手目にする理由: 単語翻訳フェーズ候補集合の責務が `collectCandidates` 1 箇所に集約しており、最初にここを閉じれば後続の confirmedTerms key 整合は同一関数内の最小 closure chain で閉じられる。
- `validation_commands`:
  - `go build ./internal/service/...`
  - `go vet ./internal/service/...`
- `completion_signal`:
  - `collectCandidates` は 13 種別に該当する候補だけを返す。
  - candidate key と confirmedTerms key が `RECORD:FIELD` 形式の REC 単位になる。
  - `NPC_:FULL` と `NPC_:SHRT` を別候補として保持できる。
  - 共通辞書一致による除外 (`confirmedTerms` 走査) は変更後も成立する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 同一ファイル内の他関数（実行ループ、結果反映）に key 形式変更の波及がある場合、最小 closure chain として同一 handoff 内で閉じる。
  - 公開接点: `collectCandidates` の返り値型、confirmedTerms key 表現、確定訳語の `dictionary_scope` 値。downstream の `B3` と単体テストはこの公開接点を前提にする。

### `B3-processing-target-sql-filter`:

- `implementation_target`: 処理対象一覧 SQL 2 種に REC 絞り込みを追加し、`dictionary_scope` 比較を `RECORD:FIELD` 形式へ変える。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-003`、`term-translation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 変更ファイル: `internal/repository/processing_target_sqlite_repository.go`
  - 変更対象: `processingTargetTermCountSQL`（`:247` 付近）、`processingTargetTermListSQL`（`:282` 付近）
  - 共通 config の REC リスト取得手段（`B1` 提供）から `IN (...)` 句を生成し、SQL に埋め込む。
  - candidate 行生成のキー（`tr.record_type` ベース）を `RECORD:FIELD` 形式の REC として組み立てる。
  - `dictionary_scope` 比較を `RECORD:FIELD` 形式 REC 値どうしの比較に揃える。
- `depends_on`: `B1-recclassification-config`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `B2-term-phase-candidate-filter`, `B4-master-dictionary-import-gate`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/repository/processing_target_sqlite_repository.go`
  - symbol: `processingTargetTermCountSQL` と `processingTargetTermListSQL`
  - 変更種別: 既存 SQL 定数の変更
  - 対応する `completion_signal` clause: 「2 SQL とも 13 種別に該当する candidate だけを返す」
  - 1 手目にする理由: 同じ責務（処理対象一覧の絞り込み）を持つ 2 SQL を同一 handoff で同時に変更しないと、件数と一覧の不整合が発生する。共通 config 由来の REC リストを SQL に差し込む構造を最初に確立する。
- `validation_commands`:
  - `go build ./internal/repository/...`
  - `go vet ./internal/repository/...`
- `completion_signal`:
  - 2 SQL とも 13 種別に該当する candidate だけを返す。
  - candidate 行のキー、`dictionary_scope` 比較が `RECORD:FIELD` 形式 REC で行われる。
  - `NPC_:FULL` と `NPC_:SHRT` を別行として返せる。
  - SQL 中に REC リテラルの直書きが残らず、共通 config からの埋め込みに統一される。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 公開接点: SQL が返す candidate キー形式、`dictionary_scope` 値の比較形式。`ST1` の整合検証はこの公開接点を前提にする。

### `B4-master-dictionary-import-gate`:

- `implementation_target`: XML 辞書取り込みの REC 許可判定を共通 config の `IsTermTarget` 呼び出しへ統合する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`（`master-dictionary-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 変更ファイル: `internal/service/master_dictionary_service.go`
  - 削除: `allowedImportREC` 固定 map（`:28` 付近）
  - 変更: `isAllowedImportREC`（`:209` 付近）を `recclassification.IsTermTarget` の呼び出しに置き換える
  - 変更: `categoryFromREC`（`:214` 付近）から `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` の分岐を削除する
- `depends_on`: `B1-recclassification-config`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `B2-term-phase-candidate-filter`, `B3-processing-target-sql-filter`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/master_dictionary_service.go`
  - symbol: `isAllowedImportREC`
  - 変更種別: 既存関数の本体置換
  - 対応する `completion_signal` clause: 「`isAllowedImportREC` が `recclassification.IsTermTarget` の戻り値を返す」
  - 1 手目にする理由: 公開接点（許可判定）の置換を最初に閉じれば、`allowedImportREC` 削除と `categoryFromREC` 整理は同一ファイル内の最小 closure chain で閉じられる。
- `validation_commands`:
  - `go build ./internal/service/...`
  - `go vet ./internal/service/...`
- `completion_signal`:
  - `isAllowedImportREC` が共通 config の `IsTermTarget` を呼ぶ実装になる。
  - `allowedImportREC` 固定 map がファイル内から削除されている。
  - `categoryFromREC` から `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` 分岐が削除されている。
  - 13 種別以外の REC に対し XML 取り込みが許可されない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 公開接点: `isAllowedImportREC` の戻り値、`categoryFromREC` の分岐集合。

### `UT1-recclassification`:

- `implementation_target`: 共通 config `IsTermTarget` の 13 種別判定単体テスト。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-008`、`master-dictionary-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規: `internal/recclassification/term_target_test.go`
  - 13 種別すべてが true を返すケース
  - `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL`、空文字、未知 REC が false を返すケース
  - `NPC_:FULL` と `NPC_:SHRT` が独立判定されること
  - REC リスト取得手段が 13 種別を重複なく返すこと
- `depends_on`: `B1-recclassification-config`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `UT2-term-phase-candidate`, `UT3-processing-target-sql`, `UT4-master-dictionary-import`, `ST1-processing-target-and-execution-rec-parity`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/recclassification/term_target_test.go`
  - symbol: `TestIsTermTarget_includes13Types`
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「13 種別 true / 13 種別外 false のテーブル駆動テストが green になる」
  - 1 手目にする理由: 後続 UT/ST の前提（共通 config の正しさ）を最小単体で固定する。
- `validation_commands`:
  - `go test ./internal/recclassification/...`
- `completion_signal`:
  - 13 種別 true / 13 種別外 false のテーブル駆動テストが green になる。
  - `NPC_:FULL` と `NPC_:SHRT` 別判定テストが green になる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`: 単体テスト。期待結果の元ネタは `spec_basis` の 13 種別仕様。

### `UT2-term-phase-candidate`:

- `implementation_target`: `collectCandidates` の REC 絞り込みと candidate key 形式の単体テスト。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-003`、`term-translation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規または既存追記: `internal/service/term_translation_phase_service_test.go` 系
  - `collectCandidates` が 13 種別のみを候補化することの検証
  - 13 種別外の REC を含む入力で対応行が候補から除外されることの検証
  - candidate key と confirmedTerms key が `RECORD:FIELD` 形式であることの検証
  - `NPC_:FULL` と `NPC_:SHRT` 同一原語が別候補として保持されることの検証
- `depends_on`: `B2-term-phase-candidate-filter`, `B1-recclassification-config`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `UT1-recclassification`, `UT3-processing-target-sql`, `UT4-master-dictionary-import`, `ST1-processing-target-and-execution-rec-parity`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/term_translation_phase_service_test.go`（既存ファイルへの追記または新規）
  - symbol: `TestCollectCandidates_filtersBy13RECTypes`
  - 変更種別: テスト関数追加
  - 対応する `completion_signal` clause: 「13 種別だけが候補に残るテストが green になる」
  - 1 手目にする理由: `B2` の中心責務（絞り込み）を最初に固定する。
- `validation_commands`:
  - `go test ./internal/service/... -run TermTranslationPhase`
- `completion_signal`:
  - 13 種別だけが候補に残るテストが green になる。
  - candidate key が `RECORD:FIELD` 形式である。
  - `NPC_:FULL` と `NPC_:SHRT` が別候補として保持される。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`: 期待結果の元ネタは `spec_basis` の REC 集合定義。

### `UT3-processing-target-sql`:

- `implementation_target`: 処理対象一覧 SQL の REC 絞り込み統合検証（in-memory SQLite を用いた repository 単体テスト）。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-003`、`term-translation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規または既存追記: `internal/repository/processing_target_sqlite_repository_test.go` 系
  - `processingTargetTermCountSQL` と `processingTargetTermListSQL` が 13 種別のみ集計／一覧化することの検証
  - 13 種別外 REC を含むデータでも対応行が結果に出ないことの検証
  - `dictionary_scope` を `RECORD:FIELD` 形式で投入したデータに対し共通辞書一致除外が成立することの検証
  - `NPC_:FULL` と `NPC_:SHRT` が別 candidate として返ることの検証
- `depends_on`: `B3-processing-target-sql-filter`, `B1-recclassification-config`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `UT1-recclassification`, `UT2-term-phase-candidate`, `UT4-master-dictionary-import`, `ST1-processing-target-and-execution-rec-parity`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/repository/processing_target_sqlite_repository_test.go`
  - symbol: `TestProcessingTargetTermSQL_filtersBy13RECTypes`
  - 変更種別: テスト関数追加
  - 対応する `completion_signal` clause: 「count SQL と list SQL の両方が 13 種別だけを返すテストが green になる」
  - 1 手目にする理由: `B3` の中心責務（SQL 絞り込み）を 2 SQL 同時に固定する。
- `validation_commands`:
  - `go test ./internal/repository/... -run ProcessingTargetTerm`
- `completion_signal`:
  - count SQL と list SQL の両方が 13 種別だけを返すテストが green になる。
  - `dictionary_scope` の `RECORD:FIELD` 形式比較で共通辞書一致除外が成立する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`: 期待結果の元ネタは `spec_basis` の 13 種別仕様。

### `UT4-master-dictionary-import`:

- `implementation_target`: XML 辞書取り込み許可判定の単体テスト。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`（`master-dictionary-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規または既存追記: `internal/service/master_dictionary_service_test.go` 系
  - `isAllowedImportREC` が 13 種別だけ true を返すことの検証
  - `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` が false を返すことの検証
  - `categoryFromREC` から DOOR/FLOR/FURN 分岐が消えていることの検証
- `depends_on`: `B4-master-dictionary-import-gate`, `B1-recclassification-config`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `UT1-recclassification`, `UT2-term-phase-candidate`, `UT3-processing-target-sql`, `ST1-processing-target-and-execution-rec-parity`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/service/master_dictionary_service_test.go`
  - symbol: `TestIsAllowedImportREC_only13Types`
  - 変更種別: テスト関数追加
  - 対応する `completion_signal` clause: 「13 種別 true / DOOR・FLOR・FURN false のテストが green になる」
  - 1 手目にする理由: `B4` の中心責務（XML 許可判定）を最初に固定する。
- `validation_commands`:
  - `go test ./internal/service/... -run MasterDictionary`
- `completion_signal`:
  - 13 種別 true / DOOR・FLOR・FURN false のテストが green になる。
  - `categoryFromREC` から DOOR/FLOR/FURN 分岐が消えていることのテストが green になる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`: 期待結果の元ネタは `spec_basis`。

### `ST1-processing-target-and-execution-rec-parity`:

- `implementation_target`: 処理対象一覧（repository）と単語翻訳フェーズ実行対象（service）の REC 集合一致をシナリオテストで証明する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`（`term-translation-phase-REQ-003`、`term-translation-phase-REQ-008`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - 新規または既存追記: 単語翻訳フェーズのシナリオテスト系（`internal/service/...` または `tests/scenario/...` の既存配置に従う）
  - 同一 xEdit 抽出データに対し、repository が返す処理対象一覧 REC 集合と service `collectCandidates` が返す候補 REC 集合が一致することを検証する。
  - 13 種別外 REC を含む入力で、両者とも対応行が存在しないことを検証する。
- `depends_on`: `B2-term-phase-candidate-filter`, `B3-processing-target-sql-filter`, `B1-recclassification-config`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `UT1-recclassification`, `UT2-term-phase-candidate`, `UT3-processing-target-sql`, `UT4-master-dictionary-import`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: 単語翻訳フェーズのシナリオテスト配置（既存配置に従い、Go test ファイルとして追加）
  - symbol: `TestScenario_ProcessingTargetAndExecutionRECParity`
  - 変更種別: テスト関数追加
  - 対応する `completion_signal` clause: 「処理対象一覧と実行対象の REC 集合一致テストが green になる」
  - 1 手目にする理由: 両者の整合という task の中心観点を 1 シナリオで閉じる。
- `validation_commands`:
  - `go test ./internal/service/... -run Scenario_ProcessingTargetAndExecutionRECParity`
- `completion_signal`:
  - 処理対象一覧と実行対象の REC 集合一致テストが green になる。
  - 13 種別外 REC が両者の対象外であることが green で確認される。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `notes`:
  - 本シナリオは backend 内部の公開接点（repository SQL の返り値と service `collectCandidates` の返り値）を起点にする APIテスト相当として扱う。UI 入口は本 task に存在しない（画面変更なし）。
  - UI 人間操作 E2E は対象外（画面変更なし）。

## Completion Packet

実装モジュールは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `N/A`（画面変更なし）
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: モジュール終了判断で読む実装事実。completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`（docs 正本化は finalization-module で扱う）
