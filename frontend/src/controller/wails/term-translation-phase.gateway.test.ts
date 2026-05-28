import { afterEach, describe, expect, test, vi } from "vitest"

import {
  GetProcessingTargetList,
  GetTermTranslationNextPhaseReadiness,
  GetTermTranslationPhaseSummary,
  SaveTermTranslationPhaseAISettings,
  StartTermTranslationPhase
} from "../../../wailsjs/go/wails/AppController.js"
import { createTermTranslationPhaseGateway } from "./term-translation-phase.gateway"

type GetTermTranslationPhaseSummaryResponse = Awaited<
  ReturnType<typeof GetTermTranslationPhaseSummary>
>
type GetTermTranslationNextPhaseReadinessResponse = Awaited<
  ReturnType<typeof GetTermTranslationNextPhaseReadiness>
>
type SaveTermTranslationPhaseAISettingsResponse = Awaited<
  ReturnType<typeof SaveTermTranslationPhaseAISettings>
>

vi.mock("../../../wailsjs/go/wails/AppController.js", () => ({
  GetProcessingTargetList: vi.fn(),
  GetTermTranslationNextPhaseReadiness: vi.fn(),
  GetTermTranslationPhaseSummary: vi.fn(),
  PauseTermTranslationPhase: vi.fn(),
  ResumeTermTranslationPhase: vi.fn(),
  RetryTermTranslationPhase: vi.fn(),
  SaveTermTranslationPhaseAISettings: vi.fn(),
  StartTermTranslationPhase: vi.fn()
}))

afterEach(() => {
  vi.restoreAllMocks()
})

describe("createTermTranslationPhaseGateway", () => {
  test("processing target list request と totalCount response をそのまま受け渡す", async () => {
    vi.mocked(GetProcessingTargetList).mockResolvedValue({
      items: [
        {
          id: "term:1",
          name: "Dragon",
          detail: "原文: Dragon",
          titleParts: [{ text: "Dragon" }],
          metadata: [{ label: "候補", value: "ドラゴン" }],
          convertValues: vi.fn()
        }
      ],
      metadata: [],
      page: 2,
      pageSize: 50,
      totalCount: 137,
      searchQuery: "Dragon",
      convertValues: vi.fn()
    })

    const gateway = createTermTranslationPhaseGateway()
    await expect(
      gateway.getProcessingTargetList?.({
        jobId: 5,
        phase: "term_translation",
        page: 2,
        pageSize: 50,
        searchQuery: "Dragon"
      })
    ).resolves.toMatchObject({
      items: [{ id: "term:1" }],
      page: 2,
      pageSize: 50,
      totalCount: 137,
      searchQuery: "Dragon"
    })

    expect(GetProcessingTargetList).toHaveBeenCalledWith({
      jobId: 5,
      phase: "term_translation",
      page: 2,
      pageSize: 50,
      searchQuery: "Dragon"
    })
  })

  test("summary request を公開 binding wrapper へ渡す", async () => {
    // 公開 seam: gateway request が wrapper へ転送されることを証明する。
    vi.mocked(GetTermTranslationPhaseSummary).mockResolvedValue({
      jobId: 5,
      currentPhase: "term_translation",
      phaseState: "ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 1,
        aiTargetCount: 1,
        currentStep: "ready"
      },
      totalTermCount: 1,
      dictionaryHitCount: 0,
      aiTargetCount: 1,
      execution: {
        credentialRef: "cred",
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      actionEnablement: {
        canStart: true,
        canPause: false,
        canResume: false,
        canRetry: false,
        canStartNextPhase: false
      }
    } as unknown as GetTermTranslationPhaseSummaryResponse)

    const gateway = createTermTranslationPhaseGateway()
    await gateway.getTermTranslationPhaseSummary({ jobId: 5 })

    expect(GetTermTranslationPhaseSummary).toHaveBeenCalledTimes(1)
    expect(GetTermTranslationPhaseSummary).toHaveBeenCalledWith({ jobId: 5 })
  })

  test("next phase readiness は公開 binding wrapper の response を返す", async () => {
    // 公開 seam: response 形を gateway contract として返す。
    vi.mocked(GetTermTranslationNextPhaseReadiness).mockResolvedValue({
      jobId: 5,
      currentPhase: "term_translation",
      phaseState: "ready",
      canStartNextPhase: false,
      blockedReason: "pending"
    } as unknown as GetTermTranslationNextPhaseReadinessResponse)

    const gateway = createTermTranslationPhaseGateway()
    await expect(
      gateway.getTermTranslationNextPhaseReadiness({ jobId: 5 })
    ).resolves.toMatchObject({
      jobId: 5,
      canStartNextPhase: false,
      blockedReason: "pending"
    })
  })

  test("binding 未接続時は wrapper 例外を返す", async () => {
    // 未接続時の error path を公開 seam で観測する。
    vi.mocked(StartTermTranslationPhase).mockRejectedValue(
      new Error("Wails binding is not wired yet: StartTermTranslationPhase")
    )

    const gateway = createTermTranslationPhaseGateway()

    await expect(
      gateway.startTermTranslationPhase({ jobId: 1 })
    ).rejects.toThrowError(
      /Wails binding is not wired yet: StartTermTranslationPhase/
    )
  })

  test("runtime shape 検証失敗時は診断へ secret 平文を出さない", async () => {
    // runtime shape 検証失敗時も secret は公開値へ含めない。
    vi.mocked(GetTermTranslationPhaseSummary).mockResolvedValue({
      jobId: 5,
      currentPhase: "term_translation",
      phaseState: "ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 1,
        aiTargetCount: 1,
        currentStep: "ready"
      },
      totalTermCount: 1,
      dictionaryHitCount: 0,
      aiTargetCount: 1,
      execution: {
        credentialRef: 123,
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch",
        credentialInput: "raw-secret-value"
      },
      actionEnablement: {
        canStart: true,
        canPause: false,
        canResume: false,
        canRetry: false,
        canStartNextPhase: false
      }
    } as unknown as GetTermTranslationPhaseSummaryResponse)

    const gateway = createTermTranslationPhaseGateway()

    await expect(
      gateway.getTermTranslationPhaseSummary({ jobId: 5 })
    ).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })

    try {
      await gateway.getTermTranslationPhaseSummary({ jobId: 5 })
    } catch (error: unknown) {
      expect(error).toBeInstanceOf(Error)
      expect(typeof (error as Error & { internalDiagnostic?: unknown }).internalDiagnostic).toBe("string")
      const diagnostic = JSON.stringify(error)
      expect(diagnostic).not.toContain("raw-secret-value")
      expect(diagnostic).not.toContain("credentialInput")
      expect(diagnostic).not.toContain("apiKey")
    }
  })

  test("save ai settings は公開フィールドだけを受け渡し secret を含まない", async () => {
    // 保存応答に secret 本体が含まれないことを確認する。
    vi.mocked(SaveTermTranslationPhaseAISettings).mockResolvedValue({
      jobId: 5,
      phaseId: "word_translation",
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled",
      credentialStatus: "configured",
      modelListStatus: "success"
    } as unknown as SaveTermTranslationPhaseAISettingsResponse)

    const gateway = createTermTranslationPhaseGateway()
    const result = await gateway.saveTermTranslationPhaseAISettings?.({
      jobId: 5,
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })

    expect(SaveTermTranslationPhaseAISettings).toHaveBeenCalledWith({
      jobId: 5,
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })
    expect(result).toMatchObject({
      jobId: 5,
      phaseId: "word_translation",
      credentialStatus: "configured",
      modelListStatus: "success"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })
})
