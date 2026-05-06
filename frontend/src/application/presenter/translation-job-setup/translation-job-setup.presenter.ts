import type {
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupPhaseId,
  TranslationJobSetupPhaseRuntimeSelection,
  TranslationJobSetupPhaseRuntimeSummary,
  TranslationJobSetupPhaseValidationResult,
  TranslationJobSetupProviderCapability,
  TranslationJobSetupProviderModelOption,
  TranslationJobSetupScreenState,
  TranslationJobSetupScreenViewModel
} from "@application/gateway-contract/translation-job-setup"

const PHASE_LABELS: Record<TranslationJobSetupPhaseId, string> = {
  word_translation: "単語翻訳",
  npc_persona_generation: "NPC ペルソナ生成",
  text_translation: "本文翻訳"
}

const PROVIDER_LABELS: Record<string, string> = {
  gemini: "Gemini",
  xai: "xAI",
  lm_studio: "LM Studio",
  openai: "OpenAI",
  openai_compatible: "OpenAI Compatible",
  "openai-compatible": "OpenAI Compatible",
  anthropic: "Anthropic"
}

export const VALIDATION_LABELS = {
  pass: "validation pass",
  fail: "validation fail",
  warning: "validation warning",
  not_ready: "未完了",
  ready: "確認済み",
  running: "確認中"
} as const

const CREATE_ERROR_LABELS: Record<string, string> = {
  phase_runtime_missing: "翻訳段階ごとの設定が不足しています。",
  required_setting_missing: "必須設定が不足しています。",
  input_not_found: "入力データを確認してください。",
  cache_missing: "入力確認へ戻って必要データを再構築してください。",
  foundation_ref_missing: "共通基盤の参照が不足しています。",
  credential_missing: "APIキー状態を確認してください。",
  model_list_credential_missing: "APIキーを登録してからモデル一覧を更新してください。",
  model_list_failed: "モデル一覧の取得に失敗しています。",
  model_selection_stale: "モデル一覧を更新したため、モデルを選び直してください。",
  provider_mode_unsupported: "選択した AI サービスでは現在の実行方法を使えません。",
  provider_unreachable: "AI サービスへ接続できませんでした。",
  duplicate_job_for_input:
    "同じ入力データの翻訳ジョブがあります。必要なら Job Management で確認してください。",
  validation_stale: "設定が変わったため、作成前確認をやり直してください。",
  partial_create_failed: "翻訳ジョブの作成途中で失敗しました。",
  ready_required: "作成前確認が完了してから実行してください。"
}

export interface TranslationJobSetupPhaseCardViewModel {
  phaseId: TranslationJobSetupPhaseId
  phaseLabel: string
  provider: string
  providerLabel: string
  providerOptions: Array<{ value: string; label: string }>
  credentialStatusLabel: string
  credentialStatusTone: "neutral" | "warning" | "success"
  showCredentialStatus: boolean
  showCredentialWarning: boolean
  modelListButtonEnabled: boolean
  modelListButtonLabel: string
  modelListButtonAriaLabel: string
  isModelListRefreshing: boolean
  modelListStatusText: string
  modelOptions: TranslationJobSetupProviderModelOption[]
  showModelSelect: boolean
  modelSelectEnabled: boolean
  selectedModel: string
  showBatchToggle: boolean
  batchEnabled: boolean
  batchHelpText: string
  statusLabel: string
  statusTone: "neutral" | "warning" | "success"
  helperText: string
}

export interface TranslationJobSetupSummaryPhaseViewModel {
  phaseId: TranslationJobSetupPhaseId
  phaseLabel: string
  providerLabel: string
  model: string
  credentialStatusLabel: string
  batchLabel: string
}

export interface TranslationJobSetupExtendedViewModel
  extends TranslationJobSetupScreenViewModel {
  phaseCards: TranslationJobSetupPhaseCardViewModel[]
  summaryPhaseCards: TranslationJobSetupSummaryPhaseViewModel[]
  createSectionTitle: string
  createSectionText: string
  globalBlockedReasons: string[]
}

function formatProviderLabel(provider: string): string {
  return PROVIDER_LABELS[provider] ?? provider
}

