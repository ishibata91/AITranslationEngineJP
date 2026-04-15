# Work Plan

- workflow: work
- status: planned
- lane_owner: orchestrate
- scope: master-persona-management
- task_id: persona-management
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml
- parent_phase: implementation-lane

## Request Summary

- マスターペルソナページの実装を開始する。
- design bundle が揃った時点で human review へ渡し、その承認前は実装へ進めない。

## Decision Basis

- `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml` に、一覧、作成、詳細確認、更新、削除、JSON からの再構築導線が完了条件として定義されている。
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md` では、マスターペルソナはマスター辞書と独立した基盤データとして定義されている。
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md` では、マスターペルソナは別 task へ切り分け済みである。
- orchestrate 契約により、`implement` task は `distill` と `design` を先行し、design bundle 完了後に HITL gate を立てる。

## Task Mode

- `task_mode`: `implement`
- `goal`: マスターペルソナ管理ページの要件、UI モック、scenario、implementation brief を task-local artifact として固定し、human review 後に実装へ進められる状態を作る。
- `constraints`: orchestrate 自身では実装しない。`docs/` 正本は human 先行で更新する。design bundle 完了後は human review 完了まで停止する。
- `close_conditions`: design bundle 完了後に human review gate を記録すること。human review 完了後にのみ implementation-scope 以降へ進むこと。最終 close では implementation-review と ui-check を必須とすること。

## Facts

- usecase の完了条件は、一覧参照、作成、詳細確認、更新、削除、`extractData.pas` 由来 JSON からの再構築導線で構成される。
- `related_screens` には `persona-management.md` が記録されているが、現時点の `docs/` には該当 screen-design / detail-spec / mock 正本が未作成である。
- frontend shell には `master-persona` ルート識別子が既に存在する。
- structure harness は 2026-04-15 に pass 済みである。

## Functional Requirements

- `summary`:
  - マスターペルソナページは独立ページとして扱い、`MASTER_PERSONA` の基盤セット要約、`MASTER_PERSONA_ENTRY` の一覧、選択中詳細、`FOUNDATION_PHASE_RUN` の最新実行状態を同一画面で観測できるようにする。
  - CRUD の主語は `MASTER_PERSONA_ENTRY` とし、`MASTER_PERSONA` はページ上部の基盤セット情報と `JSONから再構築` の対象として扱う。
  - 一覧行は `npc_name`、`race / sex / voice`、`persona_text` 要約、`npc_form_id` を最低表示粒度とし、詳細では `npc_name`、`npc_form_id`、`race`、`sex`、`voice`、`persona_text`、`MASTER_PERSONA.persona_name`、`source_type`、`built_at`、最新 `FOUNDATION_PHASE_RUN` の `phase_code / status / started_at / finished_at` を観測できるようにする。
  - create / update の最低編集項目は `npc_form_id`、`npc_name`、`race`、`sex`、`voice`、`persona_text` とし、`master_persona_id`、`persona_name`、`source_type`、`built_at` は初期設計では参照専用として扱う。
  - `JSONから再構築` は `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard.esm_Export.json` を concrete example とし、`ファイル選択 -> 再構築待ち -> 生成中 -> 完了` の staged 操作と、完了後の same-page refresh を設計前提として扱う。
- `in_scope`:
  - マスターペルソナページへ独立して到達できること。
  - 基盤セット情報、一覧、検索、選択中詳細を同一画面で確認できること。
  - 一覧ヘッダーからマスターペルソナエントリを新規作成できること。
  - 選択中エントリの詳細から更新モーダルと削除確認モーダルを開き、更新・削除できること。
  - `JSONから再構築` 導線でファイル選択、再構築待ち、生成中、完了の状態差分を確認できること。
  - 再構築完了後は検索条件を初期化し、一覧先頭へ戻し、再構築結果の代表エントリを再選択して詳細へ同期すること。
  - UI 文言とラベルは `docs/spec.md` の用語に合わせ、日本語で統一すること。
- `non_functional_requirements`:
  - `docs/architecture.md` の `View -> ScreenController -> Frontend UseCase -> GatewayContract / Store` 境界を崩さない。
  - The Ethereal Archive に従い、上端 glass header、amber / parchment 系の tonal layering、Noto Serif 系 typography、ghost border を使い、sidebar を持ち込まない。
  - 一覧、詳細、モーダル、再構築バーは同一ページ内で状態可視性を維持し、別画面遷移を要求しない。
  - マスターペルソナは基盤データとしてジョブ内ペルソナと混同せず、`MASTER_PERSONA` と `JOB_PERSONA_ENTRY` の責務分離を壊さない。
  - `FOUNDATION_PHASE_RUN` は基盤セット生成の materialization boundary として扱い、再構築導線は run 状態観測を伴う command として設計する。
