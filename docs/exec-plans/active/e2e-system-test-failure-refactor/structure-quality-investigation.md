# 構造品質調査

## 判断結果

- 判定: 完了
- 調査 mode: `構造品質調査`
- 対象: `EF-001` から `EF-005`
- 引き継ぎ先: `designer`

## 根拠参照

- `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
- `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
- [frontend/src/ui/App.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte:63)
- [frontend/src/bootstrap/app-screen-controller-factories.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/bootstrap/app-screen-controller-factories.ts:50)
- [frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:175)
- [internal/controller/wails/provider_settings_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/controller/wails/provider_settings_controller.go:116)
- [internal/service/provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go:270)
- [frontend/src/application/presenter/master-persona/master-persona.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts:112)
- [internal/service/master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go:353)
- [internal/repository/master_persona_sqlite_repository.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/master_persona_sqlite_repository.go:572)
- [frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte:61)
- [internal/service/translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:427)
- [scripts/test/seed-system-test-db/main.go](/Users/iorishibata/Repositories/AITranslationEngineJP/scripts/test/seed-system-test-db/main.go:59)
- [tests/system/support/scenario-wails-mocks.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/tests/system/support/scenario-wails-mocks.ts:434)
- [tests/system/master-persona.spec.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/tests/system/master-persona.spec.ts:19)
- [tests/system/translation-job-management.spec.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/tests/system/translation-job-management.spec.ts:184)
- [tests/system/job-run-shell.spec.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/tests/system/job-run-shell.spec.ts:13)
- [tests/system/output-management.spec.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/tests/system/output-management.spec.ts:10)

## 観測事実

- `EF-001`: root `App.svelte` は provider settings screen controller の既定生成で `null gateway` を注入する。一方で production factory は Wails gateway を生成する。`ProviderSettingsUseCase.saveSettings` は gateway 応答が返れば `saveNotice` を成功文言へ更新し、保存時に endpoint 妥当性検証を呼ばない。
- `EF-002`: master persona 画面は `state.modelOptions.length === 0` なら model select を disabled にする。backend の `LoadAISettings` は `PERSONA_GENERATION_SETTINGS` に行がなければ空設定を返す。system-test seed には master persona AI 設定と provider settings の投入がない。
- `EF-003`: translation job management の `ResumeJob` は service / usecase コメント通り状態変更を実行せず、要約だけを返す。画面側は `requestResume()` 後、job state が `Ready` でなければ job run shell を開く。
- `EF-004`: root `App.svelte` は output management screen controller の既定生成で `null gateway` を注入する。一方で production factory は `createTranslationOutputArtifactGateway()` を接続している。
- `EF-005`: system-test seed は `system-test-ready`, `system-test-paused`, `system-test-failed` など 6 job だけを投入する。job-run / translation phase 系 test と scenario mock は `system-test-term`, `system-test-persona`, `system-test-body-failed`, `system-test-body-completed`, `system-test-body-running` を前提にする。translation job management service は completed job を管理画面対象外にする。

## EF 別整理

### EF-001 AIサービス設定

#### 原因候補

- 構造設計不整合: root `App.svelte` の既定 wiring が production factory と異なり、provider settings だけ gateway 未接続の経路を持つ。`frontend/src/ui/App.svelte:63-67` と `frontend/src/bootstrap/app-screen-controller-factories.ts:83-84` が衝突している。
- 責務分離不足: 保存処理と接続確認処理が別 usecase に分かれているが、UI 成功表示が「保存成功」と「利用可能」を混同しやすい。`saveSettings` は persistence 成功だけで `saveNotice` を出し、`ValidateProviderSettings` は別操作でしか動かない。
- コーディング規約逸脱: `coding-guidelines.md` の「境界入力は使用直前に検証し、失敗を握りつぶさない」に対し、save 経路では endpoint 妥当性を使用直前に検証していない。

#### 責務境界

- frontend usecase の責務: 入力 draft と保存通知の管理。
- backend provider settings service の責務: row 保存、secret 保存、validation state 管理。
- validation adapter の責務: 接続確認だけを行う。

#### 変更候補

- `frontend/src/ui/App.svelte` の既定 wiring と production wiring の二重定義解消。
- provider settings の成功表示条件整理。保存成功と接続確認成功の文言境界を見直す候補がある。
- save と validate の契約分離を維持するか、save 前の入力妥当性検証を追加するかの構造判断材料が必要である。

#### 変更不要範囲

