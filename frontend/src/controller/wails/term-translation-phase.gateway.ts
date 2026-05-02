import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import type {
  GetTermTranslationNextPhaseReadinessRequestDto,
  GetTermTranslationNextPhaseReadinessResponseDto,
  GetTermTranslationPhaseSummaryRequestDto,
  GetTermTranslationPhaseSummaryResponseDto,
  PauseTermTranslationPhaseRequestDto,
  PauseTermTranslationPhaseResponseDto,
  ResumeTermTranslationPhaseRequestDto,
  ResumeTermTranslationPhaseResponseDto,
  RetryTermTranslationPhaseRequestDto,
  RetryTermTranslationPhaseResponseDto,
  StartTermTranslationPhaseRequestDto,
  StartTermTranslationPhaseResponseDto
} from "@controller/wails/gateway-dto/term-translation-phase"

type TermTranslationPhaseBindingName =
  | "GetTermTranslationPhaseSummary"
  | "StartTermTranslationPhase"
  | "PauseTermTranslationPhase"
  | "ResumeTermTranslationPhase"
  | "RetryTermTranslationPhase"
  | "GetTermTranslationNextPhaseReadiness"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: TermTranslationPhaseBindingName,
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
  bindingName: TermTranslationPhaseBindingName
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
    toRecord(wailsRecord["TermTranslationPhaseController"]),
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
    bindingName: TermTranslationPhaseBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-term-phase-wails-gateway 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class TermTranslationPhaseGateway implements TermTranslationPhaseGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getTermTranslationPhaseSummary(
    request: GetTermTranslationPhaseSummaryRequestDto
  ): Promise<GetTermTranslationPhaseSummaryResponseDto> {
    return this.invokeBinding("GetTermTranslationPhaseSummary", request)
  }

  startTermTranslationPhase(
    request: StartTermTranslationPhaseRequestDto
  ): Promise<StartTermTranslationPhaseResponseDto> {
    return this.invokeBinding("StartTermTranslationPhase", request)
  }

  pauseTermTranslationPhase(
    request: PauseTermTranslationPhaseRequestDto
  ): Promise<PauseTermTranslationPhaseResponseDto> {
    return this.invokeBinding("PauseTermTranslationPhase", request)
  }

  resumeTermTranslationPhase(
    request: ResumeTermTranslationPhaseRequestDto
  ): Promise<ResumeTermTranslationPhaseResponseDto> {
    return this.invokeBinding("ResumeTermTranslationPhase", request)
  }

  retryTermTranslationPhase(
    request: RetryTermTranslationPhaseRequestDto
  ): Promise<RetryTermTranslationPhaseResponseDto> {
    return this.invokeBinding("RetryTermTranslationPhase", request)
  }

  getTermTranslationNextPhaseReadiness(
    request: GetTermTranslationNextPhaseReadinessRequestDto
  ): Promise<GetTermTranslationNextPhaseReadinessResponseDto> {
    return this.invokeBinding("GetTermTranslationNextPhaseReadiness", request)
  }
}

export function createTermTranslationPhaseGateway(): TermTranslationPhaseGatewayContract {
  return new TermTranslationPhaseGateway(createBindingInvoker())
}
