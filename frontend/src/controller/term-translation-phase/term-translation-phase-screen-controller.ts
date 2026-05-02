import type {
  TermTranslationPhaseScreenControllerContract,
  TermTranslationPhaseScreenViewModelListener,
  TermTranslationPhaseScreenState,
  TermTranslationPhaseScreenViewModel
} from "@application/contract/term-translation-phase"

interface TermTranslationPhaseStoreLike {
  subscribe(listener: (state: TermTranslationPhaseScreenState) => void): () => void
  snapshot(): TermTranslationPhaseScreenState
}

interface TermTranslationPhasePresenterLike {
  toViewModel(
    state: TermTranslationPhaseScreenState,
    isGatewayConnected: boolean
  ): TermTranslationPhaseScreenViewModel
}

interface TermTranslationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  refresh(): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
}

interface TermTranslationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TermTranslationPhaseStoreLike
  presenter: TermTranslationPhasePresenterLike
  useCase: TermTranslationPhaseUseCaseLike
}

export class TermTranslationPhaseScreenController
  implements TermTranslationPhaseScreenControllerContract
{
  constructor(
    private readonly dependencies: TermTranslationPhaseScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: TermTranslationPhaseScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): TermTranslationPhaseScreenViewModel {
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
}
