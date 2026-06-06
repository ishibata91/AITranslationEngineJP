# テスト設計: master-dictionary-boundary-test

- `skill`: test-design
- `source_plan`: `./plan.md`
- `source_detail_spec_diff`: `./detail-spec-diff.md` (status: approved)
- `source_design_diff_diagram`: `./design-diff-diagram.md` (status: human-review-resolved)
- `status`: draft

## テスト設計の位置づけ

本書は、境界結合テストという新規 test 種別の pilot 設計書である。
通常の `test-design` skill が扱う UI 人間操作 E2E テスト（観点表 CSV 形式）ではなく、
MasterDictionary の Wails controller method を対象とした境界結合テストの観点を固定する。

境界結合テストは `boundary-test-REQ-001` で定義する独立した test 層であり、
「backend が応答として生成する DTO 値 semantic」と「frontend が消費する DTO 値の解釈」を、
同一 golden を介して契約として固定する責務を担う。

---

## 1. 既存テスト層との責務境界

### 1.1 境界結合テストが責任を持つ観点

| 観点 | 理由 |
| --- | --- |
| backend が実際に返す DTO の field 値 semantic（field 名、型、値域） | backend 単体テストは Go package 内振る舞いを証明するが、DTO 値の semantic は責務にしない |
| `id` field が string 形式（int64 ではなく）で返ることの backend 側証明 | frontend gateway test は frontend 解釈だけを証明し、backend が実際に string を返すことは証明しない |
| `note` field の固定文言「マスター辞書エントリ」が backend から返ることの backend 側証明 | frontend gateway test は frontend 解釈だけを証明し、backend 側の生成は責務にしない |
| `GetMasterDictionaryEntry` が不在 ID に対して `entry: null` を返すことの backend 側証明 | frontend gateway test は null 応答時の frontend 解釈を証明するが、backend が実際に null を返すことは証明しない |
| 状態遷移後（作成・更新・削除）の応答 field 値が参照系 method の field 値と semantic 整合すること | system E2E は UI 表示を証明するが、DTO 値の同一 semantic は証明しない |
| `ImportMasterDictionaryXml` 応答の `importedCount`/`updatedCount`/`skippedCount` semantic | system E2E は UI 上の「完了」表示を証明するが、各 count field の値 semantic は証明しない |
| 13 種別外 REC を含む XML の取り込み時に `skippedCount` が増加すること | backend 単体テストは `IsTermTarget` 判定を証明するが、応答 DTO の `skippedCount` semantic は責務にしない |
| frontend 側が同じ golden の field 値を受理できること（型・解釈・マッピングの両側整合） | 上記の backend 側証明と frontend 側証明を同一 golden で結びつける責務は他層に存在しない |

### 1.2 他層に委ねる観点（境界結合テストの対象外）

| 観点 | 委ねる層 | 根拠 |
| --- | --- | --- |
| `rec`/`edid` field の非包含（frontend 解釈の正確性） | frontend gateway test（`master-dictionary.gateway.test.ts`） | `boundary-test-REQ-004`：frontend 単独の型・解釈は gateway test の責務 |
| フィルタクエリを Wails binding に正しく渡すこと | frontend gateway test | binding 呼び出し引数の検証は gateway test の責務 |
| UI 上の一覧表示、詳細表示、検索結果表示 | system E2E（`master-dictionary-management.spec.ts`） | `boundary-test-REQ-004`：UI 表示は境界結合テストの対象外 |
| 複数画面にわたる業務 flow（検索 → 選択 → 詳細表示） | system E2E | `boundary-test-REQ-004`：業務 flow 全体は対象外 |
| `IsTermTarget` の 13 種別判定ロジック自体 | backend 単体テスト（`internal/service/` 配下） | `Q-007` 確定：判定ロジックは backend 単体テストの責務として残す |
| SQL 実行計画、index 利用 | backend 単体・repository テスト | `boundary-test-REQ-004`：SQL 内部挙動は対象外 |
| XML 形式不正時の UI 上のエラー表示 | system E2E | UI 表示は境界結合テストの対象外 |

---

## 2. 境界結合テスト観点表

### 2.1 凡例

