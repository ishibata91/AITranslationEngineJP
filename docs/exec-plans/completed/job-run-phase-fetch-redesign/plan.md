# Task Plan: job-run-phase-fetch-redesign

- `workflow`: design-then-implement
- `status`: planned
- `lane_owner`: implement-lane（設計レビュー承認済み 2026-05-31、実装フェーズ進行中）
- `task_id`: job-run-phase-fetch-redesign
- `task_mode`: redesign
- `request_summary`: 単語翻訳画面（JobRunPage）の初回表示で処理対象一覧が0件になる不具合を、局所修正ではなく term/persona/body 3段階を貫く取得・表示フローの共通設計で作り直す。
- `goal`: JobRunPage の各段階で、処理対象一覧が初回表示で正しく件数分表示される。term/persona/body を同型の共通設計に抽象化し、同種不具合の再発を構造的に防ぐ。
- `constraints`: 設計判断を伴うため局所修正レーンでは扱わない。設計は専用レーン整備後にそのレーンで行う。bridge 自体は正常。backend は readiness（次フェーズ開始可否）の責務再配置のため変更する（2026-05-31 人間設計レビューでスコープ拡大、下記「スコープ拡大」参照）。
- `close_conditions`: 共通設計が人間承認され、実装後に下記「不足テスト」が green、実機でジョブ6初回表示に処理対象が出る。
- `execution_branch`: `codex/job-run-phase-fetch-redesign`（作成済み、`source_branch` から分岐）
- `source_branch`: `codex/fix-lucien-target-list-empty`（現ブランチ。前タスクの調査履歴を含む）
- `target_branch`: `master`

## 経緯（前タスク fix-lucien-target-list-empty からの引き継ぎ）

前タスクは fix-lane（局所修正）で7ラウンド原因究明したが症状を解消できず、設計レベルの作り直しが必要と判断して close した。前タスクの記録は `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/` にある（fix-decision.md 追補1〜5、attempts/obs-r1〜r6 の実機観測証跡を含む）。

## 症状（未解決）

- 操作: `dictionaries/Lucien.esp_Export.json` を読み込んだジョブ（例: jobId=6）の単語翻訳画面を開く。
- 結果: 進捗パネルは「AI 翻訳対象語 4,930（4931）」を表示するが、処理対象一覧パネルは初回表示で「0 件」「処理対象がありません」。
- リロードや別経路でも0件のまま。検索操作をすると50件表示される（手動取得経路は正常）。

## 確定している事実（前タスクの実機観測 obs-r1〜r6 より）

- backend / bridge は正常。`GetProcessingTargetList` を直接呼ぶと即座に items=50, totalCount=4931 を返す。実DBも 4931 行。
- 画面経由だと bridge 応答が約15秒かかる（直呼びの約60倍）。初回マウントで term/persona/body の3 controller が各々 summary/readiness/processingTarget の3本を `Promise.all` で同時取得し、最大9本の Wails bridge 呼び出しが同時発火して IPC が詰まることが疑われる（未確定）。
- データは bridge→store→viewModel まで届く（obs-r5/r6 で items=50 を store.update / subscribe callback で確認）。
- だが画面の `JobRunPage.svelte` の `currentProcessingTargetPageState`（`$derived.by`）が viewModel 更新に再評価で反応せず undefined のまま、または15秒の遅延中に画面が unmount/再マウントされ反映が落ちる。
- 進捗の母数（`summary.aiTargetCount`）と一覧（`getProcessingTargetList` の items）は別 bridge・別集計の独立経路。一致保証は元々ない。

## 試したが解消しなかった局所修正（すべて revert 済み）

- `$effect`(setJobId) と `onMount`(load) の二重起動除去・onMount 一本化。
- 取得を表示中段階のみに限定（9本→3本）。
- `currentProcessingTargetPageState` の derived の早期 return 除去・依存束縛。

いずれも単独では症状を解消できず、原因が複数層に絡む。局所パッチを積み上げる限界に達したため作り直す。

## 作り直し方針（設計レーンで確定する）

term/persona/body の3段階を貫く取得・表示フローの共通設計を作り、各段階はそれを使う形に抽象化する。設計レーンで次を確定する。

