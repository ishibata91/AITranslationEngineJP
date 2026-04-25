import type {
  TranslationInputCommandResponse,
  TranslationInputGatewayContract,
  TranslationInputReviewItem,
  TranslationInputReviewStatus,
  TranslationInputScreenState,
  TranslationInputStagedFile,
  TranslationInputWarning
} from "@application/gateway-contract/translation-input"

interface TranslationInputStoreLike {
  snapshot(): TranslationInputScreenState
  update(mutator: (draft: TranslationInputScreenState) => void): void
}

function toErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message.trim() !== "") {
    return error.message
  }

  return fallback
}

function fileNameFromPath(filePath: string): string {
  const segments = filePath.split(/[\\/]/)
  return segments.at(-1) || filePath
}

function normalizeWarnings(
  warnings: TranslationInputWarning[] | null | undefined
): TranslationInputWarning[] {
  return Array.isArray(warnings) ? warnings : []
}

function normalizeSummary(
  summary: TranslationInputCommandResponse["summary"] | null | undefined
): TranslationInputCommandResponse["summary"] | undefined {
  if (!summary) {
    return undefined
  }

  return {
    ...summary,
    categories: Array.isArray(summary.categories) ? summary.categories : [],
    sampleFields: Array.isArray(summary.sampleFields) ? summary.sampleFields : [],
    warnings: normalizeWarnings(summary.warnings)
  }
}

function normalizeCommandResponse(
  response: TranslationInputCommandResponse
): TranslationInputCommandResponse {
  return {
    ...response,
    summary: normalizeSummary(response.summary),
    warnings: normalizeWarnings(response.warnings)
  }
}

function mergeWarnings(
  responseWarnings: TranslationInputWarning[] | null | undefined,
  summaryWarnings: TranslationInputWarning[] | null | undefined
): TranslationInputWarning[] {
  const merged = [
    ...normalizeWarnings(responseWarnings),
    ...normalizeWarnings(summaryWarnings)
  ]
  const seen = new Set<string>()
  return merged.filter((warning) => {
    const key = `${warning.kind}:${warning.recordType}:${warning.subrecordType}:${warning.message}`
    if (seen.has(key)) {
      return false
    }
    seen.add(key)
    return true
  })
}

function decideItemStatus(
  accepted: boolean,
  errorKind: string | undefined,
  warnings: TranslationInputWarning[]
): TranslationInputReviewStatus {
  if (!accepted && errorKind === "source_file_missing") {
    return "rebuild-required"
  }

  if (!accepted) {
    return "failed"
  }

  if (warnings.length > 0) {
    return "warning"
  }

  return "registered"
}

function buildReviewItemFromResponse(
  response: TranslationInputCommandResponse,
  stagedFile: TranslationInputStagedFile,
  action: "import" | "rebuild",
  existingItem?: TranslationInputReviewItem | null
): TranslationInputReviewItem {
  const normalizedResponse = normalizeCommandResponse(response)
  const summary = normalizedResponse.summary ?? existingItem?.summary ?? null
  const warnings = mergeWarnings(response.warnings, summary?.warnings ?? [])
  const filePath = summary?.input.sourceFilePath ?? existingItem?.filePath ?? stagedFile.filePath
  const fileName = existingItem?.fileName ?? summary?.input.sourceFilePath
    ? fileNameFromPath(filePath)
    : stagedFile.fileName
  const inputId = summary?.input.id ?? existingItem?.inputId ?? null
  const importTimestamp =
    summary?.input.importedAt ?? existingItem?.importTimestamp ?? new Date().toISOString()
  const errorKind = response.accepted
    ? null
    : (normalizedResponse.errorKind ?? existingItem?.errorKind ?? null)
  const localId = existingItem?.localId ?? `${inputId ?? stagedFile.fileHash}:${importTimestamp}:${action}`

  return {
    localId,
    inputId,
    fileName,
    filePath,
    fileHash: existingItem?.fileHash ?? stagedFile.fileHash,
    importTimestamp,
    status: decideItemStatus(response.accepted, response.errorKind, warnings),
    accepted: response.accepted,
    canRebuild: inputId !== null,
    lastAction: action,
    errorKind,
    warnings,
    summary
  }
}

