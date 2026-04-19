# Task Plan: 2026-04-19-legacy-schema-ui-migration-todo

- `workflow`: propose-plans
- `status`: implementation-scope-approved
- `lane_owner`: Codex owns the approved design bundle and implementation handoff. GitHub Copilot implements only from the approved implementation-scope.
- `task_id`: `2026-04-19-legacy-schema-ui-migration-todo`
- `task_mode`: backend-frontend-migration
- `request_summary`: 新 canonical schema / repository 作成後に、旧 master 名 schema / repository / service / UI を移行する設計を開始する。
- `goal`: 既存の master dictionary / master persona 画面と backend wiring を、canonical `FoundationDataRepository` と関連 repository へ寄せる。既存画面の主要導線は維持し、旧 `master_*` table への新規 read / write を止める。
- `constraints`: product code はこの plan では変更しない。旧 schema は作りかけの資産として扱い、backfill なしで drop する。
- `close_conditions`: requirements-design、scenario-design、HTML mock、ER diff、implementation-scope が human-approved handoff として整合している。

## Artifact Index

- `requirements_design`: `./requirements-design.md`
- `ui_design`: `./legacy-schema-ui-migration.ui.html`; HTML mock reviews required UI changes without a layout redesign.
- `scenario_design`: `./scenario-design.md`
- `diagramming`: `./legacy-schema-ui-migration.review-er-diff.puml`
- `implementation_scope`: `./implementation-scope.md`

## Workflow State

- `distiller`: `not-spawned`; predecessor plan と既存 code / docs の直接確認で design start に必要な事実を集めた。
- `designer`: `html-mock-policy`; Codex uses task-local HTML mock as the UI review artifact per human direction.
- `investigator`: `deferred`; runtime UI observation は implementation-scope 前または実装後 review で必要判定する。
- `human_review_gate`: approved-for-implementation-scope.

## Design Scope Notes

- 旧 `master_dictionary_entries` / `master_persona_*` table への新規 read / write を止める。
- DB 構成差分は `./legacy-schema-ui-migration.review-er-diff.puml` を review 用 ER 差分図として扱う。
- 既存画面の list / detail / create / update / delete / import / generation 導線は維持する。ただし canonical schema で出せない表示と contract field は塞ぐ。
- UI 改修は必要。対象は画面再設計ではなく、表示項目削除、modal / detail の縮小、再起動後状態の見せ方、frontend / Wails contract の縮小に限定する。
- dictionary は `DICTIONARY_ENTRY` と `XTRANSLATOR_TRANSLATION_XML` を中心に read / command model を作る。
- persona は `PERSONA`、`NPC_PROFILE`、必要に応じて `NPC_RECORD` から read / command model を作る。
- legacy data backfill はしない。旧 table の既存 row は保持対象にしない。
- 旧 `master_*` schema / table は cutover migration で drop する。
- `master_persona_ai_settings` は `PERSONA_GENERATION_SETTINGS` に置換し、`master_persona_run_status` は置換せず drop する。
- `PERSONA_GENERATION_SETTINGS` は `id = 1` の singleton page setting とし、`scope_key` は持たない。
- 再起動後は JSON 未選択状態に戻す。生成する場合は人間が JSON を手動で読み直し、preview は「新規に追加できる NPC がいるか」だけを確認する。
- 既存 `PERSONA` は上書きせず、作成済み NPC は skip する。
- persona の既存判定 identity は canonical `NPC_PROFILE` の unique key と同じ `target_plugin_name + form_id + record_type` に固定する。`target_plugin_name` は lifecycle ownership ではなく identity / filter key として扱う。
- preview / frontend / Wails contract は会話なし NPC や generic NPC の skip count を返さず、ペルソナ候補数、新規追加可能数、作成済み数を中心にする。
- persona generation は AI 生成物が揃った後に、1 NPC 単位の同一 transaction で `NPC_PROFILE` と `PERSONA` を作成する。生成途中や失敗途中の `NPC_PROFILE` / `PERSONA` は残さない。
- dictionary UI の `REC` / `EDID` は詳細、import summary、frontend / Wails contract に残さない。XML parse 中の一時情報として使うことは許可するが、永続化と UI 表示はしない。
- dictionary 登録時の重複判定は `trim(source_term) + translated_term` の完全一致とし、原文 trim は全登録経路で必須にする。
- 訳語の自動ノイズ除去は frontend の手動新規登録だけに限定し、XML import や既存更新の訳語は全体正規化しない。
- persona の `generation_source_json`、`baseline_applied`、dialogue 表示、dialogue modal、dialogue count 表示は残さない。会話が見つからない NPC は入力 JSON parse 時点で生成対象から除外し、persona UI の対象として扱わない。
- persona edit modal はペルソナ要約、話し方、ペルソナ本文の 3 項目へ縮め、canonical column は `personality_summary` / `speech_style` / `persona_description` に写像する。`NPC_PROFILE` / `NPC_RECORD` 由来の identity / snapshot field を汎用編集対象にしない。
- AI settings は canonical foundation data と混ぜず `PERSONA_GENERATION_SETTINGS` に分離する。run state は DB に保存せず、画面メモリだけで扱う。

## Routing Notes

- `required_reading`: `../../completed/2026-04-19-sqlite-migration-repositories/`, `docs/er.md`, `docs/detail-specs/master-dictionary.md`, `docs/scenario-tests/master-dictionary-management.md`, completed master persona plans, `internal/infra/sqlite/migrations/003_canonical_er_v1_tables.sql`, `internal/infra/sqlite/foundation_data_repository.go`, `internal/controller/wails/master_*`, `internal/service/master_*`, `internal/usecase/master_*`, `internal/bootstrap/app_controller.go`, `internal/repository/master_*`, `frontend/src/**/master-*`.
- `source_diagram_targets`: `docs/diagrams/conceptual/combined_perspective.puml`, `internal/infra/sqlite/migrations/003_canonical_er_v1_tables.sql`, `./legacy-schema-ui-migration.review-er-diff.puml`.
- `canonicalization_targets`: `N/A` until human-approved docs canonicalization is requested after implementation.
- `validation_commands`: draft `go test ./internal/infra/sqlite ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap`; draft `npm --prefix frontend run check`; draft `npm --prefix frontend run test`; required `python3 scripts/harness/run.py --suite structure`.

## HITL Status

- `functional_or_design_hitl`: `approved-for-implementation-scope`
- `approval_record`: `approved-by-user-2026-04-20-handoff-prompt-request`

## Copilot Result

- `completed_handoffs`: N/A
- `touched_files`: N/A
- `implemented_scope`: N/A
- `test_results`: N/A
- `implementation_investigation`: N/A
- `ui_evidence`: N/A
- `implementation_review_result`: N/A
- `sonar_gate_result`: N/A
- `residual_risks`: N/A
- `docs_changes`: N/A

## Closeout Notes

- `canonicalized_artifacts`: N/A
- `follow_up`: 人間が `implementation-scope.md` の Copilot Handoff Prompt を Copilot に渡す。実装完了後、Copilot completion packet を plan に反映し、必要なら Codex が docs 正本化を扱う。

## Outcome

- TODO stub was promoted to design-bundle draft.
- `requirements-design.md` and `scenario-design.md` were started.
- Human-approved design decisions were reflected into `implementation-scope.md`.
- Copilot handoff prompt is ready in `implementation-scope.md`.
