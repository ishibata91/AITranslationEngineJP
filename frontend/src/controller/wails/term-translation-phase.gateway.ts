import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import {
  GetProcessingTargetList,
  GetTermTranslationNextPhaseReadiness,
  GetTermTranslationPhaseSummary,
  PauseTermTranslationPhase,
  ResumeTermTranslationPhase,
  RetryTermTranslationPhase,
  SaveTermTranslationPhaseAISettings,
  StartTermTranslationPhase
} from "../../../wailsjs/go/wails/AppController.js"
import type {
  GetProcessingTargetListRequestDto,
  GetProcessingTargetListResponseDto,
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
  SaveTermTranslationPhaseAISettingsRequestDto,
  SaveTermTranslationPhaseAISettingsResponseDto,
  StartTermTranslationPhaseRequestDto,
  StartTermTranslationPhaseResponseDto
} from "@controller/wails/gateway-dto/term-translation-phase"
import {
  createGatewayResponseShapeError,
  isArrayOf,
  isBoolean,
  isNumber,
  isOptionalNumber,
  isOptionalString,
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

type BindingInvoker = <RequestDto, ResponseDto>(
  binding: WailsTermTranslationBinding<RequestDto>,
  bindingName: string,
  request: RequestDto,
  isResponseDto: RuntimeShapeValidator<ResponseDto>
) => Promise<ResponseDto>

type WailsTermTranslationBinding<RequestDto> = (
  request: RequestDto
) => Promise<unknown>

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    binding: WailsTermTranslationBinding<RequestDto>,
    bindingName: string,
    request: RequestDto,
    isResponseDto: RuntimeShapeValidator<ResponseDto>
  ): Promise<ResponseDto> => {
    return binding(request).then((response) => {
      if (isResponseDto(response)) {
        return response
      }

      throw createGatewayResponseShapeError(bindingName, responseShapeIssues)
    })
  }
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

function isProgress(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isNumber(value["percent"]) || invalid(`${path}.percent`, "number")) &&
    (isNumber(value["processedCount"]) ||
      invalid(`${path}.processedCount`, "number")) &&
    (isNumber(value["totalCount"]) ||
      invalid(`${path}.totalCount`, "number")) &&
    (isNumber(value["aiTargetCount"]) ||
      invalid(`${path}.aiTargetCount`, "number")) &&
    (isString(value["currentStep"]) || invalid(`${path}.currentStep`, "string"))
  )
}

function isExecution(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["credentialRef"]) ||
      invalid(`${path}.credentialRef`, "string")) &&
    (isString(value["provider"]) || invalid(`${path}.provider`, "string")) &&
    (isString(value["model"]) || invalid(`${path}.model`, "string")) &&
    (isString(value["executionMode"]) ||
      invalid(`${path}.executionMode`, "string")) &&
    (isOptionalString(value["snapshotDigest"]) ||
      invalid(`${path}.snapshotDigest`, "string or undefined")) &&
    (isOptionalString(value["snapshotVersion"]) ||
      invalid(`${path}.snapshotVersion`, "string or undefined"))
  )
}

function isResultSummary(value: unknown, path: string): boolean {
  if (value === undefined) {
    return true
  }
  if (!isRecord(value)) {
    return invalid(path, "object or undefined")
  }

  return (
    (isNumber(value["confirmedCount"]) ||
      invalid(`${path}.confirmedCount`, "number")) &&
    (isNumber(value["jobDictionaryAppliedCount"]) ||
      invalid(`${path}.jobDictionaryAppliedCount`, "number")) &&
    (isNumber(value["replacementTargetCount"]) ||
      invalid(`${path}.replacementTargetCount`, "number")) &&
    (isNumber(value["unmatchedCount"]) ||
      invalid(`${path}.unmatchedCount`, "number")) &&
    (isBoolean(value["providerSkipped"]) ||
      invalid(`${path}.providerSkipped`, "boolean"))
  )
}

function isErrorSummary(value: unknown, path: string): boolean {
  if (value === undefined) {
    return true
  }
  if (!isRecord(value)) {
    return invalid(path, "object or undefined")
  }

  return (
    (isString(value["errorKind"]) || invalid(`${path}.errorKind`, "string")) &&
    (isString(value["reason"]) || invalid(`${path}.reason`, "string")) &&
    (isBoolean(value["retryable"]) ||
      invalid(`${path}.retryable`, "boolean")) &&
    (isBoolean(value["isRedacted"]) || invalid(`${path}.isRedacted`, "boolean"))
  )
}

