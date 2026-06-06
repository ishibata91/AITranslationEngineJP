# plan: master-dictionary-boundary-test

## 依頼要約

frontend-backend 境界結合テストの pilot として、MasterDictionary 機能を対象に「backend apitest と frontend test が同じ golden を共有する」枠組みを確立する。

MasterDictionary は直近 commit `ca5f8480`（domain state projection を backend が expose し action enablement を frontend が導出する大規模 refactor）後に、唯一手動動作確認済みで境界を信用できる対象である。本 pilot で「frontend-backend 境界結合テスト」の層を repo に導入し、backend 側 contract 検証 test と frontend 側 fixture 消費 test が同一 golden を介して契約を固定する形を実証する。

枠組み確立後の後続作業として test 戦略 docs 整備（候補 B）を予定するが、本 task の範囲には含めない。

## 背景

- 直近 commit `ca5f8480` で「backend が状態を expose、frontend が action enablement を導出」へ大量の境界をまとめて refactor した
- 結果として、境界 DTO の値 semantic が変わった時に検出する層が repo に存在しないことが品質管理側からの指摘で明らかになった
- 既存テスト層: backend 単体 (`internal/*/*_test.go`)、backend api test (`internal/apitest/`)、frontend 単体 (`frontend/src/**/*.test.ts`)、system E2E (`tests/system/`)。frontend gateway test は wails runtime を mock しており、backend が返す DTO 値と frontend 側解釈を一対の契約として直接 assert する層が欠落している
- scenario / UC level E2E は粒度が高すぎ、境界の値変化のうち UI に露出しないものは届かない（脳筋アプローチ）

## 作業 branch

- 作業 branch: `claude/master-dictionary-boundary-test`
- 分岐元 branch: `master`
- 分岐元 commit: `ba04db513bbb44a51e78aa983805973f4f0a5a73`

## 事前変更

本 plan の入口で、不要となった内部結合テスト folder を削除する変更を branch 内に含む（preparation-module 段階で実施）。

- `internal/integrationtest/` folder 削除（4 files、2205 行）。`go test` は通っていたが、`SCN-SMR-*` という古い設計文書由来の scenario ID 名で書かれており、現状の設計判断と乖離していると user 判断で削除
- `.go-arch-lint.yml` から `integrationtest` component / deps 定義を削除
- `internal/apitest/README.md` から `internal/integrationtest` への相互参照を削除

削除目的は、「frontend-backend 境界結合テスト」を新規層として導入する際に、既存の `internal/integrationtest`（SQLite + repository の内部結合）との名前空間衝突を解消し、「結合テスト」という呼称を境界結合テストに割り当てられる状態にすることである。

## 想定 Y/N 評価（design-module 入口）

| # | 想定 | Y/N | 根拠 | 参照 |
| --- | --- | --- | --- | --- |
| 1 | 仕様変更または仕様追加がある | Y | 「frontend-backend 境界結合テスト」という新規 test 層と、その導出方法・書き方を開発規約として追加する。test 戦略の仕様追加に該当する | `docs/coding-guidelines-tests.md`、`docs/architecture.md` |
| 2 | 画面変更がある | N | プロダクト画面は触らない | - |
| 3 | 内部構造変更がある | Y | 境界結合テスト用の test 層を新規導入し、配置 path、arch lint、skill/workflow 組み込みを変える。プロダクト production code の内部構造は触らないが、test 層構造を新規に追加する | `internal/apitest/`、`frontend/src/`、`.go-arch-lint.yml`、`.claude/skills/` |
| 4 | 画面の表示変更がある | N | 表示変更なし | - |
| 5 | frontend ロジック変更がある | N | state、API、Wails bridge、ルーティング、副作用、フォーム validation を触らない | - |
| 6 | backend 変更がある | N | production Go code を触らない。apitest を追加するのみ | - |
| 7 | frontend と backend を接続する | N | 新規接続なし。既存接続を境界結合テストで検証するのみ | - |
| 8 | 実装済み責務を独立に証明したい | Y | MasterDictionary の境界既存実装を境界結合テストで証明する | `frontend/src/controller/wails/master-dictionary.gateway.test.ts`、`tests/system/master-dictionary-management.spec.ts` |
| 9 | 実行時にしか確定しない値または原因分離が要る分岐がある | N | 観測ログ追加対象ではない | - |

