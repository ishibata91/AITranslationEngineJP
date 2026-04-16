import { describe, expect, test, vi } from "vitest"

import type {
  MasterPersonaScreenState,
  MasterPersonaScreenViewModel
} from "@application/gateway-contract/master-persona"

import { MasterPersonaScreenController } from "./master-persona-screen-controller"

function createState(
  overrides: Partial<MasterPersonaScreenState> = {}
): MasterPersonaScreenState {
  return {
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
    ...overrides
  }
}

function createViewModel(state: MasterPersonaScreenState): MasterPersonaScreenViewModel {
  return {
    ...state,
    gatewayStatus: "接続準備済み",
    pluginOptions: [{ value: "", label: "すべてのプラグイン" }],
    totalPages: 1,
    pageStatusText: "1 - 0 件を表示しています。",
    selectionStatusText: "選択中のペルソナはありません。",
    listHeadline: "0 件から絞り込みます。",
    detailLockText: "更新と削除を行えます",
    detailStatusText: "一覧からペルソナを選ぶと、詳細を同じ画面で確認できます。",
    canStartPreview: false,
    canStartGeneration: false,
    canMutate: false,
    isRunActive: false,
    hasPreview: false,
    aiProviderLabel: "Gemini",
    promptTemplateDescription:
      "プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。",
    progressPercent: 0
  }
}

function createControllerHarness(initialState: MasterPersonaScreenState = createState()) {
  let state = initialState
  const listeners = new Set<(state: MasterPersonaScreenState) => void>()

  const store = {
    subscribe: vi.fn((listener: (nextState: MasterPersonaScreenState) => void) => {
      listeners.add(listener)
      return () => {
        listeners.delete(listener)
      }
    }),
    snapshot: vi.fn(() => state),
    update: vi.fn((mutator: (draft: MasterPersonaScreenState) => void) => {
      const draft = structuredClone(state)
      mutator(draft)
      state = draft
      for (const listener of listeners) {
        listener(state)
      }
    })
  }

  const presenter = {
    toViewModel: vi.fn((nextState: MasterPersonaScreenState) =>
      createViewModel(nextState)
    )
  }

  const useCase = {
    loadScreen: vi.fn(async () => {}),
    loadPage: vi.fn(async () => {}),
    selectEntry: vi.fn(async () => {}),
    loadDialogueList: vi.fn(async () => {}),
    previewGeneration: vi.fn(async () => {}),
    executeGeneration: vi.fn(async () => {}),
    loadRunStatus: vi.fn(async () => {}),
    interruptGeneration: vi.fn(async () => {}),
    cancelGeneration: vi.fn(async () => {}),
    saveAISettings: vi.fn(async () => {}),
    saveCurrentEntry: vi.fn(async () => {}),
    deleteCurrentEntry: vi.fn(async () => {}),
    closeDialogueModal: vi.fn(() => {}),
    setModalState: vi.fn(() => {})
  }

  const runtimePollingAdapter = {
    start: vi.fn(() => true),
    stop: vi.fn(() => {})
  }

  const controller = new MasterPersonaScreenController({
    isGatewayConnected: true,
    store,
    presenter,
    useCase,
    runtimePollingAdapter
  })

  return { controller, useCase, runtimePollingAdapter, getState: () => state }
}

describe("MasterPersonaScreenController", () => {
  test("mount は polling を開始して loadScreen を呼ぶ", async () => {
    const harness = createControllerHarness()

    await harness.controller.mount()

    expect(harness.runtimePollingAdapter.start).toHaveBeenCalledTimes(1)
    expect(harness.useCase.loadScreen).toHaveBeenCalledTimes(1)
  })

  test("openDialogueModal は useCase.loadDialogueList を呼ぶ", async () => {
    const harness = createControllerHarness(
      createState({
        selectedEntry: {
          identityKey: "key",
          targetPlugin: "FollowersPlus.esp",
          formId: "1",
          recordType: "NPC_",
          editorId: "edid",
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
    )

    await harness.controller.openDialogueModal()

    expect(harness.useCase.loadDialogueList).toHaveBeenCalledTimes(1)
  })

  test("handlePluginFilterChange は plugin filter を更新して loadPage を呼ぶ", () => {
    const harness = createControllerHarness(
      createState({
        pluginFilter: "OldPlugin.esp",
        page: 3,
        errorMessage: "stale"
      })
    )
    const select = document.createElement("select")
    const option = document.createElement("option")
    option.value = "FollowersPlus.esp"
    select.append(option)
    select.value = "FollowersPlus.esp"

    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", {
      value: select,
      configurable: true
    })

    harness.controller.handlePluginFilterChange(event)

    expect(harness.getState().pluginFilter).toBe("FollowersPlus.esp")
    expect(harness.getState().page).toBe(1)
    expect(harness.getState().errorMessage).toBe("")
    expect(harness.useCase.loadPage).toHaveBeenCalledTimes(1)
  })

  test("setAIProvider は canonical provider ID を state へ保持する", () => {
    const harness = createControllerHarness(createState())
    const select = document.createElement("select")
    const option = document.createElement("option")
    option.value = "lm_studio"
    select.append(option)
    select.value = "lm_studio"

    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", {
      value: select,
      configurable: true
    })

    harness.controller.setAIProvider(event)

    expect(harness.getState().aiSettings.provider).toBe("lm_studio")
    expect(harness.getState().aiSettingsMessage).toBe("")
  })

  test("stageJsonSelection は preview をクリアして file reference を保持する", () => {
    const harness = createControllerHarness(
      createState({
        preview: {
          fileName: "old.json",
          targetPlugin: "OldPlugin.esp",
          totalNpcCount: 3,
          generatableCount: 1,
          existingSkipCount: 1,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "生成可能"
        }
      })
    )
    const file = new File(["{}"], "sample.json", {
      type: "application/json"
    }) as File & { path: string }
    file.path = "/tmp/sample.json"

    harness.controller.stageJsonSelection(file)

    expect(harness.getState().selectedFileName).toBe("sample.json")
    expect(harness.getState().selectedFileReference).toBe("/tmp/sample.json")
    expect(harness.getState().preview).toBeNull()
    expect(harness.useCase.previewGeneration).toHaveBeenCalledTimes(1)
  })

  test("stageJsonSelection は file 未選択時に自動 preview を呼ばない", () => {
    const harness = createControllerHarness(createState())

    harness.controller.stageJsonSelection(null)

    expect(harness.getState().selectedFileName).toBe("未選択")
    expect(harness.getState().selectedFileReference).toBeNull()
    expect(harness.useCase.previewGeneration).not.toHaveBeenCalled()
  })
})
