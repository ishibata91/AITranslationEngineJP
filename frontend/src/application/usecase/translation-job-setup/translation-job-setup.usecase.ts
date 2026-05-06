import type {
  CreateTranslationJobResponse,
  DeleteTranslationJobSetupInputResponse,
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupBatchMode,
  TranslationJobSetupCredentialReference,
  TranslationJobSetupCredentialStatus,
  TranslationJobSetupGatewayContract,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupPhaseId,
  TranslationJobSetupPhaseRuntimeDraft,
  TranslationJobSetupPhaseRuntimeSelection,
  TranslationJobSetupProviderCapability,
  TranslationJobSetupProviderModelListStatus,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationState
} from "@application/gateway-contract/translation-job-setup"

interface TranslationJobSetupStoreLike {
  snapshot(): TranslationJobSetupScreenState
  update(mutator: (draft: TranslationJobSetupScreenState) => void): void
}

const PHASE_IDS: TranslationJobSetupPhaseId[] = [
  "word_translation",
  "npc_persona_generation",
  "text_translation"
]

function normalizeRequestTokenProvider(provider: string): string {
  const normalized = provider.trim().toLowerCase()
  if (normalized === "lm_studio") {
    return "lm"
  }

  return normalized.replaceAll(/[^a-z0-9]+/g, "-") || "provider"
}

function credentialAllowsEmptySecret(provider: string): boolean {
  return provider === "lm_studio"
}

function findRuntimeOption(
  options: TranslationJobSetupOptionsResponse | null,
  runtimeKey: string | null
): {
  provider: string
  model: string
  mode: string
} | null {
  if (!options || !runtimeKey) {
    return null
  }

  return (
    options.aiRuntimeOptions.find(
      (option) =>
        [option.provider, option.model, option.mode].join("::") === runtimeKey
    ) ?? null
  )
}

function isUsableCredentialRef(options: {
  provider: string
  isConfigured: boolean
  isMissingSecret: boolean
}): boolean {
  return (
    options.isConfigured &&
    (!options.isMissingSecret || credentialAllowsEmptySecret(options.provider))
  )
}

function resolveInitialRuntimeKey(
  options: TranslationJobSetupOptionsResponse
): string | null {
  return options.aiRuntimeOptions[0]
    ? [
        options.aiRuntimeOptions[0].provider,
        options.aiRuntimeOptions[0].model,
        options.aiRuntimeOptions[0].mode
      ].join("::")
    : null
}

function resolveLegacyCredentialRef(
  options: TranslationJobSetupOptionsResponse,
  runtimeOption: {
    provider: string
    model: string
    mode: string
  } | null,
  currentCredentialRef = ""
): string {
  const credentialRefs = runtimeOption
    ? options.credentialRefs.filter(
        (credential) => credential.provider === runtimeOption.provider
      )
    : options.credentialRefs
  const candidates =
    credentialRefs.length > 0 ? credentialRefs : options.credentialRefs

  if (
    currentCredentialRef &&
    candidates.some(
      (credential) => credential.credentialRef === currentCredentialRef
    )
  ) {
    return currentCredentialRef
  }

  return (
    candidates.find((credential) => isUsableCredentialRef(credential))
      ?.credentialRef ??
    candidates[0]?.credentialRef ??
    ""
  )
}

function providerNeedsCredential(
  capability: TranslationJobSetupProviderCapability | null
): boolean {
  return capability?.credentialRequirement === "required"
}

function normalizePhaseCredentialRef(
  options: TranslationJobSetupOptionsResponse | null,
  provider: string,
  credentialRef: string
): string {
  const capability = findCapability(options, provider)
  if (!providerNeedsCredential(capability)) {
    return ""
  }

  const credential = resolveCredentialReference(
    options,
    provider,
    credentialRef
  )
  return credential?.credentialRef ?? ""
}

