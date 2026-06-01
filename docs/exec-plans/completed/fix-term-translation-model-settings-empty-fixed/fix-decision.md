# 修正方針判断

- task-id: fix-term-translation-model-settings-empty-fixed
- 作成日: 2026-06-01
- 改訂日: 2026-06-01（人間レビュー差し戻し対応）

## 判断結果

- 判定: 完了
- 停止理由: なし

## 観測済み問題

- 問題 1: 単語翻訳フェーズの AI モデル設定パネルで、AI サービス/モデル/処理方式がすべて空欄（選択肢なし）であるにもかかわらず、状態 pill に「固定済み」と表示される。
- 問題 2: 「開始」ボタンが無効化され、「実行設定が未構成のため開始できません。」と表示される。
- 問題 3: 「中断」ボタンが無効化され、「フェーズが実行中ではありません。」と表示される。
- 問題 4: AI サービス/モデル/処理方式の select 要素が空文字値の選択肢 1 件のみで、ユーザーが選択を変更できない。
- 期待との差分: 仕様（`docs/screen-design/screens/term-translation-phase.md` 107-113 行）では、モデル未選択状態は「AI 設定不足: 認証不足またはモデル未選択の警告を表示する」状態に分岐し、ユーザーが設定できる経路が確保されるべきである。

## 画面再現確認

- Wails 接続対象: `http://localhost:34115`
- 再現手順:
  1. `http://localhost:34115/#translation-management` へ移動する。
  2. ジョブ#3（Lucien.esp_Export.json）のリンクをクリックしてジョブ実行画面へ移動する。
  3. 単語翻訳フェーズの AI モデル選択領域を確認する。
- 操作結果:
  - `AIサービス` combobox: `value=""`, options `[{value:"", text:""}]` の空1件のみ。
  - `モデル` combobox: `value=""` で未選択、options に空エントリが混在。
  - `処理方式` combobox: `value=""`, options `[{value:"", text:""}]` の空1件のみ。
  - `status-pill.status-success` に「固定済み」が表示されている。
  - 「開始」「中断」ボタンが無効化されている。
- 画面状態: 人間観測記録の観測 1〜5 をすべて再現した。
- 証跡 path:
  - `tmp/logs/screenshot-job3-ai-model-panel.png`（AI モデルパネルと「固定済み」pill が確認できるスクリーンショット）
  - `tmp/logs/screenshot-job3-term-translation.png`（ジョブ実行画面全体のスクリーンショット）

## 原因仮説と検証結果

### 仮説 A: backend が `execution.model` を空文字で返している

- 根拠: presenter の `summary?.execution.model ?? "-"` は `null`/`undefined` のみ `"-"` にフォールバックし、空文字 `""` は通り抜ける。
- 検証: 一時観測ログ（`TEMP-OBS`）を `loadExecutionContext` 末尾に追加し、`ReadSummary` 呼び出し時のログを観測した。
- 観測結果: `execution_model:""`, `execution_provider:""`, `execution_mode:""`, `execution_credential_ref:""`, `run_exists:false` が記録された。
- 判定: **確定**。backend が `execution.model` を空文字で返している。

### 仮説 B: `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` パスが原因

- 根拠: `run_exists:false` かつジョブ状態 `ready` のとき、`termTranslationExecutionBasePhase` は `applyRuntimeSnapshot=true` を返す。runtime snapshot が存在しない場合、`ErrNotFound` パスで `initial.AIProvider/ModelName/ExecutionMode/CredentialRef` をすべて空文字で上書きする（`term_translation_phase_service.go` 1468-1473 行）。
- 検証: 追加観測ログで `term_translation_runtime_snapshot_not_found`, `reason:"snapshot_not_found"` が記録された。
- 判定: **確定**。runtime snapshot が存在しないため `ErrNotFound` パスを通り、空文字が確定した。

### 仮説 C: Svelte の判定条件が空文字を検出できていない

