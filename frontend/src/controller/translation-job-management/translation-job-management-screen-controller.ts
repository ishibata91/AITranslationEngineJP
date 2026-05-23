import type {
  TranslationJobManagementFilterId,
  TranslationJobManagementScreenState,
  TranslationJobManagementScreenViewModel
} from "@application/contract/translation-job-management/translation-job-management-screen-types"
import type {
  TranslationJobManagementScreenControllerContract,
  TranslationJobManagementScreenViewModelListener
} from "@application/contract/translation-job-management/translation-job-management-screen-contract"

interface TranslationJobManagementStoreLike {
  subscribe(
    listener: (state: TranslationJobManagementScreenState) => void
  ): () => void
  snapshot(): TranslationJobManagementScreenState
}

interface TranslationJobManagementPresenterLike {
  toViewModel(
    state: TranslationJobManagementScreenState,
    isGatewayConnected: boolean
  ): TranslationJobManagementScreenViewModel
}

interface TranslationJobManagementUseCaseLike {
  load(): Promise<void>
  reload(): Promise<void>
  selectJob(jobId: number): Promise<void>
  setFilter(filterId: TranslationJobManagementFilterId): void
  setSearchQuery(searchQuery: string): void
  requestStop(): Promise<void>
  requestResume(): Promise<void>
  openDeleteConfirmation(): void
  closeDeleteConfirmation(): void
  deleteSelectedJob(): Promise<void>
}

interface TranslationJobManagementScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TranslationJobManagementStoreLike
  presenter: TranslationJobManagementPresenterLike
  useCase: TranslationJobManagementUseCaseLike
}

export class TranslationJobManagementScreenController implements TranslationJobManagementScreenControllerContract {
  constructor(
    private readonly dependencies: TranslationJobManagementScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(
    listener: TranslationJobManagementScreenViewModelListener
  ): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): TranslationJobManagementScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  async reload(): Promise<void> {
    await this.dependencies.useCase.reload()
  }

  async selectJob(jobId: number): Promise<void> {
    await this.dependencies.useCase.selectJob(jobId)
  }

  setFilter(filterId: TranslationJobManagementFilterId): void {
    this.dependencies.useCase.setFilter(filterId)
  }

  setSearchQuery(searchQuery: string): void {
    this.dependencies.useCase.setSearchQuery(searchQuery)
  }

  async requestStop(): Promise<void> {
    await this.dependencies.useCase.requestStop()
  }

  async requestResume(): Promise<void> {
    await this.dependencies.useCase.requestResume()
  }

  openDeleteConfirmation(): void {
    this.dependencies.useCase.openDeleteConfirmation()
  }

  closeDeleteConfirmation(): void {
    this.dependencies.useCase.closeDeleteConfirmation()
  }

  async deleteSelectedJob(): Promise<void> {
    await this.dependencies.useCase.deleteSelectedJob()
  }
}
