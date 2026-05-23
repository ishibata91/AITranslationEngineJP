import type {
  TranslationInputScreenControllerContract,
  TranslationInputScreenViewModelListener
} from "@application/contract/translation-input/translation-input-screen-contract"
import type {
  CreateTranslationJobFromInputResponse,
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
  startImport(importDraft?: {
    fileName?: string
    fileContent?: string
  }): Promise<void>
  rebuildSelected(): Promise<void>
  createTranslationJobFromSelected?: () => Promise<CreateTranslationJobFromInputResponse | null>
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

async function readFileContentFromArrayBuffer(file: File): Promise<string> {
  const bytes = await file.arrayBuffer()
  return new TextDecoder().decode(bytes)
}

async function readFileContent(file: File): Promise<string> {
  if (typeof file.text === "function") {
    try {
      return await file.text()
    } catch {
      return readFileContentFromArrayBuffer(file)
    }
  }

  return readFileContentFromArrayBuffer(file)
}

export class TranslationInputScreenController implements TranslationInputScreenControllerContract {
  private stagedImportDraft: {
    fileName: string
    fileContent: string
  } | null = null
  private stagedSelectionVersion = 0

  constructor(
    private readonly dependencies: TranslationInputScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return Promise.resolve()
  }

  dispose(): void {
    this.stagedImportDraft = null
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

    const selectionVersion = this.stagedSelectionVersion + 1
    this.stagedSelectionVersion = selectionVersion
    this.stagedImportDraft = null

    const stagedFile: TranslationInputStagedFile = {
      fileName: file.name,
      filePath: resolveFileReference(file),
      fileHash: "計算中"
    }

    this.dependencies.store.update((draft) => {
      draft.stagedFile = stagedFile
      draft.operationState = "idle"
      draft.errorMessage = ""
    })

    try {
      const [fileHashResult, fileContentResult] = await Promise.allSettled([
        digestFileHash(file),
        readFileContent(file)
      ])
      if (this.stagedSelectionVersion !== selectionVersion) {
        return
      }

      if (fileContentResult.status !== "fulfilled") {
        throw fileContentResult.reason
      }

      const fileHash =
        fileHashResult.status === "fulfilled"
          ? fileHashResult.value
          : "計算失敗"
      const fileContent = fileContentResult.value

      this.dependencies.store.update((draft) => {
        if (
          this.stagedSelectionVersion !== selectionVersion ||
          draft.stagedFile === null
        ) {
          return
        }

        draft.stagedFile.fileHash = fileHash
        draft.operationState = "ready"
      })
      this.stagedImportDraft = {
        fileName: file.name,
        fileContent
      }
    } catch {
      this.dependencies.store.update((draft) => {
        if (
          this.stagedSelectionVersion !== selectionVersion ||
          draft.stagedFile === null
        ) {
          return
        }

        draft.stagedFile.fileHash = "計算失敗"
        draft.operationState = "idle"
        draft.errorMessage =
          "JSON file の読み込み完了を確認できませんでした。もう一度選択してください。"
      })
    }
  }

  resetImportSelection(): void {
    this.stagedSelectionVersion += 1
    this.stagedImportDraft = null
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
    const state = this.dependencies.store.snapshot()
    if (state.stagedFile !== null && this.stagedImportDraft === null) {
      this.dependencies.store.update((draft) => {
        draft.errorMessage =
          "JSON file の読み込み完了を待ってから登録してください。"
      })
      return
    }

    await this.dependencies.useCase.startImport(
      this.stagedImportDraft ?? undefined
    )
    if (this.dependencies.store.snapshot().stagedFile === null) {
      this.stagedImportDraft = null
    }
  }

  async rebuildSelected(): Promise<void> {
    await this.dependencies.useCase.rebuildSelected()
  }

  createTranslationJobFromSelected(): Promise<CreateTranslationJobFromInputResponse | null> {
    if (!this.dependencies.useCase.createTranslationJobFromSelected) {
      return Promise.resolve(null)
    }
    return this.dependencies.useCase.createTranslationJobFromSelected()
  }
}
