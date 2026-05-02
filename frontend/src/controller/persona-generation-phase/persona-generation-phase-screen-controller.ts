import type {
  PersonaGenerationPhaseScreenControllerContract,
  PersonaGenerationPhaseScreenState,
  PersonaGenerationPhaseScreenViewModel,
  PersonaGenerationPhaseScreenViewModelListener
} from "@application/contract/persona-generation-phase"

interface PersonaGenerationPhaseStoreLike {
  subscribe(listener: (state: PersonaGenerationPhaseScreenState) => void): () => void
  snapshot(): PersonaGenerationPhaseScreenState
}

interface PersonaGenerationPhasePresenterLike {
  toViewModel(
    state: PersonaGenerationPhaseScreenState,
    isGatewayConnected: boolean
  ): PersonaGenerationPhaseScreenViewModel
}

interface PersonaGenerationPhaseUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  refresh(): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkBodyReadiness(): Promise<void>
  startBodyPhase(): Promise<void>
}

interface PersonaGenerationPhaseScreenControllerDependencies {
  isGatewayConnected: boolean
  store: PersonaGenerationPhaseStoreLike
  presenter: PersonaGenerationPhasePresenterLike
  useCase: PersonaGenerationPhaseUseCaseLike
}

export class PersonaGenerationPhaseScreenController
  implements PersonaGenerationPhaseScreenControllerContract
{
  constructor(
    private readonly dependencies: PersonaGenerationPhaseScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: PersonaGenerationPhaseScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): PersonaGenerationPhaseScreenViewModel {
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

  async checkBodyReadiness(): Promise<void> {
    await this.dependencies.useCase.checkBodyReadiness()
  }

  async startBodyPhase(): Promise<void> {
    await this.dependencies.useCase.startBodyPhase()
  }
}