## decision table 適用結果

| artifact | 要不要 | 根拠 |
| --- | --- | --- |
| 詳細仕様差分 | 要 | 想定 1 Y（境界結合テストの導出方法・書き方・skill/workflow 組み込みを仕様として追加） |
| 画面設計差分 | 省略 | 想定 2 N |
| 設計差分図 | 要 | 想定 3 Y（test 層構造、apitest と frontend test の golden 共有 flow を図示） |
| 人間設計レビュー | 要 | モジュール常時必須 |
| 実装範囲 | 要 | モジュール常時必須 |
| テスト設計 | 要 | モジュール常時必須（境界結合テストの観点表を pilot 範囲で作る） |

## 後続モジュールへの引き継ぎ予定論点

- skill/workflow への組み込みは詳細仕様差分の一部として扱う。新規 skill が必要か、既存 skill（test-design、tests-scenario、tests-unit、implementation-module、design-module）の改訂で足りるかを判定する
- skill/workflow の改訂が確定した場合、別 task として `workflow-contract-maintenance` の起動を提案する（本 task では docs 上の方針固定までを担う）
- observation-record の暫定推奨は **案 B（新規 skill `tests-boundary` 追加）** を起点に、AI service 境界 / SQL repository 境界への次適用見通しも踏まえて `workflow-contract-maintenance` 起動時に最終決定する。

## implementation-module 完了記録（2026-06-05）

### decision table 結果

| 想定（plan.md） | 評価 | artifact 機械的判定 | 本 task 暫定運用 |
| --- | --- | --- | --- |
| 想定 5: frontend ロジック変更 | N | frontend ロジック実装 不要 | 不要（test code のみ追加） |
| 想定 6: backend 変更 | N | backend 実装 / 単体テスト 不要 | 不要（test code のみ追加） |
| 想定 7: frontend と backend を接続 | N | 統合境界実装 / シナリオテスト 不要 | 要（暫定運用、user 承認済み）: 既存接続の独立証明として境界結合テストを置く |
| 想定 8: 実装済み責務を独立に証明 | Y | 単体テスト 要 | 暫定運用: tests-unit ではなく境界結合テストで証明、`tests-scenario` 寄せで実装 |
| 想定 9: 実行時にしか確定しない値 | N | 観測ログ追加 不要 | 不要 |

### 完了 handoff

| wave | handoff | 担当 agent | 完了状態 |
| --- | --- | --- | --- |
| wave-1 | handoff-shared-golden | integration_implementer | 完了（golden 11 + 両側 helper） |
| wave-1 | handoff-boundary-api-draft | integration_implementer | 完了（境界 API 仕様 draft、6 method 網羅） |
| wave-2 | handoff-boundary-test-backend | implementation_tester | 完了（10 test pass、片側書き換え検出手動確認済み） |
| wave-2 | handoff-boundary-test-frontend | implementation_tester | 完了（10 test pass、片側書き換え検出手動確認済み） |
| wave-3 | handoff-observation-record | implementation_tester | 完了（6 観察項目記録、暫定推奨 案 B） |

### 最終検証結果（`python3 scripts/harness/run.py --suite all`）

