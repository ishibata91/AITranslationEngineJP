# 修正前調査

## 判定

- 判定: 完了
- 調査 mode: `修正前調査`、`UI 根拠`
- 対象 task: `2026-05-07-master-persona-model-refresh-provider-reset`
- 推奨引き継ぎ先: `fix_lane`

## 人間観測

- 人間観測では、モデル一覧更新後に AI サービスの選択状態が維持されず、モデルを選択できないと記録されている。根拠: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/human-observation.md)
- 人間観測では、`fake-model` を選べた直前 task があり、実画面確認だけが残留リスクだった。根拠: [investigation-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/investigation-input.md), [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)

## 観測事実

- 観測: AI 設定カードの更新ボタンは、モデル一覧取得専用処理ではなく `refreshAISettings()` を呼ぶ。根拠: [GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte:67
frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte:78

- 観測: `refreshAISettings()` は `useCase.loadAISettings()` をそのまま呼ぶ。根拠: [master-persona-screen-controller.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts)
frontend/src/controller/master-persona/master-persona-screen-controller.ts:210

- 観測: `loadAISettings()` 成功時は、gateway から返った `provider/model/executionMethod` で `draft.aiSettings` を丸ごと置き換え、`modelOptions` は `mergedSettings.model` 1 件だけから再生成する。`model` が空文字なら `modelOptions` は空配列になる。根拠: [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:115
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:121

- 観測: provider を UI で変更した時点で、frontend は `aiSettings.model=""` と `modelOptions=[]` を即時適用する。根拠: [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:228
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:233

- 観測: presenter は `modelOptions.length===0` の間、警告文言を「モデル一覧を更新後に選べる状態で接続します。」にし、`canSelectModel=false` にする。画面側は `canSelectModel=false` の間、モデル select を disabled にする。根拠: [master-persona.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts), [GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
frontend/src/application/presenter/master-persona/master-persona.presenter.ts:65
frontend/src/application/presenter/master-persona/master-persona.presenter.ts:71
frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte:68
frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte:73

- 観測: backend の `LoadSettings()` は、通常 path では保存済み `record.Provider` と `record.Model` を返す。`fakeModelDefaults` が有効な path だけ、保存済み provider を補正しつつ `Model: "fake-model"` を返す。根拠: [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
internal/service/master_persona_service.go:312
internal/service/master_persona_service.go:317
internal/service/master_persona_service.go:335

- 観測: 現行ローカル UI では `http://localhost:34115/#master-persona` に到達できた。provider を `xAI` に変更した前後で、provider の選択値は `xai` のままだった。一方で model は更新前後とも disabled で、表示値は `設定が必要` のままだった。根拠: [before-refresh-xai.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/before-refresh-xai.png), [after-refresh-xai.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/after-refresh-xai.png)

- 観測: 現行ローカル UI の console には、`InvalidStateError: Failed to execute 'send' on 'WebSocket': Still in CONNECTING state.` が複数回出ている。参照スタックは `master-persona.usecase.ts:277` で、detail 取得失敗の catch 側に寄っている。さらに `wails dev Disconnected from backend` と `Connected to backend` が反復している。根拠: [console.log](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log)
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:1
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:489
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:522
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:523

## UI 証跡

- 取得済み画面: [before-refresh.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/before-refresh.png)
- 取得済み画面: [after-refresh.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/after-refresh.png)
- 取得済み画面: [before-refresh-xai.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/before-refresh-xai.png)
- 取得済み画面: [after-refresh-xai.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-master-persona-model-refresh-provider-reset/after-refresh-xai.png)
- 取得済み text: `agent-browser get text "#masterPersonaView"` では、更新後も「モデル一覧を更新してください。」と「設定が必要」が残った。根拠: 現 turn 観測

## ログ証跡

- 保存済み console: [console.log](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log)
- 保存済み errors: [errors.log](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/errors.log)

## 原因候補

- 候補1: 更新ボタンの意味が「モデル一覧取得」ではなく「保存済み AI 設定の再読込」になっている可能性が高い。`loadAISettings()` は provider と model をまとめて再投入するため、保存済み値が現在の UI 選択と異なる環境では、provider が巻き戻る余地がある。根拠: [master-persona-screen-controller.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts), [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
frontend/src/controller/master-persona/master-persona-screen-controller.ts:210
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:116

- 候補2: provider 変更時に frontend が model と modelOptions を必ず空にする一方で、更新後に復元できる model source は `loadAISettings()` の戻り値しかない。保存済み `record.Model` が空なら、更新後も model は常に disabled のまま残る。根拠: [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts), [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:121
frontend/src/application/usecase/master-persona/master-persona.usecase.ts:232
internal/service/master_persona_service.go:320

- 候補3: `fake-model` を返す backend path は `fakeModelDefaults` 有効時だけである。実行中 Wails app がその path に入っていない場合、直前 task の期待と異なり model options は空のままになる可能性がある。根拠: [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go), [implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/implementation-result.md)
internal/service/master_persona_service.go:317
internal/service/master_persona_service.go:337

- 候補4: 現行ローカル app では backend 接続が不安定で、Wails WebSocket の `CONNECTING` エラーと再接続反復がある。人間観測の「リセット」は、設定再読込そのものに加えて、接続不安定時の再初期化が重なった見え方である可能性がある。根拠: [console.log](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log)
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:1
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:522
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/console.log:523

## 影響ファイル候補

- [frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
- [frontend/src/controller/master-persona/master-persona-screen-controller.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts)
- [frontend/src/application/usecase/master-persona/master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
- [frontend/src/application/presenter/master-persona/master-persona.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts)
- [internal/service/master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
- [frontend/src/controller/wails/master-persona.gateway.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/wails/master-persona.gateway.ts)

## 残り不足

- 不足: 現行実行環境の保存済み `master persona ai settings` 実値が未確認である。provider 巻き戻りの発火条件をまだ固定できない。
- 不足: `fakeModelDefaults` が現行 Wails app で有効か無効かを、この調査では確認していない。
- 不足: `MasterPersonaLoadAISettings` 呼び出し時の実レスポンス payload を直接記録していない。
- 不足: 人間観測の「provider reset」は、現行ローカル UI では再現しなかった。現行ローカル UI で再現できたのは「更新後も model を選べない」点までである。

## 残留リスク

- リスク: 修正対象を frontend だけに寄せると、`fakeModelDefaults` や保存済み設定の drift を見逃す可能性がある。
- リスク: 現行ローカル app の接続不安定が観測を汚しているため、再現条件を固定しないまま修正へ進むと、再発条件を取りこぼす可能性がある。

## 推奨 next step

- 推奨: `fix_lane` が `修正実行入力` を作る前に、`MasterPersonaLoadAISettings` の実レスポンスと現行保存済み AI 設定の値を 1 回だけ追加確認する。
- 推奨: 追加確認では、保存済み provider/model と `fakeModelDefaults` の有効状態を分けて記録する。
- 推奨: 追加確認後は、更新ボタンが再読込なのかモデル取得なのかを起点に narrow した修正実行入力へ進める。

## 追加観測事実

- 追加観測: 実行中 Wails app の `AppController.MasterPersonaLoadAISettings()` を画面上から直接呼ぶと、実レスポンス payload は `{"provider":"","model":"","executionMethod":"single_request"}` だった。根拠: [load-ai-settings-payload.json](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/load-ai-settings-payload.json)
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/load-ai-settings-payload.json:1

- 追加観測: 既定 DB path は repo 直下の `db/master-dictionary.sqlite3` である。`PERSONA_GENERATION_SETTINGS` の読込 SQL は `id = 1` の singleton row を読む。根拠: [app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go), [master_persona_sqlite_repository.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/master_persona_sqlite_repository.go)
internal/bootstrap/app_controller.go:291
internal/repository/master_persona_sqlite_repository.go:109

- 追加観測: `db/master-dictionary.sqlite3` の `PERSONA_GENERATION_SETTINGS` には `id = 1` の row が存在しなかった。確認結果は `row_count=0, provider=\"\", model=\"\", execution_method=\"\"` である。保存済み provider/model を DB から観測できる範囲では、未保存と読める。根拠: [persona-generation-settings.csv](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/persona-generation-settings.csv)
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/persona-generation-settings.csv:1
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/persona-generation-settings.csv:2

- 追加観測: `fakeModelDefaults` は、`masterPersonaAIMode() == "fake"` の時だけ wiring される。option が有効なら `service.fakeModelDefaults = true` になり、`LoadSettings()` は保存値に関係なく `model: "fake-model"` を返す。根拠: [app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go), [master_persona_provider_transport.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_provider_transport.go), [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
internal/bootstrap/app_controller.go:135
internal/bootstrap/app_controller.go:403
internal/bootstrap/app_controller.go:460
internal/service/master_persona_provider_transport.go:60
internal/service/master_persona_service.go:317
internal/service/master_persona_service.go:337

- 追加観測: 現行 Wails app の `MasterPersonaLoadAISettings` 実レスポンスは `model=""` であり、`provider=""` でもある。上記コード契約と照合すると、現行 app では `fakeModelDefaults` が有効である根拠は観測できず、無効である観測と整合する。もし有効なら、空保存時でも `provider="gemini"` と `model="fake-model"` が返るはずである。根拠: [load-ai-settings-payload.json](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/load-ai-settings-payload.json), [master_persona_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/master_persona_service.go)
tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/load-ai-settings-payload.json:1
internal/service/master_persona_service.go:317
internal/service/master_persona_service.go:333
internal/service/master_persona_service.go:337

- 追加観測: 実行中プロセス環境を `ps` で直接確認する試行は sandbox 制限で失敗した。環境変数そのものの live 値は未確認である。根拠: 現 turn 観測

## 追加ログ証跡

- 追加保存: [load-ai-settings-payload.json](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/load-ai-settings-payload.json)
- 追加保存: [persona-generation-settings.csv](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/logs/2026-05-07-master-persona-model-refresh-provider-reset/persona-generation-settings.csv)

## 修正実行入力へ進めるか

- 判定: 進めてよい。
- 理由: `MasterPersonaLoadAISettings` の実レスポンス payload は取得済みである。保存済み AI 設定は DB row 不在として記録できた。`fakeModelDefaults` は live env 値そのものは未確認だが、現行 app の実レスポンスとコード契約の照合により、少なくとも現行 app 挙動は「無効と整合する」と観測できた。
- 注意: 人間観測の `provider reset` 自体は現行ローカル UI で再現していない。修正実行入力では「更新ボタンが設定再読込を呼ぶこと」と「保存済み設定 row 不在で model が復元されないこと」を観測事実として分離して扱う必要がある。
