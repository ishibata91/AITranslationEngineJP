import { vi } from "vitest"

import { MasterDictionaryRuntimeEventAdapter } from "./master-dictionary-runtime-event-adapter"

type RuntimeCallback = (...args: unknown[]) => void

type RuntimeRegistration = {
  eventName: string
  callback: RuntimeCallback
  maxCallbacks: number
}

type RuntimeEventLoggerFields = Record<string, string | undefined>

type RuntimeEventLogger = {
  info: (message: string, fields?: RuntimeEventLoggerFields) => void
  warn: (message: string, fields?: RuntimeEventLoggerFields) => void
}

function createRuntimeHarness(options?: { logger?: RuntimeEventLogger }) {
  const registrations: RuntimeRegistration[] = []
  const detachProgress = vi.fn()
  const detachCompleted = vi.fn()

  ;(window as Window & { runtime?: unknown }).runtime = {
    EventsOnMultiple: vi.fn(
      (eventName: string, callback: RuntimeCallback, maxCallbacks: number) => {
        registrations.push({ eventName, callback, maxCallbacks })
        return eventName === "master-dictionary:import-progress"
          ? detachProgress
          : detachCompleted
      }
    )
  }

  const onImportProgress = vi.fn()
  const onImportCompleted = vi.fn()
  const adapter = new MasterDictionaryRuntimeEventAdapter(
    {
      onImportProgress,
      onImportCompleted
    },
    options?.logger
  )

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
    expect(registrations[0]?.eventName).toBe(
      "master-dictionary:import-progress"
    )
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
    expect(registrations[1]?.eventName).toBe(
      "master-dictionary:import-completed"
    )
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
        lastEntryId: 201
      }
    })
  })

  test("completed event の page payload を onImportCompleted へ転送する", () => {
    // Arrange
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness()
    adapter.subscribe()
    const page = {
      items: [
        {
          id: 201,
          source: "source text",
          translation: "translated text",
          category: "NPC",
          origin: "master.xml",
          updatedAt: "2026-05-09T00:00:00Z"
        }
      ],
      totalCount: 1,
      page: 1,
      pageSize: 50,
      selectedId: 201
    }

    // Act
    registrations[1]?.callback({ page })

    // Assert
    expect(onImportCompleted).toHaveBeenCalledWith({ page })
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

  test("runtime 不在時は skipped reason=runtime_unavailable を warn log へ出す", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const adapter = new MasterDictionaryRuntimeEventAdapter(
      {
        onImportProgress: vi.fn(),
        onImportCompleted: vi.fn()
      },
      logger
    )

    // Act
    adapter.subscribe()

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event subscription skipped",
      expect.objectContaining({
        event: "runtime_event_subscribe",
        where: "frontend.runtime.master_dictionary",
        result: "skipped",
        reason: "runtime_unavailable"
      })
    )
  })

  test("payload parse 失敗は dropped log として扱われ store を更新しない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportProgress } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[0]?.callback(undefined)

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_progress",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "payload_parse_failed"
      })
    )
    expect(onImportProgress).not.toHaveBeenCalled()
  })

  test("progress が number でない payload は skipped reason=invalid_progress で store を更新しない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportProgress } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[0]?.callback({ progress: "invalid" })

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event skipped",
      expect.objectContaining({
        event: "runtime_event_progress",
        where: "frontend.runtime.master_dictionary",
        result: "skipped",
        reason: "invalid_progress"
      })
    )
    expect(onImportProgress).not.toHaveBeenCalled()
  })

  test("progress payload が正しい時は accepted を info log へ出す", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations } = createRuntimeHarness({ logger })
    adapter.subscribe()

    // Act
    registrations[0]?.callback({ progress: 78 })

    // Assert
    expect(logger.info).toHaveBeenCalledWith(
      "runtime event accepted",
      expect.objectContaining({
        event: "runtime_event_progress",
        where: "frontend.runtime.master_dictionary",
        result: "accepted"
      })
    )
  })

  test("detach 後は detached を info log へ出す", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter } = createRuntimeHarness({ logger })
    adapter.subscribe()

    // Act
    adapter.detach()

    // Assert
    expect(logger.info).toHaveBeenCalledWith(
      "runtime event detached",
      expect.objectContaining({
        event: "runtime_event_detach",
        where: "frontend.runtime.master_dictionary",
        result: "detached"
      })
    )
  })

  test("completed payload が object でない時は dropped log として扱われ完了 handler を呼ばない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback("invalid")

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "payload_parse_failed"
      })
    )
    expect(onImportCompleted).not.toHaveBeenCalled()
  })

  test("completed payload が空 object の時は dropped log として扱われ完了 handler を呼ばない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback({})

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "invalid_payload"
      })
    )
    expect(onImportCompleted).not.toHaveBeenCalled()
  })

  test("completed payload が未知 key だけの object の時は dropped log として扱われ完了 handler を呼ばない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback({ unknown: true })

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "invalid_payload"
      })
    )
    expect(onImportCompleted).not.toHaveBeenCalled()
  })

  test("completed payload の page が有効な page state でない時は dropped log として扱われ完了 handler を呼ばない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback({
      page: {
        items: [],
        totalCount: 1
      }
    })

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "invalid_payload"
      })
    )
    expect(onImportCompleted).not.toHaveBeenCalled()
  })

  test("completed payload の page が不正なら summary が有効でも dropped log として扱われ完了 handler を呼ばない", () => {
    // Arrange
    const logger = {
      info: vi.fn(),
      warn: vi.fn()
    }
    const { adapter, registrations, onImportCompleted } = createRuntimeHarness({
      logger
    })
    adapter.subscribe()

    // Act
    registrations[1]?.callback({
      page: {
        items: [],
        totalCount: 1
      },
      summary: {
        filePath: "master.xml",
        fileName: "master.xml",
        importedCount: 2,
        updatedCount: 0,
        skippedCount: 1,
        lastEntryId: 201
      }
    })

    // Assert
    expect(logger.warn).toHaveBeenCalledWith(
      "runtime event dropped",
      expect.objectContaining({
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "invalid_payload"
      })
    )
    expect(onImportCompleted).not.toHaveBeenCalled()
  })
})
