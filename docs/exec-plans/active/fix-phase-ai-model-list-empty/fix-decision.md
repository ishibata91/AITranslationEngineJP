# Fix Decision Report

## 判断結果

- 判定: 完了（人間レビューによる再訂正後、全面書き換え済み）
- 停止理由: なし

---

## 観測済み問題

- 問題 1: job-run 画面の 3 phase（単語翻訳 / ペルソナ生成 / 本文翻訳）AI 設定パネルで、provider を選択して「モデル一覧を更新」を押してもモデル select の選択肢が「選んでください / モデル一覧を更新してください」のままで変わらない。
- 問題 2: provider-settings 画面では Gemini が「利用可能」（接続確認済み）と表示されているのに、job-run の AI 設定パネルでは `credential_missing` が返る。この矛盾がユーザーから見ると不可解な状態になる。

---

## 画面再現確認

- Wails 接続対象: `http://localhost:34115`（`npm run dev:wails:run` で起動済み）
- 再現手順: ジョブ#3 を選択して job-run 画面へ遷移 → 単語翻訳フェーズの AI 設定パネルで Gemini を選択 → 「モデル一覧を更新」ボタンを押下
- 操作結果: `refreshModelList_response` で `status: "credential_missing"`, `modelsCount: 0`, `failureKind: "model_list_credential_missing"`（一時観測ログ msgid=1090 で確認）
- 画面状態: モデル select は「選んでください / モデル一覧を更新してください」のまま変わらない
- 証跡参照: 前回観測済み（chrome-devtools MCP による再現確認は本 task 差し戻し対応で省略）

---

## 確定原因

### 確定原因: secret store 抽象化の不徹底による判定根拠の食い違い

- 原因: `ProviderSettingsService.ListProviderModels` 経路で credential 解決が失敗して `credential_missing` を返す。根本は secret store 抽象化の不徹底にある。
  - `InMemorySecretStore` は `SecretStore` interface と `ProviderSettingsSecretStore` interface の両方を満たす単一実装であり、master_persona と provider-settings が共有している。`AI_MODE=fake` 環境でも `Load` は空文字を返すだけ。
  - `ProviderClient` と `ProviderModelListLoader` は `AI_MODE=fake` で deterministicProvider に置換され credential 不要で動く。
  - secret store だけ抽象化が貫徹していないため、サービス層の credential ガード（`resolveValidationSecret` → `secretStore.Load`）に引っかかる。
  - 結果として provider-settings 画面（DB の `validation_state` だけを参照して「利用可能」表示）と job-run 画面（secret store 実値を要求）の判定根拠が食い違い、ユーザーから見ると矛盾した状態になる。
- 観測根拠:
  - DB クエリ: `gemini` の `credential_reference_id = "provider-settings:gemini"` が non-nil（保存済み）、`validation_state = "validated"`（接続確認済み）。
  - 観測ログ: `refreshModelList_response` で `status: "credential_missing"`, `modelsCount: 0`。
  - `run-wails.sh` コード確認: `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND="in-memory"` が明示的に設定されており、プロセス再起動で store がリセットされる。
  - `app_controller.go` コード確認: `case providerSettingsSecretBackendInMemory` は `NewProviderSettingsInMemorySecretStore()` を返し、fake モードでの credential 不要化が secret store 層まで伝わっていない。

---

## 採用する修正方針

### 方針 1: 共通 `FakeSecretStore` を 1 つだけ新規実装する

- 対象: `internal/repository/fake_secret_store.go`（新規）
- 方針: `SecretStore` interface と `ProviderSettingsSecretStore` interface の両方を満たす単一実装とする。`Load(_, key)` は provider 種別を問わず固定値（`"fake-secret"`）を返す。`Save` と `Delete` は no-op とする。master_persona 専用 fake は作らない。
- 理由: `AI_MODE=fake` 環境で secret store が空文字を返すことでサービス層の credential ガードに引っかかる根本を除去する。fake な provider が credential 不要で動くように、secret store 層も fake 化して一貫性を持たせる。両インターフェースを 1 つの実装で満たすことで、master_persona と provider-settings が共有していた元の構造を維持できる。

