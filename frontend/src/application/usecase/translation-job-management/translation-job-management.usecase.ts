import type {
  TranslationJobManagementActionResponse,
  TranslationJobManagementGatewayContract,
  TranslationJobManagementJobDetail,
  TranslationJobManagementJobSummary,
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
  jobs: TranslationJobManagementJobSummary[]
  selectedJobId: number | null
  selectedJobDetail: TranslationJobManagementJobDetail | null
  detailPhase: "idle" | "loading" | "ready" | "stale"
  filterId: TranslationJobManagementFilterId
  searchQuery: string
  isReloading: boolean
  activeOperation: TranslationJobManagementOperationAvailability["kind"] | "reload" | null
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
}

interface TranslationJobManagementStoreLike {
  snapshot(): TranslationJobManagementScreenState
  update(mutator: (draft: TranslationJobManagementScreenState) => void): void
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

function createMissingGatewayMessage(): string {
  return "gateway は統合前です。fakeApi で表示確認を行ってください。"
}

function buildFeedback(
  title: string,
  message: string,
  tone: "info" | "success" | "warning" | "danger",
  category?:
    | TranslationJobManagementReasonCategory
    | "stop_requested"
    | "delete_success"
    | "resume_success"
) {
  return {
    tone,
    title,
    message,
    category
  } as const
}

export class TranslationJobManagementUseCase {
  constructor(
    private readonly gateway: TranslationJobManagementGatewayContract | null,
    private readonly store: TranslationJobManagementStoreLike
  ) {}

  async load(): Promise<void> {
    await this.reloadInternal(false)
  }

  async reload(): Promise<void> {
    await this.reloadInternal(true)
  }

  setFilter(filterId: TranslationJobManagementFilterId): void {
    this.store.update((draft) => {
      draft.filterId = filterId
    })
  }

  setSearchQuery(searchQuery: string): void {
    this.store.update((draft) => {
      draft.searchQuery = searchQuery
    })
  }

