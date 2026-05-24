import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import {
  ListProviderSettings,
  ResetProviderSettings,
  SaveProviderSettings,
  ValidateProviderSettings
} from "../../../wailsjs/go/wails/AppController.js"
import type {
  ListProviderSettingsRequestDto,
  ListProviderSettingsResponseDto,
  ResetProviderSettingsRequestDto,
  ResetProviderSettingsResponseDto,
  SaveProviderSettingsRequestDto,
  SaveProviderSettingsResponseDto,
  ValidateProviderSettingsRequestDto,
  ValidateProviderSettingsResponseDto
} from "@controller/wails/gateway-dto/provider-settings"
import {
  createGatewayResponseShapeError,
  isArrayOf,
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
  binding: WailsProviderSettingsBinding<RequestDto>,
  bindingName: string,
  request: RequestDto,
  isResponseDto: RuntimeShapeValidator<ResponseDto>
) => Promise<ResponseDto>

type WailsProviderSettingsBinding<RequestDto> = (
  request: RequestDto
) => Promise<unknown>

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    binding: WailsProviderSettingsBinding<RequestDto>,
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

function isProviderId(
  value: unknown
): value is ListProviderSettingsResponseDto["providers"][number]["providerId"] {
  return value === "gemini" || value === "lm_studio" || value === "xai"
}

function isCredentialState(value: unknown): boolean {
  return (
    value === "missing" || value === "configured" || value === "not_required"
  )
}

function isValidationState(value: unknown): boolean {
  return (
    value === "not_validated" ||
    value === "pending" ||
    value === "validated" ||
    value === "failed"
  )
}

function isSavedState(value: unknown): boolean {
  return value === "not_saved" || value === "partial" || value === "configured"
}

function isErrorKind(value: unknown): boolean {
  return (
    value === "credential_missing" ||
    value === "endpoint_missing" ||
    value === "validation_failed" ||
    value === "validation_stale" ||
    value === "provider_unreachable" ||
    value === "model_list_failed" ||
    value === "settings_not_saved"
  )
}

function isProviderSettingsRoute(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isString(value["routeId"]) || invalid(`${path}.routeId`, "string")) &&
    (isString(value["label"]) || invalid(`${path}.label`, "string")) &&
    (isString(value["currentRouteState"]) ||
      invalid(`${path}.currentRouteState`, "string")) &&
    (isString(value["dashboardEntryId"]) ||
      invalid(`${path}.dashboardEntryId`, "string"))
  )
}

function isProviderSummary(value: unknown, path: string): boolean {
  if (!isRecord(value)) {
    return invalid(path, "object")
  }

  return (
    (isProviderId(value["providerId"]) ||
      invalid(`${path}.providerId`, "known provider id")) &&
    (isString(value["label"]) || invalid(`${path}.label`, "string")) &&
    (isOptionalString(value["endpoint"]) ||
      invalid(`${path}.endpoint`, "string or undefined")) &&
    (isCredentialState(value["credentialState"]) ||
      invalid(`${path}.credentialState`, "known credential state")) &&
    (isOptionalString(value["credentialReferenceId"]) ||
      invalid(`${path}.credentialReferenceId`, "string or undefined")) &&
    (isValidationState(value["validationState"]) ||
      invalid(`${path}.validationState`, "known validation state")) &&
    (isSavedState(value["savedState"]) ||
      invalid(`${path}.savedState`, "known saved state")) &&
    (isOptionalString(value["requestToken"]) ||
      invalid(`${path}.requestToken`, "string or undefined")) &&
    (value["lastFailureKind"] === undefined ||
      isErrorKind(value["lastFailureKind"]) ||
      invalid(`${path}.lastFailureKind`, "known provider settings error kind"))
  )
}

function isListProviderSettingsResponseDto(
  value: unknown
): value is ListProviderSettingsResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    isProviderSettingsRoute(value["route"], "$.route") &&
    (isArrayOf(value["providers"], (item) =>
      isProviderSummary(item, "$.providers[]")
    ) ||
      invalid("$.providers", "provider settings summary array"))
  )
}

function isProviderMutationResponseDto(
  value: unknown
): value is SaveProviderSettingsResponseDto | ResetProviderSettingsResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return isProviderSummary(value["provider"], "$.provider")
}

function isValidateProviderSettingsResponseDto(
  value: unknown
): value is ValidateProviderSettingsResponseDto {
  resetIssues()
  if (!isRecord(value)) {
    return invalid("$", "object")
  }

  return (
    (isProviderId(value["providerId"]) ||
      invalid("$.providerId", "known provider id")) &&
    (isValidationState(value["validationState"]) ||
      invalid("$.validationState", "known validation state")) &&
    (isString(value["requestToken"]) || invalid("$.requestToken", "string")) &&
    (value["failureKind"] === undefined ||
      isErrorKind(value["failureKind"]) ||
      invalid("$.failureKind", "known provider settings error kind"))
  )
}

class ProviderSettingsGateway implements ProviderSettingsGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  ListProviderSettings(
    request: ListProviderSettingsRequestDto = {}
  ): Promise<ListProviderSettingsResponseDto> {
    return this.invokeBinding(
      ListProviderSettings,
      "ListProviderSettings",
      request,
      isListProviderSettingsResponseDto
    )
  }

  SaveProviderSettings(
    request: SaveProviderSettingsRequestDto
  ): Promise<SaveProviderSettingsResponseDto> {
    return this.invokeBinding(
      SaveProviderSettings,
      "SaveProviderSettings",
      request,
      isProviderMutationResponseDto
    )
  }

  ResetProviderSettings(
    request: ResetProviderSettingsRequestDto
  ): Promise<ResetProviderSettingsResponseDto> {
    return this.invokeBinding(
      ResetProviderSettings,
      "ResetProviderSettings",
      request,
      isProviderMutationResponseDto
    )
  }

  ValidateProviderSettings(
    request: ValidateProviderSettingsRequestDto
  ): Promise<ValidateProviderSettingsResponseDto> {
    return this.invokeBinding(
      ValidateProviderSettings,
      "ValidateProviderSettings",
      request,
      isValidateProviderSettingsResponseDto
    )
  }
}

export function createProviderSettingsGateway(): ProviderSettingsGatewayContract {
  return new ProviderSettingsGateway(createBindingInvoker())
}
