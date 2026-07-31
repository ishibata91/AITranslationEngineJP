import { describe, expect, it } from "vitest"
import {
  canOperateBatch,
  orderProviderModels,
  providerDefaults
} from "./translation-run-provider"

describe("翻訳提供元の画面状態", () => {
  const luna = "gpt-5.6-luna"
  // R-1-3: Luna を初期選択し、取得した他モデルを削除しないこと。
  it("OpenAIのモデル一覧ではLunaを先頭にして他モデルも残す", () => {
    expect(
      orderProviderModels("openai", ["gpt-other", luna, "gpt-another"])
    ).toEqual([luna, "gpt-other", "gpt-another"])
  })

  // R-1-4: OpenAI batch 以外では OpenAI の既定値を自動使用しないこと。
  it("同期とxAIはOpenAIのendpointとLunaを使わない", () => {
    expect(providerDefaults("sync")).toEqual({
      endpoint: "http://127.0.0.1:1234",
      models: [],
      model: ""
    })
    expect(providerDefaults("xai")).toEqual({
      endpoint: "https://api.x.ai",
      models: [],
      model: ""
    })
  })

  // R-1-6: OpenAI API キーが空なら batch 操作を許可しないこと。
  it("OpenAI APIキーが空ならBatch操作を許可しない", () => {
    expect(canOperateBatch("openai", "")).toBe(false)
    expect(canOperateBatch("openai", "  ")).toBe(false)
    expect(canOperateBatch("openai", "sk-test")).toBe(true)
  })

  // R-1-8: 各提供元への切替値は独立し、切替前のモデル一覧を含まないこと。
  it("提供元を切り替えると選択先のendpointとモデルだけを持つ", () => {
    expect(providerDefaults("openai")).toEqual({
      endpoint: "https://api.openai.com/v1",
      models: [luna],
      model: luna
    })
    expect(providerDefaults("sync").models).toEqual([])
    expect(providerDefaults("xai").models).toEqual([])
  })
})
