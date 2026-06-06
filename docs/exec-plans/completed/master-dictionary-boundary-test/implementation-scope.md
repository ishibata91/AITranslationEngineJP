# Implementation Scope: master-dictionary-boundary-test

- `skill`: implementation-scope
- `status`: approved
- `source_plan`: `./plan.md`
- `human_review_status`: 詳細仕様差分・本 implementation-scope ともに人間レビュー approved（2026-06-05）。
- `approval_record`: 詳細仕様差分の approval は `./detail-spec-diff.md` 内に記録。本 implementation-scope は 2026-06-05 に人間レビュー 4 観点を全て approve（うち「片側書き換え検出を手動確認に留める」は資料作成時に再検討する保留付き approval として観察項目化する）。
- `module_entry`: `.claude/skills/implementation-module/SKILL.md`
- `handoff_runtime`: `claude-module`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`（approved, Q-001..Q-008 確定済み）
- `design_diff_diagram`: `./design-diff-diagram.md`（human-review-resolved）
- `screen_design_diff`: `N/A`（plan 想定 2 N、4 N、5 N により画面変更なし）

## Fixed Decisions

### 本 task で実装範囲に含める対象

- 「境界結合テスト」層の repo 配置を新規追加し、`master-dictionary` を pilot として Wails 境界の golden 共有を実証する。
- 共有 golden 集合と、その読み出し経路（backend 側 / frontend 側）を本 task の範囲とする。
- pilot 対象 6 method（`ListMasterDictionaryEntries` / `GetMasterDictionaryEntry` / `CreateMasterDictionaryEntry` / `UpdateMasterDictionaryEntry` / `DeleteMasterDictionaryEntry` / `ImportMasterDictionaryXml`）の代表 case を境界結合テスト として固定する。
- pilot 完了後に別 task `workflow-contract-maintenance` で skill 改訂を判断するための「観察記録 docs draft」を active plan folder 内に残す。
- pilot 範囲の境界 API 仕様 draft を active plan folder 内に Markdown で作る。正本反映は後続 `finalization-module` に委ねる。

### 本 task の範囲外（明確に外す対象）

- `docs/detail-specs/boundary-integration-test.md` 正本ファイル本体の作成（detail-spec-diff の `Q-002` 反映方針 = 正本本体作成は本 task では行わない）。draft は active plan folder 内に置く。
- `docs/detail-specs/master-dictionary.md` 正本への追補（finalization-module で扱う）。
- 既存 skill（`test-design` / `tests-scenario` / `tests-unit` / `implementation-module` / `design-module`）の改訂、または新規 skill `tests-boundary` の追加（detail-spec-diff `Q-006` で別 task 委譲を確定）。
- AI service 境界 / SQL repository 境界への境界結合テスト適用（detail-spec-diff `Q-001` で将来拡張枠と明示済み）。
- `IsTermTarget` の 13 種別判定そのものの assert（detail-spec-diff `Q-007` で backend 単体テスト責務として残す）。
- 機械可読形式（OpenAPI / JSON Schema / TypeScript type）への境界 API 仕様書の昇格（detail-spec-diff `boundary-test-REQ-005` で後続検討事項）。

### 配置 path 判断

- backend 側境界結合テスト: 既存 `internal/apitest/` 配下に置く案と、新規 path（例: `internal/boundarytest/`）に切る案がある。本 task では新規 path に切らず、**`internal/apitest/` 内の独立 file 群として置く**ことを決定する。理由は次のとおり。
  - 既存 `internal/apitest/` の `README.md` は「Wails controller / bootstrap 済み controller を入口にした受け入れ条件を証明する場所」と定義しており、境界結合テストは入口を共有する。
  - `.go-arch-lint.yml` の component 定義を追加せずに済むため、本 task で arch lint 構造を触らない。
  - 観測責務（受け入れ条件 vs DTO field 値 semantic）は file 命名規約と test 関数命名規約で区別する：file は `*_boundary_test.go`、test 関数 prefix は `TestBoundary_`。
  - 別 task `workflow-contract-maintenance` で「独立 component に切るか」を判断する観察項目として残す。
- frontend 側境界結合テスト: 既存 `frontend/src/controller/wails/` 配下の `*.gateway.test.ts` とは別 file として、同じ folder に `master-dictionary.boundary.test.ts` を置く。理由は wails runtime mock 機構を gateway test と共有でき、frontend test の検出経路を CI 上で分けずに済むため。
- 共有 golden の配置: `frontend/` と `internal/` のどちらにも属さない、repo 直下または `internal/apitest/testdata/` 相当の独立 path に置く必要がある。本 task では **`internal/apitest/testdata/boundary/master_dictionary/`** に Go embed 用 JSON file として置き、frontend 側からは vitest 実行時に相対 path で読み込む形に倒す。
  - 同一の物理 file を両側が読むため、片側だけの書き換えで他方が落ちる検出経路（`boundary-test-REQ-003` 後段）が成立する。
  - JSON 形式とした根拠: 両言語から型注釈なしで読める最小公倍数の形式である。型は Go test 側で DTO 構造体に decode し、frontend test 側で TypeScript type に cast する。
- 境界 API 仕様 draft の配置: 本 task では active plan folder 内 `./boundary-api-draft.master-dictionary.md` として作る。正本反映先は detail-spec-diff `boundary-test-REQ-005` に従い、第一候補は `docs/detail-specs/master-dictionary.md` 内の「境界 API」節とする。
- 観察記録 docs draft の配置: 本 task では active plan folder 内 `./observation-record.boundary-test.md` として作る。pilot 実装中に detail-spec-diff `boundary-test-REQ-006` 末尾の観察項目 5 件（test-design 不足 / decision table 追加要否 / 既存 skill 文言衝突 / golden 更新手順独立性 / 閾値）を実装担当者が記録する。

### handoff 分割原則の遵守確認

- 本 task は production code を触らず、テスト層導入と test code 追加に閉じる。
- それでも `境界結合テスト pilot` という独立した受け入れユースケースを backend 側 test と frontend 側 test で分割する。理由は、両者は同じ golden を介して接続するが、`go test` と `vitest` という別検証 command を持ち、変更ファイル群が言語別に分かれるため、`境界規約` の「backend と frontend は同一 handoff に含めない」原則に従う。
- 共有 golden は両 handoff が依存する shared 境界であるため、**wave-1 の独立 handoff として先行確定**させる。これは `境界規約` の「統合境界 handoff は API、Wails、DTO、gateway、adapter 契約 の接続と実画面確認を扱う」を golden 共有点の確定として読み替えた扱いである。
- UI 変更がない task のため、`境界規約` の「UI がある task では frontend を backend より前」の制約は適用対象外。

- human review 済みの判断だけを書く
- frontend handoff がある場合は、承認済み `screen-design-diff.<screen-id>.md` を source にする → 本 task は画面変更なし、`N/A`
- `unanswered_questions`: `0`
- 承認済み詳細仕様差分と回答欄だけを handoff source にする
- downstream handoff が依存する public seam は各実装 handoff の完了条件として固定する → 本 task では共有 golden の path / 値構造を public seam として固定する
- secret を扱う handoff は参照値、secret 本体、secret 解決責務層、出力禁止値を分ける → 本 task は secret を扱わない
- backend、frontend、統合境界 は原則として別 handoff に分ける
- `E2E` は UI 人間操作起点だけを指す → 本 task は E2E 追加を含まない
- `APIテスト` は public seam 起点の system-level test とする → 本 task の境界結合テストは APIテストではなく独立 test 種別として扱う

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `handoff-shared-golden`, `handoff-boundary-api-draft` | `なし` | `handoff-shared-golden <-> handoff-boundary-api-draft` | `なし` |
| `wave-2` | `handoff-boundary-test-backend`, `handoff-boundary-test-frontend` | `wave-1` 完了（共有 golden の path と値構造、境界 API draft が固定済み） | `handoff-boundary-test-backend <-> handoff-boundary-test-frontend` | `なし` |
| `wave-3` | `handoff-observation-record` | `wave-2` 完了（pilot 実装中に発生した観察項目を記録するため） | `なし` | `なし` |

## Handoffs

### `handoff_id`: `handoff-shared-golden`

- `implementation_target`: pilot 対象 6 method（`ListMasterDictionaryEntries` / `GetMasterDictionaryEntry` / `CreateMasterDictionaryEntry` / `UpdateMasterDictionaryEntry` / `DeleteMasterDictionaryEntry` / `ImportMasterDictionaryXml`）× 代表 case の golden JSON file 群と、両側 test からの読み込み helper を追加する。
- `implementation_artifact`: 統合境界実装
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`（`boundary-test-REQ-003`、`master-dictionary-REQ-pilot-001`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - 共有 golden 配置 path `internal/apitest/testdata/boundary/master_dictionary/` を新規作成する。
  - pilot 対象 method × 代表 case ごとに golden JSON file を作る。代表 case は detail-spec-diff `master-dictionary-REQ-pilot-001` 仕様の次を網羅する。
    - 正常系: 1 件以上の代表エントリを含む `ListMasterDictionaryEntries` 応答
    - 空集合: エントリ 0 件の `ListMasterDictionaryEntries` 応答
    - 不在: `GetMasterDictionaryEntry` に対する不在 ID 入力の `null` 応答
    - 状態遷移 3 連鎖: 作成 → 取得、更新 → 取得、削除 → 取得
    - REC 取捨: 13 種別内 REC を含む XML / 13 種別外 REC を含む XML の `ImportMasterDictionaryXml` 応答 field 差分
  - golden JSON は「入力相当の説明」と「期待される応答 DTO field 値」を 1 ファイル 1 case で持つ構造とする。
  - backend 側読み込み helper（Go embed + JSON decode）を `internal/apitest/` 配下に追加する。
  - frontend 側読み込み helper（vitest から相対 path で JSON を fetch して TypeScript type に cast）を `frontend/src/controller/wails/` 配下に追加する。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `handoff-boundary-api-draft`
