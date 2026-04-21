import type {
  MasterPersonaEditableFieldMap,
  MasterPersonaScreenControllerContract,
  MasterPersonaScreenViewModelListener
} from "@application/contract/master-persona/master-persona-screen-contract"
import type {
  MasterPersonaScreenState,
  MasterPersonaScreenViewModel
} from "@application/gateway-contract/master-persona"

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

interface MasterPersonaStoreLike {
  subscribe(listener: (state: MasterPersonaScreenState) => void): () => void
  snapshot(): MasterPersonaScreenState
  update(mutator: (draft: MasterPersonaScreenState) => void): void
}

interface MasterPersonaPresenterLike {
  toViewModel(
    state: MasterPersonaScreenState,
    isGatewayConnected: boolean
  ): MasterPersonaScreenViewModel
}

interface MasterPersonaUseCaseLike {
  loadScreen(): Promise<void>
  loadPage(preferredIdentityKey?: string | null): Promise<void>
  selectEntry(identityKey: string): Promise<void>
  previewGeneration(): Promise<void>
  executeGeneration(): Promise<void>
  loadRunStatus(): Promise<void>
  interruptGeneration(): Promise<void>
  cancelGeneration(): Promise<void>
  saveAISettings(): Promise<void>
  saveCurrentEntry(): Promise<void>
  deleteCurrentEntry(): Promise<void>
  setModalState(modalState: "edit" | "delete" | null): void
}

interface MasterPersonaRuntimePollingAdapterLike {
  start(onTick: () => void): boolean
  stop(): void
}

interface MasterPersonaScreenControllerDependencies {
  isGatewayConnected: boolean
  store: MasterPersonaStoreLike
  presenter: MasterPersonaPresenterLike
  useCase: MasterPersonaUseCaseLike
  runtimePollingAdapter: MasterPersonaRuntimePollingAdapterLike
}

export class MasterPersonaScreenController
  implements MasterPersonaScreenControllerContract
{
  constructor(
    private readonly dependencies: MasterPersonaScreenControllerDependencies
  ) {}

  async mount(): Promise<void> {
    this.dependencies.runtimePollingAdapter.start(() => {
      void this.dependencies.useCase.loadRunStatus()
    })
    await this.dependencies.useCase.loadScreen()
  }

  dispose(): void {
    this.dependencies.runtimePollingAdapter.stop()
  }

  subscribe(listener: MasterPersonaScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): MasterPersonaScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  async selectRow(identityKey: string): Promise<void> {
    await this.dependencies.useCase.selectEntry(identityKey)
  }

  handleSearchInput(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.keyword = target.value
      draft.page = 1
      draft.errorMessage = ""
    })
    void this.dependencies.useCase.loadPage()
  }

  handlePluginFilterChange(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }

    this.dependencies.store.update((draft) => {
      draft.pluginFilter = target.value
      draft.page = 1
      draft.errorMessage = ""
    })
    void this.dependencies.useCase.loadPage()
  }

  goToPrevPage(): void {
    const state = this.dependencies.store.snapshot()
    if (state.page <= 1) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.page -= 1
    })
    void this.dependencies.useCase.loadPage()
  }

  goToNextPage(): void {
    const state = this.dependencies.store.snapshot()
    const totalPages = Math.max(1, Math.ceil(state.totalCount / state.pageSize))
    if (state.page >= totalPages) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.page += 1
    })
    void this.dependencies.useCase.loadPage()
  }

  stageJsonSelection(file: File | null): void {
    this.dependencies.store.update((draft) => {
      draft.selectedFileName = file ? file.name : "未選択"
      draft.selectedFileReference = file ? resolveFileReference(file) : null
      draft.preview = null
      draft.errorMessage = ""
      if (!file) {
        draft.runStatus = {
          ...draft.runStatus,
          runState: "入力待ち",
          message: "入力ファイルを選ぶと状態を表示します。"
        }
      }
    })

    if (!file) {
      return
    }

    void this.dependencies.useCase.previewGeneration()
  }

  resetJsonSelection(): void {
    this.dependencies.store.update((draft) => {
      draft.selectedFileName = "未選択"
      draft.selectedFileReference = null
      draft.preview = null
      draft.errorMessage = ""
    })
  }

  async previewGeneration(): Promise<void> {
    await this.dependencies.useCase.previewGeneration()
  }

  async executeGeneration(): Promise<void> {
    await this.dependencies.useCase.executeGeneration()
  }

  async interruptGeneration(): Promise<void> {
    await this.dependencies.useCase.interruptGeneration()
  }

  async cancelGeneration(): Promise<void> {
    await this.dependencies.useCase.cancelGeneration()
  }

  async saveAISettings(): Promise<void> {
    await this.dependencies.useCase.saveAISettings()
  }

  setAIProvider(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLSelectElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.aiSettings.provider = target.value
      draft.aiSettingsMessage = ""
    })
  }

  setAIModel(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.aiSettings.model = target.value
      draft.aiSettingsMessage = ""
    })
  }

  setAPIKey(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.aiSettings.apiKey = target.value
      draft.aiSettingsMessage = ""
    })
  }

  openEditModal(): void {
    const state = this.dependencies.store.snapshot()
    if (!state.selectedEntry) {
      return
    }
    this.dependencies.useCase.setModalState("edit")
  }

  closeEditModal(): void {
    this.dependencies.useCase.setModalState(null)
  }

  openDeleteModal(): void {
    const state = this.dependencies.store.snapshot()
    if (!state.selectedEntry) {
      return
    }
    this.dependencies.useCase.setModalState("delete")
  }

  closeDeleteModal(): void {
    this.dependencies.useCase.setModalState(null)
  }

  async saveCurrentEntry(): Promise<void> {
    await this.dependencies.useCase.saveCurrentEntry()
  }

  async deleteCurrentEntry(): Promise<void> {
    await this.dependencies.useCase.deleteCurrentEntry()
  }

  setEditFormField(
    field: keyof MasterPersonaEditableFieldMap,
    event: Event
  ): void {
    const target = event.currentTarget
    if (
      !(target instanceof HTMLInputElement) &&
      !(target instanceof HTMLTextAreaElement)
    ) {
      return
    }
    this.dependencies.store.update((draft) => {
      draft.editForm[field] = target.value
    })
  }
}
