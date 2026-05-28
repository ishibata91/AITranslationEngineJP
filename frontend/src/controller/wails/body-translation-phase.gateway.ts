import type { BodyTranslationPhaseGatewayContract } from "@application/gateway-contract/body-translation-phase"
import {
  CancelBodyTranslationPhase,
  GetBodyTranslationOutputReadiness,
  GetBodyTranslationPhaseSummary,
  GetProcessingTargetList,
  PauseBodyTranslationPhase,
  ResumeBodyTranslationPhase,
  RetryBodyTranslationPhase,
  SaveBodyTranslationPhaseAISettings,
  StartBodyTranslationPhase
} from "../../../wailsjs/go/wails/AppController.js"
import type {
  CancelBodyTranslationPhaseRequestDto,
  CancelBodyTranslationPhaseResponseDto,
  GetProcessingTargetListRequestDto,
  GetProcessingTargetListResponseDto,
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
  SaveBodyTranslationPhaseAISettingsRequestDto,
  SaveBodyTranslationPhaseAISettingsResponseDto,
  StartBodyTranslationPhaseRequestDto,
  StartBodyTranslationPhaseResponseDto
} from "@controller/wails/gateway-dto/body-translation-phase"
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
  binding: WailsBodyTranslationBinding<RequestDto>,
  bindingName: string,
  request: RequestDto,
  isResponseDto: RuntimeShapeValidator<ResponseDto>
) => Promise<ResponseDto>

type WailsBodyTranslationBinding<RequestDto> = (
  request: RequestDto
) => Promise<unknown>

function toRecord(value: unknown): Record<string, unknown> | null {
  if (typeof value !== "object" || value === null) {
    return null
  }

  return value as Record<string, unknown>
}

function appControllerBinding<RequestDto>(
  bindingName: string,
  fallback: WailsBodyTranslationBinding<RequestDto>
): WailsBodyTranslationBinding<RequestDto> {
  return (request: RequestDto): Promise<unknown> => {
    const globalRecord = toRecord(globalThis)
    const goRecord = toRecord(globalRecord?.["go"])
    const wailsRecord = toRecord(goRecord?.["wails"])
    const controller = toRecord(wailsRecord?.["AppController"])
    const binding = controller?.[bindingName]
    if (typeof binding === "function") {
      return Promise.resolve((binding as (arg: unknown) => unknown)(request))
    }
    return fallback(request)
  }
}

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    binding: WailsBodyTranslationBinding<RequestDto>,
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

function isOptionalNumberOrString(
  value: unknown
): value is number | string | undefined {
  return value === undefined || isNumber(value) || isString(value)
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

function isBodyPhaseErrorKind(value: unknown): boolean {
  return (
    value === "persona_phase_incomplete" ||
    value === "terminal_job" ||
    value === "active_phase_exists" ||
    value === "input_snapshot_failed" ||
    value === "provider_failure" ||
    value === "invalid_provider_response" ||
    value === "protection_validation_failed" ||
    value === "save_failed" ||
    value === "output_readiness_blocked" ||
    value === "late_response_rejected" ||
    value === "secret_redacted"
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
    (isNumber(value["targetCount"]) ||
      invalid(`${path}.targetCount`, "number")) &&
    (isNumber(value["translatedCount"]) ||
      invalid(`${path}.translatedCount`, "number")) &&
    (isNumber(value["skippedCount"]) ||
      invalid(`${path}.skippedCount`, "number")) &&
    (isString(value["currentStep"]) || invalid(`${path}.currentStep`, "string"))
  )
}

function isInputSummary(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isNumber(value["targetCount"]) ||
      invalid(`${path}.targetCount`, "number")) &&
    (value["skippedReasons"] === undefined ||
      isArrayOf(value["skippedReasons"], isString) ||
      invalid(`${path}.skippedReasons`, "string array or undefined")) &&
    (isOptionalString(value["inputSnapshotRef"]) ||
      invalid(`${path}.inputSnapshotRef`, "string or undefined")) &&
    (isString(value["dictionaryDigest"]) ||
      invalid(`${path}.dictionaryDigest`, "string")) &&
    (isString(value["personaDigest"]) ||
      invalid(`${path}.personaDigest`, "string")) &&
    (isString(value["metadataDigest"]) ||
      invalid(`${path}.metadataDigest`, "string")) &&
    (isString(value["promptDigest"]) ||
      invalid(`${path}.promptDigest`, "string"))
  )
}