- `parallel_blockers`: `なし`（境界 API draft は docs draft のみで shared file を触らない）
- `first_action`:
  - path: `internal/apitest/testdata/boundary/master_dictionary/list_normal.golden.json`
  - 対象単位: `ListMasterDictionaryEntries` の正常系代表 case 用 golden JSON file
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「pilot 対象 6 method の代表 case golden JSON が `internal/apitest/testdata/boundary/master_dictionary/` 配下に揃う」
  - 1 手目にする理由: 共有 golden は両 handoff の shared 境界であり、配置 path の合意点を最初に物理 file として固定することで、後続 backend / frontend handoff の `embed path` と `相対 path` が同じ起点を読めるようにする。
- `validation_commands`:
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && find internal/apitest/testdata/boundary/master_dictionary -name '*.golden.json' | wc -l` で golden file 数を確認する。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && jq empty internal/apitest/testdata/boundary/master_dictionary/*.golden.json` で全 golden の JSON 構文を検証する。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && go build ./internal/apitest/...` で backend 側読み込み helper の compile を確認する。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && npx tsc --noEmit -p frontend/tsconfig.json` で frontend 側読み込み helper の type 検証を行う。
- `completion_signal`:
  - pilot 対象 6 method の代表 case golden JSON が `internal/apitest/testdata/boundary/master_dictionary/` 配下に揃う。
  - 各 golden file は「入力相当の説明」と「期待 DTO field 値」を持つ。
  - backend 側 helper が embed JSON を decode して構造体として読める。
  - frontend 側 helper が同じ JSON を読み TypeScript type に cast できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`（境界結合テストは APIテスト/UI人間操作E2E のいずれにも該当しない独立 test 種別）
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff は backend 実装でも frontend 実装でもなく、両側を接続する shared 境界（golden 配置 path、helper 契約）の確定のため、`統合境界実装` に分類する。
  - 並列 pair `handoff-boundary-api-draft` は docs draft 追加のみで、touch する file が `./boundary-api-draft.master-dictionary.md` だけのため、shared 境界変更にあたらない。
  - 本番経路: 共有 golden は production code path に登場しない test 専用 asset。`internal/apitest/testdata/boundary/master_dictionary/<method>_<case>.golden.json` を backend test が `embed` で読み、frontend test が相対 path で読む。
  - golden 値の決定性のため、ID / clock / 並び順は固定値を使う。

### `handoff_id`: `handoff-boundary-api-draft`

- `implementation_target`: pilot 範囲 6 method の境界 API 仕様 draft を active plan folder 内に Markdown で作る。
- `implementation_artifact`: 単体テスト 以外の docs 追加（active plan folder 内 draft、正本反映は finalization-module 担当）→ `implementation_artifact` の選択肢に該当 category がないため、本 handoff は **`統合境界実装` の補助 docs draft 扱い**として実装モジュールに渡す。
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`（`boundary-test-REQ-005`、`master-dictionary-REQ-pilot-001`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `./boundary-api-draft.master-dictionary.md` を新規作成する。
  - 記述対象は pilot 6 method の次の項目とする。
    - method 名と入出力対応
    - 応答 DTO の field 一覧（名、型、値域、null 許容、欠落時の扱い）
    - 列挙値の意味と網羅性
    - 状態遷移の前後で変化する field と変化しない field の規約
    - エラー応答契約
  - 正本反映先案として `docs/detail-specs/master-dictionary.md` 内「境界 API」節 を冒頭に明記する。
  - 機械可読形式昇格の判断は本 task では行わず、observation-record に観察項目として渡す。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `handoff-shared-golden`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `docs/exec-plans/active/master-dictionary-boundary-test/boundary-api-draft.master-dictionary.md`
  - 対象単位: `ListMasterDictionaryEntries` の入出力対応節
  - 変更種別: 新規追加（active plan folder 内 docs draft）
  - 対応する `completion_signal` clause: 「pilot 6 method 全てについて入出力、応答 DTO field 一覧、列挙値、状態遷移規約、エラー応答契約が記述される」
  - 1 手目にする理由: 一覧 method は他 5 method の応答 DTO field 構造の起点になるため、ここを最初に固定すると `Get` / `Create` / `Update` / `Delete` の DTO 表記が一致する。
