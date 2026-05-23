import type {
  TranslationJobManagementJobDetail,
  TranslationJobManagementBlockedReason,
  TranslationJobManagementCurrentPhase,
  TranslationJobManagementOperationAvailability,
  TranslationJobManagementReasonCategory
} from "@application/gateway-contract/translation-job-management"

type TranslationJobManagementFilterId =
  | "all"
  | "Ready"
  | "Running"
  | "Paused"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"

interface TranslationJobManagementScreenState {
  phase: "idle" | "loading" | "ready" | "empty" | "error"
  detailPhase: "idle" | "loading" | "ready" | "stale"
  isReloading: boolean
  selectedJobId: number | null
  filterId: TranslationJobManagementFilterId
  searchQuery: string
  activeOperation:
    | TranslationJobManagementOperationAvailability["kind"]
    | "reload"
    | null
  isDeleteConfirmationOpen: boolean
  feedback: {
    tone: "info" | "success" | "warning" | "danger"
    title: string
    message: string
    category?:
      | TranslationJobManagementReasonCategory
      | "stop_requested"
      | "delete_success"
      | "resume_success"
  } | null
  jobs: Array<{
    jobId: number
    jobState: TranslationJobManagementJobDetail["jobState"]
    jobStateLabel: string
    stateTone: TranslationJobManagementJobDetail["stateTone"]
    inputSource: TranslationJobManagementJobDetail["inputSource"]
    progress: TranslationJobManagementJobDetail["progress"]
    canOpenPhase?: boolean
    openBlockedReason?: TranslationJobManagementBlockedReason
    stopAvailability: TranslationJobManagementOperationAvailability
    resumeAvailability: TranslationJobManagementOperationAvailability
    deleteAvailability: TranslationJobManagementOperationAvailability
  }>
  selectedJobDetail: TranslationJobManagementJobDetail | null
}

interface TranslationJobManagementFilterChipViewModel {
  id: TranslationJobManagementFilterId
  label: string
  count: number
  selected: boolean
}

interface TranslationJobManagementOperationViewModel {
  kind: TranslationJobManagementOperationAvailability["kind"]
  label: string
  helperText: string
  enabled: boolean
  reasonText: string
  busy: boolean
}

interface TranslationJobManagementJobCardViewModel {
  jobId: number
  title: string
  jobState: TranslationJobManagementJobDetail["jobState"]
  stateLabel: string
  stateDescription: string
  stateTone: TranslationJobManagementJobDetail["stateTone"]
  inputSourceLabel: string
  sourcePath: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  lastUpdatedLabel: string
  canOpenPhase: boolean
  openBlockedReason: TranslationJobManagementBlockedReason | null
  openBlockedReasonText: string
  jobRunTarget: TranslationJobManagementJobRunTarget | null
  isSelected: boolean
  stopOperation: TranslationJobManagementOperationViewModel
  resumeOperation: TranslationJobManagementOperationViewModel
  deleteOperation: TranslationJobManagementOperationViewModel
}

interface TranslationJobManagementReasonItemViewModel {
  category: TranslationJobManagementReasonCategory
  categoryLabel: string
  title: string
  detail: string
}

interface TranslationJobManagementSelectedJobViewModel {
  jobId: number
  jobIdLabel: string
  jobState: TranslationJobManagementJobDetail["jobState"]
  stateLabel: string
  stateDescription: string
  stateTone: TranslationJobManagementJobDetail["stateTone"]
  inputSourceLabel: string
  inputSourceKindLabel: string
  sourcePath: string
  pluginName: string
  extractedJsonLabel: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  lastUpdatedLabel: string
  cacheStateLabel: string
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialStateLabel: string
  stopOperation: TranslationJobManagementOperationViewModel
  resumeOperation: TranslationJobManagementOperationViewModel
  deleteOperation: TranslationJobManagementOperationViewModel
  resumeBlockedReasons: TranslationJobManagementReasonItemViewModel[]
  warnings: TranslationJobManagementReasonItemViewModel[]
  deleteImpactLines: string[]
}

interface TranslationJobManagementDeleteConfirmationViewModel {
  title: string
  lines: string[]
  confirmLabel: string
  cancelLabel: string
  busy: boolean
}

