import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
import type {
  CreateMasterDictionaryEntryRequestDto,
  CreateMasterDictionaryEntryResponseDto,
  DeleteMasterDictionaryEntryRequestDto,
  DeleteMasterDictionaryEntryResponseDto,
  GetMasterDictionaryEntryRequestDto,
  GetMasterDictionaryEntryResponseDto,
  ImportMasterDictionaryXmlRequestDto,
  ImportMasterDictionaryXmlResponseDto,
  ListMasterDictionaryEntriesRequestDto,
  ListMasterDictionaryEntriesResponseDto,
  UpdateMasterDictionaryEntryRequestDto,
  UpdateMasterDictionaryEntryResponseDto
} from "@controller/wails/gateway-dto/master-dictionary"

type MasterDictionaryBindingName =
  | "ListMasterDictionaryEntries"
  | "GetMasterDictionaryEntry"
  | "CreateMasterDictionaryEntry"
  | "UpdateMasterDictionaryEntry"
  | "DeleteMasterDictionaryEntry"
  | "ImportMasterDictionaryXml"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: MasterDictionaryBindingName,
  request: RequestDto
) => Promise<ResponseDto>

type BindingFunction = (request: unknown) => Promise<unknown>

function toRecord(value: unknown): Record<string, unknown> | null {
  if (typeof value !== "object" || value === null) {
    return null
  }
  return value as Record<string, unknown>
}

function resolveBindingFunction(
  bindingName: MasterDictionaryBindingName
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
    toRecord(wailsRecord["MasterDictionaryController"]),
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

function createBindingInvoker(): BindingInvoker {
  return <RequestDto, ResponseDto>(
    bindingName: MasterDictionaryBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. backend-crud/import-flow 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class MasterDictionaryGateway implements MasterDictionaryGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  listMasterDictionaryEntries(
    request: ListMasterDictionaryEntriesRequestDto
  ): Promise<ListMasterDictionaryEntriesResponseDto> {
    return this.invokeBinding("ListMasterDictionaryEntries", request)
  }

  getMasterDictionaryEntry(
    request: GetMasterDictionaryEntryRequestDto
  ): Promise<GetMasterDictionaryEntryResponseDto> {
    return this.invokeBinding("GetMasterDictionaryEntry", request)
  }

  createMasterDictionaryEntry(
    request: CreateMasterDictionaryEntryRequestDto
  ): Promise<CreateMasterDictionaryEntryResponseDto> {
    return this.invokeBinding("CreateMasterDictionaryEntry", request)
  }

  updateMasterDictionaryEntry(
    request: UpdateMasterDictionaryEntryRequestDto
  ): Promise<UpdateMasterDictionaryEntryResponseDto> {
    return this.invokeBinding("UpdateMasterDictionaryEntry", request)
  }

  deleteMasterDictionaryEntry(
    request: DeleteMasterDictionaryEntryRequestDto
  ): Promise<DeleteMasterDictionaryEntryResponseDto> {
    return this.invokeBinding("DeleteMasterDictionaryEntry", request)
  }

  importMasterDictionaryXml(
    request: ImportMasterDictionaryXmlRequestDto
  ): Promise<ImportMasterDictionaryXmlResponseDto> {
    return this.invokeBinding("ImportMasterDictionaryXml", request)
  }
}

export function createMasterDictionaryGateway(): MasterDictionaryGatewayContract {
  return new MasterDictionaryGateway(createBindingInvoker())
}
