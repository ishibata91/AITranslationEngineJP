import type { TranslationInputScreenState } from "@application/gateway-contract/translation-input"

type Listener = (state: TranslationInputScreenState) => void

function asArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : []
}

function createInitialState(): TranslationInputScreenState {
  return {
    items: [],
    selectedItemId: null,
    stagedFile: null,
    operationState: "idle",
    errorMessage: "",
    latestResponse: null
  }
}

export class TranslationInputStore {
  private state: TranslationInputScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): TranslationInputScreenState {
    return {
      ...this.state,
      items: asArray(this.state.items).map((item) => ({
        ...item,
        warnings: asArray(item.warnings).map((warning) => ({ ...warning })),
        summary: item.summary
          ? {
              ...item.summary,
              input: { ...item.summary.input },
              categories: asArray(item.summary.categories).map((category) => ({
                ...category
              })),
              sampleFields: asArray(item.summary.sampleFields).map((field) => ({
                ...field
              })),
              warnings: asArray(item.summary.warnings).map((warning) => ({ ...warning }))
            }
          : null
      })),
      stagedFile: this.state.stagedFile ? { ...this.state.stagedFile } : null,
      latestResponse: this.state.latestResponse
        ? {
            ...this.state.latestResponse,
            warnings: asArray(this.state.latestResponse.warnings).map((warning) => ({
              ...warning
            })),
            summary: this.state.latestResponse.summary
              ? {
                  ...this.state.latestResponse.summary,
                  input: { ...this.state.latestResponse.summary.input },
                  categories:
                    asArray(this.state.latestResponse.summary.categories).map((category) => ({
                      ...category
                    })),
                  sampleFields:
                    asArray(this.state.latestResponse.summary.sampleFields).map((field) => ({
                      ...field
                    })),
                  warnings:
                    asArray(this.state.latestResponse.summary.warnings).map((warning) => ({
                      ...warning
                    }))
                }
              : undefined
          }
        : null
    }
  }

  update(mutator: (draft: TranslationInputScreenState) => void): void {
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