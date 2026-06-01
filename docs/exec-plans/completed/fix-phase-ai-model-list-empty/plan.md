# fix-phase-ai-model-list-empty

## 依頼要約

3 phase（単語翻訳 / ペルソナ生成 / 本文翻訳）の AI 設定パネルで、AI サービス（provider）を選択した後に「モデル一覧を更新」ボタンを押してもモデル select の選択肢が更新されない。dev 環境（`npm run dev:wails:run`）の fake provider なら fake モデルが選択肢に並ぶ想定だが、実画面では「選んでください / モデル一覧を更新してください」のままで fake モデルが出ない。結果として provider → モデル → 処理方式 → 固定 → 開始有効化までの遷移が画面操作で踏めない。観測ログ駆動で確定原因と修正方針を固定する。

## 分岐元

- 分岐元 branch: `claude/fix-term-translation-model-settings-empty-fixed`（前 task の作業ブランチをそのまま継続使用）
- 分岐元 commit: 本 plan 作成時点の HEAD

理由: 前 task `fix-term-translation-model-settings-empty-fixed` の修正で生まれた残課題のため、ブランチと未 commit の変更をそのまま引き継いで作業を続ける。前 task の commit と本 task の commit は最終的に同じブランチ上に積まれる。

## 作業 branch

- `claude/fix-term-translation-model-settings-empty-fixed`（継続）

## 人間観測記録

- 対象環境: `npm run dev:wails:run`、`http://localhost:34115`、ジョブ#3、単語翻訳フェーズ
- 観測 1: AI サービス select の選択肢に「選んでください / Gemini / LM Studio / xAI」が並ぶ（provider-settings 経路は機能）
- 観測 2: 「Gemini」を選択 → DOM select の value が "gemini" になる（evaluate_script で確認）
- 観測 3: 「モデル一覧を更新」ボタン押下 → モデル select は「選んでください / モデル一覧を更新してください」のまま変わらず
- 観測 4: backend `npm run dev:wails:run` stdout に `ListTranslationJobSetupProviderModels` 呼び出し履歴なし
- 観測 5: ブラウザ console に error / warn なし
- 期待との差分: dev 環境の fake provider は ListTranslationJobSetupProviderModels に対して fake モデル一覧を返すはずなので、モデル select の選択肢に fake モデルが並ぶことが期待される。

## 前 task からの引き継ぎ

- 前 task `fix-term-translation-model-settings-empty-fixed` の主修正（snapshot 直書き廃止、JOB_PHASE_AI_SETTINGS 新設、状態 pill の「設定未完了」分岐、backend-local / frontend-local harness 通過）は完了済み。
- 前 task の途中で 2 件の追加 frontend 修正（provider 一覧の viewModel 流入、モデル一覧取得経路の追加）を実施したが、後者の経路が実画面で機能しないことが本 task の対象。
- 前 task の plan / fix-decision / detail-spec-diff / design-diff / implementation-scope はそのまま正本として残し、本 task の検討対象から除外する。

## 想定 Y/N 評価

| 想定 | Y/N | 根拠 |
| --- | --- | --- |
| 仕様変更または仕様追加がある | N | 画面仕様（`docs/screen-design/screens/term-translation-phase.md`）はすでに provider 選択 → モデル一覧取得 → モデル選択 → 固定 の遷移を定義済み。 |
| 画面変更がある | N | 画面構造・layout・文言の変更なし。 |
| 内部構造変更がある | Y | provider 選択値と controller / refreshModelList の配線が機能していないため、state と props の流れの修正が要る。 |
| 画面の表示変更がある | N | 状態 pill / select / ボタン文言は変更不要。 |
| frontend ロジック変更がある | Y | panel の onProviderChange / onRefresh、controller の refreshModelList、gateway 呼び出し経路のどこかに不具合がある。 |
| backend 変更がある | 判断保留 | backend に呼び出しが届いているか investigation-module で確定する。 |
| frontend と backend を接続する | 判断保留 | 同上。 |
| 実装済み責務を独立に証明したい | Y | controller の refreshModelList を単体テストで証明する。 |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | provider 選択時の DOM value、Svelte 5 reactivity、controller listener 通知の経路を実行時に観測して固定する。 |

「仕様変更または仕様追加がある」が N のため、本モジュールを継続する。

## Wails 接続対象

- 起動 command: `npm run dev:wails:run`
- 接続先: `http://localhost:34115`

## 修正方針確定の背景

- 初回の fix-decision では「controller が `credential_missing` 時に viewModel へ status を渡してエラーメッセージを表示する」フロントエンド修正を採用したが、人間レビューで差し戻しとなった。
- 差し戻し理由: 症状を隠す対症療法であり、根本原因（`AI_MODE=fake` 環境で secret store 抽象化が貫徹していないこと）を解消していない。正しい修正は secret store 層に `FakeProviderSettingsSecretStore` を追加して fake 環境の一貫性を持たせること。

