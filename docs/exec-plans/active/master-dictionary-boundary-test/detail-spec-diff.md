# 詳細仕様差分: master-dictionary-boundary-test

- `skill`: detail-spec-design
- `status`: approved
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/boundary-integration-test.md`（新規正本ファイル案。`docs/detail-specs/` 配下に他の機能別仕様と並べ、test 層自体の仕様を 1 件の正本として扱う）
- `screen_design_diff`: `N/A`
- `component_diagram`: `./diagrams/component.boundary-test-layer.mmd`（後続の `diagramming` で作成する。本書では仕様だけを固定し、図は別成果物とする）

本書は 2 系統の親要件を扱う。

- `boundary-test-REQ-00X`: 新規 test 層「境界結合テスト」自体の追加要件。追加要件として扱い、新規正本ファイル（案: `docs/detail-specs/boundary-integration-test.md`）への正本反映を想定する。
- `master-dictionary-REQ-pilot-001`: 既存 `master-dictionary` 機能を pilot 対象として境界結合テストで証明する追加要件。pilot 範囲の境界仕様の明文化は `docs/detail-specs/master-dictionary.md` への追補として扱う。

「境界結合テスト」は frontend / backend、AI service、SQL repository などの境界一般を対象にする概念として扱う。本 task では frontend-backend 境界（Wails controller method 単位）を pilot 適用範囲とするが、層の責務定義としては AI service の RPC method、SQL repository の query method なども同じ枠組みで扱える形にする。

詳細仕様差分は「仕様として成立する条件、判断できる状態、処理結果」だけを固定する。配置 path、命名規約、CI 連動方式、skill 改訂手順などの実装方式・作業運用は仕様に含めず、`implementation-scope` 以降で扱う。ただし、後続が再解釈すると仕様差分になる判断（test 層の責務境界、golden の単一性、pilot 範囲の選別根拠）は本書で固定する。

## 新規正本ファイル案

正本反映先として新規ファイル `docs/detail-specs/boundary-integration-test.md` を切る案を提示する。既存 `docs/detail-specs/` には機能別仕様（`master-dictionary.md`、`translation-job-management.md` 等）が並ぶ体系であり、test 層自体の仕様は他の機能別仕様と並ぶ独立した正本として置く。`docs/coding-guidelines-tests.md` は改訂せず、必要なら新規正本ファイルへの参照を 1 行追加するだけにとどめる。

正本に載せる節構成案は次のとおり。

1. 目的と層の定義（境界一般を対象にする旨、frontend-backend / AI service / SQL の各境界を含む旨）
2. 既存 test 層との責務境界（単体、API、gateway、UI 人間操作 E2E との差）
3. 導出単位（境界 method 単位を基本とする）
4. assert する責務と assert しない責務
5. 境界契約として固める情報（後述「境界 API 仕様書」追加要件の節と接続）
6. golden 集合の単一性と決定性
7. 契約変更時の改訂規約（意図的変更、壊れた時の対応 flow）
8. skill / workflow への組み込み方針（別 task に委ねる旨と判断観察項目）

本 task では正本ファイル本体の作成はしない。test 戦略 docs 化は後続の候補 B として扱う。

## 詳細仕様差分

### `boundary-test-REQ-001` 境界結合テスト層を独立した test 層として位置づける

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規）

親要件:
開発者は境界（frontend と backend、frontend と AI service、backend と SQL repository など）を接続する DTO 値 semantic と状態遷移後の値を、独立した test 層として証明できる。

仕様:
- 境界結合テストは frontend / backend、AI service、SQL repository などの境界一般を対象にする。本 task の pilot 適用範囲は frontend-backend 境界（Wails controller を入口とする backend test と、同じ Wails binding を消費する frontend test の対）に限るが、層の責務定義としては他の境界にも同じ枠組みで適用できる形にする。
- 1 単位の境界結合テストは、境界の片側が応答として生成する値集合と、もう片側が解釈に使う値集合が同一であることを成立条件にする。
- 既存 test 層との責務境界は次のとおり固定する。
  - backend 単体テスト（`internal/<package>/*_test.go`）は、Go package 内の公開振る舞い、分岐、エラー経路を証明する。境界 DTO 値の semantic は責務にしない。
  - backend API テスト（`internal/apitest/`）は、Wails controller または bootstrap 済み controller を入口にした受け入れ条件を証明する。frontend 側の解釈は責務にしない。
  - frontend 単体テスト（`frontend/src/**/*.test.ts`）は、frontend の公開振る舞い、状態、callback を証明する。backend が実際に返す値の生成は責務にしない。
  - frontend gateway テスト（`frontend/src/controller/wails/*.gateway.test.ts`）は、wails runtime を mock した上で frontend 側の型・解釈・マッピングを証明する。backend が同じ値を実際に生成することは責務にしない。
  - UI 人間操作 E2E テスト（`tests/system/`）は、人間操作の主要 flow と表示結果を証明する。UI に露出しない境界 DTO 値の差分は責務にしない。
- 境界結合テストは、上記いずれの層でも観測されない「境界 DTO の値 semantic が変わった時に検出する」責務だけを単独で担う。

回答済み:
- `Q-001`: 新規 test 層の表記名は「境界結合テスト」に統一する。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書および新規正本ファイル案 `docs/detail-specs/boundary-integration-test.md`。`frontend-backend` 限定にせず、AI service、SQL repository などを含む境界一般の概念として扱う。層の定義文に「frontend / backend、AI service、SQL repository などの境界一般を対象とする」旨を含める。
- `Q-002`: 詳細仕様正本は新規ファイルを切る。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 新規正本ファイル案 `docs/detail-specs/boundary-integration-test.md`。`docs/coding-guidelines-tests.md` の改訂はしない。本 task では正本ファイル本体の作成はせず、新規ファイル名案と節構成案を本書に提示するに留める（test 戦略 docs 化は後続の候補 B で扱う）。

### `boundary-test-REQ-002` 境界結合テストは境界 method 単位で導出する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規）

親要件:
開発者は、境界結合テストの対象を一意に決められる導出単位を持つ。

仕様:
- 境界結合テストの導出単位は、境界 method 1 件とする。frontend-backend 境界では Wails controller の公開 method 1 件、AI service 境界では RPC method 1 件、SQL repository 境界では query method 1 件が単位になる。
- 1 件の method につき、入口は「同じ method 名」、観測点は「method が返す DTO の値集合」とする。
- 同一 method に対して、次の代表 case を境界結合テストの対象 case として扱う。
  - 正常系の代表値 1 件以上。
  - 空集合・null・未選択など、UI に露出しないが消費側の分岐に影響する境界値 case。
  - 状態遷移後の値（例: 作成・更新・削除直後の値、import 直後の値）。
  - 消費側が解釈するエラー応答 case（method がエラーを応答契約として返す場合）。
- 1 件の method の DTO 全 field を 1 case で網羅する責務は持たない。代表 case の組み合わせで「DTO 値 semantic の固定」が成立する粒度に分ける。
- 同一の振る舞いがすでに backend API テスト、frontend gateway テスト、UI 人間操作 E2E テストのいずれかで証明済みの場合は、境界結合テストの対象から外してよい。除外根拠は対象 case 表に明記する。
- 機能単位（feature）や operation 単位は補助概念にとどめ、test ファイル分割の単位ではない。複数 method が 1 つの状態遷移を協調して構成する場合は、case 表の中で機能単位の見出しを補助的に置く扱いを許容するが、導出単位は method 単位のまま維持する。

回答済み:
- `Q-003`: 導出単位は method 単位で固定する。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書および新規正本ファイル案。Wails controller method、AI service RPC method、repository query method など、境界 method を基本単位とする。feature 単位や operation 単位は補助概念にとどめ、test ファイル分割の単位ではない旨を明示する。

### `boundary-test-REQ-003` 境界結合テストは境界の片側が出力した golden を反対側が消費する形で契約を固定する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規）

親要件:
開発者は、境界の片側が実際に応答する DTO 値と反対側が解釈に使う値が同一であることを、単一の参照点を介して固定できる。

仕様:
- 境界結合テストの 1 単位は、両側 test が同一の golden 集合を共有することを成立条件にする。
- golden 集合は、対象 method の代表 case ごとに「入力相当の説明」と「期待される応答 DTO 値」を含む単一の参照点とする。
- 生成側 test（frontend-backend 境界では backend 側）は、bootstrap 済み境界入口に対して対象 case を実行し、実際に応答した DTO 値が golden と一致することを assert する。
- 消費側 test（frontend-backend 境界では frontend 側）は、境界 runtime を mock した上で同じ golden 値を応答として供給し、消費側の解釈・型・マッピングが golden と矛盾しないことを assert する。
- golden 集合は test と独立した参照点として扱い、両側のどちらか一方だけが期待値を書き換えても、もう一方の test が失敗する状態を維持する。
- golden の値生成は決定性を持つ。clock、random、ID、並び順、外部応答を固定する。
- 境界に必要な情報（DTO 構造、field 名、値型、列挙値、状態遷移）は仕様として明示する（後述「境界 API 仕様書追加要件」の節を参照）。
- 契約を意図的に変更する時の手順と、golden が壊れた時の対応 flow は新規正本ファイル案の改訂規約節に含める。具体的には次を含む。
  - 意図的変更時: 仕様（境界 API 仕様書）を先に改訂し、golden を更新し、両側 test を同時に更新する順序を規約として明文化する。
  - 壊れた時: 生成側と消費側のどちらが契約逸脱しているかを切り分け、仕様 → golden → test の順で再整合する flow を規約として明文化する。

回答済み:
- `Q-004`: golden 集合は仕様として固める。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書および新規正本ファイル案。境界に必要な情報（DTO 構造、field 名、値型、列挙値、状態遷移）を仕様として明示する。表現形式（JSON / TypeScript const / Go embed）は `implementation-scope` の選択に委ねる。追加要件として「境界 API 仕様書」が必要になる可能性が高いため、次節「境界 API 仕様書追加要件」で必要性と方向性を明示する。
- `Q-005`: golden 更新手順は仕様に含める。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 新規正本ファイル案の改訂規約節。「契約を意図的に変更する時の手順」と「golden が壊れた時の対応 flow」を新規正本 docs の節として含める。`docs/coding-guidelines-tests.md` 側の運用規約には分けない。

### `boundary-test-REQ-004` 境界結合テストが assert する責務と assert しない責務を分離する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規）

親要件:
開発者は、境界結合テストが扱う観測点と扱わない観測点を再解釈なしに区別できる。

仕様:
- 境界結合テストが assert する責務は次に限る。
  - 対象 method が返す DTO の値 semantic（field 名、型、値域、enum 値、null 許容、欠落 field の扱い）。
  - 状態遷移後（作成、更新、削除、import）に method が返す値の整合。
  - 空集合・未選択・null など、UI に露出しないが消費側分岐に影響する境界値の応答。
  - エラー応答契約として method が返す値の semantic（method が応答として error を返す場合）。
  - 消費側 gateway / 解釈層が、上記 DTO 値を矛盾なく受理できること。
- 境界結合テストが assert しない責務は次のとおり固定する。
  - UI 表示（layout、文言、style）。
  - 業務 flow 全体（複数画面、複数 method を順に組み合わせた業務シナリオ）。
  - SQL 実行計画、index 利用、storage 内部挙動。
  - 観測ログのフォーマット、ログ宛先、ログ駆動の原因分離。

回答済み:
- なし

### `boundary-test-REQ-005` 境界 API 仕様書追加要件（Q-004 起因）

- `変更種別`: 追加
- `要件扱い`: 追加要件（本 task では必要性と方向性のみ仕様化、具体作成は後続 task）
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規。仕様書の位置づけ節として記載）

親要件:
開発者は、境界結合テストの golden が参照する DTO 構造・field 名・値型・列挙値・状態遷移を、test と独立した「境界 API 仕様書」として参照できる。

仕様（必要性）:
- 境界結合テストは golden の単一参照点を成立条件にするが、golden 自体は値の列であり、構造の意味（field の役割、列挙値の意味、状態遷移の規約）を持たない。
- DTO 構造・field 名・値型・列挙値・状態遷移を仕様として明示するためには、test と独立した「境界 API 仕様書」が必要になる可能性が高い。
- 境界 API 仕様書は、境界 method 単位で次を記述する想定とする。
  - method 名と入出力の対応。
  - 応答 DTO の field 一覧（名、型、値域、null 許容、欠落時の扱い）。
  - 列挙値の意味と網羅性。
  - 状態遷移の前後で変化する field と変化しない field の規約。
  - エラー応答契約。

仕様（形式案）:
- 機械可読形式（OpenAPI、JSON Schema、TypeScript type）と人間可読仕様（Markdown）の併用、または Markdown のみの 2 案を後続 task で比較する。
- pilot 範囲では、まず Markdown による仕様記述で necessary 情報を満たせるかを試し、満たせない論点が出た時に機械可読形式の導入を再検討する。

仕様（置き場所案）:
- 境界 API 仕様書は機能別詳細仕様（`docs/detail-specs/<feature>.md`）の中に節として置く案と、独立した境界仕様 docs（例: `docs/detail-specs/<feature>-boundary-api.md`）として置く案を後続 task で比較する。
- pilot 範囲（`master-dictionary`）では、`docs/detail-specs/master-dictionary.md` 内の「境界 API」節として置く案を第一候補とする。

仕様（本 task での扱い範囲）:
- 本 task では「境界 API 仕様書」の必要性と方向性を仕様差分として固定するに留める。
- 具体的な境界 API 仕様書の作成、形式の最終決定、置き場所の最終決定は本 task の範囲外とし、後続 task（test 戦略 docs 化、または pilot 実装 task）に委ねる。

回答済み:
- なし

### `boundary-test-REQ-006` 境界結合テスト層の skill / workflow への組み込みを方針として固定する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/boundary-integration-test.md`（新規）

親要件:
開発者と agent は、境界結合テストの判断・設計・実装を既存 skill 群と整合した手順で進められる。

仕様:
- 境界結合テストは、既存 test 層（単体、API、UI 人間操作 E2E）のいずれとも独立な test 種別として、test 設計と実装スコープの decision table に登場する。
- `design-module` の test-design 起動入力は、境界結合テスト観点を test 種別の 1 つとして受け取れる状態にする。
- `implementation-module` の decision table（シナリオテスト / 単体テスト / 観測ログ / 最終検証）は、境界結合テストを独立工程として扱える状態にする。
- 既存 skill（`test-design`、`tests-scenario`、`tests-unit`、`implementation-module`、`design-module`）の本文改訂、または新規 skill（仮称 `tests-boundary`）の追加のいずれかで上記を満たす。どちらを採るかは本 task では決めず、pilot 完了後に別 task `workflow-contract-maintenance` で決める。
- 本 task では skill 本体の改訂は行わない。本 task の成果物は「境界結合テスト層の存在と責務」「skill / workflow への組み込みが必要であること」「pilot 完了後に別 task で skill 改訂方針を決めること」を仕様として固定するところまでとする。
- pilot 完了後の別 task での判断に必要な観察項目を次のとおり固定する。
  - pilot 実装で `test-design` の入力フォーマットに不足項目が出たか（出た場合は何が不足したか）。
  - pilot 実装で `implementation-module` の decision table に追加工程が必要だったか（必要な場合は工程の入出力）。
  - pilot 実装で既存 skill 本文の文言が境界結合テストの責務と衝突したか（衝突した場合は具体箇所）。
  - pilot 実装で golden 更新手順が既存 skill 群（`tests-unit`、`tests-scenario`）と独立工程として成立したか。
  - 上記観察項目のうち、既存 skill 改訂で吸収しきれない数が閾値を超える場合は新規 skill `tests-boundary` を追加する判断材料にする。

回答済み:
- `Q-006`: skill 改訂は別 task で扱う。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書および新規正本ファイル案。本 task では「skill / workflow への組み込みが必要であること」「skill 本体改訂は別 task で扱うこと」「pilot 完了後に別 task で判断すること」「判断に必要な観察項目」を仕様として固定する。新規 `tests-boundary` skill の要否、既存 skill の改訂箇所、`workflow-contract-maintenance` 起動の trigger 条件は pilot 完了後に decide する。

### `master-dictionary-REQ-pilot-001` MasterDictionary 機能を境界結合テスト pilot として証明する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/detail-specs/master-dictionary.md`（pilot 対象として境界結合テストで証明する旨の追補。詳細 case 表の配置は後続成果物に委ねる）

親要件:
開発者は、直近の大規模 refactor 後に唯一手動動作確認済みの MasterDictionary 機能を pilot とし、境界結合テスト枠組みが現実に機能することを実証できる。

仕様:
- pilot 対象は、MasterDictionary 機能の Wails controller 公開 method の集合とする。
  - `ListMasterDictionaryEntries`: 一覧応答の field 集合（canonical field のみで `rec` / `edid` を含まない）と pagination 値の semantic を境界結合テストの対象にする。
  - `GetMasterDictionaryEntry`: 単一エントリ応答の field 集合、エントリ不在時の `null` 応答、`note` の固定文言の semantic を境界結合テストの対象にする。
  - `CreateMasterDictionaryEntry`: 作成直後に応答が返すエントリ値が、`ListMasterDictionaryEntries` / `GetMasterDictionaryEntry` の同一エントリ参照と同じ semantic で読めることを境界結合テストの対象にする。
  - `UpdateMasterDictionaryEntry`: 更新後の応答が、更新前後で変化する field と変化しない field を判別できる semantic を持つことを境界結合テストの対象にする。
  - `DeleteMasterDictionaryEntry`: 削除応答が削除完了状態と削除後の参照可否（`null` 応答へ遷移する）の semantic を持つことを境界結合テストの対象にする。
  - `ImportMasterDictionaryXml`: 取り込み結果応答が、取り込まれた件数、対象外として扱われた REC の集約情報、エラー区分の semantic を持つことを境界結合テストの対象にする。
- pilot で扱う代表 case は次を含む。
  - 正常系: 1 件以上の代表エントリを含む応答。
  - 空集合: エントリ 0 件の `ListMasterDictionaryEntries` 応答。
  - 不在: `GetMasterDictionaryEntry` に対する不在 ID 入力の `null` 応答。
  - 状態遷移: 作成 → 取得、更新 → 取得、削除 → 取得の 3 連鎖。
  - REC 取捨: 13 種別内 REC を含む XML 取り込み応答と、13 種別外 REC を含む XML 取り込み応答の応答側 field の差分。
- pilot 範囲の REC 取捨 case は、応答 DTO の field（取り込み件数、対象外件数、対象外区分、エラー区分など）として観測する形に統一する。`IsTermTarget` の 13 種別判定そのものは backend 単体テストの責務として残し、境界結合テストの観測点にはしない。
- 判定結果の理由を境界契約として残したい場合は、応答 DTO に `kind`（対象外区分）、`skippedReason`（対象外理由）などの field を追加する設計に倒す。backend が「条件」を expose する規約（コーディング規約「状態 / 種別 / 条件」原則）に従う。
- pilot 範囲から外す method・case は次を根拠に明記する。
  - UI に露出しない値変化がない（既存 frontend gateway test と backend API test で意味的に重複する）case は対象外。
  - SQL 内部挙動・index 効果は境界結合テストの責務外（`boundary-test-REQ-004` に従う）。
- pilot 完了の成立条件は、上記 method の代表 case について backend 側応答と frontend 側解釈が同一 golden を介して整合することを test として固定できることとする。

回答済み:
- `Q-007`: REC 取捨は応答 DTO field 観測で固定する。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書および pilot 範囲。理由は次のとおり。(1) 境界結合テストの責務（境界契約検証）に一致する。(2) frontend が消費する形と一致する。(3) コーディング規約「状態 / 種別 / 条件」原則と整合する。(4) `Q-001` で境界一般に拡張する方針と整合し、AI service / SQL でも同じ形が使える。判定結果の理由を境界契約として残したい場合は DTO に `kind`、`skippedReason` などの field を追加する設計に倒す（backend が「条件」を expose する規約に従う）。`IsTermTarget` の 13 種別判定そのものは backend 単体テストの責務として残す。
- `Q-008`: pilot case 表の配置先は後続成果物に委ねる。
  - `回答者`: 人間
  - `回答日`: 2026-06-05
  - `反映先`: 本書では pilot case 観点の方向性だけ残し、case 表本体の配置場所（`docs/detail-specs/boundary-integration-test.md` か `docs/detail-specs/master-dictionary.md` か）は `implementation-scope` または `test-design` 起動時に decide する。

## 根拠

- `source`:
  - `docs/exec-plans/active/master-dictionary-boundary-test/plan.md`
  - `docs/coding-guidelines-tests.md`
  - `docs/detail-specs/master-dictionary.md`
  - `internal/apitest/README.md`
  - `frontend/src/controller/wails/master-dictionary.gateway.test.ts`
  - `frontend/src/controller/wails/master-dictionary.gateway.ts`
  - `tests/system/master-dictionary-management.spec.ts`
  - `internal/apitest/provider_settings_contract_freeze_test.go`（既存の contract 凍結 test との責務差を確認するため参照）
  - `.claude/skills/test-design/`、`.claude/skills/tests-scenario/`、`.claude/skills/tests-unit/`、`.claude/skills/implementation-module/`、`.claude/skills/design-module/`（skill / workflow 組み込み方針の根拠）
- `review`: 人間レビュー完了（Q-001..Q-008 すべて回答固定、2026-06-05）。
- `validation`: 未実行。本書は詳細仕様差分であり、検証は `implementation-scope` 確定後に test 実装で行う。