- 根拠: `TermTranslationPhasePanel.svelte` 64 行の `viewModel.modelLabel === "-"` は、値が空文字 `""` の場合に偽になる。
- 検証: JavaScript で `aiServiceSelect.value === ""` が確認された。
- 判定: **確定（症状の一つ）**。backend が空文字を返す結果として Svelte の判定が偽になる。backend を修正すれば Svelte 側の変更は不要になる見込みであり、本仮説は症状の連鎖として扱う。根本ではない。

## 確定原因

### 確定原因（backend 層）: `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` 処理が仕様逸脱

- 原因: `applyTermTranslationRuntimeSnapshot`（`internal/service/term_translation_phase_service.go` 1468-1473 行）は、runtime snapshot が存在しない（`ErrNotFound`）場合に `initial.AIProvider`, `initial.ModelName`, `initial.ExecutionMode`, `initial.CredentialRef` をすべて空文字で上書きする。
- 仕様との乖離: `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` テーブルは「フェーズ共有 AI 設定」の保存先として機能する設計であり、`SaveAISettings` がユーザーの設定操作をここへ書き込む。snapshot が存在しない状態はユーザーがまだ設定を保存していない状態を意味する。この場合に `initial`（`translation` フェーズからの引き継ぎ値）を空文字で上書きすることは、「設定がない場合は設定前の状態を維持する」という仕様の意図に反する。正しい挙動は、snapshot が存在しない場合に `initial` をそのまま返し、「未設定」を表現する空の `JobPhaseRun` フィールドをゼロ値のまま維持することである。
- 観測根拠:
  - 観測ログ `term_translation_runtime_snapshot_not_found`, `reason:"snapshot_not_found"` が出た。
  - 続く `term_translation_execution_context_loaded` で `execution_model:""` が出た。
  - `SaveAISettings` を呼んでいない新規ジョブでは snapshot が存在せず、常にこのパスを通る。

### 追加確定原因（症状の連鎖）

- `applyTermTranslationRuntimeSnapshot` が空文字を返すことで、`execution.Provider`, `execution.Model`, `execution.ExecutionMode` がすべて空文字になる。
- presenter の `?? "-"` は `null`/`undefined` のみ `"-"` にフォールバックするため、空文字 `""` を `"-"` に変換しない。
- `TermTranslationPhasePanel.svelte` の `viewModel.modelLabel === "-"` が偽になり `aiSettingsBlockedReason = ""` となる。その結果、`aiSettingsStatusLabel = "固定済み"` が表示される。
- これらは backend の空文字出力に起因する連鎖症状であり、backend を修正すれば presenter・Svelte の変更なしに解消する見込みである。

## 採用する修正方針（恒久修正）

### 方針 1: `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` パスを「`initial` をそのまま返す」に修正する

- 方針: runtime snapshot が存在しない場合は `initial` を空文字で上書きせず、そのままの値を返す。
- 理由: snapshot が存在しないことは「ユーザーが設定を保存していない状態」であり、`initial` の値を破壊する理由はない。`initial` には `translation` フェーズからの引き継ぎ値が入っており、それをゼロ値に上書きすることは情報の破壊である。引き継ぎ値がない場合（`termTranslationInitialExecutionPhase` が `ErrNotFound` を返した場合）、`termTranslationExecutionBasePhase` はゼロ値の `JobPhaseRun` を `initial` として返すため、空文字上書きは冗長かつ有害である。
- 変更箇所: `internal/service/term_translation_phase_service.go` 1468-1473 行の `ErrNotFound` ブロック内の空文字上書きを削除し、`return initial, nil` に変更する。
- 変更後の挙動:
  - snapshot 存在: snapshot の値（`Provider`, `ModelName`, `ExecutionMode`）を `initial` に上書きして返す（現行の 1478-1481 行）。
  - snapshot 非存在: `initial` をそのまま返す（修正後）。
  - `initial` がゼロ値の場合は空文字のまま返り、presenter・Svelte は「未設定」として正しく判定できる状態になる。

