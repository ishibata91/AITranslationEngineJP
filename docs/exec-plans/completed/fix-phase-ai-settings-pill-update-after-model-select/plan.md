# fix-phase-ai-settings-pill-update-after-model-select

## 依頼要約

3 phase（単語翻訳 / ペルソナ生成 / 本文翻訳）の AI 設定パネルで、AI サービスとモデルを選択しても AI 設定 pill が「設定未完了」のまま、開始ボタンも disabled のまま固定できない。前 task `fix-phase-ai-model-list-empty` で fake モデル一覧の取得経路は修正済み、本 task はその後段のモデル選択 → 処理方式選択 → AI 設定固定 → 開始ボタン有効化までを画面操作で踏める状態にする。観測ログ駆動で確定原因と修正方針を固定する。

## 分岐元

- 分岐元 branch: `master`
- 分岐元 commit: b3532ad4d3f9abf43b3af0f725187c23e2fc2b0f

理由: 前 task の finalize 完了後に master から新規 branch を切る。前 2 task の変更は master に取り込み済み。

## 作業 branch

- `claude/fix-phase-ai-settings-pill-update-after-model-select`

## 人間観測記録（前 task 実画面確認時に固定）

- 対象環境: `npm run dev:wails:run`、`http://localhost:34115`、ジョブ#3、単語翻訳フェーズ
- 観測 1: AI サービス select に Gemini / LM Studio / xAI が表示される
- 観測 2: Gemini 選択 → 「モデル一覧を更新」ボタン押下 → モデル select に `fake-model` 選択肢が表示される（前 task で修正済み）
- 観測 3: モデル select で `fake-model` を選択 → AI 設定 pill が「設定未完了」のまま変わらない
- 観測 4: 開始ボタンが disabled のまま、blocked reason に「実行設定が未構成のため開始できません」が表示される
- 期待との差分: モデル選択後は AI 設定 pill が「固定可能」または「固定済み」に切り替わり、開始ボタンが有効化されるべき。

## 関連実装

- 前 task `fix-phase-ai-model-list-empty` で `FakeSecretStore` + `SECRET_BACKEND=fake` を導入し、credential 解決問題は解消済み。
- 前々 task `fix-term-translation-model-settings-empty-fixed` で frontend 経路（`saveAISettings`, `availableProviders`, `availableModels`, `refreshModelList`）を追加したが、provider と model を順次選択した時に saveAISettings が正しい値で呼ばれず viewModel が「設定未完了」のまま固まる仮説がある。
- 影響候補ファイル（要 investigation）:
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` の `handleProviderChange` / `handleModelChange` / `handleExecutionModeChange`
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller.ts` の `saveAISettings`
  - 3 phase 同型なので persona-generation / body-translation も対象

## 想定 Y/N 評価（investigation-module 入口時点で確定）

| 想定 | Y/N | 根拠 / 参照 |
| --- | --- | --- |
| 仕様変更または仕様追加がある | N | `docs/screen-design/screens/term-translation-phase.md` でモデル選択後に pill 固定可能・開始ボタン有効化される導線は定義済み。期待挙動が満たされていない状態であり、新規仕様判断は不要。 |
| 画面変更がある | N | layout、文言、表示構造、style の変更は依頼に含まれない。 |
| 内部構造変更がある | Y | `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte:148-185` の `saveAISettings` と handler、または `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts:210` の `saveAISettings` 経路のいずれかで viewModel 反映が失敗する。 |
| 画面の表示変更がある（storybook-module の TRIGGER 軸） | N | pill / select / ボタン文言は変更不要。表示更新の不具合は state 経路の問題に帰着する想定。 |
| frontend ロジック変更がある（implementation-module の TRIGGER 軸） | Y | panel handler、controller の `saveAISettings`、usecase の store 反映、presenter の `isExecutionConfigured` 派生のいずれかに修正が入る想定。 |
| backend 変更がある | N（暫定） | 前 task で backend `SECRET_BACKEND=fake` 経路は通過済み。本 task の症状は frontend 内の state 反映までで説明できる想定。fix_decider の観測で覆る場合は本欄を更新する。 |
| frontend と backend を接続する | N（暫定） | `saveAISettings` は usecase 内 store 更新で完結する想定。Wails bridge 追加は不要見込み。観測で覆る場合は本欄を更新する。 |
| 実装済み責務を独立に証明したい | Y | controller / presenter / usecase の `saveAISettings` 後の `isExecutionConfigured` 反映を単体テストで証明する。 |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | provider/model/executionMode の選択順序、Svelte 5 `$derived` の再評価、usecase の store 反映と presenter 派生の組み合わせは実行時にしか確定しない。 |

