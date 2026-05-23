import type {
  CancelPersonaGenerationPhaseRequest,
  GetPersonaGenerationBodyReadinessRequest,
  GetPersonaGenerationPhaseSummaryRequest,
  PausePersonaGenerationPhaseRequest,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseCommandResponse,
  PersonaGenerationPhaseSummaryResponse,
  ResumePersonaGenerationPhaseRequest,
  RetryPersonaGenerationPhaseRequest,
  StartPersonaGenerationPhaseRequest
} from "@application/gateway-contract/persona-generation-phase"
import type {
  CancelPersonaGenerationPhaseRequestDto,
  CancelPersonaGenerationPhaseResponseDto,
  GetPersonaGenerationBodyReadinessRequestDto,
  GetPersonaGenerationBodyReadinessResponseDto,
  GetPersonaGenerationPhaseSummaryRequestDto,
  GetPersonaGenerationPhaseSummaryResponseDto,
  PausePersonaGenerationPhaseRequestDto,
  PausePersonaGenerationPhaseResponseDto,
  ResumePersonaGenerationPhaseRequestDto,
  ResumePersonaGenerationPhaseResponseDto,
  RetryPersonaGenerationPhaseRequestDto,
  RetryPersonaGenerationPhaseResponseDto,
  StartPersonaGenerationPhaseRequestDto,
  StartPersonaGenerationPhaseResponseDto
} from "@controller/wails/gateway-dto/persona-generation-phase"
import { afterEach, describe, expect, expectTypeOf, test, vi } from "vitest"

import { createPersonaGenerationPhaseGateway } from "./persona-generation-phase.gateway"

