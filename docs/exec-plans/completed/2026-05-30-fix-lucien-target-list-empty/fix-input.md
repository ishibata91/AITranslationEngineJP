# 修正実行入力: fix-lucien-target-list-empty

人間修正レビューで承認済み（approve）。本書は実装 agent、シナリオテスト agent、browser_confirmation へ渡す単一入力である。

## 承認状態

- 修正方針判断: 承認済み。
- UC 差分候補: 承認済み（分類=差分なし、UC 正本改訂不要）。
- E2E テスト観点差分: 承認済み（分類=追加候補あり、観点 3 件）。

## 確定原因（fix-decision.md 追補5・obs-r6 実機観測、最終確定）

真因は `JobRunPage.svelte` 367-384行の `currentProcessingTargetPageState`（`$derived.by`）が **`viewModel` の更新に再評価で反応していない** こと（Svelte 5 reactivity / 依存追跡の不成立）。

obs-r6 実測（ジョブ6初回表示）:
- store 反映後、`viewModel.processingTargetPageState.items=50, totalCount=4931` に正しく更新される。
- だが同じ瞬間 `currentProcessingTargetPageState`（derived）は `undefined` のまま。初回評価時の undefined をキャッシュし、`viewModel` 更新で再評価されない。
- よって Panel に `processingTargetPageState=undefined` が渡り、Panel は items を出せず空表示になる。

データは bridge→store→viewModel まで完全に届いている（items=50）。落ちるのは最後の derived 1個だけ。過去ラウンドの「取得経路・二重起動・IPC飽和・連番競合」はすべて無関係（データは常に届いていた）。

注記: 過去の二重起動除去・onMount一本化・表示中段階のみ取得の変更は、真因ではなかったが副作用のない整理として残す（fix-decision.md 追補参照）。今回の真因修正は別物。

- store には4931件が正しく書き込まれ、presenter にも4931件届くが、画面は0件のまま確定する（取得失敗ではなく描画反映前の競合）。obs-r4 で `after store.update totalCount=4931`・`presenter.toViewModel totalCount=4931` を観測しつつ画面は「処理対象がありません」。
- 過去の「二重起動（6本）」は一段階で、それを除いても persona/body の取得が加わり9本同時で飽和する。表示中でない persona/body まで初回取得する点が過剰。
- 進捗パネル母数は正しい（backend 正常、実機 bridge と実DBで 4,931件を確認）。

## 採用する修正方針（恒久修正・真因確定後の最終版）

**初回マウントの bridge 同時起動数を削減する。onMount の初回取得を、表示中の段階（`resolveInitialPhasePage(selectedJobTarget)` が返す段階。ジョブ6なら term）だけに限定する。** persona/body は `currentPhasePage` がその段階へ切り替わった時点で初めて取得する。

- 初回の bridge 呼び出しが9本→3本（表示中段階の summary/readiness/processingTargetList）に減り、IPC 飽和が解消する。
- 既に適用済みの「`$effect` から取得起動を除去し onMount に集約」は維持する（二重起動除去）。その上で onMount の取得対象を表示中段階のみに絞る。
- 段階切替時の取得は既存の phase page 切替経路に沿って、その段階の controller が未取得なら取得する形にする。
- 連番競合検出（`processingTargetListRequestSequence`）は手動操作の競合防止に有効なため撤廃しない。新しい状態値は追加しない。

## 禁止する修正

1. 進捗パネル母数や backend の SQL・DTO・bind 配線の変更（backend は正常）。
2. 処理対象パネルのフィルタ初期値・検索機能の削除や、検索初期値の細工による症状回避。
3. `processingTargetListRequestSequence` による取得競合検出そのものの撤廃。
4. 新しい phase 状態値、処理対象専用の重複 state フィールドの追加。
5. 症状の表示だけを隠す対症療法（初期値を検索クエリで上書きして強制取得させる等）。
6. bridge 呼び出し総数は変えずに直列化だけで対応する対症療法。
7. persona/body の取得をスキップする条件分岐で表示問題を隠す対症療法（段階切替時に正しく取得する前提でのみ「表示中段階のみ初回取得」を許容。取得を恒久的に欠落させない）。

