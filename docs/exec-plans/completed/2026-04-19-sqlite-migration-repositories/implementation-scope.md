# Implementation Scope: 2026-04-19-sqlite-migration-repositories

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: approved-by-user-close-request-2026-04-19
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`

## Source Artifacts

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `figma_artifact`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `diagramming`: `N/A`

## Fixed Decisions

- 1 job は 1 xEdit extracted data を参照する。複数入力は複数 job で扱う。
- `JOB_PHASE_RUN` は翻訳 job の phase 実行だけを表す。基盤構築は含めない。
- phase rerun は同じ `JOB_PHASE_RUN` の状態を戻す。attempt table は作らない。
- `PERSONA` と `DICTIONARY_ENTRY` は共通 / job-local を同じ table で扱う。
- `NPC_PROFILE` は抽出スナップショットをまたいだ NPC の根で、`NPC_RECORD` と `PERSONA` が参照する。
- `PERSONA` は `NPC_PROFILE` と 1:1。共通ペルソナは job-local 生成の重複スキップ用途である。
- repository 境界は `TranslationSourceRepository`、`FoundationDataRepository`、`JobLifecycleRepository`、`JobOutputRepository`、`TranslationFieldDefinitionRepository` に固定する。
- `JobOutputRepository` は DB 上の翻訳成果状態として `JOB_TRANSLATION_FIELD` を扱う。実ファイル export は repository の責務にしない。
- この plan の実装は migration と repository 作成まで。既存 service / bootstrap / UI の接続変更は後続 TODO plan に分ける。

## Handoffs

### `canonical-schema-migration-er-v1`

- `implementation_target`: `internal/infra/sqlite/migrations`, `internal/infra/sqlite`
- `owned_scope`: 統合 ER 準拠の canonical migration、schema application test、foreign key / unique / index の確認。
- `depends_on`: none
- `validation_commands`: `go test ./internal/infra/sqlite`
- `completion_signal`: fresh DB で canonical ER schema が作成され、`NPC_PROFILE` を含む主要 table と制約が検証される。
- `notes`: 旧 `master_*` schema の削除、backfill、UI / service 切替は行わない。ここでの `master_*` は既存資産名を指す。既存 migration を書き換える必要がないなら新規 migration 追加を優先する。

### `repository-contracts-er-v1`

- `implementation_target`: `internal/repository`
- `owned_scope`: ER v1 の Go type、固定 repository interface、not found / conflict error、transaction 境界の contract。
- `depends_on`: `canonical-schema-migration-er-v1`
- `validation_commands`: `go test ./internal/repository`
- `completion_signal`: `TranslationSourceRepository`、`FoundationDataRepository`、`JobLifecycleRepository`、`JobOutputRepository`、`TranslationFieldDefinitionRepository` の contract が test 可能になる。
- `notes`: table 1 件ごとの薄い repository は作らない。`TranslationSourceRepository` は翻訳入力元、`FoundationDataRepository` はペルソナ / 辞書、`JobLifecycleRepository` は job / phase state machine、`JobOutputRepository` は `JOB_TRANSLATION_FIELD`、`TranslationFieldDefinitionRepository` は lookup / seed data を扱う。旧 master 名 repository の置換や service adapter 変更は含めない。

### `sqlite-repositories-er-v1`

- `implementation_target`: `internal/repository`, `internal/infra/sqlite`
- `owned_scope`: SQLite repository 実装、temp DB integration test、transaction rollback test、reopen persistence test。
- `depends_on`: `repository-contracts-er-v1`
- `validation_commands`: `go test ./internal/infra/sqlite ./internal/repository`
- `completion_signal`: scenario `SCN-SMR-002` から `SCN-SMR-005` が repository test で通る。
- `notes`: real AI API、real xEdit、real XML file を使わない。fixture は小さく保つ。service / bootstrap / frontend wiring は含めない。phase 完了と job output 更新を atomic にしたい場合は、repository を混ぜずに上位 transaction runner で `JobLifecycleRepository` と `JobOutputRepository` を同一 transaction に入れる。

## Explicitly Out Of Scope

- 旧 master 名 dictionary / persona UI、service、gateway の canonical schema 接続。
- 旧 `master_dictionary_entries` / `master_persona_entries` の削除、backfill、互換 adapter。
- `internal/bootstrap` の production wiring 切替。
- 旧 master 名 persona の AI settings / run status の再設計。
- 実ファイル export / writer の実装。
- `TranslationArtifactRepository` の新設。

## Deferred To TODO Plan

- `legacy-master-ui-migration`: 旧 master 名 dictionary / persona UI、service、gateway を canonical schema へ寄せる。
- `legacy-operation-state`: 旧 master 名 persona の AI settings / run status を canonical ER 外の操作状態として残すか、新 schema を足すかを決める。

## Completion Packet

Copilot は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence: N/A`
- `implementation_review_result`
- `sonar_gate_result`
- `residual_risks`
- `docs_changes: none`