function findPhaseSelection(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): TranslationJobSetupPhaseRuntimeSelection | null {
  return (
    state.phaseRuntimeSelections?.find((selection) => selection.phaseId === phaseId) ??
    null
  )
}

function findModelList(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): ListTranslationJobSetupProviderModelsResponse | null {
  return (
    state.providerModelLists?.find((item) => item.phaseId === phaseId) ?? null
  )
}

function findCapability(
  state: TranslationJobSetupScreenState,
  provider: string
): TranslationJobSetupProviderCapability | null {
  return (
    state.options?.providerCapabilities?.find(
      (capability) => capability.provider === provider
    ) ?? null
  )
}

function findPhaseValidation(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
): TranslationJobSetupPhaseValidationResult | null {
  return (
    state.validationResult?.phaseResults?.find(
      (result) => result.phaseId === phaseId
    ) ?? null
  )
}

function findPhaseSummary(
  summaries: TranslationJobSetupPhaseRuntimeSummary[] | undefined,
  phaseId: TranslationJobSetupPhaseId
): TranslationJobSetupPhaseRuntimeSummary | null {
  return summaries?.find((summary) => summary.phaseId === phaseId) ?? null
}

function isCredentialMissing(
  selection: TranslationJobSetupPhaseRuntimeSelection | null
): boolean {
  return selection?.credentialStatus === "missing"
}

function isModelListUsable(
  modelList: ListTranslationJobSetupProviderModelsResponse | null
): boolean {
  return (
    modelList?.status === "success" ||
    modelList?.status === "credential_not_required"
  )
}

function buildCredentialStatusLabel(
  selection: TranslationJobSetupPhaseRuntimeSelection | null
): string {
  if (!selection) {
    return "未設定"
  }

  if (selection.credentialStatus === "configured") {
    return "登録済み"
  }

  if (selection.credentialStatus === "not_required") {
    return "不要"
  }

  return "未設定"
}

function buildModelListStatusText(
  modelList: ListTranslationJobSetupProviderModelsResponse | null,
  selection: TranslationJobSetupPhaseRuntimeSelection | null
): string {
  if (!selection || selection.provider === "") {
    return "AI サービスを選んでください。"
  }

  if (selection.credentialStatus === "missing") {
    return "AIサービス設定が未完了です。設定が必要です。"
  }

  switch (modelList?.status) {
    case "loading":
      return "モデル一覧を更新しています。"
    case "success":
    case "credential_not_required":
      return "モデル一覧を更新しました。"
    case "credential_missing":
      return "AIサービス設定が未完了です。設定が必要です。"
    case "failed":
      return "モデル一覧を取得できませんでした。時間をおいて再実行してください。"
    case "not_updated":
    default:
      return "モデル一覧を更新してください。"
  }
}

function buildPhaseStatus(
  selection: TranslationJobSetupPhaseRuntimeSelection | null,
  modelList: ListTranslationJobSetupProviderModelsResponse | null,
  validation: TranslationJobSetupPhaseValidationResult | null
): { label: string; tone: "neutral" | "warning" | "success"; helper: string } {
  if (!selection || selection.provider === "") {
    return {
      label: "未設定",
      tone: "warning",
      helper: "AI サービスを選んでください。"
    }
  }

  if (selection.credentialStatus === "missing") {
    return {
      label: "設定が必要",
      tone: "warning",
      helper: "AIサービス設定が未完了です。設定が必要です。"
    }
  }

  if (modelList?.status === "loading") {
    return {
      label: "更新中",
      tone: "neutral",
      helper: "モデル一覧を更新しています。"
    }
  }

  if (modelList?.status === "failed") {
    return {
      label: "取得失敗",
      tone: "warning",
      helper: "モデル一覧の取得に失敗しました。"
    }
  }

  if (selection.model === "") {
    return {
      label: "モデル未選択",
      tone: "warning",
      helper: "モデル一覧を更新して、使うモデルを選んでください。"
    }
  }

  if (validation?.canCreate === false) {
    return {
      label: "要確認",
      tone: "warning",
      helper:
        validation.blockingFailureCategory ??
        "作成前確認で不足があります。"
    }
  }

  return {
    label: "設定済み",
    tone: "success",
    helper: "翻訳ジョブ作成に使う設定です。"
  }
}

