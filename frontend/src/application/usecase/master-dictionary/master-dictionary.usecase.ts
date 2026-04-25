import type {
  MasterDictionaryEntryDetail,
  MasterDictionaryGatewayContract,
  MasterDictionaryPageEntry,
  MasterDictionaryPageState
} from "@application/gateway-contract/master-dictionary"
import {
  buildRefreshPayload,
  buildUpsertPayload
} from "@application/contract/master-dictionary"
import type {
  RuntimeImportCompletedPayload,
  RuntimeImportProgressPayload
} from "@application/contract/master-dictionary/master-dictionary-screen-types"
import { MasterDictionaryStore } from "@application/store/master-dictionary"

function toErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message.trim() !== "") {
    return error.message
  }

  if (typeof error === "string" && error.trim() !== "") {
    return error
  }

  if (error && typeof error === "object") {
    const errorFields = error as {
      message?: unknown
      error?: unknown
      cause?: unknown
    }

    for (const key of ["message", "error", "cause"] as const) {
      const value = errorFields[key]
      if (typeof value === "string" && value.trim() !== "") {
        return value
      }
    }
  }

  return fallback
}

function toDetailFromPageEntry(
  entry: MasterDictionaryPageEntry
): MasterDictionaryEntryDetail {
  return {
    id: String(entry.id),
    source: entry.source,
    translation: entry.translation,
    category: entry.category,
    origin: entry.origin,
    updatedAt: entry.updatedAt,
    note: "マスター辞書エントリ"
  }
}

function chooseSelectedId(
  list: { id: string }[],
  preferredId: string | null | undefined
): string | null {
  if (list.length === 0) {
    return null
  }
  if (preferredId && list.some((entry) => entry.id === preferredId)) {
    return preferredId
  }
  return list[0]?.id ?? null
}

function chooseSelectedIdFromPage(
  pageState: MasterDictionaryPageState,
  preferredId: string | null | undefined
): string | null {
  const availableIds = new Set(pageState.items.map((item) => String(item.id)))
  if (preferredId && availableIds.has(preferredId)) {
    return preferredId
  }
  if (typeof pageState.selectedId === "number") {
    const backendSelectedId = String(pageState.selectedId)
    if (availableIds.has(backendSelectedId)) {
      return backendSelectedId
    }
  }
  return pageState.items.length > 0 ? String(pageState.items[0]?.id) : null
}

export class MasterDictionaryUseCase {
  private gateway: MasterDictionaryGatewayContract | null

  private listRequestSequence = 0

  private detailRequestSequence = 0

  constructor(
    gateway: MasterDictionaryGatewayContract | null,
    private readonly store: MasterDictionaryStore
  ) {
    this.gateway = gateway
  }