- `validation_commands`:
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && test -f docs/exec-plans/active/master-dictionary-boundary-test/boundary-api-draft.master-dictionary.md` でファイル存在確認。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && grep -c '^## ' docs/exec-plans/active/master-dictionary-boundary-test/boundary-api-draft.master-dictionary.md` で 6 method 分の節がある（最低 6 個の `## ` 見出し）ことを確認。
- `completion_signal`:
  - `./boundary-api-draft.master-dictionary.md` が active plan folder 内に存在する。
  - pilot 6 method 全てについて入出力、応答 DTO field 一覧、列挙値、状態遷移規約、エラー応答契約が記述される。
  - 正本反映先案として `docs/detail-specs/master-dictionary.md`「境界 API」節を明示する。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff は docs draft 追加のみで production code / production test を touch しない。正本反映は本 task では行わず、後続 `finalization-module` 経由で `updating-docs` skill を使う前提とする。
  - 並列 pair `handoff-shared-golden` と shared file を触らないため並列実行可能。

### `handoff_id`: `handoff-boundary-test-backend`

- `implementation_target`: pilot 6 method について、bootstrap 済み Wails controller を入口にした境界結合テスト backend 側を `internal/apitest/` 配下に追加する。
- `implementation_artifact`: シナリオテスト（既存 decision table に「境界結合テスト」工程がないため、最も近い `シナリオテスト` 区分で実装モジュールに渡し、observation-record に「decision table に独立工程として追加すべきか」の観察項目を残す）
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`（`master-dictionary-REQ-pilot-001`、`boundary-test-REQ-002`、`boundary-test-REQ-003`、`boundary-test-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `internal/apitest/master_dictionary_boundary_test.go` を新規追加する。
  - bootstrap 済み Wails controller の `ListMasterDictionaryEntries` / `GetMasterDictionaryEntry` / `CreateMasterDictionaryEntry` / `UpdateMasterDictionaryEntry` / `DeleteMasterDictionaryEntry` / `ImportMasterDictionaryXml` を入口にする。
  - 代表 case ごとに `handoff-shared-golden` で固定した golden JSON を読み、controller 応答 DTO の field 値が golden と一致することを assert する。
  - 状態遷移 3 連鎖（作成 → 取得、更新 → 取得、削除 → 取得）は同一 test 関数内で順序実行する。
  - REC 取捨は応答 DTO field 観測（detail-spec-diff `Q-007` 確定）で行う。`IsTermTarget` の 13 種別判定そのものの assert は本 handoff の責務外。
  - 決定性のため clock / random / ID / 並び順 / 外部応答を test 内で固定する。
  - test 関数命名は `TestBoundary_MasterDictionary_<Method>_<Case>` prefix に揃える。
  - production Go code（`internal/<package>/` 配下の `_test.go` を除く file）は触らない。