function buildProviderOptions(state: TranslationJobSetupScreenState): Array<{
  value: string
  label: string
}> {
  return (
    state.options?.providerCapabilities?.map((capability) => ({
      value: capability.provider,
      label: formatProviderLabel(capability.provider)
    })) ?? []
  )
}

function buildPhaseCards(
  state: TranslationJobSetupScreenState
): TranslationJobSetupPhaseCardViewModel[] {
  const providerOptions = buildProviderOptions(state)

  return (Object.keys(PHASE_LABELS) as TranslationJobSetupPhaseId[]).map(
    (phaseId) => {
      const selection = findPhaseSelection(state, phaseId)
      const modelList = findModelList(state, phaseId)
      const capability = findCapability(state, selection?.provider ?? "")
      const validation = findPhaseValidation(state, phaseId)
      const phaseStatus = buildPhaseStatus(selection, modelList, validation)

      return {
        phaseId,
        phaseLabel: PHASE_LABELS[phaseId],
        provider: selection?.provider ?? "",
        providerLabel: formatProviderLabel(selection?.provider ?? ""),
        providerOptions,
        credentialStatusLabel: buildCredentialStatusLabel(selection),
        credentialStatusTone:
          selection?.credentialStatus === "configured"
            ? "success"
            : selection?.credentialStatus === "not_required"
              ? "neutral"
              : "warning",
        showCredentialStatus: true,
        showCredentialWarning:
          isCredentialMissing(selection) ||
          modelList?.status === "credential_missing",
        modelListButtonEnabled:
          state.phase !== "creating" &&
          modelList?.status !== "loading" &&
          selection?.provider !== "" &&
          selection?.credentialStatus !== "missing",
        modelListButtonLabel: "モデル一覧を更新",
        modelListButtonAriaLabel:
          modelList?.status === "loading"
            ? `${PHASE_LABELS[phaseId]}のモデル一覧を更新中`
            : `${PHASE_LABELS[phaseId]}のモデル一覧を更新`,
        isModelListRefreshing: modelList?.status === "loading",
        modelListStatusText: buildModelListStatusText(modelList, selection),
        modelOptions: isModelListUsable(modelList) ? modelList?.models ?? [] : [],
        showModelSelect: isModelListUsable(modelList),
        modelSelectEnabled:
          state.phase !== "creating" &&
          isModelListUsable(modelList) &&
          selection?.credentialStatus !== "missing",
        selectedModel: selection?.model ?? "",
        showBatchToggle: capability?.supportsBatchMode === true,
        batchEnabled: selection?.batchMode === "enabled",
        batchHelpText: "API利用料が安くなる場合があります。",
        statusLabel: phaseStatus.label,
        statusTone: phaseStatus.tone,
        helperText: phaseStatus.helper
      }
    }
  )
}

function buildGlobalBlockedReasons(
  state: TranslationJobSetupScreenState,
  phaseCards: TranslationJobSetupPhaseCardViewModel[]
): string[] {
  if (state.summary) {
    return []
  }

  const reasons: string[] = []
  if (state.selectedInputSourceId === null) {
    reasons.push("入力データを選んでください。")
  }

  for (const phaseCard of phaseCards) {
    if (phaseCard.statusTone === "warning") {
      reasons.push(`${phaseCard.phaseLabel}: ${phaseCard.helperText}`)
    }
  }

  if (state.validationState === "running") {
    reasons.push("作成前確認を更新しています。")
  } else if (!state.validationResult) {
    reasons.push("3 つの翻訳段階が揃うと作成前確認を実行します。")
  } else if (!state.validationResult.canCreate) {
    reasons.push(
      state.validationResult.blockingFailureCategory ??
        "作成前確認に失敗しています。"
    )
  }

  if (state.createErrorKind) {
    reasons.push(
      CREATE_ERROR_LABELS[state.createErrorKind] ?? state.createErrorKind
    )
  }

  return Array.from(new Set(reasons))
}

