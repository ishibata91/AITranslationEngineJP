# Implementation Scope: 2026-04-19-legacy-schema-ui-migration-todo

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: approved-by-user-2026-04-20-usecase-scope-handoff
- `scope_revision_reason`: implementation-orchestrate rerouted broad handoffs, then the 67-step layer split was judged over-fragmented. This revision uses use-case vertical slices with one validation intent per handoff.
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`

## Source Artifacts

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `N/A`
- `ui_mock_artifact`: `./legacy-schema-ui-migration.ui.html`
- `scenario_design`: `./scenario-design.md`
- `diagramming`: `./legacy-schema-ui-migration.review-er-diff.puml`

## Fixed Decisions

- 旧 `master_*` table は backfill なしで drop する。
- 旧 master repository / service / Wails / frontend contract は canonical schema へ寄せる。
- `master_persona_ai_settings` は `PERSONA_GENERATION_SETTINGS(id = 1)` に置換する。
- `master_persona_run_status` は置換せず drop し、run state は DB に保存しない。
- 再起動後は JSON 未選択状態に戻し、JSON は人間が手動で読み直す。
- persona 既存判定 identity は `target_plugin_name + form_id + record_type` に固定する。
- `target_plugin_name` は lifecycle ownership ではなく identity / filter key として扱う。
- persona generation は AI 生成物が揃った後、1 NPC 単位の同一 transaction で `NPC_PROFILE` と `PERSONA` を作成する。
- 生成途中や失敗途中の `NPC_PROFILE` / `PERSONA` は残さない。
- 会話が見つからない NPC は入力 JSON parse 時点で除外し、preview / generation target / persona UI へ出さない。
- preview / frontend / Wails contract は会話なし NPC や generic NPC の skip count を返さない。
- preview はペルソナ候補数、新規追加可能数、作成済み数を中心にする。
- persona edit modal はペルソナ要約、話し方、ペルソナ本文の 3 項目だけを編集対象にする。
- `NPC_PROFILE` / `NPC_RECORD` 由来の identity / snapshot field は表示用 read model に留め、汎用編集対象にしない。
- `generation_source_json`、`baseline_applied`、dialogue 表示、dialogue modal、dialogue count 表示は残さない。
- dictionary の `REC` / `EDID` は UI、frontend / Wails contract、canonical mapping から外す。
- dictionary の `REC` / `EDID` は XML parse 中の一時情報としてのみ使える。
- dictionary 重複判定は `trim(source_term) + translated_term` の完全一致に固定する。
- 原文 trim は全 dictionary 登録経路で必須にする。
- 訳語の前後ノイズ除去は frontend の手動新規登録だけに限定する。

## Handoff Split Basis

- split_rule: `1 independently verifiable use-case slice x 1 validation intent`
- use_case_definition: domain 名や画面名ではなく、人間または system が開始する処理単位で切る。
- layer_policy: この plan では layer 単位の 67 分割は過剰なため、1 use case の完了に必要な backend / frontend 変更を同じ handoff に含める。
- fallback_policy: それでも context が不足する場合だけ、当該 use case を backend contract と frontend UI / state に二分して propose-plans へ戻す。
- validation_ownership_policy: 中間 handoff の validation は `completion_signal` を直接検証する focused test に限定し、全体 `go test`、frontend check / 全体 test、structure harness、UI 起動確認は `final-validation-and-review` に寄せる。
- object_names: `dictionary` / `persona` は対象 object 名としてだけ使い、分割根拠にはしない。

## Handoff Order

1. `schema-legacy-cutover`
2. `dictionary-read-detail-cutover`
3. `dictionary-create-update-delete-cutover`
4. `dictionary-xml-import-cutover`
5. `persona-read-detail-cutover`
6. `persona-ai-settings-restart-cutover`
7. `persona-json-preview-cutover`
8. `persona-generation-cutover`
9. `persona-edit-delete-cutover`
10. `final-validation-and-review`

## Handoffs

### `schema-legacy-cutover`

- `implementation_target`: schema lifecycle cutover
- `owned_scope`:
  - `internal/infra/sqlite/migrations/` legacy drop and `PERSONA_GENERATION_SETTINGS` creation
  - `internal/infra/sqlite/dbinit/` migration registration only if the local migration mechanism requires it
  - schema / migration tests under the same backend packages
- `depends_on`:
  - `./requirements-design.md`
  - `./scenario-design.md`
  - `./legacy-schema-ui-migration.review-er-diff.puml`
- `validation_commands`:
  - `go test ./internal/infra/sqlite -run 'TestSchemaCutover'`
- `completion_signal`: legacy dictionary / persona master tables are dropped without backfill, `PERSONA_GENERATION_SETTINGS(id = 1)` exists, and no persisted run status table is created.
- `notes`:
  - Use case: schema cutover.
  - Validation intent: schema / migration completion only. Downstream repository, service, controller, and bootstrap cutover is validated by later use-case handoffs and final validation.
  - Do not create dual-write compatibility tables.
  - Do not update canonical docs in this handoff.

### `dictionary-read-detail-cutover`

- `implementation_target`: dictionary list / search / filter / detail
- `owned_scope`:
  - backend repository / SQLite / service / usecase / Wails read paths for canonical `DICTIONARY_ENTRY`
  - `internal/bootstrap/app_controller.go` wiring needed only for dictionary read paths
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller read paths
  - `frontend/src/ui/screens/master-dictionary/` read / detail display
  - focused backend and frontend tests for read / detail
- `depends_on`:
  - `schema-legacy-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*DictionaryReadDetailCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-dictionary/master-dictionary.presenter.test.ts src/application/usecase/master-dictionary/master-dictionary.usecase.test.ts src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts src/controller/master-dictionary/master-dictionary-screen-controller.test.ts src/ui/App.test.ts -t 'dictionary-read-detail-cutover'`
- `completion_signal`: dictionary list / search / filter / detail reads canonical data and no product-facing contract or UI exposes `REC` / `EDID`.
- `notes`:
  - Use case: read / detail.
  - Layer exception: this vertical slice crosses layers so the read behavior can be validated end to end in one lane.
  - Do not update canonical docs in this handoff.

### `dictionary-create-update-delete-cutover`

- `implementation_target`: dictionary create / update / delete
- `owned_scope`:
  - backend repository / SQLite / service / usecase / Wails command paths for canonical `DICTIONARY_ENTRY`
  - `internal/bootstrap/app_controller.go` wiring needed only for dictionary command paths
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller command paths
  - `frontend/src/ui/screens/master-dictionary/` create / edit / delete UI
  - focused backend and frontend tests for mutation and duplicate detection
- `depends_on`:
  - `dictionary-read-detail-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*DictionaryCreateUpdateDeleteCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/usecase/master-dictionary/master-dictionary.usecase.test.ts src/controller/master-dictionary/master-dictionary-screen-controller.test.ts src/ui/App.test.ts -t 'dictionary-create-update-delete-cutover'`
- `completion_signal`: create / update / delete writes canonical dictionary data, trims source terms on all registration paths, and detects duplicates by `trim(source_term) + translated_term`.
- `notes`:
  - Use case: create / update / delete.
  - Translated-term edge cleanup is limited to frontend manual new registration.
  - Do not normalize XML import or update translated terms globally.
  - Do not update canonical docs in this handoff.

### `dictionary-xml-import-cutover`

- `implementation_target`: dictionary XML import
- `owned_scope`:
  - backend XML parse / repository / SQLite / service / usecase / Wails import paths
  - canonical `XTRANSLATOR_TRANSLATION_XML` provenance and imported `DICTIONARY_ENTRY` persistence
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller import paths
  - `frontend/src/ui/screens/master-dictionary/` import UI and summary
  - focused backend and frontend tests with small XML fixtures
- `depends_on`:
  - `dictionary-create-update-delete-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*DictionaryXMLImportCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/usecase/master-dictionary/master-dictionary.usecase.test.ts src/controller/master-dictionary/master-dictionary-screen-controller.test.ts src/controller/wails/master-dictionary.gateway.test.ts src/ui/App.test.ts -t 'dictionary-xml-import-cutover'`
- `completion_signal`: XML import stores provenance and entries canonically while keeping `REC` / `EDID` parser-only, with no `selectedRec` or `EDID` in product-facing import UI / contract.
- `notes`:
  - Use case: import.
  - This handoff closes dictionary migration; later handoffs must not modify dictionary behavior except for shared compile fixes.
  - Do not update canonical docs in this handoff.

### `persona-read-detail-cutover`

- `implementation_target`: persona list / filter / detail
- `owned_scope`:
  - backend repository / SQLite / service / usecase / Wails read paths for `PERSONA` + `NPC_PROFILE` + needed `NPC_RECORD`
  - `internal/bootstrap/app_controller.go` wiring needed only for persona read paths
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller read paths
  - `frontend/src/ui/screens/master-persona/` list / detail display
  - focused backend and frontend tests for read / detail and plugin filter
- `depends_on`:
  - `dictionary-xml-import-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*PersonaReadDetailCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-persona/master-persona.presenter.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts src/controller/master-persona/master-persona-screen-controller.test.ts src/controller/wails/master-persona.gateway.test.ts src/ui/App.test.ts -t 'persona-read-detail-cutover'`
- `completion_signal`: persona list / detail reads canonical join data and removes generation source, baseline, dialogue payload, dialogue count, and dialogue modal support.
- `notes`:
  - Use case: read / detail.
  - Identity / snapshot fields may be read-only display data, but must not become generic editable fields.
  - Do not update canonical docs in this handoff.

### `persona-ai-settings-restart-cutover`

- `implementation_target`: persona AI settings save / restore and restart state
- `owned_scope`:
  - backend repository / SQLite / service / usecase / Wails settings save / restore paths
  - existing secret-store boundary for API key handling
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller settings paths
  - `frontend/src/ui/screens/master-persona/` settings UI and restart state display
  - focused backend and frontend tests for settings persistence and restart behavior
- `depends_on`:
  - `persona-read-detail-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*PersonaAISettingsRestartCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-persona/master-persona.presenter.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts src/controller/master-persona/master-persona-screen-controller.test.ts src/controller/wails/master-persona.gateway.test.ts src/ui/App.test.ts -t 'persona-ai-settings-restart-cutover'`
- `completion_signal`: provider / model are restored after restart, API key stays outside DB, JSON selection returns to unselected state, and run state is not persisted.
- `notes`:
  - Use case: settings save / restore.
  - Do not create `PERSONA_GENERATION_RUN_STATUS`.
  - Do not update canonical docs in this handoff.

### `persona-json-preview-cutover`

- `implementation_target`: persona JSON load / parse / preview
- `owned_scope`:
  - backend JSON parse / existing identity lookup / service / usecase / Wails preview paths
  - canonical identity lookup by `target_plugin_name + form_id + record_type`
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller preview paths
  - `frontend/src/ui/screens/master-persona/` JSON selection and preview UI
  - focused backend and frontend tests with JSON fixtures
- `depends_on`:
  - `persona-ai-settings-restart-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*PersonaJSONPreviewCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-persona/master-persona.presenter.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts src/controller/master-persona/master-persona-screen-controller.test.ts src/controller/wails/master-persona.gateway.test.ts src/ui/App.test.ts -t 'persona-json-preview-cutover'`
- `completion_signal`: preview excludes no-dialogue NPCs at parse time and returns candidate count, newly addable count, and existing count without zero-dialogue / generic skip counts.
- `notes`:
  - Use case: preview.
  - `target_plugin_name` is identity / filter key only, not lifecycle ownership.
  - Do not update canonical docs in this handoff.

### `persona-generation-cutover`

- `implementation_target`: persona AI generation and canonical write
- `owned_scope`:
  - backend AI generation orchestration, existing identity re-check, transaction write, and Wails generation command
  - repository / SQLite command paths for atomic `NPC_PROFILE` + `PERSONA` creation
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller generation paths
  - `frontend/src/ui/screens/master-persona/` generation progress and result UI
  - focused backend and frontend tests using fake AI transport
- `depends_on`:
  - `persona-json-preview-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*PersonaGenerationCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-persona/master-persona.presenter.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts src/controller/master-persona/master-persona-screen-controller.test.ts src/controller/wails/master-persona.gateway.test.ts src/ui/App.test.ts -t 'persona-generation-cutover'`
- `completion_signal`: generation never overwrites existing persona, writes `NPC_PROFILE` + `PERSONA` only after complete AI output, and leaves no partial rows on failure.
- `notes`:
  - Use case: generation.
  - Run state remains screen memory only.
  - Tests must not require real AI API credentials.
  - Do not update canonical docs in this handoff.

### `persona-edit-delete-cutover`

- `implementation_target`: persona edit / delete
- `owned_scope`:
  - backend repository / SQLite / service / usecase / Wails edit / delete paths
  - frontend gateway contract / DTO mapping / state / presenter / usecase / controller edit / delete paths
  - `frontend/src/ui/screens/master-persona/` edit modal and delete UI
  - focused backend and frontend tests for editable fields and delete behavior
- `depends_on`:
  - `persona-generation-cutover`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap -run 'Test.*PersonaEditDeleteCutover'`
  - `npm --prefix frontend run test -- --runInBand src/application/presenter/master-persona/master-persona.presenter.test.ts src/application/usecase/master-persona/master-persona.usecase.test.ts src/controller/master-persona/master-persona-screen-controller.test.ts src/controller/wails/master-persona.gateway.test.ts src/ui/App.test.ts -t 'persona-edit-delete-cutover'`
