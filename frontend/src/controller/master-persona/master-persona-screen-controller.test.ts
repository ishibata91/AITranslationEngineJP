import { describe, expect, test, vi } from "vitest"

/* eslint-disable import-x/no-duplicates */
import type {
  MasterPersonaScreenState,
  MasterPersonaScreenViewModel,
  MasterPersonaGatewayContract,
  MasterPersonaUpdateRequest
} from "@application/gateway-contract/master-persona"
import type * as GatewayContractPublic from "@application/gateway-contract/master-persona"
import { buildMasterPersonaUpdateInput } from "@application/gateway-contract/master-persona"
/* eslint-enable import-x/no-duplicates */

import { MasterPersonaScreenController } from "./master-persona-screen-controller"

// ---------------------------------------------------------------------------
// 型レベル assertion: persona-read-detail-cutover
// ---------------------------------------------------------------------------

// MasterPersonaGatewayContract に dialogue 取得メソッドがないこと
type _AssertNoContractMethod<K extends string> = K extends keyof MasterPersonaGatewayContract ? never : true
const _noGetDialogueListMethod: _AssertNoContractMethod<"getMasterPersonaDialogueList"> = true
void _noGetDialogueListMethod

// MasterPersonaDialogueLine は public seam (barrel) から export されていないこと
// @ts-expect-error MasterPersonaDialogueLine は public seam から export されていないこと
// eslint-disable-next-line @typescript-eslint/no-unused-vars
type _noDialogueLine = GatewayContractPublic.MasterPersonaDialogueLine

// MasterPersonaUpdateInput に identity/snapshot optional fields がないこと
// これらの key が残っている間は compile error になる (persona-read-detail-cutover red test)
type _AssertNoUpdateInputKey<K extends string> = K extends keyof MasterPersonaUpdateRequest["entry"] ? never : true
const _noFormIdInInput: _AssertNoUpdateInputKey<"formId"> = true
void _noFormIdInInput
const _noEditorIdInInput: _AssertNoUpdateInputKey<"editorId"> = true
void _noEditorIdInInput
const _noRaceInInput: _AssertNoUpdateInputKey<"race"> = true
void _noRaceInInput
const _noSexInInput: _AssertNoUpdateInputKey<"sex"> = true
void _noSexInInput
const _noVoiceTypeInInput: _AssertNoUpdateInputKey<"voiceType"> = true
void _noVoiceTypeInInput
const _noClassNameInInput: _AssertNoUpdateInputKey<"className"> = true
void _noClassNameInInput
const _noSourcePluginInInput: _AssertNoUpdateInputKey<"sourcePlugin"> = true
void _noSourcePluginInInput

// ---------------------------------------------------------------------------
// 型レベル assertion: persona-json-preview-cutover
// ---------------------------------------------------------------------------

// MasterPersonaPreviewResult に zeroDialogueSkipCount/genericNpcCount がないこと
// これらの key が残っている間は compile error になる (persona-json-preview-cutover red test)
type _AssertNoPreviewResultKey<K extends string> = K extends keyof GatewayContractPublic.MasterPersonaPreviewResult ? never : true
const _noZeroDialogueSkipCount: _AssertNoPreviewResultKey<"zeroDialogueSkipCount"> = true
void _noZeroDialogueSkipCount
const _noGenericNpcCount: _AssertNoPreviewResultKey<"genericNpcCount"> = true
void _noGenericNpcCount
const _noTotalNpcCount: _AssertNoPreviewResultKey<"totalNpcCount"> = true
void _noTotalNpcCount
const _noGeneratableCount: _AssertNoPreviewResultKey<"generatableCount"> = true
void _noGeneratableCount
const _noExistingSkipCountInPreview: _AssertNoPreviewResultKey<"existingSkipCount"> = true
void _noExistingSkipCountInPreview

// MasterPersonaPreviewResult に candidateCount/newlyAddableCount/existingCount があること
// これらの key がない間は compile error になる (persona-json-preview-cutover red test)
type _AssertHasPreviewResultKey<K extends string> = K extends keyof GatewayContractPublic.MasterPersonaPreviewResult ? true : never
const _hasCandidateCount: _AssertHasPreviewResultKey<"candidateCount"> = true
void _hasCandidateCount
const _hasNewlyAddableCount: _AssertHasPreviewResultKey<"newlyAddableCount"> = true
void _hasNewlyAddableCount
const _hasExistingCount: _AssertHasPreviewResultKey<"existingCount"> = true
void _hasExistingCount

