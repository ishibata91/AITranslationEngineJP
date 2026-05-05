import { describe, expect, it } from "vitest"

import { TranslationInputStore } from "./translation-input.store"

function createSummary() {
  return {
    input: {
      id: 1,
      sourceFilePath: "/tmp/test-input.json",
      sourceTool: "xedit-json",
      targetPluginName: "TestInputPluginA.esp",
      targetPluginType: "esp",
      recordCount: 2,
      importedAt: "2026-05-05T00:00:00Z"
    },
    translationRecordCount: 2,
    translationFieldCount: 3,
    categories: [
      {
        category: "NPC_",
        recordCount: 1,
        fieldCount: 2
      }
    ],
    sampleFields: [
      {
        recordType: "NPC_",
        subrecordType: "FULL",
        formId: "TEST_NPC_A",
        editorId: "TestNpcA",
        sourceText: "Test NPC A",
        translatable: true
      }
    ],
    warnings: [
      {
        kind: "unknown_field_definition" as const,
        recordType: "NPC_",
        subrecordType: "DESC",
        message: "unknown"
      }
    ]
  }
}

describe("TranslationInputStore", () => {
  it("subscribe immediately notifies current state", () => {
    const store = new TranslationInputStore()
    const snapshots: unknown[] = []

    store.subscribe((state) => {
      snapshots.push(state)
    })

    expect(snapshots).toHaveLength(1)
    expect(store.snapshot()).toMatchObject({
      items: [],
      selectedItemId: null,
      stagedFile: null,
      operationState: "idle",
      errorMessage: "",
      latestResponse: null
    })
  })

  it("unsubscribe removes listener from future updates", () => {
    const store = new TranslationInputStore()
    let callCount = 0

    const unsubscribe = store.subscribe(() => {
      callCount += 1
    })
    unsubscribe()

    store.update((draft) => {
      draft.operationState = "importing"
    })

    expect(callCount).toBe(1)
  })

  it("snapshot returns defensive copies", () => {
    const store = new TranslationInputStore()

    store.update((draft) => {
      draft.items = [
        {
          localId: "local-a",
          inputId: 1,
          fileName: "test-input.json",
          filePath: "/tmp/test-input.json",
          fileHash: "hash-a",
          importTimestamp: "2026-05-05T00:00:00Z",
          status: "registered",
          accepted: true,
          canRebuild: true,
          lastAction: "import",
          errorKind: null,
          warnings: createSummary().warnings,
          summary: createSummary()
        }
      ]
      draft.selectedItemId = "local-a"
      draft.stagedFile = {
        fileName: "test-input.json",
        filePath: "/tmp/test-input.json",
        fileHash: "hash-a"
      }
      draft.latestResponse = {
        accepted: true,
        warnings: createSummary().warnings,
        summary: createSummary()
      }
    })

    const snapshot = store.snapshot()
    snapshot.items[0].fileName = "changed.json"
    snapshot.items[0].warnings[0].message = "changed"
    snapshot.items[0].summary!.input.targetPluginName = "Changed.esp"
    snapshot.items[0].summary!.categories[0].category = "CHANGED"
    snapshot.items[0].summary!.sampleFields[0].sourceText = "changed"
    snapshot.items[0].summary!.warnings[0].message = "changed"
    snapshot.stagedFile!.fileHash = "changed"
    snapshot.latestResponse!.warnings[0].message = "changed"
    snapshot.latestResponse!.summary!.input.targetPluginName = "Changed.esp"
    snapshot.latestResponse!.summary!.categories[0].category = "CHANGED"
    snapshot.latestResponse!.summary!.sampleFields[0].sourceText = "changed"
    snapshot.latestResponse!.summary!.warnings[0].message = "changed"

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.items[0]?.fileName).toBe("test-input.json")
    expect(nextSnapshot.items[0]?.warnings[0]?.message).toBe("unknown")
    expect(nextSnapshot.items[0]?.summary?.input.targetPluginName).toBe(
      "TestInputPluginA.esp"
    )
    expect(nextSnapshot.items[0]?.summary?.categories[0]?.category).toBe("NPC_")
    expect(nextSnapshot.items[0]?.summary?.sampleFields[0]?.sourceText).toBe(
      "Test NPC A"
    )
    expect(nextSnapshot.items[0]?.summary?.warnings[0]?.message).toBe("unknown")
    expect(nextSnapshot.stagedFile?.fileHash).toBe("hash-a")
    expect(nextSnapshot.latestResponse?.warnings[0]?.message).toBe("unknown")
    expect(nextSnapshot.latestResponse?.summary?.input.targetPluginName).toBe(
      "TestInputPluginA.esp"
    )
    expect(nextSnapshot.latestResponse?.summary?.categories[0]?.category).toBe(
      "NPC_"
    )
    expect(nextSnapshot.latestResponse?.summary?.sampleFields[0]?.sourceText).toBe(
      "Test NPC A"
    )
    expect(nextSnapshot.latestResponse?.summary?.warnings[0]?.message).toBe(
      "unknown"
    )
  })
})
