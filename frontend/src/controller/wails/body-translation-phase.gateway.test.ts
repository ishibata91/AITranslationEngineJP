import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseSummaryResponse,
  CancelBodyTranslationPhaseRequest,
  GetBodyTranslationOutputReadinessRequest,
  GetBodyTranslationPhaseSummaryRequest,
  PauseBodyTranslationPhaseRequest,
  ResumeBodyTranslationPhaseRequest,
  RetryBodyTranslationPhaseRequest,
  StartBodyTranslationPhaseRequest
} from "@application/gateway-contract/body-translation-phase"
import type {
  CancelBodyTranslationPhaseRequestDto,
  CancelBodyTranslationPhaseResponseDto,
  GetBodyTranslationOutputReadinessRequestDto,
  GetBodyTranslationOutputReadinessResponseDto,
  GetBodyTranslationPhaseSummaryRequestDto,
  GetBodyTranslationPhaseSummaryResponseDto,
  PauseBodyTranslationPhaseRequestDto,
  PauseBodyTranslationPhaseResponseDto,
  ResumeBodyTranslationPhaseRequestDto,
  ResumeBodyTranslationPhaseResponseDto,
  RetryBodyTranslationPhaseRequestDto,
  RetryBodyTranslationPhaseResponseDto,
  StartBodyTranslationPhaseRequestDto,
  StartBodyTranslationPhaseResponseDto
} from "@controller/wails/gateway-dto/body-translation-phase"
import { afterEach, describe, expect, expectTypeOf, test, vi } from "vitest"

import { createBodyTranslationPhaseGateway } from "./body-translation-phase.gateway"

