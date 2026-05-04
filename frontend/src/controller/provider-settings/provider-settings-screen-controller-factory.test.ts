import { describe, expect, test } from "vitest"

import { createProviderSettingsScreenControllerFactory } from "./provider-settings-screen-controller-factory"

describe("createProviderSettingsScreenControllerFactory", () => {
  test("同一 app session 中は同じ controller instance を返す", () => {
    const factory = createProviderSettingsScreenControllerFactory(null)

    const firstController = factory()
    const secondController = factory()

    expect(secondController).toBe(firstController)
  })

  test("gateway null は未接続 status を返す", () => {
    const controller = createProviderSettingsScreenControllerFactory(null)()

    expect(controller.getViewModel().gatewayStatus).toBe("未接続")
  })
})