- `out_of_scope`:
  - `docs/` 正本の更新、implementation-scope artifact の確定、product code の実装。
  - `JOB_PERSONA_ENTRY` を使う翻訳ジョブ実行時ペルソナの観測や編集。
  - 複数基盤セットの比較、ロールバック、履歴差分表示。
  - AI モデル選択、prompt 内容、生成アルゴリズムそのものの仕様追加。
- `open_questions`:
  - create / update form の項目粒度は ER 最小集合で十分か。`editor_id`、`class`、`source plugin`、補足メモなどの追加項目を初期実装へ含めるか。
  - `MASTER_PERSONA.persona_name` は `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard.esm_Export.json` の `target_plugin` と実行時刻から自動命名するか、再構築前に人手入力させるか。
  - `JSONから再構築` は既存セットを丸ごと置換する command として扱うか、新しい `MASTER_PERSONA` セットを生成して切り替える command として扱うか。
  - staged state は `ファイル選択 -> 再構築待ち -> 生成中 -> 完了` を基本とするが、失敗時に staged file を残したまま再実行させるか、選択解除へ戻すか。
  - `FOUNDATION_PHASE_RUN` は最新 1 件の状態だけを画面表示すれば十分か、過去 run の履歴一覧まで初期実装へ含めるか。
  - shell route 以外に master persona 用 frontend / backend contract は確認できていないため、DTO 名称と response shape は新設前提でよいか。
- `required_reading`: `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/foundation-master-er.d2`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/foundation-master-er.svg`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/design-system-ethereal-archive.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-dashboard-and-app-shell.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-master-dictionary-management.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts`, `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard.esm_Export.json`

## Artifacts

- `ui_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.ui.html
- `final_mock_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-persona/index.html
- `scenario_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.scenario.md
- `final_scenario_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-persona-management.md
- `implementation_scope_artifact_path`: 未作成。implementation-scope は本段階で実行しない。
- `review_diff_diagrams`: なし
- `source_diagram_targets`: なし
- `canonicalization_targets`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-persona/index.html`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-persona-management.md`

## Work Brief

- `implementation_target`: master-persona-management design bundle
- `accepted_scope`: requirements、ui-mock、scenario、implementation-brief の固定まで。implementation-scope と product code は human review 後にのみ扱う。
- `ordered_scope`:
  1. `master-persona` route 配下に独立ページ screen controller / gateway contract を置く前提を固定し、View から Wails 境界を直接参照しない。
  2. query は `MASTER_PERSONA` summary、`MASTER_PERSONA_ENTRY` list-detail、最新 `FOUNDATION_PHASE_RUN` status に分け、command は entry CRUD と set-level `JSONから再構築` に分ける。
  3. `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Dawnguard.esm_Export.json` を concrete example とし、`ファイル選択 -> 再構築待ち -> 生成中 -> 完了` の staged state、失敗時の戻り先候補、same-page refresh を screen state 責務として整理する。
  4. scenario では一覧検索、entry CRUD、run 状態観測、再構築完了後の再同期に加え、ファイル未選択ゲート、JSON 不正時の失敗系、責務境界維持を検証対象へ渡す。
- `design_review_questions`:
  - create / update form は ER 最小項目で十分か。追加の属性項目を初期実装へ含めるべきか。
  - `persona_name` の命名方法を人手入力にするか自動命名にするか。
  - `JSONから再構築` は既存セット置換か、新規セット生成か。
  - `FOUNDATION_PHASE_RUN` を最新 1 件表示に留めるか履歴表示まで広げるか。
  - master persona 用 contract 名称と DTO shape を新設前提で進めてよいか。
- `diagram_need`: なし。現行 `docs/architecture.md` と `foundation-master-er.d2` で境界説明は足りる。構造主語や依存方向の追加が確定した時だけ後続 phase で review diff を起こす。
- `implementation_scope_handoff`: `JSONから再構築` は set-level command、create / update / delete は entry-level command、一覧 / 詳細 / run 状態は query として分離する。`JOB_PERSONA_ENTRY` を画面主語へ混ぜない。set 置換か新規生成かは human 判断待ちの open question として残す。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `npm --prefix frontend run check`, `python3 scripts/harness/run.py --suite execution`

## Investigation

- `reproduction_status`: not-applicable
- `trace_hypotheses`: 
- `observation_points`: 
- `residual_risks`: latest `FOUNDATION_PHASE_RUN` だけを見せるか履歴まで出すかは未確定。

## Acceptance Checks

- design bundle に usecase の completion criteria が漏れなく反映される。
- UI モックで一覧、検索または絞り込み、詳細、作成、更新、削除、再構築導線の全体像を観測できる。
- scenario artifact で manual check steps と主要例外系を観測可能にする。

## Required Evidence

- structure harness pass 記録
- distill 結果
- design bundle 一式
- human review gate 記録

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: design bundle 完了後に記録する

## Closeout Notes

- `canonicalized_artifacts`: 未確定

## Outcome

- design bundle 生成待ち