| suite | 結果 | 備考 |
| --- | --- | --- |
| docs structure | PASS | - |
| lint:backend | PASS（修正後再実行） | 初回 gofmt 失敗 → `internal/apitest/` 配下 2 file を `gofmt -w` で整形して解消 |
| lint:frontend | PASS（修正後再実行） | 初回 knip が `frontend/src/controller/wails/boundary-golden-loader.ts` を unused 認識 → `frontend/package.json` の knip ignore に追加して解消 |
| test:backend | PASS | 境界結合テスト backend 10 件含む全 pass |
| test:frontend | PASS | 境界結合テスト frontend 10 件含む 55 file 646 test 全 pass |
| test:system | FAIL（本 task 範囲外） | 9 件 fail、本 task の base commit `ba04db513` 由来 |

### 本 task 範囲外の system test failure（user 承認 2026-06-05、別 task で扱う）

本 task の変更は production code を 1 行も触らず test code + docs のみで構成されており、failure 経路は本 task と独立する。

| failing test 群 | 由来 commit |
| --- | --- |
| `tests/system/refactor-action-enablement.spec.ts`（RAEF-E2E-001/005/006/009/014/015） | `ca5f8480 refactor: derive action enablement on frontend` |
| `tests/system/translation-phases.spec.ts` E2E-UC-052 | 同上 |
| `tests/system/fix-phase-ai-settings-pill.spec.ts` E2E-UC-057-B | 同上 |
| `tests/system/job-run-shell.spec.ts` E2E-UC-046 | 同上 |

これらは本 task の base commit ba04db513 時点で既に存在し、本 task の `internal/apitest/` 配下追加 / docs 追加 / `frontend/package.json` knip ignore 1 行追加とは independent。別 task として扱う。

### 観察事項（observation-record に集約済み、`workflow-contract-maintenance` 入力）

- 暫定運用（`シナリオテスト` 寄せ）の妥当性確認
- 6 観察項目（test-design 入力フォーマット不足 / decision table 追加工程 / 既存 skill 文言衝突 / golden 更新手順独立性 / 閾値判断 / 片側書き換え検出自動化）

### finalization-module 正本化判断（2026-06-06）

- **仕様追加対象**: 「境界結合テスト」テスト種別の定義、責務分離、構成要素、既存テスト種別との境界
- **正本反映先**: `docs/detail-specs/boundary-integration-test.md`（新規作成）
- **反映元**: `boundary-integration-test.spec.md`（active plan folder 内、5 章 / 85 行）
- **人間承認状態**: 承認済み（2026-06-06、user 「おk」明示承認、spec の 5 章構成と内容を恒久仕様として確定）
- **判断結果**: 恒久仕様として承認、正本反映を実施
- **後続課題に切り出すもの**:
  - 別 task 候補 0: 既存境界の結合テスト一斉展開（rollout、12 feature 対象）
  - 別 task 候補 1: master-dictionary 画面設計の更新（4 件）
  - 別 task 候補 2: master-dictionary 現実装と golden / test の追従（4 件）
  - workflow-contract-maintenance: observation-record の 7 観察項目、暫定推奨 案 B（新規 skill `tests-boundary` 追加）

### finalization-module 引き継ぎ事項

- **正本反映対象（唯一）**: `boundary-integration-test.spec.md`（2026-06-06、6 章構成に絞り込み）→ 新規 `docs/detail-specs/boundary-integration-test.md`。境界結合テスト種別の「定義、責務、最小規約」のみを含み、書き方・命名・ID 規約・CI 統合などの運用詳細は spec の範囲外。
- **active plan folder 内に残す（正本化しない）**:
  - `boundary-api.master-dictionary.contract.ts`: master-dictionary 個別の形式 draft。pilot 試作として保持
  - `boundary-api-draft.master-dictionary.md`: 突合結果。別 task 候補の入力として保持
  - `next-task-plan-draft.boundary-test-rollout.md`: rollout 計画 draft。別 task 起動時の入力として保持
- 検証 partial pass の扱い: 本 task 範囲内全 pass、system test 範囲外 fail は別 task。
- workflow-contract-maintenance への引き継ぎ強化: observation-record の観察項目 7（user 気付き）を含む 7 項目全てを別 task の入力にする。暫定推奨は案 B（新規 skill `tests-boundary` 追加）を起点に、形式 / 意味 / 表示の責務分離（形式 = contract.ts、意味 = UC、表示 = 画面設計）を skill 設計に組み込む方向で再判断する。