「仕様変更または仕様追加がある」が N のため、investigation-module を継続する。`backend 変更がある` と `frontend と backend を接続する` は暫定 N とし、`fix_decider` の観測ログ駆動検証で確定する。

## Wails 接続対象

- 起動 command: `npm run dev:wails:run`
- 接続先: `http://localhost:34115`

## 修正方針判断

- 成果物: [fix-decision.md](./fix-decision.md)
- 確定原因: `internal/bootstrap/app_controller.go` で 3 phase 全ての service 組み立て時に `JobPhaseAISettingsRepository` の注入が欠落している。`SaveTermTranslationPhaseAISettings` が `phase ai settings repository is not configured` エラーを返すため、summary の `aiSettings` が更新されず、pill が「設定未完了」固定になる。
- 採用する修正方針: bootstrap に `WithTermTranslationJobPhaseAISettings`, `WithPersonaGenerationPhaseAISettingsRepository`, `WithBodyTranslationJobPhaseAISettingsRepository` を追加して repository を注入する（backend のみ）。secondary として presenter の `buildModelOptions` placeholder 重複を解消する（frontend）。

## テスト設計成果物

- UC 差分候補: [uc-diff-candidates.md](./uc-diff-candidates.md)
  - 分類サマリ: `記述不足` を含む。`新規判断必要` は含まない。
  - 対象 UC: 翻訳段階を開始する（`uc-translation-management.md`）
  - 不足: AI 設定保存 → pill 変化 → 開始ボタン有効化フローの主シナリオ記述が不足。
- E2E テスト観点差分: [e2e-test-aspect-diff.md](./e2e-test-aspect-diff.md)
  - 分類サマリ: `追加候補あり`
  - 追加候補: `E2E-UC-056`（単語翻訳）、`E2E-UC-057`（NPC ペルソナ生成）、`E2E-UC-058`（本文翻訳）
  - 注記: 承認時は `E2E-UC-053/054/055` だったが、既存 spec での ID 衝突を回避して `E2E-UC-056/057/058` に採番し直した

## 人間修正レビュー

- 概要資料: [human-review-overview.md](./human-review-overview.md)
- 承認日: 2026-06-02
- 承認内容: 修正方針判断、UC 差分候補（記述不足 2 件、`新規判断必要` なし）、E2E テスト観点差分（追加候補 3 件、`判断不足` なし）の 3 成果物を一括承認。

## 修正実行入力（implementation-module へ引き継ぎ）

### 承認済み修正方針

- 主因修正（backend）: `internal/bootstrap/app_controller.go` で 3 phase 全ての service 組み立て時に `JobPhaseAISettingsRepository` を注入する。
  - `WithTermTranslationJobPhaseAISettings`
  - `WithPersonaGenerationPhaseAISettingsRepository`
  - `WithBodyTranslationJobPhaseAISettingsRepository`
- secondary 修正（frontend）: 3 phase の presenter `buildModelOptions` から `{value: "", label: "選んでください"}` placeholder を削除し、`AIModelSelectionCard` 側固定 placeholder と重複させない。主因修正後の検証で取捨選択する。

### 禁止する修正

- frontend で `isExecutionConfigured` を強制 true にする分岐の追加
- `saveAISettings` のエラー握り潰し + summary 強制再取得
- 新規状態フラグ（`localAISettingsSaved` 等）の追加
- `SaveTermTranslationPhaseAISettings` の mock / fake 挿入による backend エラー隠蔽

### 影響ファイル候補

| ファイル | 種別 |
| --- | --- |
| `internal/bootstrap/app_controller.go` | 主因（必須） |
| `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts` | secondary |
| `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts` | secondary |
| `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts` | secondary |

### 承認済み UC 差分

- 分類: `記述不足` 2 件、`新規判断必要` なし、`差分なし` 1 件。
- 補完対象 UC: 「翻訳段階を開始する」（`docs/usecases/uc-translation-management.md`）。
- 記述不足箇所: (1) 主シナリオに「モデル選択 → `saveAISettings` 成功 → summary 反映 → pill「固定済み」→ 開始ボタン有効化」の状態遷移、(2) 例外フローに backend `phase ai settings repository is not configured` 境界。
- UC 正本への恒久追記要否は `finalization-module` の `updating-docs` 起動時に判断する。

### 承認済み E2E テスト観点差分