function isRequestSummary(value: unknown, path: string): boolean {
  if (value === undefined) {
    return true
  }
  if (!isRecord(value)) {
    return invalid(path, "object or undefined")
  }

  return (
    (isOptionalNumber(value["providerTargetCount"]) ||
      invalid(`${path}.providerTargetCount`, "number or undefined")) &&
    (isOptionalNumber(value["exactDictionaryExclusionCount"]) ||
      invalid(`${path}.exactDictionaryExclusionCount`, "number or undefined")) &&
    (isOptionalNumber(value["partialDictionaryConstraintCount"]) ||
      invalid(`${path}.partialDictionaryConstraintCount`, "number or undefined"))
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
    (isNumber(value["requestUnitCount"]) ||
      invalid(`${path}.requestUnitCount`, "number")) &&
    (isNumber(value["outputCount"]) ||
      invalid(`${path}.outputCount`, "number"))
  )
}

function isFieldIdentity(value: unknown, path: string): boolean {
  if (value === undefined) {
    return true
  }
  if (!isRecord(value)) {
    return invalid(path, "object or undefined")
  }

  return (
    (isOptionalNumberOrString(value["translationFieldId"]) ||
      invalid(`${path}.translationFieldId`, "number, string, or undefined")) &&
    (isOptionalNumberOrString(value["phaseTranslationFieldId"]) ||
      invalid(`${path}.phaseTranslationFieldId`, "number, string, or undefined")) &&
    (isOptionalString(value["recordType"]) ||
      invalid(`${path}.recordType`, "string or undefined")) &&
    (isOptionalString(value["fieldType"]) ||
      invalid(`${path}.fieldType`, "string or undefined")) &&
    (isOptionalString(value["formId"]) ||
      invalid(`${path}.formId`, "string or undefined")) &&
    (isOptionalString(value["editorId"]) ||
      invalid(`${path}.editorId`, "string or undefined")) &&
    (isOptionalString(value["fieldLabel"]) ||
      invalid(`${path}.fieldLabel`, "string or undefined"))
  )
}

function isFieldResultItem(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    isFieldIdentity(value["identity"], `${path}.identity`) &&
    (isOptionalNumberOrString(value["fieldId"]) ||
      invalid(`${path}.fieldId`, "number, string, or undefined")) &&
    (isOptionalString(value["fieldLabel"]) ||
      invalid(`${path}.fieldLabel`, "string or undefined")) &&
    (isOptionalString(value["sourceExcerpt"]) ||
      invalid(`${path}.sourceExcerpt`, "string or undefined")) &&
    (isOptionalString(value["translatedText"]) ||
      invalid(`${path}.translatedText`, "string or undefined")) &&
    (isOptionalString(value["outputStatus"]) ||
      invalid(`${path}.outputStatus`, "string or undefined")) &&
    (isOptionalString(value["protectionValidationResult"]) ||
      invalid(`${path}.protectionValidationResult`, "string or undefined")) &&
    (isOptionalString(value["protectionValidationSummary"]) ||
      invalid(`${path}.protectionValidationSummary`, "string or undefined")) &&
    (isOptionalNumber(value["retryCount"]) ||
      invalid(`${path}.retryCount`, "number or undefined"))
  )
}

function isFieldResultSummary(value: unknown, path: string): boolean {
  if (value === undefined) {
    return true
  }
  if (!isRecord(value)) {
    return invalid(path, "object or undefined")
  }

  return (
    (isNumber(value["translatedCount"]) ||
      invalid(`${path}.translatedCount`, "number")) &&
    (isNumber(value["failedCount"]) ||
      invalid(`${path}.failedCount`, "number")) &&
    (isNumber(value["skippedCount"]) ||
      invalid(`${path}.skippedCount`, "number")) &&
    (isNumber(value["protectionFailedCount"]) ||
      invalid(`${path}.protectionFailedCount`, "number")) &&
    (isNumber(value["outputReadyCount"]) ||
      invalid(`${path}.outputReadyCount`, "number")) &&
    (isOptionalNumber(value["outputCount"]) ||
      invalid(`${path}.outputCount`, "number or undefined")) &&
    (value["fieldResults"] === undefined ||
      isArrayOf(value["fieldResults"], (item) =>
        isFieldResultItem(item, `${path}.fieldResults[]`)
      ) ||
      invalid(`${path}.fieldResults`, "field result array or undefined"))
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
    (isBodyPhaseErrorKind(value["errorKind"]) ||
      invalid(`${path}.errorKind`, "known body phase error kind")) &&
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
    (isBoolean(value["canCancel"]) || invalid(`${path}.canCancel`, "boolean")) &&
    (isOptionalString(value["cancelBlockedReason"]) ||
      invalid(`${path}.cancelBlockedReason`, "string or undefined")) &&
    (isBoolean(value["canCheckOutputReadiness"]) ||
      invalid(`${path}.canCheckOutputReadiness`, "boolean")) &&
    (isOptionalString(value["outputReadinessBlockedReason"]) ||
      invalid(`${path}.outputReadinessBlockedReason`, "string or undefined"))
  )
}