// MasterPersonaPreviewStateEntry (public/state seam) に legacy fields がないこと
// totalNpcCount/generatableCount/existingSkipCount/zeroDialogueSkipCount/genericNpcCount が
// optional key として残っている間は compile error になる (persona-json-preview-cutover red test)
type _AssertNoPreviewStateEntryKey<K extends string> = K extends keyof GatewayContractPublic.MasterPersonaPreviewStateEntry ? never : true
const _noTotalNpcCountInState: _AssertNoPreviewStateEntryKey<"totalNpcCount"> = true
void _noTotalNpcCountInState
const _noGeneratableCountInState: _AssertNoPreviewStateEntryKey<"generatableCount"> = true
void _noGeneratableCountInState
const _noExistingSkipCountInState: _AssertNoPreviewStateEntryKey<"existingSkipCount"> = true
void _noExistingSkipCountInState
const _noZeroDialogueSkipCountInState: _AssertNoPreviewStateEntryKey<"zeroDialogueSkipCount"> = true
void _noZeroDialogueSkipCountInState
const _noGenericNpcCountInState: _AssertNoPreviewStateEntryKey<"genericNpcCount"> = true
void _noGenericNpcCountInState

// ---------------------------------------------------------------------------
// 型レベル RED assertion: persona-generation-cutover
// ---------------------------------------------------------------------------
// MasterPersonaRunStatus に zeroDialogueSkipCount/genericNpcCount がまだ残っている
// これらの key が削除されるまで compile error になる (persona-generation-cutover RED)
type _AssertNoRunStatusKey<K extends string> = K extends keyof GatewayContractPublic.MasterPersonaRunStatus ? never : true
const _noZeroDialogueSkipCountInRunStatus: _AssertNoRunStatusKey<"zeroDialogueSkipCount"> = true
void _noZeroDialogueSkipCountInRunStatus
const _noGenericNpcCountInRunStatus: _AssertNoRunStatusKey<"genericNpcCount"> = true
void _noGenericNpcCountInRunStatus

// ---------------------------------------------------------------------------
// 型レベル GREEN assertion: persona-edit-delete-cutover
// ---------------------------------------------------------------------------
// MasterPersonaModalState に 'create' が含まれないこと (persona-edit-delete-cutover)
// 'create' が追加された場合は compile error になる (RED に反転する)
type _AssertNoCreateInModalState = "create" extends GatewayContractPublic.MasterPersonaModalState ? never : true
const _noCreateInModalState: _AssertNoCreateInModalState = true
void _noCreateInModalState

