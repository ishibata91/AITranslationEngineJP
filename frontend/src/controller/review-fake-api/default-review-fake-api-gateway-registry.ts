import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseGatewayContract,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

import type {
  MasterPersonaDeleteRequest,
  MasterPersonaDetail,
  MasterPersonaAISettingsResponse,
  MasterPersonaGatewayContract,
  MasterPersonaMutationResponse,
  MasterPersonaPageState,
  MasterPersonaPreviewRequest,
  MasterPersonaRunStatus,
  MasterPersonaUpdateRequest
} from "@application/gateway-contract/master-persona"

import type {
  CreateTranslationJobRequest,
  DeleteTranslationJobSetupInputRequest,
  DeleteTranslationJobSetupInputResponse,
  GetTranslationJobSetupSummaryRequest,
  ListTranslationJobSetupProviderModelsRequest,
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupGatewayContract,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupPhaseId,
  TranslationJobSetupPhaseRuntimeSelection,
  TranslationJobSetupPhaseRuntimeValidationSelection,
  TranslationJobSetupPhaseValidationResult,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse,
  ValidateTranslationJobSetupRequest
} from "@application/gateway-contract/translation-job-setup"

import type {
  TranslationOutputArtifactGatewayContract,
  TranslationOutputDiffPreviewResponse,
  TranslationOutputReviewResponse
} from "@application/gateway-contract/translation-output-artifact"

import type {
  ListProviderSettingsResponse,
  ProviderSettingsGatewayContract,
  ProviderSettingsProviderId,
  ProviderSettingsSummary,
  ResetProviderSettingsRequest,
  SaveProviderSettingsRequest,
  ValidateProviderSettingsRequest,
  ValidateProviderSettingsResponse
} from "@application/gateway-contract/provider-settings"

import type {
  ReviewFakeApiGatewayRegistry,
  ReviewFakeApiScenarioId
} from "./review-fake-api-runtime"

const REVIEW_PROVIDER_IDS = ["gemini", "xai", "lm_studio"] as const
const REVIEW_PHASE_IDS: TranslationJobSetupPhaseId[] = [
  "word_translation",
  "npc_persona_generation",
  "text_translation"
]

function createProvider(
  providerId: ProviderSettingsProviderId,
  overrides: Partial<ProviderSettingsSummary> = {}
): ProviderSettingsSummary {
  const baseProvider: ProviderSettingsSummary = {
    providerId,
    label:
      providerId === "lm_studio"
        ? "LM Studio"
        : providerId === "xai"
          ? "xAI"
          : "Gemini",
    endpoint:
      providerId === "lm_studio"
        ? "http://127.0.0.1:1234/v1"
        : "https://example.invalid",
    credentialState: providerId === "lm_studio" ? "not_required" : "configured",
    validationState: "validated",
    savedState: "configured",
    requestToken: `${providerId}-review`
  }

  return {
    ...baseProvider,
    ...overrides
  }
}

function createResponse(
  providers: ProviderSettingsSummary[]
): ListProviderSettingsResponse {
  return {
    route: {
      routeId: "provider-settings",
      label: "AIサービス設定",
      currentRouteState: "レビュー確認",
      dashboardEntryId: "provider-settings"
    },
    providers
  }
}

function createProvidersForScenario(
  scenarioId: ReviewFakeApiScenarioId
): ProviderSettingsSummary[] {
  if (scenarioId === "empty") {
    return []
  }

  if (scenarioId === "running") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        validationState: "pending"
      })
    )
  }

  if (scenarioId === "error") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        validationState: "failed",
        lastFailureKind: "provider_unreachable"
      })
    )
  }

  if (scenarioId === "config-missing") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        credentialState:
          providerId === "lm_studio" ? "not_required" : "missing",
        validationState:
          providerId === "lm_studio" ? "validated" : "not_validated",
        savedState: providerId === "lm_studio" ? "configured" : "partial",
        lastFailureKind:
          providerId === "lm_studio" ? undefined : "credential_missing"
      })
    )
  }

  return REVIEW_PROVIDER_IDS.map((providerId) => createProvider(providerId))
}

