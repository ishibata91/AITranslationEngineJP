# Implementation Scope: translation-input-import-non-dialogue

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved by current handoff
- `approval_record`: current user instruction states approved design artifacts are the source
- `module_entry`: `.claude/skills/implementation-module/SKILL.md`
- `handoff_runtime`: `claude-module`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `component_diagram`: `./design-diff-diagram.md`
- `screen_design_diff`: `N/A`
- `related_code`: `internal/service/translation_input_import_service.go`

## Fixed Decisions

- `unanswered_questions`: `0`
- 承認済み詳細仕様差分だけを implementation handoff の仕様根拠にする。
- 画面変更はないため、frontend handoff、Wails handoff、UI handoff は作らない。
- Wails DTO と既存 import usecase の入出力形状は変更しない。
- docs 正本本文、作業流れ、既存設計成果物は implementation-module の範囲外にする。
- secret は扱わない。
- `E2E` は UI 人間操作起点だけを指す。本 task の受け入れ証明は単体テストを主対象にする。

## Approved Implementation Scope

- `internal/service/translation_input_import_service.go` の JSON decode 対象を、`dialogue_groups` だけから会話以外の取り込み対象へ広げる。
- `translationInputDocument` は `target_plugin` と既存 `dialogue_groups` に加えて、`items`、`cells`、`locations`、`magic`、`messages`、`system`、`load_screens`、`npcs`、`quests` を保持できる。
- `npcs` は配列ではなく、FormID を key にした object として扱う。
- `quests` の `stages` と `objectives` は、最上位配列ではなく `quests` の子要素として扱う。
- `dialogue_groups` が空でも、`target_plugin` と取り込み可能な会話以外の翻訳レコードが 1 件以上あれば登録結果を作れる。
- JSON 内の `type` が `RECORD FIELD` 形式の場合は、既存の `parseRecordAndSubrecord` と同じ `RECORD:FIELD` 判断へそろえる。
- 各要素の `id`、`editor_id`、`type`、本文候補を `preparedTranslationRecord` と `preparedTranslationField` へ正規化する。
- 原文が空の翻訳フィールドは、レコード登録対象に含めても field 登録対象にはしない。
- 会話以外の要素でも、未定義 RecordType と SubrecordType の組み合わせは既存 warning 規則で扱う。
- `prepared.records`、`categories`、`warnings`、`fieldCount`、`SampleFields` は会話以外のレコードとフィールドを反映する。
- `persistPreparedImport` と transaction 境界は変更しない。
- 単語翻訳フェーズの候補判定は、取り込み元の JSON 配列名に依存しない状態を維持する。

## Out Of Scope

- frontend state、Svelte 表示コンポーネント、style、Storybook、fixture は変更しない。
- Wails DTO、controller request / response DTO、gateway、generated binding は変更しない。
- DB schema、migration、repository contract は変更しない。
- `extractData.pas` は変更しない。
- docs 正本本文、docs index、他 task folder は変更しない。
- 最上位 `responses`、`stages`、`objectives` は確認済み JSON 形状ではないため、新規 top-level 取り込み対象にしない。
- DOOR:FULL、FLOR:FULL、FURN:FULL は 13 種別外であり、単語翻訳フェーズ候補として扱わない。
- 観測ログ追加は不要とする。

## Dependencies

| handoff | depends_on | reason |
| --- | --- | --- |
| `backend-import-non-dialogue-records` | `なし` | service 内の parse と prepare 変更だけで開始できる。 |
| `unit-test-import-non-dialogue-records` | `backend-import-non-dialogue-records` | 実装済み decode と prepare の公開振る舞いを単体テストで証明する。 |

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `backend-import-non-dialogue-records` | `なし` | `なし` | `なし` |
| `wave-2` | `unit-test-import-non-dialogue-records` | `backend-import-non-dialogue-records` | `なし` | `depends_on` |

## Verification Units

- backend 実装の局所検証は `go test ./internal/service` とする。
- 単体テストの局所検証は `go test ./internal/service` とする。
- backend 全体の回帰確認は implementation-module の final validation で `python3 scripts/harness/run.py --suite backend-local` を扱う。
- frontend、Wails、UI 人間操作E2E、Storybook review は本 task の検証単位に含めない。

## Handoffs

### `backend-import-non-dialogue-records`

