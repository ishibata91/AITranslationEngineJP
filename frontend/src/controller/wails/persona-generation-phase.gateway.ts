import type { PersonaGenerationPhaseGatewayContract } from "@application/gateway-contract/persona-generation-phase"
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

type PersonaGenerationPhaseBindingName =
  | "GetPersonaGenerationPhaseSummary"
  | "StartPersonaGenerationPhase"
  | "PausePersonaGenerationPhase"
  | "ResumePersonaGenerationPhase"
  | "RetryPersonaGenerationPhase"
  | "CancelPersonaGenerationPhase"
  | "GetPersonaGenerationBodyReadiness"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: PersonaGenerationPhaseBindingName,
  request: RequestDto
) => Promise<ResponseDto>

type BindingFunction = (request: unknown) => Promise<unknown>

function toRecord(value: unknown): Record<string, unknown> | null {
  if (typeof value !== "object" || value === null) {
    return null
  }

  return value as Record<string, unknown>
}

function resolveBindingFunction(
  bindingName: PersonaGenerationPhaseBindingName
): BindingFunction | null {
  const globalRecord = toRecord(globalThis)
  const goRecord = toRecord(globalRecord?.["go"])
  if (!goRecord) {
    return null
  }

  const wailsRecord = toRecord(goRecord["wails"])
  if (!wailsRecord) {
    return null
  }

  const controllerCandidates = [
    toRecord(wailsRecord["PersonaGenerationPhaseController"]),
    toRecord(wailsRecord["AppController"])
  ]

  for (const controllerRecord of controllerCandidates) {
    if (!controllerRecord) {
      continue
    }

    const binding = controllerRecord[bindingName]
    if (typeof binding !== "function") {
      continue
    }

    return (request: unknown) =>
      Promise.resolve((binding as (arg: unknown) => unknown)(request))
  }

  return null
}

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    bindingName: PersonaGenerationPhaseBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-persona-phase-wails-gateway 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class PersonaGenerationPhaseGateway
  implements PersonaGenerationPhaseGatewayContract
{
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getPersonaGenerationPhaseSummary(
    request: GetPersonaGenerationPhaseSummaryRequestDto
  ): Promise<GetPersonaGenerationPhaseSummaryResponseDto> {
    return this.invokeBinding("GetPersonaGenerationPhaseSummary", request)
  }

  startPersonaGenerationPhase(
    request: StartPersonaGenerationPhaseRequestDto
  ): Promise<StartPersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("StartPersonaGenerationPhase", request)
  }

  pausePersonaGenerationPhase(
    request: PausePersonaGenerationPhaseRequestDto
  ): Promise<PausePersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("PausePersonaGenerationPhase", request)
  }

  resumePersonaGenerationPhase(
    request: ResumePersonaGenerationPhaseRequestDto
  ): Promise<ResumePersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("ResumePersonaGenerationPhase", request)
  }

  retryPersonaGenerationPhase(
    request: RetryPersonaGenerationPhaseRequestDto
  ): Promise<RetryPersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("RetryPersonaGenerationPhase", request)
  }

  cancelPersonaGenerationPhase(
    request: CancelPersonaGenerationPhaseRequestDto
  ): Promise<CancelPersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("CancelPersonaGenerationPhase", request)
  }

  getPersonaGenerationBodyReadiness(
    request: GetPersonaGenerationBodyReadinessRequestDto
  ): Promise<GetPersonaGenerationBodyReadinessResponseDto> {
    return this.invokeBinding("GetPersonaGenerationBodyReadiness", request)
  }
}

export function createPersonaGenerationPhaseGateway(): PersonaGenerationPhaseGatewayContract {
  return new PersonaGenerationPhaseGateway(createBindingInvoker())
}
