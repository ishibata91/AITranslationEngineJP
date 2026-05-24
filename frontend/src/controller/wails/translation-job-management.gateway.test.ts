import { afterEach, describe, expect, test, vi } from "vitest"

import {
  DeleteJob,
  GetJobDetail,
  ListIncompleteJobs,
  RequestStop
} from "../../../wailsjs/go/wails/AppController.js"
import { createTranslationJobManagementGateway } from "./translation-job-management.gateway"

type ListIncompleteJobsResponse = Awaited<ReturnType<typeof ListIncompleteJobs>>
type GetJobDetailResponse = Awaited<ReturnType<typeof GetJobDetail>>
type RequestStopResponse = Awaited<ReturnType<typeof RequestStop>>

vi.mock("../../../wailsjs/go/wails/AppController.js", () => ({
  ListIncompleteJobs: vi.fn(),
  GetJobDetail: vi.fn(),
  RequestStop: vi.fn(),
  ResumeJob: vi.fn(),
  DeleteJob: vi.fn()
}))

function createJobSummary(jobId: number) {
  return {
    jobId,
    jobState: "Ready" as const,
    jobStateLabel: "Ready",
    stateTone: "neutral" as const,
    canOpenPhase: true,
    openBlockedReason: {
      category: "cache_missing" as const,
      title: "cache missing",
      detail: "cache is not prepared"
    },
    inputSource: {
      inputSourceId: 1,
      inputSourceLabel: "Input Source",
      inputSourceKindLabel: "Plugin",
      sourcePath: "mods/source.json",
      pluginName: "plugin.esp",
      extractedJsonLabel: "extracted.json"
    },
    progress: {
      currentPhase: "term_translation" as const,
      currentPhaseLabel: "Term Translation",
      percent: 0,
      progressLabel: "0%",
      lastUpdatedLabel: "just now"
    },
    stopAvailability: {
      kind: "stop" as const,
      enabled: true,
      label: "Stop",
      helperText: "Stop current job",
      reasonCategory: "stop_requested" as const,
      reasonText: "stop is available"
    },
    resumeAvailability: {
      kind: "resume" as const,
      enabled: false,
      label: "Resume",
      helperText: "Resume paused job",
      reasonCategory: "resume_failed" as const,
      reasonText: "job is not paused"
    },
    deleteAvailability: {
      kind: "delete" as const,
      enabled: true,
      label: "Delete",
      helperText: "Delete job",
      reasonCategory: "delete_failed" as const,
      reasonText: "none"
    }
  }
}

function createJobDetail(): GetJobDetailResponse {
  return {
    ...createJobSummary(7),
    cacheState: "available",
    cacheStateLabel: "Available",
    runtimeSummary: {
      providerLabel: "OpenAI",
      modelLabel: "gpt-4.1-mini",
      executionModeLabel: "batch",
      credentialState: "configured",
      credentialStateLabel: "Configured"
    },
    resumeBlockedReasons: [
      {
        category: "resume_failed",
        title: "resume blocked",
        detail: "job is not paused"
      }
    ],
    warnings: [
      {
        category: "list_load_failure",
        title: "warning",
        detail: "transient"
      }
    ],
    deleteImpactLines: ["output cache will be removed"]
  } as GetJobDetailResponse
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe("createTranslationJobManagementGateway", () => {
  test("request passthrough: action request を公開 binding wrapper にそのまま渡す", async () => {
    const response = {
      message: "停止要求中",
      tone: "info"
    } as unknown as RequestStopResponse
    vi.mocked(RequestStop).mockResolvedValue(response)

    const gateway = createTranslationJobManagementGateway()

    await gateway.RequestStop({ jobId: 3 })

    expect(RequestStop).toHaveBeenCalledTimes(1)
    expect(RequestStop).toHaveBeenCalledWith({ jobId: 3 })
  })

  test("valid response: full DTO shape を満たす job detail を返す", async () => {
    vi.mocked(GetJobDetail).mockResolvedValue(createJobDetail())

    const gateway = createTranslationJobManagementGateway()

    await expect(gateway.GetJobDetail({ jobId: 7 })).resolves.toMatchObject({
      jobId: 7,
      runtimeSummary: { providerLabel: "OpenAI" }
    })
  })

  test("binding failure: 未接続時の wrapper 例外を返す", async () => {
    vi.mocked(DeleteJob).mockRejectedValue(
      new Error("Wails binding is not wired yet: DeleteJob")
    )

    const gateway = createTranslationJobManagementGateway()

    await expect(gateway.DeleteJob({ jobId: 9 })).rejects.toThrow(
      "Wails binding is not wired yet: DeleteJob"
    )
  })

  test("invalid response shape: GatewayResponseShapeError と secret 非露出診断を返す", async () => {
    vi.mocked(ListIncompleteJobs).mockResolvedValue({
      jobs: [
        {
          ...createJobSummary(1),
          inputSource: {
            ...createJobSummary(1).inputSource,
            inputSourceId: "broken-id",
            credentialInput: "raw-secret-value"
          }
        }
      ]
    } as unknown as ListIncompleteJobsResponse)

    const gateway = createTranslationJobManagementGateway()

    await expect(gateway.ListIncompleteJobs()).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })

    try {
      await gateway.ListIncompleteJobs()
    } catch (error: unknown) {
      expect(error).toBeInstanceOf(Error)
      expect(typeof (error as Error & { internalDiagnostic?: unknown }).internalDiagnostic).toBe("string")
      const diagnostic = JSON.stringify(error)
      expect(diagnostic).not.toContain("raw-secret-value")
      expect(diagnostic).not.toContain("credentialInput")
      expect(diagnostic).not.toContain("apiKey")
    }
  })

})
