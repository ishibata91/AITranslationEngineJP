import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
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

type ProviderSettingsBindingName =
  | "ListProviderSettings"
  | "SaveProviderSettings"
  | "ResetProviderSettings"
  | "ValidateProviderSettings"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: ProviderSettingsBindingName,
  request?: RequestDto
) => Promise<ResponseDto>

type BindingFunction = (...args: [] | [unknown]) => Promise<unknown>

function toRecord(value: unknown): Record<string, unknown> | null {
  if (typeof value !== "object" || value === null) {
    return null
  }

  return value as Record<string, unknown>
}

function resolveBindingFunction(
  bindingName: ProviderSettingsBindingName
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
    toRecord(wailsRecord["ProviderSettingsController"]),
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

    return (...args: [] | [unknown]) =>
      Promise.resolve(
        (binding as (...invokeArgs: [] | [unknown]) => unknown)(...args)
      )
  }

  return null
}

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    bindingName: ProviderSettingsBindingName,
    request?: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-provider-settings-wails-gateway 完了後に接続します。`
        )
      )
    }

    if (request === undefined) {
      return binding().then((response) => response as ResponseDto)
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class ProviderSettingsGateway implements ProviderSettingsGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  ListProviderSettings(
    request: ListProviderSettingsRequestDto = {}
  ): Promise<ListProviderSettingsResponseDto> {
    return this.invokeBinding("ListProviderSettings", request)
  }

  SaveProviderSettings(
    request: SaveProviderSettingsRequestDto
  ): Promise<SaveProviderSettingsResponseDto> {
    return this.invokeBinding("SaveProviderSettings", request)
  }

  ResetProviderSettings(
    request: ResetProviderSettingsRequestDto
  ): Promise<ResetProviderSettingsResponseDto> {
    return this.invokeBinding("ResetProviderSettings", request)
  }

  ValidateProviderSettings(
    request: ValidateProviderSettingsRequestDto
  ): Promise<ValidateProviderSettingsResponseDto> {
    return this.invokeBinding("ValidateProviderSettings", request)
  }
}

export function createProviderSettingsGateway(): ProviderSettingsGatewayContract {
  return new ProviderSettingsGateway(createBindingInvoker())
}
