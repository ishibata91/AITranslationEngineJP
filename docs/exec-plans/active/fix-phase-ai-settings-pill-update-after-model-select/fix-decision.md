# fix-decision: fix-phase-ai-settings-pill-update-after-model-select

## 判断結果

- 判定: 完了
- 停止理由: なし

## 観測済み問題

- 問題: 3 phase（単語翻訳 / ペルソナ生成 / 本文翻訳）の AI 設定パネルで、AI サービスとモデルを選択して `saveAISettings` を呼んでも、`SaveTermTranslationPhaseAISettings` が `phase ai settings repository is not configured` エラーを返す。その結果 `fetchSummaryAndReadiness` が実行されず、`summary.aiSettings` が更新されない。`isExecutionConfigured` は `summary.aiSettings` が未設定のまま false になり、AI 設定 pill が「設定未完了」のまま固定される。
- 期待との差分: モデル選択後に `saveAISettings` が成功し、`fetchSummaryAndReadiness` が summary を更新し、`isExecutionConfigured` が true になり、pill が「固定済み」へ切り替わり、開始ボタンが有効化されるべき。

## 画面再現確認

- Wails 接続対象: `http://localhost:34115`
- 再現手順:
  1. 翻訳管理画面でジョブ#3 を選択して単語翻訳フェーズ画面に遷移する
  2. AIサービス select で Gemini を選択する（DOM value: `"gemini"`）
  3. 「モデル一覧を更新」ボタンをクリックする
  4. モデル select で `fake-model` を選択する
  5. pill と開始ボタンの状態を確認する
- 操作結果:
  - ステップ3: `ListTranslationJobSetupProviderModels` が `{provider: "gemini", status: "success", models: [{modelId: "fake-model"}]}` を返す。モデル select に `fake-model` が追加される（成功）
  - ステップ4: `change` イベントで `handleModelChange` が呼ばれ、`saveAISettings({provider: "gemini", model: "fake-model", executionMode: "batch", batchMode: "disabled"})` が gateway に渡る
  - backend の `SaveTermTranslationPhaseAISettings` は `save term translation phase ai settings: save term translation phase ai settings: save term translation ai settings: phase ai settings repository is not configured` エラーを返す
  - usecase の `saveAISettings` はエラーをキャッチして `fetchSummaryAndReadiness` を呼ばない
  - pill は「設定未完了」のまま固定される
- 画面状態: pill「設定未完了」、開始ボタン disabled、エラーメッセージ「単語翻訳段階の AI 設定保存に失敗しました」（errorMessage に記録されるが表示上変化しない）
- 証跡 path:
  - `docs/exec-plans/active/fix-phase-ai-settings-pill-update-after-model-select/screenshot-02-initial-state.png`
  - `docs/exec-plans/active/fix-phase-ai-settings-pill-update-after-model-select/screenshot-03-pill-stuck-after-model-select.png`

## 原因仮説と検証結果

### 仮説 A: backend の `phaseAISettingsRepository` 未注入

- 根拠: `internal/bootstrap/app_controller.go` で `NewTermTranslationPhaseService(...)` を呼ぶ際に `WithTermTranslationJobPhaseAISettings` が呼ばれていない
- 検証: JS で `window.go.wails.AppController.SaveTermTranslationPhaseAISettings({provider: "gemini", model: "fake-model", executionMode: "batch", batchMode: "disabled"})` を直接呼んだところ、`save term translation ai settings: phase ai settings repository is not configured` エラーが返った
- 結果: 支持（確定）

### 仮説 B: frontend の `resolvedProviderValue` が空文字になり `refreshModelList` が早期リターンする

- 根拠: `fill` ツールが Svelte の `onchange` を発火させない場合、`selectedProviderId = ""` のまま `resolvedProviderValue = ""` になり、`if (!provider) { return }` で早期リターンする
- 検証: JS で provider select の `change` イベントを手動発火後にモデル一覧更新ボタンをクリックし、`ListTranslationJobSetupProviderModels` の呼び出しをインターセプト → `{provider: "gemini", requestToken: "1"}` で正しく呼ばれた
- 結果: 否定（`refreshModelList` は正常に動作し、`fake-model` が返る）

### 仮説 C: `response.provider === provider` チェックで弾かれる

- 根拠: backend から返る `provider` 値が frontend の `provider` と一致しない可能性
- 検証: `ListTranslationJobSetupProviderModels` の response は `{provider: "gemini"}` を返し、request の `provider: "gemini"` と一致する
- 結果: 否定

### 仮説 D: `AIModelSelectionCard` が `<option value="">選んでください</option>` を固定描画し、`buildModelOptions` も同じ placeholder を返すため重複する

- 根拠: `AIModelSelectionCard.svelte:195` に `<option value="">{emptyModelLabel}</option>` が固定存在し、`buildModelOptions` も `currentModel` 未設定時に `{value: "", label: "選んでください"}` を先頭に追加する
- 検証: JS で model select の options を確認 → `[{value:""}, {value:""}, {value:"fake-model"}]` と `value=""` が2つある
- 結果: 支持（secondary 問題として確定）。ただし主因は仮説 A であり、仮説 D は pill 固定とは独立した UI 問題。

## 確定原因

### 主因

