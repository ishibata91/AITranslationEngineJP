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
        provider: "gemini",
        model: "gemini-2.5-pro",
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
      provider: "gemini",
      model: "gemini-2.5-pro",
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

  test("persona-read-detail-cutover: getMasterPersonaDetail は identityKey だけを含む request で Wails binding を呼ぶ", async () => {
    // Arrange
    const getMasterPersonaGetDetail = vi.fn(() =>
      Promise.resolve({
        entry: {
          identityKey: "FollowersPlus.esp:FE01A812:NPC_",
          targetPlugin: "FollowersPlus.esp",
          formId: "FE01A812",
          recordType: "NPC_",
          editorId: "FP_LysMaren",
          displayName: "Lys Maren",
          voiceType: "FemaleYoungEager",
          className: "FPScoutClass",
          sourcePlugin: "FollowersPlus.esp",
          personaSummary: "summary",
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          runLockReason: "更新と削除を行えます"
        }
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaGetDetail: getMasterPersonaGetDetail
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    await gateway.getMasterPersonaDetail({ identityKey: "FollowersPlus.esp:FE01A812:NPC_" })

    // Assert
    expect(getMasterPersonaGetDetail).toHaveBeenCalledTimes(1)
    expect(getMasterPersonaGetDetail).toHaveBeenCalledWith({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
  })

  test("persona-read-detail-cutover: getMasterPersonaPage は refresh を含む request で Wails binding を呼ぶ", async () => {
    // Arrange
    const getMasterPersonaGetPage = vi.fn(() =>
      Promise.resolve({
        page: {
          items: [],
          pluginGroups: [],
          totalCount: 0,
          page: 1,
          pageSize: 30,
          selectedIdentityKey: ""
        }
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaGetPage: getMasterPersonaGetPage
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    await gateway.getMasterPersonaPage({
      refresh: {
        keyword: "lys",
        pluginFilter: "FollowersPlus.esp",
        page: 1,
        pageSize: 30
      }
    })

    // Assert
    expect(getMasterPersonaGetPage).toHaveBeenCalledTimes(1)
    expect(getMasterPersonaGetPage).toHaveBeenCalledWith({
      refresh: {
        keyword: "lys",
        pluginFilter: "FollowersPlus.esp",
        page: 1,
        pageSize: 30
      }
    })
  })

  test("persona-read-detail-cutover: gateway は getMasterPersonaDialogueList binding を公開しない", () => {
    // Arrange
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn()
        }
      }
    } as unknown as GoRecord)

    // Act
    const gateway = createMasterPersonaGateway()

    // Assert
    expect("getMasterPersonaDialogueList" in gateway).toBe(false)
  })

  test("persona-ai-settings-restart-cutover: loadMasterPersonaAISettings は provider と model を backend から返す", async () => {
    // Arrange
    const loadAISettings = vi.fn(() =>
      Promise.resolve({
        provider: "lm_studio",
        model: "restart-cutover-model",
        apiKey: "restart-cutover-key"
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: loadAISettings,
          MasterPersonaGetRunStatus: vi.fn()
        }
      }
    })
    const gateway = createMasterPersonaGateway()

    // Act
    const result = await gateway.loadMasterPersonaAISettings()

    // Assert
    expect(result.provider).toBe("lm_studio")
    expect(result.model).toBe("restart-cutover-model")
    expect(result.apiKey).toBe("restart-cutover-key")
    expect(loadAISettings).toHaveBeenCalledTimes(1)
    expect(loadAISettings).toHaveBeenCalledWith()
  })

  test("persona-ai-settings-restart-cutover: getMasterPersonaRunStatus は runState を backend から返す", async () => {
    // Arrange
    const getRunStatus = vi.fn(() =>
      Promise.resolve({
        runState: "完了",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 1,
        successCount: 1,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "",
        message: "生成完了"
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: getRunStatus
        }
      }
    })
    const gateway = createMasterPersonaGateway()

    // Act
    const result = await gateway.getMasterPersonaRunStatus()

    // Assert
    expect(result.runState).toBe("完了")
    expect(getRunStatus).toHaveBeenCalledTimes(1)
    expect(getRunStatus).toHaveBeenCalledWith()
  })

  test("persona-json-preview-cutover: previewMasterPersonaGeneration は candidateCount/newlyAddableCount/existingCount を Wails から受け取る", async () => {
    // Arrange
    const previewGeneration = vi.fn(() =>
      Promise.resolve({
        fileName: "FollowersPlus.json",
        targetPlugin: "FollowersPlus.esp",
        candidateCount: 840,
        newlyAddableCount: 228,
        existingCount: 612,
        status: "生成可能"
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaPreviewGeneration: previewGeneration
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    const result = await gateway.previewMasterPersonaGeneration({
      filePath: "/tmp/FollowersPlus.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "" }
    })

    // Assert
    expect(previewGeneration).toHaveBeenCalledTimes(1)
    expect((result as unknown as Record<string, unknown>).candidateCount).toBe(840)
    expect((result as unknown as Record<string, unknown>).newlyAddableCount).toBe(228)
    expect((result as unknown as Record<string, unknown>).existingCount).toBe(612)
    expect(Object.keys(result as unknown as object)).not.toContain("zeroDialogueSkipCount")
    expect(Object.keys(result as unknown as object)).not.toContain("genericNpcCount")
  })

  test("persona-generation-cutover: executeMasterPersonaGeneration は filePath と aiSettings を MasterPersonaExecuteGeneration binding へ渡し existingSkipCount を含む runStatus を返す", async () => {
    // Arrange
    const executeGeneration = vi.fn(() =>
      Promise.resolve({
        runState: "完了",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 7,
        successCount: 7,
        existingSkipCount: 2,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "",
        message: "完了"
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaExecuteGeneration: executeGeneration
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    const result = await gateway.executeMasterPersonaGeneration({
      filePath: "/tmp/FollowersPlus.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "test-key" }
    })

    // Assert
    expect(executeGeneration).toHaveBeenCalledTimes(1)
    expect(executeGeneration).toHaveBeenCalledWith({
      filePath: "/tmp/FollowersPlus.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "test-key" }
    })
    // generation never overwrites existing: existingSkipCount を runStatus で確認できる
    expect(result.existingSkipCount).toBe(2)
    expect(result.runState).toBe("完了")
  })

  test("persona-generation-cutover: executeMasterPersonaGeneration の request には zeroDialogueSkipCount / genericNpcCount は含まれない", async () => {
    // Arrange
    const executeGeneration = vi.fn(() =>
      Promise.resolve({
        runState: "完了",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 5,
        successCount: 5,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "",
        message: "完了"
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaExecuteGeneration: executeGeneration
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    await gateway.executeMasterPersonaGeneration({
      filePath: "/tmp/sample.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "" }
    })

    // Assert
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    const callArg = (executeGeneration.mock.calls as any)[0][0] as Record<string, unknown>
    expect(Object.keys(callArg)).not.toContain("zeroDialogueSkipCount")
    expect(Object.keys(callArg)).not.toContain("genericNpcCount")
  })

  test("persona-edit-delete-cutover: updateMasterPersona は MasterPersonaUpdate binding へ request を転送する", async () => {
    // Arrange
    const masterPersonaUpdate = vi.fn(() =>
      Promise.resolve({
        page: {
          items: [],
          pluginGroups: [],
          totalCount: 0,
          page: 1,
          pageSize: 30,
          selectedIdentityKey: ""
        }
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaUpdate: masterPersonaUpdate
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    const request = {
      identityKey: "FollowersPlus.esp:FE01A812:NPC_",
      entry: { displayName: "Lys Maren", personaBody: "本文" },
      refresh: { keyword: "", pluginFilter: "", page: 1, pageSize: 30 }
    }
    await gateway.updateMasterPersona(request)

    // Assert
    expect(masterPersonaUpdate).toHaveBeenCalledTimes(1)
    expect(masterPersonaUpdate).toHaveBeenCalledWith(request)
  })

  test("persona-edit-delete-cutover: deleteMasterPersona は MasterPersonaDelete binding へ identityKey と refresh を渡す", async () => {
    // Arrange
    const masterPersonaDelete = vi.fn(() =>
      Promise.resolve({
        page: {
          items: [],
          pluginGroups: [],
          totalCount: 0,
          page: 1,
          pageSize: 30,
          selectedIdentityKey: ""
        }
      })
    )
    installGo({
      wails: {
        AppController: {
          MasterPersonaLoadAISettings: vi.fn(),
          MasterPersonaGetRunStatus: vi.fn(),
          MasterPersonaDelete: masterPersonaDelete
        }
      }
    } as unknown as GoRecord)
    const gateway = createMasterPersonaGateway()

    // Act
    await gateway.deleteMasterPersona({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_",
      refresh: { keyword: "", pluginFilter: "", page: 1, pageSize: 30 }
    })

    // Assert
    expect(masterPersonaDelete).toHaveBeenCalledTimes(1)
    expect(masterPersonaDelete).toHaveBeenCalledWith({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_",
      refresh: { keyword: "", pluginFilter: "", page: 1, pageSize: 30 }
    })
  })
})