- `completion_signal`: edit accepts only persona summary, speech style, and persona body; identity / snapshot fields remain read-only and manual persona creation is not reintroduced.
- `notes`:
  - Use case: update / delete.
  - This handoff closes persona migration; later handoffs must not modify persona behavior except for shared compile fixes.
  - Do not update canonical docs in this handoff.

### `final-validation-and-review`

- `implementation_target`: final implementation review and UI evidence
- `owned_scope`:
  - final implementation review record
  - UI evidence record using existing local review conventions
  - completion packet aggregation
- `depends_on`:
  - `schema-legacy-cutover`
  - `dictionary-xml-import-cutover`
  - `persona-edit-delete-cutover`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `go test ./internal/...`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
  - `npm run dev:wails:docker-mcp`
- `completion_signal`: review confirms every use-case handoff is complete, no old `master_*` read / write path remains for migrated surfaces, UI evidence matches the HTML mock behaviorally, and completion packet includes validation, UI evidence, residual risks, and `docs_changes: none`.
- `notes`:
  - Playwright MCP should connect through `http://host.docker.internal:34115` when UI evidence is needed.
  - Use fake AI / fixture data for UI evidence.
  - Paid provider connectivity is out of scope.
  - Do not update canonical docs in this handoff.

## Explicitly Out Of Scope