function isPhaseDrivenState(state: TranslationJobSetupScreenState): boolean {
  return (
    (state.phaseRuntimeSelections?.length ?? 0) > 0 ||
    (state.options?.providerCapabilities?.length ?? 0) > 0 ||
    (state.options?.phaseRuntimeDrafts?.length ?? 0) > 0
  )
}

function findSelectedRuntimeOption(
  state: TranslationJobSetupScreenState
): { provider: string; model: string; mode: string } | null {
  return (
    state.options?.aiRuntimeOptions.find(
      (option) =>
        [option.provider, option.model, option.mode].join("::") ===
        state.selectedRuntimeKey
    ) ?? null
  )
}

function resolveAvailableCredentialRefs(state: TranslationJobSetupScreenState) {
  const selectedRuntimeOption = findSelectedRuntimeOption(state)
  const credentialRefs = state.options?.credentialRefs ?? []
  if (!selectedRuntimeOption) {
    return credentialRefs
  }

  const providerMatches = credentialRefs.filter(
    (credential) => credential.provider === selectedRuntimeOption.provider
  )

  return providerMatches.length > 0 ? providerMatches : credentialRefs
}

function buildLegacyValidationStatusText(
  state: TranslationJobSetupScreenState
): string {
  if (state.validationState === "running") {
    return "validation を実行しています。完了後に pass / fail / warning を更新します。"
  }

  if (state.validationState === "stale") {
    return "設定を変更したため validation が失効しました。create 前に再実行が必要です。"
  }

  if (state.validationState === "not-run" || !state.validationResult) {
    return "validation 未実行です。入力、runtime、credential を確認して実行してください。"
  }

  const label =
    VALIDATION_LABELS[
      state.validationResult.status as keyof typeof VALIDATION_LABELS
    ] ?? state.validationResult.status
  const sliceText =
    state.validationResult.targetSlices.length > 0
      ? `対象断面: ${state.validationResult.targetSlices.join(" / ")}`
      : "対象断面はありません。"
  const failureText = state.validationResult.blockingFailureCategory
    ? ` 失敗理由: ${state.validationResult.blockingFailureCategory}`
    : ""

  return `${label} / ${sliceText}${failureText}`
}

function buildLegacyBlockedReasons(state: TranslationJobSetupScreenState): string[] {
  const reasons: string[] = []

  if (state.summary) {
    return reasons
  }

  if (state.validationState === "not-run") {
    reasons.push("validation 未実行です。")
  }

  if (state.validationState === "stale" || state.dirty) {
    reasons.push("validation が失効しています。")
  }

  if (!state.validationResult?.canCreate) {
    reasons.push("blocking failure を解消するまで create できません。")
  }

  if (!state.selectedCredentialRef) {
    reasons.push("credential 参照を選択してください。")
  }

  return Array.from(new Set(reasons))
}

function buildLegacyCreateStatusText(state: TranslationJobSetupScreenState): string {
  if (state.phase === "creating") {
    return "translation job を作成しています。成功後は read-only summary へ切り替えます。"
  }

  if (state.summary) {
    return "create 成功済みです。ready job summary を read-only で表示しています。"
  }

  if (state.createErrorKind) {
    return CREATE_ERROR_LABELS[state.createErrorKind] ?? state.createErrorKind
  }

  return "validation が fresh かつ create 可能な時だけ job を作成できます。"
}

function buildLegacyCredentialStateText(
  availableCredentialRefs: Array<{
    provider: string
    credentialRef: string
    isConfigured: boolean
    isMissingSecret: boolean
  }>,
  selectedCredentialRef: string
): string {
  if (availableCredentialRefs.length === 0) {
    return "credential 参照はありません。"
  }

  const selectedCredential = availableCredentialRefs.find(
    (credential) => credential.credentialRef === selectedCredentialRef
  )
  if (!selectedCredential) {
    return "credential 参照を選択してください。"
  }

  if (!selectedCredential.isConfigured) {
    return "credential は未設定です。"
  }

  if (
    selectedCredential.isMissingSecret &&
    selectedCredential.provider !== "lm_studio"
  ) {
    return "credential 参照はありますが secret が不足しています。"
  }

  return "credential 参照は設定済みです。"
}

