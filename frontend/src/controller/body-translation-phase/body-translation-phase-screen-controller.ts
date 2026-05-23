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

interface BodyTranslationPhasePresenterLike {
  toViewModel(
    state: BodyTranslationPhaseScreenState,
    isGatewayConnected: boolean
  ): BodyTranslationPhaseScreenViewModel
}

interface BodyTranslationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  refresh(): Promise<void>
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

interface BodyTranslationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: BodyTranslationPhaseStoreLike
  presenter: BodyTranslationPhasePresenterLike
  useCase: BodyTranslationPhaseUseCaseLike
}

export class BodyTranslationPhaseScreenController implements BodyTranslationPhaseScreenControllerContract {
  constructor(
    private readonly dependencies: BodyTranslationPhaseScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: BodyTranslationPhaseScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): BodyTranslationPhaseScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  async setJobId(jobId: number | null): Promise<void> {
    await this.dependencies.useCase.setJobId(jobId)
  }

  async refresh(): Promise<void> {
    await this.dependencies.useCase.refresh()
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
