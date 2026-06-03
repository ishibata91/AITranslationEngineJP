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

// 有効な term phase summary response の最小 fixture を返すヘルパー
function createValidTermSummaryResponse(): GetTermTranslationPhaseSummaryResponse {
  return {
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
    projection: {
      phaseLifecycle: "ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 1,
      confirmedCount: 0
    }
  } as unknown as GetTermTranslationPhaseSummaryResponse
}

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
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        aiTargetCount: 1,
        confirmedCount: 0
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
      jobIsTerminal: false,
      totalCount: 10,
      confirmedCount: 0
    } as unknown as GetTermTranslationNextPhaseReadinessResponse)

    const gateway = createTermTranslationPhaseGateway()
    await expect(
      gateway.getTermTranslationNextPhaseReadiness({ jobId: 5 })
    ).resolves.toMatchObject({
      jobId: 5,
      jobIsTerminal: false,
      totalCount: 10
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
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        aiTargetCount: 1,
        confirmedCount: 0
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

  // RAEF-UNIT-019〜021: assertTermProjectionShape の runtime shape validator 境界値証明
  // 根拠: design-diff.md G-4-a

  // RAEF-UNIT-019: projection 必須 field が全て揃う場合は shape 検証通過する（正常パス）
  test("projection 必須 field が全て揃う場合は shape 検証が通過する（G-4-a projection 正常パス）", async () => {
    // G-4-a: projection に phaseLifecycle / jobLifecycle / errorKind / aiSettingsConfigured / aiTargetCount / confirmedCount が揃う → 検証通過
    vi.mocked(GetTermTranslationPhaseSummary).mockResolvedValue(createValidTermSummaryResponse())

    const gateway = createTermTranslationPhaseGateway()

    await expect(gateway.getTermTranslationPhaseSummary({ jobId: 5 })).resolves.toMatchObject({
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        aiTargetCount: 1,
        confirmedCount: 0
      }
    })
  })

  // RAEF-UNIT-020: actionEnablement field が response に含まれない場合でも検証通過する（削除後の非必須化確認）
  test("response に actionEnablement field が存在しない場合でも検証が通過する（G-4-a actionEnablement 削除後の非必須化）", async () => {
    // G-4-a: actionEnablement は削除済みのため response に含まれなくても検証通過する
    // motivating bug（canStartNextPhase 必須検証で落ちる）の解消を証明する
    const responseWithoutActionEnablement = createValidTermSummaryResponse()
    // actionEnablement field が存在しないことを確認する
    expect((responseWithoutActionEnablement as unknown as Record<string, unknown>)["actionEnablement"]).toBeUndefined()

    vi.mocked(GetTermTranslationPhaseSummary).mockResolvedValue(responseWithoutActionEnablement)

    const gateway = createTermTranslationPhaseGateway()

    await expect(gateway.getTermTranslationPhaseSummary({ jobId: 5 })).resolves.toBeDefined()
  })

  // RAEF-UNIT-021: phaseLifecycle が存在しない場合は検証失敗する
  test("projection の phaseLifecycle が存在しない場合は GatewayResponseShapeError を返す（G-4-a projection 必須 field 不足）", async () => {
    // G-4-a: projection.phaseLifecycle が存在しない → 必須 field 不足で検証失敗
    const invalidResponse = {
      ...createValidTermSummaryResponse(),
      projection: {
        // phaseLifecycle を意図的に除外する
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        aiTargetCount: 1,
        confirmedCount: 0
      }
    } as unknown as GetTermTranslationPhaseSummaryResponse

    vi.mocked(GetTermTranslationPhaseSummary).mockResolvedValue(invalidResponse)

    const gateway = createTermTranslationPhaseGateway()

    await expect(gateway.getTermTranslationPhaseSummary({ jobId: 5 })).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })
  })

  test("save ai settings は公開フィールドだけを受け渡し secret を含まない", async () => {
    // 保存応答に secret 本体が含まれないことを確認する。
    // SaveAISettings 入力から job_id を除外する仕様に従う。
    vi.mocked(SaveTermTranslationPhaseAISettings).mockResolvedValue({
      phaseType: "word_translation",
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    } as unknown as SaveTermTranslationPhaseAISettingsResponse)

    const gateway = createTermTranslationPhaseGateway()
    const result = await gateway.saveTermTranslationPhaseAISettings?.({
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })

    expect(SaveTermTranslationPhaseAISettings).toHaveBeenCalledWith({
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })
    expect(result).toMatchObject({
      phaseType: "word_translation",
      provider: "gemini"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })
})