function buildSummaryPhaseCards(
  state: TranslationJobSetupScreenState
): TranslationJobSetupSummaryPhaseViewModel[] {
  return (Object.keys(PHASE_LABELS) as TranslationJobSetupPhaseId[]).map(
    (phaseId) => {
      const summary = findPhaseSummary(state.summary?.phaseRuntimeSummaries, phaseId)
      return {
        phaseId,
        phaseLabel: PHASE_LABELS[phaseId],
        providerLabel: formatProviderLabel(summary?.provider ?? "-"),
        model: summary?.model ?? "-",
        credentialStatusLabel: buildCredentialStatusLabel(
          summary
            ? {
                phaseId,
                provider: summary.provider,
                model: summary.model,
                credentialRef: summary.credentialRef,
                credentialStatus: summary.credentialStatus,
                executionMode: summary.executionMode,
                batchMode: summary.batchMode,
                modelListSourceToken: summary.modelListSourceToken
              }
            : null
        ),
        batchLabel:
          summary?.batchMode === "enabled"
            ? "使う"
            : summary?.batchMode === "disabled"
              ? "使わない"
              : "対象外"
      }
    }
  )
}

function buildCreateSectionText(
  state: TranslationJobSetupScreenState,
  blockedReasons: string[]
): string {
  if (state.phase === "creating") {
    return "翻訳ジョブを作成しています。"
  }

  if (state.summary) {
    return "作成済みの設定内容を読み取り専用で表示しています。"
  }

  if (blockedReasons.length === 0) {
    return "3 つの翻訳段階の不足はありません。次へ進めます。"
  }

  return "不足を解消すると、次へ進めます。"
}

function findSelectedInputCandidate(state: TranslationJobSetupScreenState) {
  return (
    state.options?.inputCandidates.find(
      (candidate) => candidate.id === state.selectedInputSourceId
    ) ?? null
  )
}

function canCreateInReadyState(state: TranslationJobSetupScreenState): boolean {
  return (
    !state.summary &&
    state.phase === "ready" &&
    state.validationState === "fresh" &&
    !state.dirty &&
    state.validationResult?.canCreate === true
  )
}

function buildValidationStatusLabel(
  state: TranslationJobSetupScreenState,
  legacyMode: boolean
): string {
  if (legacyMode) {
    return state.validationResult?.status
      ? (VALIDATION_LABELS[
          state.validationResult.status as keyof typeof VALIDATION_LABELS
        ] ?? state.validationResult.status)
      : "validation 未実行"
  }

  if (state.validationState === "running") {
    return VALIDATION_LABELS.running
  }

  return state.validationResult?.canCreate
    ? VALIDATION_LABELS.ready
    : VALIDATION_LABELS.not_ready
}

function buildPhaseDrivenValidationStatusText(
  state: TranslationJobSetupScreenState
): string {
  return (
    state.validationResult?.blockingFailureCategory ??
    (state.validationState === "running"
      ? "作成前確認を更新しています。"
      : "不足があると作成できません。")
  )
}

interface TranslationJobSetupDerivedState {
  legacyMode: boolean
  selectedRuntimeOption: { provider: string; model: string; mode: string } | null
  availableCredentialRefs: Array<{
    provider: string
    credentialRef: string
    isConfigured: boolean
    isMissingSecret: boolean
  }>
  phaseCards: TranslationJobSetupPhaseCardViewModel[]
  globalBlockedReasons: string[]
  canCreate: boolean
  canValidate: boolean
  validationStatusLabel: string
  validationStatusText: string
  createStatusText: string
  credentialStateText: string
}

function buildLegacyDerivedState(
  state: TranslationJobSetupScreenState
): TranslationJobSetupDerivedState {
  const selectedRuntimeOption = findSelectedRuntimeOption(state)
  const availableCredentialRefs = resolveAvailableCredentialRefs(state)
  const globalBlockedReasons = buildLegacyBlockedReasons(state)

  return {
    legacyMode: true,
    selectedRuntimeOption,
    availableCredentialRefs,
    phaseCards: [],
    globalBlockedReasons,
    canCreate: canCreateInReadyState(state),
    canValidate:
      !state.summary &&
      state.phase === "ready" &&
      state.selectedInputSourceId !== null &&
      selectedRuntimeOption !== null &&
      state.selectedCredentialRef !== "",
    validationStatusLabel: buildValidationStatusLabel(state, true),
    validationStatusText: buildLegacyValidationStatusText(state),
    createStatusText: buildLegacyCreateStatusText(state),
    credentialStateText: buildLegacyCredentialStateText(
      availableCredentialRefs,
      state.selectedCredentialRef
    )
  }
}

