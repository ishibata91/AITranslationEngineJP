import type { TranslationInputGatewayContract } from "@application/gateway-contract/translation-input"
import type {
  ImportTranslationInputRequestDto,
  ImportTranslationInputResponseDto,
  RebuildTranslationInputCacheRequestDto,
  RebuildTranslationInputCacheResponseDto
} from "@controller/wails/gateway-dto/translation-input"

type TranslationInputBindingName =
  | "ImportTranslationInput"
  | "RebuildTranslationInputCache"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: TranslationInputBindingName,
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
  bindingName: TranslationInputBindingName
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
    toRecord(wailsRecord["TranslationInputController"]),
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
    bindingName: TranslationInputBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. backend-input-intake 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class TranslationInputGateway implements TranslationInputGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  importTranslationInput(
    request: ImportTranslationInputRequestDto
  ): Promise<ImportTranslationInputResponseDto> {
    return this.invokeBinding("ImportTranslationInput", request)
  }

  rebuildTranslationInputCache(
    request: RebuildTranslationInputCacheRequestDto
  ): Promise<RebuildTranslationInputCacheResponseDto> {
    return this.invokeBinding("RebuildTranslationInputCache", request)
  }
}

export function createTranslationInputGateway(): TranslationInputGatewayContract {
  return new TranslationInputGateway(createBindingInvoker())
}