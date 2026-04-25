import type {
  TranslationInputScreenControllerContract,
  TranslationInputScreenViewModelListener
} from "@application/contract/translation-input/translation-input-screen-contract"
import type {
  TranslationInputScreenState,
  TranslationInputScreenViewModel,
  TranslationInputStagedFile
} from "@application/gateway-contract/translation-input"

interface TranslationInputStoreLike {
  subscribe(listener: (state: TranslationInputScreenState) => void): () => void
  snapshot(): TranslationInputScreenState
  update(mutator: (draft: TranslationInputScreenState) => void): void
}

interface TranslationInputPresenterLike {
  toViewModel(
    state: TranslationInputScreenState,
    isGatewayConnected: boolean
  ): TranslationInputScreenViewModel
}

interface TranslationInputUseCaseLike {
  startImport(): Promise<void>
  rebuildSelected(): Promise<void>
}

interface TranslationInputScreenControllerDependencies {
  isGatewayConnected: boolean
  store: TranslationInputStoreLike
  presenter: TranslationInputPresenterLike
  useCase: TranslationInputUseCaseLike
}

function resolveFileReference(file: File): string {
  const pathRecord = file as File & {
    path?: string
    webkitRelativePath?: string
  }
  const candidates = [pathRecord.path, pathRecord.webkitRelativePath, file.name]
  for (const candidate of candidates) {
    if (typeof candidate === "string" && candidate.trim() !== "") {
      return candidate
    }
  }
  return file.name
}

async function digestFileHash(file: File): Promise<string> {
  const digestApi = globalThis.crypto?.subtle
  if (!digestApi) {
    return "計算不可"
  }

  const bytes = await file.arrayBuffer()
  const digest = await digestApi.digest("SHA-256", bytes)
  return Array.from(new Uint8Array(digest))
    .map((value) => value.toString(16).padStart(2, "0"))
    .join("")
}

export class TranslationInputScreenController
  implements TranslationInputScreenControllerContract
{
  constructor(
    private readonly dependencies: TranslationInputScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return Promise.resolve()
  }

  dispose(): void {
    return
  }

  subscribe(listener: TranslationInputScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): TranslationInputScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  selectItem(localId: string): void {
    this.dependencies.store.update((draft) => {
      draft.selectedItemId = localId
      draft.errorMessage = ""
    })
  }

  async stageJsonImport(file: File | null): Promise<void> {
    if (!file) {
      this.resetImportSelection()
      return
    }

    const stagedFile: TranslationInputStagedFile = {
      fileName: file.name,
      filePath: resolveFileReference(file),
      fileHash: "計算中"
    }

    this.dependencies.store.update((draft) => {
      draft.stagedFile = stagedFile
      draft.operationState = "ready"
      draft.errorMessage = ""
    })

    try {
      const fileHash = await digestFileHash(file)
      this.dependencies.store.update((draft) => {
        if (draft.stagedFile?.filePath !== stagedFile.filePath) {
          return
        }

        draft.stagedFile.fileHash = fileHash
      })
    } catch {
      this.dependencies.store.update((draft) => {
        if (draft.stagedFile?.filePath !== stagedFile.filePath) {
          return
        }

        draft.stagedFile.fileHash = "計算失敗"
      })
    }
  }

  resetImportSelection(): void {
    this.dependencies.store.update((draft) => {
      if (draft.operationState === "importing") {
        return
      }

      draft.stagedFile = null
      draft.operationState = "idle"
      draft.errorMessage = ""
    })
  }

  async startImport(): Promise<void> {
    await this.dependencies.useCase.startImport()
  }

  async rebuildSelected(): Promise<void> {
    await this.dependencies.useCase.rebuildSelected()
  }
}