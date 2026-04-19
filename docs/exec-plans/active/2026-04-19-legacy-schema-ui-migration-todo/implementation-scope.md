# Implementation Scope: 2026-04-19-legacy-schema-ui-migration-todo

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: approved-by-user-2026-04-20-handoff-prompt-request
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

## Handoff Order

1. `backend-dictionary-canonical-cutover`
2. `backend-persona-canonical-cutover`
3. `frontend-contract-ui-cutover`
4. `tests-legacy-schema-ui-migration`
5. `review-legacy-schema-ui-migration`

## Handoffs

### `backend-dictionary-canonical-cutover`

- `implementation_target`: backend dictionary persistence / service / Wails boundary
- `owned_scope`:
  - `internal/infra/sqlite/migrations/` legacy master dictionary drop and canonical schema wiring if needed
  - `internal/repository/` dictionary repository adapters or removals needed for canonical `DICTIONARY_ENTRY`
  - `internal/infra/sqlite/` dictionary SQLite concrete and tests
  - `internal/service/master_dictionary_*`
  - `internal/usecase/master_dictionary_*`
  - `internal/controller/wails/master_dictionary*_controller.go`
  - `internal/bootstrap/app_controller.go`
  - related backend tests under the same packages
- `depends_on`:
  - `./requirements-design.md`
  - `./scenario-design.md`
  - `./legacy-schema-ui-migration.review-er-diff.puml`
- `validation_commands`:
  - `go test ./internal/infra/sqlite ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: master dictionary read / create / update / delete / XML import reads and writes canonical `DICTIONARY_ENTRY` and `XTRANSLATOR_TRANSLATION_XML`. New writes to old `master_dictionary_entries` are gone. `REC` / `EDID` are absent from frontend / Wails contract and UI-facing response fields. Duplicate detection uses `trim(source_term) + translated_term`.
- `notes`:
  - Do not create a new REC / EDID provenance schema.
  - Do not normalize XML import translated terms globally.
  - Frontend manual new registration may strip obvious translated-term edge noise, but backend must not broaden that rule to XML import or update paths.
  - Do not update canonical docs in this handoff.

### `backend-persona-canonical-cutover`

- `implementation_target`: backend persona persistence / generation / service / Wails boundary
- `owned_scope`:
  - `internal/infra/sqlite/migrations/` legacy master persona table drop and `PERSONA_GENERATION_SETTINGS` creation
  - `internal/repository/` persona query / command / AI settings repositories and adapters
  - `internal/infra/sqlite/` persona SQLite concrete and tests
  - `internal/service/master_persona_*`
  - `internal/usecase/master_persona_*`
  - `internal/controller/wails/master_persona_controller.go`
  - `internal/bootstrap/app_controller.go`
  - related backend tests under the same packages
- `depends_on`:
  - `backend-dictionary-canonical-cutover`
- `validation_commands`:
  - `go test ./internal/infra/sqlite ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: master persona list / detail read from `PERSONA` + `NPC_PROFILE` + needed `NPC_RECORD`. `master_persona_entries`, `master_persona_ai_settings`, and `master_persona_run_status` are no longer used for new reads / writes. `PERSONA_GENERATION_SETTINGS(id = 1)` persists provider / model only. API key remains behind the existing secret-store seam. No `PERSONA_GENERATION_RUN_STATUS` table / row is created. Runtime run state stays in memory only.
- `notes`:
  - Existing persona identity must be `target_plugin_name + form_id + record_type`.
  - `target_plugin_name` is identity / filter key, not persona lifecycle ownership.
  - Preview and execute must exclude no-dialogue NPCs at parse time.
  - Preview / Wails contract must not return zero-dialogue or generic NPC skip counts.
  - Preview should expose persona candidate count, newly addable count, existing count, target plugin, file name, and status as needed.
  - Generation must re-check existing `PERSONA` immediately before write and must never overwrite existing persona.
  - Create `NPC_PROFILE` and `PERSONA` in the same transaction only after AI output is complete for that NPC.
  - A failed or partial generation must not leave `NPC_PROFILE` or `PERSONA` rows for that NPC.
  - Edit/update accepts only persona summary, speech style, and persona body fields.
  - Remove `generation_source_json`, `baseline_applied`, dialogue payload, dialogue modal support, and dialogue count from product-facing contracts.
  - Do not update canonical docs in this handoff.

