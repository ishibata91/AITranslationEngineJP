import type {
  ProcessingTargetListRequest,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"
import type {
  CancelPersonaGenerationPhaseRequest,
  GetPersonaGenerationBodyReadinessRequest,
  GetPersonaGenerationPhaseSummaryRequest,
  PausePersonaGenerationPhaseRequest,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseCommandResponse,
  PersonaGenerationPhaseFetchResponse,
  ResumePersonaGenerationPhaseRequest,
  RetryPersonaGenerationPhaseRequest,
  StartPersonaGenerationPhaseRequest
} from "@application/gateway-contract/persona-generation-phase"
import type {
  CancelPersonaGenerationPhaseRequestDto,
  CancelPersonaGenerationPhaseResponseDto,
  GetProcessingTargetListRequestDto,
  GetProcessingTargetListResponseDto,
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
    ProcessingTargetController?: Record<string, ReturnType<typeof vi.fn>>
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
    expectTypeOf<GetPersonaGenerationPhaseSummaryResponseDto>().toEqualTypeOf<PersonaGenerationPhaseFetchResponse>()
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
    expectTypeOf<GetProcessingTargetListRequestDto>().toEqualTypeOf<ProcessingTargetListRequest>()
    expectTypeOf<GetProcessingTargetListResponseDto>().toEqualTypeOf<ProcessingTargetListResponse>()
  })

  test("processing target list は request と totalCount response を保持する", async () => {
    const getProcessingTargetList = vi.fn(() =>
      Promise.resolve({
        items: [
          {
            id: "persona:1",
            name: "Aela",
            detail: "FormID: 0001A696",
            titleParts: [{ text: "Aela" }],
            metadata: [{ label: "種族", value: "Nord" }]
          }
        ],
        metadata: [],
        page: 3,
        pageSize: 25,
        totalCount: 88,
        searchQuery: "Nord"
      })
    )

    installGo({
      wails: {
        ProcessingTargetController: {
          GetProcessingTargetList: getProcessingTargetList
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()
    await expect(
      gateway.getProcessingTargetList?.({
        jobId: 8,
        phase: "persona_generation",
        page: 3,
        pageSize: 25,
        searchQuery: "Nord"
      })
    ).resolves.toMatchObject({
      items: [{ id: "persona:1" }],
      page: 3,
      pageSize: 25,
      totalCount: 88,
      searchQuery: "Nord"
    })

    expect(getProcessingTargetList).toHaveBeenCalledWith({
      jobId: 8,
      phase: "persona_generation",
      page: 3,
      pageSize: 25,
      searchQuery: "Nord"
    })
  })

  test("processing target list の DTO shape が崩れた時は検証エラーを返す", async () => {
    const getProcessingTargetList = vi.fn(() =>
      Promise.resolve({
        items: [],
        metadata: [],
        page: 1,
        pageSize: 50,
        totalCount: "88",
        searchQuery: "Nord"
      })
    )

    installGo({
      wails: {
        ProcessingTargetController: {
          GetProcessingTargetList: getProcessingTargetList
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()

    await expect(
      gateway.getProcessingTargetList?.({
        jobId: 8,
        phase: "persona_generation",
        page: 1,
        pageSize: 50,
        searchQuery: "Nord"
      })
    ).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })
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
    // SaveAISettings 入力から job_id を除外する仕様に従う。
    const saveSettings = vi.fn(() =>
      Promise.resolve({
        phaseType: "npc_persona_generation",
        provider: "xai",
        model: "grok-4",
        executionMode: "sync",
        batchMode: "disabled"
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
      provider: "xai",
      model: "grok-4",
      executionMode: "sync",
      batchMode: "disabled"
    })

    expect(saveSettings).toHaveBeenCalledWith({
      provider: "xai",
      model: "grok-4",
      executionMode: "sync",
      batchMode: "disabled"
    })
    expect(result).toMatchObject({
      phaseType: "npc_persona_generation",
      provider: "xai"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })

  // RAEF-UNIT-022: assertPersonaProjectionShape — previousPhaseLifecycle が存在しない場合は検証失敗する
  // 根拠: design-diff.md G-4-b

  // RAEF-UNIT-022: previousPhaseLifecycle が存在しない場合は GatewayResponseShapeError を返す
  test("projection の previousPhaseLifecycle が存在しない場合は GatewayResponseShapeError を返す（G-4-b projection 必須 field 不足）", async () => {
    // G-4-b: projection.previousPhaseLifecycle が存在しない → 必須 field 不足で検証失敗
    const getPersonaSummary = vi.fn(() =>
      Promise.resolve({
        projection: {
          // previousPhaseLifecycle を意図的に除外する
          phaseLifecycle: "idle_ready",
          jobLifecycle: "running",
          errorKind: "none",
          aiSettingsConfigured: true,
          targetCount: 5
          // previousPhaseLifecycle が存在しない
        },
        summary: {
          jobId: 8,
          currentPhase: "persona_generation",
          phaseState: "idle_ready"
        }
      })
    )

    installGo({
      wails: {
        PersonaGenerationPhaseController: {
          GetPersonaGenerationPhaseSummary: getPersonaSummary
        }
      }
    })

    const gateway = createPersonaGenerationPhaseGateway()

    await expect(gateway.getPersonaGenerationPhaseSummary({ jobId: 8 })).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })
  })
})