function createState(
  overrides: Partial<MasterPersonaScreenState> = {}
): MasterPersonaScreenState {
  return {
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
    ...overrides
  } as MasterPersonaScreenState
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
    previewGeneration: vi.fn(async () => {}),
    executeGeneration: vi.fn(async () => {}),
    loadRunStatus: vi.fn(async () => {}),
    interruptGeneration: vi.fn(async () => {}),
    cancelGeneration: vi.fn(async () => {}),
    saveAISettings: vi.fn(async () => {}),
    saveCurrentEntry: vi.fn(async () => {}),
    deleteCurrentEntry: vi.fn(async () => {}),
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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
    useCase: useCase as any,
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
          candidateCount: 2,
          newlyAddableCount: 1,
          existingCount: 1,
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

  test("selectRow は useCase.selectEntry を指定した identity key で呼ぶ", async () => {
    const harness = createControllerHarness()

    await harness.controller.selectRow("FollowersPlus.esp:FE01A812:NPC_")

    expect(harness.useCase.selectEntry).toHaveBeenCalledWith("FollowersPlus.esp:FE01A812:NPC_")
    expect(harness.useCase.selectEntry).toHaveBeenCalledTimes(1)
  })

  test("handleSearchInput は keyword を state に保持して loadPage を呼ぶ", () => {
    const harness = createControllerHarness(createState({ keyword: "", page: 3 }))
    const input = document.createElement("input")
    input.value = "lys"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", {
      value: input,
      configurable: true
    })

    harness.controller.handleSearchInput(event)

    expect(harness.getState().keyword).toBe("lys")
    expect(harness.getState().page).toBe(1)
    expect(harness.useCase.loadPage).toHaveBeenCalledTimes(1)
  })

  test("persona-read-detail-cutover: handlePluginFilterChange は plugin filter を更新して loadPage を呼ぶ", () => {
    const harness = createControllerHarness(
      createState({ pluginFilter: "", page: 1 })
    )
    const select = document.createElement("select")
    const option = document.createElement("option")
    option.value = "NightCourt.esp"
    select.append(option)
    select.value = "NightCourt.esp"

    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", {
      value: select,
      configurable: true
    })

    harness.controller.handlePluginFilterChange(event)

    expect(harness.getState().pluginFilter).toBe("NightCourt.esp")
    expect(harness.getState().page).toBe(1)
    expect(harness.useCase.loadPage).toHaveBeenCalledTimes(1)
  })

  test("persona-read-detail-cutover: selectRow は identity key で selectEntry を呼び detail を反映する", async () => {
    const harness = createControllerHarness()

    await harness.controller.selectRow("FollowersPlus.esp:FE01A812:NPC_")

    expect(harness.useCase.selectEntry).toHaveBeenCalledWith("FollowersPlus.esp:FE01A812:NPC_")
    expect(harness.useCase.selectEntry).toHaveBeenCalledTimes(1)
  })

  test("persona-read-detail-cutover: selectRow は useCase の dialogue modal 操作を起動しない", async () => {
    // Arrange
    const harness = createControllerHarness()
    const loadDialogueListSpy = vi.fn(async () => {})
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    ;(harness.useCase as any).loadDialogueList = loadDialogueListSpy

    // Act
    await harness.controller.selectRow("FollowersPlus.esp:FE01A812:NPC_")

    // Assert
    expect(harness.useCase.selectEntry).toHaveBeenCalledTimes(1)
    expect(loadDialogueListSpy).not.toHaveBeenCalled()
  })

  test("persona-edit-delete-cutover: setEditFormField は personaSummary を store へ反映する", () => {
    // Arrange
    const harness = createControllerHarness(createState())
    const textarea = document.createElement("textarea")
    textarea.value = "乾いた率直さで応じる。"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: textarea, configurable: true })

    // Act
    harness.controller.setEditFormField("personaSummary", event)

    // Assert
    expect(harness.getState().editForm.personaSummary).toBe("乾いた率直さで応じる。")
  })

  test("persona-read-detail-cutover: setEditFormField は personaBody を store へ反映する", () => {
    // Arrange
    const harness = createControllerHarness(createState())
    const textarea = document.createElement("textarea")
    textarea.value = "updated persona body"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: textarea, configurable: true })

    // Act
    harness.controller.setEditFormField("personaBody", event)

    // Assert
    expect(harness.getState().editForm.personaBody).toBe("updated persona body")
  })

  test("persona-read-detail-cutover: editForm の identity/snapshot fields は backend payload に含まれない", () => {
    // Arrange
    const state = createState({
      editForm: {
        formId: "FE01A812",
        editorId: "FP_LysMaren",
        displayName: "Lys Maren",
        race: "Nord",
        sex: "Female",
        voiceType: "FemaleYoungEager",
        className: "FPScoutClass",
        sourcePlugin: "FollowersPlus.esp",
        personaBody: "短く本音を置く。"
      }
    })

    // Act
    const payload = buildMasterPersonaUpdateInput(state)

    // Assert: cutover 後は personaSummary / speechStyle / personaBody のみを含む
    expect(Object.keys(payload).sort()).toEqual(["personaBody", "personaSummary", "speechStyle"])
  })

  test("persona-json-preview-cutover: stageJsonSelection は preview をクリアして file reference を更新する", () => {
    // Arrange
    const harness = createControllerHarness(createState())
    const file = new File(["{}"], "FollowersPlus.json", {
      type: "application/json"
    }) as File & { path: string }
    file.path = "/tmp/FollowersPlus.json"

    // Act
    harness.controller.stageJsonSelection(file)

    // Assert
    expect(harness.getState().selectedFileName).toBe("FollowersPlus.json")
    expect(harness.getState().selectedFileReference).toBe("/tmp/FollowersPlus.json")
    expect(harness.getState().preview).toBeNull()
    expect(harness.useCase.previewGeneration).toHaveBeenCalledTimes(1)
  })

  test("persona-generation-cutover: executeGeneration は useCase.executeGeneration を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.executeGeneration()

    // Assert
    expect(harness.useCase.executeGeneration).toHaveBeenCalledTimes(1)
  })

  test("persona-generation-cutover: interruptGeneration は useCase.interruptGeneration を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.interruptGeneration()

    // Assert
    expect(harness.useCase.interruptGeneration).toHaveBeenCalledTimes(1)
  })

  test("persona-generation-cutover: cancelGeneration は useCase.cancelGeneration を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.cancelGeneration()

    // Assert
    expect(harness.useCase.cancelGeneration).toHaveBeenCalledTimes(1)
  })

  test("persona-generation-cutover: stageJsonSelection(null) は runStatus を 入力待ち へリセットする", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        runStatus: {
          runState: "完了",
          targetPlugin: "FollowersPlus.esp",
          processedCount: 7,
          successCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "",
          message: "完了"
        }
      })
    )

    // Act
    harness.controller.stageJsonSelection(null)

    // Assert
    expect(harness.getState().runStatus.runState).toBe("入力待ち")
  })

  test("persona-edit-delete-cutover: saveCurrentEntry は useCase.saveCurrentEntry を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.saveCurrentEntry()

    // Assert
    expect(harness.useCase.saveCurrentEntry).toHaveBeenCalledTimes(1)
  })

  test("persona-edit-delete-cutover: deleteCurrentEntry は useCase.deleteCurrentEntry を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.deleteCurrentEntry()

    // Assert
    expect(harness.useCase.deleteCurrentEntry).toHaveBeenCalledTimes(1)
  })

  test("persona-edit-delete-cutover: openEditModal は selectedEntry が null のとき setModalState を呼ばない", () => {
    // Arrange
    const harness = createControllerHarness(createState({ selectedEntry: null }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.useCase.setModalState).not.toHaveBeenCalled()
  })

  test("persona-edit-delete-cutover: controller に openCreateModal は存在しない", () => {
    // Arrange
    const harness = createControllerHarness()

    // Assert: persona creation is not reintroduced
    expect("openCreateModal" in harness.controller).toBe(false)
  })
})
