import { afterEach, describe, expect, test, vi } from "vitest"

import { createMasterPersonaGateway } from "./master-persona.gateway"

type GoRecord = {
  wails: {
    AppController: {
      MasterPersonaLoadAISettings: ReturnType<typeof vi.fn>
      MasterPersonaGetRunStatus: ReturnType<typeof vi.fn>
    }
  }
}

const originalGo: unknown = Reflect.get(globalThis as object, "go")

function installGo(record: GoRecord): void {
  Object.defineProperty(globalThis, "go", {
    value: record,
    configurable: true,
    writable: true
  })
}

afterEach(() => {
  vi.restoreAllMocks()
  Object.defineProperty(globalThis, "go", {
    value: originalGo,
    configurable: true,
    writable: true
  })
})

describe("createMasterPersonaGateway", () => {
  test("no-arg bindings は request なしで Wails binding を呼ぶ", async () => {
    const loadMasterPersonaAISettings = vi.fn(() =>
      Promise.resolve({
        provider: "fake",
        model: "fake-master-persona",
        apiKey: ""
      })
    )
    const getMasterPersonaRunStatus = vi.fn(() =>
      Promise.resolve({
        runState: "入力待ち",
        targetPlugin: "",
        processedCount: 0,
        successCount: 0,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "",
        message: "入力ファイルを選ぶと状態を表示します。"
      })
    )

    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: loadMasterPersonaAISettings,
          MasterPersonaGetRunStatus: getMasterPersonaRunStatus
        }
      }
    })

    const gateway = createMasterPersonaGateway()

    await expect(gateway.loadMasterPersonaAISettings()).resolves.toEqual({
      provider: "fake",
      model: "fake-master-persona",
      apiKey: ""
    })
    await expect(gateway.getMasterPersonaRunStatus()).resolves.toEqual({
      runState: "入力待ち",
      targetPlugin: "",
      processedCount: 0,
      successCount: 0,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "",
      message: "入力ファイルを選ぶと状態を表示します。"
    })

    expect(loadMasterPersonaAISettings).toHaveBeenCalledTimes(1)
    expect(loadMasterPersonaAISettings).toHaveBeenCalledWith()
    expect(getMasterPersonaRunStatus).toHaveBeenCalledTimes(1)
    expect(getMasterPersonaRunStatus).toHaveBeenCalledWith()
  })
})