## 影響ファイル候補（確定ではない。実装 agent が責務境界で最小変更を選ぶ）

- `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`

## fail-test 戦略（人間判断: 二重起動を単体で赤にする）

症状（初回0件）は実機 bridge の IPC 並列制限・タイミング依存であり、単体でもモック E2E でも再現できないことが確認済み（fix_decider・単体テスターの一致した結論）。`fetchSummaryAndReadiness` は `Promise.all` 一括待機後に1回 store.update するため、最大 seq の結果（正常件数）が必ず残り、連番 skip だけでは 0 件にならない。0 件の真の機序は 6本同時呼び出しによる `Promise.all` ハング（症状B、実機固有）。

よって fail-test は症状ではなく**根本原因（初回ロードの二重起動）**を対象にする。

1. 単体 fail-test: 初回ロードで同一 jobId に対し `fetchSummaryAndReadiness`（および処理対象取得）が **2回起動する**ことを検出し、未修正コードで**赤**にする。
2. 実装修正: frontend_implementer が二重起動を除去（初回ロードで取得起動を1回にする）。
3. 回帰防止: 単体テストが緑（1回）になる。既存 E2E-LTLE-001/002/003（モックで green）は回帰防止として維持。

## 単体テスト追加入力（implementation_unit_tester 向け）

- 対象責務: `term-translation-phase.usecase.ts` の初回ロード経路（`load` / `setJobId` / `fetchSummaryAndReadiness`）。
- fail-test として先に追加: 初回ロード相当の呼び出し（`onMount` の `load()` と `$effect` の `setJobId(jobId)` の両起動を模す）の後、ゲートウェイの取得系メソッド（`getTermTranslationPhaseSummary` または処理対象取得）が同一 jobId に対し **2回以上呼ばれる**ことを観測点にし、「初回ロードで取得が二重起動しない（1回である）」を assert する。
  - 未修正コード: 二重起動して取得が2回 → 赤。
  - 修正後: 1回 → 緑。
- 観測手段: ゲートウェイをスパイ（呼び出し回数を記録するテストダブル）にして、初回ロードでの呼び出し回数を数える。Promise 解決順の細工は不要（回数だけ見る）。決定的。
- 注意: 手動操作（検索・ページ）での取得起動を禁止する assert にはしない。あくまで初回ロードの二重起動だけを対象にする。連番競合検出の分岐は壊さない。
- 検証コマンド: `python3 scripts/harness/run.py --suite frontend-local`。

## シナリオテスト（既存、回帰防止として維持）

- 既存 `tests/system/fix-lucien-target-list-empty.spec.ts` の E2E-LTLE-001/002/003 を回帰防止として維持する。
- モック環境では確定原因が再現しないため green のまま。fail-test の役割は単体テストが担う。
- 開始点・観測点は承認済み `test-design.csv` に従う。

## 実装後ブラウザ確認入力（browser_confirmation 向け）

- 確認 URL: `http://localhost:34115/#translation-management`。
- 操作経路（fix_decider の再現手順を共有）:
  1. 翻訳管理画面でジョブ6を選択し「現在の翻訳段階へ進む」を押下。
  2. 単語翻訳段階画面の初回表示を待つ（検索・ページ操作なし）。
  3. 処理対象パネルの件数と行、進捗パネル母数を観測。
  4. 検索語を入力 → 件数変化を確認 → リロードして初回表示を再観測。
- 修正前の問題状態: 初回表示で処理対象パネルが0件「処理対象が見つかりません。」。進捗母数は4,931。リロードで再び0件。
- 修正後に満たすべき期待状態: 初回表示で処理対象一覧に行が1件以上表示され、進捗母数と整合する。リロード後も0件に戻らない。
- 証跡置き場: `tmp/agent-browser/`。

## 進行メモ

- fix-lane DAG 順序: シナリオテスト追加証跡 → 実装修正証跡 → 単体テスト追加証跡 → 実装後ブラウザ確認 → ハーネス実行 → 作業 commit。
- 実装種別: frontend（`frontend_implementer`）。
