import type {
  GenerateXTranslatorOutputArtifactRequest,
  GetTranslationOutputDiffPreviewRequest,
  GetTranslationOutputReviewRequest,
  RegenerateXTranslatorOutputArtifactRequest,
  TranslationOutputArtifactCommandResponse,
  TranslationOutputArtifactErrorKind,
  TranslationOutputArtifactErrorSummary,
  TranslationOutputArtifactGatewayContract,
  TranslationOutputDiffPreviewResponse,
  TranslationOutputReviewResponse
} from "@application/gateway-contract/translation-output-artifact"

type TranslationOutputArtifactActionKind = "generate" | "regenerate"

type TranslationOutputArtifactViewState =
  | "loading"
  | "empty"
  | "awaiting_selection"
  | "not_ready"
  | "ready"
  | "generating"
  | "success"
  | "failed"
  | "stale"

interface TranslationOutputArtifactCommandSnapshot {
  jobId: number
  artifactId: number
  artifactStatus: string
  rowCount: number
  filePath?: string
  targetGame: string
  operationKind?: string
  replacedArtifactId?: number
  affectedFieldIds?: number[]
  duplicateRowCreated?: boolean
  errorKind?: TranslationOutputArtifactErrorKind
  errorReason?: string
  retryable?: boolean
  isRedacted?: boolean
}

interface TranslationOutputReviewSnapshot {
  completedJobs: TranslationOutputReviewResponse["completedJobs"]
  selectedJobId: number
  selectedJobStatus: string
  bodyPhaseStatus: string
  outputReady: boolean
  translatedCount: number
  rowCount: number
  inputSnapshotDigest: string
  sourceFileDigest: string
  readiness: boolean
  retryable: boolean
  artifactId: number
  artifactStatus: string
  artifactRowCount: number
  currentVersion: boolean
  rejectionKind?: TranslationOutputArtifactErrorKind
  rejectionReasons: TranslationOutputArtifactErrorSummary[]
}

interface TranslationOutputDiffPreviewSnapshot {
  artifactId: number
  rows: TranslationOutputDiffPreviewResponse["rows"]
  compatibilityPassed: boolean
  compatibilityWarningCount: number
  compatibilityRejectCount: number
}

interface TranslationOutputArtifactActionDisablement {
  actionKind: TranslationOutputArtifactActionKind
  disabled: boolean
  reason?: string
  errorKind?: TranslationOutputArtifactErrorKind
}

interface TranslationOutputArtifactScreenState {
  phase: "idle" | "loading" | "ready" | "submitting"
  viewState: TranslationOutputArtifactViewState
  completedJobs: TranslationOutputReviewResponse["completedJobs"]
  selectedJobId: number | null
  selectedArtifactId: number | null
  review: TranslationOutputReviewSnapshot | null
  diffPreview: TranslationOutputDiffPreviewSnapshot | null
  lastCommand: TranslationOutputArtifactCommandSnapshot | null
  actionDisablements: TranslationOutputArtifactActionDisablement[]
  refreshPending: boolean
  targetGame: string
  outputPath: string
  pathState: "empty" | "valid" | "invalid" | "readonly"
  pathReason: string
  errorMessage: string
  pendingAction: TranslationOutputArtifactActionKind | "refresh" | null
  hasLoaded: boolean
}