interface TranslationJobManagementJobRunTarget {
  jobId: number
  stateLabel: string
  stateDescription: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  inputSourceLabel: string
  sourcePath: string
}

interface TranslationJobManagementScreenViewModel {
  gatewayStatus: string
  pageTitle: string
  pageLead: string
  headerCountLabel: string
  listEmptyTitle: string
  listEmptyDescription: string
  listErrorTitle: string
  listErrorDescription: string
  detailPlaceholderTitle: string
  detailPlaceholderDescription: string
  phase: TranslationJobManagementScreenState["phase"]
  detailPhase: TranslationJobManagementScreenState["detailPhase"]
  isReloading: boolean
  searchQuery: string
  filterChips: TranslationJobManagementFilterChipViewModel[]
  jobs: TranslationJobManagementJobCardViewModel[]
  feedback: TranslationJobManagementScreenState["feedback"]
  selectedJob: TranslationJobManagementSelectedJobViewModel | null
  deleteConfirmation: TranslationJobManagementDeleteConfirmationViewModel | null
  jobRunTarget: TranslationJobManagementJobRunTarget | null
}

const FILTER_LABELS: Record<TranslationJobManagementFilterId, string> = {
  all: "すべて",
  Ready: "実行前",
  Running: "実行中",
  Paused: "中断中",
  RecoverableFailed: "再開可能な失敗",
  Failed: "回復不能",
  Canceled: "キャンセル済み"
}

const STATE_DESCRIPTION: Record<string, string> = {
  Ready: "実行前",
  Running: "実行中",
  Paused: "中断中",
  RecoverableFailed: "再開可能な失敗",
  Failed: "回復できない失敗",
  Canceled: "キャンセル済み"
}

function matchesFilter(
  jobState: string,
  filterId: TranslationJobManagementFilterId
): boolean {
  return filterId === "all" || jobState === filterId
}

function matchesSearch(
  job: TranslationJobManagementScreenState["jobs"][number],
  searchQuery: string
): boolean {
  const normalized = searchQuery.trim().toLowerCase()
  if (normalized === "") {
    return true
  }

  return [
    job.inputSource.inputSourceLabel,
    STATE_DESCRIPTION[job.jobState] ?? job.jobStateLabel,
    job.progress.currentPhaseLabel,
    String(job.jobId)
  ].some((value) => value.toLowerCase().includes(normalized))
}

function toOperationViewModel(
  operation: TranslationJobManagementOperationAvailability,
  activeOperation: TranslationJobManagementScreenState["activeOperation"]
): TranslationJobManagementOperationViewModel {
  return {
    kind: operation.kind,
    label: operation.label,
    helperText: operation.helperText,
    enabled: operation.enabled,
    reasonText: operation.reasonText ?? "",
    busy: activeOperation === operation.kind
  }
}

function toReasonItems(
  reasons: TranslationJobManagementJobDetail["warnings"]
): TranslationJobManagementReasonItemViewModel[] {
  return reasons.map((reason) => ({
    category: reason.category,
    categoryLabel: toReasonCategoryLabel(reason.category),
    title: reason.title,
    detail: reason.detail
  }))
}

function toReasonCategoryLabel(
  category: TranslationJobManagementReasonCategory
): string {
  switch (category) {
    case "cache_missing":
      return "入力キャッシュ欠落"
    case "terminal_state":
      return "終端状態"
    case "state_projection_inconsistent":
      return "状態投影不整合"
    case "phase_progress_aggregation_failed":
      return "進捗集約失敗"
    case "stale_selection":
      return "選択情報が古い"
    case "list_load_failure":
      return "一覧取得失敗"
    case "running_delete_blocked":
      return "実行中削除禁止"
    case "stop_requested":
      return "停止要求中"
    case "stop_failed":
      return "停止失敗"
    case "delete_failed":
      return "削除失敗"
    case "resume_failed":
      return "再開失敗"
  }
}