### 方針 2: `InMemorySecretStore` 本体を完全削除する

- 削除対象: `internal/repository/master_persona_repository.go` から以下をすべて削除する。
  - `InMemorySecretStore` 型定義
  - `NewInMemorySecretStore` コンストラクタ
  - `NewProviderSettingsInMemorySecretStore` エイリアスコンストラクタ
  - `Load` / `Save` / `Delete` メソッド

### 方針 3: Save/Load 振る舞いをアサートしているテスト関数を削除する

削除対象テスト関数（計 10 件）は以下のとおり。

| ファイル | テスト関数名 | 削除理由 |
| --- | --- | --- |
| `internal/repository/master_persona_keyring_secret_store_test.go` | `TestInMemorySecretStoreSaveLoadDelete` | `InMemorySecretStore` 本体の Save/Load/Delete 動作をアサートしており、本体削除後に存在できない |
| `internal/repository/provider_settings_keyring_secret_store_test.go` | `TestProviderSettingsInMemorySecretStoreSaveLoadDelete` | 同上 |
| `internal/repository/provider_settings_keyring_secret_store_test.go` | `TestProviderSettingsInMemorySecretStoreIsProcessLocalAcrossNewInstance` | 同上 |
| `internal/bootstrap/app_controller_test.go` | `TestNewProviderSettingsSecretStoreFromEnvInMemoryBypassesKeyringOpen` | `providerSettingsSecretBackendInMemory` 定数と `in-memory` 分岐の動作をアサートしており、方針 4 で当該定数と分岐を削除するため存在できない |
| `internal/service/provider_settings_service_test.go` | `TestProviderSettingsServiceSavePreservesSecretBoundary` | `secretStore.Load` の返り値（`loadedSecret`）を直接アサートしており、固定値返しの `FakeSecretStore` では証明できない |
| `internal/service/provider_settings_service_test.go` | `TestProviderSettingsServiceSavePreservesExistingSecretWhenInputMissing` | 同上（`loadedSecret != "xai-secret"` をアサート） |
| `internal/service/provider_settings_service_test.go` | `TestProviderSettingsServiceSaveRejectsInvalidEndpointWithoutMutatingState` | 同上（`loadedSecret != "stored-secret"` をアサート） |
| `internal/service/provider_settings_service_test.go` | `TestProviderSettingsServiceResetKeepsRowAndDeletesSecret` | 同上（`loadedSecret != ""` をアサート） |
| `internal/service/master_persona_service_test.go` | `TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverSaveSettingsDoesNotSaveAPIKey` | `secretStore.Load` の返り値（`storedKey`）を直接アサートしており、固定値返しの `FakeSecretStore` では証明できない |
| `internal/service/master_persona_service_test.go` | `TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverLoadSettingsDoesNotRestoreAPIKeyFromSecretStore` | `secretStore.Save` で実値をセットした後、`service.LoadSettings` 経由で API キーが露出しないかをアサートしており、no-op `Save` の `FakeSecretStore` では証明できない |

seed 用途のみで Save/Load をアサートしていない 5 箇所（下記）は削除せず、`FakeSecretStore` に直接置き換える。

| ファイル | 箇所 | 差し替え内容 |
| --- | --- | --- |
| `internal/bootstrap/app_controller_test.go` 135 行 | `NewInMemorySecretStore` をソース文字列検索するアサーション | アサーションを削除する（対象実装が存在しなくなるため） |
| `internal/bootstrap/app_controller_test.go` 1164 行 | `inMemorySecretStore := repository.NewInMemorySecretStore()` | `FakeSecretStore` のコンストラクタ呼び出しに置き換え |
| `internal/apitest/observability_log_test_helpers_test.go` 239 行 | `NewInMemorySecretStore()` | `FakeSecretStore` に置き換え |
| `internal/apitest/model_settings_card_fake_mode_test.go` 150 行 | `NewInMemorySecretStore()` | `FakeSecretStore` に置き換え |
| `internal/apitest/job_setup_relocation_scenario_test.go` 93, 105 行 | `secretStore` フィールド型が `*repository.InMemorySecretStore`、105 行で `secretStore.Save` を呼ぶ | フィールド型をインターフェース型に変更、`FakeSecretStore` に置き換え、`Save` 呼び出しを削除 |
| `internal/service/term_translation_phase_scenario_test.go` 310 行 | `NewInMemorySecretStore()` | `FakeSecretStore` に置き換え |