- `depends_on`: `handoff-shared-golden`（golden JSON file と backend 側 helper が揃っていること）
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `handoff-boundary-test-frontend`
- `parallel_blockers`: `なし`（両者は同じ golden を read-only で読むだけで、shared 境界を変更しない。frontend test と backend test の検証 command は別言語・別 process。）
- `first_action`:
  - path: `internal/apitest/master_dictionary_boundary_test.go`
  - symbol: `TestBoundary_MasterDictionary_ListMasterDictionaryEntries_Normal`
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「pilot 6 method の代表 case について、bootstrap 済み controller 応答が golden の field 値と一致することを assert する test 群が `internal/apitest/master_dictionary_boundary_test.go` に揃う」
  - 1 手目にする理由: 一覧 method の正常系は最も単純な field 一致 assert であり、`handoff-shared-golden` の embed helper と controller bootstrap helper を最小経路で結線できるため、後続 case の test 雛形になる。
- `validation_commands`:
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && go test -run TestBoundary_MasterDictionary ./internal/apitest/...` で全 pilot 境界結合テストが pass することを確認する。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && go vet ./internal/apitest/...` で vet 通過を確認する。
- `completion_signal`:
  - pilot 6 method の代表 case について、bootstrap 済み controller 応答が golden の field 値と一致することを assert する test 群が `internal/apitest/master_dictionary_boundary_test.go` に揃う。
  - 状態遷移 3 連鎖の test 関数が存在する。
  - REC 取捨が応答 DTO field 観測の形で固定されている。
  - 片側だけ golden を書き換えると本 test 群が失敗する状態が成立する（手動確認: 任意 golden の値を 1 つ書き換えて `go test` が落ちることを記録）。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`（APIテストでも UI人間操作E2E でもない独立 test 種別）
- `execution_stage`: `実装後`
- `notes`:
  - `implementation_artifact` を `シナリオテスト` で渡す理由: 既存 decision table（`implementation-module`）には「境界結合テスト」工程がないため、最も近い「Wails 公開 API 起点の system-level test」区分の `シナリオテスト` に寄せる。これは detail-spec-diff `Q-006` 観察項目「decision table に追加工程が必要だったか」を後続別 task で判断する材料を作るための暫定運用である。
  - 本番経路: bootstrap 済み Wails controller → service → repository（in-memory または test DB）→ 応答 DTO。controller method 1 件を 1 単位として golden と突き合わせる。
  - 並列 pair `handoff-boundary-test-frontend` は同じ golden を read するが、shared 境界（golden）の書き換えを伴わないため並列可能。検証 command は `go test` のみで失敗時の特定が容易。

### `handoff_id`: `handoff-boundary-test-frontend`

- `implementation_target`: pilot 6 method について、wails runtime を mock した上で golden を mock 応答として供給し、frontend gateway / 解釈層の field 値が golden と矛盾しないことを assert する境界結合テスト frontend 側を `frontend/src/controller/wails/` 配下に追加する。
- `implementation_artifact`: シナリオテスト（同上の理由で `シナリオテスト` 区分に寄せる。observation-record に decision table 追加工程の観察項目として残す）
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`（`master-dictionary-REQ-pilot-001`、`boundary-test-REQ-002`、`boundary-test-REQ-003`、`boundary-test-REQ-004`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`（画面変更なし、`storybook-module` を経由しない）
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `frontend/src/controller/wails/master-dictionary.boundary.test.ts` を新規追加する。
  - 既存 `master-dictionary.gateway.test.ts` の wails runtime mock 機構を共有しつつ、`handoff-shared-golden` で固定した golden JSON を mock 応答として登録する。
  - pilot 6 method 各々について、gateway / 解釈層を経由した結果の field 値が golden の想定と矛盾しないことを assert する。
  - 状態遷移 3 連鎖の case は、同一 test 内で mock 応答を順次差し替える形で実装する。
  - 既存 production frontend code（`frontend/src/controller/wails/master-dictionary.gateway.ts` を含む）は触らない。
  - test 命名は `boundary > master-dictionary > <method> > <case>` の describe / it 階層に揃える。
- `depends_on`: `handoff-shared-golden`（golden JSON file と frontend 側 helper が揃っていること）
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `handoff-boundary-test-backend`
- `parallel_blockers`: `なし`（前述同様、read-only な golden 共有のみで shared 境界を変更しない）
- `first_action`:
  - path: `frontend/src/controller/wails/master-dictionary.boundary.test.ts`
  - symbol: `describe('boundary > master-dictionary > ListMasterDictionaryEntries')`
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「pilot 6 method の代表 case について、wails runtime mock 経由で gateway / 解釈層が返す field 値が golden と矛盾しないことを assert する test 群が `frontend/src/controller/wails/master-dictionary.boundary.test.ts` に揃う」
  - 1 手目にする理由: 一覧 method の正常系は既存 gateway test と同一の mock 接続 path を使うため、`handoff-shared-golden` の frontend 読み込み helper と既存 wails runtime mock を結線する最初の test として最適。
- `validation_commands`:
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP/frontend && npx vitest run src/controller/wails/master-dictionary.boundary.test.ts` で frontend 側境界結合テストが pass することを確認する。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP/frontend && npx tsc --noEmit -p tsconfig.json` で type 検証を行う。
- `completion_signal`:
  - pilot 6 method の代表 case について、wails runtime mock 経由で gateway / 解釈層が返す field 値が golden と矛盾しないことを assert する test 群が `frontend/src/controller/wails/master-dictionary.boundary.test.ts` に揃う。
  - 状態遷移 3 連鎖の test が存在する。
  - 片側だけ golden を書き換えると本 test 群が失敗する状態が成立する（手動確認: 任意 golden の値を 1 つ書き換えて `vitest` が落ちることを記録）。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 本 handoff の `implementation_artifact` を `シナリオテスト` に寄せる理由は `handoff-boundary-test-backend` と同じく decision table の暫定運用。
  - 本 handoff は UI 表示変更を伴わないため `storybook-module` を経由しない。frontend ロジック実装（state / API / Wails bridge / ルーティング / 副作用 / フォーム validation）も触らない。test code 追加のみ。
  - 本番経路: 既存 wails runtime mock → 既存 `master-dictionary.gateway.ts` → 解釈結果 field 値。

### `handoff_id`: `handoff-observation-record`

- `implementation_target`: pilot 実装中に detail-spec-diff `boundary-test-REQ-006` 末尾の観察項目を記録した `./observation-record.boundary-test.md` を作る。
- `implementation_artifact`: シナリオテスト 以外の docs 追加（active plan folder 内 draft）→ 実装モジュールには `tests-scenario` 完了後の補助 docs 作業として渡す。
- `implementation_skill`: `tests-scenario`（暫定。docs draft 追加のため最終的に finalization-module で取り扱う）
- `spec_basis`: `./detail-spec-diff.md`（`boundary-test-REQ-006`）
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `なし`
  - `secret_values_for_provider_external_api_internal_auth`: `なし`
  - `secret_resolution_owner_layer`: `なし`
  - `forbidden_outputs`: `なし`
- `owned_scope`:
  - `./observation-record.boundary-test.md` を新規作成する。
  - 次の 6 観察項目について、pilot 実装中の事実を記録する（detail-spec-diff `boundary-test-REQ-006` 末尾、および本 implementation-scope の人間レビュー観点 4 で追加した観察項目）。
    - `test-design` の入力フォーマットに不足項目が出たか（出た場合は何が不足したか）
    - `implementation-module` の decision table に追加工程が必要だったか（必要な場合は工程の入出力）
    - 既存 skill 本文の文言が境界結合テストの責務と衝突したか（衝突した場合は具体箇所）
    - golden 更新手順が既存 skill 群（`tests-unit`、`tests-scenario`）と独立工程として成立したか
    - 上記観察項目のうち、既存 skill 改訂で吸収しきれない数が閾値を超えるか
    - 片側書き換え検出を自動化する場合の必要工程・必要場所（自動化判断は本 task では行わず、別 task `workflow-contract-maintenance` 起動時の判断材料として残す）
  - 別 task `workflow-contract-maintenance` 起動時の判断材料になる形式（観察事実 → 判断材料 → 暫定推奨）で記述する。
- `depends_on`: `handoff-shared-golden`, `handoff-boundary-api-draft`, `handoff-boundary-test-backend`, `handoff-boundary-test-frontend`（pilot 実装中の事実を記録するため、全 wave-1 / wave-2 完了を前提とする）
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`:
  - path: `docs/exec-plans/active/master-dictionary-boundary-test/observation-record.boundary-test.md`
  - 対象単位: 観察項目 5 件の見出し節を空節として並べる
  - 変更種別: 新規追加
  - 対応する `completion_signal` clause: 「観察項目 5 件全てが記録される」
  - 1 手目にする理由: 観察項目 5 件の節構造を先に固定すれば、wave-2 中の実装担当者が後追いで該当節に事実を書き加える形に倒せる。
