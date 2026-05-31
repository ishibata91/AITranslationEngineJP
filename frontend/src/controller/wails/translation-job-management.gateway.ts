import type { TranslationJobManagementGatewayContract } from "@application/gateway-contract/translation-job-management"
import {
  DeleteJob,
  GetJobDetail,
  ListIncompleteJobs,
  RequestStop,
  ResumeJob
} from "../../../wailsjs/go/wails/AppController.js"
import type {
  TranslationJobManagementActionRequestDto,
  TranslationJobManagementActionResponseDto,
  TranslationJobManagementDeleteRequestDto,
  TranslationJobManagementGetDetailRequestDto,
  TranslationJobManagementJobDetailDto,
  TranslationJobManagementListResponseDto
} from "@controller/wails/gateway-dto/translation-job-management"
import {
  createGatewayResponseShapeError,
  isArrayOf,
  isBoolean,
  isNumber,
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
  binding: WailsTranslationJobManagementBinding<RequestDto>,
  bindingName: string,
  request: RequestDto | undefined,
  isResponseDto: RuntimeShapeValidator<ResponseDto>
) => Promise<ResponseDto>

type WailsTranslationJobManagementBinding<RequestDto> = (
  request: RequestDto
) => Promise<unknown>

type WailsTranslationJobManagementNoArgBinding = () => Promise<unknown>

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    binding: WailsTranslationJobManagementBinding<RequestDto>,
    bindingName: string,
    request: RequestDto | undefined,
    isResponseDto: RuntimeShapeValidator<ResponseDto>
  ): Promise<ResponseDto> => {
    const responsePromise =
      request === undefined
        ? (binding as unknown as WailsTranslationJobManagementNoArgBinding)()
        : binding(request)

    return responsePromise.then((response) => {
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

function isJobState(value: unknown): boolean {
  return (
    value === "Ready" ||
    value === "Running" ||
    value === "Paused" ||
    value === "RecoverableFailed" ||
    value === "Failed" ||
    value === "Canceled" ||
    value === "Completed"
  )
}

function isStateTone(value: unknown): boolean {
  return value === "neutral" || value === "info" || value === "warning" || value === "danger"
}

function isActionTone(value: unknown): boolean {
  return (
    value === "info" ||
    value === "success" ||
    value === "warning" ||
    value === "danger"
  )
}

function isCredentialState(value: unknown): boolean {
  return (
    value === "configured" || value === "missing" || value === "inaccessible"
  )
}

function isCurrentPhase(value: unknown): boolean {
  return (
    value === "term_translation" ||
    value === "persona_generation" ||
    value === "body_translation"
  )
}

function isReasonCategory(value: unknown): boolean {
  return (
    value === "cache_missing" ||
    value === "terminal_state" ||
    value === "state_projection_inconsistent" ||
    value === "runtime_snapshot_missing" ||
    value === "phase_progress_aggregation_failed" ||
    value === "stale_selection" ||
    value === "list_load_failure" ||
    value === "running_delete_blocked" ||
    value === "stop_requested" ||
    value === "stop_failed" ||
    value === "delete_failed" ||
    value === "resume_failed"
  )
}

function isOperationKind(value: unknown): boolean {
  return value === "stop" || value === "resume" || value === "delete"
}

function isBlockedReason(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isReasonCategory(value["category"]) ||
      invalid(`${path}.category`, "known reason category")) &&
    (isString(value["title"]) || invalid(`${path}.title`, "string")) &&
    (isString(value["detail"]) || invalid(`${path}.detail`, "string"))
  )
}

function isInputSource(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isNumber(value["inputSourceId"]) ||
      invalid(`${path}.inputSourceId`, "number")) &&
    (isString(value["inputSourceLabel"]) ||
      invalid(`${path}.inputSourceLabel`, "string")) &&
    (isString(value["inputSourceKindLabel"]) ||
      invalid(`${path}.inputSourceKindLabel`, "string")) &&
    (isString(value["sourcePath"]) ||
      invalid(`${path}.sourcePath`, "string")) &&
    (isString(value["pluginName"]) ||
      invalid(`${path}.pluginName`, "string")) &&
    (isString(value["extractedJsonLabel"]) ||
      invalid(`${path}.extractedJsonLabel`, "string"))
  )
}

function isProgress(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isCurrentPhase(value["currentPhase"]) ||
      invalid(`${path}.currentPhase`, "known current phase")) &&
    (isString(value["currentPhaseLabel"]) ||
      invalid(`${path}.currentPhaseLabel`, "string")) &&
    (isNumber(value["percent"]) || invalid(`${path}.percent`, "number")) &&
    (isString(value["progressLabel"]) ||
      invalid(`${path}.progressLabel`, "string")) &&
    (isString(value["lastUpdatedLabel"]) ||
      invalid(`${path}.lastUpdatedLabel`, "string"))
  )
}

function isRuntimeSummary(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["providerLabel"]) ||
      invalid(`${path}.providerLabel`, "string")) &&
    (isString(value["modelLabel"]) ||
      invalid(`${path}.modelLabel`, "string")) &&
    (isString(value["executionModeLabel"]) ||
      invalid(`${path}.executionModeLabel`, "string")) &&
    (isCredentialState(value["credentialState"]) ||
      invalid(`${path}.credentialState`, "known credential state")) &&
    (isString(value["credentialStateLabel"]) ||
      invalid(`${path}.credentialStateLabel`, "string"))
  )
}

