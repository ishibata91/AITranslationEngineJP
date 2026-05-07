# xAI 実画面観測

## 観測

- 対象: ジョブセットアップの単語翻訳モデル選択。
- provider: `xai`
- 起動 URL: `http://localhost:34115/#translation-management`
- 実行: `agent-browser` で xAI を選択し、モデル一覧更新を押下した。

## 修正前の停止点

- xAI 選択後、`単語翻訳のモデル一覧を更新` ボタンが disabled のままだった。
- 原因: `model-settings-card-policy.ts` が `credentialStatus === "missing"` の時に更新ボタンを無効化していた。
- 影響: `refreshPhaseModels()` の修正後も、xAI ではクリック前に止まっていた。

## 修正後の観測

- xAI 選択後、`単語翻訳のモデル一覧を更新` ボタンが有効になった。
- 更新後、モデル select に `fake-model` が表示された。
- console では `ListTranslationJobSetupProviderModels` が `provider:"xai"`、`status:"success"`、`models:[{"modelId":"fake-model"}]` を返した。

## 証跡

- screenshot: [xai-fake-model-visible.png](/Users/iorishibata/Repositories/AITranslationEngineJP/tmp/agent-browser/2026-05-07-model-selection-fake-gate/xai-fake-model-visible.png)

## 未確認

- `npc_persona_generation` と `text_translation` の xAI クリック確認は未実施。