function isOutputReadinessSummary(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isBoolean(value["ready"]) || invalid(`${path}.ready`, "boolean")) &&
    (isOptionalString(value["blockedReason"]) ||
      invalid(`${path}.blockedReason`, "string or undefined")) &&
    (value["errorKind"] === undefined ||
      isBodyPhaseErrorKind(value["errorKind"]) ||
      invalid(`${path}.errorKind`, "known body phase error kind")) &&
    (isNumber(value["completedFieldCount"]) ||
      invalid(`${path}.completedFieldCount`, "number")) &&
    (isBoolean(value["statusConsistent"]) ||
      invalid(`${path}.statusConsistent`, "boolean"))
  )
}

function hasBodyPhaseBase(value: unknown): value is Record<string, unknown> {
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
    isInputSummary(value["inputSummary"], "$.inputSummary") &&
    isRequestSummary(value["requestSummary"], "$.requestSummary") &&
    isExecution(value["execution"], "$.execution") &&
    (value["fieldResults"] === undefined ||
      isArrayOf(value["fieldResults"], (item) =>
        isFieldResultItem(item, "$.fieldResults[]")
      ) ||
      invalid("$.fieldResults", "field result array or undefined")) &&
    isFieldResultSummary(value["resultSummary"], "$.resultSummary") &&
    isOutputReadinessSummary(value["outputReadiness"], "$.outputReadiness")
  )
}

function isBodyTranslationPhaseSummaryResponseDto(
  value: unknown
): value is GetBodyTranslationPhaseSummaryResponseDto {
  resetIssues()
  return (
    hasBodyPhaseBase(value) &&
    isErrorSummary(value["errorSummary"], "$.errorSummary") &&
    isActionEnablement(value["actionEnablement"], "$.actionEnablement")
  )
}

function isBodyTranslationPhaseCommandResponseDto(
  value: unknown
): value is StartBodyTranslationPhaseResponseDto {
  resetIssues()
  return (
    hasBodyPhaseBase(value) &&
    (isOptionalString(value["inputSnapshotDigest"]) ||
      invalid("$.inputSnapshotDigest", "string or undefined")) &&
    (isBoolean(value["retryable"]) || invalid("$.retryable", "boolean")) &&
    isErrorSummary(value["errorSummary"], "$.errorSummary")
  )
}

function isBodyTranslationAISettingsResponseDto(
  value: unknown
): value is SaveBodyTranslationPhaseAISettingsResponseDto {
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

function isBodyTranslationOutputReadinessResponseDto(
  value: unknown
): value is GetBodyTranslationOutputReadinessResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid("$.jobId", "number")) &&
    (isString(value["currentPhase"]) || invalid("$.currentPhase", "string")) &&
    (isString(value["phaseState"]) || invalid("$.phaseState", "string")) &&
    (isBoolean(value["ready"]) || invalid("$.ready", "boolean")) &&
    (isOptionalString(value["blockedReason"]) ||
      invalid("$.blockedReason", "string or undefined")) &&
    (value["errorKind"] === undefined ||
      isBodyPhaseErrorKind(value["errorKind"]) ||
      invalid("$.errorKind", "known body phase error kind")) &&
    (isNumber(value["completedFieldCount"]) ||
      invalid("$.completedFieldCount", "number")) &&
    (isBoolean(value["statusConsistent"]) ||
      invalid("$.statusConsistent", "boolean")) &&
    (isNumber(value["outputCount"]) || invalid("$.outputCount", "number"))
  )
}