function sanitizeErrorMessage(error: unknown, fallback: string): string {
  if (
    error instanceof Error &&
    error.message.startsWith("Wails binding is not wired yet:")
  ) {
    return error.message
  }

  return fallback
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

function isPhaseDrivenState(state: TranslationJobSetupScreenState): boolean {
  return (
    (state.phaseRuntimeSelections?.length ?? 0) > 0 ||
    (state.options?.providerCapabilities?.length ?? 0) > 0 ||
    (state.options?.phaseRuntimeDrafts?.length ?? 0) > 0
  )
}

function createFallbackSummary(
  response: CreateTranslationJobResponse
): TranslationJobSetupSummaryResponse | null {
  if (!response.executionSummary || response.errorKind) {
    return null
  }

  return {
    jobId: response.jobId,
    jobState: response.jobState,
    inputSource: response.inputSource,
    executionSummary: { ...response.executionSummary },
    validationPassSlices: [...response.validationPassSlices],
    canStartPhase: false,
    phaseRuntimeSummaries: response.phaseRuntimeSummaries?.map((item) => ({
      ...item
    }))
  }
}

function findCapability(
  options: TranslationJobSetupOptionsResponse | null,
  provider: string
): TranslationJobSetupProviderCapability | null {
  return (
    options?.providerCapabilities?.find(
      (capability) => capability.provider === provider
    ) ?? null
  )
}

function resolveExecutionMode(
  capability: TranslationJobSetupProviderCapability | null,
  currentMode: string
): string {
  if (
    capability?.supportedExecutionModes.includes(currentMode) &&
    currentMode !== ""
  ) {
    return currentMode
  }

  return capability?.supportedExecutionModes[0] ?? "sync"
}

function resolveCredentialReference(
  options: TranslationJobSetupOptionsResponse | null,
  provider: string,
  currentCredentialRef = ""
): TranslationJobSetupCredentialReference | null {
  const matches =
    options?.credentialRefs.filter(
      (credential) => credential.provider === provider
    ) ?? []

  if (currentCredentialRef) {
    const current = matches.find(
      (credential) => credential.credentialRef === currentCredentialRef
    )
    if (current) {
      return current
    }
  }

  return matches[0] ?? null
}

function toCredentialStatus(
  capability: TranslationJobSetupProviderCapability | null,
  credential: TranslationJobSetupCredentialReference | null
): TranslationJobSetupCredentialStatus {
  if (!providerNeedsCredential(capability)) {
    return "not_required"
  }

  if (!credential || !credential.isConfigured) {
    return "missing"
  }

  if (
    credential.isMissingSecret &&
    !credentialAllowsEmptySecret(credential.provider)
  ) {
    return "missing"
  }

  return "configured"
}

function normalizeBatchMode(
  capability: TranslationJobSetupProviderCapability | null,
  currentBatchMode: TranslationJobSetupBatchMode
): TranslationJobSetupBatchMode {
  if (!capability?.supportsBatchMode) {
    return "unsupported"
  }

  return currentBatchMode === "enabled" ? "enabled" : "disabled"
}

function buildPhaseSelection(
  options: TranslationJobSetupOptionsResponse,
  phaseId: TranslationJobSetupPhaseId,
  sourceDraft?: TranslationJobSetupPhaseRuntimeDraft
): TranslationJobSetupPhaseRuntimeSelection {
  const provider =
    sourceDraft?.provider ?? options.providerCapabilities?.[0]?.provider ?? ""
  const capability = findCapability(options, provider)
  const credential = resolveCredentialReference(
    options,
    provider,
    sourceDraft?.credentialRef ?? ""
  )

  return {
    phaseId,
    provider,
    model: sourceDraft?.model ?? "",
    credentialRef: normalizePhaseCredentialRef(
      options,
      provider,
      credential?.credentialRef ?? sourceDraft?.credentialRef ?? ""
    ),
    credentialStatus:
      sourceDraft?.credentialStatus ??
      toCredentialStatus(capability, credential),
    executionMode: resolveExecutionMode(
      capability,
      sourceDraft?.executionMode ?? ""
    ),
    batchMode: normalizeBatchMode(
      capability,
      sourceDraft?.batchMode ?? "disabled"
    ),
    modelListSourceToken: sourceDraft?.modelListSourceToken ?? ""
  }
}

function createPhaseSelections(
  options: TranslationJobSetupOptionsResponse
): TranslationJobSetupPhaseRuntimeSelection[] {
  return PHASE_IDS.map((phaseId) =>
    buildPhaseSelection(
      options,
      phaseId,
      options.phaseRuntimeDrafts?.find((draft) => draft.phaseId === phaseId)
    )
  )
}

function createModelListState(
  selection: TranslationJobSetupPhaseRuntimeSelection
): ListTranslationJobSetupProviderModelsResponse {
  return {
    phaseId: selection.phaseId,
    provider: selection.provider,
    credentialStatus: selection.credentialStatus,
    requestToken: "",
    sourceToken: selection.modelListSourceToken,
    status: "not_updated",
    models: []
  }
}

function createModelListStates(
  selections: TranslationJobSetupPhaseRuntimeSelection[]
): ListTranslationJobSetupProviderModelsResponse[] {
  return selections.map((selection) => createModelListState(selection))
}

function findPhaseSelection(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): TranslationJobSetupPhaseRuntimeSelection | null {
  return (
    state.phaseRuntimeSelections?.find(
      (selection) => selection.phaseId === phaseId
    ) ?? null
  )
}

function findModelList(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): ListTranslationJobSetupProviderModelsResponse | null {
  return (
    state.providerModelLists?.find((entry) => entry.phaseId === phaseId) ?? null
  )
}

function replacePhaseSelection(
  draft: TranslationJobSetupScreenState,
  nextSelection: TranslationJobSetupPhaseRuntimeSelection
): void {
  draft.phaseRuntimeSelections = PHASE_IDS.map((phaseId) =>
    phaseId === nextSelection.phaseId
      ? { ...nextSelection }
      : (draft.phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === phaseId
        ) ?? nextSelection)
  )
}

