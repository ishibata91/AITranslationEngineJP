import { vi } from "vitest"

import { MasterDictionaryRuntimeEventAdapter } from "./master-dictionary-runtime-event-adapter"

type RuntimeCallback = (...args: unknown[]) => void

type RuntimeRegistration = {
  eventName: string
  callback: RuntimeCallback
  maxCallbacks: number
}

function createRuntimeHarness() {
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

  return {
    adapter,
    registrations,
    detachProgress,
    detachCompleted,
    onImportProgress,
    onImportCompleted
  }
}

describe("MasterDictionaryRuntimeEventAdapter", () => {
  afterEach(() => {
    delete (window as Window & { runtime?: unknown }).runtime
  })

  test("runtime がある時は subscribe で true を返す", () => {
    // Arrange
    const { adapter } = createRuntimeHarness()

    // Act
    const subscribed = adapter.subscribe()

    // Assert
    expect(subscribed).toBe(true)
  })

  test("subscribe は progress event を登録する", () => {
    // Arrange
    const { adapter, registrations } = createRuntimeHarness()

    // Act
    adapter.subscribe()

    // Assert
    expect(registrations[0]?.eventName).toBe("master-dictionary:import-progress")
  })

  test("subscribe は progress event を無制限 callback で登録する", () => {
    // Arrange
    const { adapter, registrations } = createRuntimeHarness()

    // Act
    adapter.subscribe()

    // Assert
    expect(registrations[0]?.maxCallbacks).toBe(-1)
  })

  test("subscribe は completed event を登録する", () => {
    // Arrange
    const { adapter, registrations } = createRuntimeHarness()

    // Act
    adapter.subscribe()

    // Assert
    expect(registrations[1]?.eventName).toBe("master-dictionary:import-completed")
  })

  test("subscribe は completed event を無制限 callback で登録する", () => {
    // Arrange
    const { adapter, registrations } = createRuntimeHarness()

    // Act
    adapter.subscribe()

    // Assert
    expect(registrations[1]?.maxCallbacks).toBe(-1)
  })

  test("progress event payload を onImportProgress へ転送する", () => {
    // Arrange
    const { adapter, registrations, onImportProgress } = createRuntimeHarness()
    adapter.subscribe()

    // Act
    registrations[0]?.callback({ progress: 78 })

    // Assert
    expect(onImportProgress).toHaveBeenCalledWith({ progress: 78 })
  })

  test("completed event payload を onImportCompleted へ転送する", () => {
    // Arrange
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness()
    adapter.subscribe()

    // Act
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

    // Assert
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
  })

  test("detach は progress listener を解除する", () => {
    // Arrange
    const { adapter, detachProgress } = createRuntimeHarness()
    adapter.subscribe()

    // Act
    adapter.detach()

    // Assert
    expect(detachProgress).toHaveBeenCalledTimes(1)
  })

  test("detach は completed listener を解除する", () => {
    // Arrange
    const { adapter, detachCompleted } = createRuntimeHarness()
    adapter.subscribe()

    // Act
    adapter.detach()

    // Assert
    expect(detachCompleted).toHaveBeenCalledTimes(1)
  })

  test("runtime 不在時は subscribe で false を返す", () => {
    // Arrange
    const adapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress: vi.fn(),
      onImportCompleted: vi.fn()
    })

    // Act
    const subscribed = adapter.subscribe()

    // Assert
    expect(subscribed).toBe(false)
  })

  test("progress payload が object でない時は空 object へ正規化する", () => {
    // Arrange
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
    const adapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress,
      onImportCompleted: vi.fn()
    })
    adapter.subscribe()

    // Act
    registrations[0]?.callback(undefined)

    // Assert
    expect(onImportProgress).toHaveBeenCalledWith({})
  })

  test("completed payload が object でない時は空 object へ正規化する", () => {
    // Arrange
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
    const onImportCompleted = vi.fn()
    const adapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress: vi.fn(),
      onImportCompleted
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback("invalid")

    // Assert
    expect(onImportCompleted).toHaveBeenCalledWith({})
  })
})