- `implementation_target`: xEdit 抽出 JSON の会話以外の翻訳レコードを import 準備へ流す。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`
- `diagram_basis`: `./design-diff-diagram.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/service/translation_input_import_service.go`
  - `translationInputDocument`
  - 会話以外の JSON element を表す private struct
  - `decodeTranslationInputDocument`
  - `prepareImport`
  - 追加する private normalize / prepare helper
  - `defaultTranslationFieldDefinitions` の必要最小追加
- `scope_size_estimate`:
  - `files`: `1 product file`
  - `changed_lines`: `300-600 lines`
  - `scale`: `通常`
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `internal/service/translation_input_import_service.go` の `translationInputDocument` に会話以外の入力 shape を追加し、`decodeTranslationInputDocument` の `dialogue_groups` 必須判定を「target_plugin と importable record の存在判定」へ変える。対応する完了条件は「dialogue_groups が空でも会話以外の取り込み可能 record で登録結果を作れる」である。この変更が後続 helper の入力境界を固定するため初手にする。
- `validation_commands`:
  - `go test ./internal/service`
- `completion_signal`:
  - `translationInputDocument` が `items`、`cells`、`locations`、`magic`、`messages`、`system`、`load_screens`、`npcs`、`quests` を decode できる。
  - `npcs` の object 形状を FormID と record 本体へ正規化できる。
  - `quests[].stages` と `quests[].objectives` を親 quest と区別できる field order または record 入力へ正規化できる。
  - `type` の `RECORD FIELD` 表現が既存 `RECORD:FIELD` 判定へそろう。
  - 原文候補は `name`、`description`、`text`、`title`、stage text、objective text から既存 prepared field へ渡る。
  - 空原文は field 登録対象にしない。
  - 未定義 RecordType と SubrecordType は既存 warning に残る。
  - `persistPreparedImport`、transaction 境界、repository contract、Wails DTO は変更されない。
  - `go test ./internal/service` が通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本番経路は `ImportXEditJSON` / `ImportXEditJSONWithContent` -> `prepareImportFromContent` -> `decodeTranslationInputDocument` -> `prepareImport` -> `persistPreparedImport` である。
  - `defaultTranslationFieldDefinitions` は warning と translatable 表示の既存規則であり、必要な REC だけを追加する。
  - 単語翻訳フェーズ候補の最終抽出は既存 `recclassification.IsTermTarget` と既存 phase 側処理に委ねる。

### `unit-test-import-non-dialogue-records`

- `implementation_target`: 会話以外の翻訳レコード取り込みを単体テストで証明する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `internal/service/translation_input_import_service_test.go`
  - 必要最小限の test fixture string または helper
- `scope_size_estimate`:
  - `files`: `1 test file`
  - `changed_lines`: `120-260 lines`
  - `scale`: `通常`
- `depends_on`: `backend-import-non-dialogue-records`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`: `internal/service/translation_input_import_service_test.go` に `ImportXEditJSONWithContent` で `dialogue_groups` なしの会話以外 JSON を登録できる単体テストを追加する。対応する完了条件は「target_plugin と取り込み可能な会話以外 record があれば登録結果を作れる」である。公開入口から最重要仕様を先に固定できるため初手にする。
- `validation_commands`:
  - `go test ./internal/service`
- `completion_signal`:
  - `dialogue_groups` なし、かつ `target_plugin` と会話以外の取り込み可能 record がある JSON が成功する。
  - `items`、`magic`、`locations`、`messages`、`system`、`load_screens`、`npcs`、`quests[].stages`、`quests[].objectives` の代表 field が persisted draft または summary に反映される。
  - `NPC_:FULL` と `NPC_:SHRT` が別 field として扱われる。
  - 13 種別に該当する REC の非空原文が登録結果に残る。
  - 13 種別外または未定義 field definition の扱いが既存 warning 規則から外れない。
  - 空原文は field count と persisted field に入らない。
  - `go test ./internal/service` が通過する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 期待結果の元ネタは `./detail-spec-diff.md` の `translation-input-intake-REQ-002` と `term-translation-phase-REQ-008` である。
  - テストは frontend、Wails、SQLite integration、実機 dictionary file に依存しない。

## Implementation Module Handoff

implementation-module は次の順で実行する。

1. `backend-import-non-dialogue-records` を `implement-backend` で実装する。
2. `unit-test-import-non-dialogue-records` を `tests-unit` で実装する。
3. final validation で backend-local の回帰確認を行う。

implementation-module へ渡さない事項:

- docs 正本反映
- Storybook review
- UI 人間操作E2E
- Wails DTO 変更
- DB migration
- `extractData.pas` 変更

## Completion Packet

実装モジュールは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `N/A`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate の結果
- `harness_gate_result`
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`
