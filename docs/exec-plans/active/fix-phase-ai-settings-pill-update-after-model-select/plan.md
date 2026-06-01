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

## 想定 Y/N 評価（暫定）

| 想定 | Y/N | 根拠 |
| --- | --- | --- |
| 仕様変更または仕様追加がある | N | 画面仕様（`docs/screen-design/screens/term-translation-phase.md`）でモデル選択後の固定遷移は定義済み。 |
| 画面変更がある | N | 画面構造変更なし。 |
| 内部構造変更がある | Y | panel の handler が saveAISettings に正しい payload を渡せていないか、controller の saveAISettings が viewModel に反映されない問題。 |
| 画面の表示変更がある | N | pill / select / ボタン文言は変更不要。 |
| frontend ロジック変更がある | Y | handler 経路、controller の state 反映、presenter の派生ロジックのどこかに不具合。 |
| backend 変更がある | 判断保留 | saveAISettings の backend binding 受信状態を investigation-module で確定する。 |
| frontend と backend を接続する | 判断保留 | 同上。 |
| 実装済み責務を独立に証明したい | Y | controller / presenter の saveAISettings 後 state を単体テストで証明。 |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | provider/model/executionMode の選択順序と Svelte 5 reactivity の組み合わせを実行時に観測する必要がある。 |

「仕様変更または仕様追加がある」が N のため、investigation-module を継続できる。

## Wails 接続対象

- 起動 command: `npm run dev:wails:run`
- 接続先: `http://localhost:34115`

## 後続モジュールへの引き継ぎ

- 入口: investigation-module
- 引き継ぐ事実: 本 plan.md の人間観測記録と暫定想定 Y/N 評価。前 task の plan / fix-decision は参考として参照可能。