| 列 | 意味 |
| --- | --- |
| ID | 境界結合テスト観点を一意に識別する ID |
| 関連仕様 | 根拠とする仕様要件 ID |
| 対象 method | Wails controller の公開 method 名 |
| case 分類 | 正常 / 境界 / 状態遷移 / REC 取捨 |
| 入力条件 | method に与える入力相当の説明（golden の入力側） |
| 観測 field | assert 対象の応答 DTO field 一覧 |
| backend 期待値 | backend 側 assert の期待値（golden の値） |
| frontend 期待値 | frontend 側 assert の期待値（同じ golden を消費した結果） |
| 除外根拠 / 備考 | 除外理由または補足 |

### 2.2 観点表

#### `ListMasterDictionaryEntries`

| ID | 関連仕様 | 対象 method | case 分類 | 入力条件 | 観測 field | backend 期待値 | frontend 期待値 | 備考 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BCT-MDC-001 | master-dictionary-REQ-pilot-001 | ListMasterDictionaryEntries | 正常 | エントリ 1 件以上が登録された状態で、フィルタなし・page=1・pageSize=30 を指定する | `entries[0].id`（型）、`entries[0].source`、`entries[0].translation`、`entries[0].category`、`entries[0].origin`、`entries[0].updatedAt`、`totalCount`、`page`、`pageSize` | `entries[0].id` が string 形式（例: `"1"`）; `totalCount` が登録件数と一致する整数; `page=1`; `pageSize=30`; canonical field 全件が非空文字列 | 同じ golden の field 値を wails mock から受け取り、gateway が返す `entries[0]` の各 field が golden と一致する | `rec`/`edid` field の非包含は frontend gateway test に委ねる |
| BCT-MDC-002 | master-dictionary-REQ-pilot-001 | ListMasterDictionaryEntries | 境界 | エントリ 0 件の状態で、フィルタなし・page=1・pageSize=30 を指定する | `entries`（配列長）、`totalCount`、`page`、`pageSize` | `entries` が空配列（length=0）; `totalCount=0`; `page=1`; `pageSize=30` | 同じ golden を消費し、gateway が返す `entries` が空配列で `totalCount=0` と一致する | system E2E は UI の「該当なし」表示を証明するが、DTO 値の `entries=[]` / `totalCount=0` は本層が証明する |

#### `GetMasterDictionaryEntry`

| ID | 関連仕様 | 対象 method | case 分類 | 入力条件 | 観測 field | backend 期待値 | frontend 期待値 | 備考 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BCT-MDC-003 | master-dictionary-REQ-pilot-001 | GetMasterDictionaryEntry | 正常 | 登録済みエントリの ID を指定する | `entry.id`（型）、`entry.source`、`entry.translation`、`entry.category`、`entry.origin`、`entry.updatedAt`、`entry.note` | `entry` が non-null; `entry.id` が string 形式; `entry.note` が `"マスター辞書エントリ"` に固定; 各 canonical field が登録値と一致 | 同じ golden の `entry` object を消費し、gateway が返す各 field が golden と一致する | `note` 固定文言の backend 側生成を証明する（frontend gateway test は frontend 解釈を証明するが backend 側生成は対象外） |
| BCT-MDC-004 | master-dictionary-REQ-pilot-001 | GetMasterDictionaryEntry | 境界 | 存在しない ID（例: `"999999"`）を指定する | `entry`（値） | `entry` が `null` | 同じ golden（`entry: null`）を消費し、gateway が返す `entry` が `null` である | frontend gateway test と観点が近いが、本 case は backend が実際に null を返すことを bootstrap 済み controller で証明する |

#### `CreateMasterDictionaryEntry`

| ID | 関連仕様 | 対象 method | case 分類 | 入力条件 | 観測 field | backend 期待値 | frontend 期待値 | 備考 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BCT-MDC-005 | master-dictionary-REQ-pilot-001 | CreateMasterDictionaryEntry | 正常 | `source`/`translation`/`category`/`origin` をすべて指定して作成を呼び出す | `entry.id`（型）、`entry.source`、`entry.translation`、`entry.category`、`entry.origin`、`entry.note`、`refreshTargetId` | `entry.id` が string 形式; `entry.source` が入力値と一致; `entry.note` が `"マスター辞書エントリ"` に固定; `refreshTargetId` が `entry.id` と同じ string 値 | 同じ golden の `entry` object を消費し、gateway が返す各 field が golden と一致する | 作成直後の応答 DTO field semantic を証明する |

