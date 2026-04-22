# Requirements Design: 2026-04-19-legacy-schema-ui-migration-todo

- `skill`: requirements-design
- `status`: draft
- `source_plan`: `./plan.md`

## Capability

- `actor`: GitHub Copilot implementation lane, existing master dictionary / master persona UI users, backend service maintainers.
- `new_capability`: 既存の master dictionary / master persona 導線を canonical schema / repository 上で動かせる。
- `changed_outcome`: 旧 `master_*` table と旧 repository を新規 read / write の中心から外し、canonical `DICTIONARY_ENTRY`、`PERSONA`、`NPC_PROFILE` 系へ接続する。

## Constraints

- `business_rules`: 共通 dictionary / persona は job をまたいで残る foundation data として扱う。job-local data と混ぜる場合も lifecycle / scope / source で区別する。
- `scope_boundaries`: この design start では product code を変更しない。layout redesign は含めず、既存画面の主要導線維持を前提にする。ただし UI 改修は必要であり、canonical schema で出せない表示項目と frontend / Wails contract field は削除対象にする。
- `invariants`: `PERSONA` は `NPC_PROFILE` と 1:1。dictionary は `source_term` / `translated_term` を中核にし、重複判定は `trim(source_term) + translated_term` の完全一致に固定する。AI settings は canonical foundation data 本体に混ぜず、`PERSONA_GENERATION_SETTINGS` に分離する。run state は DB に保存せず、画面メモリだけで扱う。
- `data_ownership`: canonical foundation data は `internal/repository.FoundationDataRepository` と `internal/infra/sqlite` が持つ。画面用 read / command model は service / usecase boundary で変換する。
- `state_transitions`: dictionary CRUD / XML import、persona preview / execute / update / delete の既存画面状態を維持する。
- `failure_recovery`: legacy data backfill は行わない。旧 table の既存 row は保持対象にせず、cutover migration で旧 `master_*` table を drop する。

## Decision Points

### REQ-001 Legacy repository の扱い

- `issue`: 既存 UI / service は `MasterDictionaryRepository` と `MasterPersona*Repository` を中心に動いている。
- `background`: canonical schema / repository は作成済みだが、bootstrap はまだ旧 master 名 repository を配線している。
- `options`: 旧 repository を削除する、adapter で UI contract を維持しながら canonical repository へ寄せる、旧 schema と canonical schema の二重書きを行う。
- `recommendation`: adapter-first で existing UI / service boundary を維持し、内部永続化を canonical repository へ寄せる。二重書きは避け、旧 `master_*` table は cutover migration で drop する。
- `reasoning`: UI と backend の変更面を分離しつつ、旧 table への新規 write を止められる。
- `consequences`: 旧 repository file は参照が消えた段階で削除候補になる。legacy data backfill は不要で、旧 table の既存 row は保持対象にしない。
- `open_risks`: drop migration の順序と、旧 table 参照が完全に消えたことを確認する static / integration test が必要。

### REQ-002 Dictionary の canonical mapping

- `issue`: 旧 dictionary は `source`, `translation`, `category`, `origin`, `REC`, `EDID` を持つが、canonical `DICTIONARY_ENTRY` は `source_term`, `translated_term`, lifecycle / scope / source / term_kind を中心にする。
- `background`: ER では `REC` / `EDID` は dictionary の中核情報ではない。
- `options`: UI contract を canonical 名へ変更する、service adapter で既存 UI contract を維持する、個別 REC / EDID provenance 用 schema を足す。
- `recommendation`: 初期移行は service adapter で主要 UI 操作を維持し、`source` -> `source_term`、`translation` -> `translated_term`、`category` -> `term_kind`、`origin` -> `dictionary_source` に写像する。`REC` / `EDID` は UI、frontend / Wails contract、canonical mapping から外す。ただし XML parse 中の一時情報として使うことは許可し、永続化と UI 表示はしない。全 dictionary 登録経路は原文だけを trim し、重複判定は `trim(source_term) + translated_term` の完全一致にする。訳語の前後ノイズ除去は frontend の手動新規登録だけに限定し、XML import と既存更新では訳語を全体正規化しない。
- `reasoning`: 既存 UI の操作性を壊さず、canonical table への read / write へ移せる。
- `consequences`: `REC` / `EDID` 用の追加 provenance schema は作らない。既存 frontend / Wails DTO はこの方針に合わせて縮める。同じ原文でも訳語が違う entry は共存できる。
- `open_risks`: 既存 tests / presenter が `REC` / `EDID` 前提を持つ場合は、期待値更新が必要。frontend 手動新規登録限定の訳語ノイズ除去が XML import や更新へ広がらないよう test で固定する必要がある。

