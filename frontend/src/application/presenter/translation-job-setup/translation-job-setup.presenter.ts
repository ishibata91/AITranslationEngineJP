import {
  createTranslationJobSetupRuntimeKey,
  TranslationJobSetupCredentialReference,
  TranslationJobSetupRuntimeOption,
  TranslationJobSetupScreenState,
  TranslationJobSetupScreenViewModel
} from "@application/gateway-contract/translation-job-setup"

const VALIDATION_LABELS: Record<string, string> = {
  pass: "validation pass",
  fail: "validation fail",
  warning: "validation warning"
}

const CREATE_ERROR_LABELS: Record<string, string> = {
  validation_failed: "validation fail を解消してから create を再実行してください。"
}

function formatRegisteredAtLabel(registeredAt: string | undefined): string {
  if (!registeredAt) {
    return "-"
  }

  const date = new Date(registeredAt)
  if (Number.isNaN(date.getTime())) {
    return registeredAt
  }

  return date.toLocaleString("ja-JP")
}

function findSelectedRuntimeOption(
  state: TranslationJobSetupScreenState
): TranslationJobSetupRuntimeOption | null {
  return (
    state.options?.aiRuntimeOptions.find(
      (option) => createTranslationJobSetupRuntimeKey(option) === state.selectedRuntimeKey
    ) ?? null
  )
}

function resolveAvailableCredentialRefs(
  state: TranslationJobSetupScreenState
): TranslationJobSetupCredentialReference[] {
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

function buildValidationStatusText(state: TranslationJobSetupScreenState): string {
  if (state.validationState === "running") {
    return "validation を実行しています。完了後に pass / fail / warning を更新します。"
  }

  if (state.validationState === "stale") {
    return "設定を変更したため validation が失効しました。create 前に再実行が必要です。"
  }

  if (state.validationState === "not-run" || !state.validationResult) {
    return "validation 未実行です。入力、runtime、credential を確認して実行してください。"
  }

  const label = VALIDATION_LABELS[state.validationResult.status] ?? state.validationResult.status
  const sliceText = state.validationResult.targetSlices.length > 0
    ? `対象断面: ${state.validationResult.targetSlices.join(" / ")}`
    : "対象断面はありません。"
  const failureText = state.validationResult.blockingFailureCategory
    ? ` 失敗理由: ${state.validationResult.blockingFailureCategory}`
    : ""

  return `${label} / ${sliceText}${failureText}`
}

function hasBlockingExistingJob(state: TranslationJobSetupScreenState): boolean {
  if (state.selectedInputSourceId === null) {
    return false
  }

  const existingJob = state.options?.existingJob
  if (!existingJob) {
    return false
  }

  if ((existingJob.inputSourceId ?? 0) > 0) {
    return existingJob.inputSourceId === state.selectedInputSourceId
  }

  const selectedInputCandidate = state.options?.inputCandidates.find(
    (candidate) => candidate.id === state.selectedInputSourceId
  )

  return existingJob.inputSource === selectedInputCandidate?.label
}

function buildBlockedReasons(state: TranslationJobSetupScreenState): string[] {
  const reasons: string[] = []

  if (state.summary) {
    return reasons
  }

  if (hasBlockingExistingJob(state)) {
    reasons.push("既存 job があるため create を無効化しています。")
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

function buildCreateStatusText(state: TranslationJobSetupScreenState): string {
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

function buildExistingJobSummary(state: TranslationJobSetupScreenState): string {
  if (!state.options?.existingJob) {
    return "既存 job はありません。"
  }

  const existingJob = state.options.existingJob
  return `job #${existingJob.jobId} / ${existingJob.status} / ${existingJob.inputSource}`
}

function buildCredentialStateText(
  availableCredentialRefs: TranslationJobSetupCredentialReference[],
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

  if (selectedCredential.isMissingSecret) {
    return "credential 参照はありますが secret が不足しています。"
  }

  return "credential 参照は設定済みです。"
}

export class TranslationJobSetupPresenter {
  toViewModel(
    state: TranslationJobSetupScreenState,
    isGatewayConnected: boolean
  ): TranslationJobSetupScreenViewModel {
    const selectedInputCandidate =
      state.options?.inputCandidates.find(
        (candidate) => candidate.id === state.selectedInputSourceId
      ) ?? null
    const selectedRuntimeOption = findSelectedRuntimeOption(state)
    const availableCredentialRefs = resolveAvailableCredentialRefs(state)
    const blockedReasons = buildBlockedReasons(state)

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      selectedInputCandidate,
      selectedRuntimeOption,
      availableCredentialRefs,
      selectedInputLabel: selectedInputCandidate?.label ?? "未選択",
      selectedInputSourceKind: selectedInputCandidate?.sourceKind ?? "-",
      selectedInputRecordCountLabel:
        selectedInputCandidate ? `${selectedInputCandidate.recordCount.toLocaleString("ja-JP")} 件` : "-",
      selectedInputRegisteredAtLabel: formatRegisteredAtLabel(
        selectedInputCandidate?.registeredAt
      ),
      existingJobSummary: buildExistingJobSummary(state),
      dictionaryLabels:
        state.options?.sharedDictionaries.map((option) => option.label) ?? [],
      personaLabels: state.options?.sharedPersonas.map((option) => option.label) ?? [],
      validationStatusLabel:
        state.validationResult?.status
          ? (VALIDATION_LABELS[state.validationResult.status] ?? state.validationResult.status)
          : "validation 未実行",
      validationStatusText: buildValidationStatusText(state),
      createStatusText: buildCreateStatusText(state),
      blockedReasons,
      canValidate:
        !state.summary &&
        state.phase === "ready" &&
        state.selectedInputSourceId !== null &&
        selectedRuntimeOption !== null &&
        state.selectedCredentialRef !== "",
      canCreate:
        !state.summary &&
        state.phase === "ready" &&
        state.validationState === "fresh" &&
        !state.dirty &&
        state.validationResult?.canCreate === true &&
        !hasBlockingExistingJob(state),
      isLoading: state.phase === "loading",
      isValidating: state.phase === "validating",
      isCreating: state.phase === "creating",
      hasExistingJob: Boolean(state.options?.existingJob),
      showCacheMissingGuidance:
        state.validationResult?.blockingFailureCategory?.toLowerCase().includes("cache") ?? false,
      credentialStateText: buildCredentialStateText(
        availableCredentialRefs,
        state.selectedCredentialRef
      )
    }
  }
}

export { VALIDATION_LABELS }