- 取得の起動主体と回数（表示中段階のみ／順序づけ／`Promise.all` 濫用の回避）。
- bridge 呼び出しの同時数制御（15秒遅延＝IPC 飽和の真因を確定し、直列化または必要分のみ取得で解消するか）。
- store → viewModel → 画面の reactive 反映を、derived の依存追跡が確実に成立する素直な形にする。
- 進捗母数と一覧の独立経路を設計上明示する。
- 連番競合検出（`processingTargetListRequestSequence`）は手動操作の競合防止に必要なため、撤廃せず役割を限定する。

注: 本 plan では設計・実装しない。設計レーンの整備後に、そのレーンで設計→人間承認→実装→検証を行う。

## 不足テスト（作り直し後に green にすべき検証資産。コードは tests/ に残置済み）

前タスクで追加し、本タスクへ引き継ぐ E2E シナリオテスト。

- ファイル: `tests/system/fix-lucien-target-list-empty.spec.ts`（E2E-LTLE-001/002/003）。
- 補助: `tests/system/support/scenario-wails-mocks.ts` の `termZeroAITargetJobId` オプション追加（E2E-LTLE-003 の母数0境界用）。
- 観点（`docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/test-design.csv` 由来）:
  - E2E-LTLE-001（正常）: 進捗母数1以上のとき、初回表示で処理対象一覧に行が1件以上、空状態にならない。
  - E2E-LTLE-002（境界）: 検索後リロードで初回0件へ戻らない。
  - E2E-LTLE-003（境界）: 進捗母数0のとき空状態を保持。
- 限界: モック E2E は同期解決のため IPC 飽和（実機の15秒遅延）を再現しない。実機固有症状の検出には実機ブラウザ確認が必要。作り直しの fail-test 戦略は設計レーンで決める。
- UC 差分: `差分なし`（既存 UC「処理対象を確認する」で説明可能。`completed/.../uc-diff.md`）。
- selector 不足: 処理対象件数・空状態・検索欄の固定 selector が画面設計未確定（`completed/.../data-testid-gaps.md`）。

## 現状の作業ツリー差分（本 plan 作成時点）

- プロダクトコード変更なし（局所修正は全 revert 済み）。
- `tests/system/fix-lucien-target-list-empty.spec.ts`（新規・不足テスト、残置）。
- `tests/system/support/scenario-wails-mocks.ts`（E2E 用 mock 追加、残置）。
- `scripts/dev/run-wails-agent-browser.sh`（wails ログを tee で stdout とファイル両方へ出す改善。本不具合とは独立）。
- `.claude/agents/fix_decider.md`（model 変更。利用者による）。

## 着手計画（implement-lane）

- `task_mode`: redesign（作り直し）。新規実装ではなく、既存の取得・表示フローを共通設計へ作り直す。
- `execution_branch`（提案）: `codex/job-run-phase-fetch-redesign` を `source_branch`（`codex/fix-lucien-target-list-empty`）から切る。前タスクの残置検証資産（E2E、mock）を引き継ぐため。
- 仕様変更/追加の有無: あり（取得の起動主体・回数・同時数制御・reactive 反映の規約を設計で確定する）。ただし利用者向け UC は差分なし（既存 UC「処理対象を確認する」で説明可能）。
- 画面変更の有無: 表示挙動（初回表示で処理対象一覧が件数分出る）の確定が論点。見た目の UX デザイン変更は薄い見込み。固定 selector（件数・空状態・検索欄）が未確定。
- 構造変更の有無: あり（取得フローと store→viewModel→画面の reactive 反映の作り直し）。
- backend 変更の有無: なし（backend/bridge は正常、plan 明記）。

### 今回作る成果物の要否

