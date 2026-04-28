import type { TranslationJobSetupGatewayContract } from "@application/gateway-contract/translation-job-setup"
import type {
  CreateTranslationJobRequestDto,
  CreateTranslationJobResponseDto,
  GetTranslationJobSetupOptionsResponseDto,
  GetTranslationJobSetupSummaryRequestDto,
  GetTranslationJobSetupSummaryResponseDto,
  ValidateTranslationJobSetupRequestDto,
  ValidateTranslationJobSetupResponseDto
} from "@controller/wails/gateway-dto/translation-job-setup"

type TranslationJobSetupBindingName =
  | "GetTranslationJobSetupOptions"
  | "ValidateTranslationJobSetup"
  | "CreateTranslationJob"
  | "GetTranslationJobSetupSummary"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: TranslationJobSetupBindingName,
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
  bindingName: TranslationJobSetupBindingName
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
    toRecord(wailsRecord["TranslationJobSetupController"]),
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
      Promise.resolve((binding as (...invokeArgs: [] | [unknown]) => unknown)(...args))
  }

  return null
}

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    bindingName: TranslationJobSetupBindingName,
    request?: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. backend-job-setup-contract-freeze 完了後に接続します。`
        )
      )
    }

    if (request === undefined) {
      return binding().then((response) => response as ResponseDto)
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class TranslationJobSetupGateway implements TranslationJobSetupGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getTranslationJobSetupOptions(): Promise<GetTranslationJobSetupOptionsResponseDto> {
    return this.invokeBinding("GetTranslationJobSetupOptions")
  }

  validateTranslationJobSetup(
    request: ValidateTranslationJobSetupRequestDto
  ): Promise<ValidateTranslationJobSetupResponseDto> {
    return this.invokeBinding("ValidateTranslationJobSetup", request)
  }

  createTranslationJob(
    request: CreateTranslationJobRequestDto
  ): Promise<CreateTranslationJobResponseDto> {
    return this.invokeBinding("CreateTranslationJob", request)
  }

  getTranslationJobSetupSummary(
    request: GetTranslationJobSetupSummaryRequestDto
  ): Promise<GetTranslationJobSetupSummaryResponseDto> {
    return this.invokeBinding("GetTranslationJobSetupSummary", request)
  }
}

export function createTranslationJobSetupGateway(): TranslationJobSetupGatewayContract {
  return new TranslationJobSetupGateway(createBindingInvoker())
}