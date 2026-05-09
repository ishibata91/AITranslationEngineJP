import type {
  RuntimeImportCompletedPayload,
  RuntimeImportProgressPayload
} from "@application/contract/master-dictionary/master-dictionary-screen-types"
import type {
  ImportMasterDictionaryXmlResponse,
  MasterDictionaryPageEntry,
  MasterDictionaryPageState
} from "@application/gateway-contract/master-dictionary"

const MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT =
  "master-dictionary:import-progress"
const MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT =
  "master-dictionary:import-completed"

interface RuntimeImportEventHandlers {
  onImportProgress: (payload: RuntimeImportProgressPayload) => void
  onImportCompleted: (
    payload: RuntimeImportCompletedPayload
  ) => void | Promise<void>
}

type RuntimeEventLogFields = Record<string, string | undefined>

interface RuntimeEventDiagnosticLogger {
  info(message: string, fields?: RuntimeEventLogFields): void
  warn(message: string, fields?: RuntimeEventLogFields): void
}

interface WailsRuntimeEventBridge {
  EventsOnMultiple: (
    eventName: string,
    callback: (...data: unknown[]) => void,
    maxCallbacks: number
  ) => (() => void) | void
}

const noopRuntimeEventDiagnosticLogger: RuntimeEventDiagnosticLogger = {
  info: () => undefined,
  warn: () => undefined
}

function parseRuntimePayload<T extends object>(args: unknown[]): T {
  const [payload] = args
  if (!payload || typeof payload !== "object") {
    return {} as T
  }
  return payload as T
}

function hasObjectPayload(args: unknown[]): boolean {
  const [payload] = args
  return !!payload && typeof payload === "object"
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return !!value && typeof value === "object" && !Array.isArray(value)
}

function hasOwnProperty(
  value: Record<string, unknown>,
  propertyName: string
): boolean {
  return Object.prototype.hasOwnProperty.call(value, propertyName)
}

function isMasterDictionaryPageEntry(
  value: unknown
): value is MasterDictionaryPageEntry {
  if (!isRecord(value)) {
    return false
  }

  return (
    typeof value.id === "number" &&
    typeof value.source === "string" &&
    typeof value.translation === "string" &&
    typeof value.category === "string" &&
    typeof value.origin === "string" &&
    typeof value.updatedAt === "string"
  )
}

function isMasterDictionaryPageState(
  value: unknown
): value is MasterDictionaryPageState {
  if (!isRecord(value)) {
    return false
  }

  return (
    Array.isArray(value.items) &&
    value.items.every(isMasterDictionaryPageEntry) &&
    typeof value.totalCount === "number" &&
    typeof value.page === "number" &&
    typeof value.pageSize === "number" &&
    (value.selectedId === undefined || typeof value.selectedId === "number")
  )
}

function isImportSummary(
  value: unknown
): value is ImportMasterDictionaryXmlResponse["summary"] {
  if (!isRecord(value)) {
    return false
  }

  return (
    typeof value.filePath === "string" &&
    typeof value.fileName === "string" &&
    typeof value.importedCount === "number" &&
    typeof value.updatedCount === "number" &&
    typeof value.skippedCount === "number" &&
    typeof value.lastEntryId === "number"
  )
}

function isRuntimeImportCompletedPayload(
  value: unknown
): value is RuntimeImportCompletedPayload {
  if (!isRecord(value)) {
    return false
  }

  const hasPage = hasOwnProperty(value, "page")
  const hasSummary = hasOwnProperty(value, "summary")
  const hasValidPage = hasPage && isMasterDictionaryPageState(value.page)
  const hasValidSummary = hasSummary && isImportSummary(value.summary)

  if (hasPage && !hasValidPage) {
    return false
  }

  if (hasSummary && !hasValidSummary) {
    return false
  }

  return hasValidPage || hasValidSummary
}

function resolveRuntimeBridge(): WailsRuntimeEventBridge | null {
  if (typeof document !== "object") {
    return null
  }

  const runtime = (
    document.defaultView as (Window & { runtime?: unknown }) | null
  )?.runtime

  if (
    !runtime ||
    typeof (runtime as WailsRuntimeEventBridge).EventsOnMultiple !== "function"
  ) {
    return null
  }

  return runtime as WailsRuntimeEventBridge
}

