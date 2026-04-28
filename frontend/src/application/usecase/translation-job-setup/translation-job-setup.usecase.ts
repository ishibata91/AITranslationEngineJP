import {
  createTranslationJobSetupRuntimeKey,
  CreateTranslationJobResponse,
  TranslationJobSetupGatewayContract,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupRuntimeOption,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationState
} from "@application/gateway-contract/translation-job-setup"

interface TranslationJobSetupStoreLike {
  snapshot(): TranslationJobSetupScreenState
  update(mutator: (draft: TranslationJobSetupScreenState) => void): void
}

function sanitizeErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message.startsWith("Wails binding is not wired yet:")) {
    return error.message
  }

  return fallback
}

function resolveInitialRuntimeKey(
  options: TranslationJobSetupOptionsResponse
): string | null {
  return options.aiRuntimeOptions[0]
    ? createTranslationJobSetupRuntimeKey(options.aiRuntimeOptions[0])
    : null
}

function findRuntimeOption(
  options: TranslationJobSetupOptionsResponse | null,
  runtimeKey: string | null
): TranslationJobSetupRuntimeOption | null {
  if (!options || !runtimeKey) {
    return null
  }

  return (
    options.aiRuntimeOptions.find(
      (option) => createTranslationJobSetupRuntimeKey(option) === runtimeKey
    ) ?? null
  )
}

function resolveCredentialRef(
  options: TranslationJobSetupOptionsResponse,
  runtimeOption: TranslationJobSetupRuntimeOption | null,
  currentCredentialRef = ""
): string {
  const credentialRefs = runtimeOption
    ? options.credentialRefs.filter(
        (credential) => credential.provider === runtimeOption.provider
      )
    : options.credentialRefs
  const candidates = credentialRefs.length > 0 ? credentialRefs : options.credentialRefs

  if (
    currentCredentialRef &&
    candidates.some((credential) => credential.credentialRef === currentCredentialRef)
  ) {
    return currentCredentialRef
  }

  return (
    candidates.find(
      (credential) => credential.isConfigured && !credential.isMissingSecret
    )?.credentialRef ?? candidates[0]?.credentialRef ?? ""
  )
}

function invalidateValidation(
  draft: TranslationJobSetupScreenState,
  nextValidationState: TranslationJobSetupValidationState
): void {
  if (draft.summary) {
    return
  }

  if (draft.validationResult) {
    draft.validationState = nextValidationState
    draft.dirty = true
  } else {
    draft.validationState = "not-run"
    draft.dirty = false
  }
  draft.createErrorKind = null
  draft.errorMessage = ""
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
    canStartPhase: false
  }
}

function isExistingJobForInput(
  options: TranslationJobSetupOptionsResponse,
  inputSourceId: number
): boolean {
  const existingJob = options.existingJob
  if (!existingJob) {
    return false
  }

  if ((existingJob.inputSourceId ?? 0) > 0) {
    return existingJob.inputSourceId === inputSourceId
  }

  const inputCandidate = options.inputCandidates.find(
    (candidate) => candidate.id === inputSourceId
  )
  return existingJob.inputSource === inputCandidate?.label
}

export class TranslationJobSetupUseCase {
  constructor(
    private readonly gateway: TranslationJobSetupGatewayContract | null,
    private readonly store: TranslationJobSetupStoreLike
  ) {}

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
      const selectedInputSourceId = options.inputCandidates[0]?.id ?? null
      const selectedRuntimeKey = resolveInitialRuntimeKey(options)
      const selectedRuntimeOption = findRuntimeOption(options, selectedRuntimeKey)
      const selectedCredentialRef = resolveCredentialRef(
        options,
        selectedRuntimeOption
      )

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.options = options
        draft.selectedInputSourceId = selectedInputSourceId
        draft.selectedRuntimeKey = selectedRuntimeKey
        draft.selectedCredentialRef = selectedCredentialRef
        draft.validationResult = null
        draft.validationState = "not-run"
        draft.dirty = false
        draft.createErrorKind = null
        draft.summary = null
      })
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
      if (draft.summary || draft.selectedInputSourceId === inputSourceId) {
        return
      }

      draft.selectedInputSourceId = inputSourceId
      invalidateValidation(draft, "stale")
    })
  }

  selectRuntime(runtimeKey: string): void {
    this.store.update((draft) => {
      if (draft.summary || draft.selectedRuntimeKey === runtimeKey) {
        return
      }

      draft.selectedRuntimeKey = runtimeKey
      draft.selectedCredentialRef = draft.options
        ? resolveCredentialRef(
            draft.options,
            findRuntimeOption(draft.options, runtimeKey),
            draft.selectedCredentialRef
          )
        : ""
      invalidateValidation(draft, "stale")
    })
  }

  selectCredentialRef(credentialRef: string): void {
    this.store.update((draft) => {
      if (draft.summary || draft.selectedCredentialRef === credentialRef) {
        return
      }

      draft.selectedCredentialRef = credentialRef
      invalidateValidation(draft, "stale")
    })
  }

  async runValidation(): Promise<void> {
    const state = this.store.snapshot()
    const runtimeOption = findRuntimeOption(state.options, state.selectedRuntimeKey)
    if (state.selectedInputSourceId === null || !runtimeOption || !state.selectedCredentialRef) {
      this.store.update((draft) => {
        draft.errorMessage = "validation 対象の入力、runtime、credential を選択してください。"
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

  async createJob(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.options || state.summary) {
      return
    }

    const inputCandidate = state.options.inputCandidates.find(
      (candidate) => candidate.id === state.selectedInputSourceId
    )
    const runtimeOption = findRuntimeOption(state.options, state.selectedRuntimeKey)
    if (
      !inputCandidate ||
      !runtimeOption ||
      !state.validationResult ||
      state.validationState !== "fresh" ||
      state.dirty ||
      !state.validationResult.canCreate ||
      isExistingJobForInput(state.options, inputCandidate.id)
    ) {
      this.store.update((draft) => {
        draft.errorMessage = "create 条件を満たしていません。validation と既存 job 状態を確認してください。"
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
        summary = await this.gateway.getTranslationJobSetupSummary({ jobId: response.jobId })
      } catch {
        if (!summary) {
          throw new Error("summary fetch failed")
        }
      }

      this.store.update((draft) => {
        draft.phase = "summary"
        draft.summary = summary
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