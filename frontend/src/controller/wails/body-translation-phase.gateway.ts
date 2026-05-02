import type { BodyTranslationPhaseGatewayContract } from "@application/gateway-contract/body-translation-phase"
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

type BodyTranslationPhaseBindingName =
  | "GetBodyTranslationPhaseSummary"
  | "StartBodyTranslationPhase"
  | "PauseBodyTranslationPhase"
  | "ResumeBodyTranslationPhase"
  | "RetryBodyTranslationPhase"
  | "CancelBodyTranslationPhase"
  | "GetBodyTranslationOutputReadiness"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: BodyTranslationPhaseBindingName,
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
  bindingName: BodyTranslationPhaseBindingName
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
    toRecord(wailsRecord["BodyTranslationPhaseController"]),
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
    bindingName: BodyTranslationPhaseBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-body-phase-wails-gateway 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class BodyTranslationPhaseGateway implements BodyTranslationPhaseGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getBodyTranslationPhaseSummary(
    request: GetBodyTranslationPhaseSummaryRequestDto
  ): Promise<GetBodyTranslationPhaseSummaryResponseDto> {
    return this.invokeBinding("GetBodyTranslationPhaseSummary", request)
  }

  startBodyTranslationPhase(
    request: StartBodyTranslationPhaseRequestDto
  ): Promise<StartBodyTranslationPhaseResponseDto> {
    return this.invokeBinding("StartBodyTranslationPhase", request)
  }

  pauseBodyTranslationPhase(
    request: PauseBodyTranslationPhaseRequestDto
  ): Promise<PauseBodyTranslationPhaseResponseDto> {
    return this.invokeBinding("PauseBodyTranslationPhase", request)
  }

  resumeBodyTranslationPhase(
    request: ResumeBodyTranslationPhaseRequestDto
  ): Promise<ResumeBodyTranslationPhaseResponseDto> {
    return this.invokeBinding("ResumeBodyTranslationPhase", request)
  }

  retryBodyTranslationPhase(
    request: RetryBodyTranslationPhaseRequestDto
  ): Promise<RetryBodyTranslationPhaseResponseDto> {
    return this.invokeBinding("RetryBodyTranslationPhase", request)
  }

  cancelBodyTranslationPhase(
    request: CancelBodyTranslationPhaseRequestDto
  ): Promise<CancelBodyTranslationPhaseResponseDto> {
    return this.invokeBinding("CancelBodyTranslationPhase", request)
  }

  getBodyTranslationOutputReadiness(
    request: GetBodyTranslationOutputReadinessRequestDto
  ): Promise<GetBodyTranslationOutputReadinessResponseDto> {
    return this.invokeBinding("GetBodyTranslationOutputReadiness", request)
  }
}

export function createBodyTranslationPhaseGateway(): BodyTranslationPhaseGatewayContract {
  return new BodyTranslationPhaseGateway(createBindingInvoker())
}
