import type { TranslationOutputArtifactGatewayContract } from "@application/gateway-contract/translation-output-artifact"
import type {
  GenerateXTranslatorOutputArtifactRequestDto,
  GenerateXTranslatorOutputArtifactResponseDto,
  GetTranslationOutputDiffPreviewRequestDto,
  GetTranslationOutputDiffPreviewResponseDto,
  GetTranslationOutputReviewRequestDto,
  GetTranslationOutputReviewResponseDto,
  RegenerateXTranslatorOutputArtifactRequestDto,
  RegenerateXTranslatorOutputArtifactResponseDto
} from "@controller/wails/gateway-dto/translation-output-artifact"

type TranslationOutputArtifactBindingName =
  | "GetTranslationOutputReview"
  | "GetTranslationOutputDiffPreview"
  | "GenerateXTranslatorOutputArtifact"
  | "RegenerateXTranslatorOutputArtifact"

type BindingInvoker = <RequestDto, ResponseDto>(
  bindingName: TranslationOutputArtifactBindingName,
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
  bindingName: TranslationOutputArtifactBindingName
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
    toRecord(wailsRecord["TranslationOutputArtifactController"]),
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
    bindingName: TranslationOutputArtifactBindingName,
    request: RequestDto
  ): Promise<ResponseDto> => {
    const binding = resolveBindingFunction(bindingName)
    if (!binding) {
      return Promise.reject(
        new Error(
          `Wails binding is not wired yet: ${bindingName}. integration-output-artifact-wails-gateway 完了後に接続します。`
        )
      )
    }

    return binding(request).then((response) => response as ResponseDto)
  }
}

class TranslationOutputArtifactGateway implements TranslationOutputArtifactGatewayContract {
  constructor(private readonly invokeBinding: BindingInvoker) {}

  getTranslationOutputReview(
    request: GetTranslationOutputReviewRequestDto
  ): Promise<GetTranslationOutputReviewResponseDto> {
    return this.invokeBinding("GetTranslationOutputReview", request)
  }

  getTranslationOutputDiffPreview(
    request: GetTranslationOutputDiffPreviewRequestDto
  ): Promise<GetTranslationOutputDiffPreviewResponseDto> {
    return this.invokeBinding("GetTranslationOutputDiffPreview", request)
  }

  generateXTranslatorOutputArtifact(
    request: GenerateXTranslatorOutputArtifactRequestDto
  ): Promise<GenerateXTranslatorOutputArtifactResponseDto> {
    return this.invokeBinding("GenerateXTranslatorOutputArtifact", request)
  }

  regenerateXTranslatorOutputArtifact(
    request: RegenerateXTranslatorOutputArtifactRequestDto
  ): Promise<RegenerateXTranslatorOutputArtifactResponseDto> {
    return this.invokeBinding("RegenerateXTranslatorOutputArtifact", request)
  }
}

export function createTranslationOutputArtifactGateway(): TranslationOutputArtifactGatewayContract {
  return new TranslationOutputArtifactGateway(createBindingInvoker())
}
