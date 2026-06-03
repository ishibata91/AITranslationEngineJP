import { afterEach, describe, expect, test, vi } from "vitest"

import {
  GetProcessingTargetList,
  GetBodyTranslationPhaseSummary,
  PauseBodyTranslationPhase,
  StartBodyTranslationPhase
} from "../../../wailsjs/go/wails/AppController.js"
import { createBodyTranslationPhaseGateway } from "./body-translation-phase.gateway"

type GetBodyTranslationPhaseSummaryResponse = Awaited<
  ReturnType<typeof GetBodyTranslationPhaseSummary>
>
type PauseBodyTranslationPhaseResponse = Awaited<
  ReturnType<typeof PauseBodyTranslationPhase>
>

vi.mock("../../../wailsjs/go/wails/AppController.js", () => ({
  CancelBodyTranslationPhase: vi.fn(),
  GetBodyTranslationOutputReadiness: vi.fn(),
  GetBodyTranslationPhaseSummary: vi.fn(),
  GetProcessingTargetList: vi.fn(),
  PauseBodyTranslationPhase: vi.fn(),
  ResumeBodyTranslationPhase: vi.fn(),
  RetryBodyTranslationPhase: vi.fn(),
  SaveBodyTranslationPhaseAISettings: vi.fn(),
  StartBodyTranslationPhase: vi.fn()
}))