type GoRecord = {
  wails: {
    AppController?: Record<string, ReturnType<typeof vi.fn>>
    BodyTranslationPhaseController?: Record<string, ReturnType<typeof vi.fn>>
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

describe("createBodyTranslationPhaseGateway", () => {
  test("gateway dto aliases stay aligned with gateway contract", () => {
    expectTypeOf<GetBodyTranslationPhaseSummaryRequestDto>().toEqualTypeOf<GetBodyTranslationPhaseSummaryRequest>()
    expectTypeOf<GetBodyTranslationPhaseSummaryResponseDto>().toEqualTypeOf<BodyTranslationPhaseSummaryResponse>()
    expectTypeOf<StartBodyTranslationPhaseRequestDto>().toEqualTypeOf<StartBodyTranslationPhaseRequest>()
    expectTypeOf<StartBodyTranslationPhaseResponseDto>().toEqualTypeOf<BodyTranslationPhaseCommandResponse>()
    expectTypeOf<PauseBodyTranslationPhaseRequestDto>().toEqualTypeOf<PauseBodyTranslationPhaseRequest>()
    expectTypeOf<PauseBodyTranslationPhaseResponseDto>().toEqualTypeOf<BodyTranslationPhaseCommandResponse>()
    expectTypeOf<ResumeBodyTranslationPhaseRequestDto>().toEqualTypeOf<ResumeBodyTranslationPhaseRequest>()
    expectTypeOf<ResumeBodyTranslationPhaseResponseDto>().toEqualTypeOf<BodyTranslationPhaseCommandResponse>()
    expectTypeOf<RetryBodyTranslationPhaseRequestDto>().toEqualTypeOf<RetryBodyTranslationPhaseRequest>()
    expectTypeOf<RetryBodyTranslationPhaseResponseDto>().toEqualTypeOf<BodyTranslationPhaseCommandResponse>()
    expectTypeOf<CancelBodyTranslationPhaseRequestDto>().toEqualTypeOf<CancelBodyTranslationPhaseRequest>()
    expectTypeOf<CancelBodyTranslationPhaseResponseDto>().toEqualTypeOf<BodyTranslationPhaseCommandResponse>()
    expectTypeOf<GetBodyTranslationOutputReadinessRequestDto>().toEqualTypeOf<GetBodyTranslationOutputReadinessRequest>()
    expectTypeOf<GetBodyTranslationOutputReadinessResponseDto>().toEqualTypeOf<BodyTranslationOutputReadinessResponse>()
  })

  test("BodyTranslationPhaseController binding を優先して summary request を渡す", async () => {
    const getSummary = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        currentPhase: "body_translation",
        phaseState: "ready",
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
          dictionaryDigest: "sha256:dictionary",
          personaDigest: "sha256:persona",
          metadataDigest: "sha256:metadata",
          promptDigest: "sha256:prompt"
        },
        execution: {
          credentialRef: "cred",
          provider: "openai-compatible",
          model: "gpt-4.1-mini",
          executionMode: "batch",
          requestUnitCount: 1,
          outputCount: 0
        },
        actionEnablement: {
          canStart: true,
          canPause: false,
          canResume: false,
          canRetry: false,
          canCancel: false,
          canCheckOutputReadiness: true
        },
        outputReadiness: {
          ready: false,
          blockedReason: "pending",
          completedFieldCount: 0,
          statusConsistent: true
        }
      })
    )

    installGo({
      wails: {
        BodyTranslationPhaseController: {
          GetBodyTranslationPhaseSummary: getSummary
        },
        AppController: {
          GetBodyTranslationPhaseSummary: vi.fn(() =>
            Promise.reject(new Error("must not call"))
          )
        }
      }
    })

    const gateway = createBodyTranslationPhaseGateway()
    await gateway.getBodyTranslationPhaseSummary({ jobId: 5 })

    expect(getSummary).toHaveBeenCalledTimes(1)
    expect(getSummary).toHaveBeenCalledWith({ jobId: 5 })
  })

  test("controller binding が無い時は AppController binding に fallback する", async () => {
    const getReadiness = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        currentPhase: "body_translation",
        phaseState: "completed",
        ready: true,
        completedFieldCount: 2,
        statusConsistent: true,
        outputCount: 2
      })
    )

    installGo({
      wails: {
        AppController: {
          GetBodyTranslationOutputReadiness: getReadiness
        }
      }
    })

    const gateway = createBodyTranslationPhaseGateway()
    await gateway.getBodyTranslationOutputReadiness({ jobId: 5 })

    expect(getReadiness).toHaveBeenCalledTimes(1)
    expect(getReadiness).toHaveBeenCalledWith({ jobId: 5 })
  })

  test("binding 未接続時は Wails not wired error を返す", async () => {
    installGo({ wails: {} })

    const gateway = createBodyTranslationPhaseGateway()

    await expect(
      gateway.startBodyTranslationPhase({ jobId: 1 })
    ).rejects.toThrowError(
      /Wails binding is not wired yet: StartBodyTranslationPhase/
    )
  })

  test("save ai settings は公開参照値だけを返し secret 本体を含めない", async () => {
    const saveSettings = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        phaseId: "text_translation",
        provider: "openai",
        model: "gpt-5.4-mini",
        executionMode: "sync",
        batchMode: "unsupported",
        credentialStatus: "configured",
        modelListStatus: "success"
      })
    )
    installGo({
      wails: {
        BodyTranslationPhaseController: {
          SaveBodyTranslationPhaseAISettings: saveSettings
        }
      }
    })

    const gateway = createBodyTranslationPhaseGateway()
    const result = await gateway.saveBodyTranslationPhaseAISettings?.({
      jobId: 5,
      provider: "openai",
      model: "gpt-5.4-mini",
      executionMode: "sync",
      batchMode: "unsupported"
    })

    expect(saveSettings).toHaveBeenCalledWith({
      jobId: 5,
      provider: "openai",
      model: "gpt-5.4-mini",
      executionMode: "sync",
      batchMode: "unsupported"
    })
    expect(result).toMatchObject({
      phaseId: "text_translation",
      credentialStatus: "configured",
      modelListStatus: "success"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })
})