## 後続モジュールへの引き継ぎ

- 入口: investigation-module。
- 引き継ぐ事実: 本 plan.md の人間観測記録と想定 Y/N 評価。前 task の plan / fix-decision は参考として参照可能だが、本 task の修正範囲外として扱う。

## 修正実行入力（人間レビュー 2026-06-02 承認済み）

### 承認済み修正方針

`fix-decision.md` の 6 方針を採用する。要旨:
1. 共通 `FakeSecretStore` 1 ファイル新規（`SecretStore` / `ProviderSettingsSecretStore` 両 interface を満たす、Load 固定値、Save/Delete no-op）
2. `InMemorySecretStore` 本体完全削除
3. 影響テスト 10 関数削除、5 箇所を fake に直接置換
4. bootstrap (`app_controller.go`) の `SECRET_BACKEND` 値 `in-memory` 削除、`fake` 分岐追加
5. `scripts/dev/run-wails.sh` の env を `=fake` に変更
6. `sonar-project.properties` の `sonar.coverage.exclusions` に新規 fake ファイル追記、fake にはテスト書かない

### 禁止する修正

- サービス層に provider 種別判定を追加して credential 解決をスキップする
- ProviderModelListLoader に credential 解決責務を下ろす
- frontend で `credential_missing` を成功扱いに書き換える
- provider-settings 画面の「利用可能」判定に secret store チェックを追加する（本 task 範囲外）

### 影響ファイル候補

- 新規: `internal/repository/fake_secret_store.go`
- 編集: `internal/repository/master_persona_repository.go`（`InMemorySecretStore` 完全削除）
- 編集: `internal/bootstrap/app_controller.go`（in-memory 分岐削除、fake 分岐追加）
- 編集: `scripts/dev/run-wails.sh`（env 値変更）
- 編集: `sonar-project.properties`（exclusion 追加）
- テスト編集: fix-decision.md に列挙した 10 関数削除 + 5 箇所置換

### 承認済み UC 差分

差分なし（`uc-translation-management.md` の既存 E1 で説明可能）。

### 承認済み E2E テスト観点差分

`test-design.csv` の E2E-UC-FAKE-001〜006 を採用:
- `AI_MODE=fake / SECRET_BACKEND=fake` で provider 選択 → モデル一覧更新 → fake モデル表示（3 phase 共通、001〜003）
- fake モデル選択 → 処理方式選択 → AI 設定固定 → 開始ボタン有効化までの画面操作確認（3 phase 共通、004〜006）

### 画面再現手順と修正後の期待状態

- 環境: `npm run dev:wails:run`、http://localhost:34115、ジョブ#3 単語翻訳フェーズ
- 操作: AI サービス select で Gemini 選択 → 「モデル一覧を更新」ボタン押下
- 修正前: モデル select は「選んでください / モデル一覧を更新してください」のまま変わらない
- 修正後: モデル select に fake モデル選択肢（例: `fake-model`）が表示され選択可能になる。選択 → 処理方式 → 固定 → 開始ボタン有効化まで踏める。

### 後続入口

implementation-module（backend 中心、frontend は前 task で追加した経路がそのまま機能する想定）。

## 最終検証結果（2026-06-02）

- `python3 scripts/harness/run.py --suite backend-local`: 全パッケージ通過。
- 実画面確認（`npm run dev:wails:run`、http://localhost:34115、ジョブ#3 単語翻訳フェーズ）:
  - AI サービス select に Gemini / LM Studio / xAI が表示される（観測 1 と同じ）
  - Gemini 選択 → 「モデル一覧を更新」ボタン押下 → モデル select に `fake-model` 選択肢が表示される（**修正前は変化なし、修正後は fake モデル表示成功**）
- 本 task の達成点: fake モデルがモデル select に表示されるまで。
- 残課題（別 task へ切り出し）: モデル選択後も AI 設定 pill が「設定未完了」のまま、開始ボタンが disabled のまま。前 task で追加した frontend 経路（saveAISettings の呼び出しシーケンス）に追加修正が要る。本 task の修正範囲外として扱う。

## 正本化判断（2026-06-02）

- 仕様変更または仕様追加: なし（想定 Y/N 評価で N、`uc-translation-management.md` 既存 E1 で説明可能）。
- 詳細仕様正本反映: 不要。
- 人間承認状態: 不要（仕様変更なし）。

## 作業 commit

- commit hash: f2fbb14ab9a01eea7c2830d792eed1fe85cab42e
- branch: claude/fix-term-translation-model-settings-empty-fixed
- 前 task `fix-term-translation-model-settings-empty-fixed` と同 commit に統合。