function isActionEnablement(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isBoolean(value["canStart"]) || invalid(`${path}.canStart`, "boolean")) &&
    (isOptionalString(value["startBlockedReason"]) ||
      invalid(`${path}.startBlockedReason`, "string or undefined")) &&
    (isBoolean(value["canPause"]) || invalid(`${path}.canPause`, "boolean")) &&
    (isOptionalString(value["pauseBlockedReason"]) ||
      invalid(`${path}.pauseBlockedReason`, "string or undefined")) &&
    (isBoolean(value["canResume"]) ||
      invalid(`${path}.canResume`, "boolean")) &&
    (isOptionalString(value["resumeBlockedReason"]) ||
      invalid(`${path}.resumeBlockedReason`, "string or undefined")) &&
    (isBoolean(value["canRetry"]) || invalid(`${path}.canRetry`, "boolean")) &&
    (isOptionalString(value["retryBlockedReason"]) ||
      invalid(`${path}.retryBlockedReason`, "string or undefined")) &&
    (isBoolean(value["canStartNextPhase"]) ||
      invalid(`${path}.canStartNextPhase`, "boolean")) &&
    (isOptionalString(value["nextPhaseBlockedReason"]) ||
      invalid(`${path}.nextPhaseBlockedReason`, "string or undefined"))
  )
}

function isTermTranslationPhaseSummaryResponseDto(
  value: unknown
): value is GetTermTranslationPhaseSummaryResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid("$.jobId", "number")) &&
    (isString(value["currentPhase"]) || invalid("$.currentPhase", "string")) &&
    (isString(value["phaseState"]) || invalid("$.phaseState", "string")) &&
    (isOptionalNumber(value["phaseRunId"]) ||
      invalid("$.phaseRunId", "number or undefined")) &&
    (isOptionalString(value["startedAt"]) ||
      invalid("$.startedAt", "string or undefined")) &&
    (isOptionalString(value["finishedAt"]) ||
      invalid("$.finishedAt", "string or undefined")) &&
    isProgress(value["progress"], "$.progress") &&
    (isNumber(value["totalTermCount"]) ||
      invalid("$.totalTermCount", "number")) &&
    (isNumber(value["dictionaryHitCount"]) ||
      invalid("$.dictionaryHitCount", "number")) &&
    (isNumber(value["aiTargetCount"]) ||
      invalid("$.aiTargetCount", "number")) &&
    isExecution(value["execution"], "$.execution") &&
    isResultSummary(value["resultSummary"], "$.resultSummary") &&
    isErrorSummary(value["errorSummary"], "$.errorSummary") &&
    isActionEnablement(value["actionEnablement"], "$.actionEnablement")
  )
}

function isTermTranslationPhaseCommandResponseDto(
  value: unknown
): value is StartTermTranslationPhaseResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid("$.jobId", "number")) &&
    (isString(value["currentPhase"]) || invalid("$.currentPhase", "string")) &&
    (isString(value["phaseState"]) || invalid("$.phaseState", "string")) &&
    (isOptionalNumber(value["phaseRunId"]) ||
      invalid("$.phaseRunId", "number or undefined")) &&
    (isOptionalString(value["startedAt"]) ||
      invalid("$.startedAt", "string or undefined")) &&
    (isOptionalString(value["finishedAt"]) ||
      invalid("$.finishedAt", "string or undefined")) &&
    isProgress(value["progress"], "$.progress") &&
    (isBoolean(value["retryable"]) || invalid("$.retryable", "boolean")) &&
    (isBoolean(value["canStartNextPhase"]) ||
      invalid("$.canStartNextPhase", "boolean")) &&
    isErrorSummary(value["errorSummary"], "$.errorSummary")
  )
}

function isTermTranslationAISettingsResponseDto(
  value: unknown
): value is SaveTermTranslationPhaseAISettingsResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid("$.jobId", "number")) &&
    (isString(value["phaseId"]) || invalid("$.phaseId", "string")) &&
    (isString(value["provider"]) || invalid("$.provider", "string")) &&
    (isString(value["model"]) || invalid("$.model", "string")) &&
    (isString(value["executionMode"]) ||
      invalid("$.executionMode", "string")) &&
    (isString(value["batchMode"]) || invalid("$.batchMode", "string")) &&
    (value["credentialStatus"] === "configured" ||
      value["credentialStatus"] === "missing" ||
      value["credentialStatus"] === "not_required" ||
      invalid("$.credentialStatus", "known credential status")) &&
    (value["modelListStatus"] === "not_updated" ||
      value["modelListStatus"] === "loading" ||
      value["modelListStatus"] === "success" ||
      value["modelListStatus"] === "failed" ||
      value["modelListStatus"] === "credential_missing" ||
      value["modelListStatus"] === "credential_not_required" ||
      invalid("$.modelListStatus", "known model list status"))
  )
}

