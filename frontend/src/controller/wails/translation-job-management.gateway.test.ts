import { afterEach, describe, expect, test, vi } from "vitest"

import { createTranslationJobManagementGateway } from "./translation-job-management.gateway"

type GoRecord = {
  wails?: {
    TranslationJobManagementController?: Record<
      string,
      ReturnType<typeof vi.fn>
    >
    AppController?: Record<string, ReturnType<typeof vi.fn>>
  }
}

const originalGo: unknown = Reflect.get(globalThis as object, "go")

function installGo(record: GoRecord): void {
  Object.defineProperty(globalThis, "go", {
    value: record,
    configurable: true,
    writable: true
  })
}

afterEach(() => {
  vi.restoreAllMocks()
  Object.defineProperty(globalThis, "go", {
    value: originalGo,
    configurable: true,
    writable: true
  })
})

describe("createTranslationJobManagementGateway", () => {
  test("ListIncompleteJobs は request なしで controller binding を呼ぶ", async () => {
    const listIncompleteJobs = vi.fn(() => Promise.resolve({ jobs: [] }))
    installGo({
      wails: {
        TranslationJobManagementController: {
          ListIncompleteJobs: listIncompleteJobs
        }
      }
    })

    const gateway = createTranslationJobManagementGateway()
    await expect(gateway.ListIncompleteJobs()).resolves.toEqual({ jobs: [] })
    expect(listIncompleteJobs).toHaveBeenCalledWith()
  })

  test("controller 未接続時は AppController の binding を使う", async () => {
    const getJobDetail = vi.fn(() => Promise.resolve({ jobId: 7 }))
    installGo({
      wails: {
        AppController: {
          GetJobDetail: getJobDetail
        }
      }
    })

    const gateway = createTranslationJobManagementGateway()
    await gateway.GetJobDetail({ jobId: 7 })

    expect(getJobDetail).toHaveBeenCalledWith({ jobId: 7 })
  })

  test("action 系は request をそのまま binding へ渡す", async () => {
    const requestStop = vi.fn(() =>
      Promise.resolve({ message: "停止要求中", tone: "info" })
    )
    const resumeJob = vi.fn(() =>
      Promise.resolve({ message: "再開要求中", tone: "info" })
    )
    const deleteJob = vi.fn(() =>
      Promise.resolve({ message: "削除要求中", tone: "warning" })
    )
    installGo({
      wails: {
        TranslationJobManagementController: {
          RequestStop: requestStop,
          ResumeJob: resumeJob,
          DeleteJob: deleteJob
        }
      }
    })

    const gateway = createTranslationJobManagementGateway()
    await gateway.RequestStop({ jobId: 1 })
    await gateway.ResumeJob({ jobId: 2 })
    await gateway.DeleteJob({ jobId: 3 })

    expect(requestStop).toHaveBeenCalledWith({ jobId: 1 })
    expect(resumeJob).toHaveBeenCalledWith({ jobId: 2 })
    expect(deleteJob).toHaveBeenCalledWith({ jobId: 3 })
  })

  test("binding 未接続時は統合前エラーを返す", async () => {
    installGo({ wails: {} })
    const gateway = createTranslationJobManagementGateway()

    await expect(gateway.DeleteJob({ jobId: 3 })).rejects.toThrow(
      "Wails binding is not wired yet: DeleteJob"
    )
  })
})
