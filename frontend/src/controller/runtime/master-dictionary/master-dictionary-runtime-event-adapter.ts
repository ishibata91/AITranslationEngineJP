import type {
  RuntimeImportCompletedPayload,
  RuntimeImportProgressPayload
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

const MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT =
  "master-dictionary:import-progress"
const MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT =
  "master-dictionary:import-completed"

interface RuntimeImportEventHandlers {
  onImportProgress: (payload: RuntimeImportProgressPayload) => void
  onImportCompleted: (payload: RuntimeImportCompletedPayload) => void
}

interface WailsRuntimeEventBridge {
  EventsOnMultiple: (
    eventName: string,
    callback: (...data: unknown[]) => void,
    maxCallbacks: number
  ) => (() => void) | void
}

function parseRuntimePayload<T extends object>(args: unknown[]): T | null {
  const [payload] = args
  if (!payload || typeof payload !== "object") {
    return null
  }
  return payload as T
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

  constructor(private readonly handlers: RuntimeImportEventHandlers) {}

  subscribe(): boolean {
    this.detach()

    const runtime = resolveRuntimeBridge()
    if (!runtime) {
      return false
    }

    this.detachProgress =
      runtime.EventsOnMultiple(
        MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT,
        (...args: unknown[]) => {
          this.pushImportStateToStore(
            parseRuntimePayload<RuntimeImportProgressPayload>(args) ?? {}
          )
        },
        -1
      ) ?? null

    this.detachCompleted =
      runtime.EventsOnMultiple(
        MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT,
        (...args: unknown[]) => {
          this.handlers.onImportCompleted(
            parseRuntimePayload<RuntimeImportCompletedPayload>(args) ?? {}
          )
        },
        -1
      ) ?? null

    return true
  }

  detach(): void {
    this.detachProgress?.()
    this.detachCompleted?.()
    this.detachProgress = null
    this.detachCompleted = null
  }

  pushImportStateToStore(payload: RuntimeImportProgressPayload): void {
    this.handlers.onImportProgress(payload)
  }
}