class BodyTranslationPhaseGateway implements BodyTranslationPhaseGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getProcessingTargetList(
    request: GetProcessingTargetListRequestDto
  ): Promise<GetProcessingTargetListResponseDto> {
    return this.invokeBinding(
      appControllerBinding("GetProcessingTargetList", GetProcessingTargetList),
      "GetProcessingTargetList",
      request,
      isProcessingTargetListResponseDto
    )
  }

  getBodyTranslationPhaseSummary(
    request: GetBodyTranslationPhaseSummaryRequestDto
  ): Promise<GetBodyTranslationPhaseSummaryResponseDto> {
    return this.invokeBinding(
      appControllerBinding(
        "GetBodyTranslationPhaseSummary",
        GetBodyTranslationPhaseSummary
      ),
      "GetBodyTranslationPhaseSummary",
      request,
      isBodyTranslationPhaseSummaryResponseDto
    )
  }

  startBodyTranslationPhase(
    request: StartBodyTranslationPhaseRequestDto
  ): Promise<StartBodyTranslationPhaseResponseDto> {
    return this.invokeBinding(
      appControllerBinding("StartBodyTranslationPhase", StartBodyTranslationPhase),
      "StartBodyTranslationPhase",
      request,
      isBodyTranslationPhaseCommandResponseDto
    )
  }

  saveBodyTranslationPhaseAISettings(
    request: SaveBodyTranslationPhaseAISettingsRequestDto
  ): Promise<SaveBodyTranslationPhaseAISettingsResponseDto> {
    return this.invokeBinding(
      appControllerBinding(
        "SaveBodyTranslationPhaseAISettings",
        SaveBodyTranslationPhaseAISettings
      ),
      "SaveBodyTranslationPhaseAISettings",
      request,
      isBodyTranslationAISettingsResponseDto
    )
  }

  pauseBodyTranslationPhase(
    request: PauseBodyTranslationPhaseRequestDto
  ): Promise<PauseBodyTranslationPhaseResponseDto> {
    return this.invokeBinding(
      appControllerBinding("PauseBodyTranslationPhase", PauseBodyTranslationPhase),
      "PauseBodyTranslationPhase",
      request,
      isBodyTranslationPhaseCommandResponseDto
    )
  }

  resumeBodyTranslationPhase(
    request: ResumeBodyTranslationPhaseRequestDto
  ): Promise<ResumeBodyTranslationPhaseResponseDto> {
    return this.invokeBinding(
      appControllerBinding(
        "ResumeBodyTranslationPhase",
        ResumeBodyTranslationPhase
      ),
      "ResumeBodyTranslationPhase",
      request,
      isBodyTranslationPhaseCommandResponseDto
    )
  }

  retryBodyTranslationPhase(
    request: RetryBodyTranslationPhaseRequestDto
  ): Promise<RetryBodyTranslationPhaseResponseDto> {
    return this.invokeBinding(
      appControllerBinding("RetryBodyTranslationPhase", RetryBodyTranslationPhase),
      "RetryBodyTranslationPhase",
      request,
      isBodyTranslationPhaseCommandResponseDto
    )
  }

  cancelBodyTranslationPhase(
    request: CancelBodyTranslationPhaseRequestDto
  ): Promise<CancelBodyTranslationPhaseResponseDto> {
    return this.invokeBinding(
      appControllerBinding(
        "CancelBodyTranslationPhase",
        CancelBodyTranslationPhase
      ),
      "CancelBodyTranslationPhase",
      request,
      isBodyTranslationPhaseCommandResponseDto
    )
  }

  getBodyTranslationOutputReadiness(
    request: GetBodyTranslationOutputReadinessRequestDto
  ): Promise<GetBodyTranslationOutputReadinessResponseDto> {
    return this.invokeBinding(
      appControllerBinding(
        "GetBodyTranslationOutputReadiness",
        GetBodyTranslationOutputReadiness
      ),
      "GetBodyTranslationOutputReadiness",
      request,
      isBodyTranslationOutputReadinessResponseDto
    )
  }
}

export function createBodyTranslationPhaseGateway(): BodyTranslationPhaseGatewayContract {
  return new BodyTranslationPhaseGateway(createBindingInvoker())
}
