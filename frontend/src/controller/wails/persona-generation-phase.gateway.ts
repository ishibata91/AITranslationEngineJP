import type {
  PersonaGenerationPhaseGatewayContract,
  PhaseProviderModelsRequest,
  PhaseProviderModelsResponse
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
  SavePersonaGenerationPhaseAISettingsRequestDto,
  SavePersonaGenerationPhaseAISettingsResponseDto,
  StartPersonaGenerationPhaseRequestDto,
  StartPersonaGenerationPhaseResponseDto
} from "@controller/wails/gateway-dto/persona-generation-phase"
import {
  createGatewayResponseShapeError,
  isArrayOf,
  isBoolean,
  isNumber,
  isRecord,
  isString
} from "@controller/wails/gateway-dto/runtime-shape"

type RuntimeShapeValidator<ResponseDto> = (
  value: unknown
) => value is ResponseDto

type RuntimeShapeIssue = {
  path: string
  expected: string
}

type PersonaGenerationPhaseBindingName =
  | "GetProcessingTargetList"
  | "GetPersonaGenerationPhaseSummary"
  | "StartPersonaGenerationPhase"
  | "SavePersonaGenerationPhaseAISettings"
  | "PausePersonaGenerationPhase"
  | "ResumePersonaGenerationPhase"
  | "RetryPersonaGenerationPhase"
  | "CancelPersonaGenerationPhase"
  | "GetPersonaGenerationBodyReadiness"
  | "ListTranslationJobSetupProviderModels"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: PersonaGenerationPhaseBindingName,
  request: RequestDto,
  isResponseDto?: RuntimeShapeValidator<ResponseDto>
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
    toRecord(wailsRecord["ProcessingTargetController"]),
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

const responseShapeIssues: RuntimeShapeIssue[] = []

function resetIssues(): void {
  responseShapeIssues.length = 0
}

function invalid(path: string, expected: string): false {
  responseShapeIssues.push({ path, expected })
  return false
}

function isMetadata(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["label"]) || invalid(`${path}.label`, "string")) &&
    (isString(value["value"]) || invalid(`${path}.value`, "string"))
  )
}

function isProcessingTargetItem(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["id"]) || invalid(`${path}.id`, "string")) &&
    (isString(value["name"]) || invalid(`${path}.name`, "string")) &&
    (isString(value["detail"]) || invalid(`${path}.detail`, "string")) &&
    (isArrayOf(value["titleParts"], (item) => {
      if (!isRecord(item)) {
        return invalid(`${path}.titleParts[]`, "object")
      }
      return (
        isString(item["text"]) || invalid(`${path}.titleParts[].text`, "string")
      )
    }) ||
      invalid(`${path}.titleParts`, "title part array")) &&
    (isArrayOf(value["metadata"], (item) =>
      isMetadata(item, `${path}.metadata[]`)
    ) ||
      invalid(`${path}.metadata`, "metadata array"))
  )
}

function isProcessingTargetListResponseDto(
  value: unknown
): value is GetProcessingTargetListResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isArrayOf(value["items"], (item) =>
      isProcessingTargetItem(item, "$.items[]")
    ) ||
      invalid("$.items", "processing target item array")) &&
    (isArrayOf(value["metadata"], (item) => isMetadata(item, "$.metadata[]")) ||
      invalid("$.metadata", "metadata array")) &&
    (isNumber(value["page"]) || invalid("$.page", "number")) &&
    (isNumber(value["pageSize"]) || invalid("$.pageSize", "number")) &&
    (isNumber(value["totalCount"]) || invalid("$.totalCount", "number")) &&
    (isString(value["searchQuery"]) || invalid("$.searchQuery", "string"))
  )
}

function isProviderModelOption(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }
  return (
    (isString(value["modelId"]) || invalid(`${path}.modelId`, "string")) &&
    (isString(value["label"]) || invalid(`${path}.label`, "string"))
  )
}

