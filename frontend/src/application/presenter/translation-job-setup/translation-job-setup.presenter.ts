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
import { buildModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"

const PHASE_LABELS: Record<TranslationJobSetupPhaseId, string> = {
  word_translation: "単語翻訳",
  npc_persona_generation: "NPC ペルソナ生成",
  text_translation: "本文翻訳"
}

const PROVIDER_LABELS: Record<string, string> = {
  gemini: "Gemini",
  xai: "xAI",
  lm_studio: "LM Studio"
}

export const VALIDATION_LABELS = {
  pass: "確認済み",
  fail: "確認失敗",
  warning: "要確認",
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
  model_list_credential_missing:
    "APIキーを登録してからモデル一覧を更新してください。",
  model_list_failed: "モデル一覧の取得に失敗しています。",
  model_selection_stale:
    "モデル一覧を更新したため、モデルを選び直してください。",
  provider_mode_unsupported:
    "選択した AI サービスでは現在の実行方法を使えません。",
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
  credentialWarningText: string
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

export interface TranslationJobSetupExtendedViewModel extends TranslationJobSetupScreenViewModel {
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
    state.providerModelLists?.find((item) => item.phaseId === phaseId) ?? null
  )
}

function findModelSettingsCard(
  state: TranslationJobSetupScreenState,
  phaseId: TranslationJobSetupPhaseId
) {
  return (
    state.modelSettingsCards?.find((card) => card.referenceId === phaseId) ??
    null
  )
}

function createFallbackModelSettingsCard(options: {
  phaseId: TranslationJobSetupPhaseId
  selection: TranslationJobSetupPhaseRuntimeSelection | null
  modelList: ListTranslationJobSetupProviderModelsResponse | null
}) {
  const provider = options.selection?.provider ?? ""
  const credentialStatus = options.selection?.credentialStatus ?? "missing"
  return {
    referenceId: options.phaseId,
    provider,
    model: options.selection?.model ?? "",
    credentialStatus,
    modelList: options.modelList
      ? {
          provider: options.modelList.provider,
          credentialStatus: options.modelList.credentialStatus,
          status: options.modelList.status,
          models: options.modelList.models,
          failureKind: options.modelList.failureKind
        }
      : {
          provider,
          credentialStatus,
          status: "not_updated" as const,
          models: []
        },
    saveStatus: options.selection?.model
      ? ("saved" as const)
      : ("clean" as const),
    saveMessage: ""
  }
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
    return "設定済み"
  }

  if (selection.credentialStatus === "not_required") {
    return "APIキー不要"
  }

  return "APIキー未設定"
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
      label: "APIキー未設定",
      tone: "warning",
      helper: "APIキーを設定してください。"
    }
  }

  if (!modelList || modelList.status === "not_updated") {
    return {
      label: "モデル一覧未更新",
      tone: "warning",
      helper: "モデル一覧を更新してください。"
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
      label: "モデル一覧取得失敗",
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
        validation.blockingFailureCategory ?? "作成前確認で不足があります。"
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
  credentialStatus: "configured" | "missing" | "not_required"
}> {
  return (
    state.options?.providerCapabilities?.map((capability) => ({
      value: capability.provider,
      label: formatProviderLabel(capability.provider),
      credentialStatus:
        capability.credentialRequirement === "not_required"
          ? "not_required"
          : "missing"
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
      const modelSettingsCard =
        findModelSettingsCard(state, phaseId) ??
        createFallbackModelSettingsCard({
          phaseId,
          selection,
          modelList
        })
      const cardViewModel = buildModelSettingsCardViewModel({
        state: modelSettingsCard,
        providerOptions,
        refreshDisabled: state.phase === "creating",
        actionDisabled: true,
        titleLabel: PHASE_LABELS[phaseId]
      })

      return {
        phaseId,
        phaseLabel: PHASE_LABELS[phaseId],
        provider: cardViewModel.provider,
        providerLabel: formatProviderLabel(cardViewModel.provider),
        providerOptions: cardViewModel.providerOptions,
        credentialStatusLabel: cardViewModel.credentialStatusLabel,
        credentialStatusTone: cardViewModel.credentialStatusTone,
        showCredentialStatus: cardViewModel.showCredentialStatus,
        showCredentialWarning: cardViewModel.showCredentialWarning,
        credentialWarningText: cardViewModel.credentialWarningText,
        modelListButtonEnabled: cardViewModel.modelListButtonEnabled,
        modelListButtonLabel: cardViewModel.modelListButtonLabel,
        modelListButtonAriaLabel: cardViewModel.modelListButtonAriaLabel,
        isModelListRefreshing: cardViewModel.isModelListRefreshing,
        modelListStatusText: cardViewModel.modelListStatusText,
        modelOptions: cardViewModel.modelOptions,
        showModelSelect: isModelListUsable(modelList),
        modelSelectEnabled:
          state.phase !== "creating" && cardViewModel.modelSelectEnabled,
        selectedModel: cardViewModel.model,
        showBatchToggle: capability?.supportsBatchMode === true,
        batchEnabled: selection?.batchMode === "enabled",
        batchHelpText: "API利用料が安くなる場合があります。",
        statusLabel:
          phaseStatus.label === "設定済み"
            ? cardViewModel.statusLabel
            : phaseStatus.label,
        statusTone:
          phaseStatus.label === "設定済み"
            ? cardViewModel.statusTone
            : phaseStatus.tone,
        helperText:
          phaseStatus.label === "設定済み"
            ? cardViewModel.helperText
            : phaseStatus.helper
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
  } else if (state.validationResult && !state.validationResult.canCreate) {
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

function resolveAvailableCredentialRefs() {
  return []
}

function buildLegacyValidationStatusText(
  state: TranslationJobSetupScreenState
): string {
  if (state.validationState === "running") {
    return "作成前確認を実行しています。"
  }

  if (state.validationState === "stale") {
    return "設定を変更したため、作成前確認の再実行が必要です。"
  }

  if (state.validationState === "not-run" || !state.validationResult) {
    return "作成前確認は未実行です。入力データと AIサービスを確認してください。"
  }

  const label =
    VALIDATION_LABELS[
      state.validationResult.status as keyof typeof VALIDATION_LABELS
    ] ?? state.validationResult.status
  const sliceText =
    state.validationResult.targetSlices.length > 0
      ? `確認対象: ${state.validationResult.targetSlices.join(" / ")}`
      : "対象断面はありません。"
  const failureText = state.validationResult.blockingFailureCategory
    ? ` 失敗理由: ${state.validationResult.blockingFailureCategory}`
    : ""

  return `${label} / ${sliceText}${failureText}`
}

function buildLegacyBlockedReasons(
  state: TranslationJobSetupScreenState
): string[] {
  const reasons: string[] = []

  if (state.summary) {
    return reasons
  }

  if (state.validationState === "not-run") {
    reasons.push("作成前確認が未実行です。")
  }

  if (state.validationState === "stale" || state.dirty) {
    reasons.push("作成前確認の再実行が必要です。")
  }

  if (!state.validationResult?.canCreate) {
    reasons.push("作成できない理由を解消してください。")
  }

  return Array.from(new Set(reasons))
}

function buildLegacyCreateStatusText(
  state: TranslationJobSetupScreenState
): string {
  if (state.phase === "creating") {
    return "翻訳ジョブを作成しています。"
  }

  if (state.summary) {
    return "作成済みの設定内容を読み取り専用で表示しています。"
  }

  if (state.createErrorKind) {
    return CREATE_ERROR_LABELS[state.createErrorKind] ?? state.createErrorKind
  }

  return "作成前確認が完了した時だけ、翻訳ジョブを作成できます。"
}

function sanitizePhaseRuntimeSelectionForView(
  selection: TranslationJobSetupPhaseRuntimeSelection
): TranslationJobSetupPhaseRuntimeSelection {
  return { ...selection }
}

function sanitizePhaseRuntimeSummaryForView(
  summary: TranslationJobSetupPhaseRuntimeSummary
): TranslationJobSetupPhaseRuntimeSummary {
  return { ...summary }
}

function sanitizeSummaryForView(
  summary: TranslationJobSetupScreenState["summary"]
): TranslationJobSetupScreenState["summary"] {
  if (!summary) {
    return null
  }

  return {
    ...summary,
    executionSummary: { ...summary.executionSummary },
    validationPassSlices: [...summary.validationPassSlices],
    phaseRuntimeSummaries: summary.phaseRuntimeSummaries?.map(
      sanitizePhaseRuntimeSummaryForView
    )
  }
}

function sanitizeValidationForView(
  validation: TranslationJobSetupScreenState["validationResult"]
): TranslationJobSetupScreenState["validationResult"] {
  if (!validation) {
    return null
  }

  return {
    ...validation,
    targetSlices: [...validation.targetSlices],
    passSlices: [...validation.passSlices],
    phaseResults: validation.phaseResults?.map((result) => ({ ...result })),
    staleModelListPhaseIds: validation.staleModelListPhaseIds
      ? [...validation.staleModelListPhaseIds]
      : undefined
  }
}

function buildSummaryPhaseCards(
  state: TranslationJobSetupScreenState
): TranslationJobSetupSummaryPhaseViewModel[] {
  return (Object.keys(PHASE_LABELS) as TranslationJobSetupPhaseId[]).map(
    (phaseId) => {
      const summary = findPhaseSummary(
        state.summary?.phaseRuntimeSummaries,
        phaseId
      )
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
                credentialStatus: summary.credentialStatus,
                executionMode: summary.executionMode,
                batchMode: summary.batchMode
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
    return state.validationResult
      ? "作成前確認は完了しています。"
      : "作成前確認はまだ未完了です。"
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

function buildExistingJobSummary(
  state: TranslationJobSetupScreenState
): string {
  const selectedInputCandidate = findSelectedInputCandidate(state)
  const existingJob =
    selectedInputCandidate?.existingJob ?? state.options?.existingJob
  if (!existingJob) {
    return "既存 job はありません。"
  }

  return `job #${existingJob.jobId} / ${existingJob.status} / ${existingJob.inputSource}`
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
  selectedRuntimeOption: {
    provider: string
    model: string
    mode: string
  } | null
  availableCredentialRefs: Array<{
    provider: string
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
  const availableCredentialRefs = resolveAvailableCredentialRefs()
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
    credentialStateText: ""
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
      selectedCredentialRef: "",
      phaseRuntimeSelections: state.phaseRuntimeSelections?.map(
        sanitizePhaseRuntimeSelectionForView
      ),
      validationResult: sanitizeValidationForView(state.validationResult),
      summary: sanitizeSummaryForView(state.summary),
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      selectedInputCandidate,
      selectedRuntimeOption: derivedState.selectedRuntimeOption,
      availableCredentialRefs: [],
      phaseValidationResults: state.validationResult?.phaseResults?.map(
        (result) => ({ ...result })
      ),
      phaseRuntimeSummaries: state.summary?.phaseRuntimeSummaries?.map(
        sanitizePhaseRuntimeSummaryForView
      ),
      selectedInputLabel: selectedInputCandidate?.label ?? "未選択",
      selectedInputSourceKind: selectedInputCandidate?.sourceKind ?? "-",
      selectedInputRecordCountLabel: selectedInputCandidate
        ? `${selectedInputCandidate.recordCount.toLocaleString("ja-JP")} 件`
        : "-",
      selectedInputRegisteredAtLabel: selectedInputCandidate?.registeredAt
        ? new Date(selectedInputCandidate.registeredAt).toLocaleString("ja-JP")
        : "-",
      existingJobSummary: buildExistingJobSummary(state),
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
      hasExistingJob: Boolean(
        selectedInputCandidate?.existingJob ?? state.options?.existingJob
      ),
      showCacheMissingGuidance:
        state.validationResult?.blockingFailureCategory?.toLowerCase() ===
        "cache missing",
      credentialStateText: "",
      phaseCards: derivedState.phaseCards,
      summaryPhaseCards: buildSummaryPhaseCards(state),
      createSectionTitle: "作成前確認",
      createSectionText: derivedState.createStatusText,
      globalBlockedReasons: derivedState.globalBlockedReasons
    }
  }
}
