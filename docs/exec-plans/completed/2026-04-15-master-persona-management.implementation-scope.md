# 実装スコープ固定

- `task_id`: `persona-management`
- `task_mode`: `implement`
- `design_review_status`: `pass`
- `hitl_status`: `approved`
- `source_brief`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.md`
- `summary`: `マスターペルソナ画面を placeholder ではなく実ページとして実装する。中心概念は JSON 再構築ではなく AI生成とし、extractData.pas JSON を入力に、任意 plugin を対象に、未生成 entry のみを追加する。既存 entry は target_plugin + form_id + record_type で判定して skip し、overwrite しない。手動新規作成は置かず、run 中の 更新 / 削除 は read-only に固定する。加えて、ui-check / system test / E2E は paid な real AI API を呼ばず、fake provider path を既定とする。`

## 共通ルール

- 中心概念は `AI生成` とする。`create`、`rebuild`、`overwrite` を主導線として実装しない。
- 入力は `extractData.pas` JSON のみを受け付ける。plugin 名による制限は入れず、任意 `target_plugin` を扱う。
- 既存 entry 判定キーは必ず `target_plugin + form_id + record_type` とする。別キーや UI 補助情報で代替しない。
- AI生成対象は `まだ MASTER_PERSONA_ENTRY が存在しないものだけ` とする。既存 entry は preview と execution の両方で skip し、overwrite しない。
- `発話 0 件` は必ず skip する。preview と execution の両方で件数をそろえる。
- `race` または `sex` が欠落している場合は生成を止めず、一時基底 `敬語なしで中性的` を適用して継続する。この内部基底を欠落属性の代替 descriptor として UI 表示へ流し込まない。
- 手動新規作成 UI、create usecase、create API、create repository command は実装しない。
- 更新と削除は通常時のみ許可する。active run 中は frontend で disabled 表示にし、backend でも reject して二重に守る。
- prompt template は UI 編集項目と保存 DTO に含めない。backend の master persona generation service 直下に定数として置き、この task では画面から変更できない。
- AI 設定はこのページ専用とする。translation job や他画面の設定ストアと共有しない。
- API key 参照と保存済み secret 取得は `SecretStore` または同等 seam 越しに限定する。test mode ではこの seam が real key 利用を許可しない構造にする。
- test 環境の既定は `real key を使わない` とする。保存済み API key が存在しても、ui-check / system test / E2E / unit test / integration test は fake provider / fake generation result を使う。
- backend は test mode で `RealAIProvider` を明示的に reject する。frontend 側の選択制御だけに依存しない。
- `MASTER_PERSONA` と `JOB_PERSONA_ENTRY` は query、DTO、route、文言、テストの全てで混在させない。
- 同一ページ更新を守る。settings 保存、preview、run 完了、run 失敗、update、delete の各完了後に route 遷移なしで list / detail / status を再同期する。
- 一覧の右側 dropdown は plugin grouping / filter 専用とし、generation status filter を追加しない。
- docs 正本更新はこの task に含めない。task-local artifact と product code と test と review だけを扱う。

## データ契約

- list query は `keyword`、`plugin_filter`、`page`、`page_size` を受ける。`generation_state` のような一覧表示前提の state filter は受けない。
- list response は `items` と `plugin_groups` を返す。`plugin_groups` は dropdown 表示用の plugin 名と件数を含む。
- list item DTO は最低でも `target_plugin`、`form_id`、`record_type`、`editor_id`、`display_name`、`race`、`sex`、`voice_type`、`class_name`、`source_plugin`、`persona_summary`、`dialogue_count`、`updated_at` を返す。
- list item DTO は生成状態を一覧主表示にしない。必要な run 状態は run panel 用の別 DTO で返す。
- detail DTO は list item DTO の情報に加えて `persona_body`、`generation_source_json`、`identity_key`、`baseline_applied`、`run_lock_reason` を返す。
- detail DTO の `race` と `sex` は nullable を許容する。frontend presenter は空値を `不明` や疑似 label に変換せず、persona-facing descriptor から除外する。
- dialogue list DTO を追加し、`identity_key` または list / detail と同値の識別子に対して `dialogue_count`、`dialogues[]` を返す。`dialogues[]` の各要素は `index` と `text` を持つ。
- page AI settings DTO は `provider`、`model`、`api_key` を扱う。prompt template は DTO に含めない。
- page AI settings DTO は test mode で空の `api_key` を許容し、保存済み API key がなくても fake provider path の preview / run を開始できるようにする。
- preview request は `extractData.pas` JSON ファイル入力と page-local AI settings を受ける。preview response は最低でも `target_plugin`、`total_npc_count`、`generatable_count`、`existing_skip_count`、`zero_dialogue_skip_count`、`generic_npc_count`、`status` を返す。
- preview response の `status` は `設定未完了`、`入力待ち`、`入力検証中`、`入力エラー`、`対象なし`、`生成可能` の少なくとも 6 状態を返す。
- run status / execute response は `run_state`、`processed_count`、`success_count`、`existing_skip_count`、`zero_dialogue_skip_count`、`generic_npc_count`、`current_actor_label`、`target_plugin` を返す。
- preview / run の UI 文言は `作成済みのペルソナはスキップされます`、`汎用NPC`、`更新と削除を行えます / 行えません` を task-local mock に合わせて返せる構造にする。
- run state は `生成中`、`中断済み`、`中止済み`、`回復可能失敗`、`完了`、`失敗` を最低限扱う。UI 表示文言は task-local mock に合わせる。
- execution は preview 結果を前提にしてよいが、backend は実行直前に既存判定を再評価し、overwrite しないことを保証する。
- update / delete API は active run 中に呼ばれた場合、明示的な domain error を返す。

## 許可ファイルと handoff

### `backend-master-persona-ai-generation`

- `implementation_target`: `backend`
- `owned_scope`:
  - `internal/controller/wails/` の master persona controller / DTO / mapper 一式
  - `internal/usecase/` の master persona list / detail / dialogue list / settings / preview / execute / update / delete usecase
  - `internal/service/` の master persona query service、AI generation service、run status service、prompt template constant、`SecretStore` または同等 seam を使う provider resolution rule
  - `internal/repository/` の master persona query / command / AI settings persistence / run persistence / test-safe secret access adapter
  - `internal/bootstrap/app_controller.go`
  - `internal/bootstrap/app_controller_test.go`
- `depends_on`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.ui.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-persona-management.scenario.md`
- `validation_commands`:
  - `go test ./internal/controller/wails ./internal/usecase ./internal/service ./internal/repository ./internal/bootstrap`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: Wails 側 API と usecase が、list / detail / dialogue list / page-local AI settings load-save / extractData.pas preview / AI execute / run status polling / interrupt / cancel / update / delete を備える。create API は存在しない。execution は既存 entry を overwrite せず、identity key と skip ルールを backend 単体で守る。test mode では保存済み API key の有無に関係なく `RealAIProvider` を reject し、fake provider path だけを通す。