### REQ-003 Persona の canonical mapping

- `issue`: 旧 persona は `identity_key`、人物属性、`persona_body`、`generation_source_json`、dialogues を 1 entry に集約している。
- `background`: canonical ER は NPC identity を `NPC_PROFILE`、抽出スナップショット属性を `NPC_RECORD`、persona 本体を `PERSONA`、根拠を `PERSONA_FIELD_EVIDENCE` に分ける。
- `options`: 旧 entry model を残す、canonical join read model を作る、persona UI を大きく変更する。
- `recommendation`: 既存の主要操作は service / usecase boundary で維持し、read model は `PERSONA` + `NPC_PROFILE` + 必要な `NPC_RECORD` から組み立てる。既存判定 identity は canonical `NPC_PROFILE` の unique key と同じ `target_plugin_name + form_id + record_type` に固定する。`target_plugin_name` は lifecycle ownership ではなく identity / filter key として扱う。`generation_source_json`、`baseline_applied`、dialogue 表示、dialogue modal、dialogue count 表示は移行後に残さない。会話が見つからない NPC は入力 JSON parse 時点で生成対象から除外し、persona UI の対象として扱わない。preview / frontend / Wails contract は会話なし NPC や generic NPC の skip count を返さず、ペルソナ候補数、新規追加可能数、作成済み数を中心にする。persona edit modal はペルソナ要約、話し方、ペルソナ本文の 3 項目へ縮め、canonical column は `personality_summary` / `speech_style` / `persona_description` に写像する。`NPC_PROFILE` / `NPC_RECORD` 由来の identity / snapshot field を汎用編集対象にしない。
- `reasoning`: 画面骨格を保ちながら、ER の責務分離と UI contract の事実性を両立できる。
- `consequences`: `persona_body` は `persona_description`、summary は `personality_summary`、話し方は `speech_style` を中心に写像する。dialogue modal、dialogue count、generation source 表示は削除する。`zeroDialogueSkipCount` / `genericNpcCount` 相当の contract field は削除する。identity / snapshot field は表示用 read model に留め、編集対象にしない。AI 生成物が揃った後に `NPC_PROFILE` と `PERSONA` を 1 NPC 単位の同一 transaction で作成し、生成失敗時はどちらも残さない。
- `open_risks`: 既存 frontend の detail / modal contract を縮めるため、UI test の期待値更新が必要。

### REQ-004 AI settings と run state の置き場所