function toJobCard(
  detail: TranslationJobManagementScreenState["jobs"][number],
  selectedJobId: number | null,
  activeOperation: TranslationJobManagementScreenState["activeOperation"]
): TranslationJobManagementJobCardViewModel {
  const stateLabel = STATE_DESCRIPTION[detail.jobState] ?? detail.jobStateLabel
  const canOpenPhase = detail.canOpenPhase !== false
  const openBlockedReason = detail.openBlockedReason ?? null

  return {
    jobId: detail.jobId,
    title: detail.inputSource.inputSourceLabel,
    jobState: detail.jobState,
    stateLabel,
    stateDescription: stateLabel,
    stateTone: detail.stateTone,
    inputSourceLabel: detail.inputSource.inputSourceLabel,
    sourcePath: detail.inputSource.sourcePath,
    currentPhase: detail.progress.currentPhase,
    currentPhaseLabel: detail.progress.currentPhaseLabel,
    progressLabel: detail.progress.progressLabel,
    lastUpdatedLabel: detail.progress.lastUpdatedLabel,
    canOpenPhase,
    openBlockedReason,
    openBlockedReasonText: openBlockedReason?.detail ?? "",
    jobRunTarget: canOpenPhase
      ? {
          jobId: detail.jobId,
          stateLabel,
          stateDescription: stateLabel,
          currentPhase: detail.progress.currentPhase,
          currentPhaseLabel: detail.progress.currentPhaseLabel,
          progressLabel: detail.progress.progressLabel,
          inputSourceLabel: detail.inputSource.inputSourceLabel,
          sourcePath: detail.inputSource.sourcePath
        }
      : null,
    isSelected: selectedJobId === detail.jobId,
    stopOperation: toOperationViewModel(
      detail.stopAvailability,
      activeOperation
    ),
    resumeOperation: toOperationViewModel(
      detail.resumeAvailability,
      activeOperation
    ),
    deleteOperation: toOperationViewModel(
      detail.deleteAvailability,
      activeOperation
    )
  }
}

type JobRunTargetSource = Pick<
  TranslationJobManagementJobDetail,
  | "jobId"
  | "jobState"
  | "jobStateLabel"
  | "inputSource"
  | "progress"
  | "canOpenPhase"
>

function toJobRunTarget(
  detail: JobRunTargetSource | null
): TranslationJobManagementJobRunTarget | null {
  if (!detail) {
    return null
  }

  if (detail.canOpenPhase === false) {
    return null
  }

  const stateLabel = STATE_DESCRIPTION[detail.jobState] ?? detail.jobStateLabel

  return {
    jobId: detail.jobId,
    stateLabel,
    stateDescription: stateLabel,
    currentPhase: detail.progress.currentPhase,
    currentPhaseLabel: detail.progress.currentPhaseLabel,
    progressLabel: detail.progress.progressLabel,
    inputSourceLabel: detail.inputSource.inputSourceLabel,
    sourcePath: detail.inputSource.sourcePath
  }
}

function resolveJobRunTarget(
  state: TranslationJobManagementScreenState
): TranslationJobManagementJobRunTarget | null {
  if (state.selectedJobDetail) {
    return toJobRunTarget(state.selectedJobDetail)
  }

  if (state.detailPhase !== "loading" || state.selectedJobId === null) {
    return null
  }

  const selectedJobSummary = state.jobs.find(
    (job) => job.jobId === state.selectedJobId
  )

  return toJobRunTarget(selectedJobSummary ?? null)
}