### `frontend-contract-ui-cutover`

- `implementation_target`: frontend gateway contract, presenter, state, and UI screens
- `owned_scope`:
  - `frontend/src/application/gateway-contract/master-dictionary/`
  - `frontend/src/application/gateway-contract/master-persona/`
  - `frontend/src/application/contract/` affected dictionary / persona contracts
  - `frontend/src/application/store/` affected dictionary / persona state
  - `frontend/src/application/presenter/` affected dictionary / persona presenters
  - `frontend/src/application/usecase/` affected dictionary / persona usecases
  - `frontend/src/controller/master-dictionary/`
  - `frontend/src/controller/master-persona/`
  - `frontend/src/controller/runtime/master-*`
  - `frontend/src/controller/wails/` affected dictionary / persona gateways and DTO mapping
  - `frontend/src/ui/screens/master-dictionary/`
  - `frontend/src/ui/screens/master-persona/`
  - affected app shell route wiring
- `depends_on`:
  - `backend-dictionary-canonical-cutover`
  - `backend-persona-canonical-cutover`
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: UI matches the task-local HTML mock at the behavioral level. Dictionary UI no longer shows or submits `REC` / `EDID`. Persona preview shows candidate count, newly addable count, and existing count only. No-dialogue / generic skip count UI is gone. Persona detail no longer shows generation source, baseline, dialogue count, or dialogue modal. Persona edit form exposes only summary, speech style, and body. Restarted app state starts from JSON unselected state while saved AI provider / model is restored.
- `notes`:
  - Do not reintroduce manual persona creation.
  - Do not show implementation planning text in the UI.
  - Existing read-only NPC identity / snapshot fields may be displayed if supplied by backend, but they must not be editable.
  - Keep UI controls close to the current screen structure; this is not a layout redesign.
  - Tests must not require real AI API credentials.
  - Do not update canonical docs in this handoff.

### `tests-legacy-schema-ui-migration`

- `implementation_target`: backend and frontend tests for migration behavior
- `owned_scope`:
  - backend repository / service / usecase / controller / bootstrap tests affected by the cutover
  - frontend gateway / presenter / usecase / screen controller tests affected by the cutover
  - focused UI tests for the changed dictionary and persona views where existing structure places them
- `depends_on`:
  - `backend-dictionary-canonical-cutover`
  - `backend-persona-canonical-cutover`
  - `frontend-contract-ui-cutover`
- `validation_commands`:
  - `go test ./internal/...`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: Tests prove canonical read / write, old table drop, no backfill, no old master write, dictionary trim + duplicate rule, no REC / EDID contract, frontend-only translated-term edge noise cleanup, persona identity key, no overwrite, no-dialogue parse exclusion, no zero-dialogue / generic skip count contract, no run status persistence, settings persistence, transient run state, and no partial `NPC_PROFILE` / `PERSONA` rows on generation failure.
- `notes`:
  - Use fake AI transport only.
  - Use temp DB fixtures.
  - Include at least one failure case proving transaction rollback for persona generation.
  - Include at least one restart/bootstrap recreation case proving JSON unselected state and restored provider / model.
  - Do not update canonical docs in this handoff.

### `review-legacy-schema-ui-migration`

- `implementation_target`: implementation review and UI evidence
- `owned_scope`:
  - implementation-review record
  - UI evidence record using existing local review conventions
- `depends_on`:
  - `tests-legacy-schema-ui-migration`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `go test ./internal/...`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
  - `npm run dev:wails:docker-mcp`
- `completion_signal`: Review confirms the implemented behavior matches the approved design bundle and HTML mock. No old `master_*` read / write path remains for the migrated surfaces. `PERSONA_GENERATION_RUN_STATUS` is absent. Run state is transient. UI no longer exposes removed DB-backed fields or removed skip counts.
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
docs/exec-plans/active/2026-04-19-legacy-schema-ui-migration-todo/implementation-scope.md

Implement the handoffs in order. Treat each handoff as the owned scope for one implementation lane. Do not edit docs, .codex, .github/skills, or .github/agents. Use fake AI transport in tests. Return the completion packet defined in the implementation scope.
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