- `validation_commands`:
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && test -f docs/exec-plans/active/master-dictionary-boundary-test/observation-record.boundary-test.md` で存在確認。
  - `cd /Users/iorishibata/Repositories/AITranslationEngineJP && grep -c '^## ' docs/exec-plans/active/master-dictionary-boundary-test/observation-record.boundary-test.md` で 6 観察項目分の節がある（最低 6 個の `## ` 見出し）ことを確認。
- `completion_signal`:
  - 観察項目 6 件全てが記録される（事実が「特になし」の場合もその旨を明示）。
  - 別 task `workflow-contract-maintenance` 起動時に読める形式になっている。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `notes`:
  - 本 handoff は docs draft の追加であり、正本反映は本 task では行わない。
  - 本書 (`implementation-scope`) の各 handoff `notes` 内で発見された暫定運用（`implementation_artifact` を `シナリオテスト` に寄せた等）も本観察記録の対象に含める。

## Completion Packet

実装モジュールは完了時に次を返す。

- `completed_handoffs`: 上記 5 handoff のうち実装モジュール側で完了したもの。
- `touched_files`: 本 task では production code / production test を touch しない前提。touched は `internal/apitest/master_dictionary_boundary_test.go`、`internal/apitest/testdata/boundary/master_dictionary/*.golden.json`、`internal/apitest/` 配下の golden 読み込み helper、`frontend/src/controller/wails/master-dictionary.boundary.test.ts`、`frontend/src/controller/wails/` 配下の golden 読み込み helper、active plan folder 内 draft 2 件。
- `implemented_scope`: 「境界結合テスト pilot（master-dictionary、Wails 境界 6 method）」「共有 golden の path / 値構造の確定」「境界 API draft」「観察記録 draft」。
- `test_results`: `go test ./internal/apitest/...` および `vitest run src/controller/wails/master-dictionary.boundary.test.ts` の結果。
- `implementation_investigation`: 片側書き換え検出（golden を 1 つ書き換えて両側 test が落ちるかの手動確認）の事実。
- `ui_evidence`: `N/A`（UI 変更なし）
- `final_validation_result`: 全 wave 完了後の `go test ./...` および frontend 全 test 実行結果。
- `codex_review_result`: 実装モジュール完了後に呼び出し元から確認する。
- `coverage_gate_result`: 既存 coverage gate を境界結合テストで悪化させないことを確認する。
- `sonar_gate_result`: 既存 repo-local Sonar issue gate を維持する。
- `harness_gate_result`: 本 task は Wails 起動を伴わない test 追加であるため、PASS 前提。
- `residual_risks`:
  - 既存 decision table に「境界結合テスト」工程がないため、本 task では `implementation_artifact` を暫定的に `シナリオテスト` に寄せた。観察記録経由で別 task `workflow-contract-maintenance` へ繰り越す。
  - `implementation_skill` を `tests-scenario` に寄せた handoff があるため、実装モジュール内での skill 解釈が暫定運用である旨を観察記録に残す。
