import { describe, expect, test, vi } from "vitest"

import type { TranslationInputGatewayContract } from "@application/gateway-contract/translation-input"

import { createTranslationInputScreenControllerFactory } from "./translation-input-screen-controller-factory"

function createGateway(): TranslationInputGatewayContract {
  return {
    importTranslationInput: vi.fn(() => Promise.reject(new Error("not used"))),
    rebuildTranslationInputCache: vi.fn(() => Promise.reject(new Error("not used")))
  }
}

describe("createTranslationInputScreenControllerFactory", () => {
  test("同一 app session 中は同じ controller instance を返す", () => {
    const factory = createTranslationInputScreenControllerFactory(createGateway())

    const firstController = factory()
    const secondController = factory()

    expect(secondController).toBe(firstController)
  })
})