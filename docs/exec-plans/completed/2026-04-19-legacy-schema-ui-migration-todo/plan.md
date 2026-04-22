# Task Plan: 2026-04-19-legacy-schema-ui-migration-todo

- `workflow`: propose-plans
- `status`: closed-with-environment-validation-blocker
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
- `approval_record`: `approved-by-user-2026-04-20-usecase-scope-handoff`

## Copilot Result

- `completed_handoffs`: `schema-legacy-cutover`, `dictionary-read-detail-cutover`, `dictionary-create-update-delete-cutover`, `dictionary-xml-import-cutover`, `persona-read-detail-cutover`, `persona-ai-settings-restart-cutover`, `persona-json-preview-cutover`, `persona-generation-cutover`, `persona-edit-delete-cutover`, `final-validation-and-review`
- `touched_files`: `internal/` repository / service / usecase / controller / bootstrap、`frontend/src/` gateway contract / usecase / controller / presenter / UI、`internal/integrationtest/` integration tests、closeout 時の backend lint 補正。
- `implemented_scope`: canonical schema への cutover、legacy dictionary / persona master table drop、dictionary canonical read / mutation / XML import、persona canonical read / settings / preview / generation / edit-delete、public contract / UI field slimming。
- `test_results`: `python3 scripts/harness/run.py --suite structure` PASS、`npm run lint:backend` PASS、`npm run lint:frontend` PASS、`npm run test:backend` PASS、`npm run test:frontend` PASS、`npm run scan:sonar` PASS。`python3 scripts/harness/run.py --suite all` は最後の `npm run test:system` で未完了。
- `implementation_investigation`: stale test、arch lint、Sonar maintainability、persona generation atomicity、dictionary XML provenance、canonical / legacy read-write path の切り分けを実施。
- `ui_evidence`: Copilot final review では HTML mock と migrated UI behavior の整合を確認済み。Codex closeout では `test:system` が Wails dev server 起動前に止まり、追加 UI 実機証跡は取得できていない。
- `implementation_review_result`: all listed handoffs reached reviewer pass; final validation lane passed except Codex sandbox system-test blocker.
- `sonar_gate_result`: scanner execution PASS。Copilot report 時点で coverage 83.6%、maintainability HIGH / BLOCKER 0。
- `residual_risks`: Codex sandbox では `sysctl kern.osproductversion` が `Operation not permitted` になり、Wails CLI の OS version detection が失敗するため `npm run test:system` を完走できない。product-facing system test は sandbox 外で再実行が必要。
- `docs_changes`: product docs 正本化なし。workflow / report 義務追加は `.codex/`、`.github/`、`work_history/` に記録。

## Closeout Notes

- `canonicalized_artifacts`: none. `implementation-scope.md` は handoff 履歴であり docs 正本へ昇格しない。
- `closeout_adjustments`: Codex closeout で backend lint の機械的補正と Darwin Wails build helper を追加した。挙動変更ではなく validation unblock 用の最小修正。
- `follow_up`: port 5173 を空けた上で、sandbox 外、または Wails CLI の `sysctl` 依存を回避できる環境で `npm run test:system` と `python3 scripts/harness/run.py --suite all` を再実行する。

## Outcome

- TODO stub was promoted to design-bundle draft.
- `requirements-design.md` and `scenario-design.md` were started.
- Human-approved design decisions were reflected into `implementation-scope.md`.
- Oversized Copilot handoffs were compressed into 10 use-case vertical-slice RunSubagent units and approved for implementation.
- Copilot implementation completed the approved cutover scope and reviewer pass state was recorded.
- Codex closeout recorded the report obligation and closed this plan with a Wails system-test environment blocker.