- 分類: `追加候補あり` 3 件、`判断不足` なし。
- 追加候補: `E2E-UC-056`（単語翻訳）、`E2E-UC-057`（NPC ペルソナ生成）、`E2E-UC-058`（本文翻訳）。3 phase 同型で「未設定状態からモデル選択 → pill「固定済み」→ 開始ボタン有効化」を証明する。
- selector 不足なし（3 phase の `*-ai-model-lock-state`、`*-start-button` は画面設計書「E2E 固定 selector」表で定義済み）。
- 範囲外: bootstrap repository 注入の DI テスト、presenter `buildModelOptions` placeholder 重複解消の単体テスト → `implementation-module` の `tests-unit` で扱う。

### 画面再現確認の再現手順と修正後の期待状態

- 再現手順と現状の操作結果: `fix-decision.md` 「画面再現確認」節を引き継ぐ。
- 修正後に満たすべき期待状態:
  - 単語翻訳 / NPC ペルソナ生成 / 本文翻訳の 3 phase で、ジョブ#3 を開きフェーズ画面でモデル選択後、AI 設定 pill が「固定済み」を表示する。
  - 開始ボタンが enabled になる。
  - `SaveXxxPhaseAISettings` がエラーを返さない。
  - 直後の summary 再取得で `summary.aiSettings` が保存値を返す。

### 想定 Y/N 評価の更新

- 「backend 変更がある」: Y（bootstrap への repository 注入が必要）
- 「frontend と backend を接続する」: N（既存 Wails binding 経路を使うのみ）

## implementation-module decision table 結果

| 想定 | Y/N | 必要 artifact | 担当 |
| --- | --- | --- | --- |
| frontend ロジック変更がある | Y | frontend ロジック実装（presenter placeholder 重複解消） | `frontend_implementer` |
| backend 変更がある | Y | backend 実装（bootstrap DI 注入） + 単体テスト | `backend_implementer`, `implementation_tester`（単体） |
| frontend と backend を接続する | N | 統合境界実装 不要 / シナリオテストは E2E 観点差分追加候補のため別途要 | - |
| 実装済み責務を独立に証明したい | Y | 単体テスト | `implementation_tester`（単体） |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | 観測ログ追加（成立条件は最終検証前に判断） | 対象層 implementer |

加えて、investigation-module 出口の E2E テスト観点差分 `追加候補あり` 3 件（`E2E-UC-056/057/058`）を fail-test 経路で証明するため `シナリオテスト` を追加要とする。

## 後続モジュールへの引き継ぎ

- 入口: `finalization-module`
- 引き継ぐ事実: 本節「修正実行入力」と decision table 結果一式、`fix-decision.md`、`uc-diff-candidates.md`、`e2e-test-aspect-diff.md`、`human-review-overview.md`、下記「実装成果物」「最終検証」節。

## 実装成果物

### backend 実装

- 変更ファイル: `internal/bootstrap/app_controller.go` — 3 phase service 組み立て時に `repository.NewSQLiteJobPhaseAISettingsRepository(foundationDataDB)` で生成した repository を `WithTermTranslationJobPhaseAISettings`, `WithPersonaGenerationPhaseAISettingsRepository`, `WithBodyTranslationJobPhaseAISettingsRepository` で注入。

### frontend ロジック実装

- `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts` — `buildModelOptions` から `{value:"", label:"選んでください"}` placeholder 重複削除
- `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts` — `buildPersonaModelOptions` 同型
- `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts` — `buildBodyModelOptions` 同型

### 単体テスト

- `internal/bootstrap/app_controller_test.go` — 3 phase の `SaveXxxPhaseAISettings` が DI 注入後にエラーを返さないことを証明する 3 テスト追加
- `frontend/src/application/presenter/{term,persona,body}-translation-phase/*.presenter.test.ts` — 3 phase の `buildXxxModelOptions` が placeholder を含まないことを証明するテスト 6 件追加

### シナリオテスト

- 追加: `tests/system/fix-phase-ai-settings-pill.spec.ts` — `E2E-UC-056-A/B`, `E2E-UC-057-A/B`, `E2E-UC-058-A/B`（3 phase × 未設定 pill 確認 + 保存後 pill 切替確認）
- ID 採番変更: 当初 `E2E-UC-053/054/055` を採番していたが既存 ID 衝突のため `E2E-UC-056/057/058` に変更
- 関連: `tests/system/fix-term-ai-model-settings.spec.ts`（既存 spec の操作順序修正で `E2E-UC-FIX-MODEL-003` 復活 pass）、`tests/system/support/scenario-wails-mocks.ts`（新規 spec 用 jobId 14/15/16 の独立 mock 追加）、`tests/system/support/translation-phase-pages.ts`（persona/body page object に provider/model select getter 追加）、`docs/e2e-test-design/test-design.csv`

### 観測ログ追加

- 省略。理由: fix_decider が確定原因（bootstrap DI 注入欠落）を観測ログ駆動で固定済みであり、修正後は repository nil の早期 return 経路を通らず、恒久観測点を追加すべき実行時分岐が残らない。一時観測ログは fix-decision 完了時点で全て削除済み。

