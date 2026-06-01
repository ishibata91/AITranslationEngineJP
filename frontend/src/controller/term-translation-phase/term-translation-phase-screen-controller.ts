import type {
  TermTranslationPhaseScreenControllerContract,
  TermTranslationPhaseScreenViewModelListener,
  TermTranslationPhaseScreenState,
  TermTranslationPhaseScreenViewModel
} from "@application/contract/term-translation-phase"

interface TermTranslationPhaseStoreLike {
  subscribe(
    listener: (state: TermTranslationPhaseScreenState) => void
  ): () => void
  snapshot(): TermTranslationPhaseScreenState
}

interface SelectOption {
  value: string
  label: string
}

interface TermTranslationPhasePresenterLike {
  toViewModel(
    state: TermTranslationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders?: SelectOption[],
    availableModels?: SelectOption[]
  ): TermTranslationPhaseScreenViewModel
}

interface TermTranslationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?: (searchQuery: string) => Promise<void>
  setProcessingTargetPage?: (page: number) => Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
}

interface TermTranslationPhaseGatewayLike {
  listProviderModels?(request: {
    provider: string
    credentialStatus: string
    requestToken: string
  }): Promise<{
    provider: string
    status: string
    models: { modelId: string; label: string }[]
    failureKind?: string
  }>
}

interface TermTranslationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TermTranslationPhaseStoreLike
  presenter: TermTranslationPhasePresenterLike
  useCase: TermTranslationPhaseUseCaseLike
  availableProviders?: SelectOption[]
  gateway?: TermTranslationPhaseGatewayLike | null
}

export class TermTranslationPhaseScreenController implements TermTranslationPhaseScreenControllerContract {
  private availableProviders: SelectOption[]
  private availableModels: SelectOption[]
  private modelListRequestToken = 0
  private readonly listeners: Set<TermTranslationPhaseScreenViewModelListener> = new Set()

  constructor(
    private readonly dependencies: TermTranslationPhaseScreenControllerDependencies
  ) {
    this.availableProviders = dependencies.availableProviders ?? []
    this.availableModels = []
  }

  private notifyListeners(): void {
    const vm = this.getViewModel()
    for (const listener of this.listeners) {
      listener(vm)
    }
  }

  setAvailableProviders(providers: SelectOption[]): void {
    this.availableProviders = providers
  }

  setAvailableModels(models: SelectOption[]): void {
    this.availableModels = models
  }

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: TermTranslationPhaseScreenViewModelListener): () => void {
    this.listeners.add(listener)
    const unsubStore = this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected,
          this.availableProviders,
          this.availableModels
        )
      )
    })
    return () => {
      this.listeners.delete(listener)
      unsubStore()
    }
  }

  getViewModel(): TermTranslationPhaseScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected,
      this.availableProviders,
      this.availableModels
    )
  }

  async refreshModelList(provider: string): Promise<void> {
    const gateway = this.dependencies.gateway
    if (!gateway?.listProviderModels) {
      return
    }
    if (!provider) {
      return
    }
    const token = String(++this.modelListRequestToken)
    try {
      const response = await gateway.listProviderModels({
        provider,
        credentialStatus: "configured",
        requestToken: token
      })
      if (response.provider === provider) {
        this.availableModels = response.models.map((m) => ({
          value: m.modelId,
          label: m.label
        }))
        this.notifyListeners()
      }
    } catch {
      // モデル一覧取得失敗時は既存の選択肢を維持する
    }
  }

  async setJobId(jobId: number | null): Promise<void> {
    await this.dependencies.useCase.setJobId(jobId)
  }

  async setProcessingTargetSearchQuery(searchQuery: string): Promise<void> {
    await this.dependencies.useCase.setProcessingTargetSearchQuery?.(
      searchQuery
    )
  }

  async setProcessingTargetPage(page: number): Promise<void> {
    await this.dependencies.useCase.setProcessingTargetPage?.(page)
  }

  async startPhase(): Promise<void> {
    await this.dependencies.useCase.startPhase()
  }

  async pausePhase(): Promise<void> {
    await this.dependencies.useCase.pausePhase()
  }

  async resumePhase(): Promise<void> {
    await this.dependencies.useCase.resumePhase()
  }

  async retryPhase(): Promise<void> {
    await this.dependencies.useCase.retryPhase()
  }

  async saveAISettings(request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }): Promise<void> {
    await this.dependencies.useCase.saveAISettings?.(request)
  }
}