### 方針 4: bootstrap の `in-memory` 分岐を `fake` 分岐へ置き換える

- 対象: `internal/bootstrap/app_controller.go` の `newProviderSettingsSecretStoreFromEnv` 関数。
- 方針: `providerSettingsSecretBackendInMemory` 定数を削除し、`case providerSettingsSecretBackendInMemory` 分岐を削除する。`case "fake"` 分岐を追加して `FakeSecretStore` を返す。

### 方針 5: dev 起動 script の環境変数を更新する

- 対象: `scripts/dev/run-wails.sh`。
- 方針: `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND=in-memory` を `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND=fake` に変更する。

### 方針 6: `sonar-project.properties` の `sonar.coverage.exclusions` に新規 fake 実装ファイルを追加する

- 対象: `sonar-project.properties` の `sonar.coverage.exclusions` 末尾。
- 追加する具体行: 現在の最終行 `internal/repository/sqlite_db.go` の後に `internal/repository/fake_secret_store.go` を追記する。
- 理由: fake 実装はテスト補助用であり、カバレッジ対象から構造的に除外する。fake 実装に対する単体テストは書かない。

---

## 禁止する修正

- 禁止修正 1: サービス層（`ProviderSettingsService`）に provider 種別判定を追加して credential 解決をスキップする
  - 理由: サービス層の責務集中が崩れる。credential 解決のスキップ判定はサービス層ではなく secret store 層で行うべき。
- 禁止修正 2: `ProviderModelListLoader` 内で credential 解決を行うように責務を下ろす
  - 理由: 責務分離が崩れる。credential 解決はサービス層が担う設計であり、loader に責務を移すと将来の変更コストが上がる。
- 禁止修正 3: 既存 `InMemorySecretStore` を残して dev 起動時に実値を seed する
  - 理由: dev 起動 script から API キーの実値を注入する仕組みになり、誤設定リスクが上がる。fake store で固定値を返す方が意図が明確。
- 禁止修正 4: frontend で `credential_missing` を無視して成功扱いにする
  - 理由: backend の正しい credential チェックを無効化する対症療法。
- 禁止修正 5: provider-settings 画面の「利用可能」判定に secret store のチェックを追加する
  - 理由: provider-settings 画面の仕様変更に相当し、本 task の修正範囲外。
- 禁止修正 6: `InMemorySecretStore` 本体を残して一部テストのみ fake に切り替える
  - 理由: 人間レビューで「完全削除」が確定方針として採用されており、本体を残すことは方針に反する。
- 禁止修正 7: master_persona 専用の `FakeMasterPersonaSecretStore` を別途作る
  - 理由: 人間レビューで確定した再訂正方針により、master_persona 専用 fake は不要と確定している。`SecretStore` と `ProviderSettingsSecretStore` の両インターフェースを満たす共通 `FakeSecretStore` 1 つで足りる。

---

## 影響ファイル候補

### 新規追加
- `internal/repository/fake_secret_store.go`: `FakeSecretStore` 実装（`SecretStore` と `ProviderSettingsSecretStore` の両インターフェースを満たす）

### 変更
- `internal/repository/master_persona_repository.go`: `InMemorySecretStore` 型定義、`NewInMemorySecretStore`、`NewProviderSettingsInMemorySecretStore`、`Load`/`Save`/`Delete` メソッドを削除
- `internal/bootstrap/app_controller.go`: `providerSettingsSecretBackendInMemory` 定数削除、`in-memory` 分岐削除、`fake` 分岐追加
- `scripts/dev/run-wails.sh`: `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND` を `fake` に変更
- `sonar-project.properties`: `sonar.coverage.exclusions` に `internal/repository/fake_secret_store.go` を追加

