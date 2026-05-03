import type {
  TranslationJobSetupScreenControllerContract,
  TranslationJobSetupScreenViewModelListener
} from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
import type {
  TranslationJobSetupPhaseId,
  TranslationJobSetupScreenState,
  TranslationJobSetupScreenViewModel
} from "@application/gateway-contract/translation-job-setup"

interface TranslationJobSetupStoreLike {
  subscribe(
    listener: (state: TranslationJobSetupScreenState) => void
  ): () => void
  snapshot(): TranslationJobSetupScreenState
}

interface TranslationJobSetupPresenterLike {
  toViewModel(
    state: TranslationJobSetupScreenState,
    isGatewayConnected: boolean
  ): TranslationJobSetupScreenViewModel
}

interface TranslationJobSetupUseCaseLike {
  load(): Promise<void>
  selectInputSource(inputSourceId: number): void
  selectRuntime(runtimeKey: string): void
  selectCredentialRef(credentialRef: string): void
  selectPhaseProvider(phaseId: TranslationJobSetupPhaseId, provider: string): void
  refreshPhaseModels(phaseId: TranslationJobSetupPhaseId): Promise<void>
  selectPhaseModel(phaseId: TranslationJobSetupPhaseId, model: string): void
  togglePhaseBatchMode(phaseId: TranslationJobSetupPhaseId, enabled: boolean): void
  acknowledgeCredentialConfigured(phaseId: TranslationJobSetupPhaseId): void
  savePhaseCredential(
    phaseId: TranslationJobSetupPhaseId,
    apiKey: string
  ): Promise<void>
  runValidation(): Promise<void>
  createJob(): Promise<void>
}

interface TranslationJobSetupScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TranslationJobSetupStoreLike
  presenter: TranslationJobSetupPresenterLike
  useCase: TranslationJobSetupUseCaseLike
}

export class TranslationJobSetupScreenController implements TranslationJobSetupScreenControllerContract {
  constructor(
    private readonly dependencies: TranslationJobSetupScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: TranslationJobSetupScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): TranslationJobSetupScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  selectInputSource(inputSourceId: number): void {
    this.dependencies.useCase.selectInputSource(inputSourceId)
  }

  selectRuntime(runtimeKey: string): void {
    this.dependencies.useCase.selectRuntime(runtimeKey)
  }

  selectCredentialRef(credentialRef: string): void {
    this.dependencies.useCase.selectCredentialRef(credentialRef)
  }

  selectPhaseProvider(phaseId: TranslationJobSetupPhaseId, provider: string): void {
    this.dependencies.useCase.selectPhaseProvider(phaseId, provider)
  }

  async refreshPhaseModels(phaseId: TranslationJobSetupPhaseId): Promise<void> {
    await this.dependencies.useCase.refreshPhaseModels(phaseId)
  }

  selectPhaseModel(phaseId: TranslationJobSetupPhaseId, model: string): void {
    this.dependencies.useCase.selectPhaseModel(phaseId, model)
  }

  togglePhaseBatchMode(phaseId: TranslationJobSetupPhaseId, enabled: boolean): void {
    this.dependencies.useCase.togglePhaseBatchMode(phaseId, enabled)
  }

  acknowledgeCredentialConfigured(phaseId: TranslationJobSetupPhaseId): void {
    this.dependencies.useCase.acknowledgeCredentialConfigured(phaseId)
  }

  async savePhaseCredential(
    phaseId: TranslationJobSetupPhaseId,
    apiKey: string
  ): Promise<void> {
    await this.dependencies.useCase.savePhaseCredential(phaseId, apiKey)
  }

  async runValidation(): Promise<void> {
    await this.dependencies.useCase.runValidation()
  }

  async createJob(): Promise<void> {
    await this.dependencies.useCase.createJob()
  }
}
