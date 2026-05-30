# Task Plan: job-run-phase-fetch-redesign

- `workflow`: design-then-implement
- `status`: planned
- `lane_owner`: 未定（設計レーン整備後に設計レーンへ）
- `task_id`: job-run-phase-fetch-redesign
- `task_mode`: redesign
- `request_summary`: 単語翻訳画面（JobRunPage）の初回表示で処理対象一覧が0件になる不具合を、局所修正ではなく term/persona/body 3段階を貫く取得・表示フローの共通設計で作り直す。
- `goal`: JobRunPage の各段階で、処理対象一覧が初回表示で正しく件数分表示される。term/persona/body を同型の共通設計に抽象化し、同種不具合の再発を構造的に防ぐ。
- `constraints`: 設計判断を伴うため局所修正レーンでは扱わない。設計は専用レーン整備後にそのレーンで行う。backend/bridge は正常（変更不要）。
- `close_conditions`: 共通設計が人間承認され、実装後に下記「不足テスト」が green、実機でジョブ6初回表示に処理対象が出る。
- `execution_branch`: 未定
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

## Outcome

- 設計レーン整備後に本 plan を設計レーンへ渡す。
