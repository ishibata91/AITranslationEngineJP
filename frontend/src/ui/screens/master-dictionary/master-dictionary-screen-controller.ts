import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"

import { MasterDictionaryPresenter } from "./master-dictionary.presenter"
import {
  MasterDictionaryRuntimeEventAdapter
} from "./master-dictionary-runtime-event-adapter"
import { MasterDictionaryStore } from "./master-dictionary.store"
import {
  DEFAULT_CATEGORY,
  DEFAULT_ORIGIN
} from "./master-dictionary-screen-constants"
import type { MasterDictionaryScreenViewModel } from "./master-dictionary-screen-types"
import { MasterDictionaryUseCase } from "./master-dictionary.usecase"

type ViewModelListener = (viewModel: MasterDictionaryScreenViewModel) => void

function resolveFileReference(file: File): string {
  const pathRecord = file as File & { path?: string; webkitRelativePath?: string }
  return pathRecord.path ?? pathRecord.webkitRelativePath ?? file.name
}

class MasterDictionaryScreenController {
  private gateway: MasterDictionaryGatewayContract | null

  private readonly store = new MasterDictionaryStore()

  private readonly presenter = new MasterDictionaryPresenter()

  private readonly useCase: MasterDictionaryUseCase

  private readonly runtimeEventAdapter: MasterDictionaryRuntimeEventAdapter

  private runtimeEventSubscribed = false

  constructor(gateway: MasterDictionaryGatewayContract | null) {
    this.gateway = gateway
    this.useCase = new MasterDictionaryUseCase(gateway, this.store)
    this.runtimeEventAdapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress: (payload) => {
        this.useCase.handleImportProgress(payload)
      },
      onImportCompleted: (payload) => {
        void this.useCase.handleImportCompleted(payload)
      }
    })
  }

  async mount(): Promise<void> {
    this.runtimeEventSubscribed = this.runtimeEventAdapter.subscribe()
    await this.useCase.loadEntries()
  }

  dispose(): void {
    this.runtimeEventAdapter.detach()
  }

  updateGateway(gateway: MasterDictionaryGatewayContract | null): void {
    if (this.gateway === gateway) {
      return
    }

    this.gateway = gateway
    this.useCase.setGateway(gateway)
    void this.useCase.loadEntries()
  }

  subscribe(listener: ViewModelListener): () => void {
    return this.store.subscribe((state) => {
      listener(this.presenter.toViewModel(state, Boolean(this.gateway)))
    })
  }

  getViewModel(): MasterDictionaryScreenViewModel {
    return this.presenter.toViewModel(this.store.snapshot(), Boolean(this.gateway))
  }

  async selectRow(id: string): Promise<void> {
    await this.useCase.selectEntry(id)
  }

  openCreateModal(): void {
    this.store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = ""
      draft.formCategory = DEFAULT_CATEGORY
      draft.formOrigin = DEFAULT_ORIGIN
      draft.formTranslation = ""
      draft.errorMessage = ""
    })
  }

  openEditModal(): void {
    const state = this.store.snapshot()
    if (!state.selectedEntry) {
      return
    }

    this.store.update((draft) => {
      draft.modalState = "edit"
      draft.formSource = state.selectedEntry?.source ?? ""
      draft.formCategory = state.selectedEntry?.category ?? DEFAULT_CATEGORY
      draft.formOrigin = state.selectedEntry?.origin ?? DEFAULT_ORIGIN
      draft.formTranslation = state.selectedEntry?.translation ?? ""
      draft.errorMessage = ""
    })
  }

  openDeleteModal(): void {
    const state = this.store.snapshot()
    if (!state.selectedEntry) {
      return
    }

    this.store.update((draft) => {
      draft.modalState = "delete"
      draft.errorMessage = ""
    })
  }

  closeEditModal(): void {
    this.store.update((draft) => {
      if (draft.modalState === "create" || draft.modalState === "edit") {
        draft.modalState = null
      }
    })
  }

  closeDeleteModal(): void {
    this.store.update((draft) => {
      if (draft.modalState === "delete") {
        draft.modalState = null
      }
    })
  }

  async saveCurrentEntry(): Promise<void> {
    await this.useCase.saveCurrentEntry()
  }

  async deleteCurrentEntry(): Promise<void> {
    await this.useCase.deleteCurrentEntry()
  }

  handleSearchInput(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    this.store.update((draft) => {
      draft.query = target.value
      draft.page = 0
      draft.errorMessage = ""
    })

    void this.useCase.loadEntries()
  }

  handleCategoryChange(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }

    this.store.update((draft) => {
      draft.category = target.value
      draft.page = 0
      draft.errorMessage = ""
    })

    void this.useCase.loadEntries()
  }

  goToPrevPage(): void {
    const state = this.store.snapshot()
    if (state.page <= 0) {
      return
    }

    this.store.update((draft) => {
      draft.page -= 1
    })
    void this.useCase.loadEntries()
  }

  goToNextPage(): void {
    const state = this.store.snapshot()
    const totalPages = Math.max(1, Math.ceil(state.totalCount / 30))
    if (state.page + 1 >= totalPages) {
      return
    }

    this.store.update((draft) => {
      draft.page += 1
    })
    void this.useCase.loadEntries()
  }

  stageXmlImport(file: File | null): void {
    this.store.update((draft) => {
      draft.selectedFileName = file ? file.name : "未選択"
      draft.selectedFileReference = file ? resolveFileReference(file) : null
      draft.importStage = file ? "ready" : "idle"
      draft.importProgress = file ? 12 : 0
      draft.importSummary = null
      draft.errorMessage = ""
    })
  }

  resetImportSelection(): void {
    const state = this.store.snapshot()
    if (state.importStage === "running") {
      return
    }

    this.store.update((draft) => {
      draft.selectedFileName = "未選択"
      draft.selectedFileReference = null
      draft.importStage = "idle"
      draft.importProgress = 0
      draft.importSummary = null
    })
  }

  async startImport(): Promise<void> {
    await this.useCase.startStagedXmlImport(this.runtimeEventSubscribed)
  }

  setFormSource(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    this.store.update((draft) => {
      draft.formSource = target.value
    })
  }

  setFormCategory(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }
    this.store.update((draft) => {
      draft.formCategory = target.value
    })
  }

  setFormOrigin(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }
    this.store.update((draft) => {
      draft.formOrigin = target.value
    })
  }

  setFormTranslation(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLTextAreaElement)) {
      return
    }
    this.store.update((draft) => {
      draft.formTranslation = target.value
    })
  }
}

export function createMasterDictionaryScreenController(
  gateway: MasterDictionaryGatewayContract | null
): MasterDictionaryScreenController {
  return new MasterDictionaryScreenController(gateway)
}
