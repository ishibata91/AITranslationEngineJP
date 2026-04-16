import { describe, expect, test, vi } from "vitest"

import type {
  MasterPersonaAISettings,
  MasterPersonaScreenState
} from "@application/gateway-contract/master-persona"

import { MasterPersonaUseCase } from "./master-persona.usecase"

function createStore(initialState?: Partial<MasterPersonaScreenState>) {
  let state: MasterPersonaScreenState = {
    items: [],
    pluginGroups: [],
    selectedIdentityKey: null,
    selectedEntry: null,
    dialogueModalOpen: false,
    dialogues: [],
    keyword: "",
    pluginFilter: "",
    page: 1,
    pageSize: 30,
    totalCount: 0,
    errorMessage: "",
    aiSettings: {
      provider: "gemini",
      model: "gemini-2.5-pro",
      apiKey: ""
    },
    aiSettingsMessage: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    preview: null,
    runStatus: {
      runState: "入力待ち",
      targetPlugin: "",
      processedCount: 0,
      successCount: 0,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "",
      message: "入力ファイルを選ぶと状態を表示します。"
    },
    modalState: null,
    editForm: {
      formId: "",
      editorId: "",
      displayName: "",
      voiceType: "",
      className: "",
      sourcePlugin: "",
      personaBody: ""
    },
    ...initialState
  }

  return {
    snapshot: () => structuredClone(state),
    update(mutator: (draft: MasterPersonaScreenState) => void) {
      const draft = structuredClone(state)
      mutator(draft)
      state = draft
    }
  }
}

function createGateway() {
  return {
    getMasterPersonaPage: vi.fn(() =>
      Promise.resolve({
        page: {
          items: [
            {
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
              dialogueCount: 44,
              updatedAt: "2026-04-15T09:42:00Z"
            }
          ],
          pluginGroups: [{ targetPlugin: "FollowersPlus.esp", count: 1 }],
          totalCount: 1,
          page: 1,
          pageSize: 30,
          selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_"
        }
      })
    ),
    getMasterPersonaDetail: vi.fn(() =>
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
          dialogueCount: 44,
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          generationSourceJson: "sample.json",
          baselineApplied: false,
          runLockReason: "更新と削除を行えます"
        }
      })
    ),
    getMasterPersonaDialogueList: vi.fn(() =>
      Promise.resolve({
        identityKey: "FollowersPlus.esp:FE01A812:NPC_",
        dialogueCount: 2,
        dialogues: [
          { index: 1, text: "line1" },
          { index: 2, text: "line2" }
        ]
      })
    ),
    loadMasterPersonaAISettings: vi.fn(() =>
      Promise.resolve({
        provider: "gemini",
        model: "gemini-2.5-pro",
        apiKey: ""
      })
    ),
    saveMasterPersonaAISettings: vi.fn((request: MasterPersonaAISettings) =>
      Promise.resolve({ ...request })
    ),
    previewMasterPersonaGeneration: vi.fn(() =>
      Promise.resolve({
        fileName: "sample.json",
        targetPlugin: "FollowersPlus.esp",
        totalNpcCount: 10,
        generatableCount: 7,
        existingSkipCount: 2,
        zeroDialogueSkipCount: 1,
        genericNpcCount: 0,
        status: "生成可能"
      })
    ),
    executeMasterPersonaGeneration: vi.fn(() =>
      Promise.resolve({
        runState: "完了",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 7,
        successCount: 7,
        existingSkipCount: 2,
        zeroDialogueSkipCount: 1,
        genericNpcCount: 0,
        currentActorLabel: "",
        message: "完了"
      })
    ),
    getMasterPersonaRunStatus: vi.fn(() =>
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
    ),
    interruptMasterPersonaGeneration: vi.fn(() =>
      Promise.resolve({
        runState: "中断済み",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 2,
        successCount: 1,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "Lys Maren",
        message: "生成を中断しました"
      })
    ),
    cancelMasterPersonaGeneration: vi.fn(() =>
      Promise.resolve({
        runState: "中止済み",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 2,
        successCount: 1,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "Lys Maren",
        message: "生成を停止しました"
      })
    ),
    updateMasterPersona: vi.fn(() =>
      Promise.resolve({
        page: {
          items: [],
          pluginGroups: [],
          totalCount: 0,
          page: 1,
          pageSize: 30
        }
      })
    ),
    deleteMasterPersona: vi.fn(() =>
      Promise.resolve({
        page: {
          items: [],
          pluginGroups: [],
          totalCount: 0,
          page: 1,
          pageSize: 30
        }
      })
    )
  }
}

