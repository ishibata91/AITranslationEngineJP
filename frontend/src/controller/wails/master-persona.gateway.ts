import type { MasterPersonaGatewayContract } from "@application/gateway-contract/master-persona"
import type {
  MasterPersonaAISettingsDto,
  MasterPersonaDeleteRequestDto,
  MasterPersonaDetailResponseDto,
  MasterPersonaDialogueListResponseDto,
  MasterPersonaIdentityRequestDto,
  MasterPersonaMutationResponseDto,
  MasterPersonaPageRequestDto,
  MasterPersonaPageResponseDto,
  MasterPersonaPreviewRequestDto,
  MasterPersonaPreviewResultDto,
  MasterPersonaRunStatusDto,
  MasterPersonaUpdateRequestDto
} from "@controller/wails/gateway-dto/master-persona"

type MasterPersonaBindingName =
  | "MasterPersonaGetPage"
  | "MasterPersonaGetDetail"
  | "MasterPersonaGetDialogueList"
  | "MasterPersonaLoadAISettings"
  | "MasterPersonaSaveAISettings"
  | "MasterPersonaPreviewGeneration"
  | "MasterPersonaExecuteGeneration"
  | "MasterPersonaGetRunStatus"
  | "MasterPersonaInterruptGeneration"
  | "MasterPersonaCancelGeneration"
  | "MasterPersonaUpdate"
  | "MasterPersonaDelete"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: MasterPersonaBindingName,
  request?: RequestDto
) => Promise<ResponseDto>

type BindingFunction = (...args: [] | [request: unknown]) => Promise<unknown>

function toRecord(value: unknown): Record<string, unknown> | null {
  if (typeof value !== "object" || value === null) {
    return null
  }
  return value as Record<string, unknown>
}

function resolveBindingFunction(
  bindingName: MasterPersonaBindingName
): BindingFunction | null {
  const globalRecord = toRecord(globalThis)
  const goRecord = toRecord(globalRecord?.go)
  const wailsRecord = toRecord(goRecord?.wails)
  if (!wailsRecord) {
    return null
  }

  const controllerCandidates = [
    toRecord(wailsRecord.MasterPersonaController),
    toRecord(wailsRecord.AppController)
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
    bindingName: MasterPersonaBindingName,
    request?: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. backend master persona 完了後に接続します。`
        )
      )
    }

    if (request === undefined) {
      return binding().then((response) => response as ResponseDto)
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class MasterPersonaGateway implements MasterPersonaGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getMasterPersonaPage(
    request: MasterPersonaPageRequestDto
  ): Promise<MasterPersonaPageResponseDto> {
    return this.invokeBinding("MasterPersonaGetPage", request)
  }

  getMasterPersonaDetail(
    request: MasterPersonaIdentityRequestDto
  ): Promise<MasterPersonaDetailResponseDto> {
    return this.invokeBinding("MasterPersonaGetDetail", request)
  }

  getMasterPersonaDialogueList(
    request: MasterPersonaIdentityRequestDto
  ): Promise<MasterPersonaDialogueListResponseDto> {
    return this.invokeBinding("MasterPersonaGetDialogueList", request)
  }

  loadMasterPersonaAISettings(): Promise<MasterPersonaAISettingsDto> {
    return this.invokeBinding("MasterPersonaLoadAISettings")
  }

  saveMasterPersonaAISettings(
    request: MasterPersonaAISettingsDto
  ): Promise<MasterPersonaAISettingsDto> {
    return this.invokeBinding("MasterPersonaSaveAISettings", request)
  }

  previewMasterPersonaGeneration(
    request: MasterPersonaPreviewRequestDto
  ): Promise<MasterPersonaPreviewResultDto> {
    return this.invokeBinding("MasterPersonaPreviewGeneration", request)
  }

  executeMasterPersonaGeneration(
    request: MasterPersonaPreviewRequestDto
  ): Promise<MasterPersonaRunStatusDto> {
    return this.invokeBinding("MasterPersonaExecuteGeneration", request)
  }

  getMasterPersonaRunStatus(): Promise<MasterPersonaRunStatusDto> {
    return this.invokeBinding("MasterPersonaGetRunStatus")
  }

  interruptMasterPersonaGeneration(): Promise<MasterPersonaRunStatusDto> {
    return this.invokeBinding("MasterPersonaInterruptGeneration")
  }

  cancelMasterPersonaGeneration(): Promise<MasterPersonaRunStatusDto> {
    return this.invokeBinding("MasterPersonaCancelGeneration")
  }

  updateMasterPersona(
    request: MasterPersonaUpdateRequestDto
  ): Promise<MasterPersonaMutationResponseDto> {
    return this.invokeBinding("MasterPersonaUpdate", request)
  }

  deleteMasterPersona(
    request: MasterPersonaDeleteRequestDto
  ): Promise<MasterPersonaMutationResponseDto> {
    return this.invokeBinding("MasterPersonaDelete", request)
  }
}

export function createMasterPersonaGateway(): MasterPersonaGatewayContract {
  return new MasterPersonaGateway(createBindingInvoker())
}