function isTermTranslationNextPhaseReadinessResponseDto(
  value: unknown
): value is GetTermTranslationNextPhaseReadinessResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid("$.jobId", "number")) &&
    (isString(value["currentPhase"]) || invalid("$.currentPhase", "string")) &&
    (isString(value["phaseState"]) || invalid("$.phaseState", "string")) &&
    (isBoolean(value["canStartNextPhase"]) ||
      invalid("$.canStartNextPhase", "boolean")) &&
    (isOptionalString(value["blockedReason"]) ||
      invalid("$.blockedReason", "string or undefined")) &&
    (isOptionalString(value["errorKind"]) ||
      invalid("$.errorKind", "string or undefined"))
  )
}

class TermTranslationPhaseGateway implements TermTranslationPhaseGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getProcessingTargetList(
    request: GetProcessingTargetListRequestDto
  ): Promise<GetProcessingTargetListResponseDto> {
    return this.invokeBinding(
      GetProcessingTargetList,
      "GetProcessingTargetList",
      request,
      isProcessingTargetListResponseDto
    )
  }

  getTermTranslationPhaseSummary(
    request: GetTermTranslationPhaseSummaryRequestDto
  ): Promise<GetTermTranslationPhaseSummaryResponseDto> {
    return this.invokeBinding(
      GetTermTranslationPhaseSummary,
      "GetTermTranslationPhaseSummary",
      request,
      isTermTranslationPhaseSummaryResponseDto
    )
  }

  startTermTranslationPhase(
    request: StartTermTranslationPhaseRequestDto
  ): Promise<StartTermTranslationPhaseResponseDto> {
    return this.invokeBinding(
      StartTermTranslationPhase,
      "StartTermTranslationPhase",
      request,
      isTermTranslationPhaseCommandResponseDto
    )
  }

  saveTermTranslationPhaseAISettings(
    request: SaveTermTranslationPhaseAISettingsRequestDto
  ): Promise<SaveTermTranslationPhaseAISettingsResponseDto> {
    return this.invokeBinding(
      SaveTermTranslationPhaseAISettings,
      "SaveTermTranslationPhaseAISettings",
      request,
      isTermTranslationAISettingsResponseDto
    )
  }

  pauseTermTranslationPhase(
    request: PauseTermTranslationPhaseRequestDto
  ): Promise<PauseTermTranslationPhaseResponseDto> {
    return this.invokeBinding(
      PauseTermTranslationPhase,
      "PauseTermTranslationPhase",
      request,
      isTermTranslationPhaseCommandResponseDto
    )
  }

  resumeTermTranslationPhase(
    request: ResumeTermTranslationPhaseRequestDto
  ): Promise<ResumeTermTranslationPhaseResponseDto> {
    return this.invokeBinding(
      ResumeTermTranslationPhase,
      "ResumeTermTranslationPhase",
      request,
      isTermTranslationPhaseCommandResponseDto
    )
  }

  retryTermTranslationPhase(
    request: RetryTermTranslationPhaseRequestDto
  ): Promise<RetryTermTranslationPhaseResponseDto> {
    return this.invokeBinding(
      RetryTermTranslationPhase,
      "RetryTermTranslationPhase",
      request,
      isTermTranslationPhaseCommandResponseDto
    )
  }

  getTermTranslationNextPhaseReadiness(
    request: GetTermTranslationNextPhaseReadinessRequestDto
  ): Promise<GetTermTranslationNextPhaseReadinessResponseDto> {
    return this.invokeBinding(
      GetTermTranslationNextPhaseReadiness,
      "GetTermTranslationNextPhaseReadiness",
      request,
      isTermTranslationNextPhaseReadinessResponseDto
    )
  }
}

export function createTermTranslationPhaseGateway(): TermTranslationPhaseGatewayContract {
  return new TermTranslationPhaseGateway(createBindingInvoker())
}
