# Scenario Design: 2026-04-19-legacy-schema-ui-migration-todo

- `skill`: scenario-design
- `status`: draft
- `source_plan`: `./plan.md`
- `requirements_source`: `./requirements-design.md`
- `ui_source`: `./legacy-schema-ui-migration.ui.html`
- `ui_mock_source`: `./legacy-schema-ui-migration.ui.html`
- `diagram_source`: `./legacy-schema-ui-migration.review-er-diff.puml`
- `final_artifact_path`: `docs/scenario-tests/legacy-schema-ui-migration.md`
- `topic_abbrev`: `LSU`

## Rules

- ケース ID は `SCN-LSU-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックにする。
- 期待結果は repository / service / frontend test で観測できる内容にする。
- paid な real AI API を前提にしない。

## Scenario Matrix

### SCN-LSU-001 Dictionary list / detail reads canonical data

- `分類`: 正常系
- `観点`: 既存 master dictionary 画面が canonical `DICTIONARY_ENTRY` 由来の一覧と詳細を表示する。
- `事前条件`: canonical table に shared dictionary entry が存在する。
- `手順`:
  1. master dictionary page を開く。
  2. 一覧、検索、カテゴリ絞り込み、詳細表示を確認する。
  3. DB 上の read 元 table を test で確認する。
- `期待結果`:
  1. 一覧と詳細は canonical cutover 後の縮小 UI contract で表示される。
  2. `source_term` / `translated_term` が原文 / 訳語として表示される。
  3. read は旧 `master_dictionary_entries` ではなく `DICTIONARY_ENTRY` から行われる。
  4. `REC` / `EDID` / `note` は詳細 UI と frontend / Wails contract に残らない。
- `観測点`: service test、controller test、frontend gateway / presenter test。
- `fake_or_stub`: temp DB。
- `責務境界メモ`: visible redesign は要求しない。

### SCN-LSU-002 Dictionary mutations write canonical data

- `分類`: 正常系
- `観点`: create / update / delete が canonical dictionary table に反映される。
- `事前条件`: master dictionary page が操作できる。
- `手順`:
  1. 新規登録で辞書 entry を作る。
  2. 詳細から更新する。
  3. 詳細から削除する。
  4. canonical table と旧 table の row 変化を確認する。
- `期待結果`:
  1. create / update / delete は `DICTIONARY_ENTRY` に反映される。
  2. 旧 `master_dictionary_entries` への新規 write は発生しない。
  3. 原文は保存前に trim され、重複判定は `trim(source_term) + translated_term` の完全一致になる。
  4. frontend の手動新規登録だけ、訳語の前後にある明らかな入力ノイズを除去できる。
  5. 同一ページ内の一覧と詳細が既存通り再同期される。
- `観測点`: repository / service integration test、frontend state test。
- `fake_or_stub`: temp DB。
- `責務境界メモ`: UI contract 互換は service adapter で吸収する。

### SCN-LSU-003 XML import writes canonical provenance and entries

- `分類`: 正常系
- `観点`: XML import が `XTRANSLATOR_TRANSLATION_XML` と `DICTIONARY_ENTRY` に保存される。
- `事前条件`: import 対象 XML fixture がある。
- `手順`:
  1. master dictionary page で XML を選択する。
  2. import を実行する。
  3. import summary、一覧再同期、DB row を確認する。
- `期待結果`:
  1. XML file provenance は `XTRANSLATOR_TRANSLATION_XML` に保存される。
  2. import された term は `DICTIONARY_ENTRY` に保存される。
  3. 旧 `master_dictionary_entries` への import write は発生しない。
  4. import でも原文は trim され、重複判定は `trim(source_term) + translated_term` の完全一致になる。
  5. XML import の訳語は全体正規化されず、XML 由来の訳語として保存される。
  6. `REC` / `EDID` / `selectedRec` は XML parse 中の一時情報としてのみ使われ、import 後の UI contract と canonical table に残らない。
- `観測点`: import service test、sqlite integration test、frontend runtime event test。
- `fake_or_stub`: small XML fixture。
- `責務境界メモ`: XML parser の仕様変更は含めない。

### SCN-LSU-004 Persona list / detail reads canonical join model

- `分類`: 正常系
- `観点`: 既存 master persona 画面が `PERSONA` + `NPC_PROFILE` 由来の一覧と詳細を表示する。
- `事前条件`: canonical persona と npc profile が存在する。
- `手順`:
  1. master persona page を開く。
  2. 一覧、plugin filter、詳細を確認する。
  3. DB 上の read 元 table を test で確認する。
- `期待結果`:
  1. persona list / detail は canonical cutover 後の縮小 UI contract で表示される。
  2. identity は `NPC_PROFILE` の target plugin / form id / record type から組み立てられる。
  3. persona 本体は `PERSONA` から読み込まれる。
  4. 旧 `master_persona_entries` への read 依存は残らない。
  5. `generation_source_json`、`baseline_applied`、dialogue 表示、dialogue modal、dialogue count 表示は残らない。
  6. 会話が見つからない NPC は入力 JSON parse 時点で除外され、persona list / detail には出ない。
- `観測点`: service test、controller test、frontend gateway / presenter test。
- `fake_or_stub`: temp DB。
- `責務境界メモ`: dialogue modal は canonical cutover 後の保持対象にしない。削除後の list / detail 画面は HTML mock で review する。

### SCN-LSU-005 Persona generation writes canonical persona data

- `分類`: 正常系
- `観点`: JSON preview / execute が canonical persona data を作成し、既存判定も canonical identity で行う。
- `事前条件`: fake AI transport と extractData.pas JSON fixture がある。
- `手順`:
  1. 会話がある NPC と会話がない NPC を含む JSON を選択して preview を実行する。
  2. AI generation を fake transport で実行する。
  3. 生成中の失敗を fake transport で発生させる。
  4. 作成件数、既存件数、canonical table を確認する。
- `期待結果`:
  1. AI 生成物が揃った NPC だけ、`NPC_PROFILE` と `PERSONA` が同一 transaction で作成される。
  2. `PERSONA` は `NPC_PROFILE` と 1:1 で作成される。
  3. 既存判定は旧 `identity_key` table ではなく、`target_plugin_name + form_id + record_type` の canonical identity で行われる。
  4. 旧 `master_persona_entries` への新規 write は発生しない。
  5. 会話が見つからない NPC は preview / generation target から除外され、persona row は作られない。
  6. preview / frontend / Wails contract は、会話なし NPC や generic NPC の skip count を返さず、ペルソナ候補数、新規追加可能数、作成済み数を返す。
  7. fake transport 失敗時も、生成途中の `NPC_PROFILE` / `PERSONA` は残らない。
- `観測点`: service integration test、repository test、frontend usecase test。
- `fake_or_stub`: fake AI transport、JSON fixture、temp DB。
- `責務境界メモ`: paid real AI API は使わない。

### SCN-LSU-006 AI settings persist and run state stays transient

- `分類`: 状態遷移
- `観点`: AI settings が foundation data と混ざらず、run state は DB に保存されず、再起動後は JSON 未選択状態から手動読込へ戻せる。
- `事前条件`: `PERSONA_GENERATION_SETTINGS` が存在し、`id = 1` の singleton row として扱える。
- `手順`:
  1. AI settings を保存する。
  2. 一部 persona が canonical `PERSONA` に保存された状態を作る。
  3. app wiring を作り直して再読込する。
  4. JSON 未読込状態から、JSON を手動で選び直して preview する。
- `期待結果`:
  1. provider / model は `PERSONA_GENERATION_SETTINGS(id = 1)` から再起動後も復元される。
  2. API key は secret store seam から読み戻され、DB の foundation table に露出しない。
  3. `PERSONA_GENERATION_RUN_STATUS` table / row は作られない。
  4. `PERSONA` / `DICTIONARY_ENTRY` に UI run state が混入しない。
  5. 再起動後は JSON 未選択状態になる。
  6. JSON を手動で読み直すと、現在の JSON から新規に追加できる NPC の有無だけが preview に表示される。
  7. `target_plugin_name` は `NPC_PROFILE` の filter key であり ownership として扱わない。
- `観測点`: repository / bootstrap test。
- `fake_or_stub`: keyring は test secret store seam。
- `責務境界メモ`: AI settings は foundation data ER から分離し、run state は画面メモリに閉じる。

### SCN-LSU-007 Legacy tables are dropped during cutover

- `分類`: 主要失敗系
- `観点`: canonical cutover が legacy data backfill を要求せず、旧 `master_*` table を drop する。
- `事前条件`: 旧 table と canonical table が同居する temp DB がある。
- `手順`:
  1. migration を適用する。
  2. canonical adapter を使って read / write する。
  3. 旧 table が DB schema から消えていることを確認する。
- `期待結果`:
  1. legacy data backfill は実行されない。
  2. 新規 read / write は canonical table へ向かう。
  3. 旧 `master_dictionary_entries`、`master_persona_entries`、`master_persona_ai_settings`、`master_persona_run_status` は drop される。
  4. `PERSONA_GENERATION_SETTINGS` が `id = 1` singleton として作られ、`PERSONA_GENERATION_RUN_STATUS` は作られない。
  5. product code に旧 table / repository 参照が残らない。
- `観測点`: migration test、service integration test。
- `fake_or_stub`: temp DB。
- `責務境界メモ`: 旧 row の保持は acceptance に含めない。

### SCN-LSU-008 Bootstrap and frontend contracts stay coherent

- `分類`: 責務境界
- `観点`: production wiring、frontend gateway contract、改修後 UI が canonical cutover 後も整合する。
- `事前条件`: canonical adapters を bootstrap へ配線済み。
- `手順`:
  1. backend controller / bootstrap tests を実行する。
  2. frontend gateway / presenter / screen controller tests を実行する。
  3. structure harness を実行する。
- `期待結果`:
  1. `internal/bootstrap/app_controller.go` は canonical adapter を配線する。
  2. frontend は縮小後の gateway contract から操作できる。
  3. 改修後 UI は HTML mock の主要表示に沿い、DB 変更で出せない表示を残さない。
  4. dependency direction と layer rule は崩れない。
- `観測点`: backend test、frontend test、structure harness。
- `fake_or_stub`: fake AI transport、temp DB。
- `責務境界メモ`: product code 実装は human-approved implementation-scope 後に行う。

## Acceptance Checks

- `SCN-LSU-001` と `SCN-LSU-002` は dictionary UI の canonical read / write acceptance を満たす。
- `SCN-LSU-003` は XML import の canonical provenance acceptance を満たす。
- `SCN-LSU-004` と `SCN-LSU-005` は persona UI / generation の canonical read / write acceptance を満たす。
- `SCN-LSU-006` は singleton の `PERSONA_GENERATION_SETTINGS` による AI settings 復元、run state 非永続化、JSON 未選択状態から手動読込へ戻す acceptance を満たす。
- `SCN-LSU-007` は legacy backfill 不要と旧 table drop の acceptance を満たす。
- `SCN-LSU-008` は bootstrap / frontend contract / UI 改修後表示の acceptance を満たす。

## Validation Commands

- `go test ./internal/infra/sqlite ./internal/repository ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap`
- `npm --prefix frontend run check`
- `npm --prefix frontend run test`
- `python3 scripts/harness/run.py --suite structure`

## Open Questions

- None. AI settings table 名は `PERSONA_GENERATION_SETTINGS`、保存粒度は `id = 1` の singleton page setting に固定する。run state は DB に保存しない。dictionary 重複判定は `trim(source_term) + translated_term` の完全一致に固定する。
