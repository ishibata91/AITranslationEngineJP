import { vi } from "vitest"

import { MasterDictionaryRuntimeEventAdapter } from "./master-dictionary-runtime-event-adapter"

type RuntimeCallback = (...args: unknown[]) => void

type RuntimeRegistration = {
  eventName: string
  callback: RuntimeCallback
  maxCallbacks: number
}

describe("MasterDictionaryRuntimeEventAdapter", () => {
  afterEach(() => {
    delete (window as Window & { runtime?: unknown }).runtime
  })

  test("progress と completed を Wails runtime event へ登録し payload を forward する", () => {
    const registrations: RuntimeRegistration[] = []
    const detachProgress = vi.fn()
    const detachCompleted = vi.fn()

    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(
        (
          eventName: string,
          callback: RuntimeCallback,
          maxCallbacks: number
        ) => {
          registrations.push({ eventName, callback, maxCallbacks })
          return eventName === "master-dictionary:import-progress"
            ? detachProgress
            : detachCompleted
        }
      )
    }

    const onImportProgress = vi.fn()
    const onImportCompleted = vi.fn()
    const adapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress,
      onImportCompleted
    })

    expect(adapter.subscribe()).toBe(true)
    expect(registrations).toEqual([
      {
        eventName: "master-dictionary:import-progress",
        callback: registrations[0]?.callback,
        maxCallbacks: -1
      },
      {
        eventName: "master-dictionary:import-completed",
        callback: registrations[1]?.callback,
        maxCallbacks: -1
      }
    ])

    registrations[0]?.callback({ progress: 78 })
    registrations[1]?.callback({
      summary: {
        filePath: "master.xml",
        fileName: "master.xml",
        importedCount: 2,
        updatedCount: 0,
        skippedCount: 1,
        selectedRec: ["BOOK:FULL"],
        lastEntryId: 201
      }
    })

    expect(onImportProgress).toHaveBeenCalledWith({ progress: 78 })
    expect(onImportCompleted).toHaveBeenCalledWith({
      summary: {
        filePath: "master.xml",
        fileName: "master.xml",
        importedCount: 2,
        updatedCount: 0,
        skippedCount: 1,
        selectedRec: ["BOOK:FULL"],
        lastEntryId: 201
      }
    })

    adapter.detach()

    expect(detachProgress).toHaveBeenCalledTimes(1)
    expect(detachCompleted).toHaveBeenCalledTimes(1)
  })

  test("runtime payload が object でない時は空 object へ正規化し runtime 不在時は subscribe しない", () => {
    const missingRuntimeAdapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress: vi.fn(),
      onImportCompleted: vi.fn()
    })

    expect(missingRuntimeAdapter.subscribe()).toBe(false)

    const registrations: RuntimeRegistration[] = []
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(
        (
          eventName: string,
          callback: RuntimeCallback,
          maxCallbacks: number
        ) => {
          registrations.push({ eventName, callback, maxCallbacks })
          return vi.fn()
        }
      )
    }

    const onImportProgress = vi.fn()
    const onImportCompleted = vi.fn()
    const adapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress,
      onImportCompleted
    })

    expect(adapter.subscribe()).toBe(true)

    registrations[0]?.callback(undefined)
    registrations[1]?.callback("invalid")

    expect(onImportProgress).toHaveBeenCalledWith({})
    expect(onImportCompleted).toHaveBeenCalledWith({})
  })
})
