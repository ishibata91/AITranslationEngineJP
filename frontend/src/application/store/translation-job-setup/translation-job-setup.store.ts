import type {
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"

type Listener = (state: TranslationJobSetupScreenState) => void

function cloneOptions(
  options: TranslationJobSetupOptionsResponse | null
): TranslationJobSetupOptionsResponse | null {
  if (!options) {
    return null
  }

  return {
    ...options,
    inputCandidates: options.inputCandidates.map((candidate) => ({ ...candidate })),
    existingJob: options.existingJob ? { ...options.existingJob } : undefined,
    sharedDictionaries: options.sharedDictionaries.map((option) => ({ ...option })),
    sharedPersonas: options.sharedPersonas.map((option) => ({ ...option })),
    aiRuntimeOptions: options.aiRuntimeOptions.map((option) => ({ ...option })),
    credentialRefs: options.credentialRefs.map((credential) => ({ ...credential }))
  }
}

function cloneValidation(
  validationResult: TranslationJobSetupValidationResponse | null
): TranslationJobSetupValidationResponse | null {
  if (!validationResult) {
    return null
  }

  return {
    ...validationResult,
    targetSlices: [...validationResult.targetSlices],
    passSlices: [...validationResult.passSlices]
  }
}

function cloneSummary(
  summary: TranslationJobSetupSummaryResponse | null
): TranslationJobSetupSummaryResponse | null {
  if (!summary) {
    return null
  }

  return {
    ...summary,
    executionSummary: { ...summary.executionSummary },
    validationPassSlices: [...summary.validationPassSlices]
  }
}

function createInitialState(): TranslationJobSetupScreenState {
  return {
    phase: "idle",
    options: null,
    selectedInputSourceId: null,
    selectedRuntimeKey: null,
    selectedCredentialRef: "",
    validationResult: null,
    validationState: "not-run",
    dirty: false,
    errorMessage: "",
    createErrorKind: null,
    summary: null
  }
}

export class TranslationJobSetupStore {
  private state: TranslationJobSetupScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): TranslationJobSetupScreenState {
    return {
      ...this.state,
      options: cloneOptions(this.state.options),
      validationResult: cloneValidation(this.state.validationResult),
      summary: cloneSummary(this.state.summary)
    }
  }

  update(mutator: (draft: TranslationJobSetupScreenState) => void): void {
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