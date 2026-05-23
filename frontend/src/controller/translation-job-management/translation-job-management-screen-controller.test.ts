import { describe, expect, test, vi } from "vitest"

import type {
  TranslationJobManagementScreenState,
  TranslationJobManagementScreenViewModel
} from "@application/contract/translation-job-management/translation-job-management-screen-types"

import { TranslationJobManagementScreenController } from "./translation-job-management-screen-controller"

function createState(): TranslationJobManagementScreenState {
  return {
    phase: "idle",
    jobs: [],
    selectedJobId: null,
    selectedJobDetail: null,
    detailPhase: "idle",
    filterId: "all",
    searchQuery: "",
    isReloading: false,
    activeOperation: null,
    isDeleteConfirmationOpen: false,
    feedback: null
  }
}

const viewModel: TranslationJobManagementScreenViewModel = {
  gatewayStatus: "接続済み",
  pageTitle: "ジョブ管理",
  pageLead: "lead",
  headerCountLabel: "0 件を表示",
  listEmptyTitle: "empty",
  listEmptyDescription: "empty desc",
  listErrorTitle: "error",
  listErrorDescription: "error desc",
  detailPlaceholderTitle: "select",
  detailPlaceholderDescription: "select desc",
  phase: "idle",
  detailPhase: "idle",
  isReloading: false,
  searchQuery: "",
  filterChips: [],
  jobs: [],
  feedback: null,
  selectedJob: null,
  deleteConfirmation: null,
  jobRunTarget: null
}

describe("TranslationJobManagementScreenController", () => {
  test("mount と action を usecase へ委譲する", async () => {
    const load = vi.fn().mockResolvedValue(undefined)
    const reload = vi.fn().mockResolvedValue(undefined)
    const selectJob = vi.fn().mockResolvedValue(undefined)
    const requestStop = vi.fn().mockResolvedValue(undefined)
    const requestResume = vi.fn().mockResolvedValue(undefined)
    const deleteSelectedJob = vi.fn().mockResolvedValue(undefined)

    const controller = new TranslationJobManagementScreenController({
      isGatewayConnected: true,
      store: {
        snapshot: () => createState(),
        subscribe: () => () => {}
      },
      presenter: { toViewModel: () => viewModel },
      useCase: {
        load,
        reload,
        selectJob,
        setFilter: vi.fn(),
        setSearchQuery: vi.fn(),
        requestStop,
        requestResume,
        openDeleteConfirmation: vi.fn(),
        closeDeleteConfirmation: vi.fn(),
        deleteSelectedJob
      }
    })

    await controller.mount()
    await controller.reload()
    await controller.selectJob(9)
    await controller.requestStop()
    await controller.requestResume()
    await controller.deleteSelectedJob()

    expect(load).toHaveBeenCalledTimes(1)
    expect(reload).toHaveBeenCalledTimes(1)
    expect(selectJob).toHaveBeenCalledWith(9)
    expect(requestStop).toHaveBeenCalledTimes(1)
    expect(requestResume).toHaveBeenCalledTimes(1)
    expect(deleteSelectedJob).toHaveBeenCalledTimes(1)
  })

  test("subscribe は presenter を通した view model を通知する", () => {
    let listener:
      | ((state: TranslationJobManagementScreenState) => void)
      | undefined
    const controller = new TranslationJobManagementScreenController({
      isGatewayConnected: false,
      store: {
        snapshot: () => createState(),
        subscribe: (cb) => {
          listener = cb
          return () => {}
        }
      },
      presenter: {
        toViewModel: (_state, isGatewayConnected) => ({
          ...viewModel,
          gatewayStatus: isGatewayConnected ? "接続済み" : "未接続"
        })
      },
      useCase: {
        load: vi.fn().mockResolvedValue(undefined),
        reload: vi.fn().mockResolvedValue(undefined),
        selectJob: vi.fn().mockResolvedValue(undefined),
        setFilter: vi.fn(),
        setSearchQuery: vi.fn(),
        requestStop: vi.fn().mockResolvedValue(undefined),
        requestResume: vi.fn().mockResolvedValue(undefined),
        openDeleteConfirmation: vi.fn(),
        closeDeleteConfirmation: vi.fn(),
        deleteSelectedJob: vi.fn().mockResolvedValue(undefined)
      }
    })

    const sink = vi.fn()
    controller.subscribe(sink)
    listener?.(createState())

    expect(sink).toHaveBeenCalled()
    const lastCall = sink.mock.calls.at(-1) as
      | [TranslationJobManagementScreenViewModel]
      | undefined
    expect(lastCall?.[0].gatewayStatus).toBe("未接続")
  })
})
