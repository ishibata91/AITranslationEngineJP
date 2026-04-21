import { describe, expect, test, vi } from "vitest"

import type {
  MasterPersonaAISettings,
  MasterPersonaDetail,
  MasterPersonaScreenState
} from "@application/gateway-contract/master-persona"

import { MasterPersonaUseCase } from "./master-persona.usecase"

function createStore(initialState?: Partial<MasterPersonaScreenState>) {
  let state: MasterPersonaScreenState = {
    items: [],
    pluginGroups: [],
    selectedIdentityKey: null,
    selectedEntry: null,
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
  } as MasterPersonaScreenState

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
              dialogueCount: 0,
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
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          runLockReason: "更新と削除を行えます"
        } as MasterPersonaDetail
      })
    ),
    getMasterPersonaDialogueList: vi.fn(),
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
        candidateCount: 9,
        newlyAddableCount: 7,
        existingCount: 2,
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
        candidateCount: 9,
        newlyAddableCount: 7,
        existingCount: 2,
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

  test("loadDetail は selectedEntry を更新する", async () => {
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadDetail("FollowersPlus.esp:FE01A812:NPC_")

    expect(store.snapshot().selectedEntry?.identityKey).toBe("FollowersPlus.esp:FE01A812:NPC_")
    expect(gateway.getMasterPersonaDetail).toHaveBeenCalledWith({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
  })

  test("loadDetail 失敗時は selectedEntry を null にして errorMessage を保持する", async () => {
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      selectedEntry: {
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
      } as MasterPersonaDetail
    })
    const gateway = createGateway()
    gateway.getMasterPersonaDetail.mockRejectedValueOnce(new Error("entry not found"))
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadDetail("FollowersPlus.esp:FE01A812:NPC_")

    expect(store.snapshot().selectedEntry).toBeNull()
    expect(store.snapshot().errorMessage).toBe("entry not found")
  })

  test("selectEntry は selectedIdentityKey を更新して loadDetail を呼ぶ", async () => {
    const store = createStore({
      selectedIdentityKey: null
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.selectEntry("FollowersPlus.esp:FE01A812:NPC_")

    expect(store.snapshot().selectedIdentityKey).toBe("FollowersPlus.esp:FE01A812:NPC_")
    expect(gateway.getMasterPersonaDetail).toHaveBeenCalledWith({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
  })

  test("loadPage は items が空のとき selectedEntry を null にする", async () => {
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      selectedEntry: {
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
      } as MasterPersonaDetail
    })
    const gateway = createGateway()
    gateway.getMasterPersonaPage.mockResolvedValueOnce({
      page: {
        items: [],
        pluginGroups: [],
        totalCount: 0,
        page: 1,
        pageSize: 30,
        selectedIdentityKey: ""
      }
    })
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadPage()

    expect(store.snapshot().items).toHaveLength(0)
    expect(store.snapshot().selectedEntry).toBeNull()
  })

  test("loadPage 失敗時は items を空にして errorMessage を保持する", async () => {
    const store = createStore({})
    const gateway = createGateway()
    gateway.getMasterPersonaPage.mockRejectedValueOnce(new Error("network error"))
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadPage()

    expect(store.snapshot().items).toHaveLength(0)
    expect(store.snapshot().errorMessage).toBe("network error")
  })

  test("persona-read-detail-cutover: loadPage は plugin filter と keyword を refresh へ反映する", async () => {
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

  test("persona-read-detail-cutover: loadDetail は identity snapshot を selectedEntry へ反映する", async () => {
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_"
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    await useCase.loadDetail("FollowersPlus.esp:FE01A812:NPC_")

    const entry = store.snapshot().selectedEntry
    expect(entry?.identityKey).toBe("FollowersPlus.esp:FE01A812:NPC_")
    expect(entry?.formId).toBe("FE01A812")
    expect(entry?.editorId).toBe("FP_LysMaren")
    expect(entry?.displayName).toBe("Lys Maren")
    expect(entry?.voiceType).toBe("FemaleYoungEager")
    expect(entry?.runLockReason).toBe("更新と削除を行えます")
    expect(entry?.personaBody).toBe("body")
  })

  test("persona-read-detail-cutover: loadDetail は getMasterPersonaDialogueList を呼ばない", async () => {
    // Arrange
    const store = createStore()
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadDetail("FollowersPlus.esp:FE01A812:NPC_")

    // Assert
    expect(gateway.getMasterPersonaDialogueList).not.toHaveBeenCalled()
  })

  test("persona-read-detail-cutover: selectEntry は getMasterPersonaDialogueList を呼ばない", async () => {
    // Arrange
    const store = createStore()
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.selectEntry("FollowersPlus.esp:FE01A812:NPC_")

    // Assert
    expect(gateway.getMasterPersonaDialogueList).not.toHaveBeenCalled()
  })

  test("persona-read-detail-cutover: saveCurrentEntry は identity / snapshot fields を generic editable input として payload に含まない", async () => {
    // Arrange
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "edit",
      editForm: {
        formId: "FE01A812",
        editorId: "FP_LysMaren",
        displayName: "Lys Maren",
        race: "Nord",
        sex: "Female",
        voiceType: "FemaleYoungEager",
        className: "FPScoutClass",
        sourcePlugin: "FollowersPlus.esp",
        personaSummary: "edited summary text",
        personaBody: "edited persona body"
      }
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(gateway.updateMasterPersona).toHaveBeenCalledTimes(1)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    const entry = (gateway.updateMasterPersona.mock.calls as any)[0][0].entry as Record<string, unknown>
    expect(entry).not.toHaveProperty("formId")
    expect(entry).not.toHaveProperty("editorId")
    expect(entry).not.toHaveProperty("race")
    expect(entry).not.toHaveProperty("sex")
    expect(entry).not.toHaveProperty("voiceType")
    expect(entry).not.toHaveProperty("className")
    expect(entry).not.toHaveProperty("sourcePlugin")
  })

  test("persona-read-detail-cutover: saveCurrentEntry は read-only selectedEntry の identityKey を update request に含める", async () => {
    // Arrange
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "edit",
      editForm: {
        formId: "",
        editorId: "",
        displayName: "Lys Maren",
        voiceType: "",
        className: "",
        sourcePlugin: "",
        personaSummary: "ペルソナ概要",
        personaBody: "ペルソナ本文"
      }
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.saveCurrentEntry()

    // Assert: identity linkage from read-only selected entry flows through to update request
    expect(gateway.updateMasterPersona).toHaveBeenCalledTimes(1)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    const request = (gateway.updateMasterPersona.mock.calls as any)[0][0] as Record<string, unknown>
    expect(request.identityKey).toBe("FollowersPlus.esp:FE01A812:NPC_")
  })

  test("persona-ai-settings-restart-cutover: loadScreen は aiSettings を backend から復元する", async () => {
    // Arrange
    const store = createStore({
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "" }
    })
    const gateway = createGateway()
    gateway.loadMasterPersonaAISettings.mockResolvedValueOnce({
      provider: "lm_studio",
      model: "restored-model",
      apiKey: "restored-key"
    })
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadScreen()

    // Assert
    expect(store.snapshot().aiSettings.provider).toBe("lm_studio")
    expect(store.snapshot().aiSettings.model).toBe("restored-model")
    expect(gateway.loadMasterPersonaAISettings).toHaveBeenCalledTimes(1)
  })

  test("persona-ai-settings-restart-cutover: loadScreen は JSON ファイル選択を backend から復元しない", async () => {
    // Arrange: fresh store with default (未選択) JSON file selection
    const store = createStore({
      selectedFileName: "未選択",
      selectedFileReference: null
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadScreen()

    // Assert: AI settings are restored from backend but JSON file selection is not touched
    expect(store.snapshot().selectedFileName).toBe("未選択")
    expect(store.snapshot().selectedFileReference).toBeNull()
    expect(gateway.loadMasterPersonaAISettings).toHaveBeenCalledTimes(1)
  })

  test("persona-ai-settings-restart-cutover: loadScreen は runStatus を backend から復元しない", async () => {
    // Arrange: store starts with default "入力待ち" run state
    const store = createStore({
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
      }
    })
    const gateway = createGateway()
    // eslint-disable-next-line local/no-commented-out-code
    // Backend returns "完了" run state (as if a previous run completed before restart)
    // Contract: loadScreen must NOT restore this "完了" state — runStatus stays input-waiting/idle
    gateway.getMasterPersonaRunStatus.mockResolvedValueOnce({
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
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadScreen()

    // Assert: runStatus must not be restored from backend on restart — stays input-waiting/idle
    expect(store.snapshot().runStatus.runState).toBe("入力待ち")
    expect(store.snapshot().runStatus.processedCount).toBe(0)
    expect(store.snapshot().runStatus.successCount).toBe(0)
  })

  test("persona-json-preview-cutover: previewGeneration は candidateCount/newlyAddableCount/existingCount だけを store に保持する", async () => {
    // Arrange
    const store = createStore({
      selectedFileReference: "/tmp/sample.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "" }
    })
    const gateway = createGateway()
    gateway.previewMasterPersonaGeneration.mockResolvedValueOnce({
      fileName: "sample.json",
      targetPlugin: "FollowersPlus.esp",
      candidateCount: 840,
      newlyAddableCount: 228,
      existingCount: 612,
      status: "生成可能"
    } as never)
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.previewGeneration()

    // Assert
    const preview = store.snapshot().preview as unknown as Record<string, unknown>
    expect(preview).not.toBeNull()
    expect(preview.candidateCount).toBe(840)
    expect(preview.newlyAddableCount).toBe(228)
    expect(preview.existingCount).toBe(612)
    expect(Object.keys(preview)).not.toContain("zeroDialogueSkipCount")
    expect(Object.keys(preview)).not.toContain("genericNpcCount")
  })

  test("persona-generation-cutover: executeGeneration 失敗時は errorMessage を設定して page を再取得しない", async () => {
    // Arrange
    const store = createStore({
      selectedFileReference: "/tmp/sample.json",
      aiSettings: { provider: "gemini", model: "gemini-2.5-pro", apiKey: "" }
    })
    const gateway = createGateway()
    gateway.executeMasterPersonaGeneration.mockRejectedValueOnce(
      new Error("AI API error: credentials invalid")
    )
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.executeGeneration()

    // Assert
    expect(store.snapshot().errorMessage).toBe("AI API error: credentials invalid")
    expect(gateway.getMasterPersonaPage).not.toHaveBeenCalled()
  })

  test("persona-generation-cutover: loadRunStatus が 生成中 から 完了 へ遷移したとき page を再取得する", async () => {
    // Arrange
    const store = createStore({
      runStatus: {
        runState: "生成中",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 3,
        successCount: 3,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "Lys Maren",
        message: "ペルソナを作成中"
      }
    })
    const gateway = createGateway()
    gateway.getMasterPersonaRunStatus.mockResolvedValueOnce({
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
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadRunStatus()

    // Assert
    expect(store.snapshot().runStatus.runState).toBe("完了")
    expect(gateway.getMasterPersonaPage).toHaveBeenCalled()
  })

  test("persona-generation-cutover: loadRunStatus が 生成中 のままのとき page を再取得しない", async () => {
    // Arrange
    const store = createStore({
      runStatus: {
        runState: "生成中",
        targetPlugin: "FollowersPlus.esp",
        processedCount: 1,
        successCount: 1,
        existingSkipCount: 0,
        zeroDialogueSkipCount: 0,
        genericNpcCount: 0,
        currentActorLabel: "Lys Maren",
        message: "ペルソナを作成中"
      }
    })
    const gateway = createGateway()
    gateway.getMasterPersonaRunStatus.mockResolvedValueOnce({
      runState: "生成中",
      targetPlugin: "FollowersPlus.esp",
      processedCount: 4,
      successCount: 4,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "Another NPC",
      message: "ペルソナを作成中"
    })
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.loadRunStatus()

    // Assert
    expect(gateway.getMasterPersonaPage).not.toHaveBeenCalled()
  })

  test("persona-edit-delete-cutover: deleteCurrentEntry は deleteMasterPersona を呼び modalState を null にする", async () => {
    // Arrange
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "delete"
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(gateway.deleteMasterPersona).toHaveBeenCalledTimes(1)
    expect(gateway.deleteMasterPersona).toHaveBeenCalledWith({
      identityKey: "FollowersPlus.esp:FE01A812:NPC_",
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      refresh: expect.objectContaining({ page: 1, pageSize: 30 })
    })
    expect(store.snapshot().modalState).toBeNull()
    expect(store.snapshot().errorMessage).toBe("")
  })

  test("persona-edit-delete-cutover: deleteCurrentEntry 失敗時は errorMessage を設定して modalState を保持する", async () => {
    // Arrange
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "delete"
    })
    const gateway = createGateway()
    gateway.deleteMasterPersona.mockRejectedValueOnce(new Error("delete failed"))
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("delete failed")
    expect(store.snapshot().modalState).toBe("delete")
  })

  test("persona-edit-delete-cutover: deleteCurrentEntry は modalState が delete でないとき gateway を呼ばない", async () => {
    // Arrange
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "edit"
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(gateway.deleteMasterPersona).not.toHaveBeenCalled()
  })

  test("persona-edit-delete-cutover: saveCurrentEntry の entry は personaSummary/speechStyle/personaBody を持ち displayName を含まない", async () => {
    // Arrange: editForm に post-cutover の shape をセット (RED: 現在は displayName が payload に含まれる)
    const store = createStore({
      selectedIdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
      modalState: "edit",
      editForm: {
        displayName: "Lys Maren",
        personaBody: "ペルソナ本文",
        // post-cutover で使われるフィールド (型互換のため as unknown as でキャスト)
        personaSummary: "乾いた率直さで応じる。",
        speechStyle: "短く本音を置く。"
      } as unknown as MasterPersonaScreenState["editForm"]
    })
    const gateway = createGateway()
    const useCase = new MasterPersonaUseCase(gateway, store)

    // Act
    await useCase.saveCurrentEntry()

    // Assert: gateway が呼ばれ entry は personaSummary/speechStyle/personaBody だけを持つ (RED)
    expect(gateway.updateMasterPersona).toHaveBeenCalledTimes(1)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    const entry = (gateway.updateMasterPersona.mock.calls as any)[0][0].entry as Record<string, unknown>
    expect(entry).not.toHaveProperty("displayName")
    expect(entry).toHaveProperty("personaSummary")
    expect(entry).toHaveProperty("speechStyle")
    expect(entry).toHaveProperty("personaBody")
  })
})
