import type {
  ImportTranslationInputRequest,
  RebuildTranslationInputCacheRequest,
  TranslationInputCommandResponse
} from "@application/gateway-contract/translation-input"

export type ImportTranslationInputRequestDto = ImportTranslationInputRequest
export type ImportTranslationInputResponseDto = TranslationInputCommandResponse

export type RebuildTranslationInputCacheRequestDto =
  RebuildTranslationInputCacheRequest
export type RebuildTranslationInputCacheResponseDto =
  TranslationInputCommandResponse