### 方針 2: presenter と Svelte は修正しない

- 方針: backend が正しく「設定あり」または「設定なし（ゼロ値）」を返すようになれば、presenter の `?? "-"` は `null`/`undefined` の場合のみフォールバックするため、`initial` がゼロ値で返った場合は空文字 `""` が返る。
- 追加確認事項: presenter の `isExecutionConfigured` は `trim()` を使って空文字判定をしており、backend 修正後に `CanStart` の判定が正しく `false` になることを確認する。Svelte の `aiSettingsBlockedReason` 判定が backend 修正後に正しく「設定未完了」を返すかを実画面で確認する。
- 判断根拠: presenter・Svelte の変更は「backend の誤りを隠す対症療法」であり、人間レビューで明示的に禁止された方針である。

### 影響範囲の確定

- 本 task は単語翻訳フェーズ（`word_translation`）の修正に限定する。
- persona-generation（`npc_persona_generation`）と body-translation（`text_translation`）も同一の `ErrNotFound` 時空文字上書きパターンを持つことを観測した（`persona_generation_phase_service.go` 896-900 行、`body_translation_phase_service.go` 764-768 行）。これらは本 task の観測対象外であり、別 task として設計・修正する。

## 禁止する修正

- 禁止修正 1: presenter の `summary?.execution.model ?? "-"` を `?.trim() || "-"` に変更する修正。
  - 理由: backend の仕様逸脱を presenter で覆い隠す対症療法である。人間レビューで明示的に差し戻された方針であり、採用しない。

- 禁止修正 2: `TermTranslationPhasePanel.svelte` の判定を `viewModel.modelLabel === ""` に変更する修正。
  - 理由: presenter が空文字を `"-"` に変換しないまま Svelte の判定だけを変更する対症療法である。backend の誤りを隠す。

- 禁止修正 3: `aiSettingsStatusLabel` に新しい状態値（例: `"未設定"` など）を追加する修正。
  - 理由: 状態値の追加はモデルの乖離を生じさせ、後続の状態依存ロジックに波及する。既存の `"固定済み"` / `"設定未完了"` の 2 値で仕様を満たせる。

- 禁止修正 4: backend が空文字ではなく `null` を返すよう型を変更する修正。
  - 理由: `TermTranslationExecutionConfigReadModel.Model` は `string` 型（非ポインタ）であり、Go のゼロ値として空文字は正当である。型変更は広範な層へ影響し、今回の修正範囲を逸脱する。

- 禁止修正 5: persona-generation・body-translation の同一パターンを本 task で同時修正する修正。
  - 理由: 本 task の観測・再現・検証範囲は単語翻訳フェーズに限定されており、他フェーズへの修正は別 task で設計・観測する。

## 影響ファイル候補

| ファイル | 理由 |
| --- | --- |
| `internal/service/term_translation_phase_service.go` | `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` ブロック（1468-1473 行）を修正する。空文字上書きを削除し `return initial, nil` に変更する。 |
| `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts` | backend 修正後に `isExecutionConfigured` の判定と `providerLabel`/`modelLabel`/`executionModeLabel` の生成が正しく動くかを確認する。変更は不要の見込みだが、確認対象として把握する。 |
| `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` | backend 修正後に `aiSettingsBlockedReason` 判定（61-73 行）が正しく「設定未完了」を返すかを確認する。変更は不要の見込みだが、確認対象として把握する。 |

## 別 task として後続対応が必要な候補

| ファイル | 理由 |
| --- | --- |
| `internal/service/persona_generation_phase_service.go` | 896-900 行に同一の `ErrNotFound` 時空文字上書きパターンが存在する。本 task の観測範囲外のため別 task で設計・観測する。 |
| `internal/service/body_translation_phase_service.go` | 764-768 行に同一の `ErrNotFound` 時空文字上書きパターンが存在する。本 task の観測範囲外のため別 task で設計・観測する。 |