export class MasterDictionaryRuntimeEventAdapter {
  private detachProgress: (() => void) | null = null

  private detachCompleted: (() => void) | null = null

  constructor(
    private readonly handlers: RuntimeImportEventHandlers,
    private readonly logger: RuntimeEventDiagnosticLogger = noopRuntimeEventDiagnosticLogger
  ) {}

  subscribe(): boolean {
    this.detach()

    const runtime = resolveRuntimeBridge()
    if (!runtime) {
      this.logRuntimeEvent("warn", "runtime event subscription skipped", {
        event: "runtime_event_subscribe",
        where: "frontend.runtime.master_dictionary",
        result: "skipped",
        reason: "runtime_unavailable"
      })
      return false
    }

    this.detachProgress =
      runtime.EventsOnMultiple(
        MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT,
        (...args: unknown[]) => {
          this.handleImportProgressEvent(args)
        },
        -1
      ) ?? null

    this.detachCompleted =
      runtime.EventsOnMultiple(
        MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT,
        (...args: unknown[]) => {
          void this.handleImportCompletedEvent(args)
        },
        -1
      ) ?? null

    this.logRuntimeEvent("info", "runtime event subscribed", {
      event: "runtime_event_subscribe",
      where: "frontend.runtime.master_dictionary",
      result: "accepted"
    })

    return true
  }

  detach(): void {
    const hadSubscription =
      this.detachProgress !== null || this.detachCompleted !== null

    this.detachProgress?.()
    this.detachCompleted?.()
    this.detachProgress = null
    this.detachCompleted = null

    if (hadSubscription) {
      this.logRuntimeEvent("info", "runtime event detached", {
        event: "runtime_event_detach",
        where: "frontend.runtime.master_dictionary",
        result: "detached"
      })
    }
  }

  pushImportStateToStore(payload: RuntimeImportProgressPayload): void {
    this.handlers.onImportProgress(payload)
  }

  private handleImportProgressEvent(args: unknown[]): void {
    if (!hasObjectPayload(args)) {
      this.logRuntimeEvent("warn", "runtime event dropped", {
        event: "runtime_event_progress",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "payload_parse_failed"
      })
      return
    }

    const payload = parseRuntimePayload<RuntimeImportProgressPayload>(args)

    if (Object.keys(payload).length > 0 && typeof payload.progress !== "number") {
      this.logRuntimeEvent("warn", "runtime event skipped", {
        event: "runtime_event_progress",
        where: "frontend.runtime.master_dictionary",
        result: "skipped",
        reason: "invalid_progress"
      })
      return
    }

    this.logRuntimeEvent("info", "runtime event accepted", {
      event: "runtime_event_progress",
      where: "frontend.runtime.master_dictionary",
      result: "accepted"
    })
    this.pushImportStateToStore(payload)
  }

  private async handleImportCompletedEvent(args: unknown[]): Promise<void> {
    if (!hasObjectPayload(args)) {
      this.logRuntimeEvent("warn", "runtime event dropped", {
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "payload_parse_failed"
      })
      return
    }

    const payload = parseRuntimePayload<RuntimeImportCompletedPayload>(args)

    if (!isRuntimeImportCompletedPayload(payload)) {
      this.logRuntimeEvent("warn", "runtime event dropped", {
        event: "runtime_event_completed",
        where: "frontend.runtime.master_dictionary",
        result: "dropped",
        reason: "invalid_payload"
      })
      return
    }

    this.logRuntimeEvent("info", "runtime event accepted", {
      event: "runtime_event_completed",
      where: "frontend.runtime.master_dictionary",
      result: "accepted"
    })
    await this.handlers.onImportCompleted(payload)
  }

  private logRuntimeEvent(
    level: "info" | "warn",
    message: string,
    fields: RuntimeEventLogFields
  ): void {
    this.logger[level](message, fields)
  }
}