function isListProviderModelsResponseDto(value: unknown): value is {
  provider: string
  status: string
  models: { modelId: string; label: string }[]
  failureKind?: string
} {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }
  return (
    (isString(value["provider"]) || invalid("$.provider", "string")) &&
    (isString(value["status"]) || invalid("$.status", "string")) &&
    (isArrayOf(value["models"], (item) =>
      isProviderModelOption(item, "$.models[]")
    ) ||
      invalid("$.models", "model option array")) &&
    (value["failureKind"] === undefined ||
      isString(value["failureKind"]) ||
      invalid("$.failureKind", "string or undefined"))
  )
}

function isPersonaGenerationPhaseProjection(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["phaseLifecycle"]) || invalid(`${path}.phaseLifecycle`, "string")) &&
    (isString(value["jobLifecycle"]) || invalid(`${path}.jobLifecycle`, "string")) &&
    (isString(value["errorKind"]) || invalid(`${path}.errorKind`, "string")) &&
    (isBoolean(value["aiSettingsConfigured"]) ||
      invalid(`${path}.aiSettingsConfigured`, "boolean")) &&
    (isNumber(value["targetCount"]) || invalid(`${path}.targetCount`, "number")) &&
    (isString(value["previousPhaseLifecycle"]) ||
      invalid(`${path}.previousPhaseLifecycle`, "string"))
  )
}

function isPersonaGenerationPhaseSummaryShape(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid(`${path}.jobId`, "number")) &&
    (isString(value["currentPhase"]) || invalid(`${path}.currentPhase`, "string")) &&
    (isString(value["phaseState"]) || invalid(`${path}.phaseState`, "string"))
  )
}

function isPersonaGenerationPhaseFetchResponseDto(
  value: unknown
): value is GetPersonaGenerationPhaseSummaryResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    isPersonaGenerationPhaseProjection(value["projection"], "$.projection") &&
    isPersonaGenerationPhaseSummaryShape(value["summary"], "$.summary")
  )
}

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    bindingName: PersonaGenerationPhaseBindingName,
    request: RequestDto,
    isResponseDto?: RuntimeShapeValidator<ResponseDto>
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-persona-phase-wails-gateway 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => {
      if (!isResponseDto || isResponseDto(response)) {
        return response as ResponseDto
      }

      throw createGatewayResponseShapeError(bindingName, responseShapeIssues)
    })
  }
}

class PersonaGenerationPhaseGateway implements PersonaGenerationPhaseGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  listProviderModels(
    request: PhaseProviderModelsRequest
  ): Promise<PhaseProviderModelsResponse> {
    return this.invokeBinding(
      "ListTranslationJobSetupProviderModels",
      {
        phaseId: "npc_persona_generation",
        provider: request.provider,
        credentialStatus: request.credentialStatus,
        requestToken: request.requestToken
      },
      isListProviderModelsResponseDto
    ).then((dto) => ({
      provider: dto.provider,
      status: dto.status,
      models: dto.models.map((m) => ({ modelId: m.modelId, label: m.label })),
      failureKind: dto.failureKind
    }))
  }

  getProcessingTargetList(
    request: GetProcessingTargetListRequestDto
  ): Promise<GetProcessingTargetListResponseDto> {
    return this.invokeBinding(
      "GetProcessingTargetList",
      request,
      isProcessingTargetListResponseDto
    )
  }

  getPersonaGenerationPhaseSummary(
    request: GetPersonaGenerationPhaseSummaryRequestDto
  ): Promise<GetPersonaGenerationPhaseSummaryResponseDto> {
    return this.invokeBinding(
      "GetPersonaGenerationPhaseSummary",
      request,
      isPersonaGenerationPhaseFetchResponseDto
    )
  }

  startPersonaGenerationPhase(
    request: StartPersonaGenerationPhaseRequestDto
  ): Promise<StartPersonaGenerationPhaseResponseDto> {
    return this.invokeBinding("StartPersonaGenerationPhase", request)
  }

  savePersonaGenerationPhaseAISettings(
    request: SavePersonaGenerationPhaseAISettingsRequestDto
  ): Promise<SavePersonaGenerationPhaseAISettingsResponseDto> {
    return this.invokeBinding("SavePersonaGenerationPhaseAISettings", request)
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