  async loadEntries(preferredId?: string | null): Promise<void> {
    const activeSequence = ++this.listRequestSequence
    const state = this.store.snapshot()

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.entries = []
        draft.totalCount = 0
        draft.selectedId = null
        draft.selectedEntry = null
      })
      return
    }

    try {
      const response = await this.gateway.listMasterDictionaryEntries({
        filters: {
          query: state.query,
          category: state.category,
          page: state.page + 1,
          pageSize: 30
        }
      })
      if (activeSequence !== this.listRequestSequence) {
        return
      }

      const nextSelectedId = chooseSelectedId(
        response.entries,
        preferredId === undefined ? state.selectedId : preferredId
      )

      this.store.update((draft) => {
        draft.entries = response.entries
        draft.totalCount = response.totalCount
        draft.page = Math.max(0, response.page - 1)
        draft.selectedId = nextSelectedId
        draft.selectedEntry = nextSelectedId ? draft.selectedEntry : null
      })

      if (!nextSelectedId) {
        this.store.update((draft) => {
          draft.selectedEntry = null
        })
        return
      }

      await this.loadEntryDetail(nextSelectedId)
    } catch (error) {
      if (activeSequence !== this.listRequestSequence) {
        return
      }
      this.store.update((draft) => {
        draft.entries = []
        draft.totalCount = 0
        draft.selectedId = null
        draft.selectedEntry = null
        draft.errorMessage = toErrorMessage(error, "一覧の取得に失敗しました。")
      })
    }
  }

  async loadEntryDetail(id: string | null): Promise<void> {
    if (!this.gateway || !id) {
      this.store.update((draft) => {
        draft.selectedEntry = null
      })
      return
    }

    const activeSequence = ++this.detailRequestSequence
    try {
      const response = await this.gateway.getMasterDictionaryEntry({ id })
      if (activeSequence !== this.detailRequestSequence) {
        return
      }
      this.store.update((draft) => {
        draft.selectedEntry = response.entry
      })
    } catch (error) {
      if (activeSequence !== this.detailRequestSequence) {
        return
      }
      this.store.update((draft) => {
        draft.selectedEntry = null
        draft.errorMessage = toErrorMessage(error, "詳細の取得に失敗しました。")
      })
    }
  }

  async selectEntry(id: string): Promise<void> {
    this.store.update((draft) => {
      draft.selectedId = id
      draft.errorMessage = ""
    })
    await this.loadEntryDetail(id)
  }

  async saveCurrentEntry(): Promise<void> {
    const state = this.store.snapshot()

    if (
      !this.gateway ||
      (state.modalState !== "create" && state.modalState !== "edit")
    ) {
      return
    }

    const payload = buildUpsertPayload(state)
    if (!payload.source || !payload.translation) {
      this.store.update((draft) => {
        draft.errorMessage = "原文と訳語を入力してください。"
      })
      return
    }

    this.store.update((draft) => {
      draft.errorMessage = ""
    })

    try {
      const refresh = buildRefreshPayload(
        state.query,
        state.category,
        state.page + 1
      )

      if (state.modalState === "create") {
        const response = await this.gateway.createMasterDictionaryEntry({
          payload,
          refresh
        })
        if (response.page) {
          this.applyRefreshPage(
            response.page,
            response.entry,
            response.refreshTargetId
          )
        } else {
          await this.loadEntries(response.refreshTargetId)
        }
        this.store.update((draft) => {
          draft.modalState = null
        })
        return
      }

      if (!state.selectedId) {
        this.store.update((draft) => {
          draft.errorMessage = "更新対象が選択されていません。"
        })
        return
      }

      const response = await this.gateway.updateMasterDictionaryEntry({
        id: state.selectedId,
        payload,
        refresh
      })
      if (response.page) {
        this.applyRefreshPage(
          response.page,
          response.entry,
          response.refreshTargetId
        )
      } else {
        await this.loadEntries(response.refreshTargetId)
      }
      this.store.update((draft) => {
        draft.modalState = null
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "保存に失敗しました。")
      })
    }
  }

  async deleteCurrentEntry(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.selectedId) {
      return
    }

    this.store.update((draft) => {
      draft.errorMessage = ""
    })

    try {
      const response = await this.gateway.deleteMasterDictionaryEntry({
        id: state.selectedId,
        refresh: buildRefreshPayload(
          state.query,
          state.category,
          state.page + 1
        )
      })
      if (response.page) {
        this.applyRefreshPage(response.page, null, response.nextSelectedId)
      } else {
        await this.loadEntries(response.nextSelectedId)
      }
      this.store.update((draft) => {
        draft.modalState = null
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "削除に失敗しました。")
      })
    }
  }

  async startStagedXmlImport(waitForRuntimeCompletion: boolean): Promise<void> {
    const state = this.store.snapshot()
    if (
      !this.gateway ||
      !state.selectedFileReference ||
      state.importStage !== "ready"
    ) {
      return
    }

    this.store.update((draft) => {
      draft.errorMessage = ""
      draft.importStage = "running"
      draft.importProgress = 50
    })

    try {
      const response = await this.gateway.importMasterDictionaryXml({
        filePath: state.selectedFileName,
        fileReference: state.selectedFileReference,
        refresh: buildRefreshPayload("", "すべて", 1)
      })

      if (waitForRuntimeCompletion) {
        return
      }

      await this.handleImportCompleted({
        page: response.page,
        summary: response.summary
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.importStage = "ready"
        draft.importProgress = 0
        draft.importSummary = null
        draft.errorMessage = toErrorMessage(
          error,
          "XML 取り込みに失敗しました。"
        )
      })
    }
  }

  async handleImportCompleted(
    payload: RuntimeImportCompletedPayload
  ): Promise<void> {
    const stateBefore = this.store.snapshot()

    this.store.update((draft) => {
      draft.query = ""
      draft.category = "すべて"
    })

    if (payload.page) {
      this.applyRefreshPage(payload.page, null, null)
    } else {
      const lastEntryId = payload.summary?.lastEntryId
      const fallbackId =
        typeof lastEntryId === "number" && lastEntryId > 0
          ? String(lastEntryId)
          : stateBefore.selectedId
      this.store.update((draft) => {
        draft.page = 0
      })
      await this.loadEntries(fallbackId)
    }

    const nextState = this.store.snapshot()
    this.store.update((draft) => {
      draft.importStage = "done"
      draft.importProgress = 100
      draft.importSummary = {
        fileName: stateBefore.selectedFileName,
        importedCount:
          payload.summary?.importedCount ??
          Math.max(nextState.totalCount - stateBefore.totalCount, 0),
        updatedCount: payload.summary?.updatedCount ?? 0,
        totalCount: payload.page?.totalCount ?? nextState.totalCount,
        selectedSource: nextState.selectedEntry?.source ?? "-"
      }
    })
  }

  handleImportProgress(payload: RuntimeImportProgressPayload): void {
    if (typeof payload.progress !== "number") {
      return
    }

    const progress = payload.progress
    this.store.update((draft) => {
      draft.importStage = "running"
      draft.importProgress = Math.max(0, Math.min(100, Math.floor(progress)))
    })
  }

  private applyRefreshPage(
    pageState: MasterDictionaryPageState,
    preferredDetail: MasterDictionaryEntryDetail | null,
    preferredId: string | null | undefined
  ): void {
    const state = this.store.snapshot()
    const resolvedSelectedId = chooseSelectedIdFromPage(
      pageState,
      preferredId === undefined ? state.selectedId : preferredId
    )

    const selectedPageEntry = resolvedSelectedId
      ? pageState.items.find(
          (candidate) => String(candidate.id) === resolvedSelectedId
        )
      : null

    this.store.update((draft) => {
      draft.entries = pageState.items.map((entry) => ({
        id: String(entry.id),
        source: entry.source,
        translation: entry.translation,
        category: entry.category,
        origin: entry.origin,
        updatedAt: entry.updatedAt
      }))
      draft.totalCount = pageState.totalCount
      draft.page = Math.max(0, pageState.page - 1)
      draft.selectedId = resolvedSelectedId
      if (!resolvedSelectedId) {
        draft.selectedEntry = null
        return
      }

      if (preferredDetail?.id === resolvedSelectedId) {
        draft.selectedEntry = preferredDetail
        return
      }

      draft.selectedEntry = selectedPageEntry
        ? toDetailFromPageEntry(selectedPageEntry)
        : null
    })
  }
}