#### 状態遷移 3 連鎖

| ID | 関連仕様 | 対象 method | case 分類 | 入力条件 | 観測 field | backend 期待値 | frontend 期待値 | 備考 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BCT-MDC-006 | master-dictionary-REQ-pilot-001 | CreateMasterDictionaryEntry → GetMasterDictionaryEntry | 状態遷移 | CreateMasterDictionaryEntry で作成したエントリの ID を使い、GetMasterDictionaryEntry を呼び出す | Create 応答の `entry.id`、`entry.source`; Get 応答の `entry.id`、`entry.source`、`entry.note` | Create 応答と Get 応答の `entry.id` が同じ string 値; `entry.source` が同じ値; Get 応答の `entry.note` が `"マスター辞書エントリ"` | frontend 側は Create golden と Get golden をそれぞれ消費し、各 field が golden と一致する | 設計差分図シーケンス図の代表 case。Create → Get の field 値 semantic 整合を証明する |
| BCT-MDC-007 | master-dictionary-REQ-pilot-001 | UpdateMasterDictionaryEntry → GetMasterDictionaryEntry | 状態遷移 | 登録済みエントリを UpdateMasterDictionaryEntry で `translation` だけ変更し、GetMasterDictionaryEntry で参照する | Update 応答の `entry.translation`; Get 応答の `entry.translation`（更新前後の差分） | Update 応答の `entry.translation` が新しい値; Get 応答の `entry.translation` が同じ新しい値; Update 前後で `entry.source` が変化しない | frontend 側は Update golden と Get golden をそれぞれ消費し、各 field が golden と一致する | 更新で変化する field と変化しない field の semantic を証明する |
| BCT-MDC-008 | master-dictionary-REQ-pilot-001 | DeleteMasterDictionaryEntry → GetMasterDictionaryEntry | 状態遷移 | 登録済みエントリを DeleteMasterDictionaryEntry で削除し、同じ ID で GetMasterDictionaryEntry を呼び出す | Delete 応答の `deletedId`; Get 応答の `entry`（null 遷移） | Delete 応答の `deletedId` が削除した ID と同じ string 値; 削除後の Get 応答で `entry` が `null` | frontend 側は Delete golden と Get golden をそれぞれ消費し、`deletedId` と `null` 遷移が golden と一致する | 削除後の参照可否（null 遷移）という状態遷移後 semantic を証明する |

#### `ImportMasterDictionaryXml`（REC 取捨）

| ID | 関連仕様 | 対象 method | case 分類 | 入力条件 | 観測 field | backend 期待値 | frontend 期待値 | 備考 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BCT-MDC-009 | master-dictionary-REQ-pilot-001、boundary-test-REQ-004 | ImportMasterDictionaryXml | REC 取捨（正常） | 13 種別内 REC のみを含む XML を取り込む | `accepted`、`summary.importedCount`、`summary.updatedCount`、`summary.skippedCount`、`summary.fileName` | `accepted=true`; `summary.importedCount` が取り込まれた件数（>0）; `summary.skippedCount=0`; `summary.fileName` が入力 XML のファイル名 | 同じ golden を消費し、gateway が返す `accepted` / `summary` の各 field が golden と一致する | 13 種別内のみなので `skippedCount=0` を境界 DTO で証明する |
| BCT-MDC-010 | master-dictionary-REQ-pilot-001、boundary-test-REQ-004、Q-007 確定 | ImportMasterDictionaryXml | REC 取捨（境界） | 13 種別外 REC を含む XML を取り込む（13 種別内と外が混在する XML が理想的） | `accepted`、`summary.importedCount`、`summary.skippedCount` | `accepted=true`; `summary.importedCount` が 13 種別内 REC の件数; `summary.skippedCount` が 13 種別外 REC の件数（>0） | 同じ golden を消費し、gateway が返す `importedCount`/`skippedCount` の各 field が golden と一致する | `IsTermTarget` の 13 種別判定ロジック自体は backend 単体テストに委ねる。境界結合テストが証明するのは「応答 DTO の `skippedCount` field に件数が乗ること」の semantic である（Q-007 確定） |

---

## 3. golden 集合の設計指針

### 3.1 golden に含める情報

各観点の golden は method × case ごとに次を含む。

