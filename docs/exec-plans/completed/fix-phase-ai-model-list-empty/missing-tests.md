# 不足テスト

## 概要

- 対象: fix-phase-ai-model-list-empty の E2E テスト観点差分確認
- 判断: 追加候補あり
- 根拠: 今回の修正は `FakeProviderSettingsSecretStore` を新規実装し、bootstrap の `SECRET_BACKEND=fake` 分岐で返す変更である。既存の test-design.csv に、`AI_MODE=fake` かつ `SECRET_BACKEND=fake` 環境でモデル一覧が正常に返ることを確認する観点がない。また `FakeProviderSettingsSecretStore.Load` の固定値返却と bootstrap 分岐の単体テスト観点も不足している。

## 不足テスト一覧

| ID | 関連UC | 対象画面 | 不足観点 | 理由 | 追加候補 |
| --- | --- | --- | --- | --- | --- |
| GAP-001 | 翻訳段階を開始する | 単語翻訳 | AI_MODE=fake かつ SECRET_BACKEND=fake 環境で Gemini を選択後にモデル一覧を更新すると fake モデルが選択肢に並ぶ | 既存 E2E に fake 環境でのモデル一覧取得成功確認がない | E2E-UC-FAKE-001 |
| GAP-002 | 翻訳段階を開始する | NPC ペルソナ生成段階 | AI_MODE=fake かつ SECRET_BACKEND=fake 環境で Gemini を選択後にモデル一覧を更新すると fake モデルが選択肢に並ぶ | 既存 E2E に fake 環境でのモデル一覧取得成功確認がない | E2E-UC-FAKE-002 |
| GAP-003 | 翻訳段階を開始する | 本文翻訳段階 | AI_MODE=fake かつ SECRET_BACKEND=fake 環境で Gemini を選択後にモデル一覧を更新すると fake モデルが選択肢に並ぶ | 既存 E2E に fake 環境でのモデル一覧取得成功確認がない | E2E-UC-FAKE-003 |
| GAP-004 | 翻訳段階を開始する | 単語翻訳 | fake モデル選択 → 処理方式選択 → AI 設定固定 → 開始ボタン有効化まで画面操作で踏める | fake モデル選択後の遷移全体を確認する観点がない | E2E-UC-FAKE-004 |
| GAP-005 | 翻訳段階を開始する | NPC ペルソナ生成段階 | fake モデル選択 → 処理方式選択 → AI 設定固定 → 開始ボタン有効化まで画面操作で踏める | fake モデル選択後の遷移全体を確認する観点がない | E2E-UC-FAKE-005 |
| GAP-006 | 翻訳段階を開始する | 本文翻訳段階 | fake モデル選択 → 処理方式選択 → AI 設定固定 → 開始ボタン有効化まで画面操作で踏める | fake モデル選択後の遷移全体を確認する観点がない | E2E-UC-FAKE-006 |

## 単体テスト不足一覧

| ID | 対象 | 不足観点 | 理由 |
| --- | --- | --- | --- |
| UNIT-GAP-001 | FakeProviderSettingsSecretStore.Load | provider 種別を問わず固定値（例: "fake-secret"）を返す | FakeProviderSettingsSecretStore は新規実装であり、固定値返却の挙動を証明する単体テストがない |
| UNIT-GAP-002 | bootstrap newProviderSettingsSecretStoreFromEnv | SECRET_BACKEND=fake 時に FakeProviderSettingsSecretStore が返る | bootstrap の fake 分岐が正しく実装されていることを証明する単体テストがない |