function upsertItem(
  items: TranslationInputReviewItem[],
  nextItem: TranslationInputReviewItem
): TranslationInputReviewItem[] {
  const nextItems = [...items]
  const existingIndex = nextItems.findIndex(
    (item) =>
      item.localId === nextItem.localId ||
      (item.inputId !== null && item.inputId === nextItem.inputId)
  )

  if (existingIndex >= 0) {
    nextItems.splice(existingIndex, 1, nextItem)
    return nextItems
  }

  return [nextItem, ...nextItems]
}

function createSyntheticStagedFile(
  state: TranslationInputScreenState,
  selectedItem: TranslationInputReviewItem | null
): TranslationInputStagedFile {
  if (state.stagedFile) {
    return state.stagedFile
  }

  if (selectedItem) {
    return {
      fileName: selectedItem.fileName,
      filePath: selectedItem.filePath,
      fileHash: selectedItem.fileHash
    }
  }

  return {
    fileName: "未選択",
    filePath: "",
    fileHash: "-"
  }
}

export class TranslationInputUseCase {
  constructor(
    private readonly gateway: TranslationInputGatewayContract | null,
    private readonly store: TranslationInputStoreLike
  ) {}

  async startImport(): Promise<void> {
    const state = this.store.snapshot()
    if (!state.stagedFile) {
      this.store.update((draft) => {
        draft.errorMessage = "登録する JSON file を選択してください。"
        draft.latestResponse = {
          accepted: false,
          errorKind: "missing_required_field",
          warnings: []
        }
      })
      return
    }

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.errorMessage = "translation-input gateway が未接続です。"
      })
      return
    }

    this.store.update((draft) => {
      draft.operationState = "importing"
      draft.errorMessage = ""
    })

    try {
      const response = normalizeCommandResponse(await this.gateway.importTranslationInput({
        filePath: state.stagedFile.filePath
      }))
      const nextItem = buildReviewItemFromResponse(
        response,
        state.stagedFile,
        "import"
      )

      this.store.update((draft) => {
        draft.items = upsertItem(draft.items, nextItem)
        draft.selectedItemId = nextItem.localId
        draft.operationState = "idle"
        draft.stagedFile = null
        draft.latestResponse = response
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.operationState = "ready"
        draft.errorMessage = toErrorMessage(error, "入力データの登録に失敗しました。")
      })
    }
  }

  async rebuildSelected(): Promise<void> {
    const state = this.store.snapshot()
    const selectedItem =
      state.items.find((item) => item.localId === state.selectedItemId) ?? null
    if (!selectedItem || selectedItem.inputId === null) {
      this.store.update((draft) => {
        draft.errorMessage = "再構築対象の cache がありません。"
        draft.latestResponse = {
          accepted: false,
          errorKind: "cache_missing",
          warnings: []
        }
      })
      return
    }

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.errorMessage = "translation-input gateway が未接続です。"
      })
      return
    }

    this.store.update((draft) => {
      draft.operationState = "rebuilding"
      draft.errorMessage = ""
    })

    try {
      const response = normalizeCommandResponse(await this.gateway.rebuildTranslationInputCache({
        inputId: selectedItem.inputId
      }))
      const nextItem = buildReviewItemFromResponse(
        response,
        createSyntheticStagedFile(state, selectedItem),
        "rebuild",
        selectedItem
      )

      this.store.update((draft) => {
        draft.items = upsertItem(draft.items, nextItem)
        draft.selectedItemId = nextItem.localId
        draft.operationState = "idle"
        draft.latestResponse = response
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.operationState = "idle"
        draft.errorMessage = toErrorMessage(error, "入力データの再構築に失敗しました。")
      })
    }
  }
}