  async selectJob(jobId: number): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.selectedJobId = jobId
        draft.feedback = buildFeedback(
          "未接続",
          createMissingGatewayMessage(),
          "warning"
        )
      })
      return
    }

    this.store.update((draft) => {
      draft.selectedJobId = jobId
      draft.detailPhase = "loading"
      draft.isDeleteConfirmationOpen = false
      if (
        draft.feedback?.category === "stale_selection" ||
        draft.feedback?.category === "delete_success"
      ) {
        draft.feedback = null
      }
    })

    try {
      const detail = await this.gateway.GetJobDetail({ jobId })
      this.store.update((draft) => {
        draft.selectedJobId = jobId
        draft.selectedJobDetail = detail
        draft.detailPhase = "ready"
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.selectedJobDetail = null
        draft.detailPhase = "stale"
        draft.feedback = buildFeedback(
          "参照できません",
          sanitizeErrorMessage(
            error,
            "選択した job を再度読み込めませんでした。一覧を更新してください。"
          ),
          "warning",
          "stale_selection"
        )
      })
    }
  }

  async requestStop(): Promise<void> {
    const state = this.store.snapshot()
    const selectedJobId = state.selectedJobDetail?.jobId
    const gateway = this.gateway
    if (!selectedJobId || !gateway) {
      return
    }

    await this.runAction("stop", async () =>
      gateway.RequestStop({ jobId: selectedJobId })
    )
  }

  async requestResume(): Promise<void> {
    const state = this.store.snapshot()
    const selectedJobId = state.selectedJobDetail?.jobId
    const gateway = this.gateway
    if (!selectedJobId || !gateway) {
      return
    }

    await this.runAction("resume", async () =>
      gateway.ResumeJob({ jobId: selectedJobId })
    )
  }

  openDeleteConfirmation(): void {
    this.store.update((draft) => {
      draft.isDeleteConfirmationOpen = true
    })
  }

  closeDeleteConfirmation(): void {
    this.store.update((draft) => {
      draft.isDeleteConfirmationOpen = false
    })
  }

  async deleteSelectedJob(): Promise<void> {
    const state = this.store.snapshot()
    if (!state.selectedJobDetail || !this.gateway) {
      return
    }

    this.store.update((draft) => {
      draft.activeOperation = "delete"
      draft.feedback = null
    })

    try {
      const response = await this.gateway.DeleteJob({
        jobId: state.selectedJobDetail.jobId
      })
      this.store.update((draft) => {
        draft.activeOperation = null
        draft.isDeleteConfirmationOpen = false
        if (response.deletedJobId) {
          draft.jobs = draft.jobs.filter((job) => job.jobId !== response.deletedJobId)
          draft.selectedJobId = null
          draft.selectedJobDetail = null
          draft.detailPhase = "idle"
          draft.feedback = buildFeedback(
            "削除しました",
            response.message,
            response.tone,
            "delete_success"
          )
          draft.phase = draft.jobs.length === 0 ? "empty" : "ready"
          return
        }

        if (response.detail) {
          syncDetailIntoList(draft.jobs, response.detail)
          draft.selectedJobDetail = response.detail
          draft.detailPhase = "ready"
        }
        draft.feedback = buildFeedback(
          "削除できません",
          response.message,
          response.tone,
          response.reasonCategory ?? "delete_failed"
        )
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.activeOperation = null
        draft.feedback = buildFeedback(
          "削除できません",
          sanitizeErrorMessage(
            error,
            "削除に失敗しました。時間をおいて再実行してください。"
          ),
          "danger",
          "delete_failed"
        )
      })
    }
  }

  private async reloadInternal(isManualReload: boolean): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.phase = "error"
        draft.detailPhase = draft.selectedJobId ? "stale" : "idle"
        draft.feedback = buildFeedback(
          "未接続",
          createMissingGatewayMessage(),
          "warning",
          "list_load_failure"
        )
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = draft.jobs.length === 0 ? "loading" : draft.phase
      draft.isReloading = isManualReload
      draft.activeOperation = isManualReload ? "reload" : draft.activeOperation
      if (!isManualReload) {
        draft.feedback = null
      }
    })

    try {
      const response = await this.gateway.ListIncompleteJobs()
      this.store.update((draft) => {
        draft.jobs = response.jobs
        draft.phase = response.jobs.length === 0 ? "empty" : "ready"
        draft.isReloading = false
        draft.activeOperation = null
      })

      const current = this.store.snapshot()
      if (current.selectedJobId !== null) {
        const stillVisible = response.jobs.some(
          (job) => job.jobId === current.selectedJobId
        )
        if (!stillVisible) {
          this.store.update((draft) => {
            draft.selectedJobDetail = null
            draft.detailPhase = "stale"
            draft.feedback = buildFeedback(
              "選択が古くなりました",
              "一覧を更新した結果、選択中の job が見つかりませんでした。別の job を選択してください。",
              "warning",
              "stale_selection"
            )
          })
          return
        }

        await this.selectJob(current.selectedJobId)
      }
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "error"
        draft.isReloading = false
        draft.activeOperation = null
        draft.feedback = buildFeedback(
          "一覧を読み込めません",
          sanitizeErrorMessage(
            error,
            "未完了ジョブの一覧取得に失敗しました。再読込してください。"
          ),
          "danger",
          "list_load_failure"
        )
      })
    }
  }

  private async runAction(
    operation: "stop" | "resume",
    runner: () => Promise<TranslationJobManagementActionResponse>
  ): Promise<void> {
    this.store.update((draft) => {
      draft.activeOperation = operation
      draft.feedback = null
    })

    try {
      const response = await runner()
      this.store.update((draft) => {
        draft.activeOperation = null
        if (response.detail) {
          syncDetailIntoList(draft.jobs, response.detail)
          draft.selectedJobDetail = response.detail
          draft.detailPhase = "ready"
        }
        draft.feedback = buildFeedback(
          operation === "stop" ? "停止要求を更新しました" : "再開結果を更新しました",
          response.message,
          response.tone,
          response.reasonCategory ??
            (operation === "stop" ? "stop_requested" : "resume_success")
        )
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.activeOperation = null
        draft.feedback = buildFeedback(
          operation === "stop" ? "停止できません" : "再開できません",
          sanitizeErrorMessage(
            error,
            operation === "stop"
              ? "停止要求に失敗しました。時間をおいて再実行してください。"
              : "再開に失敗しました。時間をおいて再実行してください。"
          ),
          "danger",
          operation === "stop" ? "stop_failed" : "resume_failed"
        )
      })
    }
  }
}

function syncDetailIntoList(
  jobs: TranslationJobManagementScreenState["jobs"],
  detail: TranslationJobManagementJobDetail
): void {
  const index = jobs.findIndex((job) => job.jobId === detail.jobId)
  if (index < 0) {
    return
  }

  jobs[index] = {
    jobId: detail.jobId,
    jobState: detail.jobState,
    jobStateLabel: detail.jobStateLabel,
    stateTone: detail.stateTone,
    canOpenPhase: detail.canOpenPhase,
    openBlockedReason: detail.openBlockedReason
      ? { ...detail.openBlockedReason }
      : undefined,
    inputSource: { ...detail.inputSource },
    progress: { ...detail.progress },
    stopAvailability: { ...detail.stopAvailability },
    resumeAvailability: { ...detail.resumeAvailability },
    deleteAvailability: { ...detail.deleteAvailability }
  }
}