`internal/bootstrap/app_controller.go` の `termTranslationPhaseController` 組み立て時に `WithTermTranslationJobPhaseAISettings(jobPhaseAISettingsRepository)` が呼ばれていない。そのため `TermTranslationPhaseService.phaseAISettingsRepository` は `nil` のままで、`SaveAISettings` を呼ぶと `phase ai settings repository is not configured` エラーを返す。フロントエンドの usecase はエラーをキャッチして `fetchSummaryAndReadiness` を呼ばないため、`summary.aiSettings` が更新されない。`isExecutionConfigured` は false のまま固定され、pill が「設定未完了」から変わらない。開始ボタンも disabled のままになる。

- 観測根拠: JS から `window.go.wails.AppController.SaveTermTranslationPhaseAISettings` を直接呼び、`phase ai settings repository is not configured` エラーが返ることを確認した

同じ問題は `PersonaGenerationPhaseService`（`WithPersonaGenerationPhaseAISettingsRepository` 未呼び出し）と `BodyTranslationPhaseService`（`WithBodyTranslationJobPhaseAISettingsRepository` 未呼び出し）にも存在する可能性が高い。3 phase 共通の未注入問題と判断する。

### secondary 問題（主因とは独立）

`AIModelSelectionCard.svelte` が model select に常に `<option value="">選んでください</option>` を固定描画し、`buildModelOptions` も `currentModel` 未設定時に `{value: "", label: "選んでください"}` を先頭に追加するため、`value=""` の option が重複して描画される。この問題は pill 固定の直接原因ではないが、UI の一貫性を損なう。

## 採用する修正方針

### 主因修正（backend）

`internal/bootstrap/app_controller.go` で、3 phase 全ての service 組み立て時に `JobPhaseAISettingsRepository` を注入する。

- `TermTranslationPhaseService` に `.WithTermTranslationJobPhaseAISettings(jobPhaseAISettingsRepository)` を追加する
- `PersonaGenerationPhaseService` に `.WithPersonaGenerationPhaseAISettingsRepository(jobPhaseAISettingsRepository)` を追加する
- `BodyTranslationPhaseService` に `.WithBodyTranslationJobPhaseAISettingsRepository(jobPhaseAISettingsRepository)` を追加する

`jobPhaseAISettingsRepository` は `repository.NewSQLiteJobPhaseAISettingsRepository(db)` で生成する（既存の SQLite 実装を使う）。または `dev:wails:run` 環境では `repository.NewInMemoryJobPhaseAISettingsRepository()` を使う判断を実装 agent に委ねる。

この修正により、`SaveTermTranslationPhaseAISettings` が正常に `JOB_PHASE_AI_SETTINGS` テーブルへ upsert し、直後の `getTermTranslationPhaseSummary` で `summary.aiSettings` に保存済み設定が返るようになる。`isExecutionConfigured` が true になり、pill が「固定済み」へ切り替わる。

### secondary 問題修正（frontend）

`buildModelOptions`（`term-translation-phase.presenter.ts` および同型の他 phase presenter）が返す options から `{value: "", label: "選んでください"}` の placeholder を削除し、`AIModelSelectionCard.svelte` が固定描画する `<option value="">選んでください</option>` だけに統一する。

ただし secondary 問題は pill 固定の直接原因ではないため、主因修正後の検証で改めて判断する。実装 agent は主因修正を優先する。

### 3 phase 共通適用の確認

この修正方針は 3 phase 全てに共通して適用できる。bootstrap の1ファイル（`app_controller.go`）に 3 phase 分の repository 注入を追加するだけで解決する。フロントエンド側の修正は各 phase の presenter に対して同型で適用する。

## 禁止する修正

- frontend で `isExecutionConfigured` を常に true にする分岐を追加すること（症状を隠す対症療法）
- `saveAISettings` のエラーを無視して `fetchSummaryAndReadiness` を強制呼び出しすること（エラー状態のまま summary を取得しても `aiSettings` は更新されないため意味がない）
- frontend の `$derived` に新しい状態値（例: `localAISettingsSaved` フラグ）を追加して pill 状態を管理すること（既存の `isExecutionConfigured` 判定ロジックで説明できる）
- `SaveTermTranslationPhaseAISettings` の mock や fake 実装を挿入して backend エラーを隠すこと

## 影響ファイル候補

| ファイル | 変更内容 | 理由 |
| --- | --- | --- |
| `internal/bootstrap/app_controller.go` | `WithTermTranslationJobPhaseAISettings`, `WithPersonaGenerationPhaseAISettingsRepository`, `WithBodyTranslationJobPhaseAISettingsRepository` を追加する | 観測で確定した主因。repository 注入が欠落している |
| `internal/repository/job_phase_ai_settings_inmemory_repository.go` または `job_phase_ai_settings_sqlite_repository.go` | 参照のみ（既存実装を注入する） | `jobPhaseAISettingsRepository` の実装は既存 |
| `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts` | `buildModelOptions` の placeholder 重複を解消する（secondary） | `AIModelSelectionCard` の固定 placeholder と重複している |
| `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts` | 同型の `buildModelOptions` 修正（secondary） | 3 phase 共通の問題 |
| `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts` | 同型の `buildModelOptions` 修正（secondary） | 3 phase 共通の問題 |

実装 agent は `internal/bootstrap/app_controller.go` の修正を最優先とし、secondary 問題は主因修正後の検証結果に応じて判断する。