export class TranslationJobManagementPresenter {
  toViewModel(
    state: TranslationJobManagementScreenState,
    isGatewayConnected: boolean
  ): TranslationJobManagementScreenViewModel {
    const filterChips = buildFilterChips(state)
    const jobCards = state.jobs
      .filter(
        (job) =>
          matchesFilter(job.jobState, state.filterId) &&
          matchesSearch(job, state.searchQuery)
      )
      .map((job) => toJobCard(job, state.selectedJobId, state.activeOperation))

    const selectedJob = state.selectedJobDetail
      ? {
          jobId: state.selectedJobDetail.jobId,
          jobIdLabel: `ジョブ #${state.selectedJobDetail.jobId}`,
          jobState: state.selectedJobDetail.jobState,
          stateLabel:
            STATE_DESCRIPTION[state.selectedJobDetail.jobState] ??
            state.selectedJobDetail.jobStateLabel,
          stateDescription:
            STATE_DESCRIPTION[state.selectedJobDetail.jobState] ??
            state.selectedJobDetail.jobStateLabel,
          stateTone: state.selectedJobDetail.stateTone,
          inputSourceLabel:
            state.selectedJobDetail.inputSource.inputSourceLabel,
          inputSourceKindLabel:
            state.selectedJobDetail.inputSource.inputSourceKindLabel,
          sourcePath: state.selectedJobDetail.inputSource.sourcePath,
          pluginName: state.selectedJobDetail.inputSource.pluginName,
          extractedJsonLabel:
            state.selectedJobDetail.inputSource.extractedJsonLabel,
          currentPhase: state.selectedJobDetail.progress.currentPhase,
          currentPhaseLabel: state.selectedJobDetail.progress.currentPhaseLabel,
          progressLabel: state.selectedJobDetail.progress.progressLabel,
          lastUpdatedLabel: state.selectedJobDetail.progress.lastUpdatedLabel,
          cacheStateLabel: state.selectedJobDetail.cacheStateLabel,
          providerLabel: state.selectedJobDetail.runtimeSummary.providerLabel,
          modelLabel: state.selectedJobDetail.runtimeSummary.modelLabel,
          executionModeLabel:
            state.selectedJobDetail.runtimeSummary.executionModeLabel,
          credentialStateLabel:
            state.selectedJobDetail.runtimeSummary.credentialStateLabel,
          stopOperation: toOperationViewModel(
            state.selectedJobDetail.stopAvailability,
            state.activeOperation
          ),
          resumeOperation: toOperationViewModel(
            state.selectedJobDetail.resumeAvailability,
            state.activeOperation
          ),
          deleteOperation: toOperationViewModel(
            state.selectedJobDetail.deleteAvailability,
            state.activeOperation
          ),
          resumeBlockedReasons: toReasonItems(
            state.selectedJobDetail.resumeBlockedReasons
          ),
          warnings: toReasonItems(state.selectedJobDetail.warnings),
          deleteImpactLines: [...state.selectedJobDetail.deleteImpactLines]
        }
      : null

    return {
      gatewayStatus: isGatewayConnected ? "接続済み" : "未接続",
      pageTitle: "未完了ジョブ一覧",
      pageLead:
        "新規翻訳を開始するか、未完了ジョブを選んで現在の翻訳段階から再開します。",
      headerCountLabel: `${jobCards.length} 件を表示`,
      listEmptyTitle: "管理対象がありません",
      listEmptyDescription:
        "未完了ジョブはありません。新規翻訳を開始して入力データの登録から進めてください。",
      listErrorTitle: "一覧を読み込めません",
      listErrorDescription:
        "未完了ジョブの一覧取得に失敗しました。再読込しても直らない場合は、統合後の gateway を確認してください。",
      detailPlaceholderTitle: "ジョブを選択してください",
      detailPlaceholderDescription:
        "一覧から 1 件選ぶと、入力出自、進捗、再開不可理由、削除影響を確認できます。",
      phase: state.phase,
      detailPhase: state.detailPhase,
      isReloading: state.isReloading,
      searchQuery: state.searchQuery,
      filterChips,
      jobs: jobCards,
      feedback: state.feedback ? { ...state.feedback } : null,
      selectedJob,
      deleteConfirmation:
        state.isDeleteConfirmationOpen && selectedJob
          ? {
              title: `${selectedJob.jobIdLabel} を削除しますか`,
              lines: [...selectedJob.deleteImpactLines],
              confirmLabel: "ジョブ情報だけを削除する",
              cancelLabel: "戻る",
              busy: state.activeOperation === "delete"
            }
          : null,
      jobRunTarget: resolveJobRunTarget(state)
    }
  }
}

function buildFilterChips(
  state: TranslationJobManagementScreenState
): TranslationJobManagementFilterChipViewModel[] {
  return (Object.keys(FILTER_LABELS) as TranslationJobManagementFilterId[]).map(
    (filterId) => ({
      id: filterId,
      label: FILTER_LABELS[filterId],
      count:
        filterId === "all"
          ? state.jobs.length
          : state.jobs.filter((job) => job.jobState === filterId).length,
      selected: state.filterId === filterId
    })
  )
}