### テスト変更（削除、計 10 件）
- `internal/repository/master_persona_keyring_secret_store_test.go`: `TestInMemorySecretStoreSaveLoadDelete` を削除
- `internal/repository/provider_settings_keyring_secret_store_test.go`: `TestProviderSettingsInMemorySecretStoreSaveLoadDelete` を削除
- `internal/repository/provider_settings_keyring_secret_store_test.go`: `TestProviderSettingsInMemorySecretStoreIsProcessLocalAcrossNewInstance` を削除
- `internal/bootstrap/app_controller_test.go`: `TestNewProviderSettingsSecretStoreFromEnvInMemoryBypassesKeyringOpen` を削除
- `internal/service/provider_settings_service_test.go`: `TestProviderSettingsServiceSavePreservesSecretBoundary` を削除
- `internal/service/provider_settings_service_test.go`: `TestProviderSettingsServiceSavePreservesExistingSecretWhenInputMissing` を削除
- `internal/service/provider_settings_service_test.go`: `TestProviderSettingsServiceSaveRejectsInvalidEndpointWithoutMutatingState` を削除
- `internal/service/provider_settings_service_test.go`: `TestProviderSettingsServiceResetKeepsRowAndDeletesSecret` を削除
- `internal/service/master_persona_service_test.go`: `TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverSaveSettingsDoesNotSaveAPIKey` を削除
- `internal/service/master_persona_service_test.go`: `TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverLoadSettingsDoesNotRestoreAPIKeyFromSecretStore` を削除

### テスト変更（差し替え、計 6 箇所）
- `internal/bootstrap/app_controller_test.go` 135 行: `NewInMemorySecretStore` を検索するアサーションを削除
- `internal/bootstrap/app_controller_test.go` 1164 行: `NewInMemorySecretStore()` を `FakeSecretStore` のコンストラクタ呼び出しに置き換え
- `internal/apitest/observability_log_test_helpers_test.go` 239 行: `NewInMemorySecretStore()` を `FakeSecretStore` に置き換え
- `internal/apitest/model_settings_card_fake_mode_test.go` 150 行: `NewInMemorySecretStore()` を `FakeSecretStore` に置き換え
- `internal/apitest/job_setup_relocation_scenario_test.go` 93, 105 行: `secretStore` フィールド型をインターフェース型に変更、`FakeSecretStore` に置き換え、`Save` 呼び出しを削除
- `internal/service/term_translation_phase_scenario_test.go` 310 行: `NewInMemorySecretStore()` を `FakeSecretStore` に置き換え

---

## E2E 観点差分

| 観点 | 前回方針 | 今回方針 |
| --- | --- | --- |
| モデル一覧更新後の選択肢表示 | `FakeProviderSettingsSecretStore` が固定値を返すため `credential_missing` が発生しなくなる | 共通 `FakeSecretStore` が固定値を返すため同等の効果 |
| secret store の取りうる実装 | fake と keychain 系のみ | 同左（`in-memory` は廃止） |
| テスト補助 KVS | テスト内専用 KVS（`provider_settings_service_test.go` と `master_persona_service_test.go` に各定義）に移行 | Save/Load 振る舞いアサート関数ごと削除。seed のみ用途は共通 `FakeSecretStore` に直接置き換え |
| fake 実装ファイル数 | 2 ファイル（`FakeProviderSettingsSecretStore` と `FakeMasterPersonaSecretStore`） | 1 ファイル（共通 `FakeSecretStore`） |
| カバレッジ除外 | 2 ファイルを追加 | 1 ファイル（`fake_secret_store.go`）を追加 |
| テスト設計への影響 | 削除テスト関数の UC 観点が失われないか確認が要る | Save/Load 振る舞いアサート関数の削除により、secret store の KVS 的動作の直接証明は失われる。ただし当該テストは secret store 実装の動作証明であり、サービス層の UC 観点（credential 秘匿、credential 保存境界）は残存する別テスト関数で引き続き証明される |
