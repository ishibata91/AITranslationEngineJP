export type ProviderSettingsProviderId = "gemini" | "lm_studio" | "xai"

export type ProviderSettingsCredentialState =
  | "missing"
  | "configured"
  | "not_required"

export type ProviderSettingsValidationState =
  | "not_validated"
  | "pending"
  | "validated"
  | "failed"

export type ProviderSettingsSavedState =
  | "not_saved"
  | "partial"
  | "configured"

export type ProviderSettingsErrorKind =
  | "credential_missing"
  | "endpoint_missing"
  | "validation_failed"
  | "validation_stale"
  | "provider_unreachable"
  | "model_list_failed"
  | "settings_not_saved"

export interface ProviderSettingsRoute {
  routeId: string
  label: string
  currentRouteState: string
  dashboardEntryId: string
}

export interface ProviderSettingsSummary {
  providerId: ProviderSettingsProviderId
  label: string
  endpoint?: string
  credentialState: ProviderSettingsCredentialState
  credentialReferenceId?: string
  validationState: ProviderSettingsValidationState
  savedState: ProviderSettingsSavedState
  requestToken?: string
  lastFailureKind?: ProviderSettingsErrorKind
}

export type ListProviderSettingsRequest = Record<string, never>

export interface ListProviderSettingsResponse {
  route: ProviderSettingsRoute
  providers: ProviderSettingsSummary[]
}

export interface SaveProviderSettingsRequest {
  providerId: ProviderSettingsProviderId
  endpoint?: string
  apiKeyInputPresent: boolean
  credentialInput?: string
}

export interface SaveProviderSettingsResponse {
  provider: ProviderSettingsSummary
}

export interface ResetProviderSettingsRequest {
  providerId: ProviderSettingsProviderId
}

export interface ResetProviderSettingsResponse {
  provider: ProviderSettingsSummary
}

export interface ValidateProviderSettingsRequest {
  providerId: ProviderSettingsProviderId
  endpoint?: string
  credentialState: ProviderSettingsCredentialState
  credentialReferenceId?: string
  requestToken: string
}

export interface ValidateProviderSettingsResponse {
  providerId: ProviderSettingsProviderId
  validationState: ProviderSettingsValidationState
  requestToken: string
  failureKind?: ProviderSettingsErrorKind
}

export interface ProviderSettingsGatewayContract {
  ListProviderSettings(
    request?: ListProviderSettingsRequest
  ): Promise<ListProviderSettingsResponse>
  SaveProviderSettings(
    request: SaveProviderSettingsRequest
  ): Promise<SaveProviderSettingsResponse>
  ResetProviderSettings(
    request: ResetProviderSettingsRequest
  ): Promise<ResetProviderSettingsResponse>
  ValidateProviderSettings(
    request: ValidateProviderSettingsRequest
  ): Promise<ValidateProviderSettingsResponse>
}