- `notes`:
  - prompt template constant は backend service 配下の専用ファイルへ分離し、settings 保存対象に混ぜない。
  - `SecretStore` または同等 seam を介さずに real API key へ触れる実装を置かない。test mode の seam は fake 値または拒否だけを返す。
  - preview と execute で `existing_skip_count`、`zero_dialogue_skip_count`、`generic_npc_count` の定義をそろえる。
  - run active 中に update / delete を reject する domain rule を service 層で持つ。
  - dialogue list は detail payload へ全件同梱せず、必要時取得でもよいが、dialogue count は detail と list の両方で参照できるようにする。

### `frontend-master-persona-page`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/application/contract/` の master persona contract
  - `frontend/src/application/gateway-contract/` の master persona gateway contract
  - `frontend/src/application/store/` の master persona page-local store
  - `frontend/src/application/presenter/` の master persona presenter
  - `frontend/src/application/usecase/` の master persona usecase
  - `frontend/src/controller/master-persona/`
  - `frontend/src/controller/runtime/master-persona/`
  - `frontend/src/controller/wails/` の master persona gateway / DTO mapping
  - `frontend/src/ui/screens/master-persona/`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/App.svelte`
- `depends_on`:
  - `backend-master-persona-ai-generation`
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
- `completion_signal`: shell route が実ページを描画し、hero、AI settings、JSON preview、run monitor、list、detail の 6 ブロックが task-local mock に沿って並ぶ。検索欄の右側には plugin dropdown のみを置く。create button / create modal は存在しない。detail には `ダイアログ数` と `ダイアログ一覧` button があり、dialogue modal は閉じた初期状態から開閉できる。run 中は `更新` と `削除` を disabled 表示し、`更新と削除を行えません` を出す。run 完了後は `更新と削除を行えます` を出す。prompt template は UI 入力欄ではなく説明だけを出す。test env / ui-check / E2E では fake provider path を既定表示として扱い、保存済み API key なしで画面を操作できる。
- `notes`:
  - page-local AI settings は master persona 画面の store と controller だけで保持し、他画面の AI settings UI と共有しない。
  - ui-check / E2E 用の runtime と fixture は fake provider を既定にし、保存済み real key の有無に依存させない。
  - preview card に `target_plugin`、`total`、`generatable`、`existing skip`、`zero-dialogue skip`、`汎用NPC` を表示する。
  - list は persona summary を主表示にし、generation target state を list row 主表示へ混ぜない。
  - presenter は `race / sex` 欠落時に `不明` や推測ラベルを出さず、利用可能な name / voice / dialogue count などだけで表示を構成する。
  - run 完了後は route 遷移せず、対象 plugin の最新 list と detail を再取得して画面を更新する。

### `tests-master-persona-generation-rules`

- `implementation_target`: `tests`
- `owned_scope`:
  - `frontend/src/controller/master-persona/*.test.ts`
  - `frontend/src/application/presenter/master-persona/*.test.ts`
  - `frontend/src/application/usecase/master-persona/*.test.ts`
  - `frontend/src/ui/App.test.ts`
  - `internal/controller/wails/` master persona tests
  - `internal/usecase/` master persona tests
  - `internal/service/` master persona tests
  - `internal/repository/` master persona tests
- `depends_on`:
  - `backend-master-persona-ai-generation`
  - `frontend-master-persona-page`
- `validation_commands`:
  - `go test ./internal/...`
  - `npm --prefix frontend run test -- --runInBand`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: テストが次を証明する。1) create 導線が存在しない。2) identity key が `target_plugin + form_id + record_type` である。3) existing entry は preview / execute の両方で overwrite されず skip される。4) zero-dialogue は skip される。5) race / sex 欠落時に `敬語なしで中性的` が適用されても、UI descriptor へ欠落属性が露出しない。6) dialogue count と dialogue modal 用データが取得できる。7) run active 中は update / delete が frontend と backend の両方で lock される。8) plugin dropdown が plugin filter としてのみ機能する。9) prompt template が UI editable DTO に出ない。10) page-local AI settings の load / save が他画面へ漏れない。11) test mode では保存済み API key が存在しても `RealAIProvider` が reject される。12) ui-check / system test / E2E は fake provider と fake run / event result だけで成立し、saved API key を要求しない。
- `notes`:
  - UI test は task-local mock の主要文言と disabled state を確認する。
  - backend test は execute 直前の再判定で overwrite を防ぐケースを 1 件以上含める。
  - fake provider、fake generation result、fake run event feed を DI で差し替える経路を既定にし、paid provider を呼ぶ test を置かない。
  - dialogue modal test は初期 closed、open、close を最低 1 ケース含める。

### `review-master-persona-ai-page`

- `implementation_target`: `review`
- `owned_scope`:
  - `implementation-review` record
  - `ui-check` record
- `depends_on`:
  - `backend-master-persona-ai-generation`
  - `frontend-master-persona-page`
  - `tests-master-persona-generation-rules`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite all`
  - `npm run dev:wails:docker-mcp`
- `completion_signal`: implementation-review と ui-check が pass し、承認済み判断とのズレが残らない。特に AI生成中心、extractData.pas JSON、任意 plugin、未生成のみ、overwrite なし、create なし、plugin filter、dialogue modal、run 中 read-only、prompt 定数化に加え、fake provider path 既定、test mode の `RealAIProvider` reject、saved API key 不要の safety design が一致している。
- `notes`:
  - Playwright MCP の接続先は `http://host.docker.internal:34115` を使う。
  - ui-check / Playwright は browser surface の検証だけを行い、fake provider path の preview / run result を使って確認する。paid provider 接続の成否確認は scope 外とする。
  - mock と実装の差分は UI 文言ではなく振る舞い優先で確認する。

## 実装順

1. backend に master persona の settings / preview / execute / run status / list / detail / dialogue list / update / delete 契約を作る。create 契約は追加しない。
2. backend service に identity key 判定、zero-dialogue skip、generic NPC 集計、no-overwrite 再判定、run lock を実装する。
3. frontend に page-local store と screen controller を作り、route placeholder を task-local mock 準拠の実ページへ置き換える。
4. list / detail / dialogue modal / preview / run monitor の same-page refresh を通し、plugin filter と run active 中の disabled state を反映する。
5. scenario に沿う test と review で AI生成ルールを固定する。

## 明示的な非目標

- 手動新規作成フローの追加
- prompt template の UI 編集、永続化、外部設定化
- 既存 entry の overwrite、diff preview、rollback
- generation status dropdown の追加
- `JOB_PERSONA_ENTRY` 側の画面や API の変更
- `docs/` 正本の更新