### 付帯修正

- `scripts/test/run-system-test.sh` — `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND` を `in-memory`（bootstrap 未サポート）から `fake` に修正。最終検証 `--suite all` の system test 起動を可能にする infrastructure 修正。master でも同状態だったが、本 task で `--suite all` を通すために含めた。

## 最終検証

### harness 全 suite 結果

- backend-local: 全パッケージ pass（14 パッケージ）
- frontend-local: 53 ファイル / 579 テスト pass
- coverage: 全体 72.7%（基準 70.0% 超）
- system: **76 pass / 2 fail / 0 did not run**
  - 本 task で追加した 6 シナリオ + 既存 `E2E-UC-FIX-MODEL-001/002/003` は全て pass
  - 当初 pre-existing と判定した 7 件中 5 件は本 task で修正完了:
    - `E2E-UC-051/052/053`（`translation-phases.spec.ts`）: 各 spec で `missingTermExecution / missingPersonaExecution / missingBodyExecution` mock オプションを指定して「未設定」状態に切り替え → pass
    - `SCN-TJM-005/007`（`translation-job-management.spec.ts`）: `resumeButton` locator を厳密化 → pass
    - `E2E-UC-045`（`job-run-shell.spec.ts:128`）: `noResultQuery` 検索後の `processingTargetRows.first()` 設計修正 → pass
  - 残 2 fail（別 task に切り出し済み）:
    - `E2E-UC-048`（`job-run-shell.spec.ts:254`）と `E2E-UC-049`（`:273`）: `JobRunPage.svelte` の `$effect` が `setCurrentPhasePage` 経由で `currentPhasePage` を Svelte 5 reactivity 依存に登録してしまい、`setCurrentPhasePage("persona"|"complete")` 後に初期 phase に戻す循環。本 task 主旨（pill 修正）と独立した別バグのため follow-up task として切り出し: `docs/exec-plans/active/fix-job-run-shell-effect-untrack/plan.md`

### 実画面確認（CLAUDE.md「実画面確認」規定）

- 起動: `npm run dev:wails:run` / `http://localhost:34115`
- 対象: ジョブ#3 / 単語翻訳フェーズ（3 phase は同型修正のため代表 1 phase で確認）
- 確認手順と結果:
  1. ジョブ#3 を開き単語翻訳フェーズ画面に遷移 → pill「設定未完了」、開始ボタン disabled（期待どおり）
  2. AIサービス select で Gemini を選択 → モデル一覧更新ボタン押下 → モデル select に `fake-model` が追加（placeholder 重複なし、修正効果確認）
  3. モデル select で `fake-model` を選択 → **pill「固定済み」、開始ボタン enabled**（修正前は両者とも変わらなかった）
- 直接 backend 呼び出し確認: chrome-devtools console から `window.go.wails.AppController.SaveTermTranslationPhaseAISettings({provider:"gemini", model:"fake-model", executionMode:"batch", batchMode:"disabled"})` を直接実行 → `phase ai settings repository is not configured` エラーが返らず正常保存値が返ることを確認（fix-decision.md の確定原因が消えたことの直接証拠）
- 証跡 screenshot: `screenshot-04-after-fix-pill-locked.png`

## 正本化判断（finalization-module）

- 仕様変更または仕様追加: なし（investigation-module で N 確定済み）
- UC 正本反映: **なし**（人間判断: UC 差分候補 2 件は E2E テスト観点として扱うべきで、UC（ユースケース）レベルの記述ではないため、`docs/usecases/uc-translation-management.md` への追記は行わない）
- E2E テスト観点正本反映: 既に implementation-module 内で `docs/e2e-test-design/test-design.csv` に `E2E-UC-056/057/058` を追加済み
- 人間承認済みの恒久仕様: なし
- `詳細仕様正本反映`（`updating-docs` skill）: skip
- 承認日: 2026-06-02

## finalization-module 完了記録

- 作業 commit: `652d72d7` （branch `claude/fix-phase-ai-settings-pill-update-after-model-select` 上、28 ファイル変更 / 1543 insertions / 69 deletions）
- local merge commit: `0d79db66` （`master` 上、`--no-ff`、conflict なし）
- merge 後検証: `python3 scripts/harness/run.py --suite all` → backend / frontend / coverage 全 pass、system 76 pass / 2 fail（既知の `E2E-UC-048/049`、follow-up task `fix-job-run-shell-effect-untrack` に切り出し済み）
- completed 移動: active → completed フォルダへ git mv 完了
- remote push: 行わない（finalization-module 規約）