function replaceModelList(
  draft: TranslationJobSetupScreenState,
  nextModelList: ListTranslationJobSetupProviderModelsResponse
): void {
  draft.providerModelLists = PHASE_IDS.map((phaseId) =>
    phaseId === nextModelList.phaseId
      ? {
          ...nextModelList,
          models: nextModelList.models.map((model) => ({ ...model }))
        }
      : (draft.providerModelLists?.find(
          (entry) => entry.phaseId === phaseId
        ) ?? {
          phaseId,
          provider: "",
          credentialStatus: "missing",
          requestToken: "",
          sourceToken: "",
          status: "not_updated",
          models: []
        })
  )
}

function isModelListUsable(
  status: TranslationJobSetupProviderModelListStatus
): boolean {
  return status === "success" || status === "credential_not_required"
}

function isPhaseLocallyComplete(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): boolean {
  const selection = findPhaseSelection(state, phaseId)
  const modelList = findModelList(state, phaseId)
  if (!selection || selection.provider === "" || selection.model === "") {
    return false
  }

  if (selection.credentialStatus === "missing") {
    return false
  }

  if (!modelList || !isModelListUsable(modelList.status)) {
    return false
  }

  if (selection.modelListSourceToken === "") {
    return false
  }

  return true
}

function isLocallyReadyForCreate(
  state: TranslationJobSetupScreenState
): boolean {
  if (!state.options || state.selectedInputSourceId === null) {
    return false
  }

  return PHASE_IDS.every((phaseId) => isPhaseLocallyComplete(state, phaseId))
}

function createValidationPayload(state: TranslationJobSetupScreenState): {
  inputSourceId: number
  runtime: {
    provider: string
    model: string
    executionMode: string
  }
  credentialRef: string
  phaseRuntimeSelections: TranslationJobSetupPhaseRuntimeSelection[]
} | null {
  if (
    !state.options ||
    state.selectedInputSourceId === null ||
    !isLocallyReadyForCreate(state)
  ) {
    return null
  }

  const primarySelection =
    state.phaseRuntimeSelections?.find(
      (selection) => selection.phaseId === "word_translation"
    ) ?? state.phaseRuntimeSelections?.[0]
  if (!primarySelection) {
    return null
  }

  return {
    inputSourceId: state.selectedInputSourceId,
    runtime: {
      provider: primarySelection.provider,
      model: primarySelection.model,
      executionMode: primarySelection.executionMode
    },
    credentialRef: normalizePhaseCredentialRef(
      state.options,
      primarySelection.provider,
      primarySelection.credentialRef
    ),
    phaseRuntimeSelections:
      state.phaseRuntimeSelections?.map((selection) => ({
        ...selection,
        credentialRef: normalizePhaseCredentialRef(
          state.options,
          selection.provider,
          selection.credentialRef
        )
      })) ?? []
  }
}

