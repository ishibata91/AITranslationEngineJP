# Requirements Design: 2026-04-19-sqlite-migration-repositories

- `skill`: requirements-design
- `status`: approved
- `source_plan`: `./plan.md`

## Capability

- `actor`: GitHub Copilot implementation lane, backend repository users, SQLite-backed services.
- `new_capability`: 統合 ER の主要テーブルを SQLite schema と repository で扱える。
- `changed_outcome`: 翻訳入力元、翻訳補助データ、ジョブ状態、ジョブ内出力状態を repository 境界で永続化できる。

## Constraints

- `business_rules`: 1 job は 1 xEdit extracted data を参照する。複数入力は複数 job として扱う。
- `scope_boundaries`: この task は migration と repository 作成まで。旧 master 名 service、bootstrap、controller、frontend UI の接続変更は扱わない。
- `invariants`: `NPC_RECORD` は `TRANSLATION_RECORD` の 1:1 派生。`NPC_PROFILE` は抽出スナップショットをまたいだ NPC の根。`PERSONA` は `NPC_PROFILE` と 1:1。訳文は `JOB_TRANSLATION_FIELD` に保持する。
- `data_ownership`: schema は `internal/infra/sqlite/migrations`、repository contract と SQLite 実装は `internal/repository` が持つ。
- `state_transitions`: job state は `JOB_PHASE_RUN` 群から集約される。phase rerun は同じ `JOB_PHASE_RUN` の状態を戻す。
- `failure_recovery`: 複数 table をまたぐ repository 書き込みは transaction で扱い、失敗時に部分永続化を残さない。

## Decision Points

### REQ-001 Migration の範囲

- `issue`: 既存 `001/002` migration は旧 master 名の辞書 / ペルソナ schema であり、統合 ER と一致しない。
- `background`: 人間判断として、最初の task は migration と repository 作成までに絞る。
- `options`: 既存 schema を削除して置換する、既存 schema と共存する canonical migration を足す、旧 UI 移行まで同時に行う。
- `recommendation`: canonical migration を作り、旧 schema の削除や UI / service 切替は行わない。
- `reasoning`: 既存資産の修正まで含めると scope が広がり、migration と repository の完成条件が曖昧になる。
- `consequences`: legacy `master_*` table と新 canonical table は一時的に共存し得る。ここでの `master_*` は既存資産名を指す。
- `open_risks`: 旧 schema の削除、backfill、UI 切替は TODO plan 側で具体化が必要。

### REQ-002 Repository の粒度

- `issue`: ER の全 table を 1 repository にまとめると肥大化し、細かく切ると 1 ユースケースの transaction 境界が読みにくくなる。
- `background`: repository は table 単位ではなく、書き込み transaction と責務のまとまりで分ける。
- `options`: ER 全体 repository、table repository、transaction 境界ごとの repository。
- `recommendation`: 初期 write repository は `TranslationSourceRepository`、`FoundationDataRepository`、`JobLifecycleRepository`、`JobOutputRepository`、`TranslationFieldDefinitionRepository` に固定する。
- `reasoning`: `TranslationSource` は原文入力構造、`FoundationData` は共通 / job-local 補助データ、`JobLifecycle` は job / phase state machine、`JobOutput` は DB 上の翻訳成果状態を扱うため、書き込み単位と責務が読みやすい。
- `consequences`: `JOB_TRANSLATION_FIELD` は `JobOutputRepository` が扱う。実ファイル export は repository ではなく exporter / writer の責務にする。
- `open_risks`: read model が必要な場合は後続で query repository を足す。`TRANSLATION_ARTIFACT` / `XTRANSLATOR_OUTPUT_ROW` の永続化要否は export 仕様側で再確認する。

### REQ-003 共通と job-local の表現

- `issue`: 共通と job-local を別 table にすると、同一概念が再分裂する。
- `background`: 共通データは job-local 生成を減らすための重複スキップ用途である。
- `options`: 別 table、同 table + lifecycle/scope/source、view 互換。
- `recommendation`: `PERSONA` と `DICTIONARY_ENTRY` は同 table + lifecycle/scope/source で表す。
- `reasoning`: 作成経路と適用範囲だけが違い、概念としては同じものを扱うため。
- `consequences`: `PERSONA.npc_profile_id` は unique にし、同一 NPC プロファイルに共通と job-local を同時保持しない。
- `open_risks`: 旧 master 名 persona 画面が要求する run status / AI settings は canonical ER とは別の UI 操作状態として再設計が必要。

### REQ-004 NPC profile のライフサイクル

- `issue`: 共通ペルソナを `NPC_RECORD` に紐づけると、入力キャッシュ flush 時に再利用単位も消える。
- `background`: `NPC_RECORD` は xEdit 抽出スナップショット上の派生レコードであり、共通ペルソナの根には向かない。
- `options`: `NPC_RECORD` に紐づける、`NPC_PROFILE` を新設する、persona 側に識別キーを重複保持する。
- `recommendation`: `NPC_PROFILE` を新設し、`NPC_RECORD` と `PERSONA` は `NPC_PROFILE` を参照する。
- `reasoning`: NPC の同一性と抽出スナップショット上の属性を分けると、ジョブ完了後の flush と共通データの残存が両立する。
- `consequences`: `TranslationSourceRepository` は翻訳入力元の保存時に `NPC_PROFILE` を upsert し、`NPC_RECORD` から参照する。
- `open_risks`: `target_plugin_name + form_id + record_type` を初期同一性キーにするが、load order や ESL compacted FormID の扱いは後続で確認する。

## Functional Requirements

- `in_scope`: ER v1 canonical migration、schema application test、repository contract、SQLite repository、repository integration test。repository contract は `TranslationSourceRepository`、`FoundationDataRepository`、`JobLifecycleRepository`、`JobOutputRepository`、`TranslationFieldDefinitionRepository` を対象にする。
- `non_functional_requirements`: deterministic tests、temp DB、real AI API 不使用、foreign key 有効化、transaction rollback の検証。
- `out_of_scope`: 旧 master 名 UI の切替、legacy data backfill、旧 schema / repository の削除、service / bootstrap wiring、AI 実行、docs 正本化、xEdit / XML parser の仕様変更、実ファイル export 実装、`TranslationArtifactRepository` の新設。
- `acceptance_basis`: fresh DB で canonical schema が作られる。固定した repository 境界ごとに保存して再読込できる。unique / foreign key / transaction rollback が test で観測できる。

## Open Questions

- 旧 master 名 UI 移行 task で、run status / AI settings を canonical ER 外の操作状態として残すか、新しい UI 状態 schema を足すか。

## Required Reading

- `docs/er.md`
- `docs/diagrams/er/combined-data-model-er.d2`
- `internal/infra/sqlite/sqlite.go`
- `internal/infra/sqlite/migrations/001_master_dictionary_entries.sql`
- `internal/infra/sqlite/migrations/002_master_persona_tables.sql`
- `internal/repository/master_dictionary_repository.go`
- `internal/repository/master_persona_repository.go`