| 成果物 | 要否 | 理由 |
| --- | --- | --- |
| 詳細仕様差分（designer） | 要 | 取得フロー（起動主体・回数・同時数制御・連番ガード役割・reactive 反映規約）を仕様差分として固定する。 |
| 画面設計差分（designer） | 要 | 処理対象件数・空状態・検索欄の固定 selector と初回表示規則を確定する。 |
| 設計差分図（diagrammer） | 要 | 取得・表示フローの構造変更を図で固定し、人間設計レビュー前に揃える。 |
| 人間設計レビュー | 要 | 不変条件（常に）。 |
| 実装範囲（designer） | 要 | 不変条件（常に）。 |
| テスト設計（test_designer） | 要 | 不変条件（常に）。残置 E2E を fail-test 戦略へ組み込む。 |
| frontend 実装 | 要 | usecase 取得フローと JobRunPage.svelte の reactive 反映を作り直す。 |
| Storybook 入力確認・frontend 実装後人間レビュー | 暫定不要 | UI の裏側（取得フロー）変更が主。画面設計差分で UX デザイン変更が確定したら要に切り替える。 |
| backend 実装 | 要 | readiness の次フェーズ開始可否（`canStartNextPhase`）の責務再配置。backend は「次へ進めるか」の判断を返さず、データの事実状態だけ返す形へ見直す（スコープ拡大）。 |
| 統合境界実装（integration_implementer） | 要見込み | bridge 呼び出しの同時数制御に加え、readiness 応答 DTO の責務変更（`canStartNextPhase` の扱い）が gateway/DTO 境界に波及するため。 |
| シナリオテスト（implementation_scenario_tester） | 要 | 残置 E2E（E2E-LTLE-001/002/003）を green にする。 |
| 単体テスト（implementation_unit_tester） | 要 | 取得フロー責務（連番ガード対称化・取得回数）を証明する。 |
| 観測ログ追加（observability_implementer） | 要見込み | 実機 15 秒遅延（IPC 飽和の疑い、未確定）の原因分離。設計で真因が確定すれば省略可。 |
| 最終検証 | 要 | 不変条件（常に）。 |
| 実装後ブラウザ確認（browser_confirmation） | 要 | 実機でしか出ない症状（15 秒遅延・初回 0 件）の確認。 |
| 正本化判断 / 詳細仕様正本反映（docs_updater） | 設計次第 | 恒久仕様が確定し人間承認された分だけ正本反映。 |

注: 「設計次第」「要見込み」「暫定不要」は設計成果物の確定後に再判定する。

### スコープ拡大（2026-05-31 人間設計レビュー）

第1回の人間設計レビューで、readiness（次フェーズ開始可否）の責務配置に設計上の問題が指摘された。経緯と決定を記録する。

- 指摘: readiness 応答 `TermTranslationNextPhaseReadinessResponse`（`frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts:139`）は `canStartNextPhase: boolean` と `blockedReason` を持つ。「データ状態の事実」だけでなく「次へ進めてよいか」の遷移判断まで backend が決めて返している点が責務配置として不適切（人間判断）。
- 決定: readiness の責務再配置を今回 task に含める。backend は「次へ進めるか」を返さず、データの事実状態だけ返す形へ見直す。「次へ進めてよいか」はその事実状態から導く。backend 変更ありへスコープ拡大する。
- 第1回設計成果物への影響: 詳細仕様差分・画面設計差分・設計差分図を差し戻し、(a) 処理対象一覧の反映取りこぼし防止、(b) readiness 責務再配置、の両方を反映して更新する。「進捗の要約・準備状態・処理対象一覧の3つを同一連番ガードで反映」という表現も、粒度差（一覧は重い取得、readiness は可否判定）を踏まえて見直す。

注: 当初 plan の「backend/bridge は正常（変更不要）」と memory `job-run-phase-fetch-redesign` の記述は本決定で更新される。

#### 第2回設計レビューの追加決定（Q-004 / Q-005）

差し戻し後の詳細仕様差分で残った未決2件を人間が回答した。

- `Q-004` 回答: 操作可否までフロント導出へ広げる。次フェーズ開始可否だけでなく、操作可否（`actionEnablement`: 開始・中断・再開・リトライ・取消・成果物出力確認の各活性と理由）も backend から外し、事実状態からフロントが導出する。term/persona/body 全段階で同型に扱う。backend は操作可否・遷移可否の判断を応答に含めない。
- `Q-005` 回答: body の専用取得 `GetBodyTranslationOutputReadiness` を廃止し、成果物出力確認に必要な事実を段階要約取得へ一本化する。出力可否（`ready`）はフロント導出。取得経路を2本から1本へ統合する。

これにより本 task は「処理対象一覧0件不具合の取得・表示フロー作り直し」に加え、「フェーズ画面の可否判断（遷移可否・操作可否・出力可否）の責務を backend から frontend へ移す再配置」を含む。backend・gateway DTO・frontend application 層すべてに波及する。等価性条件（再配置前後で同じ事実入力に同じ可否結果）を実装範囲とテスト設計で担保する。

## 実装進捗（implement-lane）