function invalidateValidation(
  draft: TranslationJobSetupScreenState,
  nextValidationState: TranslationJobSetupValidationState
): void {
  if (draft.summary) {
    return
  }

  draft.validationState = nextValidationState
  draft.dirty = nextValidationState !== "not-run"
  draft.validationResult = null
  draft.createErrorKind = null
  draft.errorMessage = ""
}

function applyLoadedOptions(
  draft: TranslationJobSetupScreenState,
  options: TranslationJobSetupOptionsResponse,
  preferredInputSourceId: number | null = null
): void {
  const hasPhaseDrafts =
    (options.phaseRuntimeDrafts?.length ?? 0) > 0 ||
    (options.providerCapabilities?.length ?? 0) > 0
  const phaseRuntimeSelections = hasPhaseDrafts
    ? createPhaseSelections(options)
    : []
  const providerModelLists = hasPhaseDrafts
    ? createModelListStates(phaseRuntimeSelections)
    : []
  const selectedRuntimeKey = hasPhaseDrafts
    ? null
    : resolveInitialRuntimeKey(options)
  const selectedCredentialRef = hasPhaseDrafts
    ? ""
    : resolveLegacyCredentialRef(
        options,
        findRuntimeOption(options, selectedRuntimeKey)
      )
  const selectedInputSourceId =
    preferredInputSourceId !== null &&
    options.inputCandidates.some(
      (candidate) => candidate.id === preferredInputSourceId
    )
      ? preferredInputSourceId
      : (options.inputCandidates[0]?.id ?? null)

  draft.phase = "ready"
  draft.options = options
  draft.selectedInputSourceId = selectedInputSourceId
  draft.deletingInputSourceId = null
  draft.selectedRuntimeKey = selectedRuntimeKey
  draft.selectedCredentialRef = selectedCredentialRef
  draft.phaseRuntimeSelections = phaseRuntimeSelections
  draft.providerModelLists = providerModelLists
  draft.validationResult = null
  draft.validationState = "not-run"
  draft.dirty = false
  draft.createErrorKind = null
  draft.summary = null
}

function findNextInputSourceIdAfterDelete(
  candidates: TranslationJobSetupOptionsResponse["inputCandidates"],
  deletedInputSourceId: number
): number | null {
  const deletedIndex = candidates.findIndex(
    (candidate) => candidate.id === deletedInputSourceId
  )
  if (deletedIndex < 0) {
    return null
  }

  return candidates[deletedIndex + 1]?.id ?? candidates[deletedIndex - 1]?.id ?? null
}

function resolveDeleteInputErrorMessage(
  response: DeleteTranslationJobSetupInputResponse
): string {
  const errorKind = response.errorKind
  switch (errorKind) {
    case "input_delete_blocked":
      return "既存 job が参照している入力データは削除できません。"
    case "input_not_found":
      return "入力データが見つかりません。再読込してください。"
    default:
      return "入力データの削除に失敗しました。"
  }
}

export class TranslationJobSetupUseCase {
  private requestSerial = 0
  private latestValidationToken = ""

  constructor(
    private readonly gateway: TranslationJobSetupGatewayContract | null,
    private readonly store: TranslationJobSetupStoreLike
  ) {}

  private nextRequestToken(provider: string): string {
    this.requestSerial += 1
    return `req-${normalizeRequestTokenProvider(provider)}-${this.requestSerial}`
  }

  private nextValidationToken(): string {
    this.requestSerial += 1
    return `validation-${this.requestSerial}`
  }

