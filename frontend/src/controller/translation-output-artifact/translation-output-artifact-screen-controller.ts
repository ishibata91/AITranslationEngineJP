import type {
  TranslationOutputArtifactScreenControllerContract,
  TranslationOutputArtifactScreenViewModelListener
} from "@application/contract/translation-output-artifact"
import type {
  TranslationOutputArtifactScreenState,
  TranslationOutputArtifactScreenViewModel
} from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"

interface TranslationOutputArtifactStoreLike {
  subscribe(
    listener: (state: TranslationOutputArtifactScreenState) => void
  ): () => void
  snapshot(): TranslationOutputArtifactScreenState
}

interface TranslationOutputArtifactPresenterLike {
  toViewModel(
    state: TranslationOutputArtifactScreenState,
    isGatewayConnected: boolean
  ): TranslationOutputArtifactScreenViewModel
}

interface TranslationOutputArtifactUseCaseLike {
  load(): Promise<void>
  setJobId(jobId: number | null): Promise<void>
  setArtifactId(artifactId: number | null): Promise<void>
  setTargetGame(targetGame: string): void
  setOutputPath(outputPath: string): void
  refresh(): Promise<void>
  generateArtifact(): Promise<void>
  regenerateArtifact(): Promise<void>
}

interface TranslationOutputArtifactScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TranslationOutputArtifactStoreLike
  presenter: TranslationOutputArtifactPresenterLike
  useCase: TranslationOutputArtifactUseCaseLike
}

export class TranslationOutputArtifactScreenController implements TranslationOutputArtifactScreenControllerContract {
  constructor(
    private readonly dependencies: TranslationOutputArtifactScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(
    listener: TranslationOutputArtifactScreenViewModelListener
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

  getViewModel(): TranslationOutputArtifactScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  async setJobId(jobId: number | null): Promise<void> {
    await this.dependencies.useCase.setJobId(jobId)
  }

  async setArtifactId(artifactId: number | null): Promise<void> {
    await this.dependencies.useCase.setArtifactId(artifactId)
  }

  setTargetGame(targetGame: string): void {
    this.dependencies.useCase.setTargetGame(targetGame)
  }

  setOutputPath(outputPath: string): void {
    this.dependencies.useCase.setOutputPath(outputPath)
  }

  async refresh(): Promise<void> {
    await this.dependencies.useCase.refresh()
  }

  async generateArtifact(): Promise<void> {
    await this.dependencies.useCase.generateArtifact()
  }

  async regenerateArtifact(): Promise<void> {
    await this.dependencies.useCase.regenerateArtifact()
  }
}