- `issue`: AI settings と generation run state は canonical `PERSONA` / `DICTIONARY_ENTRY` の foundation data とは性質が違う。
- `background`: ER では foundation data 作成 run を `JOB_PHASE_RUN` に含めない。
- `options`: 旧 master persona tables を残す、settings だけを table 化して run state は画面メモリに戻す、settings も config / secret store 側へ寄せて DB table を足さない。
- `recommendation`: `PERSONA_GENERATION_SETTINGS` だけを足す。`id = 1` の singleton page setting とし、provider 設定は単一で保持する。settings は provider / model を保存し、API key は既存 secret store seam を維持する。run state / message / timestamps は DB に保存しない。実行中プロセスの再開は行わず、生成する場合は人間が JSON を手動で読み直す。preview は現在読み込んだ JSON から「新規に追加できる NPC がいるか」を確認する。`generation_source_json` と `baseline_applied` は保持しない。
- `reasoning`: 保存済み `PERSONA` だけを生成結果の正とするため、run state table は data integrity に寄与しない。再起動後は JSON 未選択へ戻せばよく、DB に残すべき画面状態は provider / model だけで足りる。
- `consequences`: migration は旧 `master_persona_ai_settings` / `master_persona_run_status` を drop し、`PERSONA_GENERATION_SETTINGS` だけを作る。再起動後は JSON 未選択状態として人間に再選択させる。生成途中や失敗途中の中途半端な persona data は canonical `PERSONA` に残さず、生成完了した persona だけを保存する。`target_plugin_name` は `NPC_PROFILE` への ownership ではなく filter key として扱う。
- `open_risks`: drop / create migration の順序、singleton row の初期化、旧 table 参照が完全に消えたことを確認する static / integration test が必要。

## Functional Requirements

- `in_scope`: 旧 master dictionary / persona repository wiring の canonical adapter 化、service / usecase mapping、Wails DTO と frontend contract の縮小、DB 変更で出せない UI 表示の削除と既存画面の改修、HTML mock による改修後画面の review artifact 作成、dictionary XML import の canonical write、persona generation の canonical write、AI settings persistence、関連 test。
- `non_functional_requirements`: paid real AI API を test で呼ばない。既存 UI の主要導線を壊さない。foreign key / transaction 境界を維持する。
- `out_of_scope`: layout redesign、legacy data backfill、docs 正本化、翻訳 job UI への接続、`TRANSLATION_ARTIFACT` / export 実装。UI 改修そのものは out of scope ではない。
- `acceptance_basis`: UI 既存主要導線が通り、canonical table に read / write され、旧 `master_*` table は drop され、AI settings は singleton の `PERSONA_GENERATION_SETTINGS` から再起動後も復元される。run state は DB に保存されない。dictionary は原文 trim と `trim(source_term) + translated_term` 完全一致の重複判定を満たし、訳語ノイズ除去は frontend の手動新規登録だけに閉じる。`REC` / `EDID`、`generation_source_json`、`baseline_applied`、dialogue 表示、dialogue modal、dialogue count 表示、会話なし / generic NPC の skip count 表示と contract field は残らない。persona 既存判定は `target_plugin_name + form_id + record_type` で行う。会話が見つからない NPC は入力時に生成対象から除外される。再起動後は JSON 未選択状態になり、人間が JSON を手動で読み直す。preview は現在読み込んだ JSON から新規に追加できる NPC の有無を示す。生成途中の中途半端な persona data は残らない。

## Open Questions

- None. AI settings table 名は `PERSONA_GENERATION_SETTINGS`、保存粒度は `id = 1` の singleton page setting に固定する。run state は DB に保存しない。dictionary 重複判定は `trim(source_term) + translated_term` の完全一致に固定する。

## Required Reading

- `./plan.md`
- `./legacy-schema-ui-migration.review-er-diff.puml`
- `../../completed/2026-04-19-sqlite-migration-repositories/requirements-design.md`
- `../../completed/2026-04-19-sqlite-migration-repositories/scenario-design.md`
- `../../completed/2026-04-19-sqlite-migration-repositories/implementation-scope.md`
- `docs/er.md`
- `internal/infra/sqlite/migrations/003_canonical_er_v1_tables.sql`
- `internal/repository/foundation_data_repository.go`
- `internal/infra/sqlite/foundation_data_repository.go`
- `internal/bootstrap/app_controller.go`
- `internal/service/master_dictionary_*.go`
- `internal/service/master_persona_*.go`
- `internal/controller/wails/master_*_controller.go`
- `frontend/src/application/gateway-contract/master-*/`
- `frontend/src/ui/screens/master-*/`