  async load(): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = "translation-job-setup gateway が未接続です。"
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.errorMessage = ""
    })

    try {
      const options = await this.gateway.getTranslationJobSetupOptions()
      this.store.update((draft) => {
        applyLoadedOptions(draft, options)
      })

      void this.revalidateIfReady()
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "Job Setup の初期データ取得に失敗しました。"
        )
      })
    }
  }

  selectInputSource(inputSourceId: number): void {
    this.store.update((draft) => {
      if (
        draft.summary ||
        draft.selectedInputSourceId === inputSourceId ||
        draft.deletingInputSourceId === inputSourceId
      ) {
        return
      }

      draft.selectedInputSourceId = inputSourceId
      invalidateValidation(draft, "stale")
    })

    void this.revalidateIfReady()
  }

  async deleteInputSource(inputSourceId: number): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.errorMessage = "translation-job-setup gateway が未接続です。"
      })
      return
    }

    const currentState = this.store.snapshot()
    if (
      currentState.summary ||
      currentState.phase === "creating" ||
      currentState.deletingInputSourceId !== null
    ) {
      return
    }

    this.store.update((draft) => {
      draft.deletingInputSourceId = inputSourceId
      draft.errorMessage = ""
    })

    try {
      const response = await this.gateway.deleteTranslationJobSetupInput({
        inputSourceId
      })
      if (response.errorKind) {
        this.store.update((draft) => {
          draft.deletingInputSourceId = null
          draft.errorMessage = resolveDeleteInputErrorMessage(response)
        })
        return
      }
      const deletedInputSourceId =
        response.deletedInputSourceId ?? inputSourceId
      let shouldRevalidate = false
      this.store.update((draft) => {
        const currentOptions = draft.options
        draft.deletingInputSourceId = null
        draft.errorMessage = ""
        if (!currentOptions) {
          return
        }

        draft.options = {
          ...currentOptions,
          inputCandidates: currentOptions.inputCandidates.filter(
            (candidate) => candidate.id !== deletedInputSourceId
          ),
          existingJob:
            currentOptions.existingJob?.inputSourceId === deletedInputSourceId
              ? undefined
              : currentOptions.existingJob
        }

        if (draft.selectedInputSourceId !== deletedInputSourceId) {
          return
        }

        draft.selectedInputSourceId = findNextInputSourceIdAfterDelete(
          currentOptions.inputCandidates,
          deletedInputSourceId
        )
        invalidateValidation(draft, "not-run")
        shouldRevalidate = draft.selectedInputSourceId !== null
      })

      if (shouldRevalidate) {
        void this.revalidateIfReady()
      }
    } catch (error) {
      this.store.update((draft) => {
        draft.deletingInputSourceId = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "入力データの削除に失敗しました。"
        )
      })
    }
  }

  selectRuntime(_runtimeKey: string): void {
    const runtimeKey = _runtimeKey
    this.store.update((draft) => {
      if (
        draft.summary ||
        isPhaseDrivenState(draft) ||
        draft.selectedRuntimeKey === runtimeKey
      ) {
        return
      }

      draft.selectedRuntimeKey = runtimeKey
      draft.selectedCredentialRef = draft.options
        ? resolveLegacyCredentialRef(
            draft.options,
            findRuntimeOption(draft.options, runtimeKey),
            draft.selectedCredentialRef
          )
        : ""
      invalidateValidation(draft, "stale")
    })
  }

  selectCredentialRef(_credentialRef: string): void {
    const credentialRef = _credentialRef
    this.store.update((draft) => {
      if (
        draft.summary ||
        isPhaseDrivenState(draft) ||
        draft.selectedCredentialRef === credentialRef
      ) {
        return
      }

      draft.selectedCredentialRef = credentialRef
      invalidateValidation(draft, "stale")
    })
  }

  selectPhaseProvider(
    phaseId: TranslationJobSetupPhaseId,
    provider: string
  ): void {
    this.store.update((draft) => {
      if (draft.summary || !draft.options) {
        return
      }

      const currentSelection = findPhaseSelection(draft, phaseId)
      if (!currentSelection || currentSelection.provider === provider) {
        return
      }

      const capability = findCapability(draft.options, provider)
      const credential = resolveCredentialReference(
        draft.options,
        provider,
        currentSelection.credentialRef
      )
      const nextSelection: TranslationJobSetupPhaseRuntimeSelection = {
        ...currentSelection,
        provider,
        model: "",
        credentialRef: normalizePhaseCredentialRef(
          draft.options,
          provider,
          credential?.credentialRef ?? ""
        ),
        credentialStatus: toCredentialStatus(capability, credential),
        executionMode: resolveExecutionMode(
          capability,
          currentSelection.executionMode
        ),
        batchMode: normalizeBatchMode(capability, currentSelection.batchMode),
        modelListSourceToken: ""
      }

      replacePhaseSelection(draft, nextSelection)
      replaceModelList(draft, createModelListState(nextSelection))
      invalidateValidation(draft, "stale")
    })

    void this.revalidateIfReady()
  }

  async refreshPhaseModels(phaseId: TranslationJobSetupPhaseId): Promise<void> {
    const state = this.store.snapshot()
    const selection = findPhaseSelection(state, phaseId)
    if (!selection || !this.gateway) {
      return
    }

    const capability = findCapability(state.options, selection.provider)
    if (
      providerNeedsCredential(capability) &&
      selection.credentialStatus !== "configured"
    ) {
      this.store.update((draft) => {
        replaceModelList(draft, {
          phaseId,
          provider: selection.provider,
          credentialStatus: selection.credentialStatus,
          requestToken: "",
          sourceToken: "",
          status: "credential_missing",
          models: []
        })
        invalidateValidation(draft, "stale")
      })
      return
    }

    const requestToken = this.nextRequestToken(selection.provider)
    this.store.update((draft) => {
      replaceModelList(draft, {
        phaseId,
        provider: selection.provider,
        credentialStatus: selection.credentialStatus,
        requestToken,
        sourceToken: "",
        status: "loading",
        models: []
      })
      invalidateValidation(draft, "stale")
    })

    try {
      const requestCredentialRef = normalizePhaseCredentialRef(
        state.options,
        selection.provider,
        selection.credentialRef
      )
      const response = await this.gateway.listTranslationJobSetupProviderModels(
        {
          phaseId,
          provider: selection.provider,
          credentialRef: requestCredentialRef,
          credentialStatus: selection.credentialStatus,
          requestToken
        }
      )

      const latestState = this.store.snapshot()
      const latestSelection = findPhaseSelection(latestState, phaseId)
      const latestModelList = findModelList(latestState, phaseId)
      if (
        !latestSelection ||
        latestSelection.provider !== response.provider ||
        latestModelList?.requestToken !== response.requestToken
      ) {
        return
      }

      this.store.update((draft) => {
        replaceModelList(draft, response)
        const currentSelection = findPhaseSelection(draft, phaseId)
        if (!currentSelection) {
          return
        }

        if (!isModelListUsable(response.status)) {
          replacePhaseSelection(draft, {
            ...currentSelection,
            model: "",
            modelListSourceToken: ""
          })
        } else if (
          !response.models.some(
            (model) => model.modelId === currentSelection.model
          )
        ) {
          replacePhaseSelection(draft, {
            ...currentSelection,
            model: "",
            modelListSourceToken: ""
          })
        }

        invalidateValidation(draft, "stale")
      })

      void this.revalidateIfReady()
    } catch (error) {
      const latestState = this.store.snapshot()
      const latestSelection = findPhaseSelection(latestState, phaseId)
      const latestModelList = findModelList(latestState, phaseId)
      if (
        !latestSelection ||
        latestSelection.provider !== selection.provider ||
        latestModelList?.requestToken !== requestToken
      ) {
        return
      }

      this.store.update((draft) => {
        replaceModelList(draft, {
          phaseId,
          provider: selection.provider,
          credentialStatus: selection.credentialStatus,
          requestToken,
          sourceToken: "",
          status: "failed",
          models: [],
          failureKind: "provider_unreachable"
        })
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "モデル一覧の取得に失敗しました。"
        )
        invalidateValidation(draft, "stale")
      })
    }
  }

  selectPhaseModel(phaseId: TranslationJobSetupPhaseId, model: string): void {
    this.store.update((draft) => {
      const currentSelection = findPhaseSelection(draft, phaseId)
      const modelList = findModelList(draft, phaseId)
      if (
        draft.summary ||
        !currentSelection ||
        !modelList ||
        !isModelListUsable(modelList.status) ||
        currentSelection.model === model
      ) {
        return
      }

      replacePhaseSelection(draft, {
        ...currentSelection,
        model,
        modelListSourceToken: modelList.sourceToken
      })
      invalidateValidation(draft, "stale")
    })

    void this.revalidateIfReady()
  }

  togglePhaseBatchMode(
    phaseId: TranslationJobSetupPhaseId,
    enabled: boolean
  ): void {
    this.store.update((draft) => {
      const currentSelection = findPhaseSelection(draft, phaseId)
      const capability = findCapability(
        draft.options,
        currentSelection?.provider ?? ""
      )
      if (
        draft.summary ||
        !currentSelection ||
        !capability?.supportsBatchMode
      ) {
        return
      }

      replacePhaseSelection(draft, {
        ...currentSelection,
        batchMode: enabled ? "enabled" : "disabled"
      })
      invalidateValidation(draft, "stale")
    })

    void this.revalidateIfReady()
  }

  async runValidation(): Promise<void> {
    const legacyState = this.store.snapshot()
    if (!isPhaseDrivenState(legacyState)) {
      await this.runLegacyValidation()
      return
    }

    await this.revalidateIfReady(true)
  }

  private async runLegacyValidation(): Promise<void> {
    const state = this.store.snapshot()
    const runtimeOption = findRuntimeOption(
      state.options,
      state.selectedRuntimeKey
    )
    if (
      state.selectedInputSourceId === null ||
      !runtimeOption ||
      !state.selectedCredentialRef
    ) {
      this.store.update((draft) => {
        draft.errorMessage =
          "validation 対象の入力、runtime、credential を選択してください。"
      })
      return
    }

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.errorMessage = "translation-job-setup gateway が未接続です。"
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "validating"
      draft.validationState = "running"
      draft.errorMessage = ""
      draft.createErrorKind = null
    })

    try {
      const validationResult = await this.gateway.validateTranslationJobSetup({
        inputSourceId: state.selectedInputSourceId,
        runtime: {
          provider: runtimeOption.provider,
          model: runtimeOption.model,
          executionMode: runtimeOption.mode
        },
        credentialRef: state.selectedCredentialRef
      })

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.validationResult = validationResult
        draft.validationState = "fresh"
        draft.dirty = false
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.validationState = "not-run"
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "Job Setup の validation に失敗しました。"
        )
      })
    }
  }

  private async revalidateIfReady(force = false): Promise<void> {
    if (!this.gateway) {
      return
    }

    const state = this.store.snapshot()
    const payload = createValidationPayload(state)
    if (!payload) {
      if (force) {
        this.store.update((draft) => {
          draft.errorMessage =
            "翻訳ジョブを作成するには、3 つの翻訳段階でモデル一覧の更新とモデル選択が必要です。"
        })
      }
      return
    }

    if (!force && state.validationState === "fresh" && !state.dirty) {
      return
    }

    const validationToken = this.nextValidationToken()
    this.latestValidationToken = validationToken

    this.store.update((draft) => {
      draft.phase = "validating"
      draft.validationState = "running"
      draft.errorMessage = ""
      draft.createErrorKind = null
    })

    try {
      const validationResult =
        await this.gateway.validateTranslationJobSetup(payload)
      if (this.latestValidationToken !== validationToken) {
        return
      }

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.validationResult = validationResult
        draft.validationState = "fresh"
        draft.dirty = false
      })
    } catch (error) {
      if (this.latestValidationToken !== validationToken) {
        return
      }

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.validationState = "not-run"
        draft.validationResult = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "Job Setup の確認に失敗しました。"
        )
      })
    }
  }

  async createJob(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.options || state.summary) {
      return
    }

    if (!isPhaseDrivenState(state)) {
      await this.createLegacyJob()
      return
    }

    if (!isLocallyReadyForCreate(state)) {
      this.store.update((draft) => {
        draft.errorMessage =
          "翻訳ジョブを作成するには、3 つの翻訳段階で不足を解消してください。"
      })
      return
    }

    if (
      state.validationState !== "fresh" ||
      state.dirty ||
      !state.validationResult?.canCreate
    ) {
      await this.runValidation()
    }

    const latestState = this.store.snapshot()
    const payload = createValidationPayload(latestState)
    const inputCandidate = latestState.options?.inputCandidates.find(
      (candidate) => candidate.id === latestState.selectedInputSourceId
    )
    if (
      !payload ||
      !latestState.options ||
      !latestState.validationResult ||
      latestState.validationState !== "fresh" ||
      latestState.dirty ||
      !latestState.validationResult.canCreate ||
      !inputCandidate
    ) {
      this.store.update((draft) => {
        draft.errorMessage =
          "翻訳ジョブを作成できません。作成前確認の不足を解消してください。"
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "creating"
      draft.errorMessage = ""
      draft.createErrorKind = null
    })

    try {
      const response = await this.gateway.createTranslationJob({
        inputSourceId: payload.inputSourceId,
        inputSource: inputCandidate.label,
        validationStatus: latestState.validationResult.status,
        validatedAt: latestState.validationResult.validatedAt,
        validationPassSlices: [...latestState.validationResult.passSlices],
        runtime: payload.runtime,
        credentialRef: payload.credentialRef,
        phaseRuntimeSelections: payload.phaseRuntimeSelections
      })

      if (response.errorKind) {
        this.store.update((draft) => {
          draft.phase = "ready"
          draft.createErrorKind = response.errorKind ?? null
        })
        return
      }

      let summary = createFallbackSummary(response)
      try {
        summary = await this.gateway.getTranslationJobSetupSummary({
          jobId: response.jobId
        })
      } catch {
        if (!summary) {
          throw new Error("summary fetch failed")
        }
      }

      this.store.update((draft) => {
        draft.phase = "summary"
        draft.summary = cloneSummary(summary)
        draft.dirty = false
        draft.createErrorKind = null
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "翻訳ジョブの作成に失敗しました。"
        )
      })
    }
  }

  private async createLegacyJob(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.options || state.summary) {
      return
    }

    const inputCandidate = state.options.inputCandidates.find(
      (candidate) => candidate.id === state.selectedInputSourceId
    )
    const runtimeOption = findRuntimeOption(
      state.options,
      state.selectedRuntimeKey
    )
    if (
      !inputCandidate ||
      !runtimeOption ||
      !state.validationResult ||
      state.validationState !== "fresh" ||
      state.dirty ||
      !state.validationResult.canCreate
    ) {
      this.store.update((draft) => {
        draft.errorMessage =
          "create 条件を満たしていません。validation を確認してください。"
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "creating"
      draft.errorMessage = ""
      draft.createErrorKind = null
    })

    try {
      const response = await this.gateway.createTranslationJob({
        inputSourceId: inputCandidate.id,
        inputSource: inputCandidate.label,
        validationStatus: state.validationResult.status,
        validatedAt: state.validationResult.validatedAt,
        validationPassSlices: [...state.validationResult.passSlices],
        runtime: {
          provider: runtimeOption.provider,
          model: runtimeOption.model,
          executionMode: runtimeOption.mode
        },
        credentialRef: state.selectedCredentialRef
      })

      if (response.errorKind) {
        this.store.update((draft) => {
          draft.phase = "ready"
          draft.createErrorKind = response.errorKind ?? null
        })
        return
      }

      let summary = createFallbackSummary(response)
      try {
        summary = await this.gateway.getTranslationJobSetupSummary({
          jobId: response.jobId
        })
      } catch {
        if (!summary) {
          throw new Error("summary fetch failed")
        }
      }

      this.store.update((draft) => {
        draft.phase = "summary"
        draft.summary = cloneSummary(summary)
        draft.dirty = false
        draft.createErrorKind = null
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "translation job の create に失敗しました。"
        )
      })
    }
  }
}