describe("MasterPersonaUseCase", () => {
  test("loadScreen は transport seam provider を page-local AI settings へ反映する", async () => {
    const store = createStore()
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadScreen()

    expect(store.snapshot().aiSettings.provider).toBe("gemini")
    expect(gateway.loadMasterPersonaAISettings).toHaveBeenCalledTimes(1)
  })

  test("loadPage は plugin filter だけを refresh request に含める", async () => {
    const store = createStore({
      keyword: "lys",
      pluginFilter: "FollowersPlus.esp",
      page: 2
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadPage()

    expect(gateway.getMasterPersonaPage).toHaveBeenCalledWith({
      refresh: {
        keyword: "lys",
        pluginFilter: "FollowersPlus.esp",
        page: 2,
        pageSize: 30
      },
      preferredIdentityKey: undefined
    })
  })

  test("loadDialogueList は closed-by-default modal を開いて dialogue を反映する", async () => {
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadDialogueList()

    expect(store.snapshot().dialogueModalOpen).toBe(true)
    expect(store.snapshot().dialogues).toHaveLength(2)
  })

  test("saveAISettings は prompt template を送らず page-local settings だけを保存する", async () => {
    const store = createStore({
      aiSettings: {
        provider: "gemini",
        model: "persona-only-model",
        apiKey: ""
      }
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.saveAISettings()

    expect(gateway.saveMasterPersonaAISettings).toHaveBeenCalledWith({
      provider: "gemini",
      model: "persona-only-model",
      apiKey: ""
    })
    expect(store.snapshot().aiSettingsMessage).toBe(
      "この画面で使う設定を保存しました。"
    )
  })

  test("executeGeneration は transport seam provider の完了結果を受け取り page を再取得する", async () => {
    const store = createStore({
      selectedFileReference: "/tmp/sample.json",
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      aiSettings: {
        provider: "gemini",
        model: "gemini-2.5-pro",
        apiKey: ""
      },
      preview: {
        fileName: "sample.json",
        targetPlugin: "FollowersPlus.esp",
        totalNpcCount: 10,
        generatableCount: 7,
        existingSkipCount: 2,
        zeroDialogueSkipCount: 1,
        genericNpcCount: 0,
        status: "生成可能"
      }
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.executeGeneration()

    expect(gateway.executeMasterPersonaGeneration).toHaveBeenCalledWith({
      filePath: "/tmp/sample.json",
      aiSettings: {
        provider: "gemini",
        model: "gemini-2.5-pro",
        apiKey: ""
      }
    })
    expect(gateway.getMasterPersonaPage).toHaveBeenCalled()
    expect(store.snapshot().runStatus.runState).toBe("完了")
  })

  test("previewGeneration は AI 設定未完了でも集計 preview を保持する", async () => {
    const store = createStore({
      selectedFileReference: "/tmp/sample.json",
      aiSettings: {
        provider: "gemini",
        model: "gemini-2.5-pro",
        apiKey: ""
      }
    })
    const gateway = createGateway()
    gateway.previewMasterPersonaGeneration.mockResolvedValueOnce({
      fileName: "sample.json",
      targetPlugin: "FollowersPlus.esp",
      totalNpcCount: 10,
      generatableCount: 7,
      existingSkipCount: 2,
      zeroDialogueSkipCount: 1,
      genericNpcCount: 0,
      status: "設定未完了"
    })
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.previewGeneration()

    expect(store.snapshot().preview).toEqual({
      fileName: "sample.json",
      targetPlugin: "FollowersPlus.esp",
      totalNpcCount: 10,
      generatableCount: 7,
      existingSkipCount: 2,
      zeroDialogueSkipCount: 1,
      genericNpcCount: 0,
      status: "設定未完了"
    })
  })

  test("previewGeneration 失敗時は preview を消して error message を保持する", async () => {
    const store = createStore({
      selectedFileReference: "/tmp/broken.json",
      preview: {
        fileName: "sample.json",
        targetPlugin: "FollowersPlus.esp",
        totalNpcCount: 10,
        generatableCount: 7,
        existingSkipCount: 2,
        zeroDialogueSkipCount: 1,
        genericNpcCount: 0,
        status: "生成可能"
      }
    })
    const gateway = createGateway()
    gateway.previewMasterPersonaGeneration.mockRejectedValueOnce(
      new Error("parse extractData json: invalid")
    )
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.previewGeneration()

    expect(store.snapshot().preview).toBeNull()
    expect(store.snapshot().errorMessage).toBe("parse extractData json: invalid")
  })
})
