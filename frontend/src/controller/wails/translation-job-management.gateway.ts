import type { TranslationJobManagementGatewayContract } from "@application/gateway-contract/translation-job-management"
import type {
  TranslationJobManagementActionRequestDto,
  TranslationJobManagementActionResponseDto,
  TranslationJobManagementDeleteRequestDto,
  TranslationJobManagementGetDetailRequestDto,
  TranslationJobManagementJobDetailDto,
  TranslationJobManagementListResponseDto
} from "@controller/wails/gateway-dto/translation-job-management"

type TranslationJobManagementBindingName =
  | "ListIncompleteJobs"
  | "GetJobDetail"
  | "RequestStop"
  | "ResumeJob"
  | "DeleteJob"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: TranslationJobManagementBindingName,
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
  bindingName: TranslationJobManagementBindingName
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
    toRecord(wailsRecord["TranslationJobManagementController"]),
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
    bindingName: TranslationJobManagementBindingName,
    request?: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-job-management-wails 完了後に接続します。`
        )
      )
    }

    if (request === undefined) {
      return binding().then((response) => response as ResponseDto)
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class TranslationJobManagementGateway
  implements TranslationJobManagementGatewayContract
{
  constructor(private readonly invokeBinding: BindingInvoker) {}

  ListIncompleteJobs(): Promise<TranslationJobManagementListResponseDto> {
    return this.invokeBinding("ListIncompleteJobs")
  }

  GetJobDetail(
    request: TranslationJobManagementGetDetailRequestDto
  ): Promise<TranslationJobManagementJobDetailDto> {
    return this.invokeBinding("GetJobDetail", request)
  }

  RequestStop(
    request: TranslationJobManagementActionRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding("RequestStop", request)
  }

  ResumeJob(
    request: TranslationJobManagementActionRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding("ResumeJob", request)
  }

  DeleteJob(
    request: TranslationJobManagementDeleteRequestDto
  ): Promise<TranslationJobManagementActionResponseDto> {
    return this.invokeBinding("DeleteJob", request)
  }
}

export function createTranslationJobManagementGateway(): TranslationJobManagementGatewayContract {
  return new TranslationJobManagementGateway(createBindingInvoker())
}