function createPendingProviderSettingsGateway(): ProviderSettingsGatewayContract {
  return {
    ListProviderSettings: () => new Promise(() => undefined),
    SaveProviderSettings: (request: SaveProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ResetProviderSettings: (request: ResetProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ValidateProviderSettings: (
      request: ValidateProviderSettingsRequest
    ): Promise<ValidateProviderSettingsResponse> =>
      Promise.resolve({
        providerId: request.providerId,
        validationState: "validated",
        requestToken: request.requestToken
      })
  }
}

function createProviderSettingsGateway(
  scenarioId: ReviewFakeApiScenarioId
): ProviderSettingsGatewayContract {
  if (scenarioId === "loading") {
    return createPendingProviderSettingsGateway()
  }

  return {
    ListProviderSettings: () =>
      Promise.resolve(createResponse(createProvidersForScenario(scenarioId))),
    SaveProviderSettings: (request: SaveProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ResetProviderSettings: (request: ResetProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ValidateProviderSettings: (
      request: ValidateProviderSettingsRequest
    ): Promise<ValidateProviderSettingsResponse> =>
      Promise.resolve({
        providerId: request.providerId,
        validationState: "validated",
        requestToken: request.requestToken
      })
  }
}

function createMasterPersonaEntry(identityKey: string): MasterPersonaDetail {
  return {
    identityKey,
    targetPlugin: "ReviewPlugin.esp",
    formId: "000A12",
    recordType: "NPC_",
    editorId: "ReviewActor",
    displayName: "レビュー対象 NPC",
    race: "NordRace",
    sex: "Female",
    voiceType: "FemaleYoungEager",
    className: "Citizen",
    sourcePlugin: "Skyrim.esm",
    personaSummary: "落ち着いた口調で話すレビュー用ペルソナ。",
    personaBody:
      "この本文は fakeAPI の人間レビュー用データです。secret や実ファイルパスは含めません。",
    speechStyle: "丁寧",
    updatedAt: "2026-05-07T00:00:00Z",
    runLockReason: ""
  }
}

function createMasterPersonaPage(
  scenarioId: ReviewFakeApiScenarioId
): MasterPersonaPageState {
  if (scenarioId === "empty") {
    return {
      items: [],
      pluginGroups: [],
      totalCount: 0,
      page: 1,
      pageSize: 30
    }
  }

  const entry = createMasterPersonaEntry("review-npc-001")
  return {
    items: [entry],
    pluginGroups: [{ targetPlugin: entry.targetPlugin, count: 1 }],
    totalCount: 1,
    page: 1,
    pageSize: 30,
    selectedIdentityKey: entry.identityKey
  }
}

function createMasterPersonaRunStatus(
  scenarioId: ReviewFakeApiScenarioId
): MasterPersonaRunStatus {
  if (scenarioId === "running") {
    return {
      runState: "生成中",
      targetPlugin: "ReviewPlugin.esp",
      processedCount: 4,
      successCount: 3,
      existingSkipCount: 1,
      currentActorLabel: "レビュー対象 NPC",
      message: "fakeAPI で生成中の状態を表示しています。",
      startedAt: "2026-05-07T00:00:00Z"
    }
  }

  return {
    runState: "入力待ち",
    targetPlugin: "ReviewPlugin.esp",
    processedCount: 0,
    successCount: 0,
    existingSkipCount: 0,
    currentActorLabel: "",
    message: "fakeAPI で確認できます。"
  }
}

function createMasterPersonaAISettingsResponse(
  scenarioId: ReviewFakeApiScenarioId
): MasterPersonaAISettingsResponse {
  const credentialStatus =
    scenarioId === "config-missing" ? "missing" : "configured"
  const model =
    credentialStatus === "configured" && scenarioId !== "empty"
      ? "fake-model"
      : ""
  return {
    aiSettings: {
      provider: "gemini",
      model,
      executionMethod: "single_request"
    },
    providerOptions: [
      { value: "gemini", label: "Gemini", credentialStatus },
      {
        value: "lm_studio",
        label: "LM Studio",
        credentialStatus: "not_required" as const
      },
      { value: "xai", label: "xAI", credentialStatus: "missing" as const }
    ],
    modelList: {
      provider: "gemini",
      credentialStatus,
      status: model
        ? ("success" as const)
        : credentialStatus === "missing"
          ? ("credential_missing" as const)
          : ("not_updated" as const),
      models: model ? [{ modelId: model, label: model }] : []
    }
  }
}

function createMasterPersonaGateway(
  scenarioId: ReviewFakeApiScenarioId
): MasterPersonaGatewayContract {
  if (scenarioId === "loading") {
    return {
      getMasterPersonaPage: () => new Promise(() => undefined),
      getMasterPersonaDetail: () => new Promise(() => undefined),
      loadMasterPersonaAISettings: () => new Promise(() => undefined),
      listMasterPersonaProviderModels: () => new Promise(() => undefined),
      saveMasterPersonaAISettings: (request) => Promise.resolve(request),
      previewMasterPersonaGeneration: () => new Promise(() => undefined),
      executeMasterPersonaGeneration: () =>
        Promise.resolve(createMasterPersonaRunStatus("running")),
      getMasterPersonaRunStatus: () =>
        Promise.resolve(createMasterPersonaRunStatus("running")),
      interruptMasterPersonaGeneration: () =>
        Promise.resolve(createMasterPersonaRunStatus("success")),
      cancelMasterPersonaGeneration: () =>
        Promise.resolve(createMasterPersonaRunStatus("success")),
      updateMasterPersona: (request: MasterPersonaUpdateRequest) =>
        Promise.resolve({
          page: createMasterPersonaPage("success"),
          changedEntry: {
            ...createMasterPersonaEntry(request.identityKey),
            personaBody: request.entry.personaBody,
            personaSummary: request.entry.personaSummary ?? "",
            speechStyle: request.entry.speechStyle
          }
        }),
      deleteMasterPersona: (request: MasterPersonaDeleteRequest) =>
        Promise.resolve({
          page: createMasterPersonaPage("empty"),
          deletedEntryId: request.identityKey
        })
    }
  }

  return {
    getMasterPersonaPage: () => {
      return Promise.resolve({ page: createMasterPersonaPage(scenarioId) })
    },
    getMasterPersonaDetail: ({ identityKey }) => {
      return Promise.resolve({ entry: createMasterPersonaEntry(identityKey) })
    },
    loadMasterPersonaAISettings: () => {
      if (scenarioId === "error") {
        return Promise.reject(new Error("モデル一覧を取得できませんでした。"))
      }
      return Promise.resolve(createMasterPersonaAISettingsResponse(scenarioId))
    },
    listMasterPersonaProviderModels: (request) =>
      Promise.resolve({
        provider: request.provider,
        credentialStatus:
          scenarioId === "config-missing"
            ? "missing"
            : request.provider === "lm_studio"
              ? "not_required"
              : "configured",
        status:
          scenarioId === "config-missing"
            ? "credential_missing"
            : request.provider === "lm_studio"
              ? "credential_not_required"
              : "success",
        models:
          scenarioId === "config-missing"
            ? []
            : [{ modelId: "fake-model", label: "fake-model" }]
      }),
    saveMasterPersonaAISettings: (request) => {
      if (scenarioId === "error") {
        return Promise.reject(new Error("AI設定の保存に失敗しました。"))
      }
      return Promise.resolve(request)
    },
    previewMasterPersonaGeneration: (request: MasterPersonaPreviewRequest) =>
      Promise.resolve({
        fileName: request.filePath.split("/").pop() ?? "review.esp",
        targetPlugin: "ReviewPlugin.esp",
        status: "previewed",
        candidateCount: 1,
        newlyAddableCount: 1,
        existingCount: 0
      }),
    executeMasterPersonaGeneration: () =>
      Promise.resolve(createMasterPersonaRunStatus("running")),
    getMasterPersonaRunStatus: () =>
      Promise.resolve(createMasterPersonaRunStatus(scenarioId)),
    interruptMasterPersonaGeneration: () =>
      Promise.resolve(createMasterPersonaRunStatus("success")),
    cancelMasterPersonaGeneration: () =>
      Promise.resolve(createMasterPersonaRunStatus("success")),
    updateMasterPersona: (
      request: MasterPersonaUpdateRequest
    ): Promise<MasterPersonaMutationResponse> =>
      Promise.resolve({
        page: createMasterPersonaPage("success"),
        changedEntry: {
          ...createMasterPersonaEntry(request.identityKey),
          personaBody: request.entry.personaBody,
          personaSummary: request.entry.personaSummary ?? "",
          speechStyle: request.entry.speechStyle
        }
      }),
    deleteMasterPersona: (
      request: MasterPersonaDeleteRequest
    ): Promise<MasterPersonaMutationResponse> =>
      Promise.resolve({
        page: createMasterPersonaPage("empty"),
        deletedEntryId: request.identityKey
      })
  }
}

function createJobSetupOptions(
  scenarioId: ReviewFakeApiScenarioId
): TranslationJobSetupOptionsResponse {
  if (scenarioId === "empty") {
    return {
      inputCandidates: [],
      sharedDictionaries: [],
      sharedPersonas: [],
      aiRuntimeOptions: [],
      providerCapabilities: [],
      phaseRuntimeDrafts: []
    }
  }

  const credentialConfigured = scenarioId !== "config-missing"
  const restoredModel = credentialConfigured ? "fake-model" : ""
  return {
    inputCandidates: [
      {
        id: 1001,
        label: "review-input.xtranslator",
        sourceKind: "xtranslator",
        registeredAt: "2026-05-07T00:00:00Z",
        recordCount: 42
      }
    ],
    sharedDictionaries: [{ id: "dict-review", label: "レビュー辞書" }],
    sharedPersonas: [{ id: "persona-review", label: "レビュー用ペルソナ" }],
    aiRuntimeOptions: [
      { provider: "gemini", model: "fake-model", mode: "single_request" }
    ],
    providerCapabilities: [
      {
        provider: "gemini",
        credentialRequirement: "required",
        supportedExecutionModes: ["single_request"],
        supportsBatchMode: true
      },
      {
        provider: "lm_studio",
        credentialRequirement: "not_required",
        supportedExecutionModes: ["single_request"],
        supportsBatchMode: false
      }
    ],
    phaseRuntimeDrafts: REVIEW_PHASE_IDS.map((phaseId) => ({
      phaseId,
      provider: "gemini",
      model: restoredModel,
      credentialStatus: credentialConfigured ? "configured" : "missing",
      executionMode: "single_request",
      batchMode: "disabled"
    }))
  }
}

function createModelListResponse(
  request: ListTranslationJobSetupProviderModelsRequest,
  scenarioId: ReviewFakeApiScenarioId
): ListTranslationJobSetupProviderModelsResponse {
  if (scenarioId === "config-missing") {
    return {
      phaseId: request.phaseId,
      provider: request.provider,
      credentialStatus: "missing",
      requestToken: request.requestToken,
      sourceToken: "",
      status: "credential_missing",
      models: [],
      failureKind: "credential_missing"
    }
  }

  if (scenarioId === "error") {
    return {
      phaseId: request.phaseId,
      provider: request.provider,
      credentialStatus: request.credentialStatus,
      requestToken: request.requestToken,
      sourceToken: "",
      status: "failed",
      models: [],
      failureKind: "provider_unreachable"
    }
  }

  return {
    phaseId: request.phaseId,
    provider: request.provider,
    credentialStatus: request.credentialStatus,
    requestToken: request.requestToken,
    sourceToken: `${request.phaseId}|${request.provider}||${request.requestToken}`,
    status: "success",
    models:
      scenarioId === "empty"
        ? []
        : [{ modelId: "fake-model", label: "fake-model" }]
  }
}

function toPhaseValidationResults(
  selections: TranslationJobSetupPhaseRuntimeValidationSelection[] | undefined
): TranslationJobSetupPhaseValidationResult[] {
  return (selections ?? []).map((selection) => ({
    phaseId: selection.phaseId,
    status:
      selection.model && selection.modelListFreshnessToken ? "pass" : "blocked",
    blockingFailureCategory:
      selection.model && selection.modelListFreshnessToken
        ? undefined
        : "model_list_failed",
    canCreate:
      selection.model !== "" && selection.modelListFreshnessToken !== "",
    modelListState: selection.model ? "success" : "not_updated",
    isModelSelectionStale:
      selection.model !== "" && selection.modelListFreshnessToken === ""
  }))
}

function createValidationResponse(
  request: ValidateTranslationJobSetupRequest
): TranslationJobSetupValidationResponse {
  const phaseResults = toPhaseValidationResults(request.phaseRuntimeSelections)
  const canCreate =
    phaseResults.length > 0 && phaseResults.every((result) => result.canCreate)

  return {
    status: canCreate ? "ready" : "blocked",
    blockingFailureCategory: canCreate ? undefined : "phase_runtime_missing",
    targetSlices: ["word", "persona", "body"],
    validatedAt: "2026-05-07T00:00:00Z",
    canCreate,
    passSlices: canCreate ? ["word", "persona", "body"] : [],
    phaseResults
  }
}

function createJobSetupSummary(
  jobId: number,
  selections: TranslationJobSetupPhaseRuntimeSelection[] = []
): TranslationJobSetupSummaryResponse {
  return {
    jobId,
    jobState: "ready",
    inputSource: "review-input.xtranslator",
    executionSummary: {
      provider: selections[0]?.provider ?? "gemini",
      model: selections[0]?.model ?? "fake-model",
      executionMode: selections[0]?.executionMode ?? "single_request"
    },
    validationPassSlices: ["word", "persona", "body"],
    canStartPhase: true,
    phaseRuntimeSummaries: selections.map((selection) => ({ ...selection }))
  }
}

function createTranslationJobSetupGateway(
  scenarioId: ReviewFakeApiScenarioId
): TranslationJobSetupGatewayContract {
  let latestCreatedSelections: TranslationJobSetupPhaseRuntimeSelection[] = []

  if (scenarioId === "loading") {
    return {
      getTranslationJobSetupOptions: () => new Promise(() => undefined),
      listTranslationJobSetupProviderModels: () => new Promise(() => undefined),
      validateTranslationJobSetup: () => new Promise(() => undefined),
      createTranslationJob: () => new Promise(() => undefined),
      deleteTranslationJobSetupInput: () => new Promise(() => undefined),
      getTranslationJobSetupSummary: () => new Promise(() => undefined)
    }
  }

  return {
    getTranslationJobSetupOptions: () => {
      return Promise.resolve(createJobSetupOptions(scenarioId))
    },
    listTranslationJobSetupProviderModels: (request) =>
      Promise.resolve(createModelListResponse(request, scenarioId)),
    validateTranslationJobSetup: (request) => {
      if (scenarioId === "error") {
        return Promise.reject(new Error("作成前確認に失敗しました。"))
      }
      return Promise.resolve(createValidationResponse(request))
    },
    createTranslationJob: (request: CreateTranslationJobRequest) => {
      const canCreate = (request.phaseRuntimeSelections ?? []).every(
        (selection) => selection.model !== ""
      )
      latestCreatedSelections = canCreate
        ? (request.phaseRuntimeSelections ?? []).map((selection) => ({
            phaseId: selection.phaseId,
            provider: selection.provider,
            model: selection.model,
            credentialStatus: selection.credentialStatus,
            executionMode: selection.executionMode,
            batchMode: selection.batchMode
          }))
        : []
      return Promise.resolve({
        jobId: 9001,
        jobState: canCreate ? "ready" : "draft",
        inputSource: request.inputSource,
        executionSummary: canCreate
          ? {
              provider: request.runtime.provider,
              model: request.runtime.model,
              executionMode: request.runtime.executionMode
            }
          : undefined,
        validationPassSlices: request.validationPassSlices,
        errorKind: canCreate ? undefined : "phase_runtime_missing",
        phaseRuntimeSummaries: request.phaseRuntimeSelections
      })
    },
    deleteTranslationJobSetupInput: (
      request: DeleteTranslationJobSetupInputRequest
    ): Promise<DeleteTranslationJobSetupInputResponse> =>
      Promise.resolve({ deletedInputSourceId: request.inputSourceId }),
    getTranslationJobSetupSummary: (
      request: GetTranslationJobSetupSummaryRequest
    ) =>
      Promise.resolve(
        createJobSetupSummary(request.jobId, latestCreatedSelections)
      )
  }
}

function createBodyTranslationSummary(
  scenarioId: ReviewFakeApiScenarioId,
  jobId: number
): BodyTranslationPhaseSummaryResponse {
  const state = createBodyTranslationState(scenarioId)

  return {
    jobId,
    currentPhase: "body_translation",
    phaseState: state.completed ? "completed" : state.running ? "running" : "ready",
    phaseRunId: 4021,
    startedAt: "2026-05-08T09:00:00Z",
    finishedAt: state.completed ? "2026-05-08T09:12:00Z" : undefined,
    progress: createBodyTranslationProgress(state),
    inputSummary: {
      targetCount: 12,
      skippedReasons: state.configMissing ? ["AI サービス設定不足"] : [],
      inputSnapshotRef: "snapshot:review-body-402",
      dictionaryDigest: "sha256:review-dictionary",
      personaDigest: "sha256:review-persona",
      metadataDigest: "sha256:review-metadata",
      promptDigest: "sha256:review-prompt"
    },
    requestSummary: {
      providerTargetCount: 12,
      exactDictionaryExclusionCount: 0,
      partialDictionaryConstraintCount: 3
    },
    execution: createBodyTranslationExecution(state),
    fieldResults: createBodyTranslationFieldResults(state),
    resultSummary: createBodyTranslationResultSummary(state),
    errorSummary: createBodyTranslationErrorSummary(state),
    actionEnablement: createBodyTranslationActionEnablement(state),
    outputReadiness: createBodyTranslationOutputReadiness(state)
  }
}

interface BodyTranslationReviewState {
  completed: boolean
  configMissing: boolean
  running: boolean
}

function createBodyTranslationState(
  scenarioId: ReviewFakeApiScenarioId
): BodyTranslationReviewState {
  const completed = scenarioId === "success"
  const configMissing = scenarioId === "config-missing"

  return {
    completed,
    configMissing,
    running: scenarioId === "running" || (!completed && !configMissing)
  }
}

function getBodyTranslationProcessedCount(
  state: BodyTranslationReviewState
): number {
  if (state.completed) {
    return 12
  }

  return state.running ? 7 : 0
}

function getBodyTranslationOutputReadinessBlockedReason(
  state: BodyTranslationReviewState
): string | undefined {
  if (state.completed) {
    return undefined
  }

  return state.configMissing
    ? "設定不足のため出力準備を確認できません。"
    : "本文翻訳が実行中です。完了後に確認してください。"
}

function createBodyTranslationProgress(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["progress"] {
  const processedCount = getBodyTranslationProcessedCount(state)

  return {
    percent: state.completed ? 100 : state.running ? 61 : 0,
    processedCount,
    totalCount: 12,
    targetCount: 12,
    translatedCount: processedCount,
    skippedCount: 0,
    currentStep: state.completed
      ? "completed"
      : state.running
        ? "provider_request"
        : "blocked"
  }
}

function createBodyTranslationExecution(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["execution"] {
  return {
    credentialRef: state.configMissing ? "未設定" : "credential 状態分類だけ表示",
    provider: state.configMissing ? "Gemini 設定不足" : "xAI",
    model: state.configMissing
      ? "model 未選択"
      : "grok-4-super-long-model-name-review",
    executionMode: state.running ? "同期実行中" : "同期実行",
    requestUnitCount: state.configMissing ? 0 : 12,
    outputCount: getBodyTranslationProcessedCount(state)
  }
}

function createBodyTranslationFieldResults(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["fieldResults"] {
  if (!state.completed) {
    return []
  }

  return [
    {
      fieldId: 1,
      fieldLabel: "DIAL 0001",
      sourceExcerpt: "We should leave before the storm reaches town.",
      translatedText: "嵐が町へ来る前に出発しましょう。",
      outputStatus: "ready",
      protectionValidationResult: "pass",
      retryCount: 0
    },
    {
      fieldId: 2,
      fieldLabel: "BOOK 0002",
      sourceExcerpt: "A short note for review.",
      translatedText: "確認用の短いメモです。",
      outputStatus: "ready",
      protectionValidationResult: "pass",
      retryCount: 0
    }
  ]
}

function createBodyTranslationResultSummary(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["resultSummary"] {
  const outputCount = getBodyTranslationProcessedCount(state)

  return {
    translatedCount: outputCount,
    failedCount: 0,
    skippedCount: 0,
    protectionFailedCount: 0,
    outputReadyCount: state.completed ? 12 : 0,
    outputCount
  }
}

function createBodyTranslationErrorSummary(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["errorSummary"] {
  if (!state.configMissing) {
    return undefined
  }

  return {
    errorKind: "secret_redacted",
    reason: "API キーと model を確認してください。secret 本体は表示しません。",
    retryable: true,
    isRedacted: true
  }
}

function createBodyTranslationActionEnablement(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["actionEnablement"] {
  return {
    canStart: state.configMissing,
    startBlockedReason: state.configMissing
      ? "API キーと model を確認してください。"
      : "本文翻訳は開始済みです。",
    canPause: state.running,
    pauseBlockedReason: state.running ? undefined : "実行中ではありません。",
    canResume: false,
    resumeBlockedReason: state.running
      ? "実行中のため再開は不要です。"
      : "再開対象ではありません。",
    canRetry: false,
    retryBlockedReason: "再試行対象ではありません。",
    canCancel: state.running,
    cancelBlockedReason: state.running ? undefined : "実行中ではありません。",
    canCheckOutputReadiness: true,
    outputReadinessBlockedReason:
      getBodyTranslationOutputReadinessBlockedReason(state)
  }
}

function createBodyTranslationOutputReadiness(
  state: BodyTranslationReviewState
): BodyTranslationPhaseSummaryResponse["outputReadiness"] {
  return {
    ready: state.completed,
    blockedReason: getBodyTranslationOutputReadinessBlockedReason(state),
    errorKind: state.completed ? undefined : "output_readiness_blocked",
    completedFieldCount: getBodyTranslationProcessedCount(state),
    statusConsistent: state.completed
  }
}

function createBodyTranslationReadiness(
  scenarioId: ReviewFakeApiScenarioId,
  jobId: number
): BodyTranslationOutputReadinessResponse {
  const summary = createBodyTranslationSummary(scenarioId, jobId)

  return {
    jobId,
    currentPhase: summary.currentPhase,
    phaseState: summary.phaseState,
    ready: summary.outputReadiness.ready,
    blockedReason: summary.outputReadiness.blockedReason,
    errorKind: summary.outputReadiness.errorKind,
    completedFieldCount: summary.outputReadiness.completedFieldCount,
    statusConsistent: summary.outputReadiness.statusConsistent,
    outputCount: summary.execution.outputCount
  }
}

function createBodyTranslationPhaseGateway(
  scenarioId: ReviewFakeApiScenarioId
): BodyTranslationPhaseGatewayContract {
  if (scenarioId === "loading") {
    return {
      getBodyTranslationPhaseSummary: () => new Promise(() => undefined),
      startBodyTranslationPhase: () => new Promise(() => undefined),
      pauseBodyTranslationPhase: () => new Promise(() => undefined),
      resumeBodyTranslationPhase: () => new Promise(() => undefined),
      retryBodyTranslationPhase: () => new Promise(() => undefined),
      cancelBodyTranslationPhase: () => new Promise(() => undefined),
      getBodyTranslationOutputReadiness: () => new Promise(() => undefined)
    }
  }

  const command = (jobId: number): BodyTranslationPhaseCommandResponse => ({
    ...createBodyTranslationSummary("running", jobId),
    retryable: true
  })

  return {
    getBodyTranslationPhaseSummary: ({ jobId }) =>
      Promise.resolve(createBodyTranslationSummary(scenarioId, jobId)),
    startBodyTranslationPhase: ({ jobId }) => Promise.resolve(command(jobId)),
    pauseBodyTranslationPhase: ({ jobId }) => Promise.resolve(command(jobId)),
    resumeBodyTranslationPhase: ({ jobId }) => Promise.resolve(command(jobId)),
    retryBodyTranslationPhase: ({ jobId }) => Promise.resolve(command(jobId)),
    cancelBodyTranslationPhase: ({ jobId }) => Promise.resolve(command(jobId)),
    getBodyTranslationOutputReadiness: ({ jobId }) =>
      Promise.resolve(createBodyTranslationReadiness(scenarioId, jobId))
  }
}

function createTranslationOutputReview(
  selectedJobId?: number
): TranslationOutputReviewResponse {
  const selectedId = selectedJobId ?? 402

  return {
    completedJobs: [
      {
        jobId: 402,
        jobStatus: "completed",
        artifactStatus: "not_generated",
        outputReady: true,
        translatedCount: 12
      }
    ],
    hasSelectedJob: selectedJobId !== undefined,
    selectedJob: {
      jobId: selectedId,
      jobStatus: "completed",
      bodyPhaseStatus: "completed",
      outputReady: true,
      resultSummary: {
        translatedCount: 12,
        rowCount: 12,
        inputProvenance: {
          inputSnapshotDigest: "sha256:review-body-402",
          sourceFileDigest: "sha256:review-source"
        }
      }
    },
    outputReadiness: {
      ready: true,
      retryable: false
    },
    artifactStatus: {
      artifactId: 0,
      status: "not_generated",
      rowCount: 0,
      currentVersion: false
    }
  }
}

function createTranslationOutputArtifactGateway(): TranslationOutputArtifactGatewayContract {
  const diffPreview: TranslationOutputDiffPreviewResponse = {
    jobId: 402,
    artifactId: 0,
    rows: [],
    compatibilitySummary: {
      passed: true,
      warningCount: 0,
      rejectCount: 0
    }
  }

  return {
    getTranslationOutputReview: ({ selectedJobId }) =>
      Promise.resolve(createTranslationOutputReview(selectedJobId)),
    getTranslationOutputDiffPreview: () => Promise.resolve(diffPreview),
    generateXTranslatorOutputArtifact: (request) =>
      Promise.resolve({
        jobId: request.jobId,
        artifactId: 9101,
        artifactStatus: "generated",
        rowCount: 12,
        targetGame: request.targetGame
      }),
    regenerateXTranslatorOutputArtifact: (request) =>
      Promise.resolve({
        jobId: request.jobId,
        artifactId: request.artifactId,
        artifactStatus: "generated",
        rowCount: 12,
        targetGame: request.targetGame
      })
  }
}

export function createDefaultReviewFakeApiGatewayRegistry(): ReviewFakeApiGatewayRegistry {
  return {
    bodyTranslationPhase: (context) =>
      createBodyTranslationPhaseGateway(context.scenarioId),
    masterPersona: (context) => createMasterPersonaGateway(context.scenarioId),
    providerSettings: (context) =>
      createProviderSettingsGateway(context.scenarioId),
    translationJobSetup: (context) =>
      createTranslationJobSetupGateway(context.scenarioId),
    translationOutputArtifact: () => createTranslationOutputArtifactGateway()
  }
}
