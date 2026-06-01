import type {
  PersonaGenerationPhaseScreenControllerContract,
  PersonaGenerationPhaseScreenState,
  PersonaGenerationPhaseScreenViewModel,
  PersonaGenerationPhaseScreenViewModelListener
} from "@application/contract/persona-generation-phase"

interface PersonaGenerationPhaseStoreLike {
  subscribe(
    listener: (state: PersonaGenerationPhaseScreenState) => void
  ): () => void
  snapshot(): PersonaGenerationPhaseScreenState
}

interface PersonaSelectOption {
  value: string
  label: string
}

interface PersonaGenerationPhasePresenterLike {
  toViewModel(
    state: PersonaGenerationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders?: PersonaSelectOption[],
    availableModels?: PersonaSelectOption[]
  ): PersonaGenerationPhaseScreenViewModel
}

interface PersonaGenerationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?: (searchQuery: string) => Promise<void>
  setProcessingTargetPage?: (page: number) => Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkBodyReadiness(): Promise<void>
  startBodyPhase(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
}

interface PersonaGenerationPhaseGatewayLike {
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

interface PersonaGenerationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: PersonaGenerationPhaseStoreLike
  presenter: PersonaGenerationPhasePresenterLike
  useCase: PersonaGenerationPhaseUseCaseLike
  gateway?: PersonaGenerationPhaseGatewayLike | null
}

export class PersonaGenerationPhaseScreenController implements PersonaGenerationPhaseScreenControllerContract {
  private availableProviders: PersonaSelectOption[] = []
  private availableModels: PersonaSelectOption[] = []
  private modelListRequestToken = 0
  private readonly listeners: Set<PersonaGenerationPhaseScreenViewModelListener> = new Set()

  constructor(
    private readonly dependencies: PersonaGenerationPhaseScreenControllerDependencies
  ) {}

  private notifyListeners(): void {
    const vm = this.getViewModel()
    for (const listener of this.listeners) {
      listener(vm)
    }
  }

  setAvailableProviders(providers: PersonaSelectOption[]): void {
    this.availableProviders = providers
  }

  setAvailableModels(models: PersonaSelectOption[]): void {
    this.availableModels = models
  }

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(
    listener: PersonaGenerationPhaseScreenViewModelListener
  ): () => void {
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

  getViewModel(): PersonaGenerationPhaseScreenViewModel {
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

  async cancelPhase(): Promise<void> {
    await this.dependencies.useCase.cancelPhase()
  }

  async checkBodyReadiness(): Promise<void> {
    await this.dependencies.useCase.checkBodyReadiness()
  }

  async startBodyPhase(): Promise<void> {
    await this.dependencies.useCase.startBodyPhase()
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
