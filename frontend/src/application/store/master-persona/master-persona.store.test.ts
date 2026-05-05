import { describe, expect, test, vi } from "vitest"

import type { MasterPersonaDetail } from "@application/gateway-contract/master-persona"

import { MasterPersonaStore } from "./master-persona.store"

function createMasterPersonaDetail(): MasterPersonaDetail {
  return {
    identityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
    targetPlugin: "TestPersonaPluginA.esp",
    formId: "FE01A812",
    recordType: "NPC_",
    editorId: "TEST_NPC_A",
    displayName: "Test NPC A",
    voiceType: "TestVoiceA",
    className: "TestClassA",
    sourcePlugin: "TestPersonaPluginA.esp",
    personaSummary: "test summary",
    updatedAt: "2026-04-15T09:42:00Z",
    personaBody: "test body",
    runLockReason: "更新と削除を行えます"
  }
}

describe("MasterPersonaStore", () => {
  test("subscribe は初期 state を即時通知する", () => {
    const store = new MasterPersonaStore()
    const listener = vi.fn()

    store.subscribe(listener)

    expect(listener).toHaveBeenCalledTimes(1)
    expect(listener.mock.calls[0]?.[0]).toMatchObject({
      items: [],
      pluginGroups: [],
      selectedIdentityKey: null,
      selectedEntry: null,
      selectedFileName: "未選択",
      selectedFileReference: null,
      preview: null,
      modalState: null
    })
  })

  test("unsubscribe 後は update しても通知されない", () => {
    const store = new MasterPersonaStore()
    const listener = vi.fn()

    const unsubscribe = store.subscribe(listener)
    unsubscribe()
    store.update((draft) => {
      draft.keyword = "test"
    })

    expect(listener).toHaveBeenCalledTimes(1)
  })

  test("snapshot は nested object の defensive copy を返す", () => {
    const store = new MasterPersonaStore()

    store.update((draft) => {
      const entry = createMasterPersonaDetail()
      draft.items = [entry]
      draft.pluginGroups = [{ targetPlugin: "TestPersonaPluginA.esp", count: 1 }]
      draft.selectedIdentityKey = entry.identityKey
      draft.selectedEntry = entry
      draft.aiSettings = {
        provider: "gemini",
        model: "gemini-2.5-pro",
        executionMethod: "single_request"
      }
      draft.modelOptions = [{ modelId: "gemini-2.5-pro", label: "Gemini Pro" }]
      draft.preview = {
        fileName: "sample.json",
        targetPlugin: "TestPersonaPluginA.esp",
        status: "生成可能"
      }
      draft.runStatus = {
        runState: "生成中",
        targetPlugin: "TestPersonaPluginA.esp",
        processedCount: 1,
        successCount: 0,
        existingSkipCount: 0,
        currentActorLabel: "Test NPC A",
        message: "running"
      }
      draft.editForm = {
        personaSummary: "summary",
        speechStyle: "style",
        personaBody: "body"
      }
    })

    const snapshot = store.snapshot()
    snapshot.items[0].displayName = "changed"
    snapshot.pluginGroups[0].count = 99
    snapshot.selectedEntry!.displayName = "changed"
    snapshot.aiSettings.provider = "changed"
    snapshot.modelOptions[0].label = "changed"
    snapshot.preview!.status = "changed"
    snapshot.runStatus.message = "changed"
    snapshot.editForm.personaBody = "changed"

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.items[0]?.displayName).toBe("Test NPC A")
    expect(nextSnapshot.pluginGroups[0]?.count).toBe(1)
    expect(nextSnapshot.selectedEntry?.displayName).toBe("Test NPC A")
    expect(nextSnapshot.aiSettings.provider).toBe("gemini")
    expect(nextSnapshot.modelOptions[0]?.label).toBe("Gemini Pro")
    expect(nextSnapshot.preview?.status).toBe("生成可能")
    expect(nextSnapshot.runStatus.message).toBe("running")
    expect(nextSnapshot.editForm.personaBody).toBe("body")
  })
})