function isAvailability(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isOperationKind(value["kind"]) ||
      invalid(`${path}.kind`, "known operation kind")) &&
    (isBoolean(value["enabled"]) ||
      invalid(`${path}.enabled`, "boolean")) &&
    (isString(value["label"]) || invalid(`${path}.label`, "string")) &&
    (isString(value["helperText"]) ||
      invalid(`${path}.helperText`, "string")) &&
    (value["reasonCategory"] === undefined ||
      isReasonCategory(value["reasonCategory"]) ||
      invalid(`${path}.reasonCategory`, "known reason category")) &&
    (isOptionalString(value["reasonText"]) ||
      invalid(`${path}.reasonText`, "string or undefined"))
  )
}

function isJobSummary(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isNumber(value["jobId"]) || invalid(`${path}.jobId`, "number")) &&
    (isJobState(value["jobState"]) ||
      invalid(`${path}.jobState`, "known job state")) &&
    (isString(value["jobStateLabel"]) ||
      invalid(`${path}.jobStateLabel`, "string")) &&
    (isStateTone(value["stateTone"]) ||
      invalid(`${path}.stateTone`, "known state tone")) &&
    (value["canOpenPhase"] === undefined ||
      isBoolean(value["canOpenPhase"]) ||
      invalid(`${path}.canOpenPhase`, "boolean or undefined")) &&
    (value["openBlockedReason"] === undefined ||
      isBlockedReason(value["openBlockedReason"], `${path}.openBlockedReason`)) &&
    isInputSource(value["inputSource"], `${path}.inputSource`) &&
    isProgress(value["progress"], `${path}.progress`) &&
    isAvailability(value["stopAvailability"], `${path}.stopAvailability`) &&
    isAvailability(value["resumeAvailability"], `${path}.resumeAvailability`) &&
    isAvailability(value["deleteAvailability"], `${path}.deleteAvailability`)
  )
}

function isJobDetail(value: unknown): value is TranslationJobManagementJobDetailDto {
  resetIssues()
  if (!isJobSummary(value, "$")) {
    return false
  }
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (value["cacheState"] === "available" ||
      value["cacheState"] === "missing" ||
      invalid("$.cacheState", "known cache state")) &&
    (isString(value["cacheStateLabel"]) ||
      invalid("$.cacheStateLabel", "string")) &&
    isRuntimeSummary(value["runtimeSummary"], "$.runtimeSummary") &&
    (isArrayOf(value["resumeBlockedReasons"], (item) =>
      isBlockedReason(item, "$.resumeBlockedReasons[]")
    ) ||
      invalid("$.resumeBlockedReasons", "blocked reason array")) &&
    (isArrayOf(value["warnings"], (item) =>
      isBlockedReason(item, "$.warnings[]")
    ) ||
      invalid("$.warnings", "blocked reason array")) &&
    (isArrayOf(value["deleteImpactLines"], isString) ||
      invalid("$.deleteImpactLines", "string array"))
  )
}

function isListResponseDto(
  value: unknown
): value is TranslationJobManagementListResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    isArrayOf(value["jobs"], (item) => isJobSummary(item, "$.jobs[]")) ||
    invalid("$.jobs", "job summary array")
  )
}

function isActionResponseDto(
  value: unknown
): value is TranslationJobManagementActionResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isString(value["message"]) || invalid("$.message", "string")) &&
    (isActionTone(value["tone"]) || invalid("$.tone", "known action tone")) &&
    (value["detail"] === undefined || isJobDetail(value["detail"])) &&
    (value["deletedJobId"] === undefined ||
      isNumber(value["deletedJobId"]) ||
      invalid("$.deletedJobId", "number or undefined")) &&
    (value["reasonCategory"] === undefined ||
      isReasonCategory(value["reasonCategory"]) ||
      invalid("$.reasonCategory", "known reason category"))
  )
}

class TranslationJobManagementGateway implements TranslationJobManagementGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  ListIncompleteJobs(): Promise<TranslationJobManagementListResponseDto> {
    return this.invokeBinding(
      ListIncompleteJobs,
      "ListIncompleteJobs",
      undefined,
      isListResponseDto
    )
  }

  GetJobDetail(
    request: TranslationJobManagementGetDetailRequestDto
  ): Promise<TranslationJobManagementJobDetailDto> {
    return this.invokeBinding(GetJobDetail, "GetJobDetail", request, isJobDetail)
  }

  RequestStop(
    request: TranslationJobManagementActionRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding(
      RequestStop,
      "RequestStop",
      request,
      isActionResponseDto
    )
  }

  ResumeJob(
    request: TranslationJobManagementActionRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding(ResumeJob, "ResumeJob", request, isActionResponseDto)
  }

  DeleteJob(
    request: TranslationJobManagementDeleteRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding(DeleteJob, "DeleteJob", request, isActionResponseDto)
  }
}

export function createTranslationJobManagementGateway(): TranslationJobManagementGatewayContract {
  return new TranslationJobManagementGateway(createBindingInvoker())
}