### 責務分離の確定（2026-06-06 user 承認）

| 持ち場 | 真とする file |
|---|---|
| 形式（型、必須、null 許容、値域） | `boundary-api.master-dictionary.contract.ts` |
| 意味（field が業務的に何を指すか） | 画面仕様 / UC（`docs/usecases/uc-master-dictionary.md`） |
| 状態遷移規約 | UC |
| 表示文言、UI 構造 | 画面設計（`docs/screen-design/screens/master-dictionary.md`） |

意味の重複を排した結果、API 仕様書（補助文書）は「突合結果のみ残す」形に縮退した。

### 別 task 候補 0: 既存境界の結合テスト一斉展開（rollout）

本 task の master-dictionary pilot を、repo 内の既存境界全体に展開する rollout 計画 draft を `next-task-plan-draft.boundary-test-rollout.md`（2026-06-06 作成）として固定した。

- 対象 feature: 13 件（master-dictionary は pilot 完了済み、12 件が新規対象）
- 進行順序（user 指示 2026-06-06）:
  1. Phase 1: 既存実装と画面設計の差異を埋める
  2. Phase 2: 整合した画面設計と実境界を突合して、結合テストを作る
- task 分割: 案 B（feature ごとに別 task）を推奨
- 親方針正本化候補: `docs/detail-specs/boundary-integration-test-rollout-plan.md`
- 起点 spec: `boundary-integration-test.spec.md`（本 task で確定）

### 別 task 候補（あるべき仕様への追従、本 task 範囲外）

本書（あるべき仕様）の確定により、画面設計と現実装の双方が乖離する状態になった。次の 2 系統を別 task で扱う前提とする。

**別 task 候補 1: 画面設計の更新**
- 対象 file: `docs/screen-design/screens/master-dictionary.md`
- 修正項目（API 仕様書 9.2 節に対応）:
  - `[6]` 辞書一覧パネル metadata から「メモ」と「登録元」を外す
  - `[10]` 詳細パネルから「由来」を外す
  - `[8]` 辞書一覧各行から「由来」を外す（見出しと 4 列構成に揃える）
  - `[11]` 新規登録・更新モーダル入力項目から「由来」を外す。E2E 固定 selector の `master-dictionary-entry-origin-input` も削除候補
- 起点モジュール: `storybook-module`（画面設計差分があるため）

**別 task 候補 2: 現実装と golden / test の追従**
- 対象 file（API 仕様書 9.3 節に対応）:
  - backend Go DTO: `MasterDictionaryEntrySummary` / `MasterDictionaryEntryDetail` / 各 payload 系 DTO から `origin` / `note` を削除
  - backend service / repository: `toEntryDetailDTO` 内の固定文言設定を削除、`origin` 列の DB schema 扱いを決定
  - frontend gateway: `master-dictionary.gateway.ts` の type 定義から `origin` / `note` を削除
  - 本 task で追加した golden 11 件: `note`、`origin` field 値を削除
  - 本 task で追加した境界結合テスト 20 件: `note`、`origin` field の assert を削除
- 起点モジュール: `design-module`（仕様変更があるため）→ `implementation-module`

順序: 本書（あるべき仕様）の正本化が先行し、次に別 task 候補 1 / 2 が並列で追従する。

### 本 task で追加した test code の位置付け

本 task で wave-1 / wave-2 で追加した golden 11 件と境界結合テスト 20 件は **現実装ベース** で書かれており、本書（あるべき仕様）から見て古い状態になる。
本 task の主目的（frontend-backend 境界結合テストの枠組み確立）には影響しない（pilot として枠組みが動作することは validation 通過で証明済み）。
別 task 候補 2 で「あるべき仕様」に追従更新する想定。