type GoRecord = {
  wails: {
    AppController?: Record<string, ReturnType<typeof vi.fn>>
    PersonaGenerationPhaseController?: Record<string, ReturnType<typeof vi.fn>>
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

describe("createPersonaGenerationPhaseGateway", () => {
  test("gateway dto aliases stay aligned with gateway contract", () => {
    expectTypeOf<GetPersonaGenerationPhaseSummaryRequestDto>().toEqualTypeOf<GetPersonaGenerationPhaseSummaryRequest>()
    expectTypeOf<GetPersonaGenerationPhaseSummaryResponseDto>().toEqualTypeOf<PersonaGenerationPhaseSummaryResponse>()
    expectTypeOf<StartPersonaGenerationPhaseRequestDto>().toEqualTypeOf<StartPersonaGenerationPhaseRequest>()
    expectTypeOf<StartPersonaGenerationPhaseResponseDto>().toEqualTypeOf<PersonaGenerationPhaseCommandResponse>()
    expectTypeOf<PausePersonaGenerationPhaseRequestDto>().toEqualTypeOf<PausePersonaGenerationPhaseRequest>()
    expectTypeOf<PausePersonaGenerationPhaseResponseDto>().toEqualTypeOf<PersonaGenerationPhaseCommandResponse>()
    expectTypeOf<ResumePersonaGenerationPhaseRequestDto>().toEqualTypeOf<ResumePersonaGenerationPhaseRequest>()
    expectTypeOf<ResumePersonaGenerationPhaseResponseDto>().toEqualTypeOf<PersonaGenerationPhaseCommandResponse>()
    expectTypeOf<RetryPersonaGenerationPhaseRequestDto>().toEqualTypeOf<RetryPersonaGenerationPhaseRequest>()
    expectTypeOf<RetryPersonaGenerationPhaseResponseDto>().toEqualTypeOf<PersonaGenerationPhaseCommandResponse>()
    expectTypeOf<CancelPersonaGenerationPhaseRequestDto>().toEqualTypeOf<CancelPersonaGenerationPhaseRequest>()
    expectTypeOf<CancelPersonaGenerationPhaseResponseDto>().toEqualTypeOf<PersonaGenerationPhaseCommandResponse>()
    expectTypeOf<GetPersonaGenerationBodyReadinessRequestDto>().toEqualTypeOf<GetPersonaGenerationBodyReadinessRequest>()
    expectTypeOf<GetPersonaGenerationBodyReadinessResponseDto>().toEqualTypeOf<PersonaGenerationBodyReadinessResponse>()
  })

  test("PersonaGenerationPhaseController binding を優先して start request を渡す", async () => {
    const startPhase = vi.fn(() =>
      Promise.resolve({
        jobId: 8,
        currentPhase: "persona_generation",
        phaseState: "running",
        phaseRunId: 13,
        progress: {
          percent: 10,
          processedCount: 1,
          totalCount: 10,
          targetCount: 10,
          currentStep: "generating"
        },
        targetSummary: {
          targetCount: 10,
          commonPersonaHitCount: 2,
          commonPersonaMissCount: 8,
          skippedCount: 0,
          skippedReasons: [],
          targetSnapshotDigest: "sha256:target"
        },
        execution: {
          credentialRef: "cred-1",
          provider: "xai",
          model: "grok-2",
          executionMode: "single_request",
          promptDigest: "sha256:prompt",
          inputCount: 10,
          outputCount: 1,
          evidenceRefs: ["evidence:npc:1"]
        },
        retryable: false,
        canStartBodyPhase: false
      })
    )

    installGo({
      wails: {
        PersonaGenerationPhaseController: {
          StartPersonaGenerationPhase: startPhase
        },
        AppController: {
          StartPersonaGenerationPhase: vi.fn(() =>
            Promise.reject(new Error("must not call"))
          )
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()
    await gateway.startPersonaGenerationPhase({ jobId: 8 })

    expect(startPhase).toHaveBeenCalledTimes(1)
    expect(startPhase).toHaveBeenCalledWith({ jobId: 8 })
  })

  test("controller binding が無い時は AppController binding に fallback する", async () => {
    const getBodyReadiness = vi.fn(() =>
      Promise.resolve({
        jobId: 8,
        currentPhase: "persona_generation",
        phaseState: "completed",
        ready: true,
        inputSummary: {
          personaCount: 5,
          missingCount: 0,
          snapshotId: "snapshot-1",
          snapshotDigest: "sha256:snapshot",
          evidenceRefs: ["evidence:npc:1"]
        }
      })
    )

    installGo({
      wails: {
        AppController: {
          GetPersonaGenerationBodyReadiness: getBodyReadiness
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()
    await gateway.getPersonaGenerationBodyReadiness({ jobId: 8 })

    expect(getBodyReadiness).toHaveBeenCalledTimes(1)
    expect(getBodyReadiness).toHaveBeenCalledWith({ jobId: 8 })
  })

  test("binding 未接続時は Wails not wired error を返す", async () => {
    installGo({ wails: {} })

    const gateway = createPersonaGenerationPhaseGateway()

    await expect(
      gateway.getPersonaGenerationPhaseSummary({ jobId: 1 })
    ).rejects.toThrowError(
      /Wails binding is not wired yet: GetPersonaGenerationPhaseSummary/
    )
  })

  test("save ai settings は公開参照値のみを転送し secret 本体を含めない", async () => {
    const saveSettings = vi.fn(() =>
      Promise.resolve({
        jobId: 8,
        phaseId: "npc_persona_generation",
        provider: "xai",
        model: "grok-4",
        executionMode: "sync",
        batchMode: "disabled",
        credentialStatus: "configured",
        modelListStatus: "success"
      })
    )
    installGo({
      wails: {
        PersonaGenerationPhaseController: {
          SavePersonaGenerationPhaseAISettings: saveSettings
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()
    const result = await gateway.savePersonaGenerationPhaseAISettings?.({
      jobId: 8,
      provider: "xai",
      model: "grok-4",
      executionMode: "sync",
      batchMode: "disabled"
    })

    expect(saveSettings).toHaveBeenCalledWith({
      jobId: 8,
      provider: "xai",
      model: "grok-4",
      executionMode: "sync",
      batchMode: "disabled"
    })
    expect(result).toMatchObject({
      phaseId: "npc_persona_generation",
      credentialStatus: "configured",
      modelListStatus: "success"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })
})
