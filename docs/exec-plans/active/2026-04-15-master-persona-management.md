# Work Plan

- workflow: work
- status: in_progress
- lane_owner: orchestrate
- scope: master-persona-management
- task_id: persona-management
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml
- parent_phase: implementation-lane

## Request Summary

- マスターペルソナページの task-local design を、承認済み判断と current UI mock に同期する。
- obsolete な create / rebuild / overwrite 前提を除去し、AI生成中心の task-local artifact に揃える。
- design bundle 完了後の reviewback として、work plan、scenario、implementation-scope の 3 文書だけを更新する。
- docs canon は更新せず、implementation handoff に必要な task-local design だけを修正する。

## Decision Basis

- current UI mock `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.ui.html` では、上段に AI設定、JSON preview、run panel、下段に persona list と detail を置く same-page 構成が確定している。
- current UI mock では、一覧は生成状態一覧ではなくペルソナ一覧であり、検索欄の右側 dropdown は plugin grouping / filter 専用である。
- current UI mock では、詳細に `ダイアログ数`、`ダイアログ一覧` button、closed-by-default の dialogue modal が含まれている。
- 承認済み判断により、入力 JSON は `extractData.pas` 前提とし、中心概念は `JSONから再構築` ではなく `AI生成` とする。
- 承認済み判断により、手動新規作成は不要であり、既存 entry は `target_plugin + form_id + record_type` で判定して skip し、overwrite しない。
- 承認済み判断により、発話 0 件は skip し、`race / sex` 欠落時は一時基底 `敬語なしで中性的` を使って生成継続するが、欠落属性を persona-facing descriptor として UI に露出しない。
- 承認済み判断により、生成中の detail action は lock され、文言は `更新と削除を行えません`、生成終了後は `更新と削除を行えます` に揃える。
- 承認済み判断により、`ui-check`、system test、E2E は paid な real AI API を呼ばない。test 環境は fake AI provider / fake generation result を既定とし、保存済み API key が存在しても test mode では real provider を拒否し、`ui-check` / E2E 実行に保存済み API key を要求しない。
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`、`/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md`、`/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md` は、AI 生成 run の観測、`MASTER_PERSONA` と `JOB_PERSONA_ENTRY` の分離、frontend の依存方向を task-local design の制約として与える。
- `python3 scripts/harness/run.py --suite structure` は 2026-04-15 に再実行で pass した。

## Task Mode

- `task_mode`: `implement`
- `goal`: マスターペルソナページの requirements、ui-mock、scenario、implementation-brief、implementation-scope を task-local artifact として固定し、承認済み scope を実装へ渡す。
- `constraints`: orchestrate 自身では product 実装をしない。`docs/` 正本は human 先行でのみ更新する。承認済み task-local design を implementation-scope に圧縮して handoff する。
- `close_conditions`: implementation-review と ui-check を通し、最終 close では review 記録と canonicalization target の扱いを確認する。

## Facts

- `frontend/src/ui/stores/shell-state.ts` には `master-persona` route があり、現状は shell 上の placeholder 導線だけが存在する。
- `frontend/src/ui/views/AppShell.svelte` では `master-dictionary` だけが実ページ描画対象であり、マスターペルソナは未接続である。
- repo には master persona 専用の product code がまだ存在しない。
- task-local artifact として `2026-04-15-master-persona-management.ui.html`、`2026-04-15-master-persona-management.scenario.md`、`2026-04-15-master-persona-management.implementation-scope.md` を保持している。
- design-review は pass した。
- human review は完了しており、current turn では approved decisions と UI mock を task-local source of truth として扱う。
- `spec.md` の `マスターペルソナ = ベースゲーム由来` と current task-local decision の `任意 plugin` は未整合であり、current turn では task-local approved design として扱う必要がある。

## Functional Requirements

- `summary`:
  - マスターペルソナページは独立ページとして扱い、一覧、検索、plugin 絞り込み、詳細、更新、削除、AI設定、extractData.pas JSON からの AI生成を同一 task の要求に含める。
  - マスターペルソナは `JOB_PERSONA_ENTRY` と同一画面や同一操作で混在させない。
  - 入力 JSON は `extractData.pas` 出力を前提とし、`target_plugin` を対象単位として扱う。
  - AI生成は `未生成のみ` を対象とし、既存 `MASTER_PERSONA_ENTRY` は `target_plugin + form_id + record_type` を既存判定キーとして上書きせず skip する。
  - 発話 0 件の NPC は `skip` し、`race / sex` 欠落 NPC は一時的な中庸ペルソナ基底 `敬語なしで中性的` を与えて生成継続する。
- `in_scope`:
  - `master-persona` route から独立ページへ到達できること。
  - 一覧は plain なペルソナ一覧として表示され、`名前`、`FormID`、`EditorID`、`種族`、`voice`、`ダイアログ数`、`収録先 plugin`、`ペルソナ要約` を観測できること。
  - 一覧上で `名前`、`FormID`、`EditorID`、`種族`、`voice` を対象に検索できること。
  - 検索欄の右側 dropdown は plugin grouping / filter 専用であり、generation status filter を置かないこと。
  - 一覧は generation target state ではなく persona summary を主表示にすること。
  - 選択中 entry の詳細で `FormID`、`EditorID`、`名前`、`voice`、`class`、`source`、`ペルソナ本文`、`ダイアログ数` を観測できること。
  - 詳細には `ダイアログ一覧` button があり、dialogue modal は closed-by-default で、必要時だけ開閉できること。
  - `更新` は `FormID`、`EditorID`、`名前`、`種族`、`性別`、`voice`、`class`、`source`、`ペルソナ本文` を編集対象に含めること。
  - `削除` は確認操作を経て実行でき、成功後は同一ページ内で一覧、選択状態、詳細表示が同期して更新されること。
  - `race / sex` が未設定で利用可能な名前属性も弱い場合でも、欠落属性を persona-facing descriptor として露出しないこと。
  - AI設定はこのページ内で保存でき、少なくとも `provider`、`model`、`APIKey` を扱えること。
  - prompt template はこのページの編集項目に含めず、将来 human が調整しやすい定数として実装側で保持すること。
  - `extractData.pas` JSON はファイル選択 UI から開始し、選択後に対象ファイル名、`target_plugin`、総 NPC 数、生成対象数、既存 skip 数、発話 0 件 skip 数、`汎用NPC` 数を同一ページで観測できること。
  - AI生成開始前に preview を確認でき、`作成済みのペルソナはスキップされます` を明示できること。
  - AI生成 run では `設定未完了`、`入力待ち`、`入力検証中`、`入力エラー`、`対象なし`、`生成可能`、`生成中`、`中断済み`、`中止済み`、`回復可能失敗`、`完了`、`失敗` を UI から追跡できること。
  - AI生成 run の進行中、`更新 / 削除` は read-only に固定され、`更新と削除を行えません` と表示されること。完了後と失敗後は `更新と削除を行えます` へ戻り、同一ページ内で一覧、詳細、現在の選択状態が破綻なく再同期されること。
- `non_functional_requirements`:
  - マスターペルソナページは、任意 plugin 単位の extractData.pas JSON 件数でも、一覧、検索、plugin 絞り込み、選択、詳細確認が継続して操作できること。
  - AI生成 run は長時間処理であっても、進捗件数、現在処理中 NPC、phase run / AI run の状態を UI から追跡できること。
  - dialogue modal は初期非表示を守り、開閉で一覧や詳細の state を壊さないこと。
  - `ui-check`、system test、E2E は paid な real AI API を呼ばず、fake AI provider と fake generation result だけで安全に完走できること。
  - test mode では保存済み API key の有無に関係なく real provider を選択不可または実行拒否とし、構造的に paid provider path を塞ぐこと。
  - `ui-check` / E2E は browser surface の検証だけで成立し、保存済み API key や real provider 接続を前提にしないこと。
  - UI 文言は `docs/spec.md` の用語と current approved task-local decision に合わせ、日本語で一貫していること。
  - frontend の画面責務は `View -> ScreenController -> Frontend UseCase -> GatewayContract/Store` の依存方向を崩さないこと。
  - visual は `The Ethereal Archive` に従い、上端 glass header、amber / parchment 系 tonal layering、`Noto Serif`、ghost border、sidebar 非採用を維持すること。
- `out_of_scope`:
  - `JOB_PERSONA_ENTRY` と job 単位ペルソナの参照、編集、削除導線。
  - 手動新規作成導線。
  - 既存 `MASTER_PERSONA_ENTRY` への overwrite / rollback / diff preview。
  - JSON 再構築や JSON からの全件再作成を前提にした workflow。
  - 複数 JSON を束ねた scheduler と履歴比較。
  - `docs/` 正本の恒久仕様変更。
- `open_questions`: なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.ui.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.scenario.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.implementation-scope.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/persona-management.yaml`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/input-data-er.d2`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/foundation-master-er.d2`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/translation-job-er.d2`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/design-system-ethereal-archive.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-dashboard-and-app-shell.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-master-dictionary-management.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`