- `internal/infra/ai/provider_settings_validation.go` の transport probe 自体は独立責務であり、保存成功表示の直接原因ではない。
- `wailsjs` generated binding は保存経路の責務過多を作っていない。

#### 追加調査が必要な箇所

- system-test 実行時に provider settings route が `App.svelte` の既定経路を通るのか、`main.ts` 側の factory 注入を通るのか。
- save 時に期待する仕様が「入力形式 rejection」なのか「保存は通すが未確認表示」なのか。

### EF-002 マスターペルソナ

#### 原因候補

- 構造設計不整合: presenter は `modelOptions.length === 0` を model select 活性条件に使うが、初期 modelOptions の供給は backend seed に依存する。seed 欠落で UI が永久に disabled になりやすい。
- 責務分離不足: master persona の画面初期表示が `LoadAISettingsState` と provider settings 保存状態の両方に依存するが、system-test seed はその前提を一つの場所で固定していない。
- テスト前提との構造不整合: test は `gemini-test` が即選択可能であることを前提にするが、実 DB では `PERSONA_GENERATION_SETTINGS` 未投入なら空設定になる。

#### 責務境界

- backend repository の責務: `PERSONA_GENERATION_SETTINGS` の保存済み行を返す。
- backend service の責務: provider settings から providerOptions と model list 初期状態を組み立てる。
- presenter の責務: 画面の活性条件を決める。

#### 変更候補

- system-test seed に master persona AI settings と対応 provider settings を追加する候補がある。
- presenter の活性条件を `modelOptions` 依存から `modelSettingsCard.modelList.status` 依存へ寄せる候補がある。
- `LoadAISettingsState` の初期 model list 契約を、空設定でも refresh 導線が明確になる形へ寄せる候補がある。

#### 変更不要範囲

- `frontend/src/controller/wails/master-persona.gateway.ts` の binding 解決は本件の主因ではない。
- `GenerationSetupPanel.svelte` の button disabled 自体は presenter から渡された状態を忠実に描画しているだけである。

#### 追加調査が必要な箇所

- system-test 実行 DB に `PERSONA_GENERATION_SETTINGS` 行が存在するか。
- provider settings summary に `gemini` の endpoint / credential が seeded されているか。

### EF-003 翻訳実行シェル

#### 原因候補

- 構造設計不整合: translation job management の `ResumeJob` は service / usecase / log 文言が一貫して「状態変更を実行しない read model」と定義しているが、UI test は resume 後の job run shell 表示を期待する。
- 責務過多: `TranslationJobManagementPage.svelte` が action result の meaning 判定と画面遷移を同時に持ち、backend が再開失敗要約しか返さなくても shell を開く分岐を含む。
- 責務分離不足: 「再開要求の受付」と「現在段階画面へ遷移できること」の判定が job state だけで結び付けられており、action result の成功契約を見ていない。

#### 責務境界

- backend service / usecase の責務: resume 可否要約の返却。現状は実行本体を持たない。
- frontend page の責務: 操作後遷移。現状は backend outcome を検証せず route 変更する。

#### 変更候補

