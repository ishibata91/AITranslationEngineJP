import type {
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupPhaseRuntimeSelection,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"
import { cloneModelSettingsCardStates } from "@application/gateway-contract/model-settings-card"

type Listener = (state: TranslationJobSetupScreenState) => void

function cloneStringArray(values: string[] | null | undefined): string[] {
  if (!Array.isArray(values)) {
    return []
  }

  return [...values]
}

function cloneOptions(
  options: TranslationJobSetupOptionsResponse | null
): TranslationJobSetupOptionsResponse | null {
  if (!options) {
    return null
  }
  const publicOptions = {
    ...options
  } as TranslationJobSetupOptionsResponse & {
    credentialRefs?: unknown
  }
  delete publicOptions.credentialRefs

  return {
    ...publicOptions,
    inputCandidates: publicOptions.inputCandidates.map((candidate) => ({
      ...candidate
    })),
    existingJob: publicOptions.existingJob
      ? { ...publicOptions.existingJob }
      : undefined,
    sharedDictionaries: publicOptions.sharedDictionaries.map((option) => ({
      ...option
    })),
    sharedPersonas: publicOptions.sharedPersonas.map((option) => ({
      ...option
    })),
    aiRuntimeOptions: publicOptions.aiRuntimeOptions.map((option) => ({
      ...option
    })),
    providerCapabilities: publicOptions.providerCapabilities?.map(
      (capability) => ({
        ...capability,
        supportedExecutionModes: [...capability.supportedExecutionModes]
      })
    ),
    phaseRuntimeDrafts: publicOptions.phaseRuntimeDrafts?.map((draft) => ({
      phaseId: draft.phaseId,
      provider: draft.provider,
      model: draft.model,
      credentialStatus: draft.credentialStatus,
      executionMode: draft.executionMode,
      batchMode: draft.batchMode
    }))
  }
}

function clonePhaseSelections(
  selections: TranslationJobSetupPhaseRuntimeSelection[] | null | undefined
): TranslationJobSetupPhaseRuntimeSelection[] {
  if (!Array.isArray(selections)) {
    return []
  }

  return selections.map((selection) => ({ ...selection }))
}

function cloneProviderModelLists(
  lists: ListTranslationJobSetupProviderModelsResponse[] | null | undefined
): ListTranslationJobSetupProviderModelsResponse[] {
  if (!Array.isArray(lists)) {
    return []
  }

  return lists.map((list) => ({
    ...list,
    models: list.models.map((model) => ({ ...model }))
  }))
}

function cloneValidation(
  validationResult: TranslationJobSetupValidationResponse | null
): TranslationJobSetupValidationResponse | null {
  if (!validationResult) {
    return null
  }

  return {
    ...validationResult,
    targetSlices: cloneStringArray(validationResult.targetSlices),
    passSlices: cloneStringArray(validationResult.passSlices),
    phaseResults: validationResult.phaseResults?.map((result) => ({
      ...result
    })),
    staleModelListPhaseIds: validationResult.staleModelListPhaseIds
      ? [...validationResult.staleModelListPhaseIds]
      : undefined
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
    validationPassSlices: [...summary.validationPassSlices],
    phaseRuntimeSummaries: summary.phaseRuntimeSummaries?.map((item) => ({
      ...item
    }))
  }
}

function createInitialState(): TranslationJobSetupScreenState {
  return {
    phase: "idle",
    options: null,
    selectedInputSourceId: null,
    deletingInputSourceId: null,
    selectedRuntimeKey: null,
    selectedCredentialRef: "",
    phaseRuntimeSelections: [],
    providerModelLists: [],
    modelSettingsCards: [],
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
      phaseRuntimeSelections: clonePhaseSelections(
        this.state.phaseRuntimeSelections
      ),
      providerModelLists: cloneProviderModelLists(this.state.providerModelLists),
      modelSettingsCards: cloneModelSettingsCardStates(
        this.state.modelSettingsCards
      ),
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