interface TranslationOutputArtifactStoreLike {
  snapshot(): TranslationOutputArtifactScreenState
  update(mutator: (draft: TranslationOutputArtifactScreenState) => void): void
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

function createGatewayDisconnectedMessage(): string {
  return "translation-output-artifact gateway が未接続です。"
}

function mapReviewResponse(
  response: TranslationOutputReviewResponse
): TranslationOutputReviewSnapshot {
  return {
    completedJobs: response.completedJobs.map((job) => ({
      ...job,
      outputStatusDistribution: job.outputStatusDistribution
        ? { ...job.outputStatusDistribution }
        : undefined
    })),
    selectedJobId: response.selectedJob.jobId,
    selectedJobStatus: response.selectedJob.jobStatus,
    bodyPhaseStatus: response.selectedJob.bodyPhaseStatus,
    outputReady: response.selectedJob.outputReady,
    translatedCount: response.selectedJob.resultSummary.translatedCount,
    rowCount: response.selectedJob.resultSummary.rowCount,
    inputSnapshotDigest:
      response.selectedJob.resultSummary.inputProvenance.inputSnapshotDigest,
    sourceFileDigest:
      response.selectedJob.resultSummary.inputProvenance.sourceFileDigest,
    readiness: response.outputReadiness.ready,
    retryable: response.outputReadiness.retryable,
    artifactId: response.artifactStatus.artifactId,
    artifactStatus: response.artifactStatus.status,
    artifactRowCount: response.artifactStatus.rowCount,
    currentVersion: response.artifactStatus.currentVersion,
    rejectionKind:
      response.outputReadiness.rejectionKind ??
      response.rejectionReasons?.[0]?.errorKind,
    rejectionReasons:
      response.rejectionReasons?.map((reason) => ({
        ...reason
      })) ?? []
  }
}

function mapCompletedJobs(
  response: TranslationOutputReviewResponse
): TranslationOutputReviewResponse["completedJobs"] {
  return response.completedJobs.map((job) => ({
    ...job,
    outputStatusDistribution: job.outputStatusDistribution
      ? { ...job.outputStatusDistribution }
      : undefined
  }))
}

function mapDiffPreviewResponse(
  response: TranslationOutputDiffPreviewResponse
): TranslationOutputDiffPreviewSnapshot {
  return {
    artifactId: response.artifactId,
    rows: response.rows.map((row) => ({ ...row })),
    compatibilityPassed: response.compatibilitySummary.passed,
    compatibilityWarningCount: response.compatibilitySummary.warningCount,
    compatibilityRejectCount: response.compatibilitySummary.rejectCount
  }
}

function mapCommandResponse(
  response: TranslationOutputArtifactCommandResponse
): TranslationOutputArtifactCommandSnapshot {
  return {
    jobId: response.jobId,
    artifactId: response.artifactId,
    artifactStatus: response.artifactStatus,
    rowCount: response.rowCount,
    filePath: response.filePath,
    targetGame: response.targetGame,
    operationKind: response.operationSummary?.operationKind,
    replacedArtifactId: response.operationSummary?.replacedArtifactId,
    affectedFieldIds: response.operationSummary?.affectedFieldIds
      ? [...response.operationSummary.affectedFieldIds]
      : undefined,
    duplicateRowCreated: response.operationSummary?.duplicateRowCreated,
    errorKind: response.errorSummary?.errorKind,
    errorReason: response.errorSummary?.reason,
    retryable: response.errorSummary?.retryable,
    isRedacted: response.errorSummary?.isRedacted
  }
}

function isReadonlyPath(outputPath: string): boolean {
  const normalized = outputPath.trim()
  return ["/System/", "/usr/", "/bin/", "/sbin/"].some((prefix) =>
    normalized.startsWith(prefix)
  )
}

function resolvePathState(outputPath: string): {
  pathState: TranslationOutputArtifactScreenState["pathState"]
  pathReason: string
} {
  const trimmed = outputPath.trim()
  if (!trimmed) {
    return {
      pathState: "empty",
      pathReason: "出力先 path を入力してください。"
    }
  }

  if (isReadonlyPath(trimmed)) {
    return {
      pathState: "readonly",
      pathReason: "read-only path には出力できません。"
    }
  }

  if (!trimmed.endsWith(".xml")) {
    return {
      pathState: "invalid",
      pathReason: "出力先 path は .xml で終える必要があります。"
    }
  }

  return {
    pathState: "valid",
    pathReason: ""
  }
}

function hasStaleRows(
  diffPreview: TranslationOutputDiffPreviewSnapshot | null
): boolean {
  return (
    diffPreview?.rows.some(
      (row) => Boolean(row.staleReason) || row.canRegenerate
    ) ?? false
  )
}

function resolveRejectionReason(
  rejectionReasons: TranslationOutputArtifactErrorSummary[],
  kind: TranslationOutputArtifactErrorKind | undefined
): string | null {
  if (!kind) {
    return null
  }

  return (
    rejectionReasons.find((reason) => reason.errorKind === kind)?.reason ??
    mapErrorKindToLabel(kind)
  )
}

function mapErrorKindToLabel(kind: TranslationOutputArtifactErrorKind): string {
  switch (kind) {
    case "not_completed":
      return "job が completed ではありません。"
    case "canceled":
      return "Canceled job は出力できません。"
    case "status_mismatch":
      return "field result または status が不整合です。"
    case "missing_required_row_field":
      return "必須 row field が不足しています。"
    case "unknown_output_status":
      return "不明な output status が含まれています。"
    case "xml_serialization_failed":
      return "XML parse failure または serialization failure が発生しました。"
    case "file_write_failed":
      return "出力先への file write に失敗しました。"
    case "artifact_save_failed":
      return "artifact 保存に失敗しました。"
    case "compatibility_rejected":
      return "compatibility reject を解消してください。"
    case "secret_redacted":
      return "公開可能な redacted summary のみ表示しています。"
  }
}

function buildActionDisablements(
  state: TranslationOutputArtifactScreenState,
  isGatewayConnected: boolean
): TranslationOutputArtifactActionDisablement[] {
  const review = state.review
  const diffPreview = state.diffPreview
  const pathInfo = resolvePathState(state.outputPath)
  const sharedReasons: string[] = []

  if (!isGatewayConnected) {
    sharedReasons.push(createGatewayDisconnectedMessage())
  }
  if (!review) {
    sharedReasons.push("completed job を選択してください。")
  }
  if (!state.targetGame.trim()) {
    sharedReasons.push("target game を選択してください。")
  }
  if (pathInfo.pathState !== "valid") {
    sharedReasons.push(pathInfo.pathReason)
  }
  if (review && !review.readiness) {
    sharedReasons.push(
      resolveRejectionReason(review.rejectionReasons, review.rejectionKind) ??
        "output readiness が false です。"
    )
  }
  if (diffPreview && diffPreview.compatibilityRejectCount > 0) {
    sharedReasons.push("compatibility reject を解消してください。")
  }

  const artifactExists = (review?.artifactId ?? 0) > 0
  const stale = Boolean(
    review && (!review.currentVersion || hasStaleRows(diffPreview))
  )
  const generateReasons = [...sharedReasons]
  const regenerateReasons = [...sharedReasons]

  if (artifactExists) {
    generateReasons.push("既存 artifact があるため再出力を使用してください。")
  }

  if (!artifactExists) {
    regenerateReasons.push("既存 artifact がありません。")
  } else if (!stale && review?.artifactStatus !== "failed") {
    regenerateReasons.push("再出力が必要な状態ではありません。")
  }

  return [
    {
      actionKind: "generate",
      disabled: generateReasons.length > 0,
      reason: generateReasons[0],
      errorKind: review?.rejectionKind
    },
    {
      actionKind: "regenerate",
      disabled: regenerateReasons.length > 0,
      reason: regenerateReasons[0],
      errorKind: review?.rejectionKind
    }
  ]
}

function resolveViewState(
  state: TranslationOutputArtifactScreenState
): TranslationOutputArtifactViewState {
  if (state.phase === "loading" && !state.hasLoaded) {
    return "loading"
  }

  if (state.phase === "submitting") {
    return "generating"
  }

  if (state.hasLoaded && state.completedJobs.length === 0) {
    return "empty"
  }

  if (state.hasLoaded && state.selectedJobId === null) {
    return "awaiting_selection"
  }

  if (
    state.lastCommand?.errorKind ||
    state.review?.artifactStatus === "failed"
  ) {
    return "failed"
  }

  if (
    state.review &&
    (!state.review.currentVersion || hasStaleRows(state.diffPreview))
  ) {
    return "stale"
  }

  if (state.review?.artifactStatus === "success") {
    return "success"
  }

  if (state.review && !state.review.readiness) {
    return "not_ready"
  }

  if (state.review?.readiness) {
    return "ready"
  }

  return state.hasLoaded ? "empty" : "loading"
}

export class TranslationOutputArtifactUseCase {
  constructor(
    private readonly gateway: TranslationOutputArtifactGatewayContract | null,
    private readonly store: TranslationOutputArtifactStoreLike
  ) {}

