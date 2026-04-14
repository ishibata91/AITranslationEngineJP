import {
  DEFAULT_CATEGORY,
  DEFAULT_ORIGIN
} from "@application/contract/master-dictionary"
import type { MasterDictionaryScreenState } from "@application/contract/master-dictionary/master-dictionary-screen-types"

type Listener = (state: MasterDictionaryScreenState) => void

function createInitialState(): MasterDictionaryScreenState {
  return {
    entries: [],
    selectedEntry: null,
    selectedId: null,
    totalCount: 0,
    query: "",
    category: "すべて",
    page: 0,
    errorMessage: "",
    modalState: null,
    formSource: "",
    formCategory: DEFAULT_CATEGORY,
    formOrigin: DEFAULT_ORIGIN,
    formTranslation: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    importStage: "idle",
    importProgress: 0,
    importSummary: null
  }
}

export class MasterDictionaryStore {
  private state: MasterDictionaryScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): MasterDictionaryScreenState {
    return {
      ...this.state,
      entries: [...this.state.entries],
      selectedEntry: this.state.selectedEntry
        ? { ...this.state.selectedEntry }
        : null,
      importSummary: this.state.importSummary
        ? { ...this.state.importSummary }
        : null
    }
  }

  update(mutator: (draft: MasterDictionaryScreenState) => void): void {
    const nextState = this.snapshot()
    mutator(nextState)
    this.state = nextState
    this.emit()
  }

  private emit(): void {
    const snapshot = this.snapshot()
    for (const listener of this.listeners) {
      listener(snapshot)
    }
  }
}