- `internal/service/translation_job_management_service.go` の resume 経路を read-only のまま保つか、実行経路へ拡張するかの人間判断が必要である。
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte` の route 遷移条件を action result の成否へ合わせる候補がある。

#### 変更不要範囲

- `JobRunPage.svelte` の phase 切替 UI は、`selectedJobTarget` が入れば動作する構造であり、resume 自体の不成立とは別問題である。
- scenario mock の body / persona / term phase summary は shell 描画の補助であり、resume 失敗の直接原因ではない。

#### 追加調査が必要な箇所

- 現行仕様で resume が product behavior として実行されるべきか、現状どおり read-only なのか。
- real backend system-test failure が「shell を開かない」のか「開くが phase summary が空」のか。

### EF-004 出力管理

#### 原因候補

- 構造設計不整合: root `App.svelte` の既定 wiring が output management だけ `null gateway` を注入している。production factory は gateway ありで生成しており、二つの composition root が食い違う。
- 責務分離不足: `TranslationOutputArtifactUseCase.refresh()` は gateway 未接続でも `phase=ready` と `hasLoaded=true` を立てるため、画面は「読込失敗」ではなく「空表示」に寄りやすい。
- テスト前提との構造不整合: scenario mock は `GetTranslationOutputReview` と diff preview を返せるが、`null gateway` 経路では一切使われない。

#### 責務境界

- composition root の責務: gateway 注入。
- output artifact usecase の責務: review と diff preview の読込。

#### 変更候補

- `frontend/src/ui/App.svelte` と production factory の DI 統合。
- gateway 未接続時の画面状態表現を `ready` ではなく接続異常へ寄せる候補がある。

#### 変更不要範囲

- `internal/controller/wails/translation_output_artifact_controller.go` の DTO 変換は候補行未表示の直接原因ではない。
- diff preview row builder や XML command 側の backend service は、候補一覧の初回表示より後段である。

#### 追加調査が必要な箇所

- system-test で output management が `App.svelte` の既定 wiring を通る再現条件。
- 本番起動で同じ route に到達した時も未接続になるか。

### EF-005 翻訳段階

#### 原因候補

- 構造設計不整合: system-test seed の job 群と、job-run / translation phase test が探す job 群が一致していない。seed は `system-test-ready` などを投入し、test と scenario mock は `system-test-term` などを期待する。
- 責務分離不足: translation job management, job run shell, translation phases がそれぞれ別の fixture source を持ち、1 つの canonical system-test seed を共有していない。
- 構造設計不整合: translation job management service は completed job を一覧対象外にするが、job-run-shell test は `system-test-body-completed` を管理画面から開く前提を持つ。

#### 責務境界

- seed script の責務: real backend system-test の canonical data を固定する。
- scenario mock の責務: frontend-only test の fake contract を返す。
- translation job management service の責務: incomplete job だけを返す。

#### 変更候補

- `scripts/test/seed-system-test-db/main.go` と `tests/system/support/scenario-wails-mocks.ts` の job family を整合させる候補がある。
- completed job の導線を job management から開くのか、output management から開くのかを仕様側で再確認する候補がある。
- translation phase test を real backend seed 準拠へ寄せるか、scenario mock 準拠へ寄せるかの整理が必要である。

#### 変更不要範囲

- `frontend/src/ui/screens/job-run/JobRunPage.svelte` の phase host 構造自体は、対象 job が選択されれば表示できる。
- `wailsjs` generated binding は対象 job 名の不一致を作っていない。

#### 追加調査が必要な箇所

- real backend の `ListIncompleteJobs` 結果に、test が要求する phase 別 job が存在するか。
- completed body job を job run shell へ導く current product rule。

## 構造品質観点別結果

### 責務過多

- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`: resume 操作と route 遷移判定を同じ関数で持つ。

### 責務分離不足

- provider settings: 保存成功と接続確認成功の利用者向け意味が分離されていない。
- master persona: provider settings 前提、page-local AI settings、model list 活性条件が seed 1 か所で固定されていない。
- system-test: real backend seed と scenario mock が別の job family を持ち、phase 系 test の正本が分散している。

### コーディング規約逸脱

- provider settings save 経路は境界入力の妥当性確認を save 直前で行わず、成功文言を返す。

### 構造設計不整合

- `frontend/src/ui/App.svelte` と `frontend/src/bootstrap/app-screen-controller-factories.ts` の composition root が二重化し、provider settings / output management の gateway 注入条件が一致していない。
- translation job management backend は read-only resume/stop を返すが、system test は状態変化を期待する。
- completed job を管理画面対象外にする backend 契約と、completed body job を管理画面から開く test 前提が衝突している。

### 未使用コード

- 今回の調査範囲では未使用コードの確証は取れていない。

## 変更不要範囲

- `internal/controller/wails/translation_output_artifact_controller.go` の DTO mapping。
- `frontend/src/controller/wails/master-persona.gateway.ts` の binding 解決。
- `frontend/src/ui/screens/job-run/JobRunPage.svelte` の phase panel 構造。

## 残り不足

- root `App.svelte` と `main.ts` のどちらが system-test 実行時の正規 entry かを、最新起動経路で未確認。
- provider settings save 仕様が「保存時 validation 必須」か「保存と validation 分離」かを、詳細仕様から未確認。
- system-test 実 DB の provider settings / master persona settings 実データを未観測。

## 残留リスク

- `EF-001` と `EF-004` は wiring 不整合と仕様判断不足が混在している可能性がある。
- `EF-002` と `EF-005` は product failure と test fixture mismatch が混在している可能性がある。
- `EF-003` は read-only 契約のまま実装修正に入ると、translation job management の責務境界を崩す可能性がある。

## 推奨 next step

- 推奨 next step: 追加調査
- 理由: system-test の実起動 entry と実 DB seed の観測が不足しており、implementation-scope 確定には追加根拠が必要である。