  async load(): Promise<void> {
    const state = this.store.snapshot()
    const pathInfo = resolvePathState(state.outputPath)
    this.store.update((draft) => {
      draft.pathState = pathInfo.pathState
      draft.pathReason = pathInfo.pathReason
      draft.actionDisablements = buildActionDisablements(
        draft,
        this.gateway !== null
      )
      draft.viewState = resolveViewState(draft)
    })

    await this.refresh()
  }

  async setJobId(jobId: number | null): Promise<void> {
    this.store.update((draft) => {
      draft.selectedJobId = jobId
      draft.selectedArtifactId = null
      draft.diffPreview = null
      draft.errorMessage = ""
      draft.lastCommand = null
      draft.actionDisablements = buildActionDisablements(
        draft,
        this.gateway !== null
      )
      draft.viewState = resolveViewState(draft)
    })

    await this.refresh()
  }

  async setArtifactId(artifactId: number | null): Promise<void> {
    this.store.update((draft) => {
      draft.selectedArtifactId = artifactId
    })

    await this.refreshDiffPreview()
  }

  setTargetGame(targetGame: string): void {
    this.store.update((draft) => {
      draft.targetGame = targetGame
      draft.actionDisablements = buildActionDisablements(
        draft,
        this.gateway !== null
      )
      draft.viewState = resolveViewState(draft)
    })
  }

