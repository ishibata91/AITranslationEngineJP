# 人間修正レビュー用 概要

investigation-module が固定した 3 成果物（修正方針判断 / UC 差分候補 / E2E テスト観点差分）の要点だけを集約した、人間レビュー用の短い資料です。詳細は各成果物を参照します。

## レビューで判断する内容

次の 3 点を 1 つずつ承認、または差し戻しするだけで足ります。

1. 採用する修正方針（恒久修正の中身）
2. UC 差分の扱い
3. E2E テスト観点の扱い

## 1. 修正方針判断

- 詳細成果物: [fix-decision.md](./fix-decision.md)

### 確定原因（観測で確定済み）

`internal/bootstrap/app_controller.go` で 3 phase（単語翻訳 / ペルソナ生成 / 本文翻訳）の service 組み立て時に、`JobPhaseAISettingsRepository` の注入が欠落しています。
そのため `SaveTermTranslationPhaseAISettings` などの保存呼び出しが backend で「phase ai settings repository is not configured」エラーを返します。
frontend usecase はエラーを捕まえて summary 取得を行わないため、`summary.aiSettings` は空のままで、`isExecutionConfigured` は false 固定、AI 設定 pill は「設定未完了」、開始ボタンは disabled に固定されます。

観測根拠: JS から `window.go.wails.AppController.SaveTermTranslationPhaseAISettings(...)` を直接呼び、上記エラーを目視確認しました（chrome-devtools 経由）。

### 採用する修正方針

- 主因修正（backend 1 ファイル）: `internal/bootstrap/app_controller.go` で 3 phase 全ての service に `JobPhaseAISettingsRepository` を注入します。
  - `WithTermTranslationJobPhaseAISettings`
  - `WithPersonaGenerationPhaseAISettingsRepository`
  - `WithBodyTranslationJobPhaseAISettingsRepository`
- secondary 修正（frontend 3 ファイル）: 3 phase の presenter `buildModelOptions` で `{value: "", label: "選んでください"}` placeholder を返さないようにします。`AIModelSelectionCard` 側が固定 placeholder を 1 つ描画するため、現状は空 option が 2 個重複しています。secondary は主因修正後の検証結果で取捨選択します。

### 禁止する修正（実装 agent に許可しない対症療法）

- frontend で `isExecutionConfigured` を強制的に true にする分岐の追加
- `saveAISettings` のエラーを握り潰して summary 取得を強行する
- 新しい状態フラグ（例: `localAISettingsSaved`）の追加
- backend エラーを隠す mock / fake の挿入

### 影響ファイル候補

| ファイル | 種別 |
| --- | --- |
| `internal/bootstrap/app_controller.go` | 主因（必須） |
| `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts` | secondary |
| `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts` | secondary |
| `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts` | secondary |

### 想定 Y/N 評価の更新

- 「backend 変更がある」: N（暫定）→ Y（bootstrap への repository 注入が必要）
- 「frontend と backend を接続する」: N のまま（Wails bridge 追加は不要、既存経路を使う）

## 2. UC 差分候補

- 詳細成果物: [uc-diff-candidates.md](./uc-diff-candidates.md)
- 分類サマリ: `記述不足` を含む。`新規判断必要` は含まない（仕様変更不要）。

### 不足箇所の要点

| No | 対象 UC | 不足内容 |
| --- | --- | --- |
| 1 | 翻訳段階を開始する | 主シナリオに「モデル選択 → `saveAISettings` 成功 → summary 反映 → pill「固定済み」→ 開始ボタン有効化」の状態遷移が未記述 |
| 2 | 翻訳段階を開始する | 例外フローに「backend `JobPhaseAISettingsRepository` 未設定で pill が変化しない」境界が未記述。本 task の恒久修正後は発生しない経路。UC 正本への恒久追記要否は finalization-module で判断 |
| 3 | 3 phase 共通 | 差分なし。3 phase は同一 UC で記述されている |

### 確認してほしい論点

- UC 正本への記述追加を本 task で実施するか（finalization-module の `updating-docs` で扱う）。investigation-module ではどちらでも進められます。

## 3. E2E テスト観点差分

- 詳細成果物: [e2e-test-aspect-diff.md](./e2e-test-aspect-diff.md)
- 分類サマリ: `追加候補あり`（`判断不足` なし）。

### 既存観点との関係

| 既存 ID | 対象画面 | 前提 AI 設定 | モデル選択 → pill 変化 操作 |
| --- | --- | --- | --- |
| `E2E-UC-045` | 単語翻訳 | 準備済み | 含まない |
| `E2E-UC-046` | NPC ペルソナ生成 | 準備済み | 含まない |
| `E2E-UC-047` | 本文翻訳 | 準備済み | 含まない |

既存 3 観点は「AI 設定済み前提で開始する」観点であり、本 task の修正対象（未設定 → モデル選択 → pill 切り替え → 開始ボタン有効化）を証明する観点が欠落しています。

### 追加候補

| 追加 ID | 対象画面 | 内容 |
| --- | --- | --- |
| `E2E-UC-053` | 単語翻訳 | 未設定状態からモデル選択 → pill「固定済み」→ 開始ボタン有効化 |
| `E2E-UC-054` | NPC ペルソナ生成 | 同型 |
| `E2E-UC-055` | 本文翻訳 | 同型 |

参照 selector（`*-ai-model-lock-state`、`*-start-button`）は 3 phase の画面設計書「E2E 固定 selector」表で全て定義済みのため、selector 追加は不要です。

### 範囲外

bootstrap の repository 注入の DI テストと、presenter `buildModelOptions` placeholder 重複解消の単体テストは、E2E テスト規約の対象外として `implementation-module` の `tests-unit` で扱います。

## レビュー後の流れ

承認時は次の `修正実行入力` が固定され、`implementation-module` へ引き継ぎます。

- 承認済み修正方針: 主因（bootstrap の repository 注入）と secondary（presenter placeholder 重複解消）
- 禁止する修正: 上記 4 項目
- 影響ファイル候補: 上記表
- 承認済み UC 差分: `記述不足` 2 件（恒久追記要否は finalization-module で判断）
- 承認済み E2E テスト観点差分: 追加候補 3 件（`E2E-UC-053` / `E2E-UC-054` / `E2E-UC-055`）
- 画面再現確認の手順と修正後の期待状態: fix-decision.md の「画面再現確認」節を引き継ぐ
