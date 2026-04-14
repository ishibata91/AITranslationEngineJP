import {
  DEFAULT_CATEGORY,
  DEFAULT_ORIGIN
} from "@application/contract/master-dictionary"
import type {
  MasterDictionaryScreenControllerContract,
  MasterDictionaryScreenViewModelListener
} from "@application/contract/master-dictionary/master-dictionary-screen-contract"
import type {
  MasterDictionaryScreenState,
  MasterDictionaryScreenViewModel
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

function resolveFileReference(file: File): string {
  const pathRecord = file as File & {
    path?: string
    webkitRelativePath?: string
  }
  const candidates = [pathRecord.path, pathRecord.webkitRelativePath, file.name]
  for (const candidate of candidates) {
    if (typeof candidate === "string" && candidate.trim() !== "") {
      return candidate
    }
  }
  return file.name
}

interface MasterDictionaryStoreLike {
  subscribe(listener: (state: MasterDictionaryScreenState) => void): () => void
  snapshot(): MasterDictionaryScreenState
  update(mutator: (draft: MasterDictionaryScreenState) => void): void
}

interface MasterDictionaryPresenterLike {
  toViewModel(
    state: MasterDictionaryScreenState,
    isGatewayConnected: boolean
  ): MasterDictionaryScreenViewModel
}

interface MasterDictionaryUseCaseLike {
  loadEntries(preferredId?: string | null): Promise<void>
  selectEntry(id: string): Promise<void>
  saveCurrentEntry(): Promise<void>
  deleteCurrentEntry(): Promise<void>
  startStagedXmlImport(waitForRuntimeCompletion: boolean): Promise<void>
}

interface MasterDictionaryRuntimeEventAdapterLike {
  subscribe(): boolean
  detach(): void
}

interface MasterDictionaryScreenControllerDependencies {
  isGatewayConnected: boolean
  store: MasterDictionaryStoreLike
  presenter: MasterDictionaryPresenterLike
  useCase: MasterDictionaryUseCaseLike
  runtimeEventAdapter: MasterDictionaryRuntimeEventAdapterLike
}

export class MasterDictionaryScreenController implements MasterDictionaryScreenControllerContract {
  private runtimeEventSubscribed = false

  constructor(
    private readonly dependencies: MasterDictionaryScreenControllerDependencies
  ) {}

  async mount(): Promise<void> {
    this.runtimeEventSubscribed =
      this.dependencies.runtimeEventAdapter.subscribe()
    await this.dependencies.useCase.loadEntries()
  }

  dispose(): void {
    this.dependencies.runtimeEventAdapter.detach()
  }

  subscribe(listener: MasterDictionaryScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): MasterDictionaryScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  async selectRow(id: string): Promise<void> {
    await this.dependencies.useCase.selectEntry(id)
  }

  openCreateModal(): void {
    this.dependencies.store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = ""
      draft.formCategory = DEFAULT_CATEGORY
      draft.formOrigin = DEFAULT_ORIGIN
      draft.formTranslation = ""
      draft.errorMessage = ""
    })
  }

  openEditModal(): void {
    const state = this.dependencies.store.snapshot()
    if (!state.selectedEntry) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.modalState = "edit"
      draft.formSource = state.selectedEntry?.source ?? ""
      draft.formCategory = state.selectedEntry?.category ?? DEFAULT_CATEGORY
      draft.formOrigin = state.selectedEntry?.origin ?? DEFAULT_ORIGIN
      draft.formTranslation = state.selectedEntry?.translation ?? ""
      draft.errorMessage = ""
    })
  }

  openDeleteModal(): void {
    const state = this.dependencies.store.snapshot()
    if (!state.selectedEntry) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.modalState = "delete"
      draft.errorMessage = ""
    })
  }

  closeEditModal(): void {
    this.dependencies.store.update((draft) => {
      if (draft.modalState === "create" || draft.modalState === "edit") {
        draft.modalState = null
      }
    })
  }

  closeDeleteModal(): void {
    this.dependencies.store.update((draft) => {
      if (draft.modalState === "delete") {
        draft.modalState = null
      }
    })
  }

  async saveCurrentEntry(): Promise<void> {
    await this.dependencies.useCase.saveCurrentEntry()
  }

  async deleteCurrentEntry(): Promise<void> {
    await this.dependencies.useCase.deleteCurrentEntry()
  }

  handleSearchInput(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.query = target.value
      draft.page = 0
      draft.errorMessage = ""
    })

    void this.dependencies.useCase.loadEntries()
  }

  handleCategoryChange(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.category = target.value
      draft.page = 0
      draft.errorMessage = ""
    })

    void this.dependencies.useCase.loadEntries()
  }

  goToPrevPage(): void {
    const state = this.dependencies.store.snapshot()
    if (state.page <= 0) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.page -= 1
    })
    void this.dependencies.useCase.loadEntries()
  }

  goToNextPage(): void {
    const state = this.dependencies.store.snapshot()
    const totalPages = Math.max(1, Math.ceil(state.totalCount / 30))
    if (state.page + 1 >= totalPages) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.page += 1
    })
    void this.dependencies.useCase.loadEntries()
  }

  stageXmlImport(file: File | null): void {
    this.dependencies.store.update((draft) => {
      draft.selectedFileName = file ? file.name : "未選択"
      draft.selectedFileReference = file ? resolveFileReference(file) : null
      draft.importStage = file ? "ready" : "idle"
      draft.importProgress = 0
      draft.importSummary = null
      draft.errorMessage = ""
    })
  }

  resetImportSelection(): void {
    const state = this.dependencies.store.snapshot()
    if (state.importStage === "running") {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.selectedFileName = "未選択"
      draft.selectedFileReference = null
      draft.importStage = "idle"
      draft.importProgress = 0
      draft.importSummary = null
    })
  }

  async startImport(): Promise<void> {
    await this.dependencies.useCase.startStagedXmlImport(
      this.runtimeEventSubscribed
    )
  }

  setFormSource(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.formSource = target.value
    })
  }

  setFormCategory(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.formCategory = target.value
    })
  }

  setFormOrigin(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.formOrigin = target.value
    })
  }

  setFormTranslation(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLTextAreaElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.formTranslation = target.value
    })
  }
}