  setOutputPath(outputPath: string): void {
    const pathInfo = resolvePathState(outputPath)
    this.store.update((draft) => {
      draft.outputPath = outputPath
      draft.pathState = pathInfo.pathState
      draft.pathReason = pathInfo.pathReason
      draft.actionDisablements = buildActionDisablements(
        draft,
        this.gateway !== null
      )
      draft.viewState = resolveViewState(draft)
    })
  }

  async refresh(): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.hasLoaded = true
        draft.errorMessage = createGatewayDisconnectedMessage()
        draft.actionDisablements = buildActionDisablements(draft, false)
        draft.viewState = resolveViewState(draft)
      })
      return
    }

    const state = this.store.snapshot()
    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = "refresh"
      draft.refreshPending = true
      draft.errorMessage = ""
    })

    try {
      const response = await this.gateway.getTranslationOutputReview({
        selectedJobId: state.selectedJobId ?? undefined
      } satisfies GetTranslationOutputReviewRequest)
      const completedJobs = mapCompletedJobs(response)

      if (response.hasSelectedJob === false || state.selectedJobId === null) {
        this.store.update((draft) => {
          draft.phase = "ready"
          draft.pendingAction = null
          draft.refreshPending = false
          draft.hasLoaded = true
          draft.completedJobs = completedJobs
          draft.review = null
          draft.selectedJobId = null
          draft.selectedArtifactId = null
          draft.diffPreview = null
          draft.errorMessage = ""
          draft.actionDisablements = buildActionDisablements(draft, true)
          draft.viewState = resolveViewState(draft)
        })
        return
      }

      const review = mapReviewResponse(response)
      const selectedArtifactId =
        state.selectedArtifactId ??
        (review.artifactId > 0 ? review.artifactId : null)

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.refreshPending = false
        draft.hasLoaded = true
        draft.completedJobs = completedJobs
        draft.review = review
        draft.selectedJobId = review.selectedJobId
        draft.selectedArtifactId = selectedArtifactId
        draft.errorMessage = ""
        draft.actionDisablements = buildActionDisablements(draft, true)
        draft.viewState = resolveViewState(draft)
      })

      await this.refreshDiffPreview()
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.refreshPending = false
        draft.hasLoaded = true
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "Output Review の取得に失敗しました。"
        )
        draft.actionDisablements = buildActionDisablements(draft, true)
        draft.viewState = resolveViewState(draft)
      })
    }
  }

  async generateArtifact(): Promise<void> {
    await this.runCommand("generate", (state) =>
      this.gateway!.generateXTranslatorOutputArtifact({
        jobId: state.selectedJobId!,
        targetGame: state.targetGame,
        outputPath: state.outputPath.trim()
      } satisfies GenerateXTranslatorOutputArtifactRequest)
    )
  }

  async regenerateArtifact(): Promise<void> {
    await this.runCommand("regenerate", (state) =>
      this.gateway!.regenerateXTranslatorOutputArtifact({
        jobId: state.selectedJobId!,
        artifactId: state.selectedArtifactId!,
        targetGame: state.targetGame,
        outputPath: state.outputPath.trim()
      } satisfies RegenerateXTranslatorOutputArtifactRequest)
    )
  }

  private async refreshDiffPreview(): Promise<void> {
    const state = this.store.snapshot()
    if (
      !this.gateway ||
      state.selectedJobId === null ||
      state.selectedArtifactId === null ||
      state.selectedArtifactId <= 0
    ) {
      this.store.update((draft) => {
        draft.diffPreview = null
        draft.actionDisablements = buildActionDisablements(
          draft,
          this.gateway !== null
        )
        draft.viewState = resolveViewState(draft)
      })
      return
    }

    try {
      const response = await this.gateway.getTranslationOutputDiffPreview({
        jobId: state.selectedJobId,
        artifactId: state.selectedArtifactId
      } satisfies GetTranslationOutputDiffPreviewRequest)
      const diffPreview = mapDiffPreviewResponse(response)

      this.store.update((draft) => {
        draft.diffPreview = diffPreview
        draft.actionDisablements = buildActionDisablements(draft, true)
        draft.viewState = resolveViewState(draft)
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.diffPreview = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "Diff Preview の取得に失敗しました。"
        )
        draft.actionDisablements = buildActionDisablements(draft, true)
        draft.viewState = resolveViewState(draft)
      })
    }
  }

  private async runCommand(
    actionKind: TranslationOutputArtifactActionKind,
    execute: (
      state: TranslationOutputArtifactScreenState
    ) => Promise<TranslationOutputArtifactCommandResponse>
  ): Promise<void> {
    const state = this.store.snapshot()
    const disablement = buildActionDisablements(
      state,
      this.gateway !== null
    ).find((item) => item.actionKind === actionKind)

    if (disablement?.disabled) {
      this.store.update((draft) => {
        draft.errorMessage = disablement.reason ?? "操作を実行できません。"
        draft.actionDisablements = buildActionDisablements(
          draft,
          this.gateway !== null
        )
        draft.viewState = resolveViewState(draft)
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "submitting"
      draft.pendingAction = actionKind
      draft.errorMessage = ""
    })

    try {
      const response = await execute(state)
      const lastCommand = mapCommandResponse(response)

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.lastCommand = lastCommand
        draft.selectedArtifactId =
          response.artifactId > 0
            ? response.artifactId
            : draft.selectedArtifactId
        draft.errorMessage =
          lastCommand.errorReason && !lastCommand.isRedacted
            ? lastCommand.errorReason
            : lastCommand.errorKind === "secret_redacted"
              ? mapErrorKindToLabel("secret_redacted")
              : ""
      })

      await this.refresh()
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "出力 command の実行に失敗しました。"
        )
        draft.viewState = resolveViewState(draft)
      })
    }
  }
}