- 入力相当の説明（fixture data の種類と値）
- 期待される応答 DTO の field 値（型含む）

### 3.2 決定性の確保

`boundary-test-REQ-003` に従い、次の要素を golden 生成時に固定する。

- clock（`updatedAt` などの timestamp は固定値にする）
- ID（自動採番で変化するため、連番ではなく相対的な一致を assert するか、fixture seed で固定する）
- 並び順（一覧の並び順を決定的にする）
- 外部応答（XML 取り込みは実ファイルをテスト fixture として固定する）

### 3.3 golden の配置候補

`Q-008` の回答に従い、本 test-design 段階では配置場所を decide してよい。

- 境界結合テスト backend 側: Go の `testdata/` または `embed` ファイルとして `internal/apitest/` 配下に置くか、または test 固有の定数として `_test.go` 内に定義する
- 境界結合テスト frontend 側: TypeScript の `const` または `*.json` として `frontend/src/` 配下の境界テスト専用 directory に置く

golden の表現形式（JSON / Go struct / TypeScript const）は `implementation-scope` で確定する。
両側 test から同一の参照点として読めることが成立条件であるため、形式を一方に揃えるか、変換層を挟む設計を `implementation-scope` で判断する。

---

## 4. 単体テスト・シナリオテストへの観点変更要否

`boundary-test-REQ-001` の責務境界と既存テストを照合した結果、次のとおりと判定する。

### 4.1 backend 単体テスト

**変更なし。**

既存 `internal/service/master_dictionary_*_test.go`、`internal/usecase/master_dictionary_usecase_test.go`、`internal/controller/wails/master_dictionary_controller_unit_test.go` は Go package 内の公開振る舞いと分岐を証明しており、境界 DTO 値の semantic は責務にしていない。
境界結合テストの追加によって既存単体テストに観点追加・削除の必要は生じない。

### 4.2 frontend 単体テスト（gateway テスト含む）

**変更なし。**

`frontend/src/controller/wails/master-dictionary.gateway.test.ts` は frontend 解釈（`rec`/`edid` 非包含、`null` 応答時の型、binding 引数渡し）を証明しており、境界結合テストとの責務は分離されている。
BCT-MDC-004 の「不在 ID への null 応答」は frontend gateway test と観点が近いが、証明する主体（frontend 解釈 vs backend 実応答）が異なるため、既存 gateway test の削除や変更は不要である。

### 4.3 シナリオテスト（system E2E）

**変更なし。**

`tests/system/master-dictionary-management.spec.ts` は UI 上の操作 flow と表示結果を証明している。境界結合テストは DTO 値 semantic の固定を担い、UI 表示は対象外（`boundary-test-REQ-004`）であるため、既存 E2E テストへの観点変更は不要である。

---

## 5. pilot 観察ポイント（skill 改訂判断材料）

`boundary-test-REQ-006` で固定した「pilot 完了後の観察項目」を、テスト観点の実施過程で観察すべき事項として次に明示する。

### 5.1 test-design skill の入力フォーマット不足

- 本書で「観測 field」「backend 期待値」「frontend 期待値」を分けて記述したが、通常の E2E テスト観点表（CSV: `ID,関連UC,対象画面,前提条件,手順,期待値,備考`）では「対象画面」「手順」列が境界結合テストに適合しない。
- **観察ポイント**: pilot 実装で `test-design` skill の CSV 形式を強制した場合に情報欠落が起きるかどうか。起きた場合は「手順」列に method 呼び出し手順を記述する拡張か、または境界結合テスト専用の観点表形式が必要になる可能性がある。

### 5.2 implementation-module の decision table への追加工程

- 現状の `implementation-module` decision table には「シナリオテスト / 単体テスト / 観測ログ / 最終検証」が並ぶが、「境界結合テスト（backend 側）」と「境界結合テスト（frontend 側）」が独立工程として存在しない。
- **観察ポイント**: pilot 実装で「backend 側 test 作成」と「frontend 側 test 作成」が分かれた工程として decision table に追加を要したかどうか。また golden 管理（作成・検証・更新）が既存工程では表現しきれない工程として浮上したかどうか。

### 5.3 既存 skill 文言との衝突