function buildPhaseDrivenDerivedState(
  state: TranslationJobSetupScreenState
): TranslationJobSetupDerivedState {
  const phaseCards = buildPhaseCards(state)
  const globalBlockedReasons = buildGlobalBlockedReasons(state, phaseCards)

  return {
    legacyMode: false,
    selectedRuntimeOption: null,
    availableCredentialRefs: [],
    phaseCards,
    globalBlockedReasons,
    canCreate:
      canCreateInReadyState(state) && globalBlockedReasons.length === 0,
    canValidate: false,
    validationStatusLabel: buildValidationStatusLabel(state, false),
    validationStatusText: buildPhaseDrivenValidationStatusText(state),
    createStatusText: buildCreateSectionText(state, globalBlockedReasons),
    credentialStateText: ""
  }
}

function buildDerivedState(
  state: TranslationJobSetupScreenState
): TranslationJobSetupDerivedState {
  return isPhaseDrivenState(state)
    ? buildPhaseDrivenDerivedState(state)
    : buildLegacyDerivedState(state)
}

export class TranslationJobSetupPresenter {
  toViewModel(
    state: TranslationJobSetupScreenState,
    isGatewayConnected: boolean
  ): TranslationJobSetupExtendedViewModel {
    const selectedInputCandidate = findSelectedInputCandidate(state)
    const derivedState = buildDerivedState(state)

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      selectedInputCandidate,
      selectedRuntimeOption: derivedState.selectedRuntimeOption,
      availableCredentialRefs: derivedState.availableCredentialRefs,
      phaseValidationResults: state.validationResult?.phaseResults?.map((result) => ({
        ...result
      })),
      phaseRuntimeSummaries: state.summary?.phaseRuntimeSummaries?.map((summary) => ({
        ...summary
      })),
      selectedInputLabel: selectedInputCandidate?.label ?? "未選択",
      selectedInputSourceKind: selectedInputCandidate?.sourceKind ?? "-",
      selectedInputRecordCountLabel: selectedInputCandidate
        ? `${selectedInputCandidate.recordCount.toLocaleString("ja-JP")} 件`
        : "-",
      selectedInputRegisteredAtLabel: selectedInputCandidate?.registeredAt
        ? new Date(selectedInputCandidate.registeredAt).toLocaleString("ja-JP")
        : "-",
      existingJobSummary: state.options?.existingJob
        ? `job #${state.options.existingJob.jobId} / ${state.options.existingJob.status} / ${state.options.existingJob.inputSource}`
        : "既存 job はありません。",
      dictionaryLabels:
        state.options?.sharedDictionaries.map((option) => option.label) ?? [],
      personaLabels:
        state.options?.sharedPersonas.map((option) => option.label) ?? [],
      validationStatusLabel: derivedState.validationStatusLabel,
      validationStatusText: derivedState.validationStatusText,
      createStatusText: derivedState.createStatusText,
      blockedReasons: derivedState.globalBlockedReasons,
      canValidate: derivedState.canValidate,
      canCreate: derivedState.canCreate,
      isLoading: state.phase === "loading",
      isValidating: state.phase === "validating",
      isCreating: state.phase === "creating",
      hasExistingJob: Boolean(state.options?.existingJob),
      showCacheMissingGuidance:
        state.validationResult?.blockingFailureCategory?.toLowerCase() ===
        "cache missing",
      credentialStateText: derivedState.credentialStateText,
      phaseCards: derivedState.phaseCards,
      summaryPhaseCards: buildSummaryPhaseCards(state),
      createSectionTitle: "作成前確認",
      createSectionText: derivedState.createStatusText,
      globalBlockedReasons: derivedState.globalBlockedReasons
    }
  }
}