- 設計レーン: 詳細仕様差分・画面設計差分・設計差分図を作成、人間設計レビュー承認済み（2026-05-31）。
- 実装範囲: `implementation-scope.md`（11引き継ぎ・4 wave）固定。テスト設計: `test-design.csv`・`test-design-unit.md`・`data-testid-gaps.md` 固定。
- wave-1（frontend 実装）完了: 取得フロー作り直し（表示中段階のみ最大2本・反映取りこぼし防止・開き直し再取得・ローディングレイヤー・derived 是正）と可否のフロント導出（term/persona/body）。frontend-local 全通過。
- `frontend 実装後人間レビュー`: 承認済み（2026-05-31、デザインレビュー ok）。Storybook レビュー中にローディング範囲を処理対象一覧領域からフェーズ画面全体オーバーレイへ拡大。
- `Storybook 後画面設計差分整合`: designer が screen-design-diff・detail-spec を実装事実へ整合、diagrammer が設計差分図 図4 を整合、test_designer が test-design の操作排他文言を整合。selector `<phase-prefix>-processing-target-loading` は据え置き。
- 確認対象 story を通常分類 `Screen Components/<段階名>/...` へ復帰済み。
- `合意済み frontend 保護` 固定（2026-05-31）。下記参照。
- wave-2（統合境界）完了: 可否値を DTO/gateway-contract から除去、body 専用取得廃止、事実状態伝送。
- wave-3（backend）完了: service/usecase/controller の可否導出撤去、事実状態返却、専用取得 endpoint 廃止。backend-local 通過。
- wave-4（テスト）完了: 型エラー解消、UT 観点（取得フロー・等価性）、E2E 8観点（fix-lucien）整合。開き直し連番破棄の実装欠陥（UT-REOPEN-001）を発見し frontend で修正。既存 E2E リグレッション（E2E-UC-048/049/052）はテスト待機不足が原因と確定しテスト側で修正。
- 観測ログ追加: 実機15秒遅延の原因分離ログ（取得起動・bridge 応答時間・反映・initialFetchDone 遷移）を usecase へ追加。
- 最終検証: backend-local 通過、frontend-local 通過（557テスト）。system test（E2E）は Wails dev サーバー未起動で環境制約（FAIL_ENVIRONMENT、再実行 `sh ./scripts/test/run-system-test.sh`）。coverage 68.7%<70% は変更前からの全体閾値で本 task では達成不能（backend 含む全体閾値）、本 task が下げたものではない。Sonar 接続エラーは環境制約。
- 現在: 実装後ブラウザ確認（実機での初回表示件数・15秒遅延解消の観測）進行。

### 合意済み frontend 保護

- 承認済み画面: JobRunPage と term/persona/body パネルの初回取得中フェーズ画面全体ローディングオーバーレイ（`phase-loading-overlay`）、件数あり/空状態表示、可否のフロント導出による操作ボタン活性・理由表示。
- 表示規則の正本: `screen-design-diff.job-run.md`（整合済み）と `detail-spec-diff.md`（整合済み）。
- 確認済み Storybook 状態: loading-layer stories 通常分類復帰済み（lint・build-storybook 通過）。
- 変更禁止範囲: wave-2/3/4 で UI 表示・画面文言・layout・style・承認済み画面設計差分を越える変更をしない。越える必要が出たら `frontend 実装` 再実行入力か人間返却を固定する。selector 値は変更しない。

### Storybook 入力確認

- 起動: `npm --prefix frontend run storybook`（`http://localhost:6008/` 固定、別 port で追加起動しない）。
- 確認対象 story（作業中分類 `Review/Changed Components/<段階名>/...`、loading-layer stories）:
  - `ProcessingTargetLoading`（初回取得中ローディングレイヤー、最前面表示と操作排他）
  - `ProcessingTargetWithItems`（件数あり表示）
  - `ProcessingTargetEmpty`（空状態表示）
  - term/persona/body の3段階同型。
- 確認観点: ローディングレイヤーの最前面表示と検索・ページ・行展開の操作排他、件数あり/空状態の表示規則、可否のフロント導出による操作ボタン活性・理由表示。
- 固定 selector: `<phase-prefix>-processing-target-loading`（新規）/`-total`/`-empty`/`-search-input`/`-row.<id>`。

## Outcome

- 設計レーン整備後に本 plan を設計レーンへ渡す。