- `tests-unit` skill は「公開振る舞い・分岐・エラー経路を証明する」と定義しており、境界結合テストの「DTO 値 semantic の両側整合」とは責務が異なる。
- `tests-scenario` skill は「業務シナリオの受け入れ条件を証明する」と定義しており、method 単位の golden 一致 assert とは異なる。
- **観察ポイント**: pilot 実装で既存 skill の文言をそのまま境界結合テストの実装に適用した際に衝突（責務の過不足、手順の不適合）が発生したかどうか。発生した具体箇所と内容。

### 5.4 golden 更新手順の独立性

- 既存 `tests-unit` や `tests-scenario` は期待値の更新を「仕様変更に伴うテスト修正」として扱うが、境界結合テストでは golden という中間成果物が存在し、その更新が両側 test に連動する。
- **観察ポイント**: pilot 実装で golden の更新手順が「仕様 → golden → backend test → frontend test」の順序で独立工程として成立したかどうか。片側だけ更新した場合に他方が失敗する状態が維持されたかどうか。

### 5.5 新規 skill 追加の閾値判断

`boundary-test-REQ-006` で定めた「既存 skill 改訂で吸収しきれない数が閾値を超える場合は新規 skill `tests-boundary` を追加する判断材料にする」について、上記 5.1〜5.4 の観察結果から次の判定を行う。

- 5.1〜5.4 のうち 2 件以下の衝突・不足: 既存 skill 改訂（案 A）で対応可能と判定する根拠になる
- 5.1〜5.4 のうち 3 件以上の衝突・不足: 新規 skill `tests-boundary` の追加（案 B）を検討する根拠になる

---

## 6. pilot case の配置先決定（Q-008）

`Q-008` の回答「case 表本体の配置場所は `test-design` 起動時に decide してよい」に従い、次のとおり配置先を決定する。

### 6.1 決定

- backend 側 test: `internal/apitest/` 配下に新規ファイルとして配置する（例: `master_dictionary_boundary_test.go`）。既存の `internal/apitest/` が Wails controller を入口とする受け入れ test の置き場所であり、bootstrap 済み controller を使う点で共通するため。
- frontend 側 test: `frontend/src/controller/wails/` 配下に新規ファイルとして配置する（例: `master-dictionary.boundary.test.ts`）。既存の `master-dictionary.gateway.test.ts` と同じ directory に置くことで、同じ wails mock 機構を参照しやすくなるため。
- golden 集合: backend 側は `internal/apitest/testdata/master-dictionary-golden/` または test ファイル内 Go struct として定義し、frontend 側は golden を import できる TypeScript ファイルとして同一 directory に置く。形式の最終決定は `implementation-scope` に委ねる。

### 6.2 根拠

- `internal/apitest/` への配置: `design-diff-diagram.md` 図 A で「backend API テストと同じ入口（Wails controller）を使うが、観測責務が異なる」と明示されており、入口の共通性から同じ directory に置く判断が整合する。
- `frontend/src/controller/wails/` への配置: 既存の gateway test と mock 機構を共有する（図 A 注記）ことを design で明示しているため、同じ directory に置くことで参照距離を最小にする。
- golden の単一性: 両側 test が同一 golden を参照する成立条件（`boundary-test-REQ-003`）を満たすため、import 可能な形式で置く。

---

## 根拠

- `source`:
  - `docs/exec-plans/active/master-dictionary-boundary-test/plan.md`
  - `docs/exec-plans/active/master-dictionary-boundary-test/detail-spec-diff.md`
  - `docs/exec-plans/active/master-dictionary-boundary-test/design-diff-diagram.md`
  - `docs/detail-specs/master-dictionary.md`
  - `docs/e2e-test-guidelines.md`
  - `docs/coding-guidelines-tests.md`
  - `frontend/src/controller/wails/master-dictionary.gateway.test.ts`
  - `frontend/src/controller/wails/master-dictionary.gateway.ts`
  - `frontend/src/application/gateway-contract/master-dictionary/master-dictionary-gateway-contract.ts`
  - `internal/controller/wails/master_dictionary_controller.go`
  - `tests/system/master-dictionary-management.spec.ts`
  - `internal/apitest/` 配下の既存 test 群（責務境界確認のため参照）
- `review`: 人間レビュー待ち
- `validation`: 未実行。`implementation-scope` 確定後に test 実装で検証する。