## Artifacts

- `ui_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.ui.html
- `final_mock_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-persona/index.html
- `scenario_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.scenario.md
- `final_scenario_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-persona-management.md
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.implementation-scope.md
- `review_diff_diagrams`: なし
- `source_diagram_targets`: なし
- `canonicalization_targets`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-persona/index.html, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-persona-management.md

## Work Brief

- `implementation_target`: マスターペルソナページ実装
- `accepted_scope`: implementation-scope、product code、tests、review までを扱う。docs 正本更新はこの turn では行わない。
- `parallel_task_groups`: frontend、backend-contract、tests、review で分ける。
- `tasks`:
  1. shell placeholder を実ページへ置き換える。
  2. 一覧、plugin filter、詳細、dialogue modal、AI設定パネル、JSON preview、run status の same-page state を frontend に実装する。
  3. persona page 用の Wails contract と usecase 境界を作り、View から直接 transport へ触れない構成にする。
  4. `未生成のみ`、`発話0件 skip`、`汎用NPC` 集計、`target_plugin + form_id + record_type` 既存判定、run 中 read-only の preview / command 契約を組み込む。
  5. DI 可能な fake AI provider / fake generation result と test-safe な secret handling seam を入れ、test mode では real provider を backend で拒否する。
  6. scenario と UI mock の主要導線を tests と review で通し、`ui-check` / E2E が保存済み API key なしで fake provider path を安全に使えることを確認する。
