import type {
  TranslationJobSetupScreenControllerContract,
  TranslationJobSetupScreenViewModelListener
} from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
import type {
  TranslationJobSetupScreenState,
  TranslationJobSetupScreenViewModel
} from "@application/gateway-contract/translation-job-setup"

interface TranslationJobSetupStoreLike {
  subscribe(listener: (state: TranslationJobSetupScreenState) => void): () => void
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
  runValidation(): Promise<void>
  createJob(): Promise<void>
}

interface TranslationJobSetupScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TranslationJobSetupStoreLike
  presenter: TranslationJobSetupPresenterLike
  useCase: TranslationJobSetupUseCaseLike
}

export class TranslationJobSetupScreenController
  implements TranslationJobSetupScreenControllerContract
{
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

  async runValidation(): Promise<void> {
    await this.dependencies.useCase.runValidation()
  }

  async createJob(): Promise<void> {
    await this.dependencies.useCase.createJob()
  }
}