- Legacy data backfill.
- Any compatibility adapter that writes old and canonical tables at the same time.
- New REC / EDID provenance schema.
- `PERSONA_GENERATION_RUN_STATUS` or other persisted run status table.
- Persona overwrite, diff preview, rollback UI, or auto-resume from previous JSON.
- Dialogue modal, dialogue count display, or dialogue payload contract in the migrated persona UI.
- Generic NPC / no-dialogue skip count display or public contract.
- Manual persona creation.
- `TRANSLATION_ARTIFACT` / export implementation.
- Canonical docs正本更新.
- Changes under `docs/`, `.codex/`, `.github/skills`, or `.github/agents` by Copilot.

## Copilot Handoff Prompt

```text
[$implementation-orchestrate](.github/skills/implementation-orchestrate/SKILL.md)

Use the approved implementation scope:
docs/exec-plans/completed/2026-04-19-legacy-schema-ui-migration-todo/implementation-scope.md

Implement the handoffs in order. Treat each handoff as one RunSubagent lane. Do not edit docs, .codex, .github/skills, or .github/agents. Use fake AI transport in tests. If a lane is still too large, stop and return a proposed backend-contract / frontend-UI split for that lane instead of continuing partially. Return the completion packet defined in the implementation scope.
```

## Completion Packet

Copilot は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `implementation_review_result`
- `sonar_gate_result`
- `residual_risks`
- `docs_changes: none`
