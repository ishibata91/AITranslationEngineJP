import type {
  BodyTranslationPhaseScreenControllerContract,
  BodyTranslationPhaseScreenState,
  BodyTranslationPhaseScreenViewModel,
  BodyTranslationPhaseScreenViewModelListener
} from "@application/contract/body-translation-phase"

interface BodyTranslationPhaseStoreLike {
  subscribe(
    listener: (state: BodyTranslationPhaseScreenState) => void
  ): () => void
  snapshot(): BodyTranslationPhaseScreenState
}

interface BodySelectOption {
  value: string
  label: string
}

interface BodyTranslationPhasePresenterLike {
  toViewModel(
    state: BodyTranslationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders?: BodySelectOption[],
    availableModels?: BodySelectOption[]
  ): BodyTranslationPhaseScreenViewModel
}

interface BodyTranslationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?: (
    searchQuery: string,
    phase?: string
  ) => Promise<void>
  setProcessingTargetPage?: (page: number, phase?: string) => Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkOutputReadiness(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
}

interface BodyTranslationPhaseGatewayLike {
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

interface BodyTranslationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: BodyTranslationPhaseStoreLike
  presenter: BodyTranslationPhasePresenterLike
  useCase: BodyTranslationPhaseUseCaseLike
  gateway?: BodyTranslationPhaseGatewayLike | null
}

export class BodyTranslationPhaseScreenController implements BodyTranslationPhaseScreenControllerContract {
  private availableProviders: BodySelectOption[] = []
  private availableModels: BodySelectOption[] = []
  private modelListRequestToken = 0
  private readonly listeners: Set<BodyTranslationPhaseScreenViewModelListener> = new Set()

  constructor(
    private readonly dependencies: BodyTranslationPhaseScreenControllerDependencies
  ) {}

  private notifyListeners(): void {
    const vm = this.getViewModel()
    for (const listener of this.listeners) {
      listener(vm)
    }
  }

  setAvailableProviders(providers: BodySelectOption[]): void {
    this.availableProviders = providers
  }

  setAvailableModels(models: BodySelectOption[]): void {
    this.availableModels = models
  }

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: BodyTranslationPhaseScreenViewModelListener): () => void {
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

  getViewModel(): BodyTranslationPhaseScreenViewModel {
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

  async setProcessingTargetSearchQuery(
    searchQuery: string,
    phase?: string
  ): Promise<void> {
    await this.dependencies.useCase.setProcessingTargetSearchQuery?.(
      searchQuery,
      phase
    )
  }

  async setProcessingTargetPage(page: number, phase?: string): Promise<void> {
    await this.dependencies.useCase.setProcessingTargetPage?.(page, phase)
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

  async checkOutputReadiness(): Promise<void> {
    await this.dependencies.useCase.checkOutputReadiness()
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