function createSummaryResponse(): GetBodyTranslationPhaseSummaryResponse {
  return {
    jobId: 5,
    currentPhase: "body_translation",
    phaseState: "ready",
    phaseRunId: 10,
    startedAt: "2026-05-25T00:00:00Z",
    finishedAt: undefined,
    progress: {
      percent: 0,
      processedCount: 0,
      totalCount: 1,
      targetCount: 1,
      translatedCount: 0,
      skippedCount: 0,
      currentStep: "ready"
    },
    inputSummary: {
      targetCount: 1,
      skippedReasons: [],
      inputSnapshotRef: "snapshot-ref",
      dictionaryDigest: "sha256:dictionary",
      personaDigest: "sha256:persona",
      metadataDigest: "sha256:metadata",
      promptDigest: "sha256:prompt"
    },
    requestSummary: {
      providerTargetCount: 1,
      exactDictionaryExclusionCount: 0,
      partialDictionaryConstraintCount: 0
    },
    execution: {
      credentialRef: "cred-ref",
      provider: "openai-compatible",
      model: "gpt-4.1-mini",
      executionMode: "batch",
      requestUnitCount: 1,
      outputCount: 0
    },
    fieldResults: [],
    resultSummary: {
      translatedCount: 0,
      failedCount: 0,
      skippedCount: 0,
      protectionFailedCount: 0,
      outputReadyCount: 0,
      outputCount: 0,
      fieldResults: []
    },
    outputReadiness: {
      completedFieldCount: 0,
      statusConsistent: true,
      outputCount: 0
    },
    errorSummary: {
      errorKind: "provider_failure",
      reason: "not started",
      retryable: true,
      isRedacted: true
    },
    projection: {
      phaseLifecycle: "ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      targetCount: 10,
      previousPhaseLifecycle: "completed",
      personaBodyReadiness: {
        bodyReadiness: true,
        snapshotReferenceStatus: "available"
      }
    }
  } as unknown as GetBodyTranslationPhaseSummaryResponse
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe("createBodyTranslationPhaseGateway", () => {
  test("processing target list request と totalCount response をそのまま受け渡す", async () => {
    vi.mocked(GetProcessingTargetList).mockResolvedValue({
      items: [
        {
          id: "field:1",
          name: "FULL",
          detail: "原文: Hello",
          titleParts: [{ text: "FULL" }],
          metadata: [{ label: "訳語", value: "こんにちは" }],
          convertValues: vi.fn()
        }
      ],
      metadata: [],
      page: 4,
      pageSize: 40,
      totalCount: 256,
      searchQuery: "Hello",
      convertValues: vi.fn()
    })

    const gateway = createBodyTranslationPhaseGateway()
    await expect(
      gateway.getProcessingTargetList?.({
        jobId: 5,
        phase: "body_translation",
        page: 4,
        pageSize: 40,
        searchQuery: "Hello"
      })
    ).resolves.toMatchObject({
      items: [{ id: "field:1" }],
      page: 4,
      pageSize: 40,
      totalCount: 256,
      searchQuery: "Hello"
    })

    expect(GetProcessingTargetList).toHaveBeenCalledWith({
      jobId: 5,
      phase: "body_translation",
      page: 4,
      pageSize: 40,
      searchQuery: "Hello"
    })
  })

  test("request passthrough: command request を公開 binding wrapper にそのまま渡す", async () => {
    const response = {
      ...createSummaryResponse(),
      inputSnapshotDigest: "sha256:snapshot",
      retryable: true
    } as unknown as PauseBodyTranslationPhaseResponse
    vi.mocked(PauseBodyTranslationPhase).mockResolvedValue(response)

    const gateway = createBodyTranslationPhaseGateway()

    await gateway.pauseBodyTranslationPhase({ jobId: 5, phaseRunId: 10 })

    expect(PauseBodyTranslationPhase).toHaveBeenCalledTimes(1)
    expect(PauseBodyTranslationPhase).toHaveBeenCalledWith({
      jobId: 5,
      phaseRunId: 10
    })
  })

  test("valid response: full DTO shape を満たす summary を返す", async () => {
    vi.mocked(GetBodyTranslationPhaseSummary).mockResolvedValue(
      createSummaryResponse()
    )

    const gateway = createBodyTranslationPhaseGateway()

    await expect(
      gateway.getBodyTranslationPhaseSummary({ jobId: 5 })
    ).resolves.toMatchObject({
      jobId: 5,
      execution: { provider: "openai-compatible" }
    })
  })

  test("binding failure: 未接続時の wrapper 例外を返す", async () => {
    vi.mocked(StartBodyTranslationPhase).mockRejectedValue(
      new Error("Wails binding is not wired yet: StartBodyTranslationPhase")
    )

    const gateway = createBodyTranslationPhaseGateway()

    await expect(
      gateway.startBodyTranslationPhase({ jobId: 5 })
    ).rejects.toThrow("Wails binding is not wired yet: StartBodyTranslationPhase")
  })

  // RAEF-UNIT-023: assertBodyProjectionShape — personaBodyReadiness object が存在しない場合は検証失敗する
  // 根拠: design-diff.md G-4-c

  // RAEF-UNIT-023: personaBodyReadiness object が存在しない場合は GatewayResponseShapeError を返す
  test("projection の personaBodyReadiness が存在しない場合は GatewayResponseShapeError を返す（G-4-c personaBodyReadiness 必須 object）", async () => {
    // G-4-c: projection.personaBodyReadiness が存在しない → 必須 object 不足で検証失敗
    const responseWithoutPersonaBodyReadiness = {
      ...createSummaryResponse(),
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        targetCount: 10,
        previousPhaseLifecycle: "completed"
        // personaBodyReadiness を意図的に除外する
      }
    } as unknown as Awaited<ReturnType<typeof GetBodyTranslationPhaseSummary>>

    vi.mocked(GetBodyTranslationPhaseSummary).mockResolvedValue(responseWithoutPersonaBodyReadiness)

    const gateway = createBodyTranslationPhaseGateway()

    await expect(gateway.getBodyTranslationPhaseSummary({ jobId: 5 })).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })
  })

  test("invalid response shape: GatewayResponseShapeError と secret 非露出診断を返す", async () => {
    vi.mocked(GetBodyTranslationPhaseSummary).mockResolvedValue({
      ...createSummaryResponse(),
      execution: {
        ...createSummaryResponse().execution,
        requestUnitCount: "invalid-count",
        credentialInput: "raw-secret-value"
      }
    } as unknown as GetBodyTranslationPhaseSummaryResponse)

    const gateway = createBodyTranslationPhaseGateway()

    await expect(
      gateway.getBodyTranslationPhaseSummary({ jobId: 5 })
    ).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })

    try {
      await gateway.getBodyTranslationPhaseSummary({ jobId: 5 })
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