- `implementation_brief_background`:
  - 現状の shell には `master-persona` route があるが、画面本体と controller wiring は未実装である。
  - current task-local design は CRUD + JSON再構築 から AI生成前提へ切り替わっており、一覧と詳細に加えて AI設定、plugin filter、preview、run 観測、dialogue modal が必要である。
  - extractData.pas JSON は `target_plugin` 単位で扱い、既存 entry は `target_plugin + form_id + record_type` を既存判定キーとして overwrite せず skip する前提へ変わった。
- `implementation_brief_recommendation`:
  - top glass header 配下に AI設定パネルと JSON preview / run card を置き、その下を一覧と詳細の 2 カラムで構成する。
  - 一覧は keyword search と plugin dropdown を持つ plain な persona list とし、persona summary を主表示にする。
  - 詳細は選択中 entry の主要属性、`ダイアログ数`、`ダイアログ一覧` button、ペルソナ本文を表示する。dialogue modal は closed-by-default とする。
  - 手動新規作成は廃止し、手動操作は update / delete のみを残す。run 中の update / delete は read-only に固定する。
  - AI設定は page-local に保存し、preview では `総NPC数`、`生成対象数`、`既存 skip 数`、`発話0件 skip 数`、`汎用NPC 数` を可視化する。
  - API key 参照は `SecretStore` または同等 seam 越しに扱い、test mode では fake provider / fake generation result を既定として保存済み API key を要求しない。
  - prompt template は画面入力ではなく実装側の定数として保持し、後から human が調整しやすい配置にする。
  - `race / sex` 欠落時の中庸ペルソナ基底は、当面 `敬語なしで中性的` を採用し、後から tuning 可能な前提で扱う。ただし UI では欠落属性を persona-facing descriptor として露出しない。
  - run card は `設定未完了`、`入力待ち`、`入力検証中`、`入力エラー`、`対象なし`、`生成可能`、`生成中`、`中断済み`、`中止済み`、`回復可能失敗`、`完了`、`失敗` を可視化する。
  - 完了時は overwrite warning ではなく skip 内訳を返し、一覧は対象 plugin を基準に結果確認しやすくする。
- `implementation_brief_unresolved_items`: なし
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## Investigation

- `reproduction_status`: not-applicable
- `trace_hypotheses`: なし
- `observation_points`: なし
- `residual_risks`: `spec.md` の `ベースゲーム由来` と current task-local decision の `任意 plugin` は未整合である。current turn では task-local approved design として扱い、将来の docs 正本同期が必要である。

## Acceptance Checks

- design bundle に AI生成前提の completion criteria と human review 承認が反映される。
- implementation-scope が frontend / backend / tests / review の ownership を曖昧さなく固定する。
- `ui-check`、system test、E2E が fake provider path と fake generation result を既定にし、paid な real AI API を呼ばない safety design が requirements と implementation-scope に反映される。
- `未生成のみ`、`発話0件 skip`、`target_plugin + form_id + record_type` 既存判定、dialogue count、dialogue modal、plugin filter、run 中 read-only の前提が scenario と implementation-scope に反映される。
- obsolete な create / rebuild / overwrite 前提が task-local artifact から除去される。

## Required Evidence

- structure harness pass 記録
- distill 結果
- requirements 更新済み active plan
- `2026-04-15-master-persona-management.ui.html`
- `2026-04-15-master-persona-management.scenario.md`
- `2026-04-15-master-persona-management.implementation-scope.md`
- design-review pass 記録
- human approval record

## HITL Status

- `functional_or_design_hitl`: `approved`
- `approval_record`: `2026-04-15 human review approved: AI generation replaces rebuild / extractData.pas JSON / arbitrary plugin / page-local AI settings / only not-yet-generated entries / existing identity key=target_plugin+form_id+record_type / no manual create / zero-dialogue skip / missing race-sex uses temporary neutral baseline=敬語なしで中性的 / prompt template is constantized not page-editable / update-delete are read-only during run / persona list uses plugin filter not generation-status filter / detail shows dialogue count and closed-by-default dialogue modal / missing race-sex are not surfaced as persona-facing descriptors`

## Closeout Notes

- `canonicalized_artifacts`: 未確定

## Outcome

- active plan を current approved decisions と current UI mock に合わせて更新した。
- requirements、scenario、implementation-scope の obsolete な create / rebuild / overwrite 前提を除去した。
- 一覧の plugin filter、detail の dialogue modal、AI生成 wording と lock wording を task-local design に固定した。