- `completion_evidence`: 全 handoff の `completion_signal` 達成記録、片側書き換え検出の手動確認結果、active plan folder 内 draft 2 件の存在。
- `telemetry_events`: `runtime: claude` の response event。
- `docs_changes: none`（正本反映は本 task では行わない。active plan folder 内 draft 2 件は task 内成果物として扱う）。

## 未決事項

なし。detail-spec-diff の `Q-001`..`Q-008` は全て回答固定済み。本 implementation-scope 固有の未決もない。

## 人間レビュー観点（2026-06-05 approve 済み）

- `implementation_artifact` を `シナリオテスト` に寄せた判断（decision table に「境界結合テスト」が存在しないための暫定運用）→ **approve**。後続別 task `workflow-contract-maintenance` 起動時に再判断する。
- 共有 golden の物理配置を `internal/apitest/testdata/boundary/master_dictionary/` に置く判断（独立 path に切らない判断）→ **approve**。
- 境界 API 仕様 draft と観察記録の正本反映を本 task の `finalization-module` 段階で扱う前提 → **approve**。
- 片側書き換え検出を手動確認とした判断（自動 test 化を本 task に含めない判断）→ **保留付き approve**。本 task 範囲では手動確認に留め、`handoff-observation-record` の観察項目として「片側書き換え検出を自動化する場合の必要工程・必要場所」を pilot 実装中に記録し、`workflow-contract-maintenance` 起動時の判断材料